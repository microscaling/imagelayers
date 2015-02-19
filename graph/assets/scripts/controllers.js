'use strict';

angular.module('iLayers')
  .controller('DashboardCtrl', ['$scope', 'registryService',
      function($scope, registryService) {
        $scope.data = registryService.inspect();
  }]);
