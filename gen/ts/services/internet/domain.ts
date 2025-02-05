// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/internet/domain.proto" (package "services.internet", syntax proto3)
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
import { Domain } from "../../resources/internet/domain";
import { PaginationResponse } from "../../resources/common/database/database";
import { PaginationRequest } from "../../resources/common/database/database";
import { TLD } from "../../resources/internet/domain";
/**
 * @generated from protobuf message services.internet.ListTLDsRequest
 */
export interface ListTLDsRequest {
    /**
     * @generated from protobuf field: optional bool internal = 1;
     */
    internal?: boolean;
}
/**
 * @generated from protobuf message services.internet.ListTLDsResponse
 */
export interface ListTLDsResponse {
    /**
     * @generated from protobuf field: repeated resources.internet.TLD tlds = 1;
     */
    tlds: TLD[];
}
/**
 * @generated from protobuf message services.internet.CheckDomainAvailabilityRequest
 */
export interface CheckDomainAvailabilityRequest {
    /**
     * @generated from protobuf field: uint64 tld_id = 1;
     */
    tldId: number;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: string name = 2;
     */
    name: string;
}
/**
 * @generated from protobuf message services.internet.CheckDomainAvailabilityResponse
 */
export interface CheckDomainAvailabilityResponse {
    /**
     * @generated from protobuf field: bool available = 1;
     */
    available: boolean;
    /**
     * @generated from protobuf field: optional bool transferable = 2;
     */
    transferable?: boolean;
}
/**
 * @generated from protobuf message services.internet.ListDomainsRequest
 */
export interface ListDomainsRequest {
    /**
     * @generated from protobuf field: resources.common.database.PaginationRequest pagination = 1;
     */
    pagination?: PaginationRequest;
}
/**
 * @generated from protobuf message services.internet.ListDomainsResponse
 */
export interface ListDomainsResponse {
    /**
     * @generated from protobuf field: resources.common.database.PaginationResponse pagination = 1;
     */
    pagination?: PaginationResponse;
    /**
     * @generated from protobuf field: repeated resources.internet.Domain domains = 2;
     */
    domains: Domain[];
}
/**
 * @generated from protobuf message services.internet.RegisterDomainRequest
 */
export interface RegisterDomainRequest {
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: string name = 1;
     */
    name: string;
    /**
     * In case a domain will be transfered
     *
     * @generated from protobuf field: optional string transfer_code = 2;
     */
    transferCode?: string;
}
/**
 * @generated from protobuf message services.internet.RegisterDomainResponse
 */
export interface RegisterDomainResponse {
    /**
     * @generated from protobuf field: resources.internet.Domain domain = 1;
     */
    domain?: Domain;
}
/**
 * @generated from protobuf message services.internet.UpdateDomainRequest
 */
export interface UpdateDomainRequest {
    /**
     * @generated from protobuf field: resources.internet.Domain domain = 1;
     */
    domain?: Domain;
}
/**
 * @generated from protobuf message services.internet.UpdateDomainResponse
 */
export interface UpdateDomainResponse {
    /**
     * @generated from protobuf field: resources.internet.Domain domain = 1;
     */
    domain?: Domain;
}
// @generated message type with reflection information, may provide speed optimized methods
class ListTLDsRequest$Type extends MessageType<ListTLDsRequest> {
    constructor() {
        super("services.internet.ListTLDsRequest", [
            { no: 1, name: "internal", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<ListTLDsRequest>): ListTLDsRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<ListTLDsRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListTLDsRequest): ListTLDsRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* optional bool internal */ 1:
                    message.internal = reader.bool();
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
    internalBinaryWrite(message: ListTLDsRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* optional bool internal = 1; */
        if (message.internal !== undefined)
            writer.tag(1, WireType.Varint).bool(message.internal);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.internet.ListTLDsRequest
 */
export const ListTLDsRequest = new ListTLDsRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ListTLDsResponse$Type extends MessageType<ListTLDsResponse> {
    constructor() {
        super("services.internet.ListTLDsResponse", [
            { no: 1, name: "tlds", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => TLD }
        ]);
    }
    create(value?: PartialMessage<ListTLDsResponse>): ListTLDsResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.tlds = [];
        if (value !== undefined)
            reflectionMergePartial<ListTLDsResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListTLDsResponse): ListTLDsResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.internet.TLD tlds */ 1:
                    message.tlds.push(TLD.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: ListTLDsResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.internet.TLD tlds = 1; */
        for (let i = 0; i < message.tlds.length; i++)
            TLD.internalBinaryWrite(message.tlds[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.internet.ListTLDsResponse
 */
export const ListTLDsResponse = new ListTLDsResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CheckDomainAvailabilityRequest$Type extends MessageType<CheckDomainAvailabilityRequest> {
    constructor() {
        super("services.internet.CheckDomainAvailabilityRequest", [
            { no: 1, name: "tld_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "3", maxLen: "60" } } } }
        ]);
    }
    create(value?: PartialMessage<CheckDomainAvailabilityRequest>): CheckDomainAvailabilityRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.tldId = 0;
        message.name = "";
        if (value !== undefined)
            reflectionMergePartial<CheckDomainAvailabilityRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CheckDomainAvailabilityRequest): CheckDomainAvailabilityRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 tld_id */ 1:
                    message.tldId = reader.uint64().toNumber();
                    break;
                case /* string name */ 2:
                    message.name = reader.string();
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
    internalBinaryWrite(message: CheckDomainAvailabilityRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 tld_id = 1; */
        if (message.tldId !== 0)
            writer.tag(1, WireType.Varint).uint64(message.tldId);
        /* string name = 2; */
        if (message.name !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.name);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.internet.CheckDomainAvailabilityRequest
 */
export const CheckDomainAvailabilityRequest = new CheckDomainAvailabilityRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CheckDomainAvailabilityResponse$Type extends MessageType<CheckDomainAvailabilityResponse> {
    constructor() {
        super("services.internet.CheckDomainAvailabilityResponse", [
            { no: 1, name: "available", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 2, name: "transferable", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<CheckDomainAvailabilityResponse>): CheckDomainAvailabilityResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.available = false;
        if (value !== undefined)
            reflectionMergePartial<CheckDomainAvailabilityResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CheckDomainAvailabilityResponse): CheckDomainAvailabilityResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* bool available */ 1:
                    message.available = reader.bool();
                    break;
                case /* optional bool transferable */ 2:
                    message.transferable = reader.bool();
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
    internalBinaryWrite(message: CheckDomainAvailabilityResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* bool available = 1; */
        if (message.available !== false)
            writer.tag(1, WireType.Varint).bool(message.available);
        /* optional bool transferable = 2; */
        if (message.transferable !== undefined)
            writer.tag(2, WireType.Varint).bool(message.transferable);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.internet.CheckDomainAvailabilityResponse
 */
export const CheckDomainAvailabilityResponse = new CheckDomainAvailabilityResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ListDomainsRequest$Type extends MessageType<ListDomainsRequest> {
    constructor() {
        super("services.internet.ListDomainsRequest", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationRequest, options: { "validate.rules": { message: { required: true } } } }
        ]);
    }
    create(value?: PartialMessage<ListDomainsRequest>): ListDomainsRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<ListDomainsRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListDomainsRequest): ListDomainsRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationRequest pagination */ 1:
                    message.pagination = PaginationRequest.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
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
    internalBinaryWrite(message: ListDomainsRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationRequest pagination = 1; */
        if (message.pagination)
            PaginationRequest.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.internet.ListDomainsRequest
 */
export const ListDomainsRequest = new ListDomainsRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ListDomainsResponse$Type extends MessageType<ListDomainsResponse> {
    constructor() {
        super("services.internet.ListDomainsResponse", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationResponse, options: { "validate.rules": { message: { required: true } } } },
            { no: 2, name: "domains", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Domain }
        ]);
    }
    create(value?: PartialMessage<ListDomainsResponse>): ListDomainsResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.domains = [];
        if (value !== undefined)
            reflectionMergePartial<ListDomainsResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListDomainsResponse): ListDomainsResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationResponse pagination */ 1:
                    message.pagination = PaginationResponse.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* repeated resources.internet.Domain domains */ 2:
                    message.domains.push(Domain.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: ListDomainsResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationResponse pagination = 1; */
        if (message.pagination)
            PaginationResponse.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.internet.Domain domains = 2; */
        for (let i = 0; i < message.domains.length; i++)
            Domain.internalBinaryWrite(message.domains[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.internet.ListDomainsResponse
 */
export const ListDomainsResponse = new ListDomainsResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class RegisterDomainRequest$Type extends MessageType<RegisterDomainRequest> {
    constructor() {
        super("services.internet.RegisterDomainRequest", [
            { no: 1, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "3", maxLen: "60" } } } },
            { no: 2, name: "transfer_code", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { len: "10", pattern: "^[0-9A-Z]{6}$" } } } }
        ]);
    }
    create(value?: PartialMessage<RegisterDomainRequest>): RegisterDomainRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.name = "";
        if (value !== undefined)
            reflectionMergePartial<RegisterDomainRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: RegisterDomainRequest): RegisterDomainRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string name */ 1:
                    message.name = reader.string();
                    break;
                case /* optional string transfer_code */ 2:
                    message.transferCode = reader.string();
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
    internalBinaryWrite(message: RegisterDomainRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string name = 1; */
        if (message.name !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.name);
        /* optional string transfer_code = 2; */
        if (message.transferCode !== undefined)
            writer.tag(2, WireType.LengthDelimited).string(message.transferCode);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.internet.RegisterDomainRequest
 */
export const RegisterDomainRequest = new RegisterDomainRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class RegisterDomainResponse$Type extends MessageType<RegisterDomainResponse> {
    constructor() {
        super("services.internet.RegisterDomainResponse", [
            { no: 1, name: "domain", kind: "message", T: () => Domain }
        ]);
    }
    create(value?: PartialMessage<RegisterDomainResponse>): RegisterDomainResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<RegisterDomainResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: RegisterDomainResponse): RegisterDomainResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.internet.Domain domain */ 1:
                    message.domain = Domain.internalBinaryRead(reader, reader.uint32(), options, message.domain);
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
    internalBinaryWrite(message: RegisterDomainResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.internet.Domain domain = 1; */
        if (message.domain)
            Domain.internalBinaryWrite(message.domain, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.internet.RegisterDomainResponse
 */
export const RegisterDomainResponse = new RegisterDomainResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UpdateDomainRequest$Type extends MessageType<UpdateDomainRequest> {
    constructor() {
        super("services.internet.UpdateDomainRequest", [
            { no: 1, name: "domain", kind: "message", T: () => Domain }
        ]);
    }
    create(value?: PartialMessage<UpdateDomainRequest>): UpdateDomainRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<UpdateDomainRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UpdateDomainRequest): UpdateDomainRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.internet.Domain domain */ 1:
                    message.domain = Domain.internalBinaryRead(reader, reader.uint32(), options, message.domain);
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
    internalBinaryWrite(message: UpdateDomainRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.internet.Domain domain = 1; */
        if (message.domain)
            Domain.internalBinaryWrite(message.domain, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.internet.UpdateDomainRequest
 */
export const UpdateDomainRequest = new UpdateDomainRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UpdateDomainResponse$Type extends MessageType<UpdateDomainResponse> {
    constructor() {
        super("services.internet.UpdateDomainResponse", [
            { no: 1, name: "domain", kind: "message", T: () => Domain }
        ]);
    }
    create(value?: PartialMessage<UpdateDomainResponse>): UpdateDomainResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<UpdateDomainResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UpdateDomainResponse): UpdateDomainResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.internet.Domain domain */ 1:
                    message.domain = Domain.internalBinaryRead(reader, reader.uint32(), options, message.domain);
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
    internalBinaryWrite(message: UpdateDomainResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.internet.Domain domain = 1; */
        if (message.domain)
            Domain.internalBinaryWrite(message.domain, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.internet.UpdateDomainResponse
 */
export const UpdateDomainResponse = new UpdateDomainResponse$Type();
/**
 * @generated ServiceType for protobuf service services.internet.DomainService
 */
export const DomainService = new ServiceType("services.internet.DomainService", [
    { name: "ListTLDs", options: {}, I: ListTLDsRequest, O: ListTLDsResponse },
    { name: "CheckDomainAvailability", options: {}, I: CheckDomainAvailabilityRequest, O: CheckDomainAvailabilityResponse },
    { name: "ListDomains", options: {}, I: ListDomainsRequest, O: ListDomainsResponse },
    { name: "RegisterDomain", options: {}, I: RegisterDomainRequest, O: RegisterDomainResponse },
    { name: "UpdateDomain", options: {}, I: UpdateDomainRequest, O: UpdateDomainResponse }
]);
