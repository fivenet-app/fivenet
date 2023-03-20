import * as jspb from 'google-protobuf'

import * as resources_users_users_pb from '../../resources/users/users_pb';


export class Vehicles extends jspb.Message {
  getPlate(): string;
  setPlate(value: string): Vehicles;

  getMode(): string;
  setMode(value: string): Vehicles;

  getType(): string;
  setType(value: string): Vehicles;

  getOwner(): resources_users_users_pb.UserShort | undefined;
  setOwner(value?: resources_users_users_pb.UserShort): Vehicles;
  hasOwner(): boolean;
  clearOwner(): Vehicles;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Vehicles.AsObject;
  static toObject(includeInstance: boolean, msg: Vehicles): Vehicles.AsObject;
  static serializeBinaryToWriter(message: Vehicles, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Vehicles;
  static deserializeBinaryFromReader(message: Vehicles, reader: jspb.BinaryReader): Vehicles;
}

export namespace Vehicles {
  export type AsObject = {
    plate: string,
    mode: string,
    type: string,
    owner?: resources_users_users_pb.UserShort.AsObject,
  }
}

