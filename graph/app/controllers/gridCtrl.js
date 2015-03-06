'use strict';

angular.module('iLayers')
  .controller('SearchCtrl', ['$scope',
    function($scope) {
        var self = this;

        //private
        self.buildGrid = function(data) {
        };

        // public
        // The grid will be a matrix -
        $scope.grid = [];
    }]);
