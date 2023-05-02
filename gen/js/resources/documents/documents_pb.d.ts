import * as jspb from 'google-protobuf'

import * as resources_documents_category_pb from '../../resources/documents/category_pb';
import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';


export class DocumentComment extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentComment;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentComment;
  hasCreatedAt(): boolean;
  clearCreatedAt(): DocumentComment;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentComment;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): DocumentComment;

  getDocumentId(): number;
  setDocumentId(value: number): DocumentComment;

  getComment(): string;
  setComment(value: string): DocumentComment;

  getCreatorId(): number;
  setCreatorId(value: number): DocumentComment;

  getCreator(): resources_users_users_pb.UserShort | undefined;
  setCreator(value?: resources_users_users_pb.UserShort): DocumentComment;
  hasCreator(): boolean;
  clearCreator(): DocumentComment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentComment.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentComment): DocumentComment.AsObject;
  static serializeBinaryToWriter(message: DocumentComment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentComment;
  static deserializeBinaryFromReader(message: DocumentComment, reader: jspb.BinaryReader): DocumentComment;
}

export namespace DocumentComment {
  export type AsObject = {
    id: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    documentId: number,
    comment: string,
    creatorId: number,
    creator?: resources_users_users_pb.UserShort.AsObject,
  }

  export enum CreatedAtCase { 
    _CREATED_AT_NOT_SET = 0,
    CREATED_AT = 2,
  }

  export enum UpdatedAtCase { 
    _UPDATED_AT_NOT_SET = 0,
    UPDATED_AT = 3,
  }

  export enum CreatorCase { 
    _CREATOR_NOT_SET = 0,
    CREATOR = 7,
  }
}

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

  getCategoryId(): number;
  setCategoryId(value: number): Document;

  getCategory(): resources_documents_category_pb.DocumentCategory | undefined;
  setCategory(value?: resources_documents_category_pb.DocumentCategory): Document;
  hasCategory(): boolean;
  clearCategory(): Document;

  getTitle(): string;
  setTitle(value: string): Document;

  getContentType(): DOC_CONTENT_TYPE;
  setContentType(value: DOC_CONTENT_TYPE): Document;

  getContent(): string;
  setContent(value: string): Document;

  getData(): string;
  setData(value: string): Document;

  getCreatorId(): number;
  setCreatorId(value: number): Document;

  getCreator(): resources_users_users_pb.UserShort | undefined;
  setCreator(value?: resources_users_users_pb.UserShort): Document;
  hasCreator(): boolean;
  clearCreator(): Document;

  getState(): string;
  setState(value: string): Document;

  getClosed(): boolean;
  setClosed(value: boolean): Document;

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
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    categoryId: number,
    category?: resources_documents_category_pb.DocumentCategory.AsObject,
    title: string,
    contentType: DOC_CONTENT_TYPE,
    content: string,
    data: string,
    creatorId: number,
    creator?: resources_users_users_pb.UserShort.AsObject,
    state: string,
    closed: boolean,
    pb_public: boolean,
  }
}

export class DocumentShort extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentShort;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentShort;
  hasCreatedAt(): boolean;
  clearCreatedAt(): DocumentShort;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentShort;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): DocumentShort;

  getCategoryId(): number;
  setCategoryId(value: number): DocumentShort;

  getCategory(): resources_documents_category_pb.DocumentCategory | undefined;
  setCategory(value?: resources_documents_category_pb.DocumentCategory): DocumentShort;
  hasCategory(): boolean;
  clearCategory(): DocumentShort;

  getTitle(): string;
  setTitle(value: string): DocumentShort;

  getCreatorId(): number;
  setCreatorId(value: number): DocumentShort;

  getCreator(): resources_users_users_pb.UserShort | undefined;
  setCreator(value?: resources_users_users_pb.UserShort): DocumentShort;
  hasCreator(): boolean;
  clearCreator(): DocumentShort;

  getState(): string;
  setState(value: string): DocumentShort;

  getClosed(): boolean;
  setClosed(value: boolean): DocumentShort;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentShort.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentShort): DocumentShort.AsObject;
  static serializeBinaryToWriter(message: DocumentShort, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentShort;
  static deserializeBinaryFromReader(message: DocumentShort, reader: jspb.BinaryReader): DocumentShort;
}

export namespace DocumentShort {
  export type AsObject = {
    id: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    categoryId: number,
    category?: resources_documents_category_pb.DocumentCategory.AsObject,
    title: string,
    creatorId: number,
    creator?: resources_users_users_pb.UserShort.AsObject,
    state: string,
    closed: boolean,
  }
}

export class DocumentAccess extends jspb.Message {
  getJobsList(): Array<DocumentJobAccess>;
  setJobsList(value: Array<DocumentJobAccess>): DocumentAccess;
  clearJobsList(): DocumentAccess;
  addJobs(value?: DocumentJobAccess, index?: number): DocumentJobAccess;

  getUsersList(): Array<DocumentUserAccess>;
  setUsersList(value: Array<DocumentUserAccess>): DocumentAccess;
  clearUsersList(): DocumentAccess;
  addUsers(value?: DocumentUserAccess, index?: number): DocumentUserAccess;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentAccess.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentAccess): DocumentAccess.AsObject;
  static serializeBinaryToWriter(message: DocumentAccess, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentAccess;
  static deserializeBinaryFromReader(message: DocumentAccess, reader: jspb.BinaryReader): DocumentAccess;
}

export namespace DocumentAccess {
  export type AsObject = {
    jobsList: Array<DocumentJobAccess.AsObject>,
    usersList: Array<DocumentUserAccess.AsObject>,
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

  getDocumentId(): number;
  setDocumentId(value: number): DocumentJobAccess;

  getJob(): string;
  setJob(value: string): DocumentJobAccess;

  getJobLabel(): string;
  setJobLabel(value: string): DocumentJobAccess;

  getMinimumgrade(): number;
  setMinimumgrade(value: number): DocumentJobAccess;

  getJobGradeLabel(): string;
  setJobGradeLabel(value: string): DocumentJobAccess;

  getAccess(): DOC_ACCESS;
  setAccess(value: DOC_ACCESS): DocumentJobAccess;

  getCreatorId(): number;
  setCreatorId(value: number): DocumentJobAccess;

  getCreator(): resources_users_users_pb.UserShort | undefined;
  setCreator(value?: resources_users_users_pb.UserShort): DocumentJobAccess;
  hasCreator(): boolean;
  clearCreator(): DocumentJobAccess;

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
    documentId: number,
    job: string,
    jobLabel: string,
    minimumgrade: number,
    jobGradeLabel: string,
    access: DOC_ACCESS,
    creatorId: number,
    creator?: resources_users_users_pb.UserShort.AsObject,
  }

  export enum CreatedAtCase { 
    _CREATED_AT_NOT_SET = 0,
    CREATED_AT = 2,
  }

  export enum UpdatedAtCase { 
    _UPDATED_AT_NOT_SET = 0,
    UPDATED_AT = 3,
  }

  export enum CreatorCase { 
    _CREATOR_NOT_SET = 0,
    CREATOR = 11,
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

  getDocumentId(): number;
  setDocumentId(value: number): DocumentUserAccess;

  getUserId(): number;
  setUserId(value: number): DocumentUserAccess;

  getUser(): resources_users_users_pb.UserShort | undefined;
  setUser(value?: resources_users_users_pb.UserShort): DocumentUserAccess;
  hasUser(): boolean;
  clearUser(): DocumentUserAccess;

  getAccess(): DOC_ACCESS;
  setAccess(value: DOC_ACCESS): DocumentUserAccess;

  getCreatorId(): number;
  setCreatorId(value: number): DocumentUserAccess;

  getCreator(): resources_users_users_pb.UserShort | undefined;
  setCreator(value?: resources_users_users_pb.UserShort): DocumentUserAccess;
  hasCreator(): boolean;
  clearCreator(): DocumentUserAccess;

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
    documentId: number,
    userId: number,
    user?: resources_users_users_pb.UserShort.AsObject,
    access: DOC_ACCESS,
    creatorId: number,
    creator?: resources_users_users_pb.UserShort.AsObject,
  }

  export enum CreatedAtCase { 
    _CREATED_AT_NOT_SET = 0,
    CREATED_AT = 2,
  }

  export enum UpdatedAtCase { 
    _UPDATED_AT_NOT_SET = 0,
    UPDATED_AT = 3,
  }

  export enum UserCase { 
    _USER_NOT_SET = 0,
    USER = 6,
  }

  export enum CreatorCase { 
    _CREATOR_NOT_SET = 0,
    CREATOR = 9,
  }
}

export class DocumentReference extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentReference;
  hasId(): boolean;
  clearId(): DocumentReference;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentReference;
  hasCreatedAt(): boolean;
  clearCreatedAt(): DocumentReference;

  getSourceDocumentId(): number;
  setSourceDocumentId(value: number): DocumentReference;

  getSourceDocument(): DocumentShort | undefined;
  setSourceDocument(value?: DocumentShort): DocumentReference;
  hasSourceDocument(): boolean;
  clearSourceDocument(): DocumentReference;

  getReference(): DOC_REFERENCE;
  setReference(value: DOC_REFERENCE): DocumentReference;

  getTargetDocumentId(): number;
  setTargetDocumentId(value: number): DocumentReference;

  getTargetDocument(): DocumentShort | undefined;
  setTargetDocument(value?: DocumentShort): DocumentReference;
  hasTargetDocument(): boolean;
  clearTargetDocument(): DocumentReference;

  getCreatorId(): number;
  setCreatorId(value: number): DocumentReference;

  getCreator(): resources_users_users_pb.UserShort | undefined;
  setCreator(value?: resources_users_users_pb.UserShort): DocumentReference;
  hasCreator(): boolean;
  clearCreator(): DocumentReference;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentReference.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentReference): DocumentReference.AsObject;
  static serializeBinaryToWriter(message: DocumentReference, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentReference;
  static deserializeBinaryFromReader(message: DocumentReference, reader: jspb.BinaryReader): DocumentReference;
}

export namespace DocumentReference {
  export type AsObject = {
    id?: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    sourceDocumentId: number,
    sourceDocument?: DocumentShort.AsObject,
    reference: DOC_REFERENCE,
    targetDocumentId: number,
    targetDocument?: DocumentShort.AsObject,
    creatorId: number,
    creator?: resources_users_users_pb.UserShort.AsObject,
  }

  export enum IdCase { 
    _ID_NOT_SET = 0,
    ID = 1,
  }

  export enum CreatedAtCase { 
    _CREATED_AT_NOT_SET = 0,
    CREATED_AT = 2,
  }

  export enum SourceDocumentCase { 
    _SOURCE_DOCUMENT_NOT_SET = 0,
    SOURCE_DOCUMENT = 4,
  }

  export enum CreatorCase { 
    _CREATOR_NOT_SET = 0,
    CREATOR = 9,
  }
}

export class DocumentRelation extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentRelation;
  hasId(): boolean;
  clearId(): DocumentRelation;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentRelation;
  hasCreatedAt(): boolean;
  clearCreatedAt(): DocumentRelation;

  getDocumentId(): number;
  setDocumentId(value: number): DocumentRelation;

  getDocument(): DocumentShort | undefined;
  setDocument(value?: DocumentShort): DocumentRelation;
  hasDocument(): boolean;
  clearDocument(): DocumentRelation;

  getSourceUserId(): number;
  setSourceUserId(value: number): DocumentRelation;

  getSourceUser(): resources_users_users_pb.UserShort | undefined;
  setSourceUser(value?: resources_users_users_pb.UserShort): DocumentRelation;
  hasSourceUser(): boolean;
  clearSourceUser(): DocumentRelation;

  getRelation(): DOC_RELATION;
  setRelation(value: DOC_RELATION): DocumentRelation;

  getTargetUserId(): number;
  setTargetUserId(value: number): DocumentRelation;

  getTargetUser(): resources_users_users_pb.UserShort | undefined;
  setTargetUser(value?: resources_users_users_pb.UserShort): DocumentRelation;
  hasTargetUser(): boolean;
  clearTargetUser(): DocumentRelation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentRelation.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentRelation): DocumentRelation.AsObject;
  static serializeBinaryToWriter(message: DocumentRelation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentRelation;
  static deserializeBinaryFromReader(message: DocumentRelation, reader: jspb.BinaryReader): DocumentRelation;
}

export namespace DocumentRelation {
  export type AsObject = {
    id?: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    documentId: number,
    document?: DocumentShort.AsObject,
    sourceUserId: number,
    sourceUser?: resources_users_users_pb.UserShort.AsObject,
    relation: DOC_RELATION,
    targetUserId: number,
    targetUser?: resources_users_users_pb.UserShort.AsObject,
  }

  export enum IdCase { 
    _ID_NOT_SET = 0,
    ID = 1,
  }

  export enum CreatedAtCase { 
    _CREATED_AT_NOT_SET = 0,
    CREATED_AT = 2,
  }

  export enum DocumentCase { 
    _DOCUMENT_NOT_SET = 0,
    DOCUMENT = 4,
  }

  export enum SourceUserCase { 
    _SOURCE_USER_NOT_SET = 0,
    SOURCE_USER = 6,
  }

  export enum TargetUserCase { 
    _TARGET_USER_NOT_SET = 0,
    TARGET_USER = 9,
  }
}

export enum DOC_CONTENT_TYPE { 
  HTML = 0,
  PLAIN = 1,
}
export enum DOC_ACCESS { 
  BLOCKED = 0,
  VIEW = 1,
  COMMENT = 2,
  ACCESS = 3,
  EDIT = 4,
}
export enum DOC_REFERENCE { 
  LINKED = 0,
  SOLVES = 1,
  CLOSES = 2,
  DEPRECATES = 3,
}
export enum DOC_RELATION { 
  MENTIONED = 0,
  TARGETS = 1,
  CAUSED = 2,
}
