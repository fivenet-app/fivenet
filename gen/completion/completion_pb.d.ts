import * as jspb from 'google-protobuf'



export class Job extends jspb.Message {
  getName(): string;
  setName(value: string): Job;

  getLabel(): string;
  setLabel(value: string): Job;

  getGradesList(): Array<JobGrade>;
  setGradesList(value: Array<JobGrade>): Job;
  clearGradesList(): Job;
  addGrades(value?: JobGrade, index?: number): JobGrade;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Job.AsObject;
  static toObject(includeInstance: boolean, msg: Job): Job.AsObject;
  static serializeBinaryToWriter(message: Job, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Job;
  static deserializeBinaryFromReader(message: Job, reader: jspb.BinaryReader): Job;
}

export namespace Job {
  export type AsObject = {
    name: string,
    label: string,
    gradesList: Array<JobGrade.AsObject>,
  }
}

export class JobGrade extends jspb.Message {
  getJob(): string;
  setJob(value: string): JobGrade;

  getGrade(): number;
  setGrade(value: number): JobGrade;

  getLabel(): string;
  setLabel(value: string): JobGrade;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JobGrade.AsObject;
  static toObject(includeInstance: boolean, msg: JobGrade): JobGrade.AsObject;
  static serializeBinaryToWriter(message: JobGrade, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JobGrade;
  static deserializeBinaryFromReader(message: JobGrade, reader: jspb.BinaryReader): JobGrade;
}

export namespace JobGrade {
  export type AsObject = {
    job: string,
    grade: number,
    label: string,
  }
}

export class CompleteJobNamesRequest extends jspb.Message {
  getSearch(): string;
  setSearch(value: string): CompleteJobNamesRequest;

  getWithgrades(): boolean;
  setWithgrades(value: boolean): CompleteJobNamesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteJobNamesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteJobNamesRequest): CompleteJobNamesRequest.AsObject;
  static serializeBinaryToWriter(message: CompleteJobNamesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteJobNamesRequest;
  static deserializeBinaryFromReader(message: CompleteJobNamesRequest, reader: jspb.BinaryReader): CompleteJobNamesRequest;
}

export namespace CompleteJobNamesRequest {
  export type AsObject = {
    search: string,
    withgrades: boolean,
  }
}

export class CompleteJobNamesResponse extends jspb.Message {
  getJobsList(): Array<Job>;
  setJobsList(value: Array<Job>): CompleteJobNamesResponse;
  clearJobsList(): CompleteJobNamesResponse;
  addJobs(value?: Job, index?: number): Job;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteJobNamesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteJobNamesResponse): CompleteJobNamesResponse.AsObject;
  static serializeBinaryToWriter(message: CompleteJobNamesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteJobNamesResponse;
  static deserializeBinaryFromReader(message: CompleteJobNamesResponse, reader: jspb.BinaryReader): CompleteJobNamesResponse;
}

export namespace CompleteJobNamesResponse {
  export type AsObject = {
    jobsList: Array<Job.AsObject>,
  }
}

export class CompleteJobGradesRequest extends jspb.Message {
  getJob(): string;
  setJob(value: string): CompleteJobGradesRequest;

  getSearch(): string;
  setSearch(value: string): CompleteJobGradesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteJobGradesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteJobGradesRequest): CompleteJobGradesRequest.AsObject;
  static serializeBinaryToWriter(message: CompleteJobGradesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteJobGradesRequest;
  static deserializeBinaryFromReader(message: CompleteJobGradesRequest, reader: jspb.BinaryReader): CompleteJobGradesRequest;
}

export namespace CompleteJobGradesRequest {
  export type AsObject = {
    job: string,
    search: string,
  }
}

export class CompleteJobGradesResponse extends jspb.Message {
  getGradesList(): Array<JobGrade>;
  setGradesList(value: Array<JobGrade>): CompleteJobGradesResponse;
  clearGradesList(): CompleteJobGradesResponse;
  addGrades(value?: JobGrade, index?: number): JobGrade;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteJobGradesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteJobGradesResponse): CompleteJobGradesResponse.AsObject;
  static serializeBinaryToWriter(message: CompleteJobGradesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteJobGradesResponse;
  static deserializeBinaryFromReader(message: CompleteJobGradesResponse, reader: jspb.BinaryReader): CompleteJobGradesResponse;
}

export namespace CompleteJobGradesResponse {
  export type AsObject = {
    gradesList: Array<JobGrade.AsObject>,
  }
}

