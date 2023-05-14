import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';


export class Permission extends jspb.Message {
  getId(): number;
  setId(value: number): Permission;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Permission;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Permission;

  getCategory(): string;
  setCategory(value: string): Permission;

  getName(): string;
  setName(value: string): Permission;

  getGuardName(): string;
  setGuardName(value: string): Permission;

  getVal(): boolean;
  setVal(value: boolean): Permission;

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
    category: string,
    name: string,
    guardName: string,
    val: boolean,
  }

  export enum CreatedAtCase { 
    _CREATED_AT_NOT_SET = 0,
    CREATED_AT = 2,
  }
}

export class Role extends jspb.Message {
  getId(): number;
  setId(value: number): Role;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Role;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Role;

  getJob(): string;
  setJob(value: string): Role;

  getJobLabel(): string;
  setJobLabel(value: string): Role;

  getGrade(): number;
  setGrade(value: number): Role;

  getJobGradeLabel(): string;
  setJobGradeLabel(value: string): Role;

  getPermissionsList(): Array<Permission>;
  setPermissionsList(value: Array<Permission>): Role;
  clearPermissionsList(): Role;
  addPermissions(value?: Permission, index?: number): Permission;

  getAttributesList(): Array<RoleAttribute>;
  setAttributesList(value: Array<RoleAttribute>): Role;
  clearAttributesList(): Role;
  addAttributes(value?: RoleAttribute, index?: number): RoleAttribute;

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
    job: string,
    jobLabel: string,
    grade: number,
    jobGradeLabel: string,
    permissionsList: Array<Permission.AsObject>,
    attributesList: Array<RoleAttribute.AsObject>,
  }

  export enum CreatedAtCase { 
    _CREATED_AT_NOT_SET = 0,
    CREATED_AT = 2,
  }
}

export class Attribute extends jspb.Message {
  getId(): number;
  setId(value: number): Attribute;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Attribute;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Attribute;

  getPermissionId(): number;
  setPermissionId(value: number): Attribute;

  getKey(): string;
  setKey(value: string): Attribute;

  getType(): string;
  setType(value: string): Attribute;

  getValue(): string;
  setValue(value: string): Attribute;

  getValidvaluesList(): Array<string>;
  setValidvaluesList(value: Array<string>): Attribute;
  clearValidvaluesList(): Attribute;
  addValidvalues(value: string, index?: number): Attribute;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Attribute.AsObject;
  static toObject(includeInstance: boolean, msg: Attribute): Attribute.AsObject;
  static serializeBinaryToWriter(message: Attribute, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Attribute;
  static deserializeBinaryFromReader(message: Attribute, reader: jspb.BinaryReader): Attribute;
}

export namespace Attribute {
  export type AsObject = {
    id: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    permissionId: number,
    key: string,
    type: string,
    value: string,
    validvaluesList: Array<string>,
  }

  export enum CreatedAtCase { 
    _CREATED_AT_NOT_SET = 0,
    CREATED_AT = 2,
  }
}

export class RoleAttribute extends jspb.Message {
  getRoleId(): number;
  setRoleId(value: number): RoleAttribute;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): RoleAttribute;
  hasCreatedAt(): boolean;
  clearCreatedAt(): RoleAttribute;

  getAttrId(): number;
  setAttrId(value: number): RoleAttribute;

  getPermissionId(): number;
  setPermissionId(value: number): RoleAttribute;

  getCategory(): string;
  setCategory(value: string): RoleAttribute;

  getName(): string;
  setName(value: string): RoleAttribute;

  getKey(): string;
  setKey(value: string): RoleAttribute;

  getType(): string;
  setType(value: string): RoleAttribute;

  getValue(): string;
  setValue(value: string): RoleAttribute;

  getValidValues(): string;
  setValidValues(value: string): RoleAttribute;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RoleAttribute.AsObject;
  static toObject(includeInstance: boolean, msg: RoleAttribute): RoleAttribute.AsObject;
  static serializeBinaryToWriter(message: RoleAttribute, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RoleAttribute;
  static deserializeBinaryFromReader(message: RoleAttribute, reader: jspb.BinaryReader): RoleAttribute;
}

export namespace RoleAttribute {
  export type AsObject = {
    roleId: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    attrId: number,
    permissionId: number,
    category: string,
    name: string,
    key: string,
    type: string,
    value: string,
    validValues: string,
  }

  export enum CreatedAtCase { 
    _CREATED_AT_NOT_SET = 0,
    CREATED_AT = 2,
  }
}

