// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/wiki/page.proto" (package "resources.wiki", syntax proto3)
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
import { PageAccess } from "./access";
/**
 * @generated from protobuf message resources.wiki.Page
 */
export interface Page {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: number;
    /**
     * @generated from protobuf field: string job = 2;
     */
    job: string;
    /**
     * @generated from protobuf field: string path = 3;
     */
    path: string;
    /**
     * @generated from protobuf field: resources.wiki.ContentType content_type = 4;
     */
    contentType: ContentType;
    /**
     * @generated from protobuf field: resources.wiki.PageMeta meta = 5;
     */
    meta?: PageMeta;
    /**
     * @generated from protobuf field: string content = 6;
     */
    content: string;
    /**
     * @generated from protobuf field: resources.wiki.PageAccess access = 7;
     */
    access?: PageAccess;
}
/**
 * @generated from protobuf message resources.wiki.PageMeta
 */
export interface PageMeta {
    /**
     * @generated from protobuf field: string title = 1;
     */
    title: string;
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp updated_at = 3;
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: repeated string tags = 4;
     */
    tags: string[];
    /**
     * @generated from protobuf field: resources.users.UserShort author = 5;
     */
    author?: UserShort;
    /**
     * @generated from protobuf field: string description = 6;
     */
    description: string;
}
/**
 * @generated from protobuf enum resources.wiki.ContentType
 */
export enum ContentType {
    /**
     * @generated from protobuf enum value: CONTENT_TYPE_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from protobuf enum value: CONTENT_TYPE_HTML = 1;
     */
    HTML = 1,
    /**
     * @generated from protobuf enum value: CONTENT_TYPE_MARKDOWN = 2;
     */
    MARKDOWN = 2
}
// @generated message type with reflection information, may provide speed optimized methods
class Page$Type extends MessageType<Page> {
    constructor() {
        super("resources.wiki.Page", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "path", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 4, name: "content_type", kind: "enum", T: () => ["resources.wiki.ContentType", ContentType, "CONTENT_TYPE_"] },
            { no: 5, name: "meta", kind: "message", T: () => PageMeta },
            { no: 6, name: "content", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 7, name: "access", kind: "message", T: () => PageAccess }
        ]);
    }
    create(value?: PartialMessage<Page>): Page {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.job = "";
        message.path = "";
        message.contentType = 0;
        message.content = "";
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
                case /* string job */ 2:
                    message.job = reader.string();
                    break;
                case /* string path */ 3:
                    message.path = reader.string();
                    break;
                case /* resources.wiki.ContentType content_type */ 4:
                    message.contentType = reader.int32();
                    break;
                case /* resources.wiki.PageMeta meta */ 5:
                    message.meta = PageMeta.internalBinaryRead(reader, reader.uint32(), options, message.meta);
                    break;
                case /* string content */ 6:
                    message.content = reader.string();
                    break;
                case /* resources.wiki.PageAccess access */ 7:
                    message.access = PageAccess.internalBinaryRead(reader, reader.uint32(), options, message.access);
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
        /* string job = 2; */
        if (message.job !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.job);
        /* string path = 3; */
        if (message.path !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.path);
        /* resources.wiki.ContentType content_type = 4; */
        if (message.contentType !== 0)
            writer.tag(4, WireType.Varint).int32(message.contentType);
        /* resources.wiki.PageMeta meta = 5; */
        if (message.meta)
            PageMeta.internalBinaryWrite(message.meta, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* string content = 6; */
        if (message.content !== "")
            writer.tag(6, WireType.LengthDelimited).string(message.content);
        /* resources.wiki.PageAccess access = 7; */
        if (message.access)
            PageAccess.internalBinaryWrite(message.access, writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.wiki.Page
 */
export const Page = new Page$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PageMeta$Type extends MessageType<PageMeta> {
    constructor() {
        super("resources.wiki.PageMeta", [
            { no: 1, name: "title", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "tags", kind: "scalar", repeat: 2 /*RepeatType.UNPACKED*/, T: 9 /*ScalarType.STRING*/ },
            { no: 5, name: "author", kind: "message", T: () => UserShort },
            { no: 6, name: "description", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<PageMeta>): PageMeta {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.title = "";
        message.tags = [];
        message.description = "";
        if (value !== undefined)
            reflectionMergePartial<PageMeta>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: PageMeta): PageMeta {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string title */ 1:
                    message.title = reader.string();
                    break;
                case /* resources.timestamp.Timestamp created_at */ 2:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* resources.timestamp.Timestamp updated_at */ 3:
                    message.updatedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.updatedAt);
                    break;
                case /* repeated string tags */ 4:
                    message.tags.push(reader.string());
                    break;
                case /* resources.users.UserShort author */ 5:
                    message.author = UserShort.internalBinaryRead(reader, reader.uint32(), options, message.author);
                    break;
                case /* string description */ 6:
                    message.description = reader.string();
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
    internalBinaryWrite(message: PageMeta, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string title = 1; */
        if (message.title !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.title);
        /* resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* resources.timestamp.Timestamp updated_at = 3; */
        if (message.updatedAt)
            Timestamp.internalBinaryWrite(message.updatedAt, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* repeated string tags = 4; */
        for (let i = 0; i < message.tags.length; i++)
            writer.tag(4, WireType.LengthDelimited).string(message.tags[i]);
        /* resources.users.UserShort author = 5; */
        if (message.author)
            UserShort.internalBinaryWrite(message.author, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* string description = 6; */
        if (message.description !== "")
            writer.tag(6, WireType.LengthDelimited).string(message.description);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.wiki.PageMeta
 */
export const PageMeta = new PageMeta$Type();
