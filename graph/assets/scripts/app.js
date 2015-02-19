'use strict';

angular.module('iLayers', ['ngRoute'])
  .config(['$httpProvider', '$locationProvider', '$routeProvider',
    function($httpProvider, $locationProvider, $routeProvider) {
      $httpProvider.defaults.withCredentials = false;

      $locationProvider.html5Mode(false);

      $routeProvider
      .when('/', {
        templateUrl: 'assets/views/dashboard.html',
        controller: 'DashboardCtrl'

      })
      .otherwise({
        redirectTo: '/'
      });

  }]);
