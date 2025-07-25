// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "resources/documents/category.proto" (package "resources.documents", syntax proto3)
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
 * @generated from protobuf message resources.documents.Category
 */
export interface Category {
    /**
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp created_at = 2
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp deleted_at = 3
     */
    deletedAt?: Timestamp;
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
     * @generated from protobuf field: optional string job = 6
     */
    job?: string;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: optional string color = 7
     */
    color?: string;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: optional string icon = 8
     */
    icon?: string;
}
// @generated message type with reflection information, may provide speed optimized methods
class Category$Type extends MessageType<Category> {
    constructor() {
        super("resources.documents.Category", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "deleted_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "3", maxLen: "128" } } } },
            { no: 5, name: "description", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "255" } } } },
            { no: 6, name: "job", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "20" } } } },
            { no: 7, name: "color", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "3", maxLen: "7" } } } },
            { no: 8, name: "icon", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "128", suffix: "Icon" } } } }
        ]);
    }
    create(value?: PartialMessage<Category>): Category {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.name = "";
        if (value !== undefined)
            reflectionMergePartial<Category>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Category): Category {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* resources.timestamp.Timestamp created_at */ 2:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* optional resources.timestamp.Timestamp deleted_at */ 3:
                    message.deletedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.deletedAt);
                    break;
                case /* string name */ 4:
                    message.name = reader.string();
                    break;
                case /* optional string description */ 5:
                    message.description = reader.string();
                    break;
                case /* optional string job */ 6:
                    message.job = reader.string();
                    break;
                case /* optional string color */ 7:
                    message.color = reader.string();
                    break;
                case /* optional string icon */ 8:
                    message.icon = reader.string();
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
    internalBinaryWrite(message: Category, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp deleted_at = 3; */
        if (message.deletedAt)
            Timestamp.internalBinaryWrite(message.deletedAt, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* string name = 4; */
        if (message.name !== "")
            writer.tag(4, WireType.LengthDelimited).string(message.name);
        /* optional string description = 5; */
        if (message.description !== undefined)
            writer.tag(5, WireType.LengthDelimited).string(message.description);
        /* optional string job = 6; */
        if (message.job !== undefined)
            writer.tag(6, WireType.LengthDelimited).string(message.job);
        /* optional string color = 7; */
        if (message.color !== undefined)
            writer.tag(7, WireType.LengthDelimited).string(message.color);
        /* optional string icon = 8; */
        if (message.icon !== undefined)
            writer.tag(8, WireType.LengthDelimited).string(message.icon);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.documents.Category
 */
export const Category = new Category$Type();
