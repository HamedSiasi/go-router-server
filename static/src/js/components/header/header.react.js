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
"use strict";

var React = require('react');
var Link = require('react-router-component').Link;

var Header = React.createClass({
    render: function() {
        return (
            <nav className="navbar navbar-default navbar-static-top" role="navigation" style={{marginBottom: 0}}>
                <div className="navbar-header">               
                    <Link href="/Display">
                        <img src="static/dist/assets/images/logo.png" alt="u-blox" style={{maxWidth: 130, padding: 5}} />
                    </Link>
                </div>
            {/* /.navbar-header */}
                <ul className="nav navbar-top-links navbar-right">
                    {/* /.dropdown User */}
                    <li className="dropdown User">
                        <a className="dropdown-toggle" data-toggle="dropdown">
                            <i className="fa fa-user fa-fw" />&nbsp;User&nbsp;<i className="fa fa-caret-down" />
                        </a>
                        <ul className="dropdown-menu dropdown-user">
                            <li className="disabled"><Link href="#/login"><i className="fa fa-sign-in fa-fw" />&nbsp;Login</Link></li>
                            <li className="disabled"><Link href="#/logout"><i className="fa fa-sign-out fa-fw" />&nbsp;Logout</Link></li>
                            <li className="divider" />
                            <li className="disabled"><Link href="#/register"><i className="fa fa-user fa-fw" />&nbsp;Add User</Link></li>
                            <li className="disabled"><a href="#"><i className="fa fa-user fa-fw" />&nbsp;User Profile</a></li>
                        </ul>
                    </li>
                    {/* /.dropdown Data */}
                    <li className="dropdown Data">
                        <a className="dropdown-toggle" data-toggle="dropdown">
                            <i className="fa fa-bar-chart-o" />&nbsp;Data&nbsp;<i className="fa fa-caret-down" />
                        </a>
                        <ul className="dropdown-menu dropdown-user">
                            <li><Link href="#/display"><i className="fa fa-tachometer" /> Dash Board</Link></li>
                            <li className="disabled"><Link href="#/query"><i className="fa fa-database" /> Query</Link></li>
                        </ul>
                    </li>
                    {/* /.dropdown Downloads */}
                    <li className="dropdown Downloads">
                        <a className="dropdown-toggle" data-toggle="dropdown">
                            <i className="fa fa-download" />&nbsp;Downloads&nbsp;<i className="fa fa-caret-down" />
                        </a>
                        <ul className="dropdown-menu dropdown-user">                        
                            <li><a href="/downloads/utils.zip"><i className="fa fa-download" /> Download Utilities</a></li>
                            {/*
                            <li className="disabled">
                                <a href="/downloads/UTM-N1_User_Manual.pdf"><i className="fa fa-book" /> Download User Manual</a>
                            </li>
                            */}
                        </ul>
                    </li>
                </ul>
                {/* /.navbar-top-links */}
            </nav>
        );
    }
});

module.exports = Header;