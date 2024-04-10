// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/rector/laws.proto" (package "services.rector", syntax proto3)
// tslint:disable
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
import { Law } from "../../resources/laws/laws";
import { LawBook } from "../../resources/laws/laws";
/**
 * @generated from protobuf message services.rector.CreateOrUpdateLawBookRequest
 */
export interface CreateOrUpdateLawBookRequest {
    /**
     * @generated from protobuf field: resources.laws.LawBook lawBook = 1;
     */
    lawBook?: LawBook;
}
/**
 * @generated from protobuf message services.rector.CreateOrUpdateLawBookResponse
 */
export interface CreateOrUpdateLawBookResponse {
    /**
     * @generated from protobuf field: resources.laws.LawBook lawBook = 1;
     */
    lawBook?: LawBook;
}
/**
 * @generated from protobuf message services.rector.DeleteLawBookRequest
 */
export interface DeleteLawBookRequest {
    /**
     * @generated from protobuf field: uint64 id = 1 [jstype = JS_STRING];
     */
    id: string;
}
/**
 * @generated from protobuf message services.rector.DeleteLawBookResponse
 */
export interface DeleteLawBookResponse {
}
/**
 * @generated from protobuf message services.rector.CreateOrUpdateLawRequest
 */
export interface CreateOrUpdateLawRequest {
    /**
     * @generated from protobuf field: resources.laws.Law law = 1;
     */
    law?: Law;
}
/**
 * @generated from protobuf message services.rector.CreateOrUpdateLawResponse
 */
export interface CreateOrUpdateLawResponse {
    /**
     * @generated from protobuf field: resources.laws.Law law = 1;
     */
    law?: Law;
}
/**
 * @generated from protobuf message services.rector.DeleteLawRequest
 */
export interface DeleteLawRequest {
    /**
     * @generated from protobuf field: uint64 id = 1 [jstype = JS_STRING];
     */
    id: string;
}
/**
 * @generated from protobuf message services.rector.DeleteLawResponse
 */
export interface DeleteLawResponse {
}
// @generated message type with reflection information, may provide speed optimized methods
class CreateOrUpdateLawBookRequest$Type extends MessageType<CreateOrUpdateLawBookRequest> {
    constructor() {
        super("services.rector.CreateOrUpdateLawBookRequest", [
            { no: 1, name: "lawBook", kind: "message", T: () => LawBook, options: { "validate.rules": { message: { required: true } } } }
        ]);
    }
    create(value?: PartialMessage<CreateOrUpdateLawBookRequest>): CreateOrUpdateLawBookRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<CreateOrUpdateLawBookRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CreateOrUpdateLawBookRequest): CreateOrUpdateLawBookRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.laws.LawBook lawBook */ 1:
                    message.lawBook = LawBook.internalBinaryRead(reader, reader.uint32(), options, message.lawBook);
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
    internalBinaryWrite(message: CreateOrUpdateLawBookRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.laws.LawBook lawBook = 1; */
        if (message.lawBook)
            LawBook.internalBinaryWrite(message.lawBook, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.rector.CreateOrUpdateLawBookRequest
 */
export const CreateOrUpdateLawBookRequest = new CreateOrUpdateLawBookRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateOrUpdateLawBookResponse$Type extends MessageType<CreateOrUpdateLawBookResponse> {
    constructor() {
        super("services.rector.CreateOrUpdateLawBookResponse", [
            { no: 1, name: "lawBook", kind: "message", T: () => LawBook }
        ]);
    }
    create(value?: PartialMessage<CreateOrUpdateLawBookResponse>): CreateOrUpdateLawBookResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<CreateOrUpdateLawBookResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CreateOrUpdateLawBookResponse): CreateOrUpdateLawBookResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.laws.LawBook lawBook */ 1:
                    message.lawBook = LawBook.internalBinaryRead(reader, reader.uint32(), options, message.lawBook);
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
    internalBinaryWrite(message: CreateOrUpdateLawBookResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.laws.LawBook lawBook = 1; */
        if (message.lawBook)
            LawBook.internalBinaryWrite(message.lawBook, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.rector.CreateOrUpdateLawBookResponse
 */
export const CreateOrUpdateLawBookResponse = new CreateOrUpdateLawBookResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteLawBookRequest$Type extends MessageType<DeleteLawBookRequest> {
    constructor() {
        super("services.rector.DeleteLawBookRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ }
        ]);
    }
    create(value?: PartialMessage<DeleteLawBookRequest>): DeleteLawBookRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = "0";
        if (value !== undefined)
            reflectionMergePartial<DeleteLawBookRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteLawBookRequest): DeleteLawBookRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id = 1 [jstype = JS_STRING];*/ 1:
                    message.id = reader.uint64().toString();
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
    internalBinaryWrite(message: DeleteLawBookRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1 [jstype = JS_STRING]; */
        if (message.id !== "0")
            writer.tag(1, WireType.Varint).uint64(message.id);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.rector.DeleteLawBookRequest
 */
export const DeleteLawBookRequest = new DeleteLawBookRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteLawBookResponse$Type extends MessageType<DeleteLawBookResponse> {
    constructor() {
        super("services.rector.DeleteLawBookResponse", []);
    }
    create(value?: PartialMessage<DeleteLawBookResponse>): DeleteLawBookResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<DeleteLawBookResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteLawBookResponse): DeleteLawBookResponse {
        return target ?? this.create();
    }
    internalBinaryWrite(message: DeleteLawBookResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.rector.DeleteLawBookResponse
 */
export const DeleteLawBookResponse = new DeleteLawBookResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateOrUpdateLawRequest$Type extends MessageType<CreateOrUpdateLawRequest> {
    constructor() {
        super("services.rector.CreateOrUpdateLawRequest", [
            { no: 1, name: "law", kind: "message", T: () => Law, options: { "validate.rules": { message: { required: true } } } }
        ]);
    }
    create(value?: PartialMessage<CreateOrUpdateLawRequest>): CreateOrUpdateLawRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<CreateOrUpdateLawRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CreateOrUpdateLawRequest): CreateOrUpdateLawRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.laws.Law law */ 1:
                    message.law = Law.internalBinaryRead(reader, reader.uint32(), options, message.law);
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
    internalBinaryWrite(message: CreateOrUpdateLawRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.laws.Law law = 1; */
        if (message.law)
            Law.internalBinaryWrite(message.law, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.rector.CreateOrUpdateLawRequest
 */
export const CreateOrUpdateLawRequest = new CreateOrUpdateLawRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateOrUpdateLawResponse$Type extends MessageType<CreateOrUpdateLawResponse> {
    constructor() {
        super("services.rector.CreateOrUpdateLawResponse", [
            { no: 1, name: "law", kind: "message", T: () => Law }
        ]);
    }
    create(value?: PartialMessage<CreateOrUpdateLawResponse>): CreateOrUpdateLawResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<CreateOrUpdateLawResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CreateOrUpdateLawResponse): CreateOrUpdateLawResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.laws.Law law */ 1:
                    message.law = Law.internalBinaryRead(reader, reader.uint32(), options, message.law);
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
    internalBinaryWrite(message: CreateOrUpdateLawResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.laws.Law law = 1; */
        if (message.law)
            Law.internalBinaryWrite(message.law, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.rector.CreateOrUpdateLawResponse
 */
export const CreateOrUpdateLawResponse = new CreateOrUpdateLawResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteLawRequest$Type extends MessageType<DeleteLawRequest> {
    constructor() {
        super("services.rector.DeleteLawRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ }
        ]);
    }
    create(value?: PartialMessage<DeleteLawRequest>): DeleteLawRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = "0";
        if (value !== undefined)
            reflectionMergePartial<DeleteLawRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteLawRequest): DeleteLawRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id = 1 [jstype = JS_STRING];*/ 1:
                    message.id = reader.uint64().toString();
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
    internalBinaryWrite(message: DeleteLawRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1 [jstype = JS_STRING]; */
        if (message.id !== "0")
            writer.tag(1, WireType.Varint).uint64(message.id);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.rector.DeleteLawRequest
 */
export const DeleteLawRequest = new DeleteLawRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteLawResponse$Type extends MessageType<DeleteLawResponse> {
    constructor() {
        super("services.rector.DeleteLawResponse", []);
    }
    create(value?: PartialMessage<DeleteLawResponse>): DeleteLawResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<DeleteLawResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteLawResponse): DeleteLawResponse {
        return target ?? this.create();
    }
    internalBinaryWrite(message: DeleteLawResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.rector.DeleteLawResponse
 */
export const DeleteLawResponse = new DeleteLawResponse$Type();
/**
 * @generated ServiceType for protobuf service services.rector.RectorLawsService
 */
export const RectorLawsService = new ServiceType("services.rector.RectorLawsService", [
    { name: "CreateOrUpdateLawBook", options: {}, I: CreateOrUpdateLawBookRequest, O: CreateOrUpdateLawBookResponse },
    { name: "DeleteLawBook", options: {}, I: DeleteLawBookRequest, O: DeleteLawBookResponse },
    { name: "CreateOrUpdateLaw", options: {}, I: CreateOrUpdateLawRequest, O: CreateOrUpdateLawResponse },
    { name: "DeleteLaw", options: {}, I: DeleteLawRequest, O: DeleteLawResponse }
]);
