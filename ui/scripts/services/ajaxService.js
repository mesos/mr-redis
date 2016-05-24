'use strict';
angular.module('mrredisApp')
	.factory('ajaxService', ['$http', '$q', 'api', function($http, $q, api){

		function call(url, method, payload){
			var defer = $q.defer();
			var endPoint = api.endPoint.url + url;
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