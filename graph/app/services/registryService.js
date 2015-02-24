'use strict';

angular.module ('iLayers')
  .factory ('registryService', ['$http',
    function ($http) {

      return {
          inspect: function (repo_list) {
            return $http.post("/registry/analyze", { "repos": repo_list });
          }
      };
  }]);
