import 'angular'
import 'ng-file-upload'

var app = angular.module('formSubmission', ['ngFileUpload']);

app.controller('FormController', ['$scope', '$http', 'Upload', function ($scope, $http, Upload) {

    $scope.uploader = {};
    $scope.form = {};
    $scope.screen = 0;

    $scope.submit = function (carrier, captcha) {
        $scope.loading = true;
        $scope.statusSecret = null;

        if (captcha !== "") {
            grecaptcha.ready(function() {
                grecaptcha.execute(captcha, {action: 'submit'}).then(function(token) {
                    $scope.form.captcha = token;
                    $scope.actualSubmit(carrier);
                });
            });
        } else {
            $scope.actualSubmit(carrier);
        }
    };

    $scope.actualSubmit = function (carrier){
        Upload.upload({
            url: '/' + carrier + '/submit',
            data: {
                'file':     $scope.uploader.file,
                'secret':   Upload.json($scope.form)
            }
        }).then(function (resp) {
            $scope.loading = false;
            $scope.screen = 1;
            $scope.uploader = {};
            $scope.form = {};
            $scope.submitError = null;
            $scope.submittedId = resp.data.id;
        }, function (resp) {
            $scope.loading = false;
            $scope.submitError = resp.data.error;
        });
    }

    $scope.statusSettings = {}
    $scope.statusLookup = function (carrier) {
        $scope.loading = true;
        $scope.statusSecret = null;
        $http({
            method: 'GET',
            url: '/' + carrier + '/secret/' + $scope.lookupSettings.id,
        }).then(function (response) {
            $scope.loading = false;
            $scope.statusError = null;
            $scope.statusResult = response.data;
        }, function (response) {
            $scope.loading = false;
            $scope.statusError = response.data.error;
        })
    };
    $scope.statusDelete = function (carrier) {
        $scope.loading = true;
        $http({
            method: 'PATCH',
            url: '/' + carrier + '/secret/' + $scope.statusResult.id,
            data: {
                action: 'delete'
            }
        }).then(function (response) {
            $scope.loading = false;
            $scope.statusError = null;
            $scope.statusResult = response.data;
        }, function (response) {
            $scope.loading = false;
            $scope.statusError = response.data.error;
        })
    };

}]);