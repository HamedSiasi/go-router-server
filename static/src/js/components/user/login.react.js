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
var Display = require('../display/display.react');
var Query = require('../db/query.react');
var AppActions = require('../../actions/app_actions.js');
var AppStore = require('../../stores/app_store.js');

var Login = React.createClass({
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
    
	authenticate: function(event) {
	    event.preventDefault();
	    event.stopPropagation();
	    
        var request = new XMLHttpRequest();
        var postData = "{\"email\": \"" + document.getElementById("email").value + "\", \"password\": \"" + document.getElementById("password").value + "\"}";

        // What to do when the server response comes back
        request.onload = function () {
        	if (request.status == 302) {  // 302 = http.StatusFound
        	    AppActions.setIsLoggedIn(true);
        	} else {
        	    AppActions.setIsLoggedIn(false);
        	    window.location = "/";
        	}
        }
        
        request.open("POST", "login", true);
        request.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
        request.send(postData);
    },

    render:function(){
    	if (this.state.okToContinue) {
    		return (
   	            <div>
                    You are logged in, <Link href="#/display">click here for Dash Board</Link>.
	            </div>
    	    );
    	} else {
            return (   
                <div className="row centered-form"><br /><br /><br />
                    <div className="col-xs-12 col-sm-8 col-md-4 col-sm-offset-2 col-md-offset-4">
                        <div className="panel panel-default">
                            <div className="panel-heading">
                                <h4 className="panel-title text-left">Please Sign In</h4>
                            </div>
                            <div className="panel-body">
                                <form role="form" action="login" method="post" >
                                    <div className="form-group">
                                        <input type="email" name="email" id="email" className="form-control input-sm" placeholder="Email Address" required autofocus />
                                    </div>
                                    <div className="form-group">
                                        <input type="password" name="password" id="password" className="form-control input-sm" placeholder="User password"  required />
                                    </div>
                                    <input type="submit" value="Login" className="btn btn-info" onClick={this.authenticate} />
                                </form>
                            </div>
                        </div>
                    </div>
                </div>
            );
    	}
    }
});

module.exports = Login;