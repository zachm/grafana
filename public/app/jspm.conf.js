System.config({
  baseURL: "public",
  defaultJSExtensions: true,
  transpiler: "none",
  paths: {
    "github:*": "vendor/jspm/github/*",
    "npm:*": "vendor/jspm/npm/*",
    "kbn": "app/components/kbn.js",
    "config": "app/components/config.js",
    "store": "app/components/store.js",
    "settings": "app/components/settings.js",
    "bootstrap": "vendor/bootstrap/bootstrap.js",
    "angular-ui": "vendor/angular-ui/angular-bootstrap.js",
    "angular-strap": "vendor/angular-other/angular-strap.js",
    "angular-dragdrop": "vendor/angular-native-dragdrop/draganddrop.js",
    "angular-bindonce": "vendor/angular-bindonce/bindonce.js",
    "spectrum": "vendor/spectrum.js",
    "filesaver": "vendor/filesaver.js",
    "bootstrap-tagsinput": "vendor/tagsinput/bootstrap-tagsinput.js"
  },

  map: {
    "angular": "github:angular/bower-angular@1.4.5",
    "angular-route": "github:angular/bower-angular-route@1.4.5",
    "angular-sanitize": "github:angular/bower-angular-sanitize@1.4.5",
    "jquery": "github:components/jquery@2.1.4",
    "lodash": "npm:lodash@3.10.1",
    "moment": "github:moment/moment@2.10.6",
    "github:angular/bower-angular-route@1.4.5": {
      "angular": "github:angular/bower-angular@1.4.5"
    },
    "github:angular/bower-angular-sanitize@1.4.5": {
      "angular": "github:angular/bower-angular@1.4.5"
    },
    "github:jspm/nodelibs-process@0.1.1": {
      "process": "npm:process@0.10.1"
    },
    "npm:lodash@3.10.1": {
      "process": "github:jspm/nodelibs-process@0.1.1"
    }
  }
});
