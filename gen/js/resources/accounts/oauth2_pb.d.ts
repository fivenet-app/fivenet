import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';


export class OAuth2Account extends jspb.Message {
  getAccountId(): number;
  setAccountId(value: number): OAuth2Account;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): OAuth2Account;
  hasCreatedAt(): boolean;
  clearCreatedAt(): OAuth2Account;

  getProviderName(): string;
  setProviderName(value: string): OAuth2Account;

  getProvider(): OAuth2Provider | undefined;
  setProvider(value?: OAuth2Provider): OAuth2Account;
  hasProvider(): boolean;
  clearProvider(): OAuth2Account;

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
    providerName: string,
    provider?: OAuth2Provider.AsObject,
    externalId: number,
    username: string,
    avatar: string,
  }

  export enum CreatedAtCase { 
    _CREATED_AT_NOT_SET = 0,
    CREATED_AT = 2,
  }
}

export class OAuth2Provider extends jspb.Message {
  getName(): string;
  setName(value: string): OAuth2Provider;

  getLabel(): string;
  setLabel(value: string): OAuth2Provider;

  getHomepage(): string;
  setHomepage(value: string): OAuth2Provider;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OAuth2Provider.AsObject;
  static toObject(includeInstance: boolean, msg: OAuth2Provider): OAuth2Provider.AsObject;
  static serializeBinaryToWriter(message: OAuth2Provider, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OAuth2Provider;
  static deserializeBinaryFromReader(message: OAuth2Provider, reader: jspb.BinaryReader): OAuth2Provider;
}

export namespace OAuth2Provider {
  export type AsObject = {
    name: string,
    label: string,
    homepage: string,
  }
}

