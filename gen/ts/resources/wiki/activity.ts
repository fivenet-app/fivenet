// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/wiki/activity.proto" (package "resources.wiki", syntax proto3)
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
import { PageUserAccess } from "./access";
import { PageJobAccess } from "./access";
import { UserShort } from "../users/users";
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.wiki.PageActivity
 */
export interface PageActivity {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: number;
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: uint64 page_id = 3;
     */
    pageId: number;
    /**
     * @generated from protobuf field: resources.wiki.PageActivityType activity_type = 4;
     */
    activityType: PageActivityType;
    /**
     * @generated from protobuf field: optional int32 creator_id = 5;
     */
    creatorId?: number;
    /**
     * @generated from protobuf field: optional resources.users.UserShort creator = 6;
     */
    creator?: UserShort; // @gotags: alias:"creator"
    /**
     * @generated from protobuf field: string creator_job = 7;
     */
    creatorJob: string;
    /**
     * @generated from protobuf field: optional string creator_job_label = 8;
     */
    creatorJobLabel?: string;
    /**
     * @generated from protobuf field: optional string reason = 9;
     */
    reason?: string;
    /**
     * @generated from protobuf field: resources.wiki.PageActivityData data = 10;
     */
    data?: PageActivityData;
}
/**
 * @dbscanner: json
 *
 * @generated from protobuf message resources.wiki.PageActivityData
 */
export interface PageActivityData {
    /**
     * @generated from protobuf oneof: data
     */
    data: {
        oneofKind: "updated";
        /**
         * @generated from protobuf field: resources.wiki.PageUpdated updated = 1;
         */
        updated: PageUpdated;
    } | {
        oneofKind: "accessUpdated";
        /**
         * @generated from protobuf field: resources.wiki.PageAccessUpdated access_updated = 2;
         */
        accessUpdated: PageAccessUpdated;
    } | {
        oneofKind: undefined;
    };
}
/**
 * @generated from protobuf message resources.wiki.PageUpdated
 */
export interface PageUpdated {
    /**
     * @generated from protobuf field: optional string title_diff = 1;
     */
    titleDiff?: string;
    /**
     * @generated from protobuf field: optional string description_diff = 2;
     */
    descriptionDiff?: string;
    /**
     * @generated from protobuf field: optional string content_diff = 3;
     */
    contentDiff?: string;
}
/**
 * @generated from protobuf message resources.wiki.PageAccessUpdated
 */
export interface PageAccessUpdated {
    /**
     * @generated from protobuf field: resources.wiki.PageAccessJobsDiff jobs = 1;
     */
    jobs?: PageAccessJobsDiff;
    /**
     * @generated from protobuf field: resources.wiki.PageAccessUsersDiff users = 2;
     */
    users?: PageAccessUsersDiff;
}
/**
 * @generated from protobuf message resources.wiki.PageAccessJobsDiff
 */
export interface PageAccessJobsDiff {
    /**
     * @generated from protobuf field: repeated resources.wiki.PageJobAccess to_create = 1;
     */
    toCreate: PageJobAccess[];
    /**
     * @generated from protobuf field: repeated resources.wiki.PageJobAccess to_update = 2;
     */
    toUpdate: PageJobAccess[];
    /**
     * @generated from protobuf field: repeated resources.wiki.PageJobAccess to_delete = 3;
     */
    toDelete: PageJobAccess[];
}
/**
 * @generated from protobuf message resources.wiki.PageAccessUsersDiff
 */
export interface PageAccessUsersDiff {
    /**
     * @generated from protobuf field: repeated resources.wiki.PageUserAccess to_create = 1;
     */
    toCreate: PageUserAccess[];
    /**
     * @generated from protobuf field: repeated resources.wiki.PageUserAccess to_update = 2;
     */
    toUpdate: PageUserAccess[];
    /**
     * @generated from protobuf field: repeated resources.wiki.PageUserAccess to_delete = 3;
     */
    toDelete: PageUserAccess[];
}
/**
 * @generated from protobuf enum resources.wiki.PageActivityType
 */
export enum PageActivityType {
    /**
     * @generated from protobuf enum value: PAGE_ACTIVITY_TYPE_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * Base
     *
     * @generated from protobuf enum value: PAGE_ACTIVITY_TYPE_CREATED = 1;
     */
    CREATED = 1,
    /**
     * @generated from protobuf enum value: PAGE_ACTIVITY_TYPE_UPDATED = 2;
     */
    UPDATED = 2,
    /**
     * @generated from protobuf enum value: PAGE_ACTIVITY_TYPE_ACCESS_UPDATED = 3;
     */
    ACCESS_UPDATED = 3,
    /**
     * @generated from protobuf enum value: PAGE_ACTIVITY_TYPE_OWNER_CHANGED = 4;
     */
    OWNER_CHANGED = 4,
    /**
     * @generated from protobuf enum value: PAGE_ACTIVITY_TYPE_DELETED = 5;
     */
    DELETED = 5
}
// @generated message type with reflection information, may provide speed optimized methods
class PageActivity$Type extends MessageType<PageActivity> {
    constructor() {
        super("resources.wiki.PageActivity", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "page_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 4, name: "activity_type", kind: "enum", T: () => ["resources.wiki.PageActivityType", PageActivityType, "PAGE_ACTIVITY_TYPE_"] },
            { no: 5, name: "creator_id", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gt: 0 } } } },
            { no: 6, name: "creator", kind: "message", T: () => UserShort },
            { no: 7, name: "creator_job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 8, name: "creator_job_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 9, name: "reason", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 10, name: "data", kind: "message", T: () => PageActivityData }
        ]);
    }
    create(value?: PartialMessage<PageActivity>): PageActivity {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.pageId = 0;
        message.activityType = 0;
        message.creatorJob = "";
        if (value !== undefined)
            reflectionMergePartial<PageActivity>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: PageActivity): PageActivity {
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
                case /* uint64 page_id */ 3:
                    message.pageId = reader.uint64().toNumber();
                    break;
                case /* resources.wiki.PageActivityType activity_type */ 4:
                    message.activityType = reader.int32();
                    break;
                case /* optional int32 creator_id */ 5:
                    message.creatorId = reader.int32();
                    break;
                case /* optional resources.users.UserShort creator */ 6:
                    message.creator = UserShort.internalBinaryRead(reader, reader.uint32(), options, message.creator);
                    break;
                case /* string creator_job */ 7:
                    message.creatorJob = reader.string();
                    break;
                case /* optional string creator_job_label */ 8:
                    message.creatorJobLabel = reader.string();
                    break;
                case /* optional string reason */ 9:
                    message.reason = reader.string();
                    break;
                case /* resources.wiki.PageActivityData data */ 10:
                    message.data = PageActivityData.internalBinaryRead(reader, reader.uint32(), options, message.data);
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
    internalBinaryWrite(message: PageActivity, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* uint64 page_id = 3; */
        if (message.pageId !== 0)
            writer.tag(3, WireType.Varint).uint64(message.pageId);
        /* resources.wiki.PageActivityType activity_type = 4; */
        if (message.activityType !== 0)
            writer.tag(4, WireType.Varint).int32(message.activityType);
        /* optional int32 creator_id = 5; */
        if (message.creatorId !== undefined)
            writer.tag(5, WireType.Varint).int32(message.creatorId);
        /* optional resources.users.UserShort creator = 6; */
        if (message.creator)
            UserShort.internalBinaryWrite(message.creator, writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        /* string creator_job = 7; */
        if (message.creatorJob !== "")
            writer.tag(7, WireType.LengthDelimited).string(message.creatorJob);
        /* optional string creator_job_label = 8; */
        if (message.creatorJobLabel !== undefined)
            writer.tag(8, WireType.LengthDelimited).string(message.creatorJobLabel);
        /* optional string reason = 9; */
        if (message.reason !== undefined)
            writer.tag(9, WireType.LengthDelimited).string(message.reason);
        /* resources.wiki.PageActivityData data = 10; */
        if (message.data)
            PageActivityData.internalBinaryWrite(message.data, writer.tag(10, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.wiki.PageActivity
 */
export const PageActivity = new PageActivity$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PageActivityData$Type extends MessageType<PageActivityData> {
    constructor() {
        super("resources.wiki.PageActivityData", [
            { no: 1, name: "updated", kind: "message", oneof: "data", T: () => PageUpdated },
            { no: 2, name: "access_updated", kind: "message", oneof: "data", T: () => PageAccessUpdated }
        ]);
    }
    create(value?: PartialMessage<PageActivityData>): PageActivityData {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.data = { oneofKind: undefined };
        if (value !== undefined)
            reflectionMergePartial<PageActivityData>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: PageActivityData): PageActivityData {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.wiki.PageUpdated updated */ 1:
                    message.data = {
                        oneofKind: "updated",
                        updated: PageUpdated.internalBinaryRead(reader, reader.uint32(), options, (message.data as any).updated)
                    };
                    break;
                case /* resources.wiki.PageAccessUpdated access_updated */ 2:
                    message.data = {
                        oneofKind: "accessUpdated",
                        accessUpdated: PageAccessUpdated.internalBinaryRead(reader, reader.uint32(), options, (message.data as any).accessUpdated)
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
    internalBinaryWrite(message: PageActivityData, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.wiki.PageUpdated updated = 1; */
        if (message.data.oneofKind === "updated")
            PageUpdated.internalBinaryWrite(message.data.updated, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* resources.wiki.PageAccessUpdated access_updated = 2; */
        if (message.data.oneofKind === "accessUpdated")
            PageAccessUpdated.internalBinaryWrite(message.data.accessUpdated, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.wiki.PageActivityData
 */
export const PageActivityData = new PageActivityData$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PageUpdated$Type extends MessageType<PageUpdated> {
    constructor() {
        super("resources.wiki.PageUpdated", [
            { no: 1, name: "title_diff", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "description_diff", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "content_diff", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<PageUpdated>): PageUpdated {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<PageUpdated>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: PageUpdated): PageUpdated {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* optional string title_diff */ 1:
                    message.titleDiff = reader.string();
                    break;
                case /* optional string description_diff */ 2:
                    message.descriptionDiff = reader.string();
                    break;
                case /* optional string content_diff */ 3:
                    message.contentDiff = reader.string();
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
    internalBinaryWrite(message: PageUpdated, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* optional string title_diff = 1; */
        if (message.titleDiff !== undefined)
            writer.tag(1, WireType.LengthDelimited).string(message.titleDiff);
        /* optional string description_diff = 2; */
        if (message.descriptionDiff !== undefined)
            writer.tag(2, WireType.LengthDelimited).string(message.descriptionDiff);
        /* optional string content_diff = 3; */
        if (message.contentDiff !== undefined)
            writer.tag(3, WireType.LengthDelimited).string(message.contentDiff);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.wiki.PageUpdated
 */
export const PageUpdated = new PageUpdated$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PageAccessUpdated$Type extends MessageType<PageAccessUpdated> {
    constructor() {
        super("resources.wiki.PageAccessUpdated", [
            { no: 1, name: "jobs", kind: "message", T: () => PageAccessJobsDiff },
            { no: 2, name: "users", kind: "message", T: () => PageAccessUsersDiff }
        ]);
    }
    create(value?: PartialMessage<PageAccessUpdated>): PageAccessUpdated {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<PageAccessUpdated>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: PageAccessUpdated): PageAccessUpdated {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.wiki.PageAccessJobsDiff jobs */ 1:
                    message.jobs = PageAccessJobsDiff.internalBinaryRead(reader, reader.uint32(), options, message.jobs);
                    break;
                case /* resources.wiki.PageAccessUsersDiff users */ 2:
                    message.users = PageAccessUsersDiff.internalBinaryRead(reader, reader.uint32(), options, message.users);
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
    internalBinaryWrite(message: PageAccessUpdated, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.wiki.PageAccessJobsDiff jobs = 1; */
        if (message.jobs)
            PageAccessJobsDiff.internalBinaryWrite(message.jobs, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* resources.wiki.PageAccessUsersDiff users = 2; */
        if (message.users)
            PageAccessUsersDiff.internalBinaryWrite(message.users, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.wiki.PageAccessUpdated
 */
export const PageAccessUpdated = new PageAccessUpdated$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PageAccessJobsDiff$Type extends MessageType<PageAccessJobsDiff> {
    constructor() {
        super("resources.wiki.PageAccessJobsDiff", [
            { no: 1, name: "to_create", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => PageJobAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } },
            { no: 2, name: "to_update", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => PageJobAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } },
            { no: 3, name: "to_delete", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => PageJobAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } }
        ]);
    }
    create(value?: PartialMessage<PageAccessJobsDiff>): PageAccessJobsDiff {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.toCreate = [];
        message.toUpdate = [];
        message.toDelete = [];
        if (value !== undefined)
            reflectionMergePartial<PageAccessJobsDiff>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: PageAccessJobsDiff): PageAccessJobsDiff {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.wiki.PageJobAccess to_create */ 1:
                    message.toCreate.push(PageJobAccess.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated resources.wiki.PageJobAccess to_update */ 2:
                    message.toUpdate.push(PageJobAccess.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated resources.wiki.PageJobAccess to_delete */ 3:
                    message.toDelete.push(PageJobAccess.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: PageAccessJobsDiff, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.wiki.PageJobAccess to_create = 1; */
        for (let i = 0; i < message.toCreate.length; i++)
            PageJobAccess.internalBinaryWrite(message.toCreate[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.wiki.PageJobAccess to_update = 2; */
        for (let i = 0; i < message.toUpdate.length; i++)
            PageJobAccess.internalBinaryWrite(message.toUpdate[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.wiki.PageJobAccess to_delete = 3; */
        for (let i = 0; i < message.toDelete.length; i++)
            PageJobAccess.internalBinaryWrite(message.toDelete[i], writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.wiki.PageAccessJobsDiff
 */
export const PageAccessJobsDiff = new PageAccessJobsDiff$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PageAccessUsersDiff$Type extends MessageType<PageAccessUsersDiff> {
    constructor() {
        super("resources.wiki.PageAccessUsersDiff", [
            { no: 1, name: "to_create", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => PageUserAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } },
            { no: 2, name: "to_update", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => PageUserAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } },
            { no: 3, name: "to_delete", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => PageUserAccess, options: { "validate.rules": { repeated: { maxItems: "20" } } } }
        ]);
    }
    create(value?: PartialMessage<PageAccessUsersDiff>): PageAccessUsersDiff {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.toCreate = [];
        message.toUpdate = [];
        message.toDelete = [];
        if (value !== undefined)
            reflectionMergePartial<PageAccessUsersDiff>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: PageAccessUsersDiff): PageAccessUsersDiff {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.wiki.PageUserAccess to_create */ 1:
                    message.toCreate.push(PageUserAccess.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated resources.wiki.PageUserAccess to_update */ 2:
                    message.toUpdate.push(PageUserAccess.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated resources.wiki.PageUserAccess to_delete */ 3:
                    message.toDelete.push(PageUserAccess.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: PageAccessUsersDiff, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.wiki.PageUserAccess to_create = 1; */
        for (let i = 0; i < message.toCreate.length; i++)
            PageUserAccess.internalBinaryWrite(message.toCreate[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.wiki.PageUserAccess to_update = 2; */
        for (let i = 0; i < message.toUpdate.length; i++)
            PageUserAccess.internalBinaryWrite(message.toUpdate[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.wiki.PageUserAccess to_delete = 3; */
        for (let i = 0; i < message.toDelete.length; i++)
            PageUserAccess.internalBinaryWrite(message.toDelete[i], writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.wiki.PageAccessUsersDiff
 */
export const PageAccessUsersDiff = new PageAccessUsersDiff$Type();
