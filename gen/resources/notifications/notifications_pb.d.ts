import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';


export class Notification extends jspb.Message {
  getId(): number;
  setId(value: number): Notification;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Notification;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Notification;

  getReadAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setReadAt(value?: resources_timestamp_timestamp_pb.Timestamp): Notification;
  hasReadAt(): boolean;
  clearReadAt(): Notification;

  getUserId(): number;
  setUserId(value: number): Notification;

  getTitle(): string;
  setTitle(value: string): Notification;

  getType(): string;
  setType(value: string): Notification;

  getContent(): string;
  setContent(value: string): Notification;

  getData(): string;
  setData(value: string): Notification;
  hasData(): boolean;
  clearData(): Notification;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Notification.AsObject;
  static toObject(includeInstance: boolean, msg: Notification): Notification.AsObject;
  static serializeBinaryToWriter(message: Notification, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Notification;
  static deserializeBinaryFromReader(message: Notification, reader: jspb.BinaryReader): Notification;
}

export namespace Notification {
  export type AsObject = {
    id: number,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    readAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    userId: number,
    title: string,
    type: string,
    content: string,
    data?: string,
  }

  export enum DataCase { 
    _DATA_NOT_SET = 0,
    DATA = 8,
  }
}

