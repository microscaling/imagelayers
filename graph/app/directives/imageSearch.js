'use strict';

angular.module ('iLayers')
  .directive('imageSearch', ['$sce', 'registryService',
    function($sce, registryService) {

    return {
      restrict: 'A',
      scope: {
        model: '='
      },
      templateUrl: 'app/views/imageSearch.html',
      controller: function($scope, $element, $attrs) {
        var self = this,
            constants = {
              max_results: 6,
              offset: 41
            };

        self.suggestImages = function(term) {
          if (term.length > 2) {
            return registryService.search(term).then(function(response){
              var data = response.data.results,
                  max = (data.length > constants.max_results) ? constants.max_results : data.length,
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

        self.attached = function(element) {
          $('.ac-container').css("top", (element[0].offsetTop + constants.offset) + "px");
        };

        self.selectImage = function(selected) {
          console.log("selected", selected);
          registryService.find_tags(selected.value).then(function(response) {
            var data = Object.keys(response.data);

            $scope.tag_list = [];
            for (var i=0; i < data.length; i++) {
              console.log("data " + i + " value: " + data[i]);
              $scope.tag_list.push({ 'label': data[i], 'tag': data[i] });
            }
          });
        };

        $scope.tag_list = []

        $scope.autocomplete_options = {
          suggest: self.suggestImages,
          on_error: console.log,
          on_attach: self.attached,
          on_select: self.selectImage
        };

      },
      link: function(scope, element) {
      }
    }
  }]);
