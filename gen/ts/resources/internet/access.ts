// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/internet/access.proto" (package "resources.internet", syntax proto3)
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
import { UserShort } from "../users/users";
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.internet.DomainAccess
 */
export interface DomainAccess {
    /**
     * @generated from protobuf field: repeated resources.internet.DomainJobAccess jobs = 1;
     */
    jobs: DomainJobAccess[]; // @gotags: alias:"job_access"
    /**
     * @generated from protobuf field: repeated resources.internet.DomainUserAccess users = 2;
     */
    users: DomainUserAccess[]; // @gotags: alias:"user_access"
}
/**
 * @generated from protobuf message resources.internet.DomainJobAccess
 */
export interface DomainJobAccess {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: number;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: uint64 target_id = 3;
     */
    targetId: number; // @gotags: alias:"page_id"
    /**
     * @generated from protobuf field: string job = 4;
     */
    job: string;
    /**
     * @generated from protobuf field: optional string job_label = 5;
     */
    jobLabel?: string;
    /**
     * @generated from protobuf field: int32 minimum_grade = 6;
     */
    minimumGrade: number;
    /**
     * @generated from protobuf field: optional string job_grade_label = 7;
     */
    jobGradeLabel?: string;
    /**
     * @generated from protobuf field: resources.internet.AccessLevel access = 8;
     */
    access: AccessLevel;
}
/**
 * @generated from protobuf message resources.internet.DomainUserAccess
 */
export interface DomainUserAccess {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: number;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: uint64 target_id = 3;
     */
    targetId: number; // @gotags: alias:"page_id"
    /**
     * @generated from protobuf field: int32 user_id = 4;
     */
    userId: number;
    /**
     * @generated from protobuf field: optional resources.users.UserShort user = 5;
     */
    user?: UserShort;
    /**
     * @generated from protobuf field: resources.internet.AccessLevel access = 6;
     */
    access: AccessLevel;
}
/**
 * @generated from protobuf enum resources.internet.AccessLevel
 */
export enum AccessLevel {
    /**
     * @generated from protobuf enum value: ACCESS_LEVEL_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from protobuf enum value: ACCESS_LEVEL_BLOCKED = 1;
     */
    BLOCKED = 1,
    /**
     * @generated from protobuf enum value: ACCESS_LEVEL_VIEW = 2;
     */
    VIEW = 2,
    /**
     * @generated from protobuf enum value: ACCESS_LEVEL_EDIT = 3;
     */
    EDIT = 3,
    /**
     * @generated from protobuf enum value: ACCESS_LEVEL_OWNER = 4;
     */
    OWNER = 4
}
// @generated message type with reflection information, may provide speed optimized methods
class DomainAccess$Type extends MessageType<DomainAccess> {
    constructor() {
        super("resources.internet.DomainAccess", [
            { no: 1, name: "jobs", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => DomainJobAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } },
            { no: 2, name: "users", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => DomainUserAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } }
        ]);
    }
    create(value?: PartialMessage<DomainAccess>): DomainAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.jobs = [];
        message.users = [];
        if (value !== undefined)
            reflectionMergePartial<DomainAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DomainAccess): DomainAccess {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.internet.DomainJobAccess jobs */ 1:
                    message.jobs.push(DomainJobAccess.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated resources.internet.DomainUserAccess users */ 2:
                    message.users.push(DomainUserAccess.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: DomainAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.internet.DomainJobAccess jobs = 1; */
        for (let i = 0; i < message.jobs.length; i++)
            DomainJobAccess.internalBinaryWrite(message.jobs[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.internet.DomainUserAccess users = 2; */
        for (let i = 0; i < message.users.length; i++)
            DomainUserAccess.internalBinaryWrite(message.users[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.internet.DomainAccess
 */
export const DomainAccess = new DomainAccess$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DomainJobAccess$Type extends MessageType<DomainJobAccess> {
    constructor() {
        super("resources.internet.DomainJobAccess", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "target_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 4, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 5, name: "job_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 6, name: "minimum_grade", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gte: 0 } } } },
            { no: 7, name: "job_grade_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 8, name: "access", kind: "enum", T: () => ["resources.internet.AccessLevel", AccessLevel, "ACCESS_LEVEL_"], options: { "validate.rules": { enum: { definedOnly: true } } } }
        ]);
    }
    create(value?: PartialMessage<DomainJobAccess>): DomainJobAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.targetId = 0;
        message.job = "";
        message.minimumGrade = 0;
        message.access = 0;
        if (value !== undefined)
            reflectionMergePartial<DomainJobAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DomainJobAccess): DomainJobAccess {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* optional resources.timestamp.Timestamp created_at */ 2:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* uint64 target_id */ 3:
                    message.targetId = reader.uint64().toNumber();
                    break;
                case /* string job */ 4:
                    message.job = reader.string();
                    break;
                case /* optional string job_label */ 5:
                    message.jobLabel = reader.string();
                    break;
                case /* int32 minimum_grade */ 6:
                    message.minimumGrade = reader.int32();
                    break;
                case /* optional string job_grade_label */ 7:
                    message.jobGradeLabel = reader.string();
                    break;
                case /* resources.internet.AccessLevel access */ 8:
                    message.access = reader.int32();
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
    internalBinaryWrite(message: DomainJobAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* optional resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* uint64 target_id = 3; */
        if (message.targetId !== 0)
            writer.tag(3, WireType.Varint).uint64(message.targetId);
        /* string job = 4; */
        if (message.job !== "")
            writer.tag(4, WireType.LengthDelimited).string(message.job);
        /* optional string job_label = 5; */
        if (message.jobLabel !== undefined)
            writer.tag(5, WireType.LengthDelimited).string(message.jobLabel);
        /* int32 minimum_grade = 6; */
        if (message.minimumGrade !== 0)
            writer.tag(6, WireType.Varint).int32(message.minimumGrade);
        /* optional string job_grade_label = 7; */
        if (message.jobGradeLabel !== undefined)
            writer.tag(7, WireType.LengthDelimited).string(message.jobGradeLabel);
        /* resources.internet.AccessLevel access = 8; */
        if (message.access !== 0)
            writer.tag(8, WireType.Varint).int32(message.access);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.internet.DomainJobAccess
 */
export const DomainJobAccess = new DomainJobAccess$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DomainUserAccess$Type extends MessageType<DomainUserAccess> {
    constructor() {
        super("resources.internet.DomainUserAccess", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "target_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 4, name: "user_id", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gt: 0 } } } },
            { no: 5, name: "user", kind: "message", T: () => UserShort },
            { no: 6, name: "access", kind: "enum", T: () => ["resources.internet.AccessLevel", AccessLevel, "ACCESS_LEVEL_"], options: { "validate.rules": { enum: { definedOnly: true } } } }
        ]);
    }
    create(value?: PartialMessage<DomainUserAccess>): DomainUserAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.targetId = 0;
        message.userId = 0;
        message.access = 0;
        if (value !== undefined)
            reflectionMergePartial<DomainUserAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DomainUserAccess): DomainUserAccess {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* optional resources.timestamp.Timestamp created_at */ 2:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* uint64 target_id */ 3:
                    message.targetId = reader.uint64().toNumber();
                    break;
                case /* int32 user_id */ 4:
                    message.userId = reader.int32();
                    break;
                case /* optional resources.users.UserShort user */ 5:
                    message.user = UserShort.internalBinaryRead(reader, reader.uint32(), options, message.user);
                    break;
                case /* resources.internet.AccessLevel access */ 6:
                    message.access = reader.int32();
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
    internalBinaryWrite(message: DomainUserAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* optional resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* uint64 target_id = 3; */
        if (message.targetId !== 0)
            writer.tag(3, WireType.Varint).uint64(message.targetId);
        /* int32 user_id = 4; */
        if (message.userId !== 0)
            writer.tag(4, WireType.Varint).int32(message.userId);
        /* optional resources.users.UserShort user = 5; */
        if (message.user)
            UserShort.internalBinaryWrite(message.user, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* resources.internet.AccessLevel access = 6; */
        if (message.access !== 0)
            writer.tag(6, WireType.Varint).int32(message.access);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.internet.DomainUserAccess
 */
export const DomainUserAccess = new DomainUserAccess$Type();
