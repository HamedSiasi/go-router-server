var React = require('react');

var ValueTTUlLenDatagram = React.createClass({
    render:function(){
        return (
            <div>
                <input className="form-control" type="number" min={1} max={100} defaultValue={100} step={1} style={{width: 80}} />
            </div>
        );
    }
});

module.exports = ValueTTUlLenDatagram;