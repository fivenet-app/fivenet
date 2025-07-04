// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "services/jobs/conduct.proto" (package "services.jobs", syntax proto3)
// tslint:disable
// @ts-nocheck
import { ServiceType } from "@protobuf-ts/runtime-rpc";
import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import { UnknownFieldHandler } from "@protobuf-ts/runtime";
import { WireType } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { reflectionMergePartial } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
import { ConductEntry } from "../../resources/jobs/conduct";
import { PaginationResponse } from "../../resources/common/database/database";
import { ConductType } from "../../resources/jobs/conduct";
import { Sort } from "../../resources/common/database/database";
import { PaginationRequest } from "../../resources/common/database/database";
// Conduct Register

/**
 * @generated from protobuf message services.jobs.ListConductEntriesRequest
 */
export interface ListConductEntriesRequest {
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
     * @generated from protobuf field: repeated resources.jobs.ConductType types = 3
     */
    types: ConductType[];
    /**
     * @generated from protobuf field: optional bool show_expired = 4
     */
    showExpired?: boolean;
    /**
     * @generated from protobuf field: repeated int32 user_ids = 5
     */
    userIds: number[];
    /**
     * @generated from protobuf field: repeated uint64 ids = 6
     */
    ids: number[];
}
/**
 * @generated from protobuf message services.jobs.ListConductEntriesResponse
 */
export interface ListConductEntriesResponse {
    /**
     * @generated from protobuf field: resources.common.database.PaginationResponse pagination = 1
     */
    pagination?: PaginationResponse;
    /**
     * @generated from protobuf field: repeated resources.jobs.ConductEntry entries = 2
     */
    entries: ConductEntry[];
}
/**
 * @generated from protobuf message services.jobs.CreateConductEntryRequest
 */
export interface CreateConductEntryRequest {
    /**
     * @generated from protobuf field: resources.jobs.ConductEntry entry = 1
     */
    entry?: ConductEntry;
}
/**
 * @generated from protobuf message services.jobs.CreateConductEntryResponse
 */
export interface CreateConductEntryResponse {
    /**
     * @generated from protobuf field: resources.jobs.ConductEntry entry = 1
     */
    entry?: ConductEntry;
}
/**
 * @generated from protobuf message services.jobs.UpdateConductEntryRequest
 */
export interface UpdateConductEntryRequest {
    /**
     * @generated from protobuf field: resources.jobs.ConductEntry entry = 1
     */
    entry?: ConductEntry;
}
/**
 * @generated from protobuf message services.jobs.UpdateConductEntryResponse
 */
export interface UpdateConductEntryResponse {
    /**
     * @generated from protobuf field: resources.jobs.ConductEntry entry = 1
     */
    entry?: ConductEntry;
}
/**
 * @generated from protobuf message services.jobs.DeleteConductEntryRequest
 */
export interface DeleteConductEntryRequest {
    /**
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
}
/**
 * @generated from protobuf message services.jobs.DeleteConductEntryResponse
 */
export interface DeleteConductEntryResponse {
}
// @generated message type with reflection information, may provide speed optimized methods
class ListConductEntriesRequest$Type extends MessageType<ListConductEntriesRequest> {
    constructor() {
        super("services.jobs.ListConductEntriesRequest", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationRequest, options: { "buf.validate.field": { required: true } } },
            { no: 2, name: "sort", kind: "message", T: () => Sort },
            { no: 3, name: "types", kind: "enum", repeat: 1 /*RepeatType.PACKED*/, T: () => ["resources.jobs.ConductType", ConductType, "CONDUCT_TYPE_"] },
            { no: 4, name: "show_expired", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ },
            { no: 5, name: "user_ids", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 5 /*ScalarType.INT32*/ },
            { no: 6, name: "ids", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ }
        ]);
    }
    create(value?: PartialMessage<ListConductEntriesRequest>): ListConductEntriesRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.types = [];
        message.userIds = [];
        message.ids = [];
        if (value !== undefined)
            reflectionMergePartial<ListConductEntriesRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListConductEntriesRequest): ListConductEntriesRequest {
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
                case /* repeated resources.jobs.ConductType types */ 3:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.types.push(reader.int32());
                    else
                        message.types.push(reader.int32());
                    break;
                case /* optional bool show_expired */ 4:
                    message.showExpired = reader.bool();
                    break;
                case /* repeated int32 user_ids */ 5:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.userIds.push(reader.int32());
                    else
                        message.userIds.push(reader.int32());
                    break;
                case /* repeated uint64 ids */ 6:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.ids.push(reader.uint64().toNumber());
                    else
                        message.ids.push(reader.uint64().toNumber());
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
    internalBinaryWrite(message: ListConductEntriesRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationRequest pagination = 1; */
        if (message.pagination)
            PaginationRequest.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.common.database.Sort sort = 2; */
        if (message.sort)
            Sort.internalBinaryWrite(message.sort, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.jobs.ConductType types = 3; */
        if (message.types.length) {
            writer.tag(3, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.types.length; i++)
                writer.int32(message.types[i]);
            writer.join();
        }
        /* optional bool show_expired = 4; */
        if (message.showExpired !== undefined)
            writer.tag(4, WireType.Varint).bool(message.showExpired);
        /* repeated int32 user_ids = 5; */
        if (message.userIds.length) {
            writer.tag(5, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.userIds.length; i++)
                writer.int32(message.userIds[i]);
            writer.join();
        }
        /* repeated uint64 ids = 6; */
        if (message.ids.length) {
            writer.tag(6, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.ids.length; i++)
                writer.uint64(message.ids[i]);
            writer.join();
        }
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.jobs.ListConductEntriesRequest
 */
export const ListConductEntriesRequest = new ListConductEntriesRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ListConductEntriesResponse$Type extends MessageType<ListConductEntriesResponse> {
    constructor() {
        super("services.jobs.ListConductEntriesResponse", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationResponse },
            { no: 2, name: "entries", kind: "message", repeat: 2 /*RepeatType.UNPACKED*/, T: () => ConductEntry }
        ]);
    }
    create(value?: PartialMessage<ListConductEntriesResponse>): ListConductEntriesResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.entries = [];
        if (value !== undefined)
            reflectionMergePartial<ListConductEntriesResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListConductEntriesResponse): ListConductEntriesResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationResponse pagination */ 1:
                    message.pagination = PaginationResponse.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* repeated resources.jobs.ConductEntry entries */ 2:
                    message.entries.push(ConductEntry.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: ListConductEntriesResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationResponse pagination = 1; */
        if (message.pagination)
            PaginationResponse.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.jobs.ConductEntry entries = 2; */
        for (let i = 0; i < message.entries.length; i++)
            ConductEntry.internalBinaryWrite(message.entries[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.jobs.ListConductEntriesResponse
 */
export const ListConductEntriesResponse = new ListConductEntriesResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateConductEntryRequest$Type extends MessageType<CreateConductEntryRequest> {
    constructor() {
        super("services.jobs.CreateConductEntryRequest", [
            { no: 1, name: "entry", kind: "message", T: () => ConductEntry, options: { "buf.validate.field": { required: true } } }
        ]);
    }
    create(value?: PartialMessage<CreateConductEntryRequest>): CreateConductEntryRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<CreateConductEntryRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CreateConductEntryRequest): CreateConductEntryRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.jobs.ConductEntry entry */ 1:
                    message.entry = ConductEntry.internalBinaryRead(reader, reader.uint32(), options, message.entry);
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
    internalBinaryWrite(message: CreateConductEntryRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.jobs.ConductEntry entry = 1; */
        if (message.entry)
            ConductEntry.internalBinaryWrite(message.entry, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.jobs.CreateConductEntryRequest
 */
export const CreateConductEntryRequest = new CreateConductEntryRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateConductEntryResponse$Type extends MessageType<CreateConductEntryResponse> {
    constructor() {
        super("services.jobs.CreateConductEntryResponse", [
            { no: 1, name: "entry", kind: "message", T: () => ConductEntry }
        ]);
    }
    create(value?: PartialMessage<CreateConductEntryResponse>): CreateConductEntryResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<CreateConductEntryResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CreateConductEntryResponse): CreateConductEntryResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.jobs.ConductEntry entry */ 1:
                    message.entry = ConductEntry.internalBinaryRead(reader, reader.uint32(), options, message.entry);
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
    internalBinaryWrite(message: CreateConductEntryResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.jobs.ConductEntry entry = 1; */
        if (message.entry)
            ConductEntry.internalBinaryWrite(message.entry, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.jobs.CreateConductEntryResponse
 */
export const CreateConductEntryResponse = new CreateConductEntryResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UpdateConductEntryRequest$Type extends MessageType<UpdateConductEntryRequest> {
    constructor() {
        super("services.jobs.UpdateConductEntryRequest", [
            { no: 1, name: "entry", kind: "message", T: () => ConductEntry, options: { "buf.validate.field": { required: true } } }
        ]);
    }
    create(value?: PartialMessage<UpdateConductEntryRequest>): UpdateConductEntryRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<UpdateConductEntryRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UpdateConductEntryRequest): UpdateConductEntryRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.jobs.ConductEntry entry */ 1:
                    message.entry = ConductEntry.internalBinaryRead(reader, reader.uint32(), options, message.entry);
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
    internalBinaryWrite(message: UpdateConductEntryRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.jobs.ConductEntry entry = 1; */
        if (message.entry)
            ConductEntry.internalBinaryWrite(message.entry, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.jobs.UpdateConductEntryRequest
 */
export const UpdateConductEntryRequest = new UpdateConductEntryRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UpdateConductEntryResponse$Type extends MessageType<UpdateConductEntryResponse> {
    constructor() {
        super("services.jobs.UpdateConductEntryResponse", [
            { no: 1, name: "entry", kind: "message", T: () => ConductEntry, options: { "buf.validate.field": { required: true } } }
        ]);
    }
    create(value?: PartialMessage<UpdateConductEntryResponse>): UpdateConductEntryResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<UpdateConductEntryResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UpdateConductEntryResponse): UpdateConductEntryResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.jobs.ConductEntry entry */ 1:
                    message.entry = ConductEntry.internalBinaryRead(reader, reader.uint32(), options, message.entry);
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
    internalBinaryWrite(message: UpdateConductEntryResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.jobs.ConductEntry entry = 1; */
        if (message.entry)
            ConductEntry.internalBinaryWrite(message.entry, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.jobs.UpdateConductEntryResponse
 */
export const UpdateConductEntryResponse = new UpdateConductEntryResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteConductEntryRequest$Type extends MessageType<DeleteConductEntryRequest> {
    constructor() {
        super("services.jobs.DeleteConductEntryRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ }
        ]);
    }
    create(value?: PartialMessage<DeleteConductEntryRequest>): DeleteConductEntryRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        if (value !== undefined)
            reflectionMergePartial<DeleteConductEntryRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteConductEntryRequest): DeleteConductEntryRequest {
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
    internalBinaryWrite(message: DeleteConductEntryRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
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
 * @generated MessageType for protobuf message services.jobs.DeleteConductEntryRequest
 */
export const DeleteConductEntryRequest = new DeleteConductEntryRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteConductEntryResponse$Type extends MessageType<DeleteConductEntryResponse> {
    constructor() {
        super("services.jobs.DeleteConductEntryResponse", []);
    }
    create(value?: PartialMessage<DeleteConductEntryResponse>): DeleteConductEntryResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<DeleteConductEntryResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteConductEntryResponse): DeleteConductEntryResponse {
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
    internalBinaryWrite(message: DeleteConductEntryResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.jobs.DeleteConductEntryResponse
 */
export const DeleteConductEntryResponse = new DeleteConductEntryResponse$Type();
/**
 * @generated ServiceType for protobuf service services.jobs.ConductService
 */
export const ConductService = new ServiceType("services.jobs.ConductService", [
    { name: "ListConductEntries", options: {}, I: ListConductEntriesRequest, O: ListConductEntriesResponse },
    { name: "CreateConductEntry", options: {}, I: CreateConductEntryRequest, O: CreateConductEntryResponse },
    { name: "UpdateConductEntry", options: {}, I: UpdateConductEntryRequest, O: UpdateConductEntryResponse },
    { name: "DeleteConductEntry", options: {}, I: DeleteConductEntryRequest, O: DeleteConductEntryResponse }
]);
