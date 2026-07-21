var goRuntimeReady = false;
var goRuntimeError = null;
var goPending = [];

function initGoRuntime() {
  if (!window.Go) {
    goRuntimeError = 'Go WASM support not available (wasm_exec.js not loaded)';
    flushPending();
    return;
  }

  const go = new Go();

  var wasmUrl = document.currentScript
    ? document.currentScript.src.replace(/gorunner\.js$/, 'gorunner.wasm')
    : 'wasm/gorunner.wasm';

  WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject)
    .then(function(result) {
      go.run(result.instance);
      goRuntimeReady = true;
      flushPending();
    })
    .catch(function(err) {
      goRuntimeError = 'Failed to load Go WASM runtime: ' + err.message;
      flushPending();
    });
}

function flushPending() {
  var queue = goPending.slice();
  goPending = [];
  queue.forEach(function(item) {
    item.fn.apply(null, item.args);
  });
}

function whenReady(fn, args) {
  if (goRuntimeError) {
    var errCb = args[args.length - 1];
    if (typeof errCb === 'function') {
      errCb(goRuntimeError);
    }
    return;
  }
  if (goRuntimeReady) {
    fn.apply(null, args);
    return;
  }
  goPending.push({ fn: fn, args: args });
}

var _runGo = function(code, onSuccess, onError) {
  whenReady(function(code, onSuccess, onError) {
    try {
      var result = window._goWasmRun(code);
      if (result.error) {
        onError(result.error);
      } else {
        onSuccess(result.stdout || '');
      }
    } catch (e) {
      onError(e.message || 'Unknown error');
    }
  }, [code, onSuccess, onError]);
};

var _runGoTest = function(funcCode, funcName, testCases, onResult) {
  whenReady(function(funcCode, funcName, testCases, onResult) {
    try {
      var testCasesJSON = JSON.stringify(testCases);
      var result = window._goWasmTest(funcCode, funcName, testCasesJSON);
      if (result.error) {
        onResult({ error: result.error });
      } else {
        onResult({ results: result.results });
      }
    } catch (e) {
      onResult({ error: e.message || 'Unknown error' });
    }
  }, [funcCode, funcName, testCases, onResult]);
};

initGoRuntime();
