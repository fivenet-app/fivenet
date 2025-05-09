// @generated by protobuf-ts 2.10.0 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/rector/cron.proto" (package "services.rector", syntax proto3)
// @ts-nocheck
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { RectorCronService } from "./cron";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { ListCronjobsResponse } from "./cron";
import type { ListCronjobsRequest } from "./cron";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service services.rector.RectorCronService
 */
export interface IRectorCronServiceClient {
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: ListCronjobs(services.rector.ListCronjobsRequest) returns (services.rector.ListCronjobsResponse);
     */
    listCronjobs(input: ListCronjobsRequest, options?: RpcOptions): UnaryCall<ListCronjobsRequest, ListCronjobsResponse>;
}
/**
 * @generated from protobuf service services.rector.RectorCronService
 */
export class RectorCronServiceClient implements IRectorCronServiceClient, ServiceInfo {
    typeName = RectorCronService.typeName;
    methods = RectorCronService.methods;
    options = RectorCronService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: ListCronjobs(services.rector.ListCronjobsRequest) returns (services.rector.ListCronjobsResponse);
     */
    listCronjobs(input: ListCronjobsRequest, options?: RpcOptions): UnaryCall<ListCronjobsRequest, ListCronjobsResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<ListCronjobsRequest, ListCronjobsResponse>("unary", this._transport, method, opt, input);
    }
}
