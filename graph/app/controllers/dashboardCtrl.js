'use strict';

angular.module('iLayers')
  .controller('DashboardCtrl', ['$scope', 'registryService', 'commandService',
      function($scope, registryService, commandService) {

        var self = this;

        //private
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

        // public
        $scope.graph = [];

        $scope.searchImages = function(images) {
          var search_terms = self.buildTerms(images);

          // Load Data
          registryService.inspect(search_terms).then(function(response){
              $scope.graph = response.data;
          });
        };

        $scope.highlightCommand = function(image, idx) {
          commandService.highlight(image.layers.slice(0, idx+1));
        };

        $scope.clearCommands = function() {
          commandService.clear();
        };
  }]);
