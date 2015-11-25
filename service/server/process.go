/* Message processing functions for the UTM server.
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
    "reflect"
    "github.com/davecgh/go-spew/spew"
	"github.com/robmeades/utm/service/globals"
)

//--------------------------------------------------------------------
// Types 
//--------------------------------------------------------------------

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

var processMsgs chan<- interface{}

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

func operateProcess() {
	deviceList := make(map[string]bool)    	
    channel := make(chan interface{})
    processMsgs = channel
    
    globals.Dbg.PrintfTrace("%s [process] --> channel created and now being serviced.\n", globals.LogTag)
    
    // Process commands on the channel
    go func() {
        for cmd := range channel {
            switch value := cmd.(type) {
            	
	            // Handle message containers holding somethings of interest,
	            // throw everything else away
	            case *MessageContainer:
    		        responseId := RESPONSE_NONE

            		// If the device is not in our list, add it and send an IntervalsGetReq
            		// if we aren't going to send one anyway lower down because of an InitIndUlmsg
            		// (or of course, if this is already an IntervalsGetCnfUlMsg)
            		if deviceList[value.DeviceUuid] == false {
            		    deviceList[value.DeviceUuid] = true
            		    if reflect.TypeOf (value.Message) != reflect.TypeOf((*InitIndUlMsg)(nil)).Elem() &&
            		       reflect.TypeOf (value.Message) != reflect.TypeOf((*IntervalsGetCnfUlMsg)(nil)).Elem() {
    	                    encodeAndEnqueue (&IntervalsGetReqDlMsg{}, value.DeviceUuid)        		           
            		       }
            		}
            		
					globals.Dbg.PrintfTrace("%s [process] --> processing message from UUID %s...\n\n%s\n\n", globals.LogTag, value.DeviceUuid, spew.Sdump(cmd))
            		
	            	switch utmMsg := value.Message.(type) {
        				
        				case *TransparentUlDatagram:
			                // Nothing to do here
        				
        				case *PingReqUlMsg:
        					// Respond
		                    encodeAndEnqueue (&PingCnfDlMsg{}, value.DeviceUuid)
        				
        				case *PingCnfUlMsg:
			                // Nothing to do here
			                responseId = RESPONSE_PING_CNF

			            case *InitIndUlMsg:
        					globals.Dbg.PrintfInfo("%s [process] --> UUID %s has protocol revision %d, which is different to this server (%d)", globals.LogTag, value.DeviceUuid, utmMsg.RevisionLevel, RevisionLevel)
        					// Get the reporting intervals for this device
		                    encodeAndEnqueue (&IntervalsGetReqDlMsg{}, value.DeviceUuid)
        					
			            case *DateTimeIndUlMsg:
			                // Send DataTimeSetReqDlMsg if out of range
			                // First check if the time is OK but the date is absent
			                // (which could be the case if the time has been set by GNSS,
			                // and we don't want to mess with that time as it will be more
			                // accurate than ours)
			                utmDays := utmMsg.UtmTime.Year() * 365 +  utmMsg.UtmTime.YearDay()
			                serverDays := time.Now().UTC().Year() * 365 + time.Now().UTC().YearDay()
		                    // Allow plenty of slack as messages might have been queued	                
			                if (utmDays < serverDays - 30) || (utmDays > serverDays + 30) {
		                        dateTimeSetReq := &DateTimeSetReqDlMsg {
                                    UtmTime:      time.Now().UTC(),
                                    SetDateOnly:  false,
                                }
		                        
			                    if utmMsg.TimeSetBy == TIME_SET_BY_GNSS {
			                        dateTimeSetReq.SetDateOnly = true;
			                    }
			                    
			                    encodeAndEnqueue (dateTimeSetReq, value.DeviceUuid)
    		                }
		                   
			            case *DateTimeSetCnfUlMsg:
			                // Nothing to do here
			                responseId = RESPONSE_DATE_TIME_SET_CNF
			                
			            case *DateTimeGetCnfUlMsg:
			                // Nothing to do here
                            responseId = RESPONSE_DATE_TIME_GET_CNF
                            
			            case *ModeSetCnfUlMsg:
	    		            // TODO start sending downlink datagrams if in a traffic test
			                responseId = RESPONSE_MODE_SET_CNF
			                
			            case *ModeGetCnfUlMsg:
			                // Nothing to do here
			                responseId = RESPONSE_MODE_GET_CNF
			                
			            case *IntervalsGetCnfUlMsg:
			                // Nothing to do here
			                responseId = RESPONSE_INTERVALS_GET_CNF
			                
			            case *ReportingIntervalSetCnfUlMsg:
			                // Nothing to do here
			                responseId = RESPONSE_REPORTING_INTERVAL_SET_CNF
			            
			            case *HeartbeatSetCnfUlMsg:
			                // Nothing to do here
			                responseId = RESPONSE_HEARTBEAT_SET_CNF

			            case *PollIndUlMsg:
			                // Nothing to do here
			            
			            case *MeasurementsIndUlMsg:
			                // Nothing to do here
			            
			            case *MeasurementsGetCnfUlMsg:
			                // Nothing to do here
			                responseId = RESPONSE_MEASUREMENTS_GET_CNF

        				// case *MeasurementsControlIndUlMsg:
        				// case *MeasurementControlSetCnfUlMsg:
        				// case *MeasurementsControlGetCnfUlMsg:
        				// case *MeasurementsControlDefaultsSetCnfUlMsg:
        				// TODO

			            case *TrafficReportIndUlMsg:
			                // Nothing to do here
			            
			            case *TrafficReportGetCnfUlMsg:
			                // Nothing to do here
			                responseId = RESPONSE_TRAFFIC_REPORT_GET_CNF

			            case *TrafficTestModeParametersSetCnfUlMsg:
    			            // TODO send ModeSetReqDlMsg if in a traffic test 
			                responseId = RESPONSE_TRAFFIC_TEST_MODE_PARAMETERS_SET_CNF
			                
			            case *TrafficTestModeParametersGetCnfUlMsg:
			                // Nothing to do here
			                responseId = RESPONSE_TRAFFIC_TEST_MODE_PARAMETERS_GET_CNF
			                
        				case *TrafficTestModeRuleBreakerUlDatagram:
        				    // TODO
        				
			            case *TrafficTestModeReportIndUlMsg:
			                // Nothing to do here
			            
			            case *TrafficTestModeReportGetCnfUlMsg:
			                // Nothing to do here
			                responseId = RESPONSE_TRAFFIC_TEST_MODE_REPORT_GET_CNF

			            case *ActivityReportIndUlMsg:
			                // Nothing to do here
			            
			            case *ActivityReportGetCnfUlMsg:
			                // Nothing to do here
			                responseId = RESPONSE_ACTIVITY_REPORT_GET_CNF

            			case *DebugIndUlMsg:
			                // Nothing to do here
			                
			            default:
    			            // Ignore any unknown UTM messages
        	                globals.Dbg.PrintfTrace("%s [process] --> unrecognised UTM message, ignoring.\n", globals.LogTag)
        	                globals.Dbg.PrintfInfo("%s [process] --> unrecognised UTM message was:\n\n%s\n", globals.LogTag, spew.Sdump(utmMsg))
	            	}
	            	
                	// If this was a response message, take it out of the expected list for this UUID
                	if responseId != RESPONSE_NONE {
    					globals.Dbg.PrintfTrace("%s [process] --> response ID %d received from UUID %s.\n", globals.LogTag, responseId, value.DeviceUuid)
                    	list := deviceExpectedMsgList[value.DeviceUuid]
                		if list != nil {
        					globals.Dbg.PrintfTrace("%s [process] --> found UUID %s in the expected message store, %d items in its list.\n", globals.LogTag, value.DeviceUuid, len(*list))
                		    for index, expectedMsg := range *list {
                		        if expectedMsg.ResponseId == responseId {
                		            *list = append((*list)[:index], (*list)[index + 1:] ...)
                					globals.Dbg.PrintfTrace("%s [process] --> response ID %d removed from list for UUID %s.\n", globals.LogTag, responseId, value.DeviceUuid)
                		            break
                		        }
                		    }
                		}
                	}
                	
					globals.Dbg.PrintfTrace("%s [process] --> processing completed.\n", globals.LogTag)
	            	
	            default:
	                globals.Dbg.PrintfTrace("%s [process] --> unrecognised command, ignoring.\n", globals.LogTag)
	                globals.Dbg.PrintfInfo("%s [process] --> unrecognised command was:\n\n%s\n", globals.LogTag, spew.Sdump(cmd))
            }
        }

        globals.Dbg.PrintfTrace("%s [process] --> command channel closed, stopping.\n", globals.LogTag)
    }()
}

func init() {
    operateProcess()
}

/* End Of File */
