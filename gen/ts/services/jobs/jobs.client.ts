// @generated by protobuf-ts 2.11.1 with parameter force_server_none,long_type_number,optimize_speed,ts_nocheck
// @generated from protobuf file "services/jobs/jobs.proto" (package "services.jobs", syntax proto3)
// tslint:disable
// @ts-nocheck
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { JobsService } from "./jobs";
import type { SetMOTDResponse } from "./jobs";
import type { SetMOTDRequest } from "./jobs";
import type { GetMOTDResponse } from "./jobs";
import type { GetMOTDRequest } from "./jobs";
import type { GetColleagueLabelsStatsResponse } from "./jobs";
import type { GetColleagueLabelsStatsRequest } from "./jobs";
import type { ManageLabelsResponse } from "./jobs";
import type { ManageLabelsRequest } from "./jobs";
import type { GetColleagueLabelsResponse } from "./jobs";
import type { GetColleagueLabelsRequest } from "./jobs";
import type { SetColleaguePropsResponse } from "./jobs";
import type { SetColleaguePropsRequest } from "./jobs";
import type { ListColleagueActivityResponse } from "./jobs";
import type { ListColleagueActivityRequest } from "./jobs";
import type { GetColleagueResponse } from "./jobs";
import type { GetColleagueRequest } from "./jobs";
import type { GetSelfResponse } from "./jobs";
import type { GetSelfRequest } from "./jobs";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { ListColleaguesResponse } from "./jobs";
import type { ListColleaguesRequest } from "./jobs";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service services.jobs.JobsService
 */
export interface IJobsServiceClient {
    /**
     * @perm
     *
     * @generated from protobuf rpc: ListColleagues
     */
    listColleagues(input: ListColleaguesRequest, options?: RpcOptions): UnaryCall<ListColleaguesRequest, ListColleaguesResponse>;
    /**
     * @perm: Name=ListColleagues
     *
     * @generated from protobuf rpc: GetSelf
     */
    getSelf(input: GetSelfRequest, options?: RpcOptions): UnaryCall<GetSelfRequest, GetSelfResponse>;
    /**
     * @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}|Types/StringList:[]string{"Note", "Labels"}
     *
     * @generated from protobuf rpc: GetColleague
     */
    getColleague(input: GetColleagueRequest, options?: RpcOptions): UnaryCall<GetColleagueRequest, GetColleagueResponse>;
    /**
     * @perm: Attrs=Types/StringList:[]string{"HIRED", "FIRED", "PROMOTED", "DEMOTED", "ABSENCE_DATE", "NOTE", "LABELS", "NAME"}
     *
     * @generated from protobuf rpc: ListColleagueActivity
     */
    listColleagueActivity(input: ListColleagueActivityRequest, options?: RpcOptions): UnaryCall<ListColleagueActivityRequest, ListColleagueActivityResponse>;
    /**
     * @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}|Types/StringList:[]string{"AbsenceDate", "Note", "Labels", "Name"}
     *
     * @generated from protobuf rpc: SetColleagueProps
     */
    setColleagueProps(input: SetColleaguePropsRequest, options?: RpcOptions): UnaryCall<SetColleaguePropsRequest, SetColleaguePropsResponse>;
    /**
     * @perm: Name=GetColleague
     *
     * @generated from protobuf rpc: GetColleagueLabels
     */
    getColleagueLabels(input: GetColleagueLabelsRequest, options?: RpcOptions): UnaryCall<GetColleagueLabelsRequest, GetColleagueLabelsResponse>;
    /**
     * @perm
     *
     * @generated from protobuf rpc: ManageLabels
     */
    manageLabels(input: ManageLabelsRequest, options?: RpcOptions): UnaryCall<ManageLabelsRequest, ManageLabelsResponse>;
    /**
     * @perm: Name=GetColleague
     *
     * @generated from protobuf rpc: GetColleagueLabelsStats
     */
    getColleagueLabelsStats(input: GetColleagueLabelsStatsRequest, options?: RpcOptions): UnaryCall<GetColleagueLabelsStatsRequest, GetColleagueLabelsStatsResponse>;
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: GetMOTD
     */
    getMOTD(input: GetMOTDRequest, options?: RpcOptions): UnaryCall<GetMOTDRequest, GetMOTDResponse>;
    /**
     * @perm
     *
     * @generated from protobuf rpc: SetMOTD
     */
    setMOTD(input: SetMOTDRequest, options?: RpcOptions): UnaryCall<SetMOTDRequest, SetMOTDResponse>;
}
/**
 * @generated from protobuf service services.jobs.JobsService
 */
export class JobsServiceClient implements IJobsServiceClient, ServiceInfo {
    typeName = JobsService.typeName;
    methods = JobsService.methods;
    options = JobsService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: ListColleagues
     */
    listColleagues(input: ListColleaguesRequest, options?: RpcOptions): UnaryCall<ListColleaguesRequest, ListColleaguesResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<ListColleaguesRequest, ListColleaguesResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=ListColleagues
     *
     * @generated from protobuf rpc: GetSelf
     */
    getSelf(input: GetSelfRequest, options?: RpcOptions): UnaryCall<GetSelfRequest, GetSelfResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetSelfRequest, GetSelfResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}|Types/StringList:[]string{"Note", "Labels"}
     *
     * @generated from protobuf rpc: GetColleague
     */
    getColleague(input: GetColleagueRequest, options?: RpcOptions): UnaryCall<GetColleagueRequest, GetColleagueResponse> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetColleagueRequest, GetColleagueResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Attrs=Types/StringList:[]string{"HIRED", "FIRED", "PROMOTED", "DEMOTED", "ABSENCE_DATE", "NOTE", "LABELS", "NAME"}
     *
     * @generated from protobuf rpc: ListColleagueActivity
     */
    listColleagueActivity(input: ListColleagueActivityRequest, options?: RpcOptions): UnaryCall<ListColleagueActivityRequest, ListColleagueActivityResponse> {
        const method = this.methods[3], opt = this._transport.mergeOptions(options);
        return stackIntercept<ListColleagueActivityRequest, ListColleagueActivityResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}|Types/StringList:[]string{"AbsenceDate", "Note", "Labels", "Name"}
     *
     * @generated from protobuf rpc: SetColleagueProps
     */
    setColleagueProps(input: SetColleaguePropsRequest, options?: RpcOptions): UnaryCall<SetColleaguePropsRequest, SetColleaguePropsResponse> {
        const method = this.methods[4], opt = this._transport.mergeOptions(options);
        return stackIntercept<SetColleaguePropsRequest, SetColleaguePropsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=GetColleague
     *
     * @generated from protobuf rpc: GetColleagueLabels
     */
    getColleagueLabels(input: GetColleagueLabelsRequest, options?: RpcOptions): UnaryCall<GetColleagueLabelsRequest, GetColleagueLabelsResponse> {
        const method = this.methods[5], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetColleagueLabelsRequest, GetColleagueLabelsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: ManageLabels
     */
    manageLabels(input: ManageLabelsRequest, options?: RpcOptions): UnaryCall<ManageLabelsRequest, ManageLabelsResponse> {
        const method = this.methods[6], opt = this._transport.mergeOptions(options);
        return stackIntercept<ManageLabelsRequest, ManageLabelsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=GetColleague
     *
     * @generated from protobuf rpc: GetColleagueLabelsStats
     */
    getColleagueLabelsStats(input: GetColleagueLabelsStatsRequest, options?: RpcOptions): UnaryCall<GetColleagueLabelsStatsRequest, GetColleagueLabelsStatsResponse> {
        const method = this.methods[7], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetColleagueLabelsStatsRequest, GetColleagueLabelsStatsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=Any
     *
     * @generated from protobuf rpc: GetMOTD
     */
    getMOTD(input: GetMOTDRequest, options?: RpcOptions): UnaryCall<GetMOTDRequest, GetMOTDResponse> {
        const method = this.methods[8], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetMOTDRequest, GetMOTDResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: SetMOTD
     */
    setMOTD(input: SetMOTDRequest, options?: RpcOptions): UnaryCall<SetMOTDRequest, SetMOTDResponse> {
        const method = this.methods[9], opt = this._transport.mergeOptions(options);
        return stackIntercept<SetMOTDRequest, SetMOTDResponse>("unary", this._transport, method, opt, input);
    }
}
