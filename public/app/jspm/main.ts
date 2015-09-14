///<reference path="../headers/common.d.ts" />

import * as _ from 'lodash';
import * as $ from 'jquery';
import kbn = require('kbn');

function bootstrap() {
  console.log('kbn', kbn);
  console.log('jquery', $);
  console.log('underscore / lodash', _);
  console.log('bootstrap');
}

bootstrap();

