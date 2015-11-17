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

const configurationFile string = "config.cfg"
var logTag string = "UTM-API"
var downlinkMessages chan<- AmqpMessage
var ueGuid string
var amqpCount int
var displayRow = &DisplayRow{}

//var UuidsMap = map[string]*DisplayRow{}

func failOnError (err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
	    panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

func getLatestState (response http.ResponseWriter, request *http.Request) *utilities.Error {
	// Ensure this is a GET request
	if (request.Method != "GET") || (request.Method == "") {
        log.Printf("%s --> received unsupported REST request %s %s.\n", logTag, request.Method, request.URL)
        return utilities.ClientError("unsupported method", http.StatusBadRequest)
	}

	// Get the latest state; only one response will come back before the requesting channel is closed
	get := make(chan LatestState)
	stateTableCmds <- &get
	state := <-get

	//Graddually group our units
	AddUuidToMap(state.LatestDisplayRow)

	//Recyle map with new state
	uuidMap = RecycleMap(uuidMap, state)

	//Copy unit into slice, for encoding
	uuidSlice = ConvertMapToSlice(uuidMap)

	// Send the requested data
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	err := json.NewEncoder(response).Encode(uuidSlice)
	if err != nil {
		log.Printf("%s --> recieved REST request %s but attempting to serialise the result:\n%s\n...yielded error %s.\n", logTag, request.URL, spew.Sdump(state), err.Error())
		return utilities.ServerError(err)
	}
	//log.Printf("%s Received rest request %s and responding with %s\n", logTag, request.URL, spew.Sdump(state))

	return nil
}

func AddUuidToMap(row *DisplayRow) *utilities.Error {
	if row != nil {
		_, ok := uuidMap[row.Uuid]
		if !ok {
			uuidMap[row.Uuid] = row
		}
	}
	return nil
}
func RecycleMap(oldMap map[string]*DisplayRow, newState LatestState) map[string]*DisplayRow {
	if len(oldMap) > 3 {
		_, ok := oldMap[newState.LatestDisplayRow.Uuid]
		if ok {
			//Delete old state
			delete(oldMap, newState.LatestDisplayRow.Uuid)
			//Add new state
			oldMap[newState.LatestDisplayRow.Uuid] = newState.LatestDisplayRow
		}
	}

	return oldMap
}
func ConvertMapToSlice(uidMap map[string]*DisplayRow) []*DisplayRow {
	var sdisplay []*DisplayRow
	for _, v := range uidMap {
		sdisplay = append(sdisplay, v)
	}

	return sdisplay
}

//DownLink (DL) messages
func setReportingInterval(response http.ResponseWriter, request *http.Request) *utilities.Error {
	// Ensure this is a POST request
	if request.Method != "POST" {
		log.Printf("%s --> received unsupported REST request %s %s.\n", logTag, request.Method, request.URL)
		return utilities.ClientError("unsupported method", http.StatusBadRequest)
	}

	// Get the minutes interval
	var mins uint32
	err := json.NewDecoder(request.Body).Decode(&mins)
	if err != nil {
		log.Printf("%s --> unable to extract the reporting interval from request %s: %s.\n", logTag, request.URL, err.Error())
		return utilities.ClientError("unable to decode reporting interval", http.StatusBadRequest)
	}

	// Encode and enqueue the requested data
	err = encodeAndEnqueueReportingInterval(uint32(mins))
	if err != nil {
		log.Printf("%s --> unable to encode and enqueue reporting interval for UTM-API %s\n", logTag, request.URL)
		return utilities.ClientError("unable to encode and enqueue reporting interval", http.StatusBadRequest)
	}

	// Success
	response.WriteHeader(http.StatusOK)
	return nil
}

func processAmqp(username, amqpAddress string) {

	//ueGuid, err := amqp.GetString("ueguid")

	// create a queue and bind it with relevant routing keys
	// amqpAddress := utilities.EnvStr("AMQP_ADDRESS")
	// username := utilities.EnvStr("UNAME")
	// ueGuid = utilities.EnvStr("UEGUID")

	q, err := OpenQueue(username, amqpAddress)

	failOnError(err, "Queue")
	defer q.Close()
	downlinkMessages = q.Downlink

	stateTableCmds <- &Connection{Status: "Disconnected"}

	connState := "Disconnected"
	connected := connState

	for {
		amqpCount = amqpCount + 1
		log.Println()
		fmt.Printf("\n=====================> processing datagram no (%v) in AMQP channel =====================================\n", amqpCount)

		//time.Sleep(time.Second * 10)
		select {
			case <-time.After(30 * time.Minute):
				connected = "Disconnected"

			case msg := <-q.Msgs:
				log.Println(logTag, "--> decoded msg:", msg)

				switch value := msg.(type) {
					case *AmqpReceiveMessage:
						log.Println(logTag, "--> is receive")

					case *AmqpResponseMessage:
						log.Println(logTag, "--> is response")
						if value.Command == "UART_data" || value.Command == "LOOPBACK_data" {
							// UART data from the UTM-API which needs to be decoded
							// NOTE: some old data extracted from logs is loopback_data. FIXME: remove this.
							decode(value.Data)
							// Get the amount of uplink data and send it on to the data table
							now := time.Now()
							stateTableCmds <- &DataVolume{
								UplinkTimestamp: &now,
								UplinkBytes:     uint64(len(value.Data)),
							}
							row.UlastMsgReceived = &now
							row.UTotalBytes = uint64(len(value.Data))
							row.TotalMsgs = row.TotalMsgs + uint64(len(value.Data))
							row.TotalBytes = row.UTotalBytes + row.DTotalBytes
							row.TotalMsgs = row.UTotalMsgs + row.DTotalMsgs + row.TotalMsgs
							connected = "CONNECTED"
						}

					case *AmqpErrorMessage:
						log.Println(logTag, "--> is error")

					default:
						log.Printf("%s --> message type: %+v.\n", logTag, msg)
						log.Fatal(logTag, "invalid message type.")
				}
		}

		if connState != connected {
			connState = connected
			log.Printf("%s --> sending new connection state: %s.\n", logTag, connState)
			stateTableCmds <- &Connection{Status: connState}
		}

	}
	amqpCount = 0
}

func Run() {
    // First, parse the configuration file
    settings, err := forge.ParseFile(configurationFile)
    if err != nil {
        panic(err)
	}

	amqp, err := settings.GetSection("amqp")
	username, err := amqp.GetString("uname")
	amqpAddress, err := amqp.GetString("amqp_address")
	//ueGuid, err := amqp.GetString("ueguid")

	host, err := settings.GetSection("host")
	port, err := host.GetString("port")

	// Process Amqp messages
	go processAmqp(username, amqpAddress)

	log.SetFlags(log.LstdFlags | log.Llongfile)

	log.Printf("UTM-API service (%s) REST interface listening on %s.\n", logTag, amqpAddress)

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