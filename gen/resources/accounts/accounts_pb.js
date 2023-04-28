// source: resources/accounts/accounts.proto
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

var resources_timestamp_timestamp_pb = require('../../resources/timestamp/timestamp_pb.js');
goog.object.extend(proto, resources_timestamp_timestamp_pb);
goog.exportSymbol('proto.resources.accounts.Account', null, global);
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
proto.resources.accounts.Account = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.resources.accounts.Account, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.resources.accounts.Account.displayName = 'proto.resources.accounts.Account';
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
proto.resources.accounts.Account.prototype.toObject = function(opt_includeInstance) {
  return proto.resources.accounts.Account.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.resources.accounts.Account} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.resources.accounts.Account.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, 0),
    createdAt: (f = msg.getCreatedAt()) && resources_timestamp_timestamp_pb.Timestamp.toObject(includeInstance, f),
    updatedAt: (f = msg.getUpdatedAt()) && resources_timestamp_timestamp_pb.Timestamp.toObject(includeInstance, f),
    username: jspb.Message.getFieldWithDefault(msg, 4, ""),
    license: jspb.Message.getFieldWithDefault(msg, 5, "")
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
 * @return {!proto.resources.accounts.Account}
 */
proto.resources.accounts.Account.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.resources.accounts.Account;
  return proto.resources.accounts.Account.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.resources.accounts.Account} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.resources.accounts.Account}
 */
proto.resources.accounts.Account.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setId(value);
      break;
    case 2:
      var value = new resources_timestamp_timestamp_pb.Timestamp;
      reader.readMessage(value,resources_timestamp_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setCreatedAt(value);
      break;
    case 3:
      var value = new resources_timestamp_timestamp_pb.Timestamp;
      reader.readMessage(value,resources_timestamp_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setUpdatedAt(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setUsername(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setLicense(value);
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
proto.resources.accounts.Account.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.resources.accounts.Account.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.resources.accounts.Account} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.resources.accounts.Account.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
  f = message.getCreatedAt();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      resources_timestamp_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getUpdatedAt();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      resources_timestamp_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getUsername();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getLicense();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
};


/**
 * optional uint64 id = 1;
 * @return {number}
 */
proto.resources.accounts.Account.prototype.getId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.resources.accounts.Account} returns this
 */
proto.resources.accounts.Account.prototype.setId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional resources.timestamp.Timestamp created_at = 2;
 * @return {?proto.resources.timestamp.Timestamp}
 */
proto.resources.accounts.Account.prototype.getCreatedAt = function() {
  return /** @type{?proto.resources.timestamp.Timestamp} */ (
    jspb.Message.getWrapperField(this, resources_timestamp_timestamp_pb.Timestamp, 2));
};


/**
 * @param {?proto.resources.timestamp.Timestamp|undefined} value
 * @return {!proto.resources.accounts.Account} returns this
*/
proto.resources.accounts.Account.prototype.setCreatedAt = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.resources.accounts.Account} returns this
 */
proto.resources.accounts.Account.prototype.clearCreatedAt = function() {
  return this.setCreatedAt(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.resources.accounts.Account.prototype.hasCreatedAt = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional resources.timestamp.Timestamp updated_at = 3;
 * @return {?proto.resources.timestamp.Timestamp}
 */
proto.resources.accounts.Account.prototype.getUpdatedAt = function() {
  return /** @type{?proto.resources.timestamp.Timestamp} */ (
    jspb.Message.getWrapperField(this, resources_timestamp_timestamp_pb.Timestamp, 3));
};


/**
 * @param {?proto.resources.timestamp.Timestamp|undefined} value
 * @return {!proto.resources.accounts.Account} returns this
*/
proto.resources.accounts.Account.prototype.setUpdatedAt = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.resources.accounts.Account} returns this
 */
proto.resources.accounts.Account.prototype.clearUpdatedAt = function() {
  return this.setUpdatedAt(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.resources.accounts.Account.prototype.hasUpdatedAt = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional string username = 4;
 * @return {string}
 */
proto.resources.accounts.Account.prototype.getUsername = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.resources.accounts.Account} returns this
 */
proto.resources.accounts.Account.prototype.setUsername = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string license = 5;
 * @return {string}
 */
proto.resources.accounts.Account.prototype.getLicense = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.resources.accounts.Account} returns this
 */
proto.resources.accounts.Account.prototype.setLicense = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


goog.object.extend(exports, proto.resources.accounts);
