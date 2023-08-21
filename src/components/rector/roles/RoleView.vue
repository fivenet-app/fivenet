<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';

import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { CheckIcon, ChevronDownIcon, CloseIcon, MinusIcon, TrashCanIcon } from 'mdi-vue3';
import Divider from '~/components/partials/elements/Divider.vue';
import RoleViewAttr from '~/components/rector/roles/RoleViewAttr.vue';
import { useNotificationsStore } from '~/store/notifications';
import { AttributeValues, Permission, Role, RoleAttribute } from '~~/gen/ts/resources/permissions/permissions';
import { Job } from '~~/gen/ts/resources/users/jobs';
import { AttrsUpdate, PermItem, PermsUpdate } from '~~/gen/ts/services/rector/rector';

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

const props = defineProps<{
    roleId: bigint;
}>();

const role = ref<Role>();

const permList = ref<Permission[]>([]);
const permCategories = ref<Set<string>>(new Set());
const permStates = ref<Map<bigint, boolean | undefined>>(new Map());

const attrList = ref<RoleAttribute[]>([]);
const attrStates = ref<Map<bigint, AttributeValues | undefined>>(new Map());

const jobs = ref<Job[]>([]);

async function getRole(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getRectorClient().getRole({
                id: props.roleId,
            });
            const { response } = await call;

            role.value = response.role;

            attrStates.value.clear();
            role.value?.attributes.forEach((attr) => {
                attrStates.value.set(attr.attrId, attr.value);
            });

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function prepareAttributeData(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCompletorClient().completeJobs({});
            const { response } = await call;

            jobs.value = response.jobs;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function deleteRole(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            await $grpc.getRectorClient().deleteRole({
                id: props.roleId,
            });

            notifications.dispatchNotification({
                title: { key: 'notifications.rector.role_deleted.title', parameters: [] },
                content: { key: 'notifications.rector.role_deleted.content', parameters: [] },
                type: 'success',
            });

            await navigateTo({ name: 'rector-roles' });
            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function getPermissions(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getRectorClient().getPermissions({
                roleId: props.roleId,
            });
            const { response } = await call;

            permList.value = response.permissions;
            attrList.value = response.attributes;

            genPermissionCategories();

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function genPermissionCategories(): Promise<void> {
    permCategories.value.clear();

    permList.value.forEach((perm) => {
        permCategories.value.add(perm.category);
    });
}

async function propogatePermissionStates(): Promise<void> {
    permStates.value.clear();

    role.value?.permissions.forEach((perm) => {
        permStates.value.set(perm.id, Boolean(perm.val));
    });
}

async function updatePermissionState(perm: bigint, state: boolean | undefined): Promise<void> {
    permStates.value.set(perm, state);
}

async function updatePermissions(): Promise<void> {
    return new Promise(async (res, rej) => {
        const currentPermissions = role.value?.permissions.map((p) => p.id) ?? [];

        const perms: PermsUpdate = {
            toRemove: [],
            toUpdate: [],
        };
        permStates.value.forEach((state, perm) => {
            if (state !== undefined) {
                const p = role.value?.permissions.find((v) => v.id === perm);

                if (p?.val !== state) {
                    const item: PermItem = {
                        id: perm,
                        val: state,
                    };

                    perms.toUpdate.push(item);
                }
            } else if (state === undefined && currentPermissions.includes(perm)) {
                perms.toRemove.push(perm);
            }
        });

        const attrs: AttrsUpdate = {
            toRemove: [],
            toUpdate: [],
        };
        attrStates.value.forEach((state, attr) => {
            if (state !== undefined) {
                attrs.toUpdate.push({
                    roleId: role.value!.id,
                    attrId: attr,
                    value: state,
                    category: '',
                    key: '',
                    name: '',
                    permissionId: 0n,
                    type: '',
                });
            } else if (state === undefined) {
                attrs.toRemove.push({
                    roleId: role.value!.id,
                    attrId: attr,
                    category: '',
                    key: '',
                    name: '',
                    permissionId: 0n,
                    type: '',
                });
            }
        });

        if (
            perms.toUpdate.length === 0 &&
            perms.toRemove.length === 0 &&
            attrs.toUpdate.length === 0 &&
            attrs.toRemove.length === 0
        )
            return res();

        try {
            await $grpc.getRectorClient().updateRolePerms({
                id: props.roleId,
                perms: perms,
                attrs: attrs,
            });

            notifications.dispatchNotification({
                title: { key: 'notifications.rector.role_updated.title', parameters: [] },
                content: { key: 'notifications.rector.role_updated.content', parameters: [] },
                type: 'success',
            });

            initializeRoleView();

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function initializeRoleView(): Promise<void> {
    await getRole();
    await prepareAttributeData();
    await getPermissions();
    await propogatePermissionStates();
}

onMounted(async () => {
    initializeRoleView();
});
</script>

<template>
    <div class="py-4 max-w-7xl mx-auto">
        <div class="px-1 sm:px-2 lg:px-4">
            <div v-if="role">
                <h2 class="text-3xl text-white">
                    {{ role?.jobLabel! }} - {{ role?.jobGradeLabel }}
                    <button v-if="can('RectorService.DeleteRole')" @click="deleteRole()">
                        <TrashCanIcon class="w-6 h-6 mx-auto text-neutral" />
                    </button>
                </h2>
                <Divider label="Permissions" />
                <div class="py-2 flex flex-col gap-4">
                    <Disclosure
                        as="div"
                        class="text-neutral hover:border-neutral/70 border-neutral/20"
                        v-for="category in permCategories"
                        :key="category"
                        v-slot="{ open }"
                    >
                        <DisclosureButton
                            :class="[
                                open ? 'rounded-t-lg border-b-0' : 'rounded-lg',
                                'flex w-full items-start justify-between text-left border-2 p-2 border-inherit transition-colors',
                            ]"
                        >
                            <span class="text-base font-semibold leading-7">
                                {{ $t(`perms.${category}.category`) }}
                            </span>
                            <span class="ml-6 flex h-7 items-center">
                                <ChevronDownIcon
                                    :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                                    aria-hidden="true"
                                />
                            </span>
                        </DisclosureButton>
                        <DisclosurePanel
                            class="px-4 pb-2 border-2 border-t-0 rounded-b-lg transition-colors border-inherit -mt-2"
                        >
                            <div class="flex flex-col gap-2 max-w-4xl mx-auto my-2">
                                <div
                                    v-for="(perm, idx) in permList.filter((p) => p.category === category)"
                                    :key="perm.id?.toString()"
                                    class="flex flex-col gap-2"
                                >
                                    <div class="flex flex-row gap-4">
                                        <div class="flex flex-1 flex-col my-auto">
                                            <span class="truncate lg:max-w-full max-w-xs">
                                                {{ $t(`perms.${perm.category}.${perm.name}.key`) }}
                                            </span>
                                            <span class="text-base-500 truncate lg:max-w-full max-w-xs">
                                                {{ $t(`perms.${perm.category}.${perm.name}.description`) }}
                                            </span>
                                        </div>
                                        <div class="flex flex-initial flex-row max-h-8 my-auto">
                                            <button
                                                :data-active="permStates.has(perm.id) ? permStates.get(perm.id) : false"
                                                @click="updatePermissionState(perm.id, true)"
                                                class="transition-colors rounded-l-lg p-1 bg-success-600/50 data-[active=true]:bg-success-600 text-base-300 data-[active=true]:text-neutral hover:bg-success-600/70"
                                            >
                                                <CheckIcon class="w-6 h-6" />
                                            </button>
                                            <button
                                                :data-active="!permStates.has(perm.id) || permStates.get(perm.id) === undefined"
                                                @click="updatePermissionState(perm.id, undefined)"
                                                class="transition-colors p-1 bg-base-700 data-[active=true]:bg-base-500 text-base-300 data-[active=true]:text-neutral hover:bg-base-600"
                                            >
                                                <MinusIcon class="w-6 h-6" />
                                            </button>
                                            <button
                                                :data-active="
                                                    permStates.has(perm.id)
                                                        ? permStates.get(perm.id) !== undefined && !permStates.get(perm.id)
                                                        : false
                                                "
                                                @click="updatePermissionState(perm.id, false)"
                                                class="transition-colors rounded-r-lg p-1 bg-error-600/50 data-[active=true]:bg-error-600 text-base-300 data-[active=true]:text-neutral hover:bg-error-600/70"
                                            >
                                                <CloseIcon class="w-6 h-6" />
                                            </button>
                                        </div>
                                    </div>
                                    <RoleViewAttr
                                        v-for="attr in attrList.filter((a) => a.permissionId === perm.id)"
                                        :attribute="attr"
                                        :permission="perm"
                                        v-model:states="attrStates"
                                        :disabled="permStates.get(perm.id) !== true"
                                        :jobs="jobs"
                                    />
                                    <div
                                        v-if="idx !== permList.filter((p) => p.category === category).length - 1"
                                        class="w-full border-t border-neutral/20"
                                    />
                                </div>
                            </div>
                        </DisclosurePanel>
                    </Disclosure>
                    <button
                        type="button"
                        @click="updatePermissions()"
                        class="inline-flex px-3 py-2 text-center justify-center transition-colors font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-success-600"
                    >
                        {{ $t('common.save', 1) }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
