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
var Moment = require('moment');

var Summary = React.createClass({
    render: function() {
	    if (this.props["SummaryData"]) {
            var summaryData = this.props["SummaryData"];
            var devicesConnected = summaryData["DevicesConnected"];
            var devicesKnown = summaryData["DevicesKnown"];
            var lastUlTime = "";
            var totalUlBytes = summaryData["TotalUlBytes"];
            var totalDlBytes = summaryData["TotalDlBytes"];
            var numExpectedMsgs = summaryData["NumExpectedMsgs"];
            if (summaryData["LastUlMsgTime"]) {
            	lastUlTime = ", last uplink " + Moment.utc(Date.parse(summaryData["LastUlMsgTime"])).fromNow();
            }
            return (
                <p className="align-right">Summary: <b>{devicesConnected}</b> device(s) connected (<b>{devicesKnown}</b> known){lastUlTime}, <b>{numExpectedMsgs}</b> confirmation(s) outstanding.</p>
            )
        } else {
            return (
                <div>
                    &nbsp;
                </div>
            )
	    }
    }
});

module.exports = Summary;