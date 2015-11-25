/* Authentication elements of the UTM server.
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

package controllers

/*
import (
    "encoding/json"
    "github.com/robmeades/utm/service/globals"
    "github.com/robmeades/utm/service/utilities"
    "github.com/robmeades/utm/service/models"
    "net/http"
)

//--------------------------------------------------------------------
// Types 
//--------------------------------------------------------------------

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

func (msg SendDlMsg) Set(response http.ResponseWriter, request *http.Request) {
     decoder := json.NewDecoder(request.Body)
     var msg ClientHeartbeatSetReqDlMsg
     err := decoder.Decode(&msg)
     if err != nil {
         globals.Dbg.PrintfTrace ("%s [dl_msgs] --> JSON says:\n\n+v\n", globals.LogTag, msg)
     }
}

/// Send a ping request to a device
func sendPingReq (response http.ResponseWriter, request *http.Request) *globals.Error {
	err := utilities.ValidateGetRequest (request)
	if err == nil {
        // Send the ping request
        //encodeAndEnqueue (&PingReqDlMsg{}, "61d25940-5307-11e5-80a4-c549fb2d6313") // TODO
        
    	// Set up the response
    	utilities.MakeOkResponse (response)
    }
	
	return err
}

/// Send a heartbeat set request to a device
func sendHeartbeatSetReq (response http.ResponseWriter, request *http.Request) *globals.Error {
	err := utilities.ValidateGetRequest (request)
	if err == nil {
        heartbeatSetReq := &HeartbeatSetReqDlMsg {
            HeartbeatSeconds:   18,    // TODO
            HeartbeatSnapToRtc: false, // TODO       
        }
        // Send the ping request
        //encodeAndEnqueue (heartbeatSetReq, "61d25940-5307-11e5-80a4-c549fb2d6313") // TODO
        
    	// Set up the response
    	utilities.MakeOkResponse (response)
    }	

	return err
}

/// Send  reporting interval set request to a device
func sendReportingIntervalSetReq (response http.ResponseWriter, request *http.Request) *globals.Error {
	err := utilities.ValidateGetRequest (request)
	if err == nil {
        reportingIntervalSetReq := &ReportingIntervalSetReqDlMsg {
            ReportingInterval:   1,    // TODO
        }
        // Send the ping request
        //encodeAndEnqueue (reportingIntervalSetReq, "61d25940-5307-11e5-80a4-c549fb2d6313") // TODO
        
    	// Set up the response
    	utilities.MakeOkResponse (response)
    }	

	return err
}

/// Send intervals get request to a device
func sendIntervalsGetReq (response http.ResponseWriter, request *http.Request) *globals.Error {
	err := utilities.ValidateGetRequest (request)
	if err == nil {

        // Send the ping request
        //encodeAndEnqueue (&IntervalsGetReqDlMsg{}, "61d25940-5307-11e5-80a4-c549fb2d6313") // TODO
        
    	// Set up the response
    	utilities.MakeOkResponse (response)
    }	

	return err
}

*/
/* End Of File */
