// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/notificator/notificator.proto" (package "services.notificator", syntax proto3)
// @ts-nocheck
import { ServiceType } from "@protobuf-ts/runtime-rpc";
import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import { UnknownFieldHandler } from "@protobuf-ts/runtime";
import { WireType } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { reflectionMergePartial } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
import { MailerEvent } from "../../resources/mailer/events";
import { SystemEvent } from "../../resources/notifications/events";
import { JobGradeEvent } from "../../resources/notifications/events";
import { JobEvent } from "../../resources/notifications/events";
import { UserEvent } from "../../resources/notifications/events";
import { Notification } from "../../resources/notifications/notifications";
import { PaginationResponse } from "../../resources/common/database/database";
import { NotificationCategory } from "../../resources/notifications/notifications";
import { PaginationRequest } from "../../resources/common/database/database";
/**
 * @generated from protobuf message services.notificator.GetNotificationsRequest
 */
export interface GetNotificationsRequest {
    /**
     * @generated from protobuf field: resources.common.database.PaginationRequest pagination = 1;
     */
    pagination?: PaginationRequest;
    /**
     * @generated from protobuf field: optional bool include_read = 2;
     */
    includeRead?: boolean;
    /**
     * @generated from protobuf field: repeated resources.notifications.NotificationCategory categories = 3;
     */
    categories: NotificationCategory[];
}
/**
 * @generated from protobuf message services.notificator.GetNotificationsResponse
 */
export interface GetNotificationsResponse {
    /**
     * @generated from protobuf field: resources.common.database.PaginationResponse pagination = 1;
     */
    pagination?: PaginationResponse;
    /**
     * @generated from protobuf field: repeated resources.notifications.Notification notifications = 2;
     */
    notifications: Notification[];
}
/**
 * @generated from protobuf message services.notificator.MarkNotificationsRequest
 */
export interface MarkNotificationsRequest {
    /**
     * @generated from protobuf field: repeated uint64 ids = 1 [jstype = JS_STRING];
     */
    ids: string[];
    /**
     * @generated from protobuf field: optional bool all = 2;
     */
    all?: boolean;
}
/**
 * @generated from protobuf message services.notificator.MarkNotificationsResponse
 */
export interface MarkNotificationsResponse {
    /**
     * @generated from protobuf field: uint64 updated = 1 [jstype = JS_STRING];
     */
    updated: string;
}
/**
 * @generated from protobuf message services.notificator.StreamRequest
 */
export interface StreamRequest {
}
/**
 * @generated from protobuf message services.notificator.StreamResponse
 */
export interface StreamResponse {
    /**
     * @generated from protobuf field: int32 notification_count = 1;
     */
    notificationCount: number;
    /**
     * @generated from protobuf field: optional bool restart = 2;
     */
    restart?: boolean;
    /**
     * @generated from protobuf oneof: data
     */
    data: {
        oneofKind: "userEvent";
        /**
         * @generated from protobuf field: resources.notifications.UserEvent user_event = 3;
         */
        userEvent: UserEvent;
    } | {
        oneofKind: "jobEvent";
        /**
         * @generated from protobuf field: resources.notifications.JobEvent job_event = 4;
         */
        jobEvent: JobEvent;
    } | {
        oneofKind: "jobGradeEvent";
        /**
         * @generated from protobuf field: resources.notifications.JobGradeEvent job_grade_event = 7;
         */
        jobGradeEvent: JobGradeEvent;
    } | {
        oneofKind: "systemEvent";
        /**
         * @generated from protobuf field: resources.notifications.SystemEvent system_event = 5;
         */
        systemEvent: SystemEvent;
    } | {
        oneofKind: "mailerEvent";
        /**
         * @generated from protobuf field: resources.mailer.MailerEvent mailer_event = 6;
         */
        mailerEvent: MailerEvent;
    } | {
        oneofKind: undefined;
    };
}
// @generated message type with reflection information, may provide speed optimized methods
class GetNotificationsRequest$Type extends MessageType<GetNotificationsRequest> {
    constructor() {
        super("services.notificator.GetNotificationsRequest", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationRequest, options: { "validate.rules": { message: { required: true } } } },
            { no: 2, name: "include_read", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ },
            { no: 3, name: "categories", kind: "enum", repeat: 1 /*RepeatType.PACKED*/, T: () => ["resources.notifications.NotificationCategory", NotificationCategory, "NOTIFICATION_CATEGORY_"], options: { "validate.rules": { repeated: { maxItems: "4", items: { enum: { definedOnly: true } } } } } }
        ]);
    }
    create(value?: PartialMessage<GetNotificationsRequest>): GetNotificationsRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.categories = [];
        if (value !== undefined)
            reflectionMergePartial<GetNotificationsRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: GetNotificationsRequest): GetNotificationsRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationRequest pagination */ 1:
                    message.pagination = PaginationRequest.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* optional bool include_read */ 2:
                    message.includeRead = reader.bool();
                    break;
                case /* repeated resources.notifications.NotificationCategory categories */ 3:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.categories.push(reader.int32());
                    else
                        message.categories.push(reader.int32());
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
    internalBinaryWrite(message: GetNotificationsRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationRequest pagination = 1; */
        if (message.pagination)
            PaginationRequest.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* optional bool include_read = 2; */
        if (message.includeRead !== undefined)
            writer.tag(2, WireType.Varint).bool(message.includeRead);
        /* repeated resources.notifications.NotificationCategory categories = 3; */
        if (message.categories.length) {
            writer.tag(3, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.categories.length; i++)
                writer.int32(message.categories[i]);
            writer.join();
        }
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.notificator.GetNotificationsRequest
 */
export const GetNotificationsRequest = new GetNotificationsRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetNotificationsResponse$Type extends MessageType<GetNotificationsResponse> {
    constructor() {
        super("services.notificator.GetNotificationsResponse", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationResponse },
            { no: 2, name: "notifications", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Notification }
        ]);
    }
    create(value?: PartialMessage<GetNotificationsResponse>): GetNotificationsResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.notifications = [];
        if (value !== undefined)
            reflectionMergePartial<GetNotificationsResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: GetNotificationsResponse): GetNotificationsResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.common.database.PaginationResponse pagination */ 1:
                    message.pagination = PaginationResponse.internalBinaryRead(reader, reader.uint32(), options, message.pagination);
                    break;
                case /* repeated resources.notifications.Notification notifications */ 2:
                    message.notifications.push(Notification.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: GetNotificationsResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.common.database.PaginationResponse pagination = 1; */
        if (message.pagination)
            PaginationResponse.internalBinaryWrite(message.pagination, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.notifications.Notification notifications = 2; */
        for (let i = 0; i < message.notifications.length; i++)
            Notification.internalBinaryWrite(message.notifications[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.notificator.GetNotificationsResponse
 */
export const GetNotificationsResponse = new GetNotificationsResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class MarkNotificationsRequest$Type extends MessageType<MarkNotificationsRequest> {
    constructor() {
        super("services.notificator.MarkNotificationsRequest", [
            { no: 1, name: "ids", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 4 /*ScalarType.UINT64*/, options: { "validate.rules": { repeated: { minItems: "1", maxItems: "20", ignoreEmpty: true } } } },
            { no: 2, name: "all", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<MarkNotificationsRequest>): MarkNotificationsRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.ids = [];
        if (value !== undefined)
            reflectionMergePartial<MarkNotificationsRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: MarkNotificationsRequest): MarkNotificationsRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated uint64 ids = 1 [jstype = JS_STRING];*/ 1:
                    if (wireType === WireType.LengthDelimited)
                        for (let e = reader.int32() + reader.pos; reader.pos < e;)
                            message.ids.push(reader.uint64().toString());
                    else
                        message.ids.push(reader.uint64().toString());
                    break;
                case /* optional bool all */ 2:
                    message.all = reader.bool();
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
    internalBinaryWrite(message: MarkNotificationsRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated uint64 ids = 1 [jstype = JS_STRING]; */
        if (message.ids.length) {
            writer.tag(1, WireType.LengthDelimited).fork();
            for (let i = 0; i < message.ids.length; i++)
                writer.uint64(message.ids[i]);
            writer.join();
        }
        /* optional bool all = 2; */
        if (message.all !== undefined)
            writer.tag(2, WireType.Varint).bool(message.all);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.notificator.MarkNotificationsRequest
 */
export const MarkNotificationsRequest = new MarkNotificationsRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class MarkNotificationsResponse$Type extends MessageType<MarkNotificationsResponse> {
    constructor() {
        super("services.notificator.MarkNotificationsResponse", [
            { no: 1, name: "updated", kind: "scalar", T: 4 /*ScalarType.UINT64*/ }
        ]);
    }
    create(value?: PartialMessage<MarkNotificationsResponse>): MarkNotificationsResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.updated = "0";
        if (value !== undefined)
            reflectionMergePartial<MarkNotificationsResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: MarkNotificationsResponse): MarkNotificationsResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 updated = 1 [jstype = JS_STRING];*/ 1:
                    message.updated = reader.uint64().toString();
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
    internalBinaryWrite(message: MarkNotificationsResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 updated = 1 [jstype = JS_STRING]; */
        if (message.updated !== "0")
            writer.tag(1, WireType.Varint).uint64(message.updated);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.notificator.MarkNotificationsResponse
 */
export const MarkNotificationsResponse = new MarkNotificationsResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class StreamRequest$Type extends MessageType<StreamRequest> {
    constructor() {
        super("services.notificator.StreamRequest", []);
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
 * @generated MessageType for protobuf message services.notificator.StreamRequest
 */
export const StreamRequest = new StreamRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class StreamResponse$Type extends MessageType<StreamResponse> {
    constructor() {
        super("services.notificator.StreamResponse", [
            { no: 1, name: "notification_count", kind: "scalar", T: 5 /*ScalarType.INT32*/ },
            { no: 2, name: "restart", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ },
            { no: 3, name: "user_event", kind: "message", oneof: "data", T: () => UserEvent },
            { no: 4, name: "job_event", kind: "message", oneof: "data", T: () => JobEvent },
            { no: 7, name: "job_grade_event", kind: "message", oneof: "data", T: () => JobGradeEvent },
            { no: 5, name: "system_event", kind: "message", oneof: "data", T: () => SystemEvent },
            { no: 6, name: "mailer_event", kind: "message", oneof: "data", T: () => MailerEvent }
        ]);
    }
    create(value?: PartialMessage<StreamResponse>): StreamResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.notificationCount = 0;
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
                case /* int32 notification_count */ 1:
                    message.notificationCount = reader.int32();
                    break;
                case /* optional bool restart */ 2:
                    message.restart = reader.bool();
                    break;
                case /* resources.notifications.UserEvent user_event */ 3:
                    message.data = {
                        oneofKind: "userEvent",
                        userEvent: UserEvent.internalBinaryRead(reader, reader.uint32(), options, (message.data as any).userEvent)
                    };
                    break;
                case /* resources.notifications.JobEvent job_event */ 4:
                    message.data = {
                        oneofKind: "jobEvent",
                        jobEvent: JobEvent.internalBinaryRead(reader, reader.uint32(), options, (message.data as any).jobEvent)
                    };
                    break;
                case /* resources.notifications.JobGradeEvent job_grade_event */ 7:
                    message.data = {
                        oneofKind: "jobGradeEvent",
                        jobGradeEvent: JobGradeEvent.internalBinaryRead(reader, reader.uint32(), options, (message.data as any).jobGradeEvent)
                    };
                    break;
                case /* resources.notifications.SystemEvent system_event */ 5:
                    message.data = {
                        oneofKind: "systemEvent",
                        systemEvent: SystemEvent.internalBinaryRead(reader, reader.uint32(), options, (message.data as any).systemEvent)
                    };
                    break;
                case /* resources.mailer.MailerEvent mailer_event */ 6:
                    message.data = {
                        oneofKind: "mailerEvent",
                        mailerEvent: MailerEvent.internalBinaryRead(reader, reader.uint32(), options, (message.data as any).mailerEvent)
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
        /* int32 notification_count = 1; */
        if (message.notificationCount !== 0)
            writer.tag(1, WireType.Varint).int32(message.notificationCount);
        /* optional bool restart = 2; */
        if (message.restart !== undefined)
            writer.tag(2, WireType.Varint).bool(message.restart);
        /* resources.notifications.UserEvent user_event = 3; */
        if (message.data.oneofKind === "userEvent")
            UserEvent.internalBinaryWrite(message.data.userEvent, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* resources.notifications.JobEvent job_event = 4; */
        if (message.data.oneofKind === "jobEvent")
            JobEvent.internalBinaryWrite(message.data.jobEvent, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* resources.notifications.JobGradeEvent job_grade_event = 7; */
        if (message.data.oneofKind === "jobGradeEvent")
            JobGradeEvent.internalBinaryWrite(message.data.jobGradeEvent, writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        /* resources.notifications.SystemEvent system_event = 5; */
        if (message.data.oneofKind === "systemEvent")
            SystemEvent.internalBinaryWrite(message.data.systemEvent, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* resources.mailer.MailerEvent mailer_event = 6; */
        if (message.data.oneofKind === "mailerEvent")
            MailerEvent.internalBinaryWrite(message.data.mailerEvent, writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message services.notificator.StreamResponse
 */
export const StreamResponse = new StreamResponse$Type();
/**
 * @generated ServiceType for protobuf service services.notificator.NotificatorService
 */
export const NotificatorService = new ServiceType("services.notificator.NotificatorService", [
    { name: "GetNotifications", options: {}, I: GetNotificationsRequest, O: GetNotificationsResponse },
    { name: "MarkNotifications", options: {}, I: MarkNotificationsRequest, O: MarkNotificationsResponse },
    { name: "Stream", serverStreaming: true, options: {}, I: StreamRequest, O: StreamResponse }
]);
