import * as jspb from 'google-protobuf'

import * as resources_common_database_database_pb from '../../resources/common/database/database_pb';
import * as resources_notifications_notifications_pb from '../../resources/notifications/notifications_pb';


export class StreamRequest extends jspb.Message {
  getLastId(): number;
  setLastId(value: number): StreamRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StreamRequest): StreamRequest.AsObject;
  static serializeBinaryToWriter(message: StreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StreamRequest;
  static deserializeBinaryFromReader(message: StreamRequest, reader: jspb.BinaryReader): StreamRequest;
}

export namespace StreamRequest {
  export type AsObject = {
    lastId: number,
  }
}

export class StreamResponse extends jspb.Message {
  getLastId(): number;
  setLastId(value: number): StreamResponse;

  getNotificationsList(): Array<resources_notifications_notifications_pb.Notification>;
  setNotificationsList(value: Array<resources_notifications_notifications_pb.Notification>): StreamResponse;
  clearNotificationsList(): StreamResponse;
  addNotifications(value?: resources_notifications_notifications_pb.Notification, index?: number): resources_notifications_notifications_pb.Notification;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StreamResponse): StreamResponse.AsObject;
  static serializeBinaryToWriter(message: StreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StreamResponse;
  static deserializeBinaryFromReader(message: StreamResponse, reader: jspb.BinaryReader): StreamResponse;
}

export namespace StreamResponse {
  export type AsObject = {
    lastId: number,
    notificationsList: Array<resources_notifications_notifications_pb.Notification.AsObject>,
  }
}

export class GetNotificationsRequest extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationRequest | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationRequest): GetNotificationsRequest;
  hasPagination(): boolean;
  clearPagination(): GetNotificationsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetNotificationsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetNotificationsRequest): GetNotificationsRequest.AsObject;
  static serializeBinaryToWriter(message: GetNotificationsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetNotificationsRequest;
  static deserializeBinaryFromReader(message: GetNotificationsRequest, reader: jspb.BinaryReader): GetNotificationsRequest;
}

export namespace GetNotificationsRequest {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationRequest.AsObject,
  }
}

export class GetNotificationsResponse extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationResponse | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationResponse): GetNotificationsResponse;
  hasPagination(): boolean;
  clearPagination(): GetNotificationsResponse;

  getNotificationsList(): Array<resources_notifications_notifications_pb.Notification>;
  setNotificationsList(value: Array<resources_notifications_notifications_pb.Notification>): GetNotificationsResponse;
  clearNotificationsList(): GetNotificationsResponse;
  addNotifications(value?: resources_notifications_notifications_pb.Notification, index?: number): resources_notifications_notifications_pb.Notification;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetNotificationsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetNotificationsResponse): GetNotificationsResponse.AsObject;
  static serializeBinaryToWriter(message: GetNotificationsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetNotificationsResponse;
  static deserializeBinaryFromReader(message: GetNotificationsResponse, reader: jspb.BinaryReader): GetNotificationsResponse;
}

export namespace GetNotificationsResponse {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationResponse.AsObject,
    notificationsList: Array<resources_notifications_notifications_pb.Notification.AsObject>,
  }
}

export class ReadNotificationsRequest extends jspb.Message {
  getIdsList(): Array<number>;
  setIdsList(value: Array<number>): ReadNotificationsRequest;
  clearIdsList(): ReadNotificationsRequest;
  addIds(value: number, index?: number): ReadNotificationsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReadNotificationsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReadNotificationsRequest): ReadNotificationsRequest.AsObject;
  static serializeBinaryToWriter(message: ReadNotificationsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReadNotificationsRequest;
  static deserializeBinaryFromReader(message: ReadNotificationsRequest, reader: jspb.BinaryReader): ReadNotificationsRequest;
}

export namespace ReadNotificationsRequest {
  export type AsObject = {
    idsList: Array<number>,
  }
}

export class ReadNotificationsResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReadNotificationsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReadNotificationsResponse): ReadNotificationsResponse.AsObject;
  static serializeBinaryToWriter(message: ReadNotificationsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReadNotificationsResponse;
  static deserializeBinaryFromReader(message: ReadNotificationsResponse, reader: jspb.BinaryReader): ReadNotificationsResponse;
}

export namespace ReadNotificationsResponse {
  export type AsObject = {
  }
}

