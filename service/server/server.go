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
    "encoding/json"
    "fmt"
    "github.com/brettlangdon/forge"
    "github.com/robmeades/utm/service/utilities"
    "log"
	"net/http"
    //"os"
    "time"
    "github.com/codegangsta/negroni"
    "github.com/davecgh/go-spew/spew"
    "github.com/goincremental/negroni-sessions"
    "github.com/goincremental/negroni-sessions/cookiestore"
    "github.com/robmeades/utm/service/routes"
)

//--------------------------------------------------------------------
// Types 
//--------------------------------------------------------------------

// A message expected back from a device
type ExpectedMsg struct {
    TimeStarted  time.Time
    ResponseId   ResponseTypeEnum
}

// The list of messages expected back from a device
type ExpectedMsgList []ExpectedMsg

// Conection details for a device
type Connection struct {
	DeviceUuid    string
	DeviceName    string
    LastHeardFrom time.Time
}

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

// Level of debug required
const Dbg utilities.DebugLevel = utilities.DEBUG_TRACE

// Server details
const configurationFile string = "config.cfg"

// Log  prefix so that we can tell who we are
var logTag string = "UTM-API"

// TODO
var displayRow = &DisplayRow{}

// A list of expected response messages against each device
var deviceExpectedMsgList map[string]ExpectedMsgList

// Downlink channel to device
var downlinkMessages chan<- AmqpMessage

// Count of AMQP messages received
var amqpCount int = 0

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Fail func
func failOnError (err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
	    panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

// TODO
func getLatestState (response http.ResponseWriter, request *http.Request) *utilities.Error {
	// Ensure this is a GET request
	if (request.Method != "GET") || (request.Method == "") {
        Dbg.PrintfError("%s --> received unsupported REST request %s %s.\n", logTag, request.Method, request.URL)
        return utilities.ClientError("unsupported method", http.StatusBadRequest)
	}

	// Get the latest state; only one response will come back before the requesting channel is closed
	get := make(chan LatestState)
	dataTableCmds <- &get
	state := <-get

	//Gradually group our units
	//AddUuidToMap(state.LatestDisplayRow)

	//Recyle map with new state
	//uuidMap = RecycleMap(uuidMap, state)

	//Copy unit into slice, for encoding
	//uuidSlice = ConvertMapToSlice(uuidMap)

	// Send the requested data
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	err := json.NewEncoder(response).Encode(uuidSlice)
	if err != nil {
		Dbg.PrintfError("%s --> received REST request %s but attempting to serialise the result:\n%s\n...yielded error %s.\n", logTag, request.URL, spew.Sdump(state), err.Error())
		return utilities.ServerError(err)
	}
	//log.Printf("%s Received rest request %s and responding with %s\n", logTag, request.URL, spew.Sdump(state))

	return nil
}

// TODO
func setReportingInterval(response http.ResponseWriter, request *http.Request) *utilities.Error {
	// Ensure this is a POST request
	if request.Method != "POST" {
		Dbg.PrintfError("%s --> received unsupported REST request %s %s.\n", logTag, request.Method, request.URL)
		return utilities.ClientError("unsupported method", http.StatusBadRequest)
	}

	// Get the minutes interval
	var mins uint32
	err := json.NewDecoder(request.Body).Decode(&mins)
	if err != nil {
		Dbg.PrintfError("%s --> unable to extract the reporting interval from request %s: %s.\n", logTag, request.URL, err.Error())
		return utilities.ClientError("unable to decode reporting interval", http.StatusBadRequest)
	}

	// Encode and enqueue the requested data
	// TODO err = encodeAndEnqueueReportingInterval(uint32(mins))
	if err != nil {
		Dbg.PrintfError("%s --> unable to encode and enqueue reporting interval for UTM-API %s\n", logTag, request.URL)
		return utilities.ClientError("unable to encode and enqueue reporting interval", http.StatusBadRequest)
	}

	// Success
	response.WriteHeader(http.StatusOK)
	return nil
}

// Process messages from the AMQP queues
func processAmqp(username, amqpAddress string) {

	// create a queue and bind it with relevant routing keys
	// amqpAddress := utilities.EnvStr("AMQP_ADDRESS")
	// username := utilities.EnvStr("UNAME")
	// ueGuid = utilities.EnvStr("UEGUID")

	q, err := OpenQueue(username, amqpAddress)

	failOnError(err, "Queue")
	defer q.Close()
	downlinkMessages = q.Downlink

	for {
		amqpCount = amqpCount + 1
		Dbg.PrintfInfo("\n=====================> processing datagram number %v in AMQP channel =====================================\n", amqpCount)

		select {
			case msg := <-q.Msgs:
				Dbg.PrintfTrace("%s --> decoded msg:\n\n%+v\n\n", logTag, msg)

				switch value := msg.(type) {
					case *AmqpReceiveMessage:
						Dbg.PrintfTrace("%s --> is receive.\n", logTag)

					case *AmqpResponseMessage:
						Dbg.PrintfTrace("%s --> is response.\n", logTag)
						if value.Command == "UART_data" {
							// UART data from the UTM-API which needs to be decoded
                			// Decode the messages
							msgs := decode(value.Data, value.DeviceUuid)

							// Pass the messages to the state table for recording
							// and pass them to processing to see if any responses
							// are required
							if msgs != nil {
    							for _, msg := range msgs {
                            		dataTableCmds <- msg							    
        							processMsgs <- msg
    							}
    						}							

							// Send the datatable a message indicating that this device
							// has been heard from							
                			dataTableCmds <- &Connection {
                			    DeviceUuid:    value.DeviceUuid,
                			    DeviceName:    value.DeviceName,
                			    LastHeardFrom: time.Now(),
              			    }
                			
							// Get the amount of uplink data and send it on to the data table
							now := time.Now()
							dataTableCmds <- &DataVolume{
								UplinkTimestamp: &now,
								UplinkBytes:     uint64(len(value.Data)),
							}
							row.UlastMsgReceived = &now
							row.UTotalBytes = uint64(len(value.Data))
							row.TotalMsgs = row.TotalMsgs + uint64(len(value.Data))
							row.TotalBytes = row.UTotalBytes + row.DTotalBytes
							row.TotalMsgs = row.UTotalMsgs + row.DTotalMsgs + row.TotalMsgs
						}

					case *AmqpErrorMessage:
						Dbg.PrintfTrace("%s--> is error.\n", logTag)

					default:
						Dbg.PrintfTrace("%s --> message type: %+v.\n", logTag, msg)
						log.Fatal(logTag, "invalid message type.")
				}
		}
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
	deviceExpectedMsgList = make(map[string]ExpectedMsgList)
	// And a timer to remove crud from it
    removeOldExpectedMsgs := time.NewTicker (time.Minute * 10)
    
    // Remove old stuff from the expected message list on a tick
    go func() {
        for _ = range removeOldExpectedMsgs.C {
            for uuid, expectedMsgList := range deviceExpectedMsgList {
                var x = 0                
                for x < len(expectedMsgList) {
                    if time.Now().After (expectedMsgList[x].TimeStarted.Add(time.Hour)) {
    		            expectedMsgList = append(expectedMsgList[:x], expectedMsgList[x + 1:] ...)
                        Dbg.PrintfTrace("%s --> giving up after waiting > 1 hour for %d from device %s.\n", logTag, expectedMsgList[x].ResponseId, uuid)
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

	fmt.Printf("######################################################################################################\n")
	fmt.Printf("UTM-API service (%s) REST interface listening on %s.\n", logTag, amqpAddress)

	store := cookiestore.New([]byte("secretkey789"))
	router := routes.LoadRoutes()

	router.Handle("/latestState", utilities.Handler(getLatestState))
	router.Handle("/reportingInterval", utilities.Handler(setReportingInterval))

	n := negroni.Classic()
	static := negroni.NewStatic(http.Dir("static"))
	static.Prefix = "/static"
	n.Use(static)
	//n.Use(negroni.HandlerFunc(system.MgoMiddleware))
	n.Use(sessions.Sessions("global_session_store", store))
	n.UseHandler(router)
	n.Run(port)
}

/* End Of File */