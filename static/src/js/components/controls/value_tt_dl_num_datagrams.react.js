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

var ValueTTDlNumDatagrams = React.createClass({
    render:function(){
        return (
            <div>
                <input className="form-control" type="number" min={1} max={10000} defaultValue={100} step={1} style={{width: 80}} />
            </div>
        );
    }
});

module.exports = ValueTTDlNumDatagrams;