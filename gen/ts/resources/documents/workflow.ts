// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "resources/documents/workflow.proto" (package "resources.documents", syntax proto3)
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
import { Duration } from "../../google/protobuf/duration";
/**
 * @dbscanner: json
 *
 * @generated from protobuf message resources.documents.Workflow
 */
export interface Workflow {
    /**
     * @generated from protobuf field: bool reminder = 1
     */
    reminder: boolean;
    /**
     * @generated from protobuf field: resources.documents.ReminderSettings reminder_settings = 2
     */
    reminderSettings?: ReminderSettings;
    /**
     * @generated from protobuf field: bool auto_close = 3
     */
    autoClose: boolean;
    /**
     * @generated from protobuf field: resources.documents.AutoCloseSettings auto_close_settings = 4
     */
    autoCloseSettings?: AutoCloseSettings;
}
/**
 * @generated from protobuf message resources.documents.ReminderSettings
 */
export interface ReminderSettings {
    /**
     * @generated from protobuf field: repeated resources.documents.Reminder reminders = 1
     */
    reminders: Reminder[];
}
/**
 * @generated from protobuf message resources.documents.Reminder
 */
export interface Reminder {
    /**
     * @generated from protobuf field: google.protobuf.Duration duration = 1
     */
    duration?: Duration;
    /**
     * @generated from protobuf field: string message = 2
     */
    message: string;
}
/**
 * @generated from protobuf message resources.documents.AutoCloseSettings
 */
export interface AutoCloseSettings {
    /**
     * @generated from protobuf field: google.protobuf.Duration duration = 1
     */
    duration?: Duration;
    /**
     * @generated from protobuf field: string message = 2
     */
    message: string;
}
/**
 * @generated from protobuf message resources.documents.WorkflowCronData
 */
export interface WorkflowCronData {
    /**
     * @generated from protobuf field: uint64 last_doc_id = 1
     */
    lastDocId: number;
}
// @generated message type with reflection information, may provide speed optimized methods
class Workflow$Type extends MessageType<Workflow> {
    constructor() {
        super("resources.documents.Workflow", [
            { no: 1, name: "reminder", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 2, name: "reminder_settings", kind: "message", T: () => ReminderSettings },
            { no: 3, name: "auto_close", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 4, name: "auto_close_settings", kind: "message", T: () => AutoCloseSettings }
        ]);
    }
    create(value?: PartialMessage<Workflow>): Workflow {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.reminder = false;
        message.autoClose = false;
        if (value !== undefined)
            reflectionMergePartial<Workflow>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Workflow): Workflow {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* bool reminder */ 1:
                    message.reminder = reader.bool();
                    break;
                case /* resources.documents.ReminderSettings reminder_settings */ 2:
                    message.reminderSettings = ReminderSettings.internalBinaryRead(reader, reader.uint32(), options, message.reminderSettings);
                    break;
                case /* bool auto_close */ 3:
                    message.autoClose = reader.bool();
                    break;
                case /* resources.documents.AutoCloseSettings auto_close_settings */ 4:
                    message.autoCloseSettings = AutoCloseSettings.internalBinaryRead(reader, reader.uint32(), options, message.autoCloseSettings);
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
    internalBinaryWrite(message: Workflow, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* bool reminder = 1; */
        if (message.reminder !== false)
            writer.tag(1, WireType.Varint).bool(message.reminder);
        /* resources.documents.ReminderSettings reminder_settings = 2; */
        if (message.reminderSettings)
            ReminderSettings.internalBinaryWrite(message.reminderSettings, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* bool auto_close = 3; */
        if (message.autoClose !== false)
            writer.tag(3, WireType.Varint).bool(message.autoClose);
        /* resources.documents.AutoCloseSettings auto_close_settings = 4; */
        if (message.autoCloseSettings)
            AutoCloseSettings.internalBinaryWrite(message.autoCloseSettings, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.documents.Workflow
 */
export const Workflow = new Workflow$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ReminderSettings$Type extends MessageType<ReminderSettings> {
    constructor() {
        super("resources.documents.ReminderSettings", [
            { no: 1, name: "reminders", kind: "message", repeat: 2 /*RepeatType.UNPACKED*/, T: () => Reminder, options: { "buf.validate.field": { repeated: { maxItems: "3" } } } }
        ]);
    }
    create(value?: PartialMessage<ReminderSettings>): ReminderSettings {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.reminders = [];
        if (value !== undefined)
            reflectionMergePartial<ReminderSettings>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: ReminderSettings): ReminderSettings {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.documents.Reminder reminders */ 1:
                    message.reminders.push(Reminder.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: ReminderSettings, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.documents.Reminder reminders = 1; */
        for (let i = 0; i < message.reminders.length; i++)
            Reminder.internalBinaryWrite(message.reminders[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.documents.ReminderSettings
 */
export const ReminderSettings = new ReminderSettings$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Reminder$Type extends MessageType<Reminder> {
    constructor() {
        super("resources.documents.Reminder", [
            { no: 1, name: "duration", kind: "message", T: () => Duration, options: { "buf.validate.field": { required: true, duration: { lt: { seconds: "7776000" }, gte: { seconds: "86400" } } } } },
            { no: 2, name: "message", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "255" } } } }
        ]);
    }
    create(value?: PartialMessage<Reminder>): Reminder {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.message = "";
        if (value !== undefined)
            reflectionMergePartial<Reminder>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Reminder): Reminder {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* google.protobuf.Duration duration */ 1:
                    message.duration = Duration.internalBinaryRead(reader, reader.uint32(), options, message.duration);
                    break;
                case /* string message */ 2:
                    message.message = reader.string();
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
    internalBinaryWrite(message: Reminder, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* google.protobuf.Duration duration = 1; */
        if (message.duration)
            Duration.internalBinaryWrite(message.duration, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* string message = 2; */
        if (message.message !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.message);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.documents.Reminder
 */
export const Reminder = new Reminder$Type();
// @generated message type with reflection information, may provide speed optimized methods
class AutoCloseSettings$Type extends MessageType<AutoCloseSettings> {
    constructor() {
        super("resources.documents.AutoCloseSettings", [
            { no: 1, name: "duration", kind: "message", T: () => Duration, options: { "buf.validate.field": { required: true, duration: { lt: { seconds: "7776000" }, gte: { seconds: "86400" } } } } },
            { no: 2, name: "message", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "buf.validate.field": { string: { maxLen: "255" } } } }
        ]);
    }
    create(value?: PartialMessage<AutoCloseSettings>): AutoCloseSettings {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.message = "";
        if (value !== undefined)
            reflectionMergePartial<AutoCloseSettings>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: AutoCloseSettings): AutoCloseSettings {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* google.protobuf.Duration duration */ 1:
                    message.duration = Duration.internalBinaryRead(reader, reader.uint32(), options, message.duration);
                    break;
                case /* string message */ 2:
                    message.message = reader.string();
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
    internalBinaryWrite(message: AutoCloseSettings, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* google.protobuf.Duration duration = 1; */
        if (message.duration)
            Duration.internalBinaryWrite(message.duration, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* string message = 2; */
        if (message.message !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.message);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.documents.AutoCloseSettings
 */
export const AutoCloseSettings = new AutoCloseSettings$Type();
// @generated message type with reflection information, may provide speed optimized methods
class WorkflowCronData$Type extends MessageType<WorkflowCronData> {
    constructor() {
        super("resources.documents.WorkflowCronData", [
            { no: 1, name: "last_doc_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ }
        ]);
    }
    create(value?: PartialMessage<WorkflowCronData>): WorkflowCronData {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.lastDocId = 0;
        if (value !== undefined)
            reflectionMergePartial<WorkflowCronData>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: WorkflowCronData): WorkflowCronData {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 last_doc_id */ 1:
                    message.lastDocId = reader.uint64().toNumber();
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
    internalBinaryWrite(message: WorkflowCronData, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 last_doc_id = 1; */
        if (message.lastDocId !== 0)
            writer.tag(1, WireType.Varint).uint64(message.lastDocId);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.documents.WorkflowCronData
 */
export const WorkflowCronData = new WorkflowCronData$Type();
