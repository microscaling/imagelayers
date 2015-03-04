'use strict';

angular.module('iLayers', ['ngRoute', 'ngAnimate', 'ngDialog'])
  .config(['$httpProvider', '$locationProvider', '$routeProvider',
    function($httpProvider, $locationProvider, $routeProvider) {
      $httpProvider.defaults.withCredentials = false;

      $locationProvider.html5Mode(false);

      $routeProvider
      .when ('/', {
        templateUrl: 'app/views/dashboard.html',
        controller: 'DashboardCtrl'

      })
      .otherwise ({
        redirectTo: '/'
      });

  }]);
