import * as jspb from 'google-protobuf'



export class DocumentCategory extends jspb.Message {
  getId(): number;
  setId(value: number): DocumentCategory;

  getName(): string;
  setName(value: string): DocumentCategory;

  getDescription(): string;
  setDescription(value: string): DocumentCategory;
  hasDescription(): boolean;
  clearDescription(): DocumentCategory;

  getJob(): string;
  setJob(value: string): DocumentCategory;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DocumentCategory.AsObject;
  static toObject(includeInstance: boolean, msg: DocumentCategory): DocumentCategory.AsObject;
  static serializeBinaryToWriter(message: DocumentCategory, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DocumentCategory;
  static deserializeBinaryFromReader(message: DocumentCategory, reader: jspb.BinaryReader): DocumentCategory;
}

export namespace DocumentCategory {
  export type AsObject = {
    id: number,
    name: string,
    description?: string,
    job: string,
  }

  export enum DescriptionCase { 
    _DESCRIPTION_NOT_SET = 0,
    DESCRIPTION = 3,
  }
}

