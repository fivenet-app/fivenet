// @generated by protobuf-ts 2.10.0 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "services/settings/settings.proto" (package "services.settings", syntax proto3)
// @ts-nocheck
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { SettingsService } from "./settings";
import type { DeleteFactionResponse } from "./settings";
import type { DeleteFactionRequest } from "./settings";
import type { UpdateJobLimitsResponse } from "./settings";
import type { UpdateJobLimitsRequest } from "./settings";
import type { GetJobLimitsResponse } from "./settings";
import type { GetJobLimitsRequest } from "./settings";
import type { GetAllPermissionsResponse } from "./settings";
import type { GetAllPermissionsRequest } from "./settings";
import type { ViewAuditLogResponse } from "./settings";
import type { ViewAuditLogRequest } from "./settings";
import type { GetEffectivePermissionsResponse } from "./settings";
import type { GetEffectivePermissionsRequest } from "./settings";
import type { GetPermissionsResponse } from "./settings";
import type { GetPermissionsRequest } from "./settings";
import type { UpdateRolePermsResponse } from "./settings";
import type { UpdateRolePermsRequest } from "./settings";
import type { DeleteRoleResponse } from "./settings";
import type { DeleteRoleRequest } from "./settings";
import type { CreateRoleResponse } from "./settings";
import type { CreateRoleRequest } from "./settings";
import type { GetRoleResponse } from "./settings";
import type { GetRoleRequest } from "./settings";
import type { GetRolesResponse } from "./settings";
import type { GetRolesRequest } from "./settings";
import type { SetJobPropsResponse } from "./settings";
import type { SetJobPropsRequest } from "./settings";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { GetJobPropsResponse } from "./settings";
import type { GetJobPropsRequest } from "./settings";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service services.settings.SettingsService
 */
export interface ISettingsServiceClient {
    /**
     * @perm
     *
     * @generated from protobuf rpc: GetJobProps(services.settings.GetJobPropsRequest) returns (services.settings.GetJobPropsResponse);
     */
    getJobProps(input: GetJobPropsRequest, options?: RpcOptions): UnaryCall<GetJobPropsRequest, GetJobPropsResponse>;
    /**
     * @perm
     *
     * @generated from protobuf rpc: SetJobProps(services.settings.SetJobPropsRequest) returns (services.settings.SetJobPropsResponse);
     */
    setJobProps(input: SetJobPropsRequest, options?: RpcOptions): UnaryCall<SetJobPropsRequest, SetJobPropsResponse>;
    /**
     * @perm
     *
     * @generated from protobuf rpc: GetRoles(services.settings.GetRolesRequest) returns (services.settings.GetRolesResponse);
     */
    getRoles(input: GetRolesRequest, options?: RpcOptions): UnaryCall<GetRolesRequest, GetRolesResponse>;
    /**
     * @perm: Name=GetRoles
     *
     * @generated from protobuf rpc: GetRole(services.settings.GetRoleRequest) returns (services.settings.GetRoleResponse);
     */
    getRole(input: GetRoleRequest, options?: RpcOptions): UnaryCall<GetRoleRequest, GetRoleResponse>;
    /**
     * @perm
     *
     * @generated from protobuf rpc: CreateRole(services.settings.CreateRoleRequest) returns (services.settings.CreateRoleResponse);
     */
    createRole(input: CreateRoleRequest, options?: RpcOptions): UnaryCall<CreateRoleRequest, CreateRoleResponse>;
    /**
     * @perm
     *
     * @generated from protobuf rpc: DeleteRole(services.settings.DeleteRoleRequest) returns (services.settings.DeleteRoleResponse);
     */
    deleteRole(input: DeleteRoleRequest, options?: RpcOptions): UnaryCall<DeleteRoleRequest, DeleteRoleResponse>;
    /**
     * @perm
     *
     * @generated from protobuf rpc: UpdateRolePerms(services.settings.UpdateRolePermsRequest) returns (services.settings.UpdateRolePermsResponse);
     */
    updateRolePerms(input: UpdateRolePermsRequest, options?: RpcOptions): UnaryCall<UpdateRolePermsRequest, UpdateRolePermsResponse>;
    /**
     * @perm: Name=GetRoles
     *
     * @generated from protobuf rpc: GetPermissions(services.settings.GetPermissionsRequest) returns (services.settings.GetPermissionsResponse);
     */
    getPermissions(input: GetPermissionsRequest, options?: RpcOptions): UnaryCall<GetPermissionsRequest, GetPermissionsResponse>;
    /**
     * @perm: Name=GetRoles
     *
     * @generated from protobuf rpc: GetEffectivePermissions(services.settings.GetEffectivePermissionsRequest) returns (services.settings.GetEffectivePermissionsResponse);
     */
    getEffectivePermissions(input: GetEffectivePermissionsRequest, options?: RpcOptions): UnaryCall<GetEffectivePermissionsRequest, GetEffectivePermissionsResponse>;
    /**
     * @perm
     *
     * @generated from protobuf rpc: ViewAuditLog(services.settings.ViewAuditLogRequest) returns (services.settings.ViewAuditLogResponse);
     */
    viewAuditLog(input: ViewAuditLogRequest, options?: RpcOptions): UnaryCall<ViewAuditLogRequest, ViewAuditLogResponse>;
    /**
     * @perm: Name=Superuser
     *
     * @generated from protobuf rpc: GetAllPermissions(services.settings.GetAllPermissionsRequest) returns (services.settings.GetAllPermissionsResponse);
     */
    getAllPermissions(input: GetAllPermissionsRequest, options?: RpcOptions): UnaryCall<GetAllPermissionsRequest, GetAllPermissionsResponse>;
    /**
     * @perm: Name=Superuser
     *
     * @generated from protobuf rpc: GetJobLimits(services.settings.GetJobLimitsRequest) returns (services.settings.GetJobLimitsResponse);
     */
    getJobLimits(input: GetJobLimitsRequest, options?: RpcOptions): UnaryCall<GetJobLimitsRequest, GetJobLimitsResponse>;
    /**
     * @perm: Name=Superuser
     *
     * @generated from protobuf rpc: UpdateJobLimits(services.settings.UpdateJobLimitsRequest) returns (services.settings.UpdateJobLimitsResponse);
     */
    updateJobLimits(input: UpdateJobLimitsRequest, options?: RpcOptions): UnaryCall<UpdateJobLimitsRequest, UpdateJobLimitsResponse>;
    /**
     * @perm: Name=Superuser
     *
     * @generated from protobuf rpc: DeleteFaction(services.settings.DeleteFactionRequest) returns (services.settings.DeleteFactionResponse);
     */
    deleteFaction(input: DeleteFactionRequest, options?: RpcOptions): UnaryCall<DeleteFactionRequest, DeleteFactionResponse>;
}
/**
 * @generated from protobuf service services.settings.SettingsService
 */
export class SettingsServiceClient implements ISettingsServiceClient, ServiceInfo {
    typeName = SettingsService.typeName;
    methods = SettingsService.methods;
    options = SettingsService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: GetJobProps(services.settings.GetJobPropsRequest) returns (services.settings.GetJobPropsResponse);
     */
    getJobProps(input: GetJobPropsRequest, options?: RpcOptions): UnaryCall<GetJobPropsRequest, GetJobPropsResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetJobPropsRequest, GetJobPropsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: SetJobProps(services.settings.SetJobPropsRequest) returns (services.settings.SetJobPropsResponse);
     */
    setJobProps(input: SetJobPropsRequest, options?: RpcOptions): UnaryCall<SetJobPropsRequest, SetJobPropsResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<SetJobPropsRequest, SetJobPropsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: GetRoles(services.settings.GetRolesRequest) returns (services.settings.GetRolesResponse);
     */
    getRoles(input: GetRolesRequest, options?: RpcOptions): UnaryCall<GetRolesRequest, GetRolesResponse> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetRolesRequest, GetRolesResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=GetRoles
     *
     * @generated from protobuf rpc: GetRole(services.settings.GetRoleRequest) returns (services.settings.GetRoleResponse);
     */
    getRole(input: GetRoleRequest, options?: RpcOptions): UnaryCall<GetRoleRequest, GetRoleResponse> {
        const method = this.methods[3], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetRoleRequest, GetRoleResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: CreateRole(services.settings.CreateRoleRequest) returns (services.settings.CreateRoleResponse);
     */
    createRole(input: CreateRoleRequest, options?: RpcOptions): UnaryCall<CreateRoleRequest, CreateRoleResponse> {
        const method = this.methods[4], opt = this._transport.mergeOptions(options);
        return stackIntercept<CreateRoleRequest, CreateRoleResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: DeleteRole(services.settings.DeleteRoleRequest) returns (services.settings.DeleteRoleResponse);
     */
    deleteRole(input: DeleteRoleRequest, options?: RpcOptions): UnaryCall<DeleteRoleRequest, DeleteRoleResponse> {
        const method = this.methods[5], opt = this._transport.mergeOptions(options);
        return stackIntercept<DeleteRoleRequest, DeleteRoleResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: UpdateRolePerms(services.settings.UpdateRolePermsRequest) returns (services.settings.UpdateRolePermsResponse);
     */
    updateRolePerms(input: UpdateRolePermsRequest, options?: RpcOptions): UnaryCall<UpdateRolePermsRequest, UpdateRolePermsResponse> {
        const method = this.methods[6], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateRolePermsRequest, UpdateRolePermsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=GetRoles
     *
     * @generated from protobuf rpc: GetPermissions(services.settings.GetPermissionsRequest) returns (services.settings.GetPermissionsResponse);
     */
    getPermissions(input: GetPermissionsRequest, options?: RpcOptions): UnaryCall<GetPermissionsRequest, GetPermissionsResponse> {
        const method = this.methods[7], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetPermissionsRequest, GetPermissionsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=GetRoles
     *
     * @generated from protobuf rpc: GetEffectivePermissions(services.settings.GetEffectivePermissionsRequest) returns (services.settings.GetEffectivePermissionsResponse);
     */
    getEffectivePermissions(input: GetEffectivePermissionsRequest, options?: RpcOptions): UnaryCall<GetEffectivePermissionsRequest, GetEffectivePermissionsResponse> {
        const method = this.methods[8], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetEffectivePermissionsRequest, GetEffectivePermissionsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm
     *
     * @generated from protobuf rpc: ViewAuditLog(services.settings.ViewAuditLogRequest) returns (services.settings.ViewAuditLogResponse);
     */
    viewAuditLog(input: ViewAuditLogRequest, options?: RpcOptions): UnaryCall<ViewAuditLogRequest, ViewAuditLogResponse> {
        const method = this.methods[9], opt = this._transport.mergeOptions(options);
        return stackIntercept<ViewAuditLogRequest, ViewAuditLogResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=Superuser
     *
     * @generated from protobuf rpc: GetAllPermissions(services.settings.GetAllPermissionsRequest) returns (services.settings.GetAllPermissionsResponse);
     */
    getAllPermissions(input: GetAllPermissionsRequest, options?: RpcOptions): UnaryCall<GetAllPermissionsRequest, GetAllPermissionsResponse> {
        const method = this.methods[10], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetAllPermissionsRequest, GetAllPermissionsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=Superuser
     *
     * @generated from protobuf rpc: GetJobLimits(services.settings.GetJobLimitsRequest) returns (services.settings.GetJobLimitsResponse);
     */
    getJobLimits(input: GetJobLimitsRequest, options?: RpcOptions): UnaryCall<GetJobLimitsRequest, GetJobLimitsResponse> {
        const method = this.methods[11], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetJobLimitsRequest, GetJobLimitsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=Superuser
     *
     * @generated from protobuf rpc: UpdateJobLimits(services.settings.UpdateJobLimitsRequest) returns (services.settings.UpdateJobLimitsResponse);
     */
    updateJobLimits(input: UpdateJobLimitsRequest, options?: RpcOptions): UnaryCall<UpdateJobLimitsRequest, UpdateJobLimitsResponse> {
        const method = this.methods[12], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateJobLimitsRequest, UpdateJobLimitsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @perm: Name=Superuser
     *
     * @generated from protobuf rpc: DeleteFaction(services.settings.DeleteFactionRequest) returns (services.settings.DeleteFactionResponse);
     */
    deleteFaction(input: DeleteFactionRequest, options?: RpcOptions): UnaryCall<DeleteFactionRequest, DeleteFactionResponse> {
        const method = this.methods[13], opt = this._transport.mergeOptions(options);
        return stackIntercept<DeleteFactionRequest, DeleteFactionResponse>("unary", this._transport, method, opt, input);
    }
}
