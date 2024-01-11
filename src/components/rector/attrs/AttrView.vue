<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useConfirmDialog } from '@vueuse/core';
import { CheckIcon, ChevronDownIcon, CloseIcon, TrashCanIcon } from 'mdi-vue3';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericDivider from '~/components/partials/elements/GenericDivider.vue';
import { useNotificatorStore } from '~/store/notificator';
import { AttributeValues, Permission, Role, RoleAttribute } from '~~/gen/ts/resources/permissions/permissions';
import { AttrsUpdate, PermItem, PermsUpdate } from '~~/gen/ts/services/rector/rector';
import AttrViewAttr from '~/components/rector/attrs/AttrViewAttr.vue';

const props = defineProps<{
    roleId: string;
}>();

const emit = defineEmits<{
    (e: 'deleted'): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const { data: role, pending, refresh, error } = useLazyAsyncData(`rector-roles-${props.roleId}`, () => getRole(props.roleId));

const changed = ref(false);

const permList = ref<Permission[]>([]);
const permCategories = ref<Set<string>>(new Set());
const permStates = ref<Map<string, boolean | undefined>>(new Map());

const attrList = ref<RoleAttribute[]>([]);
const attrStates = ref<Map<string, AttributeValues | undefined>>(new Map());

async function getRole(id: string): Promise<Role> {
    try {
        const call = $grpc.getRectorClient().getRole({
            id,
            filtered: false,
        });
        const { response } = await call;

        if (response.role === undefined) {
            throw new Error('failed to get role from server response');
        }

        return response.role;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function deleteRole(id: string): Promise<void> {
    try {
        await $grpc.getRectorClient().deleteRole({ id });

        notifications.dispatchNotification({
            title: { key: 'notifications.rector.role_deleted.title', parameters: {} },
            content: { key: 'notifications.rector.role_deleted.content', parameters: {} },
            type: 'success',
        });

        emit('deleted');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function getPermissions(roleId: string): Promise<void> {
    try {
        const call = $grpc.getRectorClient().getPermissions({
            roleId,
            filtered: false,
        });
        const { response } = await call;

        permList.value = response.permissions;
        attrList.value = response.attributes;

        genPermissionCategories();
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
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

async function updatePermissionState(perm: string, state: boolean | undefined): Promise<void> {
    changed.value = true;
    permStates.value.set(perm, state);
}

async function updatePermissions(): Promise<void> {
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
                permissionId: '0',
                type: '',
            });
        } else if (state !== undefined) {
            attrs.toUpdate.push({
                roleId: role.value!.id,
                attrId: attr,
                maxValues: state,
                category: '',
                key: '',
                name: '',
                permissionId: '0',
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
        return;
    }

    try {
        await $grpc.getRectorClient().updateRoleLimits({
            roleId: props.roleId,
            perms,
            attrs,
        });

        notifications.dispatchNotification({
            title: { key: 'notifications.rector.role_updated.title', parameters: {} },
            content: { key: 'notifications.rector.role_updated.content', parameters: {} },
            type: 'success',
        });

        changed.value = false;
        refresh();
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
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
        attrStates.value.set(attr.attrId, attr.maxValues);
    });

    role.value?.attributes.forEach((attr) => {
        attrStates.value.set(attr.attrId, attr.maxValues);
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
            <DataNoDataBlock v-else-if="role === null" :type="$t('common.role', 2)" />
            <template v-else>
                <h2 class="text-3xl text-neutral" :title="`ID: ${role.id}`">
                    {{ role?.jobLabel! }}
                    <button v-if="can('RectorService.DeleteRole')" type="button" class="ml-1" @click="reveal()">
                        <TrashCanIcon class="w-6 h-6 mx-auto text-neutral" />
                    </button>
                </h2>
                <GenericDivider :label="$t('common.permission', 2)" />
                <div class="py-2 flex flex-col gap-4">
                    <button
                        type="button"
                        :disabled="!changed"
                        class="inline-flex px-3 py-2 text-center justify-center transition-colors font-semibold rounded-md text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                        :class="
                            !changed
                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500'
                        "
                        @click="updatePermissions()"
                    >
                        {{ $t('common.save', 1) }}
                    </button>

                    <Disclosure
                        v-for="category in permCategories"
                        :key="category"
                        v-slot="{ open }"
                        as="div"
                        class="text-neutral hover:border-neutral/70 border-neutral/20"
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
                                    :key="perm.id"
                                    class="flex flex-col gap-2"
                                >
                                    <div class="flex flex-row gap-4">
                                        <div class="flex flex-1 flex-col my-auto">
                                            <span class="truncate" :title="`${$t('common.id')}: ${perm.id}`">
                                                {{ $t(`perms.${perm.category}.${perm.name}.key`) }}
                                            </span>
                                            <span class="text-base-500 truncate">
                                                {{ $t(`perms.${perm.category}.${perm.name}.description`) }}
                                            </span>
                                        </div>
                                        <div class="flex flex-initial flex-row max-h-8 my-auto">
                                            <button
                                                :data-active="permStates.has(perm.id) ? permStates.get(perm.id) : false"
                                                class="transition-colors rounded-l-lg p-1 bg-success-600/50 data-[active=true]:bg-success-600 text-base-300 data-[active=true]:text-neutral hover:bg-success-600/70"
                                                @click="updatePermissionState(perm.id, true)"
                                            >
                                                <CheckIcon class="w-6 h-6" />
                                            </button>
                                            <button
                                                :data-active="permStates.get(perm.id) === undefined || !permStates.get(perm.id)"
                                                class="transition-colors rounded-r-lg p-1 bg-error-600/50 data-[active=true]:bg-error-600 text-base-300 data-[active=true]:text-neutral hover:bg-error-600/70"
                                                @click="updatePermissionState(perm.id, false)"
                                            >
                                                <CloseIcon class="w-6 h-6" />
                                            </button>
                                        </div>
                                    </div>
                                    <AttrViewAttr
                                        v-for="attr in attrList.filter((a) => a.permissionId === perm.id)"
                                        :key="attr.attrId"
                                        v-model:states="attrStates"
                                        :attribute="attr"
                                        :permission="perm"
                                        @changed="changed = true"
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
            </template>
        </div>
    </div>
</template>