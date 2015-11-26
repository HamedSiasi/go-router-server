/* Structs to send DL messages via the client interface from the UTM 
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
    "github.com/robmeades/utm/service/globals"
    "github.com/robmeades/utm/service/utilities"
    "net/http"
    "io/ioutil"
    "time"
    "errors"
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

type ClientSendMsg struct {
    Uuid  string           `bson:"device_uuid" json:"device_uuid"`
    SendMsgType string     `bson:"type" json:"type"`
    MsgBody interface{}    `bson:"body" json:"body"`
}

type Msg struct {
    MsgType ClientSendEnum `bson:"type" json:"type"`
    MsgBody interface{}    `bson:"body" json:"body"`
}

type ClientSendEnum uint32

const (
    CLIENT_SEND_NULL                       ClientSendEnum = iota
    CLIENT_SEND_PING                               
    CLIENT_SEND_REBOOT
    CLIENT_SEND_DATE_TIME_SET
    CLIENT_SEND_DATE_TIME_GET
    CLIENT_SEND_MODE_SET
    CLIENT_SEND_MODE_GET
    CLIENT_SEND_HEARTBEAT_SET
    CLIENT_SEND_REPORTING_INTERVAL_SET
    CLIENT_SEND_INTERVALS_GET
    CLIENT_SEND_MEASUREMENTS_GET
    CLIENT_SEND_MEASUREMENT_CONTROL_SET
    CLIENT_SEND_MEASUREMENTS_CONTROL_GET
    CLIENT_SEND_MEASUREMENTS_CONTROL_DEFAULTS_SET
    CLIENT_SEND_TRAFFIC_REPORT_GET
    CLIENT_SEND_TRAFFIC_TEST_MODE_PARAMETERS_SET
    CLIENT_SEND_TRAFFIC_TEST_MODE_PARAMETERS_GET
    CLIENT_SEND_TRAFFIC_TEST_MODE_REPORT_GET
    CLIENT_SEND_ACTIVITY_REPORT_GET
)

// TODO
// type ClientMeasurementControlGeneric struct
// type ClientMeasurementControlPowerState struct
// type ClientMeasurementControlUnion union

// ClientPingReq
type ClientPingReq struct {
    // Empty
}

// ClientRebootReq
type ClientRebootReq struct {
    SdCardNotRequired bool      `bson:"sd_card_not_required" json:"sd_card_not_required"`
    DisableModemDebug bool      `bson:"disable_modem_debug" json:"disable_modem_debug"`
    DisableButton     bool      `bson:"disable_button" json:"disable_button"`
    DisableServerPing bool      `bson:"disable_server_ping" json:"disable_server_ping"`
}

// ClientDateTimeSetReq
type ClientDateTimeSetReq struct {
    UtmTime           time.Time `bson:"time" json:"time"`
    SetDateOnly       bool      `bson:"set_date_only" json:"set_date_only"`
}

// ClientDateTimeGetReq
type ClientDateTimeGetReq struct {
    // Empty    
}

type ClientModeEnum uint32

const (
    CLIENT_MODE_NULL          ClientModeEnum = iota
    CLIENT_MODE_SELF_TEST
    CLIENT_MODE_COMMISSIONING
    CLIENT_MODE_STANDARD_TRX
    CLIENT_MODE_TRAFFIC_TEST
)

// ClientModeSetReq
type ClientModeSetReq struct {
    Mode        ClientModeEnum  `bson:"mode" json:"mode"`
}

// ClientModeGetReq
type ClientModeGetReq struct {
    // Empty    
}

// ClientIntervalsGetReq
type ClientIntervalsGetReq struct {
    // Empty    
}

// ClientReportingIntervalSetReq
type ClientReportingIntervalSetReq struct {
    ReportingInterval uint32  `bson:"reporting_interval" json:"reporting_interval"`
}

// ClientHeartbeatSetReq
type ClientHeartbeatSetReq struct {
    HeartbeatSeconds   uint32  `bson:"heartbeat_seconds" json:"heartbeat_seconds"`
    HeartbeatSnapToRtc bool    `bson:"heartbeat_snap_to_rtc" json:"heartbeat_snap_to_rtc"`
}

// ClientMeasurementsGetReq
type ClientMeasurementsGetReq struct {
    // Empty
}

// TODO
// type ClientMeasurementControlSetReq struct
// type ClientMeasurementsControlDefaultsSetReq struct

// ClientTrafficReportGetReq
type ClientTrafficReportGetReq struct {
    // Empty
}

// ClientTrafficTestModeParametersSetReq
type ClientTrafficTestModeParametersSetReq struct {
    NumUlDatagrams      uint32  `bson:"num_ul_datagrams" json:"num_ul_datagrams"`
    LenUlDatagram       uint32  `bson:"len_ul_datagram" json:"len_ul_datagram"`
    NumDlDatagrams      uint32  `bson:"num_dl_datagrams" json:"num_dl_datagrams"`
    LenDlDatagram       uint32  `bson:"len_dl_datagram" json:"len_dl_datagram"`
    TimeoutSeconds      uint32  `bson:"timeout_seconds" json:"timeout_seconds"`
    NoReportsDuringTest bool
}

// ClientTrafficTestModeParametersGetReq
type ClientTrafficTestModeParametersGetReq struct {
    // Empty
}

// ClientTrafficTestModeReportGetReq
type ClientTrafficTestModeReportGetReq struct {
    // Empty
}

// ClientActivityReportGetReq
type ClientActivityReportGetReq struct {
    // Empty
}

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

var clientSendEnumString map[string]ClientSendEnum = map[string]ClientSendEnum {
    "CLIENT_SEND_NULL":                              CLIENT_SEND_NULL,
    "CLIENT_SEND_PING":                              CLIENT_SEND_PING,
    "CLIENT_SEND_REBOOT":                            CLIENT_SEND_REBOOT,
    "CLIENT_SEND_DATE_TIME_SET":                     CLIENT_SEND_DATE_TIME_SET,
    "CLIENT_SEND_DATE_TIME_GET":                     CLIENT_SEND_DATE_TIME_GET,
    "CLIENT_SEND_MODE_SET":                          CLIENT_SEND_MODE_SET,
    "CLIENT_SEND_MODE_GET":                          CLIENT_SEND_MODE_GET,
    "CLIENT_SEND_HEARTBEAT_SET":                     CLIENT_SEND_HEARTBEAT_SET,
    "CLIENT_SEND_REPORTING_INTERVAL_SET":            CLIENT_SEND_REPORTING_INTERVAL_SET,
    "CLIENT_SEND_INTERVALS_GET":                     CLIENT_SEND_INTERVALS_GET,
    "CLIENT_SEND_MEASUREMENTS_GET":                  CLIENT_SEND_MEASUREMENTS_GET,
    "CLIENT_SEND_MEASUREMENT_CONTROL_SET":           CLIENT_SEND_MEASUREMENT_CONTROL_SET,
    "CLIENT_SEND_MEASUREMENTS_CONTROL_GET":          CLIENT_SEND_MEASUREMENTS_CONTROL_GET,
    "CLIENT_SEND_MEASUREMENTS_CONTROL_DEFAULTS_SET": CLIENT_SEND_MEASUREMENTS_CONTROL_DEFAULTS_SET,
    "CLIENT_SEND_TRAFFIC_REPORT_GET":                CLIENT_SEND_TRAFFIC_REPORT_GET,
    "CLIENT_SEND_TRAFFIC_TEST_MODE_PARAMETERS_SET":  CLIENT_SEND_TRAFFIC_TEST_MODE_PARAMETERS_SET,
    "CLIENT_SEND_TRAFFIC_TEST_MODE_PARAMETERS_GET":  CLIENT_SEND_TRAFFIC_TEST_MODE_PARAMETERS_GET,
    "CLIENT_SEND_TRAFFIC_TEST_MODE_REPORT_GET":      CLIENT_SEND_TRAFFIC_TEST_MODE_REPORT_GET,
    "CLIENT_SEND_ACTIVITY_REPORT_GET":               CLIENT_SEND_ACTIVITY_REPORT_GET,
}

var clientModeEnumString map[string]ClientModeEnum = map[string]ClientModeEnum {
    "CLIENT_MODE_NULL":           CLIENT_MODE_NULL,
    "CLIENT_MODE_SELF_TEST":      CLIENT_MODE_SELF_TEST,
    "CLIENT_MODE_COMMISSIONING":  CLIENT_MODE_COMMISSIONING,
    "CLIENT_MODE_STANDARD_TRX":   CLIENT_MODE_STANDARD_TRX,
    "CLIENT_MODE_TRAFFIC_TEST":   CLIENT_MODE_TRAFFIC_TEST,
}

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Send a message to the UTM on command from the client
func (msg *ClientSendMsg) Send(response http.ResponseWriter, request *http.Request) {
	err := utilities.ValidatePostRequest (request)
	
	if err != nil {
        response.WriteHeader(403)
        return
    } else {
        globals.Dbg.PrintfTrace ("%s [dl_msgs] --> sendMsg request received from client.\n", globals.LogTag)
        utilities.DumpRequest (request)
        body, err := ioutil.ReadAll(request.Body)
        if err != nil {
            response.WriteHeader(404)
            return
        } else {
            globals.Dbg.PrintfTrace ("%s [dl_msgs] --> request body is \"%s\".\n", globals.LogTag, body)
            var msgContainer ClientSendMsg
            err = json.Unmarshal(body, &msgContainer)
            if err != nil {
                globals.Dbg.PrintfTrace("%s [dl_msgs] --> received:\n\n%+v\n\n...which is not JSON decodable:\n\n%+v\n", globals.LogTag, body, err.Error())
                response.WriteHeader(404)
                return
            } else {                
                globals.Dbg.PrintfTrace ("%s [dl_msgs] --> JSON says:\n\n%+v\n", globals.LogTag, msgContainer)
                msgEnum := clientSendEnumString[msgContainer.SendMsgType]
                
                if msgEnum == CLIENT_SEND_NULL {
                    globals.Dbg.PrintfTrace ("%s [dl_msgs] --> unknown message type: \"%s\".\n", globals.LogTag, msgContainer.SendMsgType)
                    response.WriteHeader(404)
                } else {
                    msgSender := new (Msg)
                    msgSender.MsgType = msgEnum
                    msgSender.MsgBody = msgContainer.MsgBody
                    err = msgSender.Send (msgContainer.Uuid)
                    if err != nil {
                        response.WriteHeader(404)
                    } else {
                        response.WriteHeader(http.StatusOK)                        
                    }                    
                }           
            }
        }
    }
}

// Send the messages
func (m *Msg) Send(uuid string) error {
    switch m.MsgType {
        case CLIENT_SEND_PING:
            globals.Dbg.PrintfTrace ("%s [dl_msgs] --> send ping request.\n", globals.LogTag)                    
            encodeAndEnqueue (&PingReqDlMsg{}, uuid)
        default:
            globals.Dbg.PrintfTrace ("%s [dl_msgs] --> asked to send an unknown message type: %d.\n", globals.LogTag, m.MsgType)                    
            return errors.New("Unknown message type")
    }
    
    return nil
}

/*
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
