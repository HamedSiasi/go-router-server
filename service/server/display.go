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
)

//--------------------------------------------------------------------
// Types used in messages with associated copy functions
//--------------------------------------------------------------------

type FrontPageSummaryData struct {
    TotalUlMsgs        int        `json:"TotalUlMsgs, omitempty"`
    TotalUlBytes       int        `json:"TotalUlBytes, omitempty"`
    LastUlMsgTime      *time.Time `json:"LastUlMsgTime, omitempty"`    
    TotalDlMsgs        int        `json:"TotalDlMsgs, omitempty"`
    TotalDlBytes       int        `json:"TotalDlBytes, omitempty"`
    LastDlMsgTime      *time.Time `json:"LastDlMsgTime, omitempty"`    
}

type FrontPageDeviceData struct {
    Uuid               string     `json:"Uuid, omitempty"`
    DeviceName         string     `json:"DeviceName, omitempty"`
    Interesting        bool       `json:"Interesting, omitempty"`
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
}

type FrontPageData struct {
    SummaryData FrontPageSummaryData
    DeviceData []FrontPageDeviceData
}

func displayFrontPageData () *FrontPageData {
    var data FrontPageData
    var deviceData FrontPageDeviceData
    
    data.SummaryData.TotalUlMsgs = totalUlMsgs
    data.SummaryData.TotalUlBytes = totalUlBytes
    if totalUlBytes > 0 { 
        data.SummaryData.LastUlMsgTime = &lastUlMsgTime
    }    
    data.SummaryData.TotalDlMsgs = totalDlMsgs
    data.SummaryData.TotalDlBytes = totalDlBytes
    if totalDlBytes > 0 { 
        data.SummaryData.LastDlMsgTime = &lastDlMsgTime
    }    
	// Get the latest state for all devices
	// This is very "go" syntax.  It means, create a channel called "get".  Pass the address
	// of the channel into the dataTableChannel (where it is populated).  Then, pass the
	// channel into the variable state.  All of these things are of type []LatestState.
	get := make(chan []DeviceLatestState)
	dataTableChannel <- &get
	allDevicesState := <- get
	
	for _, deviceState := range allDevicesState {
	    deviceData.Uuid             = deviceState.DeviceUuid
        deviceData.DeviceName       = deviceState.State.DeviceName
        if deviceState.State.LatestInterest != nil {
            deviceData.Interesting      = deviceState.State.LatestInterest.IsInteresting
        }
        if deviceState.State.LatestModeData != nil {
            deviceData.Mode             = deviceState.State.LatestModeData.Mode
        }    
        if deviceState.State.LatestTrafficVolumeData != nil {
            deviceData.TotalUlMsgs      = deviceState.State.LatestTrafficVolumeData.TotalUlMsgs
            deviceData.TotalUlBytes     = deviceState.State.LatestTrafficVolumeData.TotalUlBytes
            deviceData.LastUlMsgTime    = &deviceState.State.LatestTrafficVolumeData.LastUlMsgTime
            deviceData.TotalDlMsgs      = deviceState.State.LatestTrafficVolumeData.TotalDlMsgs
            deviceData.TotalDlBytes     = deviceState.State.LatestTrafficVolumeData.TotalDlBytes
            deviceData.LastDlMsgTime    = &deviceState.State.LatestTrafficVolumeData.LastDlMsgTime            
        }
        if deviceState.State.LatestSignalStrengthData != nil {
            deviceData.Rsrp             = deviceState.State.LatestSignalStrengthData.RsrpDbm            
        }
        if deviceState.State.LatestUtmStatusData != nil {
            deviceData.BatteryLevel     = deviceState.State.LatestUtmStatusData.EnergyLeft
            deviceData.DiskSpaceLeft    = deviceState.State.LatestUtmStatusData.DiskSpaceLeft
        }
        if deviceState.State.LatestIntervalsData != nil {
            deviceData.Reporting        = getReportingString (deviceState.State.LatestIntervalsData)
            deviceData.Heartbeat        = getHeartbeatString (deviceState.State.LatestIntervalsData)
        }
        data.DeviceData = append (data.DeviceData, deviceData)
	}

   return &data
}

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

/* End Of File */
