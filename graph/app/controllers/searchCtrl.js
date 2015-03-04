'use strict';

angular.module('iLayers')
  .controller('SearchCtrl', ['$scope', 'ngDialog', 'registryService',
      function($scope, ngDialog, registryService) {
        var self = this;

        $scope.showSearch = function() {
          console.log("search");
          ngDialog.open({ template: 'app/views/searchDialog.html', className: 'ngdialog-theme-layers' });
        };

      }]);
