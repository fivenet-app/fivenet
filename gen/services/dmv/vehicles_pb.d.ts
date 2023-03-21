import * as jspb from 'google-protobuf'

import * as resources_common_database_database_pb from '../../resources/common/database/database_pb';
import * as resources_vehicles_vehicles_pb from '../../resources/vehicles/vehicles_pb';


export class FindVehiclesRequest extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationRequest | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationRequest): FindVehiclesRequest;
  hasPagination(): boolean;
  clearPagination(): FindVehiclesRequest;

  getOrderbyList(): Array<resources_common_database_database_pb.OrderBy>;
  setOrderbyList(value: Array<resources_common_database_database_pb.OrderBy>): FindVehiclesRequest;
  clearOrderbyList(): FindVehiclesRequest;
  addOrderby(value?: resources_common_database_database_pb.OrderBy, index?: number): resources_common_database_database_pb.OrderBy;

  getSearch(): string;
  setSearch(value: string): FindVehiclesRequest;

  getType(): string;
  setType(value: string): FindVehiclesRequest;

  getUserId(): number;
  setUserId(value: number): FindVehiclesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindVehiclesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: FindVehiclesRequest): FindVehiclesRequest.AsObject;
  static serializeBinaryToWriter(message: FindVehiclesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindVehiclesRequest;
  static deserializeBinaryFromReader(message: FindVehiclesRequest, reader: jspb.BinaryReader): FindVehiclesRequest;
}

export namespace FindVehiclesRequest {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationRequest.AsObject,
    orderbyList: Array<resources_common_database_database_pb.OrderBy.AsObject>,
    search: string,
    type: string,
    userId: number,
  }
}

export class FindVehiclesResponse extends jspb.Message {
  getPagination(): resources_common_database_database_pb.PaginationResponse | undefined;
  setPagination(value?: resources_common_database_database_pb.PaginationResponse): FindVehiclesResponse;
  hasPagination(): boolean;
  clearPagination(): FindVehiclesResponse;

  getVehiclesList(): Array<resources_vehicles_vehicles_pb.Vehicle>;
  setVehiclesList(value: Array<resources_vehicles_vehicles_pb.Vehicle>): FindVehiclesResponse;
  clearVehiclesList(): FindVehiclesResponse;
  addVehicles(value?: resources_vehicles_vehicles_pb.Vehicle, index?: number): resources_vehicles_vehicles_pb.Vehicle;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindVehiclesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: FindVehiclesResponse): FindVehiclesResponse.AsObject;
  static serializeBinaryToWriter(message: FindVehiclesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindVehiclesResponse;
  static deserializeBinaryFromReader(message: FindVehiclesResponse, reader: jspb.BinaryReader): FindVehiclesResponse;
}

export namespace FindVehiclesResponse {
  export type AsObject = {
    pagination?: resources_common_database_database_pb.PaginationResponse.AsObject,
    vehiclesList: Array<resources_vehicles_vehicles_pb.Vehicle.AsObject>,
  }
}

