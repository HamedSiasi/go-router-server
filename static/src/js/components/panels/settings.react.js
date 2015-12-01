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
var SettingStd = require('../controls/setting_std.react');
var SettingTT = require('../controls/setting_tt.react');

var Settings = React.createClass({
    render: function() {
            var Names = [];
            var UuidMap = {}; 
            if (this.props["DeviceData"] && (this.props["DeviceData"].length > 0)) {        	
                this.props["DeviceData"].forEach(function(device, i) {
                	Names.push(device["DeviceName"]);
                	UuidMap[device["DeviceName"]] = device["Uuid"];
                });
            }
            return (
            <div className="row" >
		        <div className="col-lg-5">
		            <div className="panel panel-info" style={{height: 180, marginTop: 10}}>
			            <div className="panel-body">
                            <SettingStd Names = {Names} UuidMap = {UuidMap} />
			            </div>
		            </div>
	            </div>
		        <div className="col-lg-7">
		            <div className="panel panel-info" style={{height: 180, marginTop: 10}}>
			            <div className="panel-body">
                            <SettingTT Names = {Names} UuidMap = {UuidMap} />
			            </div>
		            </div>
                </div>
	        </div>
	    )
    }
});

module.exports = Settings;