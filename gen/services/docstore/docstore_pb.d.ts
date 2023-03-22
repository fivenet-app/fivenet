import * as jspb from 'google-protobuf'

import * as resources_common_database_database_pb from '../../resources/common/database/database_pb';
import * as resources_documents_documents_pb from '../../resources/documents/documents_pb';


export class ListTemplatesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTemplatesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListTemplatesRequest): ListTemplatesRequest.AsObject;
  static serializeBinaryToWriter(message: ListTemplatesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTemplatesRequest;
  static deserializeBinaryFromReader(message: ListTemplatesRequest, reader: jspb.BinaryReader): ListTemplatesRequest;
}

export namespace ListTemplatesRequest {
  export type AsObject = {
  }
}

export class ListTemplatesResponse extends jspb.Message {
  getTemplatesList(): Array<resources_documents_documents_pb.DocumentTemplateShort>;
  setTemplatesList(value: Array<resources_documents_documents_pb.DocumentTemplateShort>): ListTemplatesResponse;
  clearTemplatesList(): ListTemplatesResponse;
  addTemplates(value?: resources_documents_documents_pb.DocumentTemplateShort, index?: number): resources_documents_documents_pb.DocumentTemplateShort;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTemplatesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListTemplatesResponse): ListTemplatesResponse.AsObject;
  static serializeBinaryToWriter(message: ListTemplatesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTemplatesResponse;
  static deserializeBinaryFromReader(message: ListTemplatesResponse, reader: jspb.BinaryReader): ListTemplatesResponse;
}

export namespace ListTemplatesResponse {
  export type AsObject = {
    templatesList: Array<resources_documents_documents_pb.DocumentTemplateShort.AsObject>,
  }
}

export class GetTemplateRequest extends jspb.Message {
  getTemplateId(): number;
  setTemplateId(value: number): GetTemplateRequest;

  getData(): string;
  setData(value: string): GetTemplateRequest;

  getProcess(): boolean;
  setProcess(value: boolean): GetTemplateRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTemplateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetTemplateRequest): GetTemplateRequest.AsObject;
  static serializeBinaryToWriter(message: GetTemplateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTemplateRequest;
  static deserializeBinaryFromReader(message: GetTemplateRequest, reader: jspb.BinaryReader): GetTemplateRequest;
}

export namespace GetTemplateRequest {
  export type AsObject = {
    templateId: number,
    data: string,
    process: boolean,
  }
}

export class GetTemplateResponse extends jspb.Message {
  getTemplate(): resources_documents_documents_pb.DocumentTemplate | undefined;
  setTemplate(value?: resources_documents_documents_pb.DocumentTemplate): GetTemplateResponse;
  hasTemplate(): boolean;
  clearTemplate(): GetTemplateResponse;

  getProcessed(): boolean;
  setProcessed(value: boolean): GetTemplateResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTemplateResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetTemplateResponse): GetTemplateResponse.AsObject;
  static serializeBinaryToWriter(message: GetTemplateResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTemplateResponse;
  static deserializeBinaryFromReader(message: GetTemplateResponse, reader: jspb.BinaryReader): GetTemplateResponse;
}

export namespace GetTemplateResponse {
  export type AsObject = {
    template?: resources_documents_documents_pb.DocumentTemplate.AsObject,
    processed: boolean,
  }
}

export class FindDocumentsRequest extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationRequest | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationRequest): FindDocumentsRequest;
  hasPagination(): boolean;
  clearPagination(): FindDocumentsRequest;

  getOrderbyList(): Array<resources_common_database_database_pb.OrderBy>;
  setOrderbyList(value: Array<resources_common_database_database_pb.OrderBy>): FindDocumentsRequest;
  clearOrderbyList(): FindDocumentsRequest;
  addOrderby(value?: resources_common_database_database_pb.OrderBy, index?: number): resources_common_database_database_pb.OrderBy;

  getSearch(): string;
  setSearch(value: string): FindDocumentsRequest;

  getCategory(): string;
  setCategory(value: string): FindDocumentsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindDocumentsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FindDocumentsRequest): FindDocumentsRequest.AsObject;
  static serializeBinaryToWriter(message: FindDocumentsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindDocumentsRequest;
  static deserializeBinaryFromReader(message: FindDocumentsRequest, reader: jspb.BinaryReader): FindDocumentsRequest;
}

export namespace FindDocumentsRequest {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationRequest.AsObject,
    orderbyList: Array<resources_common_database_database_pb.OrderBy.AsObject>,
    search: string,
    category: string,
  }
}

export class FindDocumentsResponse extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationResponse | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationResponse): FindDocumentsResponse;
  hasPagination(): boolean;
  clearPagination(): FindDocumentsResponse;

  getDocumentsList(): Array<resources_documents_documents_pb.Document>;
  setDocumentsList(value: Array<resources_documents_documents_pb.Document>): FindDocumentsResponse;
  clearDocumentsList(): FindDocumentsResponse;
  addDocuments(value?: resources_documents_documents_pb.Document, index?: number): resources_documents_documents_pb.Document;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindDocumentsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: FindDocumentsResponse): FindDocumentsResponse.AsObject;
  static serializeBinaryToWriter(message: FindDocumentsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindDocumentsResponse;
  static deserializeBinaryFromReader(message: FindDocumentsResponse, reader: jspb.BinaryReader): FindDocumentsResponse;
}

export namespace FindDocumentsResponse {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationResponse.AsObject,
    documentsList: Array<resources_documents_documents_pb.Document.AsObject>,
  }
}

export class GetDocumentRequest extends jspb.Message {
  getDocumentId(): number;
  setDocumentId(value: number): GetDocumentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentRequest): GetDocumentRequest.AsObject;
  static serializeBinaryToWriter(message: GetDocumentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentRequest;
  static deserializeBinaryFromReader(message: GetDocumentRequest, reader: jspb.BinaryReader): GetDocumentRequest;
}

export namespace GetDocumentRequest {
  export type AsObject = {
    documentId: number,
  }
}

export class GetDocumentResponse extends jspb.Message {
  getDocument(): resources_documents_documents_pb.Document | undefined;
  setDocument(value?: resources_documents_documents_pb.Document): GetDocumentResponse;
  hasDocument(): boolean;
  clearDocument(): GetDocumentResponse;

  getAccess(): resources_documents_documents_pb.DocumentAccess | undefined;
  setAccess(value?: resources_documents_documents_pb.DocumentAccess): GetDocumentResponse;
  hasAccess(): boolean;
  clearAccess(): GetDocumentResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentResponse): GetDocumentResponse.AsObject;
  static serializeBinaryToWriter(message: GetDocumentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentResponse;
  static deserializeBinaryFromReader(message: GetDocumentResponse, reader: jspb.BinaryReader): GetDocumentResponse;
}

export namespace GetDocumentResponse {
  export type AsObject = {
    document?: resources_documents_documents_pb.Document.AsObject,
    access?: resources_documents_documents_pb.DocumentAccess.AsObject,
  }
}

export class GetDocumentReferencesRequest extends jspb.Message {
  getDocumentId(): number;
  setDocumentId(value: number): GetDocumentReferencesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentReferencesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentReferencesRequest): GetDocumentReferencesRequest.AsObject;
  static serializeBinaryToWriter(message: GetDocumentReferencesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentReferencesRequest;
  static deserializeBinaryFromReader(message: GetDocumentReferencesRequest, reader: jspb.BinaryReader): GetDocumentReferencesRequest;
}

export namespace GetDocumentReferencesRequest {
  export type AsObject = {
    documentId: number,
  }
}

export class GetDocumentReferencesResponse extends jspb.Message {
  getReferencesList(): Array<resources_documents_documents_pb.DocumentReference>;
  setReferencesList(value: Array<resources_documents_documents_pb.DocumentReference>): GetDocumentReferencesResponse;
  clearReferencesList(): GetDocumentReferencesResponse;
  addReferences(value?: resources_documents_documents_pb.DocumentReference, index?: number): resources_documents_documents_pb.DocumentReference;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentReferencesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentReferencesResponse): GetDocumentReferencesResponse.AsObject;
  static serializeBinaryToWriter(message: GetDocumentReferencesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentReferencesResponse;
  static deserializeBinaryFromReader(message: GetDocumentReferencesResponse, reader: jspb.BinaryReader): GetDocumentReferencesResponse;
}

export namespace GetDocumentReferencesResponse {
  export type AsObject = {
    referencesList: Array<resources_documents_documents_pb.DocumentReference.AsObject>,
  }
}

export class GetDocumentRelationsRequest extends jspb.Message {
  getDocumentId(): number;
  setDocumentId(value: number): GetDocumentRelationsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentRelationsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentRelationsRequest): GetDocumentRelationsRequest.AsObject;
  static serializeBinaryToWriter(message: GetDocumentRelationsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentRelationsRequest;
  static deserializeBinaryFromReader(message: GetDocumentRelationsRequest, reader: jspb.BinaryReader): GetDocumentRelationsRequest;
}

export namespace GetDocumentRelationsRequest {
  export type AsObject = {
    documentId: number,
  }
}

export class GetDocumentRelationsResponse extends jspb.Message {
  getRelationsList(): Array<resources_documents_documents_pb.DocumentRelation>;
  setRelationsList(value: Array<resources_documents_documents_pb.DocumentRelation>): GetDocumentRelationsResponse;
  clearRelationsList(): GetDocumentRelationsResponse;
  addRelations(value?: resources_documents_documents_pb.DocumentRelation, index?: number): resources_documents_documents_pb.DocumentRelation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentRelationsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentRelationsResponse): GetDocumentRelationsResponse.AsObject;
  static serializeBinaryToWriter(message: GetDocumentRelationsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentRelationsResponse;
  static deserializeBinaryFromReader(message: GetDocumentRelationsResponse, reader: jspb.BinaryReader): GetDocumentRelationsResponse;
}

export namespace GetDocumentRelationsResponse {
  export type AsObject = {
    relationsList: Array<resources_documents_documents_pb.DocumentRelation.AsObject>,
  }
}

export class AddDocumentReferenceRequest extends jspb.Message {
  getReference(): resources_documents_documents_pb.DocumentReference | undefined;
  setReference(value?: resources_documents_documents_pb.DocumentReference): AddDocumentReferenceRequest;
  hasReference(): boolean;
  clearReference(): AddDocumentReferenceRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddDocumentReferenceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddDocumentReferenceRequest): AddDocumentReferenceRequest.AsObject;
  static serializeBinaryToWriter(message: AddDocumentReferenceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddDocumentReferenceRequest;
  static deserializeBinaryFromReader(message: AddDocumentReferenceRequest, reader: jspb.BinaryReader): AddDocumentReferenceRequest;
}

export namespace AddDocumentReferenceRequest {
  export type AsObject = {
    reference?: resources_documents_documents_pb.DocumentReference.AsObject,
  }
}

export class AddDocumentReferenceResponse extends jspb.Message {
  getId(): number;
  setId(value: number): AddDocumentReferenceResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddDocumentReferenceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddDocumentReferenceResponse): AddDocumentReferenceResponse.AsObject;
  static serializeBinaryToWriter(message: AddDocumentReferenceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddDocumentReferenceResponse;
  static deserializeBinaryFromReader(message: AddDocumentReferenceResponse, reader: jspb.BinaryReader): AddDocumentReferenceResponse;
}

export namespace AddDocumentReferenceResponse {
  export type AsObject = {
    id: number,
  }
}

export class RemoveDcoumentReferenceRequest extends jspb.Message {
  getId(): number;
  setId(value: number): RemoveDcoumentReferenceRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveDcoumentReferenceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveDcoumentReferenceRequest): RemoveDcoumentReferenceRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveDcoumentReferenceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveDcoumentReferenceRequest;
  static deserializeBinaryFromReader(message: RemoveDcoumentReferenceRequest, reader: jspb.BinaryReader): RemoveDcoumentReferenceRequest;
}

export namespace RemoveDcoumentReferenceRequest {
  export type AsObject = {
    id: number,
  }
}

export class RemoveDcoumentReferenceResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveDcoumentReferenceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveDcoumentReferenceResponse): RemoveDcoumentReferenceResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveDcoumentReferenceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveDcoumentReferenceResponse;
  static deserializeBinaryFromReader(message: RemoveDcoumentReferenceResponse, reader: jspb.BinaryReader): RemoveDcoumentReferenceResponse;
}

export namespace RemoveDcoumentReferenceResponse {
  export type AsObject = {
  }
}

export class AddDocumentRelationRequest extends jspb.Message {
  getRelation(): resources_documents_documents_pb.DocumentRelation | undefined;
  setRelation(value?: resources_documents_documents_pb.DocumentRelation): AddDocumentRelationRequest;
  hasRelation(): boolean;
  clearRelation(): AddDocumentRelationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddDocumentRelationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddDocumentRelationRequest): AddDocumentRelationRequest.AsObject;
  static serializeBinaryToWriter(message: AddDocumentRelationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddDocumentRelationRequest;
  static deserializeBinaryFromReader(message: AddDocumentRelationRequest, reader: jspb.BinaryReader): AddDocumentRelationRequest;
}

export namespace AddDocumentRelationRequest {
  export type AsObject = {
    relation?: resources_documents_documents_pb.DocumentRelation.AsObject,
  }
}

export class AddDocumentRelationResponse extends jspb.Message {
  getId(): number;
  setId(value: number): AddDocumentRelationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddDocumentRelationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddDocumentRelationResponse): AddDocumentRelationResponse.AsObject;
  static serializeBinaryToWriter(message: AddDocumentRelationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddDocumentRelationResponse;
  static deserializeBinaryFromReader(message: AddDocumentRelationResponse, reader: jspb.BinaryReader): AddDocumentRelationResponse;
}

export namespace AddDocumentRelationResponse {
  export type AsObject = {
    id: number,
  }
}

export class RemoveDcoumentRelationRequest extends jspb.Message {
  getId(): number;
  setId(value: number): RemoveDcoumentRelationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveDcoumentRelationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveDcoumentRelationRequest): RemoveDcoumentRelationRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveDcoumentRelationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveDcoumentRelationRequest;
  static deserializeBinaryFromReader(message: RemoveDcoumentRelationRequest, reader: jspb.BinaryReader): RemoveDcoumentRelationRequest;
}

export namespace RemoveDcoumentRelationRequest {
  export type AsObject = {
    id: number,
  }
}

export class RemoveDcoumentRelationResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveDcoumentRelationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveDcoumentRelationResponse): RemoveDcoumentRelationResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveDcoumentRelationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveDcoumentRelationResponse;
  static deserializeBinaryFromReader(message: RemoveDcoumentRelationResponse, reader: jspb.BinaryReader): RemoveDcoumentRelationResponse;
}

export namespace RemoveDcoumentRelationResponse {
  export type AsObject = {
  }
}

export class GetDocumentCommentsRequest extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationResponse | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationResponse): GetDocumentCommentsRequest;
  hasPagination(): boolean;
  clearPagination(): GetDocumentCommentsRequest;

  getDocumentId(): number;
  setDocumentId(value: number): GetDocumentCommentsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentCommentsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentCommentsRequest): GetDocumentCommentsRequest.AsObject;
  static serializeBinaryToWriter(message: GetDocumentCommentsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentCommentsRequest;
  static deserializeBinaryFromReader(message: GetDocumentCommentsRequest, reader: jspb.BinaryReader): GetDocumentCommentsRequest;
}

export namespace GetDocumentCommentsRequest {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationResponse.AsObject,
    documentId: number,
  }
}

export class GetDocumentCommentsResponse extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationResponse | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationResponse): GetDocumentCommentsResponse;
  hasPagination(): boolean;
  clearPagination(): GetDocumentCommentsResponse;

  getCommentsList(): Array<resources_documents_documents_pb.DocumentComment>;
  setCommentsList(value: Array<resources_documents_documents_pb.DocumentComment>): GetDocumentCommentsResponse;
  clearCommentsList(): GetDocumentCommentsResponse;
  addComments(value?: resources_documents_documents_pb.DocumentComment, index?: number): resources_documents_documents_pb.DocumentComment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentCommentsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentCommentsResponse): GetDocumentCommentsResponse.AsObject;
  static serializeBinaryToWriter(message: GetDocumentCommentsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentCommentsResponse;
  static deserializeBinaryFromReader(message: GetDocumentCommentsResponse, reader: jspb.BinaryReader): GetDocumentCommentsResponse;
}

export namespace GetDocumentCommentsResponse {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationResponse.AsObject,
    commentsList: Array<resources_documents_documents_pb.DocumentComment.AsObject>,
  }
}

export class PostDocumentCommentRequest extends jspb.Message {
  getComment(): resources_documents_documents_pb.DocumentComment | undefined;
  setComment(value?: resources_documents_documents_pb.DocumentComment): PostDocumentCommentRequest;
  hasComment(): boolean;
  clearComment(): PostDocumentCommentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PostDocumentCommentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PostDocumentCommentRequest): PostDocumentCommentRequest.AsObject;
  static serializeBinaryToWriter(message: PostDocumentCommentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PostDocumentCommentRequest;
  static deserializeBinaryFromReader(message: PostDocumentCommentRequest, reader: jspb.BinaryReader): PostDocumentCommentRequest;
}

export namespace PostDocumentCommentRequest {
  export type AsObject = {
    comment?: resources_documents_documents_pb.DocumentComment.AsObject,
  }
}

export class PostDocumentCommentResponse extends jspb.Message {
  getId(): number;
  setId(value: number): PostDocumentCommentResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PostDocumentCommentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PostDocumentCommentResponse): PostDocumentCommentResponse.AsObject;
  static serializeBinaryToWriter(message: PostDocumentCommentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PostDocumentCommentResponse;
  static deserializeBinaryFromReader(message: PostDocumentCommentResponse, reader: jspb.BinaryReader): PostDocumentCommentResponse;
}

export namespace PostDocumentCommentResponse {
  export type AsObject = {
    id: number,
  }
}

export class EditDocumentCommentRequest extends jspb.Message {
  getComment(): resources_documents_documents_pb.DocumentComment | undefined;
  setComment(value?: resources_documents_documents_pb.DocumentComment): EditDocumentCommentRequest;
  hasComment(): boolean;
  clearComment(): EditDocumentCommentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EditDocumentCommentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EditDocumentCommentRequest): EditDocumentCommentRequest.AsObject;
  static serializeBinaryToWriter(message: EditDocumentCommentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EditDocumentCommentRequest;
  static deserializeBinaryFromReader(message: EditDocumentCommentRequest, reader: jspb.BinaryReader): EditDocumentCommentRequest;
}

export namespace EditDocumentCommentRequest {
  export type AsObject = {
    comment?: resources_documents_documents_pb.DocumentComment.AsObject,
  }
}

export class EditDocumentCommentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EditDocumentCommentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EditDocumentCommentResponse): EditDocumentCommentResponse.AsObject;
  static serializeBinaryToWriter(message: EditDocumentCommentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EditDocumentCommentResponse;
  static deserializeBinaryFromReader(message: EditDocumentCommentResponse, reader: jspb.BinaryReader): EditDocumentCommentResponse;
}

export namespace EditDocumentCommentResponse {
  export type AsObject = {
  }
}

export class CreateDocumentRequest extends jspb.Message {
  getCategoryId(): number;
  setCategoryId(value: number): CreateDocumentRequest;
  hasCategoryId(): boolean;
  clearCategoryId(): CreateDocumentRequest;

  getTitle(): string;
  setTitle(value: string): CreateDocumentRequest;

  getContent(): string;
  setContent(value: string): CreateDocumentRequest;

  getContentType(): resources_documents_documents_pb.DOC_CONTENT_TYPE;
  setContentType(value: resources_documents_documents_pb.DOC_CONTENT_TYPE): CreateDocumentRequest;

  getData(): string;
  setData(value: string): CreateDocumentRequest;
  hasData(): boolean;
  clearData(): CreateDocumentRequest;

  getState(): string;
  setState(value: string): CreateDocumentRequest;

  getClosed(): boolean;
  setClosed(value: boolean): CreateDocumentRequest;

  getPublic(): boolean;
  setPublic(value: boolean): CreateDocumentRequest;

  getAccess(): resources_documents_documents_pb.DocumentAccess | undefined;
  setAccess(value?: resources_documents_documents_pb.DocumentAccess): CreateDocumentRequest;
  hasAccess(): boolean;
  clearAccess(): CreateDocumentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateDocumentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateDocumentRequest): CreateDocumentRequest.AsObject;
  static serializeBinaryToWriter(message: CreateDocumentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateDocumentRequest;
  static deserializeBinaryFromReader(message: CreateDocumentRequest, reader: jspb.BinaryReader): CreateDocumentRequest;
}

export namespace CreateDocumentRequest {
  export type AsObject = {
    categoryId?: number,
    title: string,
    content: string,
    contentType: resources_documents_documents_pb.DOC_CONTENT_TYPE,
    data?: string,
    state: string,
    closed: boolean,
    pb_public: boolean,
    access?: resources_documents_documents_pb.DocumentAccess.AsObject,
  }

  export enum CategoryIdCase { 
    _CATEGORY_ID_NOT_SET = 0,
    CATEGORY_ID = 1,
  }

  export enum DataCase { 
    _DATA_NOT_SET = 0,
    DATA = 5,
  }

  export enum AccessCase { 
    _ACCESS_NOT_SET = 0,
    ACCESS = 9,
  }
}

export class CreateDocumentResponse extends jspb.Message {
  getDocumentId(): number;
  setDocumentId(value: number): CreateDocumentResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateDocumentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateDocumentResponse): CreateDocumentResponse.AsObject;
  static serializeBinaryToWriter(message: CreateDocumentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateDocumentResponse;
  static deserializeBinaryFromReader(message: CreateDocumentResponse, reader: jspb.BinaryReader): CreateDocumentResponse;
}

export namespace CreateDocumentResponse {
  export type AsObject = {
    documentId: number,
  }
}

export class UpdateDocumentRequest extends jspb.Message {
  getDocumentId(): number;
  setDocumentId(value: number): UpdateDocumentRequest;

  getCategoryId(): number;
  setCategoryId(value: number): UpdateDocumentRequest;
  hasCategoryId(): boolean;
  clearCategoryId(): UpdateDocumentRequest;

  getTitle(): string;
  setTitle(value: string): UpdateDocumentRequest;
  hasTitle(): boolean;
  clearTitle(): UpdateDocumentRequest;

  getContent(): string;
  setContent(value: string): UpdateDocumentRequest;
  hasContent(): boolean;
  clearContent(): UpdateDocumentRequest;

  getContentType(): resources_documents_documents_pb.DOC_CONTENT_TYPE;
  setContentType(value: resources_documents_documents_pb.DOC_CONTENT_TYPE): UpdateDocumentRequest;
  hasContentType(): boolean;
  clearContentType(): UpdateDocumentRequest;

  getData(): string;
  setData(value: string): UpdateDocumentRequest;
  hasData(): boolean;
  clearData(): UpdateDocumentRequest;

  getState(): string;
  setState(value: string): UpdateDocumentRequest;
  hasState(): boolean;
  clearState(): UpdateDocumentRequest;

  getClosed(): boolean;
  setClosed(value: boolean): UpdateDocumentRequest;
  hasClosed(): boolean;
  clearClosed(): UpdateDocumentRequest;

  getPublic(): boolean;
  setPublic(value: boolean): UpdateDocumentRequest;
  hasPublic(): boolean;
  clearPublic(): UpdateDocumentRequest;

  getAccess(): resources_documents_documents_pb.DocumentAccess | undefined;
  setAccess(value?: resources_documents_documents_pb.DocumentAccess): UpdateDocumentRequest;
  hasAccess(): boolean;
  clearAccess(): UpdateDocumentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateDocumentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateDocumentRequest): UpdateDocumentRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateDocumentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateDocumentRequest;
  static deserializeBinaryFromReader(message: UpdateDocumentRequest, reader: jspb.BinaryReader): UpdateDocumentRequest;
}

export namespace UpdateDocumentRequest {
  export type AsObject = {
    documentId: number,
    categoryId?: number,
    title?: string,
    content?: string,
    contentType?: resources_documents_documents_pb.DOC_CONTENT_TYPE,
    data?: string,
    state?: string,
    closed?: boolean,
    pb_public?: boolean,
    access?: resources_documents_documents_pb.DocumentAccess.AsObject,
  }

  export enum CategoryIdCase { 
    _CATEGORY_ID_NOT_SET = 0,
    CATEGORY_ID = 2,
  }

  export enum TitleCase { 
    _TITLE_NOT_SET = 0,
    TITLE = 3,
  }

  export enum ContentCase { 
    _CONTENT_NOT_SET = 0,
    CONTENT = 4,
  }

  export enum ContentTypeCase { 
    _CONTENT_TYPE_NOT_SET = 0,
    CONTENT_TYPE = 5,
  }

  export enum DataCase { 
    _DATA_NOT_SET = 0,
    DATA = 6,
  }

  export enum StateCase { 
    _STATE_NOT_SET = 0,
    STATE = 7,
  }

  export enum ClosedCase { 
    _CLOSED_NOT_SET = 0,
    CLOSED = 8,
  }

  export enum PublicCase { 
    _PUBLIC_NOT_SET = 0,
    PUBLIC = 9,
  }

  export enum AccessCase { 
    _ACCESS_NOT_SET = 0,
    ACCESS = 10,
  }
}

export class UpdateDocumentResponse extends jspb.Message {
  getDocumentId(): number;
  setDocumentId(value: number): UpdateDocumentResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateDocumentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateDocumentResponse): UpdateDocumentResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateDocumentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateDocumentResponse;
  static deserializeBinaryFromReader(message: UpdateDocumentResponse, reader: jspb.BinaryReader): UpdateDocumentResponse;
}

export namespace UpdateDocumentResponse {
  export type AsObject = {
    documentId: number,
  }
}

export class GetDocumentAccessRequest extends jspb.Message {
  getDocumentId(): number;
  setDocumentId(value: number): GetDocumentAccessRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentAccessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentAccessRequest): GetDocumentAccessRequest.AsObject;
  static serializeBinaryToWriter(message: GetDocumentAccessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentAccessRequest;
  static deserializeBinaryFromReader(message: GetDocumentAccessRequest, reader: jspb.BinaryReader): GetDocumentAccessRequest;
}

export namespace GetDocumentAccessRequest {
  export type AsObject = {
    documentId: number,
  }
}

export class GetDocumentAccessResponse extends jspb.Message {
  getAccess(): resources_documents_documents_pb.DocumentAccess | undefined;
  setAccess(value?: resources_documents_documents_pb.DocumentAccess): GetDocumentAccessResponse;
  hasAccess(): boolean;
  clearAccess(): GetDocumentAccessResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentAccessResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentAccessResponse): GetDocumentAccessResponse.AsObject;
  static serializeBinaryToWriter(message: GetDocumentAccessResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentAccessResponse;
  static deserializeBinaryFromReader(message: GetDocumentAccessResponse, reader: jspb.BinaryReader): GetDocumentAccessResponse;
}

export namespace GetDocumentAccessResponse {
  export type AsObject = {
    access?: resources_documents_documents_pb.DocumentAccess.AsObject,
  }
}

export class SetDocumentAccessRequest extends jspb.Message {
  getDocumentId(): number;
  setDocumentId(value: number): SetDocumentAccessRequest;

  getMode(): DOC_ACCESS_UPDATE_MODE;
  setMode(value: DOC_ACCESS_UPDATE_MODE): SetDocumentAccessRequest;

  getAccess(): resources_documents_documents_pb.DocumentAccess | undefined;
  setAccess(value?: resources_documents_documents_pb.DocumentAccess): SetDocumentAccessRequest;
  hasAccess(): boolean;
  clearAccess(): SetDocumentAccessRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDocumentAccessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDocumentAccessRequest): SetDocumentAccessRequest.AsObject;
  static serializeBinaryToWriter(message: SetDocumentAccessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDocumentAccessRequest;
  static deserializeBinaryFromReader(message: SetDocumentAccessRequest, reader: jspb.BinaryReader): SetDocumentAccessRequest;
}

export namespace SetDocumentAccessRequest {
  export type AsObject = {
    documentId: number,
    mode: DOC_ACCESS_UPDATE_MODE,
    access?: resources_documents_documents_pb.DocumentAccess.AsObject,
  }
}

export class SetDocumentAccessResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDocumentAccessResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetDocumentAccessResponse): SetDocumentAccessResponse.AsObject;
  static serializeBinaryToWriter(message: SetDocumentAccessResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDocumentAccessResponse;
  static deserializeBinaryFromReader(message: SetDocumentAccessResponse, reader: jspb.BinaryReader): SetDocumentAccessResponse;
}

export namespace SetDocumentAccessResponse {
  export type AsObject = {
  }
}

export class GetUserDocumentsRequest extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationRequest | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationRequest): GetUserDocumentsRequest;
  hasPagination(): boolean;
  clearPagination(): GetUserDocumentsRequest;

  getUserId(): number;
  setUserId(value: number): GetUserDocumentsRequest;

  getRelationsList(): Array<resources_documents_documents_pb.DOC_RELATION_TYPE>;
  setRelationsList(value: Array<resources_documents_documents_pb.DOC_RELATION_TYPE>): GetUserDocumentsRequest;
  clearRelationsList(): GetUserDocumentsRequest;
  addRelations(value: resources_documents_documents_pb.DOC_RELATION_TYPE, index?: number): GetUserDocumentsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserDocumentsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserDocumentsRequest): GetUserDocumentsRequest.AsObject;
  static serializeBinaryToWriter(message: GetUserDocumentsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserDocumentsRequest;
  static deserializeBinaryFromReader(message: GetUserDocumentsRequest, reader: jspb.BinaryReader): GetUserDocumentsRequest;
}

export namespace GetUserDocumentsRequest {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationRequest.AsObject,
    userId: number,
    relationsList: Array<resources_documents_documents_pb.DOC_RELATION_TYPE>,
  }
}

export class GetUserDocumentsResponse extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationResponse | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationResponse): GetUserDocumentsResponse;
  hasPagination(): boolean;
  clearPagination(): GetUserDocumentsResponse;

  getRelationsList(): Array<resources_documents_documents_pb.DocumentRelation>;
  setRelationsList(value: Array<resources_documents_documents_pb.DocumentRelation>): GetUserDocumentsResponse;
  clearRelationsList(): GetUserDocumentsResponse;
  addRelations(value?: resources_documents_documents_pb.DocumentRelation, index?: number): resources_documents_documents_pb.DocumentRelation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserDocumentsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserDocumentsResponse): GetUserDocumentsResponse.AsObject;
  static serializeBinaryToWriter(message: GetUserDocumentsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserDocumentsResponse;
  static deserializeBinaryFromReader(message: GetUserDocumentsResponse, reader: jspb.BinaryReader): GetUserDocumentsResponse;
}

export namespace GetUserDocumentsResponse {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationResponse.AsObject,
    relationsList: Array<resources_documents_documents_pb.DocumentRelation.AsObject>,
  }
}

export enum DOC_ACCESS_UPDATE_MODE { 
  UPDATE = 0,
  DELETE = 1,
  CLEAR = 2,
}
