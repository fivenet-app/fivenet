// @generated by protobuf-ts 2.9.1 with parameter optimize_code_size,long_type_bigint
// @generated from protobuf file "resources/jobs/requests.proto" (package "resources.jobs", syntax proto3)
// tslint:disable
import { MessageType } from "@protobuf-ts/runtime";
import { UserShort } from "../users/users";
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.jobs.RequestType
 */
export interface RequestType {
    /**
     * @generated from protobuf field: uint64 id = 1 [jstype = JS_STRING];
     */
    id: string;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 3;
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp deleted_at = 4;
     */
    deletedAt?: Timestamp;
    /**
     * @generated from protobuf field: string job = 5;
     */
    job: string;
    /**
     * @generated from protobuf field: string name = 6;
     */
    name: string;
    /**
     * @generated from protobuf field: optional string description = 7;
     */
    description?: string;
    /**
     * @generated from protobuf field: uint32 weight = 8;
     */
    weight: number;
}
/**
 * @generated from protobuf message resources.jobs.Request
 */
export interface Request {
    /**
     * @generated from protobuf field: uint64 id = 1 [jstype = JS_STRING];
     */
    id: string; // @gotags: sql:"primary_key"
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 3;
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp deleted_at = 4;
     */
    deletedAt?: Timestamp;
    /**
     * @generated from protobuf field: string job = 5;
     */
    job: string;
    /**
     * @generated from protobuf field: optional uint64 type_id = 6 [jstype = JS_STRING];
     */
    typeId?: string;
    /**
     * @generated from protobuf field: optional resources.jobs.RequestType type = 7;
     */
    type?: RequestType;
    /**
     * @sanitize
     *
     * @generated from protobuf field: string title = 8;
     */
    title: string;
    /**
     * @sanitize
     *
     * @generated from protobuf field: string message = 9;
     */
    message: string;
    /**
     * @sanitize
     *
     * @generated from protobuf field: optional string status = 10;
     */
    status?: string;
    /**
     * @generated from protobuf field: int32 creator_id = 11;
     */
    creatorId: number;
    /**
     * @generated from protobuf field: optional resources.users.UserShort creator = 12;
     */
    creator?: UserShort; // @gotags: alias:"creator"
    /**
     * @generated from protobuf field: optional bool approved = 13;
     */
    approved?: boolean;
    /**
     * @generated from protobuf field: optional int32 approver_id = 14;
     */
    approverId?: number;
    /**
     * @generated from protobuf field: optional resources.users.UserShort approver_user = 15;
     */
    approverUser?: UserShort; // @gotags: alias:"approver"
    /**
     * @generated from protobuf field: bool closed = 16;
     */
    closed: boolean;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp begins_at = 17;
     */
    beginsAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp ends_at = 18;
     */
    endsAt?: Timestamp;
}
/**
 * @generated from protobuf message resources.jobs.RequestComment
 */
export interface RequestComment {
    /**
     * @generated from protobuf field: uint64 id = 1 [jstype = JS_STRING];
     */
    id: string;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 3;
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp deleted_at = 4;
     */
    deletedAt?: Timestamp;
    /**
     * @generated from protobuf field: uint64 request_id = 5 [jstype = JS_STRING];
     */
    requestId: string;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: string comment = 6;
     */
    comment: string;
    /**
     * @generated from protobuf field: optional int32 creator_id = 7;
     */
    creatorId?: number;
    /**
     * @generated from protobuf field: optional resources.users.UserShort creator = 8;
     */
    creator?: UserShort; // @gotags: alias:"creator"
}
// @generated message type with reflection information, may provide speed optimized methods
class RequestType$Type extends MessageType<RequestType> {
    constructor() {
        super("resources.jobs.RequestType", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "deleted_at", kind: "message", T: () => Timestamp },
            { no: 5, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 6, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "3", maxLen: "32" } } } },
            { no: 7, name: "description", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "3", maxLen: "255" } } } },
            { no: 8, name: "weight", kind: "scalar", T: 13 /*ScalarType.UINT32*/, options: { "validate.rules": { uint32: { lt: 4294967295 } } } }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.jobs.RequestType
 */
export const RequestType = new RequestType$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Request$Type extends MessageType<Request> {
    constructor() {
        super("resources.jobs.Request", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "deleted_at", kind: "message", T: () => Timestamp },
            { no: 5, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 6, name: "type_id", kind: "scalar", opt: true, T: 4 /*ScalarType.UINT64*/ },
            { no: 7, name: "type", kind: "message", T: () => RequestType },
            { no: 8, name: "title", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "3", maxLen: "255" } } } },
            { no: 9, name: "message", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "3", maxLen: "4096" } } } },
            { no: 10, name: "status", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "24" } } } },
            { no: 11, name: "creator_id", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gt: 0 } } } },
            { no: 12, name: "creator", kind: "message", T: () => UserShort },
            { no: 13, name: "approved", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ },
            { no: 14, name: "approver_id", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/ },
            { no: 15, name: "approver_user", kind: "message", T: () => UserShort },
            { no: 16, name: "closed", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 17, name: "begins_at", kind: "message", T: () => Timestamp },
            { no: 18, name: "ends_at", kind: "message", T: () => Timestamp }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.jobs.Request
 */
export const Request = new Request$Type();
// @generated message type with reflection information, may provide speed optimized methods
class RequestComment$Type extends MessageType<RequestComment> {
    constructor() {
        super("resources.jobs.RequestComment", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "deleted_at", kind: "message", T: () => Timestamp },
            { no: 5, name: "request_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ },
            { no: 6, name: "comment", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "3", maxBytes: "2048" } } } },
            { no: 7, name: "creator_id", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/ },
            { no: 8, name: "creator", kind: "message", T: () => UserShort }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.jobs.RequestComment
 */
export const RequestComment = new RequestComment$Type();
