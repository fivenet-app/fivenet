import * as jspb from 'google-protobuf'

import * as resources_livemap_livemap_pb from '../../resources/livemap/livemap_pb';


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
  getUsersList(): Array<resources_livemap_livemap_pb.UserMarker>;
  setUsersList(value: Array<resources_livemap_livemap_pb.UserMarker>): ServerStreamResponse;
  clearUsersList(): ServerStreamResponse;
  addUsers(value?: resources_livemap_livemap_pb.UserMarker, index?: number): resources_livemap_livemap_pb.UserMarker;

  getDispatchesList(): Array<resources_livemap_livemap_pb.UserMarker>;
  setDispatchesList(value: Array<resources_livemap_livemap_pb.UserMarker>): ServerStreamResponse;
  clearDispatchesList(): ServerStreamResponse;
  addDispatches(value?: resources_livemap_livemap_pb.UserMarker, index?: number): resources_livemap_livemap_pb.UserMarker;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServerStreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ServerStreamResponse): ServerStreamResponse.AsObject;
  static serializeBinaryToWriter(message: ServerStreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServerStreamResponse;
  static deserializeBinaryFromReader(message: ServerStreamResponse, reader: jspb.BinaryReader): ServerStreamResponse;
}

export namespace ServerStreamResponse {
  export type AsObject = {
    usersList: Array<resources_livemap_livemap_pb.UserMarker.AsObject>,
    dispatchesList: Array<resources_livemap_livemap_pb.UserMarker.AsObject>,
  }
}

