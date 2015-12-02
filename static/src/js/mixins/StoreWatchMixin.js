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

var React = require('react');
var AppStore = require('../stores/app_store');

var StoreWatchMixin = function(cb) {
    return {
        getInitialState:function() {
            return cb(this)
        },
        componentWillMount:function() {
            AppStore.addChangeListener(this._onChange)
        },
        componentWillUnmount:function() {
            AppStore.removeChangeListener(this._onChange)
        },
        _onChange: function() {
            this.setState(cb(this))
        }
    }
}

module.exports = StoreWatchMixin;
