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

var AppDispatcher = require('../dispatchers/app_dispatcher');
var AppConstants = require('../constants/app_store_types');
var assign = require('react/lib/Object.assign');
var EventEmitter = require('events').EventEmitter;

var CHANGE_EVENT = 'change';

//var _uuidsCheckedList = {"61d25940-5307-11e5-80a4-c549fb2d6313": "61d25940-5307-11e5-80a4-c549fb2d6313"}
var _uuidsCheckedList = {};

//------------------------------------------------------------
// Private functions
//------------------------------------------------------------

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

// TODO
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
    emitChange: function() {
        this.emit(CHANGE_EVENT)
    },

    addChangeListener: function(callback) {
        this.on(CHANGE_EVENT, callback)
    },

    removeChangeListener: function(callback) {
        this.removeListener(CHANGE_EVENT, callback)
    },

    getAllUuidsChecked: function () {
    	return _uuidsCheckedList;
    },
    
    isUuidChecked: function(uuid) {   	
    	return (_uuidsCheckedList[uuid] != null);
    },

    dispatcherIndex: AppDispatcher.register(function(payload) {
        var action = payload.action; // this is our action from handleViewAction
    
        switch(action.actionType) {
            case AppConstants.STORE_SET_UUID_CHECKED:
                _setUuidChecked(action.uuid);
            break;
            case AppConstants.STORE_SET_UUID_UNCHECKED:
                _setUuidUnchecked(action.uuid);
            break;
            case AppConstants.STORE_IS_UUID_CHECKED:
                _isUuidUnchecked(action.uuid);
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
