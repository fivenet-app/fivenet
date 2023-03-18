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
    grade: number,
    label: string,
  }
}

