'use strict';

angular.module('iLayers')
  .controller('DashboardCtrl', ['$scope', 'registryService',
      function($scope, registryService) {

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
              // TODO get yer DOM outta my controller
              $('#graph').append("<div class='loading'>Loading...</div>");
              $scope.graph = response.data;
          });
        };
  }]);
