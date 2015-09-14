module.exports = function(config) {
  'use strict';
  return {
    options: {
      sfx: false,
      baseURL: "./public_gen",
      configFile: "./public_gen/app/jspm.conf.js",
      minify: false,
      build: {
        mangle: false
      }
    },
    dist: {
      files: [{
        "src":  "./public_gen/app/jspm/build.js",
        "dest": "./public_gen/app/bundle.js"
      }]
    }
  };
};
