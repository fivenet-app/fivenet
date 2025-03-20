// @generated by protobuf-ts 2.9.6 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/users/props.proto" (package "resources.users", syntax proto3)
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
import { CitizenLabels } from "./labels";
import { File } from "../filestore/file";
import { JobGrade } from "./jobs";
import { Job } from "./jobs";
import { Timestamp } from "../timestamp/timestamp";
/**
 * @generated from protobuf message resources.users.UserProps
 */
export interface UserProps {
    /**
     * @generated from protobuf field: int32 user_id = 1;
     */
    userId: number;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp updated_at = 2;
     */
    updatedAt?: Timestamp;
    /**
     * @generated from protobuf field: optional bool wanted = 3;
     */
    wanted?: boolean;
    /**
     * @generated from protobuf field: optional string job_name = 4;
     */
    jobName?: string; // @gotags: alias:"job"
    /**
     * @generated from protobuf field: optional resources.users.Job job = 5;
     */
    job?: Job;
    /**
     * @generated from protobuf field: optional int32 job_grade_number = 6;
     */
    jobGradeNumber?: number; // @gotags: alias:"job_grade"
    /**
     * @generated from protobuf field: optional resources.users.JobGrade job_grade = 7;
     */
    jobGrade?: JobGrade;
    /**
     * @generated from protobuf field: optional uint32 traffic_infraction_points = 8;
     */
    trafficInfractionPoints?: number;
    /**
     * @generated from protobuf field: optional int64 open_fines = 9;
     */
    openFines?: number;
    /**
     * @generated from protobuf field: optional string blood_type = 10;
     */
    bloodType?: string;
    /**
     * @generated from protobuf field: optional resources.filestore.File mug_shot = 11;
     */
    mugShot?: File;
    /**
     * @generated from protobuf field: optional resources.users.CitizenLabels labels = 12;
     */
    labels?: CitizenLabels;
    /**
     * @sanitize: method=StripTags
     *
     * @generated from protobuf field: optional string email = 19;
     */
    email?: string;
}
// @generated message type with reflection information, may provide speed optimized methods
class UserProps$Type extends MessageType<UserProps> {
    constructor() {
        super("resources.users.UserProps", [
            { no: 1, name: "user_id", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gt: 0 } } } },
            { no: 2, name: "updated_at", kind: "message", T: () => Timestamp },
            { no: 3, name: "wanted", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ },
            { no: 4, name: "job_name", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 5, name: "job", kind: "message", T: () => Job },
            { no: 6, name: "job_grade_number", kind: "scalar", opt: true, T: 5 /*ScalarType.INT32*/ },
            { no: 7, name: "job_grade", kind: "message", T: () => JobGrade },
            { no: 8, name: "traffic_infraction_points", kind: "scalar", opt: true, T: 13 /*ScalarType.UINT32*/ },
            { no: 9, name: "open_fines", kind: "scalar", opt: true, T: 3 /*ScalarType.INT64*/, L: 2 /*LongType.NUMBER*/ },
            { no: 10, name: "blood_type", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/ },
            { no: 11, name: "mug_shot", kind: "message", T: () => File },
            { no: 12, name: "labels", kind: "message", T: () => CitizenLabels },
            { no: 19, name: "email", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { minLen: "6", maxLen: "80" } } } }
        ]);
    }
    create(value?: PartialMessage<UserProps>): UserProps {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.userId = 0;
        if (value !== undefined)
            reflectionMergePartial<UserProps>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UserProps): UserProps {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int32 user_id */ 1:
                    message.userId = reader.int32();
                    break;
                case /* optional resources.timestamp.Timestamp updated_at */ 2:
                    message.updatedAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.updatedAt);
                    break;
                case /* optional bool wanted */ 3:
                    message.wanted = reader.bool();
                    break;
                case /* optional string job_name */ 4:
                    message.jobName = reader.string();
                    break;
                case /* optional resources.users.Job job */ 5:
                    message.job = Job.internalBinaryRead(reader, reader.uint32(), options, message.job);
                    break;
                case /* optional int32 job_grade_number */ 6:
                    message.jobGradeNumber = reader.int32();
                    break;
                case /* optional resources.users.JobGrade job_grade */ 7:
                    message.jobGrade = JobGrade.internalBinaryRead(reader, reader.uint32(), options, message.jobGrade);
                    break;
                case /* optional uint32 traffic_infraction_points */ 8:
                    message.trafficInfractionPoints = reader.uint32();
                    break;
                case /* optional int64 open_fines */ 9:
                    message.openFines = reader.int64().toNumber();
                    break;
                case /* optional string blood_type */ 10:
                    message.bloodType = reader.string();
                    break;
                case /* optional resources.filestore.File mug_shot */ 11:
                    message.mugShot = File.internalBinaryRead(reader, reader.uint32(), options, message.mugShot);
                    break;
                case /* optional resources.users.CitizenLabels labels */ 12:
                    message.labels = CitizenLabels.internalBinaryRead(reader, reader.uint32(), options, message.labels);
                    break;
                case /* optional string email */ 19:
                    message.email = reader.string();
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
    internalBinaryWrite(message: UserProps, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* int32 user_id = 1; */
        if (message.userId !== 0)
            writer.tag(1, WireType.Varint).int32(message.userId);
        /* optional resources.timestamp.Timestamp updated_at = 2; */
        if (message.updatedAt)
            Timestamp.internalBinaryWrite(message.updatedAt, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* optional bool wanted = 3; */
        if (message.wanted !== undefined)
            writer.tag(3, WireType.Varint).bool(message.wanted);
        /* optional string job_name = 4; */
        if (message.jobName !== undefined)
            writer.tag(4, WireType.LengthDelimited).string(message.jobName);
        /* optional resources.users.Job job = 5; */
        if (message.job)
            Job.internalBinaryWrite(message.job, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* optional int32 job_grade_number = 6; */
        if (message.jobGradeNumber !== undefined)
            writer.tag(6, WireType.Varint).int32(message.jobGradeNumber);
        /* optional resources.users.JobGrade job_grade = 7; */
        if (message.jobGrade)
            JobGrade.internalBinaryWrite(message.jobGrade, writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        /* optional uint32 traffic_infraction_points = 8; */
        if (message.trafficInfractionPoints !== undefined)
            writer.tag(8, WireType.Varint).uint32(message.trafficInfractionPoints);
        /* optional int64 open_fines = 9; */
        if (message.openFines !== undefined)
            writer.tag(9, WireType.Varint).int64(message.openFines);
        /* optional string blood_type = 10; */
        if (message.bloodType !== undefined)
            writer.tag(10, WireType.LengthDelimited).string(message.bloodType);
        /* optional resources.filestore.File mug_shot = 11; */
        if (message.mugShot)
            File.internalBinaryWrite(message.mugShot, writer.tag(11, WireType.LengthDelimited).fork(), options).join();
        /* optional resources.users.CitizenLabels labels = 12; */
        if (message.labels)
            CitizenLabels.internalBinaryWrite(message.labels, writer.tag(12, WireType.LengthDelimited).fork(), options).join();
        /* optional string email = 19; */
        if (message.email !== undefined)
            writer.tag(19, WireType.LengthDelimited).string(message.email);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.users.UserProps
 */
export const UserProps = new UserProps$Type();
