/**
 * Copyright (c) 2014, U-blox.
 * All rights reserved.
 */
var React = require('react');
var AppStore = require('../../stores/app-store.js');
var Settings = require('../panels/settings.react')
var DisplayRow = require('./displayRow.react')

var Link = require('react-router-component').Link

var arrData = [];

var currentUuidsObject = new Object();
var UuidsMap = new Map();
var totalMsg = 0;
var totalBytes = 0;

var Display = React.createClass({
    getInitialState: function(){   
        var data = {data:[
        ]}    
        return data;
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
	            <Settings />  
	            <DisplayRow DeviceData = { this.state.data["DeviceData"]} />  
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

Object.size = function(obj) {
    var size = 0, key;
    for (key in obj) {
        if (obj.hasOwnProperty(key)) size++;
    }
    return size;
};

module.exports = Display;