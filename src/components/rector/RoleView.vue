<script lang="ts" setup>
import { Role } from '@fivenet/gen/resources/permissions/permissions_pb';
import { RpcError } from 'grpc-web';
import { CreateRoleRequest, DeleteRoleRequest, GetRoleRequest } from '@fivenet/gen/services/rector/rector_pb';
import { MagnifyingGlassIcon } from '@heroicons/vue/24/solid';
import DataPendingBlock from '../partials/DataPendingBlock.vue';
import DataErrorBlock from '../partials/DataErrorBlock.vue';
import { TrashIcon } from '@heroicons/vue/20/solid';

const { $grpc } = useNuxtApp();

const props = defineProps({
    roleId: {
        type: Number,
        required: true,
    }
});

async function getRole(): Promise<Role> {
    return new Promise(async (res, rej) => {
        const req = new GetRoleRequest();
        req.setId(props.roleId);

        try {
            const resp = await $grpc.getRectorClient().
                getRole(req, null);

            return res(resp.getRole()!);
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const { data: role, pending, refresh, error } = await useLazyAsyncData(`rector-roles-${props.roleId}`, () => getRole());

async function createRole(job: string, grade: number): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new CreateRoleRequest();
        req.setJob(job);
        req.setGrade(grade);

        try {
            return $grpc.getRectorClient().
                createRole(req, null);
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function deleteRole(id: number): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new DeleteRoleRequest();
        req.setId(id);

        try {
            return $grpc.getRectorClient().
                deleteRole(req, null);
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function removePermission(id: number): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new DeleteRoleRequest();
        req.setId(id);

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
    <div v-if="role">
        <h2 class="text-3xl text-white">
            {{ toTitleCase(role?.getName()!) }}
        </h2>
        <p class="text-gray-400 text-sm">
            {{ role.getDescription() }}
        </p>
        <div class="relative">
            <div class="absolute inset-0 flex items-center" aria-hidden="true">
                <div class="w-full border-t border-gray-300"></div>
            </div>
            <div class="relative flex justify-center">
                <span class="bg-primary-900 px-3 text-base font-semibold leading-6 text-gray-200">Permissions</span>
            </div>
        </div>
        <div class="py-2">
            <!-- TODO add CreateRole logic-->
            <div class="px-2 sm:px-6 lg:px-8">
                <div class="flow-root mt-2">
                    <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                        <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                            <DataPendingBlock v-if="pending" message="Loading role permissions..." />
                            <DataErrorBlock v-else-if="error" title="Unable to load role permissions!" :retry="refresh" />
                            <button v-else-if="role.getPermissionsList().length == 0" type="button"
                                class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2">
                                <MagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
                                <span class="block mt-2 text-sm font-semibold text-gray-300">
                                    No Permissions found for this role.
                                </span>
                            </button>
                            <div v-else>
                                <table class="min-w-full divide-y divide-base-600">
                                    <thead>
                                        <tr>
                                            <th scope="col"
                                                class="py-3.5 pl-3 pr-4 sm:pr-0 text-left text-sm font-semibold text-neutral">
                                                Actions
                                            </th>
                                            <th scope="col"
                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                                Name
                                            </th>
                                            <th scope="col"
                                                class="py-3.5 px-2 text-center text-sm font-semibold text-neutral">
                                                Description
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody class="divide-y divide-base-800">
                                        <tr v-for="perm in role.getPermissionsList()" :key="perm.getId()">
                                            <td class="whitespace-nowrap py-2 pl-3 pr-4 text-sm font-medium sm:pr-0">
                                                <div class="flex flex-row">
                                                    <button v-can="'RectorService.UpdateRole'">
                                                        <TrashIcon class="w-6 h-6 mx-auto text-neutral"
                                                            @click="removePermission(perm.getId())" />
                                                    </button>
                                                </div>
                                            </td>
                                            <td
                                                class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
                                                {{ perm.getName() }}
                                            </td>
                                            <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
                                                {{ perm.getDescription() }}
                                            </td>
                                        </tr>
                                    </tbody>
                                    <thead>
                                        <tr>
                                            <th scope="col"
                                                class="py-3.5 pl-3 pr-4 sm:pr-0 text-left text-sm font-semibold text-neutral">
                                                Actions
                                            </th>
                                            <th scope="col"
                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                                Name
                                            </th>
                                            <th scope="col"
                                                class="py-3.5 px-2 text-center text-sm font-semibold text-neutral">
                                                Description
                                            </th>
                                        </tr>
                                    </thead>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
