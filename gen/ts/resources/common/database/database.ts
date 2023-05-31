// @generated by protobuf-ts 2.9.0 with parameter optimize_code_size,long_type_bigint
// @generated from protobuf file "resources/common/database/database.proto" (package "resources.common.database", syntax proto3)
// tslint:disable
import { MessageType } from "@protobuf-ts/runtime";
/**
 * @generated from protobuf message resources.common.database.PaginationRequest
 */
export interface PaginationRequest {
    /**
     * @generated from protobuf field: int64 offset = 1;
     */
    offset: bigint;
    /**
     * @generated from protobuf field: optional int64 page_size = 2;
     */
    pageSize?: bigint;
}
/**
 * @generated from protobuf message resources.common.database.PaginationResponse
 */
export interface PaginationResponse {
    /**
     * @generated from protobuf field: int64 total_count = 1;
     */
    totalCount: bigint;
    /**
     * @generated from protobuf field: int64 offset = 2;
     */
    offset: bigint;
    /**
     * @generated from protobuf field: int64 end = 3;
     */
    end: bigint;
    /**
     * @generated from protobuf field: int64 page_size = 4;
     */
    pageSize: bigint;
}
/**
 * @generated from protobuf message resources.common.database.OrderBy
 */
export interface OrderBy {
    /**
     * @generated from protobuf field: string column = 1;
     */
    column: string;
    /**
     * @generated from protobuf field: bool desc = 2;
     */
    desc: boolean;
}
// @generated message type with reflection information, may provide speed optimized methods
class PaginationRequest$Type extends MessageType<PaginationRequest> {
    constructor() {
        super("resources.common.database.PaginationRequest", [
            { no: 1, name: "offset", kind: "scalar", T: 3 /*ScalarType.INT64*/, L: 0 /*LongType.BIGINT*/, options: { "validate.rules": { int64: { gte: "0" } } } },
            { no: 2, name: "page_size", kind: "scalar", opt: true, T: 3 /*ScalarType.INT64*/, L: 0 /*LongType.BIGINT*/, options: { "validate.rules": { int64: { gte: "0" } } } }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.common.database.PaginationRequest
 */
export const PaginationRequest = new PaginationRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PaginationResponse$Type extends MessageType<PaginationResponse> {
    constructor() {
        super("resources.common.database.PaginationResponse", [
            { no: 1, name: "total_count", kind: "scalar", T: 3 /*ScalarType.INT64*/, L: 0 /*LongType.BIGINT*/ },
            { no: 2, name: "offset", kind: "scalar", T: 3 /*ScalarType.INT64*/, L: 0 /*LongType.BIGINT*/ },
            { no: 3, name: "end", kind: "scalar", T: 3 /*ScalarType.INT64*/, L: 0 /*LongType.BIGINT*/ },
            { no: 4, name: "page_size", kind: "scalar", T: 3 /*ScalarType.INT64*/, L: 0 /*LongType.BIGINT*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.common.database.PaginationResponse
 */
export const PaginationResponse = new PaginationResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class OrderBy$Type extends MessageType<OrderBy> {
    constructor() {
        super("resources.common.database.OrderBy", [
            { no: 1, name: "column", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "1", maxLen: "64" } } } },
            { no: 2, name: "desc", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.common.database.OrderBy
 */
export const OrderBy = new OrderBy$Type();
