var AppConstants = require('../constants/app-constants');
var AppDispatcher = require('../dispatchers/app-dispatcher');

var AppActions = {
  setCommissioning: function(uuid){
    AppDispatcher.handleViewAction({
      actionType: AppConstants.SET_COMMISSIONING,
      uuid: uuid
    })
  },
  setTrafficTest: function(uuid){
    AppDispatcher.handleViewAction({
      actionType: AppConstants.SET_TRAFFIC_TEST,
      uuid: uuid
    })
  },
  setStandardTrx: function(uuid){
    AppDispatcher.handleViewAction({
      actionType: AppConstants.SET_STANDARD_TRX,
      uuid: uuid
    })
  },
  setHeartBeat: function(uuid){
    AppDispatcher.handleViewAction({
      actionType: AppConstants.SET_HEARTBEAT,
      uuid: uuid
    })
  },
  setReportingInterval: function(uuid){
    AppDispatcher.handleViewAction({
      actionType: AppConstants.SET_REPORTING_INTERVAL,
      uuid: uuid
    })
  },
  reboot: function(uuid){
    AppDispatcher.handleViewAction({
      actionType: AppConstants.REBOOT,
      uuid: uuid
    })
  }
}

module.exports = AppActions;
