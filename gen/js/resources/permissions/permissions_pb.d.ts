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

export class RawAttribute extends jspb.Message {
  getRoleId(): number;
  setRoleId(value: number): RawAttribute;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): RawAttribute;
  hasCreatedAt(): boolean;
  clearCreatedAt(): RawAttribute;

  getAttrId(): number;
  setAttrId(value: number): RawAttribute;

  getPermissionId(): number;
  setPermissionId(value: number): RawAttribute;

  getCategory(): string;
  setCategory(value: string): RawAttribute;

  getName(): string;
  setName(value: string): RawAttribute;

  getKey(): string;
  setKey(value: string): RawAttribute;

  getType(): string;
  setType(value: string): RawAttribute;

  getRawValue(): string;
  setRawValue(value: string): RawAttribute;

  getRawValidValues(): string;
  setRawValidValues(value: string): RawAttribute;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RawAttribute.AsObject;
  static toObject(includeInstance: boolean, msg: RawAttribute): RawAttribute.AsObject;
  static serializeBinaryToWriter(message: RawAttribute, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RawAttribute;
  static deserializeBinaryFromReader(message: RawAttribute, reader: jspb.BinaryReader): RawAttribute;
}

export namespace RawAttribute {
  export type AsObject = {
    roleId: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    attrId: number,
    permissionId: number,
    category: string,
    name: string,
    key: string,
    type: string,
    rawValue: string,
    rawValidValues: string,
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

  getValue(): AttributeValues | undefined;
  setValue(value?: AttributeValues): RoleAttribute;
  hasValue(): boolean;
  clearValue(): RoleAttribute;

  getValidValues(): AttributeValues | undefined;
  setValidValues(value?: AttributeValues): RoleAttribute;
  hasValidValues(): boolean;
  clearValidValues(): RoleAttribute;

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
    value?: AttributeValues.AsObject,
    validValues?: AttributeValues.AsObject,
  }

  export enum CreatedAtCase { 
    _CREATED_AT_NOT_SET = 0,
    CREATED_AT = 2,
  }
}

export class AttributeValues extends jspb.Message {
  getStringList(): StringList | undefined;
  setStringList(value?: StringList): AttributeValues;
  hasStringList(): boolean;
  clearStringList(): AttributeValues;

  getJobList(): StringList | undefined;
  setJobList(value?: StringList): AttributeValues;
  hasJobList(): boolean;
  clearJobList(): AttributeValues;

  getJobGradeList(): JobGradeList | undefined;
  setJobGradeList(value?: JobGradeList): AttributeValues;
  hasJobGradeList(): boolean;
  clearJobGradeList(): AttributeValues;

  getValidValuesCase(): AttributeValues.ValidValuesCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AttributeValues.AsObject;
  static toObject(includeInstance: boolean, msg: AttributeValues): AttributeValues.AsObject;
  static serializeBinaryToWriter(message: AttributeValues, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AttributeValues;
  static deserializeBinaryFromReader(message: AttributeValues, reader: jspb.BinaryReader): AttributeValues;
}

export namespace AttributeValues {
  export type AsObject = {
    stringList?: StringList.AsObject,
    jobList?: StringList.AsObject,
    jobGradeList?: JobGradeList.AsObject,
  }

  export enum ValidValuesCase { 
    VALID_VALUES_NOT_SET = 0,
    STRING_LIST = 1,
    JOB_LIST = 2,
    JOB_GRADE_LIST = 3,
  }
}

export class StringList extends jspb.Message {
  getStringsList(): Array<string>;
  setStringsList(value: Array<string>): StringList;
  clearStringsList(): StringList;
  addStrings(value: string, index?: number): StringList;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StringList.AsObject;
  static toObject(includeInstance: boolean, msg: StringList): StringList.AsObject;
  static serializeBinaryToWriter(message: StringList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StringList;
  static deserializeBinaryFromReader(message: StringList, reader: jspb.BinaryReader): StringList;
}

export namespace StringList {
  export type AsObject = {
    stringsList: Array<string>,
  }
}

export class JobGradeList extends jspb.Message {
  getJobsMap(): jspb.Map<string, number>;
  clearJobsMap(): JobGradeList;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JobGradeList.AsObject;
  static toObject(includeInstance: boolean, msg: JobGradeList): JobGradeList.AsObject;
  static serializeBinaryToWriter(message: JobGradeList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JobGradeList;
  static deserializeBinaryFromReader(message: JobGradeList, reader: jspb.BinaryReader): JobGradeList;
}

export namespace JobGradeList {
  export type AsObject = {
    jobsMap: Array<[string, number]>,
  }
}

