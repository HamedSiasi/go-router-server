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
 * For hints as to how this works, see here:
 * 
 * https://facebook.github.io/flux/docs/todo-list.html#content
 */
"use strict";

var AppDispatcher = require('../dispatchers/app_dispatcher');
var AppConstants = require('../constants/app_store_types');
var assign = require('react/lib/Object.assign');
var EventEmitter = require('events').EventEmitter;

var CHANGE_EVENT = 'change';

var _uuidsCheckedList = {};
var _heartbeatSeconds;
var _heartbeatSnapToRtc;
var _reportingInterval;
var _ttParameters = {};
var _isLoggedIn = false;

//------------------------------------------------------------
// Private functions
//------------------------------------------------------------

//Store the isLoggedIn value
function _storeIsLoggedIn(isLoggedIn) {
	_isLoggedIn = isLoggedIn;
}

// Add a UUID to the list as being checked
function _setUuidChecked(uuid) {
    _uuidsCheckedList[uuid] = uuid;
}

// Remove a UUID from the list of checked ones
function _setUuidUnchecked(uuid) {
    if (_uuidsCheckedList) {
   	    delete _uuidsCheckedList[uuid];
    }
}

// Store the heartbeat value
function _storeHeartbeatSeconds(heartbeatSeconds) {
	_heartbeatSeconds = heartbeatSeconds;
}

//Store the heartbeat snap-to-rtc value
function _storeHeartbeatSnapToRtc(heartbeatSnapToRtc) {
	_heartbeatSnapToRtc = heartbeatSnapToRtc;
}

//Store the reporting interval value
function _storeReportingInterval(reportingInterval) {
	_reportingInterval = reportingInterval;
}

//Store the traffic test parameters
function _storeTtNumUlDatagrams(numUlDatagrams) {
	_ttParameters["numUlDatagrams"] = numUlDatagrams;
}
function _storeTtLenUlDatagram(lenUlDatagram) {
	_ttParameters["lenUlDatagram"] = lenUlDatagram;
}
function _storeTtNumDlDatagrams(numDlDatagrams) {
	_ttParameters["numDlDatagrams"] = numDlDatagrams;
}
function _storeTtLenDlDatagram(lenDlDatagram) {
	_ttParameters["lenDlDatagram"] = lenDlDatagram;
}
function _storeTtTimeoutSeconds(timeoutSeconds) {
	_ttParameters["timeoutSeconds"] = timeoutSeconds;
}
function _storeTtNoReportsDuringTest(noReportsDuringTest) {
	_ttParameters["noReportsDuringTest"] = noReportsDuringTest;
}
function _storeTtDlIntervalSeconds(dlIntervalSeconds) {
	_ttParameters["dlIntervalSeconds"] = dlIntervalSeconds;
}

// TODO Rob to understand this later 
function addUser(company, firstName, lastName, email, password) {
 obj = {};
 obj.company = company;
 obj.firstName = firstName;
 obj.lastName = lastName;
 obj.email = email;
 obj.password = password;

 $.ajax({
     url: 'http://localhost:3000/register',
     dataType: 'json',
     method: 'put',
     async: false,
     data: obj,
     success: function(data) {
         return  
     },
     error: function(xhr, status, err) {
         console.error('/', status, err.toString());
     }
 });
}

//------------------------------------------------------------
// Public functions
//------------------------------------------------------------

// Handle indexof not being defined, see:
// http://www.tutorialspoint.com/javascript/array_indexof.htm
if (!Array.prototype.indexOf)
{
    Array.prototype.indexOf = function(elt /*, from*/)
    {
        var len = this.length;
      
        var from = Number(arguments[1]) || 0;
        from = (from < 0)
        ? Math.ceil(from)
        : Math.floor(from);
      
        if (from < 0)
            from += len;
      
        for (; from < len; from++)
        {
            if (from in this && this[from] === elt)
                return from;
        }
        return -1;
    };
}

// The store itself
var AppStore = assign(EventEmitter.prototype, {
    // This is vital: I create lots of listeners, so need to remove the limit of 11
	_maxListeners: 0,
	
	emitChange: function() {
        this.emit(CHANGE_EVENT)
    },

    addChangeListener: function(callback) {
        this.on(CHANGE_EVENT, callback)
    },

    removeChangeListener: function(callback) {
        this.removeListener(CHANGE_EVENT, callback)
    },

    getIsLoggedIn: function() {   	
    	return _isLoggedIn;
    },

    getAllUuidsChecked: function () {
    	return _uuidsCheckedList;
    },
    
    isUuidChecked: function(uuid) {   	
    	return (_uuidsCheckedList[uuid] != null);
    },

    getHeartbeatSeconds: function () {
    	return _heartbeatSeconds;
    },
    
    getHeartbeatSnapToRtc: function () {
    	return _heartbeatSnapToRtc;
    },
    
    getReportingInterval: function () {
    	return _reportingInterval;
    },
    
    getTtParameters: function () {
    	return _ttParameters;
    },
    
    dispatcherIndex: AppDispatcher.register(function(payload) {
        var action = payload.action; // this is our action from handleViewAction
    
        switch(action.actionType) {
            case AppConstants.STORE_IS_LOGGED_IN:
            	_storeIsLoggedIn(action.isLoggedIn);
            break;
            case AppConstants.STORE_SET_UUID_CHECKED:
                _setUuidChecked(action.uuid);
            break;
            case AppConstants.STORE_SET_UUID_UNCHECKED:
                _setUuidUnchecked(action.uuid);
            break;
            case AppConstants.STORE_IS_UUID_CHECKED:
                _isUuidUnchecked(action.uuid);
            break;
            case AppConstants.STORE_HEARTBEAT_SECONDS:
                _storeHeartbeatSeconds(action.heartbeatSeconds);
            break;
            case AppConstants.STORE_HEARTBEAT_SNAP_TO_RTC:
                _storeHeartbeatSnapToRtc(action.heartbeatSnapToRtc);
            break;
            case AppConstants.STORE_REPORTING_INTERVAL:
                _storeReportingInterval(action.reportingInterval);
            break;
            case AppConstants.STORE_TT_NUM_UL_DATAGRAMS:
                _storeTtNumUlDatagrams(action.numUlDatagrams);
            break;
            case AppConstants.STORE_TT_LEN_UL_DATAGRAM:
                _storeTtLenUlDatagram(action.lenUlDatagram);
            break;
            case AppConstants.STORE_TT_NUM_DL_DATAGRAMS:
                _storeTtNumDlDatagrams(action.numDlDatagrams);
            break;
            case AppConstants.STORE_TT_LEN_DL_DATAGRAM:
                _storeTtLenDlDatagram(action.lenDlDatagram);
            break;
            case AppConstants.STORE_TT_TIMEOUT_SECONDS:
                _storeTtTimeoutSeconds(action.timeoutSeconds);
            break;
            case AppConstants.STORE_TT_NO_REPORTS_DURING_TEST:
                _storeTtNoReportsDuringTest(action.noReportsDuringTest);
            break;
            case AppConstants.STORE_TT_DL_INTERVAL_SECONDS:
                _storeTtDlIntervalSeconds(action.dlIntervalSeconds);
            break;
             // Insert more things here
            default:
                console.log("Unknown action type sent into app store");
            break;
        }

        AppStore.emitChange();

        return true;
    })
})

module.exports = AppStore;
