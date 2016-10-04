'use strict';
  angular.module('mrredisApp.dashboard')
    .controller('settingsDialogController', ['$rootScope', '$scope', '$mdDialog', '$mdMedia','$mdToast', '$timeout', 'dashboardServices', 
      function($rootScope, $scope, $mdDialog,  $mdMedia, $mdToast, $timeout, dashboardServices){ 
        $scope.duplicateName = false;   
        $scope.changingSettings = false;
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

        $scope.settings = {
          endpoint: $rootScope.endPoint,
          refreshInterval: 10
        };

      $scope.checkEndpointURI = function (newEndpoint, callBack, callBackError) {
        var urlPattern = /(http|https):\/\/[\w-]+(\.[\w-]+)+([\w.,@?^=%&amp;:\/~+#-]*[\w@?^=%&amp;\/~+#-])(:\d{2,4})?/;
        if (!urlPattern.test(newEndpoint)) {
          console.log('Url Test failed');
          $scope.invalidUrl = true;
          if(callBackError){
            callBackError();
          }
        } else {
          $scope.invalidUrl = false;
          if(callBack){
            callBack();
          }
        }
      };

      $scope.processSettingsForm = function () {
        $scope.changingSettings = true;
        var endpoint = $scope.settings.endpoint;
         
        $scope.checkEndpointURI(endpoint, function() {
          // Remove trailing slashes
          endpoint = endpoint.replace(/\/+$/, "");
          
          if (!endpoint.endsWith('/v1')) {
            endpoint += '/v1';
          }
          
          // Feedback sanitized URI to UI
          $scope.settings.endpoint = endpoint;
        
          // Ajax queries are sent using the global '$rootScope.endPoint'. Temporarily save
          // its old value in order to restore it in case the new endpoint is unreachable.
          var oldEndpoint = $rootScope.endPoint;
          $rootScope.endPoint = endpoint;

          // Test access to new endpoint.
          dashboardServices.getDBList().then(function(response) {
            $scope.changingSettings = false;
            var result = {
              reload: true,
              data: "Settings saved."
            }
            $mdDialog.hide(result);
            
          },function(error){
            console.log(error)
            
            $rootScope.endPoint = oldEndpoint
            $scope.invalidServer = true
            $scope.changingSettings = false;
          });

        }, function(error) {
          $scope.changingSettings = false;
        });
      }
}]);
