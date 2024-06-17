// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/rector/config.proto" (package "services.rector", syntax proto3)
// @ts-nocheck
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { RectorConfigService } from "./config";
import type { UpdateAppConfigResponse } from "./config";
import type { UpdateAppConfigRequest } from "./config";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { GetAppConfigResponse } from "./config";
import type { GetAppConfigRequest } from "./config";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service services.rector.RectorConfigService
 */
export interface IRectorConfigServiceClient {
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: GetAppConfig(services.rector.GetAppConfigRequest) returns (services.rector.GetAppConfigResponse);
     */
    getAppConfig(input: GetAppConfigRequest, options?: RpcOptions): UnaryCall<GetAppConfigRequest, GetAppConfigResponse>;
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: UpdateAppConfig(services.rector.UpdateAppConfigRequest) returns (services.rector.UpdateAppConfigResponse);
     */
    updateAppConfig(input: UpdateAppConfigRequest, options?: RpcOptions): UnaryCall<UpdateAppConfigRequest, UpdateAppConfigResponse>;
}
/**
 * @generated from protobuf service services.rector.RectorConfigService
 */
export class RectorConfigServiceClient implements IRectorConfigServiceClient, ServiceInfo {
    typeName = RectorConfigService.typeName;
    methods = RectorConfigService.methods;
    options = RectorConfigService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: GetAppConfig(services.rector.GetAppConfigRequest) returns (services.rector.GetAppConfigResponse);
     */
    getAppConfig(input: GetAppConfigRequest, options?: RpcOptions): UnaryCall<GetAppConfigRequest, GetAppConfigResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetAppConfigRequest, GetAppConfigResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: UpdateAppConfig(services.rector.UpdateAppConfigRequest) returns (services.rector.UpdateAppConfigResponse);
     */
    updateAppConfig(input: UpdateAppConfigRequest, options?: RpcOptions): UnaryCall<UpdateAppConfigRequest, UpdateAppConfigResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateAppConfigRequest, UpdateAppConfigResponse>("unary", this._transport, method, opt, input);
    }
}
