'use strict';
angular.module('mrredisApp')
	.factory('ajaxService', ['$rootScope', '$state', '$http', '$q', function($rootScope, $state, $http, $q){

		function call(url, method, payload){
			var defer = $q.defer();
			var endPoint = $rootScope.endPoint + url;
			$http({
				url : endPoint,
				method : method,
				data : payload
			}).then(function(response){
				defer.resolve(response);
			}, function(error){
				defer.reject(error);
				if(error.status === -1){							
					console.log('Uh-Oh! looks like the end point is not accessible.');
					$state.go('config');
				}
				
			});

			return defer.promise;
		}

		return {
			call : call
		};
	}]);