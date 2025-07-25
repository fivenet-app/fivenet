// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "services/settings/config.proto" (package "services.settings", syntax proto3)
// tslint:disable
// @ts-nocheck
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { ConfigService } from "./config";
import type { UpdateAppConfigResponse } from "./config";
import type { UpdateAppConfigRequest } from "./config";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { GetAppConfigResponse } from "./config";
import type { GetAppConfigRequest } from "./config";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service services.settings.ConfigService
 */
export interface IConfigServiceClient {
    /**
     * @perm: Name=Superuser
     *
     * @generated from protobuf rpc: GetAppConfig
     */
    getAppConfig(input: GetAppConfigRequest, options?: RpcOptions): UnaryCall<GetAppConfigRequest, GetAppConfigResponse>;
    /**
     * @perm: Name=Superuser
     *
     * @generated from protobuf rpc: UpdateAppConfig
     */
    updateAppConfig(input: UpdateAppConfigRequest, options?: RpcOptions): UnaryCall<UpdateAppConfigRequest, UpdateAppConfigResponse>;
}
/**
 * @generated from protobuf service services.settings.ConfigService
 */
export class ConfigServiceClient implements IConfigServiceClient, ServiceInfo {
    typeName = ConfigService.typeName;
    methods = ConfigService.methods;
    options = ConfigService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @perm: Name=Superuser
     *
     * @generated from protobuf rpc: GetAppConfig
     */
    getAppConfig(input: GetAppConfigRequest, options?: RpcOptions): UnaryCall<GetAppConfigRequest, GetAppConfigResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetAppConfigRequest, GetAppConfigResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=Superuser
     *
     * @generated from protobuf rpc: UpdateAppConfig
     */
    updateAppConfig(input: UpdateAppConfigRequest, options?: RpcOptions): UnaryCall<UpdateAppConfigRequest, UpdateAppConfigResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateAppConfigRequest, UpdateAppConfigResponse>("unary", this._transport, method, opt, input);
    }
}
