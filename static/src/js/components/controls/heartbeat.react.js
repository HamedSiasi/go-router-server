var React = require('react');

var HeartBeat = React.createClass({
  render:function(){
    return (
         <select className="form-control" style={{width: 180, height: 30, marginTop: 10, display: 'inline'}}>
              <option>Heart Beat</option>
              <option>Fast (2min)</option>
              <option>Medium (5min)t</option>
              <option>Slow(30)</option>
           </select>
    );
  }
});

module.exports = HeartBeat;