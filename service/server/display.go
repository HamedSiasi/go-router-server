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
	"fmt"
	"github.com/u-blox/utm-server/service/globals"
	"sort"
	"time"
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

type FrontPageSummaryData struct {
	TotalUlMsgs      int        `json:"TotalUlMsgs, omitempty"`
	TotalUlBytes     int        `json:"TotalUlBytes, omitempty"`
	LastUlMsgTime    *time.Time `json:"LastUlMsgTime, omitempty"`
	TotalDlMsgs      int        `json:"TotalDlMsgs, omitempty"`
	TotalDlBytes     int        `json:"TotalDlBytes, omitempty"`
	LastDlMsgTime    *time.Time `json:"LastDlMsgTime, omitempty"`
	DevicesKnown     int        `json:"DevicesKnown, omitempty"`
	DevicesConnected int        `json:"DevicesConnected, omitempty"`
	NumExpectedMsgs  int        `json:"NumExpectedMsgs, omitempty"`
}

type FrontPageDeviceData struct {
	Uuid                string               `json:"Uuid, omitempty"`
	DeviceName          string               `json:"DeviceName, omitempty"`
	Connected           bool                 `json:"Connected, omitempty"`
	UpDuration          string               `json:"UpDuration, omitempty"`
	DeviceTime          *time.Time           `json:"DeviceTime, omitempty"`
	Interesting         *time.Time           `json:"Interesting, omitempty"`
	Mode                string               `json:"Mode, omitempty"`
	TotalUlMsgs         int                  `json:"TotalUlMsgs, omitempty"`
	TotalUlBytes        int                  `json:"TotalUlBytes, omitempty"`
	LastUlMsgTime       *time.Time           `json:"LastUlMsgTime, omitempty"`
	TotalDlMsgs         int                  `json:"TotalDlMsgs, omitempty"`
	TotalDlBytes        int                  `json:"TotalDlBytes, omitempty"`
	LastDlMsgTime       *time.Time           `json:"LastDlMsgTime, omitempty"`
	BatteryLevel        string               `json:"BatteryLevel, omitempty"`
	DiskSpaceLeft       string               `json:"DiskSpaceLeft, omitempty"`
	Reporting           string               `json:"Reporting, omitempty"`
	Heartbeat           string               `json:"Heartbeat, omitempty"`
	NumExpectedMsgs     int                  `json:"NumExpectedMsgs, omitempty"`
	Rsrp                string               `json:"Rsrp, omitempty"`
	RsrpTime            *time.Time           `json:"RsrpTime, omitempty"`
	Rssi                string               `json:"Rssi, omitempty"`
	RssiTime            *time.Time           `json:"RssiTime, omitempty"`
	CellId              string               `json:"CellId, omitempty"`
	CellIdTime          *time.Time           `json:"CellIdTime, omitempty"`
	TxPower             string               `json:"TxPower, omitempty"`
	TxPowerTime         *time.Time           `json:"TxPowerTime, omitempty"`
	CoverageClass       string               `json:"CoverageClass, omitempty"`
	CoverageClassTime   *time.Time           `json:"CoverageClassTime, omitempty"`
	TxDuration          string               `json:"TxTime, omitempty"`
	RxDuration          string               `json:"RxTime, omitempty"`
	TtUlExpected        int                  `json:"TtUlExpected, omitempty"`
	TtUlLength          int                  `json:"TtUlLength, omitempty"`
	TtDlExpected        int                  `json:"TtDlExpected, omitempty"`
	TtDlInterval        int                  `json:"TtDlInterval, omitempty"`
	TtDlLength          int                  `json:"TtDlLength, omitempty"`
	TtTimeout           int                  `json:"TtTimeout, omitempty"`
	TtTimeStarted       *time.Time           `json:"TtTimeStarted, omitempty"`
	TtTimeStopped       *time.Time           `json:"TtTimeStopped, omitempty"`
	TtDlDatagramsTx     int                  `json:"TtDlDatagramsTx, omitempty"`
	TtDlDatagramsRx     int                  `json:"TtDlDatagramsRx, omitempty"`
	TtDlDatagramsMissed int                  `json:"TtDlDatagramsMissed, omitempty"`
	TtDlBytesRx         int                  `json:"TtDlBytes, omitempty"`
	TtUlDatagramsTx     int                  `json:"TtUlDatagramsTx, omitempty"`
	TtUlDatagramsRx     int                  `json:"TtUlDatagramsRx, omitempty"`
	TtUlDatagramsMissed int                  `json:"TtUlDatagramsMissed, omitempty"`
	TtUlBytesRx         int                  `json:"TtUlBytes, omitempty"`
	TtUlState           TrafficTestStateEnum `json:"TtUlState", omitempty"`
	TtDlState           TrafficTestStateEnum `json:"TtDlState", omitempty"`
}

type FrontPageData struct {
	SummaryData FrontPageSummaryData
	DeviceData  []FrontPageDeviceData
}

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

func displayFrontPageData() *FrontPageData {
	var data FrontPageData

	// Get the latest state for all devices
	// This is very "go" syntax.  It means, create a channel called "get".  Pass the address
	// of the channel into the dataTableChannel (where it is populated).  Then, pass the
	// channel into the variable state.  All of these things are of type []LatestState.
	get := make(chan []LatestState)
	dataTableChannel <- &get
	allDevicesState, isOpen := <-get

	if isOpen {
		for _, deviceState := range allDevicesState {

			var deviceData FrontPageDeviceData
			var ttTimeStopped time.Time

			data.SummaryData.DevicesKnown++
			deviceData.Uuid = deviceState.DeviceUuid
			deviceData.DeviceName = deviceState.DeviceName
			deviceData.Connected = deviceState.Connected
			if deviceData.Connected {
				data.SummaryData.DevicesConnected++
			}
			if deviceState.LatestActivityReportData != nil {
				deviceData.UpDuration = makeNiceDurationStringFromSeconds(deviceState.LatestActivityReportData.UpTimeSeconds)
			} // Put the time in if it has been reported recently
			if (deviceState.LatestDateTimeData != nil) && time.Now().UTC().Before(deviceState.LatestDateTimeData.Timestamp.Add(time.Minute)) {
				deviceTime := deviceState.LatestDateTimeData.UtmTime.UTC().Add(time.Now().UTC().Sub(deviceState.LatestDateTimeData.Timestamp))
				deviceData.DeviceTime = &deviceTime
			}
			if deviceState.LatestInterest != nil {
				deviceData.Interesting = &deviceState.LatestInterest.Timestamp
			}
			if deviceState.LatestModeData != nil {
				deviceData.Mode = deviceState.LatestModeData.Mode
			}
			if deviceState.LatestTrafficVolumeData != nil {
				deviceData.TotalUlMsgs = deviceState.LatestTrafficVolumeData.UlMsgs
				deviceData.TotalUlBytes = deviceState.LatestTrafficVolumeData.UlBytes
				deviceData.LastUlMsgTime = &deviceState.LatestTrafficVolumeData.LastUlMsgTime
				deviceData.TotalDlMsgs = deviceState.LatestTrafficVolumeData.DlMsgs
				deviceData.TotalDlBytes = deviceState.LatestTrafficVolumeData.DlBytes
				deviceData.LastDlMsgTime = &deviceState.LatestTrafficVolumeData.LastDlMsgTime
				if deviceState.LatestTrafficVolumeData.UlTotals != nil {
					data.SummaryData.TotalUlMsgs = deviceState.LatestTrafficVolumeData.UlTotals.Msgs
					data.SummaryData.TotalUlBytes = deviceState.LatestTrafficVolumeData.UlTotals.Bytes
					if data.SummaryData.TotalUlMsgs > 0 {
						if data.SummaryData.LastUlMsgTime != nil {
							if deviceState.LatestTrafficVolumeData.UlTotals.Timestamp.After(*data.SummaryData.LastUlMsgTime) {
								data.SummaryData.LastUlMsgTime = &deviceState.LatestTrafficVolumeData.UlTotals.Timestamp
							}
						} else {
							data.SummaryData.LastUlMsgTime = &deviceState.LatestTrafficVolumeData.UlTotals.Timestamp
						}
					}
				}
				if deviceState.LatestTrafficVolumeData.DlTotals != nil {
					data.SummaryData.TotalDlMsgs = deviceState.LatestTrafficVolumeData.DlTotals.Msgs
					data.SummaryData.TotalDlBytes = deviceState.LatestTrafficVolumeData.DlTotals.Bytes
					if data.SummaryData.TotalDlMsgs > 0 {
						if data.SummaryData.LastDlMsgTime != nil {
							if deviceState.LatestTrafficVolumeData.DlTotals.Timestamp.After(*data.SummaryData.LastDlMsgTime) {
								data.SummaryData.LastDlMsgTime = &deviceState.LatestTrafficVolumeData.DlTotals.Timestamp
							}
						} else {
							data.SummaryData.LastDlMsgTime = &deviceState.LatestTrafficVolumeData.DlTotals.Timestamp
						}
					}
				}
			}
			deviceData.Rsrp = "---"
			deviceData.Rssi = "---"
			if deviceState.LatestSignalStrengthData != nil {
				if deviceState.LatestSignalStrengthData.RsrpPresent {
					deviceData.Rsrp = fmt.Sprintf("%.1f", deviceState.LatestSignalStrengthData.RsrpDbm)
					deviceData.RsrpTime = deviceState.LatestSignalStrengthData.RsrpTimestamp
				}
				if deviceState.LatestSignalStrengthData.RssiPresent {
					deviceData.Rssi = fmt.Sprintf("%.1f", deviceState.LatestSignalStrengthData.RssiDbm)
					deviceData.RssiTime = deviceState.LatestSignalStrengthData.RssiTimestamp
				}
			}
			if deviceState.LatestUtmStatusData != nil {
				deviceData.BatteryLevel = deviceState.LatestUtmStatusData.EnergyLeft
				deviceData.DiskSpaceLeft = deviceState.LatestUtmStatusData.DiskSpaceLeft
			}
			if deviceState.LatestIntervalsData != nil {
				deviceData.Reporting = getReportingString(deviceState.LatestIntervalsData)
				deviceData.Heartbeat = getHeartbeatString(deviceState.LatestIntervalsData)
			}
			deviceData.NumExpectedMsgs = 0
			if deviceState.LatestExpectedMsgData != nil && deviceState.LatestExpectedMsgData.ExpectedMsgList != nil {
				deviceData.NumExpectedMsgs = len(*deviceState.LatestExpectedMsgData.ExpectedMsgList)
				data.SummaryData.NumExpectedMsgs += deviceData.NumExpectedMsgs
			}
			deviceData.CellId = "---"
			if deviceState.LatestCellIdData != nil {
				deviceData.CellId = fmt.Sprintf("%d", deviceState.LatestCellIdData.CellId)
				deviceData.CellIdTime = &deviceState.LatestCellIdData.Timestamp
			}
			deviceData.TxPower = "---"
			deviceData.CoverageClass = "---"
			deviceData.TxDuration = "---"
			deviceData.RxDuration = "---"
			if deviceState.LatestActivityReportData != nil {
				deviceData.TxDuration = makeNiceDurationStringFromMilliseconds(deviceState.LatestActivityReportData.TotalTransmitSeconds)
				deviceData.RxDuration = makeNiceDurationStringFromMilliseconds(deviceState.LatestActivityReportData.TotalReceiveSeconds)
				if deviceState.LatestActivityReportData.TxPowerPresent {
					deviceData.TxPower = fmt.Sprintf("%d", deviceState.LatestActivityReportData.TxPowerDbm)
					deviceData.TxPowerTime = &deviceState.LatestActivityReportData.TxPowerTimestamp
				}
				if deviceState.LatestActivityReportData.DlMcsPresent {
					deviceData.CoverageClassTime = &deviceState.LatestActivityReportData.DlMcsTimestamp
					if deviceState.LatestActivityReportData.DlMcs == 5 {
						deviceData.CoverageClass = "0"
					} else {
						if deviceState.LatestActivityReportData.DlMcs == 3 {
							deviceData.CoverageClass = "1"
						} else {
							deviceData.CoverageClass = "2"
						}
					}
				}
			}
			if deviceState.LatestTrafficTestModeParametersData != nil {
				deviceData.TtTimeout = int(deviceState.LatestTrafficTestModeParametersData.TimeoutSeconds)
				deviceData.TtUlExpected = int(deviceState.LatestTrafficTestModeParametersData.NumUlDatagrams)
				deviceData.TtUlLength = int(deviceState.LatestTrafficTestModeParametersData.LenUlDatagram)
				deviceData.TtDlExpected = int(deviceState.LatestTrafficTestModeParametersData.NumDlDatagrams)
				deviceData.TtDlLength = int(deviceState.LatestTrafficTestModeParametersData.LenDlDatagram)
			}
			if deviceState.LatestTrafficTestContext != nil {
				if deviceState.LatestTrafficTestContext.Parameters != nil {
					deviceData.TtDlInterval = int(deviceState.LatestTrafficTestContext.Parameters.DlIntervalSeconds)
				}
				deviceData.TtTimeStarted = &deviceState.LatestTrafficTestContext.TimeStarted
				if (deviceData.TtUlExpected > 0) && (deviceState.LatestTrafficTestContext.UlTimeStopped.After(*deviceData.TtTimeStarted)) {
					ttTimeStopped = deviceState.LatestTrafficTestContext.UlTimeStopped
				}
				if (deviceData.TtDlExpected > 0) && (deviceState.LatestTrafficTestContext.DlTimeStopped.After(*deviceData.TtTimeStarted)) {
					if deviceState.LatestTrafficTestContext.DlTimeStopped.After(ttTimeStopped) {
						ttTimeStopped = deviceState.LatestTrafficTestContext.DlTimeStopped
					}
				}
				if ttTimeStopped.After(*deviceData.TtTimeStarted) {
					deviceData.TtTimeStopped = &ttTimeStopped
				}
				deviceData.TtDlDatagramsTx = int(deviceState.LatestTrafficTestContext.DlDatagrams)
				if deviceState.LatestTrafficTestContext.DeviceTrafficReport != nil {
					deviceData.TtDlDatagramsRx = int(deviceState.LatestTrafficTestContext.DeviceTrafficReport.NumTrafficTestDatagramsDl)
					deviceData.TtDlDatagramsMissed = int(deviceState.LatestTrafficTestContext.DeviceTrafficReport.NumTrafficTestDlDatagramsMissed)
					deviceData.TtDlBytesRx = int(deviceState.LatestTrafficTestContext.DeviceTrafficReport.NumTrafficTestBytesDl)
					deviceData.TtUlDatagramsTx = int(deviceState.LatestTrafficTestContext.DeviceTrafficReport.NumTrafficTestDatagramsUl)
				}
				deviceData.TtUlDatagramsRx = int(deviceState.LatestTrafficTestContext.UlDatagrams)
				deviceData.TtUlDatagramsMissed = int(deviceState.LatestTrafficTestContext.UlDatagramsMissed)
				deviceData.TtUlBytesRx = int(deviceState.LatestTrafficTestContext.UlBytes)
				deviceData.TtUlState = deviceState.LatestTrafficTestContext.UlState
				deviceData.TtDlState = deviceState.LatestTrafficTestContext.DlState
			}

			data.DeviceData = append(data.DeviceData, deviceData)
		}

		// And finally, sort the data by whether the device is connected or not and then by friendly name
		sort.Sort(ByNameAndConnected(data.DeviceData))
		globals.Dbg.PrintfInfo("%s [display] --> Displaying this:\n\n%+v\n\n", globals.LogTag, data)
	}

	return &data
}

// Sort by name and connected
type ByNameAndConnected []FrontPageDeviceData

func BToI(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (u ByNameAndConnected) Len() int      { return len(u) }
func (u ByNameAndConnected) Swap(i, j int) { u[i], u[j] = u[j], u[i] }
func (u ByNameAndConnected) Less(i, j int) bool {
	if BToI(u[i].Connected) == BToI(u[j].Connected) {
		return u[i].DeviceName < u[j].DeviceName
	} else {
		return BToI(u[i].Connected) > BToI(u[j].Connected)
	}
}

// Return the reporting interval as a duration formed as a string
func getReportingString(intervals *IntervalsData) string {
	var reportingString string

	if intervals != nil {
		if intervals.HeartbeatPresent && intervals.ReportingIntervalPresent {
			if intervals.HeartbeatSnapToRtc {
				reportingString = makeNiceTimeStringFromSeconds(3600 * intervals.ReportingInterval)
			} else {
				reportingString = makeNiceTimeStringFromSeconds(intervals.HeartbeatSeconds * intervals.ReportingInterval)
			}
		}
	}

	return reportingString
}

func getHeartbeatString(intervals *IntervalsData) string {
	var heartbeatString string

	if intervals != nil {
		if intervals.HeartbeatPresent {
			heartbeatString = makeNiceTimeStringFromSeconds(intervals.HeartbeatSeconds)
			if intervals.HeartbeatSnapToRtc {
				heartbeatString += "*"
			}
		}
	}

	return heartbeatString
}

func makeNiceTimeStringFromSeconds(seconds uint32) string {
	var theTime time.Time

	theTime = theTime.Add(time.Duration(seconds) * time.Second)

	return theTime.Format("15:04:05")
}

func makeNiceDurationStringFromSeconds(seconds uint32) string {
	var theTime time.Time

	theTime = theTime.Add(time.Duration(seconds) * time.Second)
	if theTime.YearDay() > 1 {
		return theTime.Format("2d 15:04:05")
	} else {
		return theTime.Format("15:04:05")
	}
}

func makeNiceDurationStringFromMilliseconds(milliseconds float32) string {
	var theTime time.Time

	theTime = theTime.Add(time.Duration(milliseconds) * time.Millisecond)
	return theTime.Format("15:04:05.000")
}

/* End Of File */
