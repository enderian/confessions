import 'angular/angular';
import 'angular-route/angular-route';
import 'jquery';
import 'bootstrap';
import './app-index.scss'
import clients from './carriers.json'

const app = angular.module('indexApp', ['ngRoute']);

app.controller('ClientController', ['$scope', function ($scope) {
    $scope.clients = clients;
}]);
app.run(['$rootScope','$location', function($rootScope, $location) {
    $rootScope.$on('$routeChangeSuccess', function(e, current, pre) {
        $rootScope.currentLocation = $location.path();
    });
}]);

app.config(['$routeProvider', '$locationProvider', function ($routeProvider, $locationProvider) {
    $routeProvider
        .when('/', {
            templateUrl: 'route/index.html'
        })
        .when('/privacy', {
            templateUrl: 'route/privacy.html'
        })
        .when('/help', {
            templateUrl: 'route/help.html'
        })
        .otherwise('/')
}]);

$(function () {
    $('[data-toggle="popover"]').popover({trigger: 'focus'})
});