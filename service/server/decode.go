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
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"time"
	"unsafe"
)

//--------------------------------------------------------------------
// Enums
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

//--------------------------------------------------------------------
// Vars
//--------------------------------------------------------------------

const MaxDatagramSizeRaw uint32 = C.MAX_DATAGRAM_SIZE_RAW

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
// The decode function
//--------------------------------------------------------------------

func decode(data []byte) {

	var bytesRemaining C.uint32_t = 1
	
	// Holder for the extracted message
	var inputBuffer C.union_UlMsgUnionTag_t
	pBuffer := (*C.UlMsgUnion_t)(unsafe.Pointer(&inputBuffer))
	
	row.UTotalMsgs = row.UTotalMsgs + uint64(len(data))

	// Loop over the messages in the datagram
	pStart := (*C.char)(unsafe.Pointer(&data[0]))
	pNext := pStart
	ppNext := (**C.char)(unsafe.Pointer(&pNext))

	hexBuffer := hex.Dump(data)
	fmt.Printf("\n\n%s --> the whole input buffer:\n(%s).\n", logTag, hexBuffer)

	var decoderCount uint32

	for bytesRemaining > 0 {
		decoderCount = decoderCount + 1
		now := time.Now()
		row.LastMsgReceived = &now

		fmt.Printf("\n\n%s --> decoding message number (%v) of AMQP datagram number (%v).\n", logTag, decoderCount, amqpCount)

		fmt.Printf("\n%s -----## show buffer data ##-----\n\n%s\n\n", logTag, spew.Sdump(inputBuffer))
		used := C.pointerSub(pNext, pStart)
		fmt.Printf("%s --> message data used:\n\n%s\n", logTag, spew.Sdump(used))
		bytesRemaining = C.uint32_t(len(data)) - used
		fmt.Printf("%s --> message data remaining:\n\n%s\n", logTag, spew.Sdump(bytesRemaining))
		if bytesRemaining > 0 {
			result := C.decodeUlMsg(ppNext, bytesRemaining, pBuffer, nil, nil)
			fmt.Printf("%s --> decode received uplink message:\n\n%+v.\n\n -----## %s ##----- \n\n", logTag, result, ulDecodeTypeDisplay[int(result)])
	
			// Extract any data to be recorded; the C symbols are not available outside
			// the package so convert into concrete go types
			var data interface{} = nil
			var rawData interface{} = nil
	
			// Now decode the messages and pass them to the state table
			switch int(result) {
				case C.DECODE_RESULT_TRANSPARENT_UL_DATAGRAM:
				// TODO
				case C.DECODE_RESULT_PING_REQ_UL_MSG:
				// Empty message
				// TODO: send a pingReqDlMsg
				case C.DECODE_RESULT_PING_CNF_UL_MSG:
				// Ignored
				case C.DECODE_RESULT_INIT_IND_UL_MSG:
					value := C.getInitIndUlMsg(inputBuffer)
					rawData = value
					data = &InitIndUlMsg {
						Timestamp:         time.Now(),
						WakeUpCode:        WakeUpEnum(value.wakeUpCode),
						RevisionLevel:     uint8(value.revisionLevel),
						SdCardNotRequired: bool(value.sdCardNotRequired),
						DisableModemDebug: bool(value.disableModemDebug),
						DisableButton:     bool(value.disableButton),
						DisableServerPing: bool(value.disableServerPing),
					}
					
				case C.DECODE_RESULT_DATE_TIME_IND_UL_MSG:
					value := C.getDateTimeIndUlMsg(inputBuffer)
					rawData = value
					data = &DateTimeIndUlMsg {
						Timestamp:      time.Now(),
						UtmTime:        time.Unix(int64(value.time), 0).Local(),
						TimeSetBy:      TimeSetByEnum(value.setBy),
					}
				
				case C.DECODE_RESULT_DATE_TIME_SET_CNF_UL_MSG:
					value := C.getDateTimeSetCnfUlMsg(inputBuffer)
					rawData = value
					data = &DateTimeSetCnfUlMsg {
						Timestamp:      time.Now(),
						UtmTime:        time.Unix(int64(value.time), 0).Local(),
						TimeSetBy:      TimeSetByEnum(value.setBy),
					}
					
				case C.DECODE_RESULT_DATE_TIME_GET_CNF_UL_MSG:
					value := C.getDateTimeGetCnfUlMsg(inputBuffer)
					rawData = value
					data = &DateTimeGetCnfUlMsg {
						Timestamp:      time.Now(),
						UtmTime:        time.Unix(int64(value.time), 0).Local(),
						TimeSetBy:      TimeSetByEnum(value.setBy),
					}
					
				case C.DECODE_RESULT_MODE_SET_CNF_UL_MSG:
					value := C.getModeSetCnfUlMsg(inputBuffer)
					rawData = value
					data = &ModeSetCnfUlMsg {
						Timestamp:      time.Now(),
				        Mode:           ModeEnum(value.mode),
					}
					
				case C.DECODE_RESULT_MODE_GET_CNF_UL_MSG:
					value := C.getModeGetCnfUlMsg(inputBuffer)
					rawData = value
					data = &ModeGetCnfUlMsg {
						Timestamp:      time.Now(),
				        Mode:           ModeEnum(value.mode),
					}
					
				case C.DECODE_RESULT_HEARTBEAT_SET_CNF_UL_MSG:
					value := C.getHeartbeatSetCnfUlMsg(inputBuffer)
					rawData = value
					data = &HeartbeatSetCnfUlMsg {
						Timestamp:          time.Now(),
				        HeartbeatSeconds:   uint32(value.heartbeatSeconds),
				        HeartbeatSnapToRtc: bool(value.heartbeatSnapToRtc),
					}

				case C.DECODE_RESULT_REPORTING_INTERVAL_SET_CNF_UL_MSG:
					value := C.getReportingIntervalSetCnfUlMsg(inputBuffer)
					rawData = value
					data = &ReportingIntervalSetCnfUlMsg {
						Timestamp:          time.Now(),
				        ReportingInterval:  uint32(value.reportingInterval),
					}
					
				case C.DECODE_RESULT_INTERVALS_GET_CNF_UL_MSG:
					value := C.getIntervalsGetCnfUlMsg(inputBuffer)
					rawData = value
					data = &IntervalsGetCnfUlMsg {
						ReportingInterval:  uint32(value.reportingInterval),
						HeartbeatSeconds:   uint32(value.heartbeatSeconds),
						HeartbeatSnapToRtc: bool(value.heartbeatSnapToRtc),
					}
	
					// TODO
					row.ReportingInterval = uint32(value.reportingInterval)
					row.HeartbeatSeconds = uint32(value.heartbeatSeconds)
					row.HeartbeatSnapToRtc = bool(value.heartbeatSnapToRtc)
	
					fmt.Printf("%s RECEIVED INTERVAL REQUEST CONFIRM %s\n", logTag, spew.Sdump(row))
	
				case C.DECODE_RESULT_POLL_IND_UL_MSG:
					value := C.getPollIndUlMsg(inputBuffer)
					rawData = value
					data = &PollIndUlMsg {
						Timestamp:       time.Now(),
				        Mode:            ModeEnum(value.mode),
				        EnergyLeft:      EnergyLeftEnum(value.energyLeft),
				        DiskSpaceLeft:   DiskSpaceLeftEnum(value.diskSpaceLeft),
					}
					
				case C.DECODE_RESULT_MEASUREMENTS_IND_UL_MSG:
					value := C.getMeasurementsIndUlMsg(inputBuffer)
					rawData = value
					data = &MeasurementsIndUlMsg {  // Structure initialisation is just horrible in this language
						Timestamp:                   time.Now(),
						Measurements: MeasurementData {
							TimeMeasured:         time.Unix(int64(value.measurements.time), 0).Local(),
							GnssPositionPresent:  bool(value.measurements.gnssPositionPresent),
							GnssPosition: GnssPosition {
								Timestamp:           time.Now(),
								Latitude:            int32(value.measurements.gnssPosition.latitude),
								Longitude:           int32(value.measurements.gnssPosition.longitude),
								Elevation:           int32(value.measurements.gnssPosition.elevation),
							},
							CellIdPresent:       bool(value.measurements.cellIdPresent),
							CellId:              uint16(value.measurements.cellId),
							RsrpPresent:         bool(value.measurements.rsrpPresent),
							Rsrp: Rsrp {
								Value: Rssi {
									Timestamp:         time.Now(),
									Rssi:              int16(value.measurements.rsrp.value),
								},
								IsSyncedWithRssi:   bool(value.measurements.rsrp.isSyncedWithRssi),
							},
							RssiPresent:         bool(value.measurements.rssiPresent),
							Rssi: Rssi {
								Timestamp:          time.Now(),
								Rssi:               int16(value.measurements.rsrp.value),
							},
							TemperaturePresent:  bool(value.measurements.temperaturePresent),
							Temperature:         int8(value.measurements.temperature),
							PowerStatePresent:   bool(value.measurements.powerStatePresent),
							PowerState: PowerState {
								ChargerState:        ChargerStateEnum(value.measurements.powerState.chargerState),
								BatteryMv:           uint16(value.measurements.powerState.batteryMV),
								EnergyMwh:           uint32(value.measurements.powerState.energyMWh),
							},
						},
					}
					
				case C.DECODE_RESULT_MEASUREMENTS_GET_CNF_UL_MSG:
					value := C.getMeasurementsGetCnfUlMsg(inputBuffer)
					rawData = value
					data = &MeasurementsGetCnfUlMsg {  // Still horrible
						Timestamp:          time.Now(),
						Measurements: MeasurementData {
							TimeMeasured:         time.Unix(int64(value.measurements.time), 0).Local(),
							GnssPositionPresent:  bool(value.measurements.gnssPositionPresent),
							GnssPosition: GnssPosition {
								Timestamp:           time.Now(),
								Latitude:            int32(value.measurements.gnssPosition.latitude),
								Longitude:           int32(value.measurements.gnssPosition.longitude),
								Elevation:           int32(value.measurements.gnssPosition.elevation),
							},
							CellIdPresent:       bool(value.measurements.cellIdPresent),
							CellId:              uint16(value.measurements.cellId),
							RsrpPresent:         bool(value.measurements.rsrpPresent),
							Rsrp: Rsrp {
								Value: Rssi {
									Timestamp:         time.Now(),
									Rssi:              int16(value.measurements.rsrp.value),
								},
								IsSyncedWithRssi:   bool(value.measurements.rsrp.isSyncedWithRssi),
							},
							RssiPresent:         bool(value.measurements.rssiPresent),
							Rssi: Rssi {
								Timestamp:          time.Now(),
								Rssi:               int16(value.measurements.rsrp.value),
							},
							TemperaturePresent:  bool(value.measurements.temperaturePresent),
							Temperature:         int8(value.measurements.temperature),
							PowerStatePresent:   bool(value.measurements.powerStatePresent),
							PowerState: PowerState {
								ChargerState:        ChargerStateEnum(value.measurements.powerState.chargerState),
								BatteryMv:           uint16(value.measurements.powerState.batteryMV),
								EnergyMwh:           uint32(value.measurements.powerState.energyMWh),
							},
						},
					}
					
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
					rawData = value
					data = &TrafficReportIndUlMsg {
						Timestamp:                  time.Now(),
				        NumDatagramsUl:             uint32(value.numDatagramsUl),
				        NumBytesUl:                 uint32(value.numBytesUl),
				        NumDatagramsDl:             uint32(value.numDatagramsDl),
				        NumBytesDl:                 uint32(value.numBytesDl),
				        NumDatagramsDlBadChecksum:  uint32(value.numDatagramsDlBadChecksum),
					}
					
				case C.DECODE_RESULT_TRAFFIC_REPORT_GET_CNF_UL_MSG:
					value := C.getTrafficReportGetCnfUlMsg(inputBuffer)
					rawData = value
					data = &TrafficReportGetCnfUlMsg {
						Timestamp:                  time.Now(),
				        NumDatagramsUl:             uint32(value.numDatagramsUl),
				        NumBytesUl:                 uint32(value.numBytesUl),
				        NumDatagramsDl:             uint32(value.numDatagramsDl),
				        NumBytesDl:                 uint32(value.numBytesDl),
				        NumDatagramsDlBadChecksum:  uint32(value.numDatagramsDlBadChecksum),
					}
					
				case C.DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF_UL_MSG:
					value := C.getTrafficTestModeParametersSetCnfUlMsg(inputBuffer)
					rawData = value
					data = &TrafficTestModeParametersSetCnfUlMsg {
						Timestamp:           time.Now(),
				        NumUlDatagrams:      uint32(value.numUlDatagrams),
				        LenUlDatagram:       uint32(value.lenUlDatagram),
				        NumDlDatagrams:      uint32(value.numDlDatagrams),
				        LenDlDatagram:       uint32(value.lenDlDatagram),
					}
					
				case C.DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF_UL_MSG:
					value := C.getTrafficTestModeParametersGetCnfUlMsg(inputBuffer)
					rawData = value
					data = &TrafficTestModeParametersGetCnfUlMsg {
						Timestamp:           time.Now(),
				        NumUlDatagrams:      uint32(value.numUlDatagrams),
				        LenUlDatagram:       uint32(value.lenUlDatagram),
				        NumDlDatagrams:      uint32(value.numDlDatagrams),
				        LenDlDatagram:       uint32(value.lenDlDatagram),
					}
					
				case C.DECODE_RESULT_TRAFFIC_TEST_MODE_RULE_BREAKER_UL_DATAGRAM:
				// TODO
				case C.DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_IND_UL_MSG:
					value := C.getTrafficTestModeReportIndUlMsg(inputBuffer)
					rawData = value
					data = &TrafficTestModeReportIndUlMsg {
						Timestamp:                           time.Now(),
				        NumTrafficTestDatagramsUl:           uint32(value.numTrafficTestDatagramsUl),
				        NumTrafficTestBytesUl:               uint32(value.numTrafficTestBytesUl),
				        NumTrafficTestDatagramsDl:           uint32(value.numTrafficTestDatagramsDl),
				        NumTrafficTestBytesDl:               uint32(value.numTrafficTestBytesDl),
				        NumTrafficTestDlDatagramsOutOfOrder: uint32(value.numTrafficTestDlDatagramsOutOfOrder),
				        NumTrafficTestDlDatagramsBad:        uint32(value.numTrafficTestDlDatagramsBad),
				        NumTrafficTestDlDatagramsMissed:     uint32(value.numTrafficTestDlDatagramsMissed),
				        TimedOut:                            bool(value.timedOut),
					}
					
				case C.DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_GET_CNF_UL_MSG:
					value := C.getTrafficTestModeReportGetCnfUlMsg(inputBuffer)
					rawData = value
					data = &TrafficTestModeReportGetCnfUlMsg {
						Timestamp:                           time.Now(),
				        NumTrafficTestDatagramsUl:           uint32(value.numTrafficTestDatagramsUl),
				        NumTrafficTestBytesUl:               uint32(value.numTrafficTestBytesUl),
				        NumTrafficTestDatagramsDl:           uint32(value.numTrafficTestDatagramsDl),
				        NumTrafficTestBytesDl:               uint32(value.numTrafficTestBytesDl),
				        NumTrafficTestDlDatagramsOutOfOrder: uint32(value.numTrafficTestDlDatagramsOutOfOrder),
				        NumTrafficTestDlDatagramsBad:        uint32(value.numTrafficTestDlDatagramsBad),
				        NumTrafficTestDlDatagramsMissed:     uint32(value.numTrafficTestDlDatagramsMissed),
				        TimedOut:                            bool(value.timedOut),
					}
					
				case C.DECODE_RESULT_ACTIVITY_REPORT_IND_UL_MSG:
					value := C.getActivityReportIndUlMsg(inputBuffer)
					rawData = value
					data = &ActivityReportIndUlMsg {
						Timestamp:                  time.Now(),
				        TotalTransmitMilliseconds:  uint32(value.totalTransmitMilliseconds),
				        TotalReceiveMilliseconds:   uint32(value.totalReceiveMilliseconds),
				        UpTimeSeconds:              uint32(value.upTimeSeconds),
				        TxPowerDbmPresent:          bool(value.txPowerDbmPresent),
				        TxPowerDbm:                 int8(value.txPowerDbm),
				        UlMcsPresent:               bool(value.ulMcsPresent),
				        UlMcs:                      uint8(value.ulMcs),
				        DlMcsPresent:               bool(value.dlMcsPresent),
				        DlMcs:                      uint8(value.dlMcs),
					}
					
				case C.DECODE_RESULT_ACTIVITY_REPORT_GET_CNF_UL_MSG:
					value := C.getActivityReportGetCnfUlMsg(inputBuffer)
					rawData = value
					data = &ActivityReportGetCnfUlMsg {
						Timestamp:                  time.Now(),
				        TotalTransmitMilliseconds:  uint32(value.totalTransmitMilliseconds),
				        TotalReceiveMilliseconds:   uint32(value.totalReceiveMilliseconds),
				        UpTimeSeconds:              uint32(value.upTimeSeconds),
				        TxPowerDbmPresent:          bool(value.txPowerDbmPresent),
				        TxPowerDbm:                 int8(value.txPowerDbm),
				        UlMcsPresent:               bool(value.ulMcsPresent),
				        UlMcs:                      uint8(value.ulMcs),
				        DlMcsPresent:               bool(value.dlMcsPresent),
				        DlMcs:                      uint8(value.dlMcs),
					}
					
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
					// Can't decode the uplink; this case is here to avoid a panic,
					// the code below will handle this situation fine
			}
	
			// Send any data to be recorded
			if data != nil {
				fmt.Printf("%s --> decoded data being recorded:\n\n%+v.\n", logTag, data)
				stateTableCmds <- data
			} else if rawData != nil {
				fmt.Printf("%s --> data not being recorded:\n\n%+v.\n", logTag, rawData)
			} else {
				fmt.Printf("%s --> error: undecodable message received.\n", logTag)
				//fmt.Printf("%s --> error: undecodable message received:\n\n%s.\n", logTag, spew.Sdump(inputBuffer))
			}
		}
	}

	decoderCount = 0
}


// TODO
func encodeAndEnqueueReportingInterval(mins uint32) error {
	if downlinkMessages != nil {
		// Create a buffer that is big enough to store all
		// the encoded data and take a pointer to it's first element
		var outputBuffer [C.MAX_DATAGRAM_SIZE_RAW]byte
		outputPointer := (*C.char)(unsafe.Pointer(&outputBuffer[0]))

		// Populate the data struture to be encoded
		var data C.ReportingIntervalSetReqDlMsg_t

		//		data.reportingIntervalMinutes = C.uint32_t(mins)
		var dataPointer = (*C.ReportingIntervalSetReqDlMsg_t)(unsafe.Pointer(&data))

		// Encode the data structure into the output buffer
		cbytes := C.encodeReportingIntervalSetReqDlMsg(outputPointer, dataPointer, nil, nil)

		// Send the populated output buffer bytes to NeulNet
		payload := outputBuffer[:cbytes]
		msg := AmqpMessage{
			DeviceUuid:   ueGuid,
			EndpointUuid: 4,
			//Payload:       payload,
		}
		for _, v := range payload {
			msg.Payload = append(msg.Payload, int(v))
			row.TotalBytes = row.TotalBytes + uint64(len(payload))
		}

		log.Printf("%s --> encoded a ReportingIntervalSetReqDlMsg of %d using AMQP message:\n\n%+v\n", logTag, mins, payload, msg)

		downlinkMessages <- msg
		now := time.Now()

		// Record the downlink data volume
		stateTableCmds <- &DataVolume{
			DownlinkTimestamp: &now,
			DownlinkBytes:     uint64(len(payload)),
		}

		return nil
	}

	return errors.New("No downlink message channel available to enqueue the encoded message.\n")
}

func encodeAndEnqueueIntervalGetReq(ueGuid string) error {
	if downlinkMessages != nil {

		// Create a buffer that is big enough to store all
		// the encoded data and take a pointer to it's first element
		var outputBuffer [C.MAX_DATAGRAM_SIZE_RAW]byte
		outputPointer := (*C.char)(unsafe.Pointer(&outputBuffer[0]))

		// Encode the data structure into the output buffer
		cbytes := C.encodeIntervalsGetReqDlMsg(outputPointer, nil, nil)

		// Send the populated output buffer bytes to NeulNet
		payload := outputBuffer[:cbytes]

		msg := AmqpMessage{
			DeviceUuid:   ueGuid,
			EndpointUuid: 4,
			//Payload:       payload,
		}
		for _, v := range payload {
			msg.Payload = append(msg.Payload, int(v))
		}
		log.Printf("%s --> encoded an IntervalsGetReqDlMsg using AMQP message:\n\n%+v\n", logTag, payload, msg)

		downlinkMessages <- msg
		now := time.Now()

		// Record the downlink data volume
		stateTableCmds <- &DataVolume{
			DownlinkTimestamp: &now,
			DownlinkBytes:     uint64(len(payload)),
		}

		//Record in the DisplayRow
		row.DTotalMsgs = row.DTotalMsgs + 1
		row.DTotalBytes = row.DTotalBytes + uint64(len(payload))
		row.DlastMsgReceived = &now

		return nil
	}

	return errors.New("No downlink message channel available to enqueue the encoded message.\n")
}

/* End Of File */
