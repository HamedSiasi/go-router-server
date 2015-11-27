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
	//"github.com/gorilla/mux"
	"github.com/robmeades/utm/service/globals"
	"github.com/robmeades/utm/service/utilities"
	"log"
	"net/http"
	//"os"
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

// A message expected back from a device
type ExpectedMsg struct {
	TimeStarted time.Time
	ResponseId  ResponseTypeEnum
}

// Conection details for a device
type Connection struct {
	DeviceUuid string  `bson:"DeviceUuid" json:"DeviceUuid"`
	DeviceName string
	Timestamp  time.Time
	UlMsgs     int
	UlBytes    int
	DlMsgs     int
	DlBytes    int
}

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

// Server details
const configurationFile string = "config.cfg"

// A list of expected response messages against each device
var deviceExpectedMsgList map[string]*[]ExpectedMsg

// Downlink channel to device
var downlinkMessages chan<- AmqpMessage

// Count of AMQP messages received
var amqpMessageCount int

// Count of the number of times we've (re)starte AMQP
var amqpRetryCount int

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Process datagrams from the AMQP queue
func processDatagrams(q *Queue) {

	for {
		amqpMessageCount++
		select {
		case msg, ok := <-q.Msgs:
			if !ok {
				return
			}
			globals.Dbg.PrintfTrace("%s [server] --> decoded msg:\n\n%+v\n\n", globals.LogTag, msg)

			switch value := msg.(type) {
			case *AmqpResponseMessage:
				globals.Dbg.PrintfTrace("%s [server] --> is response.\n", globals.LogTag)
				if value.Command == "UART_data" {
					savedUlMsgs := totalUlMsgs
					savedUlBytes := totalUlBytes
					savedDlMsgs := totalDlMsgs
					savedDlBytes := totalDlBytes
					// UART data from the UTM-API which needs to be decoded
					// Decode the messages
					dlMsgs := decode(value.Data, value.DeviceUuid)

					// Pass the messages to the data table for recording
					// and pass them to processing to see if any responses
					// are required
					if dlMsgs != nil {
						for _, dlMsg := range dlMsgs {
							dataTableChannel <- dlMsg
							processMsgs <- dlMsg
						}
					}

					// Send the datatable a message with connection
					// data for this device
					dataTableChannel <- &Connection{
						DeviceUuid: value.DeviceUuid,
						DeviceName: value.DeviceName,
						Timestamp:  time.Now(),
						UlMsgs:     totalUlMsgs - savedUlMsgs,
						UlBytes:    totalUlBytes - savedUlBytes,
						DlMsgs:     totalDlMsgs - savedDlMsgs,
						DlBytes:    totalDlBytes - savedDlBytes,
					}
				}

			case *error:
				// If an error has occurred, drop out of the loop
				globals.Dbg.PrintfTrace("%s [server] --> AMQP error received (%s), dropping out...\n", globals.LogTag, (*value).Error())
				return

			default:
				globals.Dbg.PrintfTrace("%s [server] --> message type: %+v.\n", globals.LogTag, msg)
				log.Fatal(globals.LogTag, "invalid message type.")
			}
		}
	}
}

// Process messages from the AMQP queues
func processAmqp(username, amqpAddress string) {

	// Open the queue and then begin processing messages
	// If we drop out of the processing function, wait
	// a little while and try again
	for {
		fmt.Printf("######################################################################################################\n")
		fmt.Printf("UTM-API service (%s) REST interface opening %s...\n", globals.LogTag, amqpAddress)

		q, err := OpenQueue(username, amqpAddress)

		if err == nil {
			defer q.Close()

			fmt.Printf("%s [server] --> connection opened.\n", globals.LogTag)

			downlinkMessages = q.Downlink

			// The meat is in here
			processDatagrams(q)

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

	// Set up the device expected message list map
	deviceExpectedMsgList = make(map[string]*[]ExpectedMsg)
	// And a timer to remove crud from it
	removeOldExpectedMsgs := time.NewTicker(time.Minute * 10)

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

	// Process Amqp messages
	go processAmqp(username, amqpAddress)

	// Set up logging
	log.SetFlags(log.LstdFlags)

	store := cookiestore.New([]byte("secretkey789"))
	router := routes.LoadRoutes()

	router.Handle("/frontPageData", utilities.Handler(getFrontPageData))
	router.HandleFunc("/register", registerHandler)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/mode", ShowModeHandler)
	router.HandleFunc("/display", ShowDisplayHandler)
	router.HandleFunc("/logout", LogoutHandler)

	sendMsg := ClientSendMsg{}
	router.HandleFunc("/sendMsg", sendMsg.Send).Methods("POST")

	n := negroni.Classic()
	static := negroni.NewStatic(http.Dir("static"))
	static.Prefix = "/static"
	n.Use(static)
	n.Use(negroni.HandlerFunc(system.MgoMiddleware))
	n.Use(sessions.Sessions("global_session_store", store))
	n.UseHandler(router)
	n.Run(port)
}

/* End Of File */
