var React = require('react');
var SettingStd = require('../controls/setting_std.react');
var SettingTT = require('../controls/setting_tt.react');

var Settings = React.createClass({
    render:function(){
        return (
            <div className="row" >
		        <div className="col-lg-5">
		            <div className="panel panel-info" style={{height: 180, marginTop: 10}}>
			            <div className="panel-body">
                            <SettingStd />
			            </div>
		            </div>
	            </div>
		        <div className="col-lg-7">
		            <div className="panel panel-info" style={{height: 180, marginTop: 10}}>
			            <div className="panel-body">
                            <SettingTT />
			            </div>
		            </div>
                </div>
	        </div>
	    )
    }
});

module.exports = Settings;