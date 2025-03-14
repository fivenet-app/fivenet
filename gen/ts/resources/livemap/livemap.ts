// @generated by protobuf-ts 2.9.5 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/livemap/livemap.proto" (package "resources.livemap", syntax proto3)
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
import { Unit } from "../centrum/units";
import { Colleague } from "../jobs/colleagues";
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.livemap.UserMarker
 */
export interface UserMarker {
    /**
     * @generated from protobuf field: int32 user_id = 1;
     */
    userId: number;
    /**
     * @generated from protobuf field: double x = 2;
     */
    x: number;
    /**
     * @generated from protobuf field: double y = 3;
     */
    y: number;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 4;
     */
    updatedAt?: Timestamp;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: optional string postal = 5;
     */
    postal?: string;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: optional string color = 6;
     */
    color?: string;
    /**
     * @generated from protobuf field: string job = 7;
     */
    job: string;
    /**
     * @generated from protobuf field: string job_label = 8;
     */
    jobLabel: string;
    /**
     * @generated from protobuf field: resources.jobs.Colleague user = 9;
     */
    user?: Colleague; // @gotags: alias:"user"
    /**
     * @generated from protobuf field: optional uint64 unit_id = 10;
     */
    unitId?: number;
    /**
     * @generated from protobuf field: optional resources.centrum.Unit unit = 11;
     */
    unit?: Unit;
    /**
     * @generated from protobuf field: bool hidden = 12;
     */
    hidden: boolean;
}
/**
 * @generated from protobuf message resources.livemap.MarkerMarker
 */
export interface MarkerMarker {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: number;
    /**
     * @generated from protobuf field: double x = 2;
     */
    x: number;
    /**
     * @generated from protobuf field: double y = 3;
     */
    y: number;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 4;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 5;
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp expires_at = 6;
     */
    expiresAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp deleted_at = 7;
     */
    deletedAt?: Timestamp;
    /**
     * @sanitize
     *
     * @generated from protobuf field: string name = 8;
     */
    name: string;
    /**
     * @sanitize
     *
     * @generated from protobuf field: optional string description = 9;
     */
    description?: string;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: optional string postal = 10;
     */
    postal?: string;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: optional string color = 11;
     */
    color?: string;
    /**
     * @generated from protobuf field: string job = 12;
     */
    job: string;
    /**
     * @generated from protobuf field: string job_label = 13;
     */
    jobLabel: string;
    /**
     * @generated from protobuf field: resources.livemap.MarkerType type = 14;
     */
    type: MarkerType; // @gotags: alias:"markerType"
    /**
     * @generated from protobuf field: resources.livemap.MarkerData data = 15;
     */
    data?: MarkerData; // @gotags: alias:"markerData"
    /**
     * @generated from protobuf field: optional int32 creator_id = 16;
     */
    creatorId?: number;
    /**
     * @generated from protobuf field: optional resources.users.UserShort creator = 17;
     */
    creator?: UserShort;
}
/**
 * @dbscanner
 *
 * @generated from protobuf message resources.livemap.MarkerData
 */
export interface MarkerData {
    /**
     * @generated from protobuf oneof: data
     */
    data: {
        oneofKind: "circle";
        /**
         * @generated from protobuf field: resources.livemap.CircleMarker circle = 3;
         */
        circle: CircleMarker;
    } | {
        oneofKind: "icon";
        /**
         * @generated from protobuf field: resources.livemap.IconMarker icon = 4;
         */
        icon: IconMarker;
    } | {
        oneofKind: undefined;
    };
}
/**
 * @generated from protobuf message resources.livemap.CircleMarker
 */
export interface CircleMarker {
    /**
     * @generated from protobuf field: int32 radius = 1;
     */
    radius: number;
    /**
     * @generated from protobuf field: optional float opacity = 2;
     */
    opacity?: number;
}
/**
 * @generated from protobuf message resources.livemap.IconMarker
 */
export interface IconMarker {
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: string icon = 1;
     */
    icon: string;
}
/**
 * @generated from protobuf message resources.livemap.Coords
 */
export interface Coords {
    /**
     * @generated from protobuf field: double x = 1;
     */
    x: number;
    /**
     * @generated from protobuf field: double y = 2;
     */
    y: number;
}
/**
 * @generated from protobuf enum resources.livemap.MarkerType
 */
export enum MarkerType {
    /**
     * @generated from protobuf enum value: MARKER_TYPE_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from protobuf enum value: MARKER_TYPE_DOT = 1;
     */
    DOT = 1,
    /**
     * @generated from protobuf enum value: MARKER_TYPE_CIRCLE = 2;
     */
    CIRCLE = 2,
    /**
     * @generated from protobuf enum value: MARKER_TYPE_ICON = 3;
     */
    ICON = 3
}
// @generated message type with reflection information, may provide speed optimized methods
class UserMarker$Type extends MessageType<UserMarker> {
    constructor() {
        super("resources.livemap.UserMarker", [
            { no: 1, name: "user_id", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gt: 0 } } } },
            { no: 2, name: "x", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 3, name: "y", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 4, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 5, name: "postal", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "48" } } } },
            { no: 6, name: "color", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { len: "7", pattern: "^#[A-Fa-f0-9]{6}$" } } } },
            { no: 7, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 8, name: "job_label", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 9, name: "user", kind: "message", T: () => Colleague },
            { no: 10, name: "unit_id", kind: "scalar", opt: true, T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 11, name: "unit", kind: "message", T: () => Unit },
            { no: 12, name: "hidden", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<UserMarker>): UserMarker {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.userId = 0;
        message.x = 0;
        message.y = 0;
        message.job = "";
        message.jobLabel = "";
        message.hidden = false;
        if (value !== undefined)
            reflectionMergePartial<UserMarker>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UserMarker): UserMarker {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 user_id */ 1:
                    message.userId = reader.int32();
                    break;
                case /* double x */ 2:
                    message.x = reader.double();
                    break;
                case /* double y */ 3:
                    message.y = reader.double();
                    break;
                case /* optional resources.timestamp.Timestamp updated_at */ 4:
                    message.updatedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.updatedAt);
                    break;
                case /* optional string postal */ 5:
                    message.postal = reader.string();
                    break;
                case /* optional string color */ 6:
                    message.color = reader.string();
                    break;
                case /* string job */ 7:
                    message.job = reader.string();
                    break;
                case /* string job_label */ 8:
                    message.jobLabel = reader.string();
                    break;
                case /* resources.jobs.Colleague user */ 9:
                    message.user = Colleague.internalBinaryRead(reader, reader.uint32(), options, message.user);
                    break;
                case /* optional uint64 unit_id */ 10:
                    message.unitId = reader.uint64().toNumber();
                    break;
                case /* optional resources.centrum.Unit unit */ 11:
                    message.unit = Unit.internalBinaryRead(reader, reader.uint32(), options, message.unit);
                    break;
                case /* bool hidden */ 12:
                    message.hidden = reader.bool();
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
    internalBinaryWrite(message: UserMarker, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* int32 user_id = 1; */
        if (message.userId !== 0)
            writer.tag(1, WireType.Varint).int32(message.userId);
        /* double x = 2; */
        if (message.x !== 0)
            writer.tag(2, WireType.Bit64).double(message.x);
        /* double y = 3; */
        if (message.y !== 0)
            writer.tag(3, WireType.Bit64).double(message.y);
        /* optional resources.timestamp.Timestamp updated_at = 4; */
        if (message.updatedAt)
            Timestamp.internalBinaryWrite(message.updatedAt, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* optional string postal = 5; */
        if (message.postal !== undefined)
            writer.tag(5, WireType.LengthDelimited).string(message.postal);
        /* optional string color = 6; */
        if (message.color !== undefined)
            writer.tag(6, WireType.LengthDelimited).string(message.color);
        /* string job = 7; */
        if (message.job !== "")
            writer.tag(7, WireType.LengthDelimited).string(message.job);
        /* string job_label = 8; */
        if (message.jobLabel !== "")
            writer.tag(8, WireType.LengthDelimited).string(message.jobLabel);
        /* resources.jobs.Colleague user = 9; */
        if (message.user)
            Colleague.internalBinaryWrite(message.user, writer.tag(9, WireType.LengthDelimited).fork(), options).join();
        /* optional uint64 unit_id = 10; */
        if (message.unitId !== undefined)
            writer.tag(10, WireType.Varint).uint64(message.unitId);
        /* optional resources.centrum.Unit unit = 11; */
        if (message.unit)
            Unit.internalBinaryWrite(message.unit, writer.tag(11, WireType.LengthDelimited).fork(), options).join();
        /* bool hidden = 12; */
        if (message.hidden !== false)
            writer.tag(12, WireType.Varint).bool(message.hidden);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.livemap.UserMarker
 */
export const UserMarker = new UserMarker$Type();
// @generated message type with reflection information, may provide speed optimized methods
class MarkerMarker$Type extends MessageType<MarkerMarker> {
    constructor() {
        super("resources.livemap.MarkerMarker", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "x", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 3, name: "y", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 4, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 5, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 6, name: "expires_at", kind: "message", T: () => Timestamp },
            { no: 7, name: "deleted_at", kind: "message", T: () => Timestamp },
            { no: 8, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "1", maxLen: "255" } } } },
            { no: 9, name: "description", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 10, name: "postal", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "48" } } } },
            { no: 11, name: "color", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { len: "7", pattern: "^#[A-Fa-f0-9]{6}$" } } } },
            { no: 12, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 13, name: "job_label", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 14, name: "type", kind: "enum", T: () => ["resources.livemap.MarkerType", MarkerType, "MARKER_TYPE_"] },
            { no: 15, name: "data", kind: "message", T: () => MarkerData },
            { no: 16, name: "creator_id", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gt: 0 } } } },
            { no: 17, name: "creator", kind: "message", T: () => UserShort }
        ]);
    }
    create(value?: PartialMessage<MarkerMarker>): MarkerMarker {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.x = 0;
        message.y = 0;
        message.name = "";
        message.job = "";
        message.jobLabel = "";
        message.type = 0;
        if (value !== undefined)
            reflectionMergePartial<MarkerMarker>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: MarkerMarker): MarkerMarker {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* double x */ 2:
                    message.x = reader.double();
                    break;
                case /* double y */ 3:
                    message.y = reader.double();
                    break;
                case /* optional resources.timestamp.Timestamp created_at */ 4:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* optional resources.timestamp.Timestamp updated_at */ 5:
                    message.updatedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.updatedAt);
                    break;
                case /* optional resources.timestamp.Timestamp expires_at */ 6:
                    message.expiresAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.expiresAt);
                    break;
                case /* optional resources.timestamp.Timestamp deleted_at */ 7:
                    message.deletedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.deletedAt);
                    break;
                case /* string name */ 8:
                    message.name = reader.string();
                    break;
                case /* optional string description */ 9:
                    message.description = reader.string();
                    break;
                case /* optional string postal */ 10:
                    message.postal = reader.string();
                    break;
                case /* optional string color */ 11:
                    message.color = reader.string();
                    break;
                case /* string job */ 12:
                    message.job = reader.string();
                    break;
                case /* string job_label */ 13:
                    message.jobLabel = reader.string();
                    break;
                case /* resources.livemap.MarkerType type */ 14:
                    message.type = reader.int32();
                    break;
                case /* resources.livemap.MarkerData data */ 15:
                    message.data = MarkerData.internalBinaryRead(reader, reader.uint32(), options, message.data);
                    break;
                case /* optional int32 creator_id */ 16:
                    message.creatorId = reader.int32();
                    break;
                case /* optional resources.users.UserShort creator */ 17:
                    message.creator = UserShort.internalBinaryRead(reader, reader.uint32(), options, message.creator);
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
    internalBinaryWrite(message: MarkerMarker, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* double x = 2; */
        if (message.x !== 0)
            writer.tag(2, WireType.Bit64).double(message.x);
        /* double y = 3; */
        if (message.y !== 0)
            writer.tag(3, WireType.Bit64).double(message.y);
        /* optional resources.timestamp.Timestamp created_at = 4; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp updated_at = 5; */
        if (message.updatedAt)
            Timestamp.internalBinaryWrite(message.updatedAt, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp expires_at = 6; */
        if (message.expiresAt)
            Timestamp.internalBinaryWrite(message.expiresAt, writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp deleted_at = 7; */
        if (message.deletedAt)
            Timestamp.internalBinaryWrite(message.deletedAt, writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        /* string name = 8; */
        if (message.name !== "")
            writer.tag(8, WireType.LengthDelimited).string(message.name);
        /* optional string description = 9; */
        if (message.description !== undefined)
            writer.tag(9, WireType.LengthDelimited).string(message.description);
        /* optional string postal = 10; */
        if (message.postal !== undefined)
            writer.tag(10, WireType.LengthDelimited).string(message.postal);
        /* optional string color = 11; */
        if (message.color !== undefined)
            writer.tag(11, WireType.LengthDelimited).string(message.color);
        /* string job = 12; */
        if (message.job !== "")
            writer.tag(12, WireType.LengthDelimited).string(message.job);
        /* string job_label = 13; */
        if (message.jobLabel !== "")
            writer.tag(13, WireType.LengthDelimited).string(message.jobLabel);
        /* resources.livemap.MarkerType type = 14; */
        if (message.type !== 0)
            writer.tag(14, WireType.Varint).int32(message.type);
        /* resources.livemap.MarkerData data = 15; */
        if (message.data)
            MarkerData.internalBinaryWrite(message.data, writer.tag(15, WireType.LengthDelimited).fork(), options).join();
        /* optional int32 creator_id = 16; */
        if (message.creatorId !== undefined)
            writer.tag(16, WireType.Varint).int32(message.creatorId);
        /* optional resources.users.UserShort creator = 17; */
        if (message.creator)
            UserShort.internalBinaryWrite(message.creator, writer.tag(17, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.livemap.MarkerMarker
 */
export const MarkerMarker = new MarkerMarker$Type();
// @generated message type with reflection information, may provide speed optimized methods
class MarkerData$Type extends MessageType<MarkerData> {
    constructor() {
        super("resources.livemap.MarkerData", [
            { no: 3, name: "circle", kind: "message", oneof: "data", T: () => CircleMarker },
            { no: 4, name: "icon", kind: "message", oneof: "data", T: () => IconMarker }
        ]);
    }
    create(value?: PartialMessage<MarkerData>): MarkerData {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.data = { oneofKind: undefined };
        if (value !== undefined)
            reflectionMergePartial<MarkerData>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: MarkerData): MarkerData {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.livemap.CircleMarker circle */ 3:
                    message.data = {
                        oneofKind: "circle",
                        circle: CircleMarker.internalBinaryRead(reader, reader.uint32(), options, (message.data as any).circle)
                    };
                    break;
                case /* resources.livemap.IconMarker icon */ 4:
                    message.data = {
                        oneofKind: "icon",
                        icon: IconMarker.internalBinaryRead(reader, reader.uint32(), options, (message.data as any).icon)
                    };
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
    internalBinaryWrite(message: MarkerData, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.livemap.CircleMarker circle = 3; */
        if (message.data.oneofKind === "circle")
            CircleMarker.internalBinaryWrite(message.data.circle, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* resources.livemap.IconMarker icon = 4; */
        if (message.data.oneofKind === "icon")
            IconMarker.internalBinaryWrite(message.data.icon, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.livemap.MarkerData
 */
export const MarkerData = new MarkerData$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CircleMarker$Type extends MessageType<CircleMarker> {
    constructor() {
        super("resources.livemap.CircleMarker", [
            { no: 1, name: "radius", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "opacity", kind: "scalar", opt: true, T: 2 /*ScalarType.FLOAT*/, options: { "validate.rules": { float: { lte: 75, gte: 1 } } } }
        ]);
    }
    create(value?: PartialMessage<CircleMarker>): CircleMarker {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.radius = 0;
        if (value !== undefined)
            reflectionMergePartial<CircleMarker>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CircleMarker): CircleMarker {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 radius */ 1:
                    message.radius = reader.int32();
                    break;
                case /* optional float opacity */ 2:
                    message.opacity = reader.float();
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
    internalBinaryWrite(message: CircleMarker, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* int32 radius = 1; */
        if (message.radius !== 0)
            writer.tag(1, WireType.Varint).int32(message.radius);
        /* optional float opacity = 2; */
        if (message.opacity !== undefined)
            writer.tag(2, WireType.Bit32).float(message.opacity);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.livemap.CircleMarker
 */
export const CircleMarker = new CircleMarker$Type();
// @generated message type with reflection information, may provide speed optimized methods
class IconMarker$Type extends MessageType<IconMarker> {
    constructor() {
        super("resources.livemap.IconMarker", [
            { no: 1, name: "icon", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "128" } } } }
        ]);
    }
    create(value?: PartialMessage<IconMarker>): IconMarker {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.icon = "";
        if (value !== undefined)
            reflectionMergePartial<IconMarker>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: IconMarker): IconMarker {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string icon */ 1:
                    message.icon = reader.string();
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
    internalBinaryWrite(message: IconMarker, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string icon = 1; */
        if (message.icon !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.icon);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.livemap.IconMarker
 */
export const IconMarker = new IconMarker$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Coords$Type extends MessageType<Coords> {
    constructor() {
        super("resources.livemap.Coords", [
            { no: 1, name: "x", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ },
            { no: 2, name: "y", kind: "scalar", T: 1 /*ScalarType.DOUBLE*/ }
        ]);
    }
    create(value?: PartialMessage<Coords>): Coords {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.x = 0;
        message.y = 0;
        if (value !== undefined)
            reflectionMergePartial<Coords>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Coords): Coords {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* double x */ 1:
                    message.x = reader.double();
                    break;
                case /* double y */ 2:
                    message.y = reader.double();
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
    internalBinaryWrite(message: Coords, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* double x = 1; */
        if (message.x !== 0)
            writer.tag(1, WireType.Bit64).double(message.x);
        /* double y = 2; */
        if (message.y !== 0)
            writer.tag(2, WireType.Bit64).double(message.y);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.livemap.Coords
 */
export const Coords = new Coords$Type();
