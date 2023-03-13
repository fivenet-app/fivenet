import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';


export class Marker extends jspb.Message {
  getX(): number;
  setX(value: number): Marker;

  getY(): number;
  setY(value: number): Marker;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): Marker;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): Marker;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Marker.AsObject;
  static toObject(includeInstance: boolean, msg: Marker): Marker.AsObject;
  static serializeBinaryToWriter(message: Marker, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Marker;
  static deserializeBinaryFromReader(message: Marker, reader: jspb.BinaryReader): Marker;
}

export namespace Marker {
  export type AsObject = {
    x: number,
    y: number,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
  }
}

export class UserMarker extends jspb.Message {
  getX(): number;
  setX(value: number): UserMarker;

  getY(): number;
  setY(value: number): UserMarker;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): UserMarker;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): UserMarker;

  getJob(): string;
  setJob(value: string): UserMarker;

  getUserId(): number;
  setUserId(value: number): UserMarker;

  getName(): string;
  setName(value: string): UserMarker;

  getIcon(): string;
  setIcon(value: string): UserMarker;

  getPopup(): string;
  setPopup(value: string): UserMarker;

  getLink(): string;
  setLink(value: string): UserMarker;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserMarker.AsObject;
  static toObject(includeInstance: boolean, msg: UserMarker): UserMarker.AsObject;
  static serializeBinaryToWriter(message: UserMarker, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserMarker;
  static deserializeBinaryFromReader(message: UserMarker, reader: jspb.BinaryReader): UserMarker;
}

export namespace UserMarker {
  export type AsObject = {
    x: number,
    y: number,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    job: string,
    userId: number,
    name: string,
    icon: string,
    popup: string,
    link: string,
  }
}

