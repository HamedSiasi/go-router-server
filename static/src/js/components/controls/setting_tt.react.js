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
var ValueTTUlNumDatagrams = require('./value_tt_ul_num_datagrams.react');
var ValueTTUlLenDatagram = require('./value_tt_ul_len_datagram.react');
var ValueTTDlNumDatagrams = require('./value_tt_dl_num_datagrams.react');
var ValueTTDlLenDatagram = require('./value_tt_dl_len_datagram.react');
var ValueTTTimeout = require('./value_tt_timeout.react');
var SetTTParameters = require('./set_tt_parameters.react');
var SetTTStart = require('./set_tt_start.react');
var SetTTStop = require('./set_tt_stop.react');

var SettingTT = React.createClass({
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
                        <td style={{width: 100}}><ValueTTUlNumDatagrams /></td>
                        <td style={{width: 110}}>Length:</td>
                        <td style={{width: 100}}><ValueTTUlLenDatagram /></td> 
                        <td><SetTTStart Names={this.props.Names} UuidMap = {this.props.UuidMap} /></td>
                    </tr>
                    <tr >
                        <td>DL: Number:</td>
                        <td><ValueTTDlNumDatagrams /></td>
                        <td>Length:</td>
                        <td><ValueTTDlLenDatagram /></td> 
                        <td><SetTTStop Names={this.props.Names} UuidMap = {this.props.UuidMap} /></td> 
                    </tr>
                    <tr >
                        <td>Timeout (secs):</td>
                        <td><ValueTTTimeout /></td>
                        <td colSpan={2}><SetTTParameters Names={this.props.Names} UuidMap = {this.props.UuidMap} /></td> 
                    </tr>
                </table>
            </div>
        );
    }
});

module.exports = SettingTT;