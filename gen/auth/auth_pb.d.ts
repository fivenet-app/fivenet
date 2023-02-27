import * as jspb from 'google-protobuf'



export class LoginRequest extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): LoginRequest;

  getPassword(): string;
  setPassword(value: string): LoginRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LoginRequest): LoginRequest.AsObject;
  static serializeBinaryToWriter(message: LoginRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginRequest;
  static deserializeBinaryFromReader(message: LoginRequest, reader: jspb.BinaryReader): LoginRequest;
}

export namespace LoginRequest {
  export type AsObject = {
    username: string,
    password: string,
  }
}

export class LoginResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): LoginResponse;

  getCharsList(): Array<Character>;
  setCharsList(value: Array<Character>): LoginResponse;
  clearCharsList(): LoginResponse;
  addChars(value?: Character, index?: number): Character;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginResponse.AsObject;
  static toObject(includeInstance: boolean, msg: LoginResponse): LoginResponse.AsObject;
  static serializeBinaryToWriter(message: LoginResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginResponse;
  static deserializeBinaryFromReader(message: LoginResponse, reader: jspb.BinaryReader): LoginResponse;
}

export namespace LoginResponse {
  export type AsObject = {
    token: string,
    charsList: Array<Character.AsObject>,
  }
}

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

export class ChooseCharacterRequest extends jspb.Message {
  getToken(): string;
  setToken(value: string): ChooseCharacterRequest;

  getIndex(): number;
  setIndex(value: number): ChooseCharacterRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChooseCharacterRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ChooseCharacterRequest): ChooseCharacterRequest.AsObject;
  static serializeBinaryToWriter(message: ChooseCharacterRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChooseCharacterRequest;
  static deserializeBinaryFromReader(message: ChooseCharacterRequest, reader: jspb.BinaryReader): ChooseCharacterRequest;
}

export namespace ChooseCharacterRequest {
  export type AsObject = {
    token: string,
    index: number,
  }
}

export class ChooseCharacterResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): ChooseCharacterResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChooseCharacterResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ChooseCharacterResponse): ChooseCharacterResponse.AsObject;
  static serializeBinaryToWriter(message: ChooseCharacterResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChooseCharacterResponse;
  static deserializeBinaryFromReader(message: ChooseCharacterResponse, reader: jspb.BinaryReader): ChooseCharacterResponse;
}

export namespace ChooseCharacterResponse {
  export type AsObject = {
    token: string,
  }
}

export class LogoutRequest extends jspb.Message {
  getToken(): string;
  setToken(value: string): LogoutRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogoutRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LogoutRequest): LogoutRequest.AsObject;
  static serializeBinaryToWriter(message: LogoutRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogoutRequest;
  static deserializeBinaryFromReader(message: LogoutRequest, reader: jspb.BinaryReader): LogoutRequest;
}

export namespace LogoutRequest {
  export type AsObject = {
    token: string,
  }
}

export class LogoutResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogoutResponse.AsObject;
  static toObject(includeInstance: boolean, msg: LogoutResponse): LogoutResponse.AsObject;
  static serializeBinaryToWriter(message: LogoutResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogoutResponse;
  static deserializeBinaryFromReader(message: LogoutResponse, reader: jspb.BinaryReader): LogoutResponse;
}

export namespace LogoutResponse {
  export type AsObject = {
  }
}

