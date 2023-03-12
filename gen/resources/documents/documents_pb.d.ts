import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';


export class Document extends jspb.Message {
  getId(): number;
  setId(value: number): Document;

  getCreatedat(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedat(value?: resources_timestamp_timestamp_pb.Timestamp): Document;
  hasCreatedat(): boolean;
  clearCreatedat(): Document;

  getUpdatedat(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedat(value?: resources_timestamp_timestamp_pb.Timestamp): Document;
  hasUpdatedat(): boolean;
  clearUpdatedat(): Document;

  getTitle(): string;
  setTitle(value: string): Document;

  getContent(): string;
  setContent(value: string): Document;

  getContenttype(): DOCUMENT_CONTENT_TYPE;
  setContenttype(value: DOCUMENT_CONTENT_TYPE): Document;

  getClosed(): boolean;
  setClosed(value: boolean): Document;

  getState(): string;
  setState(value: string): Document;

  getCreator(): resources_users_users_pb.ShortUser | undefined;
  setCreator(value?: resources_users_users_pb.ShortUser): Document;
  hasCreator(): boolean;
  clearCreator(): Document;

  getPublic(): boolean;
  setPublic(value: boolean): Document;

  getCategoryid(): number;
  setCategoryid(value: number): Document;

  getTargetdocumentid(): number;
  setTargetdocumentid(value: number): Document;

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
    createdat?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedat?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    title: string,
    content: string,
    contenttype: DOCUMENT_CONTENT_TYPE,
    closed: boolean,
    state: string,
    creator?: resources_users_users_pb.ShortUser.AsObject,
    pb_public: boolean,
    categoryid: number,
    targetdocumentid: number,
  }
}

export class DocumentTemplate extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentTemplate;

  getJob(): string;
  setJob(value: string): DocumentTemplate;

  getJobgrade(): number;
  setJobgrade(value: number): DocumentTemplate;

  getTitle(): string;
  setTitle(value: string): DocumentTemplate;

  getDescription(): string;
  setDescription(value: string): DocumentTemplate;

  getContenttitle(): string;
  setContenttitle(value: string): DocumentTemplate;

  getContent(): string;
  setContent(value: string): DocumentTemplate;

  getAdditionaldata(): string;
  setAdditionaldata(value: string): DocumentTemplate;

  getCreatorid(): number;
  setCreatorid(value: number): DocumentTemplate;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentTemplate.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentTemplate): DocumentTemplate.AsObject;
  static serializeBinaryToWriter(message: DocumentTemplate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentTemplate;
  static deserializeBinaryFromReader(message: DocumentTemplate, reader: jspb.BinaryReader): DocumentTemplate;
}

export namespace DocumentTemplate {
  export type AsObject = {
    id: number,
    job: string,
    jobgrade: number,
    title: string,
    description: string,
    contenttitle: string,
    content: string,
    additionaldata: string,
    creatorid: number,
  }
}

export class DocumentTemplateShort extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentTemplateShort;

  getJob(): string;
  setJob(value: string): DocumentTemplateShort;

  getJobgrade(): number;
  setJobgrade(value: number): DocumentTemplateShort;

  getTitle(): string;
  setTitle(value: string): DocumentTemplateShort;

  getDescription(): string;
  setDescription(value: string): DocumentTemplateShort;

  getCreatorid(): number;
  setCreatorid(value: number): DocumentTemplateShort;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentTemplateShort.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentTemplateShort): DocumentTemplateShort.AsObject;
  static serializeBinaryToWriter(message: DocumentTemplateShort, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentTemplateShort;
  static deserializeBinaryFromReader(message: DocumentTemplateShort, reader: jspb.BinaryReader): DocumentTemplateShort;
}

export namespace DocumentTemplateShort {
  export type AsObject = {
    id: number,
    job: string,
    jobgrade: number,
    title: string,
    description: string,
    creatorid: number,
  }
}

export class DocumentAccess extends jspb.Message {
  getDocumentid(): number;
  setDocumentid(value: number): DocumentAccess;

  getJobaccessList(): Array<DocumentJobAccess>;
  setJobaccessList(value: Array<DocumentJobAccess>): DocumentAccess;
  clearJobaccessList(): DocumentAccess;
  addJobaccess(value?: DocumentJobAccess, index?: number): DocumentJobAccess;

  getUseraccessList(): Array<DocumentUserAccess>;
  setUseraccessList(value: Array<DocumentUserAccess>): DocumentAccess;
  clearUseraccessList(): DocumentAccess;
  addUseraccess(value?: DocumentUserAccess, index?: number): DocumentUserAccess;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentAccess.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentAccess): DocumentAccess.AsObject;
  static serializeBinaryToWriter(message: DocumentAccess, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentAccess;
  static deserializeBinaryFromReader(message: DocumentAccess, reader: jspb.BinaryReader): DocumentAccess;
}

export namespace DocumentAccess {
  export type AsObject = {
    documentid: number,
    jobaccessList: Array<DocumentJobAccess.AsObject>,
    useraccessList: Array<DocumentUserAccess.AsObject>,
  }
}

export class DocumentJobAccess extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentJobAccess;

  getCreatedat(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedat(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentJobAccess;
  hasCreatedat(): boolean;
  clearCreatedat(): DocumentJobAccess;

  getUpdatedat(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedat(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentJobAccess;
  hasUpdatedat(): boolean;
  clearUpdatedat(): DocumentJobAccess;

  getDocumentid(): number;
  setDocumentid(value: number): DocumentJobAccess;

  getJob(): string;
  setJob(value: string): DocumentJobAccess;

  getMinimumgrade(): number;
  setMinimumgrade(value: number): DocumentJobAccess;

  getAccess(): DOCUMENT_ACCESS;
  setAccess(value: DOCUMENT_ACCESS): DocumentJobAccess;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentJobAccess.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentJobAccess): DocumentJobAccess.AsObject;
  static serializeBinaryToWriter(message: DocumentJobAccess, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentJobAccess;
  static deserializeBinaryFromReader(message: DocumentJobAccess, reader: jspb.BinaryReader): DocumentJobAccess;
}

export namespace DocumentJobAccess {
  export type AsObject = {
    id: number,
    createdat?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedat?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    documentid: number,
    job: string,
    minimumgrade: number,
    access: DOCUMENT_ACCESS,
  }
}

export class DocumentUserAccess extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentUserAccess;

  getCreatedat(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedat(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentUserAccess;
  hasCreatedat(): boolean;
  clearCreatedat(): DocumentUserAccess;

  getUpdatedat(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedat(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentUserAccess;
  hasUpdatedat(): boolean;
  clearUpdatedat(): DocumentUserAccess;

  getDocumentid(): number;
  setDocumentid(value: number): DocumentUserAccess;

  getUserid(): number;
  setUserid(value: number): DocumentUserAccess;

  getAccess(): DOCUMENT_ACCESS;
  setAccess(value: DOCUMENT_ACCESS): DocumentUserAccess;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentUserAccess.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentUserAccess): DocumentUserAccess.AsObject;
  static serializeBinaryToWriter(message: DocumentUserAccess, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentUserAccess;
  static deserializeBinaryFromReader(message: DocumentUserAccess, reader: jspb.BinaryReader): DocumentUserAccess;
}

export namespace DocumentUserAccess {
  export type AsObject = {
    id: number,
    createdat?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedat?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    documentid: number,
    userid: number,
    access: DOCUMENT_ACCESS,
  }
}

export class DocumentCategory extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentCategory;

  getName(): string;
  setName(value: string): DocumentCategory;

  getDescription(): string;
  setDescription(value: string): DocumentCategory;

  getJob(): string;
  setJob(value: string): DocumentCategory;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentCategory.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentCategory): DocumentCategory.AsObject;
  static serializeBinaryToWriter(message: DocumentCategory, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentCategory;
  static deserializeBinaryFromReader(message: DocumentCategory, reader: jspb.BinaryReader): DocumentCategory;
}

export namespace DocumentCategory {
  export type AsObject = {
    id: number,
    name: string,
    description: string,
    job: string,
  }
}

export class DocumentRelation extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentRelation;

  getCreatedat(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedat(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentRelation;
  hasCreatedat(): boolean;
  clearCreatedat(): DocumentRelation;

  getUpdatedat(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedat(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentRelation;
  hasUpdatedat(): boolean;
  clearUpdatedat(): DocumentRelation;

  getDocumentid(): number;
  setDocumentid(value: number): DocumentRelation;

  getTargetuserid(): number;
  setTargetuserid(value: number): DocumentRelation;

  getRelation(): DOCUMENT_RELATION;
  setRelation(value: DOCUMENT_RELATION): DocumentRelation;

  getCauseuserid(): number;
  setCauseuserid(value: number): DocumentRelation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentRelation.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentRelation): DocumentRelation.AsObject;
  static serializeBinaryToWriter(message: DocumentRelation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentRelation;
  static deserializeBinaryFromReader(message: DocumentRelation, reader: jspb.BinaryReader): DocumentRelation;
}

export namespace DocumentRelation {
  export type AsObject = {
    id: number,
    createdat?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedat?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    documentid: number,
    targetuserid: number,
    relation: DOCUMENT_RELATION,
    causeuserid: number,
  }
}

export enum DOCUMENT_CONTENT_TYPE { 
  HTML = 0,
}
export enum DOCUMENT_ACCESS { 
  BLOCKED = 0,
  VIEW = 1,
  EDIT = 2,
  LEADER = 3,
  ADMIN = 4,
}
export enum DOCUMENT_RELATION { 
  CREATED = 0,
  MENTIONED = 1,
  CAUSED = 2,
}
