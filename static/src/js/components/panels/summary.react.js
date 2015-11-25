var React = require('react');

var Summary = React.createClass({
      render:function(){
              var rows = [];

              this.props.data.forEach(function(uuid, i) {


                  

                                      rows.push(  
                                    <div><br />
                                    <div className="col-lg-4">
                                    </div>
                                    <div className="col-lg-4" key={i}>
                                      <div className="panel panel-info" style={{height: 110, width: 350, marginTop: 20}}>
                                        <div className="panel-body">
                                          <p style={{fontStyle: 'italic'}}>

                                           <b>Total Msg:</b> <span className="resetColor">  {uuid["TotalMsgs"] }</span><br />
                                            <b>Total Bytes:</b> {uuid["TotalBytes"] }<br />
                                            <b>Last Msg:</b>    {uuid["LastMsgReceived"]}<br />
                                          </p> 
                                        </div>
                                      </div>
                                      </div>
                                      </div>  

                                   );   

                     
               
               


           
               });

          return (

              <div>
                {rows[0]}
             </div>
          );
      }
});

module.exports = Summary;