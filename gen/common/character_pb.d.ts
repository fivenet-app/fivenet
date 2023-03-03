import * as jspb from 'google-protobuf'



export class Character extends jspb.Message {
  getId(): number;
  setId(value: number): Character;

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

  getProps(): Props | undefined;
  setProps(value?: Props): Character;
  hasProps(): boolean;
  clearProps(): Character;

  getLicensesList(): Array<License>;
  setLicensesList(value: Array<License>): Character;
  clearLicensesList(): Character;
  addLicenses(value?: License, index?: number): License;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Character.AsObject;
  static toObject(includeInstance: boolean, msg: Character): Character.AsObject;
  static serializeBinaryToWriter(message: Character, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Character;
  static deserializeBinaryFromReader(message: Character, reader: jspb.BinaryReader): Character;
}

export namespace Character {
  export type AsObject = {
    id: number,
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
  getName(): string;
  setName(value: string): License;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): License.AsObject;
  static toObject(includeInstance: boolean, msg: License): License.AsObject;
  static serializeBinaryToWriter(message: License, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): License;
  static deserializeBinaryFromReader(message: License, reader: jspb.BinaryReader): License;
}

export namespace License {
  export type AsObject = {
    name: string,
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

