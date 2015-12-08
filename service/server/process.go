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
    "reflect"
    "github.com/davecgh/go-spew/spew"
    "github.com/robmeades/utm/service/globals"
)

//--------------------------------------------------------------------
// Types 
//--------------------------------------------------------------------

// A message expected back from a device
type ExpectedMsg struct {
    TimeStarted time.Time
    ResponseId  ResponseTypeEnum
}

// Structure to allow the latest encode state for a particular
// device to be retrieved over a channel
type DeviceEncodeStateGet struct {
    DeviceUuid   string
    State        chan DeviceTotalsState
}

// Structure to allow other functions to add to the
// decode encode state that is held in here
// for a device
type DeviceEncodeStateAdd struct {
    DeviceUuid   string
    ResponseId   ResponseTypeEnum
    State        TotalsState
}

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

// To process a downlink message, send a MessageContainer into this channel
// containing the received message.  Any downlink messages required
// will be encoded and sent off to AMQP.  The number of bytes
// encoded will be added to a map against UUIDs

// To get the status of encoded messages (e.g. number of them or number
// of bytes), send a '*chan DeviceEncodeStateChannel' into this channel
// containing the device UUID and a pointer to a channel
// down which to send the EncodeState; a copy of all quantities will
// be copied into the struct and then the channel will be closed.

var processMsgsChannel chan<- interface{}

// Keep track of the encode totals for all devices here
var totalsEncodeState TotalsState

// Keep track of expected responses
var deviceExpectedMsgList map[string]*[]ExpectedMsg

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Remove a response ID from the list if it's been received
func removeResponse (responseId ResponseTypeEnum, uuid string) {
    if (responseId != RESPONSE_NONE) {
        globals.Dbg.PrintfTrace("%s [process] --> response ID %d received from UUID %s.\n", globals.LogTag, responseId, uuid)
        list := deviceExpectedMsgList[uuid]
        if list != nil {
            globals.Dbg.PrintfTrace("%s [process] --> found UUID %s in the expected message store, %d items in its list.\n", globals.LogTag, uuid, len(*list))
            for index, expectedMsg := range *list {
                if expectedMsg.ResponseId == responseId {
                    *list = append((*list)[:index], (*list)[index + 1:] ...)
                    globals.Dbg.PrintfTrace("%s [process] --> response ID %d removed from list for UUID %s.\n", globals.LogTag, responseId, uuid)
                    break
                }
            }
        }
    }    
}

// Add a response ID to the list
func addResponse (responseId ResponseTypeEnum, uuid string) {
    if (responseId != RESPONSE_NONE) {
        globals.Dbg.PrintfTrace("%s [process] --> now expecting response ID %d from UUID %s.\n", globals.LogTag, responseId, uuid)
        list := deviceExpectedMsgList[uuid]
        if list == nil {
            var expectedMsgList []ExpectedMsg 
            list = &expectedMsgList
            deviceExpectedMsgList[uuid] = list
            globals.Dbg.PrintfTrace("%s [process] --> making a new list for UUID %s, number of lists is now %d.\n", globals.LogTag, uuid, len(deviceExpectedMsgList))
        }
        
        expectedMsg := ExpectedMsg {
            TimeStarted: time.Now(),
            ResponseId:  responseId,            
        }
        
        *list = append(*list, expectedMsg)
        globals.Dbg.PrintfTrace("%s [process] --> expected list for UUID %s is now size %d.\n", globals.LogTag, uuid, len(*deviceExpectedMsgList[uuid]))
    }    
}

// Do the processing
func operateProcess() {
    deviceEncodeStateList := make(map[string]*DeviceTotalsState)        
    channel := make(chan interface{})
    processMsgsChannel = channel
    deviceExpectedMsgList = make(map[string]*[]ExpectedMsg)
    removeOldExpectedMsgs := time.NewTicker(time.Minute)
    
    globals.Dbg.PrintfTrace("%s [process] --> channel created and now being serviced.\n", globals.LogTag)
    
    // Remove old stuff from the expected message list on a tick
    go func() {
        for _ = range removeOldExpectedMsgs.C {
            for uuid, expectedMsgList := range deviceExpectedMsgList {
                var x = 0
                for x < len(*expectedMsgList) {
                    if time.Now().After((*expectedMsgList)[x].TimeStarted.Add(time.Hour)) {
                        globals.Dbg.PrintfTrace("%s [server] --> giving up after waiting > 1 hour for %d from device %s.\n", globals.LogTag, (*expectedMsgList)[x].ResponseId, uuid)
                        *expectedMsgList = append((*expectedMsgList)[:x], (*expectedMsgList)[x+1:]...)
                    }
                    x++
                }
            }
        }
    }()

    // Process commands on the channel
    go func() {
        for cmd := range channel {
            switch value := cmd.(type) {
                // Handle message containers holding somethings of interest,
                // throw everything else away
                case *MessageContainer:
                {
                    var err error = nil
                    var byteCount int = 0
                    var msgCount int = 0
                    responseIdReceived := RESPONSE_NONE
                    responseIdExpected := RESPONSE_NONE

                    // If the device is not in our list, add it and send an IntervalsGetReq
                    // if we aren't going to send one anyway lower down because of an InitIndUlmsg
                    // (or of course, if this is already an IntervalsGetCnfUlMsg)
                    encodeState := deviceEncodeStateList[value.DeviceUuid]
                    if encodeState == nil {
                        encodeState = &DeviceTotalsState {
                            Timestamp:  time.Now().UTC(),
                            DeviceUuid: value.DeviceUuid,
                            Msgs:       0,
                            Bytes:      0,
                            Totals:     &totalsEncodeState,
                        }
                        deviceEncodeStateList[value.DeviceUuid] = encodeState;
                        if reflect.TypeOf (value.Message) != reflect.TypeOf((*InitIndUlMsg)(nil)).Elem() &&
                           reflect.TypeOf (value.Message) != reflect.TypeOf((*IntervalsGetCnfUlMsg)(nil)).Elem() {
                            var count int 
                            err, count, responseIdExpected = encodeAndEnqueue (&IntervalsGetReqDlMsg{}, value.DeviceUuid)
                            if err == nil && count > 0 {
                                byteCount += count                           
                                msgCount++
                                addResponse (responseIdExpected, value.DeviceUuid)
                                responseIdExpected = RESPONSE_NONE
                            }    
                        }
                    }
                    
                    globals.Dbg.PrintfTrace("%s [process] --> processing message from UUID %s...\n\n%s\n\n", globals.LogTag, value.DeviceUuid, spew.Sdump(cmd))
                    
                    switch utmMsg := value.Message.(type) {
                        case *TransparentUlDatagram:
                        {
                            // Nothing to do here
                        }
                        case *PingReqUlMsg:
                        {
                            // Respond
                            var count int 
                            err, count, responseIdExpected = encodeAndEnqueue (&PingCnfDlMsg{}, value.DeviceUuid)
                            if err == nil && count > 0 {
                                byteCount += count                           
                                msgCount++
                            }    
                        }
                        case *PingCnfUlMsg:
                        {
                            // Nothing to do here
                            responseIdReceived = RESPONSE_PING_CNF
                        }
                        case *InitIndUlMsg:
                        {
                            if uint8(RevisionLevel) != utmMsg.RevisionLevel {
                                globals.Dbg.PrintfInfo("%s [process] --> UUID %s has protocol revision %d, which is different to this server (%d)", globals.LogTag, value.DeviceUuid, utmMsg.RevisionLevel, RevisionLevel)
                            }    
                            // Get the reporting intervals for this device
                            var count int 
                            err, count, responseIdExpected = encodeAndEnqueue (&IntervalsGetReqDlMsg{}, value.DeviceUuid)
                            if err == nil && count > 0 {
                                byteCount += count                           
                                msgCount++
                            }    
                        }
                        case *DateTimeIndUlMsg:
                        {
                            // Send DataTimeSetReqDlMsg if out of range
                            // First check if the time is OK but the date is absent
                            // (which could be the case if the time has been set by GNSS,
                            // and we don't want to mess with that time as it will be more
                            // accurate than ours)
                            utmDays := utmMsg.UtmTime.Year() * 365 +  utmMsg.UtmTime.YearDay()
                            serverDays := time.Now().UTC().Year() * 365 + time.Now().UTC().YearDay()
                            // Allow plenty of slack as messages might have been queued                    
                            if (utmDays < serverDays - 30) || (utmDays > serverDays + 30) {
                                dateTimeSetReq := &DateTimeSetReqDlMsg {
                                    UtmTime:      time.Now().UTC(),
                                    SetDateOnly:  false,
                                }
                                
                                if utmMsg.TimeSetBy == TIME_SET_BY_GNSS {
                                    dateTimeSetReq.SetDateOnly = true;
                                }
                                
                                var count int 
                                err, count, responseIdExpected = encodeAndEnqueue (dateTimeSetReq, value.DeviceUuid)
                                if err == nil && count > 0 {
                                    byteCount += count                           
                                    msgCount++
                                }    
                            }
                        }
                        case *DateTimeSetCnfUlMsg:
                        {
                            // Nothing to do here
                            responseIdReceived = RESPONSE_DATE_TIME_SET_CNF
                        }
                        case *DateTimeGetCnfUlMsg:
                        {
                            // Nothing to do here
                            responseIdReceived = RESPONSE_DATE_TIME_GET_CNF
                        }
                        case *ModeSetCnfUlMsg:
                        {
                            // TODO start sending downlink datagrams if in a traffic test
                            responseIdReceived = RESPONSE_MODE_SET_CNF
                        }
                        case *ModeGetCnfUlMsg:
                        {
                            // Nothing to do here
                            responseIdReceived = RESPONSE_MODE_GET_CNF
                        }
                        case *IntervalsGetCnfUlMsg:
                        {
                            // Nothing to do here
                            responseIdReceived = RESPONSE_INTERVALS_GET_CNF
                        }
                        case *ReportingIntervalSetCnfUlMsg:
                        {
                            // Nothing to do here
                            responseIdReceived = RESPONSE_REPORTING_INTERVAL_SET_CNF
                        }
                        case *HeartbeatSetCnfUlMsg:
                        {
                            // Nothing to do here
                            responseIdReceived = RESPONSE_HEARTBEAT_SET_CNF
                        }
                        case *PollIndUlMsg:
                        {
                            // Nothing to do here
                        }
                        case *MeasurementsIndUlMsg:
                        {
                            // Nothing to do here
                        }
                        case *MeasurementsGetCnfUlMsg:
                        {
                            // Nothing to do here
                            responseIdReceived = RESPONSE_MEASUREMENTS_GET_CNF
                        }
                        // case *MeasurementsControlIndUlMsg:
                        // case *MeasurementControlSetCnfUlMsg:
                        // case *MeasurementsControlGetCnfUlMsg:
                        // case *MeasurementsControlDefaultsSetCnfUlMsg:
                        // TODO
                        case *TrafficReportIndUlMsg:
                        {
                            // Nothing to do here
                        }
                        case *TrafficReportGetCnfUlMsg:
                        {
                            // Nothing to do here
                            responseIdReceived = RESPONSE_TRAFFIC_REPORT_GET_CNF
                        }
                        case *TrafficTestModeParametersSetCnfUlMsg:
                        {
                            // TODO send ModeSetReqDlMsg if in a traffic test 
                            responseIdReceived = RESPONSE_TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF
                        }
                        case *TrafficTestModeParametersGetCnfUlMsg:
                        {
                            // Nothing to do here
                            responseIdReceived = RESPONSE_TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF
                        }
                        case *TrafficTestModeRuleBreakerUlDatagram:
                        {
                            // Nothing to do here, the traffic test mode channel will deal with this
                        }
                        case *BadTrafficTestModeRuleBreakerUlDatagram:
                        {
                            // Nothing to do here, the traffic test mode channel will deal with this
                        }
                        case *OutOfSequenceTrafficTestModeRuleBreakerUlDatagram:
                        {
                            // Nothing to do here, the traffic test mode channel will deal with this
                        }
                        case *TrafficTestModeReportIndUlMsg:
                        {
                            // Nothing to do here
                        }
                        case *TrafficTestModeReportGetCnfUlMsg:
                        {
                            // Nothing to do here
                            responseIdReceived = RESPONSE_TRAFFIC_TEST_MODE_REPORT_GET_CNF
                        }
                        case *ActivityReportIndUlMsg:
                        {
                            // Nothing to do here
                        }
                        case *ActivityReportGetCnfUlMsg:
                        {
                            // Nothing to do here
                            responseIdReceived = RESPONSE_ACTIVITY_REPORT_GET_CNF
                        }
                        case *DebugIndUlMsg:
                        {
                            // Nothing to do here
                        }
                        default:
                        {
                            // Ignore any unknown UTM messages
                            globals.Dbg.PrintfTrace("%s [process] --> unrecognised UTM message, ignoring.\n", globals.LogTag)
                            globals.Dbg.PrintfInfo("%s [process] --> unrecognised UTM message was:\n\n%s\n", globals.LogTag, spew.Sdump(utmMsg))
                        } // case
                    } // switch
                    
                    // Count up the downlink status for this devce
                    if msgCount > 0 {
                        encodeState.Timestamp = time.Now().UTC()
                        encodeState.Msgs += msgCount
                        encodeState.Bytes += byteCount
                        totalsEncodeState.Timestamp = encodeState.Timestamp
                        totalsEncodeState.Msgs += msgCount
                        totalsEncodeState.Bytes += byteCount
                    }    

                    // Deal with expected or received responses
                    removeResponse (responseIdReceived, value.DeviceUuid)
                    // Remove an expected response first as an addition is likely
                    // to be a re-transmission of the first
                    removeResponse (responseIdExpected, value.DeviceUuid)
                    addResponse (responseIdExpected, value.DeviceUuid)
                    
                    globals.Dbg.PrintfTrace("%s [process] --> processing completed.\n", globals.LogTag)
                } // case    
                // Add to the encode state totals in here, useful if other things send messages
                case *DeviceEncodeStateAdd:
                {
                    // Retrieve the encode state
                    if encodeState, isPresent := deviceEncodeStateList[value.DeviceUuid]; isPresent {
                        encodeState.Timestamp = value.State.Timestamp
                        encodeState.Msgs += value.State.Msgs
                        encodeState.Bytes += value.State.Bytes
                        if encodeState.Totals != nil {
                            encodeState.Totals.Timestamp = value.State.Timestamp
                            encodeState.Totals.Msgs += value.State.Msgs
                            encodeState.Totals.Bytes += value.State.Bytes
                        }
                           
                        // Add the response ID
                        addResponse (value.ResponseId, value.DeviceUuid)
                    }
                }
                // Return the encode state for a given UUID 
                case *DeviceEncodeStateGet:
                {
                    // Retrieve the encode state
                    if encodeState, isPresent := deviceEncodeStateList[value.DeviceUuid]; isPresent {
                        // Copy in the EncodeState data, post it and close the channel
                        globals.Dbg.PrintfTrace("%s [process] --> fetching encode state for UUID %s.\n", globals.LogTag, value.DeviceUuid)
                        totalsState := TotalsState {
                            Timestamp:  totalsEncodeState.Timestamp,
                            Msgs:       totalsEncodeState.Msgs,
                            Bytes:      totalsEncodeState.Bytes,                            
                        }
                        state := DeviceTotalsState {
                            DeviceUuid: encodeState.DeviceUuid,
                            Timestamp:  encodeState.Timestamp,
                            Msgs:       encodeState.Msgs,
                            Bytes:      encodeState.Bytes,
                            Totals:     &totalsState,
                        }
                        if expectedMsgList, hasExpectedMsgList := deviceExpectedMsgList[value.DeviceUuid]; hasExpectedMsgList {
                            tmp := make([]ExpectedMsg, 0)
                            state.ExpectedMsgList = &tmp
                            for _, expectedMsg := range *expectedMsgList {
                                *state.ExpectedMsgList = append(*state.ExpectedMsgList, expectedMsg)
                            }
                        }
                        value.State <- state
                        globals.Dbg.PrintfTrace("%s [process] --> provided encode state.\n", globals.LogTag)
                    } else {
                        globals.Dbg.PrintfTrace("%s [process] --> asked for encode state for unknown UUID %s.\n", globals.LogTag, value.DeviceUuid)
                    }
                    close(value.State)
                } // case       
                default:
                {
                    globals.Dbg.PrintfTrace("%s [process] --> unrecognised command, ignoring.\n", globals.LogTag)
                    globals.Dbg.PrintfInfo("%s [process] --> unrecognised command was:\n\n%s\n", globals.LogTag, spew.Sdump(cmd))
                } // case
            } // switch
        } // for
        globals.Dbg.PrintfTrace("%s [process] --> command channel closed, stopping.\n", globals.LogTag)
    }()
}

func init() {
    operateProcess()
}

/* End Of File */
