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
						}, function(error){
							console.log('Entered error block in dashboard module. No instances: ');
							console.log(error);
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