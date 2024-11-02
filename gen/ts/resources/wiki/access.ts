// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/wiki/access.proto" (package "resources.wiki", syntax proto3)
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
 * @generated from protobuf message resources.wiki.PageAccess
 */
export interface PageAccess {
    /**
     * @generated from protobuf field: repeated resources.wiki.PageJobAccess jobs = 1;
     */
    jobs: PageJobAccess[]; // @gotags: alias:"job_access"
}
/**
 * @generated from protobuf message resources.wiki.PageJobAccess
 */
export interface PageJobAccess {
    /**
     * @generated from protobuf field: uint64 id = 1 [jstype = JS_STRING];
     */
    id: string;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: uint64 page_id = 3 [jstype = JS_STRING];
     */
    pageId: string;
    /**
     * @generated from protobuf field: string job = 4;
     */
    job: string;
    /**
     * @generated from protobuf field: optional string job_label = 5;
     */
    jobLabel?: string; // @gotags: alias:"job_label"
    /**
     * @generated from protobuf field: int32 minimum_grade = 6;
     */
    minimumGrade: number;
    /**
     * @generated from protobuf field: optional string job_grade_label = 7;
     */
    jobGradeLabel?: string; // @gotags: alias:"job_grade_label"
    /**
     * @generated from protobuf field: resources.wiki.AccessLevel access = 8;
     */
    access: AccessLevel;
}
/**
 * @generated from protobuf enum resources.wiki.AccessLevel
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
     * @generated from protobuf enum value: ACCESS_LEVEL_ACCESS = 3;
     */
    ACCESS = 3,
    /**
     * @generated from protobuf enum value: ACCESS_LEVEL_EDIT = 4;
     */
    EDIT = 4,
    /**
     * @generated from protobuf enum value: ACCESS_LEVEL_OWNER = 5;
     */
    OWNER = 5
}
// @generated message type with reflection information, may provide speed optimized methods
class PageAccess$Type extends MessageType<PageAccess> {
    constructor() {
        super("resources.wiki.PageAccess", [
            { no: 1, name: "jobs", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => PageJobAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } }
        ]);
    }
    create(value?: PartialMessage<PageAccess>): PageAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.jobs = [];
        if (value !== undefined)
            reflectionMergePartial<PageAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: PageAccess): PageAccess {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.wiki.PageJobAccess jobs */ 1:
                    message.jobs.push(PageJobAccess.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: PageAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.wiki.PageJobAccess jobs = 1; */
        for (let i = 0; i < message.jobs.length; i++)
            PageJobAccess.internalBinaryWrite(message.jobs[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.wiki.PageAccess
 */
export const PageAccess = new PageAccess$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PageJobAccess$Type extends MessageType<PageJobAccess> {
    constructor() {
        super("resources.wiki.PageJobAccess", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "page_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ },
            { no: 4, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 5, name: "job_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 6, name: "minimum_grade", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gte: 0 } } } },
            { no: 7, name: "job_grade_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 8, name: "access", kind: "enum", T: () => ["resources.wiki.AccessLevel", AccessLevel, "ACCESS_LEVEL_"], options: { "validate.rules": { enum: { definedOnly: true } } } }
        ]);
    }
    create(value?: PartialMessage<PageJobAccess>): PageJobAccess {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = "0";
        message.pageId = "0";
        message.job = "";
        message.minimumGrade = 0;
        message.access = 0;
        if (value !== undefined)
            reflectionMergePartial<PageJobAccess>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: PageJobAccess): PageJobAccess {
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
                case /* uint64 page_id = 3 [jstype = JS_STRING];*/ 3:
                    message.pageId = reader.uint64().toString();
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
                case /* resources.wiki.AccessLevel access */ 8:
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
    internalBinaryWrite(message: PageJobAccess, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1 [jstype = JS_STRING]; */
        if (message.id !== "0")
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* optional resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* uint64 page_id = 3 [jstype = JS_STRING]; */
        if (message.pageId !== "0")
            writer.tag(3, WireType.Varint).uint64(message.pageId);
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
        /* resources.wiki.AccessLevel access = 8; */
        if (message.access !== 0)
            writer.tag(8, WireType.Varint).int32(message.access);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.wiki.PageJobAccess
 */
export const PageJobAccess = new PageJobAccess$Type();
