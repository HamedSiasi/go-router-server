/* UTM interface definition
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

#ifndef UTM_API_H
#define UTM_API_H

/**
 * @file utm_api.h
 * This file defines the API to the UTM
 * NOTE: when making any changes here, please
 * also update the logging values in utm_msg_codec.c.
 */

#ifdef __cplusplus
extern "C" {
#endif

/// The maximum length of a raw datagram in bytes.
// This is defined up here because utm_msgs.h
// needs it.
#define MAX_DATAGRAM_SIZE_RAW 121

#include <utm_msgs.h>

// ----------------------------------------------------------------
// GENERAL COMPILE-TIME CONSTANTS
// ----------------------------------------------------------------

/// The revision level of this API.
// 0: No revision yet, under development and still changing constantly.
// 1: First release version.
// 2: Added "noReportsDuringTest" to the TrafficTestModeParametersSetReq/SetCnf/GetCnf messages.
//    Added "txPowerDbm", "ulMcs" and "dlMcs" (plus present flags) to ActivityIndUlMsg and ActivityGetCnfUlMsg.
// 3: txPowerDbm becomes an int_t value instead of a uint16_t value.
#define REVISION_LEVEL 3

/// How often a report is sent to the network (in heartbeats).
#define DEFAULT_REPORTING_INTERVAL 1

/// How often the UTM wakes-up to take measurements.
#define DEFAULT_HEARTBEAT_SECONDS   900

/// The maximum size of the traffic test rule breaker datagram.
#define TRAFFIC_TEST_MODE_RULE_BREAKER_MAX_LENGTH 100

/// Transparent datagram ID.  If this is present at the start
// of a datagram then the whole datagram is encoded/decoded
// as a flat array by the message codec.
// THIS MUST BE ZERO.
#define TRANSPARENT_DATAGRAM_ID 0

/// Ping request/confirm ID.
#define PING_REQ_MSG_ID 1
#define PING_CNF_MSG_ID 2

// ----------------------------------------------------------------
// PUBLIC FUNCTIONS
// ----------------------------------------------------------------

#if TRANSPARENT_DATAGRAM_ID != 0
#error TRANSPARENT_DATAGRAM_ID must be zero.
#endif

/// The message IDs in the downlink direction (i.e. to the UTM)
typedef enum MsgIdDlTag_t
{
    TRANSPARENT_DL_DATAGRAM = TRANSPARENT_DATAGRAM_ID, //!< Transparent to this protocol,
                                                       //! can only appear at the start of
                                                       //! a datagram.  The whole datagram
                                                       //! is passed on transparently.
    PING_REQ_DL_MSG = PING_REQ_MSG_ID,                 //!< A ping request.
    PING_CNF_DL_MSG = PING_CNF_MSG_ID,                 //!< A ping confirm.
    REBOOT_REQ_DL_MSG,                                 //!< Reboot the UTM.
    DATE_TIME_SET_REQ_DL_MSG,                          //!< Set the date/time on the UTM.
    DATE_TIME_GET_REQ_DL_MSG,                          //!< Get the date/time from the UTM.
    MODE_SET_REQ_DL_MSG,                               //!< Set the operating mode of the UTM.
    MODE_GET_REQ_DL_MSG,                               //!< Get the operating mode of the UTM.
    HEARTBEAT_SET_REQ_DL_MSG,                          //!< Set the rate at which wake-ups to take
                                                       //! to take measurements are performed by
                                                       //! the UTM.
    REPORTING_INTERVAL_SET_REQ_DL_MSG,                 //!< Set the rate at which reports
                                                       //! are returned by the UTM.
    INTERVALS_GET_REQ_DL_MSG,                          //!< Get the rate at which wake-ups for
                                                       //! measurements and reports are made.
    MEASUREMENTS_GET_REQ_DL_MSG,                       //!< Get a set of measurements from the UTM.
    MEASUREMENT_CONTROL_SET_REQ_DL_MSG,                //!< Set the control settings for a given measurement.
    MEASUREMENTS_CONTROL_GET_REQ_DL_MSG,               //!< Get the settings for all measurements.
    MEASUREMENTS_CONTROL_DEFAULTS_SET_REQ_DL_MSG,      //!< Set the measurements controls for all
                                                       //! measurements to defaults.
    TRAFFIC_REPORT_GET_REQ_DL_MSG,                     //!< Get a traffic report from the UTM.
    TRAFFIC_TEST_MODE_PARAMETERS_SET_REQ_DL_MSG,       //!< Set the test mode settings.
    TRAFFIC_TEST_MODE_PARAMETERS_GET_REQ_DL_MSG,       //!< Get the test mode settings.
    TRAFFIC_TEST_MODE_RULE_BREAKER_DL_DATAGRAM,        //!< This indicates the presence of a
                                                       //!  whole rule-breaker datagram with special handling at
                                                       //!  both ends
    TRAFFIC_TEST_MODE_REPORT_GET_REQ_DL_MSG,           //!< Get a Traffic Test Mode report from the UTM.
    ACTIVITY_REPORT_GET_REQ_DL_MSG,                    //!< Get an Activity report from the UTM.
    MAX_NUM_DL_MSGS,                                   //!< The maximum number of downlink.
                                                       //! messages.
    MIN_TRAFFIC_TEST_MODE_FILL_DL = MAX_NUM_DL_MSGS    //!< The minimum value that a traffic test mode fill
                                                       // value can take, to ensure no overlap with valid
                                                       // message IDs in case of error
} MsgIdDl_t;

/// The message IDs in the uplink direction (i.e. from the UTM)
typedef enum MsgIdUlTag_t
{
    TRANSPARENT_UL_DATAGRAM = TRANSPARENT_DATAGRAM_ID, //!< Transparent to this protocol,
                                                       //! can only appear at the start of
                                                       //! a datagram.  The whole datagram
                                                       //! is passed on transparently.
    PING_REQ_UL_MSG = PING_REQ_MSG_ID,                 //!< A ping request.
    PING_CNF_UL_MSG = PING_CNF_MSG_ID,                 //!< A ping confirm.
    INIT_IND_UL_MSG,                                   //!< Power on of the UTM has completed.
    POLL_IND_UL_MSG,                                   //!< A poll message, sent on expiry of a
                                                       //! reporting period.
    DATE_TIME_IND_UL_MSG,                              //!< The date/time of the UTM.
    DATE_TIME_SET_CNF_UL_MSG,                          //!< Response to a date/time set request.
    DATE_TIME_GET_CNF_UL_MSG,                          //!< Response to a date/time get request.
    MODE_SET_CNF_UL_MSG,                               //!< Response to a mode set request.
    MODE_GET_CNF_UL_MSG,                               //!< Response to a mode get request.
    HEARTBEAT_SET_CNF_UL_MSG,                          //!< The rate at which sensor readings are
                                                       //! taken by the UTM.
    REPORTING_INTERVAL_SET_CNF_UL_MSG,                 //!< The rate at which sensor readings are
                                                       //! reported by the UTM.
    INTERVALS_GET_CNF_UL_MSG,                          //!< The rate at which sensor
                                                       //! readings and reports are made.
    MEASUREMENTS_IND_UL_MSG,                           //!< A periodic measurements report.
    MEASUREMENTS_GET_CNF_UL_MSG,                       //!< Response to a measurements request.
    MEASUREMENT_CONTROL_SET_CNF_UL_MSG,                //!< Response to a measurement control set request.
    MEASUREMENTS_CONTROL_IND_UL_MSG,                   //!< The control settings for all measurements.
    MEASUREMENTS_CONTROL_GET_CNF_UL_MSG,               //!< Response to a measurementS control get request.
    MEASUREMENTS_CONTROL_DEFAULTS_SET_CNF_UL_MSG,      //!< Response to measurements controls defaults set request.
    TRAFFIC_REPORT_IND_UL_MSG,                         //!< A periodic traffic report.
    TRAFFIC_REPORT_GET_CNF_UL_MSG,                     //!< Response to a traffic report request.
    TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF_UL_MSG,       //!< Response to test mode settings set request.
    TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF_UL_MSG,       //!< Response to test mode settings get request.
    TRAFFIC_TEST_MODE_RULE_BREAKER_UL_DATAGRAM,        //!< NOT A MESSAGE!  This indicates the presence of a
                                                       //!  whole rule-breaker datagram with special handling at
                                                       //!  both ends
    TRAFFIC_TEST_MODE_REPORT_IND_UL_MSG,               //!< The traffic report sent in Traffic Test mode only.
    TRAFFIC_TEST_MODE_REPORT_GET_CNF_UL_MSG,           //!< The traffic report sent in response to a Traffic Test mode report request.
    ACTIVITY_REPORT_IND_UL_MSG,                        //!< An Activity report from the UTM.
    ACTIVITY_REPORT_GET_CNF_UL_MSG,                    //!< The Activity report  sent in response to an Activity report request
    DEBUG_IND_UL_MSG,                                  //!< A debug string.
    MAX_NUM_UL_MSGS,                                   //!< The maximum number of uplink messages.
    MIN_TRAFFIC_TEST_MODE_FILL_UL = MAX_NUM_UL_MSGS    //!< The minimum value that a traffic test mode fill
                                                       // value can take, to ensure no overlap with valid
                                                       // message IDs in case of error
} MsgIdUl_t;

// ----------------------------------------------------------------
// MESSAGE ENCODING FUNCTIONS
// ----------------------------------------------------------------

/// Encode a transparent message.  This occupies an entire
// datagram, rather than the usual shorter message length,
// and the contents are a flat array.  Note that it does NOT
// include a checksum, any validation of the contents is
// up to the application.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_DATAGRAM_SIZE_RAW long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeTransparentDatagram (char * pBuffer,
                                    TransparentDatagram_t * pDatagram,
                                    char ** ppLog,
                                    uint32_t * pLogSize);

/// Encode a ping request message.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param ppLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodePingReqMsg (char * pBuffer,
                           char ** ppLog,
                           uint32_t * pLogSize);

/// Encode a ping confirm message, sent in response to a ping request.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodePingCnfMsg (char * pBuffer,
                           char ** ppLog,
                           uint32_t * pLogSize);

/// Encode an uplink message that is sent at power-on of the
// UTM.  Indicates that the device has been initialised.  After
// transmission of this message measurements  will be taken
// and reported at the indicated rates.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeInitIndUlMsg (char * pBuffer,
                             InitIndUlMsg_t * pMsg,
                             char ** ppLog,
                             uint32_t * pLogSize);

/// Encode an uplink message that is sent at the expiry of a
// reporting period, or whenever the status of the UTM
// needs to be conveyed.
// \param pBuffer  A pointer to the buffer to encode into.  The
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// buffer length must be at least MAX_MESSAGE_SIZE long
uint32_t encodePollIndUlMsg (char * pBuffer,
                             PollIndUlMsg_t * pMsg,
                             char ** ppLog,
                             uint32_t * pLogSize);

/// Encode a downlink message that reboots the device.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeRebootReqDlMsg (char * pBuffer,
                               RebootReqDlMsg_t *pMsg,
                               char ** ppLog,
                               uint32_t * pLogSize);

/// Encode an uplink message that indicates the date/time on
// the UTM.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeDateTimeIndUlMsg (char * pBuffer,
                                 DateTimeIndUlMsg_t * pMsg,
                                 char ** ppLog,
                                 uint32_t * pLogSize);

/// Encode a downlink message that sets the date/time on
// the UTM.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeDateTimeSetReqDlMsg (char * pBuffer,
                                    DateTimeSetReqDlMsg_t * pMsg,
                                    char ** ppLog,
                                    uint32_t * pLogSize);

/// Encode an uplink message that is sent in response to a
// DateTimeSetReqDlMsg.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeDateTimeSetCnfUlMsg (char * pBuffer,
                                    DateTimeSetCnfUlMsg_t * pMsg,
                                    char ** ppLog,
                                    uint32_t * pLogSize);

/// Encode a downlink message that gets the date/time from the
// UTM.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeDateTimeGetReqDlMsg (char * pBuffer,
                                    char ** ppLog,
                                    uint32_t * pLogSize);

/// Encode an uplink message that is sent in response to a
// DateTimeGetReqDlMsg.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeDateTimeGetCnfUlMsg (char * pBuffer,
                                    DateTimeGetCnfUlMsg_t * pMsg,
                                    char ** ppLog,
                                    uint32_t * pLogSize);

/// Encode a downlink message that sets the mode the device
// is operating in.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeModeSetReqDlMsg (char * pBuffer,
                                ModeSetReqDlMsg_t * pMsg,
                                char ** ppLog,
                                uint32_t * pLogSize);

/// Encode an uplink message that is sent in response to a
// ModeSetReqDlMsg.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeModeSetCnfUlMsg (char * pBuffer,
                                ModeSetCnfUlMsg_t * pMsg,
                                char ** ppLog,
                                uint32_t * pLogSize);

/// Encode a downlink message that gets the mode the device
// is operating in.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeModeGetReqDlMsg (char * pBuffer,
                                char ** ppLog,
                                uint32_t * pLogSize);

/// Encode an uplink message that is sent in response to a
// ModeGetReqDlMsg.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeModeGetCnfUlMsg (char * pBuffer,
                                ModeGetCnfUlMsg_t * pMsg,
                                char ** ppLog,
                                uint32_t * pLogSize);

/// Encode a downlink message that retrieves the reading and
// reporting intervals.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeIntervalsGetReqDlMsg (char * pBuffer,
                                     char ** ppLog,
                                     uint32_t * pLogSize);

/// Encode an uplink message that is sent as a response to a
// IntervalsGetReqDlMsg.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeIntervalsGetCnfUlMsg (char * pBuffer,
                                     IntervalsGetCnfUlMsg_t * pMsg,
                                     char ** ppLog,
                                     uint32_t * pLogSize);

/// Encode a downlink message that sets the reporting interval.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeReportingIntervalSetReqDlMsg (char * pBuffer,
                                             ReportingIntervalSetReqDlMsg_t * pMsg,
                                             char ** ppLog,
                                             uint32_t * pLogSize);

/// Encode an uplink message that is sent as a response to a
// ReportingIntervalSetReqDlMsg.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeReportingIntervalSetCnfUlMsg (char * pBuffer,
                                             ReportingIntervalSetCnfUlMsg_t * pMsg,
                                             char ** ppLog,
                                             uint32_t * pLogSize);

/// Encode a downlink message that sets the heartbeat.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeHeartbeatSetReqDlMsg (char * pBuffer,
                                     HeartbeatSetReqDlMsg_t * pMsg,
                                     char ** ppLog,
                                     uint32_t * pLogSize);

/// Encode an uplink message that is sent as a response to a
// HeartbeatSetReqDlMsg.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeHeartbeatSetCnfUlMsg (char * pBuffer,
                                     HeartbeatSetCnfUlMsg_t * pMsg,
                                     char ** ppLog,
                                     uint32_t * pLogSize);

/// Encode an uplink message containing the measurements.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeMeasurementsIndUlMsg (char * pBuffer,
                                     MeasurementsIndUlMsg_t * pMsg,
                                     char ** ppLog,
                                     uint32_t * pLogSize);

/// Encode a downlink message that retrieves measurements.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeMeasurementsGetReqDlMsg (char * pBuffer,
                                        char ** ppLog,
                                        uint32_t * pLogSize);

/// Encode an uplink message that is sent as a response to a
// MeasurementsGetReqDlMsg.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeMeasurementsGetCnfUlMsg (char * pBuffer,
                                        MeasurementsGetCnfUlMsg_t * pMsg,
                                        char ** ppLog,
                                        uint32_t * pLogSize);

/// Encode a downlink message to set controls for a given measurement.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeMeasurementControlSetReqDlMsg (char * pBuffer,
                                              MeasurementControlSetReqDlMsg_t * pMsg,
                                              char ** ppLog,
                                              uint32_t * pLogSize);

/// Encode an uplink message that is sent as a response to a
// MeasurementControlSetReqDlMsg.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeMeasurementControlSetCnfUlMsg (char * pBuffer,
                                              MeasurementControlSetCnfUlMsg_t * pMsg,
                                              char ** ppLog,
                                              uint32_t * pLogSize);

/// Encode an uplink message that gives the measurements control settings.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeMeasurementsControlIndUlMsg (char * pBuffer,
                                            MeasurementsControlIndUlMsg_t * pMsg,
                                            char ** ppLog,
                                            uint32_t * pLogSize);

/// Encode a downlink message to get the measurements control settings.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeMeasurementsControlGetReqDlMsg (char * pBuffer,
                                               char ** ppLog,
                                               uint32_t * pLogSize);

/// Encode an uplink message that is sent as a response to a
// MeasurementsControlGetReqDlMsg.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeMeasurementsControlGetCnfUlMsg (char * pBuffer,
                                               MeasurementsControlGetCnfUlMsg_t * pMsg,
                                               char ** ppLog,
                                               uint32_t * pLogSize);

/// Encode a downlink message to put the measurements control settings
// back to defaults.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeMeasurementsControlDefaultsSetReqDlMsg (char * pBuffer,
                                                       char ** ppLog,
                                                       uint32_t * pLogSize);

/// Encode an uplink message that is sent as a response to a
// MeasurementsControlDefaultsSetReqDlMsg.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeMeasurementsControlDefaultsSetCnfUlMsg (char * pBuffer,
                                                       char ** ppLog,
                                                       uint32_t * pLogSize);

/// Encode an uplink message containing a traffic report.
// Values should be those accumulated since the InitIndUlMsg was
// sent.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeTrafficReportIndUlMsg (char * pBuffer,
                                      TrafficReportIndUlMsg_t * pMsg,
                                      char ** ppLog,
                                      uint32_t * pLogSize);

/// Encode a downlink message that retrieves a traffic report.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeTrafficReportGetReqDlMsg (char * pBuffer,
                                         char ** ppLog,
                                         uint32_t * pLogSize);

/// Encode an uplink message that is sent as a response to a
// TrafficReportGetReqDlMsg.
// Values should be those accumulated since the InitIndUlMsg was
// sent.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeTrafficReportGetCnfUlMsg (char * pBuffer,
                                         TrafficReportGetCnfUlMsg_t * pMsg,
                                         char ** ppLog,
                                         uint32_t * pLogSize);

/// Encode a downlink message that sets the Traffic Test mode settings
// for the UTM.  This must precede a mode change command to
// enter traffic test mode.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeTrafficTestModeParametersSetReqDlMsg (char * pBuffer,
                                                     TrafficTestModeParametersSetReqDlMsg_t * pMsg,
                                                     char ** ppLog,
                                                     uint32_t * pLogSize);

/// Encode an uplink message that is sent in response to a
// TrafficTestModeParametersSetReqMsg.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeTrafficTestModeParametersSetCnfUlMsg (char * pBuffer,
                                                     TrafficTestModeParametersSetCnfUlMsg_t * pMsg,
                                                     char ** ppLog,
                                                     uint32_t * pLogSize);

/// Encode a downlink message that gets the Traffic Test mode settings
// from the UTM.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeTrafficTestModeParametersGetReqDlMsg (char * pBuffer,
                                                     char ** ppLog,
                                                     uint32_t * pLogSize);

/// Encode an uplink message that is sent in response to a
// TrafficTestModeParametersGetReqDlMsg.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeTrafficTestModeParametersGetCnfUlMsg (char * pBuffer,
                                                     TrafficTestModeParametersGetCnfUlMsg_t * pMsg,
                                                     char ** ppLog,
                                                     uint32_t * pLogSize);

/// Encode an uplink traffic test datagram.
// This breaks all the rules, it is not a message, it is a whole
// datagram consisting of an ID followed by fill of a specified
// length.
// \param pBuffer  A pointer to the buffer to encode into.  Must
// be TRAFFIC_TEST_RULE_BREAKER_LENGTH bytes long.
// \param pSpec a pointer to the spec for the datagram
// (i.e. fill and length).
// \param isDownlink if true this will be encoded as a downlink
// datagram, otherwise it will be encoded as an uplink datagram
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeTrafficTestModeRuleBreakerDatagram (char * pBuffer,
                                                   TrafficTestModeRuleBreakerDatagram_t * pSpec,
                                                   bool isDownlink,
                                                   char ** ppLog,
                                                   uint32_t * pLogSize);

/// Encode an uplink message containing the traffic report format
// used in Traffic Test mode.
// Values should be those accumulated since the start of Traffic
// Test mode.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeTrafficTestModeReportIndUlMsg (char * pBuffer,
                                              TrafficTestModeReportIndUlMsg_t * pMsg,
                                              char ** ppLog,
                                              uint32_t * pLogSize);

/// Encode a downlink message that retrieves a Test Mode traffic report.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeTrafficTestModeReportGetReqDlMsg (char * pBuffer,
                                                 char ** ppLog,
                                                 uint32_t * pLogSize);

/// Encode an uplink message that is sent as a response to a
// TrafficTestModeReportGetReqDlMsg.
// Values should be those accumulated since the start of the
// last Traffic Test mode.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeTrafficTestModeReportGetCnfUlMsg (char * pBuffer,
                                                 TrafficTestModeReportGetCnfUlMsg_t * pMsg,
                                                 char ** ppLog,
                                                 uint32_t * pLogSize);

/// Encode an uplink message containing the activity of the module.
// Values should be those accumulated since boot.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeActivityReportIndUlMsg (char * pBuffer,
                                       ActivityReportIndUlMsg_t * pMsg,
                                       char ** ppLog,
                                       uint32_t * pLogSize);

/// Encode a downlink message that retrieves an Activity report.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeActivityReportGetReqDlMsg (char * pBuffer,
                                          char ** ppLog,
                                          uint32_t * pLogSize);

/// Encode an uplink message that is sent as a response to a
// ActivityReportGetReqDlMsg.
// Values should be those accumulated since boot.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeActivityReportGetCnfUlMsg (char * pBuffer,
                                          ActivityReportGetCnfUlMsg_t * pMsg,
                                          char ** ppLog,
                                          uint32_t * pLogSize);

/// Encode an uplink message that contains a debug string.
// \param pBuffer  A pointer to the buffer to encode into.  The
// buffer length must be at least MAX_MESSAGE_SIZE long
// \param pMsg  A pointer to the message to send.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The number of bytes encoded.
uint32_t encodeDebugIndUlMsg (char * pBuffer,
                              DebugIndUlMsg_t * pMsg,
                              char ** ppLog,
                              uint32_t * pLogSize);

// ----------------------------------------------------------------
// MESSAGE DECODING FUNCTIONS
// ----------------------------------------------------------------

/// The outcome of message decoding.
//
// !!! When you add anything to the generic of UL sections
// here (but not the DL bits as C# never decodes a DL thing),
// align it with the DLL exported version in the dll wrapper files
// so that the C# application can decode it.
typedef enum DecodeResultTag_t
{
     DECODE_RESULT_FAILURE = 0,         //!< Generic failed decode.
     DECODE_RESULT_INPUT_TOO_SHORT,     //!< Not enough input bytes.
     DECODE_RESULT_OUTPUT_TOO_SHORT,    //!< Not enough room in the
                                        //! output.
     DECODE_RESULT_UNKNOWN_MSG_ID,      //!< Rogue message ID.
     DECODE_RESULT_BAD_MSG_FORMAT,      //!< A problem with the format of a message.
     DECODE_RESULT_BAD_TRAFFIC_TEST_MODE_DATAGRAM, //!< A traffic test datagram has been decoded
                                                   //! but the contents are not of the right length.
     DECODE_RESULT_OUT_OF_SEQUENCE_TRAFFIC_TEST_MODE_DATAGRAM, //!< A traffic test datagram has been decoded
                                                               //! but the sequence number is not as expected.
     DECODE_RESULT_BAD_CHECKSUM, //!< The checksum on a message was bad, it should be ignored.
     DECODE_RESULT_DL_MSG_BASE = 0x40,  //!< From here on are the
                                        //! downlink messages.
                                        // !!! If you add one here
                                        // update the next line !!!
     DECODE_RESULT_TRANSPARENT_DL_DATAGRAM = DECODE_RESULT_DL_MSG_BASE,
     DECODE_RESULT_PING_REQ_DL_MSG,
     DECODE_RESULT_PING_CNF_DL_MSG,
     DECODE_RESULT_REBOOT_REQ_DL_MSG,
     DECODE_RESULT_DATE_TIME_SET_REQ_DL_MSG,
     DECODE_RESULT_DATE_TIME_GET_REQ_DL_MSG,
     DECODE_RESULT_MODE_SET_REQ_DL_MSG,
     DECODE_RESULT_MODE_GET_REQ_DL_MSG,
     DECODE_RESULT_HEARTBEAT_SET_REQ_DL_MSG,
     DECODE_RESULT_REPORTING_INTERVAL_SET_REQ_DL_MSG,
     DECODE_RESULT_INTERVALS_GET_REQ_DL_MSG,
     DECODE_RESULT_MEASUREMENTS_GET_REQ_DL_MSG,
     DECODE_RESULT_MEASUREMENT_CONTROL_SET_REQ_DL_MSG,
     DECODE_RESULT_MEASUREMENTS_CONTROL_GET_REQ_DL_MSG,
     DECODE_RESULT_MEASUREMENTS_CONTROL_DEFAULTS_SET_REQ_DL_MSG,
     DECODE_RESULT_TRAFFIC_REPORT_GET_REQ_DL_MSG,
     DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_SET_REQ_DL_MSG,
     DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_GET_REQ_DL_MSG,
     DECODE_RESULT_TRAFFIC_TEST_MODE_RULE_BREAKER_DL_DATAGRAM,
     DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_GET_REQ_DL_MSG,
     DECODE_RESULT_ACTIVITY_REPORT_GET_REQ_DL_MSG,  // !!! If you add one here
                                                    // update the next line !!!
     MAX_DL_REQ_MSG = DECODE_RESULT_ACTIVITY_REPORT_GET_REQ_DL_MSG,
     DECODE_RESULT_UL_MSG_BASE = 0x80,    //!< From here on are the
                                          //! uplink messages.
     DECODE_RESULT_TRANSPARENT_UL_DATAGRAM = DECODE_RESULT_UL_MSG_BASE,
     DECODE_RESULT_PING_REQ_UL_MSG,
     DECODE_RESULT_PING_CNF_UL_MSG,
     DECODE_RESULT_INIT_IND_UL_MSG,
     DECODE_RESULT_DATE_TIME_IND_UL_MSG,
     DECODE_RESULT_DATE_TIME_SET_CNF_UL_MSG,
     DECODE_RESULT_DATE_TIME_GET_CNF_UL_MSG,
     DECODE_RESULT_MODE_SET_CNF_UL_MSG,
     DECODE_RESULT_MODE_GET_CNF_UL_MSG,
     DECODE_RESULT_HEARTBEAT_SET_CNF_UL_MSG,
     DECODE_RESULT_REPORTING_INTERVAL_SET_CNF_UL_MSG,
     DECODE_RESULT_INTERVALS_GET_CNF_UL_MSG,
     DECODE_RESULT_POLL_IND_UL_MSG,
     DECODE_RESULT_MEASUREMENTS_IND_UL_MSG,
     DECODE_RESULT_MEASUREMENTS_GET_CNF_UL_MSG,
     DECODE_RESULT_MEASUREMENTS_CONTROL_IND_UL_MSG,
     DECODE_RESULT_MEASUREMENT_CONTROL_SET_CNF_UL_MSG,
     DECODE_RESULT_MEASUREMENTS_CONTROL_GET_CNF_UL_MSG,
     DECODE_RESULT_MEASUREMENTS_CONTROL_DEFAULTS_SET_CNF_UL_MSG,
     DECODE_RESULT_TRAFFIC_REPORT_IND_UL_MSG,
     DECODE_RESULT_TRAFFIC_REPORT_GET_CNF_UL_MSG,
     DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF_UL_MSG,
     DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF_UL_MSG,
     DECODE_RESULT_TRAFFIC_TEST_MODE_RULE_BREAKER_UL_DATAGRAM,
     DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_IND_UL_MSG,
     DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_GET_CNF_UL_MSG,
     DECODE_RESULT_ACTIVITY_REPORT_IND_UL_MSG,
     DECODE_RESULT_ACTIVITY_REPORT_GET_CNF_UL_MSG,
     DECODE_RESULT_DEBUG_IND_UL_MSG,         // !!! If you add one here
                                             // update the next line !!!
     MAX_UL_REQ_MSG = DECODE_RESULT_DEBUG_IND_UL_MSG,
     MAX_NUM_DECODE_RESULTS             //!< The maximum number of
                                        //! decode results.
} DecodeResult_t;

/// Decode a downlink message. When a datagram has been received
// this function should be called iteratively to decode all the
// messages contained within it.  The result, in pOutputBuffer,
// should be cast by the calling function to DlMsgUnion_t and
// the relevant member selected according to the
// DecodeResult_t code.
// \param ppInBuffer  A pointer to the pointer to decode from.
// On completion this is pointing to the next byte that
// could be decoded, after the currently decoded message,
// in the buffer.
// \param sizeInBuffer  The number of bytes left to decode.
// \param pOutBuffer  A pointer to the buffer to write the
// result into.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The result of the decoding, which hopefully says
// what message has been decoded.
DecodeResult_t decodeDlMsg (const char ** ppInBuffer,
                            uint32_t sizeInBuffer,
                            DlMsgUnion_t * pOutBuffer,
                            char ** ppLog,
                            uint32_t * pLogSize);

/// Decode an uplink message. When a datagram has been received
// this function should be called iteratively to decode all the
// messages contained within it.  The result, in pOutputBuffer,
// should be cast by the calling function to UlMsgUnion_t and
// the relevant member selected according to the
// DecodeResult_t code.
// \param ppInBuffer  A pointer to the pointer to decode from.
// On completion this is pointing to the next byte that
// could be decoded, after the currently decoded message,
// in the buffer.
// \param sizeInBuffer  The number of bytes left to decode.
// \param pOutBuffer  A pointer to the buffer to write the
// result into.
// \param ppLog  Optionally, a pointer to a pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pLogSize  Optional pointer to the size of the log buffer.
// Should be present if pLog is present.
// \return  The result of the decoding, which hopefully says
// what message has been decoded.
DecodeResult_t decodeUlMsg (const char ** ppInBuffer,
                            uint32_t sizeInBuffer,
                            UlMsgUnion_t * pOutBuffer,
                            char ** ppLog,
                            uint32_t * pLogSize);

// ----------------------------------------------------------------
// LOGGING FUNCTIONS
// ----------------------------------------------------------------

/// Return a string representing the hex value of an uint32_t.
// \param value the value.
// \return the string.
const char * getHexAsString (uint32_t value);

/// Return a string representing a Boolean (for logging).
// \param boolValue the Boolean value.
// \return the string.
const char * getStringBoolean (bool boolValue);

/// Return a string representing a wakeup code (for logging).
// \param wakeupCode the wake-up code.
// \return the string.
const char * getStringWakeUpCode (WakeUpCode_t wakeupCode);

/// Return a string representing the mode (for logging).
// \param mode the mode.
// \return the string.
const char * getStringMode (Mode_t mode);

/// Return a string representing the thing that set
// the time on the UTM (for logging).
// \param timeSetBy the source of time setting.
// \return the string.
const char * getStringTimeSetBy (TimeSetBy_t timeSetBy);

/// Return a string representing the amount of energy left (for logging).
// \param energyLeft the energy left.
// \return the string.
const char * getStringEnergyLeft (EnergyLeft_t energyLeft);

/// Return a string representing the amount of disk space left (for logging).
// \param diskSpaceLeft the diskSpace left.
// \return the string.
const char * getStringDiskSpaceLeft (DiskSpaceLeft_t diskSpaceLeft);

/// Return a string representing the charger state (for logging).
// \param chargerState the charger state.
// \return the string.
const char * getStringChargerState (ChargerState_t chargerState);

/// Return a string representing the measurement type (for logging).
// \param chargerState the charger state.
// \return the string.
const char * getStringMeasurementType (MeasurementType_t type);

/// Calculate the byes used and the new buffer size after an snprintf().
// \param pBufferSize  a pointer to the original size of the buffer.  This
// will be modified on exit to reflect the bytes used.
// \param bytesUsed  the number of bytes that snprintf() has reported used.
// \return the actual number of bytes used, taking into account truncation.
uint32_t calcBytesUsed (uint32_t *pBufferSize, uint32_t bytesUsed);

/// Log a begin tag.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param pTag  a pointer to the tag to use.
// \return the number of bytes used.
uint32_t logBeginTag (char * pBuffer, uint32_t *pBufferSize, const char *pTag);

/// Log a begin tag with a string value.
// \param pBuffer a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param pTag  a pointer to the tag to use.
// \param pValue  a pointer to the value to use.
// \return the number of bytes used.
uint32_t logBeginTagWithStringValue (char * pBuffer, uint32_t *pBufferSize, const char *pTag, const char * pValue);

/// Log an end tag.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param pTag  a pointer to the tag to use.
// \return the number of bytes used.
uint32_t logEndTag (char * pBuffer, uint32_t *pBufferSize, const char *pTag);

/// Log just a tag with no values.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param pTag  a pointer to the tag to use.
// \return the number of bytes used.
uint32_t logFlagTag(char * pBuffer, uint32_t *pBufferSize, const char *pTag);

/// Log a tag with a simple string value.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param pTag  a pointer to the tag to use.
// \param pValue  a pointer to the value to use.
// \return the number of bytes used.
uint32_t logTagWithStringValue (char * pBuffer, uint32_t *pBufferSize, const char * pTag, const char * pValue);

/// Log a tag with a simple uint32_t value.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param pTag  a pointer to the tag to use.
// \param value  the value to use.
// \return the number of bytes used.
uint32_t logTagWithUint32Value (char * pBuffer, uint32_t *pBufferSize, const char * pTag, uint32_t value);

/// Log a tag with a uint32_t value preceded by a presence flag.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param pTag  a pointer to the tag to use.
// \param present true if the value is present.
// \param value  the value to use.
// \return the number of bytes used.
uint32_t logTagWithPresenceAndUint32Value (char * pBuffer, uint32_t *pBufferSize, const char * pTag, bool present, uint32_t value);

/// Log the contents of a TransparentDatagram
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param datagramSize  the size of the contents of the data.
// \param pContents  a pointer to the start of the contents.
// \return the number of bytes used.
uint32_t logTransparentData (char * pBuffer, uint32_t *pBufferSize, uint32_t datagramSize, const char * pContents);

/// Log a date/time value.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param timeValue  the value to use.
// \return the number of bytes used.
uint32_t logDateTime (char * pBuffer, uint32_t *pBufferSize, uint32_t timeValue);

/// Log a heartbeat value.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param heartbeatValue  a pointer to the heartbeat value.
// \param snapToRtc  true if heartbeatValue is a snapTo value.
// \return the number of bytes used.
uint32_t logHeartbeat (char * pBuffer, uint32_t *pBufferSize, uint32_t heartbeatValue, bool snapToRtc);

/// Log an RSSI value.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param rssi the value to use.
// \return the number of bytes used.
uint32_t logRssi (char * pBuffer, uint32_t *pBufferSize, Rssi_t rssi);

/// Log an RSRP value.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param rsrp the value to use.
// \param isSyncedWithRssi true if the RSRP value was taken at the samem time as
// the RSSI value in the same report
// \return the number of bytes used.
uint32_t logRsrp (char * pBuffer, uint32_t *pBufferSize, Rssi_t rsrp, bool isSyncedWithRssi);

/// Log the temperature.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param temperature the value to use.
// \return the number of bytes used.
uint32_t logTemperature (char * pBuffer, uint32_t *pBufferSize, Temperature_t temperature);

/// Log the battery voltage
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param batteryVoltage the voltage value to use.
// \return the number of bytes used.
uint32_t logBatteryVoltage (char * pBuffer, uint32_t *pBufferSize, uint16_t batteryVoltage);

/// Log the battery energy left
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param batteryEnergy the energy value to use.
// \return the number of bytes used.
uint32_t logBatteryEnergy (char * pBuffer, uint32_t *pBufferSize, uint32_t batteryEnergy);

/// Log a position value.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param latitude the latitude value.
// \param longitude the longitude value.
// \param elevation the elevation value.
// \return the number of bytes used.
uint32_t logPosition (char * pBuffer, uint32_t *pBufferSize, uint32_t latitude, uint32_t longitude, uint32_t elevation);

/// Log an "OnlyRecordIf" control.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param present true if "OnlyRecordIf" is present.
// \param value  the "OnlyRecordIf" value.
// \param aboveNotBelow the "aboveNotBelow" value.
// \param atTransitionOnly the "atTransitionOnly" value.
// \param isOneShot the "isOneShot" value.
// \return the number of bytes used.
uint32_t logOnlyRecordIf (char * pBuffer, uint32_t *pBufferSize, bool present, int32_t value, bool aboveNotBelow, bool atTransitionOnly, bool isOneShot);

/// Log the traffic report data for transmit from the UTM.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param datagrams the number of datagrams.
// \param bytes the number of bytes.
// \return the number of bytes used to code the message.
uint32_t logTrafficReportUl (char * pBuffer, uint32_t *pBufferSize, uint32_t datagrams, uint32_t bytes);

/// Log the traffic report data for receive at the UTM.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param datagrams the number of datagrams.
// \param bytes the number of bytes.
// \param badChecksum the number of datagrams with a bad checksum.
// \return the number of bytes used to code the message.
uint32_t logTrafficReportDl(char * pBuffer, uint32_t *pBufferSize, uint32_t datagrams, uint32_t bytes, uint32_t badChecksum);

/// Log the traffic test mode parameters for the UTM to transmit.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param count the number of datagrams to send.
// \param length the length of the datagrams.
// \param noReportsDuringTest whether reports are sent during the test or not.
// \return the number of bytes used to code the message.
uint32_t logTrafficTestModeParametersUl(char * pBuffer, uint32_t *pBufferSize, uint32_t count, uint32_t length, bool noReportsDuringTest);

/// Log the traffic test mode parameters for the UTM to receive.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param count the number of datagrams to expect.
// \param length the length of the datagrams.
// \return the number of bytes used to code the message.
uint32_t logTrafficTestModeParametersDl (char * pBuffer, uint32_t *pBufferSize, uint32_t count, uint32_t length);

/// Log the Traffic Test mode report transmit values from the UTM.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param datagramCount the number of datagrams transmitted.
// \param bytes the number of bytes transmitted.
// \return the number of bytes used to code the message.
uint32_t logTrafficTestModeReportUl (char * pBuffer, uint32_t *pBufferSize, uint32_t datagramCount, uint32_t bytes);

/// Log the Traffic Test mode report receive values from the UTM.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param datagramCount the number of datagrams received.
// \param bytes the number of bytes received.
// \param outOfOrder the number of datagrams received out of order.
// \param bad the number of datagrams received with short length.
// \param missed the number of datagrams missed.
// \return the number of bytes used to code the message.
uint32_t logTrafficTestModeReportDl (char * pBuffer, uint32_t *pBufferSize, uint32_t datagramCount, uint32_t bytes, uint32_t outOfOrder, uint32_t bad, uint32_t missed);

/// Log a Traffic Test mode datagram
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param fill the value that the datagram is filled with.
// \param length the length of the datagram.
// \return the number of bytes used to code the message.
uint32_t logTrafficTestModeRuleBreakerDatagram (char * pBuffer, uint32_t *pBufferSize, uint8_t fill, uint32_t length);

/// Log the Activity report receive values from the UTM.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param totalTransmitMilliseconds how long the module has been on for transmit.
// \param totalReceiveMilliseconds  how long the module has been on for receive.
// \param upTimeSeconds the time that the UTM has been up for.
// \return the number of bytes used to code the message.
uint32_t logActivityReport(char * pBuffer, uint32_t *pBufferSize, uint32_t totalTransmitMilliseconds, uint32_t totalReceiveMilliseconds, uint32_t upTimeSeconds);

/// Log the Activity report receive values from the UTM.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param txPowerDbmPresent if true then the Tx power is present.
// \param txPowerDbm  the Tx power in dBm.
// \param mcsPresent if true then the uplink MCS is present.
// \param mcs the uplink MCS.
// \return the number of bytes used to code the message.
uint32_t logUlRf(char * pBuffer, uint32_t *pBufferSize, bool txPowerDbmPresent, int8_t txPowerDbm, bool mcsPresent, uint8_t mcs);

/// Log the Activity report receive values from the UTM.
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param mcsPresent if true then the downlink MCS is present.
// \param mcs the downlink MCS.
// \return the number of bytes used to code the message.
uint32_t logDlRf(char * pBuffer, uint32_t *pBufferSize, bool mcsPresent, uint8_t mcs);

/// Log a string from a DebugInd message
// \param pBuffer  a pointer to the logging buffer.
// \param pBufferSize a pointer to the size of the logging buffer.  This
// will be modified on exit to reflect the bytes used.
// \param pString a pointer to the string from the message.
// \param length the length of the string.
// \return the number of bytes used to code the message.
uint32_t logDebug (char * pBuffer, uint32_t *pBufferSize, char * pString, uint8_t length);

// ----------------------------------------------------------------
// MISC FUNCTIONS
// ----------------------------------------------------------------

/// Only used in the DLL form, sets up the "printf()" function
// for logging.
// \param guiPrintToConsole  the printf function.
void initDll (void (*guiPrintToConsole) (const char *));

#ifdef __cplusplus
}
#endif

#endif

// End Of File
