var React = require('react');
var Setting = require('./controls/setting.react');
var Header = require('./Header/Header.react')

var Template = React.createClass({

	render:function(){
		return (
		<div>
		  <Header />
		    {this.props.children}
		 </div>
		);
		
	}
});

module.exports = Template;