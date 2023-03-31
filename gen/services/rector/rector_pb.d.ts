import * as jspb from 'google-protobuf'

import * as resources_common_database_database_pb from '../../resources/common/database/database_pb';
import * as resources_permissions_permissions_pb from '../../resources/permissions/permissions_pb';


export class GetRolesRequest extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationRequest | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationRequest): GetRolesRequest;
  hasPagination(): boolean;
  clearPagination(): GetRolesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRolesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetRolesRequest): GetRolesRequest.AsObject;
  static serializeBinaryToWriter(message: GetRolesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRolesRequest;
  static deserializeBinaryFromReader(message: GetRolesRequest, reader: jspb.BinaryReader): GetRolesRequest;
}

export namespace GetRolesRequest {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationRequest.AsObject,
  }
}

export class GetRolesResponse extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationResponse | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationResponse): GetRolesResponse;
  hasPagination(): boolean;
  clearPagination(): GetRolesResponse;

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
    pagination?: resources_common_database_database_pb.PaginationResponse.AsObject,
    rolesList: Array<resources_permissions_permissions_pb.Role.AsObject>,
  }
}

export class UpdateRoleRequest extends jspb.Message {
  getRole(): resources_permissions_permissions_pb.Role | undefined;
  setRole(value?: resources_permissions_permissions_pb.Role): UpdateRoleRequest;
  hasRole(): boolean;
  clearRole(): UpdateRoleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRoleRequest): UpdateRoleRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRoleRequest;
  static deserializeBinaryFromReader(message: UpdateRoleRequest, reader: jspb.BinaryReader): UpdateRoleRequest;
}

export namespace UpdateRoleRequest {
  export type AsObject = {
    role?: resources_permissions_permissions_pb.Role.AsObject,
  }
}

export class UpdateRoleResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateRoleResponse): UpdateRoleResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateRoleResponse;
  static deserializeBinaryFromReader(message: UpdateRoleResponse, reader: jspb.BinaryReader): UpdateRoleResponse;
}

export namespace UpdateRoleResponse {
  export type AsObject = {
  }
}

export class DeleteRoleRequest extends jspb.Message {
  getRoleId(): number;
  setRoleId(value: number): DeleteRoleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteRoleRequest): DeleteRoleRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteRoleRequest;
  static deserializeBinaryFromReader(message: DeleteRoleRequest, reader: jspb.BinaryReader): DeleteRoleRequest;
}

export namespace DeleteRoleRequest {
  export type AsObject = {
    roleId: number,
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

