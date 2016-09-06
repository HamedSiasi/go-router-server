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
var AppActions = require('../../actions/app_actions.js');
var AppConstants = require ('../../constants/app_limits')
var AppStore = require('../../stores/app_store.js');

var ValueTtLenDlDatagram = React.createClass({
    getInitialState: function() {
        AppActions.setTtLenDlDatagram(AppConstants.TT_DATAGRAM_LEN_DEFAULT);
        return {value: AppConstants.TT_DATAGRAM_LEN_DEFAULT};
    },

    componentDidMount: function() {
        AppStore.addChangeListener(this.onChange);
    },

    componentWillUnmount: function() {
        AppStore.removeChangeListener(this.onChange);
    },

    onChange: function() {
        this.setState({value: AppStore.getTtParameters()["lenDlDatagram"]});
    },

    handleChange: function(newValue) {
        this.setState ({value: newValue.target.valueAsNumber});
    },

    handleBlur: function(newValue) {
	    var tmp = newValue.target.valueAsNumber;
        if (!tmp) {
    	    tmp = AppConstants.TT_DATAGRAM_LEN_MIN;
        }
    
        if (tmp < AppConstants.TT_DATAGRAM_LEN_MIN) {
            tmp = AppConstants.TT_DATAGRAM_LEN_MIN;
        }
        if (tmp > AppConstants.TT_DATAGRAM_LEN_MAX) {
    	    tmp = AppConstants.TT_DATAGRAM_LEN_MAX;
        }
    
        this.setState ({value: tmp});
        AppActions.setTtLenDlDatagram(tmp);        	
    },

    render:function(){
        var value = this.state.value;
        return (
            <input className="form-control bfh-number" type="number" value={value} step={5}  onChange={this.handleChange} onClick={this.handleBlur} onBlur={this.handleBlur} style={{width: 80}} />
        );
    }
});

module.exports = ValueTtLenDlDatagram;