/* Message encode functions for the UTM server.
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

//--------------------------------------------------------------------
// CGO stuff
//--------------------------------------------------------------------

/*
#cgo CFLAGS: -I .
#include <stdbool.h>
#include <stdint.h>
#include <utm_api.h>

*/
import "C"  // There must be no line breaks between this and the commented-out section that comes before (bloody whitespace sensitive syntax)
import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"time"
	"unsafe"
)

//--------------------------------------------------------------------
// Enums
//--------------------------------------------------------------------

type ResponseTypeEnum uint32

const (
       RESPONSE_NONE                     ResponseTypeEnum = iota
       RESPONSE_PING_CNF
       RESPONSE_DATE_TIME_SET_CNF
       RESPONSE_DATE_TIME_GET_CNF
       RESPONSE_MODE_SET_CNF
       RESPONSE_MODE_GET_CNF
       RESPONSE_INTERVALS_GET_CNF
       RESPONSE_REPORTING_INTERVAL_SET_CNF
       RESPONSE_HEARTBEAT_SET_CNF
       RESPONSE_MEASUREMENTS_GET_CNF
       RESPONSE_TRAFFIC_REPORT_GET_CNF
       RESPONSE_TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF
       RESPONSE_TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF
       RESPONSE_TRAFFIC_TEST_MODE_REPORT_GET_CNF
       RESPONSE_ACTIVITY_REPORT_GET_CNF
)

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

var xmlEncodeBuffer [8192]C.char

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Encode and enqueue a message, return an error if there is one
func encodeAndEnqueue(msg interface{}, uuid string) error {
    
    if msg != nil {        
		// Create a buffer that is big enough to store all
		// the encoded data and take a pointer to it's first element
		var outputBuffer [MaxDatagramSizeRaw]byte
		outputPointer := (*C.char)(unsafe.Pointer(&outputBuffer[0]))
		var byteCount C.uint32_t = 0     
        responseId := RESPONSE_NONE
		
    	// A place to put the XML output from the decoder
    	pXmlBuffer := (*C.char) (unsafe.Pointer(&(xmlEncodeBuffer[0])))
    	ppXmlBuffer := (**C.char) (unsafe.Pointer(&pXmlBuffer))
        xmlBufferLen := (C.uint32_t) (len(xmlEncodeBuffer))
        pXmlBufferLen := (*C.uint32_t) (unsafe.Pointer(&xmlBufferLen))
        
		// Encode each message type
		switch value := msg.(type) {
		    case *TransparentDlDatagram:
		        // TODO
            case *PingReqDlMsg:
        		byteCount = C.encodePingReqMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded PingReqDlMsg.\n", logTag)
                responseId = RESPONSE_PING_CNF
                
            case *PingCnfDlMsg:
        		byteCount = C.encodePingCnfMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded PingCnfDlMsg.\n", logTag)
                
            case *RebootReqDlMsg:
                data := C.RebootReqDlMsg_t {
                    sdCardNotRequired: (C.bool) (value.SdCardNotRequired),
                    disableModemDebug: (C.bool) (value.DisableModemDebug),
                    disableButton:     (C.bool) (value.DisableButton),
                    disableServerPing: (C.bool) (value.DisableServerPing),                    
                }
        		dataPointer := (*C.RebootReqDlMsg_t)(unsafe.Pointer(&data))
        		byteCount = C.encodeRebootReqDlMsg(outputPointer, dataPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded RebootReqDlMsg.\n", logTag)
                
            case *DateTimeSetReqDlMsg:
                data := C.DateTimeSetReqDlMsg_t {
                    time:        (C.uint32_t) (value.UtmTime.Unix()),
                    setDateOnly: (C.bool) (value.SetDateOnly),
                }
        		dataPointer := (*C.DateTimeSetReqDlMsg_t)(unsafe.Pointer(&data))
        		byteCount = C.encodeDateTimeSetReqDlMsg(outputPointer, dataPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded DateTimeSetReqDlMsg.\n", logTag)
                responseId = RESPONSE_DATE_TIME_SET_CNF
            case *DateTimeGetReqDlMsg:
        		byteCount = C.encodeDateTimeGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded DateTimeGetReqDlMsg.\n", logTag)
                responseId = RESPONSE_DATE_TIME_GET_CNF
                
            case *ModeSetReqDlMsg:
                data := C.ModeSetReqDlMsg_t {
                    mode:  (C.Mode_t) (value.Mode),
                }
        		dataPointer := (*C.ModeSetReqDlMsg_t)(unsafe.Pointer(&data))
        		byteCount = C.encodeModeSetReqDlMsg(outputPointer, dataPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded ModeSetReqDlMsg.\n", logTag)
                responseId = RESPONSE_MODE_SET_CNF
                
            case *ModeGetReqDlMsg:
        		byteCount = C.encodeModeGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded ModeGetReqDlMsg.\n", logTag)
                responseId = RESPONSE_MODE_GET_CNF
            
            case *IntervalsGetReqDlMsg:
        		byteCount = C.encodeIntervalsGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded IntervalsGetReqDlMsg.\n", logTag)
                responseId = RESPONSE_INTERVALS_GET_CNF
                
            case *ReportingIntervalSetReqDlMsg:
                data := C.ReportingIntervalSetReqDlMsg_t {
                    reportingInterval:  (C.uint32_t) (value.ReportingInterval),
                }
        		dataPointer := (*C.ReportingIntervalSetReqDlMsg_t)(unsafe.Pointer(&data))
        		byteCount = C.encodeReportingIntervalSetReqDlMsg(outputPointer, dataPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded ReportingIntervalSetReqDlMsg.\n", logTag)
                responseId = RESPONSE_REPORTING_INTERVAL_SET_CNF
                
            case *HeartbeatSetReqDlMsg:
                data := C.HeartbeatSetReqDlMsg_t {
                    heartbeatSeconds:    (C.uint32_t) (value.HeartbeatSeconds),
                    heartbeatSnapToRtc:  (C.bool) (value.HeartbeatSnapToRtc),
                }
        		dataPointer := (*C.HeartbeatSetReqDlMsg_t)(unsafe.Pointer(&data))
        		byteCount = C.encodeHeartbeatSetReqDlMsg(outputPointer, dataPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded HeartbeatSetReqDlMsg.\n", logTag)
                responseId = RESPONSE_HEARTBEAT_SET_CNF
                
            case *MeasurementsGetReqDlMsg:
        		byteCount = C.encodeMeasurementsGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded MeasurementsReportReqDlMsg.\n", logTag)        
                responseId = RESPONSE_MEASUREMENTS_GET_CNF
                        
            // TODO
            // case *MeasurementControlSetReqDlMsg
            // case *MeasurementsControlDefaultsSetReqDlMsg
            
            case *TrafficReportGetReqDlMsg:
        		byteCount = C.encodeTrafficReportGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded TrafficReportGetReqDlMsg.\n", logTag)
                responseId = RESPONSE_TRAFFIC_REPORT_GET_CNF
                                
            case *TrafficTestModeParametersSetReqDlMsg:
                data := C.TrafficTestModeParametersSetReqDlMsg_t {
                    numUlDatagrams:      (C.uint32_t) (value.NumUlDatagrams),
                    lenUlDatagram:       (C.uint32_t) (value.LenUlDatagram),
                    numDlDatagrams:      (C.uint32_t) (value.NumDlDatagrams),
                    lenDlDatagram:       (C.uint32_t) (value.LenDlDatagram),
                    timeoutSeconds:      (C.uint32_t) (value.TimeoutSeconds),
                    noReportsDuringTest: (C.bool) (value.NoReportsDuringTest),
                }
        		dataPointer := (*C.TrafficTestModeParametersSetReqDlMsg_t)(unsafe.Pointer(&data))
        		byteCount = C.encodeTrafficTestModeParametersSetReqDlMsg(outputPointer, dataPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded TrafficTestModeParametersSetReqDlMsg.\n", logTag)
                responseId = RESPONSE_TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF
                
            case *TrafficTestModeParametersGetReqDlMsg:
        		byteCount = C.encodeTrafficTestModeParametersGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded TrafficTestModeParametersGetReqDlMsg.\n", logTag)
                responseId = RESPONSE_TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF
                                
            case *TrafficTestModeReportGetReqDlMsg:
        		byteCount = C.encodeTrafficTestModeReportGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded TrafficTestModeReportGetReqDlMsg.\n", logTag)
                responseId = RESPONSE_TRAFFIC_TEST_MODE_REPORT_GET_CNF
                                
            case *ActivityReportGetReqDlMsg:
        		byteCount = C.encodeActivityReportGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                Dbg.PrintfTrace("%s --> encoded ActivityReportGetReqDlMsg.\n", logTag)                
                responseId = RESPONSE_ACTIVITY_REPORT_GET_CNF

		    default:
                Dbg.PrintfTrace("%s --> asked to send unknown message:\n\n%s\n", logTag, spew.Sdump(msg))
		}
		
	    if byteCount > 0 {
    		// Send the output buffer to the TSW server
    		payload := outputBuffer[:byteCount]
    		msg := AmqpMessage {
    			DeviceUuid:   uuid,
    			EndpointUuid: 4,
    		}
    		
    		for _, v := range payload {
    			msg.Payload = append(msg.Payload, int(v))
    			// TODO
    			row.TotalBytes = row.TotalBytes + uint64(len(payload))
    		}
    
    		downlinkMessages <- msg
    		
    		Dbg.PrintfTrace("%s --> encoded %d bytes into AMQP message:\n\n%+v\n", logTag, byteCount, msg)
			Dbg.PrintfInfo("%s --> XML buffer pointer 0x%08x, used %d, left %d:.\n", logTag, *ppXmlBuffer, C.uint32_t(len(xmlEncodeBuffer)) - xmlBufferLen, xmlBufferLen)
    		
        	// If a response is expected, add it to the list for this device
        	if (responseId != RESPONSE_NONE) {
            	list := deviceExpectedMsgList[uuid]
        		if list == nil {
        		    list = make (ExpectedMsgList, 0)
        			deviceExpectedMsgList[uuid] = list
        		}
        		
        		expectedMsg := ExpectedMsg {
                    TimeStarted: time.Now(),
                    ResponseId:  responseId,		    
        		}
        		
    			list = append(list, expectedMsg)
            }    
            		
    		// TODO
    		now := time.Now()    
    		// Record the downlink data volume
    		dataTableCmds <- &DataVolume {
    			DownlinkTimestamp: &now,
    			DownlinkBytes:     uint64(len(payload)),
    		}
    		
    	    return nil
        }    				   
	    
   	    return nil
    }
    
	return errors.New("No downlink message channel available to enqueue the encoded message.\n")
}

/* End Of File */
