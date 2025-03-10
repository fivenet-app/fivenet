// @generated by protobuf-ts 2.9.5 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/internet/ads.proto" (package "services.internet", syntax proto3)
// @ts-nocheck
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { AdsService } from "./ads";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { GetAdsResponse } from "./ads";
import type { GetAdsRequest } from "./ads";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service services.internet.AdsService
 */
export interface IAdsServiceClient {
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: GetAds(services.internet.GetAdsRequest) returns (services.internet.GetAdsResponse);
     */
    getAds(input: GetAdsRequest, options?: RpcOptions): UnaryCall<GetAdsRequest, GetAdsResponse>;
}
/**
 * @generated from protobuf service services.internet.AdsService
 */
export class AdsServiceClient implements IAdsServiceClient, ServiceInfo {
    typeName = AdsService.typeName;
    methods = AdsService.methods;
    options = AdsService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: GetAds(services.internet.GetAdsRequest) returns (services.internet.GetAdsResponse);
     */
    getAds(input: GetAdsRequest, options?: RpcOptions): UnaryCall<GetAdsRequest, GetAdsResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetAdsRequest, GetAdsResponse>("unary", this._transport, method, opt, input);
    }
}
