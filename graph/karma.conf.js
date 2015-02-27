// Karma configuration
// Generated on Tue Feb 24 2015 13:04:12 GMT-0600 (CST)

module.exports = function(config) {
  config.set({

    // base path that will be used to resolve all patterns (eg. files, exclude)
    basePath: 'app',


    // frameworks to use
    // available frameworks: https://npmjs.org/browse/keyword/karma-adapter
    frameworks: ['jasmine'],


    // list of files / patterns to load in the browser
    files: [
      'scripts/jquery-1.11.2.min.js',
      'scripts/angular.min.js',
      'scripts/angular-mocks.js',
      'scripts/angular-route.min.js',
      'scripts/angular-animate.min.js',
      'app.js',
      'controllers/**/*.js',
      'filters/**/*.js',
      'services/**/*.js',
      'directives/**/*.js',
      'specs/**/*.js',
      '**/*.html'
    ],


    // list of files to exclude
    exclude: [
    ],


    // preprocess matching files before serving them to the browser
    // available preprocessors: https://npmjs.org/browse/keyword/karma-preprocessor
    preprocessors: {
       '**/*.html': ['ng-html2js']
    },

    ngHtml2JsPreprocessor: {
      prependPrefix: 'app/'
    },

    plugins: [
      'karma-phantomjs-launcher',
      'karma-jasmine',
      'karma-spec-reporter',
      'karma-ng-html2js-preprocessor'
    ],

    // test results reporter to use
    // possible values: 'dots', 'progress'
    // available reporters: https://npmjs.org/browse/keyword/karma-reporter
    reporters: ['spec'],


    // web server port
    port: 9876,


    // enable / disable colors in the output (reporters and logs)
    colors: true,


    // level of logging
    // possible values: config.LOG_DISABLE || config.LOG_ERROR || config.LOG_WARN || config.LOG_INFO || config.LOG_DEBUG
    logLevel: config.LOG_INFO,


    // enable / disable watching file and executing tests whenever any file changes
    autoWatch: true,


    // start these browsers
    // available browser launchers: https://npmjs.org/browse/keyword/karma-launcher
    //browsers: ['Chrome'],
    browsers: ['PhantomJS'],


    // Continuous Integration mode
    // if true, Karma captures browsers, runs the tests and exits
    singleRun: true
  });
};
