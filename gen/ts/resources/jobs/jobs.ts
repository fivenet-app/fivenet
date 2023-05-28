// @generated by protobuf-ts 2.9.0 with parameter optimize_code_size,long_type_number
// @generated from protobuf file "resources/jobs/jobs.proto" (package "resources.jobs", syntax proto3)
// tslint:disable
import { MessageType } from "@protobuf-ts/runtime";
/**
 * @generated from protobuf message resources.jobs.Job
 */
export interface Job {
    /**
     * @generated from protobuf field: string name = 1;
     */
    name: string; // @gotags: sql:"primary_key" alias:"name"
    /**
     * @generated from protobuf field: string label = 2;
     */
    label: string; // @gotags: alias:"label"
    /**
     * @generated from protobuf field: repeated resources.jobs.JobGrade grades = 3;
     */
    grades: JobGrade[];
}
/**
 * @generated from protobuf message resources.jobs.JobGrade
 */
export interface JobGrade {
    /**
     * @generated from protobuf field: optional string job_name = 1;
     */
    jobName?: string; // @gotags: alias:"job_name"
    /**
     * @generated from protobuf field: int32 grade = 2;
     */
    grade: number; // @gotags: alias:"grade"
    /**
     * @generated from protobuf field: string label = 3;
     */
    label: string; // @gotags: alias:"label"
}
/**
 * @generated from protobuf message resources.jobs.JobProps
 */
export interface JobProps {
    /**
     * @generated from protobuf field: string job = 1;
     */
    job: string; // @gotags: alias:"job"
    /**
     * @generated from protobuf field: string theme = 2;
     */
    theme: string; // @gotags: alias:"theme"
    /**
     * @generated from protobuf field: string livemap_marker_color = 3;
     */
    livemapMarkerColor: string; // @gotags: alias:"livemap_marker_color"
    /**
     * @generated from protobuf field: string quick_buttons = 4;
     */
    quickButtons: string; // @gotags: alias:"quick_buttons"
}
// @generated message type with reflection information, may provide speed optimized methods
class Job$Type extends MessageType<Job> {
    constructor() {
        super("resources.jobs.Job", [
            { no: 1, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 2, name: "label", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 3, name: "grades", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => JobGrade }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.jobs.Job
 */
export const Job = new Job$Type();
// @generated message type with reflection information, may provide speed optimized methods
class JobGrade$Type extends MessageType<JobGrade> {
    constructor() {
        super("resources.jobs.JobGrade", [
            { no: 1, name: "job_name", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } },
            { no: 2, name: "grade", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gt: 0 } } } },
            { no: 3, name: "label", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "50" } } } }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.jobs.JobGrade
 */
export const JobGrade = new JobGrade$Type();
// @generated message type with reflection information, may provide speed optimized methods
class JobProps$Type extends MessageType<JobProps> {
    constructor() {
        super("resources.jobs.JobProps", [
            { no: 1, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 2, name: "theme", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 3, name: "livemap_marker_color", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 4, name: "quick_buttons", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message resources.jobs.JobProps
 */
export const JobProps = new JobProps$Type();
