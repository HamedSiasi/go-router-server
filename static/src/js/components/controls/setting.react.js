var React = require('react');
var Apply = require('./apply.react');
var Reboot = require('./reboot.react');
var Reporting = require('./reporting.react');
var HeartBeat = require('./heartbeat.react');

var Setting = React.createClass({
  render: function() {
    return (
<div >          
   <table >
      <tr>
        <td> <Reboot />Reboot</td>
        <td></td> 
      </tr>
      <tr>
        <td>  <Reporting /></td>
        <td></td> 
      </tr>
      <tr>
        <td>  <HeartBeat /></td>
        <td></td> 
      </tr>
      <tr>
        <td><Apply /> </td>
        <td></td> 
      </tr>
  </table>
</div>

    );
  }
});

module.exports = Setting;