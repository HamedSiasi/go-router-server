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
var AppConstants = require ('../../constants/app-constants');

var ValueTTTimeout = React.createClass({
    getInitialState: function() {
        return {value: AppConstants.TT_TIMEOUT_DEFAULT};
    },

    handleChange: function(newValue) {
    	if ((newValue.target.value >= AppConstants.TT_TIMEOUT_MIN) && (newValue.target.value <= AppConstants.TT_TIMEOUT_MAX)) {
    	    this.setState ({value: newValue.target.value});
        }
    },

    render:function(){
	    var value = this.state.value;
        return (
            <input className="form-control bfh-number" type="number" min={AppConstants.TT_TIMEOUT_MIN} max={AppConstants.TT_TIMEOUT_MAX} value={value} step={60} onChange={this.handleChange} style={{width: 80, marginTop: 10}} />
        );
    }
});

module.exports = ValueTTTimeout;