'use strict';
angular.module('mrredisApp.config')
	.controller('configController', ['$rootScope', '$scope','$state', /*'$mdEditDialog', '$q',  '$mdDialog', '$mdMedia', '$mdToast' ,*/
		function ($rootScope, $scope, $state) {
			console.log('Entered configController: ');
			$scope.invalidUrl = false;			
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
					$state.go('base.dashboard');	
				}
				
			};
	}]);
