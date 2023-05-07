// source: services/docstore/docstore.proto
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

var resources_common_database_database_pb = require('../../resources/common/database/database_pb.js');
goog.object.extend(proto, resources_common_database_database_pb);
var resources_documents_category_pb = require('../../resources/documents/category_pb.js');
goog.object.extend(proto, resources_documents_category_pb);
var resources_documents_documents_pb = require('../../resources/documents/documents_pb.js');
goog.object.extend(proto, resources_documents_documents_pb);
var resources_documents_templates_pb = require('../../resources/documents/templates_pb.js');
goog.object.extend(proto, resources_documents_templates_pb);
goog.exportSymbol('proto.services.docstore.ACCESS_LEVEL_UPDATE_MODE', null, global);
goog.exportSymbol('proto.services.docstore.AddDocumentReferenceRequest', null, global);
goog.exportSymbol('proto.services.docstore.AddDocumentReferenceResponse', null, global);
goog.exportSymbol('proto.services.docstore.AddDocumentRelationRequest', null, global);
goog.exportSymbol('proto.services.docstore.AddDocumentRelationResponse', null, global);
goog.exportSymbol('proto.services.docstore.CreateDocumentCategoryRequest', null, global);
goog.exportSymbol('proto.services.docstore.CreateDocumentCategoryResponse', null, global);
goog.exportSymbol('proto.services.docstore.CreateDocumentRequest', null, global);
goog.exportSymbol('proto.services.docstore.CreateDocumentResponse', null, global);
goog.exportSymbol('proto.services.docstore.CreateTemplateRequest', null, global);
goog.exportSymbol('proto.services.docstore.CreateTemplateResponse', null, global);
goog.exportSymbol('proto.services.docstore.DeleteDocumentCategoryRequest', null, global);
goog.exportSymbol('proto.services.docstore.DeleteDocumentCategoryResponse', null, global);
goog.exportSymbol('proto.services.docstore.DeleteDocumentCommentRequest', null, global);
goog.exportSymbol('proto.services.docstore.DeleteDocumentCommentResponse', null, global);
goog.exportSymbol('proto.services.docstore.DeleteDocumentRequest', null, global);
goog.exportSymbol('proto.services.docstore.DeleteDocumentResponse', null, global);
goog.exportSymbol('proto.services.docstore.DeleteTemplateRequest', null, global);
goog.exportSymbol('proto.services.docstore.DeleteTemplateResponse', null, global);
goog.exportSymbol('proto.services.docstore.EditDocumentCommentRequest', null, global);
goog.exportSymbol('proto.services.docstore.EditDocumentCommentResponse', null, global);
goog.exportSymbol('proto.services.docstore.GetDocumentAccessRequest', null, global);
goog.exportSymbol('proto.services.docstore.GetDocumentAccessResponse', null, global);
goog.exportSymbol('proto.services.docstore.GetDocumentCommentsRequest', null, global);
goog.exportSymbol('proto.services.docstore.GetDocumentCommentsResponse', null, global);
goog.exportSymbol('proto.services.docstore.GetDocumentReferencesRequest', null, global);
goog.exportSymbol('proto.services.docstore.GetDocumentReferencesResponse', null, global);
goog.exportSymbol('proto.services.docstore.GetDocumentRelationsRequest', null, global);
goog.exportSymbol('proto.services.docstore.GetDocumentRelationsResponse', null, global);
goog.exportSymbol('proto.services.docstore.GetDocumentRequest', null, global);
goog.exportSymbol('proto.services.docstore.GetDocumentResponse', null, global);
goog.exportSymbol('proto.services.docstore.GetTemplateRequest', null, global);
goog.exportSymbol('proto.services.docstore.GetTemplateResponse', null, global);
goog.exportSymbol('proto.services.docstore.ListDocumentCategoriesRequest', null, global);
goog.exportSymbol('proto.services.docstore.ListDocumentCategoriesResponse', null, global);
goog.exportSymbol('proto.services.docstore.ListDocumentsRequest', null, global);
goog.exportSymbol('proto.services.docstore.ListDocumentsResponse', null, global);
goog.exportSymbol('proto.services.docstore.ListTemplatesRequest', null, global);
goog.exportSymbol('proto.services.docstore.ListTemplatesResponse', null, global);
goog.exportSymbol('proto.services.docstore.ListUserDocumentsRequest', null, global);
goog.exportSymbol('proto.services.docstore.ListUserDocumentsResponse', null, global);
goog.exportSymbol('proto.services.docstore.PostDocumentCommentRequest', null, global);
goog.exportSymbol('proto.services.docstore.PostDocumentCommentResponse', null, global);
goog.exportSymbol('proto.services.docstore.RemoveDocumentReferenceRequest', null, global);
goog.exportSymbol('proto.services.docstore.RemoveDocumentReferenceResponse', null, global);
goog.exportSymbol('proto.services.docstore.RemoveDocumentRelationRequest', null, global);
goog.exportSymbol('proto.services.docstore.RemoveDocumentRelationResponse', null, global);
goog.exportSymbol('proto.services.docstore.SetDocumentAccessRequest', null, global);
goog.exportSymbol('proto.services.docstore.SetDocumentAccessResponse', null, global);
goog.exportSymbol('proto.services.docstore.UpdateDocumentCategoryRequest', null, global);
goog.exportSymbol('proto.services.docstore.UpdateDocumentCategoryResponse', null, global);
goog.exportSymbol('proto.services.docstore.UpdateDocumentRequest', null, global);
goog.exportSymbol('proto.services.docstore.UpdateDocumentResponse', null, global);
goog.exportSymbol('proto.services.docstore.UpdateTemplateRequest', null, global);
goog.exportSymbol('proto.services.docstore.UpdateTemplateResponse', null, global);
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
proto.services.docstore.ListTemplatesRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.ListTemplatesRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.ListTemplatesRequest.displayName = 'proto.services.docstore.ListTemplatesRequest';
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
proto.services.docstore.ListTemplatesResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.services.docstore.ListTemplatesResponse.repeatedFields_, null);
};
goog.inherits(proto.services.docstore.ListTemplatesResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.ListTemplatesResponse.displayName = 'proto.services.docstore.ListTemplatesResponse';
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
proto.services.docstore.GetTemplateRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.GetTemplateRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.GetTemplateRequest.displayName = 'proto.services.docstore.GetTemplateRequest';
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
proto.services.docstore.GetTemplateResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.GetTemplateResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.GetTemplateResponse.displayName = 'proto.services.docstore.GetTemplateResponse';
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
proto.services.docstore.CreateTemplateRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.CreateTemplateRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.CreateTemplateRequest.displayName = 'proto.services.docstore.CreateTemplateRequest';
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
proto.services.docstore.CreateTemplateResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.CreateTemplateResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.CreateTemplateResponse.displayName = 'proto.services.docstore.CreateTemplateResponse';
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
proto.services.docstore.UpdateTemplateRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.UpdateTemplateRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.UpdateTemplateRequest.displayName = 'proto.services.docstore.UpdateTemplateRequest';
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
proto.services.docstore.UpdateTemplateResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.UpdateTemplateResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.UpdateTemplateResponse.displayName = 'proto.services.docstore.UpdateTemplateResponse';
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
proto.services.docstore.DeleteTemplateRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.DeleteTemplateRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.DeleteTemplateRequest.displayName = 'proto.services.docstore.DeleteTemplateRequest';
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
proto.services.docstore.DeleteTemplateResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.DeleteTemplateResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.DeleteTemplateResponse.displayName = 'proto.services.docstore.DeleteTemplateResponse';
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
proto.services.docstore.ListDocumentsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.services.docstore.ListDocumentsRequest.repeatedFields_, null);
};
goog.inherits(proto.services.docstore.ListDocumentsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.ListDocumentsRequest.displayName = 'proto.services.docstore.ListDocumentsRequest';
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
proto.services.docstore.ListDocumentsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.services.docstore.ListDocumentsResponse.repeatedFields_, null);
};
goog.inherits(proto.services.docstore.ListDocumentsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.ListDocumentsResponse.displayName = 'proto.services.docstore.ListDocumentsResponse';
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
proto.services.docstore.GetDocumentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.GetDocumentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.GetDocumentRequest.displayName = 'proto.services.docstore.GetDocumentRequest';
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
proto.services.docstore.GetDocumentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.GetDocumentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.GetDocumentResponse.displayName = 'proto.services.docstore.GetDocumentResponse';
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
proto.services.docstore.GetDocumentReferencesRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.GetDocumentReferencesRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.GetDocumentReferencesRequest.displayName = 'proto.services.docstore.GetDocumentReferencesRequest';
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
proto.services.docstore.GetDocumentReferencesResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.services.docstore.GetDocumentReferencesResponse.repeatedFields_, null);
};
goog.inherits(proto.services.docstore.GetDocumentReferencesResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.GetDocumentReferencesResponse.displayName = 'proto.services.docstore.GetDocumentReferencesResponse';
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
proto.services.docstore.GetDocumentRelationsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.GetDocumentRelationsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.GetDocumentRelationsRequest.displayName = 'proto.services.docstore.GetDocumentRelationsRequest';
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
proto.services.docstore.GetDocumentRelationsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.services.docstore.GetDocumentRelationsResponse.repeatedFields_, null);
};
goog.inherits(proto.services.docstore.GetDocumentRelationsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.GetDocumentRelationsResponse.displayName = 'proto.services.docstore.GetDocumentRelationsResponse';
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
proto.services.docstore.AddDocumentReferenceRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.AddDocumentReferenceRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.AddDocumentReferenceRequest.displayName = 'proto.services.docstore.AddDocumentReferenceRequest';
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
proto.services.docstore.AddDocumentReferenceResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.AddDocumentReferenceResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.AddDocumentReferenceResponse.displayName = 'proto.services.docstore.AddDocumentReferenceResponse';
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
proto.services.docstore.RemoveDocumentReferenceRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.RemoveDocumentReferenceRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.RemoveDocumentReferenceRequest.displayName = 'proto.services.docstore.RemoveDocumentReferenceRequest';
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
proto.services.docstore.RemoveDocumentReferenceResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.RemoveDocumentReferenceResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.RemoveDocumentReferenceResponse.displayName = 'proto.services.docstore.RemoveDocumentReferenceResponse';
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
proto.services.docstore.AddDocumentRelationRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.AddDocumentRelationRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.AddDocumentRelationRequest.displayName = 'proto.services.docstore.AddDocumentRelationRequest';
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
proto.services.docstore.AddDocumentRelationResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.AddDocumentRelationResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.AddDocumentRelationResponse.displayName = 'proto.services.docstore.AddDocumentRelationResponse';
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
proto.services.docstore.RemoveDocumentRelationRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.RemoveDocumentRelationRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.RemoveDocumentRelationRequest.displayName = 'proto.services.docstore.RemoveDocumentRelationRequest';
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
proto.services.docstore.RemoveDocumentRelationResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.RemoveDocumentRelationResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.RemoveDocumentRelationResponse.displayName = 'proto.services.docstore.RemoveDocumentRelationResponse';
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
proto.services.docstore.GetDocumentCommentsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.GetDocumentCommentsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.GetDocumentCommentsRequest.displayName = 'proto.services.docstore.GetDocumentCommentsRequest';
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
proto.services.docstore.GetDocumentCommentsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.services.docstore.GetDocumentCommentsResponse.repeatedFields_, null);
};
goog.inherits(proto.services.docstore.GetDocumentCommentsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.GetDocumentCommentsResponse.displayName = 'proto.services.docstore.GetDocumentCommentsResponse';
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
proto.services.docstore.PostDocumentCommentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.PostDocumentCommentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.PostDocumentCommentRequest.displayName = 'proto.services.docstore.PostDocumentCommentRequest';
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
proto.services.docstore.PostDocumentCommentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.PostDocumentCommentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.PostDocumentCommentResponse.displayName = 'proto.services.docstore.PostDocumentCommentResponse';
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
proto.services.docstore.EditDocumentCommentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.EditDocumentCommentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.EditDocumentCommentRequest.displayName = 'proto.services.docstore.EditDocumentCommentRequest';
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
proto.services.docstore.EditDocumentCommentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.EditDocumentCommentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.EditDocumentCommentResponse.displayName = 'proto.services.docstore.EditDocumentCommentResponse';
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
proto.services.docstore.DeleteDocumentCommentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.DeleteDocumentCommentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.DeleteDocumentCommentRequest.displayName = 'proto.services.docstore.DeleteDocumentCommentRequest';
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
proto.services.docstore.DeleteDocumentCommentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.DeleteDocumentCommentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.DeleteDocumentCommentResponse.displayName = 'proto.services.docstore.DeleteDocumentCommentResponse';
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
proto.services.docstore.CreateDocumentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.CreateDocumentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.CreateDocumentRequest.displayName = 'proto.services.docstore.CreateDocumentRequest';
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
proto.services.docstore.CreateDocumentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.CreateDocumentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.CreateDocumentResponse.displayName = 'proto.services.docstore.CreateDocumentResponse';
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
proto.services.docstore.UpdateDocumentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.UpdateDocumentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.UpdateDocumentRequest.displayName = 'proto.services.docstore.UpdateDocumentRequest';
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
proto.services.docstore.UpdateDocumentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.UpdateDocumentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.UpdateDocumentResponse.displayName = 'proto.services.docstore.UpdateDocumentResponse';
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
proto.services.docstore.DeleteDocumentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.DeleteDocumentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.DeleteDocumentRequest.displayName = 'proto.services.docstore.DeleteDocumentRequest';
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
proto.services.docstore.DeleteDocumentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.DeleteDocumentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.DeleteDocumentResponse.displayName = 'proto.services.docstore.DeleteDocumentResponse';
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
proto.services.docstore.GetDocumentAccessRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.GetDocumentAccessRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.GetDocumentAccessRequest.displayName = 'proto.services.docstore.GetDocumentAccessRequest';
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
proto.services.docstore.GetDocumentAccessResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.GetDocumentAccessResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.GetDocumentAccessResponse.displayName = 'proto.services.docstore.GetDocumentAccessResponse';
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
proto.services.docstore.SetDocumentAccessRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.SetDocumentAccessRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.SetDocumentAccessRequest.displayName = 'proto.services.docstore.SetDocumentAccessRequest';
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
proto.services.docstore.SetDocumentAccessResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.SetDocumentAccessResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.SetDocumentAccessResponse.displayName = 'proto.services.docstore.SetDocumentAccessResponse';
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
proto.services.docstore.ListUserDocumentsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.services.docstore.ListUserDocumentsRequest.repeatedFields_, null);
};
goog.inherits(proto.services.docstore.ListUserDocumentsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.ListUserDocumentsRequest.displayName = 'proto.services.docstore.ListUserDocumentsRequest';
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
proto.services.docstore.ListUserDocumentsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.services.docstore.ListUserDocumentsResponse.repeatedFields_, null);
};
goog.inherits(proto.services.docstore.ListUserDocumentsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.ListUserDocumentsResponse.displayName = 'proto.services.docstore.ListUserDocumentsResponse';
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
proto.services.docstore.ListDocumentCategoriesRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.ListDocumentCategoriesRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.ListDocumentCategoriesRequest.displayName = 'proto.services.docstore.ListDocumentCategoriesRequest';
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
proto.services.docstore.ListDocumentCategoriesResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.services.docstore.ListDocumentCategoriesResponse.repeatedFields_, null);
};
goog.inherits(proto.services.docstore.ListDocumentCategoriesResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.ListDocumentCategoriesResponse.displayName = 'proto.services.docstore.ListDocumentCategoriesResponse';
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
proto.services.docstore.CreateDocumentCategoryRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.CreateDocumentCategoryRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.CreateDocumentCategoryRequest.displayName = 'proto.services.docstore.CreateDocumentCategoryRequest';
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
proto.services.docstore.CreateDocumentCategoryResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.CreateDocumentCategoryResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.CreateDocumentCategoryResponse.displayName = 'proto.services.docstore.CreateDocumentCategoryResponse';
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
proto.services.docstore.UpdateDocumentCategoryRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.UpdateDocumentCategoryRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.UpdateDocumentCategoryRequest.displayName = 'proto.services.docstore.UpdateDocumentCategoryRequest';
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
proto.services.docstore.UpdateDocumentCategoryResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.UpdateDocumentCategoryResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.UpdateDocumentCategoryResponse.displayName = 'proto.services.docstore.UpdateDocumentCategoryResponse';
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
proto.services.docstore.DeleteDocumentCategoryRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.services.docstore.DeleteDocumentCategoryRequest.repeatedFields_, null);
};
goog.inherits(proto.services.docstore.DeleteDocumentCategoryRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.DeleteDocumentCategoryRequest.displayName = 'proto.services.docstore.DeleteDocumentCategoryRequest';
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
proto.services.docstore.DeleteDocumentCategoryResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.services.docstore.DeleteDocumentCategoryResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.services.docstore.DeleteDocumentCategoryResponse.displayName = 'proto.services.docstore.DeleteDocumentCategoryResponse';
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
proto.services.docstore.ListTemplatesRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.ListTemplatesRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.ListTemplatesRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListTemplatesRequest.toObject = function(includeInstance, msg) {
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
 * @return {!proto.services.docstore.ListTemplatesRequest}
 */
proto.services.docstore.ListTemplatesRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.ListTemplatesRequest;
  return proto.services.docstore.ListTemplatesRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.ListTemplatesRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.ListTemplatesRequest}
 */
proto.services.docstore.ListTemplatesRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.ListTemplatesRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.ListTemplatesRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.ListTemplatesRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListTemplatesRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.services.docstore.ListTemplatesResponse.repeatedFields_ = [1];



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
proto.services.docstore.ListTemplatesResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.ListTemplatesResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.ListTemplatesResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListTemplatesResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    templatesList: jspb.Message.toObjectList(msg.getTemplatesList(),
    resources_documents_templates_pb.TemplateShort.toObject, includeInstance)
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
 * @return {!proto.services.docstore.ListTemplatesResponse}
 */
proto.services.docstore.ListTemplatesResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.ListTemplatesResponse;
  return proto.services.docstore.ListTemplatesResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.ListTemplatesResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.ListTemplatesResponse}
 */
proto.services.docstore.ListTemplatesResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_templates_pb.TemplateShort;
      reader.readMessage(value,resources_documents_templates_pb.TemplateShort.deserializeBinaryFromReader);
      msg.addTemplates(value);
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
proto.services.docstore.ListTemplatesResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.ListTemplatesResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.ListTemplatesResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListTemplatesResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTemplatesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      resources_documents_templates_pb.TemplateShort.serializeBinaryToWriter
    );
  }
};


/**
 * repeated resources.documents.TemplateShort templates = 1;
 * @return {!Array<!proto.resources.documents.TemplateShort>}
 */
proto.services.docstore.ListTemplatesResponse.prototype.getTemplatesList = function() {
  return /** @type{!Array<!proto.resources.documents.TemplateShort>} */ (
    jspb.Message.getRepeatedWrapperField(this, resources_documents_templates_pb.TemplateShort, 1));
};


/**
 * @param {!Array<!proto.resources.documents.TemplateShort>} value
 * @return {!proto.services.docstore.ListTemplatesResponse} returns this
*/
proto.services.docstore.ListTemplatesResponse.prototype.setTemplatesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.resources.documents.TemplateShort=} opt_value
 * @param {number=} opt_index
 * @return {!proto.resources.documents.TemplateShort}
 */
proto.services.docstore.ListTemplatesResponse.prototype.addTemplates = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.resources.documents.TemplateShort, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.services.docstore.ListTemplatesResponse} returns this
 */
proto.services.docstore.ListTemplatesResponse.prototype.clearTemplatesList = function() {
  return this.setTemplatesList([]);
};





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
proto.services.docstore.GetTemplateRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.GetTemplateRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.GetTemplateRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetTemplateRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    templateId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    data: jspb.Message.getFieldWithDefault(msg, 2, ""),
    render: jspb.Message.getBooleanFieldWithDefault(msg, 3, false)
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
 * @return {!proto.services.docstore.GetTemplateRequest}
 */
proto.services.docstore.GetTemplateRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.GetTemplateRequest;
  return proto.services.docstore.GetTemplateRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.GetTemplateRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.GetTemplateRequest}
 */
proto.services.docstore.GetTemplateRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setTemplateId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setData(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setRender(value);
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
proto.services.docstore.GetTemplateRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.GetTemplateRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.GetTemplateRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetTemplateRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTemplateId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
  f = message.getData();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = /** @type {boolean} */ (jspb.Message.getField(message, 3));
  if (f != null) {
    writer.writeBool(
      3,
      f
    );
  }
};


/**
 * optional uint64 template_id = 1;
 * @return {number}
 */
proto.services.docstore.GetTemplateRequest.prototype.getTemplateId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.GetTemplateRequest} returns this
 */
proto.services.docstore.GetTemplateRequest.prototype.setTemplateId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string data = 2;
 * @return {string}
 */
proto.services.docstore.GetTemplateRequest.prototype.getData = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.services.docstore.GetTemplateRequest} returns this
 */
proto.services.docstore.GetTemplateRequest.prototype.setData = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bool render = 3;
 * @return {boolean}
 */
proto.services.docstore.GetTemplateRequest.prototype.getRender = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 3, false));
};


/**
 * @param {boolean} value
 * @return {!proto.services.docstore.GetTemplateRequest} returns this
 */
proto.services.docstore.GetTemplateRequest.prototype.setRender = function(value) {
  return jspb.Message.setField(this, 3, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.services.docstore.GetTemplateRequest} returns this
 */
proto.services.docstore.GetTemplateRequest.prototype.clearRender = function() {
  return jspb.Message.setField(this, 3, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.GetTemplateRequest.prototype.hasRender = function() {
  return jspb.Message.getField(this, 3) != null;
};





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
proto.services.docstore.GetTemplateResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.GetTemplateResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.GetTemplateResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetTemplateResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    template: (f = msg.getTemplate()) && resources_documents_templates_pb.Template.toObject(includeInstance, f),
    rendered: jspb.Message.getBooleanFieldWithDefault(msg, 2, false)
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
 * @return {!proto.services.docstore.GetTemplateResponse}
 */
proto.services.docstore.GetTemplateResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.GetTemplateResponse;
  return proto.services.docstore.GetTemplateResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.GetTemplateResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.GetTemplateResponse}
 */
proto.services.docstore.GetTemplateResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_templates_pb.Template;
      reader.readMessage(value,resources_documents_templates_pb.Template.deserializeBinaryFromReader);
      msg.setTemplate(value);
      break;
    case 2:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setRendered(value);
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
proto.services.docstore.GetTemplateResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.GetTemplateResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.GetTemplateResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetTemplateResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTemplate();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_documents_templates_pb.Template.serializeBinaryToWriter
    );
  }
  f = message.getRendered();
  if (f) {
    writer.writeBool(
      2,
      f
    );
  }
};


/**
 * optional resources.documents.Template template = 1;
 * @return {?proto.resources.documents.Template}
 */
proto.services.docstore.GetTemplateResponse.prototype.getTemplate = function() {
  return /** @type{?proto.resources.documents.Template} */ (
    jspb.Message.getWrapperField(this, resources_documents_templates_pb.Template, 1));
};


/**
 * @param {?proto.resources.documents.Template|undefined} value
 * @return {!proto.services.docstore.GetTemplateResponse} returns this
*/
proto.services.docstore.GetTemplateResponse.prototype.setTemplate = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.GetTemplateResponse} returns this
 */
proto.services.docstore.GetTemplateResponse.prototype.clearTemplate = function() {
  return this.setTemplate(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.GetTemplateResponse.prototype.hasTemplate = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional bool rendered = 2;
 * @return {boolean}
 */
proto.services.docstore.GetTemplateResponse.prototype.getRendered = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 2, false));
};


/**
 * @param {boolean} value
 * @return {!proto.services.docstore.GetTemplateResponse} returns this
 */
proto.services.docstore.GetTemplateResponse.prototype.setRendered = function(value) {
  return jspb.Message.setProto3BooleanField(this, 2, value);
};





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
proto.services.docstore.CreateTemplateRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.CreateTemplateRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.CreateTemplateRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.CreateTemplateRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    template: (f = msg.getTemplate()) && resources_documents_templates_pb.Template.toObject(includeInstance, f)
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
 * @return {!proto.services.docstore.CreateTemplateRequest}
 */
proto.services.docstore.CreateTemplateRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.CreateTemplateRequest;
  return proto.services.docstore.CreateTemplateRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.CreateTemplateRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.CreateTemplateRequest}
 */
proto.services.docstore.CreateTemplateRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_templates_pb.Template;
      reader.readMessage(value,resources_documents_templates_pb.Template.deserializeBinaryFromReader);
      msg.setTemplate(value);
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
proto.services.docstore.CreateTemplateRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.CreateTemplateRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.CreateTemplateRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.CreateTemplateRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTemplate();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_documents_templates_pb.Template.serializeBinaryToWriter
    );
  }
};


/**
 * optional resources.documents.Template template = 1;
 * @return {?proto.resources.documents.Template}
 */
proto.services.docstore.CreateTemplateRequest.prototype.getTemplate = function() {
  return /** @type{?proto.resources.documents.Template} */ (
    jspb.Message.getWrapperField(this, resources_documents_templates_pb.Template, 1));
};


/**
 * @param {?proto.resources.documents.Template|undefined} value
 * @return {!proto.services.docstore.CreateTemplateRequest} returns this
*/
proto.services.docstore.CreateTemplateRequest.prototype.setTemplate = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.CreateTemplateRequest} returns this
 */
proto.services.docstore.CreateTemplateRequest.prototype.clearTemplate = function() {
  return this.setTemplate(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.CreateTemplateRequest.prototype.hasTemplate = function() {
  return jspb.Message.getField(this, 1) != null;
};





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
proto.services.docstore.CreateTemplateResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.CreateTemplateResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.CreateTemplateResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.CreateTemplateResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.CreateTemplateResponse}
 */
proto.services.docstore.CreateTemplateResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.CreateTemplateResponse;
  return proto.services.docstore.CreateTemplateResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.CreateTemplateResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.CreateTemplateResponse}
 */
proto.services.docstore.CreateTemplateResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.CreateTemplateResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.CreateTemplateResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.CreateTemplateResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.CreateTemplateResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 id = 1;
 * @return {number}
 */
proto.services.docstore.CreateTemplateResponse.prototype.getId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.CreateTemplateResponse} returns this
 */
proto.services.docstore.CreateTemplateResponse.prototype.setId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.UpdateTemplateRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.UpdateTemplateRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.UpdateTemplateRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.UpdateTemplateRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    template: (f = msg.getTemplate()) && resources_documents_templates_pb.Template.toObject(includeInstance, f)
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
 * @return {!proto.services.docstore.UpdateTemplateRequest}
 */
proto.services.docstore.UpdateTemplateRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.UpdateTemplateRequest;
  return proto.services.docstore.UpdateTemplateRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.UpdateTemplateRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.UpdateTemplateRequest}
 */
proto.services.docstore.UpdateTemplateRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_templates_pb.Template;
      reader.readMessage(value,resources_documents_templates_pb.Template.deserializeBinaryFromReader);
      msg.setTemplate(value);
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
proto.services.docstore.UpdateTemplateRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.UpdateTemplateRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.UpdateTemplateRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.UpdateTemplateRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTemplate();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_documents_templates_pb.Template.serializeBinaryToWriter
    );
  }
};


/**
 * optional resources.documents.Template template = 1;
 * @return {?proto.resources.documents.Template}
 */
proto.services.docstore.UpdateTemplateRequest.prototype.getTemplate = function() {
  return /** @type{?proto.resources.documents.Template} */ (
    jspb.Message.getWrapperField(this, resources_documents_templates_pb.Template, 1));
};


/**
 * @param {?proto.resources.documents.Template|undefined} value
 * @return {!proto.services.docstore.UpdateTemplateRequest} returns this
*/
proto.services.docstore.UpdateTemplateRequest.prototype.setTemplate = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.UpdateTemplateRequest} returns this
 */
proto.services.docstore.UpdateTemplateRequest.prototype.clearTemplate = function() {
  return this.setTemplate(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.UpdateTemplateRequest.prototype.hasTemplate = function() {
  return jspb.Message.getField(this, 1) != null;
};





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
proto.services.docstore.UpdateTemplateResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.UpdateTemplateResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.UpdateTemplateResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.UpdateTemplateResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.UpdateTemplateResponse}
 */
proto.services.docstore.UpdateTemplateResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.UpdateTemplateResponse;
  return proto.services.docstore.UpdateTemplateResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.UpdateTemplateResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.UpdateTemplateResponse}
 */
proto.services.docstore.UpdateTemplateResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.UpdateTemplateResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.UpdateTemplateResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.UpdateTemplateResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.UpdateTemplateResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 id = 1;
 * @return {number}
 */
proto.services.docstore.UpdateTemplateResponse.prototype.getId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.UpdateTemplateResponse} returns this
 */
proto.services.docstore.UpdateTemplateResponse.prototype.setId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.DeleteTemplateRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.DeleteTemplateRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.DeleteTemplateRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteTemplateRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.DeleteTemplateRequest}
 */
proto.services.docstore.DeleteTemplateRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.DeleteTemplateRequest;
  return proto.services.docstore.DeleteTemplateRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.DeleteTemplateRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.DeleteTemplateRequest}
 */
proto.services.docstore.DeleteTemplateRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.DeleteTemplateRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.DeleteTemplateRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.DeleteTemplateRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteTemplateRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 id = 1;
 * @return {number}
 */
proto.services.docstore.DeleteTemplateRequest.prototype.getId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.DeleteTemplateRequest} returns this
 */
proto.services.docstore.DeleteTemplateRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.DeleteTemplateResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.DeleteTemplateResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.DeleteTemplateResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteTemplateResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.services.docstore.DeleteTemplateResponse}
 */
proto.services.docstore.DeleteTemplateResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.DeleteTemplateResponse;
  return proto.services.docstore.DeleteTemplateResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.DeleteTemplateResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.DeleteTemplateResponse}
 */
proto.services.docstore.DeleteTemplateResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.DeleteTemplateResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.DeleteTemplateResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.DeleteTemplateResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteTemplateResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.services.docstore.ListDocumentsRequest.repeatedFields_ = [2,4];



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
proto.services.docstore.ListDocumentsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.ListDocumentsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.ListDocumentsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListDocumentsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pagination: (f = msg.getPagination()) && resources_common_database_database_pb.PaginationRequest.toObject(includeInstance, f),
    orderbyList: jspb.Message.toObjectList(msg.getOrderbyList(),
    resources_common_database_database_pb.OrderBy.toObject, includeInstance),
    search: jspb.Message.getFieldWithDefault(msg, 3, ""),
    categoryIdsList: (f = jspb.Message.getRepeatedField(msg, 4)) == null ? undefined : f
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
 * @return {!proto.services.docstore.ListDocumentsRequest}
 */
proto.services.docstore.ListDocumentsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.ListDocumentsRequest;
  return proto.services.docstore.ListDocumentsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.ListDocumentsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.ListDocumentsRequest}
 */
proto.services.docstore.ListDocumentsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_common_database_database_pb.PaginationRequest;
      reader.readMessage(value,resources_common_database_database_pb.PaginationRequest.deserializeBinaryFromReader);
      msg.setPagination(value);
      break;
    case 2:
      var value = new resources_common_database_database_pb.OrderBy;
      reader.readMessage(value,resources_common_database_database_pb.OrderBy.deserializeBinaryFromReader);
      msg.addOrderby(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setSearch(value);
      break;
    case 4:
      var values = /** @type {!Array<number>} */ (reader.isDelimited() ? reader.readPackedUint64() : [reader.readUint64()]);
      for (var i = 0; i < values.length; i++) {
        msg.addCategoryIds(values[i]);
      }
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
proto.services.docstore.ListDocumentsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.ListDocumentsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.ListDocumentsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListDocumentsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPagination();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_common_database_database_pb.PaginationRequest.serializeBinaryToWriter
    );
  }
  f = message.getOrderbyList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      resources_common_database_database_pb.OrderBy.serializeBinaryToWriter
    );
  }
  f = message.getSearch();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getCategoryIdsList();
  if (f.length > 0) {
    writer.writePackedUint64(
      4,
      f
    );
  }
};


/**
 * optional resources.common.database.PaginationRequest pagination = 1;
 * @return {?proto.resources.common.database.PaginationRequest}
 */
proto.services.docstore.ListDocumentsRequest.prototype.getPagination = function() {
  return /** @type{?proto.resources.common.database.PaginationRequest} */ (
    jspb.Message.getWrapperField(this, resources_common_database_database_pb.PaginationRequest, 1));
};


/**
 * @param {?proto.resources.common.database.PaginationRequest|undefined} value
 * @return {!proto.services.docstore.ListDocumentsRequest} returns this
*/
proto.services.docstore.ListDocumentsRequest.prototype.setPagination = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.ListDocumentsRequest} returns this
 */
proto.services.docstore.ListDocumentsRequest.prototype.clearPagination = function() {
  return this.setPagination(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.ListDocumentsRequest.prototype.hasPagination = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * repeated resources.common.database.OrderBy orderBy = 2;
 * @return {!Array<!proto.resources.common.database.OrderBy>}
 */
proto.services.docstore.ListDocumentsRequest.prototype.getOrderbyList = function() {
  return /** @type{!Array<!proto.resources.common.database.OrderBy>} */ (
    jspb.Message.getRepeatedWrapperField(this, resources_common_database_database_pb.OrderBy, 2));
};


/**
 * @param {!Array<!proto.resources.common.database.OrderBy>} value
 * @return {!proto.services.docstore.ListDocumentsRequest} returns this
*/
proto.services.docstore.ListDocumentsRequest.prototype.setOrderbyList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.resources.common.database.OrderBy=} opt_value
 * @param {number=} opt_index
 * @return {!proto.resources.common.database.OrderBy}
 */
proto.services.docstore.ListDocumentsRequest.prototype.addOrderby = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.resources.common.database.OrderBy, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.services.docstore.ListDocumentsRequest} returns this
 */
proto.services.docstore.ListDocumentsRequest.prototype.clearOrderbyList = function() {
  return this.setOrderbyList([]);
};


/**
 * optional string search = 3;
 * @return {string}
 */
proto.services.docstore.ListDocumentsRequest.prototype.getSearch = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.services.docstore.ListDocumentsRequest} returns this
 */
proto.services.docstore.ListDocumentsRequest.prototype.setSearch = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * repeated uint64 category_ids = 4;
 * @return {!Array<number>}
 */
proto.services.docstore.ListDocumentsRequest.prototype.getCategoryIdsList = function() {
  return /** @type {!Array<number>} */ (jspb.Message.getRepeatedField(this, 4));
};


/**
 * @param {!Array<number>} value
 * @return {!proto.services.docstore.ListDocumentsRequest} returns this
 */
proto.services.docstore.ListDocumentsRequest.prototype.setCategoryIdsList = function(value) {
  return jspb.Message.setField(this, 4, value || []);
};


/**
 * @param {number} value
 * @param {number=} opt_index
 * @return {!proto.services.docstore.ListDocumentsRequest} returns this
 */
proto.services.docstore.ListDocumentsRequest.prototype.addCategoryIds = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 4, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.services.docstore.ListDocumentsRequest} returns this
 */
proto.services.docstore.ListDocumentsRequest.prototype.clearCategoryIdsList = function() {
  return this.setCategoryIdsList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.services.docstore.ListDocumentsResponse.repeatedFields_ = [2];



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
proto.services.docstore.ListDocumentsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.ListDocumentsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.ListDocumentsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListDocumentsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    pagination: (f = msg.getPagination()) && resources_common_database_database_pb.PaginationResponse.toObject(includeInstance, f),
    documentsList: jspb.Message.toObjectList(msg.getDocumentsList(),
    resources_documents_documents_pb.Document.toObject, includeInstance)
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
 * @return {!proto.services.docstore.ListDocumentsResponse}
 */
proto.services.docstore.ListDocumentsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.ListDocumentsResponse;
  return proto.services.docstore.ListDocumentsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.ListDocumentsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.ListDocumentsResponse}
 */
proto.services.docstore.ListDocumentsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_common_database_database_pb.PaginationResponse;
      reader.readMessage(value,resources_common_database_database_pb.PaginationResponse.deserializeBinaryFromReader);
      msg.setPagination(value);
      break;
    case 2:
      var value = new resources_documents_documents_pb.Document;
      reader.readMessage(value,resources_documents_documents_pb.Document.deserializeBinaryFromReader);
      msg.addDocuments(value);
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
proto.services.docstore.ListDocumentsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.ListDocumentsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.ListDocumentsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListDocumentsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPagination();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_common_database_database_pb.PaginationResponse.serializeBinaryToWriter
    );
  }
  f = message.getDocumentsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      resources_documents_documents_pb.Document.serializeBinaryToWriter
    );
  }
};


/**
 * optional resources.common.database.PaginationResponse pagination = 1;
 * @return {?proto.resources.common.database.PaginationResponse}
 */
proto.services.docstore.ListDocumentsResponse.prototype.getPagination = function() {
  return /** @type{?proto.resources.common.database.PaginationResponse} */ (
    jspb.Message.getWrapperField(this, resources_common_database_database_pb.PaginationResponse, 1));
};


/**
 * @param {?proto.resources.common.database.PaginationResponse|undefined} value
 * @return {!proto.services.docstore.ListDocumentsResponse} returns this
*/
proto.services.docstore.ListDocumentsResponse.prototype.setPagination = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.ListDocumentsResponse} returns this
 */
proto.services.docstore.ListDocumentsResponse.prototype.clearPagination = function() {
  return this.setPagination(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.ListDocumentsResponse.prototype.hasPagination = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * repeated resources.documents.Document documents = 2;
 * @return {!Array<!proto.resources.documents.Document>}
 */
proto.services.docstore.ListDocumentsResponse.prototype.getDocumentsList = function() {
  return /** @type{!Array<!proto.resources.documents.Document>} */ (
    jspb.Message.getRepeatedWrapperField(this, resources_documents_documents_pb.Document, 2));
};


/**
 * @param {!Array<!proto.resources.documents.Document>} value
 * @return {!proto.services.docstore.ListDocumentsResponse} returns this
*/
proto.services.docstore.ListDocumentsResponse.prototype.setDocumentsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.resources.documents.Document=} opt_value
 * @param {number=} opt_index
 * @return {!proto.resources.documents.Document}
 */
proto.services.docstore.ListDocumentsResponse.prototype.addDocuments = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.resources.documents.Document, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.services.docstore.ListDocumentsResponse} returns this
 */
proto.services.docstore.ListDocumentsResponse.prototype.clearDocumentsList = function() {
  return this.setDocumentsList([]);
};





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
proto.services.docstore.GetDocumentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.GetDocumentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.GetDocumentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    documentId: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.GetDocumentRequest}
 */
proto.services.docstore.GetDocumentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.GetDocumentRequest;
  return proto.services.docstore.GetDocumentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.GetDocumentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.GetDocumentRequest}
 */
proto.services.docstore.GetDocumentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setDocumentId(value);
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
proto.services.docstore.GetDocumentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.GetDocumentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.GetDocumentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDocumentId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 document_id = 1;
 * @return {number}
 */
proto.services.docstore.GetDocumentRequest.prototype.getDocumentId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.GetDocumentRequest} returns this
 */
proto.services.docstore.GetDocumentRequest.prototype.setDocumentId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.GetDocumentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.GetDocumentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.GetDocumentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    document: (f = msg.getDocument()) && resources_documents_documents_pb.Document.toObject(includeInstance, f),
    access: (f = msg.getAccess()) && resources_documents_documents_pb.DocumentAccess.toObject(includeInstance, f)
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
 * @return {!proto.services.docstore.GetDocumentResponse}
 */
proto.services.docstore.GetDocumentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.GetDocumentResponse;
  return proto.services.docstore.GetDocumentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.GetDocumentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.GetDocumentResponse}
 */
proto.services.docstore.GetDocumentResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_documents_pb.Document;
      reader.readMessage(value,resources_documents_documents_pb.Document.deserializeBinaryFromReader);
      msg.setDocument(value);
      break;
    case 2:
      var value = new resources_documents_documents_pb.DocumentAccess;
      reader.readMessage(value,resources_documents_documents_pb.DocumentAccess.deserializeBinaryFromReader);
      msg.setAccess(value);
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
proto.services.docstore.GetDocumentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.GetDocumentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.GetDocumentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDocument();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_documents_documents_pb.Document.serializeBinaryToWriter
    );
  }
  f = message.getAccess();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      resources_documents_documents_pb.DocumentAccess.serializeBinaryToWriter
    );
  }
};


/**
 * optional resources.documents.Document document = 1;
 * @return {?proto.resources.documents.Document}
 */
proto.services.docstore.GetDocumentResponse.prototype.getDocument = function() {
  return /** @type{?proto.resources.documents.Document} */ (
    jspb.Message.getWrapperField(this, resources_documents_documents_pb.Document, 1));
};


/**
 * @param {?proto.resources.documents.Document|undefined} value
 * @return {!proto.services.docstore.GetDocumentResponse} returns this
*/
proto.services.docstore.GetDocumentResponse.prototype.setDocument = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.GetDocumentResponse} returns this
 */
proto.services.docstore.GetDocumentResponse.prototype.clearDocument = function() {
  return this.setDocument(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.GetDocumentResponse.prototype.hasDocument = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional resources.documents.DocumentAccess access = 2;
 * @return {?proto.resources.documents.DocumentAccess}
 */
proto.services.docstore.GetDocumentResponse.prototype.getAccess = function() {
  return /** @type{?proto.resources.documents.DocumentAccess} */ (
    jspb.Message.getWrapperField(this, resources_documents_documents_pb.DocumentAccess, 2));
};


/**
 * @param {?proto.resources.documents.DocumentAccess|undefined} value
 * @return {!proto.services.docstore.GetDocumentResponse} returns this
*/
proto.services.docstore.GetDocumentResponse.prototype.setAccess = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.GetDocumentResponse} returns this
 */
proto.services.docstore.GetDocumentResponse.prototype.clearAccess = function() {
  return this.setAccess(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.GetDocumentResponse.prototype.hasAccess = function() {
  return jspb.Message.getField(this, 2) != null;
};





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
proto.services.docstore.GetDocumentReferencesRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.GetDocumentReferencesRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.GetDocumentReferencesRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentReferencesRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    documentId: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.GetDocumentReferencesRequest}
 */
proto.services.docstore.GetDocumentReferencesRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.GetDocumentReferencesRequest;
  return proto.services.docstore.GetDocumentReferencesRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.GetDocumentReferencesRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.GetDocumentReferencesRequest}
 */
proto.services.docstore.GetDocumentReferencesRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setDocumentId(value);
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
proto.services.docstore.GetDocumentReferencesRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.GetDocumentReferencesRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.GetDocumentReferencesRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentReferencesRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDocumentId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 document_id = 1;
 * @return {number}
 */
proto.services.docstore.GetDocumentReferencesRequest.prototype.getDocumentId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.GetDocumentReferencesRequest} returns this
 */
proto.services.docstore.GetDocumentReferencesRequest.prototype.setDocumentId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.services.docstore.GetDocumentReferencesResponse.repeatedFields_ = [1];



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
proto.services.docstore.GetDocumentReferencesResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.GetDocumentReferencesResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.GetDocumentReferencesResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentReferencesResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    referencesList: jspb.Message.toObjectList(msg.getReferencesList(),
    resources_documents_documents_pb.DocumentReference.toObject, includeInstance)
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
 * @return {!proto.services.docstore.GetDocumentReferencesResponse}
 */
proto.services.docstore.GetDocumentReferencesResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.GetDocumentReferencesResponse;
  return proto.services.docstore.GetDocumentReferencesResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.GetDocumentReferencesResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.GetDocumentReferencesResponse}
 */
proto.services.docstore.GetDocumentReferencesResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_documents_pb.DocumentReference;
      reader.readMessage(value,resources_documents_documents_pb.DocumentReference.deserializeBinaryFromReader);
      msg.addReferences(value);
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
proto.services.docstore.GetDocumentReferencesResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.GetDocumentReferencesResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.GetDocumentReferencesResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentReferencesResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getReferencesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      resources_documents_documents_pb.DocumentReference.serializeBinaryToWriter
    );
  }
};


/**
 * repeated resources.documents.DocumentReference references = 1;
 * @return {!Array<!proto.resources.documents.DocumentReference>}
 */
proto.services.docstore.GetDocumentReferencesResponse.prototype.getReferencesList = function() {
  return /** @type{!Array<!proto.resources.documents.DocumentReference>} */ (
    jspb.Message.getRepeatedWrapperField(this, resources_documents_documents_pb.DocumentReference, 1));
};


/**
 * @param {!Array<!proto.resources.documents.DocumentReference>} value
 * @return {!proto.services.docstore.GetDocumentReferencesResponse} returns this
*/
proto.services.docstore.GetDocumentReferencesResponse.prototype.setReferencesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.resources.documents.DocumentReference=} opt_value
 * @param {number=} opt_index
 * @return {!proto.resources.documents.DocumentReference}
 */
proto.services.docstore.GetDocumentReferencesResponse.prototype.addReferences = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.resources.documents.DocumentReference, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.services.docstore.GetDocumentReferencesResponse} returns this
 */
proto.services.docstore.GetDocumentReferencesResponse.prototype.clearReferencesList = function() {
  return this.setReferencesList([]);
};





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
proto.services.docstore.GetDocumentRelationsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.GetDocumentRelationsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.GetDocumentRelationsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentRelationsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    documentId: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.GetDocumentRelationsRequest}
 */
proto.services.docstore.GetDocumentRelationsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.GetDocumentRelationsRequest;
  return proto.services.docstore.GetDocumentRelationsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.GetDocumentRelationsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.GetDocumentRelationsRequest}
 */
proto.services.docstore.GetDocumentRelationsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setDocumentId(value);
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
proto.services.docstore.GetDocumentRelationsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.GetDocumentRelationsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.GetDocumentRelationsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentRelationsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDocumentId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 document_id = 1;
 * @return {number}
 */
proto.services.docstore.GetDocumentRelationsRequest.prototype.getDocumentId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.GetDocumentRelationsRequest} returns this
 */
proto.services.docstore.GetDocumentRelationsRequest.prototype.setDocumentId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.services.docstore.GetDocumentRelationsResponse.repeatedFields_ = [1];



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
proto.services.docstore.GetDocumentRelationsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.GetDocumentRelationsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.GetDocumentRelationsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentRelationsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    relationsList: jspb.Message.toObjectList(msg.getRelationsList(),
    resources_documents_documents_pb.DocumentRelation.toObject, includeInstance)
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
 * @return {!proto.services.docstore.GetDocumentRelationsResponse}
 */
proto.services.docstore.GetDocumentRelationsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.GetDocumentRelationsResponse;
  return proto.services.docstore.GetDocumentRelationsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.GetDocumentRelationsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.GetDocumentRelationsResponse}
 */
proto.services.docstore.GetDocumentRelationsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_documents_pb.DocumentRelation;
      reader.readMessage(value,resources_documents_documents_pb.DocumentRelation.deserializeBinaryFromReader);
      msg.addRelations(value);
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
proto.services.docstore.GetDocumentRelationsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.GetDocumentRelationsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.GetDocumentRelationsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentRelationsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRelationsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      resources_documents_documents_pb.DocumentRelation.serializeBinaryToWriter
    );
  }
};


/**
 * repeated resources.documents.DocumentRelation relations = 1;
 * @return {!Array<!proto.resources.documents.DocumentRelation>}
 */
proto.services.docstore.GetDocumentRelationsResponse.prototype.getRelationsList = function() {
  return /** @type{!Array<!proto.resources.documents.DocumentRelation>} */ (
    jspb.Message.getRepeatedWrapperField(this, resources_documents_documents_pb.DocumentRelation, 1));
};


/**
 * @param {!Array<!proto.resources.documents.DocumentRelation>} value
 * @return {!proto.services.docstore.GetDocumentRelationsResponse} returns this
*/
proto.services.docstore.GetDocumentRelationsResponse.prototype.setRelationsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.resources.documents.DocumentRelation=} opt_value
 * @param {number=} opt_index
 * @return {!proto.resources.documents.DocumentRelation}
 */
proto.services.docstore.GetDocumentRelationsResponse.prototype.addRelations = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.resources.documents.DocumentRelation, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.services.docstore.GetDocumentRelationsResponse} returns this
 */
proto.services.docstore.GetDocumentRelationsResponse.prototype.clearRelationsList = function() {
  return this.setRelationsList([]);
};





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
proto.services.docstore.AddDocumentReferenceRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.AddDocumentReferenceRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.AddDocumentReferenceRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.AddDocumentReferenceRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    reference: (f = msg.getReference()) && resources_documents_documents_pb.DocumentReference.toObject(includeInstance, f)
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
 * @return {!proto.services.docstore.AddDocumentReferenceRequest}
 */
proto.services.docstore.AddDocumentReferenceRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.AddDocumentReferenceRequest;
  return proto.services.docstore.AddDocumentReferenceRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.AddDocumentReferenceRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.AddDocumentReferenceRequest}
 */
proto.services.docstore.AddDocumentReferenceRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_documents_pb.DocumentReference;
      reader.readMessage(value,resources_documents_documents_pb.DocumentReference.deserializeBinaryFromReader);
      msg.setReference(value);
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
proto.services.docstore.AddDocumentReferenceRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.AddDocumentReferenceRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.AddDocumentReferenceRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.AddDocumentReferenceRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getReference();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_documents_documents_pb.DocumentReference.serializeBinaryToWriter
    );
  }
};


/**
 * optional resources.documents.DocumentReference reference = 1;
 * @return {?proto.resources.documents.DocumentReference}
 */
proto.services.docstore.AddDocumentReferenceRequest.prototype.getReference = function() {
  return /** @type{?proto.resources.documents.DocumentReference} */ (
    jspb.Message.getWrapperField(this, resources_documents_documents_pb.DocumentReference, 1));
};


/**
 * @param {?proto.resources.documents.DocumentReference|undefined} value
 * @return {!proto.services.docstore.AddDocumentReferenceRequest} returns this
*/
proto.services.docstore.AddDocumentReferenceRequest.prototype.setReference = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.AddDocumentReferenceRequest} returns this
 */
proto.services.docstore.AddDocumentReferenceRequest.prototype.clearReference = function() {
  return this.setReference(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.AddDocumentReferenceRequest.prototype.hasReference = function() {
  return jspb.Message.getField(this, 1) != null;
};





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
proto.services.docstore.AddDocumentReferenceResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.AddDocumentReferenceResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.AddDocumentReferenceResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.AddDocumentReferenceResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.AddDocumentReferenceResponse}
 */
proto.services.docstore.AddDocumentReferenceResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.AddDocumentReferenceResponse;
  return proto.services.docstore.AddDocumentReferenceResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.AddDocumentReferenceResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.AddDocumentReferenceResponse}
 */
proto.services.docstore.AddDocumentReferenceResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.AddDocumentReferenceResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.AddDocumentReferenceResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.AddDocumentReferenceResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.AddDocumentReferenceResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 id = 1;
 * @return {number}
 */
proto.services.docstore.AddDocumentReferenceResponse.prototype.getId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.AddDocumentReferenceResponse} returns this
 */
proto.services.docstore.AddDocumentReferenceResponse.prototype.setId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.RemoveDocumentReferenceRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.RemoveDocumentReferenceRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.RemoveDocumentReferenceRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.RemoveDocumentReferenceRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.RemoveDocumentReferenceRequest}
 */
proto.services.docstore.RemoveDocumentReferenceRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.RemoveDocumentReferenceRequest;
  return proto.services.docstore.RemoveDocumentReferenceRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.RemoveDocumentReferenceRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.RemoveDocumentReferenceRequest}
 */
proto.services.docstore.RemoveDocumentReferenceRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.RemoveDocumentReferenceRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.RemoveDocumentReferenceRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.RemoveDocumentReferenceRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.RemoveDocumentReferenceRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 id = 1;
 * @return {number}
 */
proto.services.docstore.RemoveDocumentReferenceRequest.prototype.getId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.RemoveDocumentReferenceRequest} returns this
 */
proto.services.docstore.RemoveDocumentReferenceRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.RemoveDocumentReferenceResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.RemoveDocumentReferenceResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.RemoveDocumentReferenceResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.RemoveDocumentReferenceResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.services.docstore.RemoveDocumentReferenceResponse}
 */
proto.services.docstore.RemoveDocumentReferenceResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.RemoveDocumentReferenceResponse;
  return proto.services.docstore.RemoveDocumentReferenceResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.RemoveDocumentReferenceResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.RemoveDocumentReferenceResponse}
 */
proto.services.docstore.RemoveDocumentReferenceResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.RemoveDocumentReferenceResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.RemoveDocumentReferenceResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.RemoveDocumentReferenceResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.RemoveDocumentReferenceResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





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
proto.services.docstore.AddDocumentRelationRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.AddDocumentRelationRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.AddDocumentRelationRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.AddDocumentRelationRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    relation: (f = msg.getRelation()) && resources_documents_documents_pb.DocumentRelation.toObject(includeInstance, f)
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
 * @return {!proto.services.docstore.AddDocumentRelationRequest}
 */
proto.services.docstore.AddDocumentRelationRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.AddDocumentRelationRequest;
  return proto.services.docstore.AddDocumentRelationRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.AddDocumentRelationRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.AddDocumentRelationRequest}
 */
proto.services.docstore.AddDocumentRelationRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_documents_pb.DocumentRelation;
      reader.readMessage(value,resources_documents_documents_pb.DocumentRelation.deserializeBinaryFromReader);
      msg.setRelation(value);
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
proto.services.docstore.AddDocumentRelationRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.AddDocumentRelationRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.AddDocumentRelationRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.AddDocumentRelationRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRelation();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_documents_documents_pb.DocumentRelation.serializeBinaryToWriter
    );
  }
};


/**
 * optional resources.documents.DocumentRelation relation = 1;
 * @return {?proto.resources.documents.DocumentRelation}
 */
proto.services.docstore.AddDocumentRelationRequest.prototype.getRelation = function() {
  return /** @type{?proto.resources.documents.DocumentRelation} */ (
    jspb.Message.getWrapperField(this, resources_documents_documents_pb.DocumentRelation, 1));
};


/**
 * @param {?proto.resources.documents.DocumentRelation|undefined} value
 * @return {!proto.services.docstore.AddDocumentRelationRequest} returns this
*/
proto.services.docstore.AddDocumentRelationRequest.prototype.setRelation = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.AddDocumentRelationRequest} returns this
 */
proto.services.docstore.AddDocumentRelationRequest.prototype.clearRelation = function() {
  return this.setRelation(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.AddDocumentRelationRequest.prototype.hasRelation = function() {
  return jspb.Message.getField(this, 1) != null;
};





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
proto.services.docstore.AddDocumentRelationResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.AddDocumentRelationResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.AddDocumentRelationResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.AddDocumentRelationResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.AddDocumentRelationResponse}
 */
proto.services.docstore.AddDocumentRelationResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.AddDocumentRelationResponse;
  return proto.services.docstore.AddDocumentRelationResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.AddDocumentRelationResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.AddDocumentRelationResponse}
 */
proto.services.docstore.AddDocumentRelationResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.AddDocumentRelationResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.AddDocumentRelationResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.AddDocumentRelationResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.AddDocumentRelationResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 id = 1;
 * @return {number}
 */
proto.services.docstore.AddDocumentRelationResponse.prototype.getId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.AddDocumentRelationResponse} returns this
 */
proto.services.docstore.AddDocumentRelationResponse.prototype.setId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.RemoveDocumentRelationRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.RemoveDocumentRelationRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.RemoveDocumentRelationRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.RemoveDocumentRelationRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.RemoveDocumentRelationRequest}
 */
proto.services.docstore.RemoveDocumentRelationRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.RemoveDocumentRelationRequest;
  return proto.services.docstore.RemoveDocumentRelationRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.RemoveDocumentRelationRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.RemoveDocumentRelationRequest}
 */
proto.services.docstore.RemoveDocumentRelationRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.RemoveDocumentRelationRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.RemoveDocumentRelationRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.RemoveDocumentRelationRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.RemoveDocumentRelationRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 id = 1;
 * @return {number}
 */
proto.services.docstore.RemoveDocumentRelationRequest.prototype.getId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.RemoveDocumentRelationRequest} returns this
 */
proto.services.docstore.RemoveDocumentRelationRequest.prototype.setId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.RemoveDocumentRelationResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.RemoveDocumentRelationResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.RemoveDocumentRelationResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.RemoveDocumentRelationResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.services.docstore.RemoveDocumentRelationResponse}
 */
proto.services.docstore.RemoveDocumentRelationResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.RemoveDocumentRelationResponse;
  return proto.services.docstore.RemoveDocumentRelationResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.RemoveDocumentRelationResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.RemoveDocumentRelationResponse}
 */
proto.services.docstore.RemoveDocumentRelationResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.RemoveDocumentRelationResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.RemoveDocumentRelationResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.RemoveDocumentRelationResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.RemoveDocumentRelationResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





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
proto.services.docstore.GetDocumentCommentsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.GetDocumentCommentsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.GetDocumentCommentsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentCommentsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pagination: (f = msg.getPagination()) && resources_common_database_database_pb.PaginationRequest.toObject(includeInstance, f),
    documentId: jspb.Message.getFieldWithDefault(msg, 2, 0)
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
 * @return {!proto.services.docstore.GetDocumentCommentsRequest}
 */
proto.services.docstore.GetDocumentCommentsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.GetDocumentCommentsRequest;
  return proto.services.docstore.GetDocumentCommentsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.GetDocumentCommentsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.GetDocumentCommentsRequest}
 */
proto.services.docstore.GetDocumentCommentsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_common_database_database_pb.PaginationRequest;
      reader.readMessage(value,resources_common_database_database_pb.PaginationRequest.deserializeBinaryFromReader);
      msg.setPagination(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setDocumentId(value);
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
proto.services.docstore.GetDocumentCommentsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.GetDocumentCommentsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.GetDocumentCommentsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentCommentsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPagination();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_common_database_database_pb.PaginationRequest.serializeBinaryToWriter
    );
  }
  f = message.getDocumentId();
  if (f !== 0) {
    writer.writeUint64(
      2,
      f
    );
  }
};


/**
 * optional resources.common.database.PaginationRequest pagination = 1;
 * @return {?proto.resources.common.database.PaginationRequest}
 */
proto.services.docstore.GetDocumentCommentsRequest.prototype.getPagination = function() {
  return /** @type{?proto.resources.common.database.PaginationRequest} */ (
    jspb.Message.getWrapperField(this, resources_common_database_database_pb.PaginationRequest, 1));
};


/**
 * @param {?proto.resources.common.database.PaginationRequest|undefined} value
 * @return {!proto.services.docstore.GetDocumentCommentsRequest} returns this
*/
proto.services.docstore.GetDocumentCommentsRequest.prototype.setPagination = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.GetDocumentCommentsRequest} returns this
 */
proto.services.docstore.GetDocumentCommentsRequest.prototype.clearPagination = function() {
  return this.setPagination(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.GetDocumentCommentsRequest.prototype.hasPagination = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional uint64 document_id = 2;
 * @return {number}
 */
proto.services.docstore.GetDocumentCommentsRequest.prototype.getDocumentId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.GetDocumentCommentsRequest} returns this
 */
proto.services.docstore.GetDocumentCommentsRequest.prototype.setDocumentId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.services.docstore.GetDocumentCommentsResponse.repeatedFields_ = [2];



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
proto.services.docstore.GetDocumentCommentsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.GetDocumentCommentsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.GetDocumentCommentsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentCommentsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    pagination: (f = msg.getPagination()) && resources_common_database_database_pb.PaginationResponse.toObject(includeInstance, f),
    commentsList: jspb.Message.toObjectList(msg.getCommentsList(),
    resources_documents_documents_pb.DocumentComment.toObject, includeInstance)
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
 * @return {!proto.services.docstore.GetDocumentCommentsResponse}
 */
proto.services.docstore.GetDocumentCommentsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.GetDocumentCommentsResponse;
  return proto.services.docstore.GetDocumentCommentsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.GetDocumentCommentsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.GetDocumentCommentsResponse}
 */
proto.services.docstore.GetDocumentCommentsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_common_database_database_pb.PaginationResponse;
      reader.readMessage(value,resources_common_database_database_pb.PaginationResponse.deserializeBinaryFromReader);
      msg.setPagination(value);
      break;
    case 2:
      var value = new resources_documents_documents_pb.DocumentComment;
      reader.readMessage(value,resources_documents_documents_pb.DocumentComment.deserializeBinaryFromReader);
      msg.addComments(value);
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
proto.services.docstore.GetDocumentCommentsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.GetDocumentCommentsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.GetDocumentCommentsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentCommentsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPagination();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_common_database_database_pb.PaginationResponse.serializeBinaryToWriter
    );
  }
  f = message.getCommentsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      resources_documents_documents_pb.DocumentComment.serializeBinaryToWriter
    );
  }
};


/**
 * optional resources.common.database.PaginationResponse pagination = 1;
 * @return {?proto.resources.common.database.PaginationResponse}
 */
proto.services.docstore.GetDocumentCommentsResponse.prototype.getPagination = function() {
  return /** @type{?proto.resources.common.database.PaginationResponse} */ (
    jspb.Message.getWrapperField(this, resources_common_database_database_pb.PaginationResponse, 1));
};


/**
 * @param {?proto.resources.common.database.PaginationResponse|undefined} value
 * @return {!proto.services.docstore.GetDocumentCommentsResponse} returns this
*/
proto.services.docstore.GetDocumentCommentsResponse.prototype.setPagination = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.GetDocumentCommentsResponse} returns this
 */
proto.services.docstore.GetDocumentCommentsResponse.prototype.clearPagination = function() {
  return this.setPagination(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.GetDocumentCommentsResponse.prototype.hasPagination = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * repeated resources.documents.DocumentComment comments = 2;
 * @return {!Array<!proto.resources.documents.DocumentComment>}
 */
proto.services.docstore.GetDocumentCommentsResponse.prototype.getCommentsList = function() {
  return /** @type{!Array<!proto.resources.documents.DocumentComment>} */ (
    jspb.Message.getRepeatedWrapperField(this, resources_documents_documents_pb.DocumentComment, 2));
};


/**
 * @param {!Array<!proto.resources.documents.DocumentComment>} value
 * @return {!proto.services.docstore.GetDocumentCommentsResponse} returns this
*/
proto.services.docstore.GetDocumentCommentsResponse.prototype.setCommentsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.resources.documents.DocumentComment=} opt_value
 * @param {number=} opt_index
 * @return {!proto.resources.documents.DocumentComment}
 */
proto.services.docstore.GetDocumentCommentsResponse.prototype.addComments = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.resources.documents.DocumentComment, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.services.docstore.GetDocumentCommentsResponse} returns this
 */
proto.services.docstore.GetDocumentCommentsResponse.prototype.clearCommentsList = function() {
  return this.setCommentsList([]);
};





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
proto.services.docstore.PostDocumentCommentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.PostDocumentCommentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.PostDocumentCommentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.PostDocumentCommentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    comment: (f = msg.getComment()) && resources_documents_documents_pb.DocumentComment.toObject(includeInstance, f)
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
 * @return {!proto.services.docstore.PostDocumentCommentRequest}
 */
proto.services.docstore.PostDocumentCommentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.PostDocumentCommentRequest;
  return proto.services.docstore.PostDocumentCommentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.PostDocumentCommentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.PostDocumentCommentRequest}
 */
proto.services.docstore.PostDocumentCommentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_documents_pb.DocumentComment;
      reader.readMessage(value,resources_documents_documents_pb.DocumentComment.deserializeBinaryFromReader);
      msg.setComment(value);
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
proto.services.docstore.PostDocumentCommentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.PostDocumentCommentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.PostDocumentCommentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.PostDocumentCommentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getComment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_documents_documents_pb.DocumentComment.serializeBinaryToWriter
    );
  }
};


/**
 * optional resources.documents.DocumentComment comment = 1;
 * @return {?proto.resources.documents.DocumentComment}
 */
proto.services.docstore.PostDocumentCommentRequest.prototype.getComment = function() {
  return /** @type{?proto.resources.documents.DocumentComment} */ (
    jspb.Message.getWrapperField(this, resources_documents_documents_pb.DocumentComment, 1));
};


/**
 * @param {?proto.resources.documents.DocumentComment|undefined} value
 * @return {!proto.services.docstore.PostDocumentCommentRequest} returns this
*/
proto.services.docstore.PostDocumentCommentRequest.prototype.setComment = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.PostDocumentCommentRequest} returns this
 */
proto.services.docstore.PostDocumentCommentRequest.prototype.clearComment = function() {
  return this.setComment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.PostDocumentCommentRequest.prototype.hasComment = function() {
  return jspb.Message.getField(this, 1) != null;
};





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
proto.services.docstore.PostDocumentCommentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.PostDocumentCommentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.PostDocumentCommentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.PostDocumentCommentResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.PostDocumentCommentResponse}
 */
proto.services.docstore.PostDocumentCommentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.PostDocumentCommentResponse;
  return proto.services.docstore.PostDocumentCommentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.PostDocumentCommentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.PostDocumentCommentResponse}
 */
proto.services.docstore.PostDocumentCommentResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.PostDocumentCommentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.PostDocumentCommentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.PostDocumentCommentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.PostDocumentCommentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 id = 1;
 * @return {number}
 */
proto.services.docstore.PostDocumentCommentResponse.prototype.getId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.PostDocumentCommentResponse} returns this
 */
proto.services.docstore.PostDocumentCommentResponse.prototype.setId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.EditDocumentCommentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.EditDocumentCommentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.EditDocumentCommentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.EditDocumentCommentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    comment: (f = msg.getComment()) && resources_documents_documents_pb.DocumentComment.toObject(includeInstance, f)
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
 * @return {!proto.services.docstore.EditDocumentCommentRequest}
 */
proto.services.docstore.EditDocumentCommentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.EditDocumentCommentRequest;
  return proto.services.docstore.EditDocumentCommentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.EditDocumentCommentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.EditDocumentCommentRequest}
 */
proto.services.docstore.EditDocumentCommentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_documents_pb.DocumentComment;
      reader.readMessage(value,resources_documents_documents_pb.DocumentComment.deserializeBinaryFromReader);
      msg.setComment(value);
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
proto.services.docstore.EditDocumentCommentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.EditDocumentCommentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.EditDocumentCommentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.EditDocumentCommentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getComment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_documents_documents_pb.DocumentComment.serializeBinaryToWriter
    );
  }
};


/**
 * optional resources.documents.DocumentComment comment = 1;
 * @return {?proto.resources.documents.DocumentComment}
 */
proto.services.docstore.EditDocumentCommentRequest.prototype.getComment = function() {
  return /** @type{?proto.resources.documents.DocumentComment} */ (
    jspb.Message.getWrapperField(this, resources_documents_documents_pb.DocumentComment, 1));
};


/**
 * @param {?proto.resources.documents.DocumentComment|undefined} value
 * @return {!proto.services.docstore.EditDocumentCommentRequest} returns this
*/
proto.services.docstore.EditDocumentCommentRequest.prototype.setComment = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.EditDocumentCommentRequest} returns this
 */
proto.services.docstore.EditDocumentCommentRequest.prototype.clearComment = function() {
  return this.setComment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.EditDocumentCommentRequest.prototype.hasComment = function() {
  return jspb.Message.getField(this, 1) != null;
};





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
proto.services.docstore.EditDocumentCommentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.EditDocumentCommentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.EditDocumentCommentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.EditDocumentCommentResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.services.docstore.EditDocumentCommentResponse}
 */
proto.services.docstore.EditDocumentCommentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.EditDocumentCommentResponse;
  return proto.services.docstore.EditDocumentCommentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.EditDocumentCommentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.EditDocumentCommentResponse}
 */
proto.services.docstore.EditDocumentCommentResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.EditDocumentCommentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.EditDocumentCommentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.EditDocumentCommentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.EditDocumentCommentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





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
proto.services.docstore.DeleteDocumentCommentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.DeleteDocumentCommentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.DeleteDocumentCommentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteDocumentCommentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    commentId: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.DeleteDocumentCommentRequest}
 */
proto.services.docstore.DeleteDocumentCommentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.DeleteDocumentCommentRequest;
  return proto.services.docstore.DeleteDocumentCommentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.DeleteDocumentCommentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.DeleteDocumentCommentRequest}
 */
proto.services.docstore.DeleteDocumentCommentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setCommentId(value);
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
proto.services.docstore.DeleteDocumentCommentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.DeleteDocumentCommentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.DeleteDocumentCommentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteDocumentCommentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCommentId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 comment_id = 1;
 * @return {number}
 */
proto.services.docstore.DeleteDocumentCommentRequest.prototype.getCommentId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.DeleteDocumentCommentRequest} returns this
 */
proto.services.docstore.DeleteDocumentCommentRequest.prototype.setCommentId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.DeleteDocumentCommentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.DeleteDocumentCommentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.DeleteDocumentCommentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteDocumentCommentResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.services.docstore.DeleteDocumentCommentResponse}
 */
proto.services.docstore.DeleteDocumentCommentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.DeleteDocumentCommentResponse;
  return proto.services.docstore.DeleteDocumentCommentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.DeleteDocumentCommentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.DeleteDocumentCommentResponse}
 */
proto.services.docstore.DeleteDocumentCommentResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.DeleteDocumentCommentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.DeleteDocumentCommentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.DeleteDocumentCommentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteDocumentCommentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





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
proto.services.docstore.CreateDocumentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.CreateDocumentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.CreateDocumentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.CreateDocumentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    categoryId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    title: jspb.Message.getFieldWithDefault(msg, 2, ""),
    content: jspb.Message.getFieldWithDefault(msg, 3, ""),
    contentType: jspb.Message.getFieldWithDefault(msg, 4, 0),
    data: jspb.Message.getFieldWithDefault(msg, 5, ""),
    state: jspb.Message.getFieldWithDefault(msg, 6, ""),
    closed: jspb.Message.getBooleanFieldWithDefault(msg, 7, false),
    pb_public: jspb.Message.getBooleanFieldWithDefault(msg, 8, false),
    access: (f = msg.getAccess()) && resources_documents_documents_pb.DocumentAccess.toObject(includeInstance, f)
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
 * @return {!proto.services.docstore.CreateDocumentRequest}
 */
proto.services.docstore.CreateDocumentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.CreateDocumentRequest;
  return proto.services.docstore.CreateDocumentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.CreateDocumentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.CreateDocumentRequest}
 */
proto.services.docstore.CreateDocumentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setCategoryId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setTitle(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setContent(value);
      break;
    case 4:
      var value = /** @type {!proto.resources.documents.DOC_CONTENT_TYPE} */ (reader.readEnum());
      msg.setContentType(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setData(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setState(value);
      break;
    case 7:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setClosed(value);
      break;
    case 8:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setPublic(value);
      break;
    case 9:
      var value = new resources_documents_documents_pb.DocumentAccess;
      reader.readMessage(value,resources_documents_documents_pb.DocumentAccess.deserializeBinaryFromReader);
      msg.setAccess(value);
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
proto.services.docstore.CreateDocumentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.CreateDocumentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.CreateDocumentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.CreateDocumentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = /** @type {number} */ (jspb.Message.getField(message, 1));
  if (f != null) {
    writer.writeUint64(
      1,
      f
    );
  }
  f = message.getTitle();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getContent();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getContentType();
  if (f !== 0.0) {
    writer.writeEnum(
      4,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 5));
  if (f != null) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getState();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getClosed();
  if (f) {
    writer.writeBool(
      7,
      f
    );
  }
  f = message.getPublic();
  if (f) {
    writer.writeBool(
      8,
      f
    );
  }
  f = message.getAccess();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      resources_documents_documents_pb.DocumentAccess.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint64 category_id = 1;
 * @return {number}
 */
proto.services.docstore.CreateDocumentRequest.prototype.getCategoryId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.CreateDocumentRequest} returns this
 */
proto.services.docstore.CreateDocumentRequest.prototype.setCategoryId = function(value) {
  return jspb.Message.setField(this, 1, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.services.docstore.CreateDocumentRequest} returns this
 */
proto.services.docstore.CreateDocumentRequest.prototype.clearCategoryId = function() {
  return jspb.Message.setField(this, 1, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.CreateDocumentRequest.prototype.hasCategoryId = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string title = 2;
 * @return {string}
 */
proto.services.docstore.CreateDocumentRequest.prototype.getTitle = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.services.docstore.CreateDocumentRequest} returns this
 */
proto.services.docstore.CreateDocumentRequest.prototype.setTitle = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string content = 3;
 * @return {string}
 */
proto.services.docstore.CreateDocumentRequest.prototype.getContent = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.services.docstore.CreateDocumentRequest} returns this
 */
proto.services.docstore.CreateDocumentRequest.prototype.setContent = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional resources.documents.DOC_CONTENT_TYPE content_type = 4;
 * @return {!proto.resources.documents.DOC_CONTENT_TYPE}
 */
proto.services.docstore.CreateDocumentRequest.prototype.getContentType = function() {
  return /** @type {!proto.resources.documents.DOC_CONTENT_TYPE} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {!proto.resources.documents.DOC_CONTENT_TYPE} value
 * @return {!proto.services.docstore.CreateDocumentRequest} returns this
 */
proto.services.docstore.CreateDocumentRequest.prototype.setContentType = function(value) {
  return jspb.Message.setProto3EnumField(this, 4, value);
};


/**
 * optional string data = 5;
 * @return {string}
 */
proto.services.docstore.CreateDocumentRequest.prototype.getData = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.services.docstore.CreateDocumentRequest} returns this
 */
proto.services.docstore.CreateDocumentRequest.prototype.setData = function(value) {
  return jspb.Message.setField(this, 5, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.services.docstore.CreateDocumentRequest} returns this
 */
proto.services.docstore.CreateDocumentRequest.prototype.clearData = function() {
  return jspb.Message.setField(this, 5, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.CreateDocumentRequest.prototype.hasData = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional string state = 6;
 * @return {string}
 */
proto.services.docstore.CreateDocumentRequest.prototype.getState = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/**
 * @param {string} value
 * @return {!proto.services.docstore.CreateDocumentRequest} returns this
 */
proto.services.docstore.CreateDocumentRequest.prototype.setState = function(value) {
  return jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional bool closed = 7;
 * @return {boolean}
 */
proto.services.docstore.CreateDocumentRequest.prototype.getClosed = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 7, false));
};


/**
 * @param {boolean} value
 * @return {!proto.services.docstore.CreateDocumentRequest} returns this
 */
proto.services.docstore.CreateDocumentRequest.prototype.setClosed = function(value) {
  return jspb.Message.setProto3BooleanField(this, 7, value);
};


/**
 * optional bool public = 8;
 * @return {boolean}
 */
proto.services.docstore.CreateDocumentRequest.prototype.getPublic = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 8, false));
};


/**
 * @param {boolean} value
 * @return {!proto.services.docstore.CreateDocumentRequest} returns this
 */
proto.services.docstore.CreateDocumentRequest.prototype.setPublic = function(value) {
  return jspb.Message.setProto3BooleanField(this, 8, value);
};


/**
 * optional resources.documents.DocumentAccess access = 9;
 * @return {?proto.resources.documents.DocumentAccess}
 */
proto.services.docstore.CreateDocumentRequest.prototype.getAccess = function() {
  return /** @type{?proto.resources.documents.DocumentAccess} */ (
    jspb.Message.getWrapperField(this, resources_documents_documents_pb.DocumentAccess, 9));
};


/**
 * @param {?proto.resources.documents.DocumentAccess|undefined} value
 * @return {!proto.services.docstore.CreateDocumentRequest} returns this
*/
proto.services.docstore.CreateDocumentRequest.prototype.setAccess = function(value) {
  return jspb.Message.setWrapperField(this, 9, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.CreateDocumentRequest} returns this
 */
proto.services.docstore.CreateDocumentRequest.prototype.clearAccess = function() {
  return this.setAccess(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.CreateDocumentRequest.prototype.hasAccess = function() {
  return jspb.Message.getField(this, 9) != null;
};





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
proto.services.docstore.CreateDocumentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.CreateDocumentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.CreateDocumentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.CreateDocumentResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    documentId: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.CreateDocumentResponse}
 */
proto.services.docstore.CreateDocumentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.CreateDocumentResponse;
  return proto.services.docstore.CreateDocumentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.CreateDocumentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.CreateDocumentResponse}
 */
proto.services.docstore.CreateDocumentResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setDocumentId(value);
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
proto.services.docstore.CreateDocumentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.CreateDocumentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.CreateDocumentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.CreateDocumentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDocumentId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 document_id = 1;
 * @return {number}
 */
proto.services.docstore.CreateDocumentResponse.prototype.getDocumentId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.CreateDocumentResponse} returns this
 */
proto.services.docstore.CreateDocumentResponse.prototype.setDocumentId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.UpdateDocumentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.UpdateDocumentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.UpdateDocumentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.UpdateDocumentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    documentId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    categoryId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    title: jspb.Message.getFieldWithDefault(msg, 3, ""),
    content: jspb.Message.getFieldWithDefault(msg, 4, ""),
    contentType: jspb.Message.getFieldWithDefault(msg, 5, 0),
    data: jspb.Message.getFieldWithDefault(msg, 6, ""),
    state: jspb.Message.getFieldWithDefault(msg, 7, ""),
    closed: jspb.Message.getBooleanFieldWithDefault(msg, 8, false),
    pb_public: jspb.Message.getBooleanFieldWithDefault(msg, 9, false),
    access: (f = msg.getAccess()) && resources_documents_documents_pb.DocumentAccess.toObject(includeInstance, f)
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
 * @return {!proto.services.docstore.UpdateDocumentRequest}
 */
proto.services.docstore.UpdateDocumentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.UpdateDocumentRequest;
  return proto.services.docstore.UpdateDocumentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.UpdateDocumentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.UpdateDocumentRequest}
 */
proto.services.docstore.UpdateDocumentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setDocumentId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setCategoryId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setTitle(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setContent(value);
      break;
    case 5:
      var value = /** @type {!proto.resources.documents.DOC_CONTENT_TYPE} */ (reader.readEnum());
      msg.setContentType(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setData(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setState(value);
      break;
    case 8:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setClosed(value);
      break;
    case 9:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setPublic(value);
      break;
    case 10:
      var value = new resources_documents_documents_pb.DocumentAccess;
      reader.readMessage(value,resources_documents_documents_pb.DocumentAccess.deserializeBinaryFromReader);
      msg.setAccess(value);
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
proto.services.docstore.UpdateDocumentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.UpdateDocumentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.UpdateDocumentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.UpdateDocumentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDocumentId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 2));
  if (f != null) {
    writer.writeUint64(
      2,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 3));
  if (f != null) {
    writer.writeString(
      3,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 4));
  if (f != null) {
    writer.writeString(
      4,
      f
    );
  }
  f = /** @type {!proto.resources.documents.DOC_CONTENT_TYPE} */ (jspb.Message.getField(message, 5));
  if (f != null) {
    writer.writeEnum(
      5,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 6));
  if (f != null) {
    writer.writeString(
      6,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 7));
  if (f != null) {
    writer.writeString(
      7,
      f
    );
  }
  f = /** @type {boolean} */ (jspb.Message.getField(message, 8));
  if (f != null) {
    writer.writeBool(
      8,
      f
    );
  }
  f = /** @type {boolean} */ (jspb.Message.getField(message, 9));
  if (f != null) {
    writer.writeBool(
      9,
      f
    );
  }
  f = message.getAccess();
  if (f != null) {
    writer.writeMessage(
      10,
      f,
      resources_documents_documents_pb.DocumentAccess.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint64 document_id = 1;
 * @return {number}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.getDocumentId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.setDocumentId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint64 category_id = 2;
 * @return {number}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.getCategoryId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.setCategoryId = function(value) {
  return jspb.Message.setField(this, 2, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.clearCategoryId = function() {
  return jspb.Message.setField(this, 2, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.hasCategoryId = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string title = 3;
 * @return {string}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.getTitle = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.setTitle = function(value) {
  return jspb.Message.setField(this, 3, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.clearTitle = function() {
  return jspb.Message.setField(this, 3, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.hasTitle = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional string content = 4;
 * @return {string}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.getContent = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.setContent = function(value) {
  return jspb.Message.setField(this, 4, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.clearContent = function() {
  return jspb.Message.setField(this, 4, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.hasContent = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional resources.documents.DOC_CONTENT_TYPE content_type = 5;
 * @return {!proto.resources.documents.DOC_CONTENT_TYPE}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.getContentType = function() {
  return /** @type {!proto.resources.documents.DOC_CONTENT_TYPE} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {!proto.resources.documents.DOC_CONTENT_TYPE} value
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.setContentType = function(value) {
  return jspb.Message.setField(this, 5, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.clearContentType = function() {
  return jspb.Message.setField(this, 5, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.hasContentType = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional string data = 6;
 * @return {string}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.getData = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/**
 * @param {string} value
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.setData = function(value) {
  return jspb.Message.setField(this, 6, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.clearData = function() {
  return jspb.Message.setField(this, 6, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.hasData = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional string state = 7;
 * @return {string}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.getState = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/**
 * @param {string} value
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.setState = function(value) {
  return jspb.Message.setField(this, 7, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.clearState = function() {
  return jspb.Message.setField(this, 7, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.hasState = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional bool closed = 8;
 * @return {boolean}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.getClosed = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 8, false));
};


/**
 * @param {boolean} value
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.setClosed = function(value) {
  return jspb.Message.setField(this, 8, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.clearClosed = function() {
  return jspb.Message.setField(this, 8, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.hasClosed = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional bool public = 9;
 * @return {boolean}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.getPublic = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 9, false));
};


/**
 * @param {boolean} value
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.setPublic = function(value) {
  return jspb.Message.setField(this, 9, value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.clearPublic = function() {
  return jspb.Message.setField(this, 9, undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.hasPublic = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * optional resources.documents.DocumentAccess access = 10;
 * @return {?proto.resources.documents.DocumentAccess}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.getAccess = function() {
  return /** @type{?proto.resources.documents.DocumentAccess} */ (
    jspb.Message.getWrapperField(this, resources_documents_documents_pb.DocumentAccess, 10));
};


/**
 * @param {?proto.resources.documents.DocumentAccess|undefined} value
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
*/
proto.services.docstore.UpdateDocumentRequest.prototype.setAccess = function(value) {
  return jspb.Message.setWrapperField(this, 10, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.UpdateDocumentRequest} returns this
 */
proto.services.docstore.UpdateDocumentRequest.prototype.clearAccess = function() {
  return this.setAccess(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.UpdateDocumentRequest.prototype.hasAccess = function() {
  return jspb.Message.getField(this, 10) != null;
};





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
proto.services.docstore.UpdateDocumentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.UpdateDocumentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.UpdateDocumentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.UpdateDocumentResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    documentId: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.UpdateDocumentResponse}
 */
proto.services.docstore.UpdateDocumentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.UpdateDocumentResponse;
  return proto.services.docstore.UpdateDocumentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.UpdateDocumentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.UpdateDocumentResponse}
 */
proto.services.docstore.UpdateDocumentResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setDocumentId(value);
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
proto.services.docstore.UpdateDocumentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.UpdateDocumentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.UpdateDocumentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.UpdateDocumentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDocumentId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 document_id = 1;
 * @return {number}
 */
proto.services.docstore.UpdateDocumentResponse.prototype.getDocumentId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.UpdateDocumentResponse} returns this
 */
proto.services.docstore.UpdateDocumentResponse.prototype.setDocumentId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.DeleteDocumentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.DeleteDocumentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.DeleteDocumentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteDocumentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    documentId: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.DeleteDocumentRequest}
 */
proto.services.docstore.DeleteDocumentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.DeleteDocumentRequest;
  return proto.services.docstore.DeleteDocumentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.DeleteDocumentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.DeleteDocumentRequest}
 */
proto.services.docstore.DeleteDocumentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setDocumentId(value);
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
proto.services.docstore.DeleteDocumentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.DeleteDocumentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.DeleteDocumentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteDocumentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDocumentId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 document_id = 1;
 * @return {number}
 */
proto.services.docstore.DeleteDocumentRequest.prototype.getDocumentId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.DeleteDocumentRequest} returns this
 */
proto.services.docstore.DeleteDocumentRequest.prototype.setDocumentId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.DeleteDocumentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.DeleteDocumentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.DeleteDocumentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteDocumentResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.services.docstore.DeleteDocumentResponse}
 */
proto.services.docstore.DeleteDocumentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.DeleteDocumentResponse;
  return proto.services.docstore.DeleteDocumentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.DeleteDocumentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.DeleteDocumentResponse}
 */
proto.services.docstore.DeleteDocumentResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.DeleteDocumentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.DeleteDocumentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.DeleteDocumentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteDocumentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





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
proto.services.docstore.GetDocumentAccessRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.GetDocumentAccessRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.GetDocumentAccessRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentAccessRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    documentId: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.GetDocumentAccessRequest}
 */
proto.services.docstore.GetDocumentAccessRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.GetDocumentAccessRequest;
  return proto.services.docstore.GetDocumentAccessRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.GetDocumentAccessRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.GetDocumentAccessRequest}
 */
proto.services.docstore.GetDocumentAccessRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setDocumentId(value);
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
proto.services.docstore.GetDocumentAccessRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.GetDocumentAccessRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.GetDocumentAccessRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentAccessRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDocumentId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 document_id = 1;
 * @return {number}
 */
proto.services.docstore.GetDocumentAccessRequest.prototype.getDocumentId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.GetDocumentAccessRequest} returns this
 */
proto.services.docstore.GetDocumentAccessRequest.prototype.setDocumentId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.GetDocumentAccessResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.GetDocumentAccessResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.GetDocumentAccessResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentAccessResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    access: (f = msg.getAccess()) && resources_documents_documents_pb.DocumentAccess.toObject(includeInstance, f)
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
 * @return {!proto.services.docstore.GetDocumentAccessResponse}
 */
proto.services.docstore.GetDocumentAccessResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.GetDocumentAccessResponse;
  return proto.services.docstore.GetDocumentAccessResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.GetDocumentAccessResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.GetDocumentAccessResponse}
 */
proto.services.docstore.GetDocumentAccessResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_documents_pb.DocumentAccess;
      reader.readMessage(value,resources_documents_documents_pb.DocumentAccess.deserializeBinaryFromReader);
      msg.setAccess(value);
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
proto.services.docstore.GetDocumentAccessResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.GetDocumentAccessResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.GetDocumentAccessResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.GetDocumentAccessResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAccess();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_documents_documents_pb.DocumentAccess.serializeBinaryToWriter
    );
  }
};


/**
 * optional resources.documents.DocumentAccess access = 1;
 * @return {?proto.resources.documents.DocumentAccess}
 */
proto.services.docstore.GetDocumentAccessResponse.prototype.getAccess = function() {
  return /** @type{?proto.resources.documents.DocumentAccess} */ (
    jspb.Message.getWrapperField(this, resources_documents_documents_pb.DocumentAccess, 1));
};


/**
 * @param {?proto.resources.documents.DocumentAccess|undefined} value
 * @return {!proto.services.docstore.GetDocumentAccessResponse} returns this
*/
proto.services.docstore.GetDocumentAccessResponse.prototype.setAccess = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.GetDocumentAccessResponse} returns this
 */
proto.services.docstore.GetDocumentAccessResponse.prototype.clearAccess = function() {
  return this.setAccess(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.GetDocumentAccessResponse.prototype.hasAccess = function() {
  return jspb.Message.getField(this, 1) != null;
};





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
proto.services.docstore.SetDocumentAccessRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.SetDocumentAccessRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.SetDocumentAccessRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.SetDocumentAccessRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    documentId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    mode: jspb.Message.getFieldWithDefault(msg, 2, 0),
    access: (f = msg.getAccess()) && resources_documents_documents_pb.DocumentAccess.toObject(includeInstance, f)
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
 * @return {!proto.services.docstore.SetDocumentAccessRequest}
 */
proto.services.docstore.SetDocumentAccessRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.SetDocumentAccessRequest;
  return proto.services.docstore.SetDocumentAccessRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.SetDocumentAccessRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.SetDocumentAccessRequest}
 */
proto.services.docstore.SetDocumentAccessRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setDocumentId(value);
      break;
    case 2:
      var value = /** @type {!proto.services.docstore.ACCESS_LEVEL_UPDATE_MODE} */ (reader.readEnum());
      msg.setMode(value);
      break;
    case 3:
      var value = new resources_documents_documents_pb.DocumentAccess;
      reader.readMessage(value,resources_documents_documents_pb.DocumentAccess.deserializeBinaryFromReader);
      msg.setAccess(value);
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
proto.services.docstore.SetDocumentAccessRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.SetDocumentAccessRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.SetDocumentAccessRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.SetDocumentAccessRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDocumentId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
  f = message.getMode();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getAccess();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      resources_documents_documents_pb.DocumentAccess.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint64 document_id = 1;
 * @return {number}
 */
proto.services.docstore.SetDocumentAccessRequest.prototype.getDocumentId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.SetDocumentAccessRequest} returns this
 */
proto.services.docstore.SetDocumentAccessRequest.prototype.setDocumentId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional ACCESS_LEVEL_UPDATE_MODE mode = 2;
 * @return {!proto.services.docstore.ACCESS_LEVEL_UPDATE_MODE}
 */
proto.services.docstore.SetDocumentAccessRequest.prototype.getMode = function() {
  return /** @type {!proto.services.docstore.ACCESS_LEVEL_UPDATE_MODE} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.services.docstore.ACCESS_LEVEL_UPDATE_MODE} value
 * @return {!proto.services.docstore.SetDocumentAccessRequest} returns this
 */
proto.services.docstore.SetDocumentAccessRequest.prototype.setMode = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional resources.documents.DocumentAccess access = 3;
 * @return {?proto.resources.documents.DocumentAccess}
 */
proto.services.docstore.SetDocumentAccessRequest.prototype.getAccess = function() {
  return /** @type{?proto.resources.documents.DocumentAccess} */ (
    jspb.Message.getWrapperField(this, resources_documents_documents_pb.DocumentAccess, 3));
};


/**
 * @param {?proto.resources.documents.DocumentAccess|undefined} value
 * @return {!proto.services.docstore.SetDocumentAccessRequest} returns this
*/
proto.services.docstore.SetDocumentAccessRequest.prototype.setAccess = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.SetDocumentAccessRequest} returns this
 */
proto.services.docstore.SetDocumentAccessRequest.prototype.clearAccess = function() {
  return this.setAccess(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.SetDocumentAccessRequest.prototype.hasAccess = function() {
  return jspb.Message.getField(this, 3) != null;
};





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
proto.services.docstore.SetDocumentAccessResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.SetDocumentAccessResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.SetDocumentAccessResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.SetDocumentAccessResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.services.docstore.SetDocumentAccessResponse}
 */
proto.services.docstore.SetDocumentAccessResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.SetDocumentAccessResponse;
  return proto.services.docstore.SetDocumentAccessResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.SetDocumentAccessResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.SetDocumentAccessResponse}
 */
proto.services.docstore.SetDocumentAccessResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.SetDocumentAccessResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.SetDocumentAccessResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.SetDocumentAccessResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.SetDocumentAccessResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.services.docstore.ListUserDocumentsRequest.repeatedFields_ = [3];



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
proto.services.docstore.ListUserDocumentsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.ListUserDocumentsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.ListUserDocumentsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListUserDocumentsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    pagination: (f = msg.getPagination()) && resources_common_database_database_pb.PaginationRequest.toObject(includeInstance, f),
    userId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    relationsList: (f = jspb.Message.getRepeatedField(msg, 3)) == null ? undefined : f
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
 * @return {!proto.services.docstore.ListUserDocumentsRequest}
 */
proto.services.docstore.ListUserDocumentsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.ListUserDocumentsRequest;
  return proto.services.docstore.ListUserDocumentsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.ListUserDocumentsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.ListUserDocumentsRequest}
 */
proto.services.docstore.ListUserDocumentsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_common_database_database_pb.PaginationRequest;
      reader.readMessage(value,resources_common_database_database_pb.PaginationRequest.deserializeBinaryFromReader);
      msg.setPagination(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setUserId(value);
      break;
    case 3:
      var values = /** @type {!Array<!proto.resources.documents.DOC_RELATION>} */ (reader.isDelimited() ? reader.readPackedEnum() : [reader.readEnum()]);
      for (var i = 0; i < values.length; i++) {
        msg.addRelations(values[i]);
      }
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
proto.services.docstore.ListUserDocumentsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.ListUserDocumentsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.ListUserDocumentsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListUserDocumentsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPagination();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_common_database_database_pb.PaginationRequest.serializeBinaryToWriter
    );
  }
  f = message.getUserId();
  if (f !== 0) {
    writer.writeInt32(
      2,
      f
    );
  }
  f = message.getRelationsList();
  if (f.length > 0) {
    writer.writePackedEnum(
      3,
      f
    );
  }
};


/**
 * optional resources.common.database.PaginationRequest pagination = 1;
 * @return {?proto.resources.common.database.PaginationRequest}
 */
proto.services.docstore.ListUserDocumentsRequest.prototype.getPagination = function() {
  return /** @type{?proto.resources.common.database.PaginationRequest} */ (
    jspb.Message.getWrapperField(this, resources_common_database_database_pb.PaginationRequest, 1));
};


/**
 * @param {?proto.resources.common.database.PaginationRequest|undefined} value
 * @return {!proto.services.docstore.ListUserDocumentsRequest} returns this
*/
proto.services.docstore.ListUserDocumentsRequest.prototype.setPagination = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.ListUserDocumentsRequest} returns this
 */
proto.services.docstore.ListUserDocumentsRequest.prototype.clearPagination = function() {
  return this.setPagination(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.ListUserDocumentsRequest.prototype.hasPagination = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional int32 user_id = 2;
 * @return {number}
 */
proto.services.docstore.ListUserDocumentsRequest.prototype.getUserId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.ListUserDocumentsRequest} returns this
 */
proto.services.docstore.ListUserDocumentsRequest.prototype.setUserId = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * repeated resources.documents.DOC_RELATION relations = 3;
 * @return {!Array<!proto.resources.documents.DOC_RELATION>}
 */
proto.services.docstore.ListUserDocumentsRequest.prototype.getRelationsList = function() {
  return /** @type {!Array<!proto.resources.documents.DOC_RELATION>} */ (jspb.Message.getRepeatedField(this, 3));
};


/**
 * @param {!Array<!proto.resources.documents.DOC_RELATION>} value
 * @return {!proto.services.docstore.ListUserDocumentsRequest} returns this
 */
proto.services.docstore.ListUserDocumentsRequest.prototype.setRelationsList = function(value) {
  return jspb.Message.setField(this, 3, value || []);
};


/**
 * @param {!proto.resources.documents.DOC_RELATION} value
 * @param {number=} opt_index
 * @return {!proto.services.docstore.ListUserDocumentsRequest} returns this
 */
proto.services.docstore.ListUserDocumentsRequest.prototype.addRelations = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 3, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.services.docstore.ListUserDocumentsRequest} returns this
 */
proto.services.docstore.ListUserDocumentsRequest.prototype.clearRelationsList = function() {
  return this.setRelationsList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.services.docstore.ListUserDocumentsResponse.repeatedFields_ = [2];



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
proto.services.docstore.ListUserDocumentsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.ListUserDocumentsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.ListUserDocumentsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListUserDocumentsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    pagination: (f = msg.getPagination()) && resources_common_database_database_pb.PaginationResponse.toObject(includeInstance, f),
    relationsList: jspb.Message.toObjectList(msg.getRelationsList(),
    resources_documents_documents_pb.DocumentRelation.toObject, includeInstance)
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
 * @return {!proto.services.docstore.ListUserDocumentsResponse}
 */
proto.services.docstore.ListUserDocumentsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.ListUserDocumentsResponse;
  return proto.services.docstore.ListUserDocumentsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.ListUserDocumentsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.ListUserDocumentsResponse}
 */
proto.services.docstore.ListUserDocumentsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_common_database_database_pb.PaginationResponse;
      reader.readMessage(value,resources_common_database_database_pb.PaginationResponse.deserializeBinaryFromReader);
      msg.setPagination(value);
      break;
    case 2:
      var value = new resources_documents_documents_pb.DocumentRelation;
      reader.readMessage(value,resources_documents_documents_pb.DocumentRelation.deserializeBinaryFromReader);
      msg.addRelations(value);
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
proto.services.docstore.ListUserDocumentsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.ListUserDocumentsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.ListUserDocumentsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListUserDocumentsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPagination();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_common_database_database_pb.PaginationResponse.serializeBinaryToWriter
    );
  }
  f = message.getRelationsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      resources_documents_documents_pb.DocumentRelation.serializeBinaryToWriter
    );
  }
};


/**
 * optional resources.common.database.PaginationResponse pagination = 1;
 * @return {?proto.resources.common.database.PaginationResponse}
 */
proto.services.docstore.ListUserDocumentsResponse.prototype.getPagination = function() {
  return /** @type{?proto.resources.common.database.PaginationResponse} */ (
    jspb.Message.getWrapperField(this, resources_common_database_database_pb.PaginationResponse, 1));
};


/**
 * @param {?proto.resources.common.database.PaginationResponse|undefined} value
 * @return {!proto.services.docstore.ListUserDocumentsResponse} returns this
*/
proto.services.docstore.ListUserDocumentsResponse.prototype.setPagination = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.ListUserDocumentsResponse} returns this
 */
proto.services.docstore.ListUserDocumentsResponse.prototype.clearPagination = function() {
  return this.setPagination(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.ListUserDocumentsResponse.prototype.hasPagination = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * repeated resources.documents.DocumentRelation relations = 2;
 * @return {!Array<!proto.resources.documents.DocumentRelation>}
 */
proto.services.docstore.ListUserDocumentsResponse.prototype.getRelationsList = function() {
  return /** @type{!Array<!proto.resources.documents.DocumentRelation>} */ (
    jspb.Message.getRepeatedWrapperField(this, resources_documents_documents_pb.DocumentRelation, 2));
};


/**
 * @param {!Array<!proto.resources.documents.DocumentRelation>} value
 * @return {!proto.services.docstore.ListUserDocumentsResponse} returns this
*/
proto.services.docstore.ListUserDocumentsResponse.prototype.setRelationsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.resources.documents.DocumentRelation=} opt_value
 * @param {number=} opt_index
 * @return {!proto.resources.documents.DocumentRelation}
 */
proto.services.docstore.ListUserDocumentsResponse.prototype.addRelations = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.resources.documents.DocumentRelation, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.services.docstore.ListUserDocumentsResponse} returns this
 */
proto.services.docstore.ListUserDocumentsResponse.prototype.clearRelationsList = function() {
  return this.setRelationsList([]);
};





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
proto.services.docstore.ListDocumentCategoriesRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.ListDocumentCategoriesRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.ListDocumentCategoriesRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListDocumentCategoriesRequest.toObject = function(includeInstance, msg) {
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
 * @return {!proto.services.docstore.ListDocumentCategoriesRequest}
 */
proto.services.docstore.ListDocumentCategoriesRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.ListDocumentCategoriesRequest;
  return proto.services.docstore.ListDocumentCategoriesRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.ListDocumentCategoriesRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.ListDocumentCategoriesRequest}
 */
proto.services.docstore.ListDocumentCategoriesRequest.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.ListDocumentCategoriesRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.ListDocumentCategoriesRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.ListDocumentCategoriesRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListDocumentCategoriesRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.services.docstore.ListDocumentCategoriesResponse.repeatedFields_ = [1];



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
proto.services.docstore.ListDocumentCategoriesResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.ListDocumentCategoriesResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.ListDocumentCategoriesResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListDocumentCategoriesResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    categoryList: jspb.Message.toObjectList(msg.getCategoryList(),
    resources_documents_category_pb.DocumentCategory.toObject, includeInstance)
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
 * @return {!proto.services.docstore.ListDocumentCategoriesResponse}
 */
proto.services.docstore.ListDocumentCategoriesResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.ListDocumentCategoriesResponse;
  return proto.services.docstore.ListDocumentCategoriesResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.ListDocumentCategoriesResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.ListDocumentCategoriesResponse}
 */
proto.services.docstore.ListDocumentCategoriesResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_category_pb.DocumentCategory;
      reader.readMessage(value,resources_documents_category_pb.DocumentCategory.deserializeBinaryFromReader);
      msg.addCategory(value);
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
proto.services.docstore.ListDocumentCategoriesResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.ListDocumentCategoriesResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.ListDocumentCategoriesResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.ListDocumentCategoriesResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCategoryList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      resources_documents_category_pb.DocumentCategory.serializeBinaryToWriter
    );
  }
};


/**
 * repeated resources.documents.DocumentCategory category = 1;
 * @return {!Array<!proto.resources.documents.DocumentCategory>}
 */
proto.services.docstore.ListDocumentCategoriesResponse.prototype.getCategoryList = function() {
  return /** @type{!Array<!proto.resources.documents.DocumentCategory>} */ (
    jspb.Message.getRepeatedWrapperField(this, resources_documents_category_pb.DocumentCategory, 1));
};


/**
 * @param {!Array<!proto.resources.documents.DocumentCategory>} value
 * @return {!proto.services.docstore.ListDocumentCategoriesResponse} returns this
*/
proto.services.docstore.ListDocumentCategoriesResponse.prototype.setCategoryList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.resources.documents.DocumentCategory=} opt_value
 * @param {number=} opt_index
 * @return {!proto.resources.documents.DocumentCategory}
 */
proto.services.docstore.ListDocumentCategoriesResponse.prototype.addCategory = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.resources.documents.DocumentCategory, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.services.docstore.ListDocumentCategoriesResponse} returns this
 */
proto.services.docstore.ListDocumentCategoriesResponse.prototype.clearCategoryList = function() {
  return this.setCategoryList([]);
};





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
proto.services.docstore.CreateDocumentCategoryRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.CreateDocumentCategoryRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.CreateDocumentCategoryRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.CreateDocumentCategoryRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    category: (f = msg.getCategory()) && resources_documents_category_pb.DocumentCategory.toObject(includeInstance, f)
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
 * @return {!proto.services.docstore.CreateDocumentCategoryRequest}
 */
proto.services.docstore.CreateDocumentCategoryRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.CreateDocumentCategoryRequest;
  return proto.services.docstore.CreateDocumentCategoryRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.CreateDocumentCategoryRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.CreateDocumentCategoryRequest}
 */
proto.services.docstore.CreateDocumentCategoryRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_category_pb.DocumentCategory;
      reader.readMessage(value,resources_documents_category_pb.DocumentCategory.deserializeBinaryFromReader);
      msg.setCategory(value);
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
proto.services.docstore.CreateDocumentCategoryRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.CreateDocumentCategoryRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.CreateDocumentCategoryRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.CreateDocumentCategoryRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCategory();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_documents_category_pb.DocumentCategory.serializeBinaryToWriter
    );
  }
};


/**
 * optional resources.documents.DocumentCategory category = 1;
 * @return {?proto.resources.documents.DocumentCategory}
 */
proto.services.docstore.CreateDocumentCategoryRequest.prototype.getCategory = function() {
  return /** @type{?proto.resources.documents.DocumentCategory} */ (
    jspb.Message.getWrapperField(this, resources_documents_category_pb.DocumentCategory, 1));
};


/**
 * @param {?proto.resources.documents.DocumentCategory|undefined} value
 * @return {!proto.services.docstore.CreateDocumentCategoryRequest} returns this
*/
proto.services.docstore.CreateDocumentCategoryRequest.prototype.setCategory = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.CreateDocumentCategoryRequest} returns this
 */
proto.services.docstore.CreateDocumentCategoryRequest.prototype.clearCategory = function() {
  return this.setCategory(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.CreateDocumentCategoryRequest.prototype.hasCategory = function() {
  return jspb.Message.getField(this, 1) != null;
};





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
proto.services.docstore.CreateDocumentCategoryResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.CreateDocumentCategoryResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.CreateDocumentCategoryResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.CreateDocumentCategoryResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, 0)
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
 * @return {!proto.services.docstore.CreateDocumentCategoryResponse}
 */
proto.services.docstore.CreateDocumentCategoryResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.CreateDocumentCategoryResponse;
  return proto.services.docstore.CreateDocumentCategoryResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.CreateDocumentCategoryResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.CreateDocumentCategoryResponse}
 */
proto.services.docstore.CreateDocumentCategoryResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.CreateDocumentCategoryResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.CreateDocumentCategoryResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.CreateDocumentCategoryResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.CreateDocumentCategoryResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
};


/**
 * optional uint64 id = 1;
 * @return {number}
 */
proto.services.docstore.CreateDocumentCategoryResponse.prototype.getId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.services.docstore.CreateDocumentCategoryResponse} returns this
 */
proto.services.docstore.CreateDocumentCategoryResponse.prototype.setId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





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
proto.services.docstore.UpdateDocumentCategoryRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.UpdateDocumentCategoryRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.UpdateDocumentCategoryRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.UpdateDocumentCategoryRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    category: (f = msg.getCategory()) && resources_documents_category_pb.DocumentCategory.toObject(includeInstance, f)
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
 * @return {!proto.services.docstore.UpdateDocumentCategoryRequest}
 */
proto.services.docstore.UpdateDocumentCategoryRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.UpdateDocumentCategoryRequest;
  return proto.services.docstore.UpdateDocumentCategoryRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.UpdateDocumentCategoryRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.UpdateDocumentCategoryRequest}
 */
proto.services.docstore.UpdateDocumentCategoryRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new resources_documents_category_pb.DocumentCategory;
      reader.readMessage(value,resources_documents_category_pb.DocumentCategory.deserializeBinaryFromReader);
      msg.setCategory(value);
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
proto.services.docstore.UpdateDocumentCategoryRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.UpdateDocumentCategoryRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.UpdateDocumentCategoryRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.UpdateDocumentCategoryRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCategory();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      resources_documents_category_pb.DocumentCategory.serializeBinaryToWriter
    );
  }
};


/**
 * optional resources.documents.DocumentCategory category = 1;
 * @return {?proto.resources.documents.DocumentCategory}
 */
proto.services.docstore.UpdateDocumentCategoryRequest.prototype.getCategory = function() {
  return /** @type{?proto.resources.documents.DocumentCategory} */ (
    jspb.Message.getWrapperField(this, resources_documents_category_pb.DocumentCategory, 1));
};


/**
 * @param {?proto.resources.documents.DocumentCategory|undefined} value
 * @return {!proto.services.docstore.UpdateDocumentCategoryRequest} returns this
*/
proto.services.docstore.UpdateDocumentCategoryRequest.prototype.setCategory = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.services.docstore.UpdateDocumentCategoryRequest} returns this
 */
proto.services.docstore.UpdateDocumentCategoryRequest.prototype.clearCategory = function() {
  return this.setCategory(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.services.docstore.UpdateDocumentCategoryRequest.prototype.hasCategory = function() {
  return jspb.Message.getField(this, 1) != null;
};





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
proto.services.docstore.UpdateDocumentCategoryResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.UpdateDocumentCategoryResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.UpdateDocumentCategoryResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.UpdateDocumentCategoryResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.services.docstore.UpdateDocumentCategoryResponse}
 */
proto.services.docstore.UpdateDocumentCategoryResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.UpdateDocumentCategoryResponse;
  return proto.services.docstore.UpdateDocumentCategoryResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.UpdateDocumentCategoryResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.UpdateDocumentCategoryResponse}
 */
proto.services.docstore.UpdateDocumentCategoryResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.UpdateDocumentCategoryResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.UpdateDocumentCategoryResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.UpdateDocumentCategoryResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.UpdateDocumentCategoryResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.services.docstore.DeleteDocumentCategoryRequest.repeatedFields_ = [1];



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
proto.services.docstore.DeleteDocumentCategoryRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.DeleteDocumentCategoryRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.DeleteDocumentCategoryRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteDocumentCategoryRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    idsList: (f = jspb.Message.getRepeatedField(msg, 1)) == null ? undefined : f
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
 * @return {!proto.services.docstore.DeleteDocumentCategoryRequest}
 */
proto.services.docstore.DeleteDocumentCategoryRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.DeleteDocumentCategoryRequest;
  return proto.services.docstore.DeleteDocumentCategoryRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.DeleteDocumentCategoryRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.DeleteDocumentCategoryRequest}
 */
proto.services.docstore.DeleteDocumentCategoryRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var values = /** @type {!Array<number>} */ (reader.isDelimited() ? reader.readPackedUint64() : [reader.readUint64()]);
      for (var i = 0; i < values.length; i++) {
        msg.addIds(values[i]);
      }
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
proto.services.docstore.DeleteDocumentCategoryRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.DeleteDocumentCategoryRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.DeleteDocumentCategoryRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteDocumentCategoryRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getIdsList();
  if (f.length > 0) {
    writer.writePackedUint64(
      1,
      f
    );
  }
};


/**
 * repeated uint64 ids = 1;
 * @return {!Array<number>}
 */
proto.services.docstore.DeleteDocumentCategoryRequest.prototype.getIdsList = function() {
  return /** @type {!Array<number>} */ (jspb.Message.getRepeatedField(this, 1));
};


/**
 * @param {!Array<number>} value
 * @return {!proto.services.docstore.DeleteDocumentCategoryRequest} returns this
 */
proto.services.docstore.DeleteDocumentCategoryRequest.prototype.setIdsList = function(value) {
  return jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {number} value
 * @param {number=} opt_index
 * @return {!proto.services.docstore.DeleteDocumentCategoryRequest} returns this
 */
proto.services.docstore.DeleteDocumentCategoryRequest.prototype.addIds = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.services.docstore.DeleteDocumentCategoryRequest} returns this
 */
proto.services.docstore.DeleteDocumentCategoryRequest.prototype.clearIdsList = function() {
  return this.setIdsList([]);
};





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
proto.services.docstore.DeleteDocumentCategoryResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.services.docstore.DeleteDocumentCategoryResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.services.docstore.DeleteDocumentCategoryResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteDocumentCategoryResponse.toObject = function(includeInstance, msg) {
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
 * @return {!proto.services.docstore.DeleteDocumentCategoryResponse}
 */
proto.services.docstore.DeleteDocumentCategoryResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.services.docstore.DeleteDocumentCategoryResponse;
  return proto.services.docstore.DeleteDocumentCategoryResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.services.docstore.DeleteDocumentCategoryResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.services.docstore.DeleteDocumentCategoryResponse}
 */
proto.services.docstore.DeleteDocumentCategoryResponse.deserializeBinaryFromReader = function(msg, reader) {
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
proto.services.docstore.DeleteDocumentCategoryResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.services.docstore.DeleteDocumentCategoryResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.services.docstore.DeleteDocumentCategoryResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.services.docstore.DeleteDocumentCategoryResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};


/**
 * @enum {number}
 */
proto.services.docstore.ACCESS_LEVEL_UPDATE_MODE = {
  UPDATE: 0,
  DELETE: 1,
  CLEAR: 2
};

goog.object.extend(exports, proto.services.docstore);
