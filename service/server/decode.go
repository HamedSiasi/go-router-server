/* Message decode functions for the UTM server.
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

uint32_t pointerSub(const char *a, const char *b)
{
	return (uint32_t) (a - b);
}

TransparentDatagram_t getTransparentDatagram(UlMsgUnion_t in)
{
	return in.transparentDatagram;
}

InitIndUlMsg_t getInitIndUlMsg(UlMsgUnion_t in)
{
	return in.initIndUlMsg;
}

DateTimeSetCnfUlMsg_t getDateTimeSetCnfUlMsg(UlMsgUnion_t in)
{
	return in.dateTimeSetCnfUlMsg;
}

DateTimeGetCnfUlMsg_t getDateTimeGetCnfUlMsg(UlMsgUnion_t in)
{
	return in.dateTimeGetCnfUlMsg;
}

DateTimeIndUlMsg_t getDateTimeIndUlMsg(UlMsgUnion_t in)
{
	return in.dateTimeIndUlMsg;
}

ModeSetCnfUlMsg_t getModeSetCnfUlMsg(UlMsgUnion_t in)
{
	return in.modeSetCnfUlMsg;
}

ModeGetCnfUlMsg_t getModeGetCnfUlMsg(UlMsgUnion_t in)
{
	return in.modeGetCnfUlMsg;
}

IntervalsGetCnfUlMsg_t getIntervalsGetCnfUlMsg(UlMsgUnion_t in)
{
	return in.intervalsGetCnfUlMsg;
}

ReportingIntervalSetCnfUlMsg_t getReportingIntervalSetCnfUlMsg(UlMsgUnion_t in)
{
	return in.reportingIntervalSetCnfUlMsg;
}

HeartbeatSetCnfUlMsg_t getHeartbeatSetCnfUlMsg(UlMsgUnion_t in)
{
	return in.heartbeatSetCnfUlMsg;
}

PollIndUlMsg_t getPollIndUlMsg(UlMsgUnion_t in)
{
	return in.pollIndUlMsg;
}

MeasurementsGetCnfUlMsg_t getMeasurementsGetCnfUlMsg(UlMsgUnion_t in)
{
	return in.measurementsGetCnfUlMsg;
}

MeasurementsIndUlMsg_t getMeasurementsIndUlMsg(UlMsgUnion_t in)
{
	return in.measurementsIndUlMsg;
}

MeasurementControlSetCnfUlMsg_t getMeasurementControlSetCnfUlMsg(UlMsgUnion_t in)
{
	return in.measurementControlSetCnfUlMsg;
}

MeasurementsControlGetCnfUlMsg_t getMeasurementsControlGetCnfUlMsg(UlMsgUnion_t in)
{
	return in.measurementsControlGetCnfUlMsg;
}

MeasurementsControlIndUlMsg_t getMeasurementsControlIndUlMsg(UlMsgUnion_t in)
{
	return in.measurementsControlIndUlMsg;
}

TrafficReportIndUlMsg_t getTrafficReportIndUlMsg(UlMsgUnion_t in)
{
	return in.trafficReportIndUlMsg;
}

TrafficReportGetCnfUlMsg_t getTrafficReportGetCnfUlMsg(UlMsgUnion_t in)
{
	return in.trafficReportGetCnfUlMsg;
}

TrafficTestModeParametersSetCnfUlMsg_t getTrafficTestModeParametersSetCnfUlMsg(UlMsgUnion_t in)
{
	return in.trafficTestModeParametersSetCnfUlMsg;
}

TrafficTestModeParametersGetCnfUlMsg_t getTrafficTestModeParametersGetCnfUlMsg(UlMsgUnion_t in)
{
	return in.trafficTestModeParametersGetCnfUlMsg;
}

TrafficTestModeRuleBreakerDatagram_t getTrafficTestModeRuleBreakerDatagram(UlMsgUnion_t in)
{
	return in.trafficTestModeRuleBreakerDatagram;
}

TrafficTestModeReportIndUlMsg_t getTrafficTestModeReportIndUlMsg(UlMsgUnion_t in)
{
	return in.trafficTestModeReportIndUlMsg;
}

TrafficTestModeReportGetCnfUlMsg_t getTrafficTestModeReportGetCnfUlMsg(UlMsgUnion_t in)
{
	return in.trafficTestModeReportGetCnfUlMsg;
}

ActivityReportIndUlMsg_t getActivityReportIndUlMsg(UlMsgUnion_t in)
{
	return in.activityReportIndUlMsg;
}

ActivityReportGetCnfUlMsg_t getActivityReportGetCnfUlMsg(UlMsgUnion_t in)
{
	return in.activityReportGetCnfUlMsg;
}

DebugIndUlMsg_t getDebugIndUlMsg(UlMsgUnion_t in)
{
	return in.debugIndUlMsg;
}

*/
import "C"  // There must be no line breaks between this and the commented-out section that comes before (bloody whitespace sensitive syntax)
import (
	"encoding/hex"
	"github.com/davecgh/go-spew/spew"
	"github.com/robmeades/utm/service/globals"
	"github.com/robmeades/utm/service/utilities"
	"time"
	"unsafe"
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

type ModeEnum uint32

const (
    MODE_NULL          ModeEnum = C.MODE_NULL
    MODE_SELF_TEST     ModeEnum = C.MODE_SELF_TEST
    MODE_COMMISSIONING ModeEnum = C.MODE_COMMISSIONING
    MODE_STANDARD_TRX  ModeEnum = C.MODE_STANDARD_TRX
    MODE_TRAFFIC_TEST  ModeEnum = C.MODE_TRAFFIC_TEST
    MAX_NUM_MODES      ModeEnum = C.MAX_NUM_MODES
)

type WakeUpEnum uint32

const (
    WAKE_UP_CODE_OK                    WakeUpEnum = C.WAKE_UP_CODE_OK
    WAKE_UP_CODE_WATCHDOG              WakeUpEnum = C.WAKE_UP_CODE_WATCHDOG
    WAKE_UP_CODE_NETWORK_PROBLEM       WakeUpEnum = C.WAKE_UP_CODE_NETWORK_PROBLEM
    WAKE_UP_CODE_SD_CARD_PROBLEM       WakeUpEnum = C.WAKE_UP_CODE_SD_CARD_PROBLEM
    WAKE_UP_CODE_SUPPLY_PROBLEM        WakeUpEnum = C.WAKE_UP_CODE_SUPPLY_PROBLEM
    WAKE_UP_CODE_PROTOCOL_PROBLEM      WakeUpEnum = C.WAKE_UP_CODE_PROTOCOL_PROBLEM
    WAKE_UP_CODE_MODULE_NOT_RESPONDING WakeUpEnum = C.WAKE_UP_CODE_MODULE_NOT_RESPONDING
    WAKE_UP_CODE_HW_PROBLEM            WakeUpEnum = C.WAKE_UP_CODE_HW_PROBLEM
    WAKE_UP_CODE_MEMORY_ALLOC_PROBLEM  WakeUpEnum = C.WAKE_UP_CODE_MEMORY_ALLOC_PROBLEM
    WAKE_UP_CODE_GENERIC_FAILURE       WakeUpEnum = C.WAKE_UP_CODE_GENERIC_FAILURE
    WAKE_UP_CODE_REBOOT                WakeUpEnum = C.WAKE_UP_CODE_REBOOT
    MAX_NUM_WAKE_UP_CODES              WakeUpEnum = C.MAX_NUM_WAKE_UP_CODES
)

type TimeSetByEnum uint32

const (
    TIME_SET_BY_NULL    TimeSetByEnum = C.TIME_SET_BY_NULL
    TIME_SET_BY_GNSS    TimeSetByEnum = C.TIME_SET_BY_GNSS
    TIME_SET_BY_PC      TimeSetByEnum = C.TIME_SET_BY_PC
    TIME_SET_BY_WEB_API TimeSetByEnum = C.TIME_SET_BY_WEB_API
    MAX_NUM_TIME_SET_BY TimeSetByEnum = C.MAX_NUM_TIME_SET_BY
)

type EnergyLeftEnum uint32

const (
    ENERGY_LEFT_LESS_THAN_5_PERCENT  EnergyLeftEnum = C.ENERGY_LEFT_LESS_THAN_5_PERCENT
    ENERGY_LEFT_LESS_THAN_10_PERCENT EnergyLeftEnum = C.ENERGY_LEFT_LESS_THAN_10_PERCENT
    ENERGY_LEFT_MORE_THAN_10_PERCENT EnergyLeftEnum = C.ENERGY_LEFT_MORE_THAN_10_PERCENT
    ENERGY_LEFT_MORE_THAN_30_PERCENT EnergyLeftEnum = C.ENERGY_LEFT_MORE_THAN_30_PERCENT
    ENERGY_LEFT_MORE_THAN_50_PERCENT EnergyLeftEnum = C.ENERGY_LEFT_MORE_THAN_50_PERCENT
    ENERGY_LEFT_MORE_THAN_70_PERCENT EnergyLeftEnum = C.ENERGY_LEFT_MORE_THAN_70_PERCENT
    ENERGY_LEFT_MORE_THAN_90_PERCENT EnergyLeftEnum = C.ENERGY_LEFT_MORE_THAN_90_PERCENT
    MAX_NUM_ENERGY_LEFT              EnergyLeftEnum = C.MAX_NUM_ENERGY_LEFT
)

type DiskSpaceLeftEnum uint32

const (
    DISK_SPACE_LEFT_LESS_THAN_1_GB DiskSpaceLeftEnum = C.DISK_SPACE_LEFT_LESS_THAN_1_GB
    DISK_SPACE_LEFT_MORE_THAN_1_GB DiskSpaceLeftEnum = C.DISK_SPACE_LEFT_MORE_THAN_1_GB
    DISK_SPACE_LEFT_MORE_THAN_2_GB DiskSpaceLeftEnum = C.DISK_SPACE_LEFT_MORE_THAN_2_GB
    DISK_SPACE_LEFT_MORE_THAN_4_GB DiskSpaceLeftEnum = C.DISK_SPACE_LEFT_MORE_THAN_4_GB
    MAX_NUM_DISK_SPACE_LEFT        DiskSpaceLeftEnum = C.MAX_NUM_DISK_SPACE_LEFT
)

type MeasurementTypeEnum uint32

const (
    MEASUREMENT_TYPE_UNKNOWN MeasurementTypeEnum = C.MEASUREMENT_TYPE_UNKNOWN
    MEASUREMENT_GNSS         MeasurementTypeEnum = C.MEASUREMENT_GNSS
    MEASUREMENT_CELL_ID      MeasurementTypeEnum = C.MEASUREMENT_CELL_ID
    MEASUREMENT_RSRP         MeasurementTypeEnum = C.MEASUREMENT_RSRP
    MEASUREMENT_RSSI         MeasurementTypeEnum = C.MEASUREMENT_RSSI
    MEASUREMENT_TEMPERATURE  MeasurementTypeEnum = C.MEASUREMENT_TEMPERATURE
    MEASUREMENT_POWER_STATE  MeasurementTypeEnum = C.MEASUREMENT_POWER_STATE
    MAX_NUM_MEASUREMENTS     MeasurementTypeEnum = C.MAX_NUM_MEASUREMENTS
)

type ChargerStateEnum uint32

const (
    CHARGER_UNKNOWN ChargerStateEnum = C.CHARGER_UNKNOWN
    CHARGER_OFF     ChargerStateEnum = C.CHARGER_OFF
    CHARGER_ON      ChargerStateEnum = C.CHARGER_ON
    CHARGER_FAULT   ChargerStateEnum = C.CHARGER_FAULT
)

type MessageContainer struct {
	DeviceUuid   string
    Timestamp    time.Time
    Message      interface{}
}

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

const MaxDatagramSizeRaw uint32 = C.MAX_DATAGRAM_SIZE_RAW

const RevisionLevel uint32 = C.REVISION_LEVEL

var totalUlMsgs int
var totalUlBytes int
var lastUlMsgTime time.Time

var ulDecodeTypeDisplay map[int]string = map[int]string {
	C.DECODE_RESULT_FAILURE:                                      "DECODE_RESULT_FAILURE",
	C.DECODE_RESULT_INPUT_TOO_SHORT:                              "DECODE_RESULT_INPUT_TOO_SHORT",
	C.DECODE_RESULT_OUTPUT_TOO_SHORT:                             "DECODE_RESULT_OUTPUT_TOO_SHORT",
	C.DECODE_RESULT_UNKNOWN_MSG_ID:                               "DECODE_RESULT_UNKNOWN_MSG_ID",
	C.DECODE_RESULT_BAD_MSG_FORMAT:                               "DECODE_RESULT_BAD_MSG_FORMAT",
	C.DECODE_RESULT_BAD_TRAFFIC_TEST_MODE_DATAGRAM:               "DECODE_RESULT_BAD_TRAFFIC_TEST_MODE_DATAGRAM",
	C.DECODE_RESULT_OUT_OF_SEQUENCE_TRAFFIC_TEST_MODE_DATAGRAM:   "DECODE_RESULT_OUT_OF_SEQUENCE_TRAFFIC_TEST_MODE_DATAGRAM",
	C.DECODE_RESULT_BAD_CHECKSUM:                                 "DECODE_RESULT_BAD_CHECKSUM",
	C.DECODE_RESULT_TRANSPARENT_UL_DATAGRAM:                      "DECODE_RESULT_TRANSPARENT_UL_DATAGRAM",
	C.DECODE_RESULT_PING_REQ_UL_MSG:                              "DECODE_RESULT_PING_REQ_UL_MSG",
	C.DECODE_RESULT_PING_CNF_UL_MSG:                              "DECODE_RESULT_PING_CNF_UL_MSG",
	C.DECODE_RESULT_INIT_IND_UL_MSG:                              "DECODE_RESULT_INIT_IND_UL_MSG",
	C.DECODE_RESULT_DATE_TIME_IND_UL_MSG:                         "DECODE_RESULT_DATE_TIME_IND_UL_MSG",
	C.DECODE_RESULT_DATE_TIME_SET_CNF_UL_MSG:                     "DECODE_RESULT_DATE_TIME_SET_CNF_UL_MSG",
	C.DECODE_RESULT_DATE_TIME_GET_CNF_UL_MSG:                     "DECODE_RESULT_DATE_TIME_GET_CNF_UL_MSG",
	C.DECODE_RESULT_MODE_SET_CNF_UL_MSG:                          "DECODE_RESULT_MODE_SET_CNF_UL_MSG",
	C.DECODE_RESULT_MODE_GET_CNF_UL_MSG:                          "DECODE_RESULT_MODE_GET_CNF_UL_MSG",
	C.DECODE_RESULT_HEARTBEAT_SET_CNF_UL_MSG:                     "DECODE_RESULT_HEARTBEAT_SET_CNF_UL_MSG",
	C.DECODE_RESULT_REPORTING_INTERVAL_SET_CNF_UL_MSG:            "DECODE_RESULT_REPORTING_INTERVAL_SET_CNF_UL_MSG",
	C.DECODE_RESULT_INTERVALS_GET_CNF_UL_MSG:                     "DECODE_RESULT_INTERVALS_GET_CNF_UL_MSG",
	C.DECODE_RESULT_POLL_IND_UL_MSG:                              "DECODE_RESULT_POLL_IND_UL_MSG",
	C.DECODE_RESULT_MEASUREMENTS_IND_UL_MSG:                      "DECODE_RESULT_MEASUREMENTS_IND_UL_MSG",
	C.DECODE_RESULT_MEASUREMENTS_GET_CNF_UL_MSG:                  "DECODE_RESULT_MEASUREMENTS_GET_CNF_UL_MSG",
	C.DECODE_RESULT_MEASUREMENTS_CONTROL_IND_UL_MSG:              "DECODE_RESULT_MEASUREMENTS_CONTROL_IND_UL_MSG",
	C.DECODE_RESULT_MEASUREMENT_CONTROL_SET_CNF_UL_MSG:           "DECODE_RESULT_MEASUREMENT_CONTROL_SET_CNF_UL_MSG",
	C.DECODE_RESULT_MEASUREMENTS_CONTROL_GET_CNF_UL_MSG:          "DECODE_RESULT_MEASUREMENTS_CONTROL_GET_CNF_UL_MSG",
	C.DECODE_RESULT_MEASUREMENTS_CONTROL_DEFAULTS_SET_CNF_UL_MSG: "DECODE_RESULT_MEASUREMENTS_CONTROL_DEFAULTS_SET_CNF_UL_MSG",
	C.DECODE_RESULT_TRAFFIC_REPORT_IND_UL_MSG:                    "DECODE_RESULT_TRAFFIC_REPORT_IND_UL_MSG",
	C.DECODE_RESULT_TRAFFIC_REPORT_GET_CNF_UL_MSG:                "DECODE_RESULT_TRAFFIC_REPORT_GET_CNF_UL_MSG",
	C.DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF_UL_MSG:  "DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF_UL_MSG",
	C.DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF_UL_MSG:  "DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF_UL_MSG",
	C.DECODE_RESULT_TRAFFIC_TEST_MODE_RULE_BREAKER_UL_DATAGRAM:   "DECODE_RESULT_TRAFFIC_TEST_MODE_RULE_BREAKER_UL_DATAGRAM",
	C.DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_IND_UL_MSG:          "DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_IND_UL_MSG",
	C.DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_GET_CNF_UL_MSG:      "DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_GET_CNF_UL_MSG",
	C.DECODE_RESULT_ACTIVITY_REPORT_IND_UL_MSG:                   "DECODE_RESULT_ACTIVITY_REPORT_IND_UL_MSG",
	C.DECODE_RESULT_ACTIVITY_REPORT_GET_CNF_UL_MSG:               "DECODE_RESULT_ACTIVITY_REPORT_GET_CNF_UL_MSG",
	C.DECODE_RESULT_DEBUG_IND_UL_MSG:                             "DECODE_RESULT_DEBUG_IND_UL_MSG",
}

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// The decode function.  Returns an array of decoded
// messages and a byte count
func decode(data []byte, uuid string) ([]interface{}, int) {

	var bytesRemaining C.uint32_t = 1
	var returnedMsgs []interface{} = nil
	var used C.uint32_t
    var inputBuffer C.union_UlMsgUnionTag_t
    var xmlDecodeBuffer = make([]byte, 8192, 8192)

	// Holder for the extracted message
	pBuffer := (*C.UlMsgUnion_t)(unsafe.Pointer(&inputBuffer))
	
	// Loop over the messages in the datagram
	pStart := (*C.char)(unsafe.Pointer(&data[0]))
	pNext := pStart
	ppNext := (**C.char)(unsafe.Pointer(&pNext))

	hexBuffer := hex.Dump(data)
	globals.Dbg.PrintfInfo("%s [decode] --> the whole input buffer:\n\n%s\n\n", globals.LogTag, hexBuffer)

	var decoderCount uint32

	for bytesRemaining > 0 {
		decoderCount = decoderCount + 1

		globals.Dbg.PrintfTrace("%s [decode] --> decoding message number %d of AMQP message number %d.\n", globals.LogTag, decoderCount, amqpMessageCount)

		globals.Dbg.PrintfInfo("%s -----## show buffer data ##-----\n\n%s\n\n", globals.LogTag, spew.Sdump(inputBuffer))
		used = C.pointerSub(pNext, pStart)
		globals.Dbg.PrintfInfo("%s [decode] --> message data used: %d.\n", globals.LogTag, used)
		bytesRemaining = C.uint32_t(len(data)) - used
		globals.Dbg.PrintfInfo("%s [decode] --> %d bytes remaining out of %d.\n", globals.LogTag, bytesRemaining, len(data))
		if bytesRemaining > 0 {
        	// A place to put the XML output from the decoder
        	pXmlBuffer := (*C.char) (unsafe.Pointer(&(xmlDecodeBuffer[0])))
        	ppXmlBuffer := (**C.char) (unsafe.Pointer(&pXmlBuffer))
            xmlBufferLen := (C.uint32_t) (len(xmlDecodeBuffer))
            pXmlBufferLen := (*C.uint32_t) (unsafe.Pointer(&xmlBufferLen))
	
			result := C.decodeUlMsg(ppNext, bytesRemaining, pBuffer, ppXmlBuffer, pXmlBufferLen)
			globals.Dbg.PrintfTrace("%s [decode] --> decode received uplink message: %+v.\n\n-----## %s ##----- \n\n", globals.LogTag, result, ulDecodeTypeDisplay[int(result)])
			globals.Dbg.PrintfInfo("%s [decode] --> XML buffer pointer 0x%08x, used %d, left %d:.\n", globals.LogTag, *ppXmlBuffer, C.uint32_t(len(xmlDecodeBuffer)) - xmlBufferLen, xmlBufferLen)
	
    		// Store XmlData in MongoDB
    		utilities.XmlDataStore(xmlDecodeBuffer, uuid)
    		globals.Dbg.PrintfInfo("%s [decode] --> the XML data is:\n\n%s\n\n", globals.LogTag, spew.Sdump(xmlDecodeBuffer))
    		
			// Now decode the messages and pass them to the state table
			switch int(result) {
				case C.DECODE_RESULT_TRANSPARENT_UL_DATAGRAM:
    				// TODO
				case C.DECODE_RESULT_PING_REQ_UL_MSG:
	    			// Empty message
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &PingReqUlMsg {},   
    					})
				case C.DECODE_RESULT_PING_CNF_UL_MSG:
	    			// Empty message
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &PingCnfUlMsg {},   
    					})
				case C.DECODE_RESULT_INIT_IND_UL_MSG:
					value := C.getInitIndUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &InitIndUlMsg {
    							WakeUpCode:        WakeUpEnum(value.wakeUpCode),
    							RevisionLevel:     uint8(value.revisionLevel),
    							SdCardNotRequired: bool(value.sdCardNotRequired),
    							DisableModemDebug: bool(value.disableModemDebug),
    							DisableButton:     bool(value.disableButton),
    							DisableServerPing: bool(value.disableServerPing),
    						},
    					})
					
				case C.DECODE_RESULT_DATE_TIME_IND_UL_MSG:
					value := C.getDateTimeIndUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &DateTimeIndUlMsg {
    							UtmTime:        time.Unix(int64(value.time), 0).Local(),
    							TimeSetBy:      TimeSetByEnum(value.setBy),
    						},
    					})
				
				case C.DECODE_RESULT_DATE_TIME_SET_CNF_UL_MSG:
					value := C.getDateTimeSetCnfUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &DateTimeSetCnfUlMsg {
    							UtmTime:        time.Unix(int64(value.time), 0).Local(),
    							TimeSetBy:      TimeSetByEnum(value.setBy),
    						},
    					})
					
				case C.DECODE_RESULT_DATE_TIME_GET_CNF_UL_MSG:
					value := C.getDateTimeGetCnfUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &DateTimeGetCnfUlMsg {
    							UtmTime:        time.Unix(int64(value.time), 0).Local(),
    							TimeSetBy:      TimeSetByEnum(value.setBy),
    						},
    					})
					
				case C.DECODE_RESULT_MODE_SET_CNF_UL_MSG:
					value := C.getModeSetCnfUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &ModeSetCnfUlMsg {
    					        Mode:           ModeEnum(value.mode),
    					    },   
    					})
					
				case C.DECODE_RESULT_MODE_GET_CNF_UL_MSG:
					value := C.getModeGetCnfUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &ModeGetCnfUlMsg {
    					        Mode:           ModeEnum(value.mode),
    					    },   
    					})
					
				case C.DECODE_RESULT_HEARTBEAT_SET_CNF_UL_MSG:
					value := C.getHeartbeatSetCnfUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &HeartbeatSetCnfUlMsg {
    					        HeartbeatSeconds:   uint32(value.heartbeatSeconds),
    					        HeartbeatSnapToRtc: bool(value.heartbeatSnapToRtc),
    					    },   
    					})

				case C.DECODE_RESULT_REPORTING_INTERVAL_SET_CNF_UL_MSG:
					value := C.getReportingIntervalSetCnfUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &ReportingIntervalSetCnfUlMsg {
    					        ReportingInterval:  uint32(value.reportingInterval),
    					    },   
    					})
					
				case C.DECODE_RESULT_INTERVALS_GET_CNF_UL_MSG:
					value := C.getIntervalsGetCnfUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &IntervalsGetCnfUlMsg {
    							ReportingInterval:  uint32(value.reportingInterval),
    							HeartbeatSeconds:   uint32(value.heartbeatSeconds),
    							HeartbeatSnapToRtc: bool(value.heartbeatSnapToRtc),
    					    },   
    					})
	
				case C.DECODE_RESULT_POLL_IND_UL_MSG:
					value := C.getPollIndUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &PollIndUlMsg {
    					        Mode:            ModeEnum(value.mode),
    					        EnergyLeft:      EnergyLeftEnum(value.energyLeft),
    					        DiskSpaceLeft:   DiskSpaceLeftEnum(value.diskSpaceLeft),
    					    },   
    					})
					
				case C.DECODE_RESULT_MEASUREMENTS_IND_UL_MSG:
					value := C.getMeasurementsIndUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &MeasurementsIndUlMsg {  // Structure initialisation is just horrible in this language
    							Measurements: MeasurementData {
    								TimeMeasured:         time.Unix(int64(value.measurements.time), 0).Local(),
    								GnssPositionPresent:  bool(value.measurements.gnssPositionPresent),
    								GnssPosition: GnssPosition {
    									Latitude:            int32(value.measurements.gnssPosition.latitude),
    									Longitude:           int32(value.measurements.gnssPosition.longitude),
    									Elevation:           int32(value.measurements.gnssPosition.elevation),
    								},
    								CellIdPresent:       bool(value.measurements.cellIdPresent),
    								CellId:              CellId(value.measurements.cellId),
    								RsrpPresent:         bool(value.measurements.rsrpPresent),
    								Rsrp: Rsrp {
    									Value:              Rssi(value.measurements.rsrp.value),
    									IsSyncedWithRssi:   bool(value.measurements.rsrp.isSyncedWithRssi),
    								},
    								RssiPresent:         bool(value.measurements.rssiPresent),
    								Rssi:                Rssi(value.measurements.rssi),
    								TemperaturePresent:  bool(value.measurements.temperaturePresent),
    								Temperature:         Temperature(value.measurements.temperature),
    								PowerStatePresent:   bool(value.measurements.powerStatePresent),
    								PowerState: PowerState {
    									ChargerState:        ChargerStateEnum(value.measurements.powerState.chargerState),
    									BatteryMv:           uint16(value.measurements.powerState.batteryMV),
    									EnergyMwh:           uint32(value.measurements.powerState.energyMWh),
    								},
    							},
    						},	
    					})
					
				case C.DECODE_RESULT_MEASUREMENTS_GET_CNF_UL_MSG:
					value := C.getMeasurementsGetCnfUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &MeasurementsIndUlMsg {  // Still horrible
    							Measurements: MeasurementData {
    								TimeMeasured:         time.Unix(int64(value.measurements.time), 0).Local(),
    								GnssPositionPresent:  bool(value.measurements.gnssPositionPresent),
    								GnssPosition: GnssPosition {
    									Latitude:            int32(value.measurements.gnssPosition.latitude),
    									Longitude:           int32(value.measurements.gnssPosition.longitude),
    									Elevation:           int32(value.measurements.gnssPosition.elevation),
    								},
    								CellIdPresent:       bool(value.measurements.cellIdPresent),
    								CellId:              CellId(value.measurements.cellId),
    								RsrpPresent:         bool(value.measurements.rsrpPresent),
    								Rsrp: Rsrp {
    									Value:              Rssi(value.measurements.rsrp.value),
    									IsSyncedWithRssi:   bool(value.measurements.rsrp.isSyncedWithRssi),
    								},
    								RssiPresent:         bool(value.measurements.rssiPresent),
    								Rssi:                Rssi(value.measurements.rssi),
    								TemperaturePresent:  bool(value.measurements.temperaturePresent),
    								Temperature:         Temperature(value.measurements.temperature),
    								PowerStatePresent:   bool(value.measurements.powerStatePresent),
    								PowerState: PowerState {
    									ChargerState:        ChargerStateEnum(value.measurements.powerState.chargerState),
    									BatteryMv:           uint16(value.measurements.powerState.batteryMV),
    									EnergyMwh:           uint32(value.measurements.powerState.energyMWh),
    								},
    							},
    						},	
    					})
					
				case C.DECODE_RESULT_MEASUREMENTS_CONTROL_IND_UL_MSG:
				// TODO
				case C.DECODE_RESULT_MEASUREMENT_CONTROL_SET_CNF_UL_MSG:
				// TODO
				case C.DECODE_RESULT_MEASUREMENTS_CONTROL_GET_CNF_UL_MSG:
				// TODO
				case C.DECODE_RESULT_MEASUREMENTS_CONTROL_DEFAULTS_SET_CNF_UL_MSG:
				// TODO
				case C.DECODE_RESULT_TRAFFIC_REPORT_IND_UL_MSG:
					value := C.getTrafficReportIndUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &TrafficReportIndUlMsg {
    					        NumDatagramsUl:             uint32(value.numDatagramsUl),
    					        NumBytesUl:                 uint32(value.numBytesUl),
    					        NumDatagramsDl:             uint32(value.numDatagramsDl),
    					        NumBytesDl:                 uint32(value.numBytesDl),
    					        NumDatagramsDlBadChecksum:  uint32(value.numDatagramsDlBadChecksum),
    					    },   
    					})
					
				case C.DECODE_RESULT_TRAFFIC_REPORT_GET_CNF_UL_MSG:
					value := C.getTrafficReportGetCnfUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &TrafficReportGetCnfUlMsg {
    					        NumDatagramsUl:             uint32(value.numDatagramsUl),
    					        NumBytesUl:                 uint32(value.numBytesUl),
    					        NumDatagramsDl:             uint32(value.numDatagramsDl),
    					        NumBytesDl:                 uint32(value.numBytesDl),
    					        NumDatagramsDlBadChecksum:  uint32(value.numDatagramsDlBadChecksum),
    					    },   
    					})
					
				case C.DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF_UL_MSG:
					value := C.getTrafficTestModeParametersSetCnfUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &TrafficTestModeParametersSetCnfUlMsg {
    					        NumUlDatagrams:      uint32(value.numUlDatagrams),
    					        LenUlDatagram:       uint32(value.lenUlDatagram),
    					        NumDlDatagrams:      uint32(value.numDlDatagrams),
    					        LenDlDatagram:       uint32(value.lenDlDatagram),
    					    },   
    					})
					
				case C.DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF_UL_MSG:
					value := C.getTrafficTestModeParametersGetCnfUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &TrafficTestModeParametersGetCnfUlMsg {
    					        NumUlDatagrams:      uint32(value.numUlDatagrams),
    					        LenUlDatagram:       uint32(value.lenUlDatagram),
    					        NumDlDatagrams:      uint32(value.numDlDatagrams),
    					        LenDlDatagram:       uint32(value.lenDlDatagram),
    					    },   
    					})
					
				case C.DECODE_RESULT_TRAFFIC_TEST_MODE_RULE_BREAKER_UL_DATAGRAM:
				// TODO
				
				case C.DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_IND_UL_MSG:
					value := C.getTrafficTestModeReportIndUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &TrafficTestModeReportIndUlMsg {
    					        NumTrafficTestDatagramsUl:           uint32(value.numTrafficTestDatagramsUl),
    					        NumTrafficTestBytesUl:               uint32(value.numTrafficTestBytesUl),
    					        NumTrafficTestDatagramsDl:           uint32(value.numTrafficTestDatagramsDl),
    					        NumTrafficTestBytesDl:               uint32(value.numTrafficTestBytesDl),
    					        NumTrafficTestDlDatagramsOutOfOrder: uint32(value.numTrafficTestDlDatagramsOutOfOrder),
    					        NumTrafficTestDlDatagramsBad:        uint32(value.numTrafficTestDlDatagramsBad),
    					        NumTrafficTestDlDatagramsMissed:     uint32(value.numTrafficTestDlDatagramsMissed),
    					        TimedOut:                            bool(value.timedOut),
    					    },   
    					})
					
				case C.DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_GET_CNF_UL_MSG:
					value := C.getTrafficTestModeReportGetCnfUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &TrafficTestModeReportGetCnfUlMsg {
    					        NumTrafficTestDatagramsUl:           uint32(value.numTrafficTestDatagramsUl),
    					        NumTrafficTestBytesUl:               uint32(value.numTrafficTestBytesUl),
    					        NumTrafficTestDatagramsDl:           uint32(value.numTrafficTestDatagramsDl),
    					        NumTrafficTestBytesDl:               uint32(value.numTrafficTestBytesDl),
    					        NumTrafficTestDlDatagramsOutOfOrder: uint32(value.numTrafficTestDlDatagramsOutOfOrder),
    					        NumTrafficTestDlDatagramsBad:        uint32(value.numTrafficTestDlDatagramsBad),
    					        NumTrafficTestDlDatagramsMissed:     uint32(value.numTrafficTestDlDatagramsMissed),
    					        TimedOut:                            bool(value.timedOut),
    					    },   
    					})
					
				case C.DECODE_RESULT_ACTIVITY_REPORT_IND_UL_MSG:
					value := C.getActivityReportIndUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
        					DeviceUuid:        uuid,
        					Timestamp:         time.Now(),
        					Message: &ActivityReportIndUlMsg {
        				        TotalTransmitMilliseconds:  uint32(value.totalTransmitMilliseconds),
        				        TotalReceiveMilliseconds:   uint32(value.totalReceiveMilliseconds),
        				        UpTimeSeconds:              uint32(value.upTimeSeconds),
        				        TxPowerDbmPresent:          bool(value.txPowerDbmPresent),
        				        TxPowerDbm:                 int8(value.txPowerDbm),
        				        UlMcsPresent:               bool(value.ulMcsPresent),
        				        UlMcs:                      uint8(value.ulMcs),
        				        DlMcsPresent:               bool(value.dlMcsPresent),
        				        DlMcs:                      uint8(value.dlMcs),
        				    },   
        				})
					
				case C.DECODE_RESULT_ACTIVITY_REPORT_GET_CNF_UL_MSG:
					value := C.getActivityReportGetCnfUlMsg(inputBuffer)
					returnedMsgs = append (returnedMsgs,
    				    &MessageContainer {
    						DeviceUuid:        uuid,
    						Timestamp:         time.Now(),
    						Message: &ActivityReportGetCnfUlMsg {
    					        TotalTransmitMilliseconds:  uint32(value.totalTransmitMilliseconds),
    					        TotalReceiveMilliseconds:   uint32(value.totalReceiveMilliseconds),
    					        UpTimeSeconds:              uint32(value.upTimeSeconds),
    					        TxPowerDbmPresent:          bool(value.txPowerDbmPresent),
    					        TxPowerDbm:                 int8(value.txPowerDbm),
    					        UlMcsPresent:               bool(value.ulMcsPresent),
    					        UlMcs:                      uint8(value.ulMcs),
    					        DlMcsPresent:               bool(value.dlMcsPresent),
    					        DlMcs:                      uint8(value.dlMcs),
    					    },   
    					})
					
				case C.DECODE_RESULT_DEBUG_IND_UL_MSG:
				// TODO
				case C.DECODE_RESULT_FAILURE:
				case C.DECODE_RESULT_INPUT_TOO_SHORT:
				case C.DECODE_RESULT_OUTPUT_TOO_SHORT:
				case C.DECODE_RESULT_UNKNOWN_MSG_ID:
				case C.DECODE_RESULT_BAD_MSG_FORMAT:
				case C.DECODE_RESULT_BAD_TRAFFIC_TEST_MODE_DATAGRAM:
				// TODO
				case C.DECODE_RESULT_OUT_OF_SEQUENCE_TRAFFIC_TEST_MODE_DATAGRAM:
				// TODO
				case C.DECODE_RESULT_BAD_CHECKSUM:
				default:
				// Can't decode the message, throw it away silently
			}	
		}
	}

	return returnedMsgs, int(used)
}

/* End Of File */
