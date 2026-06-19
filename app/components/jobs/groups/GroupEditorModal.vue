<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useCompletorStore } from '~/stores/completor';
import { getJobsGroupsClient } from '~~/gen/ts/clients';
import type { UploadFileResponse } from '~~/gen/ts/resources/file/filestore';
import type { File as FileGrpc } from '~~/gen/ts/resources/file/file';
import { type Group, GroupMembershipMode, GroupState, GroupType } from '~~/gen/ts/resources/jobs/groups/group';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import {
    groupMembershipModeItems as groupMembershipModeItemKeys,
    groupStateItems as groupStateItemKeys,
    groupTypeItems as groupTypeItemKeys,
} from './helpers';

const props = defineProps<{
    group?: Group;
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
});

const formSnapshot = computed(() => ({
    name: state.name,
    description: state.description,
    shortName: state.shortName,
    color: state.color,
    type: state.type,
    membershipMode: state.membershipMode,
    state: state.state,
}));

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(formSnapshot);

function setFormFromProps(): void {
    state.name = props.group?.name ?? '';
    state.description = props.group?.description ?? '';
    state.shortName = props.group?.shortName ?? '';
    state.color = props.group?.color ?? '';
    state.type = props.group?.type ?? GroupType.MANUAL;
    state.membershipMode = props.group?.membershipMode ?? GroupMembershipMode.FLEXIBLE;
    state.state = props.group?.state ?? GroupState.ACTIVE;
    logoFile.value = props.group?.logoFile;
    selectedLeaderUsers.value = [];

    syncSnapshot();
}

onBeforeMount(setFormFromProps);
watch(
    () => props.group,
    () => setFormFromProps(),
);

const formRef = useTemplateRef('formRef');

const canSubmit = ref<boolean>(true);

const modalTitle = computed(() =>
    props.group?.id ? t('components.jobs.groups.editor.update_title') : t('components.jobs.groups.editor.create_title'),
);
const groupTypeItems = computed(() => groupTypeItemKeys.map((item) => ({ ...item, label: t(item.labelKey) })));
const groupMembershipModeItems = computed(() =>
    groupMembershipModeItemKeys.map((item) => ({ ...item, label: t(item.labelKey) })),
);
const groupStateItems = computed(() =>
    groupStateItemKeys
        .filter((item) => item.value !== GroupState.ARCHIVED)
        .map((item) => ({ ...item, label: t(item.labelKey) })),
);

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
        const payload = {
            name: values.name.trim(),
            description: values.description.trim() ? values.description.trim() : undefined,
            shortName: values.shortName.trim() ? values.shortName.trim() : undefined,
            color: values.color.trim() ? values.color.trim() : undefined,
            type: values.type,
            membershipMode: values.membershipMode,
        };

        const call = props.group?.id
            ? jobsGroupsClient.updateGroup({
                  id: props.group.id,
                  ...payload,
                  state: values.state,
              })
            : jobsGroupsClient.createGroup({
                  ...payload,
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
                        <UFormField name="type" :label="$t('common.type')" required>
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

                        <UFormField name="membershipMode" :label="$t('components.jobs.groups.membership_mode')" required>
                            <ClientOnly>
                                <USelectMenu
                                    v-model="state.membershipMode"
                                    class="w-full"
                                    :items="groupMembershipModeItems"
                                    value-key="value"
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
