// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "services/settings/accounts.proto" (package "services.settings", syntax proto3)
// tslint:disable
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
import { Account } from "../../resources/accounts/accounts";
import { PaginationResponse } from "../../resources/common/database/database";
import { Sort } from "../../resources/common/database/database";
import { PaginationRequest } from "../../resources/common/database/database";
/**
 * @generated from protobuf message services.settings.ListAccountsRequest
 */
export interface ListAccountsRequest {
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
     * @generated from protobuf field: optional string license = 3
     */
    license?: string;
    /**
     * @generated from protobuf field: optional bool enabled = 4
     */
    enabled?: boolean;
    /**
     * @generated from protobuf field: optional string username = 5
     */
    username?: string;
    /**
     * @generated from protobuf field: optional string external_id = 6
     */
    externalId?: string;
}
/**
 * @generated from protobuf message services.settings.ListAccountsResponse
 */
export interface ListAccountsResponse {
    /**
     * @generated from protobuf field: resources.common.database.PaginationResponse pagination = 1
     */
    pagination?: PaginationResponse;
    /**
     * @generated from protobuf field: repeated resources.accounts.Account accounts = 2
     */
    accounts: Account[];
}
/**
 * @generated from protobuf message services.settings.UpdateAccountRequest
 */
export interface UpdateAccountRequest {
    /**
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
    /**
     * @generated from protobuf field: optional bool enabled = 2
     */
    enabled?: boolean;
    /**
     * @generated from protobuf field: optional int32 last_char = 3
     */
    lastChar?: number;
}
/**
 * @generated from protobuf message services.settings.UpdateAccountResponse
 */
export interface UpdateAccountResponse {
    /**
     * @generated from protobuf field: resources.accounts.Account account = 1
     */
    account?: Account;
}
/**
 * @generated from protobuf message services.settings.DisconnectOAuth2ConnectionRequest
 */
export interface DisconnectOAuth2ConnectionRequest {
    /**
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
    /**
     * @generated from protobuf field: string provider_name = 2
     */
    providerName: string;
}
/**
 * @generated from protobuf message services.settings.DisconnectOAuth2ConnectionResponse
 */
export interface DisconnectOAuth2ConnectionResponse {
}
/**
 * @generated from protobuf message services.settings.DeleteAccountRequest
 */
export interface DeleteAccountRequest {
    /**
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
}
/**
 * @generated from protobuf message services.settings.DeleteAccountResponse
 */
export interface DeleteAccountResponse {
}
// @generated message type with reflection information, may provide speed optimized methods
class ListAccountsRequest$Type extends MessageType<ListAccountsRequest> {
    constructor() {
        super("services.settings.ListAccountsRequest", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationRequest, options: { "buf.validate.field": { required: true } } },
            { no: 2, name: "sort", kind: "message", T: () => Sort },
            { no: 3, name: "license", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "64" } } } },
            { no: 4, name: "enabled", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ },
            { no: 5, name: "username", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "64" } } } },
            { no: 6, name: "external_id", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "128" } } } }
        ]);
    }
    create(value?: PartialMessage<ListAccountsRequest>): ListAccountsRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<ListAccountsRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListAccountsRequest): ListAccountsRequest {
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
                case /* optional string license */ 3:
                    message.license = reader.string();
                    break;
                case /* optional bool enabled */ 4:
                    message.enabled = reader.bool();
                    break;
                case /* optional string username */ 5:
                    message.username = reader.string();
                    break;
                case /* optional string external_id */ 6:
                    message.externalId = reader.string();
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
    internalBinaryWrite(message: ListAccountsRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationRequest pagination = 1; */
        if (message.pagination)
            PaginationRequest.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.common.database.Sort sort = 2; */
        if (message.sort)
            Sort.internalBinaryWrite(message.sort, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* optional string license = 3; */
        if (message.license !== undefined)
            writer.tag(3, WireType.LengthDelimited).string(message.license);
        /* optional bool enabled = 4; */
        if (message.enabled !== undefined)
            writer.tag(4, WireType.Varint).bool(message.enabled);
        /* optional string username = 5; */
        if (message.username !== undefined)
            writer.tag(5, WireType.LengthDelimited).string(message.username);
        /* optional string external_id = 6; */
        if (message.externalId !== undefined)
            writer.tag(6, WireType.LengthDelimited).string(message.externalId);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.settings.ListAccountsRequest
 */
export const ListAccountsRequest = new ListAccountsRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ListAccountsResponse$Type extends MessageType<ListAccountsResponse> {
    constructor() {
        super("services.settings.ListAccountsResponse", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationResponse },
            { no: 2, name: "accounts", kind: "message", repeat: 2 /*RepeatType.UNPACKED*/, T: () => Account }
        ]);
    }
    create(value?: PartialMessage<ListAccountsResponse>): ListAccountsResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.accounts = [];
        if (value !== undefined)
            reflectionMergePartial<ListAccountsResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListAccountsResponse): ListAccountsResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationResponse pagination */ 1:
                    message.pagination = PaginationResponse.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* repeated resources.accounts.Account accounts */ 2:
                    message.accounts.push(Account.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: ListAccountsResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationResponse pagination = 1; */
        if (message.pagination)
            PaginationResponse.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.accounts.Account accounts = 2; */
        for (let i = 0; i < message.accounts.length; i++)
            Account.internalBinaryWrite(message.accounts[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.settings.ListAccountsResponse
 */
export const ListAccountsResponse = new ListAccountsResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UpdateAccountRequest$Type extends MessageType<UpdateAccountRequest> {
    constructor() {
        super("services.settings.UpdateAccountRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/, options: { "buf.validate.field": { uint64: { gt: "0" } } } },
            { no: 2, name: "enabled", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ },
            { no: 3, name: "last_char", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/ }
        ]);
    }
    create(value?: PartialMessage<UpdateAccountRequest>): UpdateAccountRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        if (value !== undefined)
            reflectionMergePartial<UpdateAccountRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UpdateAccountRequest): UpdateAccountRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* optional bool enabled */ 2:
                    message.enabled = reader.bool();
                    break;
                case /* optional int32 last_char */ 3:
                    message.lastChar = reader.int32();
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
    internalBinaryWrite(message: UpdateAccountRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* optional bool enabled = 2; */
        if (message.enabled !== undefined)
            writer.tag(2, WireType.Varint).bool(message.enabled);
        /* optional int32 last_char = 3; */
        if (message.lastChar !== undefined)
            writer.tag(3, WireType.Varint).int32(message.lastChar);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.settings.UpdateAccountRequest
 */
export const UpdateAccountRequest = new UpdateAccountRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UpdateAccountResponse$Type extends MessageType<UpdateAccountResponse> {
    constructor() {
        super("services.settings.UpdateAccountResponse", [
            { no: 1, name: "account", kind: "message", T: () => Account }
        ]);
    }
    create(value?: PartialMessage<UpdateAccountResponse>): UpdateAccountResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<UpdateAccountResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UpdateAccountResponse): UpdateAccountResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.accounts.Account account */ 1:
                    message.account = Account.internalBinaryRead(reader, reader.uint32(), options, message.account);
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
    internalBinaryWrite(message: UpdateAccountResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.accounts.Account account = 1; */
        if (message.account)
            Account.internalBinaryWrite(message.account, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.settings.UpdateAccountResponse
 */
export const UpdateAccountResponse = new UpdateAccountResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DisconnectOAuth2ConnectionRequest$Type extends MessageType<DisconnectOAuth2ConnectionRequest> {
    constructor() {
        super("services.settings.DisconnectOAuth2ConnectionRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/, options: { "buf.validate.field": { uint64: { gt: "0" } } } },
            { no: 2, name: "provider_name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "255" } } } }
        ]);
    }
    create(value?: PartialMessage<DisconnectOAuth2ConnectionRequest>): DisconnectOAuth2ConnectionRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.providerName = "";
        if (value !== undefined)
            reflectionMergePartial<DisconnectOAuth2ConnectionRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DisconnectOAuth2ConnectionRequest): DisconnectOAuth2ConnectionRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* string provider_name */ 2:
                    message.providerName = reader.string();
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
    internalBinaryWrite(message: DisconnectOAuth2ConnectionRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* string provider_name = 2; */
        if (message.providerName !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.providerName);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.settings.DisconnectOAuth2ConnectionRequest
 */
export const DisconnectOAuth2ConnectionRequest = new DisconnectOAuth2ConnectionRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DisconnectOAuth2ConnectionResponse$Type extends MessageType<DisconnectOAuth2ConnectionResponse> {
    constructor() {
        super("services.settings.DisconnectOAuth2ConnectionResponse", []);
    }
    create(value?: PartialMessage<DisconnectOAuth2ConnectionResponse>): DisconnectOAuth2ConnectionResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<DisconnectOAuth2ConnectionResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DisconnectOAuth2ConnectionResponse): DisconnectOAuth2ConnectionResponse {
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
    internalBinaryWrite(message: DisconnectOAuth2ConnectionResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.settings.DisconnectOAuth2ConnectionResponse
 */
export const DisconnectOAuth2ConnectionResponse = new DisconnectOAuth2ConnectionResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteAccountRequest$Type extends MessageType<DeleteAccountRequest> {
    constructor() {
        super("services.settings.DeleteAccountRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/, options: { "buf.validate.field": { uint64: { gt: "0" } } } }
        ]);
    }
    create(value?: PartialMessage<DeleteAccountRequest>): DeleteAccountRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        if (value !== undefined)
            reflectionMergePartial<DeleteAccountRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteAccountRequest): DeleteAccountRequest {
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
    internalBinaryWrite(message: DeleteAccountRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
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
 * @generated MessageType for protobuf message services.settings.DeleteAccountRequest
 */
export const DeleteAccountRequest = new DeleteAccountRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteAccountResponse$Type extends MessageType<DeleteAccountResponse> {
    constructor() {
        super("services.settings.DeleteAccountResponse", []);
    }
    create(value?: PartialMessage<DeleteAccountResponse>): DeleteAccountResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<DeleteAccountResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteAccountResponse): DeleteAccountResponse {
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
    internalBinaryWrite(message: DeleteAccountResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.settings.DeleteAccountResponse
 */
export const DeleteAccountResponse = new DeleteAccountResponse$Type();
/**
 * @generated ServiceType for protobuf service services.settings.AccountsService
 */
export const AccountsService = new ServiceType("services.settings.AccountsService", [
    { name: "ListAccounts", options: {}, I: ListAccountsRequest, O: ListAccountsResponse },
    { name: "UpdateAccount", options: {}, I: UpdateAccountRequest, O: UpdateAccountResponse },
    { name: "DisconnectOAuth2Connection", options: {}, I: DisconnectOAuth2ConnectionRequest, O: DisconnectOAuth2ConnectionResponse },
    { name: "DeleteAccount", options: {}, I: DeleteAccountRequest, O: DeleteAccountResponse }
]);
