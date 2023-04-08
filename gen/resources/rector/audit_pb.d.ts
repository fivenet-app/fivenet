import * as jspb from 'google-protobuf'



export class AuditLogEntry extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuditLogEntry.AsObject;
  static toObject(includeInstance: boolean, msg: AuditLogEntry): AuditLogEntry.AsObject;
  static serializeBinaryToWriter(message: AuditLogEntry, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuditLogEntry;
  static deserializeBinaryFromReader(message: AuditLogEntry, reader: jspb.BinaryReader): AuditLogEntry;
}

export namespace AuditLogEntry {
  export type AsObject = {
  }
}

