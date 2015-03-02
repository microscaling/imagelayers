'use strict';

angular.module('iLayers')
  .directive('search', function() {

    return {
      templateUrl: 'app/views/search.html',
      restrict: 'E',
      scope: {},

      controller: function($scope, $location) {
         $scope.changeSearch = function(newSearch) {
           $location.search('images', newSearch);
         };
      }
    }

  });
