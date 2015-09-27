
import {describe, beforeEach, it, sinon, expect} from 'test/lib/common';

describe.only("QueryPart", () => {
  it('test', () => {
    var column = new QueryPartModel({
      type: 'column'
    });

    var agg = new QueryPartModel({
      type: 'aggregate'
      value: 'mean'
    });

    var derivate = new QueryPartModel({
      type: 'transform'
      value: 'derivate'
      params: [{type: 'interval', options: ['1s', '10s']}]
    });

    var mathPart = new QueryPartModel({
      type: 'math'
    });

    var asPart = new QueryPartModel({
      type: 'as'
    });

    var parts = [
      new QueryPart(column, 'value'),
      new QueryPart(agg, 'mean'),
      new QueryPart(derivate, '1s'),
      new QueryPart(as, 'test'),
    ];

    vat text = QueryPart.renderAll(segments);
    expect(text).to.be('derivate(mean("value"), 1s) as "test"');
  });
});

