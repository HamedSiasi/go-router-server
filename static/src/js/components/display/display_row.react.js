/**
 * Copyright (C) u-blox Melbourn Ltd
 * u-blox Melbourn Ltd, Melbourn, UK
 * 
 * All rights reserved.
 *
 * This source file is the sole property of u-blox Melbourn Ltd.
 * Reproduction or utilization of this source in whole or part is
 * forbidden without the written consent of u-blox Melbourn Ltd.
 * 
 * Parts of this file are written in JSX, not HTML.  If you want
 * to put any content in here that should be generated as HTML,
 * stuff it through:
 * 
 * https://facebook.github.io/react/html-jsx.html
 * 
 * ...to get your syntax correct.
 */
"use strict";

var React = require('react');
var ValueUuidSelected = require('../controls/value_uuid_selected.react');
var Connected = require('../controls/value_connected.react');
var TtState = require('../controls/value_tt_state.react');
var TtNumbers = require('../controls/value_tt_numbers.react');
var TtConfig = require('../controls/value_tt_config.react');
var AppActions = require('../../actions/app_actions.js');
var Moment = require('moment');
var Link = require('react-router-component').Link;

var DisplayRow = React.createClass({
    getInitialState: function(){   
        return null;
    },
    
    handleCheckAll: function(checkAll) {
        if (this.props["DeviceData"] && (this.props["DeviceData"].length > 0)) {
        	for (var i = 0; i < this.props["DeviceData"].length; i++) {
            	uuid = this.props["DeviceData"][i]["Uuid"];
                if (checkAll.target.checked == true) {
                    AppActions.setUuidChecked(uuid);
                } else {
                    AppActions.setUuidUnchecked(uuid);
                }
        	}
        }
    },
    
    render: function() {
        var rows = [];
        
        if (this.props["DeviceData"] && (this.props["DeviceData"].length > 0)) {
            
        	this.props["DeviceData"].forEach(function(device, i) {
                var deviceTime;
                var cellIdTime;
                var rsrpTime;
                var rssiTime;
                var txPowerTime;
                var coverageClassTime;
                var ttTimeStarted;
                var ttTimeUpdated;
                var ttTimeStopped;
                var ttDuration;
            	
            	if (device["DeviceTime"]) {
            		deviceTime = Moment(Date.parse(device["DeviceTime"])).format("YYYY-MM-DD HH:mm:ss");
                    if (deviceTime == NaN) {
                    	deviceTime = "";
                    }			
            	}
            	if (device["CellIdTime"]) {
                    cellIdTime = Moment(Date.parse(device["CellIdTime"])).fromNow()
                    if (cellIdTime == NaN) {
                        celldTime = "";
                    }			
            	}
            	if (device["RsrpTime"]) {
            	    rsrpTime = Moment(Date.parse(device["RsrpTime"])).fromNow();
                    if (rsrpTime == NaN) {
                	    rsrpTime = "dBm"; 
                    } else {
                 	    rsrpTime = "dBm " + rsrpTime; 
                    }
            	}    
            	if (device["RssiTime"]) {                
                    rssiTime = Moment(Date.parse(device["RssiTime"])).fromNow();
                    if (rssiTime == NaN) {
                	    rssiTime = "dBm"; 
                    } else {                	
                	    rssiTime = "dBm " + rssiTime; 
                    }
            	}
            	if (device["TxPowerTime"]) {                
                    txPowerTime = Moment(Date.parse(device["TxPowerTime"])).fromNow();
                    if (txPowerTime == NaN) {
                	    txPowerTime = "dBm"; 
                    } else {
                	    txPowerTime = "dBm " + txPowerTime; 
                    }
            	}                
            	if (device["CoverageClassTime"]) {                
                    coverageClassTime = Moment(Date.parse(device["CoverageClassTime"])).fromNow();
                    if (coverageClassTime == NaN) {
                    	coverageClassTime = "";
                    }
            	}
            	
            	if (device["TtTimeStarted"]) {
            		ttTimeStarted = Moment(Date.parse(device["TtTimeStarted"]));
            		if (Moment(ttTimeStarted).isAfter('2014-01-01', 'year')) {
                    	if (device["TtTimeStopped"]) {
                    		ttTimeStopped = Moment(Date.parse(device["TtTimeStopped"]));
                    		if (Moment(ttTimeStopped).isAfter('2014-01-01', 'year')) {
                                ttDuration = Moment(Moment(ttTimeStopped).diff(Moment(ttTimeStarted))).format("HH:mm:ss");
                                ttTimeUpdated = ttTimeStopped;
                    	    } else {
                    	    	ttTimeStopped = ""
                    	    }
                    	} else {
	                        if (device["TtRunning"]) { 
	                            ttDuration = Moment(Moment(Date.now()).diff(Moment(ttTimeStarted))).format("HH:mm:ss");
	                        }
                    	}
	                    ttTimeUpdated = Moment(ttTimeUpdated).fromNow() + " (" + Moment(ttTimeUpdated).format("HH:mm:ss") + ")";
            	    } else {
            	    	ttTimeStarted = ""
            	    }
            	}
            	
                rows.push(
                    <tr className="even gradeC" key={i}>
                        <td style={{textAlign: 'center', width: 15}}>
                            <ValueUuidSelected Uuid={device["Uuid"]} /><br/>
                            <Connected IsConnected={device["Connected"]} />
                        </td>
                        <td style={{width: 250}}>
                            Name: <b>{device["DeviceName"]}</b><br />
                            UUID: {device["Uuid"]}<br />
                            Mode: {device["Mode"]}<br />
                            Reporting: <b>{device["Reporting"]}</b><br />
                            Heartbeat: <b>{device["Heartbeat"]}</b><br />
                            Time: {deviceTime}
                        </td>
                        <td style={{width: 170}}>
                            <i className="fa fa-arrow-up" /> Msgs: <b>{device["TotalUlMsgs"]}</b><br />
                            <i className="fa fa-arrow-up" /> Bytes: <b>{device["TotalUlBytes"]}</b><br />
                            <i className="fa fa-arrow-up" /> Last Msg: {Moment(Date.parse(device["LastUlMsgTime"])).format("YYYY-MM-DD HH:mm:ss")}<br />
                            <i className="fa fa-arrow-down" /> Msgs: <b>{device["TotalDlMsgs"]}</b><br />
                            <i className="fa fa-arrow-down" /> Bytes: <b>{device["TotalDlBytes"]}</b><br />
                            <i className="fa fa-arrow-down" /> Last Msg: {Moment(Date.parse(device["LastDlMsgTime"])).format("YYYY-MM-DD HH:mm:ss")}
                        </td>
                        <td style={{width: 200}}>
                            <i className="fa fa-rss" /> Cell: <b>{device["CellId"]}</b> {cellIdTime}<br />
                            <i className="fa fa-signal" /> RSRP: <b>{device["Rsrp"]}</b> {rsrpTime}<br />
                            <i className="fa fa-signal" /> RSSI: <b>{device["Rssi"]}</b> {rssiTime}<br />
                            <i className="fa fa-bolt" /> Tx: <b>{device["TxPower"]}</b> {txPowerTime}<br />
                            <i className="fa fa-globe" /> Coverage Class: <b>{device["CoverageClass"]}</b> {coverageClassTime}<br />
                        </td>
                        <td style={{width: 80}}>
                            <i className="fa fa-floppy-o" /> {device["DiskSpaceLeft"]}<br />
                            <i className="fa fa-battery-full" /> {device["BatteryLevel"]}<br />
                            <i className="fa fa-clock-o" /> {device["UpDuration"]}<br />
                            <i className="fa fa-arrow-up" /> {device["TxTime"]}<br />
                            <i className="fa fa-arrow-down" /> {device["RxTime"]}<br />
                            <i className="fa fa-rotate-left" /> {device["NumExpectedMsgs"]}
                            </td> 
                        <td className="center" style={{width: 200}}>
                            <TtConfig UlDatagrams={device["TtUlExpected"]} UlLength={device["TtUlLength"]} DlDatagrams={device["TtDlExpected"]} DlLength={device["TtDlLength"]} DlInterval={device["TtDlInterval"]} Timeout={device["TtTimeout"]}/>
                            <TtState IsRunning={device["TtRunning"]} IsPassed={device["TtPassed"]} IsFailed={device["TtFailed"]} IsTimedOut={device["TtTimedOut"]} TimeUpdated={ttTimeUpdated} />
                            <TtNumbers IsUplink={true} Tx={device["TtUlDatagramsTx"]} Rx={device["TtUlDatagramsRx"]} Missed={device["TtUlDatagramsMissed"]} Target={device["TtUlExpected"]}/>
                            <TtNumbers IsUplink={false} Tx={device["TtDlDatagramsTx"]} Rx={device["TtDlDatagramsRx"]} Missed={device["TtDlDatagramsMissed"]} Target={device["TtDlExpected"]}/>
                            <i className="fa fa-clock-o" /> {ttDuration}
                        </td> 
                    </tr>
                );
            });
        }
        
        return (        
            <div className="row" >
                <div className="panel panel-default" >
                    <div className="_panel-heading" style={{width:'100%'}}>
                        <div className="panel-body">
                            <div className="dataTable_wrapper">
                                <table className="table table-striped table-bordered table-hover" id="dataTables-example">
                                    <thead>
                                        <tr className="info">
                                            <th style={{textAlign: 'center', width: 15}}><input type="checkbox" onChange={this.handleCheckAll} value="checkAll" defaultChecked={false} /></th>
                                            <th>Device</th>
                                            <th>Normal Traffic</th>
                                            <th>Radio</th>
                                            <th>Status</th>
                                            <th>Test Traffic</th>
                                        </tr>
                                    </thead>
                                    <tbody style={{fontSize: 12}}>
                                        {rows}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
});

module.exports = DisplayRow;