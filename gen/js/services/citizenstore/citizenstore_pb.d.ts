import * as jspb from 'google-protobuf'

import * as resources_common_database_database_pb from '../../resources/common/database/database_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';


export class ListCitizensRequest extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationRequest | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationRequest): ListCitizensRequest;
  hasPagination(): boolean;
  clearPagination(): ListCitizensRequest;

  getSearchName(): string;
  setSearchName(value: string): ListCitizensRequest;

  getWanted(): boolean;
  setWanted(value: boolean): ListCitizensRequest;

  getPhoneNumber(): string;
  setPhoneNumber(value: string): ListCitizensRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListCitizensRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListCitizensRequest): ListCitizensRequest.AsObject;
  static serializeBinaryToWriter(message: ListCitizensRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListCitizensRequest;
  static deserializeBinaryFromReader(message: ListCitizensRequest, reader: jspb.BinaryReader): ListCitizensRequest;
}

export namespace ListCitizensRequest {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationRequest.AsObject,
    searchName: string,
    wanted: boolean,
    phoneNumber: string,
  }
}

export class ListCitizensResponse extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationResponse | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationResponse): ListCitizensResponse;
  hasPagination(): boolean;
  clearPagination(): ListCitizensResponse;

  getUsersList(): Array<resources_users_users_pb.User>;
  setUsersList(value: Array<resources_users_users_pb.User>): ListCitizensResponse;
  clearUsersList(): ListCitizensResponse;
  addUsers(value?: resources_users_users_pb.User, index?: number): resources_users_users_pb.User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListCitizensResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListCitizensResponse): ListCitizensResponse.AsObject;
  static serializeBinaryToWriter(message: ListCitizensResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListCitizensResponse;
  static deserializeBinaryFromReader(message: ListCitizensResponse, reader: jspb.BinaryReader): ListCitizensResponse;
}

export namespace ListCitizensResponse {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationResponse.AsObject,
    usersList: Array<resources_users_users_pb.User.AsObject>,
  }
}

export class GetUserRequest extends jspb.Message {
  getUserId(): number;
  setUserId(value: number): GetUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserRequest): GetUserRequest.AsObject;
  static serializeBinaryToWriter(message: GetUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserRequest;
  static deserializeBinaryFromReader(message: GetUserRequest, reader: jspb.BinaryReader): GetUserRequest;
}

export namespace GetUserRequest {
  export type AsObject = {
    userId: number,
  }
}

export class GetUserResponse extends jspb.Message {
  getUser(): resources_users_users_pb.User | undefined;
  setUser(value?: resources_users_users_pb.User): GetUserResponse;
  hasUser(): boolean;
  clearUser(): GetUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserResponse): GetUserResponse.AsObject;
  static serializeBinaryToWriter(message: GetUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserResponse;
  static deserializeBinaryFromReader(message: GetUserResponse, reader: jspb.BinaryReader): GetUserResponse;
}

export namespace GetUserResponse {
  export type AsObject = {
    user?: resources_users_users_pb.User.AsObject,
  }
}

export class ListUserActivityRequest extends jspb.Message {
  getUserId(): number;
  setUserId(value: number): ListUserActivityRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserActivityRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserActivityRequest): ListUserActivityRequest.AsObject;
  static serializeBinaryToWriter(message: ListUserActivityRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserActivityRequest;
  static deserializeBinaryFromReader(message: ListUserActivityRequest, reader: jspb.BinaryReader): ListUserActivityRequest;
}

export namespace ListUserActivityRequest {
  export type AsObject = {
    userId: number,
  }
}

export class ListUserActivityResponse extends jspb.Message {
  getActivityList(): Array<resources_users_users_pb.UserActivity>;
  setActivityList(value: Array<resources_users_users_pb.UserActivity>): ListUserActivityResponse;
  clearActivityList(): ListUserActivityResponse;
  addActivity(value?: resources_users_users_pb.UserActivity, index?: number): resources_users_users_pb.UserActivity;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserActivityResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserActivityResponse): ListUserActivityResponse.AsObject;
  static serializeBinaryToWriter(message: ListUserActivityResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserActivityResponse;
  static deserializeBinaryFromReader(message: ListUserActivityResponse, reader: jspb.BinaryReader): ListUserActivityResponse;
}

export namespace ListUserActivityResponse {
  export type AsObject = {
    activityList: Array<resources_users_users_pb.UserActivity.AsObject>,
  }
}

export class SetUserPropsRequest extends jspb.Message {
  getProps(): resources_users_users_pb.UserProps | undefined;
  setProps(value?: resources_users_users_pb.UserProps): SetUserPropsRequest;
  hasProps(): boolean;
  clearProps(): SetUserPropsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetUserPropsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetUserPropsRequest): SetUserPropsRequest.AsObject;
  static serializeBinaryToWriter(message: SetUserPropsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetUserPropsRequest;
  static deserializeBinaryFromReader(message: SetUserPropsRequest, reader: jspb.BinaryReader): SetUserPropsRequest;
}

export namespace SetUserPropsRequest {
  export type AsObject = {
    props?: resources_users_users_pb.UserProps.AsObject,
  }
}

export class SetUserPropsResponse extends jspb.Message {
  getProps(): resources_users_users_pb.UserProps | undefined;
  setProps(value?: resources_users_users_pb.UserProps): SetUserPropsResponse;
  hasProps(): boolean;
  clearProps(): SetUserPropsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetUserPropsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetUserPropsResponse): SetUserPropsResponse.AsObject;
  static serializeBinaryToWriter(message: SetUserPropsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetUserPropsResponse;
  static deserializeBinaryFromReader(message: SetUserPropsResponse, reader: jspb.BinaryReader): SetUserPropsResponse;
}

export namespace SetUserPropsResponse {
  export type AsObject = {
    props?: resources_users_users_pb.UserProps.AsObject,
  }
}

