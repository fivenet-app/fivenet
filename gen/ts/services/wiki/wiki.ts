// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "services/wiki/wiki.proto" (package "services.wiki", syntax proto3)
// tslint:disable
// @ts-nocheck
import { UploadFileResponse } from "../../resources/file/filestore";
import { UploadFileRequest } from "../../resources/file/filestore";
import { ServiceType } from "@protobuf-ts/runtime-rpc";
import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import { WireType } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import { UnknownFieldHandler } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { reflectionMergePartial } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
import { PageActivity } from "../../resources/wiki/activity";
import { ContentType } from "../../resources/common/content/content";
import { Page } from "../../resources/wiki/page";
import { PageShort } from "../../resources/wiki/page";
import { PaginationResponse } from "../../resources/common/database/database";
import { Sort } from "../../resources/common/database/database";
import { PaginationRequest } from "../../resources/common/database/database";
/**
 * @generated from protobuf message services.wiki.ListPagesRequest
 */
export interface ListPagesRequest {
    /**
     * @generated from protobuf field: resources.common.database.PaginationRequest pagination = 1
     */
    pagination?: PaginationRequest;
    /**
     * @generated from protobuf field: optional resources.common.database.Sort sort = 2
     */
    sort?: Sort;
    /**
     * Search params
     *
     * @generated from protobuf field: optional string job = 3
     */
    job?: string;
    /**
     * @generated from protobuf field: optional bool root_only = 4
     */
    rootOnly?: boolean;
    /**
     * @generated from protobuf field: optional string search = 5
     */
    search?: string;
}
/**
 * @generated from protobuf message services.wiki.ListPagesResponse
 */
export interface ListPagesResponse {
    /**
     * @generated from protobuf field: resources.common.database.PaginationResponse pagination = 1
     */
    pagination?: PaginationResponse;
    /**
     * @generated from protobuf field: repeated resources.wiki.PageShort pages = 2
     */
    pages: PageShort[];
}
/**
 * @generated from protobuf message services.wiki.GetPageRequest
 */
export interface GetPageRequest {
    /**
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
}
/**
 * @generated from protobuf message services.wiki.GetPageResponse
 */
export interface GetPageResponse {
    /**
     * @generated from protobuf field: resources.wiki.Page page = 1
     */
    page?: Page;
}
/**
 * @generated from protobuf message services.wiki.CreatePageRequest
 */
export interface CreatePageRequest {
    /**
     * @generated from protobuf field: optional uint64 parent_id = 1
     */
    parentId?: number;
    /**
     * @generated from protobuf field: resources.common.content.ContentType content_type = 2
     */
    contentType: ContentType;
}
/**
 * @generated from protobuf message services.wiki.CreatePageResponse
 */
export interface CreatePageResponse {
    /**
     * @generated from protobuf field: string job = 1
     */
    job: string;
    /**
     * @generated from protobuf field: uint64 id = 2
     */
    id: number;
}
/**
 * @generated from protobuf message services.wiki.UpdatePageRequest
 */
export interface UpdatePageRequest {
    /**
     * @generated from protobuf field: resources.wiki.Page page = 1
     */
    page?: Page;
}
/**
 * @generated from protobuf message services.wiki.UpdatePageResponse
 */
export interface UpdatePageResponse {
    /**
     * @generated from protobuf field: resources.wiki.Page page = 1
     */
    page?: Page;
}
/**
 * @generated from protobuf message services.wiki.DeletePageRequest
 */
export interface DeletePageRequest {
    /**
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
}
/**
 * @generated from protobuf message services.wiki.DeletePageResponse
 */
export interface DeletePageResponse {
}
/**
 * @generated from protobuf message services.wiki.ListPageActivityRequest
 */
export interface ListPageActivityRequest {
    /**
     * @generated from protobuf field: resources.common.database.PaginationRequest pagination = 1
     */
    pagination?: PaginationRequest;
    /**
     * @generated from protobuf field: uint64 page_id = 2
     */
    pageId: number;
}
/**
 * @generated from protobuf message services.wiki.ListPageActivityResponse
 */
export interface ListPageActivityResponse {
    /**
     * @generated from protobuf field: resources.common.database.PaginationResponse pagination = 1
     */
    pagination?: PaginationResponse;
    /**
     * @generated from protobuf field: repeated resources.wiki.PageActivity activity = 2
     */
    activity: PageActivity[];
}
// @generated message type with reflection information, may provide speed optimized methods
class ListPagesRequest$Type extends MessageType<ListPagesRequest> {
    constructor() {
        super("services.wiki.ListPagesRequest", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationRequest, options: { "buf.validate.field": { required: true } } },
            { no: 2, name: "sort", kind: "message", T: () => Sort },
            { no: 3, name: "job", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "50" } } } },
            { no: 4, name: "root_only", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ },
            { no: 5, name: "search", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "64" } } } }
        ]);
    }
    create(value?: PartialMessage<ListPagesRequest>): ListPagesRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<ListPagesRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListPagesRequest): ListPagesRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationRequest pagination */ 1:
                    message.pagination = PaginationRequest.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* optional resources.common.database.Sort sort */ 2:
                    message.sort = Sort.internalBinaryRead(reader, reader.uint32(), options, message.sort);
                    break;
                case /* optional string job */ 3:
                    message.job = reader.string();
                    break;
                case /* optional bool root_only */ 4:
                    message.rootOnly = reader.bool();
                    break;
                case /* optional string search */ 5:
                    message.search = reader.string();
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
    internalBinaryWrite(message: ListPagesRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationRequest pagination = 1; */
        if (message.pagination)
            PaginationRequest.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.common.database.Sort sort = 2; */
        if (message.sort)
            Sort.internalBinaryWrite(message.sort, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* optional string job = 3; */
        if (message.job !== undefined)
            writer.tag(3, WireType.LengthDelimited).string(message.job);
        /* optional bool root_only = 4; */
        if (message.rootOnly !== undefined)
            writer.tag(4, WireType.Varint).bool(message.rootOnly);
        /* optional string search = 5; */
        if (message.search !== undefined)
            writer.tag(5, WireType.LengthDelimited).string(message.search);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.wiki.ListPagesRequest
 */
export const ListPagesRequest = new ListPagesRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ListPagesResponse$Type extends MessageType<ListPagesResponse> {
    constructor() {
        super("services.wiki.ListPagesResponse", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationResponse, options: { "buf.validate.field": { required: true } } },
            { no: 2, name: "pages", kind: "message", repeat: 2 /*RepeatType.UNPACKED*/, T: () => PageShort }
        ]);
    }
    create(value?: PartialMessage<ListPagesResponse>): ListPagesResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.pages = [];
        if (value !== undefined)
            reflectionMergePartial<ListPagesResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListPagesResponse): ListPagesResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationResponse pagination */ 1:
                    message.pagination = PaginationResponse.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* repeated resources.wiki.PageShort pages */ 2:
                    message.pages.push(PageShort.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: ListPagesResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationResponse pagination = 1; */
        if (message.pagination)
            PaginationResponse.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.wiki.PageShort pages = 2; */
        for (let i = 0; i < message.pages.length; i++)
            PageShort.internalBinaryWrite(message.pages[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.wiki.ListPagesResponse
 */
export const ListPagesResponse = new ListPagesResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetPageRequest$Type extends MessageType<GetPageRequest> {
    constructor() {
        super("services.wiki.GetPageRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ }
        ]);
    }
    create(value?: PartialMessage<GetPageRequest>): GetPageRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        if (value !== undefined)
            reflectionMergePartial<GetPageRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: GetPageRequest): GetPageRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
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
    internalBinaryWrite(message: GetPageRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.wiki.GetPageRequest
 */
export const GetPageRequest = new GetPageRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetPageResponse$Type extends MessageType<GetPageResponse> {
    constructor() {
        super("services.wiki.GetPageResponse", [
            { no: 1, name: "page", kind: "message", T: () => Page }
        ]);
    }
    create(value?: PartialMessage<GetPageResponse>): GetPageResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<GetPageResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: GetPageResponse): GetPageResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.wiki.Page page */ 1:
                    message.page = Page.internalBinaryRead(reader, reader.uint32(), options, message.page);
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
    internalBinaryWrite(message: GetPageResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.wiki.Page page = 1; */
        if (message.page)
            Page.internalBinaryWrite(message.page, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.wiki.GetPageResponse
 */
export const GetPageResponse = new GetPageResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreatePageRequest$Type extends MessageType<CreatePageRequest> {
    constructor() {
        super("services.wiki.CreatePageRequest", [
            { no: 1, name: "parent_id", kind: "scalar", opt: true, T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/, options: { "buf.validate.field": { uint64: { gt: "0" } } } },
            { no: 2, name: "content_type", kind: "enum", T: () => ["resources.common.content.ContentType", ContentType, "CONTENT_TYPE_"], options: { "buf.validate.field": { enum: { definedOnly: true } } } }
        ]);
    }
    create(value?: PartialMessage<CreatePageRequest>): CreatePageRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.contentType = 0;
        if (value !== undefined)
            reflectionMergePartial<CreatePageRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CreatePageRequest): CreatePageRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* optional uint64 parent_id */ 1:
                    message.parentId = reader.uint64().toNumber();
                    break;
                case /* resources.common.content.ContentType content_type */ 2:
                    message.contentType = reader.int32();
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
    internalBinaryWrite(message: CreatePageRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* optional uint64 parent_id = 1; */
        if (message.parentId !== undefined)
            writer.tag(1, WireType.Varint).uint64(message.parentId);
        /* resources.common.content.ContentType content_type = 2; */
        if (message.contentType !== 0)
            writer.tag(2, WireType.Varint).int32(message.contentType);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.wiki.CreatePageRequest
 */
export const CreatePageRequest = new CreatePageRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreatePageResponse$Type extends MessageType<CreatePageResponse> {
    constructor() {
        super("services.wiki.CreatePageResponse", [
            { no: 1, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ }
        ]);
    }
    create(value?: PartialMessage<CreatePageResponse>): CreatePageResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.job = "";
        message.id = 0;
        if (value !== undefined)
            reflectionMergePartial<CreatePageResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CreatePageResponse): CreatePageResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string job */ 1:
                    message.job = reader.string();
                    break;
                case /* uint64 id */ 2:
                    message.id = reader.uint64().toNumber();
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
    internalBinaryWrite(message: CreatePageResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string job = 1; */
        if (message.job !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.job);
        /* uint64 id = 2; */
        if (message.id !== 0)
            writer.tag(2, WireType.Varint).uint64(message.id);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.wiki.CreatePageResponse
 */
export const CreatePageResponse = new CreatePageResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UpdatePageRequest$Type extends MessageType<UpdatePageRequest> {
    constructor() {
        super("services.wiki.UpdatePageRequest", [
            { no: 1, name: "page", kind: "message", T: () => Page, options: { "buf.validate.field": { required: true } } }
        ]);
    }
    create(value?: PartialMessage<UpdatePageRequest>): UpdatePageRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<UpdatePageRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UpdatePageRequest): UpdatePageRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.wiki.Page page */ 1:
                    message.page = Page.internalBinaryRead(reader, reader.uint32(), options, message.page);
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
    internalBinaryWrite(message: UpdatePageRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.wiki.Page page = 1; */
        if (message.page)
            Page.internalBinaryWrite(message.page, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.wiki.UpdatePageRequest
 */
export const UpdatePageRequest = new UpdatePageRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UpdatePageResponse$Type extends MessageType<UpdatePageResponse> {
    constructor() {
        super("services.wiki.UpdatePageResponse", [
            { no: 1, name: "page", kind: "message", T: () => Page }
        ]);
    }
    create(value?: PartialMessage<UpdatePageResponse>): UpdatePageResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<UpdatePageResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UpdatePageResponse): UpdatePageResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.wiki.Page page */ 1:
                    message.page = Page.internalBinaryRead(reader, reader.uint32(), options, message.page);
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
    internalBinaryWrite(message: UpdatePageResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.wiki.Page page = 1; */
        if (message.page)
            Page.internalBinaryWrite(message.page, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.wiki.UpdatePageResponse
 */
export const UpdatePageResponse = new UpdatePageResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeletePageRequest$Type extends MessageType<DeletePageRequest> {
    constructor() {
        super("services.wiki.DeletePageRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ }
        ]);
    }
    create(value?: PartialMessage<DeletePageRequest>): DeletePageRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        if (value !== undefined)
            reflectionMergePartial<DeletePageRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeletePageRequest): DeletePageRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
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
    internalBinaryWrite(message: DeletePageRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.wiki.DeletePageRequest
 */
export const DeletePageRequest = new DeletePageRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeletePageResponse$Type extends MessageType<DeletePageResponse> {
    constructor() {
        super("services.wiki.DeletePageResponse", []);
    }
    create(value?: PartialMessage<DeletePageResponse>): DeletePageResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<DeletePageResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeletePageResponse): DeletePageResponse {
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
    internalBinaryWrite(message: DeletePageResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.wiki.DeletePageResponse
 */
export const DeletePageResponse = new DeletePageResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ListPageActivityRequest$Type extends MessageType<ListPageActivityRequest> {
    constructor() {
        super("services.wiki.ListPageActivityRequest", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationRequest, options: { "buf.validate.field": { required: true } } },
            { no: 2, name: "page_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ }
        ]);
    }
    create(value?: PartialMessage<ListPageActivityRequest>): ListPageActivityRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.pageId = 0;
        if (value !== undefined)
            reflectionMergePartial<ListPageActivityRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListPageActivityRequest): ListPageActivityRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationRequest pagination */ 1:
                    message.pagination = PaginationRequest.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* uint64 page_id */ 2:
                    message.pageId = reader.uint64().toNumber();
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
    internalBinaryWrite(message: ListPageActivityRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationRequest pagination = 1; */
        if (message.pagination)
            PaginationRequest.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* uint64 page_id = 2; */
        if (message.pageId !== 0)
            writer.tag(2, WireType.Varint).uint64(message.pageId);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.wiki.ListPageActivityRequest
 */
export const ListPageActivityRequest = new ListPageActivityRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ListPageActivityResponse$Type extends MessageType<ListPageActivityResponse> {
    constructor() {
        super("services.wiki.ListPageActivityResponse", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationResponse, options: { "buf.validate.field": { required: true } } },
            { no: 2, name: "activity", kind: "message", repeat: 2 /*RepeatType.UNPACKED*/, T: () => PageActivity }
        ]);
    }
    create(value?: PartialMessage<ListPageActivityResponse>): ListPageActivityResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.activity = [];
        if (value !== undefined)
            reflectionMergePartial<ListPageActivityResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListPageActivityResponse): ListPageActivityResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationResponse pagination */ 1:
                    message.pagination = PaginationResponse.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* repeated resources.wiki.PageActivity activity */ 2:
                    message.activity.push(PageActivity.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: ListPageActivityResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationResponse pagination = 1; */
        if (message.pagination)
            PaginationResponse.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.wiki.PageActivity activity = 2; */
        for (let i = 0; i < message.activity.length; i++)
            PageActivity.internalBinaryWrite(message.activity[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.wiki.ListPageActivityResponse
 */
export const ListPageActivityResponse = new ListPageActivityResponse$Type();
/**
 * @generated ServiceType for protobuf service services.wiki.WikiService
 */
export const WikiService = new ServiceType("services.wiki.WikiService", [
    { name: "ListPages", options: {}, I: ListPagesRequest, O: ListPagesResponse },
    { name: "GetPage", options: {}, I: GetPageRequest, O: GetPageResponse },
    { name: "CreatePage", options: {}, I: CreatePageRequest, O: CreatePageResponse },
    { name: "UpdatePage", options: {}, I: UpdatePageRequest, O: UpdatePageResponse },
    { name: "DeletePage", options: {}, I: DeletePageRequest, O: DeletePageResponse },
    { name: "ListPageActivity", options: {}, I: ListPageActivityRequest, O: ListPageActivityResponse },
    { name: "UploadFile", clientStreaming: true, options: {}, I: UploadFileRequest, O: UploadFileResponse }
]);
