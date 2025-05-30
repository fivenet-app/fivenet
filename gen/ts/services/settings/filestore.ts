// @generated by protobuf-ts 2.10.0 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/settings/filestore.proto" (package "services.settings", syntax proto3)
// @ts-nocheck
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
import { File } from "../../resources/filestore/file";
import { FileInfo } from "../../resources/filestore/file";
import { PaginationResponse } from "../../resources/common/database/database";
import { PaginationRequest } from "../../resources/common/database/database";
/**
 * @generated from protobuf message services.settings.ListFilesRequest
 */
export interface ListFilesRequest {
    /**
     * @generated from protobuf field: resources.common.database.PaginationRequest pagination = 1;
     */
    pagination?: PaginationRequest;
    /**
     * @generated from protobuf field: optional string path = 2;
     */
    path?: string;
}
/**
 * @generated from protobuf message services.settings.ListFilesResponse
 */
export interface ListFilesResponse {
    /**
     * @generated from protobuf field: resources.common.database.PaginationResponse pagination = 1;
     */
    pagination?: PaginationResponse;
    /**
     * @generated from protobuf field: repeated resources.filestore.FileInfo files = 2;
     */
    files: FileInfo[];
}
/**
 * @generated from protobuf message services.settings.UploadFileRequest
 */
export interface UploadFileRequest {
    /**
     * @generated from protobuf field: string prefix = 1;
     */
    prefix: string;
    /**
     * @generated from protobuf field: string name = 2;
     */
    name: string;
    /**
     * @generated from protobuf field: resources.filestore.File file = 3;
     */
    file?: File;
}
/**
 * @generated from protobuf message services.settings.UploadFileResponse
 */
export interface UploadFileResponse {
    /**
     * @generated from protobuf field: resources.filestore.FileInfo file = 1;
     */
    file?: FileInfo;
}
/**
 * @generated from protobuf message services.settings.DeleteFileRequest
 */
export interface DeleteFileRequest {
    /**
     * @generated from protobuf field: string path = 1;
     */
    path: string;
}
/**
 * @generated from protobuf message services.settings.DeleteFileResponse
 */
export interface DeleteFileResponse {
}
// @generated message type with reflection information, may provide speed optimized methods
class ListFilesRequest$Type extends MessageType<ListFilesRequest> {
    constructor() {
        super("services.settings.ListFilesRequest", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationRequest, options: { "validate.rules": { message: { required: true } } } },
            { no: 2, name: "path", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "128" } } } }
        ]);
    }
    create(value?: PartialMessage<ListFilesRequest>): ListFilesRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<ListFilesRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListFilesRequest): ListFilesRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationRequest pagination */ 1:
                    message.pagination = PaginationRequest.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* optional string path */ 2:
                    message.path = reader.string();
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
    internalBinaryWrite(message: ListFilesRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationRequest pagination = 1; */
        if (message.pagination)
            PaginationRequest.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* optional string path = 2; */
        if (message.path !== undefined)
            writer.tag(2, WireType.LengthDelimited).string(message.path);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.settings.ListFilesRequest
 */
export const ListFilesRequest = new ListFilesRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ListFilesResponse$Type extends MessageType<ListFilesResponse> {
    constructor() {
        super("services.settings.ListFilesResponse", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationResponse },
            { no: 2, name: "files", kind: "message", repeat: 2 /*RepeatType.UNPACKED*/, T: () => FileInfo }
        ]);
    }
    create(value?: PartialMessage<ListFilesResponse>): ListFilesResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.files = [];
        if (value !== undefined)
            reflectionMergePartial<ListFilesResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListFilesResponse): ListFilesResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationResponse pagination */ 1:
                    message.pagination = PaginationResponse.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* repeated resources.filestore.FileInfo files */ 2:
                    message.files.push(FileInfo.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: ListFilesResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationResponse pagination = 1; */
        if (message.pagination)
            PaginationResponse.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.filestore.FileInfo files = 2; */
        for (let i = 0; i < message.files.length; i++)
            FileInfo.internalBinaryWrite(message.files[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.settings.ListFilesResponse
 */
export const ListFilesResponse = new ListFilesResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UploadFileRequest$Type extends MessageType<UploadFileRequest> {
    constructor() {
        super("services.settings.UploadFileRequest", [
            { no: 1, name: "prefix", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "file", kind: "message", T: () => File, options: { "validate.rules": { message: { required: true } } } }
        ]);
    }
    create(value?: PartialMessage<UploadFileRequest>): UploadFileRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.prefix = "";
        message.name = "";
        if (value !== undefined)
            reflectionMergePartial<UploadFileRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UploadFileRequest): UploadFileRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string prefix */ 1:
                    message.prefix = reader.string();
                    break;
                case /* string name */ 2:
                    message.name = reader.string();
                    break;
                case /* resources.filestore.File file */ 3:
                    message.file = File.internalBinaryRead(reader, reader.uint32(), options, message.file);
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
    internalBinaryWrite(message: UploadFileRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string prefix = 1; */
        if (message.prefix !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.prefix);
        /* string name = 2; */
        if (message.name !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.name);
        /* resources.filestore.File file = 3; */
        if (message.file)
            File.internalBinaryWrite(message.file, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.settings.UploadFileRequest
 */
export const UploadFileRequest = new UploadFileRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UploadFileResponse$Type extends MessageType<UploadFileResponse> {
    constructor() {
        super("services.settings.UploadFileResponse", [
            { no: 1, name: "file", kind: "message", T: () => FileInfo }
        ]);
    }
    create(value?: PartialMessage<UploadFileResponse>): UploadFileResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<UploadFileResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UploadFileResponse): UploadFileResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.filestore.FileInfo file */ 1:
                    message.file = FileInfo.internalBinaryRead(reader, reader.uint32(), options, message.file);
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
    internalBinaryWrite(message: UploadFileResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.filestore.FileInfo file = 1; */
        if (message.file)
            FileInfo.internalBinaryWrite(message.file, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.settings.UploadFileResponse
 */
export const UploadFileResponse = new UploadFileResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteFileRequest$Type extends MessageType<DeleteFileRequest> {
    constructor() {
        super("services.settings.DeleteFileRequest", [
            { no: 1, name: "path", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<DeleteFileRequest>): DeleteFileRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.path = "";
        if (value !== undefined)
            reflectionMergePartial<DeleteFileRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteFileRequest): DeleteFileRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string path */ 1:
                    message.path = reader.string();
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
    internalBinaryWrite(message: DeleteFileRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string path = 1; */
        if (message.path !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.path);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.settings.DeleteFileRequest
 */
export const DeleteFileRequest = new DeleteFileRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteFileResponse$Type extends MessageType<DeleteFileResponse> {
    constructor() {
        super("services.settings.DeleteFileResponse", []);
    }
    create(value?: PartialMessage<DeleteFileResponse>): DeleteFileResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<DeleteFileResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteFileResponse): DeleteFileResponse {
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
    internalBinaryWrite(message: DeleteFileResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.settings.DeleteFileResponse
 */
export const DeleteFileResponse = new DeleteFileResponse$Type();
/**
 * @generated ServiceType for protobuf service services.settings.FilestoreService
 */
export const FilestoreService = new ServiceType("services.settings.FilestoreService", [
    { name: "ListFiles", options: {}, I: ListFilesRequest, O: ListFilesResponse },
    { name: "UploadFile", options: {}, I: UploadFileRequest, O: UploadFileResponse },
    { name: "DeleteFile", options: {}, I: DeleteFileRequest, O: DeleteFileResponse }
]);
