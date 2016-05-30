'use strict';
  angular.module('mrredisApp.dashboard')
    .controller('instanceDetailsDialogController', ['$scope', '$state', '$mdDialog', 'dbDetails', '$mdMedia','$mdToast', 'dashboardServices', 
      function($scope, $state, $mdDialog, dbDetails, $mdMedia, $mdToast, dashboardServices){            
        $scope.customFullscreen = $mdMedia('xs') || $mdMedia('sm');
        $scope.reloading = false;
        $scope.hide = function() {
          $mdDialog.hide();
        }
        $scope.close = function() {
            $mdDialog.cancel();
        }
        $scope.save = function() {
            $mdDialog.hide(answer);
        }

        console.log('Before adding to scope of details: ');
        console.log(dbDetails);
      	$scope.dbShowDetails = dbDetails;        
        console.log('The DB to show details: ');
        console.log($scope.dbShowDetails);

        //Reload the table
        $scope.reloadDetails = function(instanceName){ 
          $scope.reloading = true;
          dashboardServices.getSingleInstanceDetails(instanceName).then(function(data){
            console.log('The Single instance details fetched: ');
            console.log(data.data);
            $scope.dbShowDetails = data.data[0];
            $scope.reloading = false;
            console.log('The reload details: ');
            console.log($scope.dbShowDetails);
          },function(error){
              console.log('Cannot Reload details: ');
              error.msg="Something went wrong. Cannot reload";

            }
          );
        };


        
}]);
