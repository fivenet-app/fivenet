<script lang="ts" setup>
import { Role } from '@fivenet/gen/resources/permissions/permissions_pb';
import { RpcError } from 'grpc-web';
import { DeleteRoleRequest, GetRolesRequest } from '@fivenet/gen/services/rector/rector_pb';

const { $grpc } = useNuxtApp();

const props = defineProps({
    roleId: {
        type: Number,
        required: true,
    }
});

async function getRole(): Promise<Role> {
    return new Promise(async (res, rej) => {
        const req = new GetRolesRequest();
        req.setWithPerms(false);
        req.setRank(1);

        try {
            const resp = await $grpc.getRectorClient().
                getRoles(req, null);

            const roles = resp.getRolesList();
            if (roles.length === 0) {
                return rej('No roles found!');
            }

            return res(roles[0]);
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const { data: role, pending, refresh, error } = await useLazyAsyncData(`rector-roles-${props.roleId}`, () => getRole());

async function deleteRole(id: number): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new DeleteRoleRequest();
        req.setRoleId(id);

        try {
            return $grpc.getRectorClient().
                deleteRole(req, null);
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <div>
        <h2 class="text-3xl text-white">TEST</h2>
    </div>
</template>
