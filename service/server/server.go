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

var deviceDecodeStateList map[string]*DeviceTotalsState        

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
                                
                                // Decode the messages
                                dlMsgs, byteCount := decode(value.Data, value.DeviceUuid)
            
                                // Pass the messages to the data table for recording
                                // and pass them to processing to see if any responses
                                // are required, totalising as we go
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
                                    }
                                }
            
                                // Ask the process loop for the encode status now
                                get := &DeviceEncodeStateGet{
                                    DeviceUuid: value.DeviceUuid,
                                }
                                get.State = make(chan DeviceTotalsState)
                                processMsgsChannel <- get
                                encodeState := <- get.State
            
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
                        case *error:
                        {   // If an error has occurred, drop out of the loop
                            globals.Dbg.PrintfTrace("%s [server] --> AMQP channel error received (%s), dropping out...\n", globals.LogTag, (*value).Error())
                            return
                        }
                        default:
                        {
                            globals.Dbg.PrintfTrace("%s [server] --> unknown message type on AMQP UlMsg channel: %+v.\n", globals.LogTag, msg)
                            log.Fatal(globals.LogTag, "invalid message type.")
                        } // case    
                    } // switch
                } else {
                    // The channel has closed so leave the infinite for() loop
                    return
                }    
            } // case    
        } // select
    }
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
    }
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

    // Set up the decode state list map
    // Do this here so that if
    deviceDecodeStateList = make(map[string]*DeviceTotalsState)        

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
