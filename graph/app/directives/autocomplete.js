'use strict';

angular.module('iLayers')
  .directive('autocomplete', ['$timeout', function($timeout) {

    return {
      templateUrl: 'app/views/metrics.html',
      restrict: 'A'
    }
  }]);
