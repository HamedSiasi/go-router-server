/* UTM interface
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

/**
 * @file utm_msg_handler.c
 * This file implements the API to the UTM.
 */

#include <stdint.h>
#ifndef _MSC_VER
# include <stdbool.h>
#else
typedef unsigned char bool; // MSVC has it's own way...
#  define false 0
#  define true 1
#endif
#include <stdio.h>
#include <stdarg.h> // for va_...
#include <string.h> // for memcpy()
#include <time.h> // for time_t

#include <utm_api.h>

#ifdef __cplusplus
extern "C"
{
#endif

#define DEBUG

#ifdef DEBUG
#define MESSAGE_CODEC_LOGMSG(...)    logMsg(__VA_ARGS__)
#else
#define MESSAGE_CODEC_LOGMSG(...)
#endif

// Note on debugging GCC issues on Windows, with thanks
// to Dennis Yurichev:
//
// http://dennisyurichev.blogspot.co.uk/2013/05/warning-invalid-parameter-passed-to-c.html
//
// While compiling under GCC (for Golang) a Mingw
// run-time library uses Windows calls to do certain things,
// like printf or snprintf etc.  I was getting a segmentation
// fault and so, running under gdb, saw the message:
// "warning: Invalid parameter passed to C runtime function."
// It is apparently being dumped into the debugger in
// msvcrt.dll's __invoke_watson() call, using well-known
// OutputDebugStringA().  So to trace what his happening you
// need to add a breakpoint there, by typing:
//
// break OutputDebugStringA
//
// ...in GDB.  You should then be able to find out what's up
// with a backtrace (type "bt").
//
// In my case it turns out that I was calling strftime() with a
// C99 formatter, %F, which is a but to recent for Golang's use
// of GCC.  I needed to use the equivalent %Y-%m-%d instead.

// ----------------------------------------------------------------
// GENERAL COMPILE-TIME CONSTANTS
// ----------------------------------------------------------------

/// The max size of a debug message (including terminator)
#define MAX_DEBUG_PRINTF_LEN 128

/// The maximum number of bitmap bytes expected
#define MAX_BITMAP_BYTES 1

// For logging
#define TAG_MSG_BI "MsgBi"
#define TAG_MSG_UL "MsgUl"
#define TAG_MSG_DL "MsgDl"
#define TAG_MSG_NAME "MsgName"
#define TAG_MSG_CONTENTS "MsgContents"
#define TAG_MSG_SIZE "MsgSize"
#define TAG_MODE "Mode"
#define TAG_TIME_SET_BY "TimeSetBy"
#define TAG_ENERGY_LEFT "EnergyLeft"
#define TAG_DISK_SPACE_LEFT "DiskSpaceLeft"
#define TAG_SET_DATE_ONLY "SetDateOnly"
#define TAG_REPORTING_INTERVAL "ReportingInterval"
#define TAG_CELL_ID "CellId"
#define TAG_CHARGER_STATE "ChargerState"
#define TAG_MEASUREMENT_TYPE "MeasurementType"
#define TAG_MEASUREMENT_CONTROL_GENERIC "MeasurementControlGeneric"
#define TAG_MEASUREMENT_CONTROL_POWER_STATE "MeasurementControlPower"
#define TAG_MAX_REPORTING_INTERVAL "MaxReportingInterval"
#define TAG_HYSTERESIS "Hysteresis"
#define TAG_CHARGER_STATE_MEASUREMENT_CONTROL "ChargerMeasurementControl"
#define TAG_REPORT_IMMEDIATELY "ReportImmediately"
#define TAG_VOLTAGE_MEASUREMENT_CONTROL "VoltageMeasurementControl"
#define TAG_ENERGY_MEASUREMENT_CONTROL "EnergyMeasurementControl"
#define TAG_TIMEOUT "Timeout"
#define TAG_TIMED_OUT "TimedOut"
#define TAG_MSG_CHECKSUM_GOOD "ChecksumGood"
#define VALUE_UNKNOWN_STRING "??"

// ----------------------------------------------------------------
// TYPES
// ----------------------------------------------------------------

// ----------------------------------------------------------------
// PRIVATE VARIABLES
// ----------------------------------------------------------------

#ifdef WIN32
static void (* mp_guiPrintToConsole) (const char*) = NULL;
#endif
static const char hexTable[] =
{ '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f' };
static char gLogBuffer[255];
static const char * gStringBoolean[] =
{ "False", "True" };
static const char * gStringWakeUpCode[] =
{ "Ok", "Watchdog", "NetworkProblem", "SdCardProblem", "SupplyProblem", "ProtocolProblem", "ModuleNotResponding", "HwProblem", "MemoryAllocProblem", "GenericFailure", "CommandedReboot" };
static const char * gStringMode[] =
{ "Null", "SelfTest", "Commissioning", "StandardTrx", "TrafficTest" };
static const char * gStringTimeSetBy[] =
{ "Unknown", "GNSS", "PC", "WebApi" };
static const char * gStringEnergyLeft[] =
{ "LessThan5Percent", "LessThan10Percent", "MoreThan10Percent", "MoreThan30Percent", "MoreThan50Percent", "MoreThan70Percent", "MoreThan90Percent" };
static const char * gStringDiskSpaceLeft[] =
{ "LessThan1Gb", "MoreThan1Gb", "MoreThan2Gb", "MoreThan4Gb" };
static const char * gStringChargerState[] =
{ "Unknown", "Off", "On", "Fault" };
static const char * gStringMeasurementType[] =
{ "Unknown", "GNSS", "CellId", "RSSI", "RSRP", "Temperature", "PowerState" };

static char gValueAsHexString[] = "0x00000000";  // Enough room for an uint32_t represented as a string

// ----------------------------------------------------------------
// PRIVATE FUNCTION PROTOTYPES
// ----------------------------------------------------------------

/// Encode a boolean value.
// \param pBuffer  A pointer to where the encoded
// value should be placed.
// \param value The Boolean value.
// \return  The number of bytes encoded.
static uint32_t encodeBool(char * pBuffer, bool value);

/// Decode a Boolean value.
// \param ppBuffer  A pointer to the pointer to decode.
// On completion this points to the location after the
// bool in the input buffer.
// \return  The decoded value.
static bool decodeBool(const char ** ppBuffer);

/// Encode a uint32_t value.
// \param pBuffer  A pointer to the value to decode.
// \param value The value.
// \return  The number of bytes encoded.
static uint32_t encodeUint32(char * pBuffer, uint32_t value);

/// Decode a uint32_t value.
// \param ppBuffer  A pointer to the pointer to decode.
// On completion this points to the location after the
// uint32_t in the input buffer.
static uint32_t decodeUint32(const char ** ppBuffer);

/// Encode a uint24_t value.
// \param pBuffer  A pointer to the value to decode.
// \param value The value.
// \return  The number of bytes encoded.
static uint32_t encodeUint24(char * pBuffer, uint32_t value);

/// Decode a uint24_t value.
// \param ppBuffer  A pointer to the pointer to decode.
// On completion this points to the location after the
// uint24_t in the input buffer.
static uint32_t decodeUint24(const char ** ppBuffer);

/// Encode a uint16_t value.
// \param pBuffer  A pointer to the value to decode.
// \param value The value.
// \return  The number of bytes encoded.
static uint32_t encodeUint16(char * pBuffer, uint16_t value);

/// Decode a uint16_t value.
// \param ppBuffer  A pointer to the pointer to decode.
// On completion this points to the location after the
// uint16_t in the input buffer.
static uint32_t decodeUint16(const char ** ppBuffer);

/// Encode the measurements.
// \param pBuffer A pointer to the buffer in which to encode.
// \param pMeasurements A pointer to the measurements.
// \param ppLog  Optionally, a pointer to a pointer to
// a log buffer in which/ to write an XML log of the
// encoded message.
// The value will be adjusted as the log is consumed.
// \param pLogSize  A pointer to the size of the log
// buffer.  This must be present if *ppLog is present.
// The value will be adjusted as the log is consumed.
// \return  The number of bytes encoded.
static uint32_t encodeMeasurements(char * pBuffer, Measurements_t * pMeasurements, char ** ppLog, uint32_t * pLogSize);

/// Decode the measurements into pValue.
// \param ppBuffer A pointer to the pointer to decode.
// On completion this points to the location after the
// Measurements_t in the input buffer.
// \param pMeasurements A place to put the measurements.
// \param ppLog  Optionally, a pointer to a pointer to
// a log buffer in which/ to write an XML log of the
// encoded message.
// The value will be adjusted as the log is consumed.
// \param pLogSize  A pointer to the size of the log
// buffer.  This must be present if *ppLog is present.
// The value will be adjusted as the log is consumed.
// \return true if the decode is successful, otherwise false.
static bool decodeMeasurements(const char ** ppBuffer, Measurements_t * pMeasurements, char ** ppLog, uint32_t * pLogSize);

/// Write an XML log of the measurements structure.
// \param pMeasurements A pointer to the measurements.
// \param pBufferSize  A pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pBufferSizeSize  A pointer to the size of the log
// buffer.  This will be adjusted as the log is consumed.
// \return the number of bytes written to the log.
static uint32_t logMeasurements(char * pBufferSize, uint32_t * pBufferSizeSize, Measurements_t * pMeasurements);

/// Encode measurement control.
// \param pBuffer A pointer to the buffer in which to encode.
// \param measurementType The type of measurement control to encode.
// \param pMeasurementControlUnion A pointer to the measurement control union.
// \param ppLog  Optionally, a pointer to a pointer to
// a log buffer in which/ to write an XML log of the
// encoded message.
// The value will be adjusted as the log is consumed.
// \param pLogSize  A pointer to the size of the log
// buffer.  This must be present if *ppLog is present.
// The value will be adjusted as the log is consumed.
// \return  The number of bytes encoded.
static uint32_t encodeMeasurementControl(char * pBuffer, MeasurementType_t measurementType, MeasurementControlUnion_t * pMeasurementControlUnion, char ** ppLog, uint32_t * pLogSize);

/// Encode the generic measurement control structure.
// \param pBuffer A pointer to the buffer in which to encode.
// \param pMeasurementControlGeneric A pointer to the generic measurement control structure.
// \param ppLog  Optionally, a pointer to a pointer to
// a log buffer in which/ to write an XML log of the
// encoded message.
// The value will be adjusted as the log is consumed.
// \param pLogSize  A pointer to the size of the log
// buffer.  This must be present if *ppLog is present.
// The value will be adjusted as the log is consumed.
// \return  The number of bytes encoded.
static uint32_t encodeMeasurementControlGeneric(char * pBuffer, MeasurementControlGeneric_t * pMeasurementControlGeneric, char ** ppLog, uint32_t * pLogSize);

/// Encode measurement control.
// \param ppBuffer A pointer to the pointer to decode.
// On completion this points to the location after the
// MeasurementControlUnion_t in the input buffer.
// \param measurementType The type of measurement control to encode.
// \param pMeasurementControlUnion A pointer to the measurement control union.
// \param ppLog  Optionally, a pointer to a pointer to
// a log buffer in which/ to write an XML log of the
// encoded message.
// The value will be adjusted as the log is consumed.
// \param pLogSize  A pointer to the size of the log
// buffer.  This must be present if *ppLog is present.
// The value will be adjusted as the log is consumed.
// \return true if the decode is successful, otherwise false.
static bool decodeMeasurementControl(const char ** ppBuffer, MeasurementType_t measurementType, MeasurementControlUnion_t * pMeasurementControlUnion, char ** ppLog, uint32_t * pLogSize);

/// Decode the generic measurement control structure.
// \param ppBuffer A pointer to the pointer to decode.
// On completion this points to the location after the
// MeasurementControlGeneric_t in the input buffer.
// \param pMeasurementControlGeneric A pointer to the generic measurement control structure.
// \param ppLog  Optionally, a pointer to a pointer to
// a log buffer in which to write an XML log of the
// encoded message.
// The value will be adjusted as the log is consumed.
// \param pLogSize  A pointer to the size of the log
// buffer.  This must be present if *ppLog is present.
// The value will be adjusted as the log is consumed.
// \return true if the decode is successful, otherwise false.
static bool decodeMeasurementControlGeneric(const char ** ppBuffer, MeasurementControlGeneric_t * pMeasurementControlGeneric, char ** ppLog, uint32_t * pLogSize);

/// Decode the generic measurement control structure.
// \param pMeasurementControlGeneric A pointer to the generic measurement control structure.
// \param pBuffer  A pointer to a log buffer in which
// to write an XML log of the encoded message.
// \param pBufferSize  A pointer to the size of the log
// buffer.  This must be present if *ppLog is present.
// The value will be adjusted as the log is consumed.
// \return the number of bytes written to the log.
static uint32_t logMeasurementControlGeneric(char * pBuffer, uint32_t * pBufferSize, MeasurementControlGeneric_t * pMeasurementControlGeneric);

/// Decode a traffic test datagram.
// \param ppBuffer A pointer to the pointer to decode.
// On completion this points to the location after the
// TrafficTestModeRuleBreakerDatagram in the input buffer.
// \param pSpec  a pointer to the spec of the expected datagram (fill and length).
// The received fill and length values are copied into this after decoding.
// \param isDownlink true if this is a downlink message, else false
// \param ppLog  Optionally, a pointer to a pointer to
// a log buffer in which/ to write an XML log of the
// encoded message.
// The value will be adjusted as the log is consumed.
// \param pLogSize  A pointer to the size of the log
// buffer.  This must be present if *ppLog is present.
// The value will be adjusted as the log is consumed.
// \return true if the decode is successful and the expected
// value is present throughout, otherwise false.
static bool decodeTrafficTestModeRuleBreakerDatagram(const char ** ppBuffer, bool isDownlink, TrafficTestModeRuleBreakerDatagram_t *pSpec, char ** ppLog, uint32_t * pLogSize);

/// Limit an int to N bits with correct sign representation.
// Note that the values should be in-range for the given number
// of bits to convert correctly, out of range values will be screwy.
static __inline uint32_t limitInt(int32_t number, uint8_t bits);

/// Sign-extend a number of N bits (held inside an uin32_t) to an int32_t.
static __inline int32_t extendInt(uint32_t number, uint8_t bits);

/// Calculate a checksum.
static char calculateChecksum(const char * pBuffer, uint32_t bufferLength);

/// Log a message for debugging, "printf()" style.
// \param pFormat The printf() style parameters.
static void logMsg(const char * pFormat, ...);

// ----------------------------------------------------------------
// GENERIC PRIVATE FUNCTIONS
// ----------------------------------------------------------------

#if defined (_MSC_VER)

#define snprintf my_snprintf

__inline int32_t my_vsnprintf(char *outBuf, size_t size, const char *format, va_list ap)
{
    int32_t count = -1;

    if (size != 0)
    {
        count = _vsnprintf_s (outBuf, size, _TRUNCATE, format, ap);
    }

    if (count == -1)
    {
        count = _vscprintf(format, ap);
    }

    return count;
}

__inline int32_t my_snprintf(char *outBuf, size_t size, const char *format, ...)
{
    int32_t count;
    va_list ap;

    va_start (ap, format);
    count = my_vsnprintf (outBuf, size, format, ap);
    va_end (ap);

    return count;
}

#endif

/// Convert a string of bytes to hex
static int bytesToHexString(const char * inBuf, int size, char *outBuf, int lenOutBuf)
{
    int x = 0;
    int y = 0;

    for (x = 0; (x < size) && (y < lenOutBuf); x++)
    {
        outBuf[y] = hexTable[(inBuf[x] >> 4) & 0x0f]; // upper nibble
        y++;
        if (y < lenOutBuf)
        {
            outBuf[y] = hexTable[inBuf[x] & 0x0f]; // lower nibble
            y++;
            if (y < lenOutBuf)
            {
                outBuf[y] = 0; // Add a terminator in case this is the last bit
            }
        }
    }

    return y;
}

/// Encode a boolean value.
uint32_t encodeBool(char * pBuffer, bool value)
{
    *pBuffer = value;
    return 1;
}

/// Decode a boolean value.
bool decodeBool(const char ** ppBuffer)
{
    bool boolValue = false;

    if (**ppBuffer)
    {
        boolValue = true;
    }

    (*ppBuffer)++;

    return boolValue;
}

/// Encode a uint32_t
uint32_t encodeUint32(char * pBuffer, uint32_t value)
{
    uint32_t numBytesEncoded = 4;

    pBuffer[0] = 0xff & (value >> 24);
    pBuffer[1] = 0xff & (value >> 16);
    pBuffer[2] = 0xff & (value >> 8);
    pBuffer[3] = 0xff & value;

    return numBytesEncoded;
}

/// Decode a uint32_t
uint32_t decodeUint32(const char ** ppBuffer)
{
    uint32_t value = 0;

    value += ((**ppBuffer) & 0xFF) << 24;
    (*ppBuffer)++;
    value += ((**ppBuffer) & 0xFF) << 16;
    (*ppBuffer)++;
    value += ((**ppBuffer) & 0xFF) << 8;
    (*ppBuffer)++;
    value += ((**ppBuffer) & 0xFF);
    (*ppBuffer)++;

    return value;
}

/// Encode a uint24_t
uint32_t encodeUint24(char * pBuffer, uint32_t value)
{
    uint32_t numBytesEncoded = 3;

    pBuffer[0] = 0xff & (value >> 16);
    pBuffer[1] = 0xff & (value >> 8);
    pBuffer[2] = 0xff & value;

    return numBytesEncoded;
}

/// Decode a uint24_t
uint32_t decodeUint24(const char ** ppBuffer)
{
    uint32_t value = 0;

    value += ((**ppBuffer) & 0xFF) << 16;
    (*ppBuffer)++;
    value += ((**ppBuffer) & 0xFF) << 8;
    (*ppBuffer)++;
    value += ((**ppBuffer) & 0xFF);
    (*ppBuffer)++;

    return value;
}

/// Encode a uint16_t
uint32_t encodeUint16(char * pBuffer, uint16_t value)
{
    uint32_t numBytesEncoded = 2;

    pBuffer[0] = 0xff & (value >> 8);
    pBuffer[1] = 0xff & value;

    return numBytesEncoded;
}

/// Decode a uint16_t
uint32_t decodeUint16(const char ** ppBuffer)
{
    uint32_t value = 0;

    value += ((**ppBuffer) & 0xFF) << 8;
    (*ppBuffer)++;
    value += ((**ppBuffer) & 0xFF);
    (*ppBuffer)++;

    return value;
}

/// Encode a Measurements_t
// The measurements message is coded as follows:
//
// uint32_t    time             UTC seconds
// uint8_t     bytesToFollow    so that some checking can be done
// uint8_t     itemsBitmap0     see below
// [uint8_t    itemsBitmap1]... may repeat, see below
// [itemsUnion item0]...        zero or more items.
//
// The generic format of the itemsBitmap is:
//  bit   7 6 5 4 3 2 1 0
//  value y x x x x x x x
//
// ...where x is set to 1 if the associated logical item is
// present, otherwise 0, and y is set to 1 if another
// itemsBitmap follows, otherwise 0.
//
// There is no generic format for an item, each is defined
// explicitly and must be present in the order of the
// itemsBitmap (where the item at bit 0 of itemsBitmap0
// comes first).  The following item formats are defined:
//
// x  Name              Size      Format
// 0: GNSSPosition,     12 bytes: 32 bits lat, 32 bits long, 32 bits elevation
// 1: CellId,            2 bytes: 16 bits, logical cell ID
// 2: RSRP,              2 bytes: 15 bits, signed plus a 1 bit flag to indicate sync with RSSI
// 3: RSSI,              2 bytes: 15 bits, signed
// 4: Temperature,       1 byte:  -128 to +127 C
// 5: Power state,       4 bytes: bits 0-5: battery voltage, 0 == 0 volts, 63 = 10 volts
//                                bits 6-7: charger state
//                                bits 8-31: energy in mWh (24 bits, unsigned)

uint32_t encodeMeasurements(char * pBuffer, Measurements_t * pMeasurements, char ** ppLog, uint32_t * pLogSize)
{
    char *pBufferAtStart;
    char *pBytesToFollow;
    uint8_t bitMap = 0;

    pBufferAtStart = pBuffer;

    // Encode time
    pBuffer += encodeUint32(pBuffer, pMeasurements->time);

    // Set things up so that bytesToFollow can be filled in later
    pBytesToFollow = pBuffer;
    pBuffer++;

    // Now fill in the bit-map, determining which items are present
    // and filling in the bitmap according to the order of the items
    // in the structure
    if (pMeasurements->gnssPositionPresent)
    {
        bitMap |= 0x01;
    }
    if (pMeasurements->cellIdPresent)
    {
        bitMap |= 0x02;
    }
    if (pMeasurements->rsrpPresent)
    {
        bitMap |= 0x04;
    }
    if (pMeasurements->rssiPresent)
    {
        bitMap |= 0x08;
    }
    if (pMeasurements->temperaturePresent)
    {
        bitMap |= 0x10;
    }
    if (pMeasurements->powerStatePresent)
    {
        bitMap |= 0x20;
    }
    // That's 6 things coded up so write this bitmap value
    // but don't advance the pointer as it may be necessary
    // to set the extension bit
    *pBuffer = bitMap;

    // That's all we have but if we have more than 8 it would go as follows
    // if (pMeasurements->blahPresent || pMeasurements->blahBlahPresent)
    // {
    //     // Set the extension bit, re-write the last bitmap
    //     // value and start on the next one
    //     bitMap |= 0x80;
    //     *pBuffer = bitMap;
    //     pBuffer++;
    //     bitMap = 0;
    //
    //     if (pMeasurements->blahPresent)
    //     {
    //         bitMap |= 0x01;
    //     }
    //     if (pMeasurements->blahBlahPresent)
    //     {
    //         bitMap |= 0x02;
    //     }
    //
    //     // This is the end of the structure, so now write this
    //     // bitmap value
    //     *pBuffer = bitMap;
    // }

    // Advance the pointer
    pBuffer++;

    // Now fill in the actual values
    if (pMeasurements->gnssPositionPresent)
    {
        pBuffer += encodeUint32(pBuffer, (uint32_t) pMeasurements->gnssPosition.latitude);
        pBuffer += encodeUint32(pBuffer, (uint32_t) pMeasurements->gnssPosition.longitude);
        pBuffer += encodeUint32(pBuffer, (uint32_t) pMeasurements->gnssPosition.elevation);
    }
    if (pMeasurements->cellIdPresent)
    {
        pBuffer += encodeUint16(pBuffer, pMeasurements->cellId);
    }
    if (pMeasurements->rsrpPresent)
    {
        if (pMeasurements->rsrp.value < MIN_RSSI_RSRP)
        {
            pMeasurements->rsrp.value = MIN_RSSI_RSRP;
        }
        else
        {
            if (pMeasurements->rsrp.value > MAX_RSSI_RSRP)
            {
                pMeasurements->rsrp.value = MAX_RSSI_RSRP;
            }
        }
        pMeasurements->rsrp.value = (int16_t) limitInt((int32_t) pMeasurements->rsrp.value, 15);

        if (pMeasurements->rsrp.isSyncedWithRssi)
        {
            pMeasurements->rsrp.value |= 0x8000;
        }
        pBuffer += encodeUint16(pBuffer, pMeasurements->rsrp.value);
    }
    if (pMeasurements->rssiPresent)
    {
        if (pMeasurements->rssi < MIN_RSSI_RSRP)
        {
            pMeasurements->rssi = MIN_RSSI_RSRP;
        }
        else
        {
            if (pMeasurements->rssi > MAX_RSSI_RSRP)
            {
                pMeasurements->rssi = MAX_RSSI_RSRP;
            }
        }
        pMeasurements->rssi = (int16_t) limitInt((int32_t) pMeasurements->rssi, 15);
        pBuffer += encodeUint16(pBuffer, pMeasurements->rssi);
    }
    if (pMeasurements->temperaturePresent)
    {
        *pBuffer = (uint8_t) pMeasurements->temperature;
        pBuffer++;
    }
    if (pMeasurements->powerStatePresent)
    {
        uint8_t x = 0;

        if (pMeasurements->powerState.batteryMV > MAX_BATTERY_VOLTAGE_MV)
        {
            pMeasurements->powerState.batteryMV = MAX_BATTERY_VOLTAGE_MV;
        }
        x |= (((uint32_t) pMeasurements->powerState.batteryMV * 0x3F / 10000) & 0x3F);
        x |= ((pMeasurements->powerState.chargerState << 6) & 0xC0);
        *pBuffer = x;
        pBuffer++;
        if (pMeasurements->powerState.energyMWh > MAX_ENERGY_MWH)
        {
            pMeasurements->powerState.energyMWh = MAX_ENERGY_MWH;
        }
        pBuffer += encodeUint24(pBuffer, pMeasurements->powerState.energyMWh);
    }

    // This is as many as we currently have.  If we have more it goes
    // as follows
    // if (pMeasurements->blahPresent)
    // {
    //     Encode blah and advance pBuffer, e.g.
    //     pBuffer += encodeUint16 (pBuffer, pMeasurements->blah);
    // }
    // if (pMeasurements->blahBlahPresent)
    // {
    //     Encode blahBlah and advance pBuffer, e.g.
    //     pBuffer += encodeUint16 (pBuffer, pMeasurements->blahBlah);
    // }

    // Now fill in the value for bytesToFollow
    *pBytesToFollow = pBuffer - (pBufferAtStart + 5); // 5 for a UInt32 and bytesToFollow itself

    // Log it
    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logMeasurements(*ppLog, pLogSize, pMeasurements);
    }

    return (pBuffer - pBufferAtStart);
}

// Decode a Measurements_t
bool decodeMeasurements(const char ** ppBuffer, Measurements_t * pMeasurements, char ** ppLog, uint32_t * pLogSize)
{
    bool success = false;
    uint8_t x;
    const char *pBufferAfterEnd;
    uint8_t bitMapBytes[MAX_BITMAP_BYTES];
    bool moreBitmapBytes = true;

    memset(&(bitMapBytes[0]), 0, sizeof(bitMapBytes));
    memset(pMeasurements, 0, sizeof(*pMeasurements));

    // Decode time
    pMeasurements->time = decodeUint32(ppBuffer);

    // Decode bytesToFollow and add it to the current pointer to find the end
    pBufferAfterEnd = *ppBuffer + **ppBuffer + 1; //+1 because this is before the increment
    (*ppBuffer)++;

    // Decode the bitmap byte(s)
    for (x = 0; moreBitmapBytes && (*ppBuffer < pBufferAfterEnd); x++)
    {
        uint8_t y;

        y = (uint8_t) * *ppBuffer;
        (*ppBuffer)++;

        if (x < MAX_BITMAP_BYTES)
        {
            bitMapBytes[x] = y;
        }

        if ((y & 0x80) == 0)
        {
            moreBitmapBytes = false;
        }
    }

    // moreBitmapBytes should be false by now and,
    // if so, decode the values, in order
    if (!moreBitmapBytes)
    {
        // GNSS Position
        if (bitMapBytes[0] & 0x01)
        {
            pMeasurements->gnssPositionPresent = true;
            pMeasurements->gnssPosition.latitude = (int32_t) decodeUint32(ppBuffer);
            pMeasurements->gnssPosition.longitude = (int32_t) decodeUint32(ppBuffer);
            pMeasurements->gnssPosition.elevation = (int32_t) decodeUint32(ppBuffer);

        }
        // Cell ID
        if (bitMapBytes[0] & 0x02)
        {
            pMeasurements->cellIdPresent = true;
            pMeasurements->cellId = decodeUint16(ppBuffer);
        }
        // RSRP
        if (bitMapBytes[0] & 0x04)
        {
            uint32_t value = decodeUint16(ppBuffer);

            pMeasurements->rsrpPresent = true;
            if (value & 0x8000)
            {
                pMeasurements->rsrp.isSyncedWithRssi = true;
            }
            value &= ~0x8000;
            pMeasurements->rsrp.value = (Rssi_t) extendInt(value, 15);
        }
        // RSSI
        if (bitMapBytes[0] & 0x08)
        {
            uint16_t value = decodeUint16(ppBuffer) & ~0x8000;

            pMeasurements->rssiPresent = true;
            pMeasurements->rssi = (Rssi_t) extendInt((uint32_t) value, 15);
        }
        // Temperature
        if (bitMapBytes[0] & 0x10)
        {
            pMeasurements->temperaturePresent = true;
            pMeasurements->temperature = (Temperature_t) **ppBuffer;
            (*ppBuffer)++;
        }
        // PowerState
        if (bitMapBytes[0] & 0x20)
        {
            pMeasurements->powerStatePresent = true;
            x = (uint8_t) **ppBuffer;
            (*ppBuffer)++;

            pMeasurements->powerState.batteryMV = (uint32_t)((uint32_t) x & 0x3F) * 10000 / 0x3F;
            pMeasurements->powerState.chargerState = (ChargerState_t) ((x & 0xC0) >> 6);
            pMeasurements->powerState.energyMWh = decodeUint24(ppBuffer);
        }

        // That's all we have but if there were more than 8 it would go as follows
        // // Next bitmap byte
        // // Blah Sensor
        // if (bitMapBytes[1] & 0x01)
        // {
        //     pMeasurements->blahPresent = true;
        //     pMeasurements->blah = decodeUint16 (ppBuffer);
        // }
        // // BlahBlah Sensor
        // if (bitMapBytes[1] & 0x02)
        // {
        //     pMeasurements->blahBlahPresent = true;
        //     pMeasurements->blahBlah = decodeUint16 (ppBuffer);
        // }

        // Having done all that, the pointer must now be at
        // the end point that we established above
        if (*ppBuffer == pBufferAfterEnd)
        {
            success = true;
        }
        else
        {
            // If it isn't, the safest thing is to use the bytesToFollow
            // value from the message to get us to the next thing as
            // we've probably misinterpreted something in the middle
            *ppBuffer = pBufferAfterEnd;
        }

        // Log it
        if ((ppLog != NULL) && (*ppLog != NULL))
        {
            *ppLog += logMeasurements(*ppLog, pLogSize, pMeasurements);
        }
    }

    return success;
}

/// Log the measurements.
uint32_t logMeasurements(char * pBuffer, uint32_t * pBufferSize, Measurements_t * pMeasurements)
{
    uint32_t logSizeAtStart = *pBufferSize;

    if (pBuffer != NULL)
    {
        pBuffer += logDateTime(pBuffer, pBufferSize, pMeasurements->time);
        if (pMeasurements->gnssPositionPresent)
        {
            pBuffer += logPosition(pBuffer, pBufferSize, pMeasurements->gnssPosition.latitude, pMeasurements->gnssPosition.longitude, pMeasurements->gnssPosition.elevation);
        }
        if (pMeasurements->cellIdPresent)
        {
            pBuffer += logTagWithUint32Value(pBuffer, pBufferSize, TAG_CELL_ID, pMeasurements->cellId);
        }
        if (pMeasurements->rsrpPresent)
        {
            pBuffer += logRsrp(pBuffer, pBufferSize, pMeasurements->rsrp.value, pMeasurements->rsrp.isSyncedWithRssi);
        }
        if (pMeasurements->rssiPresent)
        {
            pBuffer += logRssi(pBuffer, pBufferSize, pMeasurements->rssi);
        }
        if (pMeasurements->temperaturePresent)
        {
            pBuffer += logTemperature(pBuffer, pBufferSize, pMeasurements->temperature);
        }
        if (pMeasurements->powerStatePresent)
        {
            pBuffer += logBatteryVoltage(pBuffer, pBufferSize, pMeasurements->powerState.batteryMV);
            pBuffer += logBatteryEnergy(pBuffer, pBufferSize, pMeasurements->powerState.energyMWh);
            pBuffer += logTagWithStringValue(pBuffer, pBufferSize, TAG_CHARGER_STATE, getStringChargerState(pMeasurements->powerState.chargerState));
        }
    }

    return logSizeAtStart - *pBufferSize;
}

/// Encode a MeasurementControlGeneric_t
// This is encoded as follows:
//
// Name                          Size/Pos  Format
// recordLocalOnly               bit 0     bool
// reportImmediately             bit 1     bool
// maxReportingIntervalPresent   bit 2     bool
// useHysteresis                 bit 3     bool
// onlyRecordIfPresent           bit 4     bool
// onlyRecordIfAboveNotBelow     bit 5     bool
// onlyRecordIfAtTransitionOnly  bit 6     bool
// onlyRecordIfIsOneShot         bit 7     bool
// measurementInterval           2 bytes   uint16_t
// maxReportingInterval          2 bytes   uint16_t, only present if maxReportingIntervalPresent is 1
// hysteresisValue               4 bytes   uint32_t, only present if useHysteresis is 1.
// onlyRecordIfValue             4 bytes   uint32_t, only present if onlyRecordIfPresent is 1.
uint32_t encodeMeasurementControlGeneric(char * pBuffer, MeasurementControlGeneric_t * pMeasurementControlGeneric, char ** ppLog, uint32_t * pLogSize)
{
    char *pBufferAtStart = pBuffer;

    // Encode the bit flags
    *pBuffer = 0;
    if (pMeasurementControlGeneric->recordLocalOnly)
    {
        *pBuffer |= 0x01;
    }
    if (pMeasurementControlGeneric->reportImmediately)
    {
        *pBuffer |= 0x02;
    }
    if (pMeasurementControlGeneric->maxReportingInterval != 0)
    {
        *pBuffer |= 0x04;
    }
    if (pMeasurementControlGeneric->useHysteresis)
    {
        *pBuffer |= 0x08;
    }
    if (pMeasurementControlGeneric->onlyRecordIfPresent)
    {
        *pBuffer |= 0x10;
    }
    if (pMeasurementControlGeneric->onlyRecordIfAboveNotBelow)
    {
        *pBuffer |= 0x20;
    }
    if (pMeasurementControlGeneric->onlyRecordIfAtTransitionOnly)
    {
        *pBuffer |= 0x40;
    }
    if (pMeasurementControlGeneric->onlyRecordIfIsOneShot)
    {
        *pBuffer |= 0x80;
    }

    // Move on to the uint16_t's
    pBuffer++;
    if (pMeasurementControlGeneric->measurementInterval > MAX_MEASUREMENT_CONTROL_MEASUREMENT_INTERVAL)
    {
        pMeasurementControlGeneric->measurementInterval = MAX_MEASUREMENT_CONTROL_MEASUREMENT_INTERVAL;
    }
    pBuffer += encodeUint16(pBuffer, pMeasurementControlGeneric->measurementInterval);
    if (pMeasurementControlGeneric->maxReportingInterval != 0)
    {
        if (pMeasurementControlGeneric->maxReportingInterval > MAX_MEASUREMENT_CONTROL_REPORTING_INTERVAL)
        {
            pMeasurementControlGeneric->maxReportingInterval = MAX_MEASUREMENT_CONTROL_REPORTING_INTERVAL;
        }
        pBuffer += encodeUint16(pBuffer, pMeasurementControlGeneric->maxReportingInterval);
    }

    // Move on to the uint32_t's
    if (pMeasurementControlGeneric->useHysteresis)
    {
        pBuffer += encodeUint32(pBuffer, pMeasurementControlGeneric->hysteresisValue);
    }
    if (pMeasurementControlGeneric->onlyRecordIfPresent)
    {
        pBuffer += encodeUint32(pBuffer, pMeasurementControlGeneric->onlyRecordIfValue);
    }

    // Log it
    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logMeasurementControlGeneric(*ppLog, pLogSize, pMeasurementControlGeneric);
    }

    return pBuffer - pBufferAtStart;
}

// Encode each of the MeasurementControlUnion_t types
uint32_t encodeMeasurementControl(char * pBuffer, MeasurementType_t measurementType, MeasurementControlUnion_t * pMeasurementControlUnion, char ** ppLog, uint32_t * pLogSize)
{
    char *pBufferAtStart = pBuffer;
    bool logIt = false;

    // Log it
    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        logIt = true;
        *ppLog += logBeginTagWithStringValue(*ppLog, pLogSize, TAG_MEASUREMENT_TYPE, getStringMeasurementType(measurementType));
    }

    switch (measurementType)
    {
        case MEASUREMENT_GNSS:
        {
            pBuffer += encodeMeasurementControlGeneric(pBuffer, &(pMeasurementControlUnion->gnss), ppLog, pLogSize);
        }
        break;
        case MEASUREMENT_CELL_ID:
        {
            pBuffer += encodeMeasurementControlGeneric(pBuffer, &(pMeasurementControlUnion->cellId), ppLog, pLogSize);
        }
        break;
        case MEASUREMENT_RSRP:
        {
            pBuffer += encodeMeasurementControlGeneric(pBuffer, &(pMeasurementControlUnion->rsrp), ppLog, pLogSize);
        }
        break;
        case MEASUREMENT_RSSI:
        {
            pBuffer += encodeMeasurementControlGeneric(pBuffer, &(pMeasurementControlUnion->rssi), ppLog, pLogSize);
        }
        break;
        case MEASUREMENT_TEMPERATURE:
        {
            pBuffer += encodeMeasurementControlGeneric(pBuffer, &(pMeasurementControlUnion->temperature), ppLog, pLogSize);
        }
        break;
        case MEASUREMENT_POWER_STATE:
        {
            *pBuffer = pMeasurementControlUnion->powerState.powerStateChargerStateReportImmediately;
            pBuffer++;
            if (logIt)
            {
                *ppLog += logBeginTag(*ppLog, pLogSize, TAG_CHARGER_STATE_MEASUREMENT_CONTROL);
                *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_REPORT_IMMEDIATELY, getStringBoolean(pMeasurementControlUnion->powerState.powerStateChargerStateReportImmediately));
                *ppLog += logEndTag(*ppLog, pLogSize, TAG_CHARGER_STATE_MEASUREMENT_CONTROL);
                *ppLog += logBeginTag(*ppLog, pLogSize, TAG_VOLTAGE_MEASUREMENT_CONTROL);
            }
            pBuffer += encodeMeasurementControlGeneric(pBuffer, &(pMeasurementControlUnion->powerState.powerStateBatteryVoltage), ppLog, pLogSize);
            if (logIt)
            {
                *ppLog += logEndTag(*ppLog, pLogSize, TAG_VOLTAGE_MEASUREMENT_CONTROL);
                *ppLog += logBeginTag(*ppLog, pLogSize, TAG_ENERGY_MEASUREMENT_CONTROL);
            }
            pBuffer += encodeMeasurementControlGeneric(pBuffer, &(pMeasurementControlUnion->powerState.powerStateBatteryEnergy), ppLog, pLogSize);
            if (logIt)
            {
                *ppLog += logEndTag(*ppLog, pLogSize, TAG_ENERGY_MEASUREMENT_CONTROL);
            }
        }
        break;
        default:
        break;
    }

    if (logIt)
    {
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MEASUREMENT_TYPE);
    }

    return pBuffer - pBufferAtStart;
}

/// Decode a MeasurementControlGeneric_t
bool decodeMeasurementControlGeneric(const char ** ppBuffer, MeasurementControlGeneric_t * pMeasurementControlGeneric, char ** ppLog, uint32_t * pLogSize)
{
    bool maxReportingIntervalPresent = false;

    // Decode the bit flags
    pMeasurementControlGeneric->recordLocalOnly = false;
    if (**ppBuffer & 0x01)
    {
        pMeasurementControlGeneric->recordLocalOnly = true;
    }
    pMeasurementControlGeneric->reportImmediately = false;
    if (**ppBuffer & 0x02)
    {
        pMeasurementControlGeneric->reportImmediately = true;
    }
    if (**ppBuffer & 0x04)
    {
        maxReportingIntervalPresent = true;
    }
    pMeasurementControlGeneric->useHysteresis = false;
    if (**ppBuffer & 0x08)
    {
        pMeasurementControlGeneric->useHysteresis = true;
    }
    pMeasurementControlGeneric->onlyRecordIfPresent = false;
    if (**ppBuffer & 0x10)
    {
        pMeasurementControlGeneric->onlyRecordIfPresent = true;
    }
    pMeasurementControlGeneric->onlyRecordIfAboveNotBelow = false;
    if (**ppBuffer & 0x20)
    {
        pMeasurementControlGeneric->onlyRecordIfAboveNotBelow = true;
    }
    pMeasurementControlGeneric->onlyRecordIfAtTransitionOnly = false;
    if (**ppBuffer & 0x40)
    {
        pMeasurementControlGeneric->onlyRecordIfAtTransitionOnly = true;
    }
    pMeasurementControlGeneric->onlyRecordIfIsOneShot = false;
    if (**ppBuffer & 0x80)
    {
        pMeasurementControlGeneric->onlyRecordIfIsOneShot = true;
    }

    // Move on to the uint16_t's
    (*ppBuffer)++;
    pMeasurementControlGeneric->measurementInterval = decodeUint16(ppBuffer);
    if (maxReportingIntervalPresent)
    {
        pMeasurementControlGeneric->maxReportingInterval = decodeUint16(ppBuffer);
    }

    // Move on to the uint32_t's
    if (pMeasurementControlGeneric->useHysteresis)
    {
        pMeasurementControlGeneric->hysteresisValue = decodeUint32(ppBuffer);
    }
    if (pMeasurementControlGeneric->onlyRecordIfPresent)
    {
        pMeasurementControlGeneric->onlyRecordIfValue = decodeUint32(ppBuffer);
    }

    // Log it
    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logMeasurementControlGeneric(*ppLog, pLogSize, pMeasurementControlGeneric);
    }

    return true;
}

/// Decode each of the MeasurementControlUnion_t types
bool decodeMeasurementControl(const char ** ppBuffer, MeasurementType_t measurementType, MeasurementControlUnion_t * pMeasurementControlUnion, char ** ppLog, uint32_t * pLogSize)
{
    bool success = false;
    bool logIt = false;

    // Log it
    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        logIt = true;
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MEASUREMENT_TYPE, getStringMeasurementType(measurementType));
    }

    memset(pMeasurementControlUnion, 0, sizeof(*pMeasurementControlUnion));

    switch (measurementType)
    {
        case MEASUREMENT_GNSS:
        {
            if (logIt)
            {
                *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MEASUREMENT_TYPE, getStringMeasurementType(measurementType));
            }
            success = decodeMeasurementControlGeneric(ppBuffer, &(pMeasurementControlUnion->gnss), ppLog, pLogSize);
        }
        break;
        case MEASUREMENT_CELL_ID:
        {
            if (logIt)
            {
                *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MEASUREMENT_TYPE, getStringMeasurementType(measurementType));
            }
            success = decodeMeasurementControlGeneric(ppBuffer, &(pMeasurementControlUnion->cellId), ppLog, pLogSize);
        }
        break;
        case MEASUREMENT_RSRP:
        {
            if (logIt)
            {
                *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MEASUREMENT_TYPE, getStringMeasurementType(measurementType));
            }
            success = decodeMeasurementControlGeneric(ppBuffer, &(pMeasurementControlUnion->rsrp), ppLog, pLogSize);
        }
        break;
        case MEASUREMENT_RSSI:
        {
            if (logIt)
            {
                *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MEASUREMENT_TYPE, getStringMeasurementType(measurementType));
            }
            success = decodeMeasurementControlGeneric(ppBuffer, &(pMeasurementControlUnion->rssi), ppLog, pLogSize);
        }
        break;
        case MEASUREMENT_TEMPERATURE:
        {
            if (logIt)
            {
                *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MEASUREMENT_TYPE, getStringMeasurementType(measurementType));
            }
            success = decodeMeasurementControlGeneric(ppBuffer, &(pMeasurementControlUnion->temperature), ppLog, pLogSize);
        }
        break;
        case MEASUREMENT_POWER_STATE:
        {
            pMeasurementControlUnion->powerState.powerStateChargerStateReportImmediately = **ppBuffer;
            (*ppBuffer)++;
            if (logIt)
            {
                *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MEASUREMENT_CONTROL_POWER_STATE);
                *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MEASUREMENT_TYPE, getStringMeasurementType(measurementType));
                *ppLog += logBeginTag(*ppLog, pLogSize, TAG_CHARGER_STATE_MEASUREMENT_CONTROL);
                *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_REPORT_IMMEDIATELY, getStringBoolean(pMeasurementControlUnion->powerState.powerStateChargerStateReportImmediately));
                *ppLog += logEndTag(*ppLog, pLogSize, TAG_CHARGER_STATE_MEASUREMENT_CONTROL);
                *ppLog += logBeginTag(*ppLog, pLogSize, TAG_VOLTAGE_MEASUREMENT_CONTROL);
            }
            success = decodeMeasurementControlGeneric(ppBuffer, &(pMeasurementControlUnion->powerState.powerStateBatteryVoltage), ppLog, pLogSize);
            if (logIt)
            {
                *ppLog += logEndTag(*ppLog, pLogSize, TAG_VOLTAGE_MEASUREMENT_CONTROL);
            }
            if (success)
            {
                if (logIt)
                {
                    *ppLog += logBeginTag(*ppLog, pLogSize, TAG_ENERGY_MEASUREMENT_CONTROL);
                }
                success = decodeMeasurementControlGeneric(ppBuffer, &(pMeasurementControlUnion->powerState.powerStateBatteryEnergy), ppLog, pLogSize);
                if (logIt)
                {
                    *ppLog += logEndTag(*ppLog, pLogSize, TAG_ENERGY_MEASUREMENT_CONTROL);
                }
            }
            if (logIt)
            {
                *ppLog += logEndTag(*ppLog, pLogSize, TAG_MEASUREMENT_CONTROL_POWER_STATE);
            }
        }
        break;
        default:
        break;
    }

    return success;
}

// Log the generic measurement control structure
uint32_t logMeasurementControlGeneric(char * pBuffer, uint32_t * pBufferSize, MeasurementControlGeneric_t * pMeasurementControlGeneric)
{
    uint32_t logSizeAtStart = *pBufferSize;

    if (pBuffer != NULL)
    {
        pBuffer += logBeginTag(pBuffer, pBufferSize, TAG_MEASUREMENT_CONTROL_GENERIC);
        pBuffer += logTagWithUint32Value(pBuffer, pBufferSize, "HeartbeatInterval", pMeasurementControlGeneric->measurementInterval);
        pBuffer += logTagWithPresenceAndUint32Value(pBuffer, pBufferSize, TAG_MAX_REPORTING_INTERVAL, (pMeasurementControlGeneric->maxReportingInterval != 0),
                pMeasurementControlGeneric->maxReportingInterval);
        pBuffer += logTagWithPresenceAndUint32Value(pBuffer, pBufferSize, TAG_HYSTERESIS, pMeasurementControlGeneric->useHysteresis, pMeasurementControlGeneric->hysteresisValue);
        pBuffer += logOnlyRecordIf(pBuffer, pBufferSize, pMeasurementControlGeneric->onlyRecordIfPresent, pMeasurementControlGeneric->onlyRecordIfValue,
                pMeasurementControlGeneric->onlyRecordIfAboveNotBelow, pMeasurementControlGeneric->onlyRecordIfAtTransitionOnly, pMeasurementControlGeneric->onlyRecordIfIsOneShot);
        pBuffer += logTagWithStringValue(pBuffer, pBufferSize, TAG_REPORT_IMMEDIATELY, getStringBoolean(pMeasurementControlGeneric->reportImmediately));
        pBuffer += logTagWithStringValue(pBuffer, pBufferSize, "RecordLocalOnly", getStringBoolean(pMeasurementControlGeneric->recordLocalOnly));
        pBuffer += logEndTag(pBuffer, pBufferSize, TAG_MEASUREMENT_CONTROL_GENERIC);
    }

    return logSizeAtStart - *pBufferSize;
}

/// Decode a TrafficTestModeRuleBreakerDatagram, checking if it's completely
// full of the expected fill
bool decodeTrafficTestModeRuleBreakerDatagram(const char ** ppBuffer, bool isDownlink, TrafficTestModeRuleBreakerDatagram_t * pSpec, char **ppLog, uint32_t * pLogSize)
{
    bool success = false;
    bool badChecksum = false;
    const char * pStart = *ppBuffer;
    const char * pEnd = pStart;
    const char * pTag = TAG_MSG_UL;
    char actualFill = 0;

    if (pSpec->length > TRAFFIC_TEST_MODE_RULE_BREAKER_MAX_LENGTH)
    {
        pSpec->length = TRAFFIC_TEST_MODE_RULE_BREAKER_MAX_LENGTH;
    }

    if (pSpec->length > 0)
    {
        pEnd = *ppBuffer + pSpec->length - CHECKSUM_SIZE - MSG_ID_SIZE;
    }
    actualFill = **ppBuffer;

    while ((*ppBuffer < pEnd) && (**ppBuffer == pSpec->fill))
    {
        (*ppBuffer)++;
    }

    // This is a bit of a cheat but I just can't be arsed to rearchitect
    // this call now.  The message ID should be included in the calculation
    // of the checksum but it's at the location before pStart so access it
    // by looking just before the start pointer
    if (calculateChecksum(pStart - 1, *ppBuffer - (pStart - 1)) != **ppBuffer)
    {
        badChecksum = true;
    }
    (*ppBuffer)++;

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        if (isDownlink)
        {
            pTag = TAG_MSG_DL;
        }

        *ppLog += logBeginTag(*ppLog, pLogSize, pTag);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestRuleBreakerDatagram");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTrafficTestModeRuleBreakerDatagram(*ppLog, pLogSize, actualFill, *ppBuffer - pStart);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, *ppBuffer - pStart);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(!badChecksum));
        *ppLog += logEndTag(*ppLog, pLogSize, pTag);
    }

    pSpec->length = *ppBuffer - pStart + MSG_ID_SIZE;  // No need to add checksum size as it's within *ppBuffer
    pSpec->fill = actualFill;

    if ((*ppBuffer == (pEnd + CHECKSUM_SIZE)) && !badChecksum)
    {
        success = true;
    }
    else
    {
        // Move the pointer to the end so as not to get out of step
        // and pass back the fill value we received so that the
        // counter can match step
        *ppBuffer = pEnd;
    }

    return success;
}

/// Limit an int to N bits with correct sign representation
static __inline uint32_t limitInt(int32_t number, uint8_t bits)
{
    uint32_t newNumber = 0;

    if (bits > 0)
    {
        newNumber = (uint32_t) number;

        // Deal with sign
        if (newNumber & (1 << ((sizeof(newNumber) * 8) - 1)))
        {
            newNumber |= 1 << (bits - 1);
        }

        // Mask off the unwanted part
        newNumber &= (1 << bits) - 1;
    }

    return newNumber;
}

/// Sign-extend a number of N bits (held inside an uin32_t) to an int32_t
// This from https://graphics.stanford.edu/~seander/bithacks.html#VariableSignExtend
static __inline int32_t extendInt(uint32_t number, uint8_t bits)
{
    uint32_t mask = 1 << (bits - 1);

    // Sign extend appropriately
    number = (number ^ mask) - mask;

    return number;
}

/// Calculate the checksum. This is not meant to be fool-proof, it is
// simply a way of determining that messages have been generated by this
// protocol from random messages that seem to land at our door for no very
// good reason and which, if interpreted as real messages, could cause all
// sorts of chaos.
static char calculateChecksum(const char * pBuffer, uint32_t bufferLength)
{
    char checkSum = 0;
    uint32_t x;

    for (x = 0; x < bufferLength; x++)
    {
        checkSum += *pBuffer;
        pBuffer++;
    }

    // For single byte buffers, which we have a few of, this will just look
    // like the byte repeated, which isn't very random, so add 128 to it to
    // make it more distinct
    checkSum += 128;

    return checkSum;
}

// ----------------------------------------------------------------
// PUBLIC LOGGING FUNCTIONS
// ----------------------------------------------------------------

/// Return a hex string equivalent of an uint32_t
const char * getHexAsString(uint32_t value)
{
    char byteArray[4];

    byteArray[0] = value >> 24;
    byteArray[1] = value >> 16;
    byteArray[2] = value >> 8;
    byteArray[3] = value;

    bytesToHexString(&(byteArray[0]), sizeof(byteArray), &(gValueAsHexString[2]), sizeof(gValueAsHexString) - 2); // 2 to leave the 0x on the front

    return &(gValueAsHexString[0]);
}

/// Return a string representing a Boolean(for logging).
const char * getStringBoolean(bool boolValue)
{
    const char * strValue = VALUE_UNKNOWN_STRING;

    if (boolValue < (sizeof(gStringBoolean) / sizeof(gStringBoolean[0])))
    {
        strValue = gStringBoolean[boolValue];
    }
    else
    {
        strValue = getHexAsString((uint32_t) boolValue);
    }

    return strValue;
}

/// Return a string representing a wakeup code (for logging).
const char * getStringWakeUpCode(WakeUpCode_t wakeupCode)
{
    const char * strValue = VALUE_UNKNOWN_STRING;

    if (wakeupCode < sizeof(gStringWakeUpCode) / sizeof(gStringWakeUpCode[0]))
    {
        strValue = gStringWakeUpCode[wakeupCode];
    }
    else
    {
        strValue = getHexAsString((uint32_t) wakeupCode);
    }

    return strValue;
}

/// Return a string representing the mode (for logging).
const char * getStringMode(Mode_t mode)
{
    const char * strValue = VALUE_UNKNOWN_STRING;

    if (mode < sizeof(gStringMode) / sizeof(gStringMode[0]))
    {
        strValue = gStringMode[mode];
    }
    else
    {
        strValue = getHexAsString((uint32_t) mode);
    }

    return strValue;
}

/// Return a string representing the thing that the
// time was set by (for logging).
const char * getStringTimeSetBy(TimeSetBy_t timeSetBy)
{
    const char * strValue = VALUE_UNKNOWN_STRING;

    if (timeSetBy < sizeof(gStringTimeSetBy) / sizeof(gStringTimeSetBy[0]))
    {
        strValue = gStringMode[timeSetBy];
    }
    else
    {
        strValue = getHexAsString((uint32_t) timeSetBy);
    }

    return strValue;
}

/// Return a string representing the amount of energy left (for logging).
const char * getStringEnergyLeft(EnergyLeft_t energyLeft)
{
    const char * strValue = VALUE_UNKNOWN_STRING;

    if (energyLeft < sizeof(gStringEnergyLeft) / sizeof(gStringEnergyLeft[0]))
    {
        strValue = gStringEnergyLeft[energyLeft];
    }
    else
    {
        strValue = getHexAsString((uint32_t) energyLeft);
    }

    return strValue;
}

/// Return a string representing the amount of disk space left (for logging).
const char * getStringDiskSpaceLeft(DiskSpaceLeft_t diskSpaceLeft)
{
    const char * strValue = VALUE_UNKNOWN_STRING;

    if (diskSpaceLeft < sizeof(gStringDiskSpaceLeft) / sizeof(gStringDiskSpaceLeft[0]))
    {
        strValue = gStringDiskSpaceLeft[diskSpaceLeft];
    }
    else
    {
        strValue = getHexAsString((uint32_t) diskSpaceLeft);
    }

    return strValue;
}

/// Return a string representing the measurement type (for logging).
const char * getStringMeasurementType(MeasurementType_t type)
{
    const char * strValue = VALUE_UNKNOWN_STRING;

    if (type < sizeof(gStringMeasurementType) / sizeof(gStringMeasurementType[0]))
    {
        strValue = gStringMeasurementType[type];
    }
    else
    {
        strValue = getHexAsString((uint32_t) type);
    }

    return strValue;
}

/// Return a string representing the charger state (for logging).
const char * getStringChargerState(ChargerState_t chargerState)
{
    const char * strValue = VALUE_UNKNOWN_STRING;

    if (chargerState < sizeof(gStringChargerState) / sizeof(gStringChargerState[0]))
    {
        strValue = gStringChargerState[chargerState];
    }
    else
    {
        strValue = getHexAsString((uint32_t) chargerState);
    }

    return strValue;
}

/// Calculate the new buffer size after an snprintf();
uint32_t calcBytesUsed(uint32_t *pBufferSize, uint32_t bytesUsed)
{
    if (bytesUsed < *pBufferSize)
    {
        *pBufferSize -= bytesUsed;
    }
    else
    {
        bytesUsed = *pBufferSize;
        *pBufferSize = 0;
    }

    return bytesUsed;
}

/// Log a start tag
uint32_t logBeginTag(char * pBuffer, uint32_t *pBufferSize, const char *pTag)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<%s>", pTag);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log a start tag with a string value
uint32_t logBeginTagWithStringValue(char * pBuffer, uint32_t *pBufferSize, const char *pTag, const char * pValue)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<%s Value=\"%s\">", pTag, pValue);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log an end tag
uint32_t logEndTag(char * pBuffer, uint32_t *pBufferSize, const char *pTag)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "</%s>", pTag);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log a flag, no actual values in it
uint32_t logFlagTag(char * pBuffer, uint32_t *pBufferSize, const char *pTag)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<%s />", pTag);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log a tag with a simple string value
uint32_t logTagWithStringValue(char * pBuffer, uint32_t *pBufferSize, const char * pTag, const char * pValue)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<%s Value=\"%s\" />", pTag, pValue);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log a tag with a simple uint32_t value
uint32_t logTagWithUint32Value(char * pBuffer, uint32_t *pBufferSize, const char * pTag, uint32_t value)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<%s Value=\"%ld\" />", pTag, value);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log a tag with a uint32_t value and a presence flag
uint32_t logTagWithPresenceAndUint32Value(char * pBuffer, uint32_t *pBufferSize, const char * pTag, bool present, uint32_t value)
{
    uint32_t bytesUsed;

    if (present)
    {
        bytesUsed = snprintf(pBuffer, *pBufferSize, "<%s Present=\"%s\" Value=\"%ld\" />", pTag, getStringBoolean(present), value);
    }
    else
    {
        bytesUsed = snprintf(pBuffer, *pBufferSize, "<%s Present=\"%s\" />", pTag, getStringBoolean(present));
    }

    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the contents of a TransparentDatagram
uint32_t logTransparentData(char * pBuffer, uint32_t *pBufferSize, uint32_t datagramSize, const char * pContents)
{
    uint32_t bytesUsed;
    uint16_t stringLen;

    stringLen = bytesToHexString(pContents, MAX_DATAGRAM_SIZE_RAW - 1, &(gLogBuffer[0]), sizeof(gLogBuffer));

    bytesUsed = snprintf(pBuffer, *pBufferSize, "<TransparentData Length=\"%d\" HexContents=\"%.*s\" />",
    MAX_DATAGRAM_SIZE_RAW - 1, stringLen, &(gLogBuffer[0]));

    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the DateTime
uint32_t logDateTime(char * pBuffer, uint32_t *pBufferSize, uint32_t timeValue)
{
    uint32_t bytesUsed;
    char timeString[32];
    struct tm * pT;
    time_t t = timeValue;

    pT = gmtime(&t);
    strftime(&(timeString[0]), sizeof(timeString), "%Y-%m-%d %X", pT);
    bytesUsed = snprintf(pBuffer, *pBufferSize, "<DateTime Value=\"%s\" TimeZone=\"UTC\" />", &(timeString[0]));

    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the Heartbeat
uint32_t logHeartbeat(char * pBuffer, uint32_t *pBufferSize, uint32_t heartbeatValue, bool snapToRtc)
{
    uint32_t bytesUsed;
    char timeString[32];
    struct tm * pT;
    time_t t = heartbeatValue;

    bytesUsed = snprintf(pBuffer, *pBufferSize, "<Heartbeat>");
    bytesUsed = calcBytesUsed(pBufferSize, bytesUsed);

    if (snapToRtc)
    {
        pT = gmtime(&t);
        strftime(&(timeString[0]), sizeof(timeString), "%M:%S", pT);
        bytesUsed += snprintf((pBuffer + bytesUsed), *pBufferSize, "<SnapToRtc Value=\"%s\" Units=\"Mins:Secs\" />", &(timeString[0]));
        bytesUsed = calcBytesUsed(pBufferSize, bytesUsed);
    }
    else
    {
        bytesUsed += snprintf((pBuffer + bytesUsed), *pBufferSize, "<Interval Value=\"%ld\" Units=\"Seconds\" />", heartbeatValue);
        bytesUsed = calcBytesUsed(pBufferSize, bytesUsed);
    }

    bytesUsed += snprintf((pBuffer + bytesUsed), *pBufferSize, "</Heartbeat>");
    bytesUsed = calcBytesUsed(pBufferSize, bytesUsed);

    return bytesUsed;
}

/// Log the RSSI
uint32_t logRssi(char * pBuffer, uint32_t *pBufferSize, Rssi_t rssi)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<Rssi Value=\"%d\" Units=\"dBm\" />", (double) rssi / 10);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the RSRP
uint32_t logRsrp(char * pBuffer, uint32_t *pBufferSize, Rssi_t rsrp, bool isSyncedWithRssi)
{
	// TODO: for reasons I have not yet determined, if the string value and the float
	// value are included in this log point I get a segmentation fault out of snprintf().
	// I've printf'ed all the values around here and everything seems in order, the string
	// isn't actually getting too large for the supplied buffer.  If I don't do the float
	// the problem goes away.  Don't understand.  Be afraid, be moderately afraid.
    //uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<Rrsp Value=\"%.1f\" Units=\"dBm\" IsSyncedWithRssi=\"%s\" />", (double) rsrp / 10, getStringBoolean(isSyncedWithRssi));
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<Rrsp Value=\"%d.%d\" Units=\"dBm\" IsSyncedWithRssi=\"%s\" />", rsrp / 10, rsrp % 10, getStringBoolean(isSyncedWithRssi));
	return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the temperature
uint32_t logTemperature(char * pBuffer, uint32_t *pBufferSize, Temperature_t temperature)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<Temperature Value=\"%d\" Units=\"Celsius\" />", temperature);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the battery voltage
uint32_t logBatteryVoltage(char * pBuffer, uint32_t *pBufferSize, uint16_t batteryVoltage)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<BatteryVoltage Value=\"%d\" Units=\"mV\" />", batteryVoltage);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the battery energy left
uint32_t logBatteryEnergy(char * pBuffer, uint32_t *pBufferSize, uint32_t batteryEnergy)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<BatteryEnergy Value=\"%ld\" Units=\"mWh\" />", batteryEnergy);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the GNSS position
uint32_t logPosition(char * pBuffer, uint32_t *pBufferSize, uint32_t latitude, uint32_t longitude, uint32_t elevation)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize,
            "<Position Latitude=\"%ld\" LatitudeUnits=\"Degrees/1000\" Longitude=\"%ld\" LongitudeUnits=\"Degrees/1000\" Elevation=\"%ld\" ElevationUnits=\"Metres\" />", latitude, longitude,
            elevation);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log an "OnlyRecordIf" control
uint32_t logOnlyRecordIf(char * pBuffer, uint32_t *pBufferSize, bool present, int32_t value, bool aboveNotBelow, bool atTransitionOnly, bool isOneShot)
{
    uint32_t bytesUsed;

    if (present)
    {
        bytesUsed = snprintf(pBuffer, *pBufferSize, "<OnlyRecordIf Present=\"%s\" Value=\"%ld\" AboveNotBelow=\"%s\" AtTransitionOnly=\"%s\" OneShot=\"%s\" />", getStringBoolean(present), value,
                getStringBoolean(aboveNotBelow), getStringBoolean(atTransitionOnly), getStringBoolean(isOneShot));
    }
    else
    {
        bytesUsed = snprintf(pBuffer, *pBufferSize, "<OnlyRecordIf Present=\"%s\" />", getStringBoolean(present));
    }

    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the Traffic Report transmit values
uint32_t logTrafficReportUl(char * pBuffer, uint32_t *pBufferSize, uint32_t datagrams, uint32_t bytes)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<TotalUplink Datagrams=\"%ld\" Bytes=\"%ld\" />", datagrams, bytes);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the Traffic Report receive values
uint32_t logTrafficReportDl(char * pBuffer, uint32_t *pBufferSize, uint32_t datagrams, uint32_t bytes, uint32_t badChecksum)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<TotalDownlink Datagrams=\"%ld\" Bytes=\"%ld\" DatagramsBadChecksum=\"%ld\" />", datagrams, bytes, badChecksum);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the Traffic Test mode parameters for transmit from the UTM.
uint32_t logTrafficTestModeParametersUl(char * pBuffer, uint32_t *pBufferSize, uint32_t count, uint32_t length, bool noReportsDuringTest)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<DatagramsUplink Count=\"%ld\" Length=\"%ld\" NoReportsDuringTest=\"%s\" />", count, length, getStringBoolean(noReportsDuringTest));
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the Traffic Test mode parameters for receive from the network.
uint32_t logTrafficTestModeParametersDl(char * pBuffer, uint32_t *pBufferSize, uint32_t count, uint32_t length)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<DatagramsDownlink Count=\"%ld\" Length=\"%ld\" />", count, length);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the Traffic Test mode report transmit values from the UTM
uint32_t logTrafficTestModeReportUl(char * pBuffer, uint32_t *pBufferSize, uint32_t datagramCount, uint32_t bytes)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<TestModeUplink Count=\"%ld\" Bytes=\"%ld\" />", datagramCount, bytes);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the Traffic Test mode report receive values from the UTM
uint32_t logTrafficTestModeReportDl(char * pBuffer, uint32_t *pBufferSize, uint32_t datagramCount, uint32_t bytes, uint32_t outOfOrder, uint32_t bad, uint32_t missed)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<TestModeDownlink Count=\"%ld\" Bytes=\"%ld\" OutOfOrder=\"%ld\" Bad=\"%ld\" Missed=\"%ld\" />", datagramCount, bytes, outOfOrder, bad,
            missed);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log a Traffic Test mode datagram
uint32_t logTrafficTestModeRuleBreakerDatagram(char * pBuffer, uint32_t *pBufferSize, uint8_t fill, uint32_t length)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<TestModeDatagram FillDec=\"%d\" FillHex=\"%02x\" Length=\"%ld\" />", fill, fill, length);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log the Activity report values from the UTM
uint32_t logActivityReport(char * pBuffer, uint32_t *pBufferSize, uint32_t totalTransmitMilliseconds, uint32_t totalReceiveMilliseconds, uint32_t upTimeSeconds)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<Activity Transmit=\"%ld\" Receive=\"%ld\" Units=\"Milliseconds\" /><UpTime Value=\"%ld\" Units=\"Seconds\" />", totalTransmitMilliseconds,
            totalReceiveMilliseconds, upTimeSeconds);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log some UL RF parameters from the UTM
uint32_t logUlRf(char * pBuffer, uint32_t *pBufferSize, bool txPowerDbmPresent, int8_t txPowerDbm, bool mcsPresent, uint8_t mcs)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<UlRf TxPowerPresent=\"%s\" TxPower=\"%d\" Units=\"dBm\" McsPresent=\"%s\", Mcs=\"%d\" />",
            getStringBoolean(txPowerDbmPresent), txPowerDbm, getStringBoolean(mcsPresent), mcs);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log some DL RF parameters from the UTM
uint32_t logDlRf(char * pBuffer, uint32_t *pBufferSize, bool mcsPresent, uint8_t mcs)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<DlRf McsPresent=\"%s\", Mcs=\"%d\" />", getStringBoolean(mcsPresent), mcs);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

/// Log a DebugInd string
uint32_t logDebugInd(char * pBuffer, uint32_t *pBufferSize, char * pString, uint8_t length)
{
    uint32_t bytesUsed = snprintf(pBuffer, *pBufferSize, "<String Value=\"%.*s\" />", length, pString);
    return calcBytesUsed(pBufferSize, bytesUsed);
}

// ----------------------------------------------------------------
// MESSAGE ENCODING FUNCTIONS
// ----------------------------------------------------------------

uint32_t encodeTransparentDatagram(char * pBuffer, TransparentDatagram_t * pDatagram, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding TransparentDatagram, ID 0x%.2x.\r\n", TRANSPARENT_DATAGRAM_ID);
    pBuffer[numBytesEncoded] = TRANSPARENT_DATAGRAM_ID;
    numBytesEncoded++;
    pBuffer++;
    memcpy(pBuffer, pDatagram->contents, sizeof(pDatagram->contents));
    numBytesEncoded += sizeof(pDatagram->contents);
    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_BI);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TransparentDatagram");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTransparentData(*ppLog, pLogSize, sizeof(pDatagram->contents), pBuffer);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_BI);
    }

    return numBytesEncoded;
}

uint32_t encodePingReqMsg(char * pBuffer, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding PingReqMsg, ID 0x%.2x.\r\n", PING_REQ_MSG_ID);
    pBuffer[numBytesEncoded] = PING_REQ_MSG_ID;
    numBytesEncoded++;
    // Empty body
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_BI);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "PingReqMsg");
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_BI);
    }

    return numBytesEncoded;
}

uint32_t encodePingCnfMsg(char * pBuffer, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding PingCnfMsg, ID 0x%.2x.\r\n", PING_CNF_MSG_ID);
    pBuffer[numBytesEncoded] = PING_CNF_MSG_ID;
    numBytesEncoded++;
    // Empty body
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_BI);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "PingCnfMsg");
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_BI);
    }

    return numBytesEncoded;
}

uint32_t encodeInitIndUlMsg(char * pBuffer, InitIndUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding InitIndUlMsg, ID 0x%.2x.\r\n", INIT_IND_UL_MSG);
    pBuffer[numBytesEncoded] = INIT_IND_UL_MSG;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = (uint8_t) pMsg->wakeUpCode;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = REVISION_LEVEL;
    numBytesEncoded++;

    pBuffer[numBytesEncoded] = 0;

    if (pMsg->sdCardNotRequired)
    {
        pBuffer[numBytesEncoded] |= 0x01;
    }
    if (pMsg->disableModemDebug)
    {
        pBuffer[numBytesEncoded] |= 0x02;
    }
    if (pMsg->disableButton)
    {
        pBuffer[numBytesEncoded] |= 0x04;
    }
    if (pMsg->disableServerPing)
    {
        pBuffer[numBytesEncoded] |= 0x08;
    }
    numBytesEncoded++;

    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "InitIndUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "WakeupCode", getStringWakeUpCode(pMsg->wakeUpCode));
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, "RevisionLevel", REVISION_LEVEL);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, "SdCardRequired", !(pMsg->sdCardNotRequired));
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "DisableModemDebug", getStringBoolean(pMsg->disableModemDebug));
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "DisableButton", getStringBoolean(pMsg->disableButton));
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "DisableServerPing", getStringBoolean(pMsg->disableServerPing));
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

/// The PollInd contains the mode, plus approximate indications of
// the battery left and disk space left, coded as follows:
//
// Name                          Bits      Format
// mode                          bits 0-2  Mode_t
// energyLeft                    bits 3-5  EnergyLeft_t
// diskSpaceLeft                 bits 6-7  DiskSpaceLeft_t

uint32_t encodePollIndUlMsg(char * pBuffer, PollIndUlMsg_t *pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding PollIndUlMsg, ID 0x%.2x.\r\n", POLL_IND_UL_MSG);
    pBuffer[numBytesEncoded] = POLL_IND_UL_MSG;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = 0;
    pBuffer[numBytesEncoded] = pMsg->mode & 0x07;
    pBuffer[numBytesEncoded] |= (pMsg->energyLeft & 0x07) << 3;
    pBuffer[numBytesEncoded] |= (pMsg->diskSpaceLeft & 0x03) << 6;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "PollIndUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MODE, getStringMode(pMsg->mode));
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_ENERGY_LEFT, getStringEnergyLeft(pMsg->energyLeft));
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_DISK_SPACE_LEFT, getStringDiskSpaceLeft(pMsg->diskSpaceLeft));
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeRebootReqDlMsg(char * pBuffer, RebootReqDlMsg_t *pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding RebootReqDlMsg, ID 0x%.2x.\r\n", REBOOT_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = REBOOT_REQ_DL_MSG;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = 0;
    if (pMsg->sdCardNotRequired)
    {
        pBuffer[numBytesEncoded] |= 0x01;
    }
    if (pMsg->disableModemDebug)
    {
        pBuffer[numBytesEncoded] |= 0x02;
    }
    if (pMsg->disableButton)
    {
        pBuffer[numBytesEncoded] |= 0x04;
    }
    if (pMsg->disableServerPing)
    {
        pBuffer[numBytesEncoded] |= 0x08;
    }
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "RebootReqDlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "SdCardRequired", getStringBoolean(!pMsg->sdCardNotRequired));
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "DisableModemDebug", getStringBoolean(pMsg->disableModemDebug));
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "DisableButton", getStringBoolean(pMsg->disableButton));
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "DisableServerPing", getStringBoolean(pMsg->disableServerPing));
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeDateTimeSetReqDlMsg(char * pBuffer, DateTimeSetReqDlMsg_t *pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding DateTimeSetReqDlMsg, ID 0x%.2x.\r\n", DATE_TIME_SET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = DATE_TIME_SET_REQ_DL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->time);
    numBytesEncoded += encodeBool(&(pBuffer[numBytesEncoded]), pMsg->setDateOnly);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "DateTimeSetReqDlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logDateTime(*ppLog, pLogSize, pMsg->time);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_SET_DATE_ONLY, getStringBoolean(pMsg->setDateOnly));
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeDateTimeSetCnfUlMsg(char * pBuffer, DateTimeSetCnfUlMsg_t *pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding DateTimeSetCnfUlMsg, ID 0x%.2x.\r\n", DATE_TIME_SET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = DATE_TIME_SET_CNF_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->time);
    pBuffer[numBytesEncoded] = pMsg->setBy;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "DateTimeSetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logDateTime(*ppLog, pLogSize, pMsg->time);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_TIME_SET_BY, getStringTimeSetBy(pMsg->setBy));
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeDateTimeGetReqDlMsg(char * pBuffer, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding DateTimeGetReqDlMsg, ID 0x%.2x.\r\n", DATE_TIME_GET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = DATE_TIME_GET_REQ_DL_MSG;
    numBytesEncoded++;
    // Empty body
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "DateTimeGetReqDlMsg");
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeDateTimeGetCnfUlMsg(char * pBuffer, DateTimeGetCnfUlMsg_t *pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding DateTimeGetCnfUlMsg, ID 0x%.2x.\r\n", DATE_TIME_GET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = DATE_TIME_GET_CNF_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->time);
    pBuffer[numBytesEncoded] = pMsg->setBy;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "DateTimeGetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logDateTime(*ppLog, pLogSize, pMsg->time);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_TIME_SET_BY, getStringTimeSetBy(pMsg->setBy));
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeDateTimeIndUlMsg(char * pBuffer, DateTimeIndUlMsg_t *pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding DateTimeIndUlMsg, ID 0x%.2x.\r\n", DATE_TIME_IND_UL_MSG);
    pBuffer[numBytesEncoded] = DATE_TIME_IND_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->time);
    pBuffer[numBytesEncoded] = pMsg->setBy;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "DateTimeIndUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logDateTime(*ppLog, pLogSize, pMsg->time);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_TIME_SET_BY, getStringTimeSetBy(pMsg->setBy));
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeModeSetReqDlMsg(char * pBuffer, ModeSetReqDlMsg_t *pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding ModeSetReqDlMsg, ID 0x%.2x.\r\n", MODE_SET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = MODE_SET_REQ_DL_MSG;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = (uint8_t) pMsg->mode;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ModeSetReqDlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MODE, getStringMode(pMsg->mode));
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeModeSetCnfUlMsg(char * pBuffer, ModeSetCnfUlMsg_t *pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding ModeSetCnfUlMsg, ID 0x%.2x.\r\n", MODE_SET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = MODE_SET_CNF_UL_MSG;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = (uint8_t) pMsg->mode;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ModeSetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MODE, getStringMode(pMsg->mode));
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeModeGetReqDlMsg(char * pBuffer, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding ModeGetReqDlMsg, ID 0x%.2x.\r\n", MODE_GET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = MODE_GET_REQ_DL_MSG;
    numBytesEncoded++;
    // Empty body
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;
    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ModeGetReqDlMsg");
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeModeGetCnfUlMsg(char * pBuffer, ModeGetCnfUlMsg_t *pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding ModeGetCnfUlMsg, ID 0x%.2x.\r\n", MODE_GET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = MODE_GET_CNF_UL_MSG;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = (uint8_t) pMsg->mode;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ModeGetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MODE, getStringMode(pMsg->mode));
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeIntervalsGetReqDlMsg(char * pBuffer, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding IntervalsGetReqDlMsg, ID 0x%.2x.\r\n", INTERVALS_GET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = INTERVALS_GET_REQ_DL_MSG;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    // Empty body
    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "IntervalsGetReqDlMsg");
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeIntervalsGetCnfUlMsg(char * pBuffer, IntervalsGetCnfUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding IntervalsGetCnfUlMsg, ID 0x%.2x.\r\n", INTERVALS_GET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = INTERVALS_GET_CNF_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->reportingInterval);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->heartbeatSeconds);
    pBuffer[numBytesEncoded] = pMsg->heartbeatSnapToRtc;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "IntervalsGetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_REPORTING_INTERVAL, pMsg->reportingInterval);
        *ppLog += logHeartbeat(*ppLog, pLogSize, pMsg->heartbeatSeconds, pMsg->heartbeatSnapToRtc);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeReportingIntervalSetReqDlMsg(char * pBuffer, ReportingIntervalSetReqDlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding ReportingIntervalSetReqDlMsg, ID 0x%.2x.\r\n", REPORTING_INTERVAL_SET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = REPORTING_INTERVAL_SET_REQ_DL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->reportingInterval);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ReportingIntervalSetReqDlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_REPORTING_INTERVAL, pMsg->reportingInterval);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeReportingIntervalSetCnfUlMsg(char * pBuffer, ReportingIntervalSetCnfUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding ReportingIntervalSetCnfUlMsg, ID 0x%.2x.\r\n", REPORTING_INTERVAL_SET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = REPORTING_INTERVAL_SET_CNF_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->reportingInterval);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ReportingIntervalSetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_REPORTING_INTERVAL, pMsg->reportingInterval);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeHeartbeatSetReqDlMsg(char * pBuffer, HeartbeatSetReqDlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding HeartbeatSetReqDlMsg, ID 0x%.2x.\r\n", HEARTBEAT_SET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = HEARTBEAT_SET_REQ_DL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->heartbeatSeconds);
    pBuffer[numBytesEncoded] = pMsg->heartbeatSnapToRtc;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "HeartbeatSetReqDlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logHeartbeat(*ppLog, pLogSize, pMsg->heartbeatSeconds, pMsg->heartbeatSnapToRtc);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeHeartbeatSetCnfUlMsg(char * pBuffer, HeartbeatSetCnfUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding HeartbeatSetCnfUlMsg, ID 0x%.2x.\r\n", HEARTBEAT_SET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = HEARTBEAT_SET_CNF_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->heartbeatSeconds);
    pBuffer[numBytesEncoded] = pMsg->heartbeatSnapToRtc;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "HeartbeatSetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logHeartbeat(*ppLog, pLogSize, pMsg->heartbeatSeconds, pMsg->heartbeatSnapToRtc);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeMeasurementsGetReqDlMsg(char * pBuffer, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding MeasurementsGetReqDlMsg, ID 0x%.2x.\r\n", MEASUREMENTS_GET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = MEASUREMENTS_GET_REQ_DL_MSG;
    numBytesEncoded++;
    // Empty body
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsGetReqDlMsg");
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeMeasurementsGetCnfUlMsg(char * pBuffer, MeasurementsGetCnfUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding MeasurementsGetCnfUlMsg, ID 0x%.2x.\r\n", MEASUREMENTS_GET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = MEASUREMENTS_GET_CNF_UL_MSG;
    numBytesEncoded++;

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsGetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
    }

    numBytesEncoded += encodeMeasurements(&(pBuffer[numBytesEncoded]), &(pMsg->measurements), ppLog, pLogSize);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeMeasurementsIndUlMsg(char * pBuffer, MeasurementsIndUlMsg_t * pMsg, char **ppLog, uint32_t *pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding MeasurementsIndUlMsg, ID 0x%.2x.\r\n", MEASUREMENTS_IND_UL_MSG);
    pBuffer[numBytesEncoded] = MEASUREMENTS_IND_UL_MSG;
    numBytesEncoded++;

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsIndUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
    }

    numBytesEncoded += encodeMeasurements(&(pBuffer[numBytesEncoded]), &(pMsg->measurements), ppLog, pLogSize);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeMeasurementControlSetReqDlMsg(char * pBuffer, MeasurementControlSetReqDlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding MeasurementControlSetReqDlMsg, ID 0x%.2x.\r\n", MEASUREMENT_CONTROL_SET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = MEASUREMENT_CONTROL_SET_REQ_DL_MSG;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = (uint8_t) pMsg->measurementType;
    numBytesEncoded++;

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementControlSetReqDlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
    }

    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), pMsg->measurementType, (MeasurementControlUnion_t *) &(pMsg->measurementControl), ppLog, pLogSize);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeMeasurementControlSetCnfUlMsg(char * pBuffer, MeasurementControlSetCnfUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding MeasurementControlSetCnfUlMsg, ID 0x%.2x.\r\n", MEASUREMENT_CONTROL_SET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = MEASUREMENT_CONTROL_SET_CNF_UL_MSG;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = (uint8_t) pMsg->measurementType;
    numBytesEncoded++;

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementControlSetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
    }

    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), pMsg->measurementType, (MeasurementControlUnion_t *) &(pMsg->measurementControl), ppLog, pLogSize);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeMeasurementsControlGetReqDlMsg(char * pBuffer, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding MeasurementsControlGetReqDlMsg, ID 0x%.2x.\r\n", MEASUREMENTS_CONTROL_GET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = MEASUREMENTS_CONTROL_GET_REQ_DL_MSG;
    numBytesEncoded++;
    // Empty body
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsControlGetReqDlMsg");
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeMeasurementsControlGetCnfUlMsg(char * pBuffer, MeasurementsControlGetCnfUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding MeasurementsControlGetCnfUlMsg, ID 0x%.2x.\r\n", MEASUREMENTS_CONTROL_GET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = MEASUREMENTS_CONTROL_GET_CNF_UL_MSG;
    numBytesEncoded++;

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsControlGetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
    }

    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), MEASUREMENT_GNSS, (MeasurementControlUnion_t *) &(pMsg->gnss), ppLog, pLogSize);
    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), MEASUREMENT_CELL_ID, (MeasurementControlUnion_t *) &(pMsg->cellId), ppLog, pLogSize);
    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), MEASUREMENT_RSRP, (MeasurementControlUnion_t *) &(pMsg->rsrp), ppLog, pLogSize);
    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), MEASUREMENT_RSSI, (MeasurementControlUnion_t *) &(pMsg->rssi), ppLog, pLogSize);
    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), MEASUREMENT_TEMPERATURE, (MeasurementControlUnion_t *) &(pMsg->temperature), ppLog, pLogSize);
    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), MEASUREMENT_POWER_STATE, (MeasurementControlUnion_t *) &(pMsg->powerState), ppLog, pLogSize);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeMeasurementsControlIndUlMsg(char * pBuffer, MeasurementsControlIndUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding MeasurementsControlIndUlMsg, ID 0x%.2x.\r\n", MEASUREMENTS_CONTROL_IND_UL_MSG);
    pBuffer[numBytesEncoded] = MEASUREMENTS_CONTROL_IND_UL_MSG;
    numBytesEncoded++;

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsControlIndUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
    }

    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), MEASUREMENT_GNSS, (MeasurementControlUnion_t *) &(pMsg->gnss), ppLog, pLogSize);
    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), MEASUREMENT_CELL_ID, (MeasurementControlUnion_t *) &(pMsg->cellId), ppLog, pLogSize);
    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), MEASUREMENT_RSRP, (MeasurementControlUnion_t *) &(pMsg->rsrp), ppLog, pLogSize);
    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), MEASUREMENT_RSSI, (MeasurementControlUnion_t *) &(pMsg->rssi), ppLog, pLogSize);
    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), MEASUREMENT_TEMPERATURE, (MeasurementControlUnion_t *) &(pMsg->temperature), ppLog, pLogSize);
    numBytesEncoded += encodeMeasurementControl(&(pBuffer[numBytesEncoded]), MEASUREMENT_POWER_STATE, (MeasurementControlUnion_t *) &(pMsg->powerState), ppLog, pLogSize);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeMeasurementsControlDefaultsSetReqDlMsg(char * pBuffer, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding MeasurementsControlDefaultsSetReqDlMsg, ID 0x%.2x.\r\n", MEASUREMENTS_CONTROL_DEFAULTS_SET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = MEASUREMENTS_CONTROL_DEFAULTS_SET_REQ_DL_MSG;
    numBytesEncoded++;
    // Empty body
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;
    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsControlDefaultsSetReqDlMsg");
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeMeasurementsControlDefaultsSetCnfUlMsg(char * pBuffer, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding MeasurementsControlDefaultsSetCnfUlMsg, ID 0x%.2x.\r\n", MEASUREMENTS_CONTROL_DEFAULTS_SET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = MEASUREMENTS_CONTROL_DEFAULTS_SET_CNF_UL_MSG;
    numBytesEncoded++;
    // Empty body
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsControlDefaultsSetCnfUlMsg");
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeTrafficReportIndUlMsg(char * pBuffer, TrafficReportIndUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding TrafficReportIndUlMsg, ID 0x%.2x.\r\n", TRAFFIC_REPORT_IND_UL_MSG);
    pBuffer[numBytesEncoded] = TRAFFIC_REPORT_IND_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numDatagramsUl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numBytesUl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numDatagramsDl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numBytesDl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numDatagramsDlBadChecksum);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficReportIndUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTrafficReportUl(*ppLog, pLogSize, pMsg->numDatagramsUl, pMsg->numBytesUl);
        *ppLog += logTrafficReportDl(*ppLog, pLogSize, pMsg->numDatagramsDl, pMsg->numBytesDl, pMsg->numDatagramsDlBadChecksum);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeTrafficReportGetReqDlMsg(char * pBuffer, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding TrafficReportGetReqDlMsg, ID 0x%.2x.\r\n", TRAFFIC_REPORT_GET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = TRAFFIC_REPORT_GET_REQ_DL_MSG;
    numBytesEncoded++;
    // Empty body
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficReportGetReqDlMsg");
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeTrafficReportGetCnfUlMsg(char * pBuffer, TrafficReportGetCnfUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding TrafficReportGetCnfUlMsg, ID 0x%.2x.\r\n", TRAFFIC_REPORT_GET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = TRAFFIC_REPORT_GET_CNF_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numDatagramsUl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numBytesUl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numDatagramsDl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numBytesDl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numDatagramsDlBadChecksum);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficReportGetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTrafficReportUl(*ppLog, pLogSize, pMsg->numDatagramsUl, pMsg->numBytesUl);
        *ppLog += logTrafficReportDl(*ppLog, pLogSize, pMsg->numDatagramsDl, pMsg->numBytesDl, pMsg->numDatagramsDlBadChecksum);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeTrafficTestModeParametersSetReqDlMsg(char * pBuffer, TrafficTestModeParametersSetReqDlMsg_t *pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding TrafficTestModeParametersSetReqDlMsg, ID 0x%.2x.\r\n", TRAFFIC_TEST_MODE_PARAMETERS_SET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = TRAFFIC_TEST_MODE_PARAMETERS_SET_REQ_DL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numUlDatagrams);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->lenUlDatagram);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numDlDatagrams);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->lenDlDatagram);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->timeoutSeconds);
    numBytesEncoded += encodeBool (&(pBuffer[numBytesEncoded]), pMsg->noReportsDuringTest);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeParametersSetReqDlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTrafficTestModeParametersUl(*ppLog, pLogSize, pMsg->numUlDatagrams, pMsg->lenUlDatagram, pMsg->noReportsDuringTest);
        *ppLog += logTrafficTestModeParametersDl(*ppLog, pLogSize, pMsg->numDlDatagrams, pMsg->lenDlDatagram);
        *ppLog += logTagWithPresenceAndUint32Value(*ppLog, pLogSize, TAG_TIMEOUT, (pMsg->timeoutSeconds != 0), pMsg->timeoutSeconds);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeTrafficTestModeParametersSetCnfUlMsg(char * pBuffer, TrafficTestModeParametersSetCnfUlMsg_t *pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding TrafficTestModeParametersSetCnfUlMsg, ID 0x%.2x.\r\n", TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numUlDatagrams);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->lenUlDatagram);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numDlDatagrams);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->lenDlDatagram);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->timeoutSeconds);
    numBytesEncoded += encodeBool (&(pBuffer[numBytesEncoded]), pMsg->noReportsDuringTest);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeParametersSetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTrafficTestModeParametersUl(*ppLog, pLogSize, pMsg->numUlDatagrams, pMsg->lenUlDatagram, pMsg->noReportsDuringTest);
        *ppLog += logTrafficTestModeParametersDl(*ppLog, pLogSize, pMsg->numDlDatagrams, pMsg->lenDlDatagram);
        *ppLog += logTagWithPresenceAndUint32Value(*ppLog, pLogSize, TAG_TIMEOUT, (pMsg->timeoutSeconds != 0), pMsg->timeoutSeconds);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeTrafficTestModeParametersGetReqDlMsg(char * pBuffer, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding TrafficTestModeParametersGetReqDlMsg, ID 0x%.2x.\r\n", TRAFFIC_TEST_MODE_PARAMETERS_GET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = TRAFFIC_TEST_MODE_PARAMETERS_GET_REQ_DL_MSG;
    numBytesEncoded++;
    // Empty body
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeParametersGetReqDlMsg");
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeTrafficTestModeParametersGetCnfUlMsg(char * pBuffer, TrafficTestModeParametersGetCnfUlMsg_t *pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding TrafficTestModeParametersGetCnfUlMsg, ID 0x%.2x.\r\n", TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numUlDatagrams);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->lenUlDatagram);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numDlDatagrams);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->lenDlDatagram);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->timeoutSeconds);
    numBytesEncoded += encodeBool (&(pBuffer[numBytesEncoded]), pMsg->noReportsDuringTest);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeParametersGetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTrafficTestModeParametersUl(*ppLog, pLogSize, pMsg->numUlDatagrams, pMsg->lenUlDatagram, pMsg->noReportsDuringTest);
        *ppLog += logTrafficTestModeParametersDl(*ppLog, pLogSize, pMsg->numDlDatagrams, pMsg->lenDlDatagram);
        *ppLog += logTagWithPresenceAndUint32Value(*ppLog, pLogSize, TAG_TIMEOUT, (pMsg->timeoutSeconds != 0), pMsg->timeoutSeconds);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

/// Encode a TrafficTestModeRuleBreakerUlDatagram
uint32_t encodeTrafficTestModeRuleBreakerDatagram(char * pBuffer, TrafficTestModeRuleBreakerDatagram_t * pSpec, bool isDownlink, char * *ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;
    char id = (char) TRAFFIC_TEST_MODE_RULE_BREAKER_UL_DATAGRAM;
    const char * pTag = TAG_MSG_UL;

    if (isDownlink)
    {
        id = (char) TRAFFIC_TEST_MODE_RULE_BREAKER_DL_DATAGRAM;
    }

    MESSAGE_CODEC_LOGMSG("Encoding TrafficTestRuleBreakerDatagram, ID 0x%.2x.\r\n", id);
    pBuffer[numBytesEncoded] = id;
    numBytesEncoded++;

    if (pSpec->length > TRAFFIC_TEST_MODE_RULE_BREAKER_MAX_LENGTH)
    {
        pSpec->length = TRAFFIC_TEST_MODE_RULE_BREAKER_MAX_LENGTH;
    }

    while (numBytesEncoded < (pSpec->length - CHECKSUM_SIZE))
    {
        pBuffer[numBytesEncoded] = (char) pSpec->fill;
        numBytesEncoded++;
    }
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        if (isDownlink)
        {
            pTag = TAG_MSG_DL;
        }

        *ppLog += logBeginTag(*ppLog, pLogSize, pTag);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeRuleBreakerDatagram");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTrafficTestModeRuleBreakerDatagram(*ppLog, pLogSize, pSpec->fill, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, pTag);
    }

    return numBytesEncoded;
}

uint32_t encodeTrafficTestModeReportIndUlMsg(char * pBuffer, TrafficTestModeReportIndUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding TrafficTestModeReportIndUlMsg, ID 0x%.2x.\r\n", TRAFFIC_TEST_MODE_REPORT_IND_UL_MSG);
    pBuffer[numBytesEncoded] = TRAFFIC_TEST_MODE_REPORT_IND_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestDatagramsUl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestBytesUl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestDatagramsDl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestBytesDl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestDlDatagramsOutOfOrder);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestDlDatagramsBad);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestDlDatagramsMissed);
    numBytesEncoded += encodeBool(&(pBuffer[numBytesEncoded]), pMsg->timedOut);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeReportIndUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTrafficTestModeReportUl(*ppLog, pLogSize, pMsg->numTrafficTestDatagramsUl, pMsg->numTrafficTestBytesUl);
        *ppLog += logTrafficTestModeReportDl(*ppLog, pLogSize, pMsg->numTrafficTestDatagramsDl, pMsg->numTrafficTestBytesDl, pMsg->numTrafficTestDlDatagramsOutOfOrder,
                pMsg->numTrafficTestDlDatagramsBad, pMsg->numTrafficTestDlDatagramsMissed);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_TIMED_OUT, getStringBoolean(pMsg->timedOut));
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeTrafficTestModeReportGetReqDlMsg(char * pBuffer, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding TrafficTestModeReportGetReqDlMsg, ID 0x%.2x.\r\n", TRAFFIC_TEST_MODE_REPORT_GET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = TRAFFIC_TEST_MODE_REPORT_GET_REQ_DL_MSG;
    numBytesEncoded++;
    // Empty body
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;
    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeReportGetReqDlMsg");
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeTrafficTestModeReportGetCnfUlMsg(char * pBuffer, TrafficTestModeReportGetCnfUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding TrafficTestModeReportGetCnfUlMsg, ID 0x%.2x.\r\n", TRAFFIC_TEST_MODE_REPORT_GET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = TRAFFIC_TEST_MODE_REPORT_GET_CNF_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestDatagramsUl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestBytesUl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestDatagramsDl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestBytesDl);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestDlDatagramsOutOfOrder);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestDlDatagramsBad);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->numTrafficTestDlDatagramsMissed);
    numBytesEncoded += encodeBool(&(pBuffer[numBytesEncoded]), pMsg->timedOut);
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeReportGetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTrafficTestModeReportUl(*ppLog, pLogSize, pMsg->numTrafficTestDatagramsUl, pMsg->numTrafficTestBytesUl);
        *ppLog += logTrafficTestModeReportDl(*ppLog, pLogSize, pMsg->numTrafficTestDatagramsDl, pMsg->numTrafficTestBytesDl, pMsg->numTrafficTestDlDatagramsOutOfOrder,
                pMsg->numTrafficTestDlDatagramsBad, pMsg->numTrafficTestDlDatagramsMissed);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_TIMED_OUT, getStringBoolean(pMsg->timedOut));
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeActivityReportIndUlMsg(char * pBuffer, ActivityReportIndUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding ActivityReportIndUlMsg, ID 0x%.2x.\r\n", ACTIVITY_REPORT_IND_UL_MSG);
    pBuffer[numBytesEncoded] = ACTIVITY_REPORT_IND_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->totalTransmitMilliseconds);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->totalReceiveMilliseconds);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->upTimeSeconds);
    if (pMsg->txPowerDbmPresent)
    {
        pBuffer[numBytesEncoded] = pMsg->txPowerDbm;
    }
    else
    {
        pBuffer[numBytesEncoded] = BYTE_NOT_PRESENT_VALUE;
    }
    numBytesEncoded++;
    if (pMsg->ulMcsPresent)
    {
        pBuffer[numBytesEncoded] = pMsg->ulMcs;
    }
    else
    {
        pBuffer[numBytesEncoded] = BYTE_NOT_PRESENT_VALUE;
    }
    numBytesEncoded++;
    if (pMsg->dlMcsPresent)
    {
        pBuffer[numBytesEncoded] = pMsg->dlMcs;
    }
    else
    {
        pBuffer[numBytesEncoded] = BYTE_NOT_PRESENT_VALUE;
    }
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ActivityReportIndUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logActivityReport(*ppLog, pLogSize, pMsg->totalTransmitMilliseconds, pMsg->totalReceiveMilliseconds, pMsg->upTimeSeconds);
        *ppLog += logUlRf(*ppLog, pLogSize, pMsg->txPowerDbmPresent, pMsg->txPowerDbm, pMsg->ulMcsPresent, pMsg->ulMcs);
        *ppLog += logDlRf(*ppLog, pLogSize, pMsg->dlMcsPresent, pMsg->dlMcs);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeActivityReportGetReqDlMsg(char * pBuffer, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding ActivityReportGetReqDlMsg, ID 0x%.2x.\r\n", ACTIVITY_REPORT_GET_REQ_DL_MSG);
    pBuffer[numBytesEncoded] = ACTIVITY_REPORT_GET_REQ_DL_MSG;
    numBytesEncoded++;
    // Empty body
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ActivityReportGetReqDlMsg");
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
    }

    return numBytesEncoded;
}

uint32_t encodeActivityReportGetCnfUlMsg(char * pBuffer, ActivityReportGetCnfUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;

    MESSAGE_CODEC_LOGMSG("Encoding ActivityReportGetCnfUlMsg, ID 0x%.2x.\r\n", ACTIVITY_REPORT_GET_CNF_UL_MSG);
    pBuffer[numBytesEncoded] = ACTIVITY_REPORT_GET_CNF_UL_MSG;
    numBytesEncoded++;
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->totalTransmitMilliseconds);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->totalReceiveMilliseconds);
    numBytesEncoded += encodeUint32(&(pBuffer[numBytesEncoded]), pMsg->upTimeSeconds);
    if (pMsg->txPowerDbmPresent)
    {
        pBuffer[numBytesEncoded] = pMsg->txPowerDbm;
    }
    else
    {
        pBuffer[numBytesEncoded] = BYTE_NOT_PRESENT_VALUE;
    }
    numBytesEncoded++;
    if (pMsg->ulMcsPresent)
    {
        pBuffer[numBytesEncoded] = pMsg->ulMcs;
    }
    else
    {
        pBuffer[numBytesEncoded] = BYTE_NOT_PRESENT_VALUE;
    }
    numBytesEncoded++;
    if (pMsg->dlMcsPresent)
    {
        pBuffer[numBytesEncoded] = pMsg->dlMcs;
    }
    else
    {
        pBuffer[numBytesEncoded] = BYTE_NOT_PRESENT_VALUE;
    }
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ActivityReportGetCnfUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logActivityReport(*ppLog, pLogSize, pMsg->totalTransmitMilliseconds, pMsg->totalReceiveMilliseconds, pMsg->upTimeSeconds);
        *ppLog += logUlRf(*ppLog, pLogSize, pMsg->txPowerDbmPresent, pMsg->txPowerDbm, pMsg->ulMcsPresent, pMsg->ulMcs);
        *ppLog += logDlRf(*ppLog, pLogSize, pMsg->dlMcsPresent, pMsg->dlMcs);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

uint32_t encodeDebugIndUlMsg(char * pBuffer, DebugIndUlMsg_t * pMsg, char **ppLog, uint32_t * pLogSize)
{
    uint32_t numBytesEncoded = 0;
    uint32_t sizeOfString = pMsg->sizeOfString;

    MESSAGE_CODEC_LOGMSG("Encoding DebugIndUlMsg, ID 0x%.2x.\r\n", DEBUG_IND_UL_MSG);
    if (sizeOfString > MAX_DEBUG_STRING_SIZE)
    {
        sizeOfString = MAX_DEBUG_STRING_SIZE;
    }
    pBuffer[numBytesEncoded] = DEBUG_IND_UL_MSG;
    numBytesEncoded++;
    pBuffer[numBytesEncoded] = pMsg->sizeOfString;
    numBytesEncoded++;
    memcpy(&(pBuffer[numBytesEncoded]), &(pMsg->string[0]), sizeOfString);
    numBytesEncoded += sizeOfString;
    pBuffer[numBytesEncoded] = calculateChecksum(&(pBuffer[0]), numBytesEncoded);
    numBytesEncoded++;

    MESSAGE_CODEC_LOGMSG("%d bytes encoded.\r\n", numBytesEncoded);

    if ((ppLog != NULL) && (*ppLog != NULL))
    {
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "DebugIndUlMsg");
        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logDebugInd(*ppLog, pLogSize, &(pMsg->string[0]), pMsg->sizeOfString);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, numBytesEncoded);
        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
    }

    return numBytesEncoded;
}

// ----------------------------------------------------------------
// MESSAGE DECODING FUNCTIONS
// ----------------------------------------------------------------
DecodeResult_t decodeDlMsg(const char ** ppInBuffer, uint32_t sizeInBuffer, DlMsgUnion_t * pOutBuffer, char * *ppLog, uint32_t * pLogSize)
{
    MsgIdDl_t msgId;
    DecodeResult_t decodeResult = DECODE_RESULT_FAILURE;
    const char * pBufferAtStart = *ppInBuffer;

    if (sizeInBuffer < MIN_MESSAGE_SIZE)
    {
        decodeResult = DECODE_RESULT_INPUT_TOO_SHORT;
        (*ppInBuffer) += sizeInBuffer;
    }
    else
    {
        decodeResult = DECODE_RESULT_UNKNOWN_MSG_ID;
        // First byte should be a valid DL message ID
        msgId = (MsgIdDl_t) **ppInBuffer;
        (*ppInBuffer)++;
        if (msgId < MAX_NUM_DL_MSGS)
        {
            switch (msgId)
            {
                case TRANSPARENT_DL_DATAGRAM:
                {
                    decodeResult = DECODE_RESULT_TRANSPARENT_DL_DATAGRAM;
                    if (pOutBuffer != NULL)
                    {
                        memcpy(&(pOutBuffer->transparentDatagram.contents[0]), *ppInBuffer, sizeInBuffer - 1);
                        (*ppInBuffer) += sizeInBuffer - 1;
                        // NOTE: there is no checksum on this message, validation of contents is up to the application
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TransparentDatagram");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTransparentData(*ppLog, pLogSize, sizeInBuffer - 1, &(pOutBuffer->transparentDatagram.contents[0]));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case PING_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_PING_REQ_DL_MSG;
                    // Empty message
                    if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                    {
                        decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                    }
                    (*ppInBuffer)++;
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "PingReqDlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case PING_CNF_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_PING_CNF_DL_MSG;
                    // Empty message
                    if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                    {
                        decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                    }
                    (*ppInBuffer)++;
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "PingCnfDlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case REBOOT_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_REBOOT_REQ_DL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->rebootReqDlMsg.sdCardNotRequired = false;
                        pOutBuffer->rebootReqDlMsg.disableModemDebug = false;
                        pOutBuffer->rebootReqDlMsg.disableButton = false;
                        pOutBuffer->rebootReqDlMsg.disableServerPing = false;

                        if (**ppInBuffer & 0x01)
                        {
                            pOutBuffer->rebootReqDlMsg.sdCardNotRequired = true;
                        }
                        if (**ppInBuffer & 0x02)
                        {
                            pOutBuffer->rebootReqDlMsg.disableModemDebug = true;
                        }
                        if (**ppInBuffer & 0x04)
                        {
                            pOutBuffer->rebootReqDlMsg.disableButton = true;
                        }
                        if (**ppInBuffer & 0x08)
                        {
                            pOutBuffer->rebootReqDlMsg.disableServerPing = true;
                        }
                        (*ppInBuffer)++;
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "RebootReqDlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "SdCardRequired", getStringBoolean(!pOutBuffer->rebootReqDlMsg.sdCardNotRequired));
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "DisableModemDebug", getStringBoolean(pOutBuffer->rebootReqDlMsg.disableModemDebug));
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "DisableButton", getStringBoolean(pOutBuffer->rebootReqDlMsg.disableButton));
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "DisableServerPing", getStringBoolean(pOutBuffer->rebootReqDlMsg.disableServerPing));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case DATE_TIME_SET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_DATE_TIME_SET_REQ_DL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->dateTimeSetReqDlMsg.time = decodeUint32(ppInBuffer);
                        pOutBuffer->dateTimeSetReqDlMsg.setDateOnly = decodeBool(ppInBuffer);
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "DateTimeSetReqDlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logDateTime(*ppLog, pLogSize, pOutBuffer->dateTimeSetReqDlMsg.time);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_SET_DATE_ONLY, getStringBoolean(pOutBuffer->dateTimeSetReqDlMsg.setDateOnly));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case DATE_TIME_GET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_DATE_TIME_GET_REQ_DL_MSG;
                    // Empty message
                    if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                    {
                        decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                    }
                    (*ppInBuffer)++;
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "DateTimeGetReqDlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case MODE_SET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_MODE_SET_REQ_DL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->modeSetReqDlMsg.mode = (Mode_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ModeSetReqDlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MODE, getStringMode(pOutBuffer->modeSetReqDlMsg.mode));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case MODE_GET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_MODE_GET_REQ_DL_MSG;
                    // Empty message
                    if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                    {
                        decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                    }
                    (*ppInBuffer)++;
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ModeGetReqDlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case HEARTBEAT_SET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_HEARTBEAT_SET_REQ_DL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->heartbeatSetReqDlMsg.heartbeatSeconds = decodeUint32(ppInBuffer);
                        pOutBuffer->heartbeatSetReqDlMsg.heartbeatSnapToRtc = decodeBool(ppInBuffer);
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "HeartbeatSetReqDlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logHeartbeat(*ppLog, pLogSize, pOutBuffer->heartbeatSetReqDlMsg.heartbeatSeconds, pOutBuffer->heartbeatSetReqDlMsg.heartbeatSnapToRtc);
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case REPORTING_INTERVAL_SET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_REPORTING_INTERVAL_SET_REQ_DL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->reportingIntervalSetReqDlMsg.reportingInterval = decodeUint32(ppInBuffer);
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ReportingIntervalSetReqDlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_REPORTING_INTERVAL, pOutBuffer->reportingIntervalSetReqDlMsg.reportingInterval);
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case INTERVALS_GET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_INTERVALS_GET_REQ_DL_MSG;
                    // Empty message
                    if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                    {
                        decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                    }
                    (*ppInBuffer)++;
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "IntervalsGetReqDlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case MEASUREMENTS_GET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_MEASUREMENTS_GET_REQ_DL_MSG;
                    // Empty message
                    if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                    {
                        decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                    }
                    (*ppInBuffer)++;
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsGetReqDlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case MEASUREMENT_CONTROL_SET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_MEASUREMENT_CONTROL_SET_REQ_DL_MSG;
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementControlSetReqDlMsg");
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                    }

                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->measurementControlSetReqDlMsg.measurementType = (MeasurementType_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (!decodeMeasurementControl(ppInBuffer, pOutBuffer->measurementControlSetReqDlMsg.measurementType,
                                (MeasurementControlUnion_t *) &(pOutBuffer->measurementControlSetReqDlMsg.measurementControl), ppLog, pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case MEASUREMENTS_CONTROL_GET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_MEASUREMENTS_CONTROL_GET_REQ_DL_MSG;
                    // Empty message
                    if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                    {
                        decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                    }
                    (*ppInBuffer)++;
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsControlGetReqDlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case MEASUREMENTS_CONTROL_DEFAULTS_SET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_MEASUREMENTS_CONTROL_DEFAULTS_SET_REQ_DL_MSG;
                    // Empty message
                    if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                    {
                        decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                    }
                    (*ppInBuffer)++;
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsControlDefaultsSetReqDlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case TRAFFIC_REPORT_GET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_TRAFFIC_REPORT_GET_REQ_DL_MSG;
                    // Empty message
                    if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                    {
                        decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                    }
                    (*ppInBuffer)++;
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficReportGetReqDlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case TRAFFIC_TEST_MODE_PARAMETERS_SET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_SET_REQ_DL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->trafficTestModeParametersSetReqDlMsg.numUlDatagrams = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersSetReqDlMsg.lenUlDatagram = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersSetReqDlMsg.numDlDatagrams = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersSetReqDlMsg.lenDlDatagram = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersSetReqDlMsg.timeoutSeconds = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersSetReqDlMsg.noReportsDuringTest = decodeBool(ppInBuffer);
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeParametersSetReqDlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTrafficTestModeParametersUl(*ppLog, pLogSize, pOutBuffer->trafficTestModeParametersSetReqDlMsg.numUlDatagrams,
                                pOutBuffer->trafficTestModeParametersSetReqDlMsg.lenUlDatagram,
                                pOutBuffer->trafficTestModeParametersSetReqDlMsg.noReportsDuringTest);
                        *ppLog += logTrafficTestModeParametersDl(*ppLog, pLogSize, pOutBuffer->trafficTestModeParametersSetReqDlMsg.numDlDatagrams,
                                pOutBuffer->trafficTestModeParametersSetReqDlMsg.lenDlDatagram);
                        *ppLog += logTagWithPresenceAndUint32Value(*ppLog, pLogSize, TAG_TIMEOUT, (pOutBuffer->trafficTestModeParametersSetReqDlMsg.timeoutSeconds != 0),
                                pOutBuffer->trafficTestModeParametersSetReqDlMsg.timeoutSeconds);
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case TRAFFIC_TEST_MODE_PARAMETERS_GET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_GET_REQ_DL_MSG;
                    // Empty message
                    if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                    {
                        decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                    }
                    (*ppInBuffer)++;
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeParametersGetReqDlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case TRAFFIC_TEST_MODE_RULE_BREAKER_DL_DATAGRAM:
                {
                    decodeResult = DECODE_RESULT_TRAFFIC_TEST_MODE_RULE_BREAKER_DL_DATAGRAM;
                    if (pOutBuffer != NULL)
                    {
                        uint8_t expectedFill = pOutBuffer->trafficTestModeRuleBreakerDatagram.fill;
                        uint8_t expectedLength = pOutBuffer->trafficTestModeRuleBreakerDatagram.length;

                        if (!decodeTrafficTestModeRuleBreakerDatagram(ppInBuffer, true, &(pOutBuffer->trafficTestModeRuleBreakerDatagram), ppLog, pLogSize))
                        {
                            if ((sizeInBuffer == expectedLength) && (pOutBuffer->trafficTestModeRuleBreakerDatagram.fill != expectedFill))
                            {
                                decodeResult = DECODE_RESULT_OUT_OF_SEQUENCE_TRAFFIC_TEST_MODE_DATAGRAM;
                            }
                            else
                            {
                                decodeResult = DECODE_RESULT_BAD_TRAFFIC_TEST_MODE_DATAGRAM;
                            }
                        }
                    }
                }
                break;
                case TRAFFIC_TEST_MODE_REPORT_GET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_GET_REQ_DL_MSG;
                    // Empty message
                    if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                    {
                        decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                    }
                    (*ppInBuffer)++;
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeReportGetReqDlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case ACTIVITY_REPORT_GET_REQ_DL_MSG:
                {
                    decodeResult = DECODE_RESULT_ACTIVITY_REPORT_GET_REQ_DL_MSG;
                    // Empty message
                    if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                    {
                        decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                    }
                    (*ppInBuffer)++;
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ActivityReportGetReqDlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                default:
                    // The decodeResult will be left as Unknown message
                break;
            }
        }
    }

    return decodeResult;
}

DecodeResult_t decodeUlMsg(const char ** ppInBuffer, uint32_t sizeInBuffer, UlMsgUnion_t * pOutBuffer, char * *ppLog, uint32_t * pLogSize)
{
    MsgIdUl_t msgId;
    DecodeResult_t decodeResult = DECODE_RESULT_FAILURE;
    const char * pBufferAtStart = *ppInBuffer;

    if (sizeInBuffer < MIN_MESSAGE_SIZE)
    {
        decodeResult = DECODE_RESULT_INPUT_TOO_SHORT;
        (*ppInBuffer) += sizeInBuffer;
    }
    else
    {
        decodeResult = DECODE_RESULT_UNKNOWN_MSG_ID;
        // First byte should be a valid UL message ID
        msgId = (MsgIdUl_t) **ppInBuffer;
        (*ppInBuffer)++;
        if (msgId < MAX_NUM_UL_MSGS)
        {
            switch (msgId)
            {
                case TRANSPARENT_UL_DATAGRAM:
                {
                    // NOTE: when the calling function receives this result
                    // it should check that were at the start of a datagram,
                    // otherwise this could copy data from beyond the end of
                    // the input buffer.
                    decodeResult = DECODE_RESULT_TRANSPARENT_UL_DATAGRAM;
                    if (pOutBuffer != NULL)
                    {
                        memcpy(&(pOutBuffer->transparentDatagram.contents[0]), *ppInBuffer, sizeInBuffer - 1);
                        (*ppInBuffer) += sizeInBuffer - 1;
                        // NOTE: no checksum on this message
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TransparentDatagram");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTransparentData(*ppLog, pLogSize, sizeInBuffer - 1, &(pOutBuffer->transparentDatagram.contents[0]));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case PING_REQ_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_PING_REQ_UL_MSG;
                    // Empty message
                    if (pOutBuffer != NULL)
                    {
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "PingReqUlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case PING_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_PING_CNF_UL_MSG;
                    // Empty message
                    if (pOutBuffer != NULL)
                    {
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_DL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "PingCnfUlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_DL);
                    }
                }
                break;
                case INIT_IND_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_INIT_IND_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->initIndUlMsg.wakeUpCode = (WakeUpCode_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        pOutBuffer->initIndUlMsg.revisionLevel = (uint8_t) * *ppInBuffer;
                        (*ppInBuffer)++;

                        pOutBuffer->initIndUlMsg.sdCardNotRequired = false;
                        pOutBuffer->initIndUlMsg.disableModemDebug = false;
                        pOutBuffer->initIndUlMsg.disableButton = false;
                        pOutBuffer->initIndUlMsg.disableServerPing = false;

                        if (**ppInBuffer & 0x01)
                        {
                            pOutBuffer->initIndUlMsg.sdCardNotRequired = true;
                        }
                        if (**ppInBuffer & 0x02)
                        {
                            pOutBuffer->initIndUlMsg.disableModemDebug = true;
                        }
                        if (**ppInBuffer & 0x04)
                        {
                            pOutBuffer->initIndUlMsg.disableButton = true;
                        }
                        if (**ppInBuffer & 0x08)
                        {
                            pOutBuffer->initIndUlMsg.disableServerPing = true;
                        }
                        (*ppInBuffer)++;
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "InitIndUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "WakeupCode", getStringWakeUpCode(pOutBuffer->initIndUlMsg.wakeUpCode));
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, "RevisionLevel", pOutBuffer->initIndUlMsg.revisionLevel);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, "SdCardRequired", !(pOutBuffer->initIndUlMsg.sdCardNotRequired));
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "DisableModemDebug", getStringBoolean(pOutBuffer->initIndUlMsg.disableModemDebug));
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "DisableButton", getStringBoolean(pOutBuffer->initIndUlMsg.disableButton));
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, "DisableServerPing", getStringBoolean(pOutBuffer->initIndUlMsg.disableServerPing));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case POLL_IND_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_POLL_IND_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->pollIndUlMsg.mode = (Mode_t) (**ppInBuffer & 0x07);
                        pOutBuffer->pollIndUlMsg.energyLeft = (EnergyLeft_t) ((**ppInBuffer >> 3) & 0x07);
                        pOutBuffer->pollIndUlMsg.diskSpaceLeft = (DiskSpaceLeft_t) ((**ppInBuffer >> 6) & 0x03);
                        (*ppInBuffer)++;
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "PollIndUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MODE, getStringMode(pOutBuffer->pollIndUlMsg.mode));
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_ENERGY_LEFT, getStringEnergyLeft(pOutBuffer->pollIndUlMsg.energyLeft));
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_DISK_SPACE_LEFT, getStringDiskSpaceLeft(pOutBuffer->pollIndUlMsg.diskSpaceLeft));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case DATE_TIME_IND_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_DATE_TIME_IND_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->dateTimeIndUlMsg.time = decodeUint32(ppInBuffer);
                        pOutBuffer->dateTimeIndUlMsg.setBy = (TimeSetBy_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "DateTimeIndUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logDateTime(*ppLog, pLogSize, pOutBuffer->dateTimeIndUlMsg.time);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_TIME_SET_BY, getStringTimeSetBy(pOutBuffer->dateTimeIndUlMsg.setBy));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case DATE_TIME_SET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_DATE_TIME_SET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->dateTimeSetCnfUlMsg.time = decodeUint32(ppInBuffer);
                        pOutBuffer->dateTimeSetCnfUlMsg.setBy = (TimeSetBy_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "DateTimeSetCnfUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logDateTime(*ppLog, pLogSize, pOutBuffer->dateTimeSetCnfUlMsg.time);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_TIME_SET_BY, getStringTimeSetBy(pOutBuffer->dateTimeSetCnfUlMsg.setBy));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case DATE_TIME_GET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_DATE_TIME_GET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->dateTimeGetCnfUlMsg.time = decodeUint32(ppInBuffer);
                        pOutBuffer->dateTimeGetCnfUlMsg.setBy = (TimeSetBy_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "DateTimeGetCnfUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logDateTime(*ppLog, pLogSize, pOutBuffer->dateTimeGetCnfUlMsg.time);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_TIME_SET_BY, getStringTimeSetBy(pOutBuffer->dateTimeGetCnfUlMsg.setBy));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case MODE_SET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_MODE_SET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->modeSetCnfUlMsg.mode = (Mode_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ModeSetCnfUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MODE, getStringMode(pOutBuffer->modeSetCnfUlMsg.mode));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case MODE_GET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_MODE_GET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->modeGetCnfUlMsg.mode = (Mode_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ModeGetCnfUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MODE, getStringMode(pOutBuffer->modeGetCnfUlMsg.mode));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case HEARTBEAT_SET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_HEARTBEAT_SET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->heartbeatSetCnfUlMsg.heartbeatSeconds = decodeUint32(ppInBuffer);
                        pOutBuffer->heartbeatSetCnfUlMsg.heartbeatSnapToRtc = decodeBool(ppInBuffer);
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "HeartbeatSetCnfUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logHeartbeat(*ppLog, pLogSize, pOutBuffer->heartbeatSetCnfUlMsg.heartbeatSeconds, pOutBuffer->heartbeatSetCnfUlMsg.heartbeatSnapToRtc);
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case REPORTING_INTERVAL_SET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_REPORTING_INTERVAL_SET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->reportingIntervalSetCnfUlMsg.reportingInterval = decodeUint32(ppInBuffer);
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ReportingIntervalSetCnfUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_REPORTING_INTERVAL, pOutBuffer->reportingIntervalSetCnfUlMsg.reportingInterval);
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case INTERVALS_GET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_INTERVALS_GET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->intervalsGetCnfUlMsg.reportingInterval = decodeUint32(ppInBuffer);
                        pOutBuffer->intervalsGetCnfUlMsg.heartbeatSeconds = decodeUint32(ppInBuffer);
                        pOutBuffer->intervalsGetCnfUlMsg.heartbeatSnapToRtc = decodeBool(ppInBuffer);
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "IntervalsGetCnfUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_REPORTING_INTERVAL, pOutBuffer->intervalsGetCnfUlMsg.reportingInterval);
                        *ppLog += logHeartbeat(*ppLog, pLogSize, pOutBuffer->intervalsGetCnfUlMsg.heartbeatSeconds, pOutBuffer->intervalsGetCnfUlMsg.heartbeatSnapToRtc);
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case MEASUREMENTS_GET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_MEASUREMENTS_GET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        if ((ppLog != NULL) && (*ppLog != NULL))
                        {
                            *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                            *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsGetCnfUlMsg");
                            *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        }
                        if (!decodeMeasurements(ppInBuffer, &(pOutBuffer->measurementsGetCnfUlMsg.measurements), ppLog, pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                        if ((ppLog != NULL) && (*ppLog != NULL))
                        {
                            *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                            *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                            *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                            *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                        }
                    }
                }
                break;
                case MEASUREMENTS_IND_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_MEASUREMENTS_IND_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        if ((ppLog != NULL) && (*ppLog != NULL))
                        {
                            *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                            *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsIndUlMsg");
                            *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        }
                        if (!decodeMeasurements(ppInBuffer, &(pOutBuffer->measurementsIndUlMsg.measurements), ppLog, pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                        if ((ppLog != NULL) && (*ppLog != NULL))
                        {
                            *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                            *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                            *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                            *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                        }
                    }
                }
                break;
                case MEASUREMENT_CONTROL_SET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_MEASUREMENT_CONTROL_SET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        if ((ppLog != NULL) && (*ppLog != NULL))
                        {
                            *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                            *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsControlSetCnfUlMsg");
                            *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        }
                        pOutBuffer->measurementControlSetCnfUlMsg.measurementType = (MeasurementType_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (!decodeMeasurementControl(ppInBuffer, pOutBuffer->measurementControlSetCnfUlMsg.measurementType,
                                (MeasurementControlUnion_t *) &(pOutBuffer->measurementControlSetCnfUlMsg.measurementControl), ppLog, pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                        if ((ppLog != NULL) && (*ppLog != NULL))
                        {
                            *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                            *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                            *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                            *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                        }
                    }
                }
                break;
                case MEASUREMENTS_CONTROL_GET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_MEASUREMENTS_CONTROL_GET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        if ((ppLog != NULL) && (*ppLog != NULL))
                        {
                            *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                            *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsControlGetCnfUlMsg");
                            *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        }
                        if (!decodeMeasurementControl(ppInBuffer, MEASUREMENT_GNSS, (MeasurementControlUnion_t *) &(pOutBuffer->measurementsControlGetCnfUlMsg.gnss), ppLog, pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if ((decodeResult != DECODE_RESULT_BAD_MSG_FORMAT)
                                && !decodeMeasurementControl(ppInBuffer, MEASUREMENT_CELL_ID, (MeasurementControlUnion_t *) &(pOutBuffer->measurementsControlGetCnfUlMsg.cellId), ppLog, pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if ((decodeResult != DECODE_RESULT_BAD_MSG_FORMAT)
                                && !decodeMeasurementControl(ppInBuffer, MEASUREMENT_RSRP, (MeasurementControlUnion_t *) &(pOutBuffer->measurementsControlGetCnfUlMsg.rsrp), ppLog, pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if ((decodeResult != DECODE_RESULT_BAD_MSG_FORMAT)
                                && !decodeMeasurementControl(ppInBuffer, MEASUREMENT_RSSI, (MeasurementControlUnion_t *) &(pOutBuffer->measurementsControlGetCnfUlMsg.rssi), ppLog, pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if ((decodeResult != DECODE_RESULT_BAD_MSG_FORMAT)
                                && !decodeMeasurementControl(ppInBuffer, MEASUREMENT_TEMPERATURE, (MeasurementControlUnion_t *) &(pOutBuffer->measurementsControlGetCnfUlMsg.temperature), ppLog,
                                        pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if ((decodeResult != DECODE_RESULT_BAD_MSG_FORMAT)
                                && !decodeMeasurementControl(ppInBuffer, MEASUREMENT_POWER_STATE, (MeasurementControlUnion_t *) &(pOutBuffer->measurementsControlGetCnfUlMsg.powerState), ppLog,
                                        pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                        if ((ppLog != NULL) && (*ppLog != NULL))
                        {
                            *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                            *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                            *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                            *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                        }
                    }
                }
                break;
                case MEASUREMENTS_CONTROL_IND_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_MEASUREMENTS_CONTROL_IND_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        if ((ppLog != NULL) && (*ppLog != NULL))
                        {
                            *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                            *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsControlIndUlMsg");
                            *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        }
                        if (!decodeMeasurementControl(ppInBuffer, MEASUREMENT_GNSS, (MeasurementControlUnion_t *) &(pOutBuffer->measurementsControlGetCnfUlMsg.gnss), ppLog, pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if ((decodeResult != DECODE_RESULT_BAD_MSG_FORMAT)
                                && !decodeMeasurementControl(ppInBuffer, MEASUREMENT_CELL_ID, (MeasurementControlUnion_t *) &(pOutBuffer->measurementsControlGetCnfUlMsg.cellId), ppLog, pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if ((decodeResult != DECODE_RESULT_BAD_MSG_FORMAT)
                                && !decodeMeasurementControl(ppInBuffer, MEASUREMENT_RSRP, (MeasurementControlUnion_t *) &(pOutBuffer->measurementsControlGetCnfUlMsg.rsrp), ppLog, pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if ((decodeResult != DECODE_RESULT_BAD_MSG_FORMAT)
                                && !decodeMeasurementControl(ppInBuffer, MEASUREMENT_RSSI, (MeasurementControlUnion_t *) &(pOutBuffer->measurementsControlGetCnfUlMsg.rssi), ppLog, pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if ((decodeResult != DECODE_RESULT_BAD_MSG_FORMAT)
                                && !decodeMeasurementControl(ppInBuffer, MEASUREMENT_TEMPERATURE, (MeasurementControlUnion_t *) &(pOutBuffer->measurementsControlGetCnfUlMsg.temperature), ppLog,
                                        pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if ((decodeResult != DECODE_RESULT_BAD_MSG_FORMAT)
                                && !decodeMeasurementControl(ppInBuffer, MEASUREMENT_POWER_STATE, (MeasurementControlUnion_t *) &(pOutBuffer->measurementsControlGetCnfUlMsg.powerState), ppLog,
                                        pLogSize))
                        {
                            decodeResult = DECODE_RESULT_BAD_MSG_FORMAT;
                        }
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                        if ((ppLog != NULL) && (*ppLog != NULL))
                        {
                            *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                            *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                            *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                            *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                        }
                    }
                }
                break;
                case MEASUREMENTS_CONTROL_DEFAULTS_SET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_MEASUREMENTS_CONTROL_DEFAULTS_SET_CNF_UL_MSG;
                    // Empty message
                    if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                    {
                        decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                    }
                    if (pOutBuffer != NULL)
                    {
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "MeasurementsControlDefaultsSetCnfUlMsg");
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case TRAFFIC_REPORT_IND_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_TRAFFIC_REPORT_IND_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->trafficReportIndUlMsg.numDatagramsUl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficReportIndUlMsg.numBytesUl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficReportIndUlMsg.numDatagramsDl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficReportIndUlMsg.numBytesDl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficReportIndUlMsg.numDatagramsDlBadChecksum = decodeUint32(ppInBuffer);
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficReportIndUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTrafficReportUl(*ppLog, pLogSize, pOutBuffer->trafficReportIndUlMsg.numDatagramsUl, pOutBuffer->trafficReportIndUlMsg.numBytesUl);
                        *ppLog += logTrafficReportDl(*ppLog, pLogSize, pOutBuffer->trafficReportIndUlMsg.numDatagramsDl, pOutBuffer->trafficReportIndUlMsg.numBytesDl,
                                pOutBuffer->trafficReportIndUlMsg.numDatagramsDlBadChecksum);
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case TRAFFIC_REPORT_GET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_TRAFFIC_REPORT_GET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->trafficReportGetCnfUlMsg.numDatagramsUl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficReportGetCnfUlMsg.numBytesUl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficReportGetCnfUlMsg.numDatagramsDl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficReportGetCnfUlMsg.numBytesDl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficReportGetCnfUlMsg.numDatagramsDlBadChecksum = decodeUint32(ppInBuffer);
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficReportGetCnfUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTrafficReportUl(*ppLog, pLogSize, pOutBuffer->trafficReportGetCnfUlMsg.numDatagramsUl, pOutBuffer->trafficReportGetCnfUlMsg.numBytesUl);
                        *ppLog += logTrafficReportDl(*ppLog, pLogSize, pOutBuffer->trafficReportGetCnfUlMsg.numDatagramsDl, pOutBuffer->trafficReportGetCnfUlMsg.numBytesDl,
                                pOutBuffer->trafficReportGetCnfUlMsg.numDatagramsDlBadChecksum);
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->trafficTestModeParametersSetCnfUlMsg.numUlDatagrams = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersSetCnfUlMsg.lenUlDatagram = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersSetCnfUlMsg.numDlDatagrams = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersSetCnfUlMsg.lenDlDatagram = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersSetCnfUlMsg.timeoutSeconds = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersSetCnfUlMsg.noReportsDuringTest = decodeBool(ppInBuffer);
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeParametersSetCnfUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTrafficTestModeParametersUl(*ppLog, pLogSize, pOutBuffer->trafficTestModeParametersSetCnfUlMsg.numUlDatagrams,
                                pOutBuffer->trafficTestModeParametersSetCnfUlMsg.lenUlDatagram,
                                pOutBuffer->trafficTestModeParametersSetCnfUlMsg.noReportsDuringTest);
                        *ppLog += logTrafficTestModeParametersDl(*ppLog, pLogSize, pOutBuffer->trafficTestModeParametersSetCnfUlMsg.numDlDatagrams,
                                pOutBuffer->trafficTestModeParametersSetCnfUlMsg.lenDlDatagram);
                        *ppLog += logTagWithPresenceAndUint32Value(*ppLog, pLogSize, TAG_TIMEOUT, (pOutBuffer->trafficTestModeParametersSetCnfUlMsg.timeoutSeconds != 0),
                                pOutBuffer->trafficTestModeParametersSetCnfUlMsg.timeoutSeconds);
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->trafficTestModeParametersGetCnfUlMsg.numUlDatagrams = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersGetCnfUlMsg.lenUlDatagram = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersGetCnfUlMsg.numDlDatagrams = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersGetCnfUlMsg.lenDlDatagram = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersGetCnfUlMsg.timeoutSeconds = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeParametersGetCnfUlMsg.noReportsDuringTest = decodeBool (ppInBuffer);
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeParametersGetCnfUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTrafficTestModeParametersUl(*ppLog, pLogSize, pOutBuffer->trafficTestModeParametersGetCnfUlMsg.numUlDatagrams,
                                pOutBuffer->trafficTestModeParametersGetCnfUlMsg.lenUlDatagram,
                                pOutBuffer->trafficTestModeParametersGetCnfUlMsg.noReportsDuringTest);
                        *ppLog += logTrafficTestModeParametersDl(*ppLog, pLogSize, pOutBuffer->trafficTestModeParametersGetCnfUlMsg.numDlDatagrams,
                                pOutBuffer->trafficTestModeParametersGetCnfUlMsg.lenDlDatagram);
                        *ppLog += logTagWithPresenceAndUint32Value(*ppLog, pLogSize, TAG_TIMEOUT, (pOutBuffer->trafficTestModeParametersGetCnfUlMsg.timeoutSeconds != 0),
                                pOutBuffer->trafficTestModeParametersGetCnfUlMsg.timeoutSeconds);
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case TRAFFIC_TEST_MODE_RULE_BREAKER_UL_DATAGRAM:
                {
                    decodeResult = DECODE_RESULT_TRAFFIC_TEST_MODE_RULE_BREAKER_UL_DATAGRAM;
                    if (pOutBuffer != NULL)
                    {
                        uint8_t expectedFill = pOutBuffer->trafficTestModeRuleBreakerDatagram.fill;
                        uint8_t expectedLength = pOutBuffer->trafficTestModeRuleBreakerDatagram.length;

                        if (!decodeTrafficTestModeRuleBreakerDatagram(ppInBuffer, false, &(pOutBuffer->trafficTestModeRuleBreakerDatagram), ppLog, pLogSize))
                        {
                            if ((sizeInBuffer == expectedLength) && (pOutBuffer->trafficTestModeRuleBreakerDatagram.fill != expectedFill))
                            {
                                decodeResult = DECODE_RESULT_OUT_OF_SEQUENCE_TRAFFIC_TEST_MODE_DATAGRAM;
                            }
                            else
                            {
                                decodeResult = DECODE_RESULT_BAD_TRAFFIC_TEST_MODE_DATAGRAM;
                            }
                        }
                    }
                }
                break;
                case TRAFFIC_TEST_MODE_REPORT_IND_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_IND_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestDatagramsUl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestBytesUl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestDatagramsDl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestBytesDl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestDlDatagramsOutOfOrder = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestDlDatagramsBad = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestDlDatagramsMissed = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportIndUlMsg.timedOut = decodeBool(ppInBuffer);
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeReportIndUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTrafficTestModeReportUl(*ppLog, pLogSize, pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestDatagramsUl,
                                pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestBytesUl);
                        *ppLog += logTrafficTestModeReportDl(*ppLog, pLogSize, pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestDatagramsDl,
                                pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestBytesDl, pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestDlDatagramsOutOfOrder,
                                pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestDlDatagramsBad, pOutBuffer->trafficTestModeReportIndUlMsg.numTrafficTestDlDatagramsMissed);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_TIMED_OUT, getStringBoolean(pOutBuffer->trafficTestModeReportIndUlMsg.timedOut));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case TRAFFIC_TEST_MODE_REPORT_GET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_TRAFFIC_TEST_MODE_REPORT_GET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestDatagramsUl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestBytesUl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestDatagramsDl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestBytesDl = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestDlDatagramsOutOfOrder = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestDlDatagramsBad = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestDlDatagramsMissed = decodeUint32(ppInBuffer);
                        pOutBuffer->trafficTestModeReportGetCnfUlMsg.timedOut = decodeBool(ppInBuffer);
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "TrafficTestModeReportGetCnfUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTrafficTestModeReportUl(*ppLog, pLogSize, pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestDatagramsUl,
                                pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestBytesUl);
                        *ppLog += logTrafficTestModeReportDl(*ppLog, pLogSize, pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestDatagramsDl,
                                pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestBytesDl, pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestDlDatagramsOutOfOrder,
                                pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestDlDatagramsBad, pOutBuffer->trafficTestModeReportGetCnfUlMsg.numTrafficTestDlDatagramsMissed);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_TIMED_OUT, getStringBoolean(pOutBuffer->trafficTestModeReportIndUlMsg.timedOut));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case ACTIVITY_REPORT_IND_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_ACTIVITY_REPORT_IND_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->activityReportIndUlMsg.txPowerDbmPresent = false;
                        pOutBuffer->activityReportIndUlMsg.ulMcsPresent = false;
                        pOutBuffer->activityReportIndUlMsg.dlMcsPresent = false;
                        pOutBuffer->activityReportIndUlMsg.totalTransmitMilliseconds = decodeUint32(ppInBuffer);
                        pOutBuffer->activityReportIndUlMsg.totalReceiveMilliseconds = decodeUint32(ppInBuffer);
                        pOutBuffer->activityReportIndUlMsg.upTimeSeconds = decodeUint32(ppInBuffer);
                        pOutBuffer->activityReportIndUlMsg.txPowerDbm = (int8_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (pOutBuffer->activityReportIndUlMsg.txPowerDbm != BYTE_NOT_PRESENT_VALUE)
                        {
                            pOutBuffer->activityReportIndUlMsg.txPowerDbmPresent = true;
                        }
                        pOutBuffer->activityReportIndUlMsg.ulMcs = (uint8_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (pOutBuffer->activityReportIndUlMsg.ulMcs != BYTE_NOT_PRESENT_VALUE)
                        {
                            pOutBuffer->activityReportIndUlMsg.ulMcsPresent = true;
                        }
                        pOutBuffer->activityReportIndUlMsg.dlMcs = (uint8_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (pOutBuffer->activityReportIndUlMsg.dlMcs != BYTE_NOT_PRESENT_VALUE)
                        {
                            pOutBuffer->activityReportIndUlMsg.dlMcsPresent = true;
                        }
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ActivityReportIndUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logActivityReport(*ppLog, pLogSize, pOutBuffer->activityReportIndUlMsg.totalTransmitMilliseconds, pOutBuffer->activityReportIndUlMsg.totalReceiveMilliseconds,
                                pOutBuffer->activityReportIndUlMsg.upTimeSeconds);
                        *ppLog += logUlRf(*ppLog, pLogSize, pOutBuffer->activityReportIndUlMsg.txPowerDbmPresent,
                                pOutBuffer->activityReportIndUlMsg.txPowerDbm,
                                pOutBuffer->activityReportIndUlMsg.ulMcsPresent,
                                pOutBuffer->activityReportIndUlMsg.ulMcs);
                        *ppLog += logDlRf(*ppLog, pLogSize, pOutBuffer->activityReportIndUlMsg.dlMcsPresent, pOutBuffer->activityReportIndUlMsg.dlMcs);
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case ACTIVITY_REPORT_GET_CNF_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_ACTIVITY_REPORT_GET_CNF_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->activityReportGetCnfUlMsg.txPowerDbmPresent = false;
                        pOutBuffer->activityReportGetCnfUlMsg.ulMcsPresent = false;
                        pOutBuffer->activityReportGetCnfUlMsg.dlMcsPresent = false;
                        pOutBuffer->activityReportGetCnfUlMsg.totalTransmitMilliseconds = decodeUint32(ppInBuffer);
                        pOutBuffer->activityReportGetCnfUlMsg.totalReceiveMilliseconds = decodeUint32(ppInBuffer);
                        pOutBuffer->activityReportGetCnfUlMsg.upTimeSeconds = decodeUint32(ppInBuffer);
                        pOutBuffer->activityReportGetCnfUlMsg.txPowerDbm = (int8_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (pOutBuffer->activityReportGetCnfUlMsg.txPowerDbm != BYTE_NOT_PRESENT_VALUE)
                        {
                            pOutBuffer->activityReportGetCnfUlMsg.txPowerDbmPresent = true;
                        }
                        pOutBuffer->activityReportGetCnfUlMsg.ulMcs = (uint8_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (pOutBuffer->activityReportGetCnfUlMsg.ulMcs != BYTE_NOT_PRESENT_VALUE)
                        {
                            pOutBuffer->activityReportGetCnfUlMsg.ulMcsPresent = true;
                        }
                        pOutBuffer->activityReportGetCnfUlMsg.dlMcs = (uint8_t) **ppInBuffer;
                        (*ppInBuffer)++;
                        if (pOutBuffer->activityReportGetCnfUlMsg.dlMcs != BYTE_NOT_PRESENT_VALUE)
                        {
                            pOutBuffer->activityReportGetCnfUlMsg.dlMcsPresent = true;
                        }
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "ActivityReportGetCnflMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logActivityReport(*ppLog, pLogSize, pOutBuffer->activityReportGetCnfUlMsg.totalTransmitMilliseconds, pOutBuffer->activityReportGetCnfUlMsg.totalReceiveMilliseconds,
                                pOutBuffer->activityReportGetCnfUlMsg.upTimeSeconds);
                        *ppLog += logUlRf(*ppLog, pLogSize, pOutBuffer->activityReportGetCnfUlMsg.txPowerDbmPresent,
                                pOutBuffer->activityReportGetCnfUlMsg.txPowerDbm,
                                pOutBuffer->activityReportGetCnfUlMsg.ulMcsPresent,
                                pOutBuffer->activityReportGetCnfUlMsg.ulMcs);
                        *ppLog += logDlRf(*ppLog, pLogSize, pOutBuffer->activityReportGetCnfUlMsg.dlMcsPresent, pOutBuffer->activityReportGetCnfUlMsg.dlMcs);
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                case DEBUG_IND_UL_MSG:
                {
                    decodeResult = DECODE_RESULT_DEBUG_IND_UL_MSG;
                    if (pOutBuffer != NULL)
                    {
                        pOutBuffer->debugIndUlMsg.sizeOfString = **ppInBuffer;
                        (*ppInBuffer)++;
                        if (pOutBuffer->debugIndUlMsg.sizeOfString > MAX_DEBUG_STRING_SIZE)
                        {
                            pOutBuffer->debugIndUlMsg.sizeOfString = MAX_DEBUG_STRING_SIZE;
                        }
                        memcpy(&(pOutBuffer->debugIndUlMsg.string[0]), *(ppInBuffer), pOutBuffer->debugIndUlMsg.sizeOfString);
                        *(ppInBuffer) += pOutBuffer->debugIndUlMsg.sizeOfString;
                        if (calculateChecksum(pBufferAtStart, *ppInBuffer - pBufferAtStart) != **ppInBuffer)
                        {
                            decodeResult = DECODE_RESULT_BAD_CHECKSUM;
                        }
                        (*ppInBuffer)++;
                    }
                    if ((ppLog != NULL) && (*ppLog != NULL))
                    {
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_UL);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_NAME, "DebugIndUlMsg");
                        *ppLog += logBeginTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logDebugInd(*ppLog, pLogSize, &(pOutBuffer->debugIndUlMsg.string[0]), pOutBuffer->debugIndUlMsg.sizeOfString);
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_CONTENTS);
                        *ppLog += logTagWithUint32Value(*ppLog, pLogSize, TAG_MSG_SIZE, pBufferAtStart - *ppInBuffer);
                        *ppLog += logTagWithStringValue(*ppLog, pLogSize, TAG_MSG_CHECKSUM_GOOD, getStringBoolean(decodeResult != DECODE_RESULT_BAD_CHECKSUM));
                        *ppLog += logEndTag(*ppLog, pLogSize, TAG_MSG_UL);
                    }
                }
                break;
                default:
                    // The decodeResult will be left as Unknown message
                break;
            }
        }
    }

    return decodeResult;
}

// ----------------------------------------------------------------
// MISC FUNCTIONS
// ----------------------------------------------------------------

// Log debug messages
void logMsg(const char * pFormat, ...)
{
    char buffer[MAX_DEBUG_PRINTF_LEN];

    va_list args;
    va_start(args, pFormat);
    vsnprintf(buffer, sizeof(buffer), pFormat, args);
    va_end(args);
#ifdef WIN32  // Leave this out as I can't figure out to stop the C# app from garbage-collecting the pointer
    if (mp_guiPrintToConsole)
    {
        (*mp_guiPrintToConsole) (buffer);
    }
#else
    // Must be on ARM
    printf(buffer);
#endif
}

void initDll(void (*guiPrintToConsole)(const char *))
{
#ifdef WIN32
    mp_guiPrintToConsole = guiPrintToConsole;
    // Comment the generic call out and put in a specific call, see reason in
    // logMsg() function above
    // logMsg ("ready.\r\n");
    // This is the signal to the GUI that we're done with initialisation
    if (mp_guiPrintToConsole)
    {
        (*mp_guiPrintToConsole) ("ready.");
    }
#endif
}

#ifdef __cplusplus
}
#endif

// End Of File
