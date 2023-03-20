// source: services/livemapper/livemap.proto
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

var resources_livemap_livemap_pb = require('../../resources/livemap/livemap_pb.js');
goog.object.extend(proto, resources_livemap_livemap_pb);
goog.exportSymbol('proto.services.livemapper.StreamRequest', null, global);
goog.exportSymbol('proto.services.livemapper.StreamResponse', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.services.livemapper.StreamRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.livemapper.StreamRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.livemapper.StreamRequest.displayName = 'proto.services.livemapper.StreamRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.services.livemapper.StreamResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.services.livemapper.StreamResponse.repeatedFields_, null);
};
goog.inherits(proto.services.livemapper.StreamResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.livemapper.StreamResponse.displayName = 'proto.services.livemapper.StreamResponse';
}



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.services.livemapper.StreamRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.livemapper.StreamRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.livemapper.StreamRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.livemapper.StreamRequest.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.services.livemapper.StreamRequest}
 */
proto.services.livemapper.StreamRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.livemapper.StreamRequest;
  return proto.services.livemapper.StreamRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.livemapper.StreamRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.livemapper.StreamRequest}
 */
proto.services.livemapper.StreamRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.services.livemapper.StreamRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.livemapper.StreamRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.livemapper.StreamRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.livemapper.StreamRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.services.livemapper.StreamResponse.repeatedFields_ = [1,2];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.services.livemapper.StreamResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.livemapper.StreamResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.livemapper.StreamResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.livemapper.StreamResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    dispatchesList: jspb.Message.toObjectList(msg.getDispatchesList(),
    resources_livemap_livemap_pb.GenericMarker.toObject, includeInstance),
    usersList: jspb.Message.toObjectList(msg.getUsersList(),
    resources_livemap_livemap_pb.UserMarker.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.services.livemapper.StreamResponse}
 */
proto.services.livemapper.StreamResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.livemapper.StreamResponse;
  return proto.services.livemapper.StreamResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.livemapper.StreamResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.livemapper.StreamResponse}
 */
proto.services.livemapper.StreamResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_livemap_livemap_pb.GenericMarker;
      reader.readMessage(value,resources_livemap_livemap_pb.GenericMarker.deserializeBinaryFromReader);
      msg.addDispatches(value);
      break;
    case 2:
      var value = new resources_livemap_livemap_pb.UserMarker;
      reader.readMessage(value,resources_livemap_livemap_pb.UserMarker.deserializeBinaryFromReader);
      msg.addUsers(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.services.livemapper.StreamResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.livemapper.StreamResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.livemapper.StreamResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.livemapper.StreamResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDispatchesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      resources_livemap_livemap_pb.GenericMarker.serializeBinaryToWriter
    );
  }
  f = message.getUsersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      resources_livemap_livemap_pb.UserMarker.serializeBinaryToWriter
    );
  }
};


/**
 * repeated resources.livemap.GenericMarker dispatches = 1;
 * @return {!Array<!proto.resources.livemap.GenericMarker>}
 */
proto.services.livemapper.StreamResponse.prototype.getDispatchesList = function() {
  return /** @type{!Array<!proto.resources.livemap.GenericMarker>} */ (
    jspb.Message.getRepeatedWrapperField(this, resources_livemap_livemap_pb.GenericMarker, 1));
};


/**
 * @param {!Array<!proto.resources.livemap.GenericMarker>} value
 * @return {!proto.services.livemapper.StreamResponse} returns this
*/
proto.services.livemapper.StreamResponse.prototype.setDispatchesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.resources.livemap.GenericMarker=} opt_value
 * @param {number=} opt_index
 * @return {!proto.resources.livemap.GenericMarker}
 */
proto.services.livemapper.StreamResponse.prototype.addDispatches = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.resources.livemap.GenericMarker, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.services.livemapper.StreamResponse} returns this
 */
proto.services.livemapper.StreamResponse.prototype.clearDispatchesList = function() {
  return this.setDispatchesList([]);
};


/**
 * repeated resources.livemap.UserMarker users = 2;
 * @return {!Array<!proto.resources.livemap.UserMarker>}
 */
proto.services.livemapper.StreamResponse.prototype.getUsersList = function() {
  return /** @type{!Array<!proto.resources.livemap.UserMarker>} */ (
    jspb.Message.getRepeatedWrapperField(this, resources_livemap_livemap_pb.UserMarker, 2));
};


/**
 * @param {!Array<!proto.resources.livemap.UserMarker>} value
 * @return {!proto.services.livemapper.StreamResponse} returns this
*/
proto.services.livemapper.StreamResponse.prototype.setUsersList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.resources.livemap.UserMarker=} opt_value
 * @param {number=} opt_index
 * @return {!proto.resources.livemap.UserMarker}
 */
proto.services.livemapper.StreamResponse.prototype.addUsers = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.resources.livemap.UserMarker, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.services.livemapper.StreamResponse} returns this
 */
proto.services.livemapper.StreamResponse.prototype.clearUsersList = function() {
  return this.setUsersList([]);
};


goog.object.extend(exports, proto.services.livemapper);
