// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/permissions/permissions.proto" (package "resources.permissions", syntax proto3)
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
 * @generated from protobuf message resources.permissions.Permission
 */
export interface Permission {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: number;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: string category = 3;
     */
    category: string;
    /**
     * @generated from protobuf field: string name = 4;
     */
    name: string;
    /**
     * @generated from protobuf field: string guard_name = 5;
     */
    guardName: string;
    /**
     * @generated from protobuf field: bool val = 6;
     */
    val: boolean;
}
/**
 * @generated from protobuf message resources.permissions.Role
 */
export interface Role {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: number;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: string job = 3;
     */
    job: string;
    /**
     * @generated from protobuf field: optional string job_label = 4;
     */
    jobLabel?: string;
    /**
     * @generated from protobuf field: int32 grade = 5;
     */
    grade: number;
    /**
     * @generated from protobuf field: optional string job_grade_label = 6;
     */
    jobGradeLabel?: string;
    /**
     * @generated from protobuf field: repeated resources.permissions.Permission permissions = 7;
     */
    permissions: Permission[];
    /**
     * @generated from protobuf field: repeated resources.permissions.RoleAttribute attributes = 8;
     */
    attributes: RoleAttribute[];
}
/**
 * @generated from protobuf message resources.permissions.RawRoleAttribute
 */
export interface RawRoleAttribute {
    /**
     * @generated from protobuf field: uint64 role_id = 1;
     */
    roleId: number;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: uint64 attr_id = 3;
     */
    attrId: number;
    /**
     * @generated from protobuf field: uint64 permission_id = 4;
     */
    permissionId: number;
    /**
     * @generated from protobuf field: string category = 5;
     */
    category: string;
    /**
     * @generated from protobuf field: string name = 6;
     */
    name: string;
    /**
     * @generated from protobuf field: string key = 7;
     */
    key: string;
    /**
     * @generated from protobuf field: string type = 8;
     */
    type: string;
    /**
     * @generated from protobuf field: resources.permissions.AttributeValues valid_values = 9;
     */
    validValues?: AttributeValues;
    /**
     * @generated from protobuf field: resources.permissions.AttributeValues value = 10;
     */
    value?: AttributeValues;
}
/**
 * @generated from protobuf message resources.permissions.RoleAttribute
 */
export interface RoleAttribute {
    /**
     * @generated from protobuf field: uint64 role_id = 1;
     */
    roleId: number;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp created_at = 2;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf field: uint64 attr_id = 3;
     */
    attrId: number;
    /**
     * @generated from protobuf field: uint64 permission_id = 4;
     */
    permissionId: number;
    /**
     * @generated from protobuf field: string category = 5;
     */
    category: string;
    /**
     * @generated from protobuf field: string name = 6;
     */
    name: string;
    /**
     * @generated from protobuf field: string key = 7;
     */
    key: string;
    /**
     * @generated from protobuf field: string type = 8;
     */
    type: string;
    /**
     * @generated from protobuf field: resources.permissions.AttributeValues valid_values = 9;
     */
    validValues?: AttributeValues;
    /**
     * @generated from protobuf field: resources.permissions.AttributeValues value = 10;
     */
    value?: AttributeValues;
    /**
     * @generated from protobuf field: optional resources.permissions.AttributeValues max_values = 11;
     */
    maxValues?: AttributeValues;
}
/**
 * @generated from protobuf message resources.permissions.AttributeValues
 */
export interface AttributeValues {
    /**
     * @generated from protobuf oneof: valid_values
     */
    validValues: {
        oneofKind: "stringList";
        /**
         * @generated from protobuf field: resources.permissions.StringList string_list = 1;
         */
        stringList: StringList;
    } | {
        oneofKind: "jobList";
        /**
         * @generated from protobuf field: resources.permissions.StringList job_list = 2;
         */
        jobList: StringList;
    } | {
        oneofKind: "jobGradeList";
        /**
         * @generated from protobuf field: resources.permissions.JobGradeList job_grade_list = 3;
         */
        jobGradeList: JobGradeList;
    } | {
        oneofKind: undefined;
    };
}
/**
 * @generated from protobuf message resources.permissions.StringList
 */
export interface StringList {
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: repeated string strings = 1;
     */
    strings: string[];
}
/**
 * @generated from protobuf message resources.permissions.JobGradeList
 */
export interface JobGradeList {
    /**
     * @generated from protobuf field: map<string, int32> jobs = 1;
     */
    jobs: {
        [key: string]: number;
    };
}
// @generated message type with reflection information, may provide speed optimized methods
class Permission$Type extends MessageType<Permission> {
    constructor() {
        super("resources.permissions.Permission", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "category", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "128" } } } },
            { no: 4, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 5, name: "guard_name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 6, name: "val", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<Permission>): Permission {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.category = "";
        message.name = "";
        message.guardName = "";
        message.val = false;
        if (value !== undefined)
            reflectionMergePartial<Permission>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Permission): Permission {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* optional resources.timestamp.Timestamp created_at */ 2:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* string category */ 3:
                    message.category = reader.string();
                    break;
                case /* string name */ 4:
                    message.name = reader.string();
                    break;
                case /* string guard_name */ 5:
                    message.guardName = reader.string();
                    break;
                case /* bool val */ 6:
                    message.val = reader.bool();
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
    internalBinaryWrite(message: Permission, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* optional resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* string category = 3; */
        if (message.category !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.category);
        /* string name = 4; */
        if (message.name !== "")
            writer.tag(4, WireType.LengthDelimited).string(message.name);
        /* string guard_name = 5; */
        if (message.guardName !== "")
            writer.tag(5, WireType.LengthDelimited).string(message.guardName);
        /* bool val = 6; */
        if (message.val !== false)
            writer.tag(6, WireType.Varint).bool(message.val);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.permissions.Permission
 */
export const Permission = new Permission$Type();
// @generated message type with reflection information, may provide speed optimized methods
class Role$Type extends MessageType<Role> {
    constructor() {
        super("resources.permissions.Role", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 4, name: "job_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 5, name: "grade", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gte: 0 } } } },
            { no: 6, name: "job_grade_label", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 7, name: "permissions", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Permission },
            { no: 8, name: "attributes", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => RoleAttribute }
        ]);
    }
    create(value?: PartialMessage<Role>): Role {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.id = 0;
        message.job = "";
        message.grade = 0;
        message.permissions = [];
        message.attributes = [];
        if (value !== undefined)
            reflectionMergePartial<Role>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Role): Role {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 id */ 1:
                    message.id = reader.uint64().toNumber();
                    break;
                case /* optional resources.timestamp.Timestamp created_at */ 2:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* string job */ 3:
                    message.job = reader.string();
                    break;
                case /* optional string job_label */ 4:
                    message.jobLabel = reader.string();
                    break;
                case /* int32 grade */ 5:
                    message.grade = reader.int32();
                    break;
                case /* optional string job_grade_label */ 6:
                    message.jobGradeLabel = reader.string();
                    break;
                case /* repeated resources.permissions.Permission permissions */ 7:
                    message.permissions.push(Permission.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated resources.permissions.RoleAttribute attributes */ 8:
                    message.attributes.push(RoleAttribute.internalBinaryRead(reader, reader.uint32(), options));
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
    internalBinaryWrite(message: Role, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 id = 1; */
        if (message.id !== 0)
            writer.tag(1, WireType.Varint).uint64(message.id);
        /* optional resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* string job = 3; */
        if (message.job !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.job);
        /* optional string job_label = 4; */
        if (message.jobLabel !== undefined)
            writer.tag(4, WireType.LengthDelimited).string(message.jobLabel);
        /* int32 grade = 5; */
        if (message.grade !== 0)
            writer.tag(5, WireType.Varint).int32(message.grade);
        /* optional string job_grade_label = 6; */
        if (message.jobGradeLabel !== undefined)
            writer.tag(6, WireType.LengthDelimited).string(message.jobGradeLabel);
        /* repeated resources.permissions.Permission permissions = 7; */
        for (let i = 0; i < message.permissions.length; i++)
            Permission.internalBinaryWrite(message.permissions[i], writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.permissions.RoleAttribute attributes = 8; */
        for (let i = 0; i < message.attributes.length; i++)
            RoleAttribute.internalBinaryWrite(message.attributes[i], writer.tag(8, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.permissions.Role
 */
export const Role = new Role$Type();
// @generated message type with reflection information, may provide speed optimized methods
class RawRoleAttribute$Type extends MessageType<RawRoleAttribute> {
    constructor() {
        super("resources.permissions.RawRoleAttribute", [
            { no: 1, name: "role_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "attr_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 4, name: "permission_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 5, name: "category", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "128" } } } },
            { no: 6, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 7, name: "key", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 8, name: "type", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 9, name: "valid_values", kind: "message", T: () => AttributeValues },
            { no: 10, name: "value", kind: "message", T: () => AttributeValues }
        ]);
    }
    create(value?: PartialMessage<RawRoleAttribute>): RawRoleAttribute {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.roleId = 0;
        message.attrId = 0;
        message.permissionId = 0;
        message.category = "";
        message.name = "";
        message.key = "";
        message.type = "";
        if (value !== undefined)
            reflectionMergePartial<RawRoleAttribute>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: RawRoleAttribute): RawRoleAttribute {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 role_id */ 1:
                    message.roleId = reader.uint64().toNumber();
                    break;
                case /* optional resources.timestamp.Timestamp created_at */ 2:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* uint64 attr_id */ 3:
                    message.attrId = reader.uint64().toNumber();
                    break;
                case /* uint64 permission_id */ 4:
                    message.permissionId = reader.uint64().toNumber();
                    break;
                case /* string category */ 5:
                    message.category = reader.string();
                    break;
                case /* string name */ 6:
                    message.name = reader.string();
                    break;
                case /* string key */ 7:
                    message.key = reader.string();
                    break;
                case /* string type */ 8:
                    message.type = reader.string();
                    break;
                case /* resources.permissions.AttributeValues valid_values */ 9:
                    message.validValues = AttributeValues.internalBinaryRead(reader, reader.uint32(), options, message.validValues);
                    break;
                case /* resources.permissions.AttributeValues value */ 10:
                    message.value = AttributeValues.internalBinaryRead(reader, reader.uint32(), options, message.value);
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
    internalBinaryWrite(message: RawRoleAttribute, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 role_id = 1; */
        if (message.roleId !== 0)
            writer.tag(1, WireType.Varint).uint64(message.roleId);
        /* optional resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* uint64 attr_id = 3; */
        if (message.attrId !== 0)
            writer.tag(3, WireType.Varint).uint64(message.attrId);
        /* uint64 permission_id = 4; */
        if (message.permissionId !== 0)
            writer.tag(4, WireType.Varint).uint64(message.permissionId);
        /* string category = 5; */
        if (message.category !== "")
            writer.tag(5, WireType.LengthDelimited).string(message.category);
        /* string name = 6; */
        if (message.name !== "")
            writer.tag(6, WireType.LengthDelimited).string(message.name);
        /* string key = 7; */
        if (message.key !== "")
            writer.tag(7, WireType.LengthDelimited).string(message.key);
        /* string type = 8; */
        if (message.type !== "")
            writer.tag(8, WireType.LengthDelimited).string(message.type);
        /* resources.permissions.AttributeValues valid_values = 9; */
        if (message.validValues)
            AttributeValues.internalBinaryWrite(message.validValues, writer.tag(9, WireType.LengthDelimited).fork(), options).join();
        /* resources.permissions.AttributeValues value = 10; */
        if (message.value)
            AttributeValues.internalBinaryWrite(message.value, writer.tag(10, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.permissions.RawRoleAttribute
 */
export const RawRoleAttribute = new RawRoleAttribute$Type();
// @generated message type with reflection information, may provide speed optimized methods
class RoleAttribute$Type extends MessageType<RoleAttribute> {
    constructor() {
        super("resources.permissions.RoleAttribute", [
            { no: 1, name: "role_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 2, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "attr_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 4, name: "permission_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 5, name: "category", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "128" } } } },
            { no: 6, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 7, name: "key", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 8, name: "type", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "255" } } } },
            { no: 9, name: "valid_values", kind: "message", T: () => AttributeValues },
            { no: 10, name: "value", kind: "message", T: () => AttributeValues },
            { no: 11, name: "max_values", kind: "message", T: () => AttributeValues }
        ]);
    }
    create(value?: PartialMessage<RoleAttribute>): RoleAttribute {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.roleId = 0;
        message.attrId = 0;
        message.permissionId = 0;
        message.category = "";
        message.name = "";
        message.key = "";
        message.type = "";
        if (value !== undefined)
            reflectionMergePartial<RoleAttribute>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: RoleAttribute): RoleAttribute {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* uint64 role_id */ 1:
                    message.roleId = reader.uint64().toNumber();
                    break;
                case /* optional resources.timestamp.Timestamp created_at */ 2:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* uint64 attr_id */ 3:
                    message.attrId = reader.uint64().toNumber();
                    break;
                case /* uint64 permission_id */ 4:
                    message.permissionId = reader.uint64().toNumber();
                    break;
                case /* string category */ 5:
                    message.category = reader.string();
                    break;
                case /* string name */ 6:
                    message.name = reader.string();
                    break;
                case /* string key */ 7:
                    message.key = reader.string();
                    break;
                case /* string type */ 8:
                    message.type = reader.string();
                    break;
                case /* resources.permissions.AttributeValues valid_values */ 9:
                    message.validValues = AttributeValues.internalBinaryRead(reader, reader.uint32(), options, message.validValues);
                    break;
                case /* resources.permissions.AttributeValues value */ 10:
                    message.value = AttributeValues.internalBinaryRead(reader, reader.uint32(), options, message.value);
                    break;
                case /* optional resources.permissions.AttributeValues max_values */ 11:
                    message.maxValues = AttributeValues.internalBinaryRead(reader, reader.uint32(), options, message.maxValues);
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
    internalBinaryWrite(message: RoleAttribute, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* uint64 role_id = 1; */
        if (message.roleId !== 0)
            writer.tag(1, WireType.Varint).uint64(message.roleId);
        /* optional resources.timestamp.Timestamp created_at = 2; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* uint64 attr_id = 3; */
        if (message.attrId !== 0)
            writer.tag(3, WireType.Varint).uint64(message.attrId);
        /* uint64 permission_id = 4; */
        if (message.permissionId !== 0)
            writer.tag(4, WireType.Varint).uint64(message.permissionId);
        /* string category = 5; */
        if (message.category !== "")
            writer.tag(5, WireType.LengthDelimited).string(message.category);
        /* string name = 6; */
        if (message.name !== "")
            writer.tag(6, WireType.LengthDelimited).string(message.name);
        /* string key = 7; */
        if (message.key !== "")
            writer.tag(7, WireType.LengthDelimited).string(message.key);
        /* string type = 8; */
        if (message.type !== "")
            writer.tag(8, WireType.LengthDelimited).string(message.type);
        /* resources.permissions.AttributeValues valid_values = 9; */
        if (message.validValues)
            AttributeValues.internalBinaryWrite(message.validValues, writer.tag(9, WireType.LengthDelimited).fork(), options).join();
        /* resources.permissions.AttributeValues value = 10; */
        if (message.value)
            AttributeValues.internalBinaryWrite(message.value, writer.tag(10, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.permissions.AttributeValues max_values = 11; */
        if (message.maxValues)
            AttributeValues.internalBinaryWrite(message.maxValues, writer.tag(11, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.permissions.RoleAttribute
 */
export const RoleAttribute = new RoleAttribute$Type();
// @generated message type with reflection information, may provide speed optimized methods
class AttributeValues$Type extends MessageType<AttributeValues> {
    constructor() {
        super("resources.permissions.AttributeValues", [
            { no: 1, name: "string_list", kind: "message", oneof: "validValues", T: () => StringList },
            { no: 2, name: "job_list", kind: "message", oneof: "validValues", T: () => StringList },
            { no: 3, name: "job_grade_list", kind: "message", oneof: "validValues", T: () => JobGradeList }
        ]);
    }
    create(value?: PartialMessage<AttributeValues>): AttributeValues {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.validValues = { oneofKind: undefined };
        if (value !== undefined)
            reflectionMergePartial<AttributeValues>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: AttributeValues): AttributeValues {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* resources.permissions.StringList string_list */ 1:
                    message.validValues = {
                        oneofKind: "stringList",
                        stringList: StringList.internalBinaryRead(reader, reader.uint32(), options, (message.validValues as any).stringList)
                    };
                    break;
                case /* resources.permissions.StringList job_list */ 2:
                    message.validValues = {
                        oneofKind: "jobList",
                        jobList: StringList.internalBinaryRead(reader, reader.uint32(), options, (message.validValues as any).jobList)
                    };
                    break;
                case /* resources.permissions.JobGradeList job_grade_list */ 3:
                    message.validValues = {
                        oneofKind: "jobGradeList",
                        jobGradeList: JobGradeList.internalBinaryRead(reader, reader.uint32(), options, (message.validValues as any).jobGradeList)
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
    internalBinaryWrite(message: AttributeValues, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* resources.permissions.StringList string_list = 1; */
        if (message.validValues.oneofKind === "stringList")
            StringList.internalBinaryWrite(message.validValues.stringList, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* resources.permissions.StringList job_list = 2; */
        if (message.validValues.oneofKind === "jobList")
            StringList.internalBinaryWrite(message.validValues.jobList, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* resources.permissions.JobGradeList job_grade_list = 3; */
        if (message.validValues.oneofKind === "jobGradeList")
            JobGradeList.internalBinaryWrite(message.validValues.jobGradeList, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.permissions.AttributeValues
 */
export const AttributeValues = new AttributeValues$Type();
// @generated message type with reflection information, may provide speed optimized methods
class StringList$Type extends MessageType<StringList> {
    constructor() {
        super("resources.permissions.StringList", [
            { no: 1, name: "strings", kind: "scalar", repeat: 2 /*RepeatType.UNPACKED*/, T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<StringList>): StringList {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.strings = [];
        if (value !== undefined)
            reflectionMergePartial<StringList>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: StringList): StringList {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated string strings */ 1:
                    message.strings.push(reader.string());
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
    internalBinaryWrite(message: StringList, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated string strings = 1; */
        for (let i = 0; i < message.strings.length; i++)
            writer.tag(1, WireType.LengthDelimited).string(message.strings[i]);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.permissions.StringList
 */
export const StringList = new StringList$Type();
// @generated message type with reflection information, may provide speed optimized methods
class JobGradeList$Type extends MessageType<JobGradeList> {
    constructor() {
        super("resources.permissions.JobGradeList", [
            { no: 1, name: "jobs", kind: "map", K: 9 /*ScalarType.STRING*/, V: { kind: "scalar", T: 5 /*ScalarType.INT32*/ } }
        ]);
    }
    create(value?: PartialMessage<JobGradeList>): JobGradeList {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.jobs = {};
        if (value !== undefined)
            reflectionMergePartial<JobGradeList>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: JobGradeList): JobGradeList {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* map<string, int32> jobs */ 1:
                    this.binaryReadMap1(message.jobs, reader, options);
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
    private binaryReadMap1(map: JobGradeList["jobs"], reader: IBinaryReader, options: BinaryReadOptions): void {
        let len = reader.uint32(), end = reader.pos + len, key: keyof JobGradeList["jobs"] | undefined, val: JobGradeList["jobs"][any] | undefined;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case 1:
                    key = reader.string();
                    break;
                case 2:
                    val = reader.int32();
                    break;
                default: throw new globalThis.Error("unknown map entry field for field resources.permissions.JobGradeList.jobs");
            }
        }
        map[key ?? ""] = val ?? 0;
    }
    internalBinaryWrite(message: JobGradeList, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* map<string, int32> jobs = 1; */
        for (let k of globalThis.Object.keys(message.jobs))
            writer.tag(1, WireType.LengthDelimited).fork().tag(1, WireType.LengthDelimited).string(k).tag(2, WireType.Varint).int32(message.jobs[k]).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.permissions.JobGradeList
 */
export const JobGradeList = new JobGradeList$Type();
