<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums, type AccessType } from '~/components/partials/access/helpers';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useCompletorStore } from '~/stores/completor';
import { getJobsGroupsClient } from '~~/gen/ts/clients';
import type { Access, JobAccess, QualificationAccess, UserAccess } from '~~/gen/ts/resources/access/access';
import type { UploadFileResponse } from '~~/gen/ts/resources/file/filestore';
import type { File as FileGrpc } from '~~/gen/ts/resources/file/file';
import { AccessLevel as GroupAccessLevel } from '~~/gen/ts/resources/jobs/groups/access/access';
import { type Group, GroupMembershipMode, GroupState, GroupType } from '~~/gen/ts/resources/jobs/groups/group';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import {
    groupTypeAllowsStrictMembershipMode,
    isLegacyGroupPolicyState,
    isValidGroupTypeMembershipMode,
    normalizeGroupMembershipMode,
} from './policy';
import {
    groupMembershipModeItems as groupMembershipModeItemKeys,
    groupStateItems as groupStateItemKeys,
    groupTypeItems as groupTypeItemKeys,
} from './helpers';

const props = defineProps<{
    group?: Group;
    access?: Access;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'created', group: Group): void;
    (e: 'updated', group: Group): void;
}>();

const notifications = useNotificationsStore();
const { t } = useI18n();
const { fileUpload } = useAppConfig();
const completorStore = useCompletorStore();

const jobsGroupsClient = await getJobsGroupsClient();

const schema = z.object({
    name: z.coerce.string().min(3).max(64),
    description: z.coerce.string().max(255).default(''),
    shortName: z.coerce.string().max(12).default(''),
    color: z.coerce.string().max(7).default(''),
    type: z.enum(GroupType),
    membershipMode: z.enum(GroupMembershipMode),
    state: z.enum(GroupState),
    access: z.object({
        jobs: z.custom<JobAccess>().array().max(20).default([]),
        users: z.custom<UserAccess>().array().max(20).default([]),
        qualifications: z.custom<QualificationAccess>().array().max(20).default([]),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    description: '',
    shortName: '',
    color: '',
    type: GroupType.MANUAL,
    membershipMode: GroupMembershipMode.FLEXIBLE,
    state: GroupState.ACTIVE,
    access: {
        jobs: [],
        users: [],
        qualifications: [],
    },
});

const formSnapshot = computed(() => ({
    name: state.name,
    description: state.description,
    shortName: state.shortName,
    color: state.color,
    type: state.type,
    membershipMode: state.membershipMode,
    state: state.state,
    access: state.access,
}));

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(formSnapshot);

const isLegacyPolicyState = computed(() => isLegacyGroupPolicyState(props.group));

function cloneAccess(access?: Access): Schema['access'] {
    return {
        jobs: access?.jobs.map((entry) => ({ ...entry })) ?? [],
        users: access?.users.map((entry) => ({ ...entry })) ?? [],
        qualifications: access?.qualifications.map((entry) => ({ ...entry })) ?? [],
    };
}

function setFormFromProps(): void {
    state.name = props.group?.name ?? '';
    state.description = props.group?.description ?? '';
    state.shortName = props.group?.shortName ?? '';
    state.color = props.group?.color ?? '';
    state.type = props.group?.type ?? GroupType.MANUAL;
    state.membershipMode = props.group?.membershipMode ?? GroupMembershipMode.FLEXIBLE;
    state.state = props.group?.state ?? GroupState.ACTIVE;
    state.access = cloneAccess(props.access);
    logoFile.value = props.group?.logoFile;
    selectedLeaderUsers.value = [];
    normalizeGroupPolicyState();

    syncSnapshot();
}

onBeforeMount(() => {
    setFormFromProps();
});
watch(
    () => props.group,
    () => {
        setFormFromProps();
    },
);
watch(
    () => props.access,
    () => {
        state.access = cloneAccess(props.access);
        syncSnapshot();
    },
);

function normalizeGroupPolicyState(): void {
    state.membershipMode = normalizeGroupMembershipMode(state.type, state.membershipMode);
}

watch(
    () => state.type,
    () => {
        normalizeGroupPolicyState();
    },
);

const formRef = useTemplateRef('formRef');

const canSubmit = ref<boolean>(true);

const modalTitle = computed(() =>
    props.group?.id ? t('components.jobs.groups.editor.update_title') : t('components.jobs.groups.editor.create_title'),
);
const groupTypeItems = computed(() => groupTypeItemKeys.map((item) => ({ ...item, label: t(item.labelKey) })));
const groupMembershipModeItems = computed(() =>
    groupMembershipModeItemKeys
        .filter((item) => groupTypeAllowsStrictMembershipMode(state.type) || item.value === GroupMembershipMode.FLEXIBLE)
        .map((item) => ({ ...item, label: t(item.labelKey) })),
);
const membershipModeDisabled = computed(() => groupMembershipModeItems.value.length <= 1);
const groupStateItems = computed(() =>
    groupStateItemKeys
        .filter((item) => item.value !== GroupState.ARCHIVED)
        .map((item) => ({ ...item, label: t(item.labelKey) })),
);
const membershipModeHintKey = computed(() =>
    groupTypeAllowsStrictMembershipMode(state.type)
        ? 'components.jobs.groups.policy.membership_mode_mixed_hint'
        : 'components.jobs.groups.policy.membership_mode_type_hint',
);
const policyStateWarning = computed(() => {
    if (!isLegacyPolicyState.value) return undefined;

    return {
        title: t('components.jobs.groups.policy.legacy_state_title'),
        description: t('components.jobs.groups.policy.legacy_state_content'),
    };
});
const groupAccessTypes: AccessType[] = [
    { label: t('common.job', 2), value: 'job' },
    { label: t('common.qualification', 2), value: 'qualification' },
];

const logoFile = ref<FileGrpc | undefined>(props.group?.logoFile);
const logoUploadGroupId = ref<number>(props.group?.id ?? 0);
const selectedLeaderUsers = ref<UserShort[]>([]);

const { resizeAndUpload } = useFileUploader(
    (opts) => jobsGroupsClient.uploadGroupLogo(opts),
    'jobgrouplogos',
    logoUploadGroupId,
);
const { uploadImages } = useImageUpload();

async function uploadGroupLogo(file: File, groupId: number): Promise<UploadFileResponse | undefined> {
    let uploaded: UploadFileResponse | undefined;
    logoUploadGroupId.value = groupId;

    await uploadImages({
        files: [file],
        uploadOne: (f) => resizeAndUpload(f),
        invalidTypeNotification: {
            title: {
                key: 'components.partials.tiptap_editor.notifications.invalid_file_type_images.title',
                parameters: {},
            },
            description: {
                key: 'components.partials.tiptap_editor.notifications.invalid_file_type_images.content',
                parameters: {},
            },
        },
        onUploaded: (resp) => {
            uploaded = resp;
            if (!resp.file) return;

            logoFile.value = resp.file;
        },
    });

    return uploaded;
}

async function handleGroupLogoUpload(file: File | null | undefined): Promise<void> {
    if (!file || !props.group?.id) return;

    canSubmit.value = false;
    try {
        await uploadGroupLogo(file, props.group.id);
    } finally {
        useTimeoutFn(() => (canSubmit.value = true), 400);
    }
}

async function clearGroupLogo(): Promise<void> {
    if (!props.group?.id) return;

    canSubmit.value = false;
    try {
        const { response } = await jobsGroupsClient.deleteGroupLogo({ id: props.group.id });
        logoFile.value = response.group?.logoFile;
        if (response.group) emit('updated', response.group);
    } finally {
        useTimeoutFn(() => (canSubmit.value = true), 400);
    }
}

async function createOrUpdateGroup(values: Schema): Promise<void> {
    try {
        values.membershipMode = normalizeGroupMembershipMode(values.type, values.membershipMode);
        normalizeGroupPolicyState();
        values.access.jobs.forEach((job) => job.id < 0 && (job.id = 0));
        values.access.users.forEach((user) => {
            if (user.id < 0) user.id = 0;
            user.user = undefined;
        });
        values.access.qualifications.forEach((qualification) => {
            if (qualification.id < 0) qualification.id = 0;
            qualification.qualificationId = qualification.qualificationId ?? 0;
        });

        const payload = {
            name: values.name.trim(),
            description: values.description.trim() ? values.description.trim() : undefined,
            shortName: values.shortName.trim() ? values.shortName.trim() : undefined,
            color: values.color.trim() ? values.color.trim() : undefined,
            type: values.type,
            membershipMode: values.membershipMode,
        };

        if (!isValidGroupTypeMembershipMode(payload.type, payload.membershipMode)) {
            throw new Error('Invalid group policy state');
        }

        const call = props.group?.id
            ? jobsGroupsClient.updateGroup({
                  id: props.group.id,
                  ...payload,
                  access: values.access,
                  state: values.state,
              })
            : jobsGroupsClient.createGroup({
                  ...payload,
                  access: values.access,
                  job: '',
                  leaderUserIds: selectedLeaderUsers.value.map((leader) => leader.userId),
                  manualMemberUserIds: [],
                  rules: [],
              });
        const { response } = await call;
        const group = response.group!;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        if (props.group?.id) {
            emit('updated', group);
        } else {
            emit('created', group);
        }

        emit('close', false);
        syncSnapshot();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateGroup(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;
    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <UModal :title="modalTitle" :close="false" :dismissible="!hasUnsavedChanges && canSubmit">
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{ modalTitle }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <div class="grid gap-4">
                    <UFormField name="name" :label="$t('common.name')" required>
                        <UInput v-model="state.name" class="w-full" name="name" type="text" :placeholder="$t('common.name')" />
                    </UFormField>

                    <UFormField name="description" :label="$t('common.description')">
                        <UTextarea
                            v-model="state.description"
                            class="w-full"
                            name="description"
                            :rows="3"
                            :placeholder="$t('common.description')"
                        />
                    </UFormField>

                    <UFormField name="shortName" :label="$t('components.jobs.groups.short_name')">
                        <UInput
                            v-model="state.shortName"
                            class="w-full"
                            name="shortName"
                            type="text"
                            :placeholder="$t('components.jobs.groups.short_name')"
                        />
                    </UFormField>

                    <UFormField v-if="props.group?.id" name="logoFile" :label="$t('common.logo')">
                        <div v-if="logoFile?.filePath" class="mb-2 flex w-full items-center justify-center">
                            <GenericImg
                                class="size-full max-h-32 min-h-32 max-w-32"
                                :src="`/api/filestore/${logoFile.filePath}`"
                                :alt="`${state.name || $t('common.group', 1)} ${$t('common.logo')}`"
                            />
                        </div>

                        <div class="flex flex-col gap-2 md:flex-row">
                            <UFileUpload
                                class="w-full flex-1 grow"
                                :disabled="!canSubmit"
                                :accept="fileUpload.types.images.join(',')"
                                :placeholder="$t('common.image')"
                                :label="$t('common.file_upload_label')"
                                :description="$t('common.allowed_file_types')"
                                @update:model-value="($event) => handleGroupLogoUpload($event)"
                            />

                            <UButton
                                v-if="logoFile?.id"
                                class="grow-0"
                                variant="outline"
                                color="red"
                                trailing-icon="i-mdi-clear"
                                :disabled="!canSubmit"
                                :label="$t('common.clear')"
                                @click="clearGroupLogo"
                            />
                        </div>
                    </UFormField>

                    <UFormField name="color" :label="$t('common.color', 1)">
                        <ColorPicker v-model="state.color" class="w-full" block />
                    </UFormField>

                    <div class="grid gap-4 sm:grid-cols-2">
                        <UFormField
                            name="type"
                            :label="$t('common.type')"
                            :help="$t('components.jobs.groups.policy.type_help')"
                            required
                        >
                            <ClientOnly>
                                <USelectMenu
                                    v-model="state.type"
                                    class="w-full"
                                    :items="groupTypeItems"
                                    value-key="value"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                >
                                    <template #empty> {{ $t('common.not_found', [$t('common.type')]) }} </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormField>

                        <UFormField
                            name="membershipMode"
                            :label="$t('components.jobs.groups.membership_mode')"
                            :help="$t(membershipModeHintKey)"
                            required
                        >
                            <ClientOnly>
                                <USelectMenu
                                    v-model="state.membershipMode"
                                    class="w-full"
                                    :items="groupMembershipModeItems"
                                    value-key="value"
                                    :disabled="!canSubmit || membershipModeDisabled"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                >
                                    <template #empty>
                                        {{ $t('common.not_found', [$t('components.jobs.groups.membership_modes.title')]) }}
                                    </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormField>

                        <UFormField v-if="props.group?.id" name="state" :label="$t('common.status')" required>
                            <ClientOnly>
                                <USelectMenu
                                    v-model="state.state"
                                    class="w-full"
                                    :items="groupStateItems"
                                    value-key="value"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                >
                                    <template #empty> {{ $t('common.not_found', [$t('common.status')]) }} </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormField>
                    </div>

                    <UAlert
                        v-if="policyStateWarning"
                        color="warning"
                        icon="i-mdi-alert-circle"
                        :title="policyStateWarning.title"
                        :description="policyStateWarning.description"
                    />

                    <UFormField name="access" :label="$t('common.access')">
                        <AccessManager
                            v-model:jobs="state.access.jobs"
                            v-model:users="state.access.users"
                            v-model:qualifications="state.access.qualifications"
                            :target-id="props.group?.id ?? 0"
                            name="access"
                            :access-types="groupAccessTypes"
                            default-access-type="job"
                            :access-roles="enumToAccessLevelEnums(GroupAccessLevel, 'enums.jobs.groups.AccessLevel')"
                            :disabled="!canSubmit"
                            hide-other-jobs
                        />
                    </UFormField>

                    <UFormField v-if="!props.group?.id" :label="$t('components.jobs.groups.leaders')">
                        <SelectMenu
                            v-model="selectedLeaderUsers"
                            class="w-full"
                            multiple
                            :searchable="
                                async (q: string) =>
                                    await completorStore.completeColleagues(
                                        q,
                                        selectedLeaderUsers.map((leader) => leader.userId),
                                    )
                            "
                            searchable-key="jobs-group-create-leaders"
                            :filter-fields="['firstname', 'lastname']"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :placeholder="$t('components.jobs.groups.leaders')"
                            :disabled="!canSubmit"
                        >
                            <template v-if="selectedLeaderUsers.length > 0" #default>
                                {{ usersToLabel(selectedLeaderUsers) }}
                            </template>
                            <template #item-label="{ item }">
                                {{ `${item?.firstname} ${item?.lastname} (${item?.dateofbirth})` }}
                            </template>
                            <template #empty>
                                {{ $t('common.not_found', [$t('common.colleague', 2)]) }}
                            </template>
                        </SelectMenu>
                    </UFormField>
                </div>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />

                <UButton
                    class="flex-1"
                    block
                    :loading="!canSubmit"
                    :disabled="!canSubmit"
                    :label="group?.id ? $t('common.update') : $t('common.create')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
