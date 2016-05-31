'use strict';
angular.module('mrredisApp.dashboard')
	.service('dashboardServices', ['$q','$timeout', 'ajaxService', 'api', function ($q, $timeout, ajaxService, api) {
		this.getDBList = function (){
			var defer  = $q.defer();
			var dbList = ajaxService.call(api.dbStatus.url, api.dbStatus.method, null);
			dbList.then(function(response){
				if('string' === typeof(response.data)){

					//Perform all calculations to make data human readable and push it with the response.
					//Calculate the memory utilization in percentage and set the UtilLevel to green, orange or red based on usage.
					//Doing this for both Master and Slaves.
					//TODO: instead of repeating write one function make a call.	
					response.noInstances = true;
					response.data = [];
					defer.resolve(response);
				}else{
					response.noInstances = (response.data.length > 0) ? false : true;
					response.data = getAllMetrics(response.data);
					defer.resolve(response);
				}
			},function(error){

				defer.reject(error);
			});
			return defer.promise;
		};
		//Get memory utilization in human readable form.
		function getMemoryUtilization(MemoryUsed, Capacity){
			var MemoryUtilized = {};
			var MemoryUsedInMB = MemoryUsed / (1024*1024);
			MemoryUtilized.MB = Math.round(MemoryUsedInMB * 100)/100;
			var MemoryUsedPercent = (MemoryUsedInMB / Capacity) * 100;
			MemoryUtilized.Percent = Math.round(MemoryUsedPercent *100)/100;
			return MemoryUtilized;
		}

		function getAllMetrics(Instance){
				var master_util_code=[]; 
				var slave_util_level=[];
				var slave_util_code=[];
				var master_uptime_hours=[];
				var slave_uptime_hours=[];
				var len = Instance.length;
				for(var i = 0; i < len; i++){
					if(Instance[i].Master && null !== Instance[i].Master){
						var master_util = getMemoryUtilization(Instance[i].Master.MemoryUsed, Instance[i].Master.MemoryCapacity)						
						master_uptime_hours[i] = getTimeInHours(Instance[i].Master.Uptime);
						if(master_util.Percent < 85){
							master_util_code[i] = "green";
						}else if(master_util.Percent > 85 && master_util.Percent < 95){
							master_util_code[i] = "orange";
						}else{
							master_util_code[i] = "red";
						}						
						Instance[i].Master.UtilCode = master_util_code[i];
						Instance[i].Master.UtilLevel = master_util.Percent;
						Instance[i].Master.UtilMB = master_util.MB;
						Instance[i].Master.UptimeHours = master_uptime_hours[i];
						//Do the same for Slaves - TODO: Make a single function and call it.
						
						if(Instance[i].Slaves && Instance[i].Slaves[0]) {
							var slave_len = Instance[i].Slaves.length;
							for (var j = 0; j < slave_len; j++) {
								var slaves_util = getMemoryUtilization(Instance[i].Slaves[j].MemoryUsed,Instance[i].Slaves[j].MemoryCapacity)

								slave_uptime_hours[j] = getTimeInHours(Instance[i].Slaves[j].Uptime);

								if(slaves_util.Percent < 85){
									slave_util_code[j] = "green";
								}else if(slaves_util.Percent > 85 && slaves_util.Percent < 95){
									slave_util_code[j] = "orange";
								}else{
									slave_util_code[j] = "red";
								}						

								Instance[i].Slaves[j].UtilCode = slave_util_code[j];
								Instance[i].Slaves[j].UtilLevel = slaves_util.Percent;
								Instance[i].Slaves[j].UtilMB = slaves_util.MB;
								Instance[i].Slaves[j].UptimeHours = slave_uptime_hours[j];
							}
						}else{
							Instance.notReady = true;
						}					}else{
						Instance.notReady = true;
					}

				}
				return Instance;

		}
		//Convert time to HH:MM:SS from seconds
		function getTimeInHours(time) {
		    var sec_num = parseInt(time, 10); // don't forget the second param
		    var hours   = Math.floor(sec_num / 3600);
		    var minutes = Math.floor((sec_num - (hours * 3600)) / 60);
		    var seconds = sec_num - (hours * 3600) - (minutes * 60);

		    if (hours   < 10) {hours   = "0"+hours;}
		    if (minutes < 10) {minutes = "0"+minutes;}
		    if (seconds < 10) {seconds = "0"+seconds;}
		    return hours+':'+minutes+':'+seconds;
		}

		this.getSingleInstanceDetails = function(instanceName){
			var defer  = $q.defer();
			var url = api.dbStatus.url + '/' + instanceName;
			var singleInstanceDetails = ajaxService.call(url, api.dbStatus.method, null);
			singleInstanceDetails.then(function(response){
				var responseArray = [];
				responseArray.push(response.data);
				response.data = getAllMetrics(responseArray);
				defer.resolve(response);
			},function(error){
				error.msg="Instance does not exist";
				defer.reject(error);
			});
			return defer.promise;
		}
		//Create a new database instance
		//TODO: Currently uses GET in the backend (URL has the parameters) needs to be changed to POST
		
		this.createInstance = function(newInstanceData, batch){
			var defer = $q.defer();
			var url = api.dbCreate.url + '/' + newInstanceData.name + '/' + newInstanceData.capacity + '/1/' + newInstanceData.slaves; //TODO : works only for one master. Add number of masters.
			var newInstance = ajaxService.call(url, api.dbCreate.method, newInstanceData);
			
			newInstance.then(function (response) {
				defer.resolve(response);
			},function(error){
				if(true === batch){
					defer.resolve(error);	
				}else{
					defer.reject(error);
				}				
			});
			return defer.promise;        		
		};


		//Delete database instance
		this.deleteInstanceService = function(databaseName){
			var defer = $q.defer();
			var url = api.dbDelete.url +'/'+ databaseName;
			var delInstance = ajaxService.call(url, api.dbDelete.method, null);
			delInstance.then(function(response){
				defer.resolve(response);
			},function(error){
				defer.reject(error);
			});
			return defer.promise;
		};

		//Delete database slave 
		/*this.deleteSlaveService = function(database){
			var defer = $q.defer();
			var url = api.dbDelete.url +'/'+ database.Name +'/'+ database.SlaveId;
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
		};*/


	}
])