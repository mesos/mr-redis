mrredisApp.directive('syncFocusWith', function($timeout, $rootScope) {
	return {
		restrict: 'A',
		scope: {
			focusValue: "=syncFocusWith"
		},
		link: function($scope, $element, attrs) {
			$scope.$watch("focusValue", function(currentValue, previousValue) {
				if (currentValue === true && !previousValue) {
					$element[0].focus();
				} else if (currentValue === false && previousValue) {
					$element[0].blur();
				}
			})
		}
	}
});