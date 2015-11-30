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
var Template = require('./app-template.js');
var Router = require('react-router-component');
var Display = require('./display/display.react');
var Login = require('./user/login.react');
var Register = require('./user/register.react');
var Index = require('./index');

var Locations = Router.Locations;
var Location  = Router.Location;

var App = React.createClass({
    render:function(){
        return (
            <Template>
                <Locations>
                    <Location path="/" handler={Index} />
                    <Location path="#/display" handler={Display} />
                    <Location path="#/login" handler={Login} />
                    <Location path="#/register" handler={Register} />
                </Locations>
            </Template>
        );
    }
});

module.exports = App;