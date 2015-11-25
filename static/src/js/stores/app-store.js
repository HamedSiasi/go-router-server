var AppDispatcher = require('../dispatchers/app-dispatcher');
var AppConstants = require('../constants/app-constants');
var assign = require('react/lib/Object.assign');
var EventEmitter = require('events').EventEmitter;

var CHANGE_EVENT = 'change';


function _setCommissioning(uuid){
}
function _setTrafficTest(uuid){
}
function _setStandardTrx(uuid){
}
function _setHeartBeat(uuid){
}
function _setReportingInterval(uuid){
}
function _reboot(uuid){
}
function addUser(company, firstName, lastName, email, password) {
  obj = {};
  obj.company = company;
  obj.firstName = firstName;
  obj.lastName = lastName;
  obj.email = email;
  obj.password = password;

  $.ajax({
        url: 'http://localhost:3000/register',
        dataType: 'json',
        method: 'put',
        async: false,
        data: obj,
        success: function(data) {
          return  
        },
        error: function(xhr, status, err) {
          console.error('/', status, err.toString());
        }
    });
  return
}

  var AppStore = assign(EventEmitter.prototype, {
  emitChange: function(){
    this.emit(CHANGE_EVENT)
  },

  addChangeListener: function(callback){
    this.on(CHANGE_EVENT, callback)
  },

  removeChangeListener: function(callback){
    this.removeListener(CHANGE_EVENT, callback)
  },

  getUtmsData: function(){
    return states
  },


  dispatcherIndex: AppDispatcher.register(function(payload){
    var action = payload.action; // this is our action from handleViewAction
    switch(action.actionType){

      case AppConstants.SET_COMMISSIONING:
        _setCommissioning(payload.action.index);
        break;
      case AppConstants.SET_TRAFFIC_TEST:
        _setTrafficTest(payload.action.index);
        break;
      case AppConstants.SET_STANDARD_TRX:
        _setStandardTrx(payload.action.index);
        break;
      case AppConstants.SET_HEARTBEAT:
        _setHeartBeat(payload.action.index);
        break;
      case AppConstants.SET_REPORTING_INTERVAL:
        _setReportingInterval(payload.action.index);
        break;
      case AppConstants.REBOOT:
        reboot(payload.action.index);
        break;
    }

    AppStore.emitChange();

    return true;
  })

})

module.exports = AppStore;


