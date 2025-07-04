// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "resources/mailer/message.proto" (package "resources.mailer", syntax proto3)
// tslint:disable
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
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
    /**
     * @generated from protobuf field: uint64 thread_id = 2
     */
    threadId: number;
    /**
     * @generated from protobuf field: uint64 sender_id = 3
     */
    senderId: number;
    /**
     * @generated from protobuf field: optional resources.mailer.Email sender = 4
     */
    sender?: Email;
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp created_at = 5
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 6
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp deleted_at = 7
     */
    deletedAt?: Timestamp;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: string title = 8
     */
    title: string;
    /**
     * @sanitize
     *
     * @generated from protobuf field: resources.common.content.Content content = 9
     */
    content?: Content;
    /**
     * @generated from protobuf field: optional resources.mailer.MessageData data = 10
     */
    data?: MessageData;
    /**
     * @generated from protobuf field: optional int32 creator_id = 11
     */
    creatorId?: number;
    /**
     * @generated from protobuf field: optional string creator_job = 12
     */
    creatorJob?: string;
}
/**
 * @dbscanner: json
 *
 * @generated from protobuf message resources.mailer.MessageData
 */
export interface MessageData {
    /**
     * @generated from protobuf field: repeated resources.mailer.MessageAttachment attachments = 1
     */
    attachments: MessageAttachment[];
}
/**
 * @generated from protobuf message resources.mailer.MessageAttachment
 */
export interface MessageAttachment {
    /**
     * @generated from protobuf oneof: data
     */
    data: {
        oneofKind: "document";
        /**
         * @generated from protobuf field: resources.mailer.MessageAttachmentDocument document = 1
         */
        document: MessageAttachmentDocument;
    } | {
        oneofKind: undefined;
    };
}
/**
 * @generated from protobuf message resources.mailer.MessageAttachmentDocument
 */
export interface MessageAttachmentDocument {
    /**
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
    /**
     * @generated from protobuf field: optional string title = 2
     */
    title?: string;
}
// @generated message type with reflection information, may provide speed optimized methods
class Message$Type extends MessageType<Message> {
    constructor() {
        super("resources.mailer.Message", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "thread_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 3, name: "sender_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 4, name: "sender", kind: "message", T: () => Email, options: { "tagger.tags": "alias:\"sender\"" } },
            { no: 5, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 6, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 7, name: "deleted_at", kind: "message", T: () => Timestamp },
            { no: 8, name: "title", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "3", maxLen: "255" } } } },
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
            { no: 1, name: "attachments", kind: "message", repeat: 2 /*RepeatType.UNPACKED*/, T: () => MessageAttachment, options: { "buf.validate.field": { repeated: { maxItems: "3" } } } }
        ]);
    }
    create(value?: PartialMessage<MessageData>): MessageData {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.attachments = [];
        if (value !== undefined)
            reflectionMergePartial<MessageData>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: MessageData): MessageData {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.mailer.MessageAttachment attachments */ 1:
                    message.attachments.push(MessageAttachment.internalBinaryRead(reader, reader.uint32(), options));
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
        /* repeated resources.mailer.MessageAttachment attachments = 1; */
        for (let i = 0; i < message.attachments.length; i++)
            MessageAttachment.internalBinaryWrite(message.attachments[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
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
class MessageAttachment$Type extends MessageType<MessageAttachment> {
    constructor() {
        super("resources.mailer.MessageAttachment", [
            { no: 1, name: "document", kind: "message", oneof: "data", T: () => MessageAttachmentDocument }
        ]);
    }
    create(value?: PartialMessage<MessageAttachment>): MessageAttachment {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.data = { oneofKind: undefined };
        if (value !== undefined)
            reflectionMergePartial<MessageAttachment>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: MessageAttachment): MessageAttachment {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.mailer.MessageAttachmentDocument document */ 1:
                    message.data = {
                        oneofKind: "document",
                        document: MessageAttachmentDocument.internalBinaryRead(reader, reader.uint32(), options, (message.data as any).document)
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
    internalBinaryWrite(message: MessageAttachment, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.mailer.MessageAttachmentDocument document = 1; */
        if (message.data.oneofKind === "document")
            MessageAttachmentDocument.internalBinaryWrite(message.data.document, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.mailer.MessageAttachment
 */
export const MessageAttachment = new MessageAttachment$Type();
// @generated message type with reflection information, may provide speed optimized methods
class MessageAttachmentDocument$Type extends MessageType<MessageAttachmentDocument> {
    constructor() {
        super("resources.mailer.MessageAttachmentDocument", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "title", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "768" } } } }
        ]);
    }
    create(value?: PartialMessage<MessageAttachmentDocument>): MessageAttachmentDocument {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        if (value !== undefined)
            reflectionMergePartial<MessageAttachmentDocument>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: MessageAttachmentDocument): MessageAttachmentDocument {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* optional string title */ 2:
                    message.title = reader.string();
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
    internalBinaryWrite(message: MessageAttachmentDocument, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* optional string title = 2; */
        if (message.title !== undefined)
            writer.tag(2, WireType.LengthDelimited).string(message.title);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.mailer.MessageAttachmentDocument
 */
export const MessageAttachmentDocument = new MessageAttachmentDocument$Type();
