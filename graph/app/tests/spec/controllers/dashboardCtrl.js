describe('DashboardCtrl', function() {
  // Load the module
  beforeEach(module('iLayers'));

  var ctrl, scope, layers;

  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();

    ctrl = $controller('DashboardCtrl', {
      $scope: scope
    });
  }));

  it('should initialize metrics', function() {
    expect(angular
             .equals(scope.metrics, { count: 0, size: 0, ave: 0, largest:0 }))
             .toBeTruthy();
  });

  it('should initialize layers', function() {
    expect(scope.layers.length).toEqual(0);
  });


  describe('buildTerms', function () {
    it('should add latest tag when empty', function() {
      var data = ctrl.buildTerms("foo");

      expect(data[0].tag).toEqual("latest");
      expect(data[0].name).toEqual("foo");
    });

    it('should return tag and name when provided', function() {
      var data = ctrl.buildTerms("foo:1.0.0");

      expect(data[0].tag).toEqual("1.0.0");
      expect(data[0].name).toEqual("foo");
    });

    it('should create terms for each image provided', function() {
      var data = ctrl.buildTerms("foo:1.0.0, baz:2.0.0");

      expect(data.length).toEqual(2);
      expect(angular
               .equals(data[0], { "name": "foo", "tag": "1.0.0" }))
               .toBeTruthy();
      expect(angular
               .equals(data[1], { "name": "baz", "tag": "2.0.0" }))
               .toBeTruthy();
    });
  });

  describe('calculateMetrics', function() {
    beforeEach(function() {
      spyOn(ctrl, "sequential");

      layers = [
        { name: 'foo', Size: 300 },
        { name: 'baz', Size: 200 },
        { name: 'bar', Size: 1000 }
      ]

    });

    it('should call sequential with the total layer count', function() {
      ctrl.calculateMetrics(layers);
      expect(ctrl.sequential).toHaveBeenCalledWith('count', 0, 3, 600);
    });

    it('should call sequential with the total layer size', function() {
      ctrl.calculateMetrics(layers);
      expect(ctrl.sequential).toHaveBeenCalledWith('size', 0, 1500, 520);
    });

    it('should call sequential with the layer average', function() {
      ctrl.calculateMetrics(layers);
      expect(ctrl.sequential).toHaveBeenCalledWith('ave', 0, 500, 520);
    });

    it('should call sequential with the largest layer size', function() {
      ctrl.calculateMetrics(layers);
      expect(ctrl.sequential).toHaveBeenCalledWith('largest', 0, 1000, 520);
    });
  });
});
