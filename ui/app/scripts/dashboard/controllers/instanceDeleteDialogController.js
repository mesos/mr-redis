'use strict';
angular.module('mrredisApp.dashboard')
  .controller('instanceDeleteDialogController', ['$scope', '$mdDialog', '$mdMedia','db', '$mdToast', 'dashboardServices', 
    function($scope, $mdDialog, $mdMedia, db, $mdToast, dashboardServices){            
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

      $scope.dbToDelete = db;
      console.log("Delete DB request: ")
      console.log($scope.dbToDelete);
      $scope.deleteInstance = function(){          
        dashboardServices.deleteInstanceService($scope.dbToDelete.Name).then(function (response) {
          console.log("This is the response after deleting instance in controller: ");
          console.log(response);
          if(response && response.status === 200){
            response.reload = true;
            $mdDialog.hide(response);
          }
        }, function(error){
          if(error && error.status === -1){
            error.msg = "Something went wrong. We could not delete the DB"; 
            $mdDialog.hide(error);
          }
        });
      };
}]);
