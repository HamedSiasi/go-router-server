var React = require('react');
var Router = require('react-router-component');
var Template = require('./app-template.js');
var Display = require('./display/display.react');
var Mode = require('./mode/mode.react');
var Login = require('./user/login.react');
var Register = require('./user/register.react');

var Locations = Router.Locations;
var Location  = Router.Location;

var App = React.createClass({
  render:function(){
    return (
      <Template>
        <Locations>
          <Location path="/" handler={Display} />
          <Location path="/mode" handler={Mode} />
          <Location path="#/login" handler={Login} />
          <Location path="#/register" handler={Register} />
        </Locations>
      </Template>
    );
  }
});

module.exports = App;

