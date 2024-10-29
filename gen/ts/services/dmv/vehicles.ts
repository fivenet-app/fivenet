// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/dmv/vehicles.proto" (package "services.dmv", syntax proto3)
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
import { Vehicle } from "../../resources/vehicles/vehicles";
import { PaginationResponse } from "../../resources/common/database/database";
import { Sort } from "../../resources/common/database/database";
import { PaginationRequest } from "../../resources/common/database/database";
/**
 * @generated from protobuf message services.dmv.ListVehiclesRequest
 */
export interface ListVehiclesRequest {
    /**
     * @generated from protobuf field: resources.common.database.PaginationRequest pagination = 1;
     */
    pagination?: PaginationRequest;
    /**
     * @generated from protobuf field: optional resources.common.database.Sort sort = 2;
     */
    sort?: Sort;
    /**
     * Search params
     *
     * @generated from protobuf field: optional string license_plate = 3;
     */
    licensePlate?: string;
    /**
     * @generated from protobuf field: optional string model = 4;
     */
    model?: string;
    /**
     * @generated from protobuf field: optional int32 user_id = 5;
     */
    userId?: number;
}
/**
 * @generated from protobuf message services.dmv.ListVehiclesResponse
 */
export interface ListVehiclesResponse {
    /**
     * @generated from protobuf field: resources.common.database.PaginationResponse pagination = 1;
     */
    pagination?: PaginationResponse;
    /**
     * @generated from protobuf field: repeated resources.vehicles.Vehicle vehicles = 2;
     */
    vehicles: Vehicle[];
}
// @generated message type with reflection information, may provide speed optimized methods
class ListVehiclesRequest$Type extends MessageType<ListVehiclesRequest> {
    constructor() {
        super("services.dmv.ListVehiclesRequest", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationRequest, options: { "validate.rules": { message: { required: true } } } },
            { no: 2, name: "sort", kind: "message", T: () => Sort },
            { no: 3, name: "license_plate", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "32" } } } },
            { no: 4, name: "model", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "32" } } } },
            { no: 5, name: "user_id", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gte: 0 } } } }
        ]);
    }
    create(value?: PartialMessage<ListVehiclesRequest>): ListVehiclesRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<ListVehiclesRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListVehiclesRequest): ListVehiclesRequest {
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
                case /* optional string license_plate */ 3:
                    message.licensePlate = reader.string();
                    break;
                case /* optional string model */ 4:
                    message.model = reader.string();
                    break;
                case /* optional int32 user_id */ 5:
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
    internalBinaryWrite(message: ListVehiclesRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationRequest pagination = 1; */
        if (message.pagination)
            PaginationRequest.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.common.database.Sort sort = 2; */
        if (message.sort)
            Sort.internalBinaryWrite(message.sort, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* optional string license_plate = 3; */
        if (message.licensePlate !== undefined)
            writer.tag(3, WireType.LengthDelimited).string(message.licensePlate);
        /* optional string model = 4; */
        if (message.model !== undefined)
            writer.tag(4, WireType.LengthDelimited).string(message.model);
        /* optional int32 user_id = 5; */
        if (message.userId !== undefined)
            writer.tag(5, WireType.Varint).int32(message.userId);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.dmv.ListVehiclesRequest
 */
export const ListVehiclesRequest = new ListVehiclesRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ListVehiclesResponse$Type extends MessageType<ListVehiclesResponse> {
    constructor() {
        super("services.dmv.ListVehiclesResponse", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationResponse },
            { no: 2, name: "vehicles", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Vehicle }
        ]);
    }
    create(value?: PartialMessage<ListVehiclesResponse>): ListVehiclesResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.vehicles = [];
        if (value !== undefined)
            reflectionMergePartial<ListVehiclesResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ListVehiclesResponse): ListVehiclesResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationResponse pagination */ 1:
                    message.pagination = PaginationResponse.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* repeated resources.vehicles.Vehicle vehicles */ 2:
                    message.vehicles.push(Vehicle.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: ListVehiclesResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationResponse pagination = 1; */
        if (message.pagination)
            PaginationResponse.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.vehicles.Vehicle vehicles = 2; */
        for (let i = 0; i < message.vehicles.length; i++)
            Vehicle.internalBinaryWrite(message.vehicles[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.dmv.ListVehiclesResponse
 */
export const ListVehiclesResponse = new ListVehiclesResponse$Type();
/**
 * @generated ServiceType for protobuf service services.dmv.DMVService
 */
export const DMVService = new ServiceType("services.dmv.DMVService", [
    { name: "ListVehicles", options: {}, I: ListVehiclesRequest, O: ListVehiclesResponse }
]);
