'use strict';

angular.module('iLayers')
  .controller('DashboardCtrl', ['$scope', '$routeParams', 'registryService',
      function($scope, $routeParams, registryService) {

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

        self.searchImages = function(route) {
          if (route.images !== undefined) {
            var search_terms = self.buildTerms(route.images);

            // Load Data
            registryService.inspect(search_terms).then(function(response){
              $scope.graph = response.data;
            });
          }
        };

        // public
        $scope.graph = [];

        // Load data from RouteParams
        self.searchImages($routeParams);
  }]);
