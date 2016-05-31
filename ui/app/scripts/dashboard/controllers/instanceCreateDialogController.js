'use strict';
  angular.module('mrredisApp.dashboard')
    .controller('instanceCreateDialogController', ['$scope', '$mdDialog', '$mdMedia','$mdToast', '$timeout', 'dashboardServices', 
      function($scope, $mdDialog,  $mdMedia, $mdToast, $timeout, dashboardServices){ 
        $scope.duplicateName = false;   
        $scope.creatingInstance = false;        
        $scope.customFullscreen = $mdMedia('xs') || $mdMedia('sm');
        $scope.hide = function() {
          $mdDialog.hide();
        }
        $scope.close = function() {
          var error = {
            status : true
          }
            $mdDialog.cancel(error);
        }
        $scope.save = function() {
            $mdDialog.hide(answer);
        }

        $scope.newInstance = {
          name: null,
          capacity: 32,
          masters: 1,
          slaves: 0
        };

      $scope.checkDBName = function (newInstanceName, callBack) {
        dashboardServices.getDBList().then(function(data){
          console.log('name being checked for returned data: ');
          console.log(data);
            if( undefined !== _.findWhere(data.data, {Name: newInstanceName})){
                $scope.duplicateName = true;
            }else{
                $scope.duplicateName = false;
                if(callBack){
                  callBack();
                }
              }
          });
        };

      //Create new database instance

      $scope.processCreateInstanceForm = function () {
        $scope.creatingInstance = true;
        $scope.checkDBName($scope.newInstance.name, function(){
          dashboardServices.createInstance($scope.newInstance).then(function(response){
            console.log('This is response from dashboardServices createInstance: ');
            console.log(response);
            if(response && response.status === 201){
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
            }else if(response && response.status === 200){
              $scope.creatingInstance = false;
              $scope.duplicateName = true;
            }
          },function(error){
            if(error && error.status === -1){
              error.msg = "Uh-oh! Something went wrong. We could not create the DB";
              $mdDialog.hide(error); 
            }
          });

        },function(error){
          console.log(error);
          //TODO Handle checkdb name API error failure
        });
      }
}]);
