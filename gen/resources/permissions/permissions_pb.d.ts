import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';


export class Permission extends jspb.Message {
  getId(): number;
  setId(value: number): Permission;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Permission;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Permission;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Permission;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): Permission;

  getName(): string;
  setName(value: string): Permission;

  getGuardName(): string;
  setGuardName(value: string): Permission;

  getDescription(): string;
  setDescription(value: string): Permission;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Permission.AsObject;
  static toObject(includeInstance: boolean, msg: Permission): Permission.AsObject;
  static serializeBinaryToWriter(message: Permission, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Permission;
  static deserializeBinaryFromReader(message: Permission, reader: jspb.BinaryReader): Permission;
}

export namespace Permission {
  export type AsObject = {
    id: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    name: string,
    guardName: string,
    description: string,
  }
}

export class Role extends jspb.Message {
  getId(): number;
  setId(value: number): Role;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Role;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Role;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Role;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): Role;

  getName(): string;
  setName(value: string): Role;

  getGuardName(): string;
  setGuardName(value: string): Role;

  getDescription(): string;
  setDescription(value: string): Role;

  getPermissionsList(): Array<Permission>;
  setPermissionsList(value: Array<Permission>): Role;
  clearPermissionsList(): Role;
  addPermissions(value?: Permission, index?: number): Permission;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Role.AsObject;
  static toObject(includeInstance: boolean, msg: Role): Role.AsObject;
  static serializeBinaryToWriter(message: Role, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Role;
  static deserializeBinaryFromReader(message: Role, reader: jspb.BinaryReader): Role;
}

export namespace Role {
  export type AsObject = {
    id: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    name: string,
    guardName: string,
    description: string,
    permissionsList: Array<Permission.AsObject>,
  }
}

