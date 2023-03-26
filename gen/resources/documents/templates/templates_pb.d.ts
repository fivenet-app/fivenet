import * as jspb from 'google-protobuf'

import * as resources_users_users_pb from '../../../resources/users/users_pb';


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
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TemplateRequirements.AsObject;
  static toObject(includeInstance: boolean, msg: TemplateRequirements): TemplateRequirements.AsObject;
  static serializeBinaryToWriter(message: TemplateRequirements, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TemplateRequirements;
  static deserializeBinaryFromReader(message: TemplateRequirements, reader: jspb.BinaryReader): TemplateRequirements;
}

export namespace TemplateRequirements {
  export type AsObject = {
  }
}

export class TemplateData extends jspb.Message {
  getActivechar(): resources_users_users_pb.User | undefined;
  setActivechar(value?: resources_users_users_pb.User): TemplateData;
  hasActivechar(): boolean;
  clearActivechar(): TemplateData;

  getUsersList(): Array<resources_users_users_pb.User>;
  setUsersList(value: Array<resources_users_users_pb.User>): TemplateData;
  clearUsersList(): TemplateData;
  addUsers(value?: resources_users_users_pb.User, index?: number): resources_users_users_pb.User;

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
    usersList: Array<resources_users_users_pb.User.AsObject>,
  }
}

