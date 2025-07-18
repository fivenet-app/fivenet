// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "resources/mailer/email.proto" (package "resources.mailer", syntax proto3)
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
import { EmailSettings } from "./settings";
import { Access } from "./access";
import { UserShort } from "../users/users";
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.mailer.Email
 */
export interface Email {
    /**
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp created_at = 2
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 3
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp deleted_at = 4
     */
    deletedAt?: Timestamp;
    /**
     * @generated from protobuf field: bool deactivated = 5
     */
    deactivated: boolean;
    /**
     * @generated from protobuf field: optional string job = 6
     */
    job?: string;
    /**
     * @generated from protobuf field: optional int32 user_id = 7
     */
    userId?: number;
    /**
     * @generated from protobuf field: optional resources.users.UserShort user = 8
     */
    user?: UserShort;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: string email = 9
     */
    email: string;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp email_changed = 10
     */
    emailChanged?: Timestamp;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: optional string label = 11
     */
    label?: string;
    /**
     * @generated from protobuf field: resources.mailer.Access access = 12
     */
    access?: Access;
    /**
     * @generated from protobuf field: optional resources.mailer.EmailSettings settings = 13
     */
    settings?: EmailSettings;
}
// @generated message type with reflection information, may provide speed optimized methods
class Email$Type extends MessageType<Email> {
    constructor() {
        super("resources.mailer.Email", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "deleted_at", kind: "message", T: () => Timestamp },
            { no: 5, name: "deactivated", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 6, name: "job", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "40" } } } },
            { no: 7, name: "user_id", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/, options: { "buf.validate.field": { int32: { gt: 0 } } } },
            { no: 8, name: "user", kind: "message", T: () => UserShort },
            { no: 9, name: "email", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "6", maxLen: "80" } } } },
            { no: 10, name: "email_changed", kind: "message", T: () => Timestamp },
            { no: 11, name: "label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "128" } } } },
            { no: 12, name: "access", kind: "message", T: () => Access },
            { no: 13, name: "settings", kind: "message", T: () => EmailSettings }
        ]);
    }
    create(value?: PartialMessage<Email>): Email {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.deactivated = false;
        message.email = "";
        if (value !== undefined)
            reflectionMergePartial<Email>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Email): Email {
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
                case /* optional resources.timestamp.Timestamp updated_at */ 3:
                    message.updatedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.updatedAt);
                    break;
                case /* optional resources.timestamp.Timestamp deleted_at */ 4:
                    message.deletedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.deletedAt);
                    break;
                case /* bool deactivated */ 5:
                    message.deactivated = reader.bool();
                    break;
                case /* optional string job */ 6:
                    message.job = reader.string();
                    break;
                case /* optional int32 user_id */ 7:
                    message.userId = reader.int32();
                    break;
                case /* optional resources.users.UserShort user */ 8:
                    message.user = UserShort.internalBinaryRead(reader, reader.uint32(), options, message.user);
                    break;
                case /* string email */ 9:
                    message.email = reader.string();
                    break;
                case /* optional resources.timestamp.Timestamp email_changed */ 10:
                    message.emailChanged = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.emailChanged);
                    break;
                case /* optional string label */ 11:
                    message.label = reader.string();
                    break;
                case /* resources.mailer.Access access */ 12:
                    message.access = Access.internalBinaryRead(reader, reader.uint32(), options, message.access);
                    break;
                case /* optional resources.mailer.EmailSettings settings */ 13:
                    message.settings = EmailSettings.internalBinaryRead(reader, reader.uint32(), options, message.settings);
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
    internalBinaryWrite(message: Email, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp updated_at = 3; */
        if (message.updatedAt)
            Timestamp.internalBinaryWrite(message.updatedAt, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp deleted_at = 4; */
        if (message.deletedAt)
            Timestamp.internalBinaryWrite(message.deletedAt, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* bool deactivated = 5; */
        if (message.deactivated !== false)
            writer.tag(5, WireType.Varint).bool(message.deactivated);
        /* optional string job = 6; */
        if (message.job !== undefined)
            writer.tag(6, WireType.LengthDelimited).string(message.job);
        /* optional int32 user_id = 7; */
        if (message.userId !== undefined)
            writer.tag(7, WireType.Varint).int32(message.userId);
        /* optional resources.users.UserShort user = 8; */
        if (message.user)
            UserShort.internalBinaryWrite(message.user, writer.tag(8, WireType.LengthDelimited).fork(), options).join();
        /* string email = 9; */
        if (message.email !== "")
            writer.tag(9, WireType.LengthDelimited).string(message.email);
        /* optional resources.timestamp.Timestamp email_changed = 10; */
        if (message.emailChanged)
            Timestamp.internalBinaryWrite(message.emailChanged, writer.tag(10, WireType.LengthDelimited).fork(), options).join();
        /* optional string label = 11; */
        if (message.label !== undefined)
            writer.tag(11, WireType.LengthDelimited).string(message.label);
        /* resources.mailer.Access access = 12; */
        if (message.access)
            Access.internalBinaryWrite(message.access, writer.tag(12, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.mailer.EmailSettings settings = 13; */
        if (message.settings)
            EmailSettings.internalBinaryWrite(message.settings, writer.tag(13, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.mailer.Email
 */
export const Email = new Email$Type();
