// @generated by protobuf-ts 2.9.5 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/completor/completor.proto" (package "services.completor", syntax proto3)
// @ts-nocheck
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { CompletorService } from "./completor";
import type { CompleteCitizenLabelsResponse } from "./completor";
import type { CompleteCitizenLabelsRequest } from "./completor";
import type { ListLawBooksResponse } from "./completor";
import type { ListLawBooksRequest } from "./completor";
import type { CompleteDocumentCategoriesResponse } from "./completor";
import type { CompleteDocumentCategoriesRequest } from "./completor";
import type { CompleteJobsResponse } from "./completor";
import type { CompleteJobsRequest } from "./completor";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { CompleteCitizensRespoonse } from "./completor";
import type { CompleteCitizensRequest } from "./completor";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service services.completor.CompletorService
 */
export interface ICompletorServiceClient {
    /**
     * @perm
     *
     * @generated from protobuf rpc: CompleteCitizens(services.completor.CompleteCitizensRequest) returns (services.completor.CompleteCitizensRespoonse);
     */
    completeCitizens(input: CompleteCitizensRequest, options?: RpcOptions): UnaryCall<CompleteCitizensRequest, CompleteCitizensRespoonse>;
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: CompleteJobs(services.completor.CompleteJobsRequest) returns (services.completor.CompleteJobsResponse);
     */
    completeJobs(input: CompleteJobsRequest, options?: RpcOptions): UnaryCall<CompleteJobsRequest, CompleteJobsResponse>;
    /**
     * @perm: Attrs=Jobs/JobList
     *
     * @generated from protobuf rpc: CompleteDocumentCategories(services.completor.CompleteDocumentCategoriesRequest) returns (services.completor.CompleteDocumentCategoriesResponse);
     */
    completeDocumentCategories(input: CompleteDocumentCategoriesRequest, options?: RpcOptions): UnaryCall<CompleteDocumentCategoriesRequest, CompleteDocumentCategoriesResponse>;
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: ListLawBooks(services.completor.ListLawBooksRequest) returns (services.completor.ListLawBooksResponse);
     */
    listLawBooks(input: ListLawBooksRequest, options?: RpcOptions): UnaryCall<ListLawBooksRequest, ListLawBooksResponse>;
    /**
     * @perm: Attrs=Jobs/JobList
     *
     * @generated from protobuf rpc: CompleteCitizenLabels(services.completor.CompleteCitizenLabelsRequest) returns (services.completor.CompleteCitizenLabelsResponse);
     */
    completeCitizenLabels(input: CompleteCitizenLabelsRequest, options?: RpcOptions): UnaryCall<CompleteCitizenLabelsRequest, CompleteCitizenLabelsResponse>;
}
/**
 * @generated from protobuf service services.completor.CompletorService
 */
export class CompletorServiceClient implements ICompletorServiceClient, ServiceInfo {
    typeName = CompletorService.typeName;
    methods = CompletorService.methods;
    options = CompletorService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: CompleteCitizens(services.completor.CompleteCitizensRequest) returns (services.completor.CompleteCitizensRespoonse);
     */
    completeCitizens(input: CompleteCitizensRequest, options?: RpcOptions): UnaryCall<CompleteCitizensRequest, CompleteCitizensRespoonse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<CompleteCitizensRequest, CompleteCitizensRespoonse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: CompleteJobs(services.completor.CompleteJobsRequest) returns (services.completor.CompleteJobsResponse);
     */
    completeJobs(input: CompleteJobsRequest, options?: RpcOptions): UnaryCall<CompleteJobsRequest, CompleteJobsResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<CompleteJobsRequest, CompleteJobsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Attrs=Jobs/JobList
     *
     * @generated from protobuf rpc: CompleteDocumentCategories(services.completor.CompleteDocumentCategoriesRequest) returns (services.completor.CompleteDocumentCategoriesResponse);
     */
    completeDocumentCategories(input: CompleteDocumentCategoriesRequest, options?: RpcOptions): UnaryCall<CompleteDocumentCategoriesRequest, CompleteDocumentCategoriesResponse> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<CompleteDocumentCategoriesRequest, CompleteDocumentCategoriesResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: ListLawBooks(services.completor.ListLawBooksRequest) returns (services.completor.ListLawBooksResponse);
     */
    listLawBooks(input: ListLawBooksRequest, options?: RpcOptions): UnaryCall<ListLawBooksRequest, ListLawBooksResponse> {
        const method = this.methods[3], opt = this._transport.mergeOptions(options);
        return stackIntercept<ListLawBooksRequest, ListLawBooksResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Attrs=Jobs/JobList
     *
     * @generated from protobuf rpc: CompleteCitizenLabels(services.completor.CompleteCitizenLabelsRequest) returns (services.completor.CompleteCitizenLabelsResponse);
     */
    completeCitizenLabels(input: CompleteCitizenLabelsRequest, options?: RpcOptions): UnaryCall<CompleteCitizenLabelsRequest, CompleteCitizenLabelsResponse> {
        const method = this.methods[4], opt = this._transport.mergeOptions(options);
        return stackIntercept<CompleteCitizenLabelsRequest, CompleteCitizenLabelsResponse>("unary", this._transport, method, opt, input);
    }
}
