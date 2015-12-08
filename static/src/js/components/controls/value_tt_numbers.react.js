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

var TtNumbers = React.createClass({
    render:function() {
        var percentComplete;
        var tx = this.props.Tx;
        var rx = this.props.Rx
        var missed = this.props.Missed;
        var directionIcon = "fa fa-arrow-down";
        var stateIcon = "static/dist/assets/images/black.png";
        
        if (this.props.IsUplink) {
        	directionIcon = "fa fa-arrow-up";
        }
        
        // The reports of how many datagrams have been transmitted
        // can be behind the number received (because we're waiting
        // for a report from the device telling us how many have
        // been sent), so adjust for that to stop things looking
        // peculiar
        if (tx < rx) {
        	tx = rx;
        }
        
        if (this.props.Target) {
            percentComplete = ((rx / this.props.Target) * 100).toFixed(0);
            if (percentComplete > 100) {
            	percentComplete = 100;
            }	
        } else {
        	percentComplete = 0;
        }
        
        switch (this.props.State) {
            case 1: //TRAFFIC_TEST_RUNNING"
                stateIcon = "static/dist/assets/images/amber.png";
            break;
            case 2: //TRAFFIC_TEST_TX_COMPLETE"
                stateIcon = "static/dist/assets/images/amber.png";
            break;
            case 3: //TRAFFIC_TEST_STOPPED"
                stateIcon = "static/dist/assets/images/black.png";
            break;
            case 4: //TRAFFIC_TEST_TIMEOUT"
                stateIcon = "static/dist/assets/images/red1.png";
            break;
            case 5: //TRAFFIC_TEST_PASS"
                stateIcon = "static/dist/assets/images/green.png";
            break;
            case 6: //TRAFFIC_TEST_FAIL"
                stateIcon = "static/dist/assets/images/red.png";
            break;
            default:
                stateIcon = "static/dist/assets/images/black.png";
            break;
        }

        return (
            <div>
            <img src={stateIcon} style={{maxWidth: 12}} /> <i className={directionIcon} /> <b>{rx}</b> out of <b>{tx}</b> (<b>{percentComplete}%</b>, <b>{missed}</b> missed)
            </div>
        );
    }
});

module.exports = TtNumbers;