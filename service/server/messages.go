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

//--------------------------------------------------------------------
// Types used in messages with associated copy functions
//--------------------------------------------------------------------

// DateTime
type DateTime struct {
    UtmTime           time.Time
    TimeSetBy         TimeSetByEnum
}
func (value *DateTime) DeepCopy() *DateTime {
    if value == nil {
        return nil
    }
    result := &DateTime {
        UtmTime:           value.UtmTime,
        TimeSetBy:         value.TimeSetBy,
    }
    return result
}

// GNSS
type GnssPosition struct {
    Latitude  int32
    Longitude int32
    Elevation int32
}
func (value *GnssPosition) DeepCopy() *GnssPosition {
    if value == nil {
        return nil
    }
    result := &GnssPosition{
        Latitude:    value.Latitude,
        Longitude:   value.Longitude,
        Elevation:   value.Elevation,
    }
    return result
}

// CellId
type CellId uint16

// RSSI
type Rssi int16

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

// Temperature
type Temperature int8

// Power/State
type PowerState struct {
    ChargerState ChargerStateEnum
    BatteryMv    uint16
    EnergyMwh    uint32
}
func (value *PowerState) DeepCopy() *PowerState {
    if value == nil {
        return nil
    }
    result := &PowerState{
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
    CellId              CellId
    RsrpPresent         bool
    Rsrp                Rsrp
    RssiPresent         bool
    Rssi                Rssi
    TemperaturePresent  bool
    Temperature         Temperature
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
    Contents     [MaxDatagramSizeRaw - 1]byte
}
func (value *TransparentUlDatagram) DeepCopy() *TransparentUlDatagram {
    if value == nil {
        return nil
    }
    result := &TransparentUlDatagram {
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

// PingReqUlMsg
type PingReqUlMsg struct {
    // Empty
}

// PingCnfDlMsg
type PingCnfDlMsg struct {
    // Empty
}

// PingReqDlMsg
type PingReqDlMsg struct {
    // Empty
}

// PingCnfUlMsg
type PingCnfUlMsg struct {
    // Empty
}

// InitIndUlMsg
type InitIndUlMsg struct {
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
    UtmTime           time.Time
    TimeSetBy         TimeSetByEnum
}
func (value *DateTimeSetCnfUlMsg) DeepCopy() *DateTimeSetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &DateTimeSetCnfUlMsg {
        UtmTime:           value.UtmTime,
        TimeSetBy:         value.TimeSetBy,
    }
    return result
}

// DateTimeGetReqDlMsg
type DateTimeGetReqDlMsg struct {
    // Empty    
}

// DateTimeGetCnfUlMsg
type DateTimeGetCnfUlMsg struct {
    UtmTime           time.Time
    TimeSetBy         TimeSetByEnum
}
func (value *DateTimeGetCnfUlMsg) DeepCopy() *DateTimeGetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &DateTimeGetCnfUlMsg {
        UtmTime:           value.UtmTime,
        TimeSetBy:         value.TimeSetBy,
    }
    return result
}

// DateTimeIndUlMsg
type DateTimeIndUlMsg struct {
    UtmTime           time.Time
    TimeSetBy         TimeSetByEnum
}
func (value *DateTimeIndUlMsg) DeepCopy() *DateTimeIndUlMsg {
    if value == nil {
        return nil
    }
    result := &DateTimeIndUlMsg {
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
    Mode        ModeEnum
}
func (value *ModeSetCnfUlMsg) DeepCopy() *ModeSetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &ModeSetCnfUlMsg {
        Mode:         value.Mode,
    }
    return result
}

// ModeGetReqDlMsg
type ModeGetReqDlMsg struct {
    // Empty    
}

// ModeGetCnfUlMsg
type ModeGetCnfUlMsg struct {
    Mode        ModeEnum
}
func (value *ModeGetCnfUlMsg) DeepCopy() *ModeGetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &ModeGetCnfUlMsg {
        Mode:         value.Mode,
    }
    return result
}

// IntervalsGetReqDlMsg
type IntervalsGetReqDlMsg struct {
    // Empty    
}

// IntervalsGetCnfUlMsg
type IntervalsGetCnfUlMsg struct {
    ReportingInterval  uint32
    HeartbeatSeconds   uint32
    HeartbeatSnapToRtc bool
}
func (value *IntervalsGetCnfUlMsg) DeepCopy() *IntervalsGetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &IntervalsGetCnfUlMsg {
        ReportingInterval:  value.ReportingInterval,
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
    ReportingInterval uint32
}
func (value *ReportingIntervalSetCnfUlMsg) DeepCopy() *ReportingIntervalSetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &ReportingIntervalSetCnfUlMsg {
        ReportingInterval:  value.ReportingInterval,
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
    HeartbeatSeconds   uint32
    HeartbeatSnapToRtc bool
}
func (value *HeartbeatSetCnfUlMsg) DeepCopy() *HeartbeatSetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &HeartbeatSetCnfUlMsg {
        HeartbeatSeconds:   value.HeartbeatSeconds,
        HeartbeatSnapToRtc: value.HeartbeatSnapToRtc,
    }
    return result
}

// PollIndUlMsg
type PollIndUlMsg struct {
    Mode          ModeEnum
    EnergyLeft    EnergyLeftEnum
    DiskSpaceLeft DiskSpaceLeftEnum
}
func (value *PollIndUlMsg) DeepCopy() *PollIndUlMsg {
    if value == nil {
        return nil
    }
    result := &PollIndUlMsg {
        Mode:          value.Mode,
        EnergyLeft:    value.EnergyLeft,
        DiskSpaceLeft: value.DiskSpaceLeft,
    }
    return result
}

// MeasurementsIndUlMsg
type MeasurementsIndUlMsg struct {
    Measurements  MeasurementData
}
func (value *MeasurementsIndUlMsg) DeepCopy() *MeasurementsIndUlMsg {
    if value == nil {
        return nil
    }
    result := &MeasurementsIndUlMsg {
        Measurements:  value.Measurements,
    }
    return result
}

// MeasurementsGetReqDlMsg
type MeasurementsGetReqDlMsg struct {
    // Empty
}

// MeasurementsGetCnfUlMsg
type MeasurementsGetCnfUlMsg struct {
    Measurements  MeasurementData
}
func (value *MeasurementsGetCnfUlMsg) DeepCopy() *MeasurementsGetCnfUlMsg {
    if value == nil {
        return nil
    }
    result := &MeasurementsGetCnfUlMsg {
        Measurements:  value.Measurements,
    }
    return result
}

// TODO
// type MeasurementControlSetReqDlMsg struct
// type MeasurementControlSetCnfUlMsg struct
// type MeasurementControlGetCnfUlMsg struct
// type MeasurementControlIndUlMsg struct
// type MeasurementsControlDefaultsSetReqDlMsg struct
// type MeasurementsControlDefaultsSetCnfUlMsg struct

// TrafficReportIndUlMsg
type TrafficReportIndUlMsg struct {
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
        NumDatagramsUl:             value.NumDatagramsUl,
        NumBytesUl:                 value.NumBytesUl,
        NumDatagramsDl:             value.NumDatagramsDl,
        NumBytesDl:                 value.NumBytesDl,
        NumDatagramsDlBadChecksum:  value.NumDatagramsDlBadChecksum,
    }
    return result
}

// TrafficReportGetReqDlMsg
type TrafficReportGetReqDlMsg struct {
    // Empty
}

// TrafficReportGetCnfUlMsg
type TrafficReportGetCnfUlMsg struct {
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
    SizeOfString  uint8
    String        [MaxDatagramSizeRaw]byte
}
func (value *DebugIndUlMsg) DeepCopy() *DebugIndUlMsg {
    if value == nil {
        return nil
    }
    result := &DebugIndUlMsg {
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
        NumUlDatagrams:      value.NumUlDatagrams,
        LenUlDatagram:       value.LenUlDatagram,
        NumDlDatagrams:      value.NumDlDatagrams,
        LenDlDatagram:       value.LenDlDatagram,
        TimeoutSeconds:      value.TimeoutSeconds,
        NoReportsDuringTest: value.NoReportsDuringTest,
    }
    return result
}

// TrafficTestModeParametersGetReqDlMsg
type TrafficTestModeParametersGetReqDlMsg struct {
    // Empty
}

// TrafficTestModeParametersGetCnfUlMsg
type TrafficTestModeParametersGetCnfUlMsg struct {
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
    Fill       byte
    Length     uint32
}
func (value *TrafficTestModeRuleBreakerUlDatagram) DeepCopy() *TrafficTestModeRuleBreakerUlDatagram {
    if value == nil {
        return nil
    }
    result := &TrafficTestModeRuleBreakerUlDatagram {
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

// TrafficTestModeReportGetReqDlMsg
type TrafficTestModeReportGetReqDlMsg struct {
    // Empty
}

// TrafficTestModeReportGetCnfUlMsg
type TrafficTestModeReportGetCnfUlMsg struct {
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

// ActivityReportGetReqDlMsg
type ActivityReportGetReqDlMsg struct {
    // Empty
}

// ActivityReportGetCnfUlMsg
type ActivityReportGetCnfUlMsg struct {
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
