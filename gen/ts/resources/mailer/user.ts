// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/mailer/user.proto" (package "resources.mailer", syntax proto3)
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
 * @generated from protobuf message resources.mailer.UserStatus
 */
export interface UserStatus {
    /**
     * @generated from protobuf field: int32 user_id = 1;
     */
    userId: number;
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp last_seen = 2;
     */
    lastSeen?: Timestamp;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: optional string status = 3;
     */
    status?: string;
}
// @generated message type with reflection information, may provide speed optimized methods
class UserStatus$Type extends MessageType<UserStatus> {
    constructor() {
        super("resources.mailer.UserStatus", [
            { no: 1, name: "user_id", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gte: 0 } } } },
            { no: 2, name: "last_seen", kind: "message", T: () => Timestamp },
            { no: 3, name: "status", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "40" } } } }
        ]);
    }
    create(value?: PartialMessage<UserStatus>): UserStatus {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.userId = 0;
        if (value !== undefined)
            reflectionMergePartial<UserStatus>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UserStatus): UserStatus {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 user_id */ 1:
                    message.userId = reader.int32();
                    break;
                case /* resources.timestamp.Timestamp last_seen */ 2:
                    message.lastSeen = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.lastSeen);
                    break;
                case /* optional string status */ 3:
                    message.status = reader.string();
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
    internalBinaryWrite(message: UserStatus, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* int32 user_id = 1; */
        if (message.userId !== 0)
            writer.tag(1, WireType.Varint).int32(message.userId);
        /* resources.timestamp.Timestamp last_seen = 2; */
        if (message.lastSeen)
            Timestamp.internalBinaryWrite(message.lastSeen, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* optional string status = 3; */
        if (message.status !== undefined)
            writer.tag(3, WireType.LengthDelimited).string(message.status);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.mailer.UserStatus
 */
export const UserStatus = new UserStatus$Type();
