import * as jspb from 'google-protobuf'



export class StreamRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StreamRequest): StreamRequest.AsObject;
  static serializeBinaryToWriter(message: StreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StreamRequest;
  static deserializeBinaryFromReader(message: StreamRequest, reader: jspb.BinaryReader): StreamRequest;
}

export namespace StreamRequest {
  export type AsObject = {
  }
}

export class ServerStreamResponse extends jspb.Message {
  getUsersList(): Array<Marker>;
  setUsersList(value: Array<Marker>): ServerStreamResponse;
  clearUsersList(): ServerStreamResponse;
  addUsers(value?: Marker, index?: number): Marker;

  getDispatchesList(): Array<Marker>;
  setDispatchesList(value: Array<Marker>): ServerStreamResponse;
  clearDispatchesList(): ServerStreamResponse;
  addDispatches(value?: Marker, index?: number): Marker;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServerStreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ServerStreamResponse): ServerStreamResponse.AsObject;
  static serializeBinaryToWriter(message: ServerStreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServerStreamResponse;
  static deserializeBinaryFromReader(message: ServerStreamResponse, reader: jspb.BinaryReader): ServerStreamResponse;
}

export namespace ServerStreamResponse {
  export type AsObject = {
    usersList: Array<Marker.AsObject>,
    dispatchesList: Array<Marker.AsObject>,
  }
}

export class Marker extends jspb.Message {
  getX(): number;
  setX(value: number): Marker;

  getY(): number;
  setY(value: number): Marker;

  getId(): string;
  setId(value: string): Marker;

  getName(): string;
  setName(value: string): Marker;

  getIcon(): string;
  setIcon(value: string): Marker;

  getPopup(): string;
  setPopup(value: string): Marker;

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
    id: string,
    name: string,
    icon: string,
    popup: string,
  }
}

