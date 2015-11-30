var React = require('react');

var ValueReporting = React.createClass({
    render:function(){
        return (
            <div>
        	    <input type="number" className="form-control bfh-number" min={1} max={10} defaultValue={1} style={{width: 145}} />
            </div>
        );
    }
});

module.exports = ValueReporting;