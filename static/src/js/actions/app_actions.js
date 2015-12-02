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
    }
}

module.exports = AppActions;