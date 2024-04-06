<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useNotificatorStore } from '~/store/notificator';
import { AttributeValues, Permission, Role, RoleAttribute } from '~~/gen/ts/resources/permissions/permissions';
import { AttrsUpdate, PermItem, PermsUpdate } from '~~/gen/ts/services/rector/rector';
import AttrViewAttr from '~/components/rector/attrs/AttrViewAttr.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';

const props = defineProps<{
    roleId: string;
}>();

const emit = defineEmits<{
    (e: 'deleted'): void;
}>();

const { t } = useI18n();

const { $grpc } = useNuxtApp();

const modal = useModal();

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

        notifications.add({
            title: { key: 'notifications.rector.role_deleted.title', parameters: {} },
            description: { key: 'notifications.rector.role_deleted.content', parameters: {} },
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

        notifications.add({
            title: { key: 'notifications.rector.role_updated.title', parameters: {} },
            description: { key: 'notifications.rector.role_updated.content', parameters: {} },
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

const accordionCategories = computed(() =>
    [...permCategories.value.entries()].map((category) => {
        return {
            label: t(`perms.${category[1]}.category`),
            category: category[0],
        };
    }),
);
</script>

<template>
    <div class="w-full">
        <div class="px-1 sm:px-2 lg:px-4">
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.role', 2)])" />
            <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.role', 2)])" :retry="refresh" />
            <DataNoDataBlock v-else-if="role === null" :type="$t('common.role', 2)" />

            <template v-else>
                <div class="flex justify-between">
                    <h2 class="text-3xl" :title="`ID: ${role.id}`">
                        {{ role?.jobLabel! }}
                    </h2>

                    <UButton
                        v-if="can('RectorService.DeleteRole')"
                        variant="link"
                        icon="i-mdi-trash-can"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteRole(role!.id),
                            })
                        "
                    />
                </div>

                <UDivider :label="$t('common.attributes', 2)" />

                <div class="flex flex-col gap-4 py-2">
                    <UButton :disabled="!changed" block @click="updatePermissions()">
                        {{ $t('common.save', 1) }}
                    </UButton>

                    <UAccordion :items="accordionCategories" multiple>
                        <template #item="{ item: category }">
                            <div class="flex flex-col gap-2 divide-y divide-gray-100 dark:divide-gray-800">
                                <div
                                    v-for="perm in permList.filter((p) => p.category === category.category)"
                                    :key="perm.id"
                                    class="flex flex-col gap-2"
                                >
                                    <div class="flex flex-row gap-4">
                                        <div class="my-auto flex flex-1 flex-col">
                                            <span :title="`${$t('common.id')}: ${perm.id}`">
                                                {{ $t(`perms.${perm.category}.${perm.name}.key`) }}
                                            </span>
                                            <span class="text-base-500">
                                                {{ $t(`perms.${perm.category}.${perm.name}.description`) }}
                                            </span>
                                        </div>
                                        <UButtonGroup class="my-auto flex max-h-8 flex-initial flex-row">
                                            <UButton
                                                :disabled="permStates.has(perm.id) ? permStates.get(perm.id) : false"
                                                color="green"
                                                icon="i-mdi-check"
                                                @click="updatePermissionState(perm.id, true)"
                                            />
                                            <UButton
                                                :disabled="
                                                    permStates.has(perm.id)
                                                        ? permStates.get(perm.id) !== undefined && !permStates.get(perm.id)
                                                        : false
                                                "
                                                color="red"
                                                icon="i-mdi-close"
                                                @click="updatePermissionState(perm.id, false)"
                                            />
                                        </UButtonGroup>
                                    </div>

                                    <AttrViewAttr
                                        v-for="attr in attrList.filter((a) => a.permissionId === perm.id)"
                                        :key="attr.attrId"
                                        v-model:states="attrStates"
                                        :attribute="attr"
                                        :permission="perm"
                                        @changed="changed = true"
                                    />
                                </div>
                            </div>
                        </template>
                    </UAccordion>
                </div>
            </template>
        </div>
    </div>
</template>
