import 'angular'
import 'particles.js'

import particlesConfig from './particlesjs-config'

const app = angular.module('indexApp', []);

app.run([ function () {
    particlesJS('particles-js', particlesConfig);
}]);

