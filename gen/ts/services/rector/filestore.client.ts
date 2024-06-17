// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/rector/filestore.proto" (package "services.rector", syntax proto3)
// @ts-nocheck
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { RectorFilestoreService } from "./filestore";
import type { DeleteFileResponse } from "./filestore";
import type { DeleteFileRequest } from "./filestore";
import type { UploadFileResponse } from "./filestore";
import type { UploadFileRequest } from "./filestore";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { ListFilesResponse } from "./filestore";
import type { ListFilesRequest } from "./filestore";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service services.rector.RectorFilestoreService
 */
export interface IRectorFilestoreServiceClient {
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: ListFiles(services.rector.ListFilesRequest) returns (services.rector.ListFilesResponse);
     */
    listFiles(input: ListFilesRequest, options?: RpcOptions): UnaryCall<ListFilesRequest, ListFilesResponse>;
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: UploadFile(services.rector.UploadFileRequest) returns (services.rector.UploadFileResponse);
     */
    uploadFile(input: UploadFileRequest, options?: RpcOptions): UnaryCall<UploadFileRequest, UploadFileResponse>;
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: DeleteFile(services.rector.DeleteFileRequest) returns (services.rector.DeleteFileResponse);
     */
    deleteFile(input: DeleteFileRequest, options?: RpcOptions): UnaryCall<DeleteFileRequest, DeleteFileResponse>;
}
/**
 * @generated from protobuf service services.rector.RectorFilestoreService
 */
export class RectorFilestoreServiceClient implements IRectorFilestoreServiceClient, ServiceInfo {
    typeName = RectorFilestoreService.typeName;
    methods = RectorFilestoreService.methods;
    options = RectorFilestoreService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: ListFiles(services.rector.ListFilesRequest) returns (services.rector.ListFilesResponse);
     */
    listFiles(input: ListFilesRequest, options?: RpcOptions): UnaryCall<ListFilesRequest, ListFilesResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<ListFilesRequest, ListFilesResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: UploadFile(services.rector.UploadFileRequest) returns (services.rector.UploadFileResponse);
     */
    uploadFile(input: UploadFileRequest, options?: RpcOptions): UnaryCall<UploadFileRequest, UploadFileResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<UploadFileRequest, UploadFileResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: DeleteFile(services.rector.DeleteFileRequest) returns (services.rector.DeleteFileResponse);
     */
    deleteFile(input: DeleteFileRequest, options?: RpcOptions): UnaryCall<DeleteFileRequest, DeleteFileResponse> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<DeleteFileRequest, DeleteFileResponse>("unary", this._transport, method, opt, input);
    }
}
