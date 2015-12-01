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
var MakeLiNameList = require('../utilities/utilities.react')

var SetTTParameters = React.createClass({
	  render: function() {
	    return (
	        <div className="btn-group">
	            <button type="button" className="btn btn-info" style={{width: 165, height: 30, marginTop: 10}}>Send Parameters</button>
	            <button type="button" className="btn btn-info dropdown-toggle" data-toggle="dropdown" style={{height: 30, marginTop: 10}} >
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