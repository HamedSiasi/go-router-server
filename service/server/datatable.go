/* Datatable storage for the UTM server.
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
    "github.com/davecgh/go-spew/spew"
	"github.com/robmeades/utm/service/globals"
)

//--------------------------------------------------------------------
// Types 
//--------------------------------------------------------------------

// NOTE: if you ever add anything here, don't forget to add it to the InitIndUlMsg handling
// and to the copying under DeviceLatestState and AllDevicesLatestState
type LatestState struct {
    Connected                            bool                           `json:"Connected,omitempty"`
    DeviceName                           string                         `json:"DeviceName,omitempty"`
    LastHeardFrom                        time.Time                      `json:"LastHeardFrom,omitempty"`
    LatestInterest                       *Interesting                   `json:"Interesting,omitempty"`
    LatestTrafficVolumeData              *TrafficVolumeData             `json:"LatestTrafficVolumeData,omitempty"`
    LatestInitIndData                    *InitIndData                   `json:"LatestInitIndData,omitempty"`
    LatestIntervalsData                  *IntervalsData                 `json:"LatestIntervalsData,omitempty"`
    LatestModeData                       *ModeData                      `json:"LatestModeData,omitempty"`
    LatestDateTimeData                   *DateTimeData                  `json:"LatestDateTimeData,omitempty"`
    LatestUtmStatusData                  *UtmStatusData                 `json:"LatestUtmStatusData,omitempty"`
    LatestGnssData                       *GnssData                      `json:"LatestGnssData,omitempty"`
    LatestCellIdData                     *CellIdData                    `json:"LatestCellIdData,omitempty"`
    LatestSignalStrengthData             *SignalStrengthData            `json:"LatestSignalStrengthData,omitempty"`
    LatestTemperatureData                *TemperatureData               `json:"LatestTemperatureData,omitempty"`
    LatestPowerStateData                 *PowerStateData                `json:"LatestPowerStateData,omitempty"`
    LatestTrafficReportData              *TrafficReportData             `json:"LatestTrafficReportData,omitempty"`
    LatestTrafficTestModeParametersData  *TrafficTestModeParametersData `json:"LatestTrafficTestModeParametersData,omitempty"`
    LatestTrafficTestModeReportData      *TrafficTestModeReportData     `json:"LatestTrafficTestModeReportData,omitempty"`
    LatestActivityReportData             *ActivityReportData            `json:"LatestActivityReportData,omitempty"`
}

// Structure to allow the latest state for a particular
// device to be retrieved over a channel
type DeviceLatestStateChannel struct {
	DeviceUuid   string
	State        chan LatestState
}

// Structure to describe the state of a device
type DeviceLatestState struct {
	DeviceUuid   string
	State        LatestState
}

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

// To update the latest values send a MessageContainer into this channel
// containing the received message; a copy their contents will be stored
// in a displayable form

// To get the latest state for a given UUID, send a '*chan DeviceLatestStateChannel'
// into this channel containing the device UUID and a pointer to a channel
// down which to send the LatestState struct; a copy of all quantities will
// be copied into the struct and then the channel will be closed.

// To get the latest state of all devices, send a '*chan []DevicesLatestState'
// into this channel and a copy of all quantities for all UUIDs will be
// copied into the struct and then the channel will be closed.

// To terminate execution simply close the channel

var dataTableChannel chan<- interface{}

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// Set the interesting flag to true and now
func (i * Interesting) Set () {
    if i != nil {
        i.IsInteresting = true
        i.Timestamp = time.Now().UTC()
    }    
}

// Unset the interesting flag to true and now
func (i * Interesting) UnSet () {
    if i != nil {
        i.IsInteresting = false
    }    
}

// Run the whole shebang
func operateDataTable() {
    channel := make(chan interface{})
    dataTableChannel = channel
	deviceLatestStateList := make(map[string]*LatestState)    	
    checkConnections := time.NewTicker (time.Second * 10)
    checkInteresting := time.NewTicker (time.Second * 30)
    
    globals.Dbg.PrintfTrace("%s [datatable] --> command channel created and now being serviced.\n", globals.LogTag)
    
    // Deal with the check connections timer
    // If we know the reporting interval we can say a device
    // is not connected if we've heard nothing for more than the
    // reporting interval
    go func() {
        for _ = range checkConnections.C {
            for uuid, state := range deviceLatestStateList {
                if state.LatestIntervalsData != nil &&
                   state.LatestIntervalsData.ReportingInterval > 0 {
                    if state.LatestIntervalsData.HeartbeatSnapToRtc {
                        if time.Now().After (state.LastHeardFrom.Add(time.Hour + time.Minute * 10)) {
                            state.Connected = false
                            globals.Dbg.PrintfTrace("%s [datatable] --> device %s is no longer connected (last heard from @ %s).\n", globals.LogTag, uuid, state.LastHeardFrom.String())
                        }
                    } else {
                        if time.Now().After (state.LastHeardFrom.Add(time.Duration(state.LatestIntervalsData.HeartbeatSeconds * (state.LatestIntervalsData.ReportingInterval + 2)) * time.Second)) {                   
                            state.Connected = false
                            globals.Dbg.PrintfTrace("%s [datatable] --> device %s is no longer connected  (last heard from @ %s).\n", globals.LogTag, uuid, state.LastHeardFrom.String())
                        }   
                    }                
                }    
            }     
        }
    }()
    
    // Deal with the check interesting timer
    // If a device has been interesting for more than
    // a minute then slap it down
    go func() {
        for _ = range checkInteresting.C {
            for uuid, state := range deviceLatestStateList {
                if state.LatestInterest != nil && state.LatestInterest.IsInteresting {
                    if time.Now().After (state.LatestInterest.Timestamp.Add(time.Minute)) {
                        state.LatestInterest.UnSet()
                        globals.Dbg.PrintfTrace("%s [datatable] --> device %s has become less interesting (last interesting @ %s).\n", globals.LogTag, uuid, state.LatestInterest.Timestamp.String())
                    }
                }    
            }     
        }
    }()
    
    // Deal with commands on the channel
    go func() {
        for cmd := range channel {
            switch value := cmd.(type) {
            	
	            // Handle connection indications
	            case *Connection:
	            	state := deviceLatestStateList[value.DeviceUuid]
            		if state != nil {
    	                state.DeviceName = value.DeviceName;
    	                state.Connected = true;
                        state.LastHeardFrom = value.Timestamp
    	                if state.LatestTrafficVolumeData == nil {
    	                    state.LatestTrafficVolumeData = &TrafficVolumeData {}
    	                }	                
    	                state.LatestTrafficVolumeData = updateTrafficVolumeData (state.LatestTrafficVolumeData, value)
    					globals.Dbg.PrintfTrace("%s [datatable] --> connection state for device %s updated.\n", globals.LogTag, value.DeviceUuid)
    	            }
			                	
	            // Handle message containers holding somethings of interest
	            case *MessageContainer:
	            
	            	// Make sure there's an entry for this device
	            	state := deviceLatestStateList[value.DeviceUuid]
            		if state == nil {
            			state = &LatestState{};
            			deviceLatestStateList[value.DeviceUuid] = state;
                        state.LastHeardFrom = time.Now().UTC()
    	                if state.LatestInterest == nil {
    	                    state.LatestInterest = &Interesting {}
    	                }    
                        state.LatestInterest.Set()
                        globals.Dbg.PrintfTrace("%s [datatable] --> heard from a new device, UUID %s.\n", globals.LogTag, value.DeviceUuid)
                        // TODO: find a way to send a get intervals request for a new device, but not from here
            		}
            		
                    globals.Dbg.PrintfTrace("%s [datatable] --> storing data for UUID %s...\n", globals.LogTag, value.DeviceUuid)
            		
	            	switch utmMsg := value.Message.(type) {
			            case *InitIndUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestInitIndData = makeInitIndData(data, value.Timestamp)
			                    // Device must have (re)booted so clear what we know
			                    // TODO: I _think_ the values that were here all get garbage collected
			                    // as there seems to be no way to do an explicit free.  But I'd really
			                    // like to check.
                                state.LastHeardFrom = time.Now().UTC()
            	                if state.LatestInterest == nil {
            	                    state.LatestInterest = &Interesting {}
            	                }    
    	                        state.LatestInterest.Set()
			                    state.LatestTrafficVolumeData = nil
                                state.LatestIntervalsData = nil
                                state.LatestModeData = nil
                                state.LatestDateTimeData = nil
                                state.LatestUtmStatusData = nil
                                state.LatestGnssData = nil
                                state.LatestCellIdData = nil
                                state.LatestSignalStrengthData = nil
                                state.LatestTemperatureData = nil
                                state.LatestPowerStateData = nil
                                state.LatestTrafficReportData = nil
                                state.LatestTrafficTestModeParametersData = nil
                                state.LatestTrafficTestModeReportData = nil
                                state.LatestActivityReportData = nil
			                }
			
			            case *IntervalsGetCnfUlMsg:
	                        state.LatestInterest.Set()
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestIntervalsData = makeIntervalsData0(data, value.Timestamp)
			                }
			
			            case *ReportingIntervalSetCnfUlMsg:
	                        state.LatestInterest.Set()
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestIntervalsData = makeIntervalsData1(data, value.Timestamp)
			                }
			
			            case *HeartbeatSetCnfUlMsg:
	                        state.LatestInterest.Set()
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestIntervalsData = makeIntervalsData2(data, value.Timestamp)
			                }
			
			            case *DateTimeIndUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestDateTimeData = makeDateTimeData0(data, value.Timestamp)
			                }
			
			            case *DateTimeSetCnfUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestDateTimeData = makeDateTimeData1(data, value.Timestamp)
			                }
			
			            case *DateTimeGetCnfUlMsg:
	                        state.LatestInterest.Set()
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestDateTimeData = makeDateTimeData2(data, value.Timestamp)
			                }
			
			            case *ModeSetCnfUlMsg:
	                        state.LatestInterest.Set()
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestModeData = makeModeData0(data, value.Timestamp)
			                }
			
			            case *ModeGetCnfUlMsg:
	                        state.LatestInterest.Set()
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestModeData = makeModeData1(data, value.Timestamp)
			                }
			                
			            case *PollIndUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestUtmStatusData = makeUtmStatusData(data, value.Timestamp)
			                    state.LatestModeData = makeModeData2(data, value.Timestamp)
			                }
			
			            case *MeasurementsIndUlMsg:
			                data := utmMsg.Measurements.DeepCopy()
			                if data != nil {
			                    state.LatestSignalStrengthData = makeSignalStrengthData(data, value.Timestamp)
			                    if (data.GnssPositionPresent) {
			                    	state.LatestGnssData = makeGnssData(data, value.Timestamp)
			                    }
			                    if (data.CellIdPresent) {
			                    	state.LatestCellIdData = makeCellIdData(data, value.Timestamp)
			                    }
			                    if (data.TemperaturePresent) {
			                    	state.LatestTemperatureData = makeTemperatureData(data, value.Timestamp)
			                    }
			                    if (data.PowerStatePresent) {
			                    	state.LatestPowerStateData = makePowerStateData(data, value.Timestamp)
			                    }
			                }
			
			            case *MeasurementsGetCnfUlMsg:
			                data := utmMsg.Measurements.DeepCopy()
			                if data != nil {
			                    state.LatestSignalStrengthData = makeSignalStrengthData(data, value.Timestamp)
			                    if (data.GnssPositionPresent) {
			                    	state.LatestGnssData = makeGnssData(data, value.Timestamp)
			                    }
			                    if (data.CellIdPresent) {
			                    	state.LatestCellIdData = makeCellIdData(data, value.Timestamp)
			                    }
			                    if (data.TemperaturePresent) {
			                    	state.LatestTemperatureData = makeTemperatureData(data, value.Timestamp)
			                    }
			                    if (data.PowerStatePresent) {
			                    	state.LatestPowerStateData = makePowerStateData(data, value.Timestamp)
			                    }
			                }
			
			            case *TrafficReportIndUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestTrafficReportData = makeTrafficReportData0(data, value.Timestamp)
			                }
			
			            case *TrafficReportGetCnfUlMsg:
	                        state.LatestInterest.Set()
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestTrafficReportData = makeTrafficReportData1(data, value.Timestamp)
			                }
			
			            case *TrafficTestModeParametersSetCnfUlMsg:
	                        state.LatestInterest.Set()
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestTrafficTestModeParametersData = makeTrafficTestModeParametersData0(data, value.Timestamp)
			                }
			
			            case *TrafficTestModeParametersGetCnfUlMsg:
	                        state.LatestInterest.Set()
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestTrafficTestModeParametersData = makeTrafficTestModeParametersData1(data, value.Timestamp)
			                }
			
			            case *TrafficTestModeReportIndUlMsg:
	                        state.LatestInterest.Set()
			                data := utmMsg.DeepCopy()
			                if data != nil {
        			            // TODO check for a pass/fail result 
			                    state.LatestTrafficTestModeReportData = makeTrafficTestModeReportData0(data, value.Timestamp)
			                }
			
			            case *TrafficTestModeReportGetCnfUlMsg:
	                        state.LatestInterest.Set()
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestTrafficTestModeReportData = makeTrafficTestModeReportData1(data, value.Timestamp)
			                }
			
			            case *ActivityReportIndUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestActivityReportData = makeActivityReportData0(data, value.Timestamp)
			                }
			
			            case *ActivityReportGetCnfUlMsg:
	                        state.LatestInterest.Set()
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestActivityReportData = makeActivityReportData1(data, value.Timestamp)
			                }
	            	}
					globals.Dbg.PrintfTrace("%s [datatable] --> storage completed.\n", globals.LogTag)
	            	
	            // Return the latest state for a given UUID 
	            case *DeviceLatestStateChannel:
	            	// Retrieve the device state
   	            	state := deviceLatestStateList[value.DeviceUuid]
   	            	
   	            	if state != nil {
		                // Duplicate the memory pointed to into a new LatestState struct,
		                // post it and close the channel
		                globals.Dbg.PrintfTrace("%s [datatable] --> fetching latest state for UUID %s.\n", globals.LogTag, value.DeviceUuid)
		                latest := LatestState{}
		                latest.Connected = state.Connected
		                latest.LastHeardFrom = state.LastHeardFrom
		                latest.LatestInterest = state.LatestInterest.DeepCopy()
                        state.LatestInterest.UnSet() // Reseting interestingness after answering a specific query
	                    latest.LatestTrafficVolumeData = state.LatestTrafficVolumeData.DeepCopy()
		                latest.LatestInitIndData = state.LatestInitIndData.DeepCopy()
		                latest.LatestIntervalsData = state.LatestIntervalsData.DeepCopy()
		                latest.LatestModeData = state.LatestModeData.DeepCopy()
		                latest.LatestDateTimeData = state.LatestDateTimeData.DeepCopy()
		                latest.LatestUtmStatusData = state.LatestUtmStatusData.DeepCopy()
		                latest.LatestGnssData = state.LatestGnssData.DeepCopy()
		                latest.LatestCellIdData = state.LatestCellIdData.DeepCopy()
		                latest.LatestSignalStrengthData = state.LatestSignalStrengthData.DeepCopy()
		                latest.LatestTemperatureData = state.LatestTemperatureData.DeepCopy()
		                latest.LatestPowerStateData = state.LatestPowerStateData.DeepCopy()
		                latest.LatestTrafficReportData = state.LatestTrafficReportData.DeepCopy()
		                latest.LatestTrafficTestModeParametersData = state.LatestTrafficTestModeParametersData.DeepCopy()
		                latest.LatestTrafficTestModeReportData = state.LatestTrafficTestModeReportData.DeepCopy()
		                latest.LatestActivityReportData = state.LatestActivityReportData.DeepCopy()
		                value.State <- latest
		                close(value.State)
		                globals.Dbg.PrintfTrace("%s [datatable] --> provided latest state and closed channel.\n", globals.LogTag)
		            } else {
		                globals.Dbg.PrintfTrace("%s [datatable] --> asked for latest state for unknown UUID %s.\n", globals.LogTag, value.DeviceUuid)
		            }
	            // Return the latest state for all UUIDs 
	            case *chan []DeviceLatestState:
	            
   	            	var allStates []DeviceLatestState
	                for uuid, state := range deviceLatestStateList {
		                // Duplicate the memory pointed to into a new LatestState struct,
		                // post it and close the channel
		                latest := DeviceLatestState{}
		                latest.State.Connected = state.Connected
		                latest.State.LastHeardFrom = state.LastHeardFrom
		                latest.State.LatestInterest = state.LatestInterest.DeepCopy()
	                    latest.State.LatestTrafficVolumeData = state.LatestTrafficVolumeData.DeepCopy()
		                latest.State.LatestInitIndData = state.LatestInitIndData.DeepCopy()
		                latest.State.LatestIntervalsData = state.LatestIntervalsData.DeepCopy()
		                latest.State.LatestModeData = state.LatestModeData.DeepCopy()
		                latest.State.LatestDateTimeData = state.LatestDateTimeData.DeepCopy()
		                latest.State.LatestUtmStatusData = state.LatestUtmStatusData.DeepCopy()
		                latest.State.LatestGnssData = state.LatestGnssData.DeepCopy()
		                latest.State.LatestCellIdData = state.LatestCellIdData.DeepCopy()
		                latest.State.LatestSignalStrengthData = state.LatestSignalStrengthData.DeepCopy()
		                latest.State.LatestTemperatureData = state.LatestTemperatureData.DeepCopy()
		                latest.State.LatestPowerStateData = state.LatestPowerStateData.DeepCopy()
		                latest.State.LatestTrafficReportData = state.LatestTrafficReportData.DeepCopy()
		                latest.State.LatestTrafficTestModeParametersData = state.LatestTrafficTestModeParametersData.DeepCopy()
		                latest.State.LatestTrafficTestModeReportData = state.LatestTrafficTestModeReportData.DeepCopy()
		                latest.State.LatestActivityReportData = state.LatestActivityReportData.DeepCopy()
		                latest.DeviceUuid = uuid
		                allStates = append(allStates, latest)
    		        }    
	                *value <- allStates
	                close(*value)
	                globals.Dbg.PrintfTrace("%s [datatable] --> provided latest state for all devices and closed channel.\n", globals.LogTag)
		                
	            default:
	                globals.Dbg.PrintfTrace("%s [datatable] --> unrecognised command into channel, ignoring:\n\n%s\n", globals.LogTag, spew.Sdump(cmd))
            }
        }

        globals.Dbg.PrintfTrace("%s [datatable] --> command channel closed, stopping.\n", globals.LogTag)
        checkConnections.Stop();
    }()
}

func init() {
    operateDataTable()
}

/* End Of File */
