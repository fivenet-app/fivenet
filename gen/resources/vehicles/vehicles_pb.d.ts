import * as jspb from 'google-protobuf'

import * as resources_users_users_pb from '../../resources/users/users_pb';


export class Vehicle extends jspb.Message {
  getPlate(): string;
  setPlate(value: string): Vehicle;

  getModel(): string;
  setModel(value: string): Vehicle;

  getType(): string;
  setType(value: string): Vehicle;

  getOwner(): resources_users_users_pb.UserShortNI | undefined;
  setOwner(value?: resources_users_users_pb.UserShortNI): Vehicle;
  hasOwner(): boolean;
  clearOwner(): Vehicle;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Vehicle.AsObject;
  static toObject(includeInstance: boolean, msg: Vehicle): Vehicle.AsObject;
  static serializeBinaryToWriter(message: Vehicle, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Vehicle;
  static deserializeBinaryFromReader(message: Vehicle, reader: jspb.BinaryReader): Vehicle;
}

export namespace Vehicle {
  export type AsObject = {
    plate: string,
    model: string,
    type: string,
    owner?: resources_users_users_pb.UserShortNI.AsObject,
  }
}

