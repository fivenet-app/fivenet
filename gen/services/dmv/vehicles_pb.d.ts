import * as jspb from 'google-protobuf'



export class FindVehiclesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindVehiclesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FindVehiclesRequest): FindVehiclesRequest.AsObject;
  static serializeBinaryToWriter(message: FindVehiclesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindVehiclesRequest;
  static deserializeBinaryFromReader(message: FindVehiclesRequest, reader: jspb.BinaryReader): FindVehiclesRequest;
}

export namespace FindVehiclesRequest {
  export type AsObject = {
  }
}

export class FindVehiclesResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindVehiclesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: FindVehiclesResponse): FindVehiclesResponse.AsObject;
  static serializeBinaryToWriter(message: FindVehiclesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindVehiclesResponse;
  static deserializeBinaryFromReader(message: FindVehiclesResponse, reader: jspb.BinaryReader): FindVehiclesResponse;
}

export namespace FindVehiclesResponse {
  export type AsObject = {
  }
}

