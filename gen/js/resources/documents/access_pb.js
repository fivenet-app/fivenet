// source: resources/documents/access.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {missingRequire} reports error on implicit type usages.
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */
// @ts-nocheck

var jspb = require('google-protobuf');
var goog = jspb;
var global =
    (typeof globalThis !== 'undefined' && globalThis) ||
    (typeof window !== 'undefined' && window) ||
    (typeof global !== 'undefined' && global) ||
    (typeof self !== 'undefined' && self) ||
    (function () { return this; }).call(null) ||
    Function('return this')();

goog.exportSymbol('proto.resources.documents.ACCESS_LEVEL', null, global);
/**
 * @enum {number}
 */
proto.resources.documents.ACCESS_LEVEL = {
  BLOCKED: 0,
  VIEW: 1,
  COMMENT: 2,
  ACCESS: 3,
  EDIT: 4
};

goog.object.extend(exports, proto.resources.documents);
