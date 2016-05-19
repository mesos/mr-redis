'use strict';
	angular.module('mrredisApp.dashboard')
			.service('dashboardServices', ['$q','$timeout', 'ajaxService', 'api', function ($q, $timeout, ajaxService, api) {
				this.getDBList = function (){
					var defer  = $q.defer();
					var dbList = ajaxService.call(api.dbStatus.url, api.dbStatus.method, null);
					dbList.then(function(response){
						console.log(response.data);
						defer.resolve(response.data);
					},function(error){
						console.log(error);
						defer.reject(error);
					});

					return defer.promise;
				};



			//Create a new database instance
			//TODO: Currently uses GET in the backend (URL has the parameters) needs to be changed to POST
			
			/*this.createInstance = function(newInstanceData){
				var defer = $q.defer();

				// simulated async function
				$timeout(function() {
					defer.resolve('data received!')
				}, 2000)
				console.log(defer.promise);
				return defer.promise
			};*/
			
			this.createInstance = function(newInstanceData){
				var defer = $q.defer();
				var url = api.dbCreate.url + '/' + newInstanceData.name + '/' + newInstanceData.capacity + '/1/' + newInstanceData.slaves; //TODO : works only for one master. Add number of masters.
				var newInstance = ajaxService.call(url, api.dbCreate.method, newInstanceData);
				
				newInstance.then(function (response) {
					console.log(response);
					defer.resolve();
				},function(error){
				console.log(error);
				defer.reject(error);
				});
				return defer.promise;        		
			};

			//Delete database instance
			this.deleteInstanceService = function(databaseName){
				var defer = $q.defer();
				var url = api.dbDelete.url +'/'+ databaseName;
				var delInstance = ajaxService.call(url, api.dbDelete.method, null);
				delInstance.then(function(response){
					console.log("Response after Delete Service: ");
					console.log(response);
					defer.resolve(response);
				},function(error){
					console.log("Response Error after Delete Service: " + error);
					defer.reject(error);
				});
				return defer.promise;
			};
		}
])