// @generated by protobuf-ts 2.9.1 with parameter optimize_code_size,long_type_bigint
// @generated from protobuf file "services/rector/rector.proto" (package "services.rector", syntax proto3)
// tslint:disable
import { ServiceType } from "@protobuf-ts/runtime-rpc";
import { MessageType } from "@protobuf-ts/runtime";
import { Law } from "../../resources/laws/laws.js";
import { LawBook } from "../../resources/laws/laws.js";
import { AuditEntry } from "../../resources/rector/audit.js";
import { PaginationResponse } from "../../resources/common/database/database.js";
import { Timestamp } from "../../resources/timestamp/timestamp.js";
import { PaginationRequest } from "../../resources/common/database/database.js";
import { Permission } from "../../resources/permissions/permissions.js";
import { RoleAttribute } from "../../resources/permissions/permissions.js";
import { Role } from "../../resources/permissions/permissions.js";
import { JobProps } from "../../resources/users/jobs.js";
/**
 * @generated from protobuf message services.rector.GetJobPropsRequest
 */
export interface GetJobPropsRequest {
}
/**
 * @generated from protobuf message services.rector.GetJobPropsResponse
 */
export interface GetJobPropsResponse {
    /**
     * @generated from protobuf field: resources.users.JobProps job_props = 1;
     */
    jobProps?: JobProps;
}
/**
 * @generated from protobuf message services.rector.SetJobPropsRequest
 */
export interface SetJobPropsRequest {
    /**
     * @generated from protobuf field: resources.users.JobProps job_props = 1;
     */
    jobProps?: JobProps;
}
/**
 * @generated from protobuf message services.rector.SetJobPropsResponse
 */
export interface SetJobPropsResponse {
}
/**
 * @generated from protobuf message services.rector.GetRolesRequest
 */
export interface GetRolesRequest {
    /**
     * @generated from protobuf field: optional bool lowest_rank = 1;
     */
    lowestRank?: boolean;
}
/**
 * @generated from protobuf message services.rector.GetRolesResponse
 */
export interface GetRolesResponse {
    /**
     * @generated from protobuf field: repeated resources.permissions.Role roles = 1;
     */
    roles: Role[];
}
/**
 * @generated from protobuf message services.rector.GetRoleRequest
 */
export interface GetRoleRequest {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: bigint;
    /**
     * @generated from protobuf field: optional bool filtered = 2;
     */
    filtered?: boolean;
}
/**
 * @generated from protobuf message services.rector.GetRoleResponse
 */
export interface GetRoleResponse {
    /**
     * @generated from protobuf field: resources.permissions.Role role = 1;
     */
    role?: Role;
}
/**
 * @generated from protobuf message services.rector.CreateRoleRequest
 */
export interface CreateRoleRequest {
    /**
     * @generated from protobuf field: string job = 1;
     */
    job: string;
    /**
     * @generated from protobuf field: int32 grade = 2;
     */
    grade: number;
}
/**
 * @generated from protobuf message services.rector.CreateRoleResponse
 */
export interface CreateRoleResponse {
    /**
     * @generated from protobuf field: resources.permissions.Role role = 1;
     */
    role?: Role;
}
/**
 * @generated from protobuf message services.rector.DeleteRoleRequest
 */
export interface DeleteRoleRequest {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: bigint;
}
/**
 * @generated from protobuf message services.rector.DeleteRoleResponse
 */
export interface DeleteRoleResponse {
}
/**
 * @generated from protobuf message services.rector.UpdateRolePermsRequest
 */
export interface UpdateRolePermsRequest {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: bigint;
    /**
     * @generated from protobuf field: optional services.rector.PermsUpdate perms = 2;
     */
    perms?: PermsUpdate;
    /**
     * @generated from protobuf field: optional services.rector.AttrsUpdate attrs = 3;
     */
    attrs?: AttrsUpdate;
}
/**
 * @generated from protobuf message services.rector.PermsUpdate
 */
export interface PermsUpdate {
    /**
     * @generated from protobuf field: repeated services.rector.PermItem to_update = 1;
     */
    toUpdate: PermItem[];
    /**
     * @generated from protobuf field: repeated uint64 to_remove = 2;
     */
    toRemove: bigint[];
}
/**
 * @generated from protobuf message services.rector.PermItem
 */
export interface PermItem {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: bigint;
    /**
     * @generated from protobuf field: bool val = 2;
     */
    val: boolean;
}
/**
 * @generated from protobuf message services.rector.AttrsUpdate
 */
export interface AttrsUpdate {
    /**
     * @generated from protobuf field: repeated resources.permissions.RoleAttribute to_update = 1;
     */
    toUpdate: RoleAttribute[];
    /**
     * @generated from protobuf field: repeated resources.permissions.RoleAttribute to_remove = 2;
     */
    toRemove: RoleAttribute[];
}
/**
 * @generated from protobuf message services.rector.UpdateRolePermsResponse
 */
export interface UpdateRolePermsResponse {
}
/**
 * @generated from protobuf message services.rector.GetPermissionsRequest
 */
export interface GetPermissionsRequest {
    /**
     * @generated from protobuf field: uint64 role_id = 1;
     */
    roleId: bigint;
    /**
     * @generated from protobuf field: optional bool filtered = 2;
     */
    filtered?: boolean;
}
/**
 * @generated from protobuf message services.rector.GetPermissionsResponse
 */
export interface GetPermissionsResponse {
    /**
     * @generated from protobuf field: repeated resources.permissions.Permission permissions = 1;
     */
    permissions: Permission[];
    /**
     * @generated from protobuf field: repeated resources.permissions.RoleAttribute attributes = 2;
     */
    attributes: RoleAttribute[];
}
/**
 * @generated from protobuf message services.rector.ViewAuditLogRequest
 */
export interface ViewAuditLogRequest {
    /**
     * @generated from protobuf field: resources.common.database.PaginationRequest pagination = 1;
     */
    pagination?: PaginationRequest;
    /**
     * @generated from protobuf field: repeated int32 user_ids = 2;
     */
    userIds: number[];
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp from = 3;
     */
    from?: Timestamp;
    /**
     * @generated from protobuf field: optional resources.timestamp.Timestamp to = 4;
     */
    to?: Timestamp;
    /**
     * @generated from protobuf field: optional string service = 5;
     */
    service?: string;
    /**
     * @generated from protobuf field: optional string method = 6;
     */
    method?: string;
    /**
     * @generated from protobuf field: optional string search = 7;
     */
    search?: string;
}
/**
 * @generated from protobuf message services.rector.ViewAuditLogResponse
 */
export interface ViewAuditLogResponse {
    /**
     * @generated from protobuf field: resources.common.database.PaginationResponse pagination = 1;
     */
    pagination?: PaginationResponse;
    /**
     * @generated from protobuf field: repeated resources.rector.AuditEntry logs = 2;
     */
    logs: AuditEntry[];
}
/**
 * @generated from protobuf message services.rector.UpdateRoleLimitsRequest
 */
export interface UpdateRoleLimitsRequest {
    /**
     * @generated from protobuf field: uint64 role_id = 1;
     */
    roleId: bigint;
    /**
     * @generated from protobuf field: optional services.rector.PermsUpdate perms = 2;
     */
    perms?: PermsUpdate;
    /**
     * @generated from protobuf field: optional services.rector.AttrsUpdate attrs = 3;
     */
    attrs?: AttrsUpdate;
}
/**
 * @generated from protobuf message services.rector.UpdateRoleLimitsResponse
 */
export interface UpdateRoleLimitsResponse {
}
// Laws

/**
 * @generated from protobuf message services.rector.CreateOrUpdateLawBookRequest
 */
export interface CreateOrUpdateLawBookRequest {
    /**
     * @generated from protobuf field: resources.laws.LawBook lawBook = 1;
     */
    lawBook?: LawBook;
}
/**
 * @generated from protobuf message services.rector.CreateOrUpdateLawBookResponse
 */
export interface CreateOrUpdateLawBookResponse {
    /**
     * @generated from protobuf field: resources.laws.LawBook lawBook = 1;
     */
    lawBook?: LawBook;
}
/**
 * @generated from protobuf message services.rector.DeleteLawBookRequest
 */
export interface DeleteLawBookRequest {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: bigint;
}
/**
 * @generated from protobuf message services.rector.DeleteLawBookResponse
 */
export interface DeleteLawBookResponse {
}
/**
 * @generated from protobuf message services.rector.CreateOrUpdateLawRequest
 */
export interface CreateOrUpdateLawRequest {
    /**
     * @generated from protobuf field: resources.laws.Law law = 1;
     */
    law?: Law;
}
/**
 * @generated from protobuf message services.rector.CreateOrUpdateLawResponse
 */
export interface CreateOrUpdateLawResponse {
    /**
     * @generated from protobuf field: resources.laws.Law law = 1;
     */
    law?: Law;
}
/**
 * @generated from protobuf message services.rector.DeleteLawRequest
 */
export interface DeleteLawRequest {
    /**
     * @generated from protobuf field: uint64 id = 1;
     */
    id: bigint;
}
/**
 * @generated from protobuf message services.rector.DeleteLawResponse
 */
export interface DeleteLawResponse {
}
// @generated message type with reflection information, may provide speed optimized methods
class GetJobPropsRequest$Type extends MessageType<GetJobPropsRequest> {
    constructor() {
        super("services.rector.GetJobPropsRequest", []);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.GetJobPropsRequest
 */
export const GetJobPropsRequest = new GetJobPropsRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetJobPropsResponse$Type extends MessageType<GetJobPropsResponse> {
    constructor() {
        super("services.rector.GetJobPropsResponse", [
            { no: 1, name: "job_props", kind: "message", T: () => JobProps }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.GetJobPropsResponse
 */
export const GetJobPropsResponse = new GetJobPropsResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class SetJobPropsRequest$Type extends MessageType<SetJobPropsRequest> {
    constructor() {
        super("services.rector.SetJobPropsRequest", [
            { no: 1, name: "job_props", kind: "message", T: () => JobProps }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.SetJobPropsRequest
 */
export const SetJobPropsRequest = new SetJobPropsRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class SetJobPropsResponse$Type extends MessageType<SetJobPropsResponse> {
    constructor() {
        super("services.rector.SetJobPropsResponse", []);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.SetJobPropsResponse
 */
export const SetJobPropsResponse = new SetJobPropsResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetRolesRequest$Type extends MessageType<GetRolesRequest> {
    constructor() {
        super("services.rector.GetRolesRequest", [
            { no: 1, name: "lowest_rank", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.GetRolesRequest
 */
export const GetRolesRequest = new GetRolesRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetRolesResponse$Type extends MessageType<GetRolesResponse> {
    constructor() {
        super("services.rector.GetRolesResponse", [
            { no: 1, name: "roles", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Role }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.GetRolesResponse
 */
export const GetRolesResponse = new GetRolesResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetRoleRequest$Type extends MessageType<GetRoleRequest> {
    constructor() {
        super("services.rector.GetRoleRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 0 /*LongType.BIGINT*/ },
            { no: 2, name: "filtered", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.GetRoleRequest
 */
export const GetRoleRequest = new GetRoleRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetRoleResponse$Type extends MessageType<GetRoleResponse> {
    constructor() {
        super("services.rector.GetRoleResponse", [
            { no: 1, name: "role", kind: "message", T: () => Role }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.GetRoleResponse
 */
export const GetRoleResponse = new GetRoleResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateRoleRequest$Type extends MessageType<CreateRoleRequest> {
    constructor() {
        super("services.rector.CreateRoleRequest", [
            { no: 1, name: "job", kind: "scalar", T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "20" } } } },
            { no: 2, name: "grade", kind: "scalar", T: 5 /*ScalarType.INT32*/, options: { "validate.rules": { int32: { gt: 0 } } } }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.CreateRoleRequest
 */
export const CreateRoleRequest = new CreateRoleRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateRoleResponse$Type extends MessageType<CreateRoleResponse> {
    constructor() {
        super("services.rector.CreateRoleResponse", [
            { no: 1, name: "role", kind: "message", T: () => Role }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.CreateRoleResponse
 */
export const CreateRoleResponse = new CreateRoleResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteRoleRequest$Type extends MessageType<DeleteRoleRequest> {
    constructor() {
        super("services.rector.DeleteRoleRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 0 /*LongType.BIGINT*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.DeleteRoleRequest
 */
export const DeleteRoleRequest = new DeleteRoleRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteRoleResponse$Type extends MessageType<DeleteRoleResponse> {
    constructor() {
        super("services.rector.DeleteRoleResponse", []);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.DeleteRoleResponse
 */
export const DeleteRoleResponse = new DeleteRoleResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UpdateRolePermsRequest$Type extends MessageType<UpdateRolePermsRequest> {
    constructor() {
        super("services.rector.UpdateRolePermsRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 0 /*LongType.BIGINT*/ },
            { no: 2, name: "perms", kind: "message", T: () => PermsUpdate },
            { no: 3, name: "attrs", kind: "message", T: () => AttrsUpdate }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.UpdateRolePermsRequest
 */
export const UpdateRolePermsRequest = new UpdateRolePermsRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PermsUpdate$Type extends MessageType<PermsUpdate> {
    constructor() {
        super("services.rector.PermsUpdate", [
            { no: 1, name: "to_update", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => PermItem },
            { no: 2, name: "to_remove", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 4 /*ScalarType.UINT64*/, L: 0 /*LongType.BIGINT*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.PermsUpdate
 */
export const PermsUpdate = new PermsUpdate$Type();
// @generated message type with reflection information, may provide speed optimized methods
class PermItem$Type extends MessageType<PermItem> {
    constructor() {
        super("services.rector.PermItem", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 0 /*LongType.BIGINT*/ },
            { no: 2, name: "val", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.PermItem
 */
export const PermItem = new PermItem$Type();
// @generated message type with reflection information, may provide speed optimized methods
class AttrsUpdate$Type extends MessageType<AttrsUpdate> {
    constructor() {
        super("services.rector.AttrsUpdate", [
            { no: 1, name: "to_update", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => RoleAttribute },
            { no: 2, name: "to_remove", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => RoleAttribute }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.AttrsUpdate
 */
export const AttrsUpdate = new AttrsUpdate$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UpdateRolePermsResponse$Type extends MessageType<UpdateRolePermsResponse> {
    constructor() {
        super("services.rector.UpdateRolePermsResponse", []);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.UpdateRolePermsResponse
 */
export const UpdateRolePermsResponse = new UpdateRolePermsResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetPermissionsRequest$Type extends MessageType<GetPermissionsRequest> {
    constructor() {
        super("services.rector.GetPermissionsRequest", [
            { no: 1, name: "role_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 0 /*LongType.BIGINT*/ },
            { no: 2, name: "filtered", kind: "scalar", opt: true, T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.GetPermissionsRequest
 */
export const GetPermissionsRequest = new GetPermissionsRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetPermissionsResponse$Type extends MessageType<GetPermissionsResponse> {
    constructor() {
        super("services.rector.GetPermissionsResponse", [
            { no: 1, name: "permissions", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => Permission },
            { no: 2, name: "attributes", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => RoleAttribute }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.GetPermissionsResponse
 */
export const GetPermissionsResponse = new GetPermissionsResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ViewAuditLogRequest$Type extends MessageType<ViewAuditLogRequest> {
    constructor() {
        super("services.rector.ViewAuditLogRequest", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationRequest, options: { "validate.rules": { message: { required: true } } } },
            { no: 2, name: "user_ids", kind: "scalar", repeat: 1 /*RepeatType.PACKED*/, T: 5 /*ScalarType.INT32*/ },
            { no: 3, name: "from", kind: "message", T: () => Timestamp },
            { no: 4, name: "to", kind: "message", T: () => Timestamp },
            { no: 5, name: "service", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "64" } } } },
            { no: 6, name: "method", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "64" } } } },
            { no: 7, name: "search", kind: "scalar", opt: true, T: 9 /*ScalarType.STRING*/, options: { "validate.rules": { string: { maxLen: "128" } } } }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.ViewAuditLogRequest
 */
export const ViewAuditLogRequest = new ViewAuditLogRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class ViewAuditLogResponse$Type extends MessageType<ViewAuditLogResponse> {
    constructor() {
        super("services.rector.ViewAuditLogResponse", [
            { no: 1, name: "pagination", kind: "message", T: () => PaginationResponse },
            { no: 2, name: "logs", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => AuditEntry }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.ViewAuditLogResponse
 */
export const ViewAuditLogResponse = new ViewAuditLogResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UpdateRoleLimitsRequest$Type extends MessageType<UpdateRoleLimitsRequest> {
    constructor() {
        super("services.rector.UpdateRoleLimitsRequest", [
            { no: 1, name: "role_id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 0 /*LongType.BIGINT*/ },
            { no: 2, name: "perms", kind: "message", T: () => PermsUpdate },
            { no: 3, name: "attrs", kind: "message", T: () => AttrsUpdate }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.UpdateRoleLimitsRequest
 */
export const UpdateRoleLimitsRequest = new UpdateRoleLimitsRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class UpdateRoleLimitsResponse$Type extends MessageType<UpdateRoleLimitsResponse> {
    constructor() {
        super("services.rector.UpdateRoleLimitsResponse", []);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.UpdateRoleLimitsResponse
 */
export const UpdateRoleLimitsResponse = new UpdateRoleLimitsResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateOrUpdateLawBookRequest$Type extends MessageType<CreateOrUpdateLawBookRequest> {
    constructor() {
        super("services.rector.CreateOrUpdateLawBookRequest", [
            { no: 1, name: "lawBook", kind: "message", T: () => LawBook, options: { "validate.rules": { message: { required: true } } } }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.CreateOrUpdateLawBookRequest
 */
export const CreateOrUpdateLawBookRequest = new CreateOrUpdateLawBookRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateOrUpdateLawBookResponse$Type extends MessageType<CreateOrUpdateLawBookResponse> {
    constructor() {
        super("services.rector.CreateOrUpdateLawBookResponse", [
            { no: 1, name: "lawBook", kind: "message", T: () => LawBook }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.CreateOrUpdateLawBookResponse
 */
export const CreateOrUpdateLawBookResponse = new CreateOrUpdateLawBookResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteLawBookRequest$Type extends MessageType<DeleteLawBookRequest> {
    constructor() {
        super("services.rector.DeleteLawBookRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 0 /*LongType.BIGINT*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.DeleteLawBookRequest
 */
export const DeleteLawBookRequest = new DeleteLawBookRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteLawBookResponse$Type extends MessageType<DeleteLawBookResponse> {
    constructor() {
        super("services.rector.DeleteLawBookResponse", []);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.DeleteLawBookResponse
 */
export const DeleteLawBookResponse = new DeleteLawBookResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateOrUpdateLawRequest$Type extends MessageType<CreateOrUpdateLawRequest> {
    constructor() {
        super("services.rector.CreateOrUpdateLawRequest", [
            { no: 1, name: "law", kind: "message", T: () => Law, options: { "validate.rules": { message: { required: true } } } }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.CreateOrUpdateLawRequest
 */
export const CreateOrUpdateLawRequest = new CreateOrUpdateLawRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateOrUpdateLawResponse$Type extends MessageType<CreateOrUpdateLawResponse> {
    constructor() {
        super("services.rector.CreateOrUpdateLawResponse", [
            { no: 1, name: "law", kind: "message", T: () => Law }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.CreateOrUpdateLawResponse
 */
export const CreateOrUpdateLawResponse = new CreateOrUpdateLawResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteLawRequest$Type extends MessageType<DeleteLawRequest> {
    constructor() {
        super("services.rector.DeleteLawRequest", [
            { no: 1, name: "id", kind: "scalar", T: 4 /*ScalarType.UINT64*/, L: 0 /*LongType.BIGINT*/ }
        ]);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.DeleteLawRequest
 */
export const DeleteLawRequest = new DeleteLawRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class DeleteLawResponse$Type extends MessageType<DeleteLawResponse> {
    constructor() {
        super("services.rector.DeleteLawResponse", []);
    }
}
/**
 * @generated MessageType for protobuf message services.rector.DeleteLawResponse
 */
export const DeleteLawResponse = new DeleteLawResponse$Type();
/**
 * @generated ServiceType for protobuf service services.rector.RectorService
 */
export const RectorService = new ServiceType("services.rector.RectorService", [
    { name: "GetJobProps", options: {}, I: GetJobPropsRequest, O: GetJobPropsResponse },
    { name: "SetJobProps", options: {}, I: SetJobPropsRequest, O: SetJobPropsResponse },
    { name: "GetRoles", options: {}, I: GetRolesRequest, O: GetRolesResponse },
    { name: "GetRole", options: {}, I: GetRoleRequest, O: GetRoleResponse },
    { name: "CreateRole", options: {}, I: CreateRoleRequest, O: CreateRoleResponse },
    { name: "DeleteRole", options: {}, I: DeleteRoleRequest, O: DeleteRoleResponse },
    { name: "UpdateRolePerms", options: {}, I: UpdateRolePermsRequest, O: UpdateRolePermsResponse },
    { name: "GetPermissions", options: {}, I: GetPermissionsRequest, O: GetPermissionsResponse },
    { name: "ViewAuditLog", options: {}, I: ViewAuditLogRequest, O: ViewAuditLogResponse },
    { name: "UpdateRoleLimits", options: {}, I: UpdateRoleLimitsRequest, O: UpdateRoleLimitsResponse },
    { name: "CreateOrUpdateLawBook", options: {}, I: CreateOrUpdateLawBookRequest, O: CreateOrUpdateLawBookResponse },
    { name: "DeleteLawBook", options: {}, I: DeleteLawBookRequest, O: DeleteLawBookResponse },
    { name: "CreateOrUpdateLaw", options: {}, I: CreateOrUpdateLawRequest, O: CreateOrUpdateLawResponse },
    { name: "DeleteLaw", options: {}, I: DeleteLawRequest, O: DeleteLawResponse }
]);
