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
)

//--------------------------------------------------------------------
// Types 
//--------------------------------------------------------------------

// NOTE: if you ever add anything here, don't forget to add it to the InitIndUlMsg handling
// and to the copying under DeviceLatestState
type LatestState struct {
    Connected                               bool                              `json:"Connected,omitempty"`
    DeviceName                              string                            `json:"DeviceName,omitempty"`
    LastHeardFrom                           time.Time                         `json:"LastHeardFrom,omitempty"`
    LatestInitIndDisplay                    *InitIndDisplay                   `json:"LatestInitIndDisplay,omitempty"`
    LatestIntervalsDisplay                  *IntervalsDisplay                 `json:"LatestIntervalsDisplay,omitempty"`
    LatestModeDisplay                       *ModeDisplay                      `json:"LatestModeDisplay,omitempty"`
    LatestDateTimeDisplay                   *DateTimeDisplay                  `json:"LatestDateTimeDisplay,omitempty"`
    LatestUtmStatusDisplay                  *UtmStatusDisplay                 `json:"LatestUtmStatusDisplay,omitempty"`
    LatestGnssDisplay                       *GnssDisplay                      `json:"LatestGnssDisplay,omitempty"`
    LatestCellIdDisplay                     *CellIdDisplay                    `json:"LatestCellIdDisplay,omitempty"`
    LatestSignalStrengthDisplay             *SignalStrengthDisplay            `json:"LatestSignalStrengthDisplay,omitempty"`
    LatestTemperatureDisplay                *TemperatureDisplay               `json:"LatestTemperatureDisplay,omitempty"`
    LatestPowerStateDisplay                 *PowerStateDisplay                `json:"LatestPowerStateDisplay,omitempty"`
    LatestTrafficReportDisplay              *TrafficReportDisplay             `json:"LatestTrafficReportDisplay,omitempty"`
    LatestTrafficTestModeParametersDisplay  *TrafficTestModeParametersDisplay `json:"LatestTrafficTestModeParametersDisplay,omitempty"`
    LatestTrafficTestModeReportDisplay      *TrafficTestModeReportDisplay     `json:"LatestTrafficTestModeReportDisplay,omitempty"`
    LatestActivityReportDisplay             *ActivityReportDisplay            `json:"LatestActivityReportDisplay,omitempty"`
    LatestDisplayRow                        *DisplayRow                       `json:"LatestDisplayRow,omitempty"`
}

type DeviceLatestState struct {
	DeviceUuid   string
	Latest       chan LatestState
}

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

// To update the latest values send a MessageContainer into this channel
// containing the received message; a copy their contents will be stored
// in a displayable form

// To get the latest state send a '*chan DeviceLatestState' into this channel
// containing the device UUID and a pointer to a channel down which to send
// the LatestState struct; a copy of all quantities will be copied into the
// struct and then the channel will be closed.

// To terminate execution simply close the channel

var dataTableCmds chan<- interface{}

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

func operateDataTable() {
    cmds := make(chan interface{})
    dataTableCmds = cmds
	deviceLatestStateList := make(map[string]*LatestState)    	
    checkConnections := time.NewTicker (time.Second * 10)
    
    Dbg.PrintfTrace("%s --> datatable command channel created and now being serviced.\n", logTag)
    
    // Deal with the check connections timer
    // If we know the reporting interval we can say a device
    // is not connected if we've heard nothing for more than the
    // reporting interval
    go func() {
        for _ = range checkConnections.C {
            for uuid, state := range deviceLatestStateList {
                if state.LatestIntervalsDisplay != nil &&
                   state.LatestIntervalsDisplay.ReportingInterval > 0 {
                    if state.LatestIntervalsDisplay.HeartbeatSnapToRtc {
                        if time.Now().After (state.LastHeardFrom.Add(time.Hour + time.Minute * 10)) {
                            state.Connected = false;
                            Dbg.PrintfTrace("%s --> device %s no longer connected (last heard from @ %s).\n", logTag, uuid, state.LastHeardFrom.String())
                        }
                    } else {
                        if time.Now().After (state.LastHeardFrom.Add(time.Duration(state.LatestIntervalsDisplay.HeartbeatSeconds * (state.LatestIntervalsDisplay.ReportingInterval + 2)) * time.Second)) {                   
                            state.Connected = false;
                            Dbg.PrintfTrace("%s --> device %s no longer connected  (last heard from @ %s).\n", logTag, uuid, state.LastHeardFrom.String())
                        }   
                    }                
                }    
            }     
        }
    }()
    
    // Deal with messages on the cmds channel
    go func() {
        for cmd := range cmds {
            switch value := cmd.(type) {
            	
	            // Handle connection indications
	            case *Connection:
	            	state := deviceLatestStateList[value.DeviceUuid]
            		if state != nil {
    	                state.LastHeardFrom = value.LastHeardFrom;
    	                state.DeviceName = value.DeviceName;
    	                state.Connected = true;
    	            }    
			                	
	            // Handle message containers holding somethings of interest
	            case *MessageContainer:
	            
	            	// Make sure there's an entry for this device
	            	state := deviceLatestStateList[value.DeviceUuid]
            		if state == nil {
            			state = &LatestState{};
            			deviceLatestStateList[value.DeviceUuid] = state;
                        state.LastHeardFrom = time.Now().Local()
                        Dbg.PrintfTrace("%s --> datatable has heard from a new device, UUID %s.\n", logTag, value.DeviceUuid)
                        // TODO: find a way to send a get intervals request for a new device, but not from here
            		}
            		
                    Dbg.PrintfTrace("%s --> datatable received a message from UUID %s.\n", logTag, value.DeviceUuid)
            		
	            	switch utmMsg := value.Message.(type) {
			            case *InitIndUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestInitIndDisplay = makeInitIndDisplay(data, value.Timestamp)
			                    // Device must have (re)booted so clear what we know
			                    // TODO: I _think_ the values that were here all get garbage collected
			                    // as there seems to be no way to do an explicit free.  But I'd really
			                    // like to check.
                                state.LastHeardFrom = time.Now().Local()
                                state.LatestIntervalsDisplay = nil
                                state.LatestModeDisplay = nil
                                state.LatestDateTimeDisplay = nil
                                state.LatestUtmStatusDisplay = nil
                                state.LatestGnssDisplay = nil
                                state.LatestCellIdDisplay = nil
                                state.LatestSignalStrengthDisplay = nil
                                state.LatestTemperatureDisplay = nil
                                state.LatestPowerStateDisplay = nil
                                state.LatestTrafficReportDisplay = nil
                                state.LatestTrafficTestModeParametersDisplay = nil
                                state.LatestTrafficTestModeReportDisplay = nil
                                state.LatestActivityReportDisplay = nil
                                state.LatestDisplayRow = nil
			                }
			
			            case *IntervalsGetCnfUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestIntervalsDisplay = makeIntervalsDisplay0(data, value.Timestamp)
			                }
			
			            case *ReportingIntervalSetCnfUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestIntervalsDisplay = makeIntervalsDisplay1(data, value.Timestamp)
			                }
			
			            case *HeartbeatSetCnfUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestIntervalsDisplay = makeIntervalsDisplay2(data, value.Timestamp)
			                }
			
			            case *DateTimeIndUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestDateTimeDisplay = makeDateTimeDisplay0(data, value.Timestamp)
			                }
			
			            case *DateTimeSetCnfUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestDateTimeDisplay = makeDateTimeDisplay1(data, value.Timestamp)
			                }
			
			            case *DateTimeGetCnfUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestDateTimeDisplay = makeDateTimeDisplay2(data, value.Timestamp)
			                }
			
			            case *ModeSetCnfUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestModeDisplay = makeModeDisplay0(data, value.Timestamp)
			                }
			
			            case *ModeGetCnfUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestModeDisplay = makeModeDisplay1(data, value.Timestamp)
			                }
			                
			            case *PollIndUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestUtmStatusDisplay = makeUtmStatusDisplay(data, value.Timestamp)
			                    state.LatestModeDisplay = makeModeDisplay2(data, value.Timestamp)
			                }
			
			            case *MeasurementsIndUlMsg:
			                data := utmMsg.Measurements.DeepCopy()
			                if data != nil {
			                    state.LatestSignalStrengthDisplay = makeSignalStrengthDisplay(data, value.Timestamp)
			                    if (data.GnssPositionPresent) {
			                    	state.LatestGnssDisplay = makeGnssDisplay(data, value.Timestamp)
			                    }
			                    if (data.CellIdPresent) {
			                    	state.LatestCellIdDisplay = makeCellIdDisplay(data, value.Timestamp)
			                    }
			                    if (data.TemperaturePresent) {
			                    	state.LatestTemperatureDisplay = makeTemperatureDisplay(data, value.Timestamp)
			                    }
			                    if (data.PowerStatePresent) {
			                    	state.LatestPowerStateDisplay = makePowerStateDisplay(data, value.Timestamp)
			                    }
			                }
			
			            case *MeasurementsGetCnfUlMsg:
			                data := utmMsg.Measurements.DeepCopy()
			                if data != nil {
			                    state.LatestSignalStrengthDisplay = makeSignalStrengthDisplay(data, value.Timestamp)
			                    if (data.GnssPositionPresent) {
			                    	state.LatestGnssDisplay = makeGnssDisplay(data, value.Timestamp)
			                    }
			                    if (data.CellIdPresent) {
			                    	state.LatestCellIdDisplay = makeCellIdDisplay(data, value.Timestamp)
			                    }
			                    if (data.TemperaturePresent) {
			                    	state.LatestTemperatureDisplay = makeTemperatureDisplay(data, value.Timestamp)
			                    }
			                    if (data.PowerStatePresent) {
			                    	state.LatestPowerStateDisplay = makePowerStateDisplay(data, value.Timestamp)
			                    }
			                }
			
			            case *TrafficReportIndUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestTrafficReportDisplay = makeTrafficReportDisplay0(data, value.Timestamp)
			                }
			
			            case *TrafficReportGetCnfUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestTrafficReportDisplay = makeTrafficReportDisplay1(data, value.Timestamp)
			                }
			
			            case *TrafficTestModeParametersSetCnfUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestTrafficTestModeParametersDisplay = makeTrafficTestModeParametersDisplay0(data, value.Timestamp)
			                }
			
			            case *TrafficTestModeParametersGetCnfUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestTrafficTestModeParametersDisplay = makeTrafficTestModeParametersDisplay1(data, value.Timestamp)
			                }
			
			            case *TrafficTestModeReportIndUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
        			            // TODO check for a pass/fail result 
			                    state.LatestTrafficTestModeReportDisplay = makeTrafficTestModeReportDisplay0(data, value.Timestamp)
			                }
			
			            case *TrafficTestModeReportGetCnfUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestTrafficTestModeReportDisplay = makeTrafficTestModeReportDisplay1(data, value.Timestamp)
			                }
			
			            case *ActivityReportIndUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestActivityReportDisplay = makeActivityReportDisplay0(data, value.Timestamp)
			                }
			
			            case *ActivityReportGetCnfUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestActivityReportDisplay = makeActivityReportDisplay1(data, value.Timestamp)
			                }
	            	}
	            	
	            // Return the latest state for a given UUID 
	            case *DeviceLatestState:
	            	// Retrieve the device state
   	            	state := deviceLatestStateList[value.DeviceUuid]
   	            	
   	            	if state != nil {
		                // Duplicate the memory pointed to into a new LatestState struct,
		                // post it and close the channel
		                Dbg.PrintfTrace("%s --> fetching latest state for UUID %s.\n", logTag, value.DeviceUuid)
		                latest := LatestState{}
		                latest.Connected = state.Connected
		                latest.LastHeardFrom = state.LastHeardFrom
		                latest.LatestInitIndDisplay = state.LatestInitIndDisplay.DeepCopy()
		                latest.LatestIntervalsDisplay = state.LatestIntervalsDisplay.DeepCopy()
		                latest.LatestModeDisplay = state.LatestModeDisplay.DeepCopy()
		                latest.LatestDateTimeDisplay = state.LatestDateTimeDisplay.DeepCopy()
		                latest.LatestUtmStatusDisplay = state.LatestUtmStatusDisplay.DeepCopy()
		                latest.LatestGnssDisplay = state.LatestGnssDisplay.DeepCopy()
		                latest.LatestCellIdDisplay = state.LatestCellIdDisplay.DeepCopy()
		                latest.LatestSignalStrengthDisplay = state.LatestSignalStrengthDisplay.DeepCopy()
		                latest.LatestTemperatureDisplay = state.LatestTemperatureDisplay.DeepCopy()
		                latest.LatestPowerStateDisplay = state.LatestPowerStateDisplay.DeepCopy()
		                latest.LatestTrafficReportDisplay = state.LatestTrafficReportDisplay.DeepCopy()
		                latest.LatestTrafficTestModeParametersDisplay = state.LatestTrafficTestModeParametersDisplay.DeepCopy()
		                latest.LatestTrafficTestModeReportDisplay = state.LatestTrafficTestModeReportDisplay.DeepCopy()
		                latest.LatestActivityReportDisplay = state.LatestActivityReportDisplay.DeepCopy()
		                latest.LatestDisplayRow = state.LatestDisplayRow.DeepCopy()
		                value.Latest <- latest
		                close(value.Latest)
		                Dbg.PrintfTrace("%s --> datatable provided latest state and closed channel.\n", logTag)
		            } else {
		                Dbg.PrintfTrace("%s --> datatable asked for latest state for unknown UUID (%s).\n", logTag, value.DeviceUuid)
		            }
		                
	            default:
	                Dbg.PrintfTrace("%s --> unrecognised datatable message, ignoring:\n\n%s\n", logTag, spew.Sdump(cmd))
            }
        }

        Dbg.PrintfTrace("%s --> datatable command channel closed, stopping.\n", logTag)
        checkConnections.Stop();
    }()
}

func init() {
    operateDataTable()
}

/* End Of File */
