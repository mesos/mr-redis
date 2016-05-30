'use strict';
  angular.module('mrredisApp.dashboard')
    .controller('instanceBatchCreateDialogController', ['$scope', '$q', '$mdDialog', '$mdMedia','$mdToast', '$timeout', 'dashboardServices', 
      function($scope, $q, $mdDialog,  $mdMedia, $mdToast, $timeout,dashboardServices){ 
        $scope.duplicateBatchName = false;
        $scope.showBatchProgress = false;
        $scope.totalNumberofInstances = 0;
        $scope.createInstancePromises = [];
        $scope.timer = 0;
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

        $scope.newBatchInstance = {
          name: null,
          capacity: 32,
          masters: 1,
          slaves: 0,
          quantity: 0
        };

      /*$scope.checkDBName = function (newInstanceName, callBack) {
        dashboardServices.getDBList().then(function(data){
            if( undefined !== _.findWhere(data, {Name: newInstanceName})){
              $scope.duplicateName = true;
            }else{
              $scope.duplicateName = false;
              if(callBack){
                callBack();
              }
            }
        });
      };*/

      //batch create the instances.

      $scope.processBatchCreateInstanceForm = function () {
        $scope.showBatchProgress = true;
        
        dashboardServices.getDBList().then(function(data){
          console.log('Existing instances: ');
          console.log(data.data.length);
          console.log('Requested instances: ');
          console.log($scope.newBatchInstance.quantity);
          $scope.totalNumberofInstances = data.data.length + $scope.newBatchInstance.quantity;
          $scope.udpateProgress();
          $scope.startTimer();
          
          for (var i = 0; i < $scope.newBatchInstance.quantity; i++){
            var instanceData = {
              name: $scope.newBatchInstance.name+'-'+i+'-'+Date.now(),
              capacity: $scope.newBatchInstance.capacity,
              masters: 1,
              slaves: $scope.newBatchInstance.slaves
            };
            $scope.createInstancePromises.push(dashboardServices.createInstance(instanceData, true));
          }

          $q.all($scope.createInstancePromises).then(function(response){
            console.log('After batch create promise');
            console.log(response);
          });

        });

      };

      $scope.udpateProgress = function(){
        
        var promise = $timeout(function(){
        console.log('^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^');
        console.log('RELOADING THE STATE');
        dashboardServices.getDBList().then(function(data){
          console.log('The total number of instances: ' + $scope.totalNumberofInstances);
          console.log('Created: ' + data.data.length);
          console.log('Remaining: ');
          console.log($scope.totalNumberofInstances - data.data.length);
          if($scope.totalNumberofInstances != data.data.length){
            $scope.udpateProgress();
          }else{
            $scope.stopTimer = true;
          }
        })
        $timeout.cancel(promise);
        }, 1000);
      }

      $scope.startTimer = function(){
        var promise1 = $timeout(function(){
          $scope.timer = $scope.timer + 1; 
          if(!$scope.stopTimer){
            $scope.startTimer();
          }
          $timeout.cancel(promise1); 
        }, 1000);
      }
}]);
