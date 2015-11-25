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
	//"github.com/gorilla/mux"
	"github.com/robmeades/utm/service/utilities"
    "log"
	"net/http"
    //"os"
    "time"
    "github.com/codegangsta/negroni"
    "github.com/davecgh/go-spew/spew"
    "github.com/goincremental/negroni-sessions"
    "github.com/goincremental/negroni-sessions/cookiestore"
    "github.com/robmeades/utm/service/models"
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
	Timestamp     time.Time
    UlMsgs        int
    UlBytes       int
    DlMsgs        int
    DlBytes       int
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

// A list of expected response messages against each device
var deviceExpectedMsgList map[string]ExpectedMsgList

// Downlink channel to device
var downlinkMessages chan<- AmqpMessage

// Count of AMQP messages received
var amqpMessageCount int

// Count of the number of times we've (re)starte AMQP
var amqpRetryCount int

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

/// Fail func
func failOnError (err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
	    panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

func loginHandler(response http.ResponseWriter, request *http.Request) {

	email := request.FormValue("email")
	password := request.FormValue("password")
	session := sessions.GetSession(request)

	db := utilities.GetDB(request)
	user := new(models.User)
	err := user.Authenticate(db, email, password)
	if err == nil {
		session.Set("user_id", user.ID.Hex())
		session.Set("user_company", user.Company)
		session.Set("user_email", user.Email)
		//Fmt.Fprintf(response, "User %s success!\n", session.Get("user_email"))
		http.Redirect(response, request, "/", 302)
	}
}

func registerHandler(response http.ResponseWriter, request *http.Request) {

	company := request.FormValue("company_name")
	firstName := request.FormValue("user_firstName")
	lastName := request.FormValue("user_lastName")
	email := request.FormValue("email")
	password := request.FormValue("password")

	db := utilities.GetDB(request)
	user := new(models.User)

	user.NewUser(db, company, firstName, lastName, email, password)
	//Fmt.Fprintf(response, "User %s created successfully!\n", firstName)
	http.Redirect(response, request, "/login", 302)
}

/// Get the summary data for the front page
func getFrontPageData (response http.ResponseWriter, request *http.Request) *utilities.Error {
	// Ensure this is a GET request
	if (request.Method != "GET") || (request.Method == "") {
        Dbg.PrintfError("%s --> received unsupported REST request %s %s.\n", logTag, request.Method, request.URL)
        return utilities.ClientError("unsupported method", http.StatusBadRequest)
	}

    displayData := displayFrontPageData()
	// Send the requested data
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	err := json.NewEncoder(response).Encode(displayData)
	if err != nil {
		Dbg.PrintfError("%s --> received REST request %s but attempting to serialise the result:\n%s\n...yielded error %s.\n", logTag, request.URL, spew.Sdump(displayData), err.Error())
		return utilities.ServerError(err)
	}

	return nil
}

// Process datagrams from the AMQP queue
func processDatagrams(q *Queue) {
    
	for {
		amqpMessageCount++
    	select {
    		case msg, ok := <-q.Msgs:
                if (!ok) {
                    return
                }
    			Dbg.PrintfTrace("%s --> decoded msg:\n\n%+v\n\n", logTag, msg)
    
    			switch value := msg.(type) {
    				case *AmqpResponseMessage:
    					Dbg.PrintfTrace("%s --> is response.\n", logTag)
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
                			dataTableChannel <- &Connection {
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
    					Dbg.PrintfTrace("%s --> AMQP error received (%s), dropping out...\n", logTag, (*value).Error())
    					return
    
    				default:
    					Dbg.PrintfTrace("%s --> message type: %+v.\n", logTag, msg)
    					log.Fatal(logTag, "invalid message type.")
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
    	fmt.Printf("UTM-API service (%s) REST interface opening %s...\n", logTag, amqpAddress)

    	q, err := OpenQueue(username, amqpAddress)

    	if err == nil {
        	defer q.Close()
    
        	fmt.Printf("%s --> connection opened.\n", logTag)
    	
        	downlinkMessages = q.Downlink
    
            processDatagrams (q)    
    	} else {
    		Dbg.PrintfTrace("%s --> error opening AMQP queue (%s).\n", logTag, err.Error())    	    
    	}
    	
		amqpRetryCount++;
		Dbg.PrintfTrace("%s --> waiting before trying again...\n", logTag)    	        	
    	time.Sleep (time.Second * 10)	
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

	store := cookiestore.New([]byte("secretkey789"))
	router := routes.LoadRoutes()

	router.Handle("/frontPageData", utilities.Handler(getFrontPageData))
	router.HandleFunc("/register", registerHandler)
	router.HandleFunc("/login", loginHandler)
	
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