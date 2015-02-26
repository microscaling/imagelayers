'use strict';

angular.module ('iLayers')
  .controller ('DashboardCtrl', ['$scope', '$timeout', 'registryService',
      function ($scope, $timeout, registryService) {
        $scope.layers = [];

        $scope.metrics = {
          count: 0,
          size: 0,
          ave: 0,
          largest: 0
        }

        var sequential = function (key, start, end, duration) {
            var range = end - start;
            var minTimer = 50;

            // calc step time to show all interediate values
            var stepTime = Math.abs(Math.floor(duration / range));

            // never go below minTimer
            stepTime = Math.max(stepTime, minTimer);

            // get current time and calculate desired end time
            var startTime = new Date().getTime();
            var endTime = startTime + duration;
            var timer;

            function run () {
                var now = new Date().getTime();
                var remaining = Math.max((endTime - now) / duration, 0);
                var value = Math.round(end - (remaining * range));
                $scope.metrics[key] = value;
                if (value != end) {
                    $timeout(run, stepTime);
                }
            }

            run();
        };

        var calculateMetrics = function (layers) {
          var count  = 0, size = 0, ave = 0, largest = 0;
          for (var i=0; i < layers.length; i++) {
            count += 1;
            size += layers[i].Size;
            ave = size / count;
            largest = Math.max(largest, layers[i].Size);
          };
          // animate the numbers
          sequential('count', 0, count, 600);
          sequential('size', 0, size, 520);
          sequential('ave', 0, ave, 520);
          sequential('largest', 0, largest, 520);
        }


        // Load Data
        registryService.inspect().then(function(response){
            $scope.layers = response.data.layers;
            calculateMetrics(response.data);
        });

  }]);
