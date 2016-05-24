'use strict';
angular.module('mrredisApp.dashboard')
	.controller('dashboardController', ['$state', '$mdEditDialog', '$q', '$scope', '$timeout','$mdDialog', '$mdMedia', '$mdToast', 'dashboardServices', 'dbList',
		function ($state, $mdEditDialog, $q, $scope, $timeout, $mdDialog, $mdMedia, $mdToast, dashboardServices, dbList) {
			//Populating the data tables source.
			//The database list is sent from the dashboard and ajax service as dbList
			$scope.dbList = dbList;				
			$scope.databases = {
				"count": 7,
				"data": $scope.dbList
			};
			
			$scope.selected = [];
			$scope.limitOptions = [5, 10, 15, {
				label: 'All',
				value: function () {					
					return $scope.databases ? $scope.databases.data.count : 0;
				}
			}];

			//Reload the table
			$scope.reload = function(){				
				$state.reload();

			};
			

			 // Toolbar search toggle
			 $scope.isHidden = true;
 			$scope.toggleSearch = function(element) {
    			$scope.isHidden = $scope.isHidden ? false : true;
  			};
  			//Set the md Data table options
			$scope.options = {
				rowSelection: false,
				multiSelect: false,
				autoSelect: false,
				decapitate: false,
				largeEditDialog: false,
				boundaryLinks: true,
				limitSelect: true,
				pageSelect: true
			};
  
			$scope.query = {
				order: 'name',
				limit: 10,
				page: 1
			};

			$scope.toggleLimitOptions = function () {
				$scope.limitOptions = $scope.limitOptions ? undefined : [5, 10, 15];
			}; 
			
		  			  
			$scope.logItem = function (item) {
				console.log(item.name, 'was selected');
			};
		  
			$scope.logOrder = function (order) {
				console.log('order: ', order);
			};
		  
			$scope.logPagination = function (page, limit) {
				console.log('page: ', page);
				console.log('limit: ', limit);
			}

			//Create new Database instance form in modal
			$scope.showCreate = function (event) {
				$mdDialog.show({
					clickOutsideToClose: false,  
					controller: 'instanceCreateDialogController',    
					focusOnOpen: false,
					targetEvent: event,
					templateUrl: 'scripts/dashboard/views/instanceCreateView.html',
				}).then(function(response) {
						console.log("entered the success state after creation");
						if(true === response.reload){
							var toast = $mdToast.simple()
				                  .textContent(response.data)
				                  .action('Ok')
				                  .hideDelay(5000)
				                  .position('bottom left');
							$mdToast.show(toast).then(function(response){
								console.log("Response from the toast promise on create success: ")
								console.log(response);
								if(response === "ok"){
									$state.reload();		
								}
							$state.reload();	
							});
			             	 	 
						}
					}, function(error) {
						console.log("Database not created. Entered Error");
						$mdToast.show(
				            $mdToast.simple()
				              .textContent(error.msg)
				              .action('Ok')
				              .hideDelay(5000)
				              .position('bottom left')
			            );
					});
			};

			//Display Single Database Instance details 
			$scope.displayInstanceDetails = function (database, event) {			          
				$mdDialog.show({
					clickOutsideToClose: false,  
					controller: 'instanceDetailsDialogController',    
					focusOnOpen: false,
					targetEvent: event,
					templateUrl: 'scripts/dashboard/views/instanceDetailsView.html',
					dbDetails: database					
				}).then(function(response) {
					$scope.alert = 'You said the information was "' + answer + '".';
				}, function() {
					$scope.alert = 'You cancelled the dialog.';
				});
			};
			
			//Delete Single Database instance
			$scope.showDeleteInstance = function(database, event) {
				console.log("database name sent to delete dialog: " + database.Name);
				$mdDialog.show({
					controller: 'instanceDeleteDialogController',
					templateUrl: 'scripts/dashboard/views/instanceDeleteView.html',
					targetEvent: event,
					db: database
				}).then(function(response) {
					if(true === response.reload){
						var toast = $mdToast.simple()
							.textContent(response.data)
							.action('Ok')
							.hideDelay(6000)
							.position('bottom left');
						$mdToast.show(toast).then(function(response){
							console.log("Response from the toast promise on create success: ")
							console.log(response);
							if(response === "ok"){
								$state.reload();		
							}
							$state.reload();		
						});

					}
				}, function(error) {
					if(error && error.status === -1){	
						$mdToast.show(
						$mdToast.simple()
							.textContent(error.msg)
							.action('Ok')
							.hideDelay(6000)
							.position('bottom left')
						);
					}
				});
			};  

			//Add slaves dynamically to the master
			$scope.addSlaves = function(database, event){
				console.log('The database to add slaves to:' + database.Name);
				$mdDialog.show({
					controller: 'addSlavesDialogController',
					templateUrl: 'scripts/dashboard/views/addSlavesView.html',
					targetEvent: event,
					db: database
				})
				.then(function(answer) {
					$scope.alert = 'You said the information was "' + answer + '".';
				}, function() {
					$scope.alert = 'You cancelled the dialog.';
				});
			}
}]);

