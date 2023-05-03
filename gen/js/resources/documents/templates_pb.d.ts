import * as jspb from 'google-protobuf'

import * as resources_documents_category_pb from '../../resources/documents/category_pb';
import * as resources_documents_documents_pb from '../../resources/documents/documents_pb';
import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';
import * as resources_vehicles_vehicles_pb from '../../resources/vehicles/vehicles_pb';


export class DocumentTemplate extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentTemplate;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentTemplate;
  hasCreatedAt(): boolean;
  clearCreatedAt(): DocumentTemplate;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentTemplate;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): DocumentTemplate;

  getJob(): string;
  setJob(value: string): DocumentTemplate;

  getJobGrade(): number;
  setJobGrade(value: number): DocumentTemplate;

  getCategory(): resources_documents_category_pb.DocumentCategory | undefined;
  setCategory(value?: resources_documents_category_pb.DocumentCategory): DocumentTemplate;
  hasCategory(): boolean;
  clearCategory(): DocumentTemplate;

  getTitle(): string;
  setTitle(value: string): DocumentTemplate;

  getDescription(): string;
  setDescription(value: string): DocumentTemplate;

  getContentTitle(): string;
  setContentTitle(value: string): DocumentTemplate;

  getContent(): string;
  setContent(value: string): DocumentTemplate;

  getSchema(): TemplateSchema | undefined;
  setSchema(value?: TemplateSchema): DocumentTemplate;
  hasSchema(): boolean;
  clearSchema(): DocumentTemplate;

  getCreatorId(): number;
  setCreatorId(value: number): DocumentTemplate;

  getCreator(): resources_users_users_pb.UserShort | undefined;
  setCreator(value?: resources_users_users_pb.UserShort): DocumentTemplate;
  hasCreator(): boolean;
  clearCreator(): DocumentTemplate;

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
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    job: string,
    jobGrade: number,
    category?: resources_documents_category_pb.DocumentCategory.AsObject,
    title: string,
    description: string,
    contentTitle: string,
    content: string,
    schema?: TemplateSchema.AsObject,
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
    CREATOR = 13,
  }
}

export class DocumentTemplateShort extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentTemplateShort;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentTemplateShort;
  hasCreatedAt(): boolean;
  clearCreatedAt(): DocumentTemplateShort;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DocumentTemplateShort;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): DocumentTemplateShort;

  getJob(): string;
  setJob(value: string): DocumentTemplateShort;

  getCategoryId(): number;
  setCategoryId(value: number): DocumentTemplateShort;

  getCategory(): resources_documents_category_pb.DocumentCategory | undefined;
  setCategory(value?: resources_documents_category_pb.DocumentCategory): DocumentTemplateShort;
  hasCategory(): boolean;
  clearCategory(): DocumentTemplateShort;

  getTitle(): string;
  setTitle(value: string): DocumentTemplateShort;

  getDescription(): string;
  setDescription(value: string): DocumentTemplateShort;

  getSchema(): TemplateSchema | undefined;
  setSchema(value?: TemplateSchema): DocumentTemplateShort;
  hasSchema(): boolean;
  clearSchema(): DocumentTemplateShort;

  getCreatorId(): number;
  setCreatorId(value: number): DocumentTemplateShort;

  getCreator(): resources_users_users_pb.UserShort | undefined;
  setCreator(value?: resources_users_users_pb.UserShort): DocumentTemplateShort;
  hasCreator(): boolean;
  clearCreator(): DocumentTemplateShort;

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
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    job: string,
    categoryId: number,
    category?: resources_documents_category_pb.DocumentCategory.AsObject,
    title: string,
    description: string,
    schema?: TemplateSchema.AsObject,
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

export class TemplateSchema extends jspb.Message {
  getRequirements(): TemplateRequirements | undefined;
  setRequirements(value?: TemplateRequirements): TemplateSchema;
  hasRequirements(): boolean;
  clearRequirements(): TemplateSchema;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TemplateSchema.AsObject;
  static toObject(includeInstance: boolean, msg: TemplateSchema): TemplateSchema.AsObject;
  static serializeBinaryToWriter(message: TemplateSchema, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TemplateSchema;
  static deserializeBinaryFromReader(message: TemplateSchema, reader: jspb.BinaryReader): TemplateSchema;
}

export namespace TemplateSchema {
  export type AsObject = {
    requirements?: TemplateRequirements.AsObject,
  }
}

export class TemplateRequirements extends jspb.Message {
  getDocuments(): ObjectSpecs | undefined;
  setDocuments(value?: ObjectSpecs): TemplateRequirements;
  hasDocuments(): boolean;
  clearDocuments(): TemplateRequirements;

  getUsers(): ObjectSpecs | undefined;
  setUsers(value?: ObjectSpecs): TemplateRequirements;
  hasUsers(): boolean;
  clearUsers(): TemplateRequirements;

  getVehicles(): ObjectSpecs | undefined;
  setVehicles(value?: ObjectSpecs): TemplateRequirements;
  hasVehicles(): boolean;
  clearVehicles(): TemplateRequirements;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TemplateRequirements.AsObject;
  static toObject(includeInstance: boolean, msg: TemplateRequirements): TemplateRequirements.AsObject;
  static serializeBinaryToWriter(message: TemplateRequirements, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TemplateRequirements;
  static deserializeBinaryFromReader(message: TemplateRequirements, reader: jspb.BinaryReader): TemplateRequirements;
}

export namespace TemplateRequirements {
  export type AsObject = {
    documents?: ObjectSpecs.AsObject,
    users?: ObjectSpecs.AsObject,
    vehicles?: ObjectSpecs.AsObject,
  }

  export enum DocumentsCase { 
    _DOCUMENTS_NOT_SET = 0,
    DOCUMENTS = 1,
  }

  export enum UsersCase { 
    _USERS_NOT_SET = 0,
    USERS = 2,
  }

  export enum VehiclesCase { 
    _VEHICLES_NOT_SET = 0,
    VEHICLES = 3,
  }
}

export class ObjectSpecs extends jspb.Message {
  getRequired(): boolean;
  setRequired(value: boolean): ObjectSpecs;
  hasRequired(): boolean;
  clearRequired(): ObjectSpecs;

  getMin(): number;
  setMin(value: number): ObjectSpecs;
  hasMin(): boolean;
  clearMin(): ObjectSpecs;

  getMax(): number;
  setMax(value: number): ObjectSpecs;
  hasMax(): boolean;
  clearMax(): ObjectSpecs;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ObjectSpecs.AsObject;
  static toObject(includeInstance: boolean, msg: ObjectSpecs): ObjectSpecs.AsObject;
  static serializeBinaryToWriter(message: ObjectSpecs, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ObjectSpecs;
  static deserializeBinaryFromReader(message: ObjectSpecs, reader: jspb.BinaryReader): ObjectSpecs;
}

export namespace ObjectSpecs {
  export type AsObject = {
    required?: boolean,
    min?: number,
    max?: number,
  }

  export enum RequiredCase { 
    _REQUIRED_NOT_SET = 0,
    REQUIRED = 1,
  }

  export enum MinCase { 
    _MIN_NOT_SET = 0,
    MIN = 2,
  }

  export enum MaxCase { 
    _MAX_NOT_SET = 0,
    MAX = 3,
  }
}

export class TemplateData extends jspb.Message {
  getActivechar(): resources_users_users_pb.User | undefined;
  setActivechar(value?: resources_users_users_pb.User): TemplateData;
  hasActivechar(): boolean;
  clearActivechar(): TemplateData;

  getDocumentsList(): Array<resources_documents_documents_pb.DocumentShort>;
  setDocumentsList(value: Array<resources_documents_documents_pb.DocumentShort>): TemplateData;
  clearDocumentsList(): TemplateData;
  addDocuments(value?: resources_documents_documents_pb.DocumentShort, index?: number): resources_documents_documents_pb.DocumentShort;

  getUsersList(): Array<resources_users_users_pb.User>;
  setUsersList(value: Array<resources_users_users_pb.User>): TemplateData;
  clearUsersList(): TemplateData;
  addUsers(value?: resources_users_users_pb.User, index?: number): resources_users_users_pb.User;

  getVehiclesList(): Array<resources_vehicles_vehicles_pb.Vehicle>;
  setVehiclesList(value: Array<resources_vehicles_vehicles_pb.Vehicle>): TemplateData;
  clearVehiclesList(): TemplateData;
  addVehicles(value?: resources_vehicles_vehicles_pb.Vehicle, index?: number): resources_vehicles_vehicles_pb.Vehicle;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TemplateData.AsObject;
  static toObject(includeInstance: boolean, msg: TemplateData): TemplateData.AsObject;
  static serializeBinaryToWriter(message: TemplateData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TemplateData;
  static deserializeBinaryFromReader(message: TemplateData, reader: jspb.BinaryReader): TemplateData;
}

export namespace TemplateData {
  export type AsObject = {
    activechar?: resources_users_users_pb.User.AsObject,
    documentsList: Array<resources_documents_documents_pb.DocumentShort.AsObject>,
    usersList: Array<resources_users_users_pb.User.AsObject>,
    vehiclesList: Array<resources_vehicles_vehicles_pb.Vehicle.AsObject>,
  }
}

