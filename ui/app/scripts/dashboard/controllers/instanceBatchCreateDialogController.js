'use strict';
  angular.module('mrredisApp.dashboard')
    .controller('instanceBatchCreateDialogController', ['$scope', '$q', '$mdDialog', '$mdMedia','$mdToast', '$timeout', 'dashboardServices', 
      function($scope, $q, $mdDialog,  $mdMedia, $mdToast, $timeout,dashboardServices){ 
        $scope.duplicateBatchName = false;
        $scope.showBatchProgress = false;
        $scope.totalNumberofInstances = 0;
        $scope.originalCount = 0;
        $scope.createInstanceDelta = 0;
        $scope.createInstancePromises = [];
        $scope.successfullyCreatedInstances = [];
        $scope.unSuccessfullyCreatedInstances = [];
        $scope.hours = 0;
        $scope.minutes = 0;
        $scope.seconds = 0;
        $scope.customFullscreen = $mdMedia('xs') || $mdMedia('sm');
        $scope.hide = function() {
          $mdDialog.hide(true);
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

      //batch create the instances.

      $scope.processBatchCreateInstanceForm = function () {
        $scope.showBatchProgress = true;
        
        dashboardServices.getDBList().then(function(data){
          $scope.originalCount = data.data.length;
          $scope.totalNumberofInstances = data.data.length + $scope.newBatchInstance.quantity;
          $scope.udpateProgress();
          $scope.startTimer();
          for (var i = 0; i < $scope.newBatchInstance.quantity; i++){
            //TOD: Change the date.Now() to something more unique
            var name = $scope.newBatchInstance.name + '-' +i+ '-' + Date.now();
            var instanceData = {
              name: name,
              capacity: $scope.newBatchInstance.capacity,
              masters: 1,
              slaves: 0
            };
            $scope.createInstancePromises.push(dashboardServices.createInstance(instanceData, true));
          }
          
          $q.all($scope.createInstancePromises).then(function(response){
            for(var x = 0, len = response.length; x < len; x++){
              if(201 === response[x].status){
                $scope.successfullyCreatedInstances.push(response[x]);
              }else{
                //TODO: Check the failure status as well as API failure and add the condition
                //Maybe the above might not be needed 
                $scope.unSuccessfullyCreatedInstances.push(response[x]);
              }
            }
          });

        });

      };

      $scope.udpateProgress = function(){
        var promise = $timeout(function(){
        dashboardServices.getDBList().then(function(data){
          console.log('Created instances:');
          console.log(data.data.length);
          $scope.createInstanceDelta = data.data.length - $scope.originalCount;
          $scope.progressIndicator = Math.floor(((data.data.length - $scope.originalCount) / $scope.newBatchInstance.quantity) * 100);
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
          $scope.seconds = $scope.seconds + 1;
          if($scope.seconds === 60){
            $scope.minutes = $scope.minutes + 1;
            $scope.seconds = 0;
            if($scope.minutes === 60){
              $scope.hours = $scope.hours + 1;
              $scope.minutes = 0;
              $scope.seconds = 0;
            }
          }
          if(!$scope.stopTimer){
            $scope.startTimer();
          }
          $timeout.cancel(promise1); 
        }, 1000);
      }
}]);
