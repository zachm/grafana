System.config({
  baseURL: "public",
  defaultJSExtensions: true,
  transpiler: "none",
  paths: {
    "github:*": "vendor/jspm/github/*",
    "npm:*": "vendor/jspm/npm/*",
    "kbn": "app/components/kbn.js"
  },

  map: {
    "jquery": "github:components/jquery@2.1.4",
    "lodash": "npm:lodash@3.10.1",
    "moment": "github:moment/moment@2.10.6",
    "github:jspm/nodelibs-process@0.1.1": {
      "process": "npm:process@0.10.1"
    },
    "npm:lodash@3.10.1": {
      "process": "github:jspm/nodelibs-process@0.1.1"
    }
  }
});
