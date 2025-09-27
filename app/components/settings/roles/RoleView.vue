<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import RoleViewAttr from '~/components/settings/roles/RoleViewAttr.vue';
import { getSettingsSettingsClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { RoleAttribute } from '~~/gen/ts/resources/permissions/attributes';
import type { Permission, Role } from '~~/gen/ts/resources/permissions/permissions';
import type { AttrsUpdate, PermsUpdate } from '~~/gen/ts/resources/settings/perms';
import EffectivePermsSlideover from './EffectivePermsSlideover.vue';
import { isEmptyAttributes } from './helpers';

const props = defineProps<{
    roleId: number;
}>();

const emit = defineEmits<{
    (e: 'deleted'): void;
}>();

const { t } = useI18n();

const { can } = useAuth();

const overlay = useOverlay();

const notifications = useNotificationsStore();

const settingsSettingsClient = await getSettingsSettingsClient();

const { data: role, status, refresh, error } = useLazyAsyncData(`settings-roles-${props.roleId}`, () => getRole(props.roleId));

const changed = ref(false);

const permList = ref<Permission[]>([]);
const permCategories = ref<Map<string, { category: string; icon: string | undefined; order: number }>>(new Map());
const permStates = ref(new Map<number, boolean | undefined>());

const attrList = ref<RoleAttribute[]>([]);

async function getRole(id: number): Promise<Role> {
    try {
        const call = settingsSettingsClient.getRole({
            id: id,
        });
        const { response } = await call;

        clearState();

        await getPermissions(id);

        await propogateRolePermissionStates(response.role!);

        return response.role!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function getPermissions(roleId: number): Promise<void> {
    try {
        const call = settingsSettingsClient.getPermissions({
            roleId: roleId,
        });
        const { response } = await call;

        permList.value = response.permissions;
        attrList.value = response.attributes;

        await genPermissionCategories();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteRole(id: number): Promise<void> {
    try {
        await settingsSettingsClient.deleteRole({
            id: id,
        });

        notifications.add({
            title: { key: 'notifications.settings.role_deleted.title', parameters: {} },
            description: { key: 'notifications.settings.role_deleted.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('deleted');
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function genPermissionCategories(): Promise<void> {
    permCategories.value.clear();

    permList.value.forEach((perm) => {
        if (permCategories.value.has(perm.category)) return;

        permCategories.value.set(perm.category, {
            category: perm.category,
            icon: perm.icon,
            order: perm.order ?? 999999,
        });
    });
}

async function propogateRolePermissionStates(role: Role): Promise<void> {
    permStates.value.clear();

    role.permissions.forEach((perm) => permStates.value.set(perm.id, Boolean(perm.val)));

    role.attributes.forEach((attr) => {
        const idx = attrList.value.findIndex((a) => a.attrId === attr.attrId);
        if (idx > -1 && attrList.value[idx]) {
            attrList.value[idx].value = attr.value;
        } else {
            attrList.value.push(attr);
        }
    });
}

async function updatePermissionState(perm: number, state: boolean | undefined): Promise<void> {
    changed.value = true;
    permStates.value.set(perm, state);
}

async function updateRolePerms(): Promise<void> {
    const currentPermissions = role.value?.permissions.map((p) => p.id) ?? [];

    const perms: PermsUpdate = {
        toRemove: [],
        toUpdate: [],
    };
    permStates.value.forEach((state, perm) => {
        if (state !== undefined) {
            const p = role.value?.permissions.find((v) => v.id === perm);

            if (p?.val !== state) {
                perms.toUpdate.push({
                    id: perm,
                    val: state,
                });
            }
        } else if (state === undefined && currentPermissions.includes(perm)) {
            perms.toRemove.push({
                id: perm,
                val: false,
            });
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
                value: attr.value,
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
        await settingsSettingsClient.updateRolePerms({
            id: props.roleId,
            perms: perms,
            attrs: attrs,
        });

        notifications.add({
            title: { key: 'notifications.settings.role_updated.title', parameters: {} },
            description: { key: 'notifications.settings.role_updated.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        changed.value = false;
        await refresh();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function clearState(): void {
    changed.value = false;
    permList.value.length = 0;
    attrList.value.length = 0;
    permCategories.value.clear();
    permStates.value.clear();
}

watch(props, async () => {
    if (!role.value || role.value?.id !== props.roleId) {
        refresh();
    }
});

type CopyRole = {
    roleId: number;
    attrList: RoleAttribute[];
    permStates: { id: number; val?: boolean | undefined }[];
};

async function copyRole(): Promise<void> {
    copyToClipboardWrapper(
        JSON.stringify({
            roleId: props.roleId,
            attrList: attrList.value.map((a) => ({
                attrId: a.attrId,
                permissionId: a.permissionId,
                value: a.value,
            })),
            permStates: [...permStates.value.entries()].map((v) => ({ id: v[0], val: v[1] })),
        } as CopyRole),
    );

    notifications.add({
        title: { key: 'notifications.clipboard.copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.copied.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

const schema = z.object({
    input: z.coerce.string(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    input: '',
});

async function pasteRole(event: FormSubmitEvent<Schema>): Promise<void> {
    let parsed: CopyRole;
    try {
        parsed = JSON.parse(event.data.input) as CopyRole;
    } catch (e) {
        console.error('Failed to parse role data from clipboard', e);
        notifications.add({
            title: { key: 'notifications.action_failed.title', parameters: {} },
            description: { key: 'notifications.action_failed.content', parameters: {} },
            type: NotificationType.ERROR,
        });
        return;
    }

    if (parsed.attrList) {
        parsed.attrList?.forEach((a) => {
            if (a.value?.validValues.oneofKind === undefined) return;

            const at = attrList.value.find((at) => at.attrId === a.attrId);
            if (at) {
                if (at.maxValues?.validValues.oneofKind === 'stringList' && a.value?.validValues.oneofKind === 'stringList') {
                    a.value.validValues.stringList.strings = a.value.validValues.stringList.strings.filter(
                        (s) =>
                            at.maxValues?.validValues.oneofKind === 'stringList' &&
                            at.maxValues.validValues.stringList.strings.includes(s),
                    );
                } else if (at.maxValues?.validValues.oneofKind === 'jobList' && a.value?.validValues.oneofKind === 'jobList') {
                    a.value.validValues.jobList.strings = a.value.validValues.jobList.strings.filter(
                        (s) =>
                            at.maxValues?.validValues.oneofKind === 'jobList' &&
                            at.maxValues.validValues.jobList.strings.includes(s),
                    );
                } else if (
                    at.maxValues?.validValues.oneofKind === 'jobGradeList' &&
                    a.value?.validValues.oneofKind === 'jobGradeList'
                ) {
                    if (!a.value.validValues.jobGradeList.fineGrained) {
                        const jobs = Object.keys(a.value.validValues.jobGradeList.jobs);
                        const maxJobs = Object.keys(at.maxValues.validValues.jobGradeList.jobs);

                        jobs.filter((j) => !maxJobs.includes(j)).forEach(
                            (j) =>
                                a.value?.validValues.oneofKind === 'jobGradeList' &&
                                // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
                                delete a.value.validValues.jobGradeList.jobs[j],
                        );
                    } else {
                        const jobs = Object.keys(a.value.validValues.jobGradeList.jobs);
                        const maxJobs = Object.keys(at.maxValues.validValues.jobGradeList.jobs);

                        jobs.filter((j) => !maxJobs.includes(j)).forEach(
                            (j) =>
                                a.value?.validValues.oneofKind === 'jobGradeList' &&
                                // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
                                delete a.value.validValues.jobGradeList.jobs[j],
                        );
                    }
                }

                at.value = a.value;
            } else {
                attrList.value.push(a);
            }
        });
    }

    parsed.permStates?.forEach((p) => {
        const pe = permList.value.find((pe) => pe.id === p.id);
        if (pe) {
            permStates.value.set(p.id, p.val);
        }
    });

    state.input = '';
    changed.value = true;
}

const accordionCategories = computed(() =>
    [...permCategories.value.entries()]
        .map((category) => {
            return {
                label: t(`perms.${category[1].category}.category`),
                category: category[1].category,
                icon: category[1].icon,
                order: category[1].order,
            };
        })
        .sort((a, b) => a.order - b.order),
);

const canUpdate = can('settings.SettingsService/UpdateRolePerms');

const effectivePermsSlideover = overlay.create(EffectivePermsSlideover, { props: { roleId: props.roleId } });
const confirmModal = overlay.create(ConfirmModal);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await updateRolePerms().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div class="w-full">
        <div class="px-1 sm:px-2">
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.role', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.role', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!role" :type="$t('common.role', 2)" :retry="refresh" />

            <template v-else>
                <div class="flex justify-between">
                    <h2 class="line-clamp-2 flex-1 text-3xl" :title="`ID: ${role.id}`">
                        {{ role?.jobLabel }} - {{ role?.jobGradeLabel }} ({{ role.grade }})
                    </h2>

                    <UButtonGroup>
                        <UTooltip :text="$t('common.effective_permissions')">
                            <UButton
                                variant="link"
                                icon="i-mdi-account-key"
                                color="primary"
                                @click="
                                    effectivePermsSlideover.open({
                                        roleId: role!.id,
                                    })
                                "
                            />
                        </UTooltip>

                        <UTooltip :text="$t('common.refresh')">
                            <UButton variant="link" icon="i-mdi-refresh" color="primary" @click="refresh()" />
                        </UTooltip>

                        <UTooltip :text="$t('common.delete')">
                            <UButton
                                v-if="can('settings.SettingsService/DeleteRole').value"
                                variant="link"
                                icon="i-mdi-delete"
                                color="error"
                                @click="
                                    confirmModal.open({
                                        confirm: async () => deleteRole(role!.id),
                                    })
                                "
                            />
                        </UTooltip>
                    </UButtonGroup>
                </div>

                <USeparator class="mb-1" :label="$t('common.permission', 2)" />

                <div class="flex flex-col gap-2">
                    <div class="flex flex-row gap-1">
                        <template v-if="canUpdate">
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
                                <UButton :disabled="changed" color="neutral" icon="i-mdi-form-textarea">
                                    {{ $t('common.paste') }}
                                </UButton>

                                <template #content>
                                    <div class="p-4">
                                        <UForm class="flex flex-col gap-1" :state="state" :schema="schema" @submit="pasteRole">
                                            <UFormField name="input">
                                                <UInput v-model="state.input" type="text" name="input" />
                                            </UFormField>

                                            <UButton type="submit" :label="$t('common.save')" />
                                        </UForm>
                                    </div>
                                </template>
                            </UPopover>
                        </template>
                        <span v-else class="flex-1" />

                        <UButton icon="i-mdi-content-copy" :disabled="changed" color="neutral" @click="copyRole">
                            {{ $t('common.copy') }}
                        </UButton>
                    </div>

                    <UAccordion :items="accordionCategories" type="multiple">
                        <template #content="{ item: category }">
                            <div class="flex flex-col divide-y divide-default">
                                <div
                                    v-for="perm in permList.filter((p) => p.category === category.category)"
                                    :key="perm.id"
                                    class="flex flex-col gap-1 py-1"
                                >
                                    <UFormField
                                        class="flex flex-1 flex-row items-center gap-2"
                                        :label="$t(`perms.${perm.category}.${perm.name}.key`)"
                                        :description="$t(`perms.${perm.category}.${perm.name}.description`)"
                                        :ui="{ wrapper: 'flex-1' }"
                                    >
                                        <UButtonGroup class="inline-flex flex-initial">
                                            <UButton
                                                color="green"
                                                :variant="permStates.get(perm.id) ? 'solid' : 'soft'"
                                                icon="i-mdi-check"
                                                :disabled="!canUpdate"
                                                @click="updatePermissionState(perm.id, true)"
                                            />

                                            <UButton
                                                color="neutral"
                                                :variant="
                                                    !permStates.has(perm.id) || permStates.get(perm.id) === undefined
                                                        ? 'solid'
                                                        : 'soft'
                                                "
                                                icon="i-mdi-minus"
                                                :disabled="!canUpdate"
                                                @click="updatePermissionState(perm.id, undefined)"
                                            />

                                            <UButton
                                                color="error"
                                                :variant="
                                                    permStates.get(perm.id) !== undefined && !permStates.get(perm.id)
                                                        ? 'solid'
                                                        : 'soft'
                                                "
                                                icon="i-mdi-close"
                                                :disabled="!canUpdate"
                                                @click="updatePermissionState(perm.id, false)"
                                            />
                                        </UButtonGroup>
                                    </UFormField>

                                    <template v-for="(attr, idx) in attrList" :key="attr.attrId">
                                        <RoleViewAttr
                                            v-if="attr.permissionId === perm.id && !isEmptyAttributes(attr.maxValues)"
                                            v-model="attrList[idx]!"
                                            :permission="perm"
                                            :disabled="!canUpdate || permStates.get(perm.id) !== true"
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
