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

var React = require('react');

var MakeLiUuidList = React.createClass({
    render: function() {
        var list = [];
        if ((this.props.Items != null) && (this.props.Items.length > 0)) {            
            this.props.Items.forEach(function(item, i) {
                list.push(
                    <li key={i}><a href="#">&nbsp;{item}&nbsp;</a></li>
                );
        });
        } else {
            list.push(
                <li key={Date.now()}><a href="#">&nbsp;Empty&nbsp;</a></li>  // Date.now{} as a random key
            );
        }

        return (
            <div >
                {list}
            </div>
        );
    }
});

Object.size = function(obj) {
    var size = 0, key;
    for (key in obj) {
        if (obj.hasOwnProperty(key)) {
            size++;
        }
    }
    return size;
};

module.exports = MakeLiUuidList, Object.size;