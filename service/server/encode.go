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
    "unsafe"
    "github.com/robmeades/utm/service/globals"
    "github.com/robmeades/utm/service/utilities"
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

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Encode and enqueue a message, return an error if there is one,
// the number of bytes encoded and the response expected (which may be
// none)
func encodeAndEnqueue(msg interface{}, uuid string) (error, int, ResponseTypeEnum) {
    
    if msg != nil {        
        // Create a buffer that is big enough to store all
        // the encoded data and take a pointer to it's first element
        var outputBuffer [MaxDatagramSizeRaw]byte
        outputPointer := (*C.char)(unsafe.Pointer(&outputBuffer[0]))
        var xmlEncodeBuffer = make([]byte, 8192, 8192)
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
                globals.Dbg.PrintfTrace("%s [encode] --> encoded PingReqDlMsg.\n", globals.LogTag)
                responseId = RESPONSE_PING_CNF
                
            case *PingCnfDlMsg:
                byteCount = C.encodePingCnfMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded PingCnfDlMsg.\n", globals.LogTag)
                
            case *RebootReqDlMsg:
                data := C.RebootReqDlMsg_t {
                    sdCardNotRequired: (C.bool) (value.SdCardNotRequired),
                    disableModemDebug: (C.bool) (value.DisableModemDebug),
                    disableButton:     (C.bool) (value.DisableButton),
                    disableServerPing: (C.bool) (value.DisableServerPing),                    
                }
                dataPointer := (*C.RebootReqDlMsg_t)(unsafe.Pointer(&data))
                byteCount = C.encodeRebootReqDlMsg(outputPointer, dataPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded RebootReqDlMsg.\n", globals.LogTag)
                
            case *DateTimeSetReqDlMsg:
                data := C.DateTimeSetReqDlMsg_t {
                    time:        (C.uint32_t) (value.UtmTime.Unix()),
                    setDateOnly: (C.bool) (value.SetDateOnly),
                }
                dataPointer := (*C.DateTimeSetReqDlMsg_t)(unsafe.Pointer(&data))
                byteCount = C.encodeDateTimeSetReqDlMsg(outputPointer, dataPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded DateTimeSetReqDlMsg.\n", globals.LogTag)
                responseId = RESPONSE_DATE_TIME_SET_CNF
            case *DateTimeGetReqDlMsg:
                byteCount = C.encodeDateTimeGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded DateTimeGetReqDlMsg.\n", globals.LogTag)
                responseId = RESPONSE_DATE_TIME_GET_CNF
                
            case *ModeSetReqDlMsg:
                data := C.ModeSetReqDlMsg_t {
                    mode:  (C.Mode_t) (value.Mode),
                }
                dataPointer := (*C.ModeSetReqDlMsg_t)(unsafe.Pointer(&data))
                byteCount = C.encodeModeSetReqDlMsg(outputPointer, dataPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded ModeSetReqDlMsg.\n", globals.LogTag)
                responseId = RESPONSE_MODE_SET_CNF
                
            case *ModeGetReqDlMsg:
                byteCount = C.encodeModeGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded ModeGetReqDlMsg.\n", globals.LogTag)
                responseId = RESPONSE_MODE_GET_CNF
            
            case *IntervalsGetReqDlMsg:
                byteCount = C.encodeIntervalsGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded IntervalsGetReqDlMsg.\n", globals.LogTag)
                responseId = RESPONSE_INTERVALS_GET_CNF
                
            case *ReportingIntervalSetReqDlMsg:
                data := C.ReportingIntervalSetReqDlMsg_t {
                    reportingInterval:  (C.uint32_t) (value.ReportingInterval),
                }
                dataPointer := (*C.ReportingIntervalSetReqDlMsg_t)(unsafe.Pointer(&data))
                byteCount = C.encodeReportingIntervalSetReqDlMsg(outputPointer, dataPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded ReportingIntervalSetReqDlMsg.\n", globals.LogTag)
                responseId = RESPONSE_REPORTING_INTERVAL_SET_CNF
                
            case *HeartbeatSetReqDlMsg:
                data := C.HeartbeatSetReqDlMsg_t {
                    heartbeatSeconds:    (C.uint32_t) (value.HeartbeatSeconds),
                    heartbeatSnapToRtc:  (C.bool) (value.HeartbeatSnapToRtc),
                }
                dataPointer := (*C.HeartbeatSetReqDlMsg_t)(unsafe.Pointer(&data))
                byteCount = C.encodeHeartbeatSetReqDlMsg(outputPointer, dataPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded HeartbeatSetReqDlMsg.\n", globals.LogTag)
                responseId = RESPONSE_HEARTBEAT_SET_CNF
                
            case *MeasurementsGetReqDlMsg:
                byteCount = C.encodeMeasurementsGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded MeasurementsReportReqDlMsg.\n", globals.LogTag)        
                responseId = RESPONSE_MEASUREMENTS_GET_CNF
                        
            // TODO
            // case *MeasurementControlSetReqDlMsg
            // case *MeasurementsControlDefaultsSetReqDlMsg
            
            case *TrafficReportGetReqDlMsg:
                byteCount = C.encodeTrafficReportGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded TrafficReportGetReqDlMsg.\n", globals.LogTag)
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
                globals.Dbg.PrintfTrace("%s [encode] --> encoded TrafficTestModeParametersSetReqDlMsg.\n", globals.LogTag)
                responseId = RESPONSE_TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF
                
            case *TrafficTestModeParametersGetReqDlMsg:
                byteCount = C.encodeTrafficTestModeParametersGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded TrafficTestModeParametersGetReqDlMsg.\n", globals.LogTag)
                responseId = RESPONSE_TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF
                                
            case *TrafficTestModeReportGetReqDlMsg:
                byteCount = C.encodeTrafficTestModeReportGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded TrafficTestModeReportGetReqDlMsg.\n", globals.LogTag)
                responseId = RESPONSE_TRAFFIC_TEST_MODE_REPORT_GET_CNF
                                
            case *ActivityReportGetReqDlMsg:
                byteCount = C.encodeActivityReportGetReqDlMsg(outputPointer, ppXmlBuffer, pXmlBufferLen)
                globals.Dbg.PrintfTrace("%s [encode] --> encoded ActivityReportGetReqDlMsg.\n", globals.LogTag)                
                responseId = RESPONSE_ACTIVITY_REPORT_GET_CNF

            default:
                globals.Dbg.PrintfTrace("%s [encode] --> asked to send unknown message:\n\n%s\n", globals.LogTag, spew.Sdump(msg))
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
            }
    
            globals.Dbg.PrintfTrace("%s [encode] --> %d byte message for AMQP downlink:\n\n%+v\n", globals.LogTag, byteCount, msg)
            if DlMsgs != nil {
                DlMsgs <- msg            
                globals.Dbg.PrintfTrace("%s [encode] --> encoded %d bytes into AMQP message:\n\n%+v\n", globals.LogTag, byteCount, msg)
    
                // Store XmlData in MongoDB
                utilities.XmlDataStore(xmlEncodeBuffer, uuid)
                globals.Dbg.PrintfInfo("%s [decode] --> the XML data is:\n\n%s\n\n", globals.LogTag, spew.Sdump(xmlEncodeBuffer))            
                globals.Dbg.PrintfInfo("%s [encode] --> XML buffer pointer 0x%08x, used %d, left %d:.\n", globals.LogTag, *ppXmlBuffer, C.uint32_t(len(xmlEncodeBuffer)) - xmlBufferLen, xmlBufferLen)
            } else {
                return errors.New("No downlink message channel available to enqueue the encoded message.\n"), 0, RESPONSE_NONE    
            }                
        }                       
        
        return nil, int(byteCount), responseId
    }
    
    return nil, 0, RESPONSE_NONE
}

/* End Of File */
