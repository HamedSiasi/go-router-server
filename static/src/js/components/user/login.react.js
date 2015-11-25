var React = require('react');
var Link = require('react-router-component').Link;

var Login = React.createClass({
  render:function(){
    return (   
          <div className="row centered-form"><br /><br /><br />
              <div className="col-xs-12 col-sm-8 col-md-4 col-sm-offset-2 col-md-offset-4">
                  <div className="panel panel-default">
                      <div className="panel-heading">
                          <h4 className="panel-title text-left">Please Sign In</h4>
                      </div>
                      <div className="panel-body">
                          <form role="form" action="/login" method="post">
                            <div className="form-group">
                                  <input type="email" name="email" id="email" className="form-control input-sm" placeholder="Email Address" required autofocus />
                              </div>
                              <div className="form-group">
                                  <input type="password" name="user_password" id="user_password" className="form-control input-sm" placeholder="User password"  required />

                              </div>
                              <input  type="submit" value="Login" className="btn btn-info" />
                          </form>
                      </div>
                  </div>
              </div>
          </div>
    );
  }
});

module.exports = Login;