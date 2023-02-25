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

export class LivemapMarker extends jspb.Message {
  getDispatches(): Marker | undefined;
  setDispatches(value?: Marker): LivemapMarker;
  hasDispatches(): boolean;
  clearDispatches(): LivemapMarker;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LivemapMarker.AsObject;
  static toObject(includeInstance: boolean, msg: LivemapMarker): LivemapMarker.AsObject;
  static serializeBinaryToWriter(message: LivemapMarker, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LivemapMarker;
  static deserializeBinaryFromReader(message: LivemapMarker, reader: jspb.BinaryReader): LivemapMarker;
}

export namespace LivemapMarker {
  export type AsObject = {
    dispatches?: Marker.AsObject,
  }
}

export class Marker extends jspb.Message {
  getX(): number;
  setX(value: number): Marker;

  getY(): number;
  setY(value: number): Marker;

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
    icon: string,
    popup: string,
  }
}

