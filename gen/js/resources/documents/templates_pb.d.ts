import * as jspb from 'google-protobuf'

import * as resources_documents_access_pb from '../../resources/documents/access_pb';
import * as resources_documents_category_pb from '../../resources/documents/category_pb';
import * as resources_documents_documents_pb from '../../resources/documents/documents_pb';
import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';
import * as resources_vehicles_vehicles_pb from '../../resources/vehicles/vehicles_pb';


export class Template extends jspb.Message {
  getId(): number;
  setId(value: number): Template;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Template;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Template;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Template;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): Template;

  getCategory(): resources_documents_category_pb.DocumentCategory | undefined;
  setCategory(value?: resources_documents_category_pb.DocumentCategory): Template;
  hasCategory(): boolean;
  clearCategory(): Template;

  getTitle(): string;
  setTitle(value: string): Template;

  getDescription(): string;
  setDescription(value: string): Template;

  getContentTitle(): string;
  setContentTitle(value: string): Template;

  getContent(): string;
  setContent(value: string): Template;

  getSchema(): TemplateSchema | undefined;
  setSchema(value?: TemplateSchema): Template;
  hasSchema(): boolean;
  clearSchema(): Template;

  getCreatorId(): number;
  setCreatorId(value: number): Template;

  getCreator(): resources_users_users_pb.UserShort | undefined;
  setCreator(value?: resources_users_users_pb.UserShort): Template;
  hasCreator(): boolean;
  clearCreator(): Template;

  getJob(): string;
  setJob(value: string): Template;

  getJobAccessList(): Array<TemplateJobAccess>;
  setJobAccessList(value: Array<TemplateJobAccess>): Template;
  clearJobAccessList(): Template;
  addJobAccess(value?: TemplateJobAccess, index?: number): TemplateJobAccess;

  getContentAccess(): resources_documents_documents_pb.DocumentAccess | undefined;
  setContentAccess(value?: resources_documents_documents_pb.DocumentAccess): Template;
  hasContentAccess(): boolean;
  clearContentAccess(): Template;

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
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    category?: resources_documents_category_pb.DocumentCategory.AsObject,
    title: string,
    description: string,
    contentTitle: string,
    content: string,
    schema?: TemplateSchema.AsObject,
    creatorId: number,
    creator?: resources_users_users_pb.UserShort.AsObject,
    job: string,
    jobAccessList: Array<TemplateJobAccess.AsObject>,
    contentAccess?: resources_documents_documents_pb.DocumentAccess.AsObject,
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

export class TemplateShort extends jspb.Message {
  getId(): number;
  setId(value: number): TemplateShort;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): TemplateShort;
  hasCreatedAt(): boolean;
  clearCreatedAt(): TemplateShort;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): TemplateShort;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): TemplateShort;

  getCategory(): resources_documents_category_pb.DocumentCategory | undefined;
  setCategory(value?: resources_documents_category_pb.DocumentCategory): TemplateShort;
  hasCategory(): boolean;
  clearCategory(): TemplateShort;

  getTitle(): string;
  setTitle(value: string): TemplateShort;

  getDescription(): string;
  setDescription(value: string): TemplateShort;

  getSchema(): TemplateSchema | undefined;
  setSchema(value?: TemplateSchema): TemplateShort;
  hasSchema(): boolean;
  clearSchema(): TemplateShort;

  getCreatorId(): number;
  setCreatorId(value: number): TemplateShort;

  getCreator(): resources_users_users_pb.UserShort | undefined;
  setCreator(value?: resources_users_users_pb.UserShort): TemplateShort;
  hasCreator(): boolean;
  clearCreator(): TemplateShort;

  getJob(): string;
  setJob(value: string): TemplateShort;

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
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    category?: resources_documents_category_pb.DocumentCategory.AsObject,
    title: string,
    description: string,
    schema?: TemplateSchema.AsObject,
    creatorId: number,
    creator?: resources_users_users_pb.UserShort.AsObject,
    job: string,
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
    CREATOR = 9,
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

export class TemplateJobAccess extends jspb.Message {
  getId(): number;
  setId(value: number): TemplateJobAccess;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): TemplateJobAccess;
  hasCreatedAt(): boolean;
  clearCreatedAt(): TemplateJobAccess;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): TemplateJobAccess;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): TemplateJobAccess;

  getTemplateId(): number;
  setTemplateId(value: number): TemplateJobAccess;

  getJob(): string;
  setJob(value: string): TemplateJobAccess;

  getJobLabel(): string;
  setJobLabel(value: string): TemplateJobAccess;

  getMinimumgrade(): number;
  setMinimumgrade(value: number): TemplateJobAccess;

  getJobGradeLabel(): string;
  setJobGradeLabel(value: string): TemplateJobAccess;

  getAccess(): resources_documents_access_pb.ACCESS_LEVEL;
  setAccess(value: resources_documents_access_pb.ACCESS_LEVEL): TemplateJobAccess;

  getCreatorId(): number;
  setCreatorId(value: number): TemplateJobAccess;

  getCreator(): resources_users_users_pb.UserShort | undefined;
  setCreator(value?: resources_users_users_pb.UserShort): TemplateJobAccess;
  hasCreator(): boolean;
  clearCreator(): TemplateJobAccess;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TemplateJobAccess.AsObject;
  static toObject(includeInstance: boolean, msg: TemplateJobAccess): TemplateJobAccess.AsObject;
  static serializeBinaryToWriter(message: TemplateJobAccess, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TemplateJobAccess;
  static deserializeBinaryFromReader(message: TemplateJobAccess, reader: jspb.BinaryReader): TemplateJobAccess;
}

export namespace TemplateJobAccess {
  export type AsObject = {
    id: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    templateId: number,
    job: string,
    jobLabel: string,
    minimumgrade: number,
    jobGradeLabel: string,
    access: resources_documents_access_pb.ACCESS_LEVEL,
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

