import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';


export class Document extends jspb.Message {
  getId(): number;
  setId(value: number): Document;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Document;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Document;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Document;
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

  getCreator(): resources_users_users_pb.ShortUser | undefined;
  setCreator(value?: resources_users_users_pb.ShortUser): Document;
  hasCreator(): boolean;
  clearCreator(): Document;

  getCreatorJob(): string;
  setCreatorJob(value: string): Document;

  getPublic(): boolean;
  setPublic(value: boolean): Document;

  getTargetDocumentId(): number;
  setTargetDocumentId(value: number): Document;

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
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    title: string,
    content: string,
    contentType: string,
    closed: boolean,
    state: string,
    creator?: resources_users_users_pb.ShortUser.AsObject,
    creatorJob: string,
    pb_public: boolean,
    targetDocumentId: number,
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

  getContentTitle(): string;
  setContentTitle(value: string): DocumentTemplate;

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
    contentTitle: string,
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

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentJobAccess;
  hasCreatedAt(): boolean;
  clearCreatedAt(): DocumentJobAccess;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentJobAccess;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): DocumentJobAccess;

  getDocumentid(): number;
  setDocumentid(value: number): DocumentJobAccess;

  getName(): string;
  setName(value: string): DocumentJobAccess;

  getMinimumgrade(): number;
  setMinimumgrade(value: number): DocumentJobAccess;

  getAccess(): string;
  setAccess(value: string): DocumentJobAccess;

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
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    documentid: number,
    name: string,
    minimumgrade: number,
    access: string,
  }
}

export class DocumentUserAccess extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentUserAccess;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentUserAccess;
  hasCreatedAt(): boolean;
  clearCreatedAt(): DocumentUserAccess;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentUserAccess;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): DocumentUserAccess;

  getDocumentid(): number;
  setDocumentid(value: number): DocumentUserAccess;

  getUserid(): number;
  setUserid(value: number): DocumentUserAccess;

  getAccess(): string;
  setAccess(value: string): DocumentUserAccess;

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
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    documentid: number,
    userid: number,
    access: string,
  }
}

