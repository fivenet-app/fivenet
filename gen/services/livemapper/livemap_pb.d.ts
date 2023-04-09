import * as jspb from 'google-protobuf'

import * as resources_jobs_jobs_pb from '../../resources/jobs/jobs_pb';
import * as resources_livemap_livemap_pb from '../../resources/livemap/livemap_pb';


export class StreamRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StreamRequest): StreamRequest.AsObject;
  static serializeBinaryToWriter(message: StreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StreamRequest;
  static deserializeBinaryFromReader(message: StreamRequest, reader: jspb.BinaryReader): StreamRequest;
}

export namespace StreamRequest {
  export type AsObject = {
  }
}

export class StreamResponse extends jspb.Message {
  getJobsList(): Array<resources_jobs_jobs_pb.Job>;
  setJobsList(value: Array<resources_jobs_jobs_pb.Job>): StreamResponse;
  clearJobsList(): StreamResponse;
  addJobs(value?: resources_jobs_jobs_pb.Job, index?: number): resources_jobs_jobs_pb.Job;

  getDispatchesList(): Array<resources_livemap_livemap_pb.DispatchMarker>;
  setDispatchesList(value: Array<resources_livemap_livemap_pb.DispatchMarker>): StreamResponse;
  clearDispatchesList(): StreamResponse;
  addDispatches(value?: resources_livemap_livemap_pb.DispatchMarker, index?: number): resources_livemap_livemap_pb.DispatchMarker;

  getUsersList(): Array<resources_livemap_livemap_pb.UserMarker>;
  setUsersList(value: Array<resources_livemap_livemap_pb.UserMarker>): StreamResponse;
  clearUsersList(): StreamResponse;
  addUsers(value?: resources_livemap_livemap_pb.UserMarker, index?: number): resources_livemap_livemap_pb.UserMarker;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StreamResponse): StreamResponse.AsObject;
  static serializeBinaryToWriter(message: StreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StreamResponse;
  static deserializeBinaryFromReader(message: StreamResponse, reader: jspb.BinaryReader): StreamResponse;
}

export namespace StreamResponse {
  export type AsObject = {
    jobsList: Array<resources_jobs_jobs_pb.Job.AsObject>,
    dispatchesList: Array<resources_livemap_livemap_pb.DispatchMarker.AsObject>,
    usersList: Array<resources_livemap_livemap_pb.UserMarker.AsObject>,
  }
}

