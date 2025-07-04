// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "resources/mailer/template.proto" (package "resources.mailer", syntax proto3)
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
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.mailer.Template
 */
export interface Template {
    /**
     * @generated from protobuf field: uint64 id = 1
     */
    id: number;
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp created_at = 3
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 4
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp deleted_at = 5
     */
    deletedAt?: Timestamp;
    /**
     * @generated from protobuf field: uint64 email_id = 6
     */
    emailId: number;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: string title = 7
     */
    title: string;
    /**
     * @sanitize
     *
     * @generated from protobuf field: string content = 8
     */
    content: string;
    /**
     * @generated from protobuf field: optional string creator_job = 9
     */
    creatorJob?: string;
    /**
     * @generated from protobuf field: optional int32 creator_id = 10
     */
    creatorId?: number;
}
// @generated message type with reflection information, may provide speed optimized methods
class Template$Type extends MessageType<Template> {
    constructor() {
        super("resources.mailer.Template", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 3, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 4, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 5, name: "deleted_at", kind: "message", T: () => Timestamp },
            { no: 6, name: "email_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 7, name: "title", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "3", maxLen: "255" } } } },
            { no: 8, name: "content", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { minLen: "3", maxLen: "10240" } } } },
            { no: 9, name: "creator_job", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "40" } } } },
            { no: 10, name: "creator_id", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/, options: { "buf.validate.field": { int32: { gt: 0 } } } }
        ]);
    }
    create(value?: PartialMessage<Template>): Template {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.emailId = 0;
        message.title = "";
        message.content = "";
        if (value !== undefined)
            reflectionMergePartial<Template>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Template): Template {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* resources.timestamp.Timestamp created_at */ 3:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* optional resources.timestamp.Timestamp updated_at */ 4:
                    message.updatedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.updatedAt);
                    break;
                case /* optional resources.timestamp.Timestamp deleted_at */ 5:
                    message.deletedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.deletedAt);
                    break;
                case /* uint64 email_id */ 6:
                    message.emailId = reader.uint64().toNumber();
                    break;
                case /* string title */ 7:
                    message.title = reader.string();
                    break;
                case /* string content */ 8:
                    message.content = reader.string();
                    break;
                case /* optional string creator_job */ 9:
                    message.creatorJob = reader.string();
                    break;
                case /* optional int32 creator_id */ 10:
                    message.creatorId = reader.int32();
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
    internalBinaryWrite(message: Template, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* resources.timestamp.Timestamp created_at = 3; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp updated_at = 4; */
        if (message.updatedAt)
            Timestamp.internalBinaryWrite(message.updatedAt, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.timestamp.Timestamp deleted_at = 5; */
        if (message.deletedAt)
            Timestamp.internalBinaryWrite(message.deletedAt, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* uint64 email_id = 6; */
        if (message.emailId !== 0)
            writer.tag(6, WireType.Varint).uint64(message.emailId);
        /* string title = 7; */
        if (message.title !== "")
            writer.tag(7, WireType.LengthDelimited).string(message.title);
        /* string content = 8; */
        if (message.content !== "")
            writer.tag(8, WireType.LengthDelimited).string(message.content);
        /* optional string creator_job = 9; */
        if (message.creatorJob !== undefined)
            writer.tag(9, WireType.LengthDelimited).string(message.creatorJob);
        /* optional int32 creator_id = 10; */
        if (message.creatorId !== undefined)
            writer.tag(10, WireType.Varint).int32(message.creatorId);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.mailer.Template
 */
export const Template = new Template$Type();
