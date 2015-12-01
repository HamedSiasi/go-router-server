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
 * This file is written in JSX, not HTML.  If you want to put any
 * content in here that should be generated as HTML, stuff it
 * through:
 * 
 * https://facebook.github.io/react/html-jsx.html
 * 
 * ...to get your syntax correct.
 */

var React = require('react');
var ValueUuidSelected = require('../controls/value_uuid_selected.react');
var Link = require('react-router-component').Link;

var DisplayRow = React.createClass({
    getInitialState: function(){   
	    return {deviceCheckedMap: {}};
    },
    
    handleCheckAll: function(checkAll) {
    	var deviceCheckedMap = this.state.deviceCheckedMap;
    	if (deviceCheckedMap && (deviceCheckedMap.length > 0)) {
    		for (var key in deviceCheckedMap) {
    		    if (deviceCheckedMap.hasOwnProperty(key)) {
    			    deviceCheckedMap[key] = checkAll.target.value.checked;
    			}
    		}
    	}
    },
    
	handleCheckOne: function(checkOne) {
        var deviceCheckedMap = this.state.deviceCheckedMap;
        if (deviceCheckedMap && (deviceCheckedMap.length > 0)) {
            var checked = deviceCheckedMap[checkOne.value];
            if (checked != null) {
        	    deviceCheckedMap[checkOne.value] = checkOne.checked;
            }
        }
    }.bind(this),

    render: function() {
    	var rows = [];
    	var deviceCheckedMap = this.state.deviceCheckedMap;
        if (this.props["DeviceData"] && (this.props["DeviceData"].length > 0)) {
            this.props["DeviceData"].forEach(function(device, i) {
            	var checked = deviceCheckedMap[device["Uuid"]];
            	if (checked == null) {
            		deviceCheckedMap[device["Uuid"]] = false;
            		checked = false;
            	}
            	rows.push(       
	                <tr className="even gradeC" key={i}>
	                    <td style={{textAlign: 'center', width: 15}}>
	                        <ValueUuidSelected Checked={checked} Uuid={device["Uuid"]} CallbackParent={this.handleCheckOne} /><br/> 
	                        <img src="static/dist/assets/images/green.png" alt="logo" style={{maxWidth: 12}} />
	                    </td>
	                    <td style={{width: 250}}>
                            <b>Name:</b> {device["DeviceName"]}<br />
                            <b>UUID:</b> {device["Uuid"]}<br />
                            <b>Mode:</b> {device["Mode"]}<br />
                            <b>Reporting:</b> {device["Reporting"]}<br />
                            <b>Heartbeat:</b> {device["Heartbeat"]}
	                    </td>
	                    <td style={{width: 170}}>
                            <b>Msgs:</b> {device["TotalUlMsgs"]}<br />
                            <b>Bytes:</b> {device["TotalUlBytes"]}<br />
                            <b>Last Msg:</b> {device["LastUlMsgTime"]}
	                    </td>
	                    <td style={{width: 170}}>
	                        <b>Msgs:</b> {device["TotalDlMsgs"]}<br />
	                        <b>Bytes:</b> {device["TotalDlBytes"]}<br />
	                        <b>Last Msg:</b> {device["LastDlMsgTime"]}
	                    </td>
	                    <td style={{width: 80}}>
	                        <i className="fa fa-signal" /> {device["Rsrp"]} dBm<br />
	                        <i className="fa fa-floppy-o" /> {device["DiskSpaceLeft"]}<br />
	                        <i className="fa fa-battery-full" /> {device["BatteryLevel"]}
	                    </td> 
	                    <td className="center" style={{width: 200}}>
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
	                                        <th style={{textAlign: 'center', width: 15}}><input type="checkbox" onClick={this.handleCheckAll} value="checkAll" defaultChecked={false} /></th>
	                                        <th>Device</th>
	                                        <th>Uplink</th>
	                                        <th>Downlink</th>
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