import * as jspb from 'google-protobuf'



export class EchoRequest extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): EchoRequest;

  getValue(): number;
  setValue(value: number): EchoRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EchoRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EchoRequest): EchoRequest.AsObject;
  static serializeBinaryToWriter(message: EchoRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EchoRequest;
  static deserializeBinaryFromReader(message: EchoRequest, reader: jspb.BinaryReader): EchoRequest;
}

export namespace EchoRequest {
  export type AsObject = {
    message: string,
    value: number,
  }
}

export class EchoResponse extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): EchoResponse;

  getValue(): string;
  setValue(value: string): EchoResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EchoResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EchoResponse): EchoResponse.AsObject;
  static serializeBinaryToWriter(message: EchoResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EchoResponse;
  static deserializeBinaryFromReader(message: EchoResponse, reader: jspb.BinaryReader): EchoResponse;
}

export namespace EchoResponse {
  export type AsObject = {
    message: string,
    value: string,
  }
}

export class ServerStreamingEchoRequest extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): ServerStreamingEchoRequest;

  getMessageCount(): number;
  setMessageCount(value: number): ServerStreamingEchoRequest;

  getMessageInterval(): number;
  setMessageInterval(value: number): ServerStreamingEchoRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServerStreamingEchoRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ServerStreamingEchoRequest): ServerStreamingEchoRequest.AsObject;
  static serializeBinaryToWriter(message: ServerStreamingEchoRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServerStreamingEchoRequest;
  static deserializeBinaryFromReader(message: ServerStreamingEchoRequest, reader: jspb.BinaryReader): ServerStreamingEchoRequest;
}

export namespace ServerStreamingEchoRequest {
  export type AsObject = {
    message: string,
    messageCount: number,
    messageInterval: number,
  }
}

export class ServerStreamingEchoResponse extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): ServerStreamingEchoResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServerStreamingEchoResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ServerStreamingEchoResponse): ServerStreamingEchoResponse.AsObject;
  static serializeBinaryToWriter(message: ServerStreamingEchoResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServerStreamingEchoResponse;
  static deserializeBinaryFromReader(message: ServerStreamingEchoResponse, reader: jspb.BinaryReader): ServerStreamingEchoResponse;
}

export namespace ServerStreamingEchoResponse {
  export type AsObject = {
    message: string,
  }
}

export class EchoStatusRequest extends jspb.Message {
  getStatus(): Status;
  setStatus(value: Status): EchoStatusRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EchoStatusRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EchoStatusRequest): EchoStatusRequest.AsObject;
  static serializeBinaryToWriter(message: EchoStatusRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EchoStatusRequest;
  static deserializeBinaryFromReader(message: EchoStatusRequest, reader: jspb.BinaryReader): EchoStatusRequest;
}

export namespace EchoStatusRequest {
  export type AsObject = {
    status: Status,
  }
}

export class EchoStatusResponse extends jspb.Message {
  getStatus(): EchoStatusResponse.InternalStatus;
  setStatus(value: EchoStatusResponse.InternalStatus): EchoStatusResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EchoStatusResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EchoStatusResponse): EchoStatusResponse.AsObject;
  static serializeBinaryToWriter(message: EchoStatusResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EchoStatusResponse;
  static deserializeBinaryFromReader(message: EchoStatusResponse, reader: jspb.BinaryReader): EchoStatusResponse;
}

export namespace EchoStatusResponse {
  export type AsObject = {
    status: EchoStatusResponse.InternalStatus,
  }

  export enum InternalStatus { 
    UNKNOWN = 0,
    SUCCESS = 1,
  }
}

export enum Status { 
  UNKNOWN = 0,
  SUCCESS = 1,
}
