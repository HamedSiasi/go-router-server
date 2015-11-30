var React = require('react');

var ValueTTDlNumDatagrams = React.createClass({
    render:function(){
        return (
            <div>
                <input className="form-control" type="number" min={1} max={10000} defaultValue={100} step={1} style={{width: 80}} />
            </div>
        );
    }
});

module.exports = ValueTTDlNumDatagrams;