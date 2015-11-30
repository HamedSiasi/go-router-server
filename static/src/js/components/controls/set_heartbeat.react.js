var React = require('react');

var SetHeartBeat = React.createClass({
	  render: function() {
	    return (
	        <div className="btn-group">
	            <button type="button" className="btn btn-info">Set Heartbeat (seconds)</button>
	            <button type="button" className="btn btn-info dropdown-toggle" data-toggle="dropdown">
	                <span className="caret" />
	            </button>
	            <ul className="dropdown-menu" role="menu">
	                <li><a href="#">blah</a></li>
	                <li><a href="#">blah1</a></li>
	            </ul>
	        </div>
	    );
	  }
	});

module.exports = SetHeartBeat;