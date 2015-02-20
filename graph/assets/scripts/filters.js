'use strict';

angular.module ('iLayers')
  .filter ('size', function () {
    var bytesToSize = function (bytes) {
       var sizes = ['bytes', 'kb', 'mb', 'gb', 'tb'];
       if (bytes == 1) return '1 byte';
       if (bytes == 0) return '0 bytes';
       var i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)));
       return Math.round(bytes / Math.pow(1024, i), 2) + ' ' + sizes[i];
    };

    return function (input) {
      return bytesToSize(input);
    };
  });
