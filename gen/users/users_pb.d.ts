import * as jspb from 'google-protobuf'

import * as common_database_database_pb from '../common/database/database_pb';
import * as common_timestamp_timestamp_pb from '../common/timestamp/timestamp_pb';
import * as common_userinfo_userinfo_pb from '../common/userinfo/userinfo_pb';


export class FindUsersRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): FindUsersRequest;

  getOrderbyList(): Array<common_database_database_pb.OrderBy>;
  setOrderbyList(value: Array<common_database_database_pb.OrderBy>): FindUsersRequest;
  clearOrderbyList(): FindUsersRequest;
  addOrderby(value?: common_database_database_pb.OrderBy, index?: number): common_database_database_pb.OrderBy;

  getFirstname(): string;
  setFirstname(value: string): FindUsersRequest;

  getLastname(): string;
  setLastname(value: string): FindUsersRequest;

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
    orderbyList: Array<common_database_database_pb.OrderBy.AsObject>,
    firstname: string,
    lastname: string,
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

  getUsersList(): Array<common_userinfo_userinfo_pb.User>;
  setUsersList(value: Array<common_userinfo_userinfo_pb.User>): FindUsersResponse;
  clearUsersList(): FindUsersResponse;
  addUsers(value?: common_userinfo_userinfo_pb.User, index?: number): common_userinfo_userinfo_pb.User;

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
    usersList: Array<common_userinfo_userinfo_pb.User.AsObject>,
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
  getUser(): common_userinfo_userinfo_pb.User | undefined;
  setUser(value?: common_userinfo_userinfo_pb.User): GetUserResponse;
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
    user?: common_userinfo_userinfo_pb.User.AsObject,
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
  getActivityList(): Array<UserActivity>;
  setActivityList(value: Array<UserActivity>): GetUserActivityResponse;
  clearActivityList(): GetUserActivityResponse;
  addActivity(value?: UserActivity, index?: number): UserActivity;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserActivityResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserActivityResponse): GetUserActivityResponse.AsObject;
  static serializeBinaryToWriter(message: GetUserActivityResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserActivityResponse;
  static deserializeBinaryFromReader(message: GetUserActivityResponse, reader: jspb.BinaryReader): GetUserActivityResponse;
}

export namespace GetUserActivityResponse {
  export type AsObject = {
    activityList: Array<UserActivity.AsObject>,
  }
}

export class SetUserPropsRequest extends jspb.Message {
  getUserid(): number;
  setUserid(value: number): SetUserPropsRequest;

  getWanted(): boolean;
  setWanted(value: boolean): SetUserPropsRequest;
  hasWanted(): boolean;
  clearWanted(): SetUserPropsRequest;

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
    wanted?: boolean,
  }

  export enum WantedCase { 
    _WANTED_NOT_SET = 0,
    WANTED = 2,
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

export class UserActivity extends jspb.Message {
  getId(): number;
  setId(value: number): UserActivity;

  getType(): string;
  setType(value: string): UserActivity;

  getCreatedat(): common_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedat(value?: common_timestamp_timestamp_pb.Timestamp): UserActivity;
  hasCreatedat(): boolean;
  clearCreatedat(): UserActivity;

  getTargetuser(): common_userinfo_userinfo_pb.ShortUser | undefined;
  setTargetuser(value?: common_userinfo_userinfo_pb.ShortUser): UserActivity;
  hasTargetuser(): boolean;
  clearTargetuser(): UserActivity;

  getCauseuser(): common_userinfo_userinfo_pb.ShortUser | undefined;
  setCauseuser(value?: common_userinfo_userinfo_pb.ShortUser): UserActivity;
  hasCauseuser(): boolean;
  clearCauseuser(): UserActivity;

  getKey(): string;
  setKey(value: string): UserActivity;

  getOldvalue(): string;
  setOldvalue(value: string): UserActivity;

  getNewvalue(): string;
  setNewvalue(value: string): UserActivity;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserActivity.AsObject;
  static toObject(includeInstance: boolean, msg: UserActivity): UserActivity.AsObject;
  static serializeBinaryToWriter(message: UserActivity, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserActivity;
  static deserializeBinaryFromReader(message: UserActivity, reader: jspb.BinaryReader): UserActivity;
}

export namespace UserActivity {
  export type AsObject = {
    id: number,
    type: string,
    createdat?: common_timestamp_timestamp_pb.Timestamp.AsObject,
    targetuser?: common_userinfo_userinfo_pb.ShortUser.AsObject,
    causeuser?: common_userinfo_userinfo_pb.ShortUser.AsObject,
    key: string,
    oldvalue: string,
    newvalue: string,
  }
}

