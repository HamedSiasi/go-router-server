/* Structs to send DL messages via the client interface from the UTM
 * NOTE: strictly speaking this functionality should be a part of the
 * controllers package, however none of the data we are manipulating
 * is in the model, it is all in the message codec/datable which have
 * been put in the server so the functionality has to be right here.
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
	"errors"
	"github.com/u-blox/utm-server/service/globals"
	"github.com/u-blox/utm-server/service/utilities"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

type ClientSendMsg struct {
	Uuid        string      `bson:"device_uuid" json:"device_uuid"`
	SendMsgType string      `bson:"type" json:"type"`
	MsgBody     interface{} `bson:"body" json:"body"`
}

type Msg struct {
	MsgType ClientSendEnum `bson:"type" json:"type"`
	MsgBody interface{}    `bson:"body" json:"body"`
}

type ClientSendEnum uint32

const (
	CLIENT_SEND_NULL ClientSendEnum = iota
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
	// This includes the downlink interval parameter
	// as well as the others which the server intercepts and
	// uses for itself
	CLIENT_SEND_TRAFFIC_TEST_MODE_PARAMETERS_SERVER_SET
	CLIENT_SEND_TRAFFIC_TEST_MODE_PARAMETERS_GET
	CLIENT_SEND_TRAFFIC_TEST_MODE_REPORT_GET
	CLIENT_SEND_ACTIVITY_REPORT_GET
)

type ClientModeEnum uint32

// IMPORTANT: these must be in the same order as the mode enum in decode.go
const (
	CLIENT_MODE_NULL ClientModeEnum = iota
	CLIENT_MODE_SELF_TEST
	CLIENT_MODE_COMMISSIONING
	CLIENT_MODE_STANDARD_TRX
	CLIENT_MODE_TRAFFIC_TEST
)

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

var clientSendEnumString map[string]ClientSendEnum = map[string]ClientSendEnum{
	"SEND_NULL":                                    CLIENT_SEND_NULL,
	"SEND_PING":                                    CLIENT_SEND_PING,
	"SEND_REBOOT":                                  CLIENT_SEND_REBOOT,
	"SEND_DATE_TIME_SET":                           CLIENT_SEND_DATE_TIME_SET,
	"SEND_DATE_TIME_GET":                           CLIENT_SEND_DATE_TIME_GET,
	"SEND_MODE_SET":                                CLIENT_SEND_MODE_SET,
	"SEND_MODE_GET":                                CLIENT_SEND_MODE_GET,
	"SEND_HEARTBEAT_SET":                           CLIENT_SEND_HEARTBEAT_SET,
	"SEND_REPORTING_INTERVAL_SET":                  CLIENT_SEND_REPORTING_INTERVAL_SET,
	"SEND_INTERVALS_GET":                           CLIENT_SEND_INTERVALS_GET,
	"SEND_MEASUREMENTS_GET":                        CLIENT_SEND_MEASUREMENTS_GET,
	"SEND_MEASUREMENT_CONTROL_SET":                 CLIENT_SEND_MEASUREMENT_CONTROL_SET,
	"SEND_MEASUREMENTS_CONTROL_GET":                CLIENT_SEND_MEASUREMENTS_CONTROL_GET,
	"SEND_MEASUREMENTS_CONTROL_DEFAULTS_SET":       CLIENT_SEND_MEASUREMENTS_CONTROL_DEFAULTS_SET,
	"SEND_TRAFFIC_REPORT_GET":                      CLIENT_SEND_TRAFFIC_REPORT_GET,
	"SEND_TRAFFIC_TEST_MODE_PARAMETERS_SERVER_SET": CLIENT_SEND_TRAFFIC_TEST_MODE_PARAMETERS_SERVER_SET,
	"SEND_TRAFFIC_TEST_MODE_PARAMETERS_GET":        CLIENT_SEND_TRAFFIC_TEST_MODE_PARAMETERS_GET,
	"SEND_TRAFFIC_TEST_MODE_REPORT_GET":            CLIENT_SEND_TRAFFIC_TEST_MODE_REPORT_GET,
	"SEND_ACTIVITY_REPORT_GET":                     CLIENT_SEND_ACTIVITY_REPORT_GET,
}

var clientModeEnumString map[string]ClientModeEnum = map[string]ClientModeEnum{
	"MODE_NULL":          CLIENT_MODE_NULL,
	"MODE_SELF_TEST":     CLIENT_MODE_SELF_TEST,
	"MODE_COMMISSIONING": CLIENT_MODE_COMMISSIONING,
	"MODE_STANDARD_TRX":  CLIENT_MODE_STANDARD_TRX,
	"MODE_TRAFFIC_TEST":  CLIENT_MODE_TRAFFIC_TEST,
}

// BootReq tags
const sdCardNotRequiredTag string = "sd_card_not_required"
const disableModemDebugTag string = "disable_modem_debug"
const disableButtonTag string = "disable_button"
const disableServerPingTag string = "disable_server_ping"

// DateTimeSetReq tags and format
const timeTag string = "time"
const setDateOnlyTag string = "set_date_only"

const dateTimeFormat = "2006-01-02 15:04:05"

// Mode tags
const modeTag string = "mode"

// Reporting Interval tags
const reportingIntervalTag string = "reporting_interval"

// Heartbeat tags
const heartbeatSecondsTag string = "heartbeat_seconds"
const heartbeatSnapToRtcTag string = "heartbeat_snap_to_rtc"

// Traffic Test Mode Parameters tags
const numUlDatagramsTag string = "num_ul_datagrams"
const lenUlDatagramTag string = "len_ul_datagram"
const numDlDatagramsTag string = "num_dl_datagrams"
const lenDlDatagramTag string = "len_dl_datagram"
const timeoutSecondsTag string = "timeout_seconds"
const noReportsDuringTestTag string = "no_reports_during_test"
const dlIntervalSecondsTag string = "dl_interval_seconds"

// TODO
// type clientMeasurementControlSetReqFields
// type clientMeasurementsControlDefaultsSetReqFields

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Get the value of a string token with the given tag
func GetValueString(body interface{}, tag string) (error, string) {

	globals.Dbg.PrintfTrace("%s [sender] --> tag is \"%s\".\n", globals.LogTag, tag)

	if body != nil && reflect.TypeOf(body).Kind() == reflect.Map {
		bodyMap := reflect.ValueOf(body).Interface().(map[string]interface{})
		globals.Dbg.PrintfTrace("%s [sender] --> map is %+v.\n", globals.LogTag, bodyMap)
		valueInterface := bodyMap[tag]
		if valueInterface != nil {
			if reflect.TypeOf(valueInterface).Kind() == reflect.String {
				valueString := reflect.ValueOf(valueInterface).Interface().(string)
				if valueString != "" {
					globals.Dbg.PrintfTrace("%s [sender] --> string value is \"%s\".\n", globals.LogTag, valueString)
					return nil, valueString
				} else {
					return errors.New("string from interface{} is null"), ""
				}
			} else {
				return errors.New("interface{} in map is not of type string"), ""
			}
		} else {
			return errors.New("map does not have tag in it"), ""
		}
	}

	return errors.New("body empty or is not a map"), ""
}

// Get the value of a UInt32 token with the given tag
func GetValueUint32(body interface{}, tag string) (error, uint32) {

	err, valueString := GetValueString(body, tag)
	if err == nil {
		value, err := strconv.ParseUint(valueString, 10, 32)
		if err == nil {
			globals.Dbg.PrintfTrace("%s [sender] --> value is %d.\n", globals.LogTag, value)
			return nil, uint32(value)
		} else {
			return errors.New("cannot convert value into uint32"), 0
		}
	} else {
		return err, 0
	}
}

// Get the value of a Boolean token with the given tag
func GetValueBool(body interface{}, tag string) (error, bool) {

	err, valueString := GetValueString(body, tag)
	if err == nil {
		value, err := strconv.ParseBool(valueString)
		if err == nil {
			globals.Dbg.PrintfTrace("%s [sender] --> value is %v.\n", globals.LogTag, value)
			return nil, value
		} else {
			return errors.New("cannot convert value into bool"), false
		}
	} else {
		return err, false
	}
}

// Get the value of a time token with the given tag
func GetValueTime(body interface{}, tag string) (error, time.Time) {

	var value time.Time
	var err error

	err, valueString := GetValueString(body, tag)
	if err == nil {
		value, err = time.Parse(dateTimeFormat, valueString)
		if err == nil {
			globals.Dbg.PrintfTrace("%s [sender] --> value is %T.\n", globals.LogTag, value)
			return nil, value
		} else {
			return errors.New("cannot convert value into time"), value
		}
	} else {
		return err, value
	}
}

// Determine what message to the UTM is requested by the REST interface
func (msg *ClientSendMsg) Send(response http.ResponseWriter, request *http.Request) {
	err := utilities.ValidatePostRequest(request)

	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	} else {
		globals.Dbg.PrintfTrace("%s [sender] --> sendMsg request received from client.\n", globals.LogTag)
		utilities.DumpRequest(request)
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			return
		} else {
			globals.Dbg.PrintfTrace("%s [sender] --> request body is \"%s\".\n", globals.LogTag, body)
			var msgContainer ClientSendMsg
			err = json.Unmarshal(body, &msgContainer)
			if err != nil {
				globals.Dbg.PrintfTrace("%s [sender] --> received:\n\n%+v\n\n...which is not JSON decodable:\n\n%+v\n", globals.LogTag, body, err.Error())
				response.WriteHeader(http.StatusNotFound)
				return
			} else {
				globals.Dbg.PrintfTrace("%s [sender] --> JSON says:\n\n%+v\n", globals.LogTag, msgContainer)
				msgEnum := clientSendEnumString[msgContainer.SendMsgType]

				if msgEnum == CLIENT_SEND_NULL {
					globals.Dbg.PrintfTrace("%s [sender] --> unknown message type: \"%s\".\n", globals.LogTag, msgContainer.SendMsgType)
					response.WriteHeader(http.StatusNotFound)
				} else {
					msgSender := new(Msg)
					msgSender.MsgType = msgEnum
					msgSender.MsgBody = msgContainer.MsgBody
					err = msgSender.Send(msgContainer.Uuid)
					if err != nil {
						response.WriteHeader(http.StatusNotFound)
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
	var err error = nil
	var byteCount int = 0
	var msgCount int = 0
	var responseId = RESPONSE_NONE

	switch m.MsgType {
	case CLIENT_SEND_PING:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send PingReq.\n", globals.LogTag)
			err, byteCount, responseId = encodeAndEnqueue(&PingReqDlMsg{}, uuid)
			if err == nil {
				msgCount++
			}
		}
	case CLIENT_SEND_REBOOT:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send RebootReq.\n", globals.LogTag)
			msg := &RebootReqDlMsg{}
			err, msg.SdCardNotRequired = GetValueBool(m.MsgBody, sdCardNotRequiredTag)
			if err == nil {
				err, msg.DisableModemDebug = GetValueBool(m.MsgBody, disableModemDebugTag)
				if err == nil {
					err, msg.DisableButton = GetValueBool(m.MsgBody, disableButtonTag)
					if err == nil {
						err, msg.DisableServerPing = GetValueBool(m.MsgBody, disableServerPingTag)
					}
				}
			}
			if err == nil {
				err, byteCount, responseId = encodeAndEnqueue(msg, uuid)
				if err == nil {
					msgCount++
				}
			}
		}
	case CLIENT_SEND_DATE_TIME_SET:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send DateTimeSetReq.\n", globals.LogTag)
			msg := &DateTimeSetReqDlMsg{}
			err, msg.UtmTime = GetValueTime(m.MsgBody, timeTag)
			if err == nil {
				err, msg.SetDateOnly = GetValueBool(m.MsgBody, setDateOnlyTag)
			}

			if err == nil {
				err, byteCount, responseId = encodeAndEnqueue(msg, uuid)
				if err == nil {
					msgCount++
				}
			}
		}
	case CLIENT_SEND_DATE_TIME_GET:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send DateTimeGetReq.\n", globals.LogTag)
			err, byteCount, responseId = encodeAndEnqueue(&DateTimeGetReqDlMsg{}, uuid)
			if err == nil {
				msgCount++
			}
		}
	case CLIENT_SEND_MODE_SET:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send ModeSetReq.\n", globals.LogTag)
			var modeString string
			err, modeString = GetValueString(m.MsgBody, modeTag)
			if err == nil {
				msg := &ModeSetReqDlMsg{}
				msg.Mode = ModeEnum(clientModeEnumString[modeString]) + MODE_NULL
				if msg.Mode != MODE_NULL {
					err, byteCount, responseId = encodeAndEnqueue(msg, uuid)
					if err == nil {
						msgCount++
						// The traffic test channel needs to know about this to determine if it
						// should end a test
						globals.Dbg.PrintfTrace("%s [sender] --> also sending message to Traffic Test channel.\n", globals.LogTag)
						msgContainer := &MessageContainer{
							DeviceUuid: uuid,
							Timestamp:  time.Now().UTC(),
							Message:    msg,
						}
						trafficTestChannel <- msgContainer
					}
				} else {
					globals.Dbg.PrintfTrace("%s [sender] --> unknown mode \"%s\".\n", globals.LogTag, modeString)
				}
			}
		}
	case CLIENT_SEND_MODE_GET:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send ModeGetReq.\n", globals.LogTag)
			err, byteCount, responseId = encodeAndEnqueue(&ModeGetReqDlMsg{}, uuid)
			if err == nil {
				msgCount++
			}
		}
	case CLIENT_SEND_HEARTBEAT_SET:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send HeartbeatSetReq.\n", globals.LogTag)
			// Retrieve the values from the map that JSON has created for the message body
			msg := &HeartbeatSetReqDlMsg{}
			err, msg.HeartbeatSeconds = GetValueUint32(m.MsgBody, heartbeatSecondsTag)
			if err == nil {
				err, msg.HeartbeatSnapToRtc = GetValueBool(m.MsgBody, heartbeatSnapToRtcTag)
			}
			if err == nil {
				err, byteCount, responseId = encodeAndEnqueue(msg, uuid)
				if err == nil {
					msgCount++
				}
			}
		}
	case CLIENT_SEND_REPORTING_INTERVAL_SET:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send ReportingIntervalSetReq.\n", globals.LogTag)
			msg := &ReportingIntervalSetReqDlMsg{}
			err, msg.ReportingInterval = GetValueUint32(m.MsgBody, reportingIntervalTag)
			if err == nil {
				err, byteCount, responseId = encodeAndEnqueue(msg, uuid)
				if err == nil {
					msgCount++
				}
			}
		}
	case CLIENT_SEND_INTERVALS_GET:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send IntervalsGetReq.\n", globals.LogTag)
			err, byteCount, responseId = encodeAndEnqueue(&IntervalsGetReqDlMsg{}, uuid)
			if err == nil {
				msgCount++
			}
		}
	case CLIENT_SEND_MEASUREMENTS_GET:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send MeasurementsGetReq.\n", globals.LogTag)
			err, byteCount, responseId = encodeAndEnqueue(&MeasurementsGetReqDlMsg{}, uuid)
			if err == nil {
				msgCount++
			}
		}
	// TODO
	//case CLIENT_SEND_MEASUREMENT_CONTROL_SET:
	//case CLIENT_SEND_MEASUREMENTS_CONTROL_GET:
	//case CLIENT_SEND_MEASUREMENTS_CONTROL_DEFAULTS_SET:

	case CLIENT_SEND_TRAFFIC_REPORT_GET:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send TrafficReportGetReq.\n", globals.LogTag)
			err, byteCount, responseId = encodeAndEnqueue(&TrafficReportGetReqDlMsg{}, uuid)
			if err == nil {
				msgCount++
			}
		}
	case CLIENT_SEND_TRAFFIC_TEST_MODE_PARAMETERS_SERVER_SET:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send TrafficTestModeParametersSetReq.\n", globals.LogTag)
			// First, assemble the server internal message
			serverMsg := &TrafficTestModeParametersServerSet{}
			serverMsg.DeviceParameters = &TrafficTestModeParametersSetReqDlMsg{}

			err, serverMsg.DeviceParameters.NumUlDatagrams = GetValueUint32(m.MsgBody, numUlDatagramsTag)
			if err == nil {
				err, serverMsg.DeviceParameters.LenUlDatagram = GetValueUint32(m.MsgBody, lenUlDatagramTag)
				if err == nil {
					err, serverMsg.DeviceParameters.NumDlDatagrams = GetValueUint32(m.MsgBody, numDlDatagramsTag)
					if err == nil {
						err, serverMsg.DeviceParameters.LenDlDatagram = GetValueUint32(m.MsgBody, lenDlDatagramTag)
						if err == nil {
							err, serverMsg.DeviceParameters.TimeoutSeconds = GetValueUint32(m.MsgBody, timeoutSecondsTag)
							if err == nil {
								err, serverMsg.DeviceParameters.NoReportsDuringTest = GetValueBool(m.MsgBody, noReportsDuringTestTag)
								if err == nil {
									err, serverMsg.DlIntervalSeconds = GetValueUint32(m.MsgBody, dlIntervalSecondsTag)
								}
							}
						}
					}
				}
			}

			if err == nil {
				// Now copy the device portion out into its message
				deviceMsg := serverMsg.DeviceParameters.DeepCopy()
				err, byteCount, responseId = encodeAndEnqueue(deviceMsg, uuid)
				if err == nil {
					msgCount++
					globals.Dbg.PrintfTrace("%s [sender] --> also sending message to Traffic Test channel.\n", globals.LogTag)
					serverMsgContainer := &MessageContainer{
						DeviceUuid: uuid,
						Timestamp:  time.Now().UTC(),
						Message:    serverMsg,
					}
					trafficTestChannel <- serverMsgContainer
				}
			}
		}
	case CLIENT_SEND_TRAFFIC_TEST_MODE_PARAMETERS_GET:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send TrafficTestModeParametersGetReq.\n", globals.LogTag)
			err, byteCount, responseId = encodeAndEnqueue(&TrafficTestModeParametersGetReqDlMsg{}, uuid)
			if err == nil {
				msgCount++
			}
		}
	case CLIENT_SEND_TRAFFIC_TEST_MODE_REPORT_GET:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send TrafficTestModeReportGetReq.\n", globals.LogTag)
			err, byteCount, responseId = encodeAndEnqueue(&TrafficTestModeReportGetReqDlMsg{}, uuid)
			if err == nil {
				msgCount++
			}
		}
	case CLIENT_SEND_ACTIVITY_REPORT_GET:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> send ActivityReportGetReq.\n", globals.LogTag)
			err, byteCount, responseId = encodeAndEnqueue(&ActivityReportGetReqDlMsg{}, uuid)
			if err == nil {
				msgCount++
			}
		}
	default:
		{
			globals.Dbg.PrintfTrace("%s [sender] --> asked to send an unknown message type: %d.\n", globals.LogTag, m.MsgType)
			return errors.New("Unknown message type")
		} // case
	} // switch

	// Send an update to the processing loop so that we don't miss
	// out on the totals
	if msgCount > 0 {
		encodeStateAdd := &DeviceEncodeStateAdd{
			DeviceUuid: uuid,
			ResponseId: responseId,
			State: TotalsState{
				Timestamp: time.Now().UTC(),
				Msgs:      msgCount,
				Bytes:     byteCount,
			},
		}
		processMsgsChannel <- encodeStateAdd
	}

	return nil
}

/* End Of File */
