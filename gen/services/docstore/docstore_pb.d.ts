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
  }
}

export class GetTemplateResponse extends jspb.Message {
  getTemplate(): resources_documents_documents_pb.DocumentTemplate | undefined;
  setTemplate(value?: resources_documents_documents_pb.DocumentTemplate): GetTemplateResponse;
  hasTemplate(): boolean;
  clearTemplate(): GetTemplateResponse;

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
  }
}

export class FindDocumentsRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): FindDocumentsRequest;

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
    offset: number,
    orderbyList: Array<resources_common_database_database_pb.OrderBy.AsObject>,
    search: string,
    category: string,
  }
}

export class FindDocumentsResponse extends jspb.Message {
  getTotalCount(): number;
  setTotalCount(value: number): FindDocumentsResponse;

  getOffset(): number;
  setOffset(value: number): FindDocumentsResponse;

  getEnd(): number;
  setEnd(value: number): FindDocumentsResponse;

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
    totalCount: number,
    offset: number,
    end: number,
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

export class GetDocumentCommentsRequest extends jspb.Message {
  getDocumentid(): number;
  setDocumentid(value: number): GetDocumentCommentsRequest;

  getOffset(): number;
  setOffset(value: number): GetDocumentCommentsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentCommentsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentCommentsRequest): GetDocumentCommentsRequest.AsObject;
  static serializeBinaryToWriter(message: GetDocumentCommentsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentCommentsRequest;
  static deserializeBinaryFromReader(message: GetDocumentCommentsRequest, reader: jspb.BinaryReader): GetDocumentCommentsRequest;
}

export namespace GetDocumentCommentsRequest {
  export type AsObject = {
    documentid: number,
    offset: number,
  }
}

export class GetDocumentCommentsResponse extends jspb.Message {
  getTotalCount(): number;
  setTotalCount(value: number): GetDocumentCommentsResponse;

  getOffset(): number;
  setOffset(value: number): GetDocumentCommentsResponse;

  getEnd(): number;
  setEnd(value: number): GetDocumentCommentsResponse;

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
    totalCount: number,
    offset: number,
    end: number,
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
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PostDocumentCommentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PostDocumentCommentResponse): PostDocumentCommentResponse.AsObject;
  static serializeBinaryToWriter(message: PostDocumentCommentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PostDocumentCommentResponse;
  static deserializeBinaryFromReader(message: PostDocumentCommentResponse, reader: jspb.BinaryReader): PostDocumentCommentResponse;
}

export namespace PostDocumentCommentResponse {
  export type AsObject = {
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
  getTitle(): string;
  setTitle(value: string): CreateDocumentRequest;

  getContent(): string;
  setContent(value: string): CreateDocumentRequest;

  getContentType(): resources_documents_documents_pb.DOC_CONTENT_TYPE;
  setContentType(value: resources_documents_documents_pb.DOC_CONTENT_TYPE): CreateDocumentRequest;

  getClosed(): boolean;
  setClosed(value: boolean): CreateDocumentRequest;

  getState(): string;
  setState(value: string): CreateDocumentRequest;

  getPublic(): boolean;
  setPublic(value: boolean): CreateDocumentRequest;

  getCategoryId(): number;
  setCategoryId(value: number): CreateDocumentRequest;

  getTargetDocumentId(): number;
  setTargetDocumentId(value: number): CreateDocumentRequest;

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
    title: string,
    content: string,
    contentType: resources_documents_documents_pb.DOC_CONTENT_TYPE,
    closed: boolean,
    state: string,
    pb_public: boolean,
    categoryId: number,
    targetDocumentId: number,
    access?: resources_documents_documents_pb.DocumentAccess.AsObject,
  }
}

export class CreateDocumentResponse extends jspb.Message {
  getId(): number;
  setId(value: number): CreateDocumentResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateDocumentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateDocumentResponse): CreateDocumentResponse.AsObject;
  static serializeBinaryToWriter(message: CreateDocumentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateDocumentResponse;
  static deserializeBinaryFromReader(message: CreateDocumentResponse, reader: jspb.BinaryReader): CreateDocumentResponse;
}

export namespace CreateDocumentResponse {
  export type AsObject = {
    id: number,
  }
}

export class UpdateDocumentRequest extends jspb.Message {
  getDocumentId(): number;
  setDocumentId(value: number): UpdateDocumentRequest;

  getTitle(): string;
  setTitle(value: string): UpdateDocumentRequest;

  getContent(): string;
  setContent(value: string): UpdateDocumentRequest;

  getContentType(): resources_documents_documents_pb.DOC_CONTENT_TYPE;
  setContentType(value: resources_documents_documents_pb.DOC_CONTENT_TYPE): UpdateDocumentRequest;

  getCategoryId(): number;
  setCategoryId(value: number): UpdateDocumentRequest;

  getClosed(): boolean;
  setClosed(value: boolean): UpdateDocumentRequest;

  getState(): string;
  setState(value: string): UpdateDocumentRequest;

  getPublic(): boolean;
  setPublic(value: boolean): UpdateDocumentRequest;

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
    title: string,
    content: string,
    contentType: resources_documents_documents_pb.DOC_CONTENT_TYPE,
    categoryId: number,
    closed: boolean,
    state: string,
    pb_public: boolean,
  }
}

export class UpdateDocumentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateDocumentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateDocumentResponse): UpdateDocumentResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateDocumentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateDocumentResponse;
  static deserializeBinaryFromReader(message: UpdateDocumentResponse, reader: jspb.BinaryReader): UpdateDocumentResponse;
}

export namespace UpdateDocumentResponse {
  export type AsObject = {
  }
}

export class GetDocumentFeedRequest extends jspb.Message {
  getDocumentId(): number;
  setDocumentId(value: number): GetDocumentFeedRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentFeedRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentFeedRequest): GetDocumentFeedRequest.AsObject;
  static serializeBinaryToWriter(message: GetDocumentFeedRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentFeedRequest;
  static deserializeBinaryFromReader(message: GetDocumentFeedRequest, reader: jspb.BinaryReader): GetDocumentFeedRequest;
}

export namespace GetDocumentFeedRequest {
  export type AsObject = {
    documentId: number,
  }
}

export class GetDocumentFeedResponse extends jspb.Message {
  getItemsList(): Array<resources_documents_documents_pb.DocumentFeed>;
  setItemsList(value: Array<resources_documents_documents_pb.DocumentFeed>): GetDocumentFeedResponse;
  clearItemsList(): GetDocumentFeedResponse;
  addItems(value?: resources_documents_documents_pb.DocumentFeed, index?: number): resources_documents_documents_pb.DocumentFeed;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentFeedResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentFeedResponse): GetDocumentFeedResponse.AsObject;
  static serializeBinaryToWriter(message: GetDocumentFeedResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentFeedResponse;
  static deserializeBinaryFromReader(message: GetDocumentFeedResponse, reader: jspb.BinaryReader): GetDocumentFeedResponse;
}

export namespace GetDocumentFeedResponse {
  export type AsObject = {
    itemsList: Array<resources_documents_documents_pb.DocumentFeed.AsObject>,
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

export enum DOC_ACCESS_UPDATE_MODE { 
  ADD = 0,
  REPLACE = 1,
  DELETE = 2,
  CLEAR = 3,
}
