import * as jspb from 'google-protobuf'

import * as resources_timestamp_timestamp_pb from '../../resources/timestamp/timestamp_pb';


export class User extends jspb.Message {
  getUserid(): number;
  setUserid(value: number): User;

  getIdentifier(): string;
  setIdentifier(value: string): User;

  getJob(): string;
  setJob(value: string): User;

  getJobgrade(): number;
  setJobgrade(value: number): User;

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

  getVisum(): number;
  setVisum(value: number): User;

  getPlaytime(): number;
  setPlaytime(value: number): User;

  getProps(): Props | undefined;
  setProps(value?: Props): User;
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
    userid: number,
    identifier: string,
    job: string,
    jobgrade: number,
    firstname: string,
    lastname: string,
    dateofbirth: string,
    sex: string,
    height: string,
    visum: number,
    playtime: number,
    props?: Props.AsObject,
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

export class Props extends jspb.Message {
  getWanted(): boolean;
  setWanted(value: boolean): Props;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Props.AsObject;
  static toObject(includeInstance: boolean, msg: Props): Props.AsObject;
  static serializeBinaryToWriter(message: Props, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Props;
  static deserializeBinaryFromReader(message: Props, reader: jspb.BinaryReader): Props;
}

export namespace Props {
  export type AsObject = {
    wanted: boolean,
  }
}

export class ShortUser extends jspb.Message {
  getUserid(): number;
  setUserid(value: number): ShortUser;

  getIdentifier(): string;
  setIdentifier(value: string): ShortUser;

  getJob(): string;
  setJob(value: string): ShortUser;

  getJobgrade(): number;
  setJobgrade(value: number): ShortUser;

  getFirstname(): string;
  setFirstname(value: string): ShortUser;

  getLastname(): string;
  setLastname(value: string): ShortUser;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ShortUser.AsObject;
  static toObject(includeInstance: boolean, msg: ShortUser): ShortUser.AsObject;
  static serializeBinaryToWriter(message: ShortUser, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ShortUser;
  static deserializeBinaryFromReader(message: ShortUser, reader: jspb.BinaryReader): ShortUser;
}

export namespace ShortUser {
  export type AsObject = {
    userid: number,
    identifier: string,
    job: string,
    jobgrade: number,
    firstname: string,
    lastname: string,
  }
}

export class UserActivity extends jspb.Message {
  getId(): number;
  setId(value: number): UserActivity;

  getType(): string;
  setType(value: string): UserActivity;

  getCreatedat(): resources_timestamp_timestamp_pb.Timestamp | undefined;
  setCreatedat(value?: resources_timestamp_timestamp_pb.Timestamp): UserActivity;
  hasCreatedat(): boolean;
  clearCreatedat(): UserActivity;

  getTargetuser(): ShortUser | undefined;
  setTargetuser(value?: ShortUser): UserActivity;
  hasTargetuser(): boolean;
  clearTargetuser(): UserActivity;

  getCauseuser(): ShortUser | undefined;
  setCauseuser(value?: ShortUser): UserActivity;
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
    type: string,
    createdat?: resources_timestamp_timestamp_pb.Timestamp.AsObject,
    targetuser?: ShortUser.AsObject,
    causeuser?: ShortUser.AsObject,
    key: string,
    oldvalue: string,
    newvalue: string,
  }
}

