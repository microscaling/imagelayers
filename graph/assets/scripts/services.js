'use strict';

angular.module ('iLayers')
  .factory ('registryService', ['$http',
    function ($http) {

      return {
          inspect: function () {
            return $http.post("/registry/analyze", { "repos": [{"name": "centurylink/image-graph", "tag": "latest"}, {"name": "centurylink/panamax-kubernetes-adapter", "tag": "latest"}] });
          }
      };
  }]);
