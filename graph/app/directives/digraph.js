'use strict';

angular.module ('iLayers')
  .directive('digraph', ['$timeout', '$sce', function($timeout, $sce) {
    return {
      restrict: 'AEC',
      scope: {
        graph: '='
      },
      replace: true,
      controller: function($scope) {
        var self = this;

        $scope.composeDigraph = function(graph) {
          var layerData = digraphHelper(graph, "", 0);
          var digraph = "digraph docker { " + layerData + "\n base [style=invisible]\n}";
          console.log('composeDigraph');
          return digraph;
        };


        var digraphHelper = function(layers, digraph, i) {
          if (layers[i] == null) return digraph;

          var id = layers[i].id.substring(0, 11);
          var tags = null;
          var parent_id = layers[i].parent.substring(0, 11);
          var line_break = "\n ";

          if (parent_id == "") {
            //digraph += "base -> \"" + id + "\" [style=invis];"; 
            digraph += line_break + "base -> \"511136ea3c5\" [style=invis]"
          }

          else {
            digraph += line_break + "\"" + parent_id +  "\"    -> \"" + id +"\"";
          }

          if (tags == null) {
            digraph += line_break + "\"" + id + "\" [label=\"" + id + " " + tags + "\",shape=box,fillcolor=\"paleturquoise\",style=\"filled,rounded\"];";
          }
            return digraphHelper(layers, digraph, i+1);
        };

      },
      link: function(scope, element) {
        scope.$watch('graph', function(graph) {
          var layers = [];
          var dotfile = "";
          for (var i=0; i< scope.graph.length; i++) {
            layers = layers.concat(scope.graph[i].layers);
            dotfile = scope.composeDigraph(layers);
            $timeout(function() { $('#graph').replaceWith(displayGraph(dotfile)); }, 2000);
          }

          function displayGraph(source, engine) {
            var result;
            try {
              result = Viz(source, 'svg', engine);
              return result;
              console.log("result is: " + result);
            } catch(e) {
              // TODO figure out error handling for this app
            }
          }

        });
      }
    }
  }]);
