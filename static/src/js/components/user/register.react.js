var React = require('react');
var Link = require('react-router-component').Link;



var Register = React.createClass({


  render:function(){
    return (
          <div className="row centered-form"><br /><br /><br />
              <div className="col-xs-12 col-sm-8 col-md-4 col-sm-offset-2 col-md-offset-4">
                  <div className="panel panel-default">
                      <div className="panel-heading">
                          <h4 className="panel-title text-left">Create User</h4>
                      </div>
                        <div className="panel-body">
                            <form role="form"action="/register" method="post">
                              <div className="form-group">
                                    <input type="text" name="company_name" id="company_name" className="form-control input-sm" placeholder="Company Name" required />
                                </div>
                                <div className="form-group">
                                    <input type="text" name="user_firstName" id="user_firstName" className="form-control input-sm" placeholder="First Name" required />
                                </div>
                                <div className="form-group">
                                    <input type="text" name="user_lastName" id="user_lastName" className="form-control input-sm" placeholder="Last Name" required />
                                </div>
                             
                            

                                <div className="form-group">
                                    <input type="email" ng-model="user.email" name="email" id="email" className="form-control input-sm" placeholder="Email Address" required/>
                                </div>

                                
                            
                                        <div className="form-group">
                                            <input type="password" ng-model="user.password" name="password" id="password" className="form-control input-sm" placeholder="Password" required/>
                                        </div>
                                   
                                  
                                        <div className="form-group">
                                            <input type="password" name="password_confirmation" id="password_confirmation" className="form-control input-sm" placeholder="Confirm Password" required/>
                                        </div>
                                   
                             

                                <input type="Submit" value="Submit" className="btn btn-info" />
                            </form>
                        </div>
                    </div>
   
            </div>
        </div>
    );
  }
});

module.exports = Register;