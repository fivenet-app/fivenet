// source: resources/notifications/notifications.proto
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
goog.exportSymbol('proto.resources.notifications.Notification', null, global);
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
proto.resources.notifications.Notification = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.resources.notifications.Notification, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.resources.notifications.Notification.displayName = 'proto.resources.notifications.Notification';
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
proto.resources.notifications.Notification.prototype.toObject = function(opt_includeInstance) {
  return proto.resources.notifications.Notification.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.resources.notifications.Notification} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.resources.notifications.Notification.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, 0),
    createdAt: (f = msg.getCreatedAt()) && resources_timestamp_timestamp_pb.Timestamp.toObject(includeInstance, f),
    readAt: (f = msg.getReadAt()) && resources_timestamp_timestamp_pb.Timestamp.toObject(includeInstance, f),
    userId: jspb.Message.getFieldWithDefault(msg, 4, 0),
    title: jspb.Message.getFieldWithDefault(msg, 5, ""),
    type: jspb.Message.getFieldWithDefault(msg, 6, ""),
    content: jspb.Message.getFieldWithDefault(msg, 7, ""),
    data: jspb.Message.getFieldWithDefault(msg, 8, "")
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
 * @return {!proto.resources.notifications.Notification}
 */
proto.resources.notifications.Notification.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.resources.notifications.Notification;
  return proto.resources.notifications.Notification.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.resources.notifications.Notification} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.resources.notifications.Notification}
 */
proto.resources.notifications.Notification.deserializeBinaryFromReader = function(msg, reader) {
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
      msg.setReadAt(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setUserId(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setTitle(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setType(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setContent(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setData(value);
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
proto.resources.notifications.Notification.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.resources.notifications.Notification.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.resources.notifications.Notification} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.resources.notifications.Notification.serializeBinaryToWriter = function(message, writer) {
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
  f = message.getReadAt();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      resources_timestamp_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getUserId();
  if (f !== 0) {
    writer.writeInt32(
      4,
      f
    );
  }
  f = message.getTitle();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getType();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getContent();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 8));
  if (f != null) {
    writer.writeString(
      8,
      f
    );
  }
};


/**
 * optional uint64 id = 1;
 * @return {number}
 */
proto.resources.notifications.Notification.prototype.getId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.resources.notifications.Notification} returns this
 */
proto.resources.notifications.Notification.prototype.setId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional resources.timestamp.Timestamp created_at = 2;
 * @return {?proto.resources.timestamp.Timestamp}
 */
proto.resources.notifications.Notification.prototype.getCreatedAt = function() {
  return /** @type{?proto.resources.timestamp.Timestamp} */ (
    jspb.Message.getWrapperField(this, resources_timestamp_timestamp_pb.Timestamp, 2));
};


/**
 * @param {?proto.resources.timestamp.Timestamp|undefined} value
 * @return {!proto.resources.notifications.Notification} returns this
*/
proto.resources.notifications.Notification.prototype.setCreatedAt = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.resources.notifications.Notification} returns this
 */
proto.resources.notifications.Notification.prototype.clearCreatedAt = function() {
  return this.setCreatedAt(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.resources.notifications.Notification.prototype.hasCreatedAt = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional resources.timestamp.Timestamp read_at = 3;
 * @return {?proto.resources.timestamp.Timestamp}
 */
proto.resources.notifications.Notification.prototype.getReadAt = function() {
  return /** @type{?proto.resources.timestamp.Timestamp} */ (
    jspb.Message.getWrapperField(this, resources_timestamp_timestamp_pb.Timestamp, 3));
};


/**
 * @param {?proto.resources.timestamp.Timestamp|undefined} value
 * @return {!proto.resources.notifications.Notification} returns this
*/
proto.resources.notifications.Notification.prototype.setReadAt = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.resources.notifications.Notification} returns this
 */
proto.resources.notifications.Notification.prototype.clearReadAt = function() {
  return this.setReadAt(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.resources.notifications.Notification.prototype.hasReadAt = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional int32 user_id = 4;
 * @return {number}
 */
proto.resources.notifications.Notification.prototype.getUserId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.resources.notifications.Notification} returns this
 */
proto.resources.notifications.Notification.prototype.setUserId = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional string title = 5;
 * @return {string}
 */
proto.resources.notifications.Notification.prototype.getTitle = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.resources.notifications.Notification} returns this
 */
proto.resources.notifications.Notification.prototype.setTitle = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string type = 6;
 * @return {string}
 */
proto.resources.notifications.Notification.prototype.getType = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/**
 * @param {string} value
 * @return {!proto.resources.notifications.Notification} returns this
 */
proto.resources.notifications.Notification.prototype.setType = function(value) {
  return jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional string content = 7;
 * @return {string}
 */
proto.resources.notifications.Notification.prototype.getContent = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/**
 * @param {string} value
 * @return {!proto.resources.notifications.Notification} returns this
 */
proto.resources.notifications.Notification.prototype.setContent = function(value) {
  return jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional string data = 8;
 * @return {string}
 */
proto.resources.notifications.Notification.prototype.getData = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/**
 * @param {string} value
 * @return {!proto.resources.notifications.Notification} returns this
 */
proto.resources.notifications.Notification.prototype.setData = function(value) {
  return jspb.Message.setField(this, 8, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.resources.notifications.Notification} returns this
 */
proto.resources.notifications.Notification.prototype.clearData = function() {
  return jspb.Message.setField(this, 8, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.resources.notifications.Notification.prototype.hasData = function() {
  return jspb.Message.getField(this, 8) != null;
};


goog.object.extend(exports, proto.resources.notifications);
