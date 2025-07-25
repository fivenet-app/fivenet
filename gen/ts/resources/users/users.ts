// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "resources/users/users.proto" (package "resources.users", syntax proto3)
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
import { License } from "./licenses";
import { UserProps } from "./props";
/**
 * @generated from protobuf message resources.users.UserShort
 */
export interface UserShort {
    /**
     * @generated from protobuf field: int32 user_id = 1
     */
    userId: number;
    /**
     * @generated from protobuf field: optional string identifier = 2
     */
    identifier?: string;
    /**
     * @generated from protobuf field: string job = 3
     */
    job: string;
    /**
     * @generated from protobuf field: optional string job_label = 4
     */
    jobLabel?: string;
    /**
     * @generated from protobuf field: int32 job_grade = 5
     */
    jobGrade: number;
    /**
     * @generated from protobuf field: optional string job_grade_label = 6
     */
    jobGradeLabel?: string;
    /**
     * @generated from protobuf field: string firstname = 7
     */
    firstname: string;
    /**
     * @generated from protobuf field: string lastname = 8
     */
    lastname: string;
    /**
     * @generated from protobuf field: string dateofbirth = 9
     */
    dateofbirth: string;
    /**
     * @generated from protobuf field: optional string phone_number = 12
     */
    phoneNumber?: string;
    /**
     * @generated from protobuf field: optional uint64 avatar_file_id = 17
     */
    avatarFileId?: number;
    /**
     * @generated from protobuf field: optional string avatar = 18
     */
    avatar?: string;
}
/**
 * @generated from protobuf message resources.users.User
 */
export interface User {
    /**
     * @generated from protobuf field: int32 user_id = 1
     */
    userId: number;
    /**
     * @generated from protobuf field: optional string identifier = 2
     */
    identifier?: string;
    /**
     * @generated from protobuf field: string job = 3
     */
    job: string;
    /**
     * @generated from protobuf field: optional string job_label = 4
     */
    jobLabel?: string;
    /**
     * @generated from protobuf field: int32 job_grade = 5
     */
    jobGrade: number;
    /**
     * @generated from protobuf field: optional string job_grade_label = 6
     */
    jobGradeLabel?: string;
    /**
     * @generated from protobuf field: string firstname = 7
     */
    firstname: string;
    /**
     * @generated from protobuf field: string lastname = 8
     */
    lastname: string;
    /**
     * @generated from protobuf field: string dateofbirth = 9
     */
    dateofbirth: string;
    /**
     * @generated from protobuf field: optional string sex = 10
     */
    sex?: string;
    /**
     * @generated from protobuf field: optional string height = 11
     */
    height?: string;
    /**
     * @generated from protobuf field: optional string phone_number = 12
     */
    phoneNumber?: string;
    /**
     * @generated from protobuf field: optional int32 visum = 13
     */
    visum?: number;
    /**
     * @generated from protobuf field: optional int32 playtime = 14
     */
    playtime?: number;
    /**
     * @generated from protobuf field: resources.users.UserProps props = 15
     */
    props?: UserProps;
    /**
     * @generated from protobuf field: repeated resources.users.License licenses = 16
     */
    licenses: License[];
    /**
     * @generated from protobuf field: optional uint64 avatar_file_id = 17
     */
    avatarFileId?: number;
    /**
     * @generated from protobuf field: optional string avatar = 18
     */
    avatar?: string;
    /**
     * @generated from protobuf field: optional string group = 20
     */
    group?: string;
}
// @generated message type with reflection information, may provide speed optimized methods
class UserShort$Type extends MessageType<UserShort> {
    constructor() {
        super("resources.users.UserShort", [
            { no: 1, name: "user_id", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "buf.validate.field": { int32: { gt: 0 } }, "tagger.tags": "alias:\"id\"" } },
            { no: 2, name: "identifier", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "64" } } } },
            { no: 3, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { ignore: "IGNORE_IF_DEFAULT_VALUE", string: { maxLen: "20" } } } },
            { no: 4, name: "job_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "50" } } } },
            { no: 5, name: "job_grade", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "buf.validate.field": { ignore: "IGNORE_IF_DEFAULT_VALUE", int32: { gte: 0 } } } },
            { no: 6, name: "job_grade_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "50" } } } },
            { no: 7, name: "firstname", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "1", maxLen: "50" } } } },
            { no: 8, name: "lastname", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { ignore: "IGNORE_IF_DEFAULT_VALUE", string: { minLen: "1", maxLen: "50" } } } },
            { no: 9, name: "dateofbirth", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { ignore: "IGNORE_IF_DEFAULT_VALUE", string: { maxLen: "10" } } } },
            { no: 12, name: "phone_number", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "20" } } } },
            { no: 17, name: "avatar_file_id", kind: "scalar", opt: true, T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 18, name: "avatar", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<UserShort>): UserShort {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.userId = 0;
        message.job = "";
        message.jobGrade = 0;
        message.firstname = "";
        message.lastname = "";
        message.dateofbirth = "";
        if (value !== undefined)
            reflectionMergePartial<UserShort>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UserShort): UserShort {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 user_id */ 1:
                    message.userId = reader.int32();
                    break;
                case /* optional string identifier */ 2:
                    message.identifier = reader.string();
                    break;
                case /* string job */ 3:
                    message.job = reader.string();
                    break;
                case /* optional string job_label */ 4:
                    message.jobLabel = reader.string();
                    break;
                case /* int32 job_grade */ 5:
                    message.jobGrade = reader.int32();
                    break;
                case /* optional string job_grade_label */ 6:
                    message.jobGradeLabel = reader.string();
                    break;
                case /* string firstname */ 7:
                    message.firstname = reader.string();
                    break;
                case /* string lastname */ 8:
                    message.lastname = reader.string();
                    break;
                case /* string dateofbirth */ 9:
                    message.dateofbirth = reader.string();
                    break;
                case /* optional string phone_number */ 12:
                    message.phoneNumber = reader.string();
                    break;
                case /* optional uint64 avatar_file_id */ 17:
                    message.avatarFileId = reader.uint64().toNumber();
                    break;
                case /* optional string avatar */ 18:
                    message.avatar = reader.string();
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
    internalBinaryWrite(message: UserShort, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* int32 user_id = 1; */
        if (message.userId !== 0)
            writer.tag(1, WireType.Varint).int32(message.userId);
        /* optional string identifier = 2; */
        if (message.identifier !== undefined)
            writer.tag(2, WireType.LengthDelimited).string(message.identifier);
        /* string job = 3; */
        if (message.job !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.job);
        /* optional string job_label = 4; */
        if (message.jobLabel !== undefined)
            writer.tag(4, WireType.LengthDelimited).string(message.jobLabel);
        /* int32 job_grade = 5; */
        if (message.jobGrade !== 0)
            writer.tag(5, WireType.Varint).int32(message.jobGrade);
        /* optional string job_grade_label = 6; */
        if (message.jobGradeLabel !== undefined)
            writer.tag(6, WireType.LengthDelimited).string(message.jobGradeLabel);
        /* string firstname = 7; */
        if (message.firstname !== "")
            writer.tag(7, WireType.LengthDelimited).string(message.firstname);
        /* string lastname = 8; */
        if (message.lastname !== "")
            writer.tag(8, WireType.LengthDelimited).string(message.lastname);
        /* string dateofbirth = 9; */
        if (message.dateofbirth !== "")
            writer.tag(9, WireType.LengthDelimited).string(message.dateofbirth);
        /* optional string phone_number = 12; */
        if (message.phoneNumber !== undefined)
            writer.tag(12, WireType.LengthDelimited).string(message.phoneNumber);
        /* optional uint64 avatar_file_id = 17; */
        if (message.avatarFileId !== undefined)
            writer.tag(17, WireType.Varint).uint64(message.avatarFileId);
        /* optional string avatar = 18; */
        if (message.avatar !== undefined)
            writer.tag(18, WireType.LengthDelimited).string(message.avatar);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.UserShort
 */
export const UserShort = new UserShort$Type();
// @generated message type with reflection information, may provide speed optimized methods
class User$Type extends MessageType<User> {
    constructor() {
        super("resources.users.User", [
            { no: 1, name: "user_id", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "buf.validate.field": { int32: { gt: 0 } }, "tagger.tags": "alias:\"id\"" } },
            { no: 2, name: "identifier", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "64" } } } },
            { no: 3, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { ignore: "IGNORE_IF_DEFAULT_VALUE", string: { maxLen: "20" } } } },
            { no: 4, name: "job_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "50" } } } },
            { no: 5, name: "job_grade", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "buf.validate.field": { ignore: "IGNORE_IF_DEFAULT_VALUE", int32: { gte: 0 } } } },
            { no: 6, name: "job_grade_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "50" } } } },
            { no: 7, name: "firstname", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "1", maxLen: "50" } } } },
            { no: 8, name: "lastname", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { ignore: "IGNORE_IF_DEFAULT_VALUE", string: { minLen: "1", maxLen: "50" } } } },
            { no: 9, name: "dateofbirth", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { ignore: "IGNORE_IF_DEFAULT_VALUE", string: { maxLen: "10" } } } },
            { no: 10, name: "sex", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "1", maxLen: "2" } } } },
            { no: 11, name: "height", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 12, name: "phone_number", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "20" } } } },
            { no: 13, name: "visum", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/, options: { "buf.validate.field": { int32: { gte: 0 } } } },
            { no: 14, name: "playtime", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/, options: { "buf.validate.field": { int32: { gte: 0 } } } },
            { no: 15, name: "props", kind: "message", T: () => UserProps, options: { "tagger.tags": "alias:\"fivenet_user_props\"" } },
            { no: 16, name: "licenses", kind: "message", repeat: 2 /*RepeatType.UNPACKED*/, T: () => License, options: { "tagger.tags": "alias:\"user_licenses\"" } },
            { no: 17, name: "avatar_file_id", kind: "scalar", opt: true, T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 18, name: "avatar", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 20, name: "group", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "50" } } } }
        ]);
    }
    create(value?: PartialMessage<User>): User {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.userId = 0;
        message.job = "";
        message.jobGrade = 0;
        message.firstname = "";
        message.lastname = "";
        message.dateofbirth = "";
        message.licenses = [];
        if (value !== undefined)
            reflectionMergePartial<User>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: User): User {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 user_id */ 1:
                    message.userId = reader.int32();
                    break;
                case /* optional string identifier */ 2:
                    message.identifier = reader.string();
                    break;
                case /* string job */ 3:
                    message.job = reader.string();
                    break;
                case /* optional string job_label */ 4:
                    message.jobLabel = reader.string();
                    break;
                case /* int32 job_grade */ 5:
                    message.jobGrade = reader.int32();
                    break;
                case /* optional string job_grade_label */ 6:
                    message.jobGradeLabel = reader.string();
                    break;
                case /* string firstname */ 7:
                    message.firstname = reader.string();
                    break;
                case /* string lastname */ 8:
                    message.lastname = reader.string();
                    break;
                case /* string dateofbirth */ 9:
                    message.dateofbirth = reader.string();
                    break;
                case /* optional string sex */ 10:
                    message.sex = reader.string();
                    break;
                case /* optional string height */ 11:
                    message.height = reader.string();
                    break;
                case /* optional string phone_number */ 12:
                    message.phoneNumber = reader.string();
                    break;
                case /* optional int32 visum */ 13:
                    message.visum = reader.int32();
                    break;
                case /* optional int32 playtime */ 14:
                    message.playtime = reader.int32();
                    break;
                case /* resources.users.UserProps props */ 15:
                    message.props = UserProps.internalBinaryRead(reader, reader.uint32(), options, message.props);
                    break;
                case /* repeated resources.users.License licenses */ 16:
                    message.licenses.push(License.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* optional uint64 avatar_file_id */ 17:
                    message.avatarFileId = reader.uint64().toNumber();
                    break;
                case /* optional string avatar */ 18:
                    message.avatar = reader.string();
                    break;
                case /* optional string group */ 20:
                    message.group = reader.string();
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
    internalBinaryWrite(message: User, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* int32 user_id = 1; */
        if (message.userId !== 0)
            writer.tag(1, WireType.Varint).int32(message.userId);
        /* optional string identifier = 2; */
        if (message.identifier !== undefined)
            writer.tag(2, WireType.LengthDelimited).string(message.identifier);
        /* string job = 3; */
        if (message.job !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.job);
        /* optional string job_label = 4; */
        if (message.jobLabel !== undefined)
            writer.tag(4, WireType.LengthDelimited).string(message.jobLabel);
        /* int32 job_grade = 5; */
        if (message.jobGrade !== 0)
            writer.tag(5, WireType.Varint).int32(message.jobGrade);
        /* optional string job_grade_label = 6; */
        if (message.jobGradeLabel !== undefined)
            writer.tag(6, WireType.LengthDelimited).string(message.jobGradeLabel);
        /* string firstname = 7; */
        if (message.firstname !== "")
            writer.tag(7, WireType.LengthDelimited).string(message.firstname);
        /* string lastname = 8; */
        if (message.lastname !== "")
            writer.tag(8, WireType.LengthDelimited).string(message.lastname);
        /* string dateofbirth = 9; */
        if (message.dateofbirth !== "")
            writer.tag(9, WireType.LengthDelimited).string(message.dateofbirth);
        /* optional string sex = 10; */
        if (message.sex !== undefined)
            writer.tag(10, WireType.LengthDelimited).string(message.sex);
        /* optional string height = 11; */
        if (message.height !== undefined)
            writer.tag(11, WireType.LengthDelimited).string(message.height);
        /* optional string phone_number = 12; */
        if (message.phoneNumber !== undefined)
            writer.tag(12, WireType.LengthDelimited).string(message.phoneNumber);
        /* optional int32 visum = 13; */
        if (message.visum !== undefined)
            writer.tag(13, WireType.Varint).int32(message.visum);
        /* optional int32 playtime = 14; */
        if (message.playtime !== undefined)
            writer.tag(14, WireType.Varint).int32(message.playtime);
        /* resources.users.UserProps props = 15; */
        if (message.props)
            UserProps.internalBinaryWrite(message.props, writer.tag(15, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.users.License licenses = 16; */
        for (let i = 0; i < message.licenses.length; i++)
            License.internalBinaryWrite(message.licenses[i], writer.tag(16, WireType.LengthDelimited).fork(), options).join();
        /* optional uint64 avatar_file_id = 17; */
        if (message.avatarFileId !== undefined)
            writer.tag(17, WireType.Varint).uint64(message.avatarFileId);
        /* optional string avatar = 18; */
        if (message.avatar !== undefined)
            writer.tag(18, WireType.LengthDelimited).string(message.avatar);
        /* optional string group = 20; */
        if (message.group !== undefined)
            writer.tag(20, WireType.LengthDelimited).string(message.group);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.User
 */
export const User = new User$Type();
