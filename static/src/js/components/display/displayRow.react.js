/**
 * Copyright (c) 2014, U-blox.
 * All rights reserved.
 */
var React = require('react');
var Link = require('react-router-component').Link;

var DisplayRow = React.createClass({

    render: function() {

        var rows = [];

          if(this.props["DeviceData"] !== undefined){   

            this.props["DeviceData"].forEach(function(uuid, i) {

             
            rows.push(       
    
              <tr className="even gradeC" key={i}>
                        <td style={{width: 15}}>
                          <a tabIndex={-1} href="#/standardtwo"> <b className="fa fa-cogs" /></a><br />
                          <input type="checkbox" style={{width: 15}} /><br />
                          <img src="static/dist/assets/images/green.png" alt="logo" style={{maxWidth: 12}} />
                        </td>
                        <td>
                          <ul className="SmallPadding">
                            <li><b>Uuid:</b> {uuid["Uuid"]}</li>
                            <li><b>Mode:</b> {uuid["Mode"]}</li>
                            <li><b>Name:</b> {uuid["UnitName"]}</li>
                             <li><b>Reporting Interval:</b> {uuid["Reporting"]}</li>
                              <li><b>Heart Beat:</b> {uuid["Heartbeat"]}</li>
                          </ul>   
                        </td>
                        <td>
                          <ul className="SmallPadding">
                            <li><b>Total Msg:</b> {uuid["TotalUlMsgs"]}</li>
                            <li><b>Total Bytes:</b> {uuid["TotalUlBytes"]}</li>
                            <li><b>Last Msg RX:</b> {uuid["LastUlMsgTime"]}</li>
                          </ul> 
                        </td>
                        <td className="center">
                          <li><b>Total Msg:</b> {uuid["TotalDlMsgs"]}</li>
                          <li><b>Total Bytes:</b> {uuid["TotalDlBytes"]}</li>
                          <li><b>Last Msg RX:</b> {uuid["LastDlMsgTime"]}</li>
                        </td>
                        <td className="center">{uuid["Rsrp"]}</td>
                        <td className="center" style={{width: 105}}>
                          <i className="fa fa-floppy-o" /> {uuid["DiskSpaceLeft"]}<br />
                          <i className="fa fa-battery-full" /> {uuid["BatteryLevel"]}
                        </td> 
                      </tr>

                      );

                
         
   
       });
}
                return (
                        <div className="row" >
                          <div className="panel panel-default" >
                            <div className="_panel-heading" style={{width:'100%'}}>
                              <div className="panel-body">
                                <div className="dataTable_wrapper">
                                  <table className="table table-striped table-bordered table-hover" id="dataTables-example">
                                    <thead>
                                      <tr className="info">
                                        <th> <input type="checkbox" style={{width: 15}} /> All</th>
                                        <th>Name/Uuid</th>
                                        <th>Uplink</th>
                                        <th>Downlink</th>
                                        <th>RSRP </th>
                                        <th>
                                          Others
                                        </th>
                                      </tr>
                                    </thead>
                                     <tbody style={{fontSize: 12}}>
                                            {rows}
                                    </tbody>
                                  </table>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                );


    }
});

module.exports = DisplayRow;