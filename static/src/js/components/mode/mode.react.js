/**
 * Copyright (c) 2014, U-blox.
 * All rights reserved.
 */
var React = require('react');
var Setting = require('../controls/setting.react');

var Mode = React.createClass({

  /**
   * @return {object}
   */
  render: function() {
   
       return (

      <div className="row" ><br /> 
        <div className="col-lg-8" style={{width:'400px', height: '180px'}}>
          <div className="panel panel-info">
            <div className="panel-body" style={{width:'400px', height: '180px'}}>
             <div style={{float:'left'}}>
                <select id="example-multiple-selected" multiple="multiple" style={{width:'110px', height: '150px'}}>
                  <option value={1}>UE-P</option>
                  <option value={2} selected="selected">UE-TW</option>
                   <option value={3}>UE-WW</option>
                  <option value={4}>UE-T</option>
                  <option value={5}>UE-R</option>
                  <option value={6}>UE-E</option>
                </select>
              </div>

              <div style={{float:'right', marginRight:'30px'}}>
                <Setting />
              </div>
              <p />
            </div>
          </div>
        </div> 
        <div className="col-lg-4">

              <div className="ScrolledContent" style={{float:'right', width:'300px', height:'100px'}}>
    a<br />a<br />a<br />a<br />a<br />a<br />a<br />a<br />a<br />a<br />a<br />a
              </div>

        </div>
        <div className="row">
        <div className="collg-4">
        </div>
          <div className="col-lg-12" style={{marginTop:'10px'}}>
         
            <div className="panel panel-info" >
              <div className="panel-heading">
                Operating Modes
              </div>
              {/* /.panel-heading */}
              <div className="panel-body">
                {/* Nav tabs */}
                <ul className="nav nav-tabs">
                  <li className="active"><a href="#standard" data-toggle="tab">Standard</a>
                  </li>
                  <li><a href="#commissioning" data-toggle="tab">Commissioning</a>
                  </li>
                  <li><a href="#traffictest" data-toggle="tab">Traffic Test</a>
                  </li>
                  <li><a href="#standalonebasic" data-toggle="tab">Standalone TRX Basic</a>
                  </li>
                </ul>
                {/* Tab panes */}
                <div className="tab-content">
                  <div className="tab-pane fade in active" id="standard">
                    <h4 style={{marginLeft: 20}}>Basic Mode</h4>
                    <p>
                    </p><div className="col-lg-12">
                      <div className="panel panel-info" style={{height: 250}}>
                        <i className="fa fa-spinner fa-spin" style={{float: 'right', padding: 10, color: 'red'}} />
                        <div className="panel-body"><p>
                          </p><div className="col-lg-8">
                            <div className="col-1-3">
                              <div className="cell red">Status: Disconnected
                              </div>
                            </div>
                          </div>
                          <div className="col-lg-8">
                            <div className="col-1-3">
                              <div className="cell grey">WAKE UP CODE: -
                              </div>
                            </div>
                            <div className="col-1-3">
                              <div className="cell grey">REVISION LEVEL: -
                              </div>
                            </div>
                            <div className="col-1-3 push-1-3">
                            </div>
                            <div className="col-1-3">
                              <div className="cell grey">LAST SEEN: - 04/10/2015
                              </div>
                            </div>
                            <div className="col-1-3">
                              <div className="cell grey">REPORTING INTERVAL: 
                              </div>
                            </div>
                            <div className="col-1-3 push-1-3">
                            </div>
                            <div className="col-1-3">
                              <div className="cell grey">HEART BEAT: -
                              </div>
                            </div>
                            <div className="col-1-3">
                              <div className="cell grey">DATE TIME: 
                              </div>
                            </div>
                          </div>
                          <p />
                        </div>
                      </div>
                    </div>
                    <p />
                  </div>
                  <div className="tab-pane fade" id="commissioning">
                    <h4 style={{marginLeft: 20}}>Commissioning Mode</h4>  
                    <p>
                    </p><div className="col-lg-12">
                      <div className="panel panel-info" style={{height: 250}}>
                        <i className="fa fa-spinner fa-spin" style={{float: 'right', padding: 10, color: 'red'}} />
                        <div className="panel-body"><p>
                          </p><div className="col-lg-8">
                            <div className="col-1-3">
                              <div className="cell red">Status: Disconnected
                              </div>
                            </div>
                          </div>
                          <div className="col-lg-8">
                            <div className="col-1-3">
                              <div className="cell grey">RSSI: -
                              </div>
                            </div>
                            <div className="col-1-3">
                              <div className="cell grey">RSRP: -
                              </div>
                            </div>
                            <div className="col-1-3 push-1-3">
                            </div>
                            <div className="col-1-3">
                              <div className="cell grey">CELL ID: - 
                              </div>
                            </div>
                            <div className="col-1-3">
                              <div className="cell grey">SNR: 
                              </div>
                            </div>
                          </div>
                          <p />
                        </div>
                      </div>
                    </div>
                    <p />
                  </div>
                  <div className="tab-pane fade" id="traffictest">
                    <h4 style={{marginLeft: 20}}>Traffic Test Mode</h4>
                    <p>
                    </p><div className="col-lg-4">
                      <div className="panel panel-info" style={{height: 300}}>
                        <div className="panel-body"><p>
                            <input className="form-control" placeholder=" No Of UL Datagrams" style={{width: 180, height: 30, margin: 5}} />
                            <input className="form-control" placeholder=" Size Of UL Datagrams" style={{width: 180, height: 30, margin: 5}} />
                            <input className="form-control" placeholder=" No Of DL Datagrams" style={{width: 180, height: 30, margin: 5}} />
                            <input className="form-control" placeholder=" Size Of DL Datagrams" style={{width: 180, height: 30, margin: 5}} /><br />
                            <button type="button" className="btn btn-info" style={{width: 100, height: 30, marginTop: 10}}>Set up Test</button>
                            <button type="button" className="btn btn-default" style={{width: 100, height: 30, marginTop: 10}}>Start Test</button>
                          </p>
                        </div>
                      </div>
                    </div>
                    <div className="panel panel-info" style={{height: 300}}>
                      <i className="fa fa-spinner fa-spin" style={{float: 'right', padding: 10, color: 'red'}} />
                      <div className="panel-body"><p>
                        </p><div className="col-lg-8" style={{float: 'right'}}>
                          <div className="col-1-3">
                            <div className="cell red">Status: Disconnected
                            </div> 
                          </div>
                        </div>
                        <div className="col-lg-8" style={{float: 'right'}}>
                          <div className="col-1-3">
                            <div className="cell grey">No of UL Datagrams: -
                            </div>
                          </div>
                          <div className="col-1-3">
                            <div className="cell grey">No of UL Bytes: -
                            </div>
                          </div>
                          <div className="col-1-3 push-1-3">
                          </div>
                          <div className="col-1-3">
                            <div className="cell grey">No of DL Datagrams: - 
                            </div>
                          </div>
                          <div className="col-1-3">
                            <div className="cell grey">No of DL Bytes: 
                            </div>
                          </div>
                          <div className="col-1-3 push-1-3">
                          </div>
                          <div className="col-1-3">
                            <div className="cell grey">No of DL Datagrams Missed: -
                            </div>
                          </div>
                        </div>
                        <p />
                      </div>
                      <div className="tab-pane fade" id="standalonebasic">
                        <h4 style={{marginLeft: 20}}>Basic Mode</h4>
                        <p>
                        </p><div className="col-lg-12">
                          <div className="panel panel-info" style={{height: 250}}>
                            <i className="fa fa-spinner fa-spin" style={{float: 'right', padding: 10, color: 'red'}} />
                            <div className="panel-body"><p>
                              </p><div className="col-lg-8">
                                <div className="col-1-3">
                                  <div className="cell red">Status: Disconnected
                                  </div>
                                </div>
                              </div>
                              <div className="col-lg-8">
                                <div className="col-1-3">
                                  <div className="cell grey">No of UL Datagrams: -
                                  </div>
                                </div>
                                <div className="col-1-3">
                                  <div className="cell grey">No of UL Bytes: -
                                  </div>
                                </div>
                                <div className="col-1-3 push-1-3">
                                </div>
                                <div className="col-1-3">
                                  <div className="cell grey">No of DL Datagrams: - 
                                  </div>
                                </div>
                                <div className="col-1-3">
                                  <div className="cell grey">No of DL Bytes: 
                                  </div>
                                </div>
                              </div>
                              <p />
                            </div>
                          </div>
                        </div>
                        <p />
                        <h4 style={{marginLeft: 20}}>Advanced Mode</h4>
                        <p>
                        </p><div className="col-lg-12">
                          <div className="panel panel-info" style={{height: 250}}>
                            <div className="panel-body"><p>
                              </p><div className="col-lg-8">
                                <div className="col-1-3">
                                  <div className="cell red">Status: Disconnected
                                  </div>
                                </div>
                              </div>
                              <div className="col-lg-8">
                                <div className="col-1-3">
                                  <div className="cell grey">No of UL Datagrams: -
                                  </div>
                                </div>
                                <div className="col-1-3">
                                  <div className="cell grey">No of UL Bytes: -
                                  </div>
                                </div>
                                <div className="col-1-3 push-1-3">
                                </div>
                                <div className="col-1-3">
                                  <div className="cell grey">No of DL Datagrams: - 
                                  </div>
                                </div>
                                <div className="col-1-3">
                                  <div className="cell grey">No of DL Bytes: 
                                  </div>
                                </div>
                              </div>
                              <p />
                            </div>
                          </div>
                        </div>
                        <p />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              {/* /.panel-body */}
            </div>
            {/* /.panel */}
            </div>
        </div>

      </div>
    );


  }
});

module.exports = Mode;