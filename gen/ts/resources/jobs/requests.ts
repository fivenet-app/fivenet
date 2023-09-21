// @generated by protobuf-ts 2.9.1 with parameter optimize_code_size,long_type_bigint
// @generated from protobuf file "resources/jobs/requests.proto" (package "resources.jobs", syntax proto3)
// tslint:disable
import { MessageType } from "@protobuf-ts/runtime";
import { UserShort } from "../users/users.js";
import { Timestamp } from "../timestamp/timestamp.js";
/**
 * @generated from protobuf message resources.jobs.RequestEntry
 */
export interface RequestEntry {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: bigint; // @gotags: sql:"primary_key" alias:"id"
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 3;
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: string job = 4;
     */
    job: string;
    /**
     * @generated from protobuf field: resources.jobs.RequestType type = 5;
     */
    type: RequestType;
    /**
     * @sanitize
     *
     * @generated from protobuf field: string message = 6;
     */
    message: string;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp begins_at = 7;
     */
    beginsAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp ends_at = 8;
     */
    endsAt?: Timestamp;
    /**
     * @sanitize
     *
     * @generated from protobuf field: optional string status = 9;
     */
    status?: string;
    /**
     * @generated from protobuf field: int32 creator_id = 10;
     */
    creatorId: number;
    /**
     * @generated from protobuf field: optional resources.users.UserShort creator = 11;
     */
    creator?: UserShort; // @gotags: alias:"creator"
    /**
     * @generated from protobuf field: optional bool approved = 12;
     */
    approved?: boolean;
    /**
     * @generated from protobuf field: optional int32 approver_user_id = 13;
     */
    approverUserId?: number;
    /**
     * @generated from protobuf field: optional resources.users.UserShort approver_user = 14;
     */
    approverUser?: UserShort; // @gotags: alias:"approver"
}
/**
 * @generated from protobuf enum resources.jobs.RequestType
 */
export enum RequestType {
    /**
     * @generated from protobuf enum value: REQUEST_TYPE_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from protobuf enum value: REQUEST_TYPE_ABSENCE = 1;
     */
    ABSENCE = 1
}
// @generated message type with reflection information, may provide speed optimized methods
class RequestEntry$Type extends MessageType<RequestEntry> {
    constructor() {
        super("resources.jobs.RequestEntry", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 0 /*LongType.BIGINT*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 5, name: "type", kind: "enum", T: () => ["resources.jobs.RequestType", RequestType, "REQUEST_TYPE_"], options: { "validate.rules": { enum: { definedOnly: true } } } },
            { no: 6, name: "message", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "3", maxLen: "2048" } } } },
            { no: 7, name: "begins_at", kind: "message", T: () => Timestamp },
            { no: 8, name: "ends_at", kind: "message", T: () => Timestamp },
            { no: 9, name: "status", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 10, name: "creator_id", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gt: 0 } } } },
            { no: 11, name: "creator", kind: "message", T: () => UserShort },
            { no: 12, name: "approved", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ },
            { no: 13, name: "approver_user_id", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/ },
            { no: 14, name: "approver_user", kind: "message", T: () => UserShort }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.jobs.RequestEntry
 */
export const RequestEntry = new RequestEntry$Type();
