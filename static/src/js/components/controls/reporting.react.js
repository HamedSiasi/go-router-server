var React = require('react');

var Reporting = React.createClass({
  render:function(){
    return (
    	    <select className="form-control" style={{width: 180, height: 30, marginTop: 10, display: 'inline'}}>
              <option>Reporting Interval</option>
              <option>Fast (2min)</option>
              <option>Medium (5min)t</option>
              <option>Slow(30)</option>
            </select>     
    );
  }
});

module.exports = Reporting;