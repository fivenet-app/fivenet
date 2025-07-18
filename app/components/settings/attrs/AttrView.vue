<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import AttrViewAttr from '~/components/settings/attrs/AttrViewAttr.vue';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { RoleAttribute } from '~~/gen/ts/resources/permissions/attributes';
import type { Permission } from '~~/gen/ts/resources/permissions/permissions';
import type { AttrsUpdate, GetJobLimitsResponse, PermsUpdate } from '~~/gen/ts/services/settings/settings';

const props = defineProps<{
    job: string;
}>();

const emit = defineEmits<{
    (e: 'deleted', job: string): void;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { isSuperuser } = useAuth();

const modal = useModal();

const notifications = useNotificationsStore();

const {
    data: jobLimits,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`settings-limiter-${props.job}`, () => getJobLimits(props.job));

const changed = ref(false);

const permList = ref<Permission[]>([]);
const permCategories = ref<Set<string>>(new Set());
const permStates = ref(new Map<number, boolean | undefined>());

const attrList = ref<RoleAttribute[]>([]);

async function getJobLimits(job: string): Promise<GetJobLimitsResponse> {
    try {
        const call = $grpc.settings.settings.getJobLimits({
            job: job,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function getAllPermissions(job: string): Promise<void> {
    try {
        const call = $grpc.settings.settings.getAllPermissions({
            job: job,
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

async function genPermissionCategories(): Promise<void> {
    permCategories.value.clear();

    permList.value.forEach((perm) => {
        permCategories.value.add(perm.category);
    });
}

async function propogatePermissionStates(): Promise<void> {
    permStates.value.clear();

    jobLimits.value?.permissions.forEach((perm) => {
        permStates.value.set(perm.id, Boolean(perm.val));
    });

    jobLimits.value?.attributes.forEach((attr) => {
        const idx = attrList.value.findIndex((a) => a.attrId === attr.attrId);
        if (idx > -1 && attrList.value[idx]) {
            attrList.value[idx].maxValues = attr.maxValues;
        } else {
            attrList.value.push(attr);
        }
    });
}

async function updatePermissionState(perm: number, state: boolean | undefined): Promise<void> {
    changed.value = true;
    permStates.value.set(perm, state);
}

async function updateJobLimits(): Promise<void> {
    const currentPermissions = jobLimits.value?.permissions.map((p) => p.id) ?? [];

    const perms: PermsUpdate = {
        toRemove: [],
        toUpdate: [],
    };
    permStates.value.forEach((state, perm) => {
        if (state !== undefined) {
            const p = jobLimits.value?.permissions.find((v) => v.id === perm);
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

        if (perm === undefined || attr.maxValues === undefined || attr.maxValues.validValues.oneofKind === undefined) {
            attrs.toRemove.push({
                roleId: 0,
                attrId: attr.attrId,
                category: '',
                key: '',
                name: '',
                permissionId: attr.permissionId,
                type: '',
            });
        } else if (attr.maxValues !== undefined) {
            attrs.toUpdate.push({
                roleId: 0,
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
        await $grpc.settings.settings.updateJobLimits({
            job: props.job,
            perms: perms,
            attrs: attrs,
        });

        notifications.add({
            title: { key: 'notifications.settings.role_updated.title', parameters: {} },
            description: { key: 'notifications.settings.role_updated.content', parameters: {} },
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

    await getAllPermissions(props.job);
    await propogatePermissionStates();
}

watch(jobLimits, async () => initializeRoleView());

watch(props, async (value) => {
    if (!jobLimits.value || value.job !== props.job) {
        refresh();
    }
});

type CopyRole = {
    job: string;
    attrList: RoleAttribute[];
    permStates: { id: number; val?: boolean | undefined }[];
};

async function copyRole(): Promise<void> {
    copyToClipboardWrapper(
        JSON.stringify({
            job: props.job,
            attrList: attrList.value,
            permStates: [...permStates.value.entries()].map((v) => ({ id: v[0], val: v[1] })),
        } as CopyRole),
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
            if (a.maxValues?.validValues.oneofKind === undefined) {
                return;
            }

            const at = attrList.value.find((at) => at.attrId === a.attrId);
            if (at) {
                if (
                    at.maxValues?.validValues.oneofKind === 'stringList' &&
                    a.maxValues?.validValues.oneofKind === 'stringList'
                ) {
                    a.maxValues.validValues.stringList.strings = a.maxValues.validValues.stringList.strings.filter(
                        (s) =>
                            at.maxValues?.validValues.oneofKind === 'stringList' &&
                            at.maxValues.validValues.stringList.strings.includes(s),
                    );
                } else if (
                    at.maxValues?.validValues.oneofKind === 'jobList' &&
                    a.maxValues?.validValues.oneofKind === 'jobList'
                ) {
                    a.maxValues.validValues.jobList.strings = a.maxValues.validValues.jobList.strings.filter(
                        (s) =>
                            at.maxValues?.validValues.oneofKind === 'jobList' &&
                            at.maxValues.validValues.jobList.strings.includes(s),
                    );
                } else if (
                    at.maxValues?.validValues.oneofKind === 'jobGradeList' &&
                    a.maxValues?.validValues.oneofKind === 'jobGradeList'
                ) {
                    if (!a.maxValues.validValues.jobGradeList.fineGrained) {
                        const jobs = Object.keys(a.maxValues.validValues.jobGradeList.jobs);
                        const maxJobs = Object.keys(at.maxValues.validValues.jobGradeList.jobs);

                        jobs.filter((j) => !maxJobs.includes(j)).forEach(
                            (j) =>
                                a.maxValues?.validValues.oneofKind === 'jobGradeList' &&
                                // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
                                delete a.maxValues.validValues.jobGradeList.jobs[j],
                        );
                    } else {
                        const jobs = Object.keys(a.maxValues.validValues.jobGradeList.jobs);
                        const maxJobs = Object.keys(at.maxValues.validValues.jobGradeList.jobs);

                        jobs.filter((j) => !maxJobs.includes(j)).forEach(
                            (j) =>
                                a.maxValues?.validValues.oneofKind === 'jobGradeList' &&
                                // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
                                delete a.maxValues.validValues.jobGradeList.jobs[j],
                        );
                    }
                }

                at.maxValues = a.maxValues;
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
    [...permCategories.value.entries()].map((category) => {
        return {
            label: t(`perms.${category[1]}.category`),
            category: category[0],
        };
    }),
);

async function deleteFaction(job: string): Promise<void> {
    try {
        await $grpc.settings.settings.deleteFaction({
            job: job,
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('deleted', job);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await updateJobLimits().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div class="w-full">
        <div class="px-1 sm:px-2">
            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.job', 1)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.job', 1)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!jobLimits" :type="$t('common.job', 1)" :retry="refresh" />

            <template v-else>
                <div class="flex justify-between">
                    <h2 class="line-clamp-2 flex-1 text-3xl" :title="`${$t('common.job')}: ${jobLimits?.job}`">
                        {{ jobLimits?.jobLabel! }}
                    </h2>

                    <UTooltip :text="$t('common.refresh')">
                        <UButton variant="link" icon="i-mdi-refresh" color="primary" @click="refresh()" />
                    </UTooltip>

                    <UTooltip v-if="isSuperuser" :text="$t('common.delete')">
                        <UButton
                            variant="link"
                            icon="i-mdi-delete"
                            color="error"
                            @click="
                                modal.open(ConfirmModal, {
                                    confirm: async () => deleteFaction(jobLimits!.job),
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
