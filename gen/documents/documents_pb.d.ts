import * as jspb from 'google-protobuf'

import * as common_database_database_pb from '../common/database/database_pb';
import * as common_timestamp_timestamp_pb from '../common/timestamp/timestamp_pb';
import * as common_userinfo_userinfo_pb from '../common/userinfo/userinfo_pb';


export class FindDocumentsRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): FindDocumentsRequest;

  getOrderbyList(): Array<common_database_database_pb.OrderBy>;
  setOrderbyList(value: Array<common_database_database_pb.OrderBy>): FindDocumentsRequest;
  clearOrderbyList(): FindDocumentsRequest;
  addOrderby(value?: common_database_database_pb.OrderBy, index?: number): common_database_database_pb.OrderBy;

  getSearch(): string;
  setSearch(value: string): FindDocumentsRequest;

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
    orderbyList: Array<common_database_database_pb.OrderBy.AsObject>,
    search: string,
  }
}

export class FindDocumentsResponse extends jspb.Message {
  getDocumentsList(): Array<Document>;
  setDocumentsList(value: Array<Document>): FindDocumentsResponse;
  clearDocumentsList(): FindDocumentsResponse;
  addDocuments(value?: Document, index?: number): Document;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindDocumentsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: FindDocumentsResponse): FindDocumentsResponse.AsObject;
  static serializeBinaryToWriter(message: FindDocumentsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindDocumentsResponse;
  static deserializeBinaryFromReader(message: FindDocumentsResponse, reader: jspb.BinaryReader): FindDocumentsResponse;
}

export namespace FindDocumentsResponse {
  export type AsObject = {
    documentsList: Array<Document.AsObject>,
  }
}

export class Document extends jspb.Message {
  getId(): number;
  setId(value: number): Document;

  getCreatedAt(): common_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: common_timestamp_timestamp_pb.Timestamp): Document;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Document;

  getUpdatedAt(): common_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: common_timestamp_timestamp_pb.Timestamp): Document;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): Document;

  getTitle(): string;
  setTitle(value: string): Document;

  getContent(): string;
  setContent(value: string): Document;

  getContentType(): string;
  setContentType(value: string): Document;

  getClosed(): boolean;
  setClosed(value: boolean): Document;

  getState(): string;
  setState(value: string): Document;

  getCreator(): common_userinfo_userinfo_pb.ShortUser | undefined;
  setCreator(value?: common_userinfo_userinfo_pb.ShortUser): Document;
  hasCreator(): boolean;
  clearCreator(): Document;

  getCreatorJob(): string;
  setCreatorJob(value: string): Document;

  getPublic(): boolean;
  setPublic(value: boolean): Document;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Document.AsObject;
  static toObject(includeInstance: boolean, msg: Document): Document.AsObject;
  static serializeBinaryToWriter(message: Document, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Document;
  static deserializeBinaryFromReader(message: Document, reader: jspb.BinaryReader): Document;
}

export namespace Document {
  export type AsObject = {
    id: number,
    createdAt?: common_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: common_timestamp_timestamp_pb.Timestamp.AsObject,
    title: string,
    content: string,
    contentType: string,
    closed: boolean,
    state: string,
    creator?: common_userinfo_userinfo_pb.ShortUser.AsObject,
    creatorJob: string,
    pb_public: boolean,
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
  getDocument(): Document | undefined;
  setDocument(value?: Document): GetDocumentResponse;
  hasDocument(): boolean;
  clearDocument(): GetDocumentResponse;

  getResponsesList(): Array<Document>;
  setResponsesList(value: Array<Document>): GetDocumentResponse;
  clearResponsesList(): GetDocumentResponse;
  addResponses(value?: Document, index?: number): Document;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentResponse): GetDocumentResponse.AsObject;
  static serializeBinaryToWriter(message: GetDocumentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentResponse;
  static deserializeBinaryFromReader(message: GetDocumentResponse, reader: jspb.BinaryReader): GetDocumentResponse;
}

export namespace GetDocumentResponse {
  export type AsObject = {
    document?: Document.AsObject,
    responsesList: Array<Document.AsObject>,
  }
}

export class CreateOrEditDocumentRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateOrEditDocumentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateOrEditDocumentRequest): CreateOrEditDocumentRequest.AsObject;
  static serializeBinaryToWriter(message: CreateOrEditDocumentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateOrEditDocumentRequest;
  static deserializeBinaryFromReader(message: CreateOrEditDocumentRequest, reader: jspb.BinaryReader): CreateOrEditDocumentRequest;
}

export namespace CreateOrEditDocumentRequest {
  export type AsObject = {
  }
}

export class CreateOrEditDocumentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateOrEditDocumentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateOrEditDocumentResponse): CreateOrEditDocumentResponse.AsObject;
  static serializeBinaryToWriter(message: CreateOrEditDocumentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateOrEditDocumentResponse;
  static deserializeBinaryFromReader(message: CreateOrEditDocumentResponse, reader: jspb.BinaryReader): CreateOrEditDocumentResponse;
}

export namespace CreateOrEditDocumentResponse {
  export type AsObject = {
  }
}

export class Template extends jspb.Message {
  getId(): number;
  setId(value: number): Template;

  getJob(): string;
  setJob(value: string): Template;

  getJobgrade(): number;
  setJobgrade(value: number): Template;

  getTitle(): string;
  setTitle(value: string): Template;

  getDescription(): string;
  setDescription(value: string): Template;

  getContentTitle(): string;
  setContentTitle(value: string): Template;

  getContent(): string;
  setContent(value: string): Template;

  getAdditionaldata(): string;
  setAdditionaldata(value: string): Template;

  getCreatorid(): number;
  setCreatorid(value: number): Template;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Template.AsObject;
  static toObject(includeInstance: boolean, msg: Template): Template.AsObject;
  static serializeBinaryToWriter(message: Template, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Template;
  static deserializeBinaryFromReader(message: Template, reader: jspb.BinaryReader): Template;
}

export namespace Template {
  export type AsObject = {
    id: number,
    job: string,
    jobgrade: number,
    title: string,
    description: string,
    contentTitle: string,
    content: string,
    additionaldata: string,
    creatorid: number,
  }
}

export class TemplateShort extends jspb.Message {
  getId(): number;
  setId(value: number): TemplateShort;

  getJob(): string;
  setJob(value: string): TemplateShort;

  getJobgrade(): number;
  setJobgrade(value: number): TemplateShort;

  getTitle(): string;
  setTitle(value: string): TemplateShort;

  getDescription(): string;
  setDescription(value: string): TemplateShort;

  getCreatorid(): number;
  setCreatorid(value: number): TemplateShort;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TemplateShort.AsObject;
  static toObject(includeInstance: boolean, msg: TemplateShort): TemplateShort.AsObject;
  static serializeBinaryToWriter(message: TemplateShort, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TemplateShort;
  static deserializeBinaryFromReader(message: TemplateShort, reader: jspb.BinaryReader): TemplateShort;
}

export namespace TemplateShort {
  export type AsObject = {
    id: number,
    job: string,
    jobgrade: number,
    title: string,
    description: string,
    creatorid: number,
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
  getTemplatesList(): Array<TemplateShort>;
  setTemplatesList(value: Array<TemplateShort>): ListTemplatesResponse;
  clearTemplatesList(): ListTemplatesResponse;
  addTemplates(value?: TemplateShort, index?: number): TemplateShort;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTemplatesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListTemplatesResponse): ListTemplatesResponse.AsObject;
  static serializeBinaryToWriter(message: ListTemplatesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTemplatesResponse;
  static deserializeBinaryFromReader(message: ListTemplatesResponse, reader: jspb.BinaryReader): ListTemplatesResponse;
}

export namespace ListTemplatesResponse {
  export type AsObject = {
    templatesList: Array<TemplateShort.AsObject>,
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
  getTemplate(): Template | undefined;
  setTemplate(value?: Template): GetTemplateResponse;
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
    template?: Template.AsObject,
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
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDocumentAccessResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDocumentAccessResponse): GetDocumentAccessResponse.AsObject;
  static serializeBinaryToWriter(message: GetDocumentAccessResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDocumentAccessResponse;
  static deserializeBinaryFromReader(message: GetDocumentAccessResponse, reader: jspb.BinaryReader): GetDocumentAccessResponse;
}

export namespace GetDocumentAccessResponse {
  export type AsObject = {
  }
}

export class SetDocumentAccessRequest extends jspb.Message {
  getId(): number;
  setId(value: number): SetDocumentAccessRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDocumentAccessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDocumentAccessRequest): SetDocumentAccessRequest.AsObject;
  static serializeBinaryToWriter(message: SetDocumentAccessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDocumentAccessRequest;
  static deserializeBinaryFromReader(message: SetDocumentAccessRequest, reader: jspb.BinaryReader): SetDocumentAccessRequest;
}

export namespace SetDocumentAccessRequest {
  export type AsObject = {
    id: number,
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

