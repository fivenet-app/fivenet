// @generated by protobuf-ts 2.10.0 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/wiki/wiki.proto" (package "services.wiki", syntax proto3)
// @ts-nocheck
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { WikiService } from "./wiki";
import type { ListPageActivityResponse } from "./wiki";
import type { ListPageActivityRequest } from "./wiki";
import type { DeletePageResponse } from "./wiki";
import type { DeletePageRequest } from "./wiki";
import type { UpdatePageResponse } from "./wiki";
import type { UpdatePageRequest } from "./wiki";
import type { CreatePageResponse } from "./wiki";
import type { CreatePageRequest } from "./wiki";
import type { GetPageResponse } from "./wiki";
import type { GetPageRequest } from "./wiki";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { ListPagesResponse } from "./wiki";
import type { ListPagesRequest } from "./wiki";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service services.wiki.WikiService
 */
export interface IWikiServiceClient {
    /**
     * @perm
     *
     * @generated from protobuf rpc: ListPages(services.wiki.ListPagesRequest) returns (services.wiki.ListPagesResponse);
     */
    listPages(input: ListPagesRequest, options?: RpcOptions): UnaryCall<ListPagesRequest, ListPagesResponse>;
    /**
     * @perm: Name=ListPages
     *
     * @generated from protobuf rpc: GetPage(services.wiki.GetPageRequest) returns (services.wiki.GetPageResponse);
     */
    getPage(input: GetPageRequest, options?: RpcOptions): UnaryCall<GetPageRequest, GetPageResponse>;
    /**
     * @perm: Attrs=Fields/StringList:[]string{"Public"}
     *
     * @generated from protobuf rpc: CreatePage(services.wiki.CreatePageRequest) returns (services.wiki.CreatePageResponse);
     */
    createPage(input: CreatePageRequest, options?: RpcOptions): UnaryCall<CreatePageRequest, CreatePageResponse>;
    /**
     * @perm: Name=ListPages
     *
     * @generated from protobuf rpc: UpdatePage(services.wiki.UpdatePageRequest) returns (services.wiki.UpdatePageResponse);
     */
    updatePage(input: UpdatePageRequest, options?: RpcOptions): UnaryCall<UpdatePageRequest, UpdatePageResponse>;
    /**
     * @perm
     *
     * @generated from protobuf rpc: DeletePage(services.wiki.DeletePageRequest) returns (services.wiki.DeletePageResponse);
     */
    deletePage(input: DeletePageRequest, options?: RpcOptions): UnaryCall<DeletePageRequest, DeletePageResponse>;
    /**
     * @perm
     *
     * @generated from protobuf rpc: ListPageActivity(services.wiki.ListPageActivityRequest) returns (services.wiki.ListPageActivityResponse);
     */
    listPageActivity(input: ListPageActivityRequest, options?: RpcOptions): UnaryCall<ListPageActivityRequest, ListPageActivityResponse>;
}
/**
 * @generated from protobuf service services.wiki.WikiService
 */
export class WikiServiceClient implements IWikiServiceClient, ServiceInfo {
    typeName = WikiService.typeName;
    methods = WikiService.methods;
    options = WikiService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: ListPages(services.wiki.ListPagesRequest) returns (services.wiki.ListPagesResponse);
     */
    listPages(input: ListPagesRequest, options?: RpcOptions): UnaryCall<ListPagesRequest, ListPagesResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<ListPagesRequest, ListPagesResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=ListPages
     *
     * @generated from protobuf rpc: GetPage(services.wiki.GetPageRequest) returns (services.wiki.GetPageResponse);
     */
    getPage(input: GetPageRequest, options?: RpcOptions): UnaryCall<GetPageRequest, GetPageResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetPageRequest, GetPageResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Attrs=Fields/StringList:[]string{"Public"}
     *
     * @generated from protobuf rpc: CreatePage(services.wiki.CreatePageRequest) returns (services.wiki.CreatePageResponse);
     */
    createPage(input: CreatePageRequest, options?: RpcOptions): UnaryCall<CreatePageRequest, CreatePageResponse> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<CreatePageRequest, CreatePageResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=ListPages
     *
     * @generated from protobuf rpc: UpdatePage(services.wiki.UpdatePageRequest) returns (services.wiki.UpdatePageResponse);
     */
    updatePage(input: UpdatePageRequest, options?: RpcOptions): UnaryCall<UpdatePageRequest, UpdatePageResponse> {
        const method = this.methods[3], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdatePageRequest, UpdatePageResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: DeletePage(services.wiki.DeletePageRequest) returns (services.wiki.DeletePageResponse);
     */
    deletePage(input: DeletePageRequest, options?: RpcOptions): UnaryCall<DeletePageRequest, DeletePageResponse> {
        const method = this.methods[4], opt = this._transport.mergeOptions(options);
        return stackIntercept<DeletePageRequest, DeletePageResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: ListPageActivity(services.wiki.ListPageActivityRequest) returns (services.wiki.ListPageActivityResponse);
     */
    listPageActivity(input: ListPageActivityRequest, options?: RpcOptions): UnaryCall<ListPageActivityRequest, ListPageActivityResponse> {
        const method = this.methods[5], opt = this._transport.mergeOptions(options);
        return stackIntercept<ListPageActivityRequest, ListPageActivityResponse>("unary", this._transport, method, opt, input);
    }
}
