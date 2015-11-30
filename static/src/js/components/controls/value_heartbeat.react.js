var React = require('react');

var ValueHeartbeat = React.createClass({
    render:function(){
        return (
            <div>
                <input className="form-control" type="number" min={1} max={3599} defaultValue={900} step={1} style={{width: 145}} />
            </div>
        );
    }
});

module.exports = ValueHeartbeat;