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
var Link = require('react-router-component').Link;

var DisplayRow = React.createClass({
    render: function() {
        var rows = [];
        if (this.props["DeviceData"] && (this.props["DeviceData"].length > 0)) {
            this.props["DeviceData"].forEach(function(uuid, i) {
	            rows.push(       
	                <tr className="even gradeC" key={i}>
	                    <td style={{textAlign: 'center', width: 15}}>
	                        <input type="checkbox" /><br />
	                        <img src="static/dist/assets/images/green.png" alt="logo" style={{maxWidth: 12}} />
	                    </td>
	                    <td style={{width: 250}}>
                            <b>Name:</b> {uuid["DeviceName"]}<br />
                            <b>UUID:</b> {uuid["Uuid"]}<br />
                            <b>Mode:</b> {uuid["Mode"]}<br />
                            <b>Reporting:</b> {uuid["Reporting"]}<br />
                            <b>Heartbeat:</b> {uuid["Heartbeat"]}
	                    </td>
	                    <td style={{width: 170}}>
                            <b>Msgs:</b> {uuid["TotalUlMsgs"]}<br />
                            <b>Bytes:</b> {uuid["TotalUlBytes"]}<br />
                            <b>Last Msg:</b> {uuid["LastUlMsgTime"]}
	                    </td>
	                    <td style={{width: 170}}>
	                        <b>Msgs:</b> {uuid["TotalDlMsgs"]}<br />
	                        <b>Bytes:</b> {uuid["TotalDlBytes"]}<br />
	                        <b>Last Msg:</b> {uuid["LastDlMsgTime"]}
	                    </td>
	                    <td style={{width: 80}}>
	                        <i className="fa fa-signal" /> {uuid["Rsrp"]} dBm<br />
	                        <i className="fa fa-floppy-o" /> {uuid["DiskSpaceLeft"]}<br />
	                        <i className="fa fa-battery-full" /> {uuid["BatteryLevel"]}
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
	                                        <th style={{textAlign: 'center', width: 15}}><input type="checkbox" /></th>
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