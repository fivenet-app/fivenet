import * as jspb from 'google-protobuf'



export class PaginationRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): PaginationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PaginationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PaginationRequest): PaginationRequest.AsObject;
  static serializeBinaryToWriter(message: PaginationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PaginationRequest;
  static deserializeBinaryFromReader(message: PaginationRequest, reader: jspb.BinaryReader): PaginationRequest;
}

export namespace PaginationRequest {
  export type AsObject = {
    offset: number,
  }
}

export class PaginationResponse extends jspb.Message {
  getTotalCount(): number;
  setTotalCount(value: number): PaginationResponse;

  getOffset(): number;
  setOffset(value: number): PaginationResponse;

  getEnd(): number;
  setEnd(value: number): PaginationResponse;

  getPageSize(): number;
  setPageSize(value: number): PaginationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PaginationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PaginationResponse): PaginationResponse.AsObject;
  static serializeBinaryToWriter(message: PaginationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PaginationResponse;
  static deserializeBinaryFromReader(message: PaginationResponse, reader: jspb.BinaryReader): PaginationResponse;
}

export namespace PaginationResponse {
  export type AsObject = {
    totalCount: number,
    offset: number,
    end: number,
    pageSize: number,
  }
}

export class OrderBy extends jspb.Message {
  getColumn(): string;
  setColumn(value: string): OrderBy;

  getDesc(): boolean;
  setDesc(value: boolean): OrderBy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrderBy.AsObject;
  static toObject(includeInstance: boolean, msg: OrderBy): OrderBy.AsObject;
  static serializeBinaryToWriter(message: OrderBy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrderBy;
  static deserializeBinaryFromReader(message: OrderBy, reader: jspb.BinaryReader): OrderBy;
}

export namespace OrderBy {
  export type AsObject = {
    column: string,
    desc: boolean,
  }
}

