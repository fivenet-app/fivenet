<script lang="ts" setup>
import type { UForm } from '#components';
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { type AccessType, enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import RequirementEntry from '~/components/qualifications/RequirementEntry.vue';
import type { Content } from '~/types/history';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { File } from '~~/gen/ts/resources/file/file';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { AccessLevel } from '~~/gen/ts/resources/qualifications/access';
import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam';
import {
    type Qualification,
    type QualificationRequirement,
    type QualificationShort,
    AutoGradeMode,
    QualificationExamMode,
} from '~~/gen/ts/resources/qualifications/qualifications';
import type { UpdateQualificationRequest, UpdateQualificationResponse } from '~~/gen/ts/services/qualifications/qualifications';
import { jobAccessEntry } from '~~/shared/types/validation';
import BackButton from '../partials/BackButton.vue';
import ConfirmModal from '../partials/ConfirmModal.vue';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';
import FormatBuilder from '../partials/FormatBuilder.vue';
import ExamEditor, { examSettings } from './exam/ExamEditor.vue';

const props = defineProps<{
    qualificationId: number;
}>();

const { t } = useI18n();

const overlay = useOverlay();

const { attr, can, activeChar } = useAuth();

const notifications = useNotificationsStore();

const historyStore = useHistoryStore();

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

const { maxAccessEntries } = useAppConfig();

const schema = z.object({
    weight: z.coerce.number(),
    abbreviation: z.coerce.string().min(3).max(20),
    title: z.coerce.string().min(3).max(255),
    description: z.union([z.coerce.string().min(3).max(512), z.coerce.string().length(0).optional()]),
    content: z.coerce.string().min(3).max(750000),
    closed: z.coerce.boolean(),
    draft: z.coerce.boolean(),
    public: z.coerce.boolean(),
    discordSyncEnabled: z.coerce.boolean(),
    discordSettings: z.object({
        roleName: z.coerce.string().max(64).optional(),
        roleFormat: z.coerce.string().max(64).optional(),
    }),
    examMode: z.nativeEnum(QualificationExamMode),
    examSettings: examSettings,
    exam: z.object({
        questions: z.custom<ExamQuestion>().array().max(100).default([]),
    }),
    access: z.object({
        jobs: jobAccessEntry.array().max(maxAccessEntries).default([]),
    }),
    labelSyncEnabled: z.coerce.boolean(),
    labelSyncFormat: z.coerce.string().max(128).optional(),
    files: z.custom<File>().array().max(5).default([]),
    requirements: z.custom<QualificationRequirement>().array().max(10).default([]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    weight: 0,
    abbreviation: '',
    title: '',
    description: '',
    content: '',
    closed: false,
    draft: false,
    public: false,
    discordSyncEnabled: false,
    discordSettings: {
        roleName: '',
        roleFormat: '',
    },
    examMode: QualificationExamMode.DISABLED,
    examSettings: {
        time: 360,
        autoGrade: false,
        autoGradeMode: AutoGradeMode.STRICT,
        minimumPoints: 0,
    },
    exam: {
        questions: [],
    },
    access: {
        jobs: [],
    },
    labelSyncEnabled: false,
    labelSyncFormat: '%abbr%: %name%',
    files: [],
    requirements: [],
});

const changed = ref(false);
const saving = ref(false);

// Track last saved string and timestamp
let lastSavedString = '';
let lastSaveTimestamp = 0;

async function saveHistory(values: Schema, name: string | undefined = undefined, type = 'qualification'): Promise<void> {
    if (saving.value) {
        return;
    }

    const now = Date.now();
    // Skip if identical to last saved or if within MIN_GAP
    if (state.content === lastSavedString || now - lastSaveTimestamp < 5000) {
        return;
    }

    saving.value = true;

    historyStore.addVersion<Content>(
        type,
        props.qualificationId,
        {
            content: values.content,
            files: values.files,
        },
        name,
    );

    useTimeoutFn(() => {
        saving.value = false;
    }, 1750);

    lastSavedString = state.content;
    lastSaveTimestamp = now;
}

historyStore.handleRefresh(() => saveHistory(state, 'qualification'));

watchDebounced(
    state,
    () => {
        if (changed.value) {
            const now = Date.now();
            // Skip if identical to last saved or if within MIN_GAP
            if (state.content === lastSavedString || now - lastSaveTimestamp < 5000) {
                return;
            }

            saveHistory(state);

            lastSavedString = state.content;
            lastSaveTimestamp = now;
        } else {
            changed.value = true;
        }
    },
    {
        debounce: 1000,
        maxWait: 2500,
    },
);

const {
    data: qualification,
    status,
    error,
    refresh,
} = useLazyAsyncData(`qualification-${props.qualificationId}-editor`, () => getQualification(props.qualificationId));

async function getQualification(qualificationId: number): Promise<Qualification> {
    try {
        const call = qualificationsQualificationsClient.getQualification({
            qualificationId: qualificationId,
            withExam: true,
        });
        const { response } = await call;

        return response.qualification!;
    } catch (e) {
        handleGRPCError(e as RpcError);

        await navigateTo({ name: 'qualifications' });
        throw e;
    }
}

useHead({
    title: () =>
        qualification.value?.title
            ? `${qualification.value.abbreviation}: ${qualification.value.title} - ${t('pages.qualifications.edit.title')}`
            : t('pages.qualifications.edit.title'),
});

function setFromProps(): void {
    if (!qualification.value) return;

    state.abbreviation = qualification.value.abbreviation;
    state.title = qualification.value.title;
    state.description = qualification.value.description;
    state.content = qualification.value.content?.rawContent ?? '';
    state.closed = qualification.value.closed;
    state.public = qualification.value.public;
    state.abbreviation = qualification.value.abbreviation;
    state.discordSyncEnabled = qualification.value.discordSyncEnabled;
    state.discordSettings = qualification.value.discordSettings ?? {
        roleName: '',
        roleFormat: '',
    };
    state.examMode = qualification.value.examMode;
    if (qualification.value.examSettings) {
        state.examSettings = {
            time: qualification.value.examSettings.time?.seconds ?? 360,
            autoGrade: qualification.value.examSettings.autoGrade,
            autoGradeMode: qualification.value.examSettings.autoGradeMode,
            minimumPoints: qualification.value.examSettings.minimumPoints,
        };
    }
    if (qualification.value.exam) {
        qualification.value.exam.questions.forEach((q) => {
            if (q.answer === undefined) {
                q.answer = {
                    answerKey: '',
                    answer: {
                        oneofKind: undefined,
                    },
                };
            }
        });
        state.exam = qualification.value.exam;
    }
    if (qualification.value.access) {
        state.access = qualification.value.access;
    }
    state.files = qualification.value.files;
    state.requirements = qualification.value.requirements;
}

watch(qualification, () => setFromProps());

async function updateQualification(values: Schema): Promise<UpdateQualificationResponse> {
    values.access.jobs.forEach((job) => job.id < 0 && (job.id = 0));

    const req: UpdateQualificationRequest = {
        qualification: {
            id: props.qualificationId!,
            job: '',
            weight: 0,
            closed: values.closed,
            draft: values.draft,
            public: values.public,
            abbreviation: values.abbreviation,
            title: values.title,
            description: values.description,
            content: {
                rawContent: values.content,
            },
            creatorId: activeChar.value!.userId,
            creatorJob: activeChar.value!.job,
            requirements: state.requirements,
            discordSyncEnabled: values.discordSyncEnabled,
            discordSettings: values.discordSettings,
            examMode: values.examMode,
            examSettings: {
                time: { seconds: values.examSettings.time ?? 360, nanos: 0 },
                autoGrade: values.examSettings.autoGrade,
                autoGradeMode: values.examSettings.autoGradeMode,
                minimumPoints: values.examSettings.minimumPoints,
            },
            exam: {
                questions: values.exam.questions.slice().map((q, idx) => {
                    if (q.answer?.answer.oneofKind === 'singleChoice') {
                        if (q.answer.answer.singleChoice.choice === '__UNDEFINED__') {
                            q.answer.answer.singleChoice.choice = '';
                        }
                    }
                    q.order = idx + 1; // Ensure order is set correctly
                    return q;
                }),
            },
            access: values.access,
            labelSyncEnabled: values.labelSyncEnabled,
            labelSyncFormat: values.labelSyncFormat,
            files: values.files,
        },
    };

    try {
        const call = qualificationsQualificationsClient.updateQualification(req);
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (event.submitter?.getAttribute('role') === 'tab') {
        return;
    }

    canSubmit.value = false;
    await updateQualification(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const accessTypes: AccessType[] = [{ label: t('common.job', 2), value: 'job' }];

function updateQualificationRequirement(idx: number, qualification?: QualificationShort): void {
    if (!qualification || !state.requirements[idx]) {
        return;
    }

    state.requirements[idx]!.qualificationId = props.qualificationId ?? 0;
    state.requirements[idx]!.targetQualificationId = qualification.id;
}

const canDo = computed(() => ({
    edit: can('qualifications.QualificationsService/UpdateQualification').value,
    access: true,
    public: attr('qualifications.QualificationsService/UpdateQualification', 'Fields', 'Public').value,
}));

const items = [
    {
        slot: 'content' as const,
        label: t('common.edit'),
        icon: 'i-mdi-pencil',
        value: 'content',
    },
    {
        slot: 'details' as const,
        label: t('common.detail', 2),
        icon: 'i-mdi-details',
        value: 'details',
    },
    {
        slot: 'access' as const,
        label: t('common.access'),
        icon: 'i-mdi-key',
        value: 'access',
    },
    {
        slot: 'exam' as const,
        label: t('common.exam', 1),
        icon: 'i-mdi-school',
        value: 'exam',
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        return (route.query.tab as string) || 'content';
    },
    set(tab) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.push({ query: { tab: tab }, hash: '#control-active-item' });
    },
});

const confirmModal = overlay.create(ConfirmModal);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0 overflow-y-hidden' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.qualifications.edit.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <BackButton :disabled="!canSubmit" />

                    <UButton
                        trailing-icon="i-mdi-content-save"
                        :disabled="!canDo.edit || !canSubmit"
                        :loading="!canSubmit"
                        @click="() => formRef?.submit()"
                    >
                        <span class="hidden truncate sm:block">
                            {{ $t('common.save') }}
                        </span>
                    </UButton>

                    <UButton
                        v-if="qualification?.draft"
                        color="info"
                        trailing-icon="i-mdi-publish"
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                        @click="
                            confirmModal.open({
                                title: $t('common.publish_confirm.title', { type: $t('common.qualification', 1) }),
                                description: $t('common.publish_confirm.description'),
                                color: 'info',
                                iconClass: 'text-info-500 dark:text-info-400',
                                icon: 'i-mdi-publish',
                                confirm: () => {
                                    state.draft = false;
                                    formRef?.submit();
                                },
                            })
                        "
                    >
                        <span class="hidden truncate sm:block">
                            {{ $t('common.publish') }}
                        </span>
                    </UButton>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <UForm
                ref="formRef"
                :schema="schema"
                :state="state"
                class="flex min-h-0 w-full flex-1 flex-col overflow-y-hidden"
                @submit="onSubmitThrottle"
            >
                <DataPendingBlock
                    v-if="isRequestPending(status)"
                    :message="$t('common.loading', [$t('common.qualification', 1)])"
                />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.qualification', 1)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock
                    v-else-if="!qualification"
                    icon="i-mdi-file-search"
                    :message="$t('common.not_found', [$t('common.qualification', 1)])"
                    :retry="refresh"
                />

                <UTabs
                    v-else
                    v-model="selectedTab"
                    class="flex-1 flex-col overflow-y-hidden"
                    :items="items"
                    variant="link"
                    :unmount-on-hide="false"
                    :ui="{ content: 'flex flex-col flex-1 min-h-0 max-h-full overflow-y-hidden' }"
                >
                    <template #content>
                        <div v-if="isRequestPending(status)" class="flex flex-col gap-2">
                            <USkeleton v-for="idx in 6" :key="idx" class="size-24 w-full" />
                        </div>

                        <template v-else>
                            <UDashboardToolbar>
                                <template #default>
                                    <div class="mx-auto my-2 flex w-full max-w-(--breakpoint-xl) flex-col gap-2">
                                        <div class="flex w-full flex-row gap-2">
                                            <UFormField
                                                class="max-w-48 shrink"
                                                name="abbreviation"
                                                :label="$t('common.abbreviation')"
                                                required
                                            >
                                                <UInput
                                                    v-model="state.abbreviation"
                                                    name="abbreviation"
                                                    type="text"
                                                    size="xl"
                                                    :placeholder="$t('common.abbreviation')"
                                                    :disabled="!canDo.edit"
                                                    class="w-full"
                                                />
                                            </UFormField>

                                            <UFormField class="flex-1" name="title" :label="$t('common.title')" required>
                                                <UInput
                                                    v-model="state.title"
                                                    name="title"
                                                    type="text"
                                                    size="xl"
                                                    :placeholder="$t('common.title')"
                                                    :disabled="!canDo.edit"
                                                    class="w-full"
                                                />
                                            </UFormField>
                                        </div>

                                        <div class="flex w-full flex-row gap-2">
                                            <UFormField class="flex-1" name="description" :label="$t('common.description')">
                                                <UTextarea
                                                    v-model="state.description"
                                                    name="description"
                                                    block
                                                    :placeholder="$t('common.description')"
                                                    :disabled="!canDo.edit"
                                                    class="w-full"
                                                />
                                            </UFormField>

                                            <div class="flex flex-initial flex-col">
                                                <UFormField
                                                    class="flex-initial"
                                                    name="closed"
                                                    :label="`${$t('common.close', 2)}?`"
                                                >
                                                    <USwitch v-model="state.closed" :disabled="!canDo.edit" />
                                                </UFormField>

                                                <UFormField class="flex-initial" name="public" :label="$t('common.public')">
                                                    <USwitch v-model="state.public" :disabled="!canDo.public" />
                                                </UFormField>
                                            </div>
                                        </div>
                                    </div>
                                </template>
                            </UDashboardToolbar>

                            <ClientOnly v-if="canDo.edit">
                                <TiptapEditor
                                    v-model="state.content"
                                    v-model:files="state.files"
                                    class="m-2 mx-auto w-full max-w-(--breakpoint-xl) flex-1"
                                    :disabled="!canDo.edit"
                                    :saving="saving"
                                    history-type="qualification"
                                    :target-id="props.qualificationId ?? 0"
                                    filestore-namespace="qualifications"
                                    :filestore-service="(opts) => qualificationsQualificationsClient.uploadFile(opts)"
                                />
                            </ClientOnly>
                        </template>
                    </template>

                    <template #access>
                        <UDashboardPanel :ui="{ root: 'min-h-0 overflow-y-auto' }">
                            <template #body>
                                <UContainer>
                                    <UPageCard :title="$t('common.access')">
                                        <AccessManager
                                            v-model:jobs="state.access.jobs"
                                            :target-id="qualificationId ?? 0"
                                            :disabled="!canDo.access"
                                            :access-types="accessTypes"
                                            :access-roles="
                                                enumToAccessLevelEnums(AccessLevel, 'enums.qualifications.AccessLevel')
                                            "
                                            name="access"
                                        />
                                    </UPageCard>
                                </UContainer>
                            </template>
                        </UDashboardPanel>
                    </template>

                    <template #details>
                        <UDashboardPanel :ui="{ root: 'min-h-0 overflow-y-auto' }">
                            <template #body>
                                <UContainer class="mb-4 flex flex-col gap-4">
                                    <UPageCard :title="$t('common.requirements', 2)">
                                        <RequirementEntry
                                            v-for="(requirement, idx) in state.requirements"
                                            :key="requirement.id"
                                            :requirement="requirement"
                                            @update-qualification="updateQualificationRequirement(idx, $event)"
                                            @remove="state.requirements.splice(idx, 1)"
                                        />

                                        <div>
                                            <UTooltip :text="$t('components.qualifications.add_requirement')">
                                                <UButton
                                                    :disabled="!canSubmit"
                                                    icon="i-mdi-plus"
                                                    @click="
                                                        state.requirements.push({
                                                            id: 0,
                                                            qualificationId: 0,
                                                            targetQualificationId: 0,
                                                        })
                                                    "
                                                />
                                            </UTooltip>
                                        </div>
                                    </UPageCard>

                                    <UPageCard :title="$t('common.discord')">
                                        <UFormField
                                            class="grid grid-cols-2 items-center gap-2"
                                            name="discordSettings.enabled"
                                            :label="$t('common.enabled')"
                                            :ui="{ container: '' }"
                                        >
                                            <USwitch v-model="state.discordSyncEnabled" :disabled="!canDo.edit" />
                                        </UFormField>

                                        <UFormField name="discordSettings.roleName" :label="$t('common.role')">
                                            <UInput
                                                v-model="state.discordSettings.roleName"
                                                name="discordSettings.roleName"
                                                type="text"
                                                :placeholder="$t('common.role')"
                                                :disabled="!canDo.edit"
                                                class="w-full"
                                            />
                                        </UFormField>

                                        <UFormField
                                            name="discordSettings.roleFormat"
                                            :label="
                                                $t(
                                                    'components.settings.job_props.discord_sync_settings.qualifications_role_format.title',
                                                )
                                            "
                                            :description="
                                                $t(
                                                    'components.settings.job_props.discord_sync_settings.qualifications_role_format.description',
                                                )
                                            "
                                        >
                                            <FormatBuilder
                                                v-model="state.discordSettings.roleFormat!"
                                                :extensions="[
                                                    { label: $t('common.abbreviation'), value: 'abbr' },
                                                    { label: $t('common.name'), value: 'name' },
                                                ]"
                                                :disabled="!canDo.edit"
                                            />
                                        </UFormField>
                                    </UPageCard>

                                    <UPageCard :title="$t('common.label', 1)">
                                        <UFormField
                                            class="grid grid-cols-2 items-center gap-2"
                                            name="labelSyncEnabled"
                                            :label="$t('common.enabled')"
                                            :ui="{ container: '' }"
                                        >
                                            <USwitch v-model="state.labelSyncEnabled" :disabled="!canDo.edit" />
                                        </UFormField>

                                        <UFormField
                                            name="labelSyncFormat"
                                            :label="
                                                $t('components.qualifications.qualification_editor.label_sync_format.label')
                                            "
                                            :description="
                                                $t(
                                                    'components.qualifications.qualification_editor.label_sync_format.description',
                                                )
                                            "
                                        >
                                            <FormatBuilder
                                                v-model="state.labelSyncFormat!"
                                                :extensions="[
                                                    { label: $t('common.abbreviation'), value: 'abbr' },
                                                    { label: $t('common.name'), value: 'name' },
                                                ]"
                                                :disabled="!canDo.edit"
                                            />
                                        </UFormField>
                                    </UPageCard>
                                </UContainer>
                            </template>
                        </UDashboardPanel>
                    </template>

                    <template #exam>
                        <div v-if="isRequestPending(status)" class="flex flex-col gap-2">
                            <USkeleton v-for="idx in 6" :key="idx" class="size-24 w-full" />
                        </div>

                        <div v-else class="overflow-y-auto">
                            <UContainer class="mb-4 flex flex-col gap-4">
                                <ExamEditor
                                    v-model:exam-mode="state.examMode"
                                    v-model:settings="state.examSettings"
                                    v-model:exam="state.exam"
                                    :disabled="!canDo.edit"
                                    class="mx-auto max-w-(--breakpoint-xl)"
                                    :qualification-id="props.qualificationId"
                                    @file-uploaded="(file) => state.files.push(file)"
                                />
                            </UContainer>
                        </div>
                    </template>
                </UTabs>
            </UForm>
        </template>
    </UDashboardPanel>
</template>
