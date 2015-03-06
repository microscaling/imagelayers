'use strict';

angular.module('iLayers')
  .directive('autocomplete', ['$timeout', 'registryService', function($timeout, registryService) {
    var updateAutoComplete = function(element) {
      var value = element.val();

      registryService.search(value).then(function(response){
        var data = response.data.results,
            list = [];

        for (var i=0; i < data.length; i++) {
          list.push(data[i].name);
        };

        console.log("list ", list);
        console.log("raw  ", element);

        element.autocomplete({
          source: list,
          select: function() {
            $timeout(function() {
              element.trigger('input');
              return false;
            }, 0);
          }
        });
      });


    }

    return {
      restrict: 'A',
      link: function(scope, element, attrs) {
            element.bind('keyup', function (e) { updateAutoComplete(element); });
      }
    }
  }]);
