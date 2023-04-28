import * as jspb from 'google-protobuf'

import * as resources_accounts_accounts_pb from '../../resources/accounts/accounts_pb';
import * as resources_jobs_jobs_pb from '../../resources/jobs/jobs_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';


export class CreateAccountRequest extends jspb.Message {
  getRegToken(): string;
  setRegToken(value: string): CreateAccountRequest;

  getUsername(): string;
  setUsername(value: string): CreateAccountRequest;

  getPassword(): string;
  setPassword(value: string): CreateAccountRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAccountRequest): CreateAccountRequest.AsObject;
  static serializeBinaryToWriter(message: CreateAccountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountRequest;
  static deserializeBinaryFromReader(message: CreateAccountRequest, reader: jspb.BinaryReader): CreateAccountRequest;
}

export namespace CreateAccountRequest {
  export type AsObject = {
    regToken: string,
    username: string,
    password: string,
  }
}

export class CreateAccountResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAccountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAccountResponse): CreateAccountResponse.AsObject;
  static serializeBinaryToWriter(message: CreateAccountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAccountResponse;
  static deserializeBinaryFromReader(message: CreateAccountResponse, reader: jspb.BinaryReader): CreateAccountResponse;
}

export namespace CreateAccountResponse {
  export type AsObject = {
  }
}

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

export class ChangePasswordRequest extends jspb.Message {
  getCurrent(): string;
  setCurrent(value: string): ChangePasswordRequest;

  getNew(): string;
  setNew(value: string): ChangePasswordRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangePasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ChangePasswordRequest): ChangePasswordRequest.AsObject;
  static serializeBinaryToWriter(message: ChangePasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangePasswordRequest;
  static deserializeBinaryFromReader(message: ChangePasswordRequest, reader: jspb.BinaryReader): ChangePasswordRequest;
}

export namespace ChangePasswordRequest {
  export type AsObject = {
    current: string,
    pb_new: string,
  }
}

export class ChangePasswordResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): ChangePasswordResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangePasswordResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ChangePasswordResponse): ChangePasswordResponse.AsObject;
  static serializeBinaryToWriter(message: ChangePasswordResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangePasswordResponse;
  static deserializeBinaryFromReader(message: ChangePasswordResponse, reader: jspb.BinaryReader): ChangePasswordResponse;
}

export namespace ChangePasswordResponse {
  export type AsObject = {
    token: string,
  }
}

export class GetAccountInfoRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAccountInfoRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAccountInfoRequest): GetAccountInfoRequest.AsObject;
  static serializeBinaryToWriter(message: GetAccountInfoRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAccountInfoRequest;
  static deserializeBinaryFromReader(message: GetAccountInfoRequest, reader: jspb.BinaryReader): GetAccountInfoRequest;
}

export namespace GetAccountInfoRequest {
  export type AsObject = {
  }
}

export class GetAccountInfoResponse extends jspb.Message {
  getAccount(): resources_accounts_accounts_pb.Account | undefined;
  setAccount(value?: resources_accounts_accounts_pb.Account): GetAccountInfoResponse;
  hasAccount(): boolean;
  clearAccount(): GetAccountInfoResponse;

  getOauth2ConnectionsList(): Array<resources_accounts_accounts_pb.OAuth2Account>;
  setOauth2ConnectionsList(value: Array<resources_accounts_accounts_pb.OAuth2Account>): GetAccountInfoResponse;
  clearOauth2ConnectionsList(): GetAccountInfoResponse;
  addOauth2Connections(value?: resources_accounts_accounts_pb.OAuth2Account, index?: number): resources_accounts_accounts_pb.OAuth2Account;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAccountInfoResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAccountInfoResponse): GetAccountInfoResponse.AsObject;
  static serializeBinaryToWriter(message: GetAccountInfoResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAccountInfoResponse;
  static deserializeBinaryFromReader(message: GetAccountInfoResponse, reader: jspb.BinaryReader): GetAccountInfoResponse;
}

export namespace GetAccountInfoResponse {
  export type AsObject = {
    account?: resources_accounts_accounts_pb.Account.AsObject,
    oauth2ConnectionsList: Array<resources_accounts_accounts_pb.OAuth2Account.AsObject>,
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

  getChar(): resources_users_users_pb.User | undefined;
  setChar(value?: resources_users_users_pb.User): SetJobResponse;
  hasChar(): boolean;
  clearChar(): SetJobResponse;

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
    pb_char?: resources_users_users_pb.User.AsObject,
  }
}

export class OAuth2DisconnectRequest extends jspb.Message {
  getProvider(): string;
  setProvider(value: string): OAuth2DisconnectRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OAuth2DisconnectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: OAuth2DisconnectRequest): OAuth2DisconnectRequest.AsObject;
  static serializeBinaryToWriter(message: OAuth2DisconnectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OAuth2DisconnectRequest;
  static deserializeBinaryFromReader(message: OAuth2DisconnectRequest, reader: jspb.BinaryReader): OAuth2DisconnectRequest;
}

export namespace OAuth2DisconnectRequest {
  export type AsObject = {
    provider: string,
  }
}

export class OAuth2DisconnectResponse extends jspb.Message {
  getSuccess(): boolean;
  setSuccess(value: boolean): OAuth2DisconnectResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OAuth2DisconnectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: OAuth2DisconnectResponse): OAuth2DisconnectResponse.AsObject;
  static serializeBinaryToWriter(message: OAuth2DisconnectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OAuth2DisconnectResponse;
  static deserializeBinaryFromReader(message: OAuth2DisconnectResponse, reader: jspb.BinaryReader): OAuth2DisconnectResponse;
}

export namespace OAuth2DisconnectResponse {
  export type AsObject = {
    success: boolean,
  }
}

