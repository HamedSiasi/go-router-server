var React = require('react');

var SetTTStop = React.createClass({
	  render: function() {
	    return (
	        <div className="btn-group">
	            <button type="button" className="btn btn-info" style={{width: 170, height: 30}}>Stop Test</button>
	            <button type="button" className="btn btn-info dropdown-toggle" data-toggle="dropdown" style={{height: 30}} >
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

module.exports = SetTTStop;