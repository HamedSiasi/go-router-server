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

var AppDispatcher = require('../dispatchers/app-dispatcher');
var AppConstants = require('../constants/app-constants');
var assign = require('react/lib/Object.assign');
var EventEmitter = require('events').EventEmitter;

var CHANGE_EVENT = 'change';

function _setUuidChecked(uuid){

}

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

var AppStore = assign(EventEmitter.prototype, {
            emitChange: function() {
            this.emit(CHANGE_EVENT)
        },

        addChangeListener: function(callback) {
           this.on(CHANGE_EVENT, callback)
        },

        removeChangeListener: function(callback){
            this.removeListener(CHANGE_EVENT, callback)
        },

        getUtmsData: function(){
            return states
        },

        dispatcherIndex: AppDispatcher.register(function(payload){
        var action = payload.action; // this is our action from handleViewAction
    
        switch(action.actionType) {
            case AppConstants.STORE_SET_UUID_CHECKED:
                _setUuidChecked(payload.action.index);
            break;
            // Insert more things here
        }

        AppStore.emitChange();

        return true;
    })
})

module.exports = AppStore;
