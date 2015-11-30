var React = require('react');

var SetReporting = React.createClass({
    render:function(){
        return (
            <div className="btn-group">
                <button type="button" className="btn btn-info">Set Reporting (heartbeats)</button>
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

module.exports = SetReporting;