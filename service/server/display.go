/* Functions that populate displayable values for the UTM server.
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

// Display InitInd
type InitIndDisplay struct {
    Timestamp         time.Time
    WakeUpCode        string
    RevisionLevel     uint8
    SdCardNotRequired bool
    DisableModemDebug bool
    DisableButton     bool
    DisableServerPing bool
}
func (value *InitIndDisplay) DeepCopy() *InitIndDisplay {
    if value == nil {
        return nil
    }
    result := &InitIndDisplay {
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
func makeInitIndDisplay(newData *InitIndUlMsg) *InitIndDisplay {
    display := InitIndDisplay {
        Timestamp:  	   newData.Timestamp,
        WakeUpCode:    	   WakeUpCodeLookUp[newData.WakeUpCode],
	    RevisionLevel:     newData.RevisionLevel,
	    SdCardNotRequired: newData.SdCardNotRequired,
	    DisableModemDebug: newData.DisableModemDebug,
	    DisableButton:     newData.DisableButton,
	    DisableServerPing: newData.DisableServerPing,
    }
    return &display
}

// Display intervals
type IntervalsDisplay struct {
    Timestamp                time.Time
    ReportingIntervalPresent bool
    ReportingInterval        uint32
    HeartbeatPresent         bool
    HeartbeatSeconds         uint32
    HeartbeatSnapToRtc       bool
}
func (value *IntervalsDisplay) DeepCopy() *IntervalsDisplay {
    if value == nil {
        return nil
    }
    result := &IntervalsDisplay {
        Timestamp:                value.Timestamp,
        ReportingIntervalPresent: value.ReportingIntervalPresent,
        ReportingInterval:        value.ReportingInterval,
        HeartbeatPresent:         value.HeartbeatPresent,
        HeartbeatSeconds:         value.HeartbeatSeconds,
        HeartbeatSnapToRtc:       value.HeartbeatSnapToRtc,
    }
    return result
}
func makeIntervalsDisplay0(newData *IntervalsGetCnfUlMsg) *IntervalsDisplay {
    display := IntervalsDisplay {
        Timestamp:                newData.Timestamp,
        ReportingIntervalPresent: true,
        ReportingInterval:        newData.ReportingInterval,
        HeartbeatPresent:         true,
        HeartbeatSeconds:         newData.HeartbeatSeconds,
        HeartbeatSnapToRtc:       newData.HeartbeatSnapToRtc,
    }
    return &display
}
func makeIntervalsDisplay1(newData *ReportingIntervalSetCnfUlMsg) *IntervalsDisplay {
    display := IntervalsDisplay {
        Timestamp:                newData.Timestamp,
        ReportingIntervalPresent: true,
        ReportingInterval:        newData.ReportingInterval,
    }
    return &display
}
func makeIntervalsDisplay2(newData *HeartbeatSetCnfUlMsg) *IntervalsDisplay {
    display := IntervalsDisplay {
        Timestamp:            newData.Timestamp,
        HeartbeatPresent:     true,
        HeartbeatSeconds:     newData.HeartbeatSeconds,
        HeartbeatSnapToRtc:   newData.HeartbeatSnapToRtc,
    }
    return &display
}

// Display mode
type ModeDisplay struct {
    Timestamp    time.Time
    Mode         string
}
func (value *ModeDisplay) DeepCopy() *ModeDisplay {
    if value == nil {
        return nil
    }
    result := &ModeDisplay {
        Timestamp:    value.Timestamp,
        Mode:         value.Mode,
    }
    return result
}
func makeModeDisplay0(newData *ModeSetCnfUlMsg) *ModeDisplay {
    display := ModeDisplay {
        Timestamp:    newData.Timestamp,
        Mode:         ModeLookUp[newData.Mode],
    }
    return &display
}
func makeModeDisplay1(newData *ModeGetCnfUlMsg) *ModeDisplay {
    display := ModeDisplay {
        Timestamp:    newData.Timestamp,
        Mode:         ModeLookUp[newData.Mode],
    }
    return &display
}
func makeModeDisplay2(newData *PollIndUlMsg) *ModeDisplay {
    display := ModeDisplay {
        Timestamp:    newData.Timestamp,
        Mode:         ModeLookUp[newData.Mode],
    }
    return &display
}

// Display date/time
type DateTimeDisplay struct {
    Timestamp       time.Time
    UtmTime         time.Time
    TimeSetBy       string
}
func (value *DateTimeDisplay) DeepCopy() *DateTimeDisplay {
    if value == nil {
        return nil
    }
    result := &DateTimeDisplay {
        Timestamp:       value.Timestamp,
        UtmTime:         value.UtmTime,
        TimeSetBy:       value.TimeSetBy,
    }
    return result
}
func makeDateTimeDisplay0(newData *DateTimeIndUlMsg) *DateTimeDisplay {
    display := DateTimeDisplay {
        Timestamp:       newData.Timestamp,
        UtmTime:         newData.UtmTime,
        TimeSetBy:       TimeSetByLookUp[newData.TimeSetBy],
    }
    return &display
}
func makeDateTimeDisplay1(newData *DateTimeSetCnfUlMsg) *DateTimeDisplay {
    display := DateTimeDisplay {
        Timestamp:       newData.Timestamp,
        UtmTime:         newData.UtmTime,
        TimeSetBy:       TimeSetByLookUp[newData.TimeSetBy],
    }
    return &display
}
func makeDateTimeDisplay2(newData *DateTimeGetCnfUlMsg) *DateTimeDisplay {
    display := DateTimeDisplay {
        Timestamp:       newData.Timestamp,
        UtmTime:         newData.UtmTime,
        TimeSetBy:       TimeSetByLookUp[newData.TimeSetBy],
    }
    return &display
}

// Display UtmStatus
type UtmStatusDisplay struct {
    Timestamp     time.Time
    EnergyLeft    string
    DiskSpaceLeft string
}
func (value *UtmStatusDisplay) DeepCopy() *UtmStatusDisplay {
    if value == nil {
        return nil
    }
    result := &UtmStatusDisplay {
        Timestamp:         value.Timestamp,
        EnergyLeft:        value.EnergyLeft,
        DiskSpaceLeft:     value.DiskSpaceLeft,
    }
    return result
}
func makeUtmStatusDisplay(newData *PollIndUlMsg) *UtmStatusDisplay {
    display := UtmStatusDisplay {
        Timestamp:  	   newData.Timestamp,
        EnergyLeft:    	   EnergyLeftLookUp[newData.EnergyLeft],
        DiskSpaceLeft:     DiskSpaceLeftLookUp[newData.DiskSpaceLeft],
    }
    return &display
}

// Display GNSS
type GnssDisplay struct {
    Timestamp    time.Time
    Latitude     float32
    Longitude    float32
    Elevation    float32
}
func (value *GnssDisplay) DeepCopy() *GnssDisplay {
    if value == nil {
        return nil
    }
    result := &GnssDisplay {
        Timestamp:     value.Timestamp,
        Latitude:      value.Latitude,
        Longitude:     value.Longitude,
        Elevation:     value.Elevation,
    }
    return result
}
func makeGnssDisplay(newData *MeasurementData) *GnssDisplay {
    display := GnssDisplay{
        Timestamp:    newData.Timestamp,
        Latitude:     float32(newData.GnssPosition.Latitude) / 1000,
        Longitude:    float32(newData.GnssPosition.Longitude) / 1000,
        Elevation:    float32(newData.GnssPosition.Elevation),
    }
    return &display
}

// Display cell ID
type CellIdDisplay struct {
    Timestamp    time.Time
    CellId       uint32
}
func (value *CellIdDisplay) DeepCopy() *CellIdDisplay {
    if value == nil {
        return nil
    }
    result := &CellIdDisplay {
        Timestamp:  value.Timestamp,
        CellId:     value.CellId,
    }
    return result
}
func makeCellIdDisplay(newData *MeasurementData) *CellIdDisplay {
    display := CellIdDisplay{
        Timestamp:    newData.Timestamp,
        CellId:       uint32(newData.CellId),
    }
    return &display
}

// Display temperature
type TemperatureDisplay struct {
    Timestamp    time.Time
    TemperatureC float32
}
func (value *TemperatureDisplay) DeepCopy() *TemperatureDisplay {
    if value == nil {
        return nil
    }
    result := &TemperatureDisplay {
        Timestamp:     value.Timestamp,
        TemperatureC:  value.TemperatureC,
    }
    return result
}
func makeTemperatureDisplay(newData *MeasurementData) *TemperatureDisplay {
    display := TemperatureDisplay{
        Timestamp:    newData.Timestamp,
        TemperatureC: float32(newData.Temperature),
    }
    return &display
}

// Display signal strength items
type SignalStrengthDisplay struct {
    Timestamp         time.Time
    RsrpDbm           float32
    Mcl               float32
    HighestMcl        float32
    InsideGsmCoverage bool
    RssiDbmPresent    bool
    RssiDbm           float32
    SnrPresent        bool
    Snr               float32
}
func (value *SignalStrengthDisplay) DeepCopy() *SignalStrengthDisplay {
    if value == nil {
        return nil
    }
    result := &SignalStrengthDisplay{
        Timestamp:         value.Timestamp,
        RsrpDbm:           value.RsrpDbm,
        Mcl:               value.Mcl,
        HighestMcl:        value.HighestMcl,
        InsideGsmCoverage: value.InsideGsmCoverage,
        RssiDbmPresent:    value.RssiDbmPresent,
        RssiDbm:           value.RssiDbm,
        SnrPresent:        value.SnrPresent,
        Snr:               value.Snr,
    }
    return result
} 
func makeSignalStrengthDisplay(newData *MeasurementData) *SignalStrengthDisplay {
    display := SignalStrengthDisplay{
        Timestamp:      newData.Timestamp,
        RsrpDbm:        float32 (newData.Rsrp.Value.Rssi) / 10,
        RssiDbmPresent: newData.RssiPresent,
        RssiDbm:        float32 (newData.Rssi.Rssi) / 10,
    }
    
    // TODO: calculate the answer

    return &display
}

// Display power/state items
type PowerStateDisplay struct {
    Timestamp       time.Time
    ChargeState     string
    BatteryVoltageV float32
    EnergyLeftWh    float32
}

func (value *PowerStateDisplay) DeepCopy() *PowerStateDisplay {
    if value == nil {
        return nil
    }
    result := &PowerStateDisplay{
        Timestamp:       value.Timestamp,
        ChargeState:     value.ChargeState,
        BatteryVoltageV: value.BatteryVoltageV,
        EnergyLeftWh:    value.EnergyLeftWh,
    }
    return result
}
func makePowerStateDisplay(newData *MeasurementData) *PowerStateDisplay {
    // Change nothing if the new data is nil
    if newData == nil {
        return nil
    }
    // Overwrite if there is no stored data
    display := &PowerStateDisplay{
        Timestamp:       newData.Timestamp,
        ChargeState:     ChargerStateEnumLookUp[newData.PowerState.ChargerState],
        BatteryVoltageV: float32(newData.PowerState.BatteryMv) / 1000,
        EnergyLeftWh:    float32(newData.PowerState.EnergyMwh) / 1000,
    }
    return display
}

// Display traffic information
type TrafficReportDisplay struct {
    Timestamp                 time.Time
    NumDatagramsUl            uint32
    NumBytesUl                uint32
    NumDatagramsDl            uint32
    NumBytesDl                uint32
    NumDatagramsDlBadChecksum uint32
}
func (value *TrafficReportDisplay) DeepCopy() *TrafficReportDisplay {
    if value == nil {
        return nil
    }
    result := &TrafficReportDisplay {
        Timestamp:                  value.Timestamp,
        NumDatagramsUl:             value.NumDatagramsUl,
        NumBytesUl:                 value.NumBytesUl,
        NumDatagramsDl:             value.NumDatagramsDl,
        NumBytesDl:                 value.NumBytesDl,
        NumDatagramsDlBadChecksum:  value.NumDatagramsDlBadChecksum,
    }
    return result
}
func makeTrafficReportDisplay0(newData *TrafficReportIndUlMsg) *TrafficReportDisplay {
    display := TrafficReportDisplay {
        Timestamp:                  newData.Timestamp,
        NumDatagramsUl:             newData.NumDatagramsUl,
        NumBytesUl:                 newData.NumBytesUl,
        NumDatagramsDl:             newData.NumDatagramsDl,
        NumBytesDl:                 newData.NumBytesDl,
        NumDatagramsDlBadChecksum:  newData.NumDatagramsDlBadChecksum,
    }
    return &display
}
func makeTrafficReportDisplay1(newData *TrafficReportGetCnfUlMsg) *TrafficReportDisplay {
    display := TrafficReportDisplay {
        Timestamp:                  newData.Timestamp,
        NumDatagramsUl:             newData.NumDatagramsUl,
        NumBytesUl:                 newData.NumBytesUl,
        NumDatagramsDl:             newData.NumDatagramsDl,
        NumBytesDl:                 newData.NumBytesDl,
        NumDatagramsDlBadChecksum:  newData.NumDatagramsDlBadChecksum,
    }
    return &display
}

// Display traffic test mode parameter information
type TrafficTestModeParametersDisplay struct {
    Timestamp           time.Time
    NumUlDatagrams      uint32
    LenUlDatagram       uint32
    NumDlDatagrams      uint32
    LenDlDatagram       uint32
    TimeoutSeconds      uint32
    NoReportsDuringTest bool
}
func (value *TrafficTestModeParametersDisplay) DeepCopy() *TrafficTestModeParametersDisplay {
    if value == nil {
        return nil
    }
    result := &TrafficTestModeParametersDisplay {
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
func makeTrafficTestModeParametersDisplay0(newData *TrafficTestModeParametersSetCnfUlMsg) *TrafficTestModeParametersDisplay {
    display := TrafficTestModeParametersDisplay {
        Timestamp:           newData.Timestamp,
        NumUlDatagrams:      newData.NumUlDatagrams,
        LenUlDatagram:       newData.LenUlDatagram,
        NumDlDatagrams:      newData.NumDlDatagrams,
        LenDlDatagram:       newData.LenDlDatagram,
        TimeoutSeconds:      newData.TimeoutSeconds,
        NoReportsDuringTest: newData.NoReportsDuringTest,
    }
    return &display
}
func makeTrafficTestModeParametersDisplay1(newData *TrafficTestModeParametersGetCnfUlMsg) *TrafficTestModeParametersDisplay {
    display := TrafficTestModeParametersDisplay {
        Timestamp:           newData.Timestamp,
        NumUlDatagrams:      newData.NumUlDatagrams,
        LenUlDatagram:       newData.LenUlDatagram,
        NumDlDatagrams:      newData.NumDlDatagrams,
        LenDlDatagram:       newData.LenDlDatagram,
        TimeoutSeconds:      newData.TimeoutSeconds,
        NoReportsDuringTest: newData.NoReportsDuringTest,
    }
    return &display
}

// Display traffic test mode information
type TrafficTestModeReportDisplay struct {
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
func (value *TrafficTestModeReportDisplay) DeepCopy() *TrafficTestModeReportDisplay {
    if value == nil {
        return nil
    }
    result := &TrafficTestModeReportDisplay {
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
func makeTrafficTestModeReportDisplay0(newData *TrafficTestModeReportIndUlMsg) *TrafficTestModeReportDisplay {
    display := TrafficTestModeReportDisplay {
        Timestamp:                           newData.Timestamp,
        NumTrafficTestDatagramsUl:           newData.NumTrafficTestDatagramsUl,
        NumTrafficTestBytesUl:               newData.NumTrafficTestBytesUl,
        NumTrafficTestDatagramsDl:           newData.NumTrafficTestDatagramsDl,
        NumTrafficTestBytesDl:               newData.NumTrafficTestBytesDl,
        NumTrafficTestDlDatagramsOutOfOrder: newData.NumTrafficTestDlDatagramsOutOfOrder,
        NumTrafficTestDlDatagramsBad:        newData.NumTrafficTestDlDatagramsBad,
        NumTrafficTestDlDatagramsMissed:     newData.NumTrafficTestDlDatagramsMissed,
        TimedOut:                            newData.TimedOut,
    }
    return &display
}
func makeTrafficTestModeReportDisplay1(newData *TrafficTestModeReportGetCnfUlMsg) *TrafficTestModeReportDisplay {
    display := TrafficTestModeReportDisplay {
        Timestamp:                           newData.Timestamp,
        NumTrafficTestDatagramsUl:           newData.NumTrafficTestDatagramsUl,
        NumTrafficTestBytesUl:               newData.NumTrafficTestBytesUl,
        NumTrafficTestDatagramsDl:           newData.NumTrafficTestDatagramsDl,
        NumTrafficTestBytesDl:               newData.NumTrafficTestBytesDl,
        NumTrafficTestDlDatagramsOutOfOrder: newData.NumTrafficTestDlDatagramsOutOfOrder,
        NumTrafficTestDlDatagramsBad:        newData.NumTrafficTestDlDatagramsBad,
        NumTrafficTestDlDatagramsMissed:     newData.NumTrafficTestDlDatagramsMissed,
        TimedOut:                            newData.TimedOut,
    }
    return &display
}

// Display activity report information
type ActivityReportDisplay struct {
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
func (value *ActivityReportDisplay) DeepCopy() *ActivityReportDisplay {
    if value == nil {
        return nil
    }
    result := &ActivityReportDisplay {
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
func makeActivityReportDisplay0(newData *ActivityReportIndUlMsg) *ActivityReportDisplay {
    display := ActivityReportDisplay {
        Timestamp:                  newData.Timestamp,
        TotalTransmitMilliseconds:  newData.TotalTransmitMilliseconds,
        TotalReceiveMilliseconds:   newData.TotalReceiveMilliseconds,
        UpTimeSeconds:              newData.UpTimeSeconds,
        TxPowerDbmPresent:          newData.TxPowerDbmPresent,
        TxPowerDbm:                 newData.TxPowerDbm,
        UlMcsPresent:               newData.UlMcsPresent,
        UlMcs:                      newData.UlMcs,
        DlMcsPresent:               newData.DlMcsPresent,
        DlMcs:                      newData.DlMcs,
    }
    return &display
}
func makeActivityReportDisplay1(newData *ActivityReportGetCnfUlMsg) *ActivityReportDisplay {
    display := ActivityReportDisplay {
        Timestamp:                  newData.Timestamp,
        TotalTransmitMilliseconds:  newData.TotalTransmitMilliseconds,
        TotalReceiveMilliseconds:   newData.TotalReceiveMilliseconds,
        UpTimeSeconds:              newData.UpTimeSeconds,
        TxPowerDbmPresent:          newData.TxPowerDbmPresent,
        TxPowerDbm:                 newData.TxPowerDbm,
        UlMcsPresent:               newData.UlMcsPresent,
        UlMcs:                      newData.UlMcs,
        DlMcsPresent:               newData.DlMcsPresent,
        DlMcs:                      newData.DlMcs,
    }
    return &display
}

/* End Of File */
