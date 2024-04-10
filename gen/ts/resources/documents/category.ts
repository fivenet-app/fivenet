// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/documents/category.proto" (package "resources.documents", syntax proto3)
// tslint:disable
import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import { WireType } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import { UnknownFieldHandler } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { reflectionMergePartial } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
/**
 * @generated from protobuf message resources.documents.Category
 */
export interface Category {
    /**
     * @generated from protobuf field: uint64 id = 1 [jstype = JS_STRING];
     */
    id: string;
    /**
     * @sanitize
     *
     * @generated from protobuf field: string name = 2;
     */
    name: string;
    /**
     * @sanitize
     *
     * @generated from protobuf field: optional string description = 3;
     */
    description?: string;
    /**
     * @generated from protobuf field: optional string job = 4;
     */
    job?: string;
}
// @generated message type with reflection information, may provide speed optimized methods
class Category$Type extends MessageType<Category> {
    constructor() {
        super("resources.documents.Category", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ },
            { no: 2, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "3", maxLen: "128" } } } },
            { no: 3, name: "description", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 4, name: "job", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } }
        ]);
    }
    create(value?: PartialMessage<Category>): Category {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = "0";
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
                case /* uint64 id = 1 [jstype = JS_STRING];*/ 1:
                    message.id = reader.uint64().toString();
                    break;
                case /* string name */ 2:
                    message.name = reader.string();
                    break;
                case /* optional string description */ 3:
                    message.description = reader.string();
                    break;
                case /* optional string job */ 4:
                    message.job = reader.string();
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
        /* uint64 id = 1 [jstype = JS_STRING]; */
        if (message.id !== "0")
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* string name = 2; */
        if (message.name !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.name);
        /* optional string description = 3; */
        if (message.description !== undefined)
            writer.tag(3, WireType.LengthDelimited).string(message.description);
        /* optional string job = 4; */
        if (message.job !== undefined)
            writer.tag(4, WireType.LengthDelimited).string(message.job);
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
