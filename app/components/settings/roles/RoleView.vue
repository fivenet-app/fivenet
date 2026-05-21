<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import {
    buildPermissionGroups,
    getPermissionNamespaceLabel,
    getPermissionServiceLabel,
    type PermissionNamespaceGroup,
    type PermissionServiceGroup,
} from '~/components/settings/permissions';
import RoleViewAttr from '~/components/settings/roles/RoleViewAttr.vue';
import { getSettingsSettingsClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { AttributeValues, type RoleAttribute } from '~~/gen/ts/resources/permissions/attributes/attributes';
import type { Permission, Role } from '~~/gen/ts/resources/permissions/permissions/permissions';
import type { AttrsUpdate, PermsUpdate } from '~~/gen/ts/resources/settings/perms';
import EffectivePermsSlideover from './EffectivePermsSlideover.vue';
import { isEmptyAttributes } from './helpers';

const props = withDefaults(
    defineProps<{
        roleId: number;
        roleCount?: number;
    }>(),
    {
        roleCount: 0,
    },
);

const emit = defineEmits<{
    (e: 'deleted'): void;
}>();

const { t, te } = useI18n();

const { activeChar, can, isSuperuser } = useAuth();

const overlay = useOverlay();

const notifications = useNotificationsStore();

const authStore = useAuthStore();
const { impersonateJob } = authStore;

const settingsSettingsClient = await getSettingsSettingsClient();

const { data: role, status, refresh, error } = useLazyAsyncData(`settings-roles-${props.roleId}`, () => getRole(props.roleId));

const changed = ref(false);

const permList = ref<Permission[]>([]);
const permNamespaces = ref<PermissionNamespaceGroup[]>([]);
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
    permNamespaces.value = buildPermissionGroups(permList.value);
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
                namespace: '',
                service: '',
                name: '',
                key: '',
                permissionId: attr.permissionId,
                type: '',
            });
        } else if (attr.value !== undefined) {
            attrs.toUpdate.push({
                roleId: role.value!.id,
                attrId: attr.attrId,
                value: attr.value,
                namespace: '',
                service: '',
                name: '',
                key: '',
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
    permNamespaces.value.length = 0;
    permStates.value.clear();
}

watch(props, async () => {
    if (!role.value || role.value?.id !== props.roleId) {
        refresh();
    }
});

function resetRole(): void {
    permStates.value.forEach((_, key) => permStates.value.set(key, undefined));
    attrList.value.forEach((attr) => {
        if (!attr.value) return;

        switch (attr.value.validValues.oneofKind) {
            case 'stringList':
                attr.value.validValues.stringList.strings = [];
                break;
            case 'jobList':
                attr.value.validValues.jobList.strings = [];
                break;
            case 'jobGradeList':
                attr.value.validValues.jobGradeList.fineGrained = false;
                attr.value.validValues.jobGradeList.jobs = {};
                attr.value.validValues.jobGradeList.grades = {};
                break;
        }
    });
    changed.value = true;
}

function toggleAll(): void {
    permList.value.forEach((perm) => permStates.value.set(perm.id, true));

    attrList.value.forEach((_, idx) => {
        if (!attrList.value[idx]?.validValues) return;

        attrList.value[idx].value = AttributeValues.clone(attrList.value[idx].validValues);
    });

    changed.value = true;
}

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

async function impersonateRole(grade: number): Promise<void> {
    try {
        await impersonateJob(grade);

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as Error);

        notifications.add({
            title: { key: 'notifications.action_failed.title', parameters: {} },
            description: { key: 'notifications.action_failed.content', parameters: {} },
            type: NotificationType.ERROR,
        });
    }
}

const accordionCategories = computed(() =>
    permNamespaces.value.map((namespace) => {
        const services = namespace.services.map((service) => ({
            ...service,
            label: getPermissionServiceLabel(service.namespace, service.service, t, te),
        }));
        const singleService = services.length === 1 ? services[0] : undefined;
        if (singleService) {
            const serviceNamespaceKey = `perms.${singleService.namespace}.namespace`;
            if (te(serviceNamespaceKey)) singleService.label = t(serviceNamespaceKey) + ` - ${singleService.label}`;
        }

        return {
            ...namespace,
            services,
            singleService,
            label: singleService?.label ?? getPermissionNamespaceLabel(namespace.namespace, t, te),
        };
    }),
);

function getPermissionsForService(service: PermissionServiceGroup): Permission[] {
    return permList.value.filter((perm) => perm.namespace === service.namespace && perm.service === service.service);
}

const canUpdate = can('settings.SettingsService/UpdateRolePerms');

const effectivePermsSlideover = overlay.create(EffectivePermsSlideover, { props: { roleId: props.roleId } });
const confirmModal = overlay.create(ConfirmModal);
const confirmImpersonateModal = overlay.create(ConfirmModal, {
    props: {
        title: t('components.settings.role_view.impersonate_confirm.title'),
        description: t('components.settings.role_view.impersonate_confirm.description'),
        color: 'warning',
        confirm: async () => role.value && impersonateRole(role.value.grade),
    },
});

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

                    <UFieldGroup>
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

                        <UTooltip v-if="canUpdate" :text="$t('common.impersonate')">
                            <UButton
                                variant="link"
                                icon="i-mdi-drama-masks"
                                color="primary"
                                :disabled="!activeChar || activeChar.jobGrade >= role.grade"
                                @click="
                                    confirmImpersonateModal.open({
                                        confirm: async () => role && impersonateRole(role.grade),
                                    })
                                "
                            />
                        </UTooltip>

                        <RefreshButton :loading="isRequestPending(status)" icon-only @click="() => refresh()" />

                        <UTooltip :text="$t('common.reset')">
                            <UButton
                                variant="link"
                                icon="i-mdi-clear"
                                color="error"
                                @click="
                                    confirmModal.open({
                                        confirm: () => resetRole(),
                                    })
                                "
                            />
                        </UTooltip>

                        <UTooltip :text="$t('common.delete')">
                            <!-- Only allow deletion if there is more than one role -->
                            <UButton
                                v-if="can('settings.SettingsService/DeleteRole').value"
                                variant="link"
                                icon="i-mdi-delete"
                                color="error"
                                :disabled="roleCount <= 1"
                                @click="
                                    confirmModal.open({
                                        confirm: async () => role && deleteRole(role.id),
                                    })
                                "
                            />
                        </UTooltip>
                    </UFieldGroup>
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
                                :label="$t('common.save', 1)"
                                @click="onSubmitThrottle"
                            />

                            <UPopover>
                                <UButton
                                    :disabled="changed"
                                    color="neutral"
                                    icon="i-mdi-form-textarea"
                                    :label="$t('common.paste')"
                                />

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

                        <UButton
                            icon="i-mdi-content-copy"
                            :disabled="changed"
                            color="neutral"
                            :label="$t('common.copy')"
                            @click="copyRole"
                        />
                    </div>

                    <UAccordion
                        :items="accordionCategories"
                        type="multiple"
                        :ui="{
                            trigger:
                                'data-[state=open]:border-l data-[state=open]:border-primary data-[state=open]:text-primary pl-2',
                            content: 'data-[state=open]:border-l data-[state=open]:border-primary',
                        }"
                    >
                        <template #content="{ item: namespace }">
                            <div v-if="namespace.singleService" class="flex flex-col divide-y divide-default">
                                <div
                                    v-for="perm in getPermissionsForService(namespace.singleService)"
                                    :key="perm.id"
                                    class="flex flex-col gap-1 py-1 pl-4"
                                >
                                    <UFormField
                                        class="flex flex-1 flex-row items-center gap-2"
                                        :label="$t(`perms.${perm.namespace}.${perm.service}.${perm.name}.key`)"
                                        :description="$t(`perms.${perm.namespace}.${perm.service}.${perm.name}.description`)"
                                        :ui="{ wrapper: 'flex-1' }"
                                    >
                                        <UFieldGroup class="inline-flex flex-initial">
                                            <UButton
                                                color="success"
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
                                        </UFieldGroup>
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

                            <UAccordion
                                v-else
                                class="p-1"
                                :items="namespace.services"
                                type="multiple"
                                :ui="{
                                    trigger: 'data-[state=open]:text-primary pl-4',
                                    content: 'pl-4 pb-2',
                                }"
                            >
                                <template #content="{ item: service }">
                                    <div class="flex flex-col divide-y divide-default">
                                        <div
                                            v-for="perm in getPermissionsForService(service)"
                                            :key="perm.id"
                                            class="flex flex-col gap-1 py-1"
                                        >
                                            <UFormField
                                                class="flex flex-1 flex-row items-center gap-2"
                                                :label="$t(`perms.${perm.namespace}.${perm.service}.${perm.name}.key`)"
                                                :description="
                                                    $t(`perms.${perm.namespace}.${perm.service}.${perm.name}.description`)
                                                "
                                                :ui="{ wrapper: 'flex-1' }"
                                            >
                                                <UFieldGroup class="inline-flex flex-initial">
                                                    <UButton
                                                        color="success"
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
                                                </UFieldGroup>
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
                        </template>
                    </UAccordion>

                    <template v-if="isSuperuser && canUpdate">
                        <USeparator class="my-2" />

                        <UButton
                            class="self-end"
                            size="xs"
                            color="neutral"
                            icon="i-mdi-check-all"
                            :label="$t('common.check_all')"
                            @click="toggleAll()"
                        />
                    </template>
                </div>
            </template>
        </div>
    </div>
</template>
