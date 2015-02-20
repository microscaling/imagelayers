'use strict';

angular.module ('iLayers')
  .controller ('DashboardCtrl', ['$scope', 'registryService',
      function ($scope, registryService) {
        $scope.layers = [];

        $scope.metrics = {
          count: 0,
          size: 0,
          ave: 0,
          largest: 0
        }

         var calculateMetrics = function(layers) {
            for (var i=0; i < layers.length; i++) {
              $scope.metrics.count += 1;
              $scope.metrics.size += layers[i].Size;
              $scope.metrics.ave = $scope.metrics.size / $scope.metrics.count;
              $scope.metrics.largest = Math.max($scope.metrics.largest, layers[i].Size);
            }
        }


        // Load Data
        registryService.inspect().then(function(response){
            $scope.layers = response.data;
            calculateMetrics(response.data);
        });

  }]);
