'use strict';

angular.module ('iLayers')
  .directive('graph', ['$timeout', function($timeout) {
    return {
      restrict: 'E',
      scope: {
        graph: '='
      },
      controller: function($scope) {
        var self = this;

        var composeDigraph = function (graph) {
          var layerData = digraphHelper(graph, "", 0); 
          var digraph = "digraph docker { " + layerData + "\n base [style=invisible]\n}";  
          $scope.digraph = digraph;
        };


        var digraphHelper = function (layers, digraph, i) { 
          if (layers[i] == null) return digraph;  

          var id = layers[i].id.substring(0, 11); 
          var tags = null; 
          var parent_id = layers[i].parent.substring(0, 11); 
          var line_break = "\n "  ;

          if (parent_id == "") { 
            //digraph += "base -> \"" + id + "\" [style=invis];"; 
            digraph += line_break + "base -> \"511136ea3c5\" [style=invis]" 
          } 
          else { 
            digraph += line_break + "\"" + parent_id +  "\"    -> \"" + id +"\""; 
          }  
          if (tags != null) { 
            digraph += line_break + "\"" + id + "\" [label=\"" + id + " " + tags + "\",shape=box,fillcolor=\"paleturquoise\",style=\"filled,rounded\"];" 
          }  
            return digraphHelper(layers, digraph, i+1);
          };

      }
    }
  }]);
