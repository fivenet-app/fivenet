// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/internet/internet.proto" (package "services.internet", syntax proto3)
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
import { SearchResult } from "../../resources/internet/search";
/**
 * @generated from protobuf message services.internet.SearchRequest
 */
export interface SearchRequest {
    /**
     * @generated from protobuf field: string search = 1;
     */
    search: string;
}
/**
 * @generated from protobuf message services.internet.SearchResponse
 */
export interface SearchResponse {
    /**
     * @generated from protobuf field: repeated resources.internet.SearchResult results = 1;
     */
    results: SearchResult[];
}
// @generated message type with reflection information, may provide speed optimized methods
class SearchRequest$Type extends MessageType<SearchRequest> {
    constructor() {
        super("services.internet.SearchRequest", [
            { no: 1, name: "search", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<SearchRequest>): SearchRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.search = "";
        if (value !== undefined)
            reflectionMergePartial<SearchRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SearchRequest): SearchRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string search */ 1:
                    message.search = reader.string();
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
    internalBinaryWrite(message: SearchRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string search = 1; */
        if (message.search !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.search);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.internet.SearchRequest
 */
export const SearchRequest = new SearchRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class SearchResponse$Type extends MessageType<SearchResponse> {
    constructor() {
        super("services.internet.SearchResponse", [
            { no: 1, name: "results", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => SearchResult }
        ]);
    }
    create(value?: PartialMessage<SearchResponse>): SearchResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.results = [];
        if (value !== undefined)
            reflectionMergePartial<SearchResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: SearchResponse): SearchResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.internet.SearchResult results */ 1:
                    message.results.push(SearchResult.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: SearchResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.internet.SearchResult results = 1; */
        for (let i = 0; i < message.results.length; i++)
            SearchResult.internalBinaryWrite(message.results[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.internet.SearchResponse
 */
export const SearchResponse = new SearchResponse$Type();
/**
 * @generated ServiceType for protobuf service services.internet.InternetService
 */
export const InternetService = new ServiceType("services.internet.InternetService", [
    { name: "Search", options: {}, I: SearchRequest, O: SearchResponse }
]);