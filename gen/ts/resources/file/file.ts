// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "resources/file/file.proto" (package "resources.file", syntax proto3)
// tslint:disable
// @ts-nocheck
import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import { WireType } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import { UnknownFieldHandler } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { reflectionMergePartial } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
import { FileMeta } from "./meta";
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.file.File
 */
export interface File {
    /**
     * @generated from protobuf field: optional uint64 parent_id = 1
     */
    parentId?: number;
    /**
     * @generated from protobuf field: uint64 id = 2
     */
    id: number;
    /**
     * @generated from protobuf field: string file_path = 3
     */
    filePath: string;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 4
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: int64 byte_size = 5
     */
    byteSize: number; // Bytes stored
    /**
     * @generated from protobuf field: string content_type = 6
     */
    contentType: string;
    /**
     * @generated from protobuf field: optional resources.file.FileMeta meta = 7
     */
    meta?: FileMeta;
}
// @generated message type with reflection information, may provide speed optimized methods
class File$Type extends MessageType<File> {
    constructor() {
        super("resources.file.File", [
            { no: 1, name: "parent_id", kind: "scalar", opt: true, T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/, options: { "buf.validate.field": { uint64: { gt: "0" } } } },
            { no: 2, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/, options: { "buf.validate.field": { uint64: { gt: "0" } } } },
            { no: 3, name: "file_path", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 4, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 5, name: "byte_size", kind: "scalar", T: 3 /*ScalarType.INT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 6, name: "content_type", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 7, name: "meta", kind: "message", T: () => FileMeta }
        ]);
    }
    create(value?: PartialMessage<File>): File {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.filePath = "";
        message.byteSize = 0;
        message.contentType = "";
        if (value !== undefined)
            reflectionMergePartial<File>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: File): File {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* optional uint64 parent_id */ 1:
                    message.parentId = reader.uint64().toNumber();
                    break;
                case /* uint64 id */ 2:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* string file_path */ 3:
                    message.filePath = reader.string();
                    break;
                case /* optional resources.timestamp.Timestamp created_at */ 4:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* int64 byte_size */ 5:
                    message.byteSize = reader.int64().toNumber();
                    break;
                case /* string content_type */ 6:
                    message.contentType = reader.string();
                    break;
                case /* optional resources.file.FileMeta meta */ 7:
                    message.meta = FileMeta.internalBinaryRead(reader, reader.uint32(), options, message.meta);
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: File, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* optional uint64 parent_id = 1; */
        if (message.parentId !== undefined)
            writer.tag(1, WireType.Varint).uint64(message.parentId);
        /* uint64 id = 2; */
        if (message.id !== 0)
            writer.tag(2, WireType.Varint).uint64(message.id);
        /* string file_path = 3; */
        if (message.filePath !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.filePath);
        /* optional resources.timestamp.Timestamp created_at = 4; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* int64 byte_size = 5; */
        if (message.byteSize !== 0)
            writer.tag(5, WireType.Varint).int64(message.byteSize);
        /* string content_type = 6; */
        if (message.contentType !== "")
            writer.tag(6, WireType.LengthDelimited).string(message.contentType);
        /* optional resources.file.FileMeta meta = 7; */
        if (message.meta)
            FileMeta.internalBinaryWrite(message.meta, writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.file.File
 */
export const File = new File$Type();
