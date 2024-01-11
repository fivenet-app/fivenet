// @generated by protobuf-ts 2.9.3 with parameter optimize_code_size,long_type_bigint
// @generated from protobuf file "resources/users/jobs.proto" (package "resources.users", syntax proto3)
// tslint:disable
import { MessageType } from "@protobuf-ts/runtime";
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.users.Job
 */
export interface Job {
    /**
     * @generated from protobuf field: string name = 1;
     */
    name: string; // @gotags: sql:"primary_key" alias:"name"
    /**
     * @generated from protobuf field: string label = 2;
     */
    label: string;
    /**
     * @generated from protobuf field: repeated resources.users.JobGrade grades = 3;
     */
    grades: JobGrade[];
}
/**
 * @generated from protobuf message resources.users.JobGrade
 */
export interface JobGrade {
    /**
     * @generated from protobuf field: optional string job_name = 1;
     */
    jobName?: string;
    /**
     * @generated from protobuf field: int32 grade = 2;
     */
    grade: number;
    /**
     * @generated from protobuf field: string label = 3;
     */
    label: string;
}
/**
 * @generated from protobuf message resources.users.JobProps
 */
export interface JobProps {
    /**
     * @generated from protobuf field: string job = 1;
     */
    job: string;
    /**
     * @generated from protobuf field: string theme = 2;
     */
    theme: string;
    /**
     * @generated from protobuf field: string livemap_marker_color = 3;
     */
    livemapMarkerColor: string;
    /**
     * @generated from protobuf field: resources.users.QuickButtons quick_buttons = 4;
     */
    quickButtons?: QuickButtons;
    /**
     * @generated from protobuf field: optional string radio_frequency = 5;
     */
    radioFrequency?: string;
    /**
     * @generated from protobuf field: optional uint64 discord_guild_id = 6 [jstype = JS_STRING];
     */
    discordGuildId?: string;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp discord_last_sync = 7;
     */
    discordLastSync?: Timestamp;
    /**
     * @generated from protobuf field: resources.users.DiscordSyncSettings discord_sync_settings = 8;
     */
    discordSyncSettings?: DiscordSyncSettings;
}
/**
 * @generated from protobuf message resources.users.QuickButtons
 */
export interface QuickButtons {
    /**
     * @generated from protobuf field: bool penalty_calculator = 1;
     */
    penaltyCalculator: boolean;
    /**
     * @generated from protobuf field: bool body_checkup = 2;
     */
    bodyCheckup: boolean;
}
/**
 * @generated from protobuf message resources.users.DiscordSyncSettings
 */
export interface DiscordSyncSettings {
    /**
     * @generated from protobuf field: bool user_info_sync = 1;
     */
    userInfoSync: boolean;
    /**
     * @generated from protobuf field: optional resources.users.UserInfoSyncSettings user_info_sync_settings = 2;
     */
    userInfoSyncSettings?: UserInfoSyncSettings;
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
     * @generated from protobuf field: optional string employee_role_format = 2;
     */
    employeeRoleFormat?: string;
}
// @generated message type with reflection information, may provide speed optimized methods
class Job$Type extends MessageType<Job> {
    constructor() {
        super("resources.users.Job", [
            { no: 1, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 2, name: "label", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 3, name: "grades", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => JobGrade }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.users.Job
 */
export const Job = new Job$Type();
// @generated message type with reflection information, may provide speed optimized methods
class JobGrade$Type extends MessageType<JobGrade> {
    constructor() {
        super("resources.users.JobGrade", [
            { no: 1, name: "job_name", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 2, name: "grade", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gt: 0 } } } },
            { no: 3, name: "label", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.users.JobGrade
 */
export const JobGrade = new JobGrade$Type();
// @generated message type with reflection information, may provide speed optimized methods
class JobProps$Type extends MessageType<JobProps> {
    constructor() {
        super("resources.users.JobProps", [
            { no: 1, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 2, name: "theme", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 3, name: "livemap_marker_color", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { len: "6", pattern: "^[A-Fa-f0-9]{6}$" } } } },
            { no: 4, name: "quick_buttons", kind: "message", T: () => QuickButtons },
            { no: 5, name: "radio_frequency", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "6" } } } },
            { no: 6, name: "discord_guild_id", kind: "scalar", opt: true, T: 4 /*ScalarType.UINT64*/ },
            { no: 7, name: "discord_last_sync", kind: "message", T: () => Timestamp },
            { no: 8, name: "discord_sync_settings", kind: "message", T: () => DiscordSyncSettings }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.users.JobProps
 */
export const JobProps = new JobProps$Type();
// @generated message type with reflection information, may provide speed optimized methods
class QuickButtons$Type extends MessageType<QuickButtons> {
    constructor() {
        super("resources.users.QuickButtons", [
            { no: 1, name: "penalty_calculator", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 2, name: "body_checkup", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.users.QuickButtons
 */
export const QuickButtons = new QuickButtons$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DiscordSyncSettings$Type extends MessageType<DiscordSyncSettings> {
    constructor() {
        super("resources.users.DiscordSyncSettings", [
            { no: 1, name: "user_info_sync", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 2, name: "user_info_sync_settings", kind: "message", T: () => UserInfoSyncSettings }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.users.DiscordSyncSettings
 */
export const DiscordSyncSettings = new DiscordSyncSettings$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UserInfoSyncSettings$Type extends MessageType<UserInfoSyncSettings> {
    constructor() {
        super("resources.users.UserInfoSyncSettings", [
            { no: 1, name: "employee_role_enabled", kind: "scalar", T: 8 /*ScalarType.BOOL*/ },
            { no: 2, name: "employee_role_format", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.users.UserInfoSyncSettings
 */
export const UserInfoSyncSettings = new UserInfoSyncSettings$Type();