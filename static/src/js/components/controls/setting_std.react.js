var React = require('react');
var SetHeartbeat = require('./set_heartbeat.react');
var ValueHeartbeat = require('./value_heartbeat.react');
var SetReporting = require('./set_reporting.react');
var ValueReporting = require('./value_reporting.react');
var GetIntervals = require('./get_intervals.react');
var GetTime = require('./get_time.react');
var GetPing = require('./get_ping.react');

var SettingStd = React.createClass({
    render: function() {
        return (
            <div >
                <table >
	                <thead>
			            <tr >
			                <th colSpan={3} style={{textAlign: 'center'}}>Device Settings</th>
			            </tr>
			        </thead>
                    <tr style={{height: 50}}>
                        <td style={{width: 170}}> <ValueHeartbeat /></td>
                        <td style={{width: 300}}> <SetHeartbeat /></td> 
                        <td></td> 
                    </tr>
                    <tr >
                        <td><ValueReporting /></td>
                        <td><SetReporting /></td> 
                        <td></td> 
                    </tr>
                    <tr >
                        <td colSpan={3}><GetIntervals /> <GetTime /> <GetPing /></td>
                    </tr>
                </table>
            </div>
        );
    }
});

module.exports = SettingStd;