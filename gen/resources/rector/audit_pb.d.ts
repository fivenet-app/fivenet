import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';


export class AuditEntry extends jspb.Message {
  getId(): number;
  setId(value: number): AuditEntry;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): AuditEntry;
  hasCreatedAt(): boolean;
  clearCreatedAt(): AuditEntry;

  getUserId(): number;
  setUserId(value: number): AuditEntry;

  getService(): string;
  setService(value: string): AuditEntry;

  getMethod(): string;
  setMethod(value: string): AuditEntry;

  getState(): EVENT_TYPE;
  setState(value: EVENT_TYPE): AuditEntry;

  getData(): string;
  setData(value: string): AuditEntry;

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
    service: string,
    method: string,
    state: EVENT_TYPE,
    data: string,
  }
}

export enum EVENT_TYPE { 
  UNKNOWN = 0,
  CREATE = 1,
  UPDATE = 2,
  DELETE = 3,
}
