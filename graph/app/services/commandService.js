'use strict';

angular.module('iLayers')
  .factory('commandService', ['$rootScope',
    function($rootScope) {
      return {
          highlight: function (layers) {
            var cmds = [];

            for(var i=0; i < layers.length; i++) {
              var command = layers[i].container_config.Cmd;
              cmds.push(command.join(" "));
            };

            $rootScope.$broadcast('command-change', { 'commands': cmds });
          },

          clear: function() {
            $rootScope.$broadcast('command-change', { 'commands': [] });
          }
      };
  }]);
