// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "resources/accounts/oauth2.proto" (package "resources.accounts", syntax proto3)
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
 * @generated from protobuf message resources.accounts.OAuth2Account
 */
export interface OAuth2Account {
    /**
     * @generated from protobuf field: uint64 account_id = 1
     */
    accountId: number;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: string provider_name = 3
     */
    providerName: string;
    /**
     * @generated from protobuf field: resources.accounts.OAuth2Provider provider = 4
     */
    provider?: OAuth2Provider;
    /**
     * @generated from protobuf field: string external_id = 5
     */
    externalId: string;
    /**
     * @generated from protobuf field: string username = 6
     */
    username: string;
    /**
     * @generated from protobuf field: string avatar = 7
     */
    avatar: string;
}
/**
 * @generated from protobuf message resources.accounts.OAuth2Provider
 */
export interface OAuth2Provider {
    /**
     * @generated from protobuf field: string name = 1
     */
    name: string;
    /**
     * @generated from protobuf field: string label = 2
     */
    label: string;
    /**
     * @generated from protobuf field: string homepage = 3
     */
    homepage: string;
    /**
     * @generated from protobuf field: optional string icon = 4
     */
    icon?: string;
}
// @generated message type with reflection information, may provide speed optimized methods
class OAuth2Account$Type extends MessageType<OAuth2Account> {
    constructor() {
        super("resources.accounts.OAuth2Account", [
            { no: 1, name: "account_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "provider_name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "255" } }, "tagger.tags": "sql:\"primary_key\" alias:\"provider_name\"" } },
            { no: 4, name: "provider", kind: "message", T: () => OAuth2Provider },
            { no: 5, name: "external_id", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "128" } } } },
            { no: 6, name: "username", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "255" } } } },
            { no: 7, name: "avatar", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "255" } } } }
        ]);
    }
    create(value?: PartialMessage<OAuth2Account>): OAuth2Account {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.accountId = 0;
        message.providerName = "";
        message.externalId = "";
        message.username = "";
        message.avatar = "";
        if (value !== undefined)
            reflectionMergePartial<OAuth2Account>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: OAuth2Account): OAuth2Account {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 account_id */ 1:
                    message.accountId = reader.uint64().toNumber();
                    break;
                case /* optional resources.timestamp.Timestamp created_at */ 2:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* string provider_name */ 3:
                    message.providerName = reader.string();
                    break;
                case /* resources.accounts.OAuth2Provider provider */ 4:
                    message.provider = OAuth2Provider.internalBinaryRead(reader, reader.uint32(), options, message.provider);
                    break;
                case /* string external_id */ 5:
                    message.externalId = reader.string();
                    break;
                case /* string username */ 6:
                    message.username = reader.string();
                    break;
                case /* string avatar */ 7:
                    message.avatar = reader.string();
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
    internalBinaryWrite(message: OAuth2Account, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 account_id = 1; */
        if (message.accountId !== 0)
            writer.tag(1, WireType.Varint).uint64(message.accountId);
        /* optional resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* string provider_name = 3; */
        if (message.providerName !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.providerName);
        /* resources.accounts.OAuth2Provider provider = 4; */
        if (message.provider)
            OAuth2Provider.internalBinaryWrite(message.provider, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* string external_id = 5; */
        if (message.externalId !== "")
            writer.tag(5, WireType.LengthDelimited).string(message.externalId);
        /* string username = 6; */
        if (message.username !== "")
            writer.tag(6, WireType.LengthDelimited).string(message.username);
        /* string avatar = 7; */
        if (message.avatar !== "")
            writer.tag(7, WireType.LengthDelimited).string(message.avatar);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.accounts.OAuth2Account
 */
export const OAuth2Account = new OAuth2Account$Type();
// @generated message type with reflection information, may provide speed optimized methods
class OAuth2Provider$Type extends MessageType<OAuth2Provider> {
    constructor() {
        super("resources.accounts.OAuth2Provider", [
            { no: 1, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "label", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "homepage", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 4, name: "icon", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<OAuth2Provider>): OAuth2Provider {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.name = "";
        message.label = "";
        message.homepage = "";
        if (value !== undefined)
            reflectionMergePartial<OAuth2Provider>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: OAuth2Provider): OAuth2Provider {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string name */ 1:
                    message.name = reader.string();
                    break;
                case /* string label */ 2:
                    message.label = reader.string();
                    break;
                case /* string homepage */ 3:
                    message.homepage = reader.string();
                    break;
                case /* optional string icon */ 4:
                    message.icon = reader.string();
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
    internalBinaryWrite(message: OAuth2Provider, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string name = 1; */
        if (message.name !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.name);
        /* string label = 2; */
        if (message.label !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.label);
        /* string homepage = 3; */
        if (message.homepage !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.homepage);
        /* optional string icon = 4; */
        if (message.icon !== undefined)
            writer.tag(4, WireType.LengthDelimited).string(message.icon);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.accounts.OAuth2Provider
 */
export const OAuth2Provider = new OAuth2Provider$Type();
