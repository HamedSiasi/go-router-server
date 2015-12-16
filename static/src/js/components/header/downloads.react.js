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
var Link = require('react-router-component').Link;
var AppStore = require('../../stores/app_store.js');

var Downloads = React.createClass({	
	getInitialState: function () {
	    return {okToContinue: AppStore.getIsLoggedIn()}
    },
	
    componentDidMount: function() {
        AppStore.addChangeListener(this.onChange);
    },

    componentWillUnmount: function() {
        AppStore.removeChangeListener(this.onChange);
    },

    onChange: function() {
        this.setState({okToContinue: AppStore.getIsLoggedIn()});
    },
    
    render: function() {
    	if (this.state.okToContinue) {
            return (
                <ul className="dropdown-menu dropdown-user">
                    <li><a href="/downloads/utils.zip"><i className="fa fa-download" />&nbsp;Utilities</a></li>
                    <li><a href="/downloads/UTM-N1_User_Manual.zip"><i className="fa fa-book" />&nbsp;User Manual</a></li>
                </ul>
            );
    	} else {
            return (
                <ul className="dropdown-menu dropdown-user">
                    [Not logged in]
                </ul>
            );
    	}
    }
});

module.exports = Downloads;