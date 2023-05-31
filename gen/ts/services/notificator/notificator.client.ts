// @generated by protobuf-ts 2.9.0 with parameter optimize_code_size,long_type_bigint
// @generated from protobuf file "services/notificator/notificator.proto" (package "services.notificator", syntax proto3)
// tslint:disable
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { NotificatorService } from "./notificator.js";
import type { StreamResponse } from "./notificator.js";
import type { StreamRequest } from "./notificator.js";
import type { ServerStreamingCall } from "@protobuf-ts/runtime-rpc";
import type { ReadNotificationsResponse } from "./notificator.js";
import type { ReadNotificationsRequest } from "./notificator.js";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { GetNotificationsResponse } from "./notificator.js";
import type { GetNotificationsRequest } from "./notificator.js";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service services.notificator.NotificatorService
 */
export interface INotificatorServiceClient {
    /**
     * @generated from protobuf rpc: GetNotifications(services.notificator.GetNotificationsRequest) returns (services.notificator.GetNotificationsResponse);
     */
    getNotifications(input: GetNotificationsRequest, options?: RpcOptions): UnaryCall<GetNotificationsRequest, GetNotificationsResponse>;
    /**
     * @generated from protobuf rpc: ReadNotifications(services.notificator.ReadNotificationsRequest) returns (services.notificator.ReadNotificationsResponse);
     */
    readNotifications(input: ReadNotificationsRequest, options?: RpcOptions): UnaryCall<ReadNotificationsRequest, ReadNotificationsResponse>;
    /**
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
     * @generated from protobuf rpc: GetNotifications(services.notificator.GetNotificationsRequest) returns (services.notificator.GetNotificationsResponse);
     */
    getNotifications(input: GetNotificationsRequest, options?: RpcOptions): UnaryCall<GetNotificationsRequest, GetNotificationsResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetNotificationsRequest, GetNotificationsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: ReadNotifications(services.notificator.ReadNotificationsRequest) returns (services.notificator.ReadNotificationsResponse);
     */
    readNotifications(input: ReadNotificationsRequest, options?: RpcOptions): UnaryCall<ReadNotificationsRequest, ReadNotificationsResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<ReadNotificationsRequest, ReadNotificationsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: Stream(services.notificator.StreamRequest) returns (stream services.notificator.StreamResponse);
     */
    stream(input: StreamRequest, options?: RpcOptions): ServerStreamingCall<StreamRequest, StreamResponse> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<StreamRequest, StreamResponse>("serverStreaming", this._transport, method, opt, input);
    }
}
