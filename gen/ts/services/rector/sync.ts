// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/rector/sync.proto" (package "services.rector", syntax proto3)
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
/**
 * @generated from protobuf message services.rector.SyncRequest
 */
export interface SyncRequest {
    /**
     * @generated from protobuf oneof: data
     */
    data: {
        oneofKind: "test";
        /**
         * @generated from protobuf field: bool test = 1;
         */
        test: boolean; // TODO create per table message and add oneof entry
    } | {
        oneofKind: undefined;
    };
}
/**
 * @generated from protobuf message services.rector.SyncResponse
 */
export interface SyncResponse {
    /**
     * @generated from protobuf field: int64 affected_rows = 1;
     */
    affectedRows: number;
}
// @generated message type with reflection information, may provide speed optimized methods
class SyncRequest$Type extends MessageType<SyncRequest> {
    constructor() {
        super("services.rector.SyncRequest", [
            { no: 1, name: "test", kind: "scalar", oneof: "data", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<SyncRequest>): SyncRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.data = { oneofKind: undefined };
        if (value !== undefined)
            reflectionMergePartial<SyncRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SyncRequest): SyncRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* bool test */ 1:
                    message.data = {
                        oneofKind: "test",
                        test: reader.bool()
                    };
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
    internalBinaryWrite(message: SyncRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* bool test = 1; */
        if (message.data.oneofKind === "test")
            writer.tag(1, WireType.Varint).bool(message.data.test);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.rector.SyncRequest
 */
export const SyncRequest = new SyncRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class SyncResponse$Type extends MessageType<SyncResponse> {
    constructor() {
        super("services.rector.SyncResponse", [
            { no: 1, name: "affected_rows", kind: "scalar", T: 3 /*ScalarType.INT64*/, L: 2 /*LongType.NUMBER*/ }
        ]);
    }
    create(value?: PartialMessage<SyncResponse>): SyncResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.affectedRows = 0;
        if (value !== undefined)
            reflectionMergePartial<SyncResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SyncResponse): SyncResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int64 affected_rows */ 1:
                    message.affectedRows = reader.int64().toNumber();
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
    internalBinaryWrite(message: SyncResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* int64 affected_rows = 1; */
        if (message.affectedRows !== 0)
            writer.tag(1, WireType.Varint).int64(message.affectedRows);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.rector.SyncResponse
 */
export const SyncResponse = new SyncResponse$Type();
/**
 * @generated ServiceType for protobuf service services.rector.SyncService
 */
export const SyncService = new ServiceType("services.rector.SyncService", [
    { name: "Sync", options: {}, I: SyncRequest, O: SyncResponse }
]);
