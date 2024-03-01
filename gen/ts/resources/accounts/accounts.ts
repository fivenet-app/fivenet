// @generated by protobuf-ts 2.9.3 with parameter optimize_speed,long_type_bigint
// @generated from protobuf file "resources/accounts/accounts.proto" (package "resources.accounts", syntax proto3)
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
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.accounts.Account
 */
export interface Account {
    /**
     * @generated from protobuf field: uint64 id = 1 [jstype = JS_STRING];
     */
    id: string;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 3;
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: string username = 4;
     */
    username: string;
    /**
     * @generated from protobuf field: string license = 5;
     */
    license: string;
}
// @generated message type with reflection information, may provide speed optimized methods
class Account$Type extends MessageType<Account> {
    constructor() {
        super("resources.accounts.Account", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "username", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "24" } } } },
            { no: 5, name: "license", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "64" } } } }
        ]);
    }
    create(value?: PartialMessage<Account>): Account {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = "0";
        message.username = "";
        message.license = "";
        if (value !== undefined)
            reflectionMergePartial<Account>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Account): Account {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id = 1 [jstype = JS_STRING];*/ 1:
                    message.id = reader.uint64().toString();
                    break;
                case /* optional resources.timestamp.Timestamp created_at */ 2:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* optional resources.timestamp.Timestamp updated_at */ 3:
                    message.updatedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.updatedAt);
                    break;
                case /* string username */ 4:
                    message.username = reader.string();
                    break;
                case /* string license */ 5:
                    message.license = reader.string();
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
    internalBinaryWrite(message: Account, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1 [jstype = JS_STRING]; */
        if (message.id !== "0")
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* optional resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp updated_at = 3; */
        if (message.updatedAt)
            Timestamp.internalBinaryWrite(message.updatedAt, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* string username = 4; */
        if (message.username !== "")
            writer.tag(4, WireType.LengthDelimited).string(message.username);
        /* string license = 5; */
        if (message.license !== "")
            writer.tag(5, WireType.LengthDelimited).string(message.license);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.accounts.Account
 */
export const Account = new Account$Type();