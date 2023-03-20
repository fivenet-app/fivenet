import * as jspb from 'google-protobuf'

import * as resources_common_database_database_pb from '../../resources/common/database/database_pb';
import * as resources_vehicles_vehicles_pb from '../../resources/vehicles/vehicles_pb';


export class FindVehiclesRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): FindVehiclesRequest;

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
    offset: number,
    orderbyList: Array<resources_common_database_database_pb.OrderBy.AsObject>,
    search: string,
    type: string,
    userId: number,
  }
}

export class FindVehiclesResponse extends jspb.Message {
  getTotalCount(): number;
  setTotalCount(value: number): FindVehiclesResponse;

  getOffset(): number;
  setOffset(value: number): FindVehiclesResponse;

  getEnd(): number;
  setEnd(value: number): FindVehiclesResponse;

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
    totalCount: number,
    offset: number,
    end: number,
    vehiclesList: Array<resources_vehicles_vehicles_pb.Vehicle.AsObject>,
  }
}

