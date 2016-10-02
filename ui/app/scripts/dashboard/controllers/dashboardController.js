'use strict';
angular.module('mrredisApp.dashboard')
	.controller('dashboardController', ['$state', '$mdEditDialog', '$q', '$scope', '$timeout','$mdDialog', '$mdMedia', '$mdToast', 'dashboardServices', 'dbList',
		function ($state, $mdEditDialog, $q, $scope, $timeout, $mdDialog, $mdMedia, $mdToast, dashboardServices, dbList) {
			
			//Populating the data tables source.
			//The database list is sent from the dashboard and ajax service as dbList
			$scope.reloading = false;
			$scope.dbList = dbList;

			$scope.databases = {				
				"count": $scope.dbList.data.length,
				"data": $scope.dbList.data,
				"noInstances" : $scope.dbList.noInstances,
				"search" : ''
			};
			
			
			
			
			$scope.limitOptions = [5, 10, 15, {
				label: 'All',
				value: function () {					
					return $scope.databases ? $scope.databases.count : 0;
				}
			}];

			//Reload the table
			$scope.reload = function(){	
				$scope.reloading = true;			
				$state.reload();
			};
			

			 // Toolbar search toggle
			 $scope.isHidden = true;
 			$scope.toggleSearch = function() {
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
				order: 'Name',
				limit: 10,
				page: 1
			};

			$scope.toggleLimitOptions = function () {
				$scope.limitOptions = $scope.limitOptions ? undefined : [5, 10, 15];
			};
			$scope.onPaginate = function(page, limit) {
			    console.log('Scope Page: ' + $scope.query.page + ' Scope Limit: ' + $scope.query.limit);
			    console.log('Page: ' + page + ' Limit: ' + limit);

			    $scope.promise = $timeout(function () {

			    }, 2000);
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
						if(true === response.reload){
							var toast = $mdToast.simple()
			                	.textContent(response.data)
			                	.action('Ok')
			                	.hideDelay(5000)
			                	.position('bottom left');
							$mdToast.show(toast).then(function(response){
								console.log(response);
							});
							$state.reload();
						}
					}, function(error) {
						if(false === error.status){
							$mdToast.show(
					            $mdToast.simple()
								.textContent(error.msg)
								.action('Ok')
								.hideDelay(5000)
								.position('bottom left')
			            	);
						}
				});
			};

			//Show batch Create new instances form in modal
			$scope.batchCreate = function (event) {
				$mdDialog.show({
					clickOutsideToClose: false,  
					controller: 'instanceBatchCreateDialogController',    
					focusOnOpen: false,
					targetEvent: event,
					templateUrl: 'scripts/dashboard/views/batchCreateInstanceView.html',
				}).then(function(response) {
						if(true === response){
							var toast = $mdToast.simple()
				                  .textContent('Requested instances have been created')
				                  .action('Ok')
				                  .hideDelay(5000)
				                  .position('bottom left');
							$mdToast.show(toast);
							$state.reload();
						}
					}, function(error) {
						$mdToast.show(
				            $mdToast.simple()
				              .textContent(error.msg)
				              .action('Ok')
				              .hideDelay(5000)
				              .position('bottom left')
			            );
					});
			};
			
			//Show UI settings
			$scope.showSettings = function (event) {
				$mdDialog.show({
					clickOutsideToClose: false,  
					controller: 'settingsDialogController',    
					focusOnOpen: false,
					targetEvent: event,
					templateUrl: 'scripts/dashboard/views/settingsView.html',
				}).then(function(response) {
					if(true === response.reload) {
						var toast = $mdToast.simple()
							.textContent(response.data)
							.action('Ok')
							.hideDelay(5000)
							.position('bottom left');
						$mdToast.show(toast).then(function(response){
							console.log(response);
						});
						$state.reload();
					}
				}, function(error) {
						if(false === error.status) {
						$mdToast.show(
							$mdToast.simple()
								.textContent(error.msg)
								.action('Ok')
								.hideDelay(5000)
								.position('bottom left')
						);
					}
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
							console.log(response);
						});
						$state.reload();
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

