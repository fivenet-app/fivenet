import * as jspb from 'google-protobuf'

import * as resources_common_database_database_pb from '../../resources/common/database/database_pb';
import * as resources_vehicles_vehicles_pb from '../../resources/vehicles/vehicles_pb';


export class ListVehiclesRequest extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationRequest | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationRequest): ListVehiclesRequest;
  hasPagination(): boolean;
  clearPagination(): ListVehiclesRequest;

  getOrderbyList(): Array<resources_common_database_database_pb.OrderBy>;
  setOrderbyList(value: Array<resources_common_database_database_pb.OrderBy>): ListVehiclesRequest;
  clearOrderbyList(): ListVehiclesRequest;
  addOrderby(value?: resources_common_database_database_pb.OrderBy, index?: number): resources_common_database_database_pb.OrderBy;

  getSearch(): string;
  setSearch(value: string): ListVehiclesRequest;

  getModel(): string;
  setModel(value: string): ListVehiclesRequest;

  getUserId(): number;
  setUserId(value: number): ListVehiclesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListVehiclesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListVehiclesRequest): ListVehiclesRequest.AsObject;
  static serializeBinaryToWriter(message: ListVehiclesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListVehiclesRequest;
  static deserializeBinaryFromReader(message: ListVehiclesRequest, reader: jspb.BinaryReader): ListVehiclesRequest;
}

export namespace ListVehiclesRequest {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationRequest.AsObject,
    orderbyList: Array<resources_common_database_database_pb.OrderBy.AsObject>,
    search: string,
    model: string,
    userId: number,
  }
}

export class ListVehiclesResponse extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationResponse | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationResponse): ListVehiclesResponse;
  hasPagination(): boolean;
  clearPagination(): ListVehiclesResponse;

  getVehiclesList(): Array<resources_vehicles_vehicles_pb.Vehicle>;
  setVehiclesList(value: Array<resources_vehicles_vehicles_pb.Vehicle>): ListVehiclesResponse;
  clearVehiclesList(): ListVehiclesResponse;
  addVehicles(value?: resources_vehicles_vehicles_pb.Vehicle, index?: number): resources_vehicles_vehicles_pb.Vehicle;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListVehiclesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListVehiclesResponse): ListVehiclesResponse.AsObject;
  static serializeBinaryToWriter(message: ListVehiclesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListVehiclesResponse;
  static deserializeBinaryFromReader(message: ListVehiclesResponse, reader: jspb.BinaryReader): ListVehiclesResponse;
}

export namespace ListVehiclesResponse {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationResponse.AsObject,
    vehiclesList: Array<resources_vehicles_vehicles_pb.Vehicle.AsObject>,
  }
}

