/**
 * @fileoverview gRPC-Web generated client stub for services.docstore
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.4.2
// 	protoc              v3.21.12
// source: services/docstore/docstore.proto


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as services_docstore_docstore_pb from '../../services/docstore/docstore_pb';


export class DocStoreServiceClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: any; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'text';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname.replace(/\/+$/, '');
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodDescriptorListTemplates = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/ListTemplates',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.ListTemplatesRequest,
    services_docstore_docstore_pb.ListTemplatesResponse,
    (request: services_docstore_docstore_pb.ListTemplatesRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.ListTemplatesResponse.deserializeBinary
  );

  listTemplates(
    request: services_docstore_docstore_pb.ListTemplatesRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.ListTemplatesResponse>;

  listTemplates(
    request: services_docstore_docstore_pb.ListTemplatesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.ListTemplatesResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.ListTemplatesResponse>;

  listTemplates(
    request: services_docstore_docstore_pb.ListTemplatesRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.ListTemplatesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/ListTemplates',
        request,
        metadata || {},
        this.methodDescriptorListTemplates,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/ListTemplates',
    request,
    metadata || {},
    this.methodDescriptorListTemplates);
  }

  methodDescriptorGetTemplate = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/GetTemplate',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.GetTemplateRequest,
    services_docstore_docstore_pb.GetTemplateResponse,
    (request: services_docstore_docstore_pb.GetTemplateRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.GetTemplateResponse.deserializeBinary
  );

  getTemplate(
    request: services_docstore_docstore_pb.GetTemplateRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.GetTemplateResponse>;

  getTemplate(
    request: services_docstore_docstore_pb.GetTemplateRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetTemplateResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.GetTemplateResponse>;

  getTemplate(
    request: services_docstore_docstore_pb.GetTemplateRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetTemplateResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/GetTemplate',
        request,
        metadata || {},
        this.methodDescriptorGetTemplate,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/GetTemplate',
    request,
    metadata || {},
    this.methodDescriptorGetTemplate);
  }

  methodDescriptorFindDocuments = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/FindDocuments',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.FindDocumentsRequest,
    services_docstore_docstore_pb.FindDocumentsResponse,
    (request: services_docstore_docstore_pb.FindDocumentsRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.FindDocumentsResponse.deserializeBinary
  );

  findDocuments(
    request: services_docstore_docstore_pb.FindDocumentsRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.FindDocumentsResponse>;

  findDocuments(
    request: services_docstore_docstore_pb.FindDocumentsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.FindDocumentsResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.FindDocumentsResponse>;

  findDocuments(
    request: services_docstore_docstore_pb.FindDocumentsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.FindDocumentsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/FindDocuments',
        request,
        metadata || {},
        this.methodDescriptorFindDocuments,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/FindDocuments',
    request,
    metadata || {},
    this.methodDescriptorFindDocuments);
  }

  methodDescriptorGetDocument = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/GetDocument',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.GetDocumentRequest,
    services_docstore_docstore_pb.GetDocumentResponse,
    (request: services_docstore_docstore_pb.GetDocumentRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.GetDocumentResponse.deserializeBinary
  );

  getDocument(
    request: services_docstore_docstore_pb.GetDocumentRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.GetDocumentResponse>;

  getDocument(
    request: services_docstore_docstore_pb.GetDocumentRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetDocumentResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.GetDocumentResponse>;

  getDocument(
    request: services_docstore_docstore_pb.GetDocumentRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetDocumentResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/GetDocument',
        request,
        metadata || {},
        this.methodDescriptorGetDocument,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/GetDocument',
    request,
    metadata || {},
    this.methodDescriptorGetDocument);
  }

  methodDescriptorCreateDocument = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/CreateDocument',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.CreateDocumentRequest,
    services_docstore_docstore_pb.CreateDocumentResponse,
    (request: services_docstore_docstore_pb.CreateDocumentRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.CreateDocumentResponse.deserializeBinary
  );

  createDocument(
    request: services_docstore_docstore_pb.CreateDocumentRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.CreateDocumentResponse>;

  createDocument(
    request: services_docstore_docstore_pb.CreateDocumentRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.CreateDocumentResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.CreateDocumentResponse>;

  createDocument(
    request: services_docstore_docstore_pb.CreateDocumentRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.CreateDocumentResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/CreateDocument',
        request,
        metadata || {},
        this.methodDescriptorCreateDocument,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/CreateDocument',
    request,
    metadata || {},
    this.methodDescriptorCreateDocument);
  }

  methodDescriptorUpdateDocument = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/UpdateDocument',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.UpdateDocumentRequest,
    services_docstore_docstore_pb.UpdateDocumentResponse,
    (request: services_docstore_docstore_pb.UpdateDocumentRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.UpdateDocumentResponse.deserializeBinary
  );

  updateDocument(
    request: services_docstore_docstore_pb.UpdateDocumentRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.UpdateDocumentResponse>;

  updateDocument(
    request: services_docstore_docstore_pb.UpdateDocumentRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.UpdateDocumentResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.UpdateDocumentResponse>;

  updateDocument(
    request: services_docstore_docstore_pb.UpdateDocumentRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.UpdateDocumentResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/UpdateDocument',
        request,
        metadata || {},
        this.methodDescriptorUpdateDocument,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/UpdateDocument',
    request,
    metadata || {},
    this.methodDescriptorUpdateDocument);
  }

  methodDescriptorGetDocumentReferences = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/GetDocumentReferences',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.GetDocumentReferencesRequest,
    services_docstore_docstore_pb.GetDocumentReferencesResponse,
    (request: services_docstore_docstore_pb.GetDocumentReferencesRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.GetDocumentReferencesResponse.deserializeBinary
  );

  getDocumentReferences(
    request: services_docstore_docstore_pb.GetDocumentReferencesRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.GetDocumentReferencesResponse>;

  getDocumentReferences(
    request: services_docstore_docstore_pb.GetDocumentReferencesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetDocumentReferencesResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.GetDocumentReferencesResponse>;

  getDocumentReferences(
    request: services_docstore_docstore_pb.GetDocumentReferencesRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetDocumentReferencesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/GetDocumentReferences',
        request,
        metadata || {},
        this.methodDescriptorGetDocumentReferences,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/GetDocumentReferences',
    request,
    metadata || {},
    this.methodDescriptorGetDocumentReferences);
  }

  methodDescriptorGetDocumentRelations = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/GetDocumentRelations',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.GetDocumentRelationsRequest,
    services_docstore_docstore_pb.GetDocumentRelationsResponse,
    (request: services_docstore_docstore_pb.GetDocumentRelationsRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.GetDocumentRelationsResponse.deserializeBinary
  );

  getDocumentRelations(
    request: services_docstore_docstore_pb.GetDocumentRelationsRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.GetDocumentRelationsResponse>;

  getDocumentRelations(
    request: services_docstore_docstore_pb.GetDocumentRelationsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetDocumentRelationsResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.GetDocumentRelationsResponse>;

  getDocumentRelations(
    request: services_docstore_docstore_pb.GetDocumentRelationsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetDocumentRelationsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/GetDocumentRelations',
        request,
        metadata || {},
        this.methodDescriptorGetDocumentRelations,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/GetDocumentRelations',
    request,
    metadata || {},
    this.methodDescriptorGetDocumentRelations);
  }

  methodDescriptorAddDocumentReference = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/AddDocumentReference',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.AddDocumentReferenceRequest,
    services_docstore_docstore_pb.AddDocumentReferenceResponse,
    (request: services_docstore_docstore_pb.AddDocumentReferenceRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.AddDocumentReferenceResponse.deserializeBinary
  );

  addDocumentReference(
    request: services_docstore_docstore_pb.AddDocumentReferenceRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.AddDocumentReferenceResponse>;

  addDocumentReference(
    request: services_docstore_docstore_pb.AddDocumentReferenceRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.AddDocumentReferenceResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.AddDocumentReferenceResponse>;

  addDocumentReference(
    request: services_docstore_docstore_pb.AddDocumentReferenceRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.AddDocumentReferenceResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/AddDocumentReference',
        request,
        metadata || {},
        this.methodDescriptorAddDocumentReference,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/AddDocumentReference',
    request,
    metadata || {},
    this.methodDescriptorAddDocumentReference);
  }

  methodDescriptorRemoveDocumentReference = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/RemoveDocumentReference',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.RemoveDocumentReferenceRequest,
    services_docstore_docstore_pb.RemoveDocumentReferenceResponse,
    (request: services_docstore_docstore_pb.RemoveDocumentReferenceRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.RemoveDocumentReferenceResponse.deserializeBinary
  );

  removeDocumentReference(
    request: services_docstore_docstore_pb.RemoveDocumentReferenceRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.RemoveDocumentReferenceResponse>;

  removeDocumentReference(
    request: services_docstore_docstore_pb.RemoveDocumentReferenceRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.RemoveDocumentReferenceResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.RemoveDocumentReferenceResponse>;

  removeDocumentReference(
    request: services_docstore_docstore_pb.RemoveDocumentReferenceRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.RemoveDocumentReferenceResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/RemoveDocumentReference',
        request,
        metadata || {},
        this.methodDescriptorRemoveDocumentReference,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/RemoveDocumentReference',
    request,
    metadata || {},
    this.methodDescriptorRemoveDocumentReference);
  }

  methodDescriptorAddDocumentRelation = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/AddDocumentRelation',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.AddDocumentRelationRequest,
    services_docstore_docstore_pb.AddDocumentRelationResponse,
    (request: services_docstore_docstore_pb.AddDocumentRelationRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.AddDocumentRelationResponse.deserializeBinary
  );

  addDocumentRelation(
    request: services_docstore_docstore_pb.AddDocumentRelationRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.AddDocumentRelationResponse>;

  addDocumentRelation(
    request: services_docstore_docstore_pb.AddDocumentRelationRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.AddDocumentRelationResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.AddDocumentRelationResponse>;

  addDocumentRelation(
    request: services_docstore_docstore_pb.AddDocumentRelationRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.AddDocumentRelationResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/AddDocumentRelation',
        request,
        metadata || {},
        this.methodDescriptorAddDocumentRelation,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/AddDocumentRelation',
    request,
    metadata || {},
    this.methodDescriptorAddDocumentRelation);
  }

  methodDescriptorRemoveDocumentRelation = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/RemoveDocumentRelation',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.RemoveDocumentRelationRequest,
    services_docstore_docstore_pb.RemoveDocumentRelationResponse,
    (request: services_docstore_docstore_pb.RemoveDocumentRelationRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.RemoveDocumentRelationResponse.deserializeBinary
  );

  removeDocumentRelation(
    request: services_docstore_docstore_pb.RemoveDocumentRelationRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.RemoveDocumentRelationResponse>;

  removeDocumentRelation(
    request: services_docstore_docstore_pb.RemoveDocumentRelationRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.RemoveDocumentRelationResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.RemoveDocumentRelationResponse>;

  removeDocumentRelation(
    request: services_docstore_docstore_pb.RemoveDocumentRelationRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.RemoveDocumentRelationResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/RemoveDocumentRelation',
        request,
        metadata || {},
        this.methodDescriptorRemoveDocumentRelation,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/RemoveDocumentRelation',
    request,
    metadata || {},
    this.methodDescriptorRemoveDocumentRelation);
  }

  methodDescriptorGetDocumentComments = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/GetDocumentComments',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.GetDocumentCommentsRequest,
    services_docstore_docstore_pb.GetDocumentCommentsResponse,
    (request: services_docstore_docstore_pb.GetDocumentCommentsRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.GetDocumentCommentsResponse.deserializeBinary
  );

  getDocumentComments(
    request: services_docstore_docstore_pb.GetDocumentCommentsRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.GetDocumentCommentsResponse>;

  getDocumentComments(
    request: services_docstore_docstore_pb.GetDocumentCommentsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetDocumentCommentsResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.GetDocumentCommentsResponse>;

  getDocumentComments(
    request: services_docstore_docstore_pb.GetDocumentCommentsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetDocumentCommentsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/GetDocumentComments',
        request,
        metadata || {},
        this.methodDescriptorGetDocumentComments,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/GetDocumentComments',
    request,
    metadata || {},
    this.methodDescriptorGetDocumentComments);
  }

  methodDescriptorPostDocumentComment = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/PostDocumentComment',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.PostDocumentCommentRequest,
    services_docstore_docstore_pb.PostDocumentCommentResponse,
    (request: services_docstore_docstore_pb.PostDocumentCommentRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.PostDocumentCommentResponse.deserializeBinary
  );

  postDocumentComment(
    request: services_docstore_docstore_pb.PostDocumentCommentRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.PostDocumentCommentResponse>;

  postDocumentComment(
    request: services_docstore_docstore_pb.PostDocumentCommentRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.PostDocumentCommentResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.PostDocumentCommentResponse>;

  postDocumentComment(
    request: services_docstore_docstore_pb.PostDocumentCommentRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.PostDocumentCommentResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/PostDocumentComment',
        request,
        metadata || {},
        this.methodDescriptorPostDocumentComment,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/PostDocumentComment',
    request,
    metadata || {},
    this.methodDescriptorPostDocumentComment);
  }

  methodDescriptorEditDocumentComment = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/EditDocumentComment',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.EditDocumentCommentRequest,
    services_docstore_docstore_pb.EditDocumentCommentResponse,
    (request: services_docstore_docstore_pb.EditDocumentCommentRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.EditDocumentCommentResponse.deserializeBinary
  );

  editDocumentComment(
    request: services_docstore_docstore_pb.EditDocumentCommentRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.EditDocumentCommentResponse>;

  editDocumentComment(
    request: services_docstore_docstore_pb.EditDocumentCommentRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.EditDocumentCommentResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.EditDocumentCommentResponse>;

  editDocumentComment(
    request: services_docstore_docstore_pb.EditDocumentCommentRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.EditDocumentCommentResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/EditDocumentComment',
        request,
        metadata || {},
        this.methodDescriptorEditDocumentComment,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/EditDocumentComment',
    request,
    metadata || {},
    this.methodDescriptorEditDocumentComment);
  }

  methodDescriptorGetDocumentAccess = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/GetDocumentAccess',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.GetDocumentAccessRequest,
    services_docstore_docstore_pb.GetDocumentAccessResponse,
    (request: services_docstore_docstore_pb.GetDocumentAccessRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.GetDocumentAccessResponse.deserializeBinary
  );

  getDocumentAccess(
    request: services_docstore_docstore_pb.GetDocumentAccessRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.GetDocumentAccessResponse>;

  getDocumentAccess(
    request: services_docstore_docstore_pb.GetDocumentAccessRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetDocumentAccessResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.GetDocumentAccessResponse>;

  getDocumentAccess(
    request: services_docstore_docstore_pb.GetDocumentAccessRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetDocumentAccessResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/GetDocumentAccess',
        request,
        metadata || {},
        this.methodDescriptorGetDocumentAccess,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/GetDocumentAccess',
    request,
    metadata || {},
    this.methodDescriptorGetDocumentAccess);
  }

  methodDescriptorSetDocumentAccess = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/SetDocumentAccess',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.SetDocumentAccessRequest,
    services_docstore_docstore_pb.SetDocumentAccessResponse,
    (request: services_docstore_docstore_pb.SetDocumentAccessRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.SetDocumentAccessResponse.deserializeBinary
  );

  setDocumentAccess(
    request: services_docstore_docstore_pb.SetDocumentAccessRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.SetDocumentAccessResponse>;

  setDocumentAccess(
    request: services_docstore_docstore_pb.SetDocumentAccessRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.SetDocumentAccessResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.SetDocumentAccessResponse>;

  setDocumentAccess(
    request: services_docstore_docstore_pb.SetDocumentAccessRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.SetDocumentAccessResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/SetDocumentAccess',
        request,
        metadata || {},
        this.methodDescriptorSetDocumentAccess,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/SetDocumentAccess',
    request,
    metadata || {},
    this.methodDescriptorSetDocumentAccess);
  }

  methodDescriptorGetUserDocuments = new grpcWeb.MethodDescriptor(
    '/services.docstore.DocStoreService/GetUserDocuments',
    grpcWeb.MethodType.UNARY,
    services_docstore_docstore_pb.GetUserDocumentsRequest,
    services_docstore_docstore_pb.GetUserDocumentsResponse,
    (request: services_docstore_docstore_pb.GetUserDocumentsRequest) => {
      return request.serializeBinary();
    },
    services_docstore_docstore_pb.GetUserDocumentsResponse.deserializeBinary
  );

  getUserDocuments(
    request: services_docstore_docstore_pb.GetUserDocumentsRequest,
    metadata: grpcWeb.Metadata | null): Promise<services_docstore_docstore_pb.GetUserDocumentsResponse>;

  getUserDocuments(
    request: services_docstore_docstore_pb.GetUserDocumentsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetUserDocumentsResponse) => void): grpcWeb.ClientReadableStream<services_docstore_docstore_pb.GetUserDocumentsResponse>;

  getUserDocuments(
    request: services_docstore_docstore_pb.GetUserDocumentsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: services_docstore_docstore_pb.GetUserDocumentsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/services.docstore.DocStoreService/GetUserDocuments',
        request,
        metadata || {},
        this.methodDescriptorGetUserDocuments,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/services.docstore.DocStoreService/GetUserDocuments',
    request,
    metadata || {},
    this.methodDescriptorGetUserDocuments);
  }

}

