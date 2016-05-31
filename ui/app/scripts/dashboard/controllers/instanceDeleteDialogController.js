'use strict';
angular.module('mrredisApp.dashboard')
  .controller('instanceDeleteDialogController', ['$scope', '$timeout', '$mdDialog', '$mdMedia','db', '$mdToast', 'dashboardServices', 
    function($scope, $timeout, $mdDialog, $mdMedia, db, $mdToast, dashboardServices){
      $scope.deletingInstance = false;        
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
        $scope.deletingInstance = true;          
        dashboardServices.deleteInstanceService($scope.dbToDelete.Name).then(function (response) {
          console.log("This is the response after deleting instance in controller: ");
          console.log(response);
          if(response && response.status === 200){
              response.reload = true;
              (function(){
              var promise = $timeout(function(){
                  $scope.creatingInstance = false;
                  $timeout.cancel(promise);
                  console.log('^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^');
                  console.log('RELOADING THE STATE');
                  console.log('Killed the promise');
                  $mdDialog.hide(response);
                }, 4000);
              })();
            }
        }, function(error){
          if(error && error.status === -1){
            error.msg = "Something went wrong. We could not delete the Instance"; 
            $mdDialog.hide(error);
          }
        });
      };
}]);
