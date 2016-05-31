'use strict';
angular.module('mrredisApp.base', [])
		.config(['$stateProvider', '$urlRouterProvider', function ($stateProvider, $urlRouterProvider) {			
			$urlRouterProvider.otherwise('/config');
			$stateProvider
					.state('base',{
						url: '/redis',
						templateUrl: 'views/redisView.html',
						controller: 'baseController'
					})																					
					.state('error', {
						url: '/error',
						controller: 'baseController',
						templateUrl: 'views/error.html'
					});
					
		}]);