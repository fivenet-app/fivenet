// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "resources/laws/laws.proto" (package "resources.laws", syntax proto3)
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
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.laws.LawBook
 */
export interface LawBook {
    /**
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 3
     */
    updatedAt?: Timestamp;
    /**
     * @sanitize
     *
     * @generated from protobuf field: string name = 4
     */
    name: string;
    /**
     * @sanitize
     *
     * @generated from protobuf field: optional string description = 5
     */
    description?: string;
    /**
     * @generated from protobuf field: repeated resources.laws.Law laws = 6
     */
    laws: Law[];
}
/**
 * @generated from protobuf message resources.laws.Law
 */
export interface Law {
    /**
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 3
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: uint64 lawbook_id = 4
     */
    lawbookId: number;
    /**
     * @sanitize
     *
     * @generated from protobuf field: string name = 5
     */
    name: string;
    /**
     * @sanitize
     *
     * @generated from protobuf field: optional string description = 6
     */
    description?: string;
    /**
     * @sanitize
     *
     * @generated from protobuf field: optional string hint = 7
     */
    hint?: string;
    /**
     * @generated from protobuf field: optional uint32 fine = 8
     */
    fine?: number;
    /**
     * @generated from protobuf field: optional uint32 detention_time = 9
     */
    detentionTime?: number;
    /**
     * @generated from protobuf field: optional uint32 stvo_points = 10
     */
    stvoPoints?: number;
}
// @generated message type with reflection information, may provide speed optimized methods
class LawBook$Type extends MessageType<LawBook> {
    constructor() {
        super("resources.laws.LawBook", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/, options: { "tagger.tags": "sql:\"primary_key\" alias:\"id\"" } },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "3", maxLen: "128" } } } },
            { no: 5, name: "description", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "255" } } } },
            { no: 6, name: "laws", kind: "message", repeat: 2 /*RepeatType.UNPACKED*/, T: () => Law }
        ]);
    }
    create(value?: PartialMessage<LawBook>): LawBook {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.name = "";
        message.laws = [];
        if (value !== undefined)
            reflectionMergePartial<LawBook>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: LawBook): LawBook {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* optional resources.timestamp.Timestamp created_at */ 2:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* optional resources.timestamp.Timestamp updated_at */ 3:
                    message.updatedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.updatedAt);
                    break;
                case /* string name */ 4:
                    message.name = reader.string();
                    break;
                case /* optional string description */ 5:
                    message.description = reader.string();
                    break;
                case /* repeated resources.laws.Law laws */ 6:
                    message.laws.push(Law.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: LawBook, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* optional resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp updated_at = 3; */
        if (message.updatedAt)
            Timestamp.internalBinaryWrite(message.updatedAt, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* string name = 4; */
        if (message.name !== "")
            writer.tag(4, WireType.LengthDelimited).string(message.name);
        /* optional string description = 5; */
        if (message.description !== undefined)
            writer.tag(5, WireType.LengthDelimited).string(message.description);
        /* repeated resources.laws.Law laws = 6; */
        for (let i = 0; i < message.laws.length; i++)
            Law.internalBinaryWrite(message.laws[i], writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.laws.LawBook
 */
export const LawBook = new LawBook$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Law$Type extends MessageType<Law> {
    constructor() {
        super("resources.laws.Law", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/, options: { "tagger.tags": "sql:\"primary_key\" alias:\"law.id\"" } },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "lawbook_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 5, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "3", maxLen: "128" } } } },
            { no: 6, name: "description", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "1024" } } } },
            { no: 7, name: "hint", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "512" } } } },
            { no: 8, name: "fine", kind: "scalar", opt: true, T: 13 /*ScalarType.UINT32*/ },
            { no: 9, name: "detention_time", kind: "scalar", opt: true, T: 13 /*ScalarType.UINT32*/ },
            { no: 10, name: "stvo_points", kind: "scalar", opt: true, T: 13 /*ScalarType.UINT32*/ }
        ]);
    }
    create(value?: PartialMessage<Law>): Law {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.lawbookId = 0;
        message.name = "";
        if (value !== undefined)
            reflectionMergePartial<Law>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Law): Law {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* optional resources.timestamp.Timestamp created_at */ 2:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* optional resources.timestamp.Timestamp updated_at */ 3:
                    message.updatedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.updatedAt);
                    break;
                case /* uint64 lawbook_id */ 4:
                    message.lawbookId = reader.uint64().toNumber();
                    break;
                case /* string name */ 5:
                    message.name = reader.string();
                    break;
                case /* optional string description */ 6:
                    message.description = reader.string();
                    break;
                case /* optional string hint */ 7:
                    message.hint = reader.string();
                    break;
                case /* optional uint32 fine */ 8:
                    message.fine = reader.uint32();
                    break;
                case /* optional uint32 detention_time */ 9:
                    message.detentionTime = reader.uint32();
                    break;
                case /* optional uint32 stvo_points */ 10:
                    message.stvoPoints = reader.uint32();
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
    internalBinaryWrite(message: Law, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* optional resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp updated_at = 3; */
        if (message.updatedAt)
            Timestamp.internalBinaryWrite(message.updatedAt, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* uint64 lawbook_id = 4; */
        if (message.lawbookId !== 0)
            writer.tag(4, WireType.Varint).uint64(message.lawbookId);
        /* string name = 5; */
        if (message.name !== "")
            writer.tag(5, WireType.LengthDelimited).string(message.name);
        /* optional string description = 6; */
        if (message.description !== undefined)
            writer.tag(6, WireType.LengthDelimited).string(message.description);
        /* optional string hint = 7; */
        if (message.hint !== undefined)
            writer.tag(7, WireType.LengthDelimited).string(message.hint);
        /* optional uint32 fine = 8; */
        if (message.fine !== undefined)
            writer.tag(8, WireType.Varint).uint32(message.fine);
        /* optional uint32 detention_time = 9; */
        if (message.detentionTime !== undefined)
            writer.tag(9, WireType.Varint).uint32(message.detentionTime);
        /* optional uint32 stvo_points = 10; */
        if (message.stvoPoints !== undefined)
            writer.tag(10, WireType.Varint).uint32(message.stvoPoints);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.laws.Law
 */
export const Law = new Law$Type();
