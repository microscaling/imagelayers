'use strict';

angular.module('iLayers')
  .controller('SearchCtrl', ['$scope', '$location', '$sce', 'ngDialog', 'registryService',
      function($scope, $location, $sce, ngDialog, registryService) {
        var self = this;

        self.buildQueryParams = function(list) {
          var params = [];
          for (var i=0; i < list.length; i++) {
            if (list[i].tag === '') {
              params.push(list[i].name);
            } else {
              params.push(list[i].name + ':' + list[i].tag);
            }
          };

          return params.join(',');
        };

        self.suggestImages = function(term) {
          if (term.length > 2) {
            return registryService.search(term).then(function(response){
              var data = response.data.results,
                  list = [];

              for (var i=0; i < data.length; i++) {
                list.push({ 'label': $sce.trustAsHtml(data[i].name), 'value': data[i].name});
              };
              return list;
            });
          } else {
             return []
          }

        };



        $scope.searchList = [];

        $scope.autocomplete_options = {
          suggest: self.suggestImages,
          on_error: console.log
        };

        $scope.showSearch = function() {
          $scope.searchList = [];
          ngDialog.open({
            template: 'app/views/searchDialog.html',
            className: 'ngdialog-theme-layers',
            controller: 'SearchCtrl' });
        };

        $scope.addRow = function() {
          var rand = Math.random() * 10;
          var item = {
            'name': '',
            'tag': 'latest'
          };
          $scope.searchList.push(item);
        };

        $scope.closeDialog = function() {
          ngDialog.closeAll();
        };

        $scope.addImages = function() {
          $location.search('images', self.buildQueryParams($scope.searchList));
          ngDialog.closeAll();
        };

      }]);
