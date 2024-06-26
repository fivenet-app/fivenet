// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/messenger/thread.proto" (package "resources.messenger", syntax proto3)
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
import { ThreadAccess } from "./access";
import { UserShort } from "../users/users";
import { Message } from "./message";
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.messenger.Thread
 */
export interface Thread {
    /**
     * @generated from protobuf field: uint64 id = 1 [jstype = JS_STRING];
     */
    id: string;
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 3;
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp deleted_at = 4;
     */
    deletedAt?: Timestamp;
    /**
     * @sanitize
     *
     * @generated from protobuf field: string title = 5;
     */
    title: string;
    /**
     * @generated from protobuf field: bool archived = 6;
     */
    archived: boolean;
    /**
     * @generated from protobuf field: optional resources.messenger.Message last_message = 7;
     */
    lastMessage?: Message;
    /**
     * @generated from protobuf field: resources.messenger.ThreadUserState user_state = 8;
     */
    userState?: ThreadUserState;
    /**
     * @generated from protobuf field: string creator_job = 9;
     */
    creatorJob: string;
    /**
     * @generated from protobuf field: optional int32 creator_id = 10;
     */
    creatorId?: number;
    /**
     * @generated from protobuf field: optional resources.users.UserShort creator = 11;
     */
    creator?: UserShort; // @gotags: alias:"creator"
    /**
     * @generated from protobuf field: resources.messenger.ThreadAccess access = 12;
     */
    access?: ThreadAccess;
}
/**
 * @generated from protobuf message resources.messenger.ThreadUserState
 */
export interface ThreadUserState {
    /**
     * @generated from protobuf field: uint64 thread_id = 1 [jstype = JS_STRING];
     */
    threadId: string;
    /**
     * @generated from protobuf field: int32 user_id = 2;
     */
    userId: number;
    /**
     * @generated from protobuf field: bool unread = 3;
     */
    unread: boolean;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp last_read = 4;
     */
    lastRead?: Timestamp;
    /**
     * @generated from protobuf field: bool important = 5;
     */
    important: boolean;
    /**
     * @generated from protobuf field: bool favorite = 6;
     */
    favorite: boolean;
    /**
     * @generated from protobuf field: bool muted = 7;
     */
    muted: boolean;
}
// @generated message type with reflection information, may provide speed optimized methods
class Thread$Type extends MessageType<Thread> {
    constructor() {
        super("resources.messenger.Thread", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "deleted_at", kind: "message", T: () => Timestamp },
            { no: 5, name: "title", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "3", maxLen: "255" } } } },
            { no: 6, name: "archived", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 7, name: "last_message", kind: "message", T: () => Message },
            { no: 8, name: "user_state", kind: "message", T: () => ThreadUserState },
            { no: 9, name: "creator_job", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 10, name: "creator_id", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/ },
            { no: 11, name: "creator", kind: "message", T: () => UserShort },
            { no: 12, name: "access", kind: "message", T: () => ThreadAccess }
        ]);
    }
    create(value?: PartialMessage<Thread>): Thread {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = "0";
        message.title = "";
        message.archived = false;
        message.creatorJob = "";
        if (value !== undefined)
            reflectionMergePartial<Thread>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Thread): Thread {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id = 1 [jstype = JS_STRING];*/ 1:
                    message.id = reader.uint64().toString();
                    break;
                case /* resources.timestamp.Timestamp created_at */ 2:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* optional resources.timestamp.Timestamp updated_at */ 3:
                    message.updatedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.updatedAt);
                    break;
                case /* optional resources.timestamp.Timestamp deleted_at */ 4:
                    message.deletedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.deletedAt);
                    break;
                case /* string title */ 5:
                    message.title = reader.string();
                    break;
                case /* bool archived */ 6:
                    message.archived = reader.bool();
                    break;
                case /* optional resources.messenger.Message last_message */ 7:
                    message.lastMessage = Message.internalBinaryRead(reader, reader.uint32(), options, message.lastMessage);
                    break;
                case /* resources.messenger.ThreadUserState user_state */ 8:
                    message.userState = ThreadUserState.internalBinaryRead(reader, reader.uint32(), options, message.userState);
                    break;
                case /* string creator_job */ 9:
                    message.creatorJob = reader.string();
                    break;
                case /* optional int32 creator_id */ 10:
                    message.creatorId = reader.int32();
                    break;
                case /* optional resources.users.UserShort creator */ 11:
                    message.creator = UserShort.internalBinaryRead(reader, reader.uint32(), options, message.creator);
                    break;
                case /* resources.messenger.ThreadAccess access */ 12:
                    message.access = ThreadAccess.internalBinaryRead(reader, reader.uint32(), options, message.access);
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
    internalBinaryWrite(message: Thread, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1 [jstype = JS_STRING]; */
        if (message.id !== "0")
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp updated_at = 3; */
        if (message.updatedAt)
            Timestamp.internalBinaryWrite(message.updatedAt, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp deleted_at = 4; */
        if (message.deletedAt)
            Timestamp.internalBinaryWrite(message.deletedAt, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* string title = 5; */
        if (message.title !== "")
            writer.tag(5, WireType.LengthDelimited).string(message.title);
        /* bool archived = 6; */
        if (message.archived !== false)
            writer.tag(6, WireType.Varint).bool(message.archived);
        /* optional resources.messenger.Message last_message = 7; */
        if (message.lastMessage)
            Message.internalBinaryWrite(message.lastMessage, writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        /* resources.messenger.ThreadUserState user_state = 8; */
        if (message.userState)
            ThreadUserState.internalBinaryWrite(message.userState, writer.tag(8, WireType.LengthDelimited).fork(), options).join();
        /* string creator_job = 9; */
        if (message.creatorJob !== "")
            writer.tag(9, WireType.LengthDelimited).string(message.creatorJob);
        /* optional int32 creator_id = 10; */
        if (message.creatorId !== undefined)
            writer.tag(10, WireType.Varint).int32(message.creatorId);
        /* optional resources.users.UserShort creator = 11; */
        if (message.creator)
            UserShort.internalBinaryWrite(message.creator, writer.tag(11, WireType.LengthDelimited).fork(), options).join();
        /* resources.messenger.ThreadAccess access = 12; */
        if (message.access)
            ThreadAccess.internalBinaryWrite(message.access, writer.tag(12, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.messenger.Thread
 */
export const Thread = new Thread$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ThreadUserState$Type extends MessageType<ThreadUserState> {
    constructor() {
        super("resources.messenger.ThreadUserState", [
            { no: 1, name: "thread_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/ },
            { no: 2, name: "user_id", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 3, name: "unread", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 4, name: "last_read", kind: "message", T: () => Timestamp },
            { no: 5, name: "important", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 6, name: "favorite", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 7, name: "muted", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<ThreadUserState>): ThreadUserState {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.threadId = "0";
        message.userId = 0;
        message.unread = false;
        message.important = false;
        message.favorite = false;
        message.muted = false;
        if (value !== undefined)
            reflectionMergePartial<ThreadUserState>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ThreadUserState): ThreadUserState {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 thread_id = 1 [jstype = JS_STRING];*/ 1:
                    message.threadId = reader.uint64().toString();
                    break;
                case /* int32 user_id */ 2:
                    message.userId = reader.int32();
                    break;
                case /* bool unread */ 3:
                    message.unread = reader.bool();
                    break;
                case /* optional resources.timestamp.Timestamp last_read */ 4:
                    message.lastRead = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.lastRead);
                    break;
                case /* bool important */ 5:
                    message.important = reader.bool();
                    break;
                case /* bool favorite */ 6:
                    message.favorite = reader.bool();
                    break;
                case /* bool muted */ 7:
                    message.muted = reader.bool();
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
    internalBinaryWrite(message: ThreadUserState, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 thread_id = 1 [jstype = JS_STRING]; */
        if (message.threadId !== "0")
            writer.tag(1, WireType.Varint).uint64(message.threadId);
        /* int32 user_id = 2; */
        if (message.userId !== 0)
            writer.tag(2, WireType.Varint).int32(message.userId);
        /* bool unread = 3; */
        if (message.unread !== false)
            writer.tag(3, WireType.Varint).bool(message.unread);
        /* optional resources.timestamp.Timestamp last_read = 4; */
        if (message.lastRead)
            Timestamp.internalBinaryWrite(message.lastRead, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* bool important = 5; */
        if (message.important !== false)
            writer.tag(5, WireType.Varint).bool(message.important);
        /* bool favorite = 6; */
        if (message.favorite !== false)
            writer.tag(6, WireType.Varint).bool(message.favorite);
        /* bool muted = 7; */
        if (message.muted !== false)
            writer.tag(7, WireType.Varint).bool(message.muted);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.messenger.ThreadUserState
 */
export const ThreadUserState = new ThreadUserState$Type();
