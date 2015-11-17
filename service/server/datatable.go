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
    "log"
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

// To update the latest values send any received message into this channel; a copy
// their contents will be stored in a displayable form

// To get the latest state send a '*chan LatestState' into this channel; a LatestState struct containing
// a copy of all quantities (the memory contents pointed to will never be changed by the
// state table) will be put onto the channel and then the channel closed

// To terminate execution simply close chan
var stateTableCmds chan<- interface{}

func operateStateTable() {
    cmds := make(chan interface{})
    stateTableCmds = cmds
    log.Printf("%s --> datatable command channel created and now being serviced.\n", logTag)
    
    go func() {
        state := LatestState{}
        for msg := range cmds {
            log.Printf("%#v --> datatable command:\n\n%+v\n\n", logTag, msg)

            switch value := msg.(type) {
	            case *Connection:
	                state.Connection = *value
	                log.Printf("%s --> setting connection state in datatable: %+v\n", logTag, value)
	
	            case *InitIndUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestInitIndDisplay = makeInitIndDisplay(data)
	                }
	
	            case *IntervalsGetCnfUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestIntervalsDisplay = makeIntervalsDisplay0(data)
	                }
	
	            case *ReportingIntervalSetCnfUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestIntervalsDisplay = makeIntervalsDisplay1(data)
	                }
	
	            case *HeartbeatSetCnfUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestIntervalsDisplay = makeIntervalsDisplay2(data)
	                }
	
	            case *DateTimeIndUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestDateTimeDisplay = makeDateTimeDisplay0(data)
	                }
	
	            case *DateTimeSetCnfUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestDateTimeDisplay = makeDateTimeDisplay1(data)
	                }
	
	            case *DateTimeGetCnfUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestDateTimeDisplay = makeDateTimeDisplay2(data)
	                }
	
	            case *ModeSetCnfUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestModeDisplay = makeModeDisplay0(data)
	                }
	
	            case *ModeGetCnfUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestModeDisplay = makeModeDisplay1(data)
	                }
	                
	            case *PollIndUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestUtmStatusDisplay = makeUtmStatusDisplay(data)
	                    state.LatestModeDisplay = makeModeDisplay2(data)
	                }
	
	            case *MeasurementsIndUlMsg:
	                data := value.Measurements.DeepCopy()
	                if data != nil {
	                    state.LatestSignalStrengthDisplay = makeSignalStrengthDisplay(data)
	                    if (data.GnssPositionPresent) {
	                    	state.LatestGnssDisplay = makeGnssDisplay(data)
	                    }
	                    if (data.CellIdPresent) {
	                    	state.LatestCellIdDisplay = makeCellIdDisplay(data)
	                    }
	                    if (data.TemperaturePresent) {
	                    	state.LatestTemperatureDisplay = makeTemperatureDisplay(data)
	                    }
	                    if (data.PowerStatePresent) {
	                    	state.LatestPowerStateDisplay = makePowerStateDisplay(data)
	                    }
	                }
	
	            case *MeasurementsGetCnfUlMsg:
	                data := value.Measurements.DeepCopy()
	                if data != nil {
	                    state.LatestSignalStrengthDisplay = makeSignalStrengthDisplay(data)
	                    if (data.GnssPositionPresent) {
	                    	state.LatestGnssDisplay = makeGnssDisplay(data)
	                    }
	                    if (data.CellIdPresent) {
	                    	state.LatestCellIdDisplay = makeCellIdDisplay(data)
	                    }
	                    if (data.TemperaturePresent) {
	                    	state.LatestTemperatureDisplay = makeTemperatureDisplay(data)
	                    }
	                    if (data.PowerStatePresent) {
	                    	state.LatestPowerStateDisplay = makePowerStateDisplay(data)
	                    }
	                }
	
	            case *TrafficReportIndUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestTrafficReportDisplay = makeTrafficReportDisplay0(data)
	                }
	
	            case *TrafficReportGetCnfUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestTrafficReportDisplay = makeTrafficReportDisplay1(data)
	                }
	
	            case *TrafficTestModeParametersSetCnfUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestTrafficTestModeParametersDisplay = makeTrafficTestModeParametersDisplay0(data)
	                }
	
	            case *TrafficTestModeParametersGetCnfUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestTrafficTestModeParametersDisplay = makeTrafficTestModeParametersDisplay1(data)
	                }
	
	            case *ActivityReportIndUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestActivityReportDisplay = makeActivityReportDisplay0(data)
	                }
	
	            case *ActivityReportGetCnfUlMsg:
	                data := value.DeepCopy()
	                if data != nil {
	                    state.LatestActivityReportDisplay = makeActivityReportDisplay1(data)
	                }
	
	            case *chan LatestState:
	                // Duplicate the memory pointed to into a new LatestState struct,
	                // post it and close the channel
	                log.Printf("%s --> datatable fetching latest state.\n", logTag)
	                latest := LatestState{}
	                latest.Connection = state.Connection
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
	                *value <- latest
	                close(*value)
	                log.Print("%s --> datatable provided latest state and closed channel.\n", logTag)
	
	            default:
	                log.Printf("%s --> unrecognised datatable message, ignoring:\n\n%s\n", logTag, spew.Sdump(msg))
            }
        }

        log.Printf("%s --> datatable command channel closed, stopping.\n", logTag)
    }()
}

func init() {
    operateStateTable()
}

/* End Of File */
