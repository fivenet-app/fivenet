<script lang="ts" setup>
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import RoleViewAttr from '~/components/rector/roles/RoleViewAttr.vue';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { AttributeValues, Permission, Role, RoleAttribute } from '~~/gen/ts/resources/permissions/permissions';
import type { AttrsUpdate, PermItem, PermsUpdate } from '~~/gen/ts/services/rector/rector';

const props = defineProps<{
    roleId: number;
}>();

const emit = defineEmits<{
    (e: 'deleted'): void;
}>();

const { t } = useI18n();

const { can } = useAuth();

const modal = useModal();

const notifications = useNotificatorStore();

const {
    data: role,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`rector-roles-${props.roleId}`, () => getRole(props.roleId));

const changed = ref(false);

const permList = ref<Permission[]>([]);
const permCategories = ref<Set<string>>(new Set());
const permStates = ref(new Map<number, boolean | undefined>());

const attrList = ref<RoleAttribute[]>([]);
const attrStates = ref(new Map<number, AttributeValues | undefined>());

async function getRole(id: number): Promise<Role> {
    try {
        const call = getGRPCRectorClient().getRole({
            id,
            filtered: true,
        });
        const { response } = await call;

        return response.role!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteRole(id: number): Promise<void> {
    try {
        await getGRPCRectorClient().deleteRole({
            id,
        });

        notifications.add({
            title: { key: 'notifications.rector.role_deleted.title', parameters: {} },
            description: { key: 'notifications.rector.role_deleted.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('deleted');
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function getPermissions(roleId: number): Promise<void> {
    try {
        const call = getGRPCRectorClient().getPermissions({
            roleId,
            filtered: true,
        });
        const { response } = await call;

        permList.value = response.permissions;
        attrList.value = response.attributes;

        genPermissionCategories();
    } catch (e) {
        handleGRPCError(e as RpcError);
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

async function updatePermissionState(perm: number, state: boolean | undefined): Promise<void> {
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
                permissionId: 0,
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
                permissionId: 0,
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
        await getGRPCRectorClient().updateRolePerms({
            id: props.roleId,
            perms: perms,
            attrs: attrs,
        });

        notifications.add({
            title: { key: 'notifications.rector.role_updated.title', parameters: {} },
            description: { key: 'notifications.rector.role_updated.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        changed.value = false;
        refresh();
    } catch (e) {
        handleGRPCError(e as RpcError);
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

    attrList.value.forEach((attr) => {
        attrStates.value.set(attr.attrId, attr.value);
    });

    role.value?.attributes.forEach((attr) => {
        attrStates.value.set(attr.attrId, attr.value);
    });
}

watch(role, async () => initializeRoleView());

watch(props, () => {
    if (!role.value || role.value?.id !== props.roleId) {
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

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await updatePermissions().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div class="w-full">
        <div class="px-1 sm:px-2">
            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.role', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.role', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!role" :type="$t('common.role', 2)" />

            <template v-else>
                <div class="flex justify-between">
                    <h2 class="text-3xl" :title="`ID: ${role.id}`">
                        {{ role?.jobLabel! }} - {{ role?.jobGradeLabel }} ({{ role.grade }})
                    </h2>

                    <UButton
                        v-if="can('RectorService.DeleteRole').value"
                        variant="link"
                        icon="i-mdi-trash-can"
                        color="red"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteRole(role!.id),
                            })
                        "
                    />
                </div>

                <UDivider :label="$t('common.permission', 2)" />

                <div class="flex flex-col gap-4 py-2">
                    <UButton
                        v-if="can('RectorService.UpdateRolePerms').value"
                        block
                        :disabled="!changed || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle"
                    >
                        {{ $t('common.save', 1) }}
                    </UButton>

                    <UAccordion :items="accordionCategories" multiple :unmount="true">
                        <template #item="{ item: category }">
                            <div class="flex flex-col gap-2 divide-y divide-gray-100 dark:divide-gray-800">
                                <div
                                    v-for="perm in permList.filter((p) => p.category === category.category)"
                                    :key="perm.id"
                                    class="flex flex-col gap-2"
                                >
                                    <div class="flex flex-row gap-2">
                                        <div class="my-auto flex flex-1 flex-col">
                                            <span
                                                :title="`${$t('common.id')}: ${perm.id}`"
                                                class="text-gray-900 dark:text-white"
                                            >
                                                {{ $t(`perms.${perm.category}.${perm.name}.key`) }}
                                            </span>
                                            <span class="text-base-500">
                                                {{ $t(`perms.${perm.category}.${perm.name}.description`) }}
                                            </span>
                                        </div>

                                        <UButtonGroup class="inline-flex flex-initial">
                                            <UButton
                                                color="green"
                                                :variant="permStates.get(perm.id) ? 'solid' : 'soft'"
                                                icon="i-mdi-check"
                                                @click="updatePermissionState(perm.id, true)"
                                            />

                                            <UButton
                                                color="black"
                                                :variant="
                                                    !permStates.has(perm.id) || permStates.get(perm.id) === undefined
                                                        ? 'solid'
                                                        : 'soft'
                                                "
                                                icon="i-mdi-minus"
                                                @click="updatePermissionState(perm.id, undefined)"
                                            />

                                            <UButton
                                                color="red"
                                                :variant="
                                                    permStates.get(perm.id) !== undefined && !permStates.get(perm.id)
                                                        ? 'solid'
                                                        : 'soft'
                                                "
                                                icon="i-mdi-close"
                                                @click="updatePermissionState(perm.id, false)"
                                            />
                                        </UButtonGroup>
                                    </div>

                                    <RoleViewAttr
                                        v-for="attr in attrList.filter((a) => a.permissionId === perm.id)"
                                        :key="attr.attrId"
                                        v-model:states="attrStates"
                                        :attribute="attr"
                                        :permission="perm"
                                        :disabled="permStates.get(perm.id) !== true"
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
