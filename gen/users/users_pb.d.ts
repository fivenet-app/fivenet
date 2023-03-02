import * as jspb from 'google-protobuf'

import * as common_character_pb from '../common/character_pb';
import * as common_database_pb from '../common/database_pb';


export class FindUsersRequest extends jspb.Message {
  getCurrent(): number;
  setCurrent(value: number): FindUsersRequest;

  getOrderbyList(): Array<common_database_pb.OrderBy>;
  setOrderbyList(value: Array<common_database_pb.OrderBy>): FindUsersRequest;
  clearOrderbyList(): FindUsersRequest;
  addOrderby(value?: common_database_pb.OrderBy, index?: number): common_database_pb.OrderBy;

  getFirstname(): string;
  setFirstname(value: string): FindUsersRequest;

  getLastname(): string;
  setLastname(value: string): FindUsersRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FindUsersRequest): FindUsersRequest.AsObject;
  static serializeBinaryToWriter(message: FindUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindUsersRequest;
  static deserializeBinaryFromReader(message: FindUsersRequest, reader: jspb.BinaryReader): FindUsersRequest;
}

export namespace FindUsersRequest {
  export type AsObject = {
    current: number,
    orderbyList: Array<common_database_pb.OrderBy.AsObject>,
    firstname: string,
    lastname: string,
  }
}

export class FindUsersResponse extends jspb.Message {
  getTotalcount(): number;
  setTotalcount(value: number): FindUsersResponse;

  getCurrent(): number;
  setCurrent(value: number): FindUsersResponse;

  getEnd(): number;
  setEnd(value: number): FindUsersResponse;

  getUsersList(): Array<common_character_pb.Character>;
  setUsersList(value: Array<common_character_pb.Character>): FindUsersResponse;
  clearUsersList(): FindUsersResponse;
  addUsers(value?: common_character_pb.Character, index?: number): common_character_pb.Character;

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
    current: number,
    end: number,
    usersList: Array<common_character_pb.Character.AsObject>,
  }
}

export class GetUserRequest extends jspb.Message {
  getIdentifier(): string;
  setIdentifier(value: string): GetUserRequest;

  getDbid(): string;
  setDbid(value: string): GetUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserRequest): GetUserRequest.AsObject;
  static serializeBinaryToWriter(message: GetUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserRequest;
  static deserializeBinaryFromReader(message: GetUserRequest, reader: jspb.BinaryReader): GetUserRequest;
}

export namespace GetUserRequest {
  export type AsObject = {
    identifier: string,
    dbid: string,
  }
}

export class GetUserResponse extends jspb.Message {
  getUser(): common_character_pb.Character | undefined;
  setUser(value?: common_character_pb.Character): GetUserResponse;
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
    user?: common_character_pb.Character.AsObject,
  }
}

export class SetUserPropsRequest extends jspb.Message {
  getIdentifier(): string;
  setIdentifier(value: string): SetUserPropsRequest;

  getWanted(): boolean;
  setWanted(value: boolean): SetUserPropsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetUserPropsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetUserPropsRequest): SetUserPropsRequest.AsObject;
  static serializeBinaryToWriter(message: SetUserPropsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetUserPropsRequest;
  static deserializeBinaryFromReader(message: SetUserPropsRequest, reader: jspb.BinaryReader): SetUserPropsRequest;
}

export namespace SetUserPropsRequest {
  export type AsObject = {
    identifier: string,
    wanted: boolean,
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

