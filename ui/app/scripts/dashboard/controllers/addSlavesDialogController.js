'use strict';
  angular.module('mrredisApp.dashboard')
    .controller('addSlavesDialogController', ['$scope', '$mdDialog', 'db', '$mdMedia','$mdToast', 'dashboardServices', 
      function($scope, $mdDialog, db, $mdMedia, $mdToast, dashboardServices){            
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

      	$scope.dbToAddSlaves = db;
}]);
