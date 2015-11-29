var React = require('react');

var Measurements = React.createClass({
  render:function(){
    return (
    	 <div className="panel panel-default">
            <div className="panel-heading">
    		   <div className="panel-body">
                <div className="dataTable_wrapper">
                  <table className="table table-striped table-bordered table-hover" id="dataTables-example">
                    <thead>
                      <tr className="info">
                        <th> <input type="checkbox" style={{width: 15}} /> All</th>
                        <th>Name/Uuid</th>
                       <th>Upstream</th>
                        <th>Downsteam</th>
                        <th>RSRP </th>
                        <th>
                          Battery - <i className="fa fa-floppy-o" />
                        </th>
                      </tr>
                    </thead>
                    <tbody style={{fontSize: 12}}>
                      <tr className="even gradeA">
                        <td style={{width: 15}}>
                          <input type="checkbox" style={{width: 15}} /><br />
                          <img src="static/Images/green.png" alt="logo" style={{maxWidth: 12}} />
                        </td>
                        <td>
                          <ul>
                            <li><b>Uuid:</b> {/*this.state["LatestDisplayRow"]["Uuid"]*/}</li>
                            <li><b>Mode:</b> {/*this.state["LatestDisplayRow"]["Mode"]*/}</li>
                            <li><b>Name:</b> {/*this.state["LatestDisplayRow"]["UnitName"]*/}</li>
                          </ul>   
                        </td>
                        <td>
                          <ul>
                            <li><b>Total Msg:</b> {/*this.state["LatestDisplayRow"]["TotalMsgs"]*/}</li>
                            <li><b>Total Bytes:</b> {/*this.state["LatestDataVolume"]["UplinkBytes"]*/}</li>
                            <li><b>Last Msg RX:</b> {/*this.state["LatestDisplayRow"]["UlastMsgReceived"]*/}</li>
                          </ul> 
                        </td>
                        <td className="center">
                          <ul>
                            <li><b>Total Msg:</b> {/*this.state["LatestDisplayRow"]["DTotalMsgs"] ? "" : "--"*/}</li>
                            <li><b>Total Bytes:</b> {/*this.state["LatestDisplayRow"]["DTotalBytes"] ? "" : "--"*/}</li>
                            <li><b>Last Msg RX:</b> {/*this.state["LatestDisplayRow"]["DlastMsgReceived"] ? "" : "--"*/}</li>
                          </ul>
                        </td>
                        <td className="center">{/*this.state["LatestDisplayRow"]["BatteryLevel"] ? "" : "0"*/}</td>
                        <td className="center">
                          {/*this.state["LatestDisplayRow"]["BatteryLevel"] ? "" : "0"*/} -
                          {/*this.state["LatestDisplayRow"]["DiskSpaceLeft"] ? "" : " 0"*/}
                        </td> 
                      </tr>
                    <tr className="even gradeC">
                        <td style={{width: 15}}>
                          <a tabIndex={-1} href="#/standardtwo"> <b className="fa fa-cogs" /></a><br />
                          <input type="checkbox" style={{width: 15}} /><br />
                          <img src="ehlo.jpg" alt="logo" style={{maxWidth: 12}} />
                        </td>
                        <td style={{width: 205}}>
                          <ul>
                            <li><b>Uuid:</b> {/*this.state["LatestDisplayRow"]["Uuid"]*/}</li>
                            <li><b>Mode:</b> {/*this.state["LatestDisplayRow"]["Mode"]*/}</li>
                            <li><b>Name:</b> {/*this.state["LatestDisplayRow"]["UnitName"]*/}</li>
                          </ul>   
                        </td>
                        <td style={{width: 75}}>
                          <ul>
                            <li><b>Total Msg:</b> {/*this.state["LatestDisplayRow"]["TotalMsgs"]*/}</li>
                            <li><b>Total Bytes:</b> {/*this.state["LatestDataVolume"]["UplinkBytes"]*/}</li>
                            <li><b>Last Msg RX:</b> {/*this.state["LatestDisplayRow"]["UlastMsgReceived"]*/}</li>
                          </ul> 
                        </td>
                        <td className="center" style={{width: 75}}>
                          <ul>
                            <li><b>Total Msg:</b> {/*this.state["LatestDisplayRow"]["DTotalMsgs"] ? "" : "--"*/}</li>
                            <li><b>Total Bytes:</b> {/*this.state["LatestDisplayRow"]["DTotalBytes"] ? "" : "--"*/}</li>
                            <li><b>Last Msg RX:</b> {/*this.state["LatestDisplayRow"]["DlastMsgReceived"] ? "" : "--"*/}</li>
                          </ul>
                        </td>
                        <td className="center" style={{width: 25}}>{/*this.state["LatestDisplayRow"]["BatteryLevel"] ? "" : "0"*/}</td>
                        <td className="center" style={{width: 25}}>
                          {/*this.state["LatestDisplayRow"]["BatteryLevel"] ? "" : "0"*/ } -
                          {/*this.state["LatestDisplayRow"]["DiskSpaceLeft"] ? "" : " 0"*/ }
                        </td> 
                      </tr>
                    </tbody> 
                  </table>
                </div>
              </div>
            </div>
      </div>
     
    );
  }
});

module.exports = Measurements;