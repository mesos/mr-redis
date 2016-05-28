'use strict';

/**
 * @ngdoc overview
 * @name mrredisApp
 * @description
 * # mrredisApp
 *
 * Main module of the application.
 */
angular
  .module('mrredisApp', ['ngAnimate','ngCookies','ngResource','ngRoute','ngSanitize','ngMaterial', 'ui.router','ngMdIcons', 'md.data.table','ngMessages',
                        'pascalprecht.translate','highcharts-ng','mrredisApp.base', 'mrredisApp.config', 'mrredisApp.dashboard'])

  .config(function($mdThemingProvider) {
    $mdThemingProvider.theme('default')
    .primaryPalette('cyan',{
      'default': '900', // by default use shade 400 from the pink palette for primary intentions
      'hue-1': 'A700', // use shade 100 for the <code>md-hue-1</code> class
      'hue-2': '600', // use shade 600 for the <code>md-hue-2</code> class
      'hue-3': '700' // use shade A100 for the <code>md-hue-3</code> class
    })
    .accentPalette('blue')    
    .warnPalette('red', {
      'default': '900', // by default use shade 400 from the pink palette for primary intentions
      'hue-1': 'A700', // use shade 100 for the <code>md-hue-1</code> class
      'hue-2': '700', // use shade 600 for the <code>md-hue-2</code> class
      'hue-3': '600' // use shade A100 for the <code>md-hue-3</code> class
    });
  })
  .config(['$httpProvider', function($httpProvider) {
        $httpProvider.defaults.useXDomain = true;
        delete $httpProvider.defaults.headers.common['X-Requested-With'];
    }])
    .constant('_',
      window._
    )
  .run(['$state', '$cookies', '$rootScope', function($state, $cookies, $rootScope){
        $rootScope.$on('$stateChangeStart',function(e, toState, toParams, fromState, fromParams){
            if(undefined === $rootScope.endPoint && null === window.localStorage.getItem('endPoint')){
              if(toState.name !== 'config'){
                e.preventDefault();
                $state.go('config');
              }
            }else if(null !== window.localStorage.getItem('endPoint')){
              $rootScope.endPoint = window.localStorage.getItem('endPoint');
            }
        });
  }]);
