var React = require('react');

var Apply = React.createClass({
  render:function(){
    return (
      <button type="button"   className="btn btn-info" style={{width: 100, height: 30, float: 'left', marginTop: 10}} >
      Apply
      </button>
    );
  }
});

module.exports = Apply;