'use strict';
angular.module('mrredisApp.base')
	.controller('baseController',['$rootScope', '$scope', '$state', 
		function($rootScope, $scope, $state){
			console.log('Base Controller');
			console.log($state.current.name);
		    if($state.current.name === 'base'){
		        $state.go('base.dashboard');
		    }
		}
	]);			
