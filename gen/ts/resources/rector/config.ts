// @generated by protobuf-ts 2.9.3 with parameter optimize_speed,long_type_bigint
// @generated from protobuf file "resources/rector/config.proto" (package "resources.rector", syntax proto3)
// tslint:disable
import type { BinaryReadOptions, BinaryWriteOptions, IBinaryReader, IBinaryWriter, PartialMessage } from '@protobuf-ts/runtime';
import { MessageType, UnknownFieldHandler, WireType, reflectionMergePartial } from '@protobuf-ts/runtime';
import { Duration } from '../../google/protobuf/duration';
/**
 * @generated from protobuf message resources.rector.AppConfig
 */
export interface AppConfig {
    /**
     * @generated from protobuf field: resources.rector.Auth auth = 1;
     */
    auth?: Auth;
    /**
     * @generated from protobuf field: resources.rector.Perms perms = 2;
     */
    perms?: Perms;
    /**
     * @generated from protobuf field: resources.rector.Website website = 3;
     */
    website?: Website;
    /**
     * @generated from protobuf field: resources.rector.JobInfo job_info = 4;
     */
    jobInfo?: JobInfo;
    /**
     * @generated from protobuf field: resources.rector.UserTracker user_tracker = 5;
     */
    userTracker?: UserTracker;
    /**
     * @generated from protobuf field: resources.rector.Discord discord = 6;
     */
    discord?: Discord;
}
/**
 * @generated from protobuf message resources.rector.Auth
 */
export interface Auth {
    /**
     * @generated from protobuf field: bool signup_enabled = 1;
     */
    signupEnabled: boolean;
}
/**
 * @generated from protobuf message resources.rector.Perms
 */
export interface Perms {
    /**
     * @generated from protobuf field: repeated resources.rector.Perm default = 1;
     */
    default: Perm[];
}
/**
 * @generated from protobuf message resources.rector.Perm
 */
export interface Perm {
    /**
     * @generated from protobuf field: string category = 1;
     */
    category: string;
    /**
     * @generated from protobuf field: string name = 2;
     */
    name: string;
}
/**
 * @generated from protobuf message resources.rector.Website
 */
export interface Website {
    /**
     * @generated from protobuf field: resources.rector.Links links = 1;
     */
    links?: Links;
}
/**
 * @generated from protobuf message resources.rector.Links
 */
export interface Links {
    /**
     * @generated from protobuf field: optional string privacy_policy = 1;
     */
    privacyPolicy?: string;
    /**
     * @generated from protobuf field: optional string imprint = 2;
     */
    imprint?: string;
}
/**
 * @generated from protobuf message resources.rector.JobInfo
 */
export interface JobInfo {
    /**
     * @generated from protobuf field: resources.rector.UnemployedJob unemployed_job = 1;
     */
    unemployedJob?: UnemployedJob;
    /**
     * @generated from protobuf field: repeated string public_jobs = 2;
     */
    publicJobs: string[];
    /**
     * @generated from protobuf field: repeated string hidden_jobs = 3;
     */
    hiddenJobs: string[];
}
/**
 * @generated from protobuf message resources.rector.UnemployedJob
 */
export interface UnemployedJob {
    /**
     * @generated from protobuf field: string name = 1;
     */
    name: string;
    /**
     * @generated from protobuf field: int32 grade = 2;
     */
    grade: number;
}
/**
 * @generated from protobuf message resources.rector.UserTracker
 */
export interface UserTracker {
    /**
     * @generated from protobuf field: google.protobuf.Duration refresh_time = 1;
     */
    refreshTime?: Duration;
    /**
     * @generated from protobuf field: google.protobuf.Duration db_refresh_time = 2;
     */
    dbRefreshTime?: Duration;
    /**
     * @generated from protobuf field: repeated string livemap_jobs = 3;
     */
    livemapJobs: string[];
}
/**
 * @generated from protobuf message resources.rector.Discord
 */
export interface Discord {
    /**
     * @generated from protobuf field: bool enabled = 1;
     */
    enabled: boolean;
    /**
     * @generated from protobuf field: google.protobuf.Duration sync_interval = 2;
     */
    syncInterval?: Duration;
    /**
     * @generated from protobuf field: optional string invite_url = 3;
     */
    inviteUrl?: string;
}
// @generated message type with reflection information, may provide speed optimized methods
class AppConfig$Type extends MessageType<AppConfig> {
    constructor() {
        super('resources.rector.AppConfig', [
            { no: 1, name: 'auth', kind: 'message', T: () => Auth },
            { no: 2, name: 'perms', kind: 'message', T: () => Perms },
            { no: 3, name: 'website', kind: 'message', T: () => Website },
            { no: 4, name: 'job_info', kind: 'message', T: () => JobInfo },
            { no: 5, name: 'user_tracker', kind: 'message', T: () => UserTracker },
            { no: 6, name: 'discord', kind: 'message', T: () => Discord },
        ]);
    }
    create(value?: PartialMessage<AppConfig>): AppConfig {
        const message = globalThis.Object.create(this.messagePrototype!);
        if (value !== undefined) reflectionMergePartial<AppConfig>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: AppConfig): AppConfig {
        let message = target ?? this.create(),
            end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.rector.Auth auth */ 1:
                    message.auth = Auth.internalBinaryRead(reader, reader.uint32(), options, message.auth);
                    break;
                case /* resources.rector.Perms perms */ 2:
                    message.perms = Perms.internalBinaryRead(reader, reader.uint32(), options, message.perms);
                    break;
                case /* resources.rector.Website website */ 3:
                    message.website = Website.internalBinaryRead(reader, reader.uint32(), options, message.website);
                    break;
                case /* resources.rector.JobInfo job_info */ 4:
                    message.jobInfo = JobInfo.internalBinaryRead(reader, reader.uint32(), options, message.jobInfo);
                    break;
                case /* resources.rector.UserTracker user_tracker */ 5:
                    message.userTracker = UserTracker.internalBinaryRead(reader, reader.uint32(), options, message.userTracker);
                    break;
                case /* resources.rector.Discord discord */ 6:
                    message.discord = Discord.internalBinaryRead(reader, reader.uint32(), options, message.discord);
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === 'throw')
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: AppConfig, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.rector.Auth auth = 1; */
        if (message.auth)
            Auth.internalBinaryWrite(message.auth, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* resources.rector.Perms perms = 2; */
        if (message.perms)
            Perms.internalBinaryWrite(message.perms, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* resources.rector.Website website = 3; */
        if (message.website)
            Website.internalBinaryWrite(message.website, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* resources.rector.JobInfo job_info = 4; */
        if (message.jobInfo)
            JobInfo.internalBinaryWrite(message.jobInfo, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* resources.rector.UserTracker user_tracker = 5; */
        if (message.userTracker)
            UserTracker.internalBinaryWrite(
                message.userTracker,
                writer.tag(5, WireType.LengthDelimited).fork(),
                options,
            ).join();
        /* resources.rector.Discord discord = 6; */
        if (message.discord)
            Discord.internalBinaryWrite(message.discord, writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false) (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.rector.AppConfig
 */
export const AppConfig = new AppConfig$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Auth$Type extends MessageType<Auth> {
    constructor() {
        super('resources.rector.Auth', [{ no: 1, name: 'signup_enabled', kind: 'scalar', T: 8 /*ScalarType.BOOL*/ }]);
    }
    create(value?: PartialMessage<Auth>): Auth {
        const message = globalThis.Object.create(this.messagePrototype!);
        message.signupEnabled = false;
        if (value !== undefined) reflectionMergePartial<Auth>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Auth): Auth {
        let message = target ?? this.create(),
            end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* bool signup_enabled */ 1:
                    message.signupEnabled = reader.bool();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === 'throw')
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: Auth, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* bool signup_enabled = 1; */
        if (message.signupEnabled !== false) writer.tag(1, WireType.Varint).bool(message.signupEnabled);
        let u = options.writeUnknownFields;
        if (u !== false) (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.rector.Auth
 */
export const Auth = new Auth$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Perms$Type extends MessageType<Perms> {
    constructor() {
        super('resources.rector.Perms', [
            {
                no: 1,
                name: 'default',
                kind: 'message',
                repeat: 1 /*RepeatType.PACKED*/,
                T: () => Perm,
                options: { 'validate.rules': { repeated: { maxItems: '100' } } },
            },
        ]);
    }
    create(value?: PartialMessage<Perms>): Perms {
        const message = globalThis.Object.create(this.messagePrototype!);
        message.default = [];
        if (value !== undefined) reflectionMergePartial<Perms>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Perms): Perms {
        let message = target ?? this.create(),
            end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.rector.Perm default */ 1:
                    message.default.push(Perm.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === 'throw')
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: Perms, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.rector.Perm default = 1; */
        for (let i = 0; i < message.default.length; i++)
            Perm.internalBinaryWrite(message.default[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false) (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.rector.Perms
 */
export const Perms = new Perms$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Perm$Type extends MessageType<Perm> {
    constructor() {
        super('resources.rector.Perm', [
            {
                no: 1,
                name: 'category',
                kind: 'scalar',
                T: 9 /*ScalarType.STRING*/,
                options: { 'validate.rules': { string: { maxLen: '128' } } },
            },
            {
                no: 2,
                name: 'name',
                kind: 'scalar',
                T: 9 /*ScalarType.STRING*/,
                options: { 'validate.rules': { string: { maxLen: '255' } } },
            },
        ]);
    }
    create(value?: PartialMessage<Perm>): Perm {
        const message = globalThis.Object.create(this.messagePrototype!);
        message.category = '';
        message.name = '';
        if (value !== undefined) reflectionMergePartial<Perm>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Perm): Perm {
        let message = target ?? this.create(),
            end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string category */ 1:
                    message.category = reader.string();
                    break;
                case /* string name */ 2:
                    message.name = reader.string();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === 'throw')
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: Perm, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string category = 1; */
        if (message.category !== '') writer.tag(1, WireType.LengthDelimited).string(message.category);
        /* string name = 2; */
        if (message.name !== '') writer.tag(2, WireType.LengthDelimited).string(message.name);
        let u = options.writeUnknownFields;
        if (u !== false) (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.rector.Perm
 */
export const Perm = new Perm$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Website$Type extends MessageType<Website> {
    constructor() {
        super('resources.rector.Website', [{ no: 1, name: 'links', kind: 'message', T: () => Links }]);
    }
    create(value?: PartialMessage<Website>): Website {
        const message = globalThis.Object.create(this.messagePrototype!);
        if (value !== undefined) reflectionMergePartial<Website>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Website): Website {
        let message = target ?? this.create(),
            end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.rector.Links links */ 1:
                    message.links = Links.internalBinaryRead(reader, reader.uint32(), options, message.links);
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === 'throw')
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: Website, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.rector.Links links = 1; */
        if (message.links)
            Links.internalBinaryWrite(message.links, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false) (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.rector.Website
 */
export const Website = new Website$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Links$Type extends MessageType<Links> {
    constructor() {
        super('resources.rector.Links', [
            {
                no: 1,
                name: 'privacy_policy',
                kind: 'scalar',
                opt: true,
                T: 9 /*ScalarType.STRING*/,
                options: { 'validate.rules': { string: { maxLen: '255' } } },
            },
            {
                no: 2,
                name: 'imprint',
                kind: 'scalar',
                opt: true,
                T: 9 /*ScalarType.STRING*/,
                options: { 'validate.rules': { string: { maxLen: '255' } } },
            },
        ]);
    }
    create(value?: PartialMessage<Links>): Links {
        const message = globalThis.Object.create(this.messagePrototype!);
        if (value !== undefined) reflectionMergePartial<Links>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Links): Links {
        let message = target ?? this.create(),
            end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* optional string privacy_policy */ 1:
                    message.privacyPolicy = reader.string();
                    break;
                case /* optional string imprint */ 2:
                    message.imprint = reader.string();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === 'throw')
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: Links, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* optional string privacy_policy = 1; */
        if (message.privacyPolicy !== undefined) writer.tag(1, WireType.LengthDelimited).string(message.privacyPolicy);
        /* optional string imprint = 2; */
        if (message.imprint !== undefined) writer.tag(2, WireType.LengthDelimited).string(message.imprint);
        let u = options.writeUnknownFields;
        if (u !== false) (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.rector.Links
 */
export const Links = new Links$Type();
// @generated message type with reflection information, may provide speed optimized methods
class JobInfo$Type extends MessageType<JobInfo> {
    constructor() {
        super('resources.rector.JobInfo', [
            {
                no: 1,
                name: 'unemployed_job',
                kind: 'message',
                T: () => UnemployedJob,
                options: { 'validate.rules': { message: { required: true } } },
            },
            {
                no: 2,
                name: 'public_jobs',
                kind: 'scalar',
                repeat: 2 /*RepeatType.UNPACKED*/,
                T: 9 /*ScalarType.STRING*/,
                options: { 'validate.rules': { repeated: { maxItems: '100' } } },
            },
            {
                no: 3,
                name: 'hidden_jobs',
                kind: 'scalar',
                repeat: 2 /*RepeatType.UNPACKED*/,
                T: 9 /*ScalarType.STRING*/,
                options: { 'validate.rules': { repeated: { maxItems: '100' } } },
            },
        ]);
    }
    create(value?: PartialMessage<JobInfo>): JobInfo {
        const message = globalThis.Object.create(this.messagePrototype!);
        message.publicJobs = [];
        message.hiddenJobs = [];
        if (value !== undefined) reflectionMergePartial<JobInfo>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: JobInfo): JobInfo {
        let message = target ?? this.create(),
            end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.rector.UnemployedJob unemployed_job */ 1:
                    message.unemployedJob = UnemployedJob.internalBinaryRead(
                        reader,
                        reader.uint32(),
                        options,
                        message.unemployedJob,
                    );
                    break;
                case /* repeated string public_jobs */ 2:
                    message.publicJobs.push(reader.string());
                    break;
                case /* repeated string hidden_jobs */ 3:
                    message.hiddenJobs.push(reader.string());
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === 'throw')
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: JobInfo, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.rector.UnemployedJob unemployed_job = 1; */
        if (message.unemployedJob)
            UnemployedJob.internalBinaryWrite(
                message.unemployedJob,
                writer.tag(1, WireType.LengthDelimited).fork(),
                options,
            ).join();
        /* repeated string public_jobs = 2; */
        for (let i = 0; i < message.publicJobs.length; i++)
            writer.tag(2, WireType.LengthDelimited).string(message.publicJobs[i]);
        /* repeated string hidden_jobs = 3; */
        for (let i = 0; i < message.hiddenJobs.length; i++)
            writer.tag(3, WireType.LengthDelimited).string(message.hiddenJobs[i]);
        let u = options.writeUnknownFields;
        if (u !== false) (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.rector.JobInfo
 */
export const JobInfo = new JobInfo$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UnemployedJob$Type extends MessageType<UnemployedJob> {
    constructor() {
        super('resources.rector.UnemployedJob', [
            {
                no: 1,
                name: 'name',
                kind: 'scalar',
                T: 9 /*ScalarType.STRING*/,
                options: { 'validate.rules': { string: { maxLen: '20' } } },
            },
            {
                no: 2,
                name: 'grade',
                kind: 'scalar',
                T: 5 /*ScalarType.INT32*/,
                options: { 'validate.rules': { int32: { gt: 0 } } },
            },
        ]);
    }
    create(value?: PartialMessage<UnemployedJob>): UnemployedJob {
        const message = globalThis.Object.create(this.messagePrototype!);
        message.name = '';
        message.grade = 0;
        if (value !== undefined) reflectionMergePartial<UnemployedJob>(this, message, value);
        return message;
    }
    internalBinaryRead(
        reader: IBinaryReader,
        length: number,
        options: BinaryReadOptions,
        target?: UnemployedJob,
    ): UnemployedJob {
        let message = target ?? this.create(),
            end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string name */ 1:
                    message.name = reader.string();
                    break;
                case /* int32 grade */ 2:
                    message.grade = reader.int32();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === 'throw')
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: UnemployedJob, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string name = 1; */
        if (message.name !== '') writer.tag(1, WireType.LengthDelimited).string(message.name);
        /* int32 grade = 2; */
        if (message.grade !== 0) writer.tag(2, WireType.Varint).int32(message.grade);
        let u = options.writeUnknownFields;
        if (u !== false) (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.rector.UnemployedJob
 */
export const UnemployedJob = new UnemployedJob$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UserTracker$Type extends MessageType<UserTracker> {
    constructor() {
        super('resources.rector.UserTracker', [
            {
                no: 1,
                name: 'refresh_time',
                kind: 'message',
                T: () => Duration,
                options: {
                    'validate.rules': { duration: { required: true, lt: { seconds: '60' }, gte: { nanos: 500000000 } } },
                },
            },
            {
                no: 2,
                name: 'db_refresh_time',
                kind: 'message',
                T: () => Duration,
                options: {
                    'validate.rules': { duration: { required: true, lt: { seconds: '60' }, gte: { nanos: 500000000 } } },
                },
            },
            {
                no: 3,
                name: 'livemap_jobs',
                kind: 'scalar',
                repeat: 2 /*RepeatType.UNPACKED*/,
                T: 9 /*ScalarType.STRING*/,
                options: { 'validate.rules': { repeated: { maxItems: '100' } } },
            },
        ]);
    }
    create(value?: PartialMessage<UserTracker>): UserTracker {
        const message = globalThis.Object.create(this.messagePrototype!);
        message.livemapJobs = [];
        if (value !== undefined) reflectionMergePartial<UserTracker>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UserTracker): UserTracker {
        let message = target ?? this.create(),
            end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* google.protobuf.Duration refresh_time */ 1:
                    message.refreshTime = Duration.internalBinaryRead(reader, reader.uint32(), options, message.refreshTime);
                    break;
                case /* google.protobuf.Duration db_refresh_time */ 2:
                    message.dbRefreshTime = Duration.internalBinaryRead(
                        reader,
                        reader.uint32(),
                        options,
                        message.dbRefreshTime,
                    );
                    break;
                case /* repeated string livemap_jobs */ 3:
                    message.livemapJobs.push(reader.string());
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === 'throw')
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: UserTracker, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* google.protobuf.Duration refresh_time = 1; */
        if (message.refreshTime)
            Duration.internalBinaryWrite(message.refreshTime, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* google.protobuf.Duration db_refresh_time = 2; */
        if (message.dbRefreshTime)
            Duration.internalBinaryWrite(message.dbRefreshTime, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* repeated string livemap_jobs = 3; */
        for (let i = 0; i < message.livemapJobs.length; i++)
            writer.tag(3, WireType.LengthDelimited).string(message.livemapJobs[i]);
        let u = options.writeUnknownFields;
        if (u !== false) (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.rector.UserTracker
 */
export const UserTracker = new UserTracker$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Discord$Type extends MessageType<Discord> {
    constructor() {
        super('resources.rector.Discord', [
            { no: 1, name: 'enabled', kind: 'scalar', T: 8 /*ScalarType.BOOL*/ },
            {
                no: 2,
                name: 'sync_interval',
                kind: 'message',
                T: () => Duration,
                options: {
                    'validate.rules': { duration: { required: true, lt: { seconds: '180000000' }, gte: { seconds: '6000' } } },
                },
            },
            {
                no: 3,
                name: 'invite_url',
                kind: 'scalar',
                opt: true,
                T: 9 /*ScalarType.STRING*/,
                options: { 'validate.rules': { string: { maxLen: '255' } } },
            },
        ]);
    }
    create(value?: PartialMessage<Discord>): Discord {
        const message = globalThis.Object.create(this.messagePrototype!);
        message.enabled = false;
        if (value !== undefined) reflectionMergePartial<Discord>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Discord): Discord {
        let message = target ?? this.create(),
            end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* bool enabled */ 1:
                    message.enabled = reader.bool();
                    break;
                case /* google.protobuf.Duration sync_interval */ 2:
                    message.syncInterval = Duration.internalBinaryRead(reader, reader.uint32(), options, message.syncInterval);
                    break;
                case /* optional string invite_url */ 3:
                    message.inviteUrl = reader.string();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === 'throw')
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: Discord, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* bool enabled = 1; */
        if (message.enabled !== false) writer.tag(1, WireType.Varint).bool(message.enabled);
        /* google.protobuf.Duration sync_interval = 2; */
        if (message.syncInterval)
            Duration.internalBinaryWrite(message.syncInterval, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* optional string invite_url = 3; */
        if (message.inviteUrl !== undefined) writer.tag(3, WireType.LengthDelimited).string(message.inviteUrl);
        let u = options.writeUnknownFields;
        if (u !== false) (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.rector.Discord
 */
export const Discord = new Discord$Type();
