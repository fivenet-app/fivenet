// @generated by protobuf-ts 2.9.6 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/vehicles/vehicles.proto" (package "resources.vehicles", syntax proto3)
// @ts-nocheck
import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import { WireType } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import { UnknownFieldHandler } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { reflectionMergePartial } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
import { UserShort } from "../users/users";
/**
 * @generated from protobuf message resources.vehicles.Vehicle
 */
export interface Vehicle {
    /**
     * @generated from protobuf field: string plate = 1;
     */
    plate: string;
    /**
     * @generated from protobuf field: optional string model = 2;
     */
    model?: string;
    /**
     * @generated from protobuf field: string type = 3;
     */
    type: string;
    /**
     * @generated from protobuf field: optional int32 owner_id = 4;
     */
    ownerId?: number;
    /**
     * @generated from protobuf field: optional string owner_identifier = 6;
     */
    ownerIdentifier?: string;
    /**
     * @generated from protobuf field: optional resources.users.UserShort owner = 5;
     */
    owner?: UserShort;
    /**
     * @generated from protobuf field: optional string job = 7;
     */
    job?: string;
    /**
     * @generated from protobuf field: optional string job_label = 8;
     */
    jobLabel?: string;
}
// @generated message type with reflection information, may provide speed optimized methods
class Vehicle$Type extends MessageType<Vehicle> {
    constructor() {
        super("resources.vehicles.Vehicle", [
            { no: 1, name: "plate", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "32" } } } },
            { no: 2, name: "model", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "64" } } } },
            { no: 3, name: "type", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "32" } } } },
            { no: 4, name: "owner_id", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/ },
            { no: 6, name: "owner_identifier", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "64" } } } },
            { no: 5, name: "owner", kind: "message", T: () => UserShort },
            { no: 7, name: "job", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 8, name: "job_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } }
        ]);
    }
    create(value?: PartialMessage<Vehicle>): Vehicle {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.plate = "";
        message.type = "";
        if (value !== undefined)
            reflectionMergePartial<Vehicle>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Vehicle): Vehicle {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string plate */ 1:
                    message.plate = reader.string();
                    break;
                case /* optional string model */ 2:
                    message.model = reader.string();
                    break;
                case /* string type */ 3:
                    message.type = reader.string();
                    break;
                case /* optional int32 owner_id */ 4:
                    message.ownerId = reader.int32();
                    break;
                case /* optional string owner_identifier */ 6:
                    message.ownerIdentifier = reader.string();
                    break;
                case /* optional resources.users.UserShort owner */ 5:
                    message.owner = UserShort.internalBinaryRead(reader, reader.uint32(), options, message.owner);
                    break;
                case /* optional string job */ 7:
                    message.job = reader.string();
                    break;
                case /* optional string job_label */ 8:
                    message.jobLabel = reader.string();
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
    internalBinaryWrite(message: Vehicle, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string plate = 1; */
        if (message.plate !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.plate);
        /* optional string model = 2; */
        if (message.model !== undefined)
            writer.tag(2, WireType.LengthDelimited).string(message.model);
        /* string type = 3; */
        if (message.type !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.type);
        /* optional int32 owner_id = 4; */
        if (message.ownerId !== undefined)
            writer.tag(4, WireType.Varint).int32(message.ownerId);
        /* optional string owner_identifier = 6; */
        if (message.ownerIdentifier !== undefined)
            writer.tag(6, WireType.LengthDelimited).string(message.ownerIdentifier);
        /* optional resources.users.UserShort owner = 5; */
        if (message.owner)
            UserShort.internalBinaryWrite(message.owner, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* optional string job = 7; */
        if (message.job !== undefined)
            writer.tag(7, WireType.LengthDelimited).string(message.job);
        /* optional string job_label = 8; */
        if (message.jobLabel !== undefined)
            writer.tag(8, WireType.LengthDelimited).string(message.jobLabel);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.vehicles.Vehicle
 */
export const Vehicle = new Vehicle$Type();
