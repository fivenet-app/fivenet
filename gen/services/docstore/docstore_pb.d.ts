import * as jspb from 'google-protobuf'

import * as resources_common_database_database_pb from '../../resources/common/database/database_pb';
import * as resources_documents_documents_pb from '../../resources/documents/documents_pb';


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
  getTotalcount(): number;
  setTotalcount(value: number): FindDocumentsResponse;

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
    totalcount: number,
    offset: number,
    end: number,
    documentsList: Array<resources_documents_documents_pb.Document.AsObject>,
  }
}

export class GetDocumentRequest extends jspb.Message {
  getId(): number;
  setId(value: number): GetDocumentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentRequest): GetDocumentRequest.AsObject;
  static serializeBinaryToWriter(message: GetDocumentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentRequest;
  static deserializeBinaryFromReader(message: GetDocumentRequest, reader: jspb.BinaryReader): GetDocumentRequest;
}

export namespace GetDocumentRequest {
  export type AsObject = {
    id: number,
  }
}

export class GetDocumentResponse extends jspb.Message {
  getDocument(): resources_documents_documents_pb.Document | undefined;
  setDocument(value?: resources_documents_documents_pb.Document): GetDocumentResponse;
  hasDocument(): boolean;
  clearDocument(): GetDocumentResponse;

  getResponsesList(): Array<resources_documents_documents_pb.Document>;
  setResponsesList(value: Array<resources_documents_documents_pb.Document>): GetDocumentResponse;
  clearResponsesList(): GetDocumentResponse;
  addResponses(value?: resources_documents_documents_pb.Document, index?: number): resources_documents_documents_pb.Document;

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
    responsesList: Array<resources_documents_documents_pb.Document.AsObject>,
  }
}

export class CreateDocumentRequest extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): CreateDocumentRequest;

  getContent(): string;
  setContent(value: string): CreateDocumentRequest;

  getContentType(): resources_documents_documents_pb.DOCUMENT_CONTENT_TYPE;
  setContentType(value: resources_documents_documents_pb.DOCUMENT_CONTENT_TYPE): CreateDocumentRequest;

  getClosed(): boolean;
  setClosed(value: boolean): CreateDocumentRequest;

  getState(): string;
  setState(value: string): CreateDocumentRequest;

  getPublic(): boolean;
  setPublic(value: boolean): CreateDocumentRequest;

  getCategoryid(): number;
  setCategoryid(value: number): CreateDocumentRequest;

  getTargetdocumentid(): number;
  setTargetdocumentid(value: number): CreateDocumentRequest;

  getJobsaccessList(): Array<resources_documents_documents_pb.DocumentJobAccess>;
  setJobsaccessList(value: Array<resources_documents_documents_pb.DocumentJobAccess>): CreateDocumentRequest;
  clearJobsaccessList(): CreateDocumentRequest;
  addJobsaccess(value?: resources_documents_documents_pb.DocumentJobAccess, index?: number): resources_documents_documents_pb.DocumentJobAccess;

  getUsersaccessList(): Array<resources_documents_documents_pb.DocumentUserAccess>;
  setUsersaccessList(value: Array<resources_documents_documents_pb.DocumentUserAccess>): CreateDocumentRequest;
  clearUsersaccessList(): CreateDocumentRequest;
  addUsersaccess(value?: resources_documents_documents_pb.DocumentUserAccess, index?: number): resources_documents_documents_pb.DocumentUserAccess;

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
    contentType: resources_documents_documents_pb.DOCUMENT_CONTENT_TYPE,
    closed: boolean,
    state: string,
    pb_public: boolean,
    categoryid: number,
    targetdocumentid: number,
    jobsaccessList: Array<resources_documents_documents_pb.DocumentJobAccess.AsObject>,
    usersaccessList: Array<resources_documents_documents_pb.DocumentUserAccess.AsObject>,
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
  getId(): number;
  setId(value: number): UpdateDocumentRequest;

  getTitle(): string;
  setTitle(value: string): UpdateDocumentRequest;

  getContent(): string;
  setContent(value: string): UpdateDocumentRequest;

  getContentType(): resources_documents_documents_pb.DOCUMENT_CONTENT_TYPE;
  setContentType(value: resources_documents_documents_pb.DOCUMENT_CONTENT_TYPE): UpdateDocumentRequest;

  getCategoryid(): number;
  setCategoryid(value: number): UpdateDocumentRequest;

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
    id: number,
    title: string,
    content: string,
    contentType: resources_documents_documents_pb.DOCUMENT_CONTENT_TYPE,
    categoryid: number,
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

export class GetDocumentResponsesRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): GetDocumentResponsesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentResponsesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentResponsesRequest): GetDocumentResponsesRequest.AsObject;
  static serializeBinaryToWriter(message: GetDocumentResponsesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentResponsesRequest;
  static deserializeBinaryFromReader(message: GetDocumentResponsesRequest, reader: jspb.BinaryReader): GetDocumentResponsesRequest;
}

export namespace GetDocumentResponsesRequest {
  export type AsObject = {
    offset: number,
  }
}

export class GetDocumentResponsesResponse extends jspb.Message {
  getTotalcount(): number;
  setTotalcount(value: number): GetDocumentResponsesResponse;

  getOffset(): number;
  setOffset(value: number): GetDocumentResponsesResponse;

  getEnd(): number;
  setEnd(value: number): GetDocumentResponsesResponse;

  getResponsesList(): Array<resources_documents_documents_pb.Document>;
  setResponsesList(value: Array<resources_documents_documents_pb.Document>): GetDocumentResponsesResponse;
  clearResponsesList(): GetDocumentResponsesResponse;
  addResponses(value?: resources_documents_documents_pb.Document, index?: number): resources_documents_documents_pb.Document;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentResponsesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentResponsesResponse): GetDocumentResponsesResponse.AsObject;
  static serializeBinaryToWriter(message: GetDocumentResponsesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentResponsesResponse;
  static deserializeBinaryFromReader(message: GetDocumentResponsesResponse, reader: jspb.BinaryReader): GetDocumentResponsesResponse;
}

export namespace GetDocumentResponsesResponse {
  export type AsObject = {
    totalcount: number,
    offset: number,
    end: number,
    responsesList: Array<resources_documents_documents_pb.Document.AsObject>,
  }
}

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
  getId(): number;
  setId(value: number): GetTemplateRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTemplateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetTemplateRequest): GetTemplateRequest.AsObject;
  static serializeBinaryToWriter(message: GetTemplateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTemplateRequest;
  static deserializeBinaryFromReader(message: GetTemplateRequest, reader: jspb.BinaryReader): GetTemplateRequest;
}

export namespace GetTemplateRequest {
  export type AsObject = {
    id: number,
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

export class GetDocumentAccessRequest extends jspb.Message {
  getId(): number;
  setId(value: number): GetDocumentAccessRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentAccessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentAccessRequest): GetDocumentAccessRequest.AsObject;
  static serializeBinaryToWriter(message: GetDocumentAccessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentAccessRequest;
  static deserializeBinaryFromReader(message: GetDocumentAccessRequest, reader: jspb.BinaryReader): GetDocumentAccessRequest;
}

export namespace GetDocumentAccessRequest {
  export type AsObject = {
    id: number,
  }
}

export class GetDocumentAccessResponse extends jspb.Message {
  getJobsList(): Array<resources_documents_documents_pb.DocumentJobAccess>;
  setJobsList(value: Array<resources_documents_documents_pb.DocumentJobAccess>): GetDocumentAccessResponse;
  clearJobsList(): GetDocumentAccessResponse;
  addJobs(value?: resources_documents_documents_pb.DocumentJobAccess, index?: number): resources_documents_documents_pb.DocumentJobAccess;

  getUsersList(): Array<resources_documents_documents_pb.DocumentUserAccess>;
  setUsersList(value: Array<resources_documents_documents_pb.DocumentUserAccess>): GetDocumentAccessResponse;
  clearUsersList(): GetDocumentAccessResponse;
  addUsers(value?: resources_documents_documents_pb.DocumentUserAccess, index?: number): resources_documents_documents_pb.DocumentUserAccess;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentAccessResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentAccessResponse): GetDocumentAccessResponse.AsObject;
  static serializeBinaryToWriter(message: GetDocumentAccessResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentAccessResponse;
  static deserializeBinaryFromReader(message: GetDocumentAccessResponse, reader: jspb.BinaryReader): GetDocumentAccessResponse;
}

export namespace GetDocumentAccessResponse {
  export type AsObject = {
    jobsList: Array<resources_documents_documents_pb.DocumentJobAccess.AsObject>,
    usersList: Array<resources_documents_documents_pb.DocumentUserAccess.AsObject>,
  }
}

export class SetDocumentAccessRequest extends jspb.Message {
  getDocumentid(): number;
  setDocumentid(value: number): SetDocumentAccessRequest;

  getMode(): DOCUMENT_ACCESS_UPDATE_MODE;
  setMode(value: DOCUMENT_ACCESS_UPDATE_MODE): SetDocumentAccessRequest;

  getJobsList(): Array<resources_documents_documents_pb.DocumentJobAccess>;
  setJobsList(value: Array<resources_documents_documents_pb.DocumentJobAccess>): SetDocumentAccessRequest;
  clearJobsList(): SetDocumentAccessRequest;
  addJobs(value?: resources_documents_documents_pb.DocumentJobAccess, index?: number): resources_documents_documents_pb.DocumentJobAccess;

  getUsersList(): Array<resources_documents_documents_pb.DocumentUserAccess>;
  setUsersList(value: Array<resources_documents_documents_pb.DocumentUserAccess>): SetDocumentAccessRequest;
  clearUsersList(): SetDocumentAccessRequest;
  addUsers(value?: resources_documents_documents_pb.DocumentUserAccess, index?: number): resources_documents_documents_pb.DocumentUserAccess;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDocumentAccessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDocumentAccessRequest): SetDocumentAccessRequest.AsObject;
  static serializeBinaryToWriter(message: SetDocumentAccessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDocumentAccessRequest;
  static deserializeBinaryFromReader(message: SetDocumentAccessRequest, reader: jspb.BinaryReader): SetDocumentAccessRequest;
}

export namespace SetDocumentAccessRequest {
  export type AsObject = {
    documentid: number,
    mode: DOCUMENT_ACCESS_UPDATE_MODE,
    jobsList: Array<resources_documents_documents_pb.DocumentJobAccess.AsObject>,
    usersList: Array<resources_documents_documents_pb.DocumentUserAccess.AsObject>,
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

export enum DOCUMENT_ACCESS_UPDATE_MODE { 
  ADD = 0,
  REPLACE = 1,
  DELETE = 2,
  CLEAR = 3,
}
