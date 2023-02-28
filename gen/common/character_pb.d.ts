import * as jspb from 'google-protobuf'



export class Character extends jspb.Message {
  getIdentifier(): string;
  setIdentifier(value: string): Character;

  getJob(): string;
  setJob(value: string): Character;

  getJobgrade(): number;
  setJobgrade(value: number): Character;

  getFirstname(): string;
  setFirstname(value: string): Character;

  getLastname(): string;
  setLastname(value: string): Character;

  getDateofbirth(): string;
  setDateofbirth(value: string): Character;

  getSex(): string;
  setSex(value: string): Character;

  getHeight(): string;
  setHeight(value: string): Character;

  getVisum(): number;
  setVisum(value: number): Character;

  getPlaytime(): number;
  setPlaytime(value: number): Character;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Character.AsObject;
  static toObject(includeInstance: boolean, msg: Character): Character.AsObject;
  static serializeBinaryToWriter(message: Character, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Character;
  static deserializeBinaryFromReader(message: Character, reader: jspb.BinaryReader): Character;
}

export namespace Character {
  export type AsObject = {
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
  }
}

