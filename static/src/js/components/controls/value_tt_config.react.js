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

var TtConfig = React.createClass({
    render:function() {
	    if ((this.props.UlLength > 0) || (this.props.DlLength > 0)) {
            //<i className="fa fa-cog" /> <i className="fa fa-arrow-up" /> N:{this.props.UlDatagrams}/L:{this.props.UlLength} <i className="fa fa-arrow-down" /> N:{this.props.DlDatagrams}/L:{this.props.DlLength}/G:{this.props.DlInterval} T:{this.props.Timeout}
            return (
                <div>
                    <i className="fa fa-cog" /> <i className="fa fa-arrow-up" /> N:<b>{this.props.UlDatagrams}</b>/L:<b>{this.props.UlLength}</b> <i className="fa fa-arrow-down" /> N:<b>{this.props.DlDatagrams}</b>/L:<b>{this.props.DlLength}</b>/G:<b>{this.props.DlInterval}</b>
                </div>
            );
        } else {
            return (
                <div>
                    <i className="fa fa-cog" />
        	    </div>
	    	);
	    }
    }
});

module.exports = TtConfig;