describe('DashboardCtrl', function() {
  // Load the module with MainController
  beforeEach(module('iLayers'));

  var ctrl, scope;

  beforeEach(inject(function($controller, $rootScope) {
    scope = $rootScope.$new();

    ctrl = $controller('DashboardCtrl', {
      $scope: scope
    });
  }));

  it('should create $scope.greeting when calling sayHello',
    function() {
      expect(false).toBeTruthy();
  });
})
