import 'angular'
import 'ng-file-upload'
import 'angular-recaptcha'

var app = angular.module('formSubmission', ['ngFileUpload', 'vcRecaptcha']);

app.controller('FormController', ['$scope', 'Upload', 'vcRecaptchaService', function ($scope, Upload, vcRecaptchaService) {

    $scope.uploader = {};
    $scope.form = {};
    $scope.submitted = false;
    $scope.loading = false;

    $scope.submit = function (carrier) {
        $scope.loading = true;
        Upload.upload({
            url: '/submit',
            data: {
                'file':     $scope.uploader.file,
                'carrier':  carrier,
                'form':     Upload.json($scope.form)
            }
        }).then(function (resp) {
            $scope.loading = false;
            $scope.submitted = true;
            $scope.submittedId = resp.data.id;
            $scope.uploader = {};
            $scope.form = {};
        }, function (resp) {
            $scope.loading = false;
            console.log('Error status: ' + resp.status);
        });
    };

}]);