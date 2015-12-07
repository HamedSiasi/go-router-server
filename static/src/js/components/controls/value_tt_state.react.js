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

var TtState = React.createClass({
    render:function() {
        if (this.props.IsRunning) {
            return (
                <div>
                    <img src="static/dist/assets/images/amber.png" style={{maxWidth: 12}} /> {this.props.TimeUpdated} 
                </div>
            );
        } else if (this.props.IsPassed) {
            return (
                <div>
                    <img src="static/dist/assets/images/green.png" style={{maxWidth: 12}} /> {this.props.TimeUpdated} 
                </div>
            );
        } else if (this.props.IsFailed) {
            return (
                <div>
                    <img src="static/dist/assets/images/red.png" style={{maxWidth: 12}} /> {this.props.TimeUpdated} 
                </div>
            );
        } else if (this.props.IsTimedOut) {
            return (
                <div>
                    <img src="static/dist/assets/images/red1.png" style={{maxWidth: 12}} /> {this.props.TimeUpdated} 
                </div>
            );
        } else {
            return (
                <div>
                    <img src="static/dist/assets/images/black.png" style={{maxWidth: 12}} /> {this.props.TimeUpdated}
                </div>
            );
        }
    }
});

module.exports = TtState;