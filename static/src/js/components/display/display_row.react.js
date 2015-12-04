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

var React = require('react');
var ValueUuidSelected = require('../controls/value_uuid_selected.react');
var Connected = require('../controls/value_connected.react');
var AppActions = require('../../actions/app_actions.js');
var Moment = require('moment');
var Link = require('react-router-component').Link;

var DisplayRow = React.createClass({
    getInitialState: function(){   
        return null;
    },
    
    handleCheckAll: function(checkAll) {
        if (this.props["DeviceData"] && (this.props["DeviceData"].length > 0)) {
        	for (i = 0; i < this.props["DeviceData"].length; i++) {
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
        var cellIdTime;
        var rsrpTime;
        var rssiTime;
        var txPowerTime;
        var coverageClassTime;
        
        if (this.props["DeviceData"] && (this.props["DeviceData"].length > 0)) {
            this.props["DeviceData"].forEach(function(device, i) {
            	
            	if (device["CellIdTime"]) {
                    cellIdTime = Moment.utc(Date.parse(device["CellIdTime"])).fromNow()
                    if (cellIdTime == NaN) {
                        celldTime = "";
                    }			
            	}
            	if (device["RsrpTime"]) {
            	    rsrpTime = Moment.utc(Date.parse(device["RsrpTime"])).fromNow();
                    if (rsrpTime == NaN) {
                	    rsrpTime = "dBm"; 
                    } else {
                 	    rsrpTime = "dBm " + rsrpTime; 
                    }
            	}    
            	if (device["RssiTime"]) {                
                    rssiTime = Moment.utc(Date.parse(device["RssiTime"])).fromNow();
                    if (rssiTime == NaN) {
                	    rssiTime = "dBm"; 
                    } else {                	
                	    rssiTime = "dBm " + rssiTime; 
                    }
            	}
            	if (device["TxPowerTime"]) {                
                    txPowerTime = Moment.utc(Date.parse(device["TxPowerTime"])).fromNow();
                    if (txPowerTime == NaN) {
                	    txPowerTime = "dBm"; 
                    } else {
                	    txPowerTime = "dBm " + txPowerTime; 
                    }
            	}                
            	if (device["CoverageClassTime"]) {                
                    coverageClassTime = Moment.isDate(device["CoverageClassTime"])
                    if (coverageClassTime == NaN) {
                    	coverageClassTime = "";
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
                        </td>
                        <td style={{width: 170}}>
                            <i className="fa fa-arrow-up" /> Msgs: <b>{device["TotalUlMsgs"]}</b><br />
                            <i className="fa fa-arrow-up" /> Bytes: <b>{device["TotalUlBytes"]}</b><br />
                            <i className="fa fa-arrow-up" /> Last Msg: {Moment.utc(Date.parse(device["LastUlMsgTime"])).local().format("YYYY-MM-DD HH:mm:ss")}<br />
                            <i className="fa fa-arrow-down" /> Msgs: <b>{device["TotalDlMsgs"]}</b><br />
                            <i className="fa fa-arrow-down" /> Bytes: <b>{device["TotalDlBytes"]}</b><br />
                            <i className="fa fa-arrow-down" /> Last Msg: {Moment.utc(Date.parse(device["LastDlMsgTime"])).local().format("YYYY-MM-DD HH:mm:ss")}
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
                        <td className="center" style={{width: 170}}>
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
                                            <th>Traffic</th>
                                            <th>Radio</th>
                                            <th>Status</th>
                                            <th>Test Results</th>
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