/* Message processing functions for the UTM server.
 *
 * Copyright (C) u-blox Melbourn Ltd
 * u-blox Melbourn Ltd, Melbourn, UK
 *
 * All rights reserved.
 *
 * This source file is the sole property of u-blox Melbourn Ltd.
 * Reproduction or utilization of this source in whole or part is
 * forbidden without the written consent of u-blox Melbourn Ltd.
 */

package server

import (
    "time"
    "github.com/davecgh/go-spew/spew"
    "github.com/robmeades/utm/service/globals"
)

//--------------------------------------------------------------------
// Types 
//--------------------------------------------------------------------

type TrafficTestStateEnum uint32

// IMPORTANT: the values here are used directly
// in the client (in component value_tt_numbers.react.js),
// make sure to keep them aligned
const (
    TRAFFIC_TEST_NOT_RUNNING  TrafficTestStateEnum = iota
    TRAFFIC_TEST_RUNNING
    TRAFFIC_TEST_TX_COMPLETE
    TRAFFIC_TEST_STOPPED
    TRAFFIC_TEST_TIMEOUT
    TRAFFIC_TEST_PASS
    TRAFFIC_TEST_FAIL
)

// The traffic test mode parameter messages for the server
type TrafficTestModeParametersServerSet struct {
    DeviceParameters   *TrafficTestModeParametersSetReqDlMsg
    DlIntervalSeconds   uint32
}

// The context for each traffic test, DeepCopy functions
// further down the page
// IMPORTANT: if you change this structure, change that function
// as well
type TrafficTestContext struct {
    DeviceUuid          string
    Parameters          *TrafficTestModeParametersServerSet
    DeviceTrafficReport *TrafficTestModeReportIndUlMsg
    TimeStarted         time.Time
    DlTimeStopped       time.Time
    UlTimeStopped       time.Time
    TimeLastDl          time.Time
    TimeLastUl          time.Time
    DlDatagramsTotal    uint32 // These two variables are used to keep track of...
    DlDatagrams         uint32
    DlBytesTotal        uint32 // ...what we're sending for the encode totals.
    DlBytes             uint32
    UlDatagrams         uint32
    UlBytes             uint32
    UlDatagramsMissed   uint32
    UlDatagramsBad      uint32
    UlDatagramsOOS      uint32
    DlFill              byte
    UlFill              byte
    DlState             TrafficTestStateEnum
    UlState             TrafficTestStateEnum
}

// Structure to allow the latest traffic test
// context for a particular device to be retrieved
// over a channel
type DeviceTrafficTestContextGet struct {
    DeviceUuid   string
    Context      chan TrafficTestContext
}

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

var trafficTestChannel chan<- interface{}

// Keep track of each test
var trafficTestList map[string]*TrafficTestContext

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// This function, like those in data.go,
// are used in datatable.go.  They are kept here
// to be close to the structure definition
func (value *TrafficTestContext) DeepCopy() *TrafficTestContext {
    if value == nil {
        return nil
    }
    result := &TrafficTestContext {
        DeviceUuid:          value.DeviceUuid,
        TimeStarted:         value.TimeStarted,
        DlTimeStopped:       value.DlTimeStopped,
        UlTimeStopped:       value.UlTimeStopped,
        TimeLastDl:          value.TimeLastDl,
        TimeLastUl:          value.TimeLastUl,
        DlDatagramsTotal:    value.DlDatagramsTotal,
        DlDatagrams:         value.DlDatagrams,
        DlBytesTotal:        value.DlBytesTotal,
        DlBytes:             value.DlBytes,
        UlDatagrams:         value.UlDatagrams,
        UlBytes:             value.UlBytes,
        UlDatagramsMissed:   value.UlDatagramsMissed,
        UlDatagramsBad:      value.UlDatagramsBad,
        UlDatagramsOOS:      value.UlDatagramsOOS,
        DlFill:              value.DlFill,
        UlFill:              value.UlFill,
        DlState:             value.DlState,
        UlState:             value.UlState,
    }
    if value.Parameters != nil {
        result.Parameters = &TrafficTestModeParametersServerSet {
            DlIntervalSeconds: value.Parameters.DlIntervalSeconds,
        }
        result.Parameters.DeviceParameters = value.Parameters.DeviceParameters.DeepCopy()                            
    }
    result.DeviceTrafficReport = value.DeviceTrafficReport.DeepCopy()                            

    return result
}

// Make a set of traffic test mode parameters for us
func makeParameters(parameters *TrafficTestModeParametersServerSet) *TrafficTestModeParametersServerSet {

    data := TrafficTestModeParametersServerSet {}    
    data.DeviceParameters = parameters.DeviceParameters.DeepCopy()
    data.DlIntervalSeconds = parameters.DlIntervalSeconds

    return &data
}

// Do the traffic test stuff
func operateTrafficTest() {
    channel := make(chan interface{})
    trafficTestChannel = channel
    trafficTestList = make(map[string]*TrafficTestContext)
    sendDlMsg := time.NewTicker(time.Second)
    
    globals.Dbg.PrintfTrace("%s [traffic_test] --> channel created and now being serviced.\n", globals.LogTag)
    
    // Send downlink messages and assess the state of the uplink on a timer
    go func() {
        for _ = range sendDlMsg.C {
            for _, context := range trafficTestList {
                if context.Parameters != nil {
                    // The DL part
                    if context.DlState == TRAFFIC_TEST_RUNNING {
                        if context.DlDatagrams < context.Parameters.DeviceParameters.NumDlDatagrams {
                            if time.Now().UTC().After (context.TimeLastDl.Add (time.Duration(context.Parameters.DlIntervalSeconds) * time.Second)) {
                                // Time to send a DL traffic test datagram
                                msg := &TrafficTestModeRuleBreakerDlDatagram {
                                        Fill:      context.DlFill,
                                        Length:    context.Parameters.DeviceParameters.LenDlDatagram,
                                    }
                                err, byteCount, _ := encodeAndEnqueue (msg, context.DeviceUuid)
                                if err == nil {
                                    context.TimeLastDl = time.Now().UTC()
                                    context.DlFill++
                                    context.DlDatagramsTotal++
                                    context.DlDatagrams++
                                    context.DlBytesTotal += uint32(byteCount)
                                    context.DlBytes += uint32(byteCount)
                                    globals.Dbg.PrintfTrace("%s [traffic_test] --> sent DL datagram %d (fill %d) to device %s.\n",
                                        globals.LogTag, context.DlDatagrams, context.DlFill, context.DeviceUuid)
                                }
                            }    
                        } else {
                            context.UlTimeStopped = time.Now().UTC()
                            context.DlState = TRAFFIC_TEST_TX_COMPLETE;
                            globals.Dbg.PrintfTrace("%s [traffic_test] --> DL completed for device %s (%d DL datagrams sent).\n",
                                globals.LogTag, context.DeviceUuid, context.DlDatagrams)
                        }
                    }   
                    // The UL part
                    if context.UlState == TRAFFIC_TEST_RUNNING {
                        // Assess whether we're donn in the downlink direction from the report
                        if context.UlDatagrams + context.UlDatagramsMissed >= context.Parameters.DeviceParameters.NumUlDatagrams {
                            if context.UlDatagramsMissed == 0 {
                                context.UlState = TRAFFIC_TEST_PASS
                                context.UlTimeStopped = time.Now().UTC()
                                globals.Dbg.PrintfTrace("%s [traffic_test] --> UL PASS on traffic test with device %s, %d datagrams received.\n",
                                    globals.LogTag, context.DeviceUuid, context.UlDatagrams)
                            } else {
                                context.UlState = TRAFFIC_TEST_FAIL
                                context.UlTimeStopped = time.Now().UTC()
                                globals.Dbg.PrintfTrace("%s [traffic_test] --> DL FAIL on traffic test with device %s (%d datagrams missed out of %d).\n",
                                    globals.LogTag, context.DeviceUuid, context.UlDatagramsMissed, context.UlDatagrams + context.UlDatagramsMissed)
                            }
                        }                        
                    } // UL test running, check for timeout (don't worry about the DL as the UTM-end will timeout
                      // on the DL taking too long)
                    if (context.UlState == TRAFFIC_TEST_RUNNING) || (context.UlState == TRAFFIC_TEST_TX_COMPLETE) {
                        if context.Parameters.DeviceParameters.TimeoutSeconds > 0 {
                            if time.Now().UTC().After (context.TimeStarted.Add (time.Duration(context.Parameters.DeviceParameters.TimeoutSeconds) * time.Second)) {
                                context.UlState = TRAFFIC_TEST_TIMEOUT;
                                context.UlTimeStopped = time.Now().UTC()
                                globals.Dbg.PrintfTrace("%s [traffic_test] --> UL timed out for %s (started %s, %d DL datagrams sent).\n",
                                    globals.LogTag, context.DeviceUuid, context.TimeStarted, context.DlDatagrams)
                            }
                        }                                
                    }
                } // Params != nil
            } // for each UUID 
        }
    }()

    // Process commands on the channel
    // This channel receives copies of messages going to
    // the device so that it can keep track of what it should be
    // doing.  It does not sent any messages to the device, the
    // timer above does that.
    go func() {
        for cmd := range channel {
            switch value := cmd.(type) { 
                case *MessageContainer:
                {
                    /// Set up a context for this device if not already done
                    context := trafficTestList[value.DeviceUuid]
                    if context == nil {
                        context = &TrafficTestContext{};
                        trafficTestList[value.DeviceUuid] = context;
                        context.DeviceUuid = value.DeviceUuid
                        globals.Dbg.PrintfTrace("%s [traffic_test] --> created context for device %s.\n", globals.LogTag, value.DeviceUuid)
                    }
                    
                    switch utmMsg := value.Message.(type) {
                        case *TrafficTestModeParametersServerSet:
                        {
                            globals.Dbg.PrintfTrace("%s [traffic_test] --> parameters received for device  %s:\n\n%+v\n",
                                globals.LogTag, value.DeviceUuid, utmMsg.DeviceParameters)
                            context.Parameters = makeParameters(utmMsg)
                        }                        
                        case *ModeSetReqDlMsg:
                        {
                            // If the mode set request is to switch out of traffic test mode
                            // then end the test, otherwise do nothing, since we only go
                            // INTO traffic test mode when the device has confirmed that it is
                            // in traffic test mode
                            if utmMsg.Mode != MODE_TRAFFIC_TEST {
                                context.DlState = TRAFFIC_TEST_STOPPED;
                                context.UlState = TRAFFIC_TEST_STOPPED;
                                context.DlTimeStopped = time.Now().UTC()
                                context.UlTimeStopped = time.Now().UTC()
                                globals.Dbg.PrintfTrace("%s [traffic_test] --> stopped traffic test with device %s.\n", globals.LogTag, value.DeviceUuid)
                            }
                        }
                        case *ModeSetCnfUlMsg:
                        {
                            // If the device is confirming a mode change into
                            // traffic test mode then start doing our stuff
                            if (utmMsg.Mode == MODE_TRAFFIC_TEST) && (context.Parameters != nil) {
                                context.TimeLastDl = time.Time{}
                                context.DlFill = DlFillMinValue
                                context.UlFill = UlFillMinValue
                                context.DlDatagrams = 0
                                context.DlBytes = 0
                                context.UlDatagrams = 0
                                context.UlBytes = 0
                                context.UlDatagramsMissed = 0
                                context.UlDatagramsBad = 0
                                context.UlDatagramsOOS = 0
                                context.TimeStarted = time.Now().UTC()
                                context.DlState = TRAFFIC_TEST_RUNNING
                                context.UlState = TRAFFIC_TEST_RUNNING
                                context.DeviceTrafficReport = nil
                                globals.Dbg.PrintfTrace("%s [traffic_test] --> started traffic test with device %s.\n", globals.LogTag, value.DeviceUuid)
                            }
                        }
                        case *TrafficTestModeReportIndUlMsg:
                        {
                            globals.Dbg.PrintfTrace("%s [traffic_test] --> received report from %s, assessing DL state.\n",
                                 globals.LogTag, value.DeviceUuid)
                            // First, store this
                            context.DeviceTrafficReport = utmMsg.DeepCopy()

                            // Assess whether we're done in the downlink direction from the report
                            if utmMsg.TimedOut {
                                context.DlState = TRAFFIC_TEST_TIMEOUT
                                context.DlTimeStopped = time.Now().UTC()
                                globals.Dbg.PrintfTrace("%s [traffic_test] --> DL TIMEOUT on traffic test with device %s (%d datagrams missed).\n",
                                    globals.LogTag, value.DeviceUuid, context.DeviceTrafficReport.NumTrafficTestDlDatagramsMissed)
                            } else {
                                if context.DeviceTrafficReport.NumTrafficTestDatagramsDl + context.DeviceTrafficReport.NumTrafficTestDlDatagramsMissed >= context.DlDatagrams {
                                    if context.DeviceTrafficReport.NumTrafficTestDlDatagramsMissed == 0 {
                                        context.DlState = TRAFFIC_TEST_PASS
                                        context.DlTimeStopped = time.Now().UTC()
                                        globals.Dbg.PrintfTrace("%s [traffic_test] --> DL PASS on traffic test with device %s, %d datagrams received.\n",
                                            globals.LogTag, value.DeviceUuid, context.DeviceTrafficReport.NumTrafficTestDatagramsDl)
                                    } else {
                                        context.DlState = TRAFFIC_TEST_FAIL
                                        context.DlTimeStopped = time.Now().UTC()
                                        globals.Dbg.PrintfTrace("%s [traffic_test] --> DL FAIL on traffic test with device %s (%d datagrams missed out of %d).\n",
                                            globals.LogTag, value.DeviceUuid, context.DeviceTrafficReport.NumTrafficTestDlDatagramsMissed,
                                            context.DeviceTrafficReport.NumTrafficTestDlDatagramsMissed +  context.DeviceTrafficReport.NumTrafficTestDatagramsDl)
                                    }
                                }
                            }
                        }
                        case *TrafficTestModeRuleBreakerUlDatagram:
                        {
                            if context.UlState == TRAFFIC_TEST_RUNNING {
                                context.UlDatagrams++
                                context.TimeLastUl = time.Now().UTC()
                                if context.Parameters != nil {
                                    context.UlBytes += context.Parameters.DeviceParameters.LenUlDatagram
                                }
                                context.UlFill++
                                globals.Dbg.PrintfTrace("%s [traffic_test] --> received good UL traffic test mode datagram %d from  %s, incremented expected fill to %d.\n",
                                    globals.LogTag, context.UlDatagrams, value.DeviceUuid, context.UlFill)                                
                            }
                        }
                        case *BadTrafficTestModeRuleBreakerUlDatagram:
                        {
                            if context.UlState == TRAFFIC_TEST_RUNNING {
                                context.UlDatagramsBad++
                                context.UlDatagramsMissed++
                                context.UlFill++
                                context.TimeLastUl = time.Now().UTC()
                                globals.Dbg.PrintfTrace("%s [traffic_test] --> received %d BAD UL traffic test mode datagram(s) from  %s, missed count is now %d, incremented expected fill to %d.\n",
                                    globals.LogTag, context.UlDatagramsBad, value.DeviceUuid, context.UlDatagramsMissed, context.UlFill)
                            }
                        }
                        case *OutOfSequenceTrafficTestModeRuleBreakerUlDatagram:
                        {
                            if context.UlState == TRAFFIC_TEST_RUNNING {
                                context.UlDatagramsOOS++
                                // Account for the gap in the fill and resynchronise ourselves
                                context.UlDatagramsMissed += uint32(utmMsg.Fill - context.UlFill)
                                context.UlFill = utmMsg.Fill + 1
                                context.TimeLastUl = time.Now().UTC()
                                globals.Dbg.PrintfTrace("%s [traffic_test] --> received %d OOS UL traffic test mode datagram(s) from %s, missed count is now, incremented expected fill to %d.\n",
                                    globals.LogTag, context.UlDatagramsOOS, value.DeviceUuid, context.UlDatagramsMissed, context.UlFill)
                            }
                        }
                        default:
                        {
                            // Do nothing
                        } // case
                    } // switch                        
                } // case
                // Return the traffic test context for a given UUID 
                case *DeviceTrafficTestContextGet:
                {
                    // Retrieve the context
                    if context, isPresent := trafficTestList[value.DeviceUuid]; isPresent {
                        // Copy in the context data, post it and close the channel
                        globals.Dbg.PrintfTrace("%s [traffic_test] --> fetching context for UUID %s.\n", globals.LogTag, value.DeviceUuid)
                        contextCopy := context.DeepCopy()
                        value.Context <- *contextCopy
                        globals.Dbg.PrintfTrace("%s [traffic_test] --> provided context.\n", globals.LogTag)
                    } else {
                        globals.Dbg.PrintfTrace("%s [traffic_test] --> asked for context for unknown UUID %s.\n", globals.LogTag, value.DeviceUuid)
                    }
                    close(value.Context)
                }       
                default:
                {
                    globals.Dbg.PrintfTrace("%s [traffic_test] --> unrecognised command, ignoring.\n", globals.LogTag)
                    globals.Dbg.PrintfInfo("%s [traffic_test] --> unrecognised command was:\n\n%s\n", globals.LogTag, spew.Sdump(cmd))
                } // case    
            } // switch
        } // for

        globals.Dbg.PrintfTrace("%s [traffic_test] --> command channel closed, stopping.\n", globals.LogTag)
    }()
}

func init() {
    operateTrafficTest()
}

/* End Of File */
