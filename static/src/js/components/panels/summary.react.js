var React = require('react');

var Summary = React.createClass({
      render:function(){
              var rows = [];

                  if(this.props["SummaryData"] !== undefined){                  

                                      rows.push(  
                                    <div><br />
                                    <div className="col-lg-4">
                                    </div>
                                    <div className="col-lg-4">
                                      <div className="panel panel-info" style={{height: 110, width: 350, marginTop: 20}}>
                                        <div className="panel-body">
                                          <p style={{fontStyle: 'italic'}}>

                                           <b>Total Uplink Msgs:</b> <span className="resetColor">  {this.props["SummaryData"]['TotalUlMsgs'] }</span><br />
                                            <b>Total Downlink Msgs:</b> <span className="resetColor">  {this.props["SummaryData"]["TotalDlMsgs"] }</span><br />
                                            <b>Total Bytes:</b> {this.props["SummaryData"]["TotalDlBytes"] }<br />
                                            <b>Last Msg:</b>    {this.props["SummaryData"]["LastDlMsgTime"]}<br />
                                          </p> 
                                        </div>
                                      </div>
                                      </div>
                                      </div>  

                                   );   

                     
               }
               




          return (

              <div>
                {rows[0]}
             </div>
          );
      }
});

module.exports = Summary;