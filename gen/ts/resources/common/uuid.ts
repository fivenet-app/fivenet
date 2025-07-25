// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "resources/common/uuid.proto" (package "resources.common", syntax proto3)
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
/**
 * @generated from protobuf message resources.common.UUID
 */
export interface UUID {
    /**
     * @generated from protobuf field: string uuid = 1
     */
    uuid: string;
}
// @generated message type with reflection information, may provide speed optimized methods
class UUID$Type extends MessageType<UUID> {
    constructor() {
        super("resources.common.UUID", [
            { no: 1, name: "uuid", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { uuid: true } } } }
        ]);
    }
    create(value?: PartialMessage<UUID>): UUID {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.uuid = "";
        if (value !== undefined)
            reflectionMergePartial<UUID>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UUID): UUID {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string uuid */ 1:
                    message.uuid = reader.string();
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
    internalBinaryWrite(message: UUID, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string uuid = 1; */
        if (message.uuid !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.uuid);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.common.UUID
 */
export const UUID = new UUID$Type();
