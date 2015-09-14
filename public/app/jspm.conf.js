System.config({
  defaultJSExtensions: true,
  transpiler: "none",
  paths: {
    "github:*": "public/jspm_packages/github/*",
    "npm:*": "public/jspm_packages/npm/*"
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
