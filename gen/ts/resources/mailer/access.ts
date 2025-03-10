// @generated by protobuf-ts 2.9.5 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/mailer/access.proto" (package "resources.mailer", syntax proto3)
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
import { QualificationShort } from "../qualifications/qualifications";
import { UserShort } from "../users/users";
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.mailer.Access
 */
export interface Access {
    /**
     * @generated from protobuf field: repeated resources.mailer.JobAccess jobs = 1;
     */
    jobs: JobAccess[]; // @gotags: alias:"job_access"
    /**
     * @generated from protobuf field: repeated resources.mailer.UserAccess users = 2;
     */
    users: UserAccess[]; // @gotags: alias:"user_access"
    /**
     * @generated from protobuf field: repeated resources.mailer.QualificationAccess qualifications = 3;
     */
    qualifications: QualificationAccess[]; // @gotags: alias:"qualification_access"
}
/**
 * @generated from protobuf message resources.mailer.JobAccess
 */
export interface JobAccess {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: number; // @gotags: sql:"primary_key" alias:"id"
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: uint64 target_id = 4;
     */
    targetId: number; // @gotags: alias:"email_id"
    /**
     * @generated from protobuf field: string job = 5;
     */
    job: string;
    /**
     * @generated from protobuf field: optional string job_label = 6;
     */
    jobLabel?: string;
    /**
     * @generated from protobuf field: int32 minimum_grade = 7;
     */
    minimumGrade: number;
    /**
     * @generated from protobuf field: optional string job_grade_label = 8;
     */
    jobGradeLabel?: string;
    /**
     * @generated from protobuf field: resources.mailer.AccessLevel access = 9;
     */
    access: AccessLevel;
}
/**
 * @generated from protobuf message resources.mailer.UserAccess
 */
export interface UserAccess {
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
    targetId: number; // @gotags: alias:"thread_id"
    /**
     * @generated from protobuf field: int32 user_id = 4;
     */
    userId: number;
    /**
     * @generated from protobuf field: optional resources.users.UserShort user = 5;
     */
    user?: UserShort;
    /**
     * @generated from protobuf field: resources.mailer.AccessLevel access = 6;
     */
    access: AccessLevel;
}
/**
 * @generated from protobuf message resources.mailer.QualificationAccess
 */
export interface QualificationAccess {
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
    targetId: number; // @gotags: alias:"thread_id"
    /**
     * @generated from protobuf field: uint64 qualification_id = 4;
     */
    qualificationId: number;
    /**
     * @generated from protobuf field: optional resources.qualifications.QualificationShort qualification = 5;
     */
    qualification?: QualificationShort;
    /**
     * @generated from protobuf field: resources.mailer.AccessLevel access = 6;
     */
    access: AccessLevel;
}
/**
 * @generated from protobuf enum resources.mailer.AccessLevel
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
     * @generated from protobuf enum value: ACCESS_LEVEL_READ = 2;
     */
    READ = 2,
    /**
     * @generated from protobuf enum value: ACCESS_LEVEL_WRITE = 3;
     */
    WRITE = 3,
    /**
     * @generated from protobuf enum value: ACCESS_LEVEL_MANAGE = 4;
     */
    MANAGE = 4
}
// @generated message type with reflection information, may provide speed optimized methods
class Access$Type extends MessageType<Access> {
    constructor() {
        super("resources.mailer.Access", [
            { no: 1, name: "jobs", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => JobAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } },
            { no: 2, name: "users", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => UserAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } },
            { no: 3, name: "qualifications", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => QualificationAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } }
        ]);
    }
    create(value?: PartialMessage<Access>): Access {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.jobs = [];
        message.users = [];
        message.qualifications = [];
        if (value !== undefined)
            reflectionMergePartial<Access>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Access): Access {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.mailer.JobAccess jobs */ 1:
                    message.jobs.push(JobAccess.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated resources.mailer.UserAccess users */ 2:
                    message.users.push(UserAccess.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated resources.mailer.QualificationAccess qualifications */ 3:
                    message.qualifications.push(QualificationAccess.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: Access, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.mailer.JobAccess jobs = 1; */
        for (let i = 0; i < message.jobs.length; i++)
            JobAccess.internalBinaryWrite(message.jobs[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.mailer.UserAccess users = 2; */
        for (let i = 0; i < message.users.length; i++)
            UserAccess.internalBinaryWrite(message.users[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.mailer.QualificationAccess qualifications = 3; */
        for (let i = 0; i < message.qualifications.length; i++)
            QualificationAccess.internalBinaryWrite(message.qualifications[i], writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.mailer.Access
 */
export const Access = new Access$Type();
// @generated message type with reflection information, may provide speed optimized methods
class JobAccess$Type extends MessageType<JobAccess> {
    constructor() {
        super("resources.mailer.JobAccess", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "target_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 5, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 6, name: "job_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 7, name: "minimum_grade", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gte: 0 } } } },
            { no: 8, name: "job_grade_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 9, name: "access", kind: "enum", T: () => ["resources.mailer.AccessLevel", AccessLevel, "ACCESS_LEVEL_"], options: { "validate.rules": { enum: { definedOnly: true } } } }
        ]);
    }
    create(value?: PartialMessage<JobAccess>): JobAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.targetId = 0;
        message.job = "";
        message.minimumGrade = 0;
        message.access = 0;
        if (value !== undefined)
            reflectionMergePartial<JobAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: JobAccess): JobAccess {
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
                case /* uint64 target_id */ 4:
                    message.targetId = reader.uint64().toNumber();
                    break;
                case /* string job */ 5:
                    message.job = reader.string();
                    break;
                case /* optional string job_label */ 6:
                    message.jobLabel = reader.string();
                    break;
                case /* int32 minimum_grade */ 7:
                    message.minimumGrade = reader.int32();
                    break;
                case /* optional string job_grade_label */ 8:
                    message.jobGradeLabel = reader.string();
                    break;
                case /* resources.mailer.AccessLevel access */ 9:
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
    internalBinaryWrite(message: JobAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* optional resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* uint64 target_id = 4; */
        if (message.targetId !== 0)
            writer.tag(4, WireType.Varint).uint64(message.targetId);
        /* string job = 5; */
        if (message.job !== "")
            writer.tag(5, WireType.LengthDelimited).string(message.job);
        /* optional string job_label = 6; */
        if (message.jobLabel !== undefined)
            writer.tag(6, WireType.LengthDelimited).string(message.jobLabel);
        /* int32 minimum_grade = 7; */
        if (message.minimumGrade !== 0)
            writer.tag(7, WireType.Varint).int32(message.minimumGrade);
        /* optional string job_grade_label = 8; */
        if (message.jobGradeLabel !== undefined)
            writer.tag(8, WireType.LengthDelimited).string(message.jobGradeLabel);
        /* resources.mailer.AccessLevel access = 9; */
        if (message.access !== 0)
            writer.tag(9, WireType.Varint).int32(message.access);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.mailer.JobAccess
 */
export const JobAccess = new JobAccess$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UserAccess$Type extends MessageType<UserAccess> {
    constructor() {
        super("resources.mailer.UserAccess", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "target_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 4, name: "user_id", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gte: 0 } } } },
            { no: 5, name: "user", kind: "message", T: () => UserShort },
            { no: 6, name: "access", kind: "enum", T: () => ["resources.mailer.AccessLevel", AccessLevel, "ACCESS_LEVEL_"], options: { "validate.rules": { enum: { definedOnly: true } } } }
        ]);
    }
    create(value?: PartialMessage<UserAccess>): UserAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.targetId = 0;
        message.userId = 0;
        message.access = 0;
        if (value !== undefined)
            reflectionMergePartial<UserAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UserAccess): UserAccess {
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
                case /* resources.mailer.AccessLevel access */ 6:
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
    internalBinaryWrite(message: UserAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
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
        /* resources.mailer.AccessLevel access = 6; */
        if (message.access !== 0)
            writer.tag(6, WireType.Varint).int32(message.access);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.mailer.UserAccess
 */
export const UserAccess = new UserAccess$Type();
// @generated message type with reflection information, may provide speed optimized methods
class QualificationAccess$Type extends MessageType<QualificationAccess> {
    constructor() {
        super("resources.mailer.QualificationAccess", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "target_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 4, name: "qualification_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 5, name: "qualification", kind: "message", T: () => QualificationShort },
            { no: 6, name: "access", kind: "enum", T: () => ["resources.mailer.AccessLevel", AccessLevel, "ACCESS_LEVEL_"], options: { "validate.rules": { enum: { definedOnly: true } } } }
        ]);
    }
    create(value?: PartialMessage<QualificationAccess>): QualificationAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.targetId = 0;
        message.qualificationId = 0;
        message.access = 0;
        if (value !== undefined)
            reflectionMergePartial<QualificationAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: QualificationAccess): QualificationAccess {
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
                case /* uint64 qualification_id */ 4:
                    message.qualificationId = reader.uint64().toNumber();
                    break;
                case /* optional resources.qualifications.QualificationShort qualification */ 5:
                    message.qualification = QualificationShort.internalBinaryRead(reader, reader.uint32(), options, message.qualification);
                    break;
                case /* resources.mailer.AccessLevel access */ 6:
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
    internalBinaryWrite(message: QualificationAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* optional resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* uint64 target_id = 3; */
        if (message.targetId !== 0)
            writer.tag(3, WireType.Varint).uint64(message.targetId);
        /* uint64 qualification_id = 4; */
        if (message.qualificationId !== 0)
            writer.tag(4, WireType.Varint).uint64(message.qualificationId);
        /* optional resources.qualifications.QualificationShort qualification = 5; */
        if (message.qualification)
            QualificationShort.internalBinaryWrite(message.qualification, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* resources.mailer.AccessLevel access = 6; */
        if (message.access !== 0)
            writer.tag(6, WireType.Varint).int32(message.access);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.mailer.QualificationAccess
 */
export const QualificationAccess = new QualificationAccess$Type();
