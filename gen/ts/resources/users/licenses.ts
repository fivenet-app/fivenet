// @generated by protobuf-ts 2.9.6 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/users/licenses.proto" (package "resources.users", syntax proto3)
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
 * @generated from protobuf message resources.users.License
 */
export interface License {
    /**
     * @generated from protobuf field: string type = 1;
     */
    type: string;
    /**
     * @generated from protobuf field: string label = 2;
     */
    label: string;
}
/**
 * @generated from protobuf message resources.users.UserLicenses
 */
export interface UserLicenses {
    /**
     * @generated from protobuf field: int32 user_id = 1;
     */
    userId: number;
    /**
     * @generated from protobuf field: repeated resources.users.License licenses = 2;
     */
    licenses: License[];
}
// @generated message type with reflection information, may provide speed optimized methods
class License$Type extends MessageType<License> {
    constructor() {
        super("resources.users.License", [
            { no: 1, name: "type", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "3", maxLen: "60" } } } },
            { no: 2, name: "label", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<License>): License {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.type = "";
        message.label = "";
        if (value !== undefined)
            reflectionMergePartial<License>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: License): License {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string type */ 1:
                    message.type = reader.string();
                    break;
                case /* string label */ 2:
                    message.label = reader.string();
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
    internalBinaryWrite(message: License, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string type = 1; */
        if (message.type !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.type);
        /* string label = 2; */
        if (message.label !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.label);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.License
 */
export const License = new License$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UserLicenses$Type extends MessageType<UserLicenses> {
    constructor() {
        super("resources.users.UserLicenses", [
            { no: 1, name: "user_id", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gt: 0 } } } },
            { no: 2, name: "licenses", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => License }
        ]);
    }
    create(value?: PartialMessage<UserLicenses>): UserLicenses {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.userId = 0;
        message.licenses = [];
        if (value !== undefined)
            reflectionMergePartial<UserLicenses>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UserLicenses): UserLicenses {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 user_id */ 1:
                    message.userId = reader.int32();
                    break;
                case /* repeated resources.users.License licenses */ 2:
                    message.licenses.push(License.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: UserLicenses, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* int32 user_id = 1; */
        if (message.userId !== 0)
            writer.tag(1, WireType.Varint).int32(message.userId);
        /* repeated resources.users.License licenses = 2; */
        for (let i = 0; i < message.licenses.length; i++)
            License.internalBinaryWrite(message.licenses[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.UserLicenses
 */
export const UserLicenses = new UserLicenses$Type();
