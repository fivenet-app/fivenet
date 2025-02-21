// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/mailer/message.proto" (package "resources.mailer", syntax proto3)
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
import { Content } from "../common/content/content";
import { Timestamp } from "../timestamp/timestamp";
import { Email } from "./email";
/**
 * @generated from protobuf message resources.mailer.Message
 */
export interface Message {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: number;
    /**
     * @generated from protobuf field: uint64 thread_id = 2;
     */
    threadId: number;
    /**
     * @generated from protobuf field: uint64 sender_id = 3;
     */
    senderId: number;
    /**
     * @generated from protobuf field: optional resources.mailer.Email sender = 4;
     */
    sender?: Email; // @gotags: alias:"sender"
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp created_at = 5;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 6;
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp deleted_at = 7;
     */
    deletedAt?: Timestamp;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: string title = 8;
     */
    title: string;
    /**
     * @sanitize
     *
     * @generated from protobuf field: resources.common.content.Content content = 9;
     */
    content?: Content;
    /**
     * @generated from protobuf field: optional resources.mailer.MessageData data = 10;
     */
    data?: MessageData;
    /**
     * @generated from protobuf field: optional int32 creator_id = 11;
     */
    creatorId?: number;
    /**
     * @generated from protobuf field: optional string creator_job = 12;
     */
    creatorJob?: string;
}
/**
 * @generated from protobuf message resources.mailer.MessageData
 */
export interface MessageData {
    /**
     * @generated from protobuf field: repeated resources.mailer.MessageDataEntry entry = 1;
     */
    entry: MessageDataEntry[];
}
/**
 * @generated from protobuf message resources.mailer.MessageDataEntry
 */
export interface MessageDataEntry {
}
// @generated message type with reflection information, may provide speed optimized methods
class Message$Type extends MessageType<Message> {
    constructor() {
        super("resources.mailer.Message", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "thread_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 3, name: "sender_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 4, name: "sender", kind: "message", T: () => Email },
            { no: 5, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 6, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 7, name: "deleted_at", kind: "message", T: () => Timestamp },
            { no: 8, name: "title", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "3", maxLen: "255" } } } },
            { no: 9, name: "content", kind: "message", T: () => Content },
            { no: 10, name: "data", kind: "message", T: () => MessageData },
            { no: 11, name: "creator_id", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/ },
            { no: 12, name: "creator_job", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<Message>): Message {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.threadId = 0;
        message.senderId = 0;
        message.title = "";
        if (value !== undefined)
            reflectionMergePartial<Message>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Message): Message {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* uint64 thread_id */ 2:
                    message.threadId = reader.uint64().toNumber();
                    break;
                case /* uint64 sender_id */ 3:
                    message.senderId = reader.uint64().toNumber();
                    break;
                case /* optional resources.mailer.Email sender */ 4:
                    message.sender = Email.internalBinaryRead(reader, reader.uint32(), options, message.sender);
                    break;
                case /* resources.timestamp.Timestamp created_at */ 5:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* optional resources.timestamp.Timestamp updated_at */ 6:
                    message.updatedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.updatedAt);
                    break;
                case /* optional resources.timestamp.Timestamp deleted_at */ 7:
                    message.deletedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.deletedAt);
                    break;
                case /* string title */ 8:
                    message.title = reader.string();
                    break;
                case /* resources.common.content.Content content */ 9:
                    message.content = Content.internalBinaryRead(reader, reader.uint32(), options, message.content);
                    break;
                case /* optional resources.mailer.MessageData data */ 10:
                    message.data = MessageData.internalBinaryRead(reader, reader.uint32(), options, message.data);
                    break;
                case /* optional int32 creator_id */ 11:
                    message.creatorId = reader.int32();
                    break;
                case /* optional string creator_job */ 12:
                    message.creatorJob = reader.string();
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
    internalBinaryWrite(message: Message, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* uint64 thread_id = 2; */
        if (message.threadId !== 0)
            writer.tag(2, WireType.Varint).uint64(message.threadId);
        /* uint64 sender_id = 3; */
        if (message.senderId !== 0)
            writer.tag(3, WireType.Varint).uint64(message.senderId);
        /* optional resources.mailer.Email sender = 4; */
        if (message.sender)
            Email.internalBinaryWrite(message.sender, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* resources.timestamp.Timestamp created_at = 5; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp updated_at = 6; */
        if (message.updatedAt)
            Timestamp.internalBinaryWrite(message.updatedAt, writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp deleted_at = 7; */
        if (message.deletedAt)
            Timestamp.internalBinaryWrite(message.deletedAt, writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        /* string title = 8; */
        if (message.title !== "")
            writer.tag(8, WireType.LengthDelimited).string(message.title);
        /* resources.common.content.Content content = 9; */
        if (message.content)
            Content.internalBinaryWrite(message.content, writer.tag(9, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.mailer.MessageData data = 10; */
        if (message.data)
            MessageData.internalBinaryWrite(message.data, writer.tag(10, WireType.LengthDelimited).fork(), options).join();
        /* optional int32 creator_id = 11; */
        if (message.creatorId !== undefined)
            writer.tag(11, WireType.Varint).int32(message.creatorId);
        /* optional string creator_job = 12; */
        if (message.creatorJob !== undefined)
            writer.tag(12, WireType.LengthDelimited).string(message.creatorJob);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.mailer.Message
 */
export const Message = new Message$Type();
// @generated message type with reflection information, may provide speed optimized methods
class MessageData$Type extends MessageType<MessageData> {
    constructor() {
        super("resources.mailer.MessageData", [
            { no: 1, name: "entry", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => MessageDataEntry, options: { "validate.rules": { repeated: { maxItems: "3" } } } }
        ]);
    }
    create(value?: PartialMessage<MessageData>): MessageData {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.entry = [];
        if (value !== undefined)
            reflectionMergePartial<MessageData>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: MessageData): MessageData {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.mailer.MessageDataEntry entry */ 1:
                    message.entry.push(MessageDataEntry.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: MessageData, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.mailer.MessageDataEntry entry = 1; */
        for (let i = 0; i < message.entry.length; i++)
            MessageDataEntry.internalBinaryWrite(message.entry[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.mailer.MessageData
 */
export const MessageData = new MessageData$Type();
// @generated message type with reflection information, may provide speed optimized methods
class MessageDataEntry$Type extends MessageType<MessageDataEntry> {
    constructor() {
        super("resources.mailer.MessageDataEntry", []);
    }
    create(value?: PartialMessage<MessageDataEntry>): MessageDataEntry {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<MessageDataEntry>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: MessageDataEntry): MessageDataEntry {
        return target ?? this.create();
    }
    internalBinaryWrite(message: MessageDataEntry, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.mailer.MessageDataEntry
 */
export const MessageDataEntry = new MessageDataEntry$Type();
