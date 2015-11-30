var React = require('react');

var GetTime = React.createClass({
    render:function(){
        return (
  	        <div className="btn-group">
  	            <button type="button" className="btn btn-info" style={{width: 90, height: 30, marginLeft: 10, marginTop: 10}} >
  	                Get Time
  	            </button>
  	            <button type="button" className="btn btn-info dropdown-toggle" data-toggle="dropdown" style={{height: 30, marginTop: 10}} >
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

module.exports = GetTime;