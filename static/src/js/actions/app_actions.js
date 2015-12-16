/**
 * Copyright (C) u-blox Melbourn Ltd
 * u-blox Melbourn Ltd, Melbourn, UK
 * 
 * All rights reserved.
 *
 * This source file is the sole property of u-blox Melbourn Ltd.
 * Reproduction or utilization of this source in whole or part is
 * forbidden without the written consent of u-blox Melbourn Ltd.
 */

var AppConstants = require('../constants/app_store_types');
var AppDispatcher = require('../dispatchers/app_dispatcher');

var AppActions = {
    addUser: function(company, firstName, lastName, email, password) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.ADD_USER,
            company: company,
            firstName: firstName,
            lastName: lastName,
            email: email,
            password: password
        })
    },
    setIsLoggedIn: function(isLoggedIn) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.STORE_IS_LOGGED_IN,
            isLoggedIn: isLoggedIn
        })
    },
    setUuidChecked: function(uuid) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.STORE_SET_UUID_CHECKED,
            uuid: uuid
        })
    },
    setUuidUnchecked: function(uuid) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.STORE_SET_UUID_UNCHECKED,
            uuid: uuid
        })
    },
    setHeartbeatSeconds: function(hearbeatSeconds) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.STORE_HEARTBEAT_SECONDS,
            heartbeatSeconds: hearbeatSeconds
        })
    },
    setHeartbeatSnapToRtc: function(hearbeatSnapToRtc) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.STORE_HEARTBEAT_SNAP_TO_RTC,
            heartbeatSnapToRtc: hearbeatSnapToRtc
        })
    },
    setReportingInterval: function(reportingInterval) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.STORE_REPORTING_INTERVAL,
            reportingInterval: reportingInterval
        })
    },
    setTtNumUlDatagrams: function(value) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.STORE_TT_NUM_UL_DATAGRAMS,
            numUlDatagrams: value
        })
    },
    setTtLenUlDatagram: function(value) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.STORE_TT_LEN_UL_DATAGRAM,
            lenUlDatagram: value
        })
    },
    setTtNumDlDatagrams: function(value) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.STORE_TT_NUM_DL_DATAGRAMS,
            numDlDatagrams: value
        })
    },
    setTtLenDlDatagram: function(value) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.STORE_TT_LEN_DL_DATAGRAM,
            lenDlDatagram: value
         })
    },
    setTtTimeoutSeconds: function(value) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.STORE_TT_TIMEOUT_SECONDS,
            timeoutSeconds: value
         })
    },
    setTtDlIntervalSeconds: function(value) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.STORE_TT_DL_INTERVAL_SECONDS,
            dlIntervalSeconds: value
         })
    },
    setTtNoReportsDuringTest: function(value) {
        AppDispatcher.handleViewAction({
            actionType: AppConstants.STORE_TT_NO_REPORTS_DURING_TEST,
            noReportsDuringTest: value
        })
    }

}

module.exports = AppActions;