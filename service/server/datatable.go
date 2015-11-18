/* Datatable definitions for the UTM server.
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
    "github.com/davecgh/go-spew/spew"
)

type LatestState struct {
    Connection                              Connection
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

// To update the latest values send a MessageContainer into this channel
// containing the received message; a copy their contents will be stored
// in a displayable form

// To get the latest state send a '*chan DeviceLatestState' into this channel
// containing the device UUID and a pointer to a channel down which to send
// the LatestState struct; a copy of all quantities will be copied into the
// struct and then the channel will be closed.

// To terminate execution simply close the channel

var stateTableCmds chan<- interface{}

func OperateStateTable() {
    cmds := make(chan interface{})
    stateTableCmds = cmds
    Dbg.PrintfTrace("%s --> datatable command channel created and now being serviced.\n", logTag)
    
    go func() {
    	deviceLatestStateList := make(map[string]*LatestState)
    	
        for msg := range cmds {
            Dbg.PrintfTrace("%s --> datatable command:\n\n%+v\n\n", logTag, msg)

            switch value := msg.(type) {
            	
	            // Handle message containers holding somethings of interest
	            case *MessageContainer:
	            
	            	// Make sure there's an entry for this device
	            	state := deviceLatestStateList[value.DeviceUuid]
            		if state == nil {
            			state = &LatestState{};
            			deviceLatestStateList[value.DeviceUuid] = state;
                        Dbg.PrintfTrace("%s --> datatable has heard from a new device, UUID %s.\n", logTag, value.DeviceUuid)
            		}
                    Dbg.PrintfTrace("%s --> datatable received a message from device with UUID %s.\n", logTag, value.DeviceUuid)
            		
	            	switch utmMsg := value.Message.(type) {
			            case *Connection:
			            	// TODO sort this
			                //state.Connection = *value
			                //Dbg.PrintfTrace("%s --> setting connection state for all devices in datatable: %+v\n", logTag, value)
			                	
			            case *InitIndUlMsg:
			                data := utmMsg.DeepCopy()
			                if data != nil {
			                    state.LatestInitIndDisplay = makeInitIndDisplay(data, value.Timestamp)
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
		                state.Connection = state.Connection
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
	                Dbg.PrintfTrace("%s --> unrecognised datatable message, ignoring:\n\n%s\n", logTag, spew.Sdump(msg))
            }
        }

        Dbg.PrintfTrace("%s --> datatable command channel closed, stopping.\n", logTag)
    }()
}

func init() {
    OperateStateTable()
}

/* End Of File */
