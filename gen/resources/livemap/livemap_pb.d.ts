import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';


export class DispatchMarker extends jspb.Message {
  getX(): number;
  setX(value: number): DispatchMarker;

  getY(): number;
  setY(value: number): DispatchMarker;

  getUpdatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): DispatchMarker;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): DispatchMarker;

  getId(): number;
  setId(value: number): DispatchMarker;

  getJob(): string;
  setJob(value: string): DispatchMarker;

  getJobLabel(): string;
  setJobLabel(value: string): DispatchMarker;

  getName(): string;
  setName(value: string): DispatchMarker;

  getIcon(): string;
  setIcon(value: string): DispatchMarker;

  getIconColor(): string;
  setIconColor(value: string): DispatchMarker;

  getPopup(): string;
  setPopup(value: string): DispatchMarker;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DispatchMarker.AsObject;
  static toObject(includeInstance: boolean, msg: DispatchMarker): DispatchMarker.AsObject;
  static serializeBinaryToWriter(message: DispatchMarker, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DispatchMarker;
  static deserializeBinaryFromReader(message: DispatchMarker, reader: jspb.BinaryReader): DispatchMarker;
}

export namespace DispatchMarker {
  export type AsObject = {
    x: number,
    y: number,
    updatedAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    id: number,
    job: string,
    jobLabel: string,
    name: string,
    icon: string,
    iconColor: string,
    popup: string,
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

  getIconColor(): string;
  setIconColor(value: string): UserMarker;

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
    iconColor: string,
    user?: resources_users_users_pb.UserShort.AsObject,
  }
}

