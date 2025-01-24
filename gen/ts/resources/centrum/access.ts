// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/centrum/access.proto" (package "resources.centrum", syntax proto3)
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
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.centrum.UnitAccess
 */
export interface UnitAccess {
    /**
     * @generated from protobuf field: repeated resources.centrum.UnitJobAccess jobs = 1;
     */
    jobs: UnitJobAccess[]; // @gotags: alias:"job_access"
    /**
     * @generated from protobuf field: repeated resources.centrum.UnitQualificationAccess qualifications = 3;
     */
    qualifications: UnitQualificationAccess[]; // @gotags: alias:"qualification_access"
}
/**
 * @generated from protobuf message resources.centrum.UnitJobAccess
 */
export interface UnitJobAccess {
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
    targetId: number; // @gotags: alias:"calendar_id"
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
     * @generated from protobuf field: resources.centrum.UnitAccessLevel access = 8;
     */
    access: UnitAccessLevel;
}
/**
 * @generated from protobuf message resources.centrum.UnitUserAccess
 */
export interface UnitUserAccess {
}
/**
 * @generated from protobuf message resources.centrum.UnitQualificationAccess
 */
export interface UnitQualificationAccess {
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
     * @generated from protobuf field: resources.centrum.UnitAccessLevel access = 6;
     */
    access: UnitAccessLevel;
}
/**
 * @generated from protobuf enum resources.centrum.UnitAccessLevel
 */
export enum UnitAccessLevel {
    /**
     * @generated from protobuf enum value: UNIT_ACCESS_LEVEL_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from protobuf enum value: UNIT_ACCESS_LEVEL_BLOCKED = 1;
     */
    BLOCKED = 1,
    /**
     * @generated from protobuf enum value: UNIT_ACCESS_LEVEL_JOIN = 2;
     */
    JOIN = 2
}
// @generated message type with reflection information, may provide speed optimized methods
class UnitAccess$Type extends MessageType<UnitAccess> {
    constructor() {
        super("resources.centrum.UnitAccess", [
            { no: 1, name: "jobs", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => UnitJobAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } },
            { no: 3, name: "qualifications", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => UnitQualificationAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } }
        ]);
    }
    create(value?: PartialMessage<UnitAccess>): UnitAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.jobs = [];
        message.qualifications = [];
        if (value !== undefined)
            reflectionMergePartial<UnitAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UnitAccess): UnitAccess {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.centrum.UnitJobAccess jobs */ 1:
                    message.jobs.push(UnitJobAccess.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated resources.centrum.UnitQualificationAccess qualifications */ 3:
                    message.qualifications.push(UnitQualificationAccess.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: UnitAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.centrum.UnitJobAccess jobs = 1; */
        for (let i = 0; i < message.jobs.length; i++)
            UnitJobAccess.internalBinaryWrite(message.jobs[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.centrum.UnitQualificationAccess qualifications = 3; */
        for (let i = 0; i < message.qualifications.length; i++)
            UnitQualificationAccess.internalBinaryWrite(message.qualifications[i], writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.centrum.UnitAccess
 */
export const UnitAccess = new UnitAccess$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UnitJobAccess$Type extends MessageType<UnitJobAccess> {
    constructor() {
        super("resources.centrum.UnitJobAccess", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "target_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 4, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 5, name: "job_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 6, name: "minimum_grade", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gte: 0 } } } },
            { no: 7, name: "job_grade_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 8, name: "access", kind: "enum", T: () => ["resources.centrum.UnitAccessLevel", UnitAccessLevel, "UNIT_ACCESS_LEVEL_"], options: { "validate.rules": { enum: { definedOnly: true } } } }
        ]);
    }
    create(value?: PartialMessage<UnitJobAccess>): UnitJobAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.targetId = 0;
        message.job = "";
        message.minimumGrade = 0;
        message.access = 0;
        if (value !== undefined)
            reflectionMergePartial<UnitJobAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UnitJobAccess): UnitJobAccess {
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
                case /* resources.centrum.UnitAccessLevel access */ 8:
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
    internalBinaryWrite(message: UnitJobAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
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
        /* resources.centrum.UnitAccessLevel access = 8; */
        if (message.access !== 0)
            writer.tag(8, WireType.Varint).int32(message.access);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.centrum.UnitJobAccess
 */
export const UnitJobAccess = new UnitJobAccess$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UnitUserAccess$Type extends MessageType<UnitUserAccess> {
    constructor() {
        super("resources.centrum.UnitUserAccess", []);
    }
    create(value?: PartialMessage<UnitUserAccess>): UnitUserAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<UnitUserAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UnitUserAccess): UnitUserAccess {
        return target ?? this.create();
    }
    internalBinaryWrite(message: UnitUserAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.centrum.UnitUserAccess
 */
export const UnitUserAccess = new UnitUserAccess$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UnitQualificationAccess$Type extends MessageType<UnitQualificationAccess> {
    constructor() {
        super("resources.centrum.UnitQualificationAccess", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "target_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 4, name: "qualification_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 5, name: "qualification", kind: "message", T: () => QualificationShort },
            { no: 6, name: "access", kind: "enum", T: () => ["resources.centrum.UnitAccessLevel", UnitAccessLevel, "UNIT_ACCESS_LEVEL_"], options: { "validate.rules": { enum: { definedOnly: true } } } }
        ]);
    }
    create(value?: PartialMessage<UnitQualificationAccess>): UnitQualificationAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.targetId = 0;
        message.qualificationId = 0;
        message.access = 0;
        if (value !== undefined)
            reflectionMergePartial<UnitQualificationAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UnitQualificationAccess): UnitQualificationAccess {
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
                case /* resources.centrum.UnitAccessLevel access */ 6:
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
    internalBinaryWrite(message: UnitQualificationAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
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
        /* resources.centrum.UnitAccessLevel access = 6; */
        if (message.access !== 0)
            writer.tag(6, WireType.Varint).int32(message.access);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.centrum.UnitQualificationAccess
 */
export const UnitQualificationAccess = new UnitQualificationAccess$Type();
