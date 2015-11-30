var React = require('react');

var ValueTTTimeout = React.createClass({
    render:function(){
        return (
            <input className="form-control" type="number" min={1} max={86400} defaultValue={6000} step={1} style={{width: 80, marginTop: 10}} />
        );
    }
});

module.exports = ValueTTTimeout;