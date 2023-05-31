// @generated by protobuf-ts 2.9.0 with parameter optimize_code_size,long_type_bigint
// @generated from protobuf file "resources/notifications/notifications.proto" (package "resources.notifications", syntax proto3)
// tslint:disable
import { MessageType } from "@protobuf-ts/runtime";
import { Timestamp } from "../timestamp/timestamp.js";
/**
 * @generated from protobuf message resources.notifications.Notification
 */
export interface Notification {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: bigint; // @gotags: alias:"id"
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp; // @gotags: alias:"created_at"
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp read_at = 3;
     */
    readAt?: Timestamp; // @gotags: alias:"read_at"
    /**
     * @generated from protobuf field: int32 user_id = 4;
     */
    userId: number; // @gotags: alias:"user_id"
    /**
     * @generated from protobuf field: string title = 5;
     */
    title: string; // @gotags: alias:"title"
    /**
     * @generated from protobuf field: optional string type = 6;
     */
    type?: string; // @gotags: alias:"type"
    /**
     * @generated from protobuf field: string content = 7;
     */
    content: string; // @gotags: alias:"content"
    /**
     * @generated from protobuf field: optional string data = 8;
     */
    data?: string; // @gotags: alias:"data"
}
// @generated message type with reflection information, may provide speed optimized methods
class Notification$Type extends MessageType<Notification> {
    constructor() {
        super("resources.notifications.Notification", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 0 /*LongType.BIGINT*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "read_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "user_id", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 5, name: "title", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 6, name: "type", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "128" } } } },
            { no: 7, name: "content", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "256" } } } },
            { no: 8, name: "data", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "256" } } } }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.notifications.Notification
 */
export const Notification = new Notification$Type();
