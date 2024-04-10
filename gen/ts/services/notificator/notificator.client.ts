// @generated by protobuf-ts 2.9.4 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/notificator/notificator.proto" (package "services.notificator", syntax proto3)
// tslint:disable
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { NotificatorService } from "./notificator";
import type { StreamResponse } from "./notificator";
import type { StreamRequest } from "./notificator";
import type { ServerStreamingCall } from "@protobuf-ts/runtime-rpc";
import type { MarkNotificationsResponse } from "./notificator";
import type { MarkNotificationsRequest } from "./notificator";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { GetNotificationsResponse } from "./notificator";
import type { GetNotificationsRequest } from "./notificator";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service services.notificator.NotificatorService
 */
export interface INotificatorServiceClient {
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: GetNotifications(services.notificator.GetNotificationsRequest) returns (services.notificator.GetNotificationsResponse);
     */
    getNotifications(input: GetNotificationsRequest, options?: RpcOptions): UnaryCall<GetNotificationsRequest, GetNotificationsResponse>;
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: MarkNotifications(services.notificator.MarkNotificationsRequest) returns (services.notificator.MarkNotificationsResponse);
     */
    markNotifications(input: MarkNotificationsRequest, options?: RpcOptions): UnaryCall<MarkNotificationsRequest, MarkNotificationsResponse>;
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: Stream(services.notificator.StreamRequest) returns (stream services.notificator.StreamResponse);
     */
    stream(input: StreamRequest, options?: RpcOptions): ServerStreamingCall<StreamRequest, StreamResponse>;
}
/**
 * @generated from protobuf service services.notificator.NotificatorService
 */
export class NotificatorServiceClient implements INotificatorServiceClient, ServiceInfo {
    typeName = NotificatorService.typeName;
    methods = NotificatorService.methods;
    options = NotificatorService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: GetNotifications(services.notificator.GetNotificationsRequest) returns (services.notificator.GetNotificationsResponse);
     */
    getNotifications(input: GetNotificationsRequest, options?: RpcOptions): UnaryCall<GetNotificationsRequest, GetNotificationsResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetNotificationsRequest, GetNotificationsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: MarkNotifications(services.notificator.MarkNotificationsRequest) returns (services.notificator.MarkNotificationsResponse);
     */
    markNotifications(input: MarkNotificationsRequest, options?: RpcOptions): UnaryCall<MarkNotificationsRequest, MarkNotificationsResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<MarkNotificationsRequest, MarkNotificationsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: Stream(services.notificator.StreamRequest) returns (stream services.notificator.StreamResponse);
     */
    stream(input: StreamRequest, options?: RpcOptions): ServerStreamingCall<StreamRequest, StreamResponse> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<StreamRequest, StreamResponse>("serverStreaming", this._transport, method, opt, input);
    }
}
