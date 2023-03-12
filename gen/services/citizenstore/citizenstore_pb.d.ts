import * as jspb from 'google-protobuf'

import * as resources_common_database_database_pb from '../../resources/common/database/database_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';


export class FindUsersRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): FindUsersRequest;

  getOrderbyList(): Array<resources_common_database_database_pb.OrderBy>;
  setOrderbyList(value: Array<resources_common_database_database_pb.OrderBy>): FindUsersRequest;
  clearOrderbyList(): FindUsersRequest;
  addOrderby(value?: resources_common_database_database_pb.OrderBy, index?: number): resources_common_database_database_pb.OrderBy;

  getSearchname(): string;
  setSearchname(value: string): FindUsersRequest;

  getWanted(): boolean;
  setWanted(value: boolean): FindUsersRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FindUsersRequest): FindUsersRequest.AsObject;
  static serializeBinaryToWriter(message: FindUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindUsersRequest;
  static deserializeBinaryFromReader(message: FindUsersRequest, reader: jspb.BinaryReader): FindUsersRequest;
}

export namespace FindUsersRequest {
  export type AsObject = {
    offset: number,
    orderbyList: Array<resources_common_database_database_pb.OrderBy.AsObject>,
    searchname: string,
    wanted: boolean,
  }
}

export class FindUsersResponse extends jspb.Message {
  getTotalcount(): number;
  setTotalcount(value: number): FindUsersResponse;

  getOffset(): number;
  setOffset(value: number): FindUsersResponse;

  getEnd(): number;
  setEnd(value: number): FindUsersResponse;

  getUsersList(): Array<resources_users_users_pb.User>;
  setUsersList(value: Array<resources_users_users_pb.User>): FindUsersResponse;
  clearUsersList(): FindUsersResponse;
  addUsers(value?: resources_users_users_pb.User, index?: number): resources_users_users_pb.User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindUsersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: FindUsersResponse): FindUsersResponse.AsObject;
  static serializeBinaryToWriter(message: FindUsersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindUsersResponse;
  static deserializeBinaryFromReader(message: FindUsersResponse, reader: jspb.BinaryReader): FindUsersResponse;
}

export namespace FindUsersResponse {
  export type AsObject = {
    totalcount: number,
    offset: number,
    end: number,
    usersList: Array<resources_users_users_pb.User.AsObject>,
  }
}

export class GetUserRequest extends jspb.Message {
  getUserid(): number;
  setUserid(value: number): GetUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserRequest): GetUserRequest.AsObject;
  static serializeBinaryToWriter(message: GetUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserRequest;
  static deserializeBinaryFromReader(message: GetUserRequest, reader: jspb.BinaryReader): GetUserRequest;
}

export namespace GetUserRequest {
  export type AsObject = {
    userid: number,
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

export class GetUserActivityRequest extends jspb.Message {
  getUserid(): number;
  setUserid(value: number): GetUserActivityRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserActivityRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserActivityRequest): GetUserActivityRequest.AsObject;
  static serializeBinaryToWriter(message: GetUserActivityRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserActivityRequest;
  static deserializeBinaryFromReader(message: GetUserActivityRequest, reader: jspb.BinaryReader): GetUserActivityRequest;
}

export namespace GetUserActivityRequest {
  export type AsObject = {
    userid: number,
  }
}

export class GetUserActivityResponse extends jspb.Message {
  getActivityList(): Array<resources_users_users_pb.UserActivity>;
  setActivityList(value: Array<resources_users_users_pb.UserActivity>): GetUserActivityResponse;
  clearActivityList(): GetUserActivityResponse;
  addActivity(value?: resources_users_users_pb.UserActivity, index?: number): resources_users_users_pb.UserActivity;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserActivityResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserActivityResponse): GetUserActivityResponse.AsObject;
  static serializeBinaryToWriter(message: GetUserActivityResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserActivityResponse;
  static deserializeBinaryFromReader(message: GetUserActivityResponse, reader: jspb.BinaryReader): GetUserActivityResponse;
}

export namespace GetUserActivityResponse {
  export type AsObject = {
    activityList: Array<resources_users_users_pb.UserActivity.AsObject>,
  }
}

export class SetUserPropsRequest extends jspb.Message {
  getUserid(): number;
  setUserid(value: number): SetUserPropsRequest;

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
    userid: number,
    props?: resources_users_users_pb.UserProps.AsObject,
  }
}

export class SetUserPropsResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetUserPropsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetUserPropsResponse): SetUserPropsResponse.AsObject;
  static serializeBinaryToWriter(message: SetUserPropsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetUserPropsResponse;
  static deserializeBinaryFromReader(message: SetUserPropsResponse, reader: jspb.BinaryReader): SetUserPropsResponse;
}

export namespace SetUserPropsResponse {
  export type AsObject = {
  }
}

