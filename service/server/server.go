/* Main processing file for UTM server.
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
    "fmt"
    "github.com/brettlangdon/forge"
    "github.com/robmeades/utm/service/globals"
    "github.com/robmeades/utm/service/utilities"
    "log"
    "net/http"
    "github.com/codegangsta/negroni"
    "github.com/goincremental/negroni-sessions"
    "github.com/goincremental/negroni-sessions/cookiestore"
    "github.com/robmeades/utm/service/routes"
    "github.com/robmeades/utm/service/system"
    "time"
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

// A struct representing the state of a device
// in a single direction (UL or DL)
type DeviceTotalsState struct {
    Timestamp       time.Time
    DeviceUuid      string
    Msgs            int
    Bytes           int 
    Totals          *TotalsState 
    ExpectedMsgList *[]ExpectedMsg // Will be nil for uplink
}

// A struct representing the state of all devices
// in a single direction (UL or DL)
type TotalsState struct {
    Timestamp    time.Time
    Msgs         int
    Bytes        int    
}

// A struct to hold some parameters
// needed to track traffic test
type TtValues struct {
    DeviceUuid  string  `bson:"DeviceUuid" json:"DeviceUuid"`
    UlFill      byte
    UlLength    uint32
}

// Conection details for a device
type Connection struct {
    DeviceUuid      string  `bson:"DeviceUuid" json:"DeviceUuid"`
    DeviceName      string
    UlDevice        TotalsState
    DlDevice        TotalsState
    ExpectedMsgList *[]ExpectedMsg
    UlTotals        *TotalsState 
    DlTotals        *TotalsState    
}

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

// Server details
const configurationFile string = "config.cfg"

// Downlink channel to device
var DlMsgs chan<- AmqpMessage

// Count of AMQP messages received
var amqpMessageCount int

// Count of the number of times we've (re)starte AMQP
var amqpRetryCount int

// Keep track of the decode totals for all devices here
var totalsDecodeState TotalsState

// Keep track of the totals on the uplink for each device
var deviceDecodeStateList map[string]*DeviceTotalsState        

// Keep hold of the values needed on the uplink for traffic test mode
// and the number of traffic test bytes encoded
var deviceTtValuesList map[string]*TtValues        

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Process uplink messages from the AMQP queue until it is closed or
// an error is flagged
func processUlAmqpMsgs(q *Queue) {

    for {
        amqpMessageCount++
        select {            
            case msg, isOpen := <-q.UlAmqpMsgs:
            {
                if isOpen {
                    // Deal with a message on the AMQP Uplink queue
                    switch value := msg.(type) {
                        case *AmqpResponseMessage:
                        {
                            var ttDlBytes uint32
                            var ttDlDatagrams uint32
                            var ttTimeLastDl time.Time
                            
                            globals.Dbg.PrintfTrace("%s [server] --> decoded AMQP JSON msg:\n\n%+v\n\n", globals.LogTag, msg)
                            if value.Command == "UART_data" {
                                // UART data from the UTM-API which needs to be decoded                            
                                // If the device is not known, add it
                                decodeState := deviceDecodeStateList[value.DeviceUuid]
                                if decodeState == nil {
                                    decodeState = &DeviceTotalsState {
                                        Timestamp:       time.Now().UTC(),
                                        DeviceUuid:      value.DeviceUuid,
                                        Msgs:            0,
                                        Bytes:           0,
                                        Totals:          &totalsDecodeState,
                                        ExpectedMsgList: nil,
                                    }
                                    deviceDecodeStateList[value.DeviceUuid] = decodeState
                                }
                                
                                // Similarly, for the traffic test values we need to pass from
                                // the traffic test channel back to the decoder
                                ttValues := deviceTtValuesList[value.DeviceUuid]
                                if ttValues == nil {
                                    ttValues = &TtValues {
                                        DeviceUuid:  value.DeviceUuid,
                                        UlFill:      0,
                                        UlLength:    0,
                                    }
                                    deviceTtValuesList[value.DeviceUuid] = ttValues
                                }
                                
                                // Decode the messages
                                dlMsgs, byteCount := decode(value.Data, value.DeviceUuid, ttValues.UlFill, ttValues.UlLength)
            
                                // Pass the messages to the data table for recording
                                // and pass them to processing to see if any responses
                                // are required, totalising as we go.
                                // Also send them to the traffic test channel in case
                                // there is anything to do there
                                if dlMsgs != nil {
                                    decodeState.Timestamp = time.Now().UTC()
                                    decodeState.Bytes += byteCount
                                    totalsDecodeState.Timestamp = time.Now().UTC()
                                    totalsDecodeState.Bytes += byteCount
                                    for _, dlMsg := range dlMsgs {
                                        decodeState.Msgs++
                                        totalsDecodeState.Msgs++
                                        dataTableChannel <- dlMsg
                                        processMsgsChannel <- dlMsg
                                        trafficTestChannel <- dlMsg
                                    }
                                }
            
                                // Ask the traffic test channel for the expected
                                // UL fill value as we may need that for the next
                                // decode, and in any case we should keep the 
                                // datatable up to date with progress
                                getTrafficTestContext := &DeviceTrafficTestContextGet{
                                    DeviceUuid: value.DeviceUuid,
                                }
                                getTrafficTestContext.Context = make(chan TrafficTestContext)
                                trafficTestChannel <- getTrafficTestContext
                                trafficTestContext, isOpen := <- getTrafficTestContext.Context
                                if isOpen {
                                    if trafficTestContext.Parameters != nil {
                                        // The only things we need locally from this are the next
                                        // expected uplink fill value and the uplink datagram length
                                        ttValues.UlFill = trafficTestContext.UlFill
                                        ttValues.UlLength = trafficTestContext.Parameters.DeviceParameters.LenUlDatagram
                                    }
                                    // Also keep track of the downlink encode additions so that we can
                                    // add them to the totals
                                    ttDlBytes = trafficTestContext.DlBytesTotal 
                                    ttDlDatagrams = trafficTestContext.DlDatagramsTotal
                                    ttTimeLastDl = trafficTestContext.TimeLastDl
                                    // Send it all off to the datatable
                                    dataTableChannel <- &trafficTestContext
                                }
                                
                                // Ask the process channel for the encode status now
                                getEncodeState := &DeviceEncodeStateGet{
                                    DeviceUuid: value.DeviceUuid,
                                }
                                getEncodeState.State = make(chan DeviceTotalsState)
                                processMsgsChannel <- getEncodeState
                                encodeState, isOpen := <- getEncodeState.State
                                
                                // Only do this if the channel wasn't closed prematurely
                                // (e.g. due to an unknown UE being requested, which might
                                // happen if it's sending rubbish and so no messages
                                // will have been decoded)
                                if isOpen {
                                    // Add in the traffic test encoded stuff, if there is
                                    // any (the decode counts for the traffic test stuff
                                    // is already taken into account by the decoder) 
                                    if ttTimeLastDl.After (encodeState.Timestamp) {
                                        encodeState.Timestamp = ttTimeLastDl
                                    }
                                    encodeState.Totals.Msgs += int(ttDlDatagrams)
                                    encodeState.Msgs += int(ttDlDatagrams)
                                    encodeState.Totals.Bytes += int(ttDlBytes)
                                    encodeState.Bytes += int(ttDlBytes)
                                    // Send the datatable a message with connection
                                    // data for this device plus the totals for all
                                    ulTotals := TotalsState {
                                        Timestamp:    totalsDecodeState.Timestamp,
                                        Msgs:         totalsDecodeState.Msgs,
                                        Bytes:        totalsDecodeState.Bytes,
                                    }
                                    dlTotals := TotalsState {
                                        Timestamp:    encodeState.Timestamp,
                                        Msgs:         encodeState.Totals.Msgs,
                                        Bytes:        encodeState.Totals.Bytes,
                                    }
                                    // Send it off
                                    dataTableChannel <- &Connection {
                                        DeviceUuid:    value.DeviceUuid,
                                        DeviceName:    value.DeviceName,
                                        UlDevice: TotalsState {
                                            Timestamp: decodeState.Timestamp,
                                            Msgs:      decodeState.Msgs,
                                            Bytes:     decodeState.Bytes,
                                        },
                                        DlDevice: TotalsState {
                                            Timestamp: encodeState.Timestamp,
                                            Msgs:      encodeState.Msgs,
                                            Bytes:     encodeState.Bytes,
                                        },
                                        ExpectedMsgList: encodeState.ExpectedMsgList,
                                        UlTotals:        &ulTotals, 
                                        DlTotals:        &dlTotals,    
                                    }
                                }    
                            }
                        } 
                        case *error:
                        {   // If an error has occurred, drop out of the loop
                            globals.Dbg.PrintfTrace("%s [server] --> AMQP channel error received (%s), dropping out...\n", globals.LogTag, (*value).Error())
                            return
                        }
                        default:
                        {
                            globals.Dbg.PrintfError("%s [server] --> unknown message type on AMQP UlMsg channel: %+v.\n", globals.LogTag, msg)
                        } // case    
                    } // switch
                } else {
                    // The channel has closed so leave the infinite for() loop
                    return
                }    
            } // case    
        } // select
    } // for
}

// Handle all the AMQP queues
func doAmqp(username, amqpAddress string) {

    // Open the queue and then begin processing messages
    // If we drop out of the processing function, wait
    // a little while and try again
    for {
        fmt.Printf("######################################################################################################\n")
        fmt.Printf("UTM-API service (%s) REST interface opening %s...\n", globals.LogTag, amqpAddress)

        q, err := OpenQueue(username, amqpAddress)

        if err == nil {

            fmt.Printf("%s [server] --> connection opened.\n", globals.LogTag)

            DlMsgs = q.DlAmqpMsgs

            // The meat is in here
            processUlAmqpMsgs(q)
            
            // If we get to here, the AMQP channel has been closed,
            // which we never want.  Try again.
            q.Close()
        } else {
            globals.Dbg.PrintfTrace("%s [server] --> error opening AMQP queue (%s).\n", globals.LogTag, err.Error())
        }
        amqpRetryCount++
        globals.Dbg.PrintfTrace("%s [server] --> waiting before trying again...\n", globals.LogTag)
        time.Sleep(time.Second * 10)
    } // for
}

// Entry point
func Run() {

    // First, parse the configuration file
    settings, err := forge.ParseFile(configurationFile)
    if err != nil {
        panic(err)
    }
    amqp, err := settings.GetSection("amqp")
    username, err := amqp.GetString("uname")
    amqpAddress, err := amqp.GetString("amqp_address")
    host, err := settings.GetSection("host")
    port, err := host.GetString("port")

    // Set up the maps here. Do this here so that if we restart
    // the AMQP queue the data is retained
    deviceDecodeStateList = make(map[string]*DeviceTotalsState)        
    deviceTtValuesList = make(map[string]*TtValues)        

    // Handle all the Amqp messages
    go doAmqp(username, amqpAddress)

    // Set up logging
    log.SetFlags(log.LstdFlags)

    store := cookiestore.New([]byte("secretkey789"))
    router := routes.LoadRoutes()

    router.Handle("/frontPageData", utilities.Handler(GetFrontPageData))
    router.HandleFunc("/register", RegisterHandler)
    router.HandleFunc("/login", LoginHandler)
    router.HandleFunc("/display", ShowDisplayHandler)
    router.HandleFunc("/logout", LogoutHandler)

    sendMsg := ClientSendMsg{}
    router.HandleFunc("/sendMsg", sendMsg.Send).Methods("POST")

    n := negroni.Classic()
    static := negroni.NewStatic(http.Dir("static"))
    static.Prefix = "/static"
    n.Use(static)
    n.Use(negroni.HandlerFunc(system.MgoOpenDbSession))
    n.Use(sessions.Sessions("global_session_store", store))
    n.UseHandler(router)
    n.Use(negroni.HandlerFunc(system.MgoCloseDbSession))
    defer system.MgoCleanup()
    n.Run(port)
}

/* End Of File */
