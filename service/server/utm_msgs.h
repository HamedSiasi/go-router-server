/* UTM message definitions.
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

#ifndef UTM_MSGS_H
#define UTM_MSGS_H

/**
 * @file utm_msgs.h
 * This file defines the messages sent between the
 * UTM and a web server.
 *
 * The message format is a 1 byte ID followed by a message
 * of length determined by the specific message type in question.
 */
 
// ----------------------------------------------------------------
// GENERAL COMPILE-TIME CONSTANTS
// ----------------------------------------------------------------

/// The maximum length of a messages in bytes
#define MAX_MESSAGE_SIZE MAX_DATAGRAM_SIZE_RAW

// The size of the checksum (which is at the end of all messages
// except the Transparent Datagram)
#define CHECKSUM_SIZE  1

// The size of a message ID
#define MSG_ID_SIZE 1

/// The minimum length of a messages in bytes
#define MIN_MESSAGE_SIZE 1 + CHECKSUM_SIZE

/// The maximum debug string size
#define MAX_DEBUG_STRING_SIZE MAX_MESSAGE_SIZE - MSG_ID_SIZE - sizeof (uint8_t) - CHECKSUM_SIZE

/// The lower limits for reading and reporting intervals
// There is also an upper limit set by the memory capacity
// of the target (to store readings before
// the are reported) which is not known here
#define MIN_HEARTBEAT_SECONDS  15 // Set to get inside the 20 second CRNTI time-out
#define MAX_HEARTBEAT_SECONDS  3599
#define MAX_HEARTBEAT_SNAP_TO_RTC_VALUE 3599
#define MIN_REPORTING_INTERVAL 1

/// Other limits
#define MAX_BATTERY_VOLTAGE_MV 10000
#define MAX_ENERGY_MWH 0xFFFFFF
#define MIN_RSSI_RSRP ((int16_t) 0xC000) // 15 bits negative number
#define MAX_RSSI_RSRP 0
#define MAX_MEASUREMENT_CONTROL_MEASUREMENT_INTERVAL 0xFFFF // 16 bits
#define MAX_MEASUREMENT_CONTROL_REPORTING_INTERVAL 0xFFFF // 16 bits

// Flags that can indicate empty values (top bit set, not 0xFF as negative
// values can be used).
#define BYTE_NOT_PRESENT_VALUE (char) 0x80

// ----------------------------------------------------------------
// TYPES
// ----------------------------------------------------------------

/// The wake up code sent from the device.
// If you add anything here, don't forget to update the related logging strings
typedef enum WakeUpCodeTag_t
{
	WAKE_UP_CODE_OK,                          //!< A good wake-up, no problems.
	WAKE_UP_CODE_WATCHDOG,                    //!< Wake-up due to the watchdog
									          //! firing.
    WAKE_UP_CODE_NETWORK_PROBLEM,             //!< Wake-up after assert due to
                                              //! problems with the network.
    WAKE_UP_CODE_SD_CARD_PROBLEM,             //!< Wake-up after assert due to
                                              //! a problem with the SD card.
    WAKE_UP_CODE_SUPPLY_PROBLEM,              //!< Wake-up after assert due to
                                              //! a problem with the power supply.
    WAKE_UP_CODE_PROTOCOL_PROBLEM,            //!< Wake-up after assert due to
                                              //! a protocol problem.
    WAKE_UP_CODE_MODULE_NOT_RESPONDING,       //!< Wake-up after assert due to
                                              //! the modem not waking up.
    WAKE_UP_CODE_HW_PROBLEM,                  //!< Wake-up after assert due to
                                              //! problems with the HW.
	WAKE_UP_CODE_MEMORY_ALLOC_PROBLEM,        //!< Wake-up after assert due to
						       			      //! memory allocation issues.
	WAKE_UP_CODE_GENERIC_FAILURE,             //!< Wake-up after a generic failure.
	WAKE_UP_CODE_REBOOT,                      //!< Waking up after a commanded reboot.
	MAX_NUM_WAKE_UP_CODES                     //!< The maximum number of
									          //! decode results.
} WakeUpCode_t;

/// The operating modes, most be codeable into 3 bits.
// If you add anything here, don't forget to update the related logging strings
typedef enum ModeTag_t
{
    MODE_NULL,
    MODE_SELF_TEST,
    MODE_COMMISSIONING,
    MODE_STANDARD_TRX,
    MODE_TRAFFIC_TEST,
    MAX_NUM_MODES
} Mode_t;

/// The source that the time was set by.
typedef enum  TimeSetByTag_t
{
    TIME_SET_BY_NULL,
    TIME_SET_BY_GNSS,
    TIME_SET_BY_PC,
    TIME_SET_BY_WEB_API,
    MAX_NUM_TIME_SET_BY
} TimeSetBy_t;

/// The approximate amount of energy left, can be coded into 3 bits.
typedef enum EnergyLeftTag_t
{
    ENERGY_LEFT_LESS_THAN_5_PERCENT,
    ENERGY_LEFT_LESS_THAN_10_PERCENT,
    ENERGY_LEFT_MORE_THAN_10_PERCENT,
    ENERGY_LEFT_MORE_THAN_30_PERCENT,
    ENERGY_LEFT_MORE_THAN_50_PERCENT,
    ENERGY_LEFT_MORE_THAN_70_PERCENT,
    ENERGY_LEFT_MORE_THAN_90_PERCENT,
    MAX_NUM_ENERGY_LEFT
} EnergyLeft_t;

/// The approximate amount of disk space left, can be coded into 2 bits.
typedef enum DiskSpaceLeftTag_t
{
    DISK_SPACE_LEFT_LESS_THAN_1_GB,
    DISK_SPACE_LEFT_MORE_THAN_1_GB,
    DISK_SPACE_LEFT_MORE_THAN_2_GB,
    DISK_SPACE_LEFT_MORE_THAN_4_GB,
    MAX_NUM_DISK_SPACE_LEFT
} DiskSpaceLeft_t;

// ----------------------------------------------------------------
// TYPES FOR MEASUREMENTS
// ----------------------------------------------------------------

/// The measurement types.
// If you add anything here, don't forget to update the related logging strings
// IMPORTANT: if you are using the C# DLL wrapper, there is a copy of this
// enum in there. Make sure it is identical to this one.
typedef enum MeasurementTypeTag_t
{
    MEASUREMENT_TYPE_UNKNOWN,
    MEASUREMENT_GNSS,
    MEASUREMENT_CELL_ID,
    MEASUREMENT_RSRP,
    MEASUREMENT_RSSI,
    MEASUREMENT_TEMPERATURE,
    MEASUREMENT_POWER_STATE,
    MAX_NUM_MEASUREMENTS
} MeasurementType_t;

/// Structure to hold the GNSS Position.
typedef struct GnssPositionTag_t
{
    int32_t latitude;  //!< In thousandths of a minute of arc (so divide by 60,000 to get degrees)
    int32_t longitude; //!< In thousandths of a minute of arc (so divide by 60,000 to get degrees)
    int32_t elevation; //!< In metres
} GnssPosition_t;

/// Temperature -128 to +127, in Centigrade.
typedef int8_t Temperature_t;

/// RSSI, in 10ths of a dBm.
typedef int16_t Rssi_t;

/// RSRP, value as parsed from DI stream, in dBm.
typedef struct RsrpTag_t
{
    Rssi_t value;             //!< The RSRP value in 10ths of a dBm.
    bool isSyncedWithRssi;    //!< If true then the RSRP value was
                              //! taken at the same time as the RSSI
                              //! value in a report.  In this case
                              //! SNR can be computed.
} Rsrp_t;

/// CellId.
typedef uint16_t CellId_t;

/// Enum to hold the charger state (must be codeable into two bits).
// If you add anything here, don't forget to update the related logging strings
typedef enum ChargerStateTag_t
{
    CHARGER_UNKNOWN = 0x00,
    CHARGER_OFF = 0x01,
    CHARGER_ON = 0x02,
    CHARGER_FAULT = 0x03,
    MAX_NUM_CHARGER
} ChargerState_t;

/// Structure to hold the power state.
typedef struct PowerStateTag_t
{
    ChargerState_t chargerState;
    uint16_t batteryMV;  //!< battery voltage in mV, max 10,000 mV, will be coded to a resolution of 159 mV (i.e. in 5 bits).
    uint32_t energyMWh;  //!< mWh left in the battery, only up to the first 24 bits are valid (i.e. max value is 16777215).
} PowerState_t;

/// The overall measurements structure.
typedef struct MeasurementsTag_t
{
    uint32_t time;              //!< Time in UTC seconds.
    bool gnssPositionPresent;
    GnssPosition_t gnssPosition;
    bool cellIdPresent;
    CellId_t cellId;
    bool rsrpPresent;
    Rsrp_t rsrp;
    bool rssiPresent;
    Rssi_t rssi;
    bool temperaturePresent;
    Temperature_t temperature;
    bool powerStatePresent;
    PowerState_t powerState;
} Measurements_t;

/// Generic control structure for any measurement.
typedef struct MeasurementControlGenericTag_t
{
    uint32_t measurementInterval;     //!< How often, in heartbeats, to take a measurement.
                                      //! While this is represented as a 32 bit number, only
                                      //! the first 16 bits are encoded
    uint32_t maxReportingInterval;    //!< How often, in reporting intervals, a reading must
                                      //! be taken and reported, irrespective of any other
                                      //! settings here.
                                      //! While this is represented as a 32 bit number, only
                                      //! the first 16 bits are encoded
                                      //! The above two values are 16 bit rather than 32 bit
                                      //! as otherwise the max measurementsControl message,
                                      //! with all the elements in, doesn't fit into our maximum
                                      //! message size of 121 bytes
    bool useHysteresis;               //!< If true, only report if the value changes by
                                      //! +/-hysteresisValue.
    uint32_t hysteresisValue;         //!< hysteresisValue.
    bool onlyRecordIfPresent;         //!< If true, hysteresis is ignored and reports
                                      //! are sent based on the onlyRecordIfValue.
    int32_t onlyRecordIfValue;        //!< The value for onlyRecordIf.
    bool onlyRecordIfAboveNotBelow;   //!< If true, only record if above the onlyRecordIfValue,
                                      //! otherwise only record if below the onlyRecordIfValue.
    bool onlyRecordIfAtTransitionOnly;//!< If true, only record when the value crosses
                                      //! onlyRecordIfValue.
    bool onlyRecordIfIsOneShot;       //!< If true, return to normal recording mode after one
                                      //! onlyRecordIfValue.
    bool reportImmediately;           //!< If true, don't wait for the reporting period
                                      //! to expire when there's a reading that meets
                                      //! the hysteresis or onlyRecordIf triggers.
    bool recordLocalOnly;             //!< If true, do not report this measurement over
                                      //! the air, just record it locally (i.e. to SD card).
} MeasurementControlGeneric_t;

/// The measurement control structure specifically for PowerState.
typedef struct MeasurementControlPowerStateTag_t
{
    bool powerStateChargerStateReportImmediately;
    MeasurementControlGeneric_t powerStateBatteryVoltage;
    MeasurementControlGeneric_t powerStateBatteryEnergy;
} MeasurementControlPowerState_t;

/// The union of all measurement control structures.
typedef union MeasurementControlUnionTag_t
{
    MeasurementControlGeneric_t gnss;
    MeasurementControlGeneric_t cellId;
    MeasurementControlGeneric_t rsrp;
    MeasurementControlGeneric_t rssi;
    MeasurementControlGeneric_t temperature;
    MeasurementControlPowerState_t powerState;
} MeasurementControlUnion_t;

// ----------------------------------------------------------------
// MESSAGE STRUCTURES
// ----------------------------------------------------------------

/// A whole datagram, for transparent mode
typedef struct TransparentDatagramTag_t
{
     char contents[MAX_DATAGRAM_SIZE_RAW - 1];
} TransparentDatagram_t;

/// PingReqMsg_t. A ping request, which requires a respone.
// No structure for this, it's an empty message.

/// PingCnfMsg_t. The confirm sent in response to a ping request,
// No structure for this, it's an empty message.

/// InitIndUlMsg_t.  Sent at power on of the UTM, indicating that it
// has initialised.  The revision level field should will be populated
// automatically by the message codec.
typedef struct InitIndUlMsgTag_t
{
    WakeUpCode_t wakeUpCode;         //!< A wake-up code from the UTM.
    uint8_t      revisionLevel;      //!< Revision level of this messaging protocol.
    bool         sdCardNotRequired;  //!< If true the UTM will ignore SD card errors.
    bool         disableModemDebug;  //!< If true then modem debug is not being
                                     //! written to SD card, saving heartbeat time.
    bool         disableButton;      //!< If true then the button on the side of the
                                     //! UTM has been disabled.
    bool         disableServerPing;  //!< If true then the periodic "are you there"
                                     //! ping request and subsequent reboot on no
                                     //! response is disabled.
} InitIndUlMsg_t;

/// RebootReqDlMsg_t. Sent to reboot the UTM and set various flags.
typedef struct RebootReqDlMsgTag_t
{
    bool sdCardNotRequired; //!< If true SD card errors are ignored,
                            //! the UTM can run without one.
    bool disableModemDebug; //!< By default the UTM writes all the modem
                            //! debug it can to SD card.  This may take
                            //! time and disk space.  If this flag is set
                            //! to true then the capture and storage of
                            //! the modem debug is disabled.
    bool disableButton;     //!< Disable the button on the side of the UTM.
    bool disableServerPing; //!< If true then the periodic "are you there"
                            //! ping request and subsequent reboot on no
                            //! response is disabled.
} RebootReqDlMsg_t;

/// DateTimeSetReqDlMsg_t. Sent to set the date/time on the UTM.
typedef struct DateTimeSetReqDlMsg_t
{
    uint32_t time;    //!< The time in UTC seconds.
    bool setDateOnly; //!< If true then the time value should be used
                      //< to set the date only and not the time.
                      //< This can be used if the time has already been
                      //< set by GNSS (in which case the date is unknown
                      //< as GNSS time does not include the date).
} DateTimeSetReqDlMsg_t;

/// DateTimeSetCnfUlMsg_t. The date/time on the UTM. Sent
// in response to a DateTimeSetReqDlMsg_t.
typedef struct DateTimeSetCnfUlMsg_t
{
    uint32_t time;     //!< The time in UTC seconds.
    TimeSetBy_t setBy; //!< The source the time was set by
} DateTimeSetCnfUlMsg_t;

/// DateTimeGetReqDlMsg_t. Get the date/time from the UTM.
// No structure for this, it's an empty message.

/// DateTimeGetCnfUlMsg_t. The date/time on the UTM. Sent
// in response to a DateTimeGetReqDlMsg_t.
typedef struct DateTimeGetCnfUlMsg_t
{
    uint32_t time;     //!< The time in UTC seconds.
    TimeSetBy_t setBy; //!< The source the time was set by
} DateTimeGetCnfUlMsg_t;

/// DateTimeIndUlMsg_t. The date/time on the UTM. Sent
// after power on to assure synchronisation.
typedef struct DateTimeIndUlMsg_t
{
    uint32_t time;     //!< The time in UTC seconds.
    TimeSetBy_t setBy; //!< The source the time was set by
} DateTimeIndUlMsg_t;

/// ModeSetReqDlMsg_t. Sent to set the mode that the device is operating in.
// NOTE: if switching to Traffic Test mode, ensure that a TestModeParametersSetReqDlMsg
// has preceded this message.
typedef struct ModeSetReqDlMsgTag_t
{
    Mode_t mode; //!< The mode to use.
} ModeSetReqDlMsg_t;

/// ModeSetCnfUlMsg_t. The mode that the device is operating in. Sent
// in response to a ModeSetReqDlMsg_t.
typedef struct ModeSetCnfUlMsgTag_t
{
    Mode_t mode; //!< The mode in use.
} ModeSetCnfUlMsg_t;

// ModeGetReqDlMsg_t. Sent to get the mode that the device is operating in.
// No structure for this message as it has no contents.

/// ModeGetCnfUlMsg_t. The mode that the device is operating in. Sent
// in response to a ModeGetReqDlMsg_t.
typedef struct ModeGetCnfUlMsgTag_t
{
    Mode_t mode; //!< The mode in use.
} ModeGetCnfUlMsg_t;

/// IntervalsGetReqDlMsg_t. Get the intervals at which the UTM reads and
// reports measurements.
// No structure for this, it's an empty message.

/// IntervalsGetCnfUlMsg_t. The intervals at which the UTM reads and
// reports measurements.  Sent in response to IntervalGetReqDlMsg_t.
typedef struct IntervalsGetCnfUlMsgTag_t
{
  uint32_t reportingInterval; //!< The interval at which the UTM
                              //! sends reports, in heartbeats.
  uint32_t heartbeatSeconds;  //!< The interval at which the UTM
                              //! wakes up to make measurements.
  bool heartbeatSnapToRtc;    //!< If true, the heartbeatSeconds represents
                              //! the number of seconds into an hour at
                              //! which the UTC should perform a heartbeat
                              //! action, rather than the delay in seconds
                              //! from the previous heartbeat.
} IntervalsGetCnfUlMsg_t;

/// ReportingIntervalSetReqDlMsg_t. Set the interval at which the UTM sends
// reports.
typedef struct ReportingIntervalSetReqDlMsgTag_t
{
  uint32_t reportingInterval; //!< The interval at which the UTM
                              //! should send reports, in heartbeats.
} ReportingIntervalSetReqDlMsg_t;

/// ReportingIntervalSetCnfUlMsg_t. The interval at which the UTM sends
// reports. Sent in response to ReportingIntervalSetReqDlMsg_t.
typedef struct ReportingIntervalSetCnfUlMsgTag_t
{
  uint32_t reportingInterval; //!< The interval at which the UTM
                              //! sends reports, in heartbeats.
} ReportingIntervalSetCnfUlMsg_t;

/// HeartbeatSetReqDlMsg_t. Set the interval at which the UTM makes
// measurements.
typedef struct HeartbeatSetReqDlMsgTag_t
{
  uint32_t heartbeatSeconds;  //!< The interval at which the UTM
                              //! should wake up to take measurements.
  bool heartbeatSnapToRtc;    //!< If true, the heartbeatSeconds represents
                              //! the number of seconds into an hour at
                              //! which the UTC should perform a heartbeat
                              //! action, rather than the delay in seconds
                              //! from the previous heartbeat.
} HeartbeatSetReqDlMsg_t;

/// HeartbeatSetCnfUlMsg_t. The interval at which the UTM wakes-up to
// readings. Sent in response to HeartbeatSetReqDlMsg_t.
typedef struct HeartbeatSetCnfUlMsgTag_t
{
  uint32_t heartbeatSeconds;  //!< The interval at which the UTM
                              //! wakes-up to make measurements.
  bool heartbeatSnapToRtc;    //!< If true, the heartbeatSeconds represents
                              //! the number of seconds into an hour at
                              //! which the UTC should perform a heartbeat
                              //! action, rather than the delay in seconds
                              //! from the previous heartbeat.
} HeartbeatSetCnfUlMsg_t;

/// PollIndUlMsg_t.  Sent at the start of every reporting period (in case there are
// no measurements to report).
typedef struct PollIndUlMsgTag_t
{
    Mode_t mode; //!< The current operating mode.
    EnergyLeft_t energyLeft; //!< Approximate indication of how much energy is left in the battery.
    DiskSpaceLeft_t diskSpaceLeft; //!< Approximate indication of how much disk space is left.
} PollIndUlMsg_t;

/// MeasurementsReportIndUlMsg_t.  A set of measurements, sent either periodically
// or as a result of some local trigger on the UTM.
typedef struct MeasurementsIndUlMsgTag_t
{
    Measurements_t measurements; //!< All the measurements.
} MeasurementsIndUlMsg_t;

/// MeasurementsReportReqDlMsg_t.  Request a set of measurements.
// No structure for this, it's an empty message.

/// MeasurementsGetCnfUlMsg_t.  A set of measurements, sent in response to a
// MeasurementsGetReqDlMsg_t.
typedef struct MeasurementsGetCnfUlMsgTag_t
{
    Measurements_t measurements; //!< All the measurements.
} MeasurementsGetCnfUlMsg_t;

/// MeasurementControlSetReqDlMsg_t.  Set the measurement control settings for
// a given measurement type.
typedef struct MeasurementControlSetReqDlMsgTag_t
{
    MeasurementType_t measurementType;            //!< The measurement type to set the controls for.
    MeasurementControlUnion_t measurementControl; //!< The settings.
} MeasurementControlSetReqDlMsg_t;

/// MeasurementControlSetCnfUlMsg_t.  Sent in response to a MeasurementControlSetReqDlMsg_t.
typedef struct MeasurementControlSetCnfUlMsgTag_t
{
    MeasurementType_t measurementType;            //!< The measurement type that the controls were
                                                  //! set for.
    MeasurementControlUnion_t measurementControl; //!< The settings.
} MeasurementControlSetCnfUlMsg_t;

/// MeasurementsControlGetReqDlMsg_t.  Get the control settings for all measurements.
// No structure for this, it's an empty message.

/// MeasurementControlGetCnfUlMsg_t.  The measurement control settings for all measurements.
// Sent in response to a MeasurementsControlGetReqDlMsg_t.
typedef struct MeasurementsControlGetCnfUlMsgTag_t
{
    MeasurementControlGeneric_t gnss;            //!< GNSS measurement control settings (values applies to lat/long only).
    MeasurementControlGeneric_t cellId;          //!< Cell ID measurement control settings.
    MeasurementControlGeneric_t rsrp;            //!< RSRP measurement control settings.
    MeasurementControlGeneric_t rssi;            //!< RSSI measurement control settings.
    MeasurementControlGeneric_t temperature;     //!< Temperature measurement control settings.
    MeasurementControlPowerState_t powerState;   //!< Power state measurement control settings.
} MeasurementsControlGetCnfUlMsg_t;

/// MeasurementControlIndUlMsg_t.  The measurement control settings for all measurements.
typedef struct MeasurementsControlIndUlMsgTag_t
{
    MeasurementControlGeneric_t gnss;            //!< GNSS measurement control settings (values applies to lat/long only).
    MeasurementControlGeneric_t cellId;          //!< Cell ID measurement control settings.
    MeasurementControlGeneric_t rsrp;            //!< RSRP measurement control settings.
    MeasurementControlGeneric_t rssi;            //!< RSSI measurement control settings.
    MeasurementControlGeneric_t temperature;     //!< Temperature measurement control settings.
    MeasurementControlPowerState_t powerState;   //!< Power state measurement control settings.
} MeasurementsControlIndUlMsg_t;

/// MeasurementsControlDefaultsSetReqDlMsg_t.  Set all measurement control settings to defaults.
// No structure for this, it's an empty message.

/// MeasurementsControlDefaultsSetCnfUlMsg_t.  Response to a MeasurementsControlDefaultsSetReqDlMsg_t.
// No structure for this, it's an empty message.

/// TrafficReportIndUlMsg_t.  A report of the traffic data that has occurred
// since the last InitIndUlMsg_t.
typedef struct TrafficReportIndUlMsgTag_t
{
    uint32_t numDatagramsUl;
    uint32_t numBytesUl;
    uint32_t numDatagramsDl;
    uint32_t numBytesDl;
    uint32_t numDatagramsDlBadChecksum;
} TrafficReportIndUlMsg_t;

/// TrafficReportGetReqDlMsg_t.  Request a traffic report.
// No structure for this, it's an empty message.

/// TrafficReportGetCnfUlMsg_t.  A report of the traffic data that has occurred
// since InitIndUlMsg_t, sent in response to a TrafficReportGetReqDlMsg_t.
typedef struct TrafficReportGetCnfUlMsgTag_t
{
    uint32_t numDatagramsUl;
    uint32_t numBytesUl;
    uint32_t numDatagramsDl;
    uint32_t numBytesDl;
    uint32_t numDatagramsDlBadChecksum;
} TrafficReportGetCnfUlMsg_t;

/// DebugIndUlMsg_t.  A generic message containing a debug string.
typedef struct DebugIndUlMsgTag_t
{
    uint8_t sizeOfString;               //!< String size in bytes
    char string[MAX_DEBUG_STRING_SIZE]; //!< The string (not NULL terminated).
} DebugIndUlMsg_t;

/// TrafficTestModeParametersSetReqDlMsg_t. Sent to set parameters for
// test mode.
typedef struct TrafficTestModeParametersSetReqDlMsgTag_t
{
    uint32_t numUlDatagrams; //!< The number of datagrams to send from the UTM.
    uint32_t lenUlDatagram; //!< The size of each uplink datagram.
    uint32_t numDlDatagrams; //!< The number of datagrams to expect from the network.
    uint32_t lenDlDatagram; //!< The size of each downlink datagram.
    uint32_t timeoutSeconds; //!< A guard timer for the test.  Zero means no timeout.
    bool noReportsDuringTest; //!< If true, no TrafficTestModeReportInds as sent during the test.
} TrafficTestModeParametersSetReqDlMsg_t;

/// TrafficTestModeParametersSetCnfUlMsg_t. Confirmation of the traffic test settings. Sent
// in response to a TestModeParametersSetReqDlMsg_t.
typedef struct TrafficTestModeParametersSetCnfUlMsgTag_t
{
    uint32_t numUlDatagrams; //!< The number of datagrams to send from the UTM.
    uint32_t lenUlDatagram; //!< The size of each uplink datagram.
    uint32_t numDlDatagrams; //!< The number of datagrams to expect from the network.
    uint32_t lenDlDatagram; //!< The size of each downnlink datagram.
    uint32_t timeoutSeconds; //!< The guard timer for the test.  Zero means no timeout.
    bool noReportsDuringTest; //!< If true, no TrafficTestModeReportInds as sent during the test.
} TrafficTestModeParametersSetCnfUlMsg_t;

/// TrafficTestModeParametersGetReqDlMsg_t. Sent to get the traffic test settings
// the device was last sent.
// No structure for this message as it has no contents.

/// TrafficTestModeParametersGetCnfUlMsg_t. The traffic test settings. Sent
// in response to a TestModeParametersGetReqDlMsg_t.
typedef struct TrafficTestModeParametersGetCnfUlMsgTag_t
{
    uint32_t numUlDatagrams; //!< The number of datagrams to send from the UTM.
    uint32_t lenUlDatagram; //!< The size of each uplink datagram.
    uint32_t numDlDatagrams; //!< The number of datagrams to expect from the network.
    uint32_t lenDlDatagram; //!< The size of each downnlink datagram.
    uint32_t timeoutSeconds; //!< The guard timer for the test.  Zero means no timeout.
    bool noReportsDuringTest; //!< If true, no TrafficTestModeReportInds as sent during the test.
} TrafficTestModeParametersGetCnfUlMsg_t;

/// The spec for the Traffic Test mode rule breaker datagram
typedef struct TrafficTestModeRuleBreakerDatagramTag_t
{
    char fill;
    uint32_t length;
} TrafficTestModeRuleBreakerDatagram_t;

/// TrafficTestModeReportIndUlMsg_t.  A report of the traffic data
// that has been sent and received since the start of Traffic Test
// mode.
typedef struct TrafficTestModeReportIndUlMsgTag_t
{
    uint32_t numTrafficTestDatagramsUl;
    uint32_t numTrafficTestBytesUl;
    uint32_t numTrafficTestDatagramsDl;
    uint32_t numTrafficTestBytesDl;
    uint32_t numTrafficTestDlDatagramsOutOfOrder;
    uint32_t numTrafficTestDlDatagramsBad;
    uint32_t numTrafficTestDlDatagramsMissed;
    bool timedOut;
} TrafficTestModeReportIndUlMsg_t;

/// TrafficTestModeReportGetReqDlMsg_t.  Request a test mode traffic report.
// No structure for this, it's an empty message.

/// TrafficTestModeReportGetCnfUlMsg_t.  A report of the Traffic Test mode
// traffic data that has occurred entering Traffic Test mode.
// Sent in response to a TrafficTestModeReportGetReqDlMsg_t.
typedef struct TrafficTestModeReportGetCnfUlMsgTag_t
{
    uint32_t numTrafficTestDatagramsUl;
    uint32_t numTrafficTestBytesUl;
    uint32_t numTrafficTestDatagramsDl;
    uint32_t numTrafficTestBytesDl;
    uint32_t numTrafficTestDlDatagramsOutOfOrder;
    uint32_t numTrafficTestDlDatagramsBad;
    uint32_t numTrafficTestDlDatagramsMissed;
    bool timedOut;
} TrafficTestModeReportGetCnfUlMsg_t;

/// ActivityReportIndUlMsg_t.  A report of the activity of the
// module since the InitInd.
typedef struct ActivityReportIndUlMsgTag_t
{
    uint32_t totalTransmitMilliseconds;
    uint32_t totalReceiveMilliseconds;
    uint32_t upTimeSeconds;
    bool txPowerDbmPresent;
    int8_t txPowerDbm;  // Note that this is signed
    bool ulMcsPresent;
    uint8_t ulMcs;
    bool dlMcsPresent;
    uint8_t dlMcs;
} ActivityReportIndUlMsg_t;

/// ActivityReportGetReqDlMsg_t.  Request an activity report.
// No structure for this, it's an empty message.

/// ActivityReportGetCnfUlMsg_t.  A report of the activity of the
// module since the InitInd.
// Sent in response to a ActivityReportGetReqDlMsg_t.
typedef struct ActivityReportGetCnfUlMsgTag_t
{
    uint32_t totalTransmitMilliseconds;
    uint32_t totalReceiveMilliseconds;
    uint32_t upTimeSeconds;
    bool txPowerDbmPresent;
    int8_t txPowerDbm;  // Note that this is signed
    bool ulMcsPresent;
    uint8_t ulMcs;
    bool dlMcsPresent;
    uint8_t dlMcs;
} ActivityReportGetCnfUlMsg_t;

// ----------------------------------------------------------------
// MESSAGE UNIONS
// ----------------------------------------------------------------

/// Union of all downlink messages.
typedef union DlMsgUnionTag_t
{
    TransparentDatagram_t transparentDatagram;
    RebootReqDlMsg_t rebootReqDlMsg;
    DateTimeSetReqDlMsg_t dateTimeSetReqDlMsg;
    ModeSetReqDlMsg_t modeSetReqDlMsg;
    ReportingIntervalSetReqDlMsg_t reportingIntervalSetReqDlMsg;
    HeartbeatSetReqDlMsg_t heartbeatSetReqDlMsg;
    MeasurementControlSetReqDlMsg_t measurementControlSetReqDlMsg;
    TrafficTestModeParametersSetReqDlMsg_t trafficTestModeParametersSetReqDlMsg;
    TrafficTestModeRuleBreakerDatagram_t trafficTestModeRuleBreakerDatagram;
    ActivityReportIndUlMsg_t activityReportIndUlMsg;
} DlMsgUnion_t;

/// Union of all uplink messages.
typedef union UlMsgUnionTag_t
{
    TransparentDatagram_t transparentDatagram;
    InitIndUlMsg_t initIndUlMsg;
    DateTimeSetCnfUlMsg_t dateTimeSetCnfUlMsg;
    DateTimeGetCnfUlMsg_t dateTimeGetCnfUlMsg;
    DateTimeIndUlMsg_t dateTimeIndUlMsg;
    ModeSetCnfUlMsg_t modeSetCnfUlMsg;
    ModeGetCnfUlMsg_t modeGetCnfUlMsg;
    IntervalsGetCnfUlMsg_t intervalsGetCnfUlMsg;
    ReportingIntervalSetCnfUlMsg_t reportingIntervalSetCnfUlMsg;
    HeartbeatSetCnfUlMsg_t heartbeatSetCnfUlMsg;
    PollIndUlMsg_t pollIndUlMsg;
    MeasurementsGetCnfUlMsg_t measurementsGetCnfUlMsg;
    MeasurementsIndUlMsg_t measurementsIndUlMsg;
    MeasurementControlSetCnfUlMsg_t measurementControlSetCnfUlMsg;
    MeasurementsControlGetCnfUlMsg_t measurementsControlGetCnfUlMsg;
    MeasurementsControlIndUlMsg_t measurementsControlIndUlMsg;
    TrafficReportIndUlMsg_t trafficReportIndUlMsg;
    TrafficReportGetCnfUlMsg_t trafficReportGetCnfUlMsg;
    TrafficTestModeParametersSetCnfUlMsg_t trafficTestModeParametersSetCnfUlMsg;
    TrafficTestModeParametersGetCnfUlMsg_t trafficTestModeParametersGetCnfUlMsg;
    TrafficTestModeRuleBreakerDatagram_t trafficTestModeRuleBreakerDatagram;
    TrafficTestModeReportIndUlMsg_t trafficTestModeReportIndUlMsg;
    TrafficTestModeReportGetCnfUlMsg_t trafficTestModeReportGetCnfUlMsg;
    ActivityReportIndUlMsg_t activityReportIndUlMsg;
    ActivityReportGetCnfUlMsg_t activityReportGetCnfUlMsg;
    DebugIndUlMsg_t debugIndUlMsg;
} UlMsgUnion_t;

#endif

// End Of File
