// @generated by protobuf-ts 2.9.5 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/common/tests/objects.proto" (package "resources.common.tests", syntax proto3)
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
/**
 * **INTERNAL ONLY** SimpleObject is used as a test object where proto-based messages are required.
 *
 * @generated from protobuf message resources.common.tests.SimpleObject
 */
export interface SimpleObject {
    /**
     * @generated from protobuf field: string field1 = 1;
     */
    field1: string;
    /**
     * @generated from protobuf field: bool field2 = 2;
     */
    field2: boolean;
}
// @generated message type with reflection information, may provide speed optimized methods
class SimpleObject$Type extends MessageType<SimpleObject> {
    constructor() {
        super("resources.common.tests.SimpleObject", [
            { no: 1, name: "field1", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "field2", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<SimpleObject>): SimpleObject {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.field1 = "";
        message.field2 = false;
        if (value !== undefined)
            reflectionMergePartial<SimpleObject>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SimpleObject): SimpleObject {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string field1 */ 1:
                    message.field1 = reader.string();
                    break;
                case /* bool field2 */ 2:
                    message.field2 = reader.bool();
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
    internalBinaryWrite(message: SimpleObject, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string field1 = 1; */
        if (message.field1 !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.field1);
        /* bool field2 = 2; */
        if (message.field2 !== false)
            writer.tag(2, WireType.Varint).bool(message.field2);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.common.tests.SimpleObject
 */
export const SimpleObject = new SimpleObject$Type();
