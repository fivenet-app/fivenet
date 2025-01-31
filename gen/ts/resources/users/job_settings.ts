// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/users/job_settings.proto" (package "resources.users", syntax proto3)
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
 * @generated from protobuf message resources.users.DiscordSyncSettings
 */
export interface DiscordSyncSettings {
    /**
     * @generated from protobuf field: bool dry_run = 1;
     */
    dryRun: boolean;
    /**
     * @generated from protobuf field: bool user_info_sync = 2;
     */
    userInfoSync: boolean;
    /**
     * @generated from protobuf field: resources.users.UserInfoSyncSettings user_info_sync_settings = 3;
     */
    userInfoSyncSettings?: UserInfoSyncSettings;
    /**
     * @generated from protobuf field: bool status_log = 4;
     */
    statusLog: boolean;
    /**
     * @generated from protobuf field: resources.users.StatusLogSettings status_log_settings = 5;
     */
    statusLogSettings?: StatusLogSettings;
    /**
     * @generated from protobuf field: bool jobs_absence = 6;
     */
    jobsAbsence: boolean;
    /**
     * @generated from protobuf field: resources.users.JobsAbsenceSettings jobs_absence_settings = 7;
     */
    jobsAbsenceSettings?: JobsAbsenceSettings;
    /**
     * @generated from protobuf field: resources.users.GroupSyncSettings group_sync_settings = 8;
     */
    groupSyncSettings?: GroupSyncSettings;
    /**
     * @generated from protobuf field: string qualifications_role_format = 9;
     */
    qualificationsRoleFormat: string;
}
/**
 * @generated from protobuf message resources.users.DiscordSyncChanges
 */
export interface DiscordSyncChanges {
    /**
     * @generated from protobuf field: repeated resources.users.DiscordSyncChange changes = 1;
     */
    changes: DiscordSyncChange[];
}
/**
 * @generated from protobuf message resources.users.DiscordSyncChange
 */
export interface DiscordSyncChange {
    /**
     * @generated from protobuf field: resources.timestamp.Timestamp time = 1;
     */
    time?: Timestamp;
    /**
     * @generated from protobuf field: string plan = 2;
     */
    plan: string;
}
/**
 * @generated from protobuf message resources.users.UserInfoSyncSettings
 */
export interface UserInfoSyncSettings {
    /**
     * @generated from protobuf field: bool employee_role_enabled = 1;
     */
    employeeRoleEnabled: boolean;
    /**
     * @generated from protobuf field: string employee_role_format = 2;
     */
    employeeRoleFormat: string;
    /**
     * @generated from protobuf field: string grade_role_format = 3;
     */
    gradeRoleFormat: string;
    /**
     * @generated from protobuf field: bool unemployed_enabled = 4;
     */
    unemployedEnabled: boolean;
    /**
     * @generated from protobuf field: resources.users.UserInfoSyncUnemployedMode unemployed_mode = 5;
     */
    unemployedMode: UserInfoSyncUnemployedMode;
    /**
     * @generated from protobuf field: string unemployed_role_name = 6;
     */
    unemployedRoleName: string;
    /**
     * @generated from protobuf field: bool sync_nicknames = 7;
     */
    syncNicknames: boolean;
    /**
     * @generated from protobuf field: repeated resources.users.GroupMapping group_mapping = 8;
     */
    groupMapping: GroupMapping[];
}
/**
 * @generated from protobuf message resources.users.GroupMapping
 */
export interface GroupMapping {
    /**
     * @generated from protobuf field: string name = 1;
     */
    name: string;
    /**
     * @generated from protobuf field: int32 from_grade = 2;
     */
    fromGrade: number;
    /**
     * @generated from protobuf field: int32 to_grade = 3;
     */
    toGrade: number;
}
/**
 * @generated from protobuf message resources.users.StatusLogSettings
 */
export interface StatusLogSettings {
    /**
     * @generated from protobuf field: string channel_id = 1;
     */
    channelId: string;
}
/**
 * @generated from protobuf message resources.users.JobsAbsenceSettings
 */
export interface JobsAbsenceSettings {
    /**
     * @generated from protobuf field: string absence_role = 1;
     */
    absenceRole: string;
}
/**
 * @generated from protobuf message resources.users.GroupSyncSettings
 */
export interface GroupSyncSettings {
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: repeated string ignored_role_ids = 1;
     */
    ignoredRoleIds: string[];
}
/**
 * @generated from protobuf message resources.users.JobSettings
 */
export interface JobSettings {
    /**
     * @generated from protobuf field: int32 absence_past_days = 1;
     */
    absencePastDays: number;
}
/**
 * @generated from protobuf enum resources.users.UserInfoSyncUnemployedMode
 */
export enum UserInfoSyncUnemployedMode {
    /**
     * @generated from protobuf enum value: USER_INFO_SYNC_UNEMPLOYED_MODE_UNSPECIFIED = 0;
     */
    UNSPECIFIED = 0,
    /**
     * @generated from protobuf enum value: USER_INFO_SYNC_UNEMPLOYED_MODE_GIVE_ROLE = 1;
     */
    GIVE_ROLE = 1,
    /**
     * @generated from protobuf enum value: USER_INFO_SYNC_UNEMPLOYED_MODE_KICK = 2;
     */
    KICK = 2
}
// @generated message type with reflection information, may provide speed optimized methods
class DiscordSyncSettings$Type extends MessageType<DiscordSyncSettings> {
    constructor() {
        super("resources.users.DiscordSyncSettings", [
            { no: 1, name: "dry_run", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 2, name: "user_info_sync", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 3, name: "user_info_sync_settings", kind: "message", T: () => UserInfoSyncSettings },
            { no: 4, name: "status_log", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 5, name: "status_log_settings", kind: "message", T: () => StatusLogSettings },
            { no: 6, name: "jobs_absence", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 7, name: "jobs_absence_settings", kind: "message", T: () => JobsAbsenceSettings },
            { no: 8, name: "group_sync_settings", kind: "message", T: () => GroupSyncSettings },
            { no: 9, name: "qualifications_role_format", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "64" } } } }
        ]);
    }
    create(value?: PartialMessage<DiscordSyncSettings>): DiscordSyncSettings {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.dryRun = false;
        message.userInfoSync = false;
        message.statusLog = false;
        message.jobsAbsence = false;
        message.qualificationsRoleFormat = "";
        if (value !== undefined)
            reflectionMergePartial<DiscordSyncSettings>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DiscordSyncSettings): DiscordSyncSettings {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* bool dry_run */ 1:
                    message.dryRun = reader.bool();
                    break;
                case /* bool user_info_sync */ 2:
                    message.userInfoSync = reader.bool();
                    break;
                case /* resources.users.UserInfoSyncSettings user_info_sync_settings */ 3:
                    message.userInfoSyncSettings = UserInfoSyncSettings.internalBinaryRead(reader, reader.uint32(), options, message.userInfoSyncSettings);
                    break;
                case /* bool status_log */ 4:
                    message.statusLog = reader.bool();
                    break;
                case /* resources.users.StatusLogSettings status_log_settings */ 5:
                    message.statusLogSettings = StatusLogSettings.internalBinaryRead(reader, reader.uint32(), options, message.statusLogSettings);
                    break;
                case /* bool jobs_absence */ 6:
                    message.jobsAbsence = reader.bool();
                    break;
                case /* resources.users.JobsAbsenceSettings jobs_absence_settings */ 7:
                    message.jobsAbsenceSettings = JobsAbsenceSettings.internalBinaryRead(reader, reader.uint32(), options, message.jobsAbsenceSettings);
                    break;
                case /* resources.users.GroupSyncSettings group_sync_settings */ 8:
                    message.groupSyncSettings = GroupSyncSettings.internalBinaryRead(reader, reader.uint32(), options, message.groupSyncSettings);
                    break;
                case /* string qualifications_role_format */ 9:
                    message.qualificationsRoleFormat = reader.string();
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
    internalBinaryWrite(message: DiscordSyncSettings, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* bool dry_run = 1; */
        if (message.dryRun !== false)
            writer.tag(1, WireType.Varint).bool(message.dryRun);
        /* bool user_info_sync = 2; */
        if (message.userInfoSync !== false)
            writer.tag(2, WireType.Varint).bool(message.userInfoSync);
        /* resources.users.UserInfoSyncSettings user_info_sync_settings = 3; */
        if (message.userInfoSyncSettings)
            UserInfoSyncSettings.internalBinaryWrite(message.userInfoSyncSettings, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* bool status_log = 4; */
        if (message.statusLog !== false)
            writer.tag(4, WireType.Varint).bool(message.statusLog);
        /* resources.users.StatusLogSettings status_log_settings = 5; */
        if (message.statusLogSettings)
            StatusLogSettings.internalBinaryWrite(message.statusLogSettings, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* bool jobs_absence = 6; */
        if (message.jobsAbsence !== false)
            writer.tag(6, WireType.Varint).bool(message.jobsAbsence);
        /* resources.users.JobsAbsenceSettings jobs_absence_settings = 7; */
        if (message.jobsAbsenceSettings)
            JobsAbsenceSettings.internalBinaryWrite(message.jobsAbsenceSettings, writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        /* resources.users.GroupSyncSettings group_sync_settings = 8; */
        if (message.groupSyncSettings)
            GroupSyncSettings.internalBinaryWrite(message.groupSyncSettings, writer.tag(8, WireType.LengthDelimited).fork(), options).join();
        /* string qualifications_role_format = 9; */
        if (message.qualificationsRoleFormat !== "")
            writer.tag(9, WireType.LengthDelimited).string(message.qualificationsRoleFormat);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.DiscordSyncSettings
 */
export const DiscordSyncSettings = new DiscordSyncSettings$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DiscordSyncChanges$Type extends MessageType<DiscordSyncChanges> {
    constructor() {
        super("resources.users.DiscordSyncChanges", [
            { no: 1, name: "changes", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => DiscordSyncChange }
        ]);
    }
    create(value?: PartialMessage<DiscordSyncChanges>): DiscordSyncChanges {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.changes = [];
        if (value !== undefined)
            reflectionMergePartial<DiscordSyncChanges>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DiscordSyncChanges): DiscordSyncChanges {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.users.DiscordSyncChange changes */ 1:
                    message.changes.push(DiscordSyncChange.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: DiscordSyncChanges, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.users.DiscordSyncChange changes = 1; */
        for (let i = 0; i < message.changes.length; i++)
            DiscordSyncChange.internalBinaryWrite(message.changes[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.DiscordSyncChanges
 */
export const DiscordSyncChanges = new DiscordSyncChanges$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DiscordSyncChange$Type extends MessageType<DiscordSyncChange> {
    constructor() {
        super("resources.users.DiscordSyncChange", [
            { no: 1, name: "time", kind: "message", T: () => Timestamp },
            { no: 2, name: "plan", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<DiscordSyncChange>): DiscordSyncChange {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.plan = "";
        if (value !== undefined)
            reflectionMergePartial<DiscordSyncChange>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: DiscordSyncChange): DiscordSyncChange {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.timestamp.Timestamp time */ 1:
                    message.time = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.time);
                    break;
                case /* string plan */ 2:
                    message.plan = reader.string();
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
    internalBinaryWrite(message: DiscordSyncChange, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.timestamp.Timestamp time = 1; */
        if (message.time)
            Timestamp.internalBinaryWrite(message.time, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* string plan = 2; */
        if (message.plan !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.plan);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.DiscordSyncChange
 */
export const DiscordSyncChange = new DiscordSyncChange$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UserInfoSyncSettings$Type extends MessageType<UserInfoSyncSettings> {
    constructor() {
        super("resources.users.UserInfoSyncSettings", [
            { no: 1, name: "employee_role_enabled", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 2, name: "employee_role_format", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "64" } } } },
            { no: 3, name: "grade_role_format", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "64" } } } },
            { no: 4, name: "unemployed_enabled", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 5, name: "unemployed_mode", kind: "enum", T: () => ["resources.users.UserInfoSyncUnemployedMode", UserInfoSyncUnemployedMode, "USER_INFO_SYNC_UNEMPLOYED_MODE_"], options: { "validate.rules": { enum: { definedOnly: true } } } },
            { no: 6, name: "unemployed_role_name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "64" } } } },
            { no: 7, name: "sync_nicknames", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 8, name: "group_mapping", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => GroupMapping }
        ]);
    }
    create(value?: PartialMessage<UserInfoSyncSettings>): UserInfoSyncSettings {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.employeeRoleEnabled = false;
        message.employeeRoleFormat = "";
        message.gradeRoleFormat = "";
        message.unemployedEnabled = false;
        message.unemployedMode = 0;
        message.unemployedRoleName = "";
        message.syncNicknames = false;
        message.groupMapping = [];
        if (value !== undefined)
            reflectionMergePartial<UserInfoSyncSettings>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UserInfoSyncSettings): UserInfoSyncSettings {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* bool employee_role_enabled */ 1:
                    message.employeeRoleEnabled = reader.bool();
                    break;
                case /* string employee_role_format */ 2:
                    message.employeeRoleFormat = reader.string();
                    break;
                case /* string grade_role_format */ 3:
                    message.gradeRoleFormat = reader.string();
                    break;
                case /* bool unemployed_enabled */ 4:
                    message.unemployedEnabled = reader.bool();
                    break;
                case /* resources.users.UserInfoSyncUnemployedMode unemployed_mode */ 5:
                    message.unemployedMode = reader.int32();
                    break;
                case /* string unemployed_role_name */ 6:
                    message.unemployedRoleName = reader.string();
                    break;
                case /* bool sync_nicknames */ 7:
                    message.syncNicknames = reader.bool();
                    break;
                case /* repeated resources.users.GroupMapping group_mapping */ 8:
                    message.groupMapping.push(GroupMapping.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: UserInfoSyncSettings, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* bool employee_role_enabled = 1; */
        if (message.employeeRoleEnabled !== false)
            writer.tag(1, WireType.Varint).bool(message.employeeRoleEnabled);
        /* string employee_role_format = 2; */
        if (message.employeeRoleFormat !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.employeeRoleFormat);
        /* string grade_role_format = 3; */
        if (message.gradeRoleFormat !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.gradeRoleFormat);
        /* bool unemployed_enabled = 4; */
        if (message.unemployedEnabled !== false)
            writer.tag(4, WireType.Varint).bool(message.unemployedEnabled);
        /* resources.users.UserInfoSyncUnemployedMode unemployed_mode = 5; */
        if (message.unemployedMode !== 0)
            writer.tag(5, WireType.Varint).int32(message.unemployedMode);
        /* string unemployed_role_name = 6; */
        if (message.unemployedRoleName !== "")
            writer.tag(6, WireType.LengthDelimited).string(message.unemployedRoleName);
        /* bool sync_nicknames = 7; */
        if (message.syncNicknames !== false)
            writer.tag(7, WireType.Varint).bool(message.syncNicknames);
        /* repeated resources.users.GroupMapping group_mapping = 8; */
        for (let i = 0; i < message.groupMapping.length; i++)
            GroupMapping.internalBinaryWrite(message.groupMapping[i], writer.tag(8, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.UserInfoSyncSettings
 */
export const UserInfoSyncSettings = new UserInfoSyncSettings$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GroupMapping$Type extends MessageType<GroupMapping> {
    constructor() {
        super("resources.users.GroupMapping", [
            { no: 1, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "64" } } } },
            { no: 2, name: "from_grade", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gte: 0 } } } },
            { no: 3, name: "to_grade", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gte: 0 } } } }
        ]);
    }
    create(value?: PartialMessage<GroupMapping>): GroupMapping {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.name = "";
        message.fromGrade = 0;
        message.toGrade = 0;
        if (value !== undefined)
            reflectionMergePartial<GroupMapping>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: GroupMapping): GroupMapping {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string name */ 1:
                    message.name = reader.string();
                    break;
                case /* int32 from_grade */ 2:
                    message.fromGrade = reader.int32();
                    break;
                case /* int32 to_grade */ 3:
                    message.toGrade = reader.int32();
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
    internalBinaryWrite(message: GroupMapping, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string name = 1; */
        if (message.name !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.name);
        /* int32 from_grade = 2; */
        if (message.fromGrade !== 0)
            writer.tag(2, WireType.Varint).int32(message.fromGrade);
        /* int32 to_grade = 3; */
        if (message.toGrade !== 0)
            writer.tag(3, WireType.Varint).int32(message.toGrade);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.GroupMapping
 */
export const GroupMapping = new GroupMapping$Type();
// @generated message type with reflection information, may provide speed optimized methods
class StatusLogSettings$Type extends MessageType<StatusLogSettings> {
    constructor() {
        super("resources.users.StatusLogSettings", [
            { no: 1, name: "channel_id", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<StatusLogSettings>): StatusLogSettings {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.channelId = "";
        if (value !== undefined)
            reflectionMergePartial<StatusLogSettings>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: StatusLogSettings): StatusLogSettings {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string channel_id */ 1:
                    message.channelId = reader.string();
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
    internalBinaryWrite(message: StatusLogSettings, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string channel_id = 1; */
        if (message.channelId !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.channelId);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.StatusLogSettings
 */
export const StatusLogSettings = new StatusLogSettings$Type();
// @generated message type with reflection information, may provide speed optimized methods
class JobsAbsenceSettings$Type extends MessageType<JobsAbsenceSettings> {
    constructor() {
        super("resources.users.JobsAbsenceSettings", [
            { no: 1, name: "absence_role", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "64" } } } }
        ]);
    }
    create(value?: PartialMessage<JobsAbsenceSettings>): JobsAbsenceSettings {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.absenceRole = "";
        if (value !== undefined)
            reflectionMergePartial<JobsAbsenceSettings>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: JobsAbsenceSettings): JobsAbsenceSettings {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string absence_role */ 1:
                    message.absenceRole = reader.string();
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
    internalBinaryWrite(message: JobsAbsenceSettings, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string absence_role = 1; */
        if (message.absenceRole !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.absenceRole);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.JobsAbsenceSettings
 */
export const JobsAbsenceSettings = new JobsAbsenceSettings$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GroupSyncSettings$Type extends MessageType<GroupSyncSettings> {
    constructor() {
        super("resources.users.GroupSyncSettings", [
            { no: 1, name: "ignored_role_ids", kind: "scalar", repeat: 2 /*RepeatType.UNPACKED*/, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { repeated: { maxItems: "25", items: { string: { maxLen: "24" } } } } } }
        ]);
    }
    create(value?: PartialMessage<GroupSyncSettings>): GroupSyncSettings {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.ignoredRoleIds = [];
        if (value !== undefined)
            reflectionMergePartial<GroupSyncSettings>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: GroupSyncSettings): GroupSyncSettings {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated string ignored_role_ids */ 1:
                    message.ignoredRoleIds.push(reader.string());
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
    internalBinaryWrite(message: GroupSyncSettings, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated string ignored_role_ids = 1; */
        for (let i = 0; i < message.ignoredRoleIds.length; i++)
            writer.tag(1, WireType.LengthDelimited).string(message.ignoredRoleIds[i]);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.GroupSyncSettings
 */
export const GroupSyncSettings = new GroupSyncSettings$Type();
// @generated message type with reflection information, may provide speed optimized methods
class JobSettings$Type extends MessageType<JobSettings> {
    constructor() {
        super("resources.users.JobSettings", [
            { no: 1, name: "absence_past_days", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gte: 0 } } } }
        ]);
    }
    create(value?: PartialMessage<JobSettings>): JobSettings {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.absencePastDays = 0;
        if (value !== undefined)
            reflectionMergePartial<JobSettings>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: JobSettings): JobSettings {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 absence_past_days */ 1:
                    message.absencePastDays = reader.int32();
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
    internalBinaryWrite(message: JobSettings, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* int32 absence_past_days = 1; */
        if (message.absencePastDays !== 0)
            writer.tag(1, WireType.Varint).int32(message.absencePastDays);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.JobSettings
 */
export const JobSettings = new JobSettings$Type();
