import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';


export class AuditEntry extends jspb.Message {
  getId(): number;
  setId(value: number): AuditEntry;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): AuditEntry;
  hasCreatedAt(): boolean;
  clearCreatedAt(): AuditEntry;

  getUserId(): number;
  setUserId(value: number): AuditEntry;

  getUser(): resources_users_users_pb.UserShort | undefined;
  setUser(value?: resources_users_users_pb.UserShort): AuditEntry;
  hasUser(): boolean;
  clearUser(): AuditEntry;

  getUserJob(): string;
  setUserJob(value: string): AuditEntry;

  getTargetJob(): string;
  setTargetJob(value: string): AuditEntry;
  hasTargetJob(): boolean;
  clearTargetJob(): AuditEntry;

  getService(): string;
  setService(value: string): AuditEntry;

  getMethod(): string;
  setMethod(value: string): AuditEntry;

  getState(): EVENT_TYPE;
  setState(value: EVENT_TYPE): AuditEntry;

  getData(): string;
  setData(value: string): AuditEntry;
  hasData(): boolean;
  clearData(): AuditEntry;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuditEntry.AsObject;
  static toObject(includeInstance: boolean, msg: AuditEntry): AuditEntry.AsObject;
  static serializeBinaryToWriter(message: AuditEntry, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuditEntry;
  static deserializeBinaryFromReader(message: AuditEntry, reader: jspb.BinaryReader): AuditEntry;
}

export namespace AuditEntry {
  export type AsObject = {
    id: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    userId: number,
    user?: resources_users_users_pb.UserShort.AsObject,
    userJob: string,
    targetJob?: string,
    service: string,
    method: string,
    state: EVENT_TYPE,
    data?: string,
  }

  export enum UserCase { 
    _USER_NOT_SET = 0,
    USER = 4,
  }

  export enum TargetJobCase { 
    _TARGET_JOB_NOT_SET = 0,
    TARGET_JOB = 6,
  }

  export enum DataCase { 
    _DATA_NOT_SET = 0,
    DATA = 10,
  }
}

export enum EVENT_TYPE { 
  UNKNOWN = 0,
  CREATE = 1,
  UPDATE = 2,
  DELETE = 3,
}
