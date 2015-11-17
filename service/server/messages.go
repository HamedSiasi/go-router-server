/* Go versions of all the messages for the UTM server.
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
    "time"
)

// TODO
type DisplayRow struct {
    TotalMsgs          uint64     `json:"TotalMsgs,omitempty"`
    TotalBytes         uint64     `json:"TotalBytes,omitempty"`
    LastMsgReceived    *time.Time `json:"LastMsgReceived,omitempty"`
    Uuid               string     `json:"Uuid,omitempty"`
    UnitName           string     `json:"UnitName, omitempty"`
    Mode               string     `json:"Mode, omitempty"`
    UTotalMsgs         uint64     `json:"UTotalMsgs, omitempty"`
    UTotalBytes        uint64     `json:"UTotalBytes, omitempty"`
    UlastMsgReceived   *time.Time `json:"UlastMsgReceived, omitempty"`
    DTotalMsgs         uint64     `json:"DTotalMsgs, omitempty"`
    DTotalBytes        uint64     `json:"DTotalBytes, omitempty"`
    DlastMsgReceived   *time.Time `json:"DlastMsgReceived, omitempty"`
    RSRP               int32      `json:"RSRP, omitempty"`
    BatteryLevel       string     `json:"BatteryLevel, omitempty"`
    DiskSpaceLeft      string     `json:"DiskSpaceLeft, omitempty"`
    ReportingInterval  uint32     `json:"ReportingInterval, omitempty"`
    HeartbeatSeconds   uint32     `json:"HeartbeatSeconds, omitempty"`
    HeartbeatSnapToRtc bool       `json:"HeartbeatSnapToRtc, omitempty"`
}

type Connection struct {
    Status string
}

//--------------------------------------------------------------------
// Types used in messages with associated copy functions
//--------------------------------------------------------------------

// DateTime
type DateTime struct {
    Timestamp         time.Time
    UtmTime           time.Time
    TimeSetBy         TimeSetByEnum
}
func (value *DateTime) DeepCopy() *DateTime {
    if value == nil {
        return nil
    }
    result := &DateTime {
        Timestamp:         value.Timestamp,
        UtmTime:           value.UtmTime,
        TimeSetBy:         value.TimeSetBy,
    }
    return result
}

// GNSS
type GnssPosition struct {
    Timestamp time.Time
    Latitude  int32
    Longitude int32
    Elevation int32
}
func (value *GnssPosition) DeepCopy() *GnssPosition {
    if value == nil {
        return nil
    }
    result := &GnssPosition{
        Timestamp:   value.Timestamp,
        Latitude:    value.Latitude,
        Longitude:   value.Longitude,
        Elevation:   value.Elevation,
    }
    return result
}

// RSSI
type Rssi struct {
    Timestamp time.Time
    Rssi      int16
}
func (value *Rssi) DeepCopy() *Rssi {
    if value == nil {
        return nil
    }
    result := &Rssi{
        Timestamp: value.Timestamp,
        Rssi:      value.Rssi,
    }
    return result
}

// RSRP
type Rsrp struct {
    Value            Rssi
    IsSyncedWithRssi bool
}
func (value *Rsrp) DeepCopy() *Rsrp {
    if value == nil {
        return nil
    }
    result := &Rsrp{
        Value:            value.Value,
        IsSyncedWithRssi: value.IsSyncedWithRssi,
    }
    return result
}

// Power/State
type PowerState struct {
    Timestamp    time.Time
    ChargerState ChargerStateEnum
    BatteryMv    uint16
    EnergyMwh    uint32
}
func (value *PowerState) DeepCopy() *PowerState {
    if value == nil {
        return nil
    }
    result := &PowerState{
        Timestamp:    value.Timestamp,
        ChargerState: value.ChargerState,
        BatteryMv:    value.BatteryMv,
        EnergyMwh:    value.EnergyMwh,
    }
    return result
}

// Measurements data
type MeasurementData struct {
    Timestamp           time.Time
    TimeMeasured        time.Time
    GnssPositionPresent bool
    GnssPosition        GnssPosition
    CellIdPresent       bool
    CellId              uint16
    RsrpPresent         bool
    Rsrp                Rsrp
    RssiPresent         bool
    Rssi                Rssi
    TemperaturePresent  bool
    Temperature         int8
    PowerStatePresent   bool
    PowerState          PowerState
}
func (value *MeasurementData) DeepCopy() *MeasurementData {
    if value == nil {
        return nil
    }
    result := &MeasurementData{
        Timestamp:           value.Timestamp,
        TimeMeasured:        value.TimeMeasured,
        GnssPositionPresent: value.GnssPositionPresent,
        GnssPosition:        value.GnssPosition,
        CellIdPresent:       value.CellIdPresent,
        CellId:              value.CellId,
        RsrpPresent:         value.RsrpPresent,
        Rsrp:                value.Rsrp,
        RssiPresent:         value.RssiPresent,
        Rssi:                value.Rssi,
        TemperaturePresent:  value.TemperaturePresent,
        Temperature:         value.Temperature,
        PowerStatePresent:   value.PowerStatePresent,
        PowerState:          value.PowerState,
    }
    return result
}

// TODO
// type MeasurementControlGeneric struct
// type MeasurementControlPowerState struct
// type MeasurementControlUnion union

//--------------------------------------------------------------------
// Messages with associated copy functions
//--------------------------------------------------------------------

// TransparentUlDatagram
type TransparentUlDatagram struct {
    Timestamp    time.Time
    Contents     [MaxDatagramSizeRaw - 1]byte
}
func (value *TransparentUlDatagram) DeepCopy() *TransparentUlDatagram {
    if value == nil {
        return nil
    }
    result := &TransparentUlDatagram {
        Timestamp:   value.Timestamp,
        Contents:    value.Contents,
    }
    return result
}

// TransparentDlDatagram
type TransparentDlDatagram struct {
    Contents     [MaxDatagramSizeRaw - 1]byte
}
func (value *TransparentDlDatagram) DeepCopy() *TransparentDlDatagram {
    if value == nil {
        return nil
    }
    result := &TransparentDlDatagram {
        Contents:    value.Contents,
    }
    return result
}

// InitIndUlMsg
type InitIndUlMsg struct {
    Timestamp         time.Time
    WakeUpCode        WakeUpEnum
    RevisionLevel     uint8
    SdCardNotRequired bool
    DisableModemDebug bool
    DisableButton     bool
    DisableServerPing bool
}
func (value *InitIndUlMsg) DeepCopy() *InitIndUlMsg {
    if value == nil {
        return nil
    }
    result := &InitIndUlMsg {
        Timestamp:         value.Timestamp,
        WakeUpCode:        value.WakeUpCode,
        RevisionLevel:     value.RevisionLevel,
        SdCardNotRequired: value.SdCardNotRequired,
        DisableModemDebug: value.DisableModemDebug,
        DisableButton:     value.DisableButton,
        DisableServerPing: value.DisableServerPing,
    }
    return result
}

// RebootReqDlMsg
type RebootReqDlMsg struct {
    SdCardNotRequired bool
    DisableModemDebug bool
    DisableButton     bool
    DisableServerPing bool
}
func (value *RebootReqDlMsg) DeepCopy() *RebootReqDlMsg {
    if value == nil {
        return nil
    }
    result := &RebootReqDlMsg {
        SdCardNotRequired: value.SdCardNotRequired,
        DisableModemDebug: value.DisableModemDebug,
        DisableButton:     value.DisableButton,
        DisableServerPing: value.DisableServerPing,
    }
    return result
}

// DateTimeSetReqDlMsg
type DateTimeSetReqDlMsg struct {
    UtmTime           time.Time
    SetDateOnly       bool
}
func (value *DateTimeSetReqDlMsg) DeepCopy() *DateTimeSetReqDlMsg {
    if value == nil {
        return nil
    }
    result := &DateTimeSetReqDlMsg {
        UtmTime:           value.UtmTime,
        SetDateOnly:       value.SetDateOnly,
    }
    return result
}

// DateTimeSetCnfUlMsg
type DateTimeSetCnfUlMsg struct {
    Timestamp         time.Time
    UtmTime           time.Time
    TimeSetBy         TimeSetByEnum
}
func (value *DateTimeSetCnfUlMsg) DeepCopy() *DateTimeSetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &DateTimeSetCnfUlMsg {
        Timestamp:         value.Timestamp,
        UtmTime:           value.UtmTime,
        TimeSetBy:         value.TimeSetBy,
    }
    return result
}

// DateTimeGetCnfUlMsg
type DateTimeGetCnfUlMsg struct {
    Timestamp         time.Time
    UtmTime           time.Time
    TimeSetBy         TimeSetByEnum
}
func (value *DateTimeGetCnfUlMsg) DeepCopy() *DateTimeGetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &DateTimeGetCnfUlMsg {
        Timestamp:         value.Timestamp,
        UtmTime:           value.UtmTime,
        TimeSetBy:         value.TimeSetBy,
    }
    return result
}

// DateTimeIndUlMsg
type DateTimeIndUlMsg struct {
    Timestamp         time.Time
    UtmTime           time.Time
    TimeSetBy         TimeSetByEnum
}
func (value *DateTimeIndUlMsg) DeepCopy() *DateTimeIndUlMsg {
    if value == nil {
        return nil
    }
    result := &DateTimeIndUlMsg {
        Timestamp:         value.Timestamp,
        UtmTime:           value.UtmTime,
        TimeSetBy:         value.TimeSetBy,
    }
    return result
}

// ModeSetReqDlMsg
type ModeSetReqDlMsg struct {
    Mode        ModeEnum
}
func (value *ModeSetReqDlMsg) DeepCopy() *ModeSetReqDlMsg {
    if value == nil {
        return nil
    }
    result := &ModeSetReqDlMsg {
        Mode:         value.Mode,
    }
    return result
}

// ModeSetCnfUlMsg
type ModeSetCnfUlMsg struct {
    Timestamp   time.Time
    Mode        ModeEnum
}
func (value *ModeSetCnfUlMsg) DeepCopy() *ModeSetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &ModeSetCnfUlMsg {
        Timestamp:    value.Timestamp,
        Mode:         value.Mode,
    }
    return result
}

// ModeGetCnfUlMsg
type ModeGetCnfUlMsg struct {
    Timestamp   time.Time
    Mode        ModeEnum
}
func (value *ModeGetCnfUlMsg) DeepCopy() *ModeGetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &ModeGetCnfUlMsg {
        Timestamp:    value.Timestamp,
        Mode:         value.Mode,
    }
    return result
}

// HeartbeatSetReqDlMsg
type HeartbeatSetReqDlMsg struct {
    HeartbeatSeconds   uint32
    HeartbeatSnapToRtc bool
}
func (value *HeartbeatSetReqDlMsg) DeepCopy() *HeartbeatSetReqDlMsg {
    if value == nil {
        return nil
    }
    result := &HeartbeatSetReqDlMsg {
        HeartbeatSeconds:   value.HeartbeatSeconds,
        HeartbeatSnapToRtc: value.HeartbeatSnapToRtc,
    }
    return result
}

// HeartbeatSetCnfUlMsg
type HeartbeatSetCnfUlMsg struct {
    Timestamp          time.Time
    HeartbeatSeconds   uint32
    HeartbeatSnapToRtc bool
}
func (value *HeartbeatSetCnfUlMsg) DeepCopy() *HeartbeatSetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &HeartbeatSetCnfUlMsg {
        Timestamp:          value.Timestamp,
        HeartbeatSeconds:   value.HeartbeatSeconds,
        HeartbeatSnapToRtc: value.HeartbeatSnapToRtc,
    }
    return result
}

// ReportingIntervalSetReqDlMsg
type ReportingIntervalSetReqDlMsg struct {
    ReportingInterval uint32
}
func (value *ReportingIntervalSetReqDlMsg) DeepCopy() *ReportingIntervalSetReqDlMsg {
    if value == nil {
        return nil
    }
    result := &ReportingIntervalSetReqDlMsg {
        ReportingInterval:  value.ReportingInterval,
    }
    return result
}

// ReportingIntervalSetCnfUlMsg
type ReportingIntervalSetCnfUlMsg struct {
    Timestamp         time.Time
    ReportingInterval uint32
}
func (value *ReportingIntervalSetCnfUlMsg) DeepCopy() *ReportingIntervalSetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &ReportingIntervalSetCnfUlMsg {
        Timestamp:          value.Timestamp,
        ReportingInterval:  value.ReportingInterval,
    }
    return result
}

// IntervalsGetCnfUlMsg
type IntervalsGetCnfUlMsg struct {
    Timestamp          time.Time
    ReportingInterval  uint32
    HeartbeatSeconds   uint32
    HeartbeatSnapToRtc bool
}
func (value *IntervalsGetCnfUlMsg) DeepCopy() *IntervalsGetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &IntervalsGetCnfUlMsg {
        Timestamp:          value.Timestamp,
        ReportingInterval:  value.ReportingInterval,
        HeartbeatSeconds:   value.HeartbeatSeconds,
        HeartbeatSnapToRtc: value.HeartbeatSnapToRtc,
    }
    return result
}

// PollIndUlMsg
type PollIndUlMsg struct {
    Timestamp     time.Time
    Mode          ModeEnum
    EnergyLeft    EnergyLeftEnum
    DiskSpaceLeft DiskSpaceLeftEnum
}
func (value *PollIndUlMsg) DeepCopy() *PollIndUlMsg {
    if value == nil {
        return nil
    }
    result := &PollIndUlMsg {
        Timestamp:     value.Timestamp,
        Mode:          value.Mode,
        EnergyLeft:    value.EnergyLeft,
        DiskSpaceLeft: value.DiskSpaceLeft,
    }
    return result
}

// MeasurementsIndUlMsg
type MeasurementsIndUlMsg struct {
    Timestamp     time.Time
    Measurements  MeasurementData
}
func (value *MeasurementsIndUlMsg) DeepCopy() *MeasurementsIndUlMsg {
    if value == nil {
        return nil
    }
    result := &MeasurementsIndUlMsg {
        Timestamp:     value.Timestamp,
        Measurements:  value.Measurements,
    }
    return result
}

// MeasurementsGetCnfUlMsg
type MeasurementsGetCnfUlMsg struct {
    Timestamp     time.Time
    Measurements  MeasurementData
}
func (value *MeasurementsGetCnfUlMsg) DeepCopy() *MeasurementsGetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &MeasurementsGetCnfUlMsg {
        Timestamp:     value.Timestamp,
        Measurements:  value.Measurements,
    }
    return result
}

// TODO
// type MeasurementControlSetReqDlMsg struct
// type MeasurementControlSetCnfUlMsg struct
// type MeasurementControlGetCnfUlMsg struct
// type MeasurementControlIndUlMsg struct

// TrafficReportIndUlMsg
type TrafficReportIndUlMsg struct {
    Timestamp                 time.Time
    NumDatagramsUl            uint32
    NumBytesUl                uint32
    NumDatagramsDl            uint32
    NumBytesDl                uint32
    NumDatagramsDlBadChecksum uint32
}
func (value *TrafficReportIndUlMsg) DeepCopy() *TrafficReportIndUlMsg {
    if value == nil {
        return nil
    }
    result := &TrafficReportIndUlMsg {
        Timestamp:                  value.Timestamp,
        NumDatagramsUl:             value.NumDatagramsUl,
        NumBytesUl:                 value.NumBytesUl,
        NumDatagramsDl:             value.NumDatagramsDl,
        NumBytesDl:                 value.NumBytesDl,
        NumDatagramsDlBadChecksum:  value.NumDatagramsDlBadChecksum,
    }
    return result
}

// TrafficReportGetCnfUlMsg
type TrafficReportGetCnfUlMsg struct {
    Timestamp                 time.Time
    NumDatagramsUl            uint32
    NumBytesUl                uint32
    NumDatagramsDl            uint32
    NumBytesDl                uint32
    NumDatagramsDlBadChecksum uint32
}
func (value *TrafficReportGetCnfUlMsg) DeepCopy() *TrafficReportGetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &TrafficReportGetCnfUlMsg {
        Timestamp:                  value.Timestamp,
        NumDatagramsUl:             value.NumDatagramsUl,
        NumBytesUl:                 value.NumBytesUl,
        NumDatagramsDl:             value.NumDatagramsDl,
        NumBytesDl:                 value.NumBytesDl,
        NumDatagramsDlBadChecksum:  value.NumDatagramsDlBadChecksum,
    }
    return result
}

// DebugIndUlMsg
type DebugIndUlMsg struct {
    Timestamp     time.Time
    SizeOfString  uint8
    String        [MaxDatagramSizeRaw]byte
}
func (value *DebugIndUlMsg) DeepCopy() *DebugIndUlMsg {
    if value == nil {
        return nil
    }
    result := &DebugIndUlMsg {
        Timestamp:      value.Timestamp,
        SizeOfString:   value.SizeOfString,
        String:         value.String,
    }
    return result
}

// TrafficTestModeParametersSetReqDlMsg
type TrafficTestModeParametersSetReqDlMsg struct {
    NumUlDatagrams      uint32
    LenUlDatagram       uint32
    NumDlDatagrams      uint32
    LenDlDatagram       uint32
    TimeoutSeconds      uint32
    NoReportsDuringTest bool
}
func (value *TrafficTestModeParametersSetReqDlMsg) DeepCopy() *TrafficTestModeParametersSetReqDlMsg {
    if value == nil {
        return nil
    }
    result := &TrafficTestModeParametersSetReqDlMsg {
        NumUlDatagrams:      value.NumUlDatagrams,
        LenUlDatagram:       value.LenUlDatagram,
        NumDlDatagrams:      value.NumDlDatagrams,
        LenDlDatagram:       value.LenDlDatagram,
        TimeoutSeconds:      value.TimeoutSeconds,
        NoReportsDuringTest: value.NoReportsDuringTest,
    }
    return result
}

// TrafficTestModeParametersSetCnfUlMsg
type TrafficTestModeParametersSetCnfUlMsg struct {
    Timestamp           time.Time
    NumUlDatagrams      uint32
    LenUlDatagram       uint32
    NumDlDatagrams      uint32
    LenDlDatagram       uint32
    TimeoutSeconds      uint32
    NoReportsDuringTest bool
}
func (value *TrafficTestModeParametersSetCnfUlMsg) DeepCopy() *TrafficTestModeParametersSetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &TrafficTestModeParametersSetCnfUlMsg {
        Timestamp:           value.Timestamp,
        NumUlDatagrams:      value.NumUlDatagrams,
        LenUlDatagram:       value.LenUlDatagram,
        NumDlDatagrams:      value.NumDlDatagrams,
        LenDlDatagram:       value.LenDlDatagram,
        TimeoutSeconds:      value.TimeoutSeconds,
        NoReportsDuringTest: value.NoReportsDuringTest,
    }
    return result
}

// TrafficTestModeParametersGetCnfUlMsg
type TrafficTestModeParametersGetCnfUlMsg struct {
    Timestamp           time.Time
    NumUlDatagrams      uint32
    LenUlDatagram       uint32
    NumDlDatagrams      uint32
    LenDlDatagram       uint32
    TimeoutSeconds      uint32
    NoReportsDuringTest bool
}
func (value *TrafficTestModeParametersGetCnfUlMsg) DeepCopy() *TrafficTestModeParametersGetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &TrafficTestModeParametersGetCnfUlMsg {
        Timestamp:           value.Timestamp,
        NumUlDatagrams:      value.NumUlDatagrams,
        LenUlDatagram:       value.LenUlDatagram,
        NumDlDatagrams:      value.NumDlDatagrams,
        LenDlDatagram:       value.LenDlDatagram,
        TimeoutSeconds:      value.TimeoutSeconds,
        NoReportsDuringTest: value.NoReportsDuringTest,
    }
    return result
}

// TrafficTestModeRuleBreakerUlDatagram
type TrafficTestModeRuleBreakerUlDatagram struct {
    Timestamp  time.Time
    Fill       byte
    Length     uint32
}
func (value *TrafficTestModeRuleBreakerUlDatagram) DeepCopy() *TrafficTestModeRuleBreakerUlDatagram {
    if value == nil {
        return nil
    }
    result := &TrafficTestModeRuleBreakerUlDatagram {
        Timestamp:   value.Timestamp,
        Fill:        value.Fill,
        Length:      value.Length,
    }
    return result
}

// TrafficTestModeRuleBreakerDlDatagram
type TrafficTestModeRuleBreakerDlDatagram struct {
    Fill       byte
    Length     uint32
}
func (value *TrafficTestModeRuleBreakerDlDatagram) DeepCopy() *TrafficTestModeRuleBreakerDlDatagram {
    if value == nil {
        return nil
    }
    result := &TrafficTestModeRuleBreakerDlDatagram {
        Fill:        value.Fill,
        Length:      value.Length,
    }
    return result
}

// TrafficTestModeReportIndUlMsg
type TrafficTestModeReportIndUlMsg struct {
    Timestamp                            time.Time
    NumTrafficTestDatagramsUl            uint32
    NumTrafficTestBytesUl                uint32
    NumTrafficTestDatagramsDl            uint32
    NumTrafficTestBytesDl                uint32
    NumTrafficTestDlDatagramsOutOfOrder  uint32
    NumTrafficTestDlDatagramsBad         uint32
    NumTrafficTestDlDatagramsMissed      uint32
    TimedOut bool
}
func (value *TrafficTestModeReportIndUlMsg) DeepCopy() *TrafficTestModeReportIndUlMsg {
    if value == nil {
        return nil
    }
    result := &TrafficTestModeReportIndUlMsg {
        Timestamp:                           value.Timestamp,
        NumTrafficTestDatagramsUl:           value.NumTrafficTestDatagramsUl,
        NumTrafficTestBytesUl:               value.NumTrafficTestBytesUl,
        NumTrafficTestDatagramsDl:           value.NumTrafficTestDatagramsDl,
        NumTrafficTestBytesDl:               value.NumTrafficTestBytesDl,
        NumTrafficTestDlDatagramsOutOfOrder: value.NumTrafficTestDlDatagramsOutOfOrder,
        NumTrafficTestDlDatagramsBad:        value.NumTrafficTestDlDatagramsBad,
        NumTrafficTestDlDatagramsMissed:     value.NumTrafficTestDlDatagramsMissed,
        TimedOut:                            value.TimedOut,
    }
    return result
}

// TrafficTestModeReportGetCnfUlMsg
type TrafficTestModeReportGetCnfUlMsg struct {
    Timestamp                            time.Time
    NumTrafficTestDatagramsUl            uint32
    NumTrafficTestBytesUl                uint32
    NumTrafficTestDatagramsDl            uint32
    NumTrafficTestBytesDl                uint32
    NumTrafficTestDlDatagramsOutOfOrder  uint32
    NumTrafficTestDlDatagramsBad         uint32
    NumTrafficTestDlDatagramsMissed      uint32
    TimedOut bool
}
func (value *TrafficTestModeReportGetCnfUlMsg) DeepCopy() *TrafficTestModeReportGetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &TrafficTestModeReportGetCnfUlMsg {
        Timestamp:                           value.Timestamp,
        NumTrafficTestDatagramsUl:           value.NumTrafficTestDatagramsUl,
        NumTrafficTestBytesUl:               value.NumTrafficTestBytesUl,
        NumTrafficTestDatagramsDl:           value.NumTrafficTestDatagramsDl,
        NumTrafficTestBytesDl:               value.NumTrafficTestBytesDl,
        NumTrafficTestDlDatagramsOutOfOrder: value.NumTrafficTestDlDatagramsOutOfOrder,
        NumTrafficTestDlDatagramsBad:        value.NumTrafficTestDlDatagramsBad,
        NumTrafficTestDlDatagramsMissed:     value.NumTrafficTestDlDatagramsMissed,
        TimedOut:                            value.TimedOut,
    }
    return result
}

// ActivityReportIndUlMsg
type ActivityReportIndUlMsg struct {
    Timestamp                   time.Time
    TotalTransmitMilliseconds   uint32
    TotalReceiveMilliseconds    uint32
    UpTimeSeconds               uint32
    TxPowerDbmPresent           bool
    TxPowerDbm                  int8
    UlMcsPresent                bool
    UlMcs                       uint8
    DlMcsPresent                bool
    DlMcs                       uint8
}
func (value *ActivityReportIndUlMsg) DeepCopy() *ActivityReportIndUlMsg {
    if value == nil {
        return nil
    }
    result := &ActivityReportIndUlMsg {
        Timestamp:                  value.Timestamp,
        TotalTransmitMilliseconds:  value.TotalTransmitMilliseconds,
        TotalReceiveMilliseconds:   value.TotalReceiveMilliseconds,
        UpTimeSeconds:              value.UpTimeSeconds,
        TxPowerDbmPresent:          value.TxPowerDbmPresent,
        TxPowerDbm:                 value.TxPowerDbm,
        UlMcsPresent:               value.UlMcsPresent,
        UlMcs:                      value.UlMcs,
        DlMcsPresent:               value.DlMcsPresent,
        DlMcs:                      value.DlMcs,
    }
    return result
}

// ActivityReportGetCnfUlMsg
type ActivityReportGetCnfUlMsg struct {
    Timestamp                   time.Time
    TotalTransmitMilliseconds   uint32
    TotalReceiveMilliseconds    uint32
    UpTimeSeconds               uint32
    TxPowerDbmPresent           bool
    TxPowerDbm                  int8
    UlMcsPresent                bool
    UlMcs                       uint8
    DlMcsPresent                bool
    DlMcs                       uint8
}
func (value *ActivityReportGetCnfUlMsg) DeepCopy() *ActivityReportGetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &ActivityReportGetCnfUlMsg {
        Timestamp:                  value.Timestamp,
        TotalTransmitMilliseconds:  value.TotalTransmitMilliseconds,
        TotalReceiveMilliseconds:   value.TotalReceiveMilliseconds,
        UpTimeSeconds:              value.UpTimeSeconds,
        TxPowerDbmPresent:          value.TxPowerDbmPresent,
        TxPowerDbm:                 value.TxPowerDbm,
        UlMcsPresent:               value.UlMcsPresent,
        UlMcs:                      value.UlMcs,
        DlMcsPresent:               value.DlMcsPresent,
        DlMcs:                      value.DlMcs,
    }
    return result
}


// TODO
func (value *DisplayRow) DeepCopy() *DisplayRow {
    if value == nil {
        return nil
    }
    result := &DisplayRow{
        TotalMsgs:          value.TotalMsgs,
        TotalBytes:         value.TotalBytes,
        LastMsgReceived:    value.LastMsgReceived,
        Uuid:               value.Uuid,
        UnitName:           value.UnitName,
        Mode:               value.Mode,
        UTotalMsgs:         value.UTotalMsgs,
        UTotalBytes:        value.UTotalBytes,
        UlastMsgReceived:   value.UlastMsgReceived,
        DTotalMsgs:         value.DTotalMsgs,
        DTotalBytes:        value.DTotalBytes,
        DlastMsgReceived:   value.DlastMsgReceived,
        RSRP:               value.RSRP,
        BatteryLevel:       value.BatteryLevel,
        DiskSpaceLeft:      value.DiskSpaceLeft,
        ReportingInterval:  value.ReportingInterval,
        HeartbeatSeconds:   value.HeartbeatSeconds,
        HeartbeatSnapToRtc: value.HeartbeatSnapToRtc,
    }
    return result
}

/* // TODO
func (value *MeasurementsIndUlMsg) DeepCopy() *MeasurementsIndUlMsg {
    if value == nil {
        return nil
    }
    result := &MeasurementsIndUlMsg{
    // rsrp: value.Timestamp,
    // rssi: value.ReportingIntervalMinutes,
    }
    return result
}
*/

type DataVolume struct {
    UplinkTimestamp   *time.Time
    DownlinkTimestamp *time.Time
    UplinkBytes       uint64
    DownlinkBytes     uint64
}


func (value *DataVolume) DeepCopy() *DataVolume {
    if value == nil {
        return nil
    }
    result := &DataVolume{
        UplinkTimestamp:   value.UplinkTimestamp,
        DownlinkTimestamp: value.DownlinkTimestamp,
        UplinkBytes:       value.UplinkBytes,
        DownlinkBytes:     value.DownlinkBytes,
    }
    return result
}

/* End Of File */
