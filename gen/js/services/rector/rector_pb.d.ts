import * as jspb from 'google-protobuf'

import * as resources_common_database_database_pb from '../../resources/common/database/database_pb';
import * as resources_jobs_jobs_pb from '../../resources/jobs/jobs_pb';
import * as resources_permissions_permissions_pb from '../../resources/permissions/permissions_pb';
import * as resources_rector_audit_pb from '../../resources/rector/audit_pb';
import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';


export class GetJobPropsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetJobPropsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetJobPropsRequest): GetJobPropsRequest.AsObject;
  static serializeBinaryToWriter(message: GetJobPropsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetJobPropsRequest;
  static deserializeBinaryFromReader(message: GetJobPropsRequest, reader: jspb.BinaryReader): GetJobPropsRequest;
}

export namespace GetJobPropsRequest {
  export type AsObject = {
  }
}

export class GetJobPropsResponse extends jspb.Message {
  getJobProps(): resources_jobs_jobs_pb.JobProps | undefined;
  setJobProps(value?: resources_jobs_jobs_pb.JobProps): GetJobPropsResponse;
  hasJobProps(): boolean;
  clearJobProps(): GetJobPropsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetJobPropsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetJobPropsResponse): GetJobPropsResponse.AsObject;
  static serializeBinaryToWriter(message: GetJobPropsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetJobPropsResponse;
  static deserializeBinaryFromReader(message: GetJobPropsResponse, reader: jspb.BinaryReader): GetJobPropsResponse;
}

export namespace GetJobPropsResponse {
  export type AsObject = {
    jobProps?: resources_jobs_jobs_pb.JobProps.AsObject,
  }
}

export class SetJobPropsRequest extends jspb.Message {
  getJobProps(): resources_jobs_jobs_pb.JobProps | undefined;
  setJobProps(value?: resources_jobs_jobs_pb.JobProps): SetJobPropsRequest;
  hasJobProps(): boolean;
  clearJobProps(): SetJobPropsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetJobPropsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetJobPropsRequest): SetJobPropsRequest.AsObject;
  static serializeBinaryToWriter(message: SetJobPropsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetJobPropsRequest;
  static deserializeBinaryFromReader(message: SetJobPropsRequest, reader: jspb.BinaryReader): SetJobPropsRequest;
}

export namespace SetJobPropsRequest {
  export type AsObject = {
    jobProps?: resources_jobs_jobs_pb.JobProps.AsObject,
  }
}

export class SetJobPropsResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetJobPropsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetJobPropsResponse): SetJobPropsResponse.AsObject;
  static serializeBinaryToWriter(message: SetJobPropsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetJobPropsResponse;
  static deserializeBinaryFromReader(message: SetJobPropsResponse, reader: jspb.BinaryReader): SetJobPropsResponse;
}

export namespace SetJobPropsResponse {
  export type AsObject = {
  }
}

export class GetRolesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRolesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetRolesRequest): GetRolesRequest.AsObject;
  static serializeBinaryToWriter(message: GetRolesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRolesRequest;
  static deserializeBinaryFromReader(message: GetRolesRequest, reader: jspb.BinaryReader): GetRolesRequest;
}

export namespace GetRolesRequest {
  export type AsObject = {
  }
}

export class GetRolesResponse extends jspb.Message {
  getRolesList(): Array<resources_permissions_permissions_pb.Role>;
  setRolesList(value: Array<resources_permissions_permissions_pb.Role>): GetRolesResponse;
  clearRolesList(): GetRolesResponse;
  addRoles(value?: resources_permissions_permissions_pb.Role, index?: number): resources_permissions_permissions_pb.Role;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRolesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetRolesResponse): GetRolesResponse.AsObject;
  static serializeBinaryToWriter(message: GetRolesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRolesResponse;
  static deserializeBinaryFromReader(message: GetRolesResponse, reader: jspb.BinaryReader): GetRolesResponse;
}

export namespace GetRolesResponse {
  export type AsObject = {
    rolesList: Array<resources_permissions_permissions_pb.Role.AsObject>,
  }
}

export class GetRoleRequest extends jspb.Message {
  getId(): number;
  setId(value: number): GetRoleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetRoleRequest): GetRoleRequest.AsObject;
  static serializeBinaryToWriter(message: GetRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRoleRequest;
  static deserializeBinaryFromReader(message: GetRoleRequest, reader: jspb.BinaryReader): GetRoleRequest;
}

export namespace GetRoleRequest {
  export type AsObject = {
    id: number,
  }
}

export class GetRoleResponse extends jspb.Message {
  getRole(): resources_permissions_permissions_pb.Role | undefined;
  setRole(value?: resources_permissions_permissions_pb.Role): GetRoleResponse;
  hasRole(): boolean;
  clearRole(): GetRoleResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetRoleResponse): GetRoleResponse.AsObject;
  static serializeBinaryToWriter(message: GetRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRoleResponse;
  static deserializeBinaryFromReader(message: GetRoleResponse, reader: jspb.BinaryReader): GetRoleResponse;
}

export namespace GetRoleResponse {
  export type AsObject = {
    role?: resources_permissions_permissions_pb.Role.AsObject,
  }
}

export class CreateRoleRequest extends jspb.Message {
  getJob(): string;
  setJob(value: string): CreateRoleRequest;

  getGrade(): number;
  setGrade(value: number): CreateRoleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRoleRequest): CreateRoleRequest.AsObject;
  static serializeBinaryToWriter(message: CreateRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRoleRequest;
  static deserializeBinaryFromReader(message: CreateRoleRequest, reader: jspb.BinaryReader): CreateRoleRequest;
}

export namespace CreateRoleRequest {
  export type AsObject = {
    job: string,
    grade: number,
  }
}

export class CreateRoleResponse extends jspb.Message {
  getRole(): resources_permissions_permissions_pb.Role | undefined;
  setRole(value?: resources_permissions_permissions_pb.Role): CreateRoleResponse;
  hasRole(): boolean;
  clearRole(): CreateRoleResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRoleResponse): CreateRoleResponse.AsObject;
  static serializeBinaryToWriter(message: CreateRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRoleResponse;
  static deserializeBinaryFromReader(message: CreateRoleResponse, reader: jspb.BinaryReader): CreateRoleResponse;
}

export namespace CreateRoleResponse {
  export type AsObject = {
    role?: resources_permissions_permissions_pb.Role.AsObject,
  }
}

export class DeleteRoleRequest extends jspb.Message {
  getId(): number;
  setId(value: number): DeleteRoleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRoleRequest): DeleteRoleRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRoleRequest;
  static deserializeBinaryFromReader(message: DeleteRoleRequest, reader: jspb.BinaryReader): DeleteRoleRequest;
}

export namespace DeleteRoleRequest {
  export type AsObject = {
    id: number,
  }
}

export class DeleteRoleResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRoleResponse): DeleteRoleResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRoleResponse;
  static deserializeBinaryFromReader(message: DeleteRoleResponse, reader: jspb.BinaryReader): DeleteRoleResponse;
}

export namespace DeleteRoleResponse {
  export type AsObject = {
  }
}

export class UpdateRolePermsRequest extends jspb.Message {
  getId(): number;
  setId(value: number): UpdateRolePermsRequest;

  getPerms(): PermsUpdate | undefined;
  setPerms(value?: PermsUpdate): UpdateRolePermsRequest;
  hasPerms(): boolean;
  clearPerms(): UpdateRolePermsRequest;

  getAttrs(): AttrsUpdate | undefined;
  setAttrs(value?: AttrsUpdate): UpdateRolePermsRequest;
  hasAttrs(): boolean;
  clearAttrs(): UpdateRolePermsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRolePermsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRolePermsRequest): UpdateRolePermsRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateRolePermsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRolePermsRequest;
  static deserializeBinaryFromReader(message: UpdateRolePermsRequest, reader: jspb.BinaryReader): UpdateRolePermsRequest;
}

export namespace UpdateRolePermsRequest {
  export type AsObject = {
    id: number,
    perms?: PermsUpdate.AsObject,
    attrs?: AttrsUpdate.AsObject,
  }

  export enum PermsCase { 
    _PERMS_NOT_SET = 0,
    PERMS = 2,
  }

  export enum AttrsCase { 
    _ATTRS_NOT_SET = 0,
    ATTRS = 3,
  }
}

export class PermsUpdate extends jspb.Message {
  getToUpdateList(): Array<PermItem>;
  setToUpdateList(value: Array<PermItem>): PermsUpdate;
  clearToUpdateList(): PermsUpdate;
  addToUpdate(value?: PermItem, index?: number): PermItem;

  getToRemoveList(): Array<number>;
  setToRemoveList(value: Array<number>): PermsUpdate;
  clearToRemoveList(): PermsUpdate;
  addToRemove(value: number, index?: number): PermsUpdate;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PermsUpdate.AsObject;
  static toObject(includeInstance: boolean, msg: PermsUpdate): PermsUpdate.AsObject;
  static serializeBinaryToWriter(message: PermsUpdate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PermsUpdate;
  static deserializeBinaryFromReader(message: PermsUpdate, reader: jspb.BinaryReader): PermsUpdate;
}

export namespace PermsUpdate {
  export type AsObject = {
    toUpdateList: Array<PermItem.AsObject>,
    toRemoveList: Array<number>,
  }
}

export class PermItem extends jspb.Message {
  getId(): number;
  setId(value: number): PermItem;

  getVal(): boolean;
  setVal(value: boolean): PermItem;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PermItem.AsObject;
  static toObject(includeInstance: boolean, msg: PermItem): PermItem.AsObject;
  static serializeBinaryToWriter(message: PermItem, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PermItem;
  static deserializeBinaryFromReader(message: PermItem, reader: jspb.BinaryReader): PermItem;
}

export namespace PermItem {
  export type AsObject = {
    id: number,
    val: boolean,
  }
}

export class AttrsUpdate extends jspb.Message {
  getToUpdateList(): Array<resources_permissions_permissions_pb.RoleAttribute>;
  setToUpdateList(value: Array<resources_permissions_permissions_pb.RoleAttribute>): AttrsUpdate;
  clearToUpdateList(): AttrsUpdate;
  addToUpdate(value?: resources_permissions_permissions_pb.RoleAttribute, index?: number): resources_permissions_permissions_pb.RoleAttribute;

  getToRemoveList(): Array<resources_permissions_permissions_pb.RoleAttribute>;
  setToRemoveList(value: Array<resources_permissions_permissions_pb.RoleAttribute>): AttrsUpdate;
  clearToRemoveList(): AttrsUpdate;
  addToRemove(value?: resources_permissions_permissions_pb.RoleAttribute, index?: number): resources_permissions_permissions_pb.RoleAttribute;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AttrsUpdate.AsObject;
  static toObject(includeInstance: boolean, msg: AttrsUpdate): AttrsUpdate.AsObject;
  static serializeBinaryToWriter(message: AttrsUpdate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AttrsUpdate;
  static deserializeBinaryFromReader(message: AttrsUpdate, reader: jspb.BinaryReader): AttrsUpdate;
}

export namespace AttrsUpdate {
  export type AsObject = {
    toUpdateList: Array<resources_permissions_permissions_pb.RoleAttribute.AsObject>,
    toRemoveList: Array<resources_permissions_permissions_pb.RoleAttribute.AsObject>,
  }
}

export class UpdateRolePermsResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRolePermsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRolePermsResponse): UpdateRolePermsResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateRolePermsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRolePermsResponse;
  static deserializeBinaryFromReader(message: UpdateRolePermsResponse, reader: jspb.BinaryReader): UpdateRolePermsResponse;
}

export namespace UpdateRolePermsResponse {
  export type AsObject = {
  }
}

export class GetPermissionsRequest extends jspb.Message {
  getSearch(): string;
  setSearch(value: string): GetPermissionsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPermissionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPermissionsRequest): GetPermissionsRequest.AsObject;
  static serializeBinaryToWriter(message: GetPermissionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPermissionsRequest;
  static deserializeBinaryFromReader(message: GetPermissionsRequest, reader: jspb.BinaryReader): GetPermissionsRequest;
}

export namespace GetPermissionsRequest {
  export type AsObject = {
    search: string,
  }
}

export class GetPermissionsResponse extends jspb.Message {
  getPermissionsList(): Array<resources_permissions_permissions_pb.Permission>;
  setPermissionsList(value: Array<resources_permissions_permissions_pb.Permission>): GetPermissionsResponse;
  clearPermissionsList(): GetPermissionsResponse;
  addPermissions(value?: resources_permissions_permissions_pb.Permission, index?: number): resources_permissions_permissions_pb.Permission;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPermissionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPermissionsResponse): GetPermissionsResponse.AsObject;
  static serializeBinaryToWriter(message: GetPermissionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPermissionsResponse;
  static deserializeBinaryFromReader(message: GetPermissionsResponse, reader: jspb.BinaryReader): GetPermissionsResponse;
}

export namespace GetPermissionsResponse {
  export type AsObject = {
    permissionsList: Array<resources_permissions_permissions_pb.Permission.AsObject>,
  }
}

export class ViewAuditLogRequest extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationRequest | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationRequest): ViewAuditLogRequest;
  hasPagination(): boolean;
  clearPagination(): ViewAuditLogRequest;

  getUserIdsList(): Array<number>;
  setUserIdsList(value: Array<number>): ViewAuditLogRequest;
  clearUserIdsList(): ViewAuditLogRequest;
  addUserIds(value: number, index?: number): ViewAuditLogRequest;

  getFrom(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setFrom(value?: resources_timestamp_timestamp_pb.Timestamp): ViewAuditLogRequest;
  hasFrom(): boolean;
  clearFrom(): ViewAuditLogRequest;

  getTo(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setTo(value?: resources_timestamp_timestamp_pb.Timestamp): ViewAuditLogRequest;
  hasTo(): boolean;
  clearTo(): ViewAuditLogRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ViewAuditLogRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ViewAuditLogRequest): ViewAuditLogRequest.AsObject;
  static serializeBinaryToWriter(message: ViewAuditLogRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ViewAuditLogRequest;
  static deserializeBinaryFromReader(message: ViewAuditLogRequest, reader: jspb.BinaryReader): ViewAuditLogRequest;
}

export namespace ViewAuditLogRequest {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationRequest.AsObject,
    userIdsList: Array<number>,
    from?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    to?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
  }

  export enum FromCase { 
    _FROM_NOT_SET = 0,
    FROM = 3,
  }

  export enum ToCase { 
    _TO_NOT_SET = 0,
    TO = 4,
  }
}

export class ViewAuditLogResponse extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationResponse | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationResponse): ViewAuditLogResponse;
  hasPagination(): boolean;
  clearPagination(): ViewAuditLogResponse;

  getLogsList(): Array<resources_rector_audit_pb.AuditEntry>;
  setLogsList(value: Array<resources_rector_audit_pb.AuditEntry>): ViewAuditLogResponse;
  clearLogsList(): ViewAuditLogResponse;
  addLogs(value?: resources_rector_audit_pb.AuditEntry, index?: number): resources_rector_audit_pb.AuditEntry;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ViewAuditLogResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ViewAuditLogResponse): ViewAuditLogResponse.AsObject;
  static serializeBinaryToWriter(message: ViewAuditLogResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ViewAuditLogResponse;
  static deserializeBinaryFromReader(message: ViewAuditLogResponse, reader: jspb.BinaryReader): ViewAuditLogResponse;
}

export namespace ViewAuditLogResponse {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationResponse.AsObject,
    logsList: Array<resources_rector_audit_pb.AuditEntry.AsObject>,
  }
}

