// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/sync/activity.proto" (package "resources.sync", syntax proto3)
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
import { JobsUserProps as JobsUserProps$ } from "../jobs/colleagues";
import { UserProps as UserProps$ } from "../users/props";
/**
 * Connect an identifier/license to the provider with the specified external id
 * (e.g., auto discord social connect on server join)
 *
 * @generated from protobuf message resources.sync.UserOAuth2Conn
 */
export interface UserOAuth2Conn {
    /**
     * @generated from protobuf field: string provider_name = 1;
     */
    providerName: string;
    /**
     * @generated from protobuf field: string identifier = 2;
     */
    identifier: string;
    /**
     * @generated from protobuf field: string external_id = 3;
     */
    externalId: string;
    /**
     * @generated from protobuf field: string username = 4;
     */
    username: string;
}
/**
 * @generated from protobuf message resources.sync.UserProps
 */
export interface UserProps {
    /**
     * @generated from protobuf field: string reason = 1;
     */
    reason: string;
    /**
     * @generated from protobuf field: resources.users.UserProps props = 2;
     */
    props?: UserProps$;
}
/**
 * @generated from protobuf message resources.sync.JobsUserProps
 */
export interface JobsUserProps {
    /**
     * @generated from protobuf field: string reason = 1;
     */
    reason: string;
    /**
     * @generated from protobuf field: resources.jobs.JobsUserProps props = 2;
     */
    props?: JobsUserProps$;
}
// @generated message type with reflection information, may provide speed optimized methods
class UserOAuth2Conn$Type extends MessageType<UserOAuth2Conn> {
    constructor() {
        super("resources.sync.UserOAuth2Conn", [
            { no: 1, name: "provider_name", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "identifier", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "external_id", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 4, name: "username", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<UserOAuth2Conn>): UserOAuth2Conn {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.providerName = "";
        message.identifier = "";
        message.externalId = "";
        message.username = "";
        if (value !== undefined)
            reflectionMergePartial<UserOAuth2Conn>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UserOAuth2Conn): UserOAuth2Conn {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string provider_name */ 1:
                    message.providerName = reader.string();
                    break;
                case /* string identifier */ 2:
                    message.identifier = reader.string();
                    break;
                case /* string external_id */ 3:
                    message.externalId = reader.string();
                    break;
                case /* string username */ 4:
                    message.username = reader.string();
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
    internalBinaryWrite(message: UserOAuth2Conn, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string provider_name = 1; */
        if (message.providerName !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.providerName);
        /* string identifier = 2; */
        if (message.identifier !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.identifier);
        /* string external_id = 3; */
        if (message.externalId !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.externalId);
        /* string username = 4; */
        if (message.username !== "")
            writer.tag(4, WireType.LengthDelimited).string(message.username);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.sync.UserOAuth2Conn
 */
export const UserOAuth2Conn = new UserOAuth2Conn$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UserProps$Type extends MessageType<UserProps> {
    constructor() {
        super("resources.sync.UserProps", [
            { no: 1, name: "reason", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 2, name: "props", kind: "message", T: () => UserProps$, options: { "validate.rules": { message: { required: true } } } }
        ]);
    }
    create(value?: PartialMessage<UserProps>): UserProps {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.reason = "";
        if (value !== undefined)
            reflectionMergePartial<UserProps>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UserProps): UserProps {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string reason */ 1:
                    message.reason = reader.string();
                    break;
                case /* resources.users.UserProps props */ 2:
                    message.props = UserProps$.internalBinaryRead(reader, reader.uint32(), options, message.props);
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
    internalBinaryWrite(message: UserProps, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string reason = 1; */
        if (message.reason !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.reason);
        /* resources.users.UserProps props = 2; */
        if (message.props)
            UserProps$.internalBinaryWrite(message.props, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.sync.UserProps
 */
export const UserProps = new UserProps$Type();
// @generated message type with reflection information, may provide speed optimized methods
class JobsUserProps$Type extends MessageType<JobsUserProps> {
    constructor() {
        super("resources.sync.JobsUserProps", [
            { no: 1, name: "reason", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 2, name: "props", kind: "message", T: () => JobsUserProps$, options: { "validate.rules": { message: { required: true } } } }
        ]);
    }
    create(value?: PartialMessage<JobsUserProps>): JobsUserProps {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.reason = "";
        if (value !== undefined)
            reflectionMergePartial<JobsUserProps>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: JobsUserProps): JobsUserProps {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string reason */ 1:
                    message.reason = reader.string();
                    break;
                case /* resources.jobs.JobsUserProps props */ 2:
                    message.props = JobsUserProps$.internalBinaryRead(reader, reader.uint32(), options, message.props);
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
    internalBinaryWrite(message: JobsUserProps, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string reason = 1; */
        if (message.reason !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.reason);
        /* resources.jobs.JobsUserProps props = 2; */
        if (message.props)
            JobsUserProps$.internalBinaryWrite(message.props, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.sync.JobsUserProps
 */
export const JobsUserProps = new JobsUserProps$Type();
