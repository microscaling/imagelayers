'use strict';

angular.module('iLayers')
  .controller('DashboardCtrl', ['$scope', '$timeout', 'registryService',
      function($scope, $timeout, registryService) {

        var self = this;

        // public
        $scope.layers = [];

        $scope.metrics = {
          count: 0,
          size: 0,
          ave: 0,
          largest: 0
        };

        $scope.searchImages = function(images) {
          var search_terms = self.buildTerms(images);

          // Load Data
          registryService.inspect(search_terms).then(function(response){
              $scope.layers = response.data;
              self.calculateMetrics(response.data);
          });
        };


        // private
        self.sequential = function(key, start, end, duration) {
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

            function run() {
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

        self.calculateMetrics = function(layers) {
          var count  = 0, size = 0, ave = 0, largest = 0;
          for (var i=0; i < layers.length; i++) {
            count += 1;
            size += layers[i].Size;
            ave = Math.floor(size / count);
            largest = Math.max(largest, layers[i].Size);
          };
          // animate the numbers
          self.sequential ('count', 0, count, 600);
          self.sequential ('size', 0, size, 520);
          self.sequential ('ave', 0, ave, 520);
          self.sequential ('largest', 0, largest, 520);
        };

        self.buildTerms = function(data) {
          var terms = data.split(','),
              search_terms = [];

          for (var i=0; i < terms.length; i++) {
            var name = terms[i].split(":")[0],
                tag = "latest";
            if (terms[i].lastIndexOf(':') != -1) {
              tag = terms[i].split(":")[1]
            }
            search_terms.push({
              "name": name.trim(),
              "tag": tag
            });
          }

          return search_terms;
        };
  }]);
