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

var BatteryLevel = React.createClass({
    render:function() {
        if (this.props.Percentage.search("> 90") >= 0) {
            return (
                <div>
                    <i className="fa fa-battery-full" /> {this.props.Percentage}
                </div>
            );
	    } else if (this.props.Percentage.search("> 70") >= 0) {
	        return (
                <div>
	                <i className="fa fa-battery-three-quarters" /> {this.props.Percentage}
                </div>
	        );
	    } else if (this.props.Percentage.search("> 50") >= 0) {
	        return (
                <div>
	                <i className="fa fa-battery-half" /> {this.props.Percentage}
                </div>
	        );
	    } else if (this.props.Percentage.search("> 30") >= 0) {
	        return (
                <div>
	                <i className="fa fa-battery-quarter" /> {this.props.Percentage}
                </div>
	        );
	    } else if (this.props.Percentage.search("> 10") >= 0) {
	        return (
                <div style={{color: 'red'}}>
	                <i className="fa fa-battery-empty" /> {this.props.Percentage}
                </div>
	        );
	    } else if (this.props.Percentage.search("< 10") >= 0) {
	        return (
                <div style={{color: 'red'}}>
                    <i className="fa fa-battery-empty" /> {this.props.Percentage}
                </div>
	        );
	    } else if (this.props.Percentage.search("< 5") >= 0) {
	        return (
                <div style={{color: 'red'}}>
                    <i className="fa fa-battery-empty" /> <b>{this.props.Percentage}</b>
                </div>
	        );
        } else { // Just in case
            return (
                <div>
                    <i className="fa fa-battery-full" /> {this.props.Percentage}
                </div>
            );
        }
    }    
});

module.exports = BatteryLevel;