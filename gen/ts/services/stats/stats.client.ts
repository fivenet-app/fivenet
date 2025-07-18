// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "services/stats/stats.proto" (package "services.stats", syntax proto3)
// tslint:disable
// @ts-nocheck
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { StatsService } from "./stats";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { GetStatsResponse } from "./stats";
import type { GetStatsRequest } from "./stats";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service services.stats.StatsService
 */
export interface IStatsServiceClient {
    /**
     * @generated from protobuf rpc: GetStats
     */
    getStats(input: GetStatsRequest, options?: RpcOptions): UnaryCall<GetStatsRequest, GetStatsResponse>;
}
/**
 * @generated from protobuf service services.stats.StatsService
 */
export class StatsServiceClient implements IStatsServiceClient, ServiceInfo {
    typeName = StatsService.typeName;
    methods = StatsService.methods;
    options = StatsService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @generated from protobuf rpc: GetStats
     */
    getStats(input: GetStatsRequest, options?: RpcOptions): UnaryCall<GetStatsRequest, GetStatsResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetStatsRequest, GetStatsResponse>("unary", this._transport, method, opt, input);
    }
}
