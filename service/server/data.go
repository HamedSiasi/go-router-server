/* Functions that populate data items values for the UTM server.
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

// Storage for the "interesting" flag
// This shoudl be set if something other than a boring RF or
// status update has arrived from the device, something which
// requires the user to go beyond the front page
type Interesting struct {
    Timestamp      time.Time    
    IsInteresting  bool
}
func (value *Interesting) DeepCopy() *Interesting {
    if value == nil {
        return nil
    }
    result := &Interesting {
        Timestamp:       value.Timestamp,  
        IsInteresting:   value.IsInteresting,
    }
    return result
}

// Storage for the expected message list
type ExpectedMsgData struct {    
    Timestamp      time.Time    
    ExpectedMsgList  *[]ExpectedMsg
}
func (value *ExpectedMsgData) DeepCopy() *ExpectedMsgData {
    if value == nil {
        return nil
    }
    expectedMsgList := make ([]ExpectedMsg, 0)
    result := &ExpectedMsgData {
        Timestamp:         value.Timestamp,  
        ExpectedMsgList:   &expectedMsgList,
    }
    if value.ExpectedMsgList != nil {
        for _, expectedMsg := range *value.ExpectedMsgList {
            expectedMsgList = append (expectedMsgList, expectedMsg)    
        }
    }
        
    return result
}
func makeExpectedMsgData(newData *[]ExpectedMsg, Time time.Time) *ExpectedMsgData {
    expectedMsgList := make ([]ExpectedMsg, 0)
    data := ExpectedMsgData {
        Timestamp:         Time,
        ExpectedMsgList:   &expectedMsgList,
    }
    if newData != nil {
        for _, expectedMsg := range *newData {
            expectedMsgList = append (expectedMsgList, expectedMsg)    
        }
    }
    return &data
}

// Storage for traffic volume data
type TrafficVolumeData struct {
    UlMsgs        int
    UlBytes       int
    LastUlMsgTime time.Time
    DlMsgs        int
    DlBytes       int
    LastDlMsgTime time.Time
    UlTotals      *TotalsState 
    DlTotals      *TotalsState    
}
func (value *TrafficVolumeData) DeepCopy() *TrafficVolumeData {
    if value == nil {
        return nil
    }
    result := &TrafficVolumeData {
        UlMsgs:             value.UlMsgs,
        UlBytes:            value.UlBytes,
        LastUlMsgTime:      value.LastUlMsgTime,
        DlMsgs:             value.DlMsgs,
        DlBytes:            value.DlBytes,
        LastDlMsgTime:      value.LastDlMsgTime,
    }
    if value.UlTotals != nil {
        result.UlTotals = &TotalsState {
            Timestamp: value.UlTotals.Timestamp,
            Msgs:      value.UlTotals.Msgs,
            Bytes:     value.UlTotals.Bytes,       
        }
    }    
    if value.DlTotals != nil {
        result.DlTotals = &TotalsState {
            Timestamp: value.DlTotals.Timestamp,
            Msgs:      value.DlTotals.Msgs,
            Bytes:     value.DlTotals.Bytes,       
        }
    }    
    return result
}
func updateTrafficVolumeData(trafficData *TrafficVolumeData, connection *Connection) *TrafficVolumeData {
    trafficData.UlMsgs   = connection.UlDevice.Msgs
    trafficData.UlBytes  = connection.UlDevice.Bytes
    if connection.UlDevice.Msgs > 0 {
        trafficData.LastUlMsgTime = connection.UlDevice.Timestamp
    }
    trafficData.DlMsgs   = connection.DlDevice.Msgs
    trafficData.DlBytes  = connection.DlDevice.Bytes
    if connection.DlDevice.Msgs > 0 {
        trafficData.LastDlMsgTime = connection.DlDevice.Timestamp
    }
    if connection.UlTotals != nil {
        if trafficData.UlTotals == nil {
            trafficData.UlTotals = &TotalsState{}
        }
        trafficData.UlTotals.Timestamp = connection.UlTotals.Timestamp
        trafficData.UlTotals.Msgs      = connection.UlTotals.Msgs
        trafficData.UlTotals.Bytes     = connection.UlTotals.Bytes
    }
    if connection.DlTotals != nil {
        if trafficData.DlTotals == nil {
            trafficData.DlTotals = &TotalsState{}
        }
        trafficData.DlTotals.Timestamp = connection.DlTotals.Timestamp
        trafficData.DlTotals.Msgs      = connection.DlTotals.Msgs
        trafficData.DlTotals.Bytes     = connection.DlTotals.Bytes
    }
    
    return trafficData
}

// Storage for InitInd data
type InitIndData struct {
    Timestamp         time.Time
    WakeUpCode        string
    RevisionLevel     uint8
    SdCardNotRequired bool
    DisableModemDebug bool
    DisableButton     bool
    DisableServerPing bool
}
func (value *InitIndData) DeepCopy() *InitIndData {
    if value == nil {
        return nil
    }
    result := &InitIndData {
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
func makeInitIndData(newData *InitIndUlMsg, Time time.Time) *InitIndData {
    data := InitIndData {
        Timestamp:         Time,
        WakeUpCode:           WakeUpCodeLookUp[newData.WakeUpCode],
        RevisionLevel:     newData.RevisionLevel,
        SdCardNotRequired: newData.SdCardNotRequired,
        DisableModemDebug: newData.DisableModemDebug,
        DisableButton:     newData.DisableButton,
        DisableServerPing: newData.DisableServerPing,
    }
    return &data
}

//  Storage for intervals data
type IntervalsData struct {
    Timestamp                time.Time
    ReportingIntervalPresent bool
    ReportingInterval        uint32
    HeartbeatPresent         bool
    HeartbeatSeconds         uint32
    HeartbeatSnapToRtc       bool
}
func (value *IntervalsData) DeepCopy() *IntervalsData {
    if value == nil {
        return nil
    }
    result := &IntervalsData {
        Timestamp:                value.Timestamp,
        ReportingIntervalPresent: value.ReportingIntervalPresent,
        ReportingInterval:        value.ReportingInterval,
        HeartbeatPresent:         value.HeartbeatPresent,
        HeartbeatSeconds:         value.HeartbeatSeconds,
        HeartbeatSnapToRtc:       value.HeartbeatSnapToRtc,
    }
    return result
}
func makeIntervalsData(newData *IntervalsGetCnfUlMsg, Time time.Time) *IntervalsData {
    data := IntervalsData {
        Timestamp:                Time,
        ReportingIntervalPresent: true,
        ReportingInterval:        newData.ReportingInterval,
        HeartbeatPresent:         true,
        HeartbeatSeconds:         newData.HeartbeatSeconds,
        HeartbeatSnapToRtc:       newData.HeartbeatSnapToRtc,
    }
    return &data
}
func updateIntervalsData0(oldData *IntervalsData, newData *ReportingIntervalSetCnfUlMsg, Time time.Time) *IntervalsData {
    data := IntervalsData {
        Timestamp:                Time,
        ReportingIntervalPresent: true,
        ReportingInterval:        newData.ReportingInterval,
    }
    if oldData != nil {
        data.HeartbeatPresent = oldData.HeartbeatPresent
        data.HeartbeatSeconds = oldData.HeartbeatSeconds
        data.HeartbeatSnapToRtc = oldData.HeartbeatSnapToRtc
    }
    return &data
}
func updateIntervalsData1(oldData *IntervalsData, newData *HeartbeatSetCnfUlMsg, Time time.Time) *IntervalsData {
    data := IntervalsData {
        Timestamp:            Time,
        HeartbeatPresent:     true,
        HeartbeatSeconds:     newData.HeartbeatSeconds,
        HeartbeatSnapToRtc:   newData.HeartbeatSnapToRtc,
    }
    if oldData != nil {
        data.ReportingIntervalPresent = oldData.ReportingIntervalPresent
        data.ReportingInterval = oldData.ReportingInterval
    }
    return &data
}

//  Storage for mode data
type ModeData struct {
    Timestamp    time.Time
    Mode         string
}
func (value *ModeData) DeepCopy() *ModeData {
    if value == nil {
        return nil
    }
    result := &ModeData {
        Timestamp:    value.Timestamp,
        Mode:         value.Mode,
    }
    return result
}
func makeModeData0(newData *ModeSetCnfUlMsg, Time time.Time) *ModeData {
    data := ModeData {
        Timestamp:    Time,
        Mode:         ModeLookUp[newData.Mode],
    }
    return &data
}
func makeModeData1(newData *ModeGetCnfUlMsg, Time time.Time) *ModeData {
    data := ModeData {
        Timestamp:    Time,
        Mode:         ModeLookUp[newData.Mode],
    }
    return &data
}
func makeModeData2(newData *PollIndUlMsg, Time time.Time) *ModeData {
    data := ModeData {
        Timestamp:    Time,
        Mode:         ModeLookUp[newData.Mode],
    }
    return &data
}

//  Storage for date/time data
type DateTimeData struct {
    Timestamp       time.Time
    UtmTime         time.Time
    TimeSetBy       string
}
func (value *DateTimeData) DeepCopy() *DateTimeData {
    if value == nil {
        return nil
    }
    result := &DateTimeData {
        Timestamp:       value.Timestamp,
        UtmTime:         value.UtmTime,
        TimeSetBy:       value.TimeSetBy,
    }
    return result
}
func makeDateTimeData0(newData *DateTimeIndUlMsg, Time time.Time) *DateTimeData {
    data := DateTimeData {
        Timestamp:       Time,
        UtmTime:         newData.UtmTime,
        TimeSetBy:       TimeSetByLookUp[newData.TimeSetBy],
    }
    return &data
}
func makeDateTimeData1(newData *DateTimeSetCnfUlMsg, Time time.Time) *DateTimeData {
    data := DateTimeData {
        Timestamp:       Time,
        UtmTime:         newData.UtmTime,
        TimeSetBy:       TimeSetByLookUp[newData.TimeSetBy],
    }
    return &data
}
func makeDateTimeData2(newData *DateTimeGetCnfUlMsg, Time time.Time) *DateTimeData {
    data := DateTimeData {
        Timestamp:       Time,
        UtmTime:         newData.UtmTime,
        TimeSetBy:       TimeSetByLookUp[newData.TimeSetBy],
    }
    return &data
}

//  Storage for UtmStatus data
type UtmStatusData struct {
    Timestamp     time.Time
    EnergyLeft    string
    DiskSpaceLeft string
}
func (value *UtmStatusData) DeepCopy() *UtmStatusData {
    if value == nil {
        return nil
    }
    result := &UtmStatusData {
        Timestamp:         value.Timestamp,
        EnergyLeft:        value.EnergyLeft,
        DiskSpaceLeft:     value.DiskSpaceLeft,
    }
    return result
}
func makeUtmStatusData(newData *PollIndUlMsg, Time time.Time) *UtmStatusData {
    data := UtmStatusData {
        Timestamp:         Time,
        EnergyLeft:           EnergyLeftLookUp[newData.EnergyLeft],
        DiskSpaceLeft:     DiskSpaceLeftLookUp[newData.DiskSpaceLeft],
    }
    return &data
}

//  Storage for GNSS data
type GnssData struct {
    Timestamp    time.Time
    Latitude     float32
    Longitude    float32
    Elevation    float32
}
func (value *GnssData) DeepCopy() *GnssData {
    if value == nil {
        return nil
    }
    result := &GnssData {
        Timestamp:     value.Timestamp,
        Latitude:      value.Latitude,
        Longitude:     value.Longitude,
        Elevation:     value.Elevation,
    }
    return result
}
func makeGnssData(newData *MeasurementData, Time time.Time) *GnssData {
    data := GnssData{
        Timestamp:    Time,
        Latitude:     float32(newData.GnssPosition.Latitude) / 1000,
        Longitude:    float32(newData.GnssPosition.Longitude) / 1000,
        Elevation:    float32(newData.GnssPosition.Elevation),
    }
    return &data
}

//  Storage for cell ID data
type CellIdData struct {
    Timestamp    time.Time
    CellId       uint32
}
func (value *CellIdData) DeepCopy() *CellIdData {
    if value == nil {
        return nil
    }
    result := &CellIdData {
        Timestamp:  value.Timestamp,
        CellId:     value.CellId,
    }
    return result
}
func makeCellIdData(newData *MeasurementData, Time time.Time) *CellIdData {
    data := CellIdData{
        Timestamp:    Time,
        CellId:       uint32(newData.CellId),
    }
    return &data
}

//  Storage for temperature data
type TemperatureData struct {
    Timestamp    time.Time
    TemperatureC float32
}
func (value *TemperatureData) DeepCopy() *TemperatureData {
    if value == nil {
        return nil
    }
    result := &TemperatureData {
        Timestamp:     value.Timestamp,
        TemperatureC:  value.TemperatureC,
    }
    return result
}
func makeTemperatureData(newData *MeasurementData, Time time.Time) *TemperatureData {
    data := TemperatureData{
        Timestamp:    Time,
        TemperatureC: float32(newData.Temperature),
    }
    return &data
}

//  Storage for signal strength data
type SignalStrengthData struct {
    RsrpTimestamp     *time.Time
    RsrpDbm           float32
    Mcl               float32
    HighestMcl        float32
    InsideGsmCoverage bool
    RssiPresent       bool
    RssiTimestamp     *time.Time
    RssiDbm           float32
    SnrPresent        bool
    Snr               float32
}
func (value *SignalStrengthData) DeepCopy() *SignalStrengthData {
    if value == nil {
        return nil
    }
    result := &SignalStrengthData{
        RsrpTimestamp:     value.RsrpTimestamp,
        RsrpDbm:           value.RsrpDbm,
        Mcl:               value.Mcl,
        HighestMcl:        value.HighestMcl,
        InsideGsmCoverage: value.InsideGsmCoverage,
        RssiPresent:       value.RssiPresent,
        RssiTimestamp:     value.RssiTimestamp,
        RssiDbm:           value.RssiDbm,
        SnrPresent:        value.SnrPresent,
        Snr:               value.Snr,
    }
    return result
} 
func updateSignalStrengthData(newData *MeasurementData, oldData *SignalStrengthData, Time time.Time) *SignalStrengthData {
    data := oldData
    
    if data == nil {
        data = &SignalStrengthData{}
    }
    
    data.RsrpTimestamp = &newData.TimeMeasured
    data.RsrpDbm       = float32 (newData.Rsrp.Value) / 10
    // TODO: calculate the other answers
     if newData.RssiPresent {
        data.RssiPresent = true
        data.RssiTimestamp = &newData.TimeMeasured
        data.RssiDbm = float32 (newData.Rssi) / 10
        // TODO: calculate the other answers
    }
         
    return data
}

//  Storage for power/state data
type PowerStateData struct {
    Timestamp       time.Time
    ChargeState     string
    BatteryVoltageV float32
    EnergyLeftWh    float32
}

func (value *PowerStateData) DeepCopy() *PowerStateData {
    if value == nil {
        return nil
    }
    result := &PowerStateData{
        Timestamp:       value.Timestamp,
        ChargeState:     value.ChargeState,
        BatteryVoltageV: value.BatteryVoltageV,
        EnergyLeftWh:    value.EnergyLeftWh,
    }
    return result
}
func makePowerStateData(newData *MeasurementData, Time time.Time) *PowerStateData {
    // Change nothing if the new data is nil
    if newData == nil {
        return nil
    }
    // Overwrite if there is no stored data
    data := &PowerStateData{
        Timestamp:       Time,
        ChargeState:     ChargerStateEnumLookUp[newData.PowerState.ChargerState],
        BatteryVoltageV: float32(newData.PowerState.BatteryMv) / 1000,
        EnergyLeftWh:    float32(newData.PowerState.EnergyMwh) / 1000,
    }
    return data
}

//  Storage for traffic data
type TrafficReportData struct {
    Timestamp                 time.Time
    NumDatagramsUl            uint32
    NumBytesUl                uint32
    NumDatagramsDl            uint32
    NumBytesDl                uint32
    NumDatagramsDlBadChecksum uint32
}
func (value *TrafficReportData) DeepCopy() *TrafficReportData {
    if value == nil {
        return nil
    }
    result := &TrafficReportData {
        Timestamp:                  value.Timestamp,
        NumDatagramsUl:             value.NumDatagramsUl,
        NumBytesUl:                 value.NumBytesUl,
        NumDatagramsDl:             value.NumDatagramsDl,
        NumBytesDl:                 value.NumBytesDl,
        NumDatagramsDlBadChecksum:  value.NumDatagramsDlBadChecksum,
    }
    return result
}
func makeTrafficReportData0(newData *TrafficReportIndUlMsg, Time time.Time) *TrafficReportData {
    data := TrafficReportData {
        Timestamp:                  Time,
        NumDatagramsUl:             newData.NumDatagramsUl,
        NumBytesUl:                 newData.NumBytesUl,
        NumDatagramsDl:             newData.NumDatagramsDl,
        NumBytesDl:                 newData.NumBytesDl,
        NumDatagramsDlBadChecksum:  newData.NumDatagramsDlBadChecksum,
    }
    return &data
}
func makeTrafficReportData1(newData *TrafficReportGetCnfUlMsg, Time time.Time) *TrafficReportData {
    data := TrafficReportData {
        Timestamp:                  Time,
        NumDatagramsUl:             newData.NumDatagramsUl,
        NumBytesUl:                 newData.NumBytesUl,
        NumDatagramsDl:             newData.NumDatagramsDl,
        NumBytesDl:                 newData.NumBytesDl,
        NumDatagramsDlBadChecksum:  newData.NumDatagramsDlBadChecksum,
    }
    return &data
}

//  Storage for traffic test mode parameter data
type TrafficTestModeParametersData struct {
    Timestamp           time.Time
    NumUlDatagrams      uint32
    LenUlDatagram       uint32
    NumDlDatagrams      uint32
    LenDlDatagram       uint32
    TimeoutSeconds      uint32
    NoReportsDuringTest bool
}
func (value *TrafficTestModeParametersData) DeepCopy() *TrafficTestModeParametersData {
    if value == nil {
        return nil
    }
    result := &TrafficTestModeParametersData {
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
func makeTrafficTestModeParametersData0(newData *TrafficTestModeParametersSetCnfUlMsg, Time time.Time) *TrafficTestModeParametersData {
    data := TrafficTestModeParametersData {
        Timestamp:           Time,
        NumUlDatagrams:      newData.NumUlDatagrams,
        LenUlDatagram:       newData.LenUlDatagram,
        NumDlDatagrams:      newData.NumDlDatagrams,
        LenDlDatagram:       newData.LenDlDatagram,
        TimeoutSeconds:      newData.TimeoutSeconds,
        NoReportsDuringTest: newData.NoReportsDuringTest,
    }
    return &data
}
func makeTrafficTestModeParametersData1(newData *TrafficTestModeParametersGetCnfUlMsg, Time time.Time) *TrafficTestModeParametersData {
    data := TrafficTestModeParametersData {
        Timestamp:           Time,
        NumUlDatagrams:      newData.NumUlDatagrams,
        LenUlDatagram:       newData.LenUlDatagram,
        NumDlDatagrams:      newData.NumDlDatagrams,
        LenDlDatagram:       newData.LenDlDatagram,
        TimeoutSeconds:      newData.TimeoutSeconds,
        NoReportsDuringTest: newData.NoReportsDuringTest,
    }
    return &data
}

//  Storage for traffic test mode data
type TrafficTestModeReportData struct {
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
func (value *TrafficTestModeReportData) DeepCopy() *TrafficTestModeReportData {
    if value == nil {
        return nil
    }
    result := &TrafficTestModeReportData {
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
func makeTrafficTestModeReportData0(newData *TrafficTestModeReportIndUlMsg, Time time.Time) *TrafficTestModeReportData {
    data := TrafficTestModeReportData {
        Timestamp:                           Time,
        NumTrafficTestDatagramsUl:           newData.NumTrafficTestDatagramsUl,
        NumTrafficTestBytesUl:               newData.NumTrafficTestBytesUl,
        NumTrafficTestDatagramsDl:           newData.NumTrafficTestDatagramsDl,
        NumTrafficTestBytesDl:               newData.NumTrafficTestBytesDl,
        NumTrafficTestDlDatagramsOutOfOrder: newData.NumTrafficTestDlDatagramsOutOfOrder,
        NumTrafficTestDlDatagramsBad:        newData.NumTrafficTestDlDatagramsBad,
        NumTrafficTestDlDatagramsMissed:     newData.NumTrafficTestDlDatagramsMissed,
        TimedOut:                            newData.TimedOut,
    }
    return &data
}
func makeTrafficTestModeReportData1(newData *TrafficTestModeReportGetCnfUlMsg, Time time.Time) *TrafficTestModeReportData {
    data := TrafficTestModeReportData {
        Timestamp:                           Time,
        NumTrafficTestDatagramsUl:           newData.NumTrafficTestDatagramsUl,
        NumTrafficTestBytesUl:               newData.NumTrafficTestBytesUl,
        NumTrafficTestDatagramsDl:           newData.NumTrafficTestDatagramsDl,
        NumTrafficTestBytesDl:               newData.NumTrafficTestBytesDl,
        NumTrafficTestDlDatagramsOutOfOrder: newData.NumTrafficTestDlDatagramsOutOfOrder,
        NumTrafficTestDlDatagramsBad:        newData.NumTrafficTestDlDatagramsBad,
        NumTrafficTestDlDatagramsMissed:     newData.NumTrafficTestDlDatagramsMissed,
        TimedOut:                            newData.TimedOut,
    }
    return &data
}

//  Storage for activity report data
type ActivityReportData struct {
    Timestamp                   time.Time
    TotalTransmitSeconds        float32
    TotalReceiveSeconds         float32
    UpTimeSeconds               uint32
    TxPowerPresent              bool
    TxPowerDbm                  int8
    TxPowerTimestamp            time.Time
    UlMcsPresent                bool
    UlMcs                       uint8
    UlMcsTimestamp              time.Time
    DlMcsPresent                bool
    DlMcs                       uint8
    DlMcsTimestamp              time.Time
}
func (value *ActivityReportData) DeepCopy() *ActivityReportData {
    if value == nil {
        return nil
    }
    result := &ActivityReportData {
        Timestamp:                  value.Timestamp,
        TotalTransmitSeconds:       value.TotalTransmitSeconds,
        TotalReceiveSeconds:        value.TotalReceiveSeconds,
        UpTimeSeconds:              value.UpTimeSeconds,
        TxPowerPresent:             value.TxPowerPresent,
        TxPowerDbm:                 value.TxPowerDbm,
        TxPowerTimestamp:           value.TxPowerTimestamp,
        UlMcsPresent:               value.UlMcsPresent,
        UlMcs:                      value.UlMcs,
        UlMcsTimestamp:             value.UlMcsTimestamp,
        DlMcsPresent:               value.DlMcsPresent,
        DlMcs:                      value.DlMcs,
        DlMcsTimestamp:             value.DlMcsTimestamp,
    }
    return result
}
func updateActivityReportData0(newData *ActivityReportIndUlMsg, oldData *ActivityReportData, Time time.Time) *ActivityReportData {
    data := oldData
    
    if oldData == nil {
        data = &ActivityReportData{}
    }
    
    data.Timestamp            = Time
    data.TotalTransmitSeconds = float32(newData.TotalTransmitMilliseconds) / 1000
    data.TotalReceiveSeconds  = float32(newData.TotalReceiveMilliseconds) / 1000
    data.UpTimeSeconds        = newData.UpTimeSeconds
    
    if newData.TxPowerPresent {
        data.TxPowerPresent = true
        data.TxPowerDbm = newData.TxPowerDbm
        data.TxPowerTimestamp = Time
    }
    if (newData.UlMcsPresent) {
        data.UlMcsPresent = true
        data.UlMcs = newData.UlMcs
        data.UlMcsTimestamp = Time
    }
    if (newData.DlMcsPresent) {
        data.DlMcsPresent = true
        data.DlMcs = newData.DlMcs
        data.DlMcsTimestamp = Time
    }
        
    return data
}
func updateActivityReportData1(newData *ActivityReportGetCnfUlMsg, oldData *ActivityReportData, Time time.Time) *ActivityReportData {
    data := oldData
    
    if oldData == nil {
        data = &ActivityReportData{}
    }
    
    data.Timestamp            = Time
    data.TotalTransmitSeconds = float32(newData.TotalTransmitMilliseconds) / 1000
    data.TotalReceiveSeconds  = float32(newData.TotalReceiveMilliseconds) / 1000
    data.UpTimeSeconds        = newData.UpTimeSeconds
    
    if newData.TxPowerPresent {
        data.TxPowerPresent = true
        data.TxPowerDbm = newData.TxPowerDbm
        data.TxPowerTimestamp = Time
    }
    if (newData.UlMcsPresent) {
        data.UlMcsPresent = true
        data.UlMcs = newData.UlMcs
        data.UlMcsTimestamp = Time
    }
    if (newData.DlMcsPresent) {
        data.DlMcsPresent = true
        data.DlMcs = newData.DlMcs
        data.DlMcsTimestamp = Time
    }
        
    return data
}

/* End Of File */
