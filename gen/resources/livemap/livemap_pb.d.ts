import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';


export class GenericMarker extends jspb.Message {
  getX(): number;
  setX(value: number): GenericMarker;

  getY(): number;
  setY(value: number): GenericMarker;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): GenericMarker;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): GenericMarker;

  getId(): number;
  setId(value: number): GenericMarker;

  getName(): string;
  setName(value: string): GenericMarker;

  getIcon(): string;
  setIcon(value: string): GenericMarker;

  getPopup(): string;
  setPopup(value: string): GenericMarker;

  getLink(): string;
  setLink(value: string): GenericMarker;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GenericMarker.AsObject;
  static toObject(includeInstance: boolean, msg: GenericMarker): GenericMarker.AsObject;
  static serializeBinaryToWriter(message: GenericMarker, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GenericMarker;
  static deserializeBinaryFromReader(message: GenericMarker, reader: jspb.BinaryReader): GenericMarker;
}

export namespace GenericMarker {
  export type AsObject = {
    x: number,
    y: number,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    id: number,
    name: string,
    icon: string,
    popup: string,
    link: string,
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

  getId(): number;
  setId(value: number): UserMarker;

  getName(): string;
  setName(value: string): UserMarker;

  getIcon(): string;
  setIcon(value: string): UserMarker;

  getPopup(): string;
  setPopup(value: string): UserMarker;

  getLink(): string;
  setLink(value: string): UserMarker;

  getJob(): string;
  setJob(value: string): UserMarker;

  getJoblabel(): string;
  setJoblabel(value: string): UserMarker;

  getUser(): resources_users_users_pb.UserShort | undefined;
  setUser(value?: resources_users_users_pb.UserShort): UserMarker;
  hasUser(): boolean;
  clearUser(): UserMarker;

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
    id: number,
    name: string,
    icon: string,
    popup: string,
    link: string,
    job: string,
    joblabel: string,
    user?: resources_users_users_pb.UserShort.AsObject,
  }
}

