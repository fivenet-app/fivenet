// @generated by protobuf-ts 2.9.5 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/common/grpcws/grpcws.proto" (package "resources.common.grpcws", syntax proto3)
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
/**
 * @generated from protobuf message resources.common.grpcws.GrpcFrame
 */
export interface GrpcFrame {
    /**
     * @generated from protobuf field: uint32 streamId = 1;
     */
    streamId: number;
    /**
     * @generated from protobuf oneof: payload
     */
    payload: {
        oneofKind: "ping";
        /**
         * @generated from protobuf field: resources.common.grpcws.Ping ping = 2;
         */
        ping: Ping;
    } | {
        oneofKind: "header";
        /**
         * @generated from protobuf field: resources.common.grpcws.Header header = 3;
         */
        header: Header;
    } | {
        oneofKind: "body";
        /**
         * @generated from protobuf field: resources.common.grpcws.Body body = 4;
         */
        body: Body;
    } | {
        oneofKind: "complete";
        /**
         * @generated from protobuf field: resources.common.grpcws.Complete complete = 5;
         */
        complete: Complete;
    } | {
        oneofKind: "failure";
        /**
         * @generated from protobuf field: resources.common.grpcws.Failure failure = 6;
         */
        failure: Failure;
    } | {
        oneofKind: "cancel";
        /**
         * @generated from protobuf field: resources.common.grpcws.Cancel cancel = 7;
         */
        cancel: Cancel;
    } | {
        oneofKind: undefined;
    };
}
/**
 * @generated from protobuf message resources.common.grpcws.Ping
 */
export interface Ping {
    /**
     * @generated from protobuf field: bool pong = 1;
     */
    pong: boolean;
}
/**
 * @generated from protobuf message resources.common.grpcws.Header
 */
export interface Header {
    /**
     * @generated from protobuf field: string operation = 1;
     */
    operation: string;
    /**
     * @generated from protobuf field: map<string, resources.common.grpcws.HeaderValue> headers = 2;
     */
    headers: {
        [key: string]: HeaderValue;
    };
    /**
     * @generated from protobuf field: int32 status = 3;
     */
    status: number;
}
/**
 * @generated from protobuf message resources.common.grpcws.HeaderValue
 */
export interface HeaderValue {
    /**
     * @generated from protobuf field: repeated string value = 1;
     */
    value: string[];
}
/**
 * @generated from protobuf message resources.common.grpcws.Body
 */
export interface Body {
    /**
     * @generated from protobuf field: bytes data = 1;
     */
    data: Uint8Array;
    /**
     * @generated from protobuf field: bool complete = 2;
     */
    complete: boolean;
}
/**
 * @generated from protobuf message resources.common.grpcws.Complete
 */
export interface Complete {
}
/**
 * @generated from protobuf message resources.common.grpcws.Failure
 */
export interface Failure {
    /**
     * @generated from protobuf field: string error_message = 1;
     */
    errorMessage: string;
    /**
     * @generated from protobuf field: string error_status = 2;
     */
    errorStatus: string;
    /**
     * @generated from protobuf field: map<string, resources.common.grpcws.HeaderValue> headers = 3;
     */
    headers: {
        [key: string]: HeaderValue;
    };
}
/**
 * @generated from protobuf message resources.common.grpcws.Cancel
 */
export interface Cancel {
}
// @generated message type with reflection information, may provide speed optimized methods
class GrpcFrame$Type extends MessageType<GrpcFrame> {
    constructor() {
        super("resources.common.grpcws.GrpcFrame", [
            { no: 1, name: "streamId", kind: "scalar", T: 13 /*ScalarType.UINT32*/ },
            { no: 2, name: "ping", kind: "message", oneof: "payload", T: () => Ping },
            { no: 3, name: "header", kind: "message", oneof: "payload", T: () => Header },
            { no: 4, name: "body", kind: "message", oneof: "payload", T: () => Body },
            { no: 5, name: "complete", kind: "message", oneof: "payload", T: () => Complete },
            { no: 6, name: "failure", kind: "message", oneof: "payload", T: () => Failure },
            { no: 7, name: "cancel", kind: "message", oneof: "payload", T: () => Cancel }
        ]);
    }
    create(value?: PartialMessage<GrpcFrame>): GrpcFrame {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.streamId = 0;
        message.payload = { oneofKind: undefined };
        if (value !== undefined)
            reflectionMergePartial<GrpcFrame>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: GrpcFrame): GrpcFrame {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint32 streamId */ 1:
                    message.streamId = reader.uint32();
                    break;
                case /* resources.common.grpcws.Ping ping */ 2:
                    message.payload = {
                        oneofKind: "ping",
                        ping: Ping.internalBinaryRead(reader, reader.uint32(), options, (message.payload as any).ping)
                    };
                    break;
                case /* resources.common.grpcws.Header header */ 3:
                    message.payload = {
                        oneofKind: "header",
                        header: Header.internalBinaryRead(reader, reader.uint32(), options, (message.payload as any).header)
                    };
                    break;
                case /* resources.common.grpcws.Body body */ 4:
                    message.payload = {
                        oneofKind: "body",
                        body: Body.internalBinaryRead(reader, reader.uint32(), options, (message.payload as any).body)
                    };
                    break;
                case /* resources.common.grpcws.Complete complete */ 5:
                    message.payload = {
                        oneofKind: "complete",
                        complete: Complete.internalBinaryRead(reader, reader.uint32(), options, (message.payload as any).complete)
                    };
                    break;
                case /* resources.common.grpcws.Failure failure */ 6:
                    message.payload = {
                        oneofKind: "failure",
                        failure: Failure.internalBinaryRead(reader, reader.uint32(), options, (message.payload as any).failure)
                    };
                    break;
                case /* resources.common.grpcws.Cancel cancel */ 7:
                    message.payload = {
                        oneofKind: "cancel",
                        cancel: Cancel.internalBinaryRead(reader, reader.uint32(), options, (message.payload as any).cancel)
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
    internalBinaryWrite(message: GrpcFrame, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint32 streamId = 1; */
        if (message.streamId !== 0)
            writer.tag(1, WireType.Varint).uint32(message.streamId);
        /* resources.common.grpcws.Ping ping = 2; */
        if (message.payload.oneofKind === "ping")
            Ping.internalBinaryWrite(message.payload.ping, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* resources.common.grpcws.Header header = 3; */
        if (message.payload.oneofKind === "header")
            Header.internalBinaryWrite(message.payload.header, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* resources.common.grpcws.Body body = 4; */
        if (message.payload.oneofKind === "body")
            Body.internalBinaryWrite(message.payload.body, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* resources.common.grpcws.Complete complete = 5; */
        if (message.payload.oneofKind === "complete")
            Complete.internalBinaryWrite(message.payload.complete, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* resources.common.grpcws.Failure failure = 6; */
        if (message.payload.oneofKind === "failure")
            Failure.internalBinaryWrite(message.payload.failure, writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        /* resources.common.grpcws.Cancel cancel = 7; */
        if (message.payload.oneofKind === "cancel")
            Cancel.internalBinaryWrite(message.payload.cancel, writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.common.grpcws.GrpcFrame
 */
export const GrpcFrame = new GrpcFrame$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Ping$Type extends MessageType<Ping> {
    constructor() {
        super("resources.common.grpcws.Ping", [
            { no: 1, name: "pong", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<Ping>): Ping {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.pong = false;
        if (value !== undefined)
            reflectionMergePartial<Ping>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Ping): Ping {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* bool pong */ 1:
                    message.pong = reader.bool();
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
    internalBinaryWrite(message: Ping, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* bool pong = 1; */
        if (message.pong !== false)
            writer.tag(1, WireType.Varint).bool(message.pong);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.common.grpcws.Ping
 */
export const Ping = new Ping$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Header$Type extends MessageType<Header> {
    constructor() {
        super("resources.common.grpcws.Header", [
            { no: 1, name: "operation", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "headers", kind: "map", K: 9 /*ScalarType.STRING*/, V: { kind: "message", T: () => HeaderValue } },
            { no: 3, name: "status", kind: "scalar", T: 5 /*ScalarType.INT32*/ }
        ]);
    }
    create(value?: PartialMessage<Header>): Header {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.operation = "";
        message.headers = {};
        message.status = 0;
        if (value !== undefined)
            reflectionMergePartial<Header>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Header): Header {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string operation */ 1:
                    message.operation = reader.string();
                    break;
                case /* map<string, resources.common.grpcws.HeaderValue> headers */ 2:
                    this.binaryReadMap2(message.headers, reader, options);
                    break;
                case /* int32 status */ 3:
                    message.status = reader.int32();
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
    private binaryReadMap2(map: Header["headers"], reader: IBinaryReader, options: BinaryReadOptions): void {
        let len = reader.uint32(), end = reader.pos + len, key: keyof Header["headers"] | undefined, val: Header["headers"][any] | undefined;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case 1:
                    key = reader.string();
                    break;
                case 2:
                    val = HeaderValue.internalBinaryRead(reader, reader.uint32(), options);
                    break;
                default: throw new globalThis.Error("unknown map entry field for field resources.common.grpcws.Header.headers");
            }
        }
        map[key ?? ""] = val ?? HeaderValue.create();
    }
    internalBinaryWrite(message: Header, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string operation = 1; */
        if (message.operation !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.operation);
        /* map<string, resources.common.grpcws.HeaderValue> headers = 2; */
        for (let k of globalThis.Object.keys(message.headers)) {
            writer.tag(2, WireType.LengthDelimited).fork().tag(1, WireType.LengthDelimited).string(k);
            writer.tag(2, WireType.LengthDelimited).fork();
            HeaderValue.internalBinaryWrite(message.headers[k], writer, options);
            writer.join().join();
        }
        /* int32 status = 3; */
        if (message.status !== 0)
            writer.tag(3, WireType.Varint).int32(message.status);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.common.grpcws.Header
 */
export const Header = new Header$Type();
// @generated message type with reflection information, may provide speed optimized methods
class HeaderValue$Type extends MessageType<HeaderValue> {
    constructor() {
        super("resources.common.grpcws.HeaderValue", [
            { no: 1, name: "value", kind: "scalar", repeat: 2 /*RepeatType.UNPACKED*/, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<HeaderValue>): HeaderValue {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.value = [];
        if (value !== undefined)
            reflectionMergePartial<HeaderValue>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: HeaderValue): HeaderValue {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated string value */ 1:
                    message.value.push(reader.string());
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
    internalBinaryWrite(message: HeaderValue, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated string value = 1; */
        for (let i = 0; i < message.value.length; i++)
            writer.tag(1, WireType.LengthDelimited).string(message.value[i]);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.common.grpcws.HeaderValue
 */
export const HeaderValue = new HeaderValue$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Body$Type extends MessageType<Body> {
    constructor() {
        super("resources.common.grpcws.Body", [
            { no: 1, name: "data", kind: "scalar", T: 12 /*ScalarType.BYTES*/ },
            { no: 2, name: "complete", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<Body>): Body {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.data = new Uint8Array(0);
        message.complete = false;
        if (value !== undefined)
            reflectionMergePartial<Body>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Body): Body {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* bytes data */ 1:
                    message.data = reader.bytes();
                    break;
                case /* bool complete */ 2:
                    message.complete = reader.bool();
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
    internalBinaryWrite(message: Body, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* bytes data = 1; */
        if (message.data.length)
            writer.tag(1, WireType.LengthDelimited).bytes(message.data);
        /* bool complete = 2; */
        if (message.complete !== false)
            writer.tag(2, WireType.Varint).bool(message.complete);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.common.grpcws.Body
 */
export const Body = new Body$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Complete$Type extends MessageType<Complete> {
    constructor() {
        super("resources.common.grpcws.Complete", []);
    }
    create(value?: PartialMessage<Complete>): Complete {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<Complete>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Complete): Complete {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
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
    internalBinaryWrite(message: Complete, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.common.grpcws.Complete
 */
export const Complete = new Complete$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Failure$Type extends MessageType<Failure> {
    constructor() {
        super("resources.common.grpcws.Failure", [
            { no: 1, name: "error_message", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "error_status", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "headers", kind: "map", K: 9 /*ScalarType.STRING*/, V: { kind: "message", T: () => HeaderValue } }
        ]);
    }
    create(value?: PartialMessage<Failure>): Failure {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.errorMessage = "";
        message.errorStatus = "";
        message.headers = {};
        if (value !== undefined)
            reflectionMergePartial<Failure>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Failure): Failure {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string error_message */ 1:
                    message.errorMessage = reader.string();
                    break;
                case /* string error_status */ 2:
                    message.errorStatus = reader.string();
                    break;
                case /* map<string, resources.common.grpcws.HeaderValue> headers */ 3:
                    this.binaryReadMap3(message.headers, reader, options);
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
    private binaryReadMap3(map: Failure["headers"], reader: IBinaryReader, options: BinaryReadOptions): void {
        let len = reader.uint32(), end = reader.pos + len, key: keyof Failure["headers"] | undefined, val: Failure["headers"][any] | undefined;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case 1:
                    key = reader.string();
                    break;
                case 2:
                    val = HeaderValue.internalBinaryRead(reader, reader.uint32(), options);
                    break;
                default: throw new globalThis.Error("unknown map entry field for field resources.common.grpcws.Failure.headers");
            }
        }
        map[key ?? ""] = val ?? HeaderValue.create();
    }
    internalBinaryWrite(message: Failure, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string error_message = 1; */
        if (message.errorMessage !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.errorMessage);
        /* string error_status = 2; */
        if (message.errorStatus !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.errorStatus);
        /* map<string, resources.common.grpcws.HeaderValue> headers = 3; */
        for (let k of globalThis.Object.keys(message.headers)) {
            writer.tag(3, WireType.LengthDelimited).fork().tag(1, WireType.LengthDelimited).string(k);
            writer.tag(2, WireType.LengthDelimited).fork();
            HeaderValue.internalBinaryWrite(message.headers[k], writer, options);
            writer.join().join();
        }
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.common.grpcws.Failure
 */
export const Failure = new Failure$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Cancel$Type extends MessageType<Cancel> {
    constructor() {
        super("resources.common.grpcws.Cancel", []);
    }
    create(value?: PartialMessage<Cancel>): Cancel {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<Cancel>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Cancel): Cancel {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
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
    internalBinaryWrite(message: Cancel, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.common.grpcws.Cancel
 */
export const Cancel = new Cancel$Type();
