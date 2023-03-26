import * as jspb from 'google-protobuf'

import * as resources_jobs_jobs_pb from '../../resources/jobs/jobs_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';


export class LoginRequest extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): LoginRequest;

  getPassword(): string;
  setPassword(value: string): LoginRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LoginRequest): LoginRequest.AsObject;
  static serializeBinaryToWriter(message: LoginRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginRequest;
  static deserializeBinaryFromReader(message: LoginRequest, reader: jspb.BinaryReader): LoginRequest;
}

export namespace LoginRequest {
  export type AsObject = {
    username: string,
    password: string,
  }
}

export class LoginResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): LoginResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginResponse.AsObject;
  static toObject(includeInstance: boolean, msg: LoginResponse): LoginResponse.AsObject;
  static serializeBinaryToWriter(message: LoginResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginResponse;
  static deserializeBinaryFromReader(message: LoginResponse, reader: jspb.BinaryReader): LoginResponse;
}

export namespace LoginResponse {
  export type AsObject = {
    token: string,
  }
}

export class GetCharactersRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCharactersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCharactersRequest): GetCharactersRequest.AsObject;
  static serializeBinaryToWriter(message: GetCharactersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCharactersRequest;
  static deserializeBinaryFromReader(message: GetCharactersRequest, reader: jspb.BinaryReader): GetCharactersRequest;
}

export namespace GetCharactersRequest {
  export type AsObject = {
  }
}

export class GetCharactersResponse extends jspb.Message {
  getCharsList(): Array<resources_users_users_pb.User>;
  setCharsList(value: Array<resources_users_users_pb.User>): GetCharactersResponse;
  clearCharsList(): GetCharactersResponse;
  addChars(value?: resources_users_users_pb.User, index?: number): resources_users_users_pb.User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCharactersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCharactersResponse): GetCharactersResponse.AsObject;
  static serializeBinaryToWriter(message: GetCharactersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCharactersResponse;
  static deserializeBinaryFromReader(message: GetCharactersResponse, reader: jspb.BinaryReader): GetCharactersResponse;
}

export namespace GetCharactersResponse {
  export type AsObject = {
    charsList: Array<resources_users_users_pb.User.AsObject>,
  }
}

export class ChooseCharacterRequest extends jspb.Message {
  getCharId(): number;
  setCharId(value: number): ChooseCharacterRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChooseCharacterRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ChooseCharacterRequest): ChooseCharacterRequest.AsObject;
  static serializeBinaryToWriter(message: ChooseCharacterRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChooseCharacterRequest;
  static deserializeBinaryFromReader(message: ChooseCharacterRequest, reader: jspb.BinaryReader): ChooseCharacterRequest;
}

export namespace ChooseCharacterRequest {
  export type AsObject = {
    charId: number,
  }
}

export class ChooseCharacterResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): ChooseCharacterResponse;

  getPermissionsList(): Array<string>;
  setPermissionsList(value: Array<string>): ChooseCharacterResponse;
  clearPermissionsList(): ChooseCharacterResponse;
  addPermissions(value: string, index?: number): ChooseCharacterResponse;

  getJobProps(): resources_jobs_jobs_pb.JobProps | undefined;
  setJobProps(value?: resources_jobs_jobs_pb.JobProps): ChooseCharacterResponse;
  hasJobProps(): boolean;
  clearJobProps(): ChooseCharacterResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChooseCharacterResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ChooseCharacterResponse): ChooseCharacterResponse.AsObject;
  static serializeBinaryToWriter(message: ChooseCharacterResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChooseCharacterResponse;
  static deserializeBinaryFromReader(message: ChooseCharacterResponse, reader: jspb.BinaryReader): ChooseCharacterResponse;
}

export namespace ChooseCharacterResponse {
  export type AsObject = {
    token: string,
    permissionsList: Array<string>,
    jobProps?: resources_jobs_jobs_pb.JobProps.AsObject,
  }
}

export class LogoutRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogoutRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LogoutRequest): LogoutRequest.AsObject;
  static serializeBinaryToWriter(message: LogoutRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogoutRequest;
  static deserializeBinaryFromReader(message: LogoutRequest, reader: jspb.BinaryReader): LogoutRequest;
}

export namespace LogoutRequest {
  export type AsObject = {
  }
}

export class LogoutResponse extends jspb.Message {
  getSuccess(): boolean;
  setSuccess(value: boolean): LogoutResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogoutResponse.AsObject;
  static toObject(includeInstance: boolean, msg: LogoutResponse): LogoutResponse.AsObject;
  static serializeBinaryToWriter(message: LogoutResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogoutResponse;
  static deserializeBinaryFromReader(message: LogoutResponse, reader: jspb.BinaryReader): LogoutResponse;
}

export namespace LogoutResponse {
  export type AsObject = {
    success: boolean,
  }
}

export class SetJobRequest extends jspb.Message {
  getCharId(): number;
  setCharId(value: number): SetJobRequest;

  getJob(): string;
  setJob(value: string): SetJobRequest;

  getJobGrade(): number;
  setJobGrade(value: number): SetJobRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetJobRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetJobRequest): SetJobRequest.AsObject;
  static serializeBinaryToWriter(message: SetJobRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetJobRequest;
  static deserializeBinaryFromReader(message: SetJobRequest, reader: jspb.BinaryReader): SetJobRequest;
}

export namespace SetJobRequest {
  export type AsObject = {
    charId: number,
    job: string,
    jobGrade: number,
  }
}

export class SetJobResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): SetJobResponse;

  getJobProps(): resources_jobs_jobs_pb.JobProps | undefined;
  setJobProps(value?: resources_jobs_jobs_pb.JobProps): SetJobResponse;
  hasJobProps(): boolean;
  clearJobProps(): SetJobResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetJobResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetJobResponse): SetJobResponse.AsObject;
  static serializeBinaryToWriter(message: SetJobResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetJobResponse;
  static deserializeBinaryFromReader(message: SetJobResponse, reader: jspb.BinaryReader): SetJobResponse;
}

export namespace SetJobResponse {
  export type AsObject = {
    token: string,
    jobProps?: resources_jobs_jobs_pb.JobProps.AsObject,
  }
}

