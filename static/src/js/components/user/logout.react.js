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
var AppActions = require('../../actions/app_actions.js');

var Logout = React.createClass({
    render:function(){
        AppActions.setIsLoggedIn(false);
        return (   
            <div>
                <br/>
                You have been logged out. <Link href="#/login">Click here to log back in again</Link>.
            </div>
        );
    }
});

module.exports = Logout;