// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "resources/qualifications/access.proto" (package "resources.qualifications", syntax proto3)
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
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.qualifications.QualificationAccess
 */
export interface QualificationAccess {
    /**
     * @generated from protobuf field: repeated resources.qualifications.QualificationJobAccess jobs = 1
     */
    jobs: QualificationJobAccess[];
}
/**
 * @generated from protobuf message resources.qualifications.QualificationJobAccess
 */
export interface QualificationJobAccess {
    /**
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: uint64 target_id = 4
     */
    targetId: number;
    /**
     * @generated from protobuf field: string job = 5
     */
    job: string;
    /**
     * @generated from protobuf field: optional string job_label = 6
     */
    jobLabel?: string;
    /**
     * @generated from protobuf field: int32 minimum_grade = 7
     */
    minimumGrade: number;
    /**
     * @generated from protobuf field: optional string job_grade_label = 8
     */
    jobGradeLabel?: string;
    /**
     * @generated from protobuf field: resources.qualifications.AccessLevel access = 9
     */
    access: AccessLevel;
}
/**
 * Dummy - DO NOT USE!
 *
 * @generated from protobuf message resources.qualifications.QualificationUserAccess
 */
export interface QualificationUserAccess {
}
/**
 * @generated from protobuf enum resources.qualifications.AccessLevel
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
     * @generated from protobuf enum value: ACCESS_LEVEL_REQUEST = 3;
     */
    REQUEST = 3,
    /**
     * @generated from protobuf enum value: ACCESS_LEVEL_TAKE = 4;
     */
    TAKE = 4,
    /**
     * @generated from protobuf enum value: ACCESS_LEVEL_GRADE = 5;
     */
    GRADE = 5,
    /**
     * @generated from protobuf enum value: ACCESS_LEVEL_EDIT = 6;
     */
    EDIT = 6
}
// @generated message type with reflection information, may provide speed optimized methods
class QualificationAccess$Type extends MessageType<QualificationAccess> {
    constructor() {
        super("resources.qualifications.QualificationAccess", [
            { no: 1, name: "jobs", kind: "message", repeat: 2 /*RepeatType.UNPACKED*/, T: () => QualificationJobAccess }
        ]);
    }
    create(value?: PartialMessage<QualificationAccess>): QualificationAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.jobs = [];
        if (value !== undefined)
            reflectionMergePartial<QualificationAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: QualificationAccess): QualificationAccess {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.qualifications.QualificationJobAccess jobs */ 1:
                    message.jobs.push(QualificationJobAccess.internalBinaryRead(reader, reader.uint32(), options));
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
        /* repeated resources.qualifications.QualificationJobAccess jobs = 1; */
        for (let i = 0; i < message.jobs.length; i++)
            QualificationJobAccess.internalBinaryWrite(message.jobs[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.qualifications.QualificationAccess
 */
export const QualificationAccess = new QualificationAccess$Type();
// @generated message type with reflection information, may provide speed optimized methods
class QualificationJobAccess$Type extends MessageType<QualificationJobAccess> {
    constructor() {
        super("resources.qualifications.QualificationJobAccess", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/, options: { "tagger.tags": "sql:\"primary_key\" alias:\"id\"" } },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "target_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 5, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "20" } } } },
            { no: 6, name: "job_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "50" } } } },
            { no: 7, name: "minimum_grade", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "buf.validate.field": { int32: { gte: 0 } } } },
            { no: 8, name: "job_grade_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "50" } } } },
            { no: 9, name: "access", kind: "enum", T: () => ["resources.qualifications.AccessLevel", AccessLevel, "ACCESS_LEVEL_"], options: { "buf.validate.field": { enum: { definedOnly: true } } } }
        ]);
    }
    create(value?: PartialMessage<QualificationJobAccess>): QualificationJobAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.targetId = 0;
        message.job = "";
        message.minimumGrade = 0;
        message.access = 0;
        if (value !== undefined)
            reflectionMergePartial<QualificationJobAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: QualificationJobAccess): QualificationJobAccess {
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
                case /* resources.qualifications.AccessLevel access */ 9:
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
    internalBinaryWrite(message: QualificationJobAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
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
        /* resources.qualifications.AccessLevel access = 9; */
        if (message.access !== 0)
            writer.tag(9, WireType.Varint).int32(message.access);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.qualifications.QualificationJobAccess
 */
export const QualificationJobAccess = new QualificationJobAccess$Type();
// @generated message type with reflection information, may provide speed optimized methods
class QualificationUserAccess$Type extends MessageType<QualificationUserAccess> {
    constructor() {
        super("resources.qualifications.QualificationUserAccess", []);
    }
    create(value?: PartialMessage<QualificationUserAccess>): QualificationUserAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<QualificationUserAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: QualificationUserAccess): QualificationUserAccess {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
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
    internalBinaryWrite(message: QualificationUserAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.qualifications.QualificationUserAccess
 */
export const QualificationUserAccess = new QualificationUserAccess$Type();
