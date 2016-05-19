/*'use strict';
	angular.module('mrredisApp.dashboard')
		.controller('toastController', ['$scope', '$mdDialog','$mdMedia', '$mdToast', 'dashboardServices', 
			function($scope, $mdDialog, $mdMedia, $mdToast, dashboardServices){
				$scope.closeToast = function() {
				$mdToast.hide();
			};
				

				//Success message on Database provisioning
				var last = {
			      bottom: false,
			      top: true,
			      left: false,
			      right: true
			    };
				$scope.toastPosition = angular.extend({},last);
				$scope.getToastPosition = function() {
				sanitizePosition();
				return Object.keys($scope.toastPosition)
					.filter(function(pos) { return $scope.toastPosition[pos]; })
					.join(' ');
				};
				function sanitizePosition() {
				var current = $scope.toastPosition;
					if ( current.bottom && last.top ) current.top = false;
					if ( current.top && last.bottom ) current.bottom = false;
					if ( current.right && last.left ) current.left = false;
					if ( current.left && last.right ) current.right = false;
					last = angular.extend({},current);
				}
			  $scope.showToast1 = function() {
					$mdToast.show(
						$mdToast.simple()
							.textContent('toastMessage')                       
							.hideDelay(3000)
					);
				};

               /*$scope.showToast2 = function() {
                  var toast = $mdToast.simple()
                     .textContent('Hello World!')
                     .action('OK')
                     .highlightAction(false);                     
                  $mdToast.show(toast).then(function(response) {
                     if ( response == 'ok' ) {
                        alert('You clicked \'OK\'.');
                     }
                  });			   
               };

 	}]); */*/