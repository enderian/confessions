ender confessions
===

This project contains all of the source code necessary to run the ender confessions client-side application on deployment hardware.

Software that was used
---

The project is mainly written in [GoLang](https://golang.org/), and supported by these open-source and free libraries and projects:
* https://angularjs.org/ (Handles our single-page landing webpage)
* https://getbootstrap.com/ (Provides their beautiful theme and components)
* https://jquery.com/ (Lightweight DOM manipulation and AJAX)
* https://github.com/CodeSeven/toastr (delivers notifications)
* http://underscorejs.org/ (powerful data manipulations)
* https://webpack.js.org/ (compiles our JavaScript and CSS into a single file)

You can find the respected licenses on each of these project's website.

Running confessions
---
This project was not created with running in mind, so prepare yourself for some heavy configuration before getting it to actually work!
So if you are determined, *brace thyself!*

To start running confessions, you will need:
* A running MongoDB instance locally on your machine.
* A `config.js` file containing the following parameters:
  * `port`, containing the host and port to run confessions on (for example: `":8080"`)
  * `confessions_images`, containing the directory to store all images into.
  * `recaptcha_key`, which is the site's ReCaptcha key, and: 
  * `recaptcha_secret`, for verifying ReCaptcha requests.
* You need to have compiled the JavaScript and CSS:
  * Make sure you have NodeJS and NPM installed.
  * Do `npm install` to download the nessesarry dependencies.
  * Finally, do `npm run build` to run webpack and compile the JS and CSS.
* With the [dep](https://golang.github.io/dep/) tool, make sure all dependencies are available by doing `dep ensure` whilst in the `src/ender.gr/confessions` directory.
* Configure the `GOPATH` and `GOROOT` environment variables.

After all of these requirements are met, you can do:
```bash
go build ender.gr/confessions
./confessions
```
Or to make it more simple, you can run:
```bash
chmod u+x start.sh
./start.sh #or
./start.sh --skip-webpack #to skip building webpack
```
And that's it! You should now have an operational confessions instance online!