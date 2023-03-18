import * as jspb from 'google-protobuf'

import * as resources_documents_documents_pb from '../../resources/documents/documents_pb';
import * as resources_jobs_jobs_pb from '../../resources/jobs/jobs_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';


export class CompleteCharNamesRequest extends jspb.Message {
  getSearch(): string;
  setSearch(value: string): CompleteCharNamesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteCharNamesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteCharNamesRequest): CompleteCharNamesRequest.AsObject;
  static serializeBinaryToWriter(message: CompleteCharNamesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteCharNamesRequest;
  static deserializeBinaryFromReader(message: CompleteCharNamesRequest, reader: jspb.BinaryReader): CompleteCharNamesRequest;
}

export namespace CompleteCharNamesRequest {
  export type AsObject = {
    search: string,
  }
}

export class CompleteCharNamesRespoonse extends jspb.Message {
  getUsersList(): Array<resources_users_users_pb.UserShort>;
  setUsersList(value: Array<resources_users_users_pb.UserShort>): CompleteCharNamesRespoonse;
  clearUsersList(): CompleteCharNamesRespoonse;
  addUsers(value?: resources_users_users_pb.UserShort, index?: number): resources_users_users_pb.UserShort;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteCharNamesRespoonse.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteCharNamesRespoonse): CompleteCharNamesRespoonse.AsObject;
  static serializeBinaryToWriter(message: CompleteCharNamesRespoonse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteCharNamesRespoonse;
  static deserializeBinaryFromReader(message: CompleteCharNamesRespoonse, reader: jspb.BinaryReader): CompleteCharNamesRespoonse;
}

export namespace CompleteCharNamesRespoonse {
  export type AsObject = {
    usersList: Array<resources_users_users_pb.UserShort.AsObject>,
  }
}

export class CompleteJobNamesRequest extends jspb.Message {
  getSearch(): string;
  setSearch(value: string): CompleteJobNamesRequest;

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
  }
}

export class CompleteJobNamesResponse extends jspb.Message {
  getJobsList(): Array<resources_jobs_jobs_pb.Job>;
  setJobsList(value: Array<resources_jobs_jobs_pb.Job>): CompleteJobNamesResponse;
  clearJobsList(): CompleteJobNamesResponse;
  addJobs(value?: resources_jobs_jobs_pb.Job, index?: number): resources_jobs_jobs_pb.Job;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteJobNamesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteJobNamesResponse): CompleteJobNamesResponse.AsObject;
  static serializeBinaryToWriter(message: CompleteJobNamesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteJobNamesResponse;
  static deserializeBinaryFromReader(message: CompleteJobNamesResponse, reader: jspb.BinaryReader): CompleteJobNamesResponse;
}

export namespace CompleteJobNamesResponse {
  export type AsObject = {
    jobsList: Array<resources_jobs_jobs_pb.Job.AsObject>,
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
  getGradesList(): Array<resources_jobs_jobs_pb.JobGrade>;
  setGradesList(value: Array<resources_jobs_jobs_pb.JobGrade>): CompleteJobGradesResponse;
  clearGradesList(): CompleteJobGradesResponse;
  addGrades(value?: resources_jobs_jobs_pb.JobGrade, index?: number): resources_jobs_jobs_pb.JobGrade;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteJobGradesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteJobGradesResponse): CompleteJobGradesResponse.AsObject;
  static serializeBinaryToWriter(message: CompleteJobGradesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteJobGradesResponse;
  static deserializeBinaryFromReader(message: CompleteJobGradesResponse, reader: jspb.BinaryReader): CompleteJobGradesResponse;
}

export namespace CompleteJobGradesResponse {
  export type AsObject = {
    gradesList: Array<resources_jobs_jobs_pb.JobGrade.AsObject>,
  }
}

export class CompleteDocumentCategoryRequest extends jspb.Message {
  getSearch(): string;
  setSearch(value: string): CompleteDocumentCategoryRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteDocumentCategoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteDocumentCategoryRequest): CompleteDocumentCategoryRequest.AsObject;
  static serializeBinaryToWriter(message: CompleteDocumentCategoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteDocumentCategoryRequest;
  static deserializeBinaryFromReader(message: CompleteDocumentCategoryRequest, reader: jspb.BinaryReader): CompleteDocumentCategoryRequest;
}

export namespace CompleteDocumentCategoryRequest {
  export type AsObject = {
    search: string,
  }
}

export class CompleteDocumentCategoryResponse extends jspb.Message {
  getCategoriesList(): Array<resources_documents_documents_pb.DocumentCategory>;
  setCategoriesList(value: Array<resources_documents_documents_pb.DocumentCategory>): CompleteDocumentCategoryResponse;
  clearCategoriesList(): CompleteDocumentCategoryResponse;
  addCategories(value?: resources_documents_documents_pb.DocumentCategory, index?: number): resources_documents_documents_pb.DocumentCategory;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteDocumentCategoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteDocumentCategoryResponse): CompleteDocumentCategoryResponse.AsObject;
  static serializeBinaryToWriter(message: CompleteDocumentCategoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteDocumentCategoryResponse;
  static deserializeBinaryFromReader(message: CompleteDocumentCategoryResponse, reader: jspb.BinaryReader): CompleteDocumentCategoryResponse;
}

export namespace CompleteDocumentCategoryResponse {
  export type AsObject = {
    categoriesList: Array<resources_documents_documents_pb.DocumentCategory.AsObject>,
  }
}

