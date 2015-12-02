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
var AppActions = require('../../actions/app_actions.js');
var AppStore = require('../../stores/app_store.js');

var ValueUuidSelected = React.createClass({
    getInitialState: function() {
        return {value: this.props.Uuid, checked: AppStore.isUuidChecked(this.props.Uuid)};
    },

    componentDidMount: function() {
        AppStore.addChangeListener(this.onChange);
    },

    componentWillUnmount: function() {
        AppStore.removeChangeListener(this.onChange);
    },

    onChange: function() {
        this.setState({checked: AppStore.isUuidChecked(this.state.value)});
        //this.forceUpdate();
    },
    
    handleChange: function(newValue) {
    	if (newValue.target.checked == true) {
            AppActions.setUuidChecked(this.state.value);
        } else {
            AppActions.setUuidUnchecked(this.state.value);
        }
    },

    render:function(){
        return (
            <input type="checkbox" value={this.state.value} checked={this.state.checked} onChange={this.handleChange} style={{width: 15}} />
        );
    }
});

module.exports = ValueUuidSelected;