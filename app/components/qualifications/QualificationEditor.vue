<script lang="ts" setup>
import type { UForm } from '#components';
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { type AccessType, enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import QualificationRequirementEntry from '~/components/qualifications/QualificationRequirementEntry.vue';
import type { Content } from '~/types/history';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { File } from '~~/gen/ts/resources/file/file';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { type QualificationJobAccess, AccessLevel } from '~~/gen/ts/resources/qualifications/access';
import type { ExamQuestions } from '~~/gen/ts/resources/qualifications/exam';
import {
    type Qualification,
    type QualificationExamSettings,
    type QualificationRequirement,
    type QualificationShort,
    AutoGradeMode,
    QualificationExamMode,
} from '~~/gen/ts/resources/qualifications/qualifications';
import type { UpdateQualificationRequest, UpdateQualificationResponse } from '~~/gen/ts/services/qualifications/qualifications';
import BackButton from '../partials/BackButton.vue';
import ConfirmModal from '../partials/ConfirmModal.vue';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';
import FormatBuilder from '../partials/FormatBuilder.vue';
import ExamEditor from './exam/ExamEditor.vue';

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

const examModes = ref<{ mode: QualificationExamMode; selected?: boolean }[]>([
    { mode: QualificationExamMode.DISABLED },
    { mode: QualificationExamMode.REQUEST_NEEDED },
    { mode: QualificationExamMode.ENABLED },
]);

const schema = z.object({
    weight: z.coerce.number(),
    abbreviation: z.string().min(3).max(20),
    title: z.string().min(3).max(255),
    description: z.union([z.string().min(3).max(512), z.string().length(0).optional()]),
    content: z.string().min(3).max(750000),
    closed: z.coerce.boolean(),
    draft: z.coerce.boolean(),
    public: z.coerce.boolean(),
    discordSyncEnabled: z.coerce.boolean(),
    discordSettings: z.object({
        roleName: z.string().max(64).optional(),
        roleFormat: z.string().max(64).optional(),
    }),
    examMode: z.nativeEnum(QualificationExamMode),
    examSettings: z.custom<QualificationExamSettings>(),
    exam: z.custom<ExamQuestions>(),
    access: z.object({
        jobs: z.custom<QualificationJobAccess>().array().max(maxAccessEntries).default([]),
    }),
    labelSyncEnabled: z.coerce.boolean(),
    labelSyncFormat: z.string().max(128).optional(),
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
        time: {
            seconds: 360,
            nanos: 0,
        },
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
        state.examSettings = qualification.value.examSettings;
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
            examSettings: values.examSettings,
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

const accessTypes: AccessType[] = [{ type: 'job', name: t('common.job', 2) }];

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

const formRef = useTemplateRef<typeof UForm>('formRef');
</script>

<template>
    <UForm
        ref="formRef"
        class="flex min-h-dvh w-full max-w-full flex-1 flex-col overflow-y-auto"
        :schema="schema"
        :state="state"
        @submit="onSubmitThrottle"
    >
        <UDashboardNavbar :title="$t('pages.qualifications.edit.title')">
            <template #right>
                <BackButton :disabled="!canSubmit" />

                <UButton
                    type="submit"
                    trailing-icon="i-mdi-content-save"
                    :disabled="!canDo.edit || !canSubmit"
                    :loading="!canSubmit"
                >
                    <span class="hidden truncate sm:block">
                        {{ $t('common.save') }}
                    </span>
                </UButton>

                <UButton
                    v-if="qualification?.draft"
                    type="submit"
                    color="info"
                    trailing-icon="i-mdi-publish"
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    @click.prevent="
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

        <UDashboardPanelContent class="p-0 sm:pb-0">
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
                class="flex flex-1 flex-col"
                :items="items"
                :ui="{
                    container: 'flex flex-1 flex-col overflow-y-hidden',
                    base: 'flex flex-1 flex-col overflow-y-hidden',
                }"
            >
                <template #content>
                    <div v-if="isRequestPending(status)" class="flex flex-col gap-2">
                        <USkeleton v-for="idx in 6" :key="idx" class="size-24 w-full" />
                    </div>

                    <template v-else>
                        <UDashboardToolbar>
                            <template #default>
                                <div class="flex w-full flex-col gap-2">
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
                                            />
                                        </UFormField>

                                        <div class="flex flex-initial flex-col">
                                            <UFormField class="flex-initial" name="closed" :label="`${$t('common.close', 2)}?`">
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

                        <UFormField
                            v-if="canDo.edit"
                            class="flex flex-1 overflow-y-hidden"
                            name="content"
                            :ui="{ container: 'flex flex-1 flex-col mt-0 overflow-y-hidden', label: { wrapper: 'hidden' } }"
                            label="&nbsp;"
                        >
                            <ClientOnly>
                                <TiptapEditor
                                    v-model="state.content"
                                    v-model:files="state.files"
                                    class="mx-auto w-full max-w-(--breakpoint-xl) flex-1 overflow-y-hidden"
                                    :disabled="!canDo.edit"
                                    :saving="saving"
                                    history-type="qualification"
                                    :target-id="props.qualificationId ?? 0"
                                    filestore-namespace="qualifications"
                                    :filestore-service="(opts) => qualificationsQualificationsClient.uploadFile(opts)"
                                />
                            </ClientOnly>
                        </UFormField>
                    </template>
                </template>

                <template #access>
                    <div class="flex flex-col gap-2 overflow-y-auto px-2">
                        <div>
                            <h2 class="text-highlighted">
                                {{ $t('common.access') }}
                            </h2>

                            <AccessManager
                                v-model:jobs="state.access.jobs"
                                :target-id="qualificationId ?? 0"
                                :disabled="!canDo.access"
                                :access-types="accessTypes"
                                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.qualifications.AccessLevel')"
                            />
                        </div>
                    </div>
                </template>

                <template #details>
                    <div class="flex flex-col gap-2 overflow-y-auto px-2">
                        <div>
                            <h2 class="text-highlighted">
                                {{ $t('common.requirements', 2) }}
                            </h2>

                            <QualificationRequirementEntry
                                v-for="(requirement, idx) in state.requirements"
                                :key="requirement.id"
                                :requirement="requirement"
                                @update-qualification="updateQualificationRequirement(idx, $event)"
                                @remove="state.requirements.splice(idx, 1)"
                            />

                            <UTooltip :text="$t('components.qualifications.add_requirement')">
                                <UButton
                                    :disabled="!canSubmit"
                                    icon="i-mdi-plus"
                                    @click="state.requirements.push({ id: 0, qualificationId: 0, targetQualificationId: 0 })"
                                />
                            </UTooltip>
                        </div>

                        <div>
                            <UAccordion
                                :items="[
                                    { slot: 'discord' as const, label: $t('common.discord'), icon: 'i-simple-icons-discord' },
                                    { slot: 'label' as const, label: $t('common.label', 1), icon: 'i-mdi-tag' },
                                ]"
                            >
                                <template #discord>
                                    <UContainer>
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
                                    </UContainer>
                                </template>

                                <template #item-label>
                                    <UContainer>
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
                                    </UContainer>
                                </template>
                            </UAccordion>
                        </div>

                        <div>
                            <h2 class="text-highlighted">
                                {{ $t('common.exam', 1) }}
                            </h2>

                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="examMode"
                                :label="$t('components.qualifications.exam_mode')"
                                :ui="{ container: '' }"
                            >
                                <ClientOnly>
                                    <USelectMenu v-model="state.examMode" :items="examModes" value-key="mode">
                                        <template #item-label>
                                            <span class="truncate">
                                                {{
                                                    $t(
                                                        `enums.qualifications.QualificationExamMode.${QualificationExamMode[state.examMode]}`,
                                                    )
                                                }}
                                            </span>
                                        </template>

                                        <template #item="{ option }">
                                            <span class="truncate">
                                                {{
                                                    $t(
                                                        `enums.qualifications.QualificationExamMode.${QualificationExamMode[option.mode]}`,
                                                    )
                                                }}
                                            </span>
                                        </template>

                                        <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>
                        </div>
                    </div>
                </template>

                <template #exam>
                    <div v-if="isRequestPending(status)" class="flex flex-col gap-2">
                        <USkeleton v-for="idx in 6" :key="idx" class="size-24 w-full" />
                    </div>

                    <ExamEditor
                        v-else
                        v-model:settings="state.examSettings"
                        v-model:questions="state.exam"
                        :disabled="!canDo.edit"
                        class="overflow-y-auto"
                        :qualification-id="props.qualificationId"
                        @file-uploaded="(file) => state.files.push(file)"
                    />
                </template>
            </UTabs>
        </UDashboardPanelContent>
    </UForm>
</template>
