<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import AttrViewAttr from '~/components/rector/attrs/AttrViewAttr.vue';
import { useNotificatorStore } from '~/stores/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Permission, Role, RoleAttribute } from '~~/gen/ts/resources/permissions/permissions';
import type { AttrsUpdate, PermItem, PermsUpdate } from '~~/gen/ts/services/rector/rector';

const props = defineProps<{
    roleId: number;
}>();

const emit = defineEmits<{
    (e: 'deleted'): void;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { isSuperuser } = useAuth();

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

async function getRole(id: number): Promise<Role> {
    try {
        const call = $grpc.rector.rector.getRole({
            id: id,
            filtered: false,
        });
        const { response } = await call;

        if (response.role === undefined) {
            throw new Error('failed to get role from server response');
        }

        return response.role;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function getPermissions(roleId: number): Promise<void> {
    try {
        const call = $grpc.rector.rector.getPermissions({
            roleId,
            filtered: false,
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
    attrList.value.forEach((attr) => {
        // Make sure the permission is enabled, otherwise attr needs to be removed
        const perm = permStates.value.get(attr.permissionId);

        if (perm === undefined || attr.value === undefined) {
            attrs.toRemove.push({
                roleId: role.value!.id,
                attrId: attr.attrId,
                category: '',
                key: '',
                name: '',
                permissionId: attr.permissionId,
                type: '',
            });
        } else if (attr.value !== undefined) {
            attrs.toUpdate.push({
                roleId: role.value!.id,
                attrId: attr.attrId,
                maxValues: attr.maxValues,
                category: '',
                key: '',
                name: '',
                permissionId: attr.permissionId,
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
        await $grpc.rector.rector.updateRoleLimits({
            roleId: props.roleId,
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
}

async function initializeRoleView(): Promise<void> {
    clearState();

    await getPermissions(props.roleId);
    await propogatePermissionStates();
}

watch(role, async () => {
    initializeRoleView();
});

watch(props, async () => {
    if (!role.value || role.value?.id !== props.roleId) {
        refresh();
    }
});

async function copyRole(): Promise<void> {
    copyToClipboardWrapper(
        JSON.stringify({
            role: role.value,
            attrList: attrList.value,
        }),
    );

    notifications.add({
        title: { key: 'notifications.clipboard.copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.copied.content', parameters: {} },
        timeout: 3250,
        type: NotificationType.INFO,
    });
}

const schema = z.object({
    input: z.string(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    input: '',
});

type CopyRole = {
    role: Role;
    attrList: RoleAttribute[];
};

async function pasteRole(event: FormSubmitEvent<Schema>): Promise<void> {
    const parsed = JSON.parse(event.data.input) as CopyRole;

    if (parsed.attrList) {
        attrList.value.length = 0;
        attrList.value.push(...parsed.attrList);
    }

    state.input = '';
    changed.value = true;
}

const accordionCategories = computed(() =>
    [...permCategories.value.entries()].map((category) => {
        return {
            label: t(`perms.${category[1]}.category`),
            category: category[0],
        };
    }),
);

async function deleteFaction(id: number): Promise<void> {
    try {
        await $grpc.rector.rector.deleteFaction({
            roleId: id,
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
                        {{ role?.jobLabel! }}
                    </h2>

                    <UTooltip v-if="isSuperuser" :text="$t('common.delete')">
                        <UButton
                            variant="link"
                            icon="i-mdi-delete"
                            color="error"
                            @click="
                                modal.open(ConfirmModal, {
                                    confirm: async () => deleteFaction(role!.id),
                                })
                            "
                        />
                    </UTooltip>
                </div>

                <UDivider class="mb-1" :label="$t('common.attributes', 2)" />

                <div class="flex flex-col gap-2">
                    <div class="flex flex-row gap-1">
                        <UButton
                            class="flex-1"
                            :disabled="!changed || !canSubmit"
                            :loading="!canSubmit"
                            icon="i-mdi-content-save"
                            @click="onSubmitThrottle"
                        >
                            {{ $t('common.save', 1) }}
                        </UButton>

                        <UPopover>
                            <UButton :disabled="changed" color="gray" icon="i-mdi-form-textarea">
                                {{ $t('common.paste') }}
                            </UButton>

                            <template #panel>
                                <div class="p-4">
                                    <UForm class="flex flex-col gap-1" :state="state" :schema="schema" @submit="pasteRole">
                                        <UFormGroup name="input">
                                            <UInput v-model="state.input" type="text" name="input" />
                                        </UFormGroup>

                                        <UButton type="submit">
                                            {{ $t('common.save') }}
                                        </UButton>
                                    </UForm>
                                </div>
                            </template>
                        </UPopover>

                        <UButton icon="i-mdi-content-copy" :disabled="changed" color="white" @click="copyRole">
                            {{ $t('common.copy') }}
                        </UButton>
                    </div>

                    <UAccordion :items="accordionCategories" multiple :unmount="true">
                        <template #item="{ item: category }">
                            <div class="flex flex-col divide-y divide-gray-100 dark:divide-gray-800">
                                <div
                                    v-for="perm in permList.filter((p) => p.category === category.category)"
                                    :key="perm.id"
                                    class="flex flex-col gap-1"
                                >
                                    <div class="flex flex-row items-center gap-2">
                                        <div class="flex-1">
                                            <p class="text-gray-900 dark:text-white" :title="`${$t('common.id')}: ${perm.id}`">
                                                {{ $t(`perms.${perm.category}.${perm.name}.key`) }}
                                            </p>
                                            <p class="text-base-500">
                                                {{ $t(`perms.${perm.category}.${perm.name}.description`) }}
                                            </p>
                                        </div>

                                        <UButtonGroup class="inline-flex flex-initial">
                                            <UButton
                                                color="green"
                                                :variant="permStates.get(perm.id) ? 'solid' : 'soft'"
                                                icon="i-mdi-check"
                                                @click="updatePermissionState(perm.id, true)"
                                            />
                                            <UButton
                                                color="error"
                                                :variant="
                                                    permStates.get(perm.id) === undefined || !permStates.get(perm.id)
                                                        ? 'solid'
                                                        : 'soft'
                                                "
                                                icon="i-mdi-close"
                                                @click="updatePermissionState(perm.id, false)"
                                            />
                                        </UButtonGroup>
                                    </div>

                                    <template v-for="(attr, idx) in attrList" :key="attr.attrId">
                                        <AttrViewAttr
                                            v-if="attr.permissionId === perm.id"
                                            v-model="attrList[idx]!"
                                            :permission="perm"
                                            @changed="changed = true"
                                        />
                                    </template>
                                </div>
                            </div>
                        </template>
                    </UAccordion>
                </div>
            </template>
        </div>
    </div>
</template>
