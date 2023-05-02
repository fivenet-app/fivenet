import * as jspb from 'google-protobuf'

import * as resources_documents_category_pb from '../../resources/documents/category_pb';
import * as resources_jobs_jobs_pb from '../../resources/jobs/jobs_pb';
import * as resources_users_users_pb from '../../resources/users/users_pb';


export class CompleteCitizensRequest extends jspb.Message {
  getSearch(): string;
  setSearch(value: string): CompleteCitizensRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteCitizensRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteCitizensRequest): CompleteCitizensRequest.AsObject;
  static serializeBinaryToWriter(message: CompleteCitizensRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteCitizensRequest;
  static deserializeBinaryFromReader(message: CompleteCitizensRequest, reader: jspb.BinaryReader): CompleteCitizensRequest;
}

export namespace CompleteCitizensRequest {
  export type AsObject = {
    search: string,
  }
}

export class CompleteCitizensRespoonse extends jspb.Message {
  getUsersList(): Array<resources_users_users_pb.UserShort>;
  setUsersList(value: Array<resources_users_users_pb.UserShort>): CompleteCitizensRespoonse;
  clearUsersList(): CompleteCitizensRespoonse;
  addUsers(value?: resources_users_users_pb.UserShort, index?: number): resources_users_users_pb.UserShort;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteCitizensRespoonse.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteCitizensRespoonse): CompleteCitizensRespoonse.AsObject;
  static serializeBinaryToWriter(message: CompleteCitizensRespoonse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteCitizensRespoonse;
  static deserializeBinaryFromReader(message: CompleteCitizensRespoonse, reader: jspb.BinaryReader): CompleteCitizensRespoonse;
}

export namespace CompleteCitizensRespoonse {
  export type AsObject = {
    usersList: Array<resources_users_users_pb.UserShort.AsObject>,
  }
}

export class CompleteJobsRequest extends jspb.Message {
  getSearch(): string;
  setSearch(value: string): CompleteJobsRequest;
  hasSearch(): boolean;
  clearSearch(): CompleteJobsRequest;

  getExactMatch(): boolean;
  setExactMatch(value: boolean): CompleteJobsRequest;

  getCurrentJob(): boolean;
  setCurrentJob(value: boolean): CompleteJobsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteJobsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteJobsRequest): CompleteJobsRequest.AsObject;
  static serializeBinaryToWriter(message: CompleteJobsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteJobsRequest;
  static deserializeBinaryFromReader(message: CompleteJobsRequest, reader: jspb.BinaryReader): CompleteJobsRequest;
}

export namespace CompleteJobsRequest {
  export type AsObject = {
    search?: string,
    exactMatch: boolean,
    currentJob: boolean,
  }

  export enum SearchCase { 
    _SEARCH_NOT_SET = 0,
    SEARCH = 1,
  }
}

export class CompleteJobsResponse extends jspb.Message {
  getJobsList(): Array<resources_jobs_jobs_pb.Job>;
  setJobsList(value: Array<resources_jobs_jobs_pb.Job>): CompleteJobsResponse;
  clearJobsList(): CompleteJobsResponse;
  addJobs(value?: resources_jobs_jobs_pb.Job, index?: number): resources_jobs_jobs_pb.Job;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteJobsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteJobsResponse): CompleteJobsResponse.AsObject;
  static serializeBinaryToWriter(message: CompleteJobsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteJobsResponse;
  static deserializeBinaryFromReader(message: CompleteJobsResponse, reader: jspb.BinaryReader): CompleteJobsResponse;
}

export namespace CompleteJobsResponse {
  export type AsObject = {
    jobsList: Array<resources_jobs_jobs_pb.Job.AsObject>,
  }
}

export class CompleteDocumentCategoriesRequest extends jspb.Message {
  getSearch(): string;
  setSearch(value: string): CompleteDocumentCategoriesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteDocumentCategoriesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteDocumentCategoriesRequest): CompleteDocumentCategoriesRequest.AsObject;
  static serializeBinaryToWriter(message: CompleteDocumentCategoriesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteDocumentCategoriesRequest;
  static deserializeBinaryFromReader(message: CompleteDocumentCategoriesRequest, reader: jspb.BinaryReader): CompleteDocumentCategoriesRequest;
}

export namespace CompleteDocumentCategoriesRequest {
  export type AsObject = {
    search: string,
  }
}

export class CompleteDocumentCategoriesResponse extends jspb.Message {
  getCategoriesList(): Array<resources_documents_category_pb.DocumentCategory>;
  setCategoriesList(value: Array<resources_documents_category_pb.DocumentCategory>): CompleteDocumentCategoriesResponse;
  clearCategoriesList(): CompleteDocumentCategoriesResponse;
  addCategories(value?: resources_documents_category_pb.DocumentCategory, index?: number): resources_documents_category_pb.DocumentCategory;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteDocumentCategoriesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteDocumentCategoriesResponse): CompleteDocumentCategoriesResponse.AsObject;
  static serializeBinaryToWriter(message: CompleteDocumentCategoriesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteDocumentCategoriesResponse;
  static deserializeBinaryFromReader(message: CompleteDocumentCategoriesResponse, reader: jspb.BinaryReader): CompleteDocumentCategoriesResponse;
}

export namespace CompleteDocumentCategoriesResponse {
  export type AsObject = {
    categoriesList: Array<resources_documents_category_pb.DocumentCategory.AsObject>,
  }
}

