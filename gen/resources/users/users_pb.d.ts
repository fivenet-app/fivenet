import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';


export class UserJob extends jspb.Message {
  getJob(): string;
  setJob(value: string): UserJob;

  getGrade(): number;
  setGrade(value: number): UserJob;

  getJobLabel(): string;
  setJobLabel(value: string): UserJob;

  getGradeLabel(): string;
  setGradeLabel(value: string): UserJob;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserJob.AsObject;
  static toObject(includeInstance: boolean, msg: UserJob): UserJob.AsObject;
  static serializeBinaryToWriter(message: UserJob, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserJob;
  static deserializeBinaryFromReader(message: UserJob, reader: jspb.BinaryReader): UserJob;
}

export namespace UserJob {
  export type AsObject = {
    job: string,
    grade: number,
    jobLabel: string,
    gradeLabel: string,
  }
}

export class UserShort extends jspb.Message {
  getUserId(): number;
  setUserId(value: number): UserShort;

  getIdentifier(): string;
  setIdentifier(value: string): UserShort;

  getJob(): string;
  setJob(value: string): UserShort;

  getJobGrade(): number;
  setJobGrade(value: number): UserShort;

  getFirstname(): string;
  setFirstname(value: string): UserShort;

  getLastname(): string;
  setLastname(value: string): UserShort;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserShort.AsObject;
  static toObject(includeInstance: boolean, msg: UserShort): UserShort.AsObject;
  static serializeBinaryToWriter(message: UserShort, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserShort;
  static deserializeBinaryFromReader(message: UserShort, reader: jspb.BinaryReader): UserShort;
}

export namespace UserShort {
  export type AsObject = {
    userId: number,
    identifier: string,
    job: string,
    jobGrade: number,
    firstname: string,
    lastname: string,
  }
}

export class User extends jspb.Message {
  getUserId(): number;
  setUserId(value: number): User;

  getIdentifier(): string;
  setIdentifier(value: string): User;

  getJob(): string;
  setJob(value: string): User;

  getJobGrade(): number;
  setJobGrade(value: number): User;

  getFirstname(): string;
  setFirstname(value: string): User;

  getLastname(): string;
  setLastname(value: string): User;

  getDateofbirth(): string;
  setDateofbirth(value: string): User;

  getSex(): string;
  setSex(value: string): User;

  getHeight(): string;
  setHeight(value: string): User;

  getPhonenumber(): string;
  setPhonenumber(value: string): User;

  getVisum(): number;
  setVisum(value: number): User;

  getPlaytime(): number;
  setPlaytime(value: number): User;

  getProps(): UserProps | undefined;
  setProps(value?: UserProps): User;
  hasProps(): boolean;
  clearProps(): User;

  getLicensesList(): Array<License>;
  setLicensesList(value: Array<License>): User;
  clearLicensesList(): User;
  addLicenses(value?: License, index?: number): License;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
  export type AsObject = {
    userId: number,
    identifier: string,
    job: string,
    jobGrade: number,
    firstname: string,
    lastname: string,
    dateofbirth: string,
    sex: string,
    height: string,
    phonenumber: string,
    visum: number,
    playtime: number,
    props?: UserProps.AsObject,
    licensesList: Array<License.AsObject>,
  }
}

export class License extends jspb.Message {
  getType(): string;
  setType(value: string): License;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): License.AsObject;
  static toObject(includeInstance: boolean, msg: License): License.AsObject;
  static serializeBinaryToWriter(message: License, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): License;
  static deserializeBinaryFromReader(message: License, reader: jspb.BinaryReader): License;
}

export namespace License {
  export type AsObject = {
    type: string,
  }
}

export class UserProps extends jspb.Message {
  getWanted(): boolean;
  setWanted(value: boolean): UserProps;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserProps.AsObject;
  static toObject(includeInstance: boolean, msg: UserProps): UserProps.AsObject;
  static serializeBinaryToWriter(message: UserProps, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserProps;
  static deserializeBinaryFromReader(message: UserProps, reader: jspb.BinaryReader): UserProps;
}

export namespace UserProps {
  export type AsObject = {
    wanted: boolean,
  }
}

export class UserActivity extends jspb.Message {
  getId(): number;
  setId(value: number): UserActivity;

  getType(): USER_ACTIVITY_TYPE;
  setType(value: USER_ACTIVITY_TYPE): UserActivity;

  getCreatedAt(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: resources_timestamp_timestamp_pb.Timestamp): UserActivity;
  hasCreatedAt(): boolean;
  clearCreatedAt(): UserActivity;

  getTargetuser(): UserShort | undefined;
  setTargetuser(value?: UserShort): UserActivity;
  hasTargetuser(): boolean;
  clearTargetuser(): UserActivity;

  getCauseuser(): UserShort | undefined;
  setCauseuser(value?: UserShort): UserActivity;
  hasCauseuser(): boolean;
  clearCauseuser(): UserActivity;

  getKey(): string;
  setKey(value: string): UserActivity;

  getOldvalue(): string;
  setOldvalue(value: string): UserActivity;

  getNewvalue(): string;
  setNewvalue(value: string): UserActivity;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserActivity.AsObject;
  static toObject(includeInstance: boolean, msg: UserActivity): UserActivity.AsObject;
  static serializeBinaryToWriter(message: UserActivity, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserActivity;
  static deserializeBinaryFromReader(message: UserActivity, reader: jspb.BinaryReader): UserActivity;
}

export namespace UserActivity {
  export type AsObject = {
    id: number,
    type: USER_ACTIVITY_TYPE,
    createdAt?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    targetuser?: UserShort.AsObject,
    causeuser?: UserShort.AsObject,
    key: string,
    oldvalue: string,
    newvalue: string,
  }
}

export enum USER_ACTIVITY_TYPE { 
  CHANGED = 0,
  MENTIONED = 1,
  CREATED = 2,
}
