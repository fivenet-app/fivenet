<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { useConfirmDialog } from '@vueuse/core';
import { CheckIcon, ChevronDownIcon, CloseIcon, MinusIcon, TrashCanIcon } from 'mdi-vue3';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Divider from '~/components/partials/elements/Divider.vue';
import RoleViewAttr from '~/components/rector/roles/RoleViewAttr.vue';
import { useNotificatorStore } from '~/store/notificator';
import { AttributeValues, Permission, Role, RoleAttribute } from '~~/gen/ts/resources/permissions/permissions';
import { AttrsUpdate, PermItem, PermsUpdate } from '~~/gen/ts/services/rector/rector';

const props = defineProps<{
    roleId: bigint;
}>();

const emit = defineEmits<{
    (e: 'deleted'): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const {
    data: role,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`rector-roles-${props.roleId.toString()}`, () => getRole(props.roleId));

const changed = ref(false);

const permList = ref<Permission[]>([]);
const permCategories = ref<Set<string>>(new Set());
const permStates = ref<Map<bigint, boolean | undefined>>(new Map());

const attrList = ref<RoleAttribute[]>([]);
const attrStates = ref<Map<bigint, AttributeValues | undefined>>(new Map());

async function getRole(id: bigint): Promise<Role> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getRectorClient().getRole({
                id: id,
                filtered: true,
            });
            const { response } = await call;

            if (response.role === undefined) return rej();

            return res(response.role!);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function deleteRole(id: bigint): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            await $grpc.getRectorClient().deleteRole({
                id: id,
            });

            notifications.dispatchNotification({
                title: { key: 'notifications.rector.role_deleted.title', parameters: {} },
                content: { key: 'notifications.rector.role_deleted.content', parameters: {} },
                type: 'success',
            });

            emit('deleted');
            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function getPermissions(roleId: bigint): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getRectorClient().getPermissions({
                roleId: roleId,
                filtered: true,
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
    changed.value = true;
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
            // Make sure the permission exists and is enabled, otherwise attr needs to be removed
            const a = attrList.value.find((a) => a.attrId === attr);
            if (a === undefined) {
                return;
            }
            const perm = permStates.value.get(a.permissionId);

            if (perm === undefined || state === undefined) {
                attrs.toRemove.push({
                    roleId: role.value!.id,
                    attrId: attr,
                    category: '',
                    key: '',
                    name: '',
                    permissionId: 0n,
                    type: '',
                });
            } else if (state !== undefined) {
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
            }
        });

        if (
            perms.toUpdate.length === 0 &&
            perms.toRemove.length === 0 &&
            attrs.toUpdate.length === 0 &&
            attrs.toRemove.length === 0
        ) {
            changed.value = false;
            return res();
        }

        try {
            await $grpc.getRectorClient().updateRolePerms({
                id: props.roleId,
                perms: perms,
                attrs: attrs,
            });

            notifications.dispatchNotification({
                title: { key: 'notifications.rector.role_updated.title', parameters: {} },
                content: { key: 'notifications.rector.role_updated.content', parameters: {} },
                type: 'success',
            });

            changed.value = false;
            refresh();

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

function clearState(): void {
    changed.value = false;
    permList.value.length = 0;
    permCategories.value.clear();
    permStates.value.clear();
    attrList.value.length = 0;
    attrStates.value.clear();
}

async function initializeRoleView(): Promise<void> {
    clearState();

    await getPermissions(props.roleId);
    await propogatePermissionStates();

    attrStates.value.clear();
    attrList.value.forEach((attr) => {
        attrStates.value.set(attr.attrId, attr.value);
    });

    role.value?.attributes.forEach((attr) => {
        attrStates.value.set(attr.attrId, attr.value);
    });
}

watch(role, async () => {
    initializeRoleView();
});

watch(props, () => {
    if (role.value === null || role.value.id !== props.roleId) {
        refresh();
    }
});

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();

onConfirm(async (id) => deleteRole(id));
</script>

<template>
    <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(role!.id)" />

    <div class="py-4 w-full">
        <div class="px-1 sm:px-2 lg:px-4">
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.role', 2)])" />
            <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.role', 2)])" :retry="refresh" />
            <DataNoDataBlock v-else-if="!role" :type="$t('common.role', 2)" />
            <div v-else>
                <h2 class="text-3xl text-neutral" :title="`ID: ${role.id}`">
                    {{ role?.jobLabel! }} - {{ role?.jobGradeLabel }} ({{ role.grade }})
                    <button v-if="can('RectorService.DeleteRole')" type="button" @click="reveal()" class="ml-1">
                        <TrashCanIcon class="w-6 h-6 mx-auto text-neutral" />
                    </button>
                </h2>
                <Divider label="Permissions" />
                <div class="py-2 flex flex-col gap-4">
                    <button
                        type="button"
                        @click="updatePermissions()"
                        :disabled="!changed"
                        class="inline-flex px-3 py-2 text-center justify-center transition-colors font-semibold rounded-md text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                        :class="
                            !changed
                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500'
                        "
                    >
                        {{ $t('common.save', 1) }}
                    </button>

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
                            <div class="flex flex-col gap-2 mx-auto my-2">
                                <div
                                    v-for="(perm, idx) in permList.filter((p) => p.category === category)"
                                    :key="perm.id?.toString()"
                                    class="flex flex-col gap-2"
                                >
                                    <div class="flex flex-row gap-4">
                                        <div class="flex flex-1 flex-col my-auto">
                                            <span class="truncate" :title="`${$t('common.id')}: ${perm.id.toString()}`">
                                                {{ $t(`perms.${perm.category}.${perm.name}.key`) }}
                                            </span>
                                            <span class="text-base-500 truncate">
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
                                        @changed="changed = true"
                                        :disabled="permStates.get(perm.id) !== true"
                                    />
                                    <div
                                        v-if="idx !== permList.filter((p) => p.category === category).length - 1"
                                        class="w-full border-t border-neutral/20"
                                    />
                                </div>
                            </div>
                        </DisclosurePanel>
                    </Disclosure>
                </div>
            </div>
        </div>
    </div>
</template>
