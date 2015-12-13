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
var Moment = require('moment');

var Query = React.createClass({
    formatDate: function(event) {
        if (this.isMounted()) {    
	        document.getElementById("startDateTime").value = Moment(Date.parse(document.getElementById("startDateTime").value)).format("YYYY-MM-DDTHH:mm");
        }
    },

    render: function() {
        return (   
            <div className="row centered-form"><br /><br /><br />
                <div className="col-xs-12 col-sm-8 col-md-4 col-sm-offset-2 col-md-offset-4">
                    <div className="panel panel-default">
                        <div className="panel-heading">
                            <h4 className="panel-title text-left">Database Message Query</h4>
                        </div>
                        <div className="panel-body">
                            <form role="form" onSubmit={this.formatDate} action="/query" method="post">
                                <div className="form-group">
                                    <input type="text" name="uuid" id="uuid" className="form-control input-sm" placeholder="UUID (e.g. b7afe031-9c1d-46b8-bf09-5dcb520003b4)" autofocus />
                               </div>
                                 <div className="form-group">
                                    <input type="datetime-local" name="startDateTime" id="startDateTime" defaultValue="2015-12-01T00:00:00" className="form-control input-sm" placeholder="Start Date/Time (in UTC)"  />
                                </div>
                                <div className="form-group">
                                    <input type="number" name="duration" id="duration" className="form-control input-sm" placeholder="Duration (in minutes), leave blank for all" />
                                </div>
                                <input  type="submit" value="Submit Query" className="btn btn-info" />
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
});

module.exports = Query;