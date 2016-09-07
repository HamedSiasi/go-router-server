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
	"github.com/dustin/go-coap"
	"github.com/u-blox/utm-server/service/globals"
	"log"
	"strings"
	"time"
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

// A struct representing the state of a device
// in a single direction (UL or DL)
/*
type DeviceTotalsState struct {
	Timestamp       time.Time
	DeviceUuid      string
	DeviceName      string
	Msgs            int
	Bytes           int
	Totals          *TotalsState
	ExpectedMsgList *[]ExpectedMsg // Will be nil for uplink
}*/

// A struct representing the state of all devices
// in a single direction (UL or DL)
type TotalsState struct {
	Timestamp time.Time
	Msgs      int
	Bytes     int
}

// A struct to hold some parameters
// needed to track traffic test
type TtValues struct {
	DeviceUuid string `bson:"DeviceUuid" json:"DeviceUuid"`
	UlFill     byte
	UlLength   uint32
}

// Conection details for a device
/*
type Connection struct {
	DeviceUuid      string `bson:"DeviceUuid" json:"DeviceUuid"`
	DeviceName      string
	RecentUl        bool
	UlDevice        TotalsState
	DlDevice        TotalsState
	ExpectedMsgList *[]ExpectedMsg
	UlTotals        *TotalsState
	DlTotals        *TotalsState
}*/

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
//var deviceDecodeStateList map[string]*DeviceTotalsState

// Keep hold of the values needed on the uplink for traffic test mode
// and the number of traffic test bytes encoded
var deviceTtValuesList map[string]*TtValues

var gRouter = make(map[string]string)

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------
// process RECEIVED Msgs
// Process uplink messages from the AMQP queue until it is closed or
// an error is flagged
func processUlAmqpMsgs(q *Queue) {
	//statusTick := time.NewTicker(time.Minute)

	// Establish the status of the process loop and of traffic test
	// for all devices on a background timer, just in case there are
	// no datagrams arriving from them
	/*
		go func() {
			for _ = range statusTick.C {
				for uuid, decodeState := range deviceDecodeStateList {
					updateDatatableState(uuid, decodeState, false)
				}
			}
		}()*/

	// The uplink processing loop
	for {
		amqpMessageCount++
		select {
		case msg, isOpen := <-q.UlAmqpMsgs:
			{
				if isOpen {
					// Deal with a message on the AMQP Uplink queue
					// type switch
					switch value := msg.(type) {
					case *AmqpResponseMessage:
						{
							if value.Command == "UART_data" {
								// UART data from the UTM-API which needs to be decoded
								// If the device is not known, add it
								value.DeviceUuid = strings.ToLower(value.DeviceUuid)
								//----------------------------------------------------------------------------------------------
								// (1) Receive readingData from SoftRadio
								//receivedStr := string(value.Data)
								//s := strings.Split(receivedStr, ":")
								//ip, readingData := s[0]+":"+s[1], s[2]

								globals.Dbg.PrintfTrace("%s -->  [AMQP] DeviceUUID:     %+v \n", globals.LogTag, value.DeviceUuid)
								globals.Dbg.PrintfTrace("%s -->  [AMQP] EndpointUUID:   %+v \n", globals.LogTag, value.EndpointUuid)
								globals.Dbg.PrintfTrace("%s -->  [AMQP] Payload:        %+v \n", globals.LogTag, value.Payload)
								globals.Dbg.PrintfTrace("%s -->  [AMQP] DeviceName:     %+v \n", globals.LogTag, value.DeviceName)
								globals.Dbg.PrintfTrace("%s -->  [AMQP] Command:        %+v \n", globals.LogTag, value.Command)
								globals.Dbg.PrintfTrace("%s -->  [AMQP] CoapData:       %+v \n", globals.LogTag, value.Data)
								globals.Dbg.PrintfTrace("%s -->  [AMQP] MessageID:      %+v \n", globals.LogTag, amqpMessageCount)
								//Parsing coap data
								var x coap.Message
								err := x.UnmarshalBinary(value.Data)

								globals.Dbg.PrintfTrace("%s -->  [COAP] type:           %+v \n", globals.LogTag, x.Type)
								globals.Dbg.PrintfTrace("%s -->  [COAP] code:           %+v \n", globals.LogTag, x.Code)
								globals.Dbg.PrintfTrace("%s -->  [COAP] MessageID:      %+v \n", globals.LogTag, x.MessageID)
								globals.Dbg.PrintfTrace("%s -->  [COAP] payload:        %+v \n", globals.LogTag, x.Payload)
								globals.Dbg.PrintfTrace("%s -->  [COAP] Token:          %+v \n", globals.LogTag, string(x.Token))
								globals.Dbg.PrintfTrace("%s -->  [COAP] Destination:    %+v \n\n", globals.LogTag, x.Option(11))
								//var dest string
								ip := x.Option(11).(string)
								// (2) Send readingData to coap server
								req := coap.Message{
									Type:      x.Type,
									Code:      x.Code,
									MessageID: x.MessageID,
									Payload:   []byte(x.Payload),
									Token:     []byte(x.Token),
								}
								//add to gRouter
								gRouter[string(x.MessageID)] = value.DeviceUuid

								req.SetOption(coap.ETag, "SARA")
								req.SetOption(coap.MaxAge, 3)
								req.SetPathString("/test")
								fmt.Println("COAP Request ---> ", req)

								c, err := coap.Dial("udp", ip)
								if err != nil {
									log.Fatalf("Error dialing: %v", err)
								}

								// (3) Receive coap server reply
								coapServerReply, err := c.Send(req)

								if err != nil {
									log.Fatalf("Error sending request: %v", err)
								}
								if coapServerReply != nil {
									fmt.Println("COAP Reply <--- ", coapServerReply, "\n")

									//Finding the UUID ot rge destination device
									thisUUID, ok := gRouter[string(coapServerReply.MessageID)]

									if !ok {
										log.Fatalf("deviceUUID not found !!!")
									} else {
										delete(gRouter, string(coapServerReply.MessageID))
									}
									globals.Dbg.PrintfTrace("%s <--  messageID:   %+v \n", globals.LogTag, coapServerReply.MessageID)
									globals.Dbg.PrintfTrace("%s <--  DeviceUUID:  %+v \n", globals.LogTag, thisUUID)
									globals.Dbg.PrintfTrace("%s <--  Payload:     %+v \n\n", globals.LogTag, string(coapServerReply.Payload))

									// (4) Send to the device
									msgToDevice := AmqpMessage{
										DeviceUuid:   thisUUID,
										EndpointUuid: 4,
									}
									for _, v := range coapServerReply.Payload {
										msgToDevice.Payload = append(msgToDevice.Payload, int(v))
									}
									// Send the TSW server
									DlMsgs <- msgToDevice
								}
							}
						}
					case *error:
						{ // If an error has occurred, drop out of the loop
							globals.Dbg.PrintfTrace("%s [server] --> AMQP channel error received (%s), dropping out...\n", globals.LogTag, (*value).Error())
							return
						}
					default:
						{
							globals.Dbg.PrintfError("%s [server] --> unknown message type on AMQP UlMsg channel, dropping out...\n\n%+v\n", globals.LogTag, msg)
							return
						} // case
					} // switch
					//if isOpen
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
	//fmt.Println("username: ", username)
	//fmt.Println("amqpAddress: ", amqpAddress)

	// Open the queue and then begin processing messages
	// If we drop out of the processing function, wait
	// a little while and try again
	for {
		fmt.Printf("----------------------------------------------------------------------------------\n")
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

	// Set up the maps here. Do this here so that if we restart
	// the AMQP queue the data is retained
	//deviceDecodeStateList = make(map[string]*DeviceTotalsState)
	deviceTtValuesList = make(map[string]*TtValues)

	// Handle all the Amqp messages
	go doAmqp(username, amqpAddress)

	// Set up logging
	log.SetFlags(log.LstdFlags)

	// --- HAMED's FOWARDING THINGY ---
	forever := make(chan bool)
	<-forever
}

/* End Of File */
