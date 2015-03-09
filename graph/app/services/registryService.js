'use strict';

angular.module('iLayers')
  .factory('registryService', ['$http',
    function($http) {

      return {
          inspect: function (repo_list) {
            return $http.post("/registry/analyze", { "repos": repo_list });
          },
          search: function(name) {
            return $http.get("/registry/search?name="+name);
          },
          find_tags: function(name) {
            return $http.get("/registry/images/"+name+"/tags");
          }
      };
  }]);
