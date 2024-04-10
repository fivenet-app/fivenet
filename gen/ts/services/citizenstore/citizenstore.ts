// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/citizenstore/citizenstore.proto" (package "services.citizenstore", syntax proto3)
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
import { File } from "../../resources/filestore/file";
import { UserProps } from "../../resources/users/users";
import { UserActivity } from "../../resources/users/users";
import { User } from "../../resources/users/users";
import { PaginationResponse } from "../../resources/common/database/database";
import { PaginationRequest } from "../../resources/common/database/database";
/**
 * @generated from protobuf message services.citizenstore.ListCitizensRequest
 */
export interface ListCitizensRequest {
    /**
     * @generated from protobuf field: resources.common.database.PaginationRequest pagination = 1;
     */
    pagination?: PaginationRequest;
    /**
     * Search params
     *
     * @generated from protobuf field: string search = 2;
     */
    search: string;
    /**
     * @generated from protobuf field: optional bool wanted = 3;
     */
    wanted?: boolean;
    /**
     * @generated from protobuf field: optional string phone_number = 4;
     */
    phoneNumber?: string;
    /**
     * @generated from protobuf field: optional uint32 traffic_infraction_points = 5;
     */
    trafficInfractionPoints?: number;
    /**
     * @generated from protobuf field: optional string dateofbirth = 6;
     */
    dateofbirth?: string;
    /**
     * @generated from protobuf field: optional uint64 open_fines = 7;
     */
    openFines?: number;
}
/**
 * @generated from protobuf message services.citizenstore.ListCitizensResponse
 */
export interface ListCitizensResponse {
    /**
     * @generated from protobuf field: resources.common.database.PaginationResponse pagination = 1;
     */
    pagination?: PaginationResponse;
    /**
     * @generated from protobuf field: repeated resources.users.User users = 2;
     */
    users: User[];
}
/**
 * @generated from protobuf message services.citizenstore.GetUserRequest
 */
export interface GetUserRequest {
    /**
     * @generated from protobuf field: int32 user_id = 1;
     */
    userId: number;
}
/**
 * @generated from protobuf message services.citizenstore.GetUserResponse
 */
export interface GetUserResponse {
    /**
     * @generated from protobuf field: resources.users.User user = 1;
     */
    user?: User;
}
/**
 * @generated from protobuf message services.citizenstore.ListUserActivityRequest
 */
export interface ListUserActivityRequest {
    /**
     * @generated from protobuf field: resources.common.database.PaginationRequest pagination = 1;
     */
    pagination?: PaginationRequest;
    /**
     * @generated from protobuf field: int32 user_id = 2;
     */
    userId: number;
}
/**
 * @generated from protobuf message services.citizenstore.ListUserActivityResponse
 */
export interface ListUserActivityResponse {
    /**
     * @generated from protobuf field: resources.common.database.PaginationResponse pagination = 1;
     */
    pagination?: PaginationResponse;
    /**
     * @generated from protobuf field: repeated resources.users.UserActivity activity = 2;
     */
    activity: UserActivity[];
}
/**
 * @generated from protobuf message services.citizenstore.SetUserPropsRequest
 */
export interface SetUserPropsRequest {
    /**
     * @generated from protobuf field: resources.users.UserProps props = 1;
     */
    props?: UserProps;
    /**
     * @sanitize
     *
     * @generated from protobuf field: string reason = 2;
     */
    reason: string;
}
/**
 * @generated from protobuf message services.citizenstore.SetUserPropsResponse
 */
export interface SetUserPropsResponse {
    /**
     * @generated from protobuf field: resources.users.UserProps props = 1;
     */
    props?: UserProps;
}
/**
 * @generated from protobuf message services.citizenstore.SetProfilePictureRequest
 */
export interface SetProfilePictureRequest {
    /**
     * @generated from protobuf field: resources.filestore.File avatar = 1;
     */
    avatar?: File;
}
/**
 * @generated from protobuf message services.citizenstore.SetProfilePictureResponse
 */
export interface SetProfilePictureResponse {
    /**
     * @generated from protobuf field: resources.filestore.File avatar = 1;
     */
    avatar?: File;
}
// @generated message type with reflection information, may provide speed optimized methods
class ListCitizensRequest$Type extends MessageType<ListCitizensRequest> {
    constructor() {
        super("services.citizenstore.ListCitizensRequest", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationRequest, options: { "validate.rules": { message: { required: true } } } },
            { no: 2, name: "search", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 3, name: "wanted", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ },
            { no: 4, name: "phone_number", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 5, name: "traffic_infraction_points", kind: "scalar", opt: true, T: 13 /*ScalarType.UINT32*/ },
            { no: 6, name: "dateofbirth", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "10" } } } },
            { no: 7, name: "open_fines", kind: "scalar", opt: true, T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ }
        ]);
    }
    create(value?: PartialMessage<ListCitizensRequest>): ListCitizensRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.search = "";
        if (value !== undefined)
            reflectionMergePartial<ListCitizensRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListCitizensRequest): ListCitizensRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationRequest pagination */ 1:
                    message.pagination = PaginationRequest.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* string search */ 2:
                    message.search = reader.string();
                    break;
                case /* optional bool wanted */ 3:
                    message.wanted = reader.bool();
                    break;
                case /* optional string phone_number */ 4:
                    message.phoneNumber = reader.string();
                    break;
                case /* optional uint32 traffic_infraction_points */ 5:
                    message.trafficInfractionPoints = reader.uint32();
                    break;
                case /* optional string dateofbirth */ 6:
                    message.dateofbirth = reader.string();
                    break;
                case /* optional uint64 open_fines */ 7:
                    message.openFines = reader.uint64().toNumber();
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
    internalBinaryWrite(message: ListCitizensRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationRequest pagination = 1; */
        if (message.pagination)
            PaginationRequest.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* string search = 2; */
        if (message.search !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.search);
        /* optional bool wanted = 3; */
        if (message.wanted !== undefined)
            writer.tag(3, WireType.Varint).bool(message.wanted);
        /* optional string phone_number = 4; */
        if (message.phoneNumber !== undefined)
            writer.tag(4, WireType.LengthDelimited).string(message.phoneNumber);
        /* optional uint32 traffic_infraction_points = 5; */
        if (message.trafficInfractionPoints !== undefined)
            writer.tag(5, WireType.Varint).uint32(message.trafficInfractionPoints);
        /* optional string dateofbirth = 6; */
        if (message.dateofbirth !== undefined)
            writer.tag(6, WireType.LengthDelimited).string(message.dateofbirth);
        /* optional uint64 open_fines = 7; */
        if (message.openFines !== undefined)
            writer.tag(7, WireType.Varint).uint64(message.openFines);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.citizenstore.ListCitizensRequest
 */
export const ListCitizensRequest = new ListCitizensRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ListCitizensResponse$Type extends MessageType<ListCitizensResponse> {
    constructor() {
        super("services.citizenstore.ListCitizensResponse", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationResponse },
            { no: 2, name: "users", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => User }
        ]);
    }
    create(value?: PartialMessage<ListCitizensResponse>): ListCitizensResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.users = [];
        if (value !== undefined)
            reflectionMergePartial<ListCitizensResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListCitizensResponse): ListCitizensResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationResponse pagination */ 1:
                    message.pagination = PaginationResponse.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* repeated resources.users.User users */ 2:
                    message.users.push(User.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: ListCitizensResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationResponse pagination = 1; */
        if (message.pagination)
            PaginationResponse.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.users.User users = 2; */
        for (let i = 0; i < message.users.length; i++)
            User.internalBinaryWrite(message.users[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.citizenstore.ListCitizensResponse
 */
export const ListCitizensResponse = new ListCitizensResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetUserRequest$Type extends MessageType<GetUserRequest> {
    constructor() {
        super("services.citizenstore.GetUserRequest", [
            { no: 1, name: "user_id", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gt: 0 } } } }
        ]);
    }
    create(value?: PartialMessage<GetUserRequest>): GetUserRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.userId = 0;
        if (value !== undefined)
            reflectionMergePartial<GetUserRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: GetUserRequest): GetUserRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 user_id */ 1:
                    message.userId = reader.int32();
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
    internalBinaryWrite(message: GetUserRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* int32 user_id = 1; */
        if (message.userId !== 0)
            writer.tag(1, WireType.Varint).int32(message.userId);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.citizenstore.GetUserRequest
 */
export const GetUserRequest = new GetUserRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetUserResponse$Type extends MessageType<GetUserResponse> {
    constructor() {
        super("services.citizenstore.GetUserResponse", [
            { no: 1, name: "user", kind: "message", T: () => User }
        ]);
    }
    create(value?: PartialMessage<GetUserResponse>): GetUserResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<GetUserResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: GetUserResponse): GetUserResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.users.User user */ 1:
                    message.user = User.internalBinaryRead(reader, reader.uint32(), options, message.user);
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
    internalBinaryWrite(message: GetUserResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.users.User user = 1; */
        if (message.user)
            User.internalBinaryWrite(message.user, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.citizenstore.GetUserResponse
 */
export const GetUserResponse = new GetUserResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ListUserActivityRequest$Type extends MessageType<ListUserActivityRequest> {
    constructor() {
        super("services.citizenstore.ListUserActivityRequest", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationRequest, options: { "validate.rules": { message: { required: true } } } },
            { no: 2, name: "user_id", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gt: 0 } } } }
        ]);
    }
    create(value?: PartialMessage<ListUserActivityRequest>): ListUserActivityRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.userId = 0;
        if (value !== undefined)
            reflectionMergePartial<ListUserActivityRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListUserActivityRequest): ListUserActivityRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationRequest pagination */ 1:
                    message.pagination = PaginationRequest.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* int32 user_id */ 2:
                    message.userId = reader.int32();
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
    internalBinaryWrite(message: ListUserActivityRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationRequest pagination = 1; */
        if (message.pagination)
            PaginationRequest.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* int32 user_id = 2; */
        if (message.userId !== 0)
            writer.tag(2, WireType.Varint).int32(message.userId);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.citizenstore.ListUserActivityRequest
 */
export const ListUserActivityRequest = new ListUserActivityRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ListUserActivityResponse$Type extends MessageType<ListUserActivityResponse> {
    constructor() {
        super("services.citizenstore.ListUserActivityResponse", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationResponse },
            { no: 2, name: "activity", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => UserActivity }
        ]);
    }
    create(value?: PartialMessage<ListUserActivityResponse>): ListUserActivityResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.activity = [];
        if (value !== undefined)
            reflectionMergePartial<ListUserActivityResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListUserActivityResponse): ListUserActivityResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationResponse pagination */ 1:
                    message.pagination = PaginationResponse.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* repeated resources.users.UserActivity activity */ 2:
                    message.activity.push(UserActivity.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: ListUserActivityResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationResponse pagination = 1; */
        if (message.pagination)
            PaginationResponse.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.users.UserActivity activity = 2; */
        for (let i = 0; i < message.activity.length; i++)
            UserActivity.internalBinaryWrite(message.activity[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.citizenstore.ListUserActivityResponse
 */
export const ListUserActivityResponse = new ListUserActivityResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class SetUserPropsRequest$Type extends MessageType<SetUserPropsRequest> {
    constructor() {
        super("services.citizenstore.SetUserPropsRequest", [
            { no: 1, name: "props", kind: "message", T: () => UserProps, options: { "validate.rules": { message: { required: true } } } },
            { no: 2, name: "reason", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "3", maxLen: "255", ignoreEmpty: true } } } }
        ]);
    }
    create(value?: PartialMessage<SetUserPropsRequest>): SetUserPropsRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.reason = "";
        if (value !== undefined)
            reflectionMergePartial<SetUserPropsRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SetUserPropsRequest): SetUserPropsRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.users.UserProps props */ 1:
                    message.props = UserProps.internalBinaryRead(reader, reader.uint32(), options, message.props);
                    break;
                case /* string reason */ 2:
                    message.reason = reader.string();
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
    internalBinaryWrite(message: SetUserPropsRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.users.UserProps props = 1; */
        if (message.props)
            UserProps.internalBinaryWrite(message.props, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* string reason = 2; */
        if (message.reason !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.reason);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.citizenstore.SetUserPropsRequest
 */
export const SetUserPropsRequest = new SetUserPropsRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class SetUserPropsResponse$Type extends MessageType<SetUserPropsResponse> {
    constructor() {
        super("services.citizenstore.SetUserPropsResponse", [
            { no: 1, name: "props", kind: "message", T: () => UserProps }
        ]);
    }
    create(value?: PartialMessage<SetUserPropsResponse>): SetUserPropsResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<SetUserPropsResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SetUserPropsResponse): SetUserPropsResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.users.UserProps props */ 1:
                    message.props = UserProps.internalBinaryRead(reader, reader.uint32(), options, message.props);
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
    internalBinaryWrite(message: SetUserPropsResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.users.UserProps props = 1; */
        if (message.props)
            UserProps.internalBinaryWrite(message.props, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.citizenstore.SetUserPropsResponse
 */
export const SetUserPropsResponse = new SetUserPropsResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class SetProfilePictureRequest$Type extends MessageType<SetProfilePictureRequest> {
    constructor() {
        super("services.citizenstore.SetProfilePictureRequest", [
            { no: 1, name: "avatar", kind: "message", T: () => File }
        ]);
    }
    create(value?: PartialMessage<SetProfilePictureRequest>): SetProfilePictureRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<SetProfilePictureRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SetProfilePictureRequest): SetProfilePictureRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.filestore.File avatar */ 1:
                    message.avatar = File.internalBinaryRead(reader, reader.uint32(), options, message.avatar);
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
    internalBinaryWrite(message: SetProfilePictureRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.filestore.File avatar = 1; */
        if (message.avatar)
            File.internalBinaryWrite(message.avatar, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.citizenstore.SetProfilePictureRequest
 */
export const SetProfilePictureRequest = new SetProfilePictureRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class SetProfilePictureResponse$Type extends MessageType<SetProfilePictureResponse> {
    constructor() {
        super("services.citizenstore.SetProfilePictureResponse", [
            { no: 1, name: "avatar", kind: "message", T: () => File }
        ]);
    }
    create(value?: PartialMessage<SetProfilePictureResponse>): SetProfilePictureResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<SetProfilePictureResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SetProfilePictureResponse): SetProfilePictureResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.filestore.File avatar */ 1:
                    message.avatar = File.internalBinaryRead(reader, reader.uint32(), options, message.avatar);
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
    internalBinaryWrite(message: SetProfilePictureResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.filestore.File avatar = 1; */
        if (message.avatar)
            File.internalBinaryWrite(message.avatar, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.citizenstore.SetProfilePictureResponse
 */
export const SetProfilePictureResponse = new SetProfilePictureResponse$Type();
/**
 * @generated ServiceType for protobuf service services.citizenstore.CitizenStoreService
 */
export const CitizenStoreService = new ServiceType("services.citizenstore.CitizenStoreService", [
    { name: "ListCitizens", options: {}, I: ListCitizensRequest, O: ListCitizensResponse },
    { name: "GetUser", options: {}, I: GetUserRequest, O: GetUserResponse },
    { name: "ListUserActivity", options: {}, I: ListUserActivityRequest, O: ListUserActivityResponse },
    { name: "SetUserProps", options: {}, I: SetUserPropsRequest, O: SetUserPropsResponse },
    { name: "SetProfilePicture", options: {}, I: SetProfilePictureRequest, O: SetProfilePictureResponse }
]);
