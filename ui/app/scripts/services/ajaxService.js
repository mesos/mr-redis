'use strict';
angular.module('mrredisApp')
	.factory('ajaxService', ['$rootScope','$http', '$q', function($rootScope, $http, $q){

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
			});

			return defer.promise;
		}

		return {
			call : call
		};
	}]);