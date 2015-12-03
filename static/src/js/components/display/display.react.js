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
var Settings = require('../panels/settings.react')
var DisplayRow = require('./display_row.react')
var Summary = require('../panels/summary.react')
var Link = require('react-router-component').Link

var Display = React.createClass({
    getInitialState: function(){   
        return {
        	data:[]
       	};
    },

    componentDidMount: function() {
        pollState(function(data) {
            // fixup missing state properties to avoid muliple levels of missing attribute tests
            [
                "SummaryData",
                "DeviceData"
            ].map(function(property) {
                if (!data[property]) {
                    data[property] = {};
                }
            });
            this.setState({data: data});
        }.bind(this), 10000);
    },

    render:function(){
        return (
            <div>
                <br />
                <Settings DeviceData = {this.state.data["DeviceData"]} />  
                <Summary SummaryData = {this.state.data["SummaryData"]} />  
                <DisplayRow DeviceData = {this.state.data["DeviceData"]} />  
            </div>
        );
    }
});

function pollState(updateState) {
    function pollLoop() {
        var x = new XMLHttpRequest();
        x.onreadystatechange = function() {
            if (x.readyState == 4) {
                if (x.status == 200) {
                    var data = JSON.parse(x.responseText);
                    updateState(data);
                }
                window.setTimeout(pollLoop, 1000);
            }
        };
        x.open("GET", "frontPageData", true);
        x.send();
    }
    pollLoop();
}

module.exports = Display;