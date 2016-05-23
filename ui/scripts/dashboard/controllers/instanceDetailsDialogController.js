'use strict';
  angular.module('mrredisApp.dashboard')
    .controller('instanceDetailsDialogController', ['$scope', '$mdDialog', 'dbDetails', '$mdMedia','$mdToast', 'dashboardServices', 
      function($scope, $mdDialog, dbDetails, $mdMedia, $mdToast, dashboardServices){            
        $scope.customFullscreen = $mdMedia('xs') || $mdMedia('sm');
        $scope.hide = function() {
          $mdDialog.hide();
        }
        $scope.close = function() {
            $mdDialog.cancel();
        }
        $scope.save = function() {
            $mdDialog.hide(answer);
        }
        
        
      	$scope.dbShowDetails = dbDetails;        
        console.log('The DB to show details: ');
        console.log($scope.dbShowDetails);
        
}]);
