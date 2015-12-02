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
var SetHeartbeat = require('./set_heartbeat.react');
var ValueHeartbeat = require('./value_heartbeat.react');
var SetReporting = require('./set_reporting.react');
var ValueReporting = require('./value_reporting.react');
var GetIntervals = require('./get_intervals.react');
var GetTime = require('./get_time.react');
var GetPing = require('./get_ping.react');

var SettingStd = React.createClass({
    render: function() {
        return (
            <div >
                <table >
                    <thead>
                        <tr >
                            <th colSpan={3} style={{textAlign: 'center'}}>General Settings</th>
                        </tr>
                    </thead>
                    <tr style={{height: 50}}>
                        <td style={{width: 170}}> <ValueHeartbeat /></td>
                        <td style={{width: 300}}> <SetHeartbeat Names={this.props.Names} UuidMap = {this.props.UuidMap} /></td> 
                        <td></td> 
                    </tr>
                    <tr >
                        <td><ValueReporting /></td>
                        <td><SetReporting Names={this.props.Names} UuidMap = {this.props.UuidMap} /></td> 
                        <td></td> 
                    </tr>
                    <tr >
                        <td colSpan={3}>
                            <GetIntervals Names={this.props.Names} UuidMap = {this.props.UuidMap} />
                            <GetTime Names={this.props.Names} UuidMap = {this.props.UuidMap} />
                            <GetPing Names={this.props.Names} UuidMap = {this.props.UuidMap} />
                        </td>
                    </tr>
                </table>
            </div>
        );
    }
});

module.exports = SettingStd;