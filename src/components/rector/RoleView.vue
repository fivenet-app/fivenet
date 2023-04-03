<script lang="ts" setup>
import { Permission, Role } from '@fivenet/gen/resources/permissions/permissions_pb';
import { RpcError } from 'grpc-web';
import { AddPermToRoleRequest, DeleteRoleRequest, GetPermissionsRequest, GetRoleRequest, RemovePermFromRoleRequest } from '@fivenet/gen/services/rector/rector_pb';
import { MagnifyingGlassIcon } from '@heroicons/vue/24/solid';
import DataPendingBlock from '../partials/DataPendingBlock.vue';
import DataErrorBlock from '../partials/DataErrorBlock.vue';
import { TrashIcon } from '@heroicons/vue/20/solid';
import Divider from '../partials/Divider.vue';
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { CheckIcon } from '@heroicons/vue/20/solid';
import { watchDebounced } from '@vueuse/shared';
import { dispatchNotification } from '../notification';

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

const { data: role, pending, refresh, error } = await useLazyAsyncData(`rector-role-${props.roleId}`, () => getRole());

async function deleteRole(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new DeleteRoleRequest();
        req.setId(props.roleId);

        try {
            return $grpc.getRectorClient().
                deleteRole(req, null);
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const entriesPerms = ref<Array<Permission>>([]);
const queryPerm = ref('');

async function getPermissions(): Promise<Array<Permission>> {
    return new Promise(async (res, rej) => {
        const req = new GetPermissionsRequest();
        req.setSearch(queryPerm.value);

        try {
            const resp = await $grpc.getRectorClient().
                getPermissions(req, null);

            return res(resp.getPermissionsList());
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const selectedPerm = ref<undefined | Permission>(undefined);
const permsToAdd = ref<Array<Permission>>([]);
const permsToRemove = ref<Array<Permission>>([]);

function addPermission(): void {
    if (!selectedPerm.value) {
        return;
    }

    // Remove perm from "to be removed" list
    const idx = permsToRemove.value.indexOf(selectedPerm.value);
    if (idx > -1) {
        permsToRemove.value.splice(idx, 1);
    }

    if (permsToAdd.value.indexOf(selectedPerm.value) === -1) {
        permsToAdd.value.push(selectedPerm.value);
    }

    const pIdx = role.value?.getPermissionsList().indexOf(selectedPerm.value) ?? -1;
    if (pIdx === -1) {
        role.value?.getPermissionsList().push(selectedPerm.value);
    }

    console.log(role.value?.getPermissionsList());
}

function removePermission(perm: Permission): void {
    // Remove perm from "to be added" list
    const idx = permsToAdd.value.indexOf(perm);
    if (idx > -1) {
        permsToAdd.value.splice(idx, 1);
    }

    if (permsToRemove.value.indexOf(perm) === -1) {
        permsToRemove.value.push(perm);
    }

    const pIdx = role.value?.getPermissionsList().indexOf(perm) ?? -1;
    if (pIdx > -1) {
        role.value?.getPermissionsList().splice(pIdx, 1);
    }

    console.log(role.value?.getPermissionsList());
}

async function saveAddPermissions(): Promise<void> {
    if (permsToAdd.value.length === 0) {
        return;
    }

    return new Promise(async (res, rej) => {
        const req = new AddPermToRoleRequest();
        req.setId(props.roleId);
        const permIds = new Array<number>();
        permsToAdd.value.forEach((v) => permIds.push(v.getId()));
        req.setPermissionsList(permIds);

        try {
            return $grpc.getRectorClient().
                addPermToRole(req, null);
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function saveRemovePermissions(): Promise<void> {
    if (permsToRemove.value.length === 0) {
        return;
    }

    return new Promise(async (res, rej) => {
        const req = new RemovePermFromRoleRequest();
        req.setId(props.roleId);
        const permIds = new Array<number>();
        permsToRemove.value.forEach((v) => permIds.push(v.getId()));
        req.setPermissionsList(permIds);

        try {
            return $grpc.getRectorClient().
                removePermFromRole(req, null);
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function saveRolePermissions(): Promise<void> {
    console.log("saveRolePermissions");
    console.log(permsToAdd.value);
    console.log(permsToRemove.value);
    await Promise.all([saveAddPermissions(), saveRemovePermissions()]);
    dispatchNotification({ title: 'Saving Role Permissions', content: '', type: 'success' });
}

onMounted(async () => {
    entriesPerms.value = await getPermissions();
});
watchDebounced(queryPerm, async () => await getPermissions(), { debounce: 750, maxWait: 1250 });
</script>

<template>
    <div v-if="role">
        <h2 class="text-3xl text-white">
            {{ toTitleCase(role?.getName()!) }}
            <button v-can="'RectorService.DeleteRole'" @click="deleteRole()">
                <TrashIcon class="w-6 h-6 mx-auto text-neutral" />
            </button>
        </h2>
        <p class="text-gray-400 text-sm">
            {{ role.getDescription() }}
        </p>
        <Divider label="Permissions" />
        <div class="py-2">
            <div class="px-2 sm:px-6 lg:px-8">
                <div v-can="'RectorService.AddPermToRole'" class="sm:flex-auto">
                    <form @submit.prevent="addPermission()">
                        <div class="flex flex-row gap-4 mx-auto">
                            <div class="flex-1 form-control">
                                <label for="owner"
                                    class="block text-sm font-medium leading-6 text-neutral">Permission</label>
                                <div class="relative items-center mt-2">
                                    <Combobox as="div" v-model="selectedPerm" nullable>
                                        <div class="relative">
                                            <ComboboxButton as="div">
                                                <ComboboxInput
                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    @change="queryPerm = $event.target.value"
                                                    :display-value="(perm: any) => perm ? `${perm?.getName()}: ${perm?.getDescription()}` : ''"
                                                    placeholder="Permission" />
                                            </ComboboxButton>

                                            <ComboboxOptions v-if="entriesPerms.length > 0"
                                                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                                                <ComboboxOption v-for="perm in entriesPerms" :key="perm?.getId()"
                                                    :value="perm" as="perm" v-slot="{ active, selected }">
                                                    <li
                                                        :class="['relative cursor-default select-none py-2 pl-8 pr-4 text-neutral', active ? 'bg-primary-500' : '']">
                                                        <span :class="['block truncate', selected && 'font-semibold']">
                                                            {{ perm?.getName() }}: {{ perm?.getDescription() }}
                                                        </span>

                                                        <span v-if="selected"
                                                            :class="[active ? 'text-neutral' : 'text-primary-500', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
                                                            <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                                        </span>
                                                    </li>
                                                </ComboboxOption>
                                            </ComboboxOptions>
                                        </div>
                                    </Combobox>
                                </div>
                            </div>
                            <div class="flex-initial form-control">
                                <button type="submit" :disabled="!selectedPerm"
                                    class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500">
                                    Add Permission
                                </button>
                            </div>
                            <div class="flex-initial form-control">
                                <button @click="saveRolePermissions()"
                                    :disabled="permsToAdd.length === 0 && permsToRemove.length === 0"
                                    class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-success-700 text-neutral hover:bg-success-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-success-700">
                                    Save Changes
                                </button>
                            </div>
                        </div>
                    </form>
                </div>
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
                                                    <button v-can="'RectorService.AddPermToRole'">
                                                        <TrashIcon class="w-6 h-6 mx-auto text-neutral"
                                                            @click="removePermission(perm)" />
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
