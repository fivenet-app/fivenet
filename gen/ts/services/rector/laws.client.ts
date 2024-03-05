// @generated by protobuf-ts 2.9.3 with parameter optimize_speed,long_type_bigint
// @generated from protobuf file "services/rector/laws.proto" (package "services.rector", syntax proto3)
// tslint:disable
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { RectorLawsService } from "./laws";
import type { DeleteLawResponse } from "./laws";
import type { DeleteLawRequest } from "./laws";
import type { CreateOrUpdateLawResponse } from "./laws";
import type { CreateOrUpdateLawRequest } from "./laws";
import type { DeleteLawBookResponse } from "./laws";
import type { DeleteLawBookRequest } from "./laws";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { CreateOrUpdateLawBookResponse } from "./laws";
import type { CreateOrUpdateLawBookRequest } from "./laws";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service services.rector.RectorLawsService
 */
export interface IRectorLawsServiceClient {
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: CreateOrUpdateLawBook(services.rector.CreateOrUpdateLawBookRequest) returns (services.rector.CreateOrUpdateLawBookResponse);
     */
    createOrUpdateLawBook(input: CreateOrUpdateLawBookRequest, options?: RpcOptions): UnaryCall<CreateOrUpdateLawBookRequest, CreateOrUpdateLawBookResponse>;
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: DeleteLawBook(services.rector.DeleteLawBookRequest) returns (services.rector.DeleteLawBookResponse);
     */
    deleteLawBook(input: DeleteLawBookRequest, options?: RpcOptions): UnaryCall<DeleteLawBookRequest, DeleteLawBookResponse>;
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: CreateOrUpdateLaw(services.rector.CreateOrUpdateLawRequest) returns (services.rector.CreateOrUpdateLawResponse);
     */
    createOrUpdateLaw(input: CreateOrUpdateLawRequest, options?: RpcOptions): UnaryCall<CreateOrUpdateLawRequest, CreateOrUpdateLawResponse>;
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: DeleteLaw(services.rector.DeleteLawRequest) returns (services.rector.DeleteLawResponse);
     */
    deleteLaw(input: DeleteLawRequest, options?: RpcOptions): UnaryCall<DeleteLawRequest, DeleteLawResponse>;
}
/**
 * @generated from protobuf service services.rector.RectorLawsService
 */
export class RectorLawsServiceClient implements IRectorLawsServiceClient, ServiceInfo {
    typeName = RectorLawsService.typeName;
    methods = RectorLawsService.methods;
    options = RectorLawsService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: CreateOrUpdateLawBook(services.rector.CreateOrUpdateLawBookRequest) returns (services.rector.CreateOrUpdateLawBookResponse);
     */
    createOrUpdateLawBook(input: CreateOrUpdateLawBookRequest, options?: RpcOptions): UnaryCall<CreateOrUpdateLawBookRequest, CreateOrUpdateLawBookResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<CreateOrUpdateLawBookRequest, CreateOrUpdateLawBookResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: DeleteLawBook(services.rector.DeleteLawBookRequest) returns (services.rector.DeleteLawBookResponse);
     */
    deleteLawBook(input: DeleteLawBookRequest, options?: RpcOptions): UnaryCall<DeleteLawBookRequest, DeleteLawBookResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<DeleteLawBookRequest, DeleteLawBookResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: CreateOrUpdateLaw(services.rector.CreateOrUpdateLawRequest) returns (services.rector.CreateOrUpdateLawResponse);
     */
    createOrUpdateLaw(input: CreateOrUpdateLawRequest, options?: RpcOptions): UnaryCall<CreateOrUpdateLawRequest, CreateOrUpdateLawResponse> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<CreateOrUpdateLawRequest, CreateOrUpdateLawResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=SuperUser
     *
     * @generated from protobuf rpc: DeleteLaw(services.rector.DeleteLawRequest) returns (services.rector.DeleteLawResponse);
     */
    deleteLaw(input: DeleteLawRequest, options?: RpcOptions): UnaryCall<DeleteLawRequest, DeleteLawResponse> {
        const method = this.methods[3], opt = this._transport.mergeOptions(options);
        return stackIntercept<DeleteLawRequest, DeleteLawResponse>("unary", this._transport, method, opt, input);
    }
}
