// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/centrum/disponents.proto" (package "resources.centrum", syntax proto3)
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
import { Colleague } from "../jobs/colleagues";
/**
 * @generated from protobuf message resources.centrum.Disponents
 */
export interface Disponents {
    /**
     * @generated from protobuf field: string job = 1;
     */
    job: string;
    /**
     * @generated from protobuf field: repeated resources.jobs.Colleague disponents = 2;
     */
    disponents: Colleague[];
}
// @generated message type with reflection information, may provide speed optimized methods
class Disponents$Type extends MessageType<Disponents> {
    constructor() {
        super("resources.centrum.Disponents", [
            { no: 1, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 2, name: "disponents", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Colleague }
        ]);
    }
    create(value?: PartialMessage<Disponents>): Disponents {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.job = "";
        message.disponents = [];
        if (value !== undefined)
            reflectionMergePartial<Disponents>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Disponents): Disponents {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string job */ 1:
                    message.job = reader.string();
                    break;
                case /* repeated resources.jobs.Colleague disponents */ 2:
                    message.disponents.push(Colleague.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: Disponents, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string job = 1; */
        if (message.job !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.job);
        /* repeated resources.jobs.Colleague disponents = 2; */
        for (let i = 0; i < message.disponents.length; i++)
            Colleague.internalBinaryWrite(message.disponents[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.centrum.Disponents
 */
export const Disponents = new Disponents$Type();
