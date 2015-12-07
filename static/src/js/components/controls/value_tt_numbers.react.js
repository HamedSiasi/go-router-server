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
        
        if (rx > tx) {
        	rx = tx;
        }
        
        if (missed > tx) {
        	missed = tx;
        }
        
        if (this.props.Target) {
            percentComplete = ((rx / this.props.Target) * 100).toFixed(0);
            if (percentComplete > 100) {
            	percentComplete = 100;
            }	
        } else {
        	percentComplete = 0;
        }

        if (this.props.IsUplink) {
            return (
                <div>
                <i className="fa fa-arrow-up" /> <b>{rx}</b> out of <b>{tx}</b> (<b>{percentComplete}%</b>, <b>{missed}</b> missed)
                </div>
            );
        } else {
            return (
                <div>
                <i className="fa fa-arrow-down" /> <b>{rx}</b> out of <b>{tx}</b> (<b>{percentComplete}%</b>, <b>{missed}</b> missed)
                </div>
            );
        }
    }
});

module.exports = TtNumbers;