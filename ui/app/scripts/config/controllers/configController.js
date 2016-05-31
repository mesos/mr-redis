'use strict';
angular.module('mrredisApp.config')
	.controller('configController', ['$rootScope', '$scope','$state', 'ajaxService', 'api', /*'$mdEditDialog', '$q',  '$mdDialog', '$mdMedia', '$mdToast' ,*/
		function ($rootScope, $scope, $state, ajaxService, api) {
			console.log('Entered configController: ');
			$scope.invalidUrl = false;
			$scope.endPointNotReachable = false;
			$scope.showEndPointLoading = false;			
			$scope.checkUrl = function(){
				var urlPattern = /(http|https):\/\/[\w-]+(\.[\w-]+)+([\w.,@?^=%&amp;:\/~+#-]*[\w@?^=%&amp;\/~+#-])(:\d{2,4})?/ ;
				if(!urlPattern.test($scope.endPoint)){
					console.log('Url Test failed');
					$scope.invalidUrl = true;
				}else{
					$scope.invalidUrl =false;
				}
			}
			$scope.setEndPoint = function(){
				$scope.endPointNotReachable = false;
				$scope.showEndPointLoading = true;
				var endPointTest = $scope.endPoint;				
				var lastChar = endPointTest.substr(-1);
				if(lastChar === '/'){
					endPointTest = endPointTest.substr(0, endPointTest.length-1);				
					$scope.endPoint = endPointTest;
				}
				$scope.checkUrl();
				if(!$scope.invalidUrl){
					window.localStorage.setItem('endPoint', $scope.endPoint);
					$rootScope.endPoint = $scope.endPoint;
					var dbList = ajaxService.call(api.dbStatus.url, api.dbStatus.method, null);
					dbList.then(function(response){
						$state.go('base.dashboard');
					},function(error){
						if(error.status === -1){
							$scope.endPointNotReachable = true;
							$scope.showEndPointLoading = false;			
							console.log('Uh-Oh! looks like the end point is not accessible.');
							

						}
					});
						
				}
				
			};
	}]);
