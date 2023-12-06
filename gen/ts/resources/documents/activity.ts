// @generated by protobuf-ts 2.9.2 with parameter optimize_code_size,long_type_bigint
// @generated from protobuf file "resources/documents/activity.proto" (package "resources.documents", syntax proto3)
// tslint:disable
import { MessageType } from "@protobuf-ts/runtime";
import { UserShort } from "../users/users";
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.documents.DocActivity
 */
export interface DocActivity {
    /**
     * @generated from protobuf field: uint64 id = 1 [jstype = JS_STRING];
     */
    id: string;
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: uint64 document_id = 3 [jstype = JS_STRING];
     */
    documentId: string;
    /**
     * @generated from protobuf field: resources.documents.DocActivityType activity_type = 4;
     */
    activityType: DocActivityType;
    /**
     * @generated from protobuf field: optional int32 creator_id = 5;
     */
    creatorId?: number;
    /**
     * @generated from protobuf field: optional resources.users.UserShort creator = 6;
     */
    creator?: UserShort; // @gotags: alias:"creator"
    /**
     * @generated from protobuf field: string creator_job = 7;
     */
    creatorJob: string;
    /**
     * @generated from protobuf field: optional string creator_job_label = 8;
     */
    creatorJobLabel?: string;
    /**
     * @generated from protobuf field: optional string reason = 9;
     */
    reason?: string;
    /**
     * @generated from protobuf field: resources.documents.DocActivityData data = 10;
     */
    data?: DocActivityData;
}
/**
 * @generated from protobuf message resources.documents.DocActivityData
 */
export interface DocActivityData {
    /**
     * @generated from protobuf oneof: data
     */
    data: {
        oneofKind: "updated";
        /**
         * @generated from protobuf field: resources.documents.DocUpdated updated = 1;
         */
        updated: DocUpdated;
    } | {
        oneofKind: "ownerChanged";
        /**
         * @generated from protobuf field: resources.documents.DocOwnerChanged owner_changed = 2;
         */
        ownerChanged: DocOwnerChanged;
    } | {
        oneofKind: "requestInfo";
        /**
         * @generated from protobuf field: resources.documents.DocRequestActivity request_info = 3;
         */
        requestInfo: DocRequestActivity;
    } | {
        oneofKind: undefined;
    };
}
/**
 * @generated from protobuf message resources.documents.DocUpdated
 */
export interface DocUpdated {
    /**
     * @generated from protobuf field: optional string title_diff = 1;
     */
    titleDiff?: string;
    /**
     * @generated from protobuf field: optional string content_diff = 2;
     */
    contentDiff?: string;
    /**
     * @generated from protobuf field: optional string state_diff = 3;
     */
    stateDiff?: string;
}
/**
 * @generated from protobuf message resources.documents.DocOwnerChanged
 */
export interface DocOwnerChanged {
    /**
     * @generated from protobuf field: int32 new_owner_id = 1;
     */
    newOwnerId: number;
    /**
     * @generated from protobuf field: resources.users.UserShort new_owner = 2;
     */
    newOwner?: UserShort;
}
/**
 * @generated from protobuf message resources.documents.DocRequestActivity
 */
export interface DocRequestActivity {
    /**
     * @generated from protobuf field: bool accepted = 1;
     */
    accepted: boolean;
}
/**
 * @generated from protobuf enum resources.documents.DocActivityType
 */
export enum DocActivityType {
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * Base
     *
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_CREATED = 1;
     */
    CREATED = 1,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_STATUS_OPEN = 2;
     */
    STATUS_OPEN = 2,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_STATUS_CLOSED = 3;
     */
    STATUS_CLOSED = 3,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_UPDATED = 4;
     */
    UPDATED = 4,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_RELATIONS_UPDATED = 5;
     */
    RELATIONS_UPDATED = 5,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_REFERENCES_UPDATED = 6;
     */
    REFERENCES_UPDATED = 6,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_ACCESS_UPDATED = 7;
     */
    ACCESS_UPDATED = 7,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_OWNER_CHANGED = 8;
     */
    OWNER_CHANGED = 8,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_DELETED = 9;
     */
    DELETED = 9,
    /**
     * Comments
     *
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_COMMENT_ADDED = 10;
     */
    COMMENT_ADDED = 10,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_COMMENT_UPDATED = 11;
     */
    COMMENT_UPDATED = 11,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_COMMENT_DELETED = 12;
     */
    COMMENT_DELETED = 12,
    /**
     * Requests
     *
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_REQUESTED_ACCESS = 13;
     */
    REQUESTED_ACCESS = 13,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_REQUESTED_CLOSURE = 14;
     */
    REQUESTED_CLOSURE = 14,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_REQUESTED_OPENING = 15;
     */
    REQUESTED_OPENING = 15,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_REQUESTED_UPDATE = 16;
     */
    REQUESTED_UPDATE = 16,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_REQUESTED_OWNER_CHANGE = 17;
     */
    REQUESTED_OWNER_CHANGE = 17,
    /**
     * @generated from protobuf enum value: DOC_ACTIVITY_TYPE_REQUESTED_DELETION = 18;
     */
    REQUESTED_DELETION = 18
}
// @generated message type with reflection information, may provide speed optimized methods
class DocActivity$Type extends MessageType<DocActivity> {
    constructor() {
        super("resources.documents.DocActivity", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "document_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ },
            { no: 4, name: "activity_type", kind: "enum", T: () => ["resources.documents.DocActivityType", DocActivityType, "DOC_ACTIVITY_TYPE_"] },
            { no: 5, name: "creator_id", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/ },
            { no: 6, name: "creator", kind: "message", T: () => UserShort },
            { no: 7, name: "creator_job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 8, name: "creator_job_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 9, name: "reason", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 10, name: "data", kind: "message", T: () => DocActivityData }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.documents.DocActivity
 */
export const DocActivity = new DocActivity$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DocActivityData$Type extends MessageType<DocActivityData> {
    constructor() {
        super("resources.documents.DocActivityData", [
            { no: 1, name: "updated", kind: "message", oneof: "data", T: () => DocUpdated },
            { no: 2, name: "owner_changed", kind: "message", oneof: "data", T: () => DocOwnerChanged },
            { no: 3, name: "request_info", kind: "message", oneof: "data", T: () => DocRequestActivity }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.documents.DocActivityData
 */
export const DocActivityData = new DocActivityData$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DocUpdated$Type extends MessageType<DocUpdated> {
    constructor() {
        super("resources.documents.DocUpdated", [
            { no: 1, name: "title_diff", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "content_diff", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "state_diff", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.documents.DocUpdated
 */
export const DocUpdated = new DocUpdated$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DocOwnerChanged$Type extends MessageType<DocOwnerChanged> {
    constructor() {
        super("resources.documents.DocOwnerChanged", [
            { no: 1, name: "new_owner_id", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "new_owner", kind: "message", T: () => UserShort }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.documents.DocOwnerChanged
 */
export const DocOwnerChanged = new DocOwnerChanged$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DocRequestActivity$Type extends MessageType<DocRequestActivity> {
    constructor() {
        super("resources.documents.DocRequestActivity", [
            { no: 1, name: "accepted", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.documents.DocRequestActivity
 */
export const DocRequestActivity = new DocRequestActivity$Type();
