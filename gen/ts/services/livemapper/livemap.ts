// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/livemapper/livemap.proto" (package "services.livemapper", syntax proto3)
// tslint:disable
import { ServiceType } from "@protobuf-ts/runtime-rpc";
import { WireType } from "@protobuf-ts/runtime";
import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import { UnknownFieldHandler } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { reflectionMergePartial } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
import { UserMarker } from "../../resources/livemap/livemap";
import { MarkerMarker } from "../../resources/livemap/livemap";
import { Job } from "../../resources/users/jobs";
/**
 * @generated from protobuf message services.livemapper.StreamRequest
 */
export interface StreamRequest {
}
/**
 * @generated from protobuf message services.livemapper.StreamResponse
 */
export interface StreamResponse {
    /**
     * @generated from protobuf oneof: data
     */
    data: {
        oneofKind: "jobs";
        /**
         * @generated from protobuf field: services.livemapper.JobsList jobs = 1;
         */
        jobs: JobsList;
    } | {
        oneofKind: "markers";
        /**
         * @generated from protobuf field: services.livemapper.MarkerMarkersUpdates markers = 2;
         */
        markers: MarkerMarkersUpdates;
    } | {
        oneofKind: "users";
        /**
         * @generated from protobuf field: services.livemapper.UserMarkersUpdates users = 3;
         */
        users: UserMarkersUpdates;
    } | {
        oneofKind: undefined;
    };
}
/**
 * @generated from protobuf message services.livemapper.JobsList
 */
export interface JobsList {
    /**
     * @generated from protobuf field: repeated resources.users.Job users = 1;
     */
    users: Job[];
    /**
     * @generated from protobuf field: repeated resources.users.Job markers = 2;
     */
    markers: Job[];
}
/**
 * @generated from protobuf message services.livemapper.MarkerMarkersUpdates
 */
export interface MarkerMarkersUpdates {
    /**
     * @generated from protobuf field: repeated resources.livemap.MarkerMarker markers = 1;
     */
    markers: MarkerMarker[];
}
/**
 * @generated from protobuf message services.livemapper.UserMarkersUpdates
 */
export interface UserMarkersUpdates {
    /**
     * @generated from protobuf field: repeated resources.livemap.UserMarker users = 1;
     */
    users: UserMarker[];
    /**
     * @generated from protobuf field: int32 part = 2;
     */
    part: number;
}
/**
 * @generated from protobuf message services.livemapper.CreateOrUpdateMarkerRequest
 */
export interface CreateOrUpdateMarkerRequest {
    /**
     * @generated from protobuf field: resources.livemap.MarkerMarker marker = 1;
     */
    marker?: MarkerMarker;
}
/**
 * @generated from protobuf message services.livemapper.CreateOrUpdateMarkerResponse
 */
export interface CreateOrUpdateMarkerResponse {
    /**
     * @generated from protobuf field: resources.livemap.MarkerMarker marker = 1;
     */
    marker?: MarkerMarker;
}
/**
 * @generated from protobuf message services.livemapper.DeleteMarkerRequest
 */
export interface DeleteMarkerRequest {
    /**
     * @generated from protobuf field: uint64 id = 1 [jstype = JS_STRING];
     */
    id: string;
}
/**
 * @generated from protobuf message services.livemapper.DeleteMarkerResponse
 */
export interface DeleteMarkerResponse {
}
// @generated message type with reflection information, may provide speed optimized methods
class StreamRequest$Type extends MessageType<StreamRequest> {
    constructor() {
        super("services.livemapper.StreamRequest", []);
    }
    create(value?: PartialMessage<StreamRequest>): StreamRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<StreamRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: StreamRequest): StreamRequest {
        return target ?? this.create();
    }
    internalBinaryWrite(message: StreamRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.livemapper.StreamRequest
 */
export const StreamRequest = new StreamRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class StreamResponse$Type extends MessageType<StreamResponse> {
    constructor() {
        super("services.livemapper.StreamResponse", [
            { no: 1, name: "jobs", kind: "message", oneof: "data", T: () => JobsList },
            { no: 2, name: "markers", kind: "message", oneof: "data", T: () => MarkerMarkersUpdates },
            { no: 3, name: "users", kind: "message", oneof: "data", T: () => UserMarkersUpdates }
        ]);
    }
    create(value?: PartialMessage<StreamResponse>): StreamResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.data = { oneofKind: undefined };
        if (value !== undefined)
            reflectionMergePartial<StreamResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: StreamResponse): StreamResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* services.livemapper.JobsList jobs */ 1:
                    message.data = {
                        oneofKind: "jobs",
                        jobs: JobsList.internalBinaryRead(reader, reader.uint32(), options, (message.data as any).jobs)
                    };
                    break;
                case /* services.livemapper.MarkerMarkersUpdates markers */ 2:
                    message.data = {
                        oneofKind: "markers",
                        markers: MarkerMarkersUpdates.internalBinaryRead(reader, reader.uint32(), options, (message.data as any).markers)
                    };
                    break;
                case /* services.livemapper.UserMarkersUpdates users */ 3:
                    message.data = {
                        oneofKind: "users",
                        users: UserMarkersUpdates.internalBinaryRead(reader, reader.uint32(), options, (message.data as any).users)
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
    internalBinaryWrite(message: StreamResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* services.livemapper.JobsList jobs = 1; */
        if (message.data.oneofKind === "jobs")
            JobsList.internalBinaryWrite(message.data.jobs, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* services.livemapper.MarkerMarkersUpdates markers = 2; */
        if (message.data.oneofKind === "markers")
            MarkerMarkersUpdates.internalBinaryWrite(message.data.markers, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* services.livemapper.UserMarkersUpdates users = 3; */
        if (message.data.oneofKind === "users")
            UserMarkersUpdates.internalBinaryWrite(message.data.users, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.livemapper.StreamResponse
 */
export const StreamResponse = new StreamResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class JobsList$Type extends MessageType<JobsList> {
    constructor() {
        super("services.livemapper.JobsList", [
            { no: 1, name: "users", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Job },
            { no: 2, name: "markers", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Job }
        ]);
    }
    create(value?: PartialMessage<JobsList>): JobsList {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.users = [];
        message.markers = [];
        if (value !== undefined)
            reflectionMergePartial<JobsList>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: JobsList): JobsList {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.users.Job users */ 1:
                    message.users.push(Job.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated resources.users.Job markers */ 2:
                    message.markers.push(Job.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: JobsList, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.users.Job users = 1; */
        for (let i = 0; i < message.users.length; i++)
            Job.internalBinaryWrite(message.users[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.users.Job markers = 2; */
        for (let i = 0; i < message.markers.length; i++)
            Job.internalBinaryWrite(message.markers[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.livemapper.JobsList
 */
export const JobsList = new JobsList$Type();
// @generated message type with reflection information, may provide speed optimized methods
class MarkerMarkersUpdates$Type extends MessageType<MarkerMarkersUpdates> {
    constructor() {
        super("services.livemapper.MarkerMarkersUpdates", [
            { no: 1, name: "markers", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => MarkerMarker }
        ]);
    }
    create(value?: PartialMessage<MarkerMarkersUpdates>): MarkerMarkersUpdates {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.markers = [];
        if (value !== undefined)
            reflectionMergePartial<MarkerMarkersUpdates>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: MarkerMarkersUpdates): MarkerMarkersUpdates {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.livemap.MarkerMarker markers */ 1:
                    message.markers.push(MarkerMarker.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: MarkerMarkersUpdates, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.livemap.MarkerMarker markers = 1; */
        for (let i = 0; i < message.markers.length; i++)
            MarkerMarker.internalBinaryWrite(message.markers[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.livemapper.MarkerMarkersUpdates
 */
export const MarkerMarkersUpdates = new MarkerMarkersUpdates$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UserMarkersUpdates$Type extends MessageType<UserMarkersUpdates> {
    constructor() {
        super("services.livemapper.UserMarkersUpdates", [
            { no: 1, name: "users", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => UserMarker },
            { no: 2, name: "part", kind: "scalar", T: 5 /*ScalarType.INT32*/ }
        ]);
    }
    create(value?: PartialMessage<UserMarkersUpdates>): UserMarkersUpdates {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.users = [];
        message.part = 0;
        if (value !== undefined)
            reflectionMergePartial<UserMarkersUpdates>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UserMarkersUpdates): UserMarkersUpdates {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.livemap.UserMarker users */ 1:
                    message.users.push(UserMarker.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* int32 part */ 2:
                    message.part = reader.int32();
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
    internalBinaryWrite(message: UserMarkersUpdates, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.livemap.UserMarker users = 1; */
        for (let i = 0; i < message.users.length; i++)
            UserMarker.internalBinaryWrite(message.users[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* int32 part = 2; */
        if (message.part !== 0)
            writer.tag(2, WireType.Varint).int32(message.part);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.livemapper.UserMarkersUpdates
 */
export const UserMarkersUpdates = new UserMarkersUpdates$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateOrUpdateMarkerRequest$Type extends MessageType<CreateOrUpdateMarkerRequest> {
    constructor() {
        super("services.livemapper.CreateOrUpdateMarkerRequest", [
            { no: 1, name: "marker", kind: "message", T: () => MarkerMarker, options: { "validate.rules": { message: { required: true } } } }
        ]);
    }
    create(value?: PartialMessage<CreateOrUpdateMarkerRequest>): CreateOrUpdateMarkerRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<CreateOrUpdateMarkerRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CreateOrUpdateMarkerRequest): CreateOrUpdateMarkerRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.livemap.MarkerMarker marker */ 1:
                    message.marker = MarkerMarker.internalBinaryRead(reader, reader.uint32(), options, message.marker);
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
    internalBinaryWrite(message: CreateOrUpdateMarkerRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.livemap.MarkerMarker marker = 1; */
        if (message.marker)
            MarkerMarker.internalBinaryWrite(message.marker, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.livemapper.CreateOrUpdateMarkerRequest
 */
export const CreateOrUpdateMarkerRequest = new CreateOrUpdateMarkerRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateOrUpdateMarkerResponse$Type extends MessageType<CreateOrUpdateMarkerResponse> {
    constructor() {
        super("services.livemapper.CreateOrUpdateMarkerResponse", [
            { no: 1, name: "marker", kind: "message", T: () => MarkerMarker }
        ]);
    }
    create(value?: PartialMessage<CreateOrUpdateMarkerResponse>): CreateOrUpdateMarkerResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<CreateOrUpdateMarkerResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CreateOrUpdateMarkerResponse): CreateOrUpdateMarkerResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.livemap.MarkerMarker marker */ 1:
                    message.marker = MarkerMarker.internalBinaryRead(reader, reader.uint32(), options, message.marker);
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
    internalBinaryWrite(message: CreateOrUpdateMarkerResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.livemap.MarkerMarker marker = 1; */
        if (message.marker)
            MarkerMarker.internalBinaryWrite(message.marker, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.livemapper.CreateOrUpdateMarkerResponse
 */
export const CreateOrUpdateMarkerResponse = new CreateOrUpdateMarkerResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteMarkerRequest$Type extends MessageType<DeleteMarkerRequest> {
    constructor() {
        super("services.livemapper.DeleteMarkerRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ }
        ]);
    }
    create(value?: PartialMessage<DeleteMarkerRequest>): DeleteMarkerRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = "0";
        if (value !== undefined)
            reflectionMergePartial<DeleteMarkerRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteMarkerRequest): DeleteMarkerRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id = 1 [jstype = JS_STRING];*/ 1:
                    message.id = reader.uint64().toString();
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
    internalBinaryWrite(message: DeleteMarkerRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1 [jstype = JS_STRING]; */
        if (message.id !== "0")
            writer.tag(1, WireType.Varint).uint64(message.id);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.livemapper.DeleteMarkerRequest
 */
export const DeleteMarkerRequest = new DeleteMarkerRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteMarkerResponse$Type extends MessageType<DeleteMarkerResponse> {
    constructor() {
        super("services.livemapper.DeleteMarkerResponse", []);
    }
    create(value?: PartialMessage<DeleteMarkerResponse>): DeleteMarkerResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<DeleteMarkerResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DeleteMarkerResponse): DeleteMarkerResponse {
        return target ?? this.create();
    }
    internalBinaryWrite(message: DeleteMarkerResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.livemapper.DeleteMarkerResponse
 */
export const DeleteMarkerResponse = new DeleteMarkerResponse$Type();
/**
 * @generated ServiceType for protobuf service services.livemapper.LivemapperService
 */
export const LivemapperService = new ServiceType("services.livemapper.LivemapperService", [
    { name: "Stream", serverStreaming: true, options: {}, I: StreamRequest, O: StreamResponse },
    { name: "CreateOrUpdateMarker", options: {}, I: CreateOrUpdateMarkerRequest, O: CreateOrUpdateMarkerResponse },
    { name: "DeleteMarker", options: {}, I: DeleteMarkerRequest, O: DeleteMarkerResponse }
]);
