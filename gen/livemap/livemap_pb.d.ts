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

export class StreamResponse extends jspb.Message {
  getDispatchesList(): Array<Marker>;
  setDispatchesList(value: Array<Marker>): StreamResponse;
  clearDispatchesList(): StreamResponse;
  addDispatches(value?: Marker, index?: number): Marker;

  getUsersList(): Array<Marker>;
  setUsersList(value: Array<Marker>): StreamResponse;
  clearUsersList(): StreamResponse;
  addUsers(value?: Marker, index?: number): Marker;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StreamResponse): StreamResponse.AsObject;
  static serializeBinaryToWriter(message: StreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StreamResponse;
  static deserializeBinaryFromReader(message: StreamResponse, reader: jspb.BinaryReader): StreamResponse;
}

export namespace StreamResponse {
  export type AsObject = {
    dispatchesList: Array<Marker.AsObject>,
    usersList: Array<Marker.AsObject>,
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

