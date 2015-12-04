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
var ValueTtNumUlDatagrams = require('./value_tt_num_ul_datagrams.react');
var ValueTtLenUlDatagram = require('./value_tt_len_ul_datagram.react');
var ValueTtNumDlDatagrams = require('./value_tt_num_dl_datagrams.react');
var ValueTtLenDlDatagram = require('./value_tt_len_dl_datagram.react');
var ValueTtTimeout = require('./value_tt_timeout.react');
var ValueTtDlInterval = require('./value_tt_dl_interval.react');
var SetTtParameters = require('./set_tt_parameters.react');
var SetTtStart = require('./set_tt_start.react');
var SetTtStop = require('./set_tt_stop.react');

var SettingTt = React.createClass({
    render: function() {
        return (
            <div >
                <table >
                    <thead>
                        <tr >
                            <th colSpan={5} style={{textAlign: 'center'}}>Traffic Test Mode Settings</th>
                        </tr>
                    </thead>
                    <tr style={{height: 50}}>
                        <td style={{width: 110}}>UL: Number:</td>
                        <td style={{width: 100}}><ValueTtNumUlDatagrams /></td>
                        <td style={{width: 110}}>Length:</td>
                        <td style={{width: 100}}><ValueTtLenUlDatagram /></td> 
                        <td><SetTtStart Names={this.props.Names} UuidMap = {this.props.UuidMap} /></td>
                    </tr>
                    <tr >
                        <td>DL: Number:</td>
                        <td><ValueTtNumDlDatagrams /></td>
                        <td>Length:</td>
                        <td><ValueTtLenDlDatagram /></td> 
                        <td><SetTtStop Names={this.props.Names} UuidMap = {this.props.UuidMap} /></td> 
                    </tr>
                    <tr >
                        <td>Timeout (secs):</td>
                        <td><ValueTtTimeout /></td>
                        <td>DL Gap (secs):</td>
                        <td><ValueTtDlInterval /></td>
                        <td><SetTtParameters Names={this.props.Names} UuidMap = {this.props.UuidMap} /></td> 
                    </tr>
                </table>
            </div>
        );
    }
});

module.exports = SettingTt;