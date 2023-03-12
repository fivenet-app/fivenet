import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';


export class Marker extends jspb.Message {
  getUserid(): number;
  setUserid(value: number): Marker;

  getJob(): string;
  setJob(value: string): Marker;

  getX(): number;
  setX(value: number): Marker;

  getY(): number;
  setY(value: number): Marker;

  getUpdatedat(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedat(value?: resources_timestamp_timestamp_pb.Timestamp): Marker;
  hasUpdatedat(): boolean;
  clearUpdatedat(): Marker;

  getName(): string;
  setName(value: string): Marker;

  getIcon(): string;
  setIcon(value: string): Marker;

  getPopup(): string;
  setPopup(value: string): Marker;

  getLink(): string;
  setLink(value: string): Marker;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Marker.AsObject;
  static toObject(includeInstance: boolean, msg: Marker): Marker.AsObject;
  static serializeBinaryToWriter(message: Marker, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Marker;
  static deserializeBinaryFromReader(message: Marker, reader: jspb.BinaryReader): Marker;
}

export namespace Marker {
  export type AsObject = {
    userid: number,
    job: string,
    x: number,
    y: number,
    updatedat?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    name: string,
    icon: string,
    popup: string,
    link: string,
  }
}

