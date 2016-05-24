'use strict';
angular.module('mrredisApp.base', [])
		.config(['$stateProvider', '$urlRouterProvider', function ($stateProvider, $urlRouterProvider) {
			// body...
			$urlRouterProvider.otherwise('/redis/dashboard');
			$stateProvider
					.state('base',{
						url: '/redis',
						templateUrl: 'views/redisView.html',
						controller: 'baseController'
					})					
					.state('error', {
						url: '/error',
						templateUrl: 'views/error.html'
					});
					
		}]);