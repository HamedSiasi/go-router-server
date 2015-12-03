/* Functions that create displays items values for the UTM server.
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
    "sort"
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

type FrontPageSummaryData struct {
    TotalUlMsgs        int        `json:"TotalUlMsgs, omitempty"`
    TotalUlBytes       int        `json:"TotalUlBytes, omitempty"`
    LastUlMsgTime      *time.Time `json:"LastUlMsgTime, omitempty"`    
    TotalDlMsgs        int        `json:"TotalDlMsgs, omitempty"`
    TotalDlBytes       int        `json:"TotalDlBytes, omitempty"`
    LastDlMsgTime      *time.Time `json:"LastDlMsgTime, omitempty"`
    DevicesKnown       int        `json:"DevicesKnown, omitempty"`
    DevicesConnected   int        `json:"DevicesConnected, omitempty"`
    NumExpectedMsgs    int        `json:"NumExpectedMsgs, omitempty"`
}

type FrontPageDeviceData struct {
    Uuid               string     `json:"Uuid, omitempty"`
    DeviceName         string     `json:"DeviceName, omitempty"`
    Connected          bool       `json:"Connected, omitempty"`
    UpDuration         string     `json:"UpDuration, omitempty"`
    Interesting        *time.Time `json:"Interesting, omitempty"`
    Mode               string     `json:"Mode, omitempty"`
    TotalUlMsgs        int        `json:"TotalUlMsgs, omitempty"`
    TotalUlBytes       int        `json:"TotalUlBytes, omitempty"`
    LastUlMsgTime      *time.Time `json:"LastUlMsgTime, omitempty"`
    TotalDlMsgs        int        `json:"TotalDlMsgs, omitempty"`
    TotalDlBytes       int        `json:"TotalDlBytes, omitempty"`
    LastDlMsgTime      *time.Time `json:"LastDlMsgTime, omitempty"`
    Rsrp               float32    `json:"Rsrp, omitempty"`
    BatteryLevel       string     `json:"BatteryLevel, omitempty"`
    DiskSpaceLeft      string     `json:"DiskSpaceLeft, omitempty"`
    Reporting          string     `json:"Reporting, omitempty"`
    Heartbeat          string     `json:"Heartbeat, omitempty"`
    NumExpectedMsgs    int        `json:"NumExpectedMsgs, omitempty"`
}

type FrontPageData struct {
    SummaryData FrontPageSummaryData
    DeviceData []FrontPageDeviceData
}

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

func displayFrontPageData () *FrontPageData {
    var data FrontPageData
    var deviceData FrontPageDeviceData
    
    // Get the latest state for all devices
    // This is very "go" syntax.  It means, create a channel called "get".  Pass the address
    // of the channel into the dataTableChannel (where it is populated).  Then, pass the
    // channel into the variable state.  All of these things are of type []LatestState.
    get := make(chan []LatestState)
    dataTableChannel <- &get
    allDevicesState := <- get
    
    for _, deviceState := range allDevicesState {
        
        data.SummaryData.DevicesKnown++;
        deviceData.Uuid             = deviceState.DeviceUuid
        deviceData.DeviceName       = deviceState.DeviceName
        deviceData.Connected        = deviceState.Connected
        if (deviceData.Connected) {
            data.SummaryData.DevicesConnected++
        }    
        if deviceState.LatestActivityReportData != nil {
            deviceData.UpDuration       = makeNiceDurationStringFromSeconds(deviceState.LatestActivityReportData.UpTimeSeconds)
        }    
        if deviceState.LatestInterest != nil {
            deviceData.Interesting = &deviceState.LatestInterest.Timestamp
        }
        if deviceState.LatestModeData != nil {
            deviceData.Mode             = deviceState.LatestModeData.Mode
        }    
        if deviceState.LatestTrafficVolumeData != nil {
            deviceData.TotalUlMsgs      = deviceState.LatestTrafficVolumeData.UlMsgs
            deviceData.TotalUlBytes     = deviceState.LatestTrafficVolumeData.UlBytes
            deviceData.LastUlMsgTime = &deviceState.LatestTrafficVolumeData.LastUlMsgTime
            deviceData.TotalDlMsgs      = deviceState.LatestTrafficVolumeData.DlMsgs
            deviceData.TotalDlBytes     = deviceState.LatestTrafficVolumeData.DlBytes
            deviceData.LastDlMsgTime = &deviceState.LatestTrafficVolumeData.LastDlMsgTime
            if deviceState.LatestTrafficVolumeData.UlTotals != nil {
                data.SummaryData.TotalUlMsgs = deviceState.LatestTrafficVolumeData.UlTotals.Msgs
                data.SummaryData.TotalUlBytes = deviceState.LatestTrafficVolumeData.UlTotals.Bytes
                if data.SummaryData.TotalUlMsgs > 0 { 
                    data.SummaryData.LastUlMsgTime = &deviceState.LatestTrafficVolumeData.UlTotals.Timestamp
                }    
            }
            if deviceState.LatestTrafficVolumeData.DlTotals != nil {
                data.SummaryData.TotalDlMsgs = deviceState.LatestTrafficVolumeData.DlTotals.Msgs
                data.SummaryData.TotalDlBytes = deviceState.LatestTrafficVolumeData.DlTotals.Bytes
                if data.SummaryData.TotalDlMsgs > 0 { 
                    data.SummaryData.LastDlMsgTime = &deviceState.LatestTrafficVolumeData.DlTotals.Timestamp
                }    
            }    
        }
        if deviceState.LatestSignalStrengthData != nil {
            deviceData.Rsrp             = deviceState.LatestSignalStrengthData.RsrpDbm            
        }
        if deviceState.LatestUtmStatusData != nil {
            deviceData.BatteryLevel     = deviceState.LatestUtmStatusData.EnergyLeft
            deviceData.DiskSpaceLeft    = deviceState.LatestUtmStatusData.DiskSpaceLeft
        }
        if deviceState.LatestIntervalsData != nil {
            deviceData.Reporting        = getReportingString (deviceState.LatestIntervalsData)
            deviceData.Heartbeat        = getHeartbeatString (deviceState.LatestIntervalsData)
        }
        deviceData.NumExpectedMsgs = 0
        if deviceState.LatestExpectedMsgData != nil && deviceState.LatestExpectedMsgData.ExpectedMsgList != nil {
            deviceData.NumExpectedMsgs = len (*deviceState.LatestExpectedMsgData.ExpectedMsgList)
            data.SummaryData.NumExpectedMsgs += deviceData.NumExpectedMsgs
        }
        data.DeviceData = append (data.DeviceData, deviceData)
    }

    // And finally, sort the data by friendly name
    sort.Sort(ByName(data.DeviceData))
    
    return &data
}

// Sort by name
type ByName []FrontPageDeviceData

func (u ByName) Len() int           { return len(u) }
func (u ByName) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }
func (u ByName) Less(i, j int) bool { return u[i].DeviceName < u[j].DeviceName }

// Return the reporting interval as a duration formed as a string
func getReportingString (intervals *IntervalsData) string {
    var reportingString string
    
    if intervals != nil {
        if intervals.HeartbeatPresent && intervals.ReportingIntervalPresent {
            if intervals.HeartbeatSnapToRtc {
                reportingString = makeNiceTimeStringFromSeconds (3600 * intervals.ReportingInterval)
            } else {
                reportingString = makeNiceTimeStringFromSeconds (intervals.HeartbeatSeconds * intervals.ReportingInterval)
            }
        }
    }
        
    return reportingString
}

func getHeartbeatString (intervals *IntervalsData) string {
    var heartbeatString string
    
    if intervals != nil {
        if intervals.HeartbeatPresent {
            heartbeatString = makeNiceTimeStringFromSeconds (intervals.HeartbeatSeconds)
            if intervals.HeartbeatSnapToRtc {
                heartbeatString += "*"
            }
        }
    }    
    
    return heartbeatString
}

func makeNiceTimeStringFromSeconds (seconds uint32) string {
    var theTime time.Time

    theTime = theTime.Add (time.Duration (seconds) * time.Second)
    
    return theTime.Format ("15:04:05")
}

func makeNiceDurationStringFromSeconds (seconds uint32) string {
    var theTime time.Time

    theTime = theTime.Add (time.Duration (seconds) * time.Second)
    if theTime.YearDay() > 1 {
        return theTime.Format ("2d 15:04:05")
    } else {
        return theTime.Format ("15:04:05")
    }
}

/* End Of File */
