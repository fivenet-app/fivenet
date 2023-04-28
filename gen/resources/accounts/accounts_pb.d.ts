import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';


export class Account extends jspb.Message {
  getId(): number;
  setId(value: number): Account;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Account;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Account;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Account;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): Account;

  getUsername(): string;
  setUsername(value: string): Account;

  getLicense(): string;
  setLicense(value: string): Account;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Account.AsObject;
  static toObject(includeInstance: boolean, msg: Account): Account.AsObject;
  static serializeBinaryToWriter(message: Account, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Account;
  static deserializeBinaryFromReader(message: Account, reader: jspb.BinaryReader): Account;
}

export namespace Account {
  export type AsObject = {
    id: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    username: string,
    license: string,
  }

  export enum CreatedAtCase { 
    _CREATED_AT_NOT_SET = 0,
    CREATED_AT = 2,
  }

  export enum UpdatedAtCase { 
    _UPDATED_AT_NOT_SET = 0,
    UPDATED_AT = 3,
  }
}

export class OAuth2Account extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): OAuth2Account;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): OAuth2Account;
  hasCreatedAt(): boolean;
  clearCreatedAt(): OAuth2Account;

  getProvider(): string;
  setProvider(value: string): OAuth2Account;

  getExternalId(): number;
  setExternalId(value: number): OAuth2Account;

  getUsername(): string;
  setUsername(value: string): OAuth2Account;

  getAvatar(): string;
  setAvatar(value: string): OAuth2Account;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OAuth2Account.AsObject;
  static toObject(includeInstance: boolean, msg: OAuth2Account): OAuth2Account.AsObject;
  static serializeBinaryToWriter(message: OAuth2Account, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OAuth2Account;
  static deserializeBinaryFromReader(message: OAuth2Account, reader: jspb.BinaryReader): OAuth2Account;
}

export namespace OAuth2Account {
  export type AsObject = {
    accountId: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    provider: string,
    externalId: number,
    username: string,
    avatar: string,
  }

  export enum CreatedAtCase { 
    _CREATED_AT_NOT_SET = 0,
    CREATED_AT = 2,
  }
}

