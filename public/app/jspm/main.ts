///<reference path="../headers/common.d.ts" />
///<amd-dependency path="bootstrap" />
///<amd-dependency path="angular-strap" />
///<amd-dependency path="angular-route" />
///<amd-dependency path="angular-sanitize" />
///<amd-dependency path="angular-dragdrop" />
///<amd-dependency path="angular-bindonce" />
///<amd-dependency path="angular-ui" />

import _ = require('lodash');
import $ = require('jquery');
import bootstrap = require('bootstrap');
import kbn = require('kbn');
import angular = require('angular');

function initApp() {

  var app = angular.module('grafana', []);
  var register_fns: any = {};

  app.constant('grafanaVersion', "@grafanaVersion@");

  function useModule(module) {
    _.extend(module, register_fns);
    return module;
  }

  app.config(function($locationProvider, $controllerProvider, $compileProvider, $filterProvider, $provide) {
    // this is how the internet told me to dynamically add modules :/
    register_fns.controller = $controllerProvider.register;
    register_fns.directive  = $compileProvider.directive;
    register_fns.factory    = $provide.factory;
    register_fns.service    = $provide.service;
    register_fns.filter     = $filterProvider.register;
  });

  var apps_deps = [
    'ngRoute',
    'ngSanitize',
    '$strap.directives',
    'ang-drag-drop',
    'grafana',
    'pasvaz.bindonce',
    'ui.bootstrap.tabs',
  ];

  var module_types = ['controllers', 'directives', 'factories', 'services', 'filters', 'routes'];

  _.each(module_types, function (type) {
    var module_name = 'grafana.'+type;
    // create the module
    useModule(angular.module(module_name, []));
    // push it into the apps dependencies
    apps_deps.push(module_name);
  });

  var preBootRequires = [
    '../core/core',
    '../services/all',
    '../features/all',
    '../controllers/all',
    '../directives/all',
    '../components/partials',
    '../routes/all',
  ];

  require(preBootRequires, function () {
    // disable tool tip animation
    $.fn.tooltip.defaults.animation = false;

    // bootstrap the app
    angular.bootstrap(document, apps_deps);
  });
}

initApp();

