'use strict';
angular.module('mrredisApp.dashboard', [])
	.config(['$stateProvider', function($stateProvider){
		$stateProvider
			.state('base.dashboard', {
				url: '/dashboard',
				resolve : {
					dbList :  function(dashboardServices){
						return dashboardServices.getDBList().then(function(data){
							return data;
						});
					}
				},
				views:{
					'redis' :{
						controller : 'dashboardController as dashboard',
						templateUrl : 'scripts/dashboard/views/dashboardView.html'
					}
				}
			});
	}]);