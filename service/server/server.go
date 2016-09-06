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
	"github.com/u-blox/utm/service/globals"
	"log"
	"strings"
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
	DeviceName      string
	Msgs            int
	Bytes           int
	Totals          *TotalsState
	ExpectedMsgList *[]ExpectedMsg // Will be nil for uplink
}

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
type Connection struct {
	DeviceUuid      string `bson:"DeviceUuid" json:"DeviceUuid"`
	DeviceName      string
	RecentUl        bool
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

var gRouter = make(map[string]string)

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Given a UUID and a pointer to the decode state for that UUID, get the
// encode state and the traffic test state and pass all of this along
// to the datatable so that it is kept up to date.
func updateDatatableState(uuid string, decodeState *DeviceTotalsState, recentUl bool) {
	var ttDlBytes uint32
	var ttDlDatagrams uint32
	var ttTimeLastDl time.Time

	// Keep track of traffic test progress
	getTrafficTestContext := &DeviceTrafficTestContextGet{
		DeviceUuid: uuid,
	}
	getTrafficTestContext.Context = make(chan TrafficTestContext)
	trafficTestChannel <- getTrafficTestContext
	trafficTestContext, isOpen := <-getTrafficTestContext.Context
	if isOpen {
		ttDlBytes = trafficTestContext.DlBytesTotal
		ttDlDatagrams = trafficTestContext.DlDatagramsTotal
		ttTimeLastDl = trafficTestContext.TimeLastDl
		// Send it all off to the datatable
		dataTableChannel <- &trafficTestContext
	}

	// Ask the process channel for the encode status now
	getEncodeState := &DeviceEncodeStateGet{
		DeviceUuid: uuid,
	}
	getEncodeState.State = make(chan DeviceTotalsState)
	processMsgsChannel <- getEncodeState
	encodeState, isOpen := <-getEncodeState.State

	// Only do this if the channel wasn't closed prematurely
	// (e.g. due to an unknown UE being requested, which might
	// happen if it's sending rubbish and so no messages
	// will have been decoded)
	if isOpen {
		// Add in the traffic test encoded stuff, if there is
		// any (the decode counts for the traffic test stuff
		// is already taken into account by the decoder)
		if ttTimeLastDl.After(encodeState.Timestamp) {
			encodeState.Timestamp = ttTimeLastDl
		}
		encodeState.Totals.Msgs += int(ttDlDatagrams)
		encodeState.Msgs += int(ttDlDatagrams)
		encodeState.Totals.Bytes += int(ttDlBytes)
		encodeState.Bytes += int(ttDlBytes)
		// Send the datatable a message with connection
		// data for this device plus the totals for all
		ulTotals := TotalsState{
			Timestamp: totalsDecodeState.Timestamp,
			Msgs:      totalsDecodeState.Msgs,
			Bytes:     totalsDecodeState.Bytes,
		}
		dlTotals := TotalsState{
			Timestamp: encodeState.Timestamp,
			Msgs:      encodeState.Totals.Msgs,
			Bytes:     encodeState.Totals.Bytes,
		}
		// Send it off
		dataTableChannel <- &Connection{
			DeviceUuid: uuid,
			DeviceName: decodeState.DeviceName,
			RecentUl:   recentUl,
			UlDevice: TotalsState{
				Timestamp: decodeState.Timestamp,
				Msgs:      decodeState.Msgs,
				Bytes:     decodeState.Bytes,
			},
			DlDevice: TotalsState{
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

// process RECEIVED Msgs
// Process uplink messages from the AMQP queue until it is closed or
// an error is flagged
func processUlAmqpMsgs(q *Queue) {
	fmt.Println("... processUlAmqpMsgs")
	statusTick := time.NewTicker(time.Minute)

	// Establish the status of the process loop and of traffic test
	// for all devices on a background timer, just in case there are
	// no datagrams arriving from them
	go func() {
		for _ = range statusTick.C {
			for uuid, decodeState := range deviceDecodeStateList {
				updateDatatableState(uuid, decodeState, false)
			}
		}
	}()

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
							globals.Dbg.PrintfTrace("%s [server] --> decoded AMQP JSON msg:\n\n%+v\n\n", globals.LogTag, msg)
							if value.Command == "UART_data" {
								// UART data from the UTM-API which needs to be decoded
								// If the device is not known, add it
								value.DeviceUuid = strings.ToLower(value.DeviceUuid) // force to lower case from the very beginning

								// HAMED SIASI
								globals.Dbg.PrintfTrace("Device UUID: %s\n", value.DeviceUuid)
								globals.Dbg.PrintfTrace("Device NAME: %s\n", value.DeviceName)
								globals.Dbg.PrintfTrace("Device DATA: %s\n", value.Data)
								//gRouter["DestinationIP"]="UUID"
								//gRouter["0.0.0.0"] = value.DeviceUuid

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
	fmt.Println("... doAmqp")
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

	fmt.Println("... RUN")
	// First, parse the configuration file
	settings, err := forge.ParseFile(configurationFile)
	if err != nil {
		panic(err)
	}
	amqp, err := settings.GetSection("amqp")
	username, err := amqp.GetString("uname")
	amqpAddress, err := amqp.GetString("amqp_address")
	//host, err := settings.GetSection("host")
	//port, err := host.GetString("port")

	// Set up the maps here. Do this here so that if we restart
	// the AMQP queue the data is retained
	deviceDecodeStateList = make(map[string]*DeviceTotalsState)
	deviceTtValuesList = make(map[string]*TtValues)

	// Handle all the Amqp messages
	go doAmqp(username, amqpAddress)

	// Set up logging
	log.SetFlags(log.LstdFlags)

	// --- HAMED's FOWARDING THINGY ---

	fmt.Println("... forever")
	forever := make(chan bool)
	<-forever
}

/* End Of File */
