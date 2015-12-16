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

var pass;
var repass;

var Register = React.createClass({
    componentWillMount: function() {
        if (!AppStore.getIsLoggedIn()) {
            window.location = "/";
        }
    },

    handlePasswordChange: function(e) {
        pass = e.target.value;
    },
        
    handleRePasswordChange: function(e) {
        repass = e.target.value;
        if (pass !== repass) {
            alert('Password do not much, Please try again');
            e.target.value = '';
        }
    },

    render:function(){
    return (
        <div className="row centered-form"><br /><br /><br />
            <div className="col-xs-12 col-sm-8 col-md-4 col-sm-offset-2 col-md-offset-4">
                <div className="panel panel-default">
                    <div className="panel-heading">
                        <h4 className="panel-title text-left">Create User</h4>
                    </div>
                    <div className="panel-body">
                        <form role="form" action="/register" method="post" data-toggle="validator" >
                            <div className="form-group">
                                <input type="text" name="user_firstName" id="user_firstName" className="form-control input-sm" placeholder="First Name" required />
                                    <div className="help-block with-errors"></div>
                            </div>
                            <div className="form-group">
                                <input type="text" name="user_lastName" id="user_lastName" className="form-control input-sm" placeholder="Last Name" required />
                            </div>
                            <div className="form-group">
                                <input type="email"  name="email" id="inputEmail" data-error="Bruh, that email address is invalid" className="form-control input-sm" placeholder="Email Address" required/>
                                <div className="help-block with-errors"></div>     
                            </div>
                            <div className="form-group">
                            <input type="text" name="company_name" id="company_name" className="form-control input-sm" placeholder="Company Name" required />
                                <div className="help-block with-errors"></div>
                            </div>
                            <div className="form-group">
                                <input type="password" data-minlength="6"  name="password" onMouseLeave ={this.handlePasswordChange} id="inputPassword" className="form-control input-sm" placeholder="Password" required/>
                                    <div className="help-block with-errors"></div>
                            </div>
                            <div className="form-group">
                                <input type="password" name="password_confirmation" onBlur ={this.handleRePasswordChange} data-match="#inputPassword" data-match-error="Whoops, these don't match" id="password_confirmation" className="form-control input-sm" placeholder="Confirm Password" required/>
                                    <div className="help-block with-errors"></div>
                            </div>
                                <input type="Submit" value="Submit" className="btn btn-info" />
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
});

module.exports = Register;