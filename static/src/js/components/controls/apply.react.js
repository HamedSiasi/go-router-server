var React = require('react');

var GetIntervals = React.createClass({
    render:function(){
        return (
            <button type="button"   className="btn btn-info" style={{width: 100, height: 30, float: 'left', marginTop: 10}} >
                Get Intervals
            </button>
        );
    }
});

module.exports = GetIntervals;