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

var ValueUuidSelected = React.createClass({
    getInitialState: function() {
        return {value: this.props.Uuid, checked: this.props.Checked};
    },

    handleChange: function(newValue) {
    	this.setState({checked: newValue.target.checked});
    	if (this.props.CallbackParent) {
            this.props.CallbackParent(this.state);
    	}
    },

	render:function(){
        return (
            <input type="checkbox" value={this.state.value} checked={this.state.checked} onChange={this.handleChange} style={{width: 15}} />
        );
    }
});

module.exports = ValueUuidSelected;