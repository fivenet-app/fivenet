// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "resources/internet/page.proto" (package "resources.internet", syntax proto3)
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
import { NodeType } from "../common/content/content";
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.internet.Page
 */
export interface Page {
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
     * @generated from protobuf field: uint64 domain_id = 5
     */
    domainId: number;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: string path = 6
     */
    path: string;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: string title = 7
     */
    title: string;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: string description = 8
     */
    description: string;
    /**
     * @generated from protobuf field: resources.internet.PageData data = 9
     */
    data?: PageData;
    /**
     * @generated from protobuf field: optional string creator_job = 10
     */
    creatorJob?: string;
    /**
     * @generated from protobuf field: optional int32 creator_id = 11
     */
    creatorId?: number;
}
/**
 * @dbscanner: json
 *
 * @generated from protobuf message resources.internet.PageData
 */
export interface PageData {
    /**
     * @generated from protobuf field: resources.internet.PageLayoutType layout_type = 1
     */
    layoutType: PageLayoutType;
    /**
     * @generated from protobuf field: optional resources.internet.ContentNode node = 2
     */
    node?: ContentNode;
}
/**
 * @generated from protobuf message resources.internet.ContentNode
 */
export interface ContentNode {
    /**
     * @generated from protobuf field: resources.common.content.NodeType type = 1
     */
    type: NodeType;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: optional string id = 2
     */
    id?: string;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: string tag = 3
     */
    tag: string;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: map<string, string> attrs = 4
     */
    attrs: {
        [key: string]: string;
    };
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: optional string text = 5
     */
    text?: string;
    /**
     * @generated from protobuf field: repeated resources.internet.ContentNode content = 6
     */
    content: ContentNode[];
    /**
     * @generated from protobuf field: repeated resources.internet.ContentNode slots = 7
     */
    slots: ContentNode[];
}
/**
 * @generated from protobuf enum resources.internet.PageLayoutType
 */
export enum PageLayoutType {
    /**
     * @generated from protobuf enum value: PAGE_LAYOUT_TYPE_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from protobuf enum value: PAGE_LAYOUT_TYPE_BASIC_PAGE = 1;
     */
    BASIC_PAGE = 1,
    /**
     * @generated from protobuf enum value: PAGE_LAYOUT_TYPE_LANDING_PAGE = 2;
     */
    LANDING_PAGE = 2
}
// @generated message type with reflection information, may provide speed optimized methods
class Page$Type extends MessageType<Page> {
    constructor() {
        super("resources.internet.Page", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "deleted_at", kind: "message", T: () => Timestamp },
            { no: 5, name: "domain_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 6, name: "path", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "128" } } } },
            { no: 7, name: "title", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "1", maxLen: "255" } } } },
            { no: 8, name: "description", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "3", maxLen: "512" } } } },
            { no: 9, name: "data", kind: "message", T: () => PageData, options: { "buf.validate.field": { required: true } } },
            { no: 10, name: "creator_job", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 11, name: "creator_id", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/ }
        ]);
    }
    create(value?: PartialMessage<Page>): Page {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.domainId = 0;
        message.path = "";
        message.title = "";
        message.description = "";
        if (value !== undefined)
            reflectionMergePartial<Page>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Page): Page {
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
                case /* uint64 domain_id */ 5:
                    message.domainId = reader.uint64().toNumber();
                    break;
                case /* string path */ 6:
                    message.path = reader.string();
                    break;
                case /* string title */ 7:
                    message.title = reader.string();
                    break;
                case /* string description */ 8:
                    message.description = reader.string();
                    break;
                case /* resources.internet.PageData data */ 9:
                    message.data = PageData.internalBinaryRead(reader, reader.uint32(), options, message.data);
                    break;
                case /* optional string creator_job */ 10:
                    message.creatorJob = reader.string();
                    break;
                case /* optional int32 creator_id */ 11:
                    message.creatorId = reader.int32();
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
    internalBinaryWrite(message: Page, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
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
        /* uint64 domain_id = 5; */
        if (message.domainId !== 0)
            writer.tag(5, WireType.Varint).uint64(message.domainId);
        /* string path = 6; */
        if (message.path !== "")
            writer.tag(6, WireType.LengthDelimited).string(message.path);
        /* string title = 7; */
        if (message.title !== "")
            writer.tag(7, WireType.LengthDelimited).string(message.title);
        /* string description = 8; */
        if (message.description !== "")
            writer.tag(8, WireType.LengthDelimited).string(message.description);
        /* resources.internet.PageData data = 9; */
        if (message.data)
            PageData.internalBinaryWrite(message.data, writer.tag(9, WireType.LengthDelimited).fork(), options).join();
        /* optional string creator_job = 10; */
        if (message.creatorJob !== undefined)
            writer.tag(10, WireType.LengthDelimited).string(message.creatorJob);
        /* optional int32 creator_id = 11; */
        if (message.creatorId !== undefined)
            writer.tag(11, WireType.Varint).int32(message.creatorId);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.internet.Page
 */
export const Page = new Page$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PageData$Type extends MessageType<PageData> {
    constructor() {
        super("resources.internet.PageData", [
            { no: 1, name: "layout_type", kind: "enum", T: () => ["resources.internet.PageLayoutType", PageLayoutType, "PAGE_LAYOUT_TYPE_"] },
            { no: 2, name: "node", kind: "message", T: () => ContentNode }
        ]);
    }
    create(value?: PartialMessage<PageData>): PageData {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.layoutType = 0;
        if (value !== undefined)
            reflectionMergePartial<PageData>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: PageData): PageData {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.internet.PageLayoutType layout_type */ 1:
                    message.layoutType = reader.int32();
                    break;
                case /* optional resources.internet.ContentNode node */ 2:
                    message.node = ContentNode.internalBinaryRead(reader, reader.uint32(), options, message.node);
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
    internalBinaryWrite(message: PageData, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.internet.PageLayoutType layout_type = 1; */
        if (message.layoutType !== 0)
            writer.tag(1, WireType.Varint).int32(message.layoutType);
        /* optional resources.internet.ContentNode node = 2; */
        if (message.node)
            ContentNode.internalBinaryWrite(message.node, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.internet.PageData
 */
export const PageData = new PageData$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ContentNode$Type extends MessageType<ContentNode> {
    constructor() {
        super("resources.internet.ContentNode", [
            { no: 1, name: "type", kind: "enum", T: () => ["resources.common.content.NodeType", NodeType, "NODE_TYPE_"], options: { "buf.validate.field": { enum: { definedOnly: true } } } },
            { no: 2, name: "id", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "tag", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 4, name: "attrs", kind: "map", K: 9 /*ScalarType.STRING*/, V: { kind: "scalar", T: 9 /*ScalarType.STRING*/ } },
            { no: 5, name: "text", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 6, name: "content", kind: "message", repeat: 2 /*RepeatType.UNPACKED*/, T: () => ContentNode },
            { no: 7, name: "slots", kind: "message", repeat: 2 /*RepeatType.UNPACKED*/, T: () => ContentNode }
        ]);
    }
    create(value?: PartialMessage<ContentNode>): ContentNode {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.type = 0;
        message.tag = "";
        message.attrs = {};
        message.content = [];
        message.slots = [];
        if (value !== undefined)
            reflectionMergePartial<ContentNode>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ContentNode): ContentNode {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.content.NodeType type */ 1:
                    message.type = reader.int32();
                    break;
                case /* optional string id */ 2:
                    message.id = reader.string();
                    break;
                case /* string tag */ 3:
                    message.tag = reader.string();
                    break;
                case /* map<string, string> attrs */ 4:
                    this.binaryReadMap4(message.attrs, reader, options);
                    break;
                case /* optional string text */ 5:
                    message.text = reader.string();
                    break;
                case /* repeated resources.internet.ContentNode content */ 6:
                    message.content.push(ContentNode.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated resources.internet.ContentNode slots */ 7:
                    message.slots.push(ContentNode.internalBinaryRead(reader, reader.uint32(), options));
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
    private binaryReadMap4(map: ContentNode["attrs"], reader: IBinaryReader, options: BinaryReadOptions): void {
        let len = reader.uint32(), end = reader.pos + len, key: keyof ContentNode["attrs"] | undefined, val: ContentNode["attrs"][any] | undefined;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case 1:
                    key = reader.string();
                    break;
                case 2:
                    val = reader.string();
                    break;
                default: throw new globalThis.Error("unknown map entry field for resources.internet.ContentNode.attrs");
            }
        }
        map[key ?? ""] = val ?? "";
    }
    internalBinaryWrite(message: ContentNode, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.content.NodeType type = 1; */
        if (message.type !== 0)
            writer.tag(1, WireType.Varint).int32(message.type);
        /* optional string id = 2; */
        if (message.id !== undefined)
            writer.tag(2, WireType.LengthDelimited).string(message.id);
        /* string tag = 3; */
        if (message.tag !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.tag);
        /* map<string, string> attrs = 4; */
        for (let k of globalThis.Object.keys(message.attrs))
            writer.tag(4, WireType.LengthDelimited).fork().tag(1, WireType.LengthDelimited).string(k).tag(2, WireType.LengthDelimited).string(message.attrs[k]).join();
        /* optional string text = 5; */
        if (message.text !== undefined)
            writer.tag(5, WireType.LengthDelimited).string(message.text);
        /* repeated resources.internet.ContentNode content = 6; */
        for (let i = 0; i < message.content.length; i++)
            ContentNode.internalBinaryWrite(message.content[i], writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.internet.ContentNode slots = 7; */
        for (let i = 0; i < message.slots.length; i++)
            ContentNode.internalBinaryWrite(message.slots[i], writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.internet.ContentNode
 */
export const ContentNode = new ContentNode$Type();
