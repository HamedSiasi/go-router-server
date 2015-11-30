var React = require('react');

var ValueTTUlNumDatagrams = React.createClass({
    render:function(){
        return (
        	<div>
	            <input className="form-control" type="number" name="UlNum" min={1} max={10000} defaultValue={100} step={1} style={{width: 80}} />
            </div>
        );
    }
});

module.exports = ValueTTUlNumDatagrams;