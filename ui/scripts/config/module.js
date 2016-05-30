'use strict';
angular.module('mrredisApp.config', [])
	.config(['$stateProvider', function($stateProvider){
		$stateProvider
			.state('config', {
				url: '/config',				
				controller : 'configController as config',
				templateUrl : 'scripts/config/views/configView.html'
					
			});
	}]);