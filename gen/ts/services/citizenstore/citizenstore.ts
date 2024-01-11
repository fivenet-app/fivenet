// @generated by protobuf-ts 2.9.3 with parameter optimize_code_size,long_type_bigint
// @generated from protobuf file "services/citizenstore/citizenstore.proto" (package "services.citizenstore", syntax proto3)
// tslint:disable
import { ServiceType } from "@protobuf-ts/runtime-rpc";
import { MessageType } from "@protobuf-ts/runtime";
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
     * @generated from protobuf field: string search_name = 2;
     */
    searchName: string;
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
    openFines?: bigint;
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
// @generated message type with reflection information, may provide speed optimized methods
class ListCitizensRequest$Type extends MessageType<ListCitizensRequest> {
    constructor() {
        super("services.citizenstore.ListCitizensRequest", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationRequest, options: { "validate.rules": { message: { required: true } } } },
            { no: 2, name: "search_name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 3, name: "wanted", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ },
            { no: 4, name: "phone_number", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 5, name: "traffic_infraction_points", kind: "scalar", opt: true, T: 13 /*ScalarType.UINT32*/ },
            { no: 6, name: "dateofbirth", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "10" } } } },
            { no: 7, name: "open_fines", kind: "scalar", opt: true, T: 4 /*ScalarType.UINT64*/, L: 0 /*LongType.BIGINT*/ }
        ]);
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
}
/**
 * @generated MessageType for protobuf message services.citizenstore.SetUserPropsResponse
 */
export const SetUserPropsResponse = new SetUserPropsResponse$Type();
/**
 * @generated ServiceType for protobuf service services.citizenstore.CitizenStoreService
 */
export const CitizenStoreService = new ServiceType("services.citizenstore.CitizenStoreService", [
    { name: "ListCitizens", options: {}, I: ListCitizensRequest, O: ListCitizensResponse },
    { name: "GetUser", options: {}, I: GetUserRequest, O: GetUserResponse },
    { name: "ListUserActivity", options: {}, I: ListUserActivityRequest, O: ListUserActivityResponse },
    { name: "SetUserProps", options: {}, I: SetUserPropsRequest, O: SetUserPropsResponse }
]);