'use strict';
  angular.module('mrredisApp.dashboard')
    .controller('instanceCreateDialogController', ['$scope', '$mdDialog', '$mdMedia','$mdToast', 'dashboardServices', 
      function($scope, $mdDialog,  $mdMedia, $mdToast, dashboardServices){ 
        $scope.duplicateName = false;           
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
        $scope.checkDBName($scope.newInstance.name, function(){
          dashboardServices.createInstance($scope.newInstance).then(function(response){
            console.log('This is response from dashboardServices createInstance: ');
            console.log(response);
            if(response && response.status === 201){
              response.reload = true;
              response.noInstances = false;
              $mdDialog.hide(response);                              
            }
          },function(error){
            if(error && error.status === -1){
              error.msg = "Uh-oh! Something went wrong. We could not create the DB";
              $mdDialog.hide(error); 
            }
          });

        });
      }
}]);
