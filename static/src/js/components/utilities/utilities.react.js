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

var MakeLiUuidList = React.createClass({
    render: function() {
        var list = [];
        if ((this.props.Items != null) && (this.props.Items.length > 0)) {        	
            this.props.Items.forEach(function(item, i) {
            	list.push(
                    <li><a href="#">&nbsp;{item}&nbsp;</a></li>
                );
        });
        } else {
        	list.push(
	            <li><a href="#">&nbsp;Empty&nbsp;</a></li>
	        );
        }

        return (
            <div >
                {list}
            </div>
        );
    }
});

module.exports = MakeLiUuidList;