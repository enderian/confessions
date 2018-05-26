import 'angular'
import 'ng-file-upload'
import 'angular-recaptcha'

var app = angular.module('formSubmission', ['ngFileUpload', 'vcRecaptcha']);

app.controller('FormController', ['$scope', 'Upload', 'vcRecaptchaService', function ($scope, Upload, vcRecaptchaService) {

    $scope.uploader = {};
    $scope.form = {};
    $scope.submitted = false;

    $scope.submit = function (carrier) {
        Upload.upload({
            url: '/submit',
            data: {
                'file':     $scope.uploader.file,
                'carrier':  carrier,
                'form':     Upload.json($scope.form)
            }
        }).then(function (resp) {
            $scope.submitted = true;
            $scope.submittedId = resp.data.id;
        }, function (resp) {
            console.log('Error status: ' + resp.status);
        });
    };

    $scope.upload = function (carrier, file) {
        console.log(file)

    };

}]);