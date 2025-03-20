// @generated by protobuf-ts 2.9.6 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/sync/data.proto" (package "resources.sync", syntax proto3)
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
import { Coords } from "../livemap/livemap";
import { License } from "../users/licenses";
import { Vehicle } from "../vehicles/vehicles";
import { User } from "../users/users";
import { Job } from "../users/jobs";
/**
 * @generated from protobuf message resources.sync.DataStatus
 */
export interface DataStatus {
    /**
     * @generated from protobuf field: int64 count = 1;
     */
    count: number;
}
/**
 * @generated from protobuf message resources.sync.DataJobs
 */
export interface DataJobs {
    /**
     * @generated from protobuf field: repeated resources.users.Job jobs = 1;
     */
    jobs: Job[];
}
/**
 * @generated from protobuf message resources.sync.DataUsers
 */
export interface DataUsers {
    /**
     * @generated from protobuf field: repeated resources.users.User users = 1;
     */
    users: User[];
}
/**
 * @generated from protobuf message resources.sync.DataVehicles
 */
export interface DataVehicles {
    /**
     * @generated from protobuf field: repeated resources.vehicles.Vehicle vehicles = 1;
     */
    vehicles: Vehicle[];
}
/**
 * @generated from protobuf message resources.sync.DataLicenses
 */
export interface DataLicenses {
    /**
     * @generated from protobuf field: repeated resources.users.License licenses = 1;
     */
    licenses: License[];
}
/**
 * @generated from protobuf message resources.sync.DataUserLocations
 */
export interface DataUserLocations {
    /**
     * @generated from protobuf field: repeated resources.sync.UserLocation users = 1;
     */
    users: UserLocation[];
    /**
     * @generated from protobuf field: optional bool clear_all = 2;
     */
    clearAll?: boolean;
}
/**
 * @generated from protobuf message resources.sync.UserLocation
 */
export interface UserLocation {
    /**
     * @generated from protobuf field: string identifier = 1;
     */
    identifier: string;
    /**
     * @generated from protobuf field: string job = 2;
     */
    job: string;
    /**
     * @generated from protobuf field: resources.livemap.Coords coords = 3;
     */
    coords?: Coords;
    /**
     * @generated from protobuf field: bool hidden = 4;
     */
    hidden: boolean;
    /**
     * @generated from protobuf field: bool remove = 5;
     */
    remove: boolean;
}
/**
 * @generated from protobuf message resources.sync.DeleteUsers
 */
export interface DeleteUsers {
    /**
     * @generated from protobuf field: repeated int32 user_ids = 1;
     */
    userIds: number[];
}
/**
 * @generated from protobuf message resources.sync.DeleteVehicles
 */
export interface DeleteVehicles {
    /**
     * @generated from protobuf field: repeated string plates = 1;
     */
    plates: string[];
}
// @generated message type with reflection information, may provide speed optimized methods
class DataStatus$Type extends MessageType<DataStatus> {
    constructor() {
        super("resources.sync.DataStatus", [
            { no: 1, name: "count", kind: "scalar", T: 3 /*ScalarType.INT64*/, L: 2 /*LongType.NUMBER*/ }
        ]);
    }
    create(value?: PartialMessage<DataStatus>): DataStatus {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.count = 0;
        if (value !== undefined)
            reflectionMergePartial<DataStatus>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DataStatus): DataStatus {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int64 count */ 1:
                    message.count = reader.int64().toNumber();
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
    internalBinaryWrite(message: DataStatus, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* int64 count = 1; */
        if (message.count !== 0)
            writer.tag(1, WireType.Varint).int64(message.count);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.sync.DataStatus
 */
export const DataStatus = new DataStatus$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DataJobs$Type extends MessageType<DataJobs> {
    constructor() {
        super("resources.sync.DataJobs", [
            { no: 1, name: "jobs", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Job, options: { "validate.rules": { repeated: { maxItems: "200" } } } }
        ]);
    }
    create(value?: PartialMessage<DataJobs>): DataJobs {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.jobs = [];
        if (value !== undefined)
            reflectionMergePartial<DataJobs>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DataJobs): DataJobs {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.users.Job jobs */ 1:
                    message.jobs.push(Job.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: DataJobs, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.users.Job jobs = 1; */
        for (let i = 0; i < message.jobs.length; i++)
            Job.internalBinaryWrite(message.jobs[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.sync.DataJobs
 */
export const DataJobs = new DataJobs$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DataUsers$Type extends MessageType<DataUsers> {
    constructor() {
        super("resources.sync.DataUsers", [
            { no: 1, name: "users", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => User, options: { "validate.rules": { repeated: { maxItems: "500" } } } }
        ]);
    }
    create(value?: PartialMessage<DataUsers>): DataUsers {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.users = [];
        if (value !== undefined)
            reflectionMergePartial<DataUsers>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DataUsers): DataUsers {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.users.User users */ 1:
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
    internalBinaryWrite(message: DataUsers, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.users.User users = 1; */
        for (let i = 0; i < message.users.length; i++)
            User.internalBinaryWrite(message.users[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.sync.DataUsers
 */
export const DataUsers = new DataUsers$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DataVehicles$Type extends MessageType<DataVehicles> {
    constructor() {
        super("resources.sync.DataVehicles", [
            { no: 1, name: "vehicles", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Vehicle, options: { "validate.rules": { repeated: { maxItems: "1000" } } } }
        ]);
    }
    create(value?: PartialMessage<DataVehicles>): DataVehicles {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.vehicles = [];
        if (value !== undefined)
            reflectionMergePartial<DataVehicles>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DataVehicles): DataVehicles {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.vehicles.Vehicle vehicles */ 1:
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
    internalBinaryWrite(message: DataVehicles, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.vehicles.Vehicle vehicles = 1; */
        for (let i = 0; i < message.vehicles.length; i++)
            Vehicle.internalBinaryWrite(message.vehicles[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.sync.DataVehicles
 */
export const DataVehicles = new DataVehicles$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DataLicenses$Type extends MessageType<DataLicenses> {
    constructor() {
        super("resources.sync.DataLicenses", [
            { no: 1, name: "licenses", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => License, options: { "validate.rules": { repeated: { maxItems: "200" } } } }
        ]);
    }
    create(value?: PartialMessage<DataLicenses>): DataLicenses {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.licenses = [];
        if (value !== undefined)
            reflectionMergePartial<DataLicenses>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DataLicenses): DataLicenses {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.users.License licenses */ 1:
                    message.licenses.push(License.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: DataLicenses, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.users.License licenses = 1; */
        for (let i = 0; i < message.licenses.length; i++)
            License.internalBinaryWrite(message.licenses[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.sync.DataLicenses
 */
export const DataLicenses = new DataLicenses$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DataUserLocations$Type extends MessageType<DataUserLocations> {
    constructor() {
        super("resources.sync.DataUserLocations", [
            { no: 1, name: "users", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => UserLocation, options: { "validate.rules": { repeated: { maxItems: "2000" } } } },
            { no: 2, name: "clear_all", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<DataUserLocations>): DataUserLocations {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.users = [];
        if (value !== undefined)
            reflectionMergePartial<DataUserLocations>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DataUserLocations): DataUserLocations {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.sync.UserLocation users */ 1:
                    message.users.push(UserLocation.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* optional bool clear_all */ 2:
                    message.clearAll = reader.bool();
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
    internalBinaryWrite(message: DataUserLocations, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.sync.UserLocation users = 1; */
        for (let i = 0; i < message.users.length; i++)
            UserLocation.internalBinaryWrite(message.users[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* optional bool clear_all = 2; */
        if (message.clearAll !== undefined)
            writer.tag(2, WireType.Varint).bool(message.clearAll);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.sync.DataUserLocations
 */
export const DataUserLocations = new DataUserLocations$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UserLocation$Type extends MessageType<UserLocation> {
    constructor() {
        super("resources.sync.UserLocation", [
            { no: 1, name: "identifier", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "64" } } } },
            { no: 2, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 3, name: "coords", kind: "message", T: () => Coords, options: { "validate.rules": { message: { required: true } } } },
            { no: 4, name: "hidden", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 5, name: "remove", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<UserLocation>): UserLocation {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.identifier = "";
        message.job = "";
        message.hidden = false;
        message.remove = false;
        if (value !== undefined)
            reflectionMergePartial<UserLocation>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UserLocation): UserLocation {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string identifier */ 1:
                    message.identifier = reader.string();
                    break;
                case /* string job */ 2:
                    message.job = reader.string();
                    break;
                case /* resources.livemap.Coords coords */ 3:
                    message.coords = Coords.internalBinaryRead(reader, reader.uint32(), options, message.coords);
                    break;
                case /* bool hidden */ 4:
                    message.hidden = reader.bool();
                    break;
                case /* bool remove */ 5:
                    message.remove = reader.bool();
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
    internalBinaryWrite(message: UserLocation, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string identifier = 1; */
        if (message.identifier !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.identifier);
        /* string job = 2; */
        if (message.job !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.job);
        /* resources.livemap.Coords coords = 3; */
        if (message.coords)
            Coords.internalBinaryWrite(message.coords, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* bool hidden = 4; */
        if (message.hidden !== false)
            writer.tag(4, WireType.Varint).bool(message.hidden);
        /* bool remove = 5; */
        if (message.remove !== false)
            writer.tag(5, WireType.Varint).bool(message.remove);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.sync.UserLocation
 */
export const UserLocation = new UserLocation$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteUsers$Type extends MessageType<DeleteUsers> {
    constructor() {
        super("resources.sync.DeleteUsers", [
            { no: 1, name: "user_ids", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { repeated: { maxItems: "100" } } } }
        ]);
    }
    create(value?: PartialMessage<DeleteUsers>): DeleteUsers {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.userIds = [];
        if (value !== undefined)
            reflectionMergePartial<DeleteUsers>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteUsers): DeleteUsers {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated int32 user_ids */ 1:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.userIds.push(reader.int32());
                    else
                        message.userIds.push(reader.int32());
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
    internalBinaryWrite(message: DeleteUsers, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated int32 user_ids = 1; */
        if (message.userIds.length) {
            writer.tag(1, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.userIds.length; i++)
                writer.int32(message.userIds[i]);
            writer.join();
        }
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.sync.DeleteUsers
 */
export const DeleteUsers = new DeleteUsers$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteVehicles$Type extends MessageType<DeleteVehicles> {
    constructor() {
        super("resources.sync.DeleteVehicles", [
            { no: 1, name: "plates", kind: "scalar", repeat: 2 /*RepeatType.UNPACKED*/, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { repeated: { maxItems: "100" } } } }
        ]);
    }
    create(value?: PartialMessage<DeleteVehicles>): DeleteVehicles {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.plates = [];
        if (value !== undefined)
            reflectionMergePartial<DeleteVehicles>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteVehicles): DeleteVehicles {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated string plates */ 1:
                    message.plates.push(reader.string());
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
    internalBinaryWrite(message: DeleteVehicles, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated string plates = 1; */
        for (let i = 0; i < message.plates.length; i++)
            writer.tag(1, WireType.LengthDelimited).string(message.plates[i]);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.sync.DeleteVehicles
 */
export const DeleteVehicles = new DeleteVehicles$Type();
