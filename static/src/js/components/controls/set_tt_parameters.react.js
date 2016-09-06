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
var MakeLiNameList = require('../utilities/make_li_name_list.react')
var AppStore = require('../../stores/app_store.js');

var SetTTParameters = React.createClass({

	sendAllMsg: function() {
        var uuidList = AppStore.getAllUuidsChecked();

        for (var name in this.props.UuidMap) {
            if ((this.props.UuidMap.hasOwnProperty(name)) && (uuidList [this.props.UuidMap[name]])) {
                var request = new XMLHttpRequest();
                var ttParameters = AppStore.getTtParameters();
                var postData = "{\"device_uuid\": \"" + this.props.UuidMap[name] +
                		"\", \"type\": \"SEND_TRAFFIC_TEST_MODE_PARAMETERS_SERVER_SET\", \"body\": {\"num_ul_datagrams\": \"" + ttParameters["numUlDatagrams"] +
                        "\", \"len_ul_datagram\": \"" + ttParameters["lenUlDatagram"] +
                        "\", \"num_dl_datagrams\": \"" + ttParameters["numDlDatagrams"] +
                        "\", \"len_dl_datagram\": \"" + ttParameters["lenDlDatagram"] +
                        "\", \"timeout_seconds\": \"" + ttParameters["timeoutSeconds"] +
                        "\", \"no_reports_during_test\": \"" + ttParameters["noReportsDuringTest"] +
                        "\", \"dl_interval_seconds\": \"" +	ttParameters["dlIntervalSeconds"] + "\"} }";

                // What to do when the server response comes back
                // Note this is not the device response, just the server
                // HTTP confirm.
                request.onload = function () {
                    // TODO something here?
                }
                request.open("POST", "sendMsg", true);
                request.setRequestHeader("Content-Type", "application/json;charset=UTF-8");

                // Actually sends the request to the server.
                request.send(postData);
            }
        }
    },
    
	render: function() {
        return (
            <div className="btn-group">
                <button type="button" className="btn btn-info" onClick={this.sendAllMsg} style={{width: 170, marginTop: 10}} >Send Parameters</button>
                <button type="button" className="btn btn-info dropdown-toggle disabled" data-toggle="dropdown" style={{marginTop: 10}} >
                    <span className="caret" />
                </button>
                <ul className="dropdown-menu" role="menu">
                    <MakeLiNameList Items={this.props.Names} />
                </ul>
            </div>
        );
      }
    });

module.exports = SetTTParameters;