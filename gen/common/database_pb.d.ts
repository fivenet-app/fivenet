import * as jspb from 'google-protobuf'



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

