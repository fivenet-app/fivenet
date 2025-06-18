<script lang="ts" setup>
import type { UForm } from '#components';
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums, type AccessType } from '~/components/partials/access/helpers';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import QualificationRequirementEntry from '~/components/qualifications/QualificationRequirementEntry.vue';
import type { Content } from '~/types/history';
import type { File } from '~~/gen/ts/resources/file/file';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { QualificationJobAccess } from '~~/gen/ts/resources/qualifications/access';
import { AccessLevel } from '~~/gen/ts/resources/qualifications/access';
import type { ExamQuestions } from '~~/gen/ts/resources/qualifications/exam';
import type {
    Qualification,
    QualificationExamSettings,
    QualificationRequirement,
    QualificationShort,
} from '~~/gen/ts/resources/qualifications/qualifications';
import { AutoGradeMode, QualificationExamMode } from '~~/gen/ts/resources/qualifications/qualifications';
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

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const modal = useModal();

const { attr, can, activeChar } = useAuth();

const notifications = useNotificationsStore();

const historyStore = useHistoryStore();

const { maxAccessEntries } = useAppConfig();

const examModes = ref<{ mode: QualificationExamMode; selected?: boolean }[]>([
    { mode: QualificationExamMode.DISABLED },
    { mode: QualificationExamMode.REQUEST_NEEDED },
    { mode: QualificationExamMode.ENABLED },
]);

const schema = z.object({
    weight: z.number(),
    abbreviation: z.string().min(3).max(20),
    title: z.string().min(3).max(255),
    description: z.union([z.string().min(3).max(512), z.string().length(0).optional()]),
    content: z.string().min(3).max(750000),
    closed: z.boolean(),
    draft: z.boolean(),
    public: z.boolean(),
    discordSyncEnabled: z.boolean(),
    discordSettings: z.object({
        roleName: z.string().max(64).optional(),
        roleFormat: z.string().max(64).optional(),
    }),
    examMode: z.nativeEnum(QualificationExamMode),
    examSettings: z.custom<QualificationExamSettings>(),
    exam: z.custom<ExamQuestions>(),
    access: z.object({
        jobs: z.custom<QualificationJobAccess>().array().max(maxAccessEntries),
    }),
    labelSyncEnabled: z.boolean(),
    labelSyncFormat: z.string().max(128).optional(),
    files: z.custom<File>().array().max(5),
    requirements: z.custom<QualificationRequirement>().array().max(10),
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
    pending: loading,
    error,
    refresh,
} = useLazyAsyncData(`qualification-${props.qualificationId}-editor`, () => getQualification(props.qualificationId));

async function getQualification(qualificationId: number): Promise<Qualification> {
    try {
        const call = $grpc.qualifications.qualifications.getQualification({
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
                questions: values.exam.questions.slice().map((q) => {
                    if (q.answer?.answer.oneofKind === 'singleChoice') {
                        if (q.answer.answer.singleChoice.choice === '__UNDEFINED__') {
                            q.answer.answer.singleChoice.choice = '';
                        }
                    }
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
        const call = $grpc.qualifications.qualifications.updateQualification(req);
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
        slot: 'content',
        label: t('common.edit'),
        icon: 'i-mdi-pencil',
    },
    {
        slot: 'details',
        label: t('common.detail', 2),
        icon: 'i-mdi-details',
    },
    {
        slot: 'access',
        label: t('common.access'),
        icon: 'i-mdi-key',
    },
    {
        slot: 'exam',
        label: t('common.exam', 1),
        icon: 'i-mdi-school',
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        const index = items.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items[value]?.slot }, hash: '#' });
    },
});

const formRef = useTemplateRef<typeof UForm>('formRef');
</script>

<template>
    <UForm
        ref="formRef"
        class="min-h-dscreen flex w-full max-w-full flex-1 flex-col overflow-y-auto"
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
                        modal.open(ConfirmModal, {
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
            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.qualification', 1)])" />
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
                    wrapper: 'space-y-0 overflow-y-hidden',
                    container: 'flex flex-1 flex-col overflow-y-hidden',
                    base: 'flex flex-1 flex-col overflow-y-hidden',
                    list: { rounded: '' },
                }"
            >
                <template #content>
                    <div v-if="loading" class="flex flex-col gap-2">
                        <USkeleton v-for="idx in 6" :key="idx" class="size-24 w-full" />
                    </div>

                    <template v-else>
                        <UDashboardToolbar>
                            <template #default>
                                <div class="flex w-full flex-col gap-2">
                                    <div class="flex w-full flex-row gap-2">
                                        <UFormGroup
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
                                        </UFormGroup>

                                        <UFormGroup class="flex-1" name="title" :label="$t('common.title')" required>
                                            <UInput
                                                v-model="state.title"
                                                name="title"
                                                type="text"
                                                size="xl"
                                                :placeholder="$t('common.title')"
                                                :disabled="!canDo.edit"
                                            />
                                        </UFormGroup>
                                    </div>

                                    <div class="flex w-full flex-row gap-2">
                                        <UFormGroup class="flex-1" name="description" :label="$t('common.description')">
                                            <UTextarea
                                                v-model="state.description"
                                                name="description"
                                                block
                                                :placeholder="$t('common.description')"
                                                :disabled="!canDo.edit"
                                            />
                                        </UFormGroup>

                                        <div class="flex flex-initial flex-col">
                                            <UFormGroup class="flex-initial" name="closed" :label="`${$t('common.close', 2)}?`">
                                                <UToggle v-model="state.closed" :disabled="!canDo.edit" />
                                            </UFormGroup>

                                            <UFormGroup class="flex-initial" name="public" :label="$t('common.public')">
                                                <UToggle v-model="state.public" :disabled="!canDo.public" />
                                            </UFormGroup>
                                        </div>
                                    </div>
                                </div>
                            </template>
                        </UDashboardToolbar>

                        <UFormGroup
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
                                    class="mx-auto w-full max-w-screen-xl flex-1 overflow-y-hidden"
                                    :disabled="!canDo.edit"
                                    :saving="saving"
                                    history-type="qualification"
                                    :target-id="props.qualificationId ?? 0"
                                    filestore-namespace="qualifications"
                                    :filestore-service="(opts) => $grpc.qualifications.qualifications.uploadFile(opts)"
                                />
                            </ClientOnly>
                        </UFormGroup>
                    </template>
                </template>

                <template #access>
                    <div class="flex flex-col gap-2 overflow-y-auto px-2">
                        <div>
                            <h2 class="text-gray-900 dark:text-white">
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
                            <h2 class="text-gray-900 dark:text-white">
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
                                    :ui="{ rounded: 'rounded-full' }"
                                    :disabled="!canSubmit"
                                    icon="i-mdi-plus"
                                    @click="state.requirements.push({ id: 0, qualificationId: 0, targetQualificationId: 0 })"
                                />
                            </UTooltip>
                        </div>

                        <div>
                            <UAccordion
                                :items="[
                                    { slot: 'discord', label: $t('common.discord'), icon: 'i-simple-icons-discord' },
                                    { slot: 'label', label: $t('common.label', 1), icon: 'i-mdi-tag' },
                                ]"
                            >
                                <template #discord>
                                    <UContainer>
                                        <UFormGroup
                                            class="grid grid-cols-2 items-center gap-2"
                                            name="discordSettings.enabled"
                                            :label="$t('common.enabled')"
                                            :ui="{ container: '' }"
                                        >
                                            <UToggle v-model="state.discordSyncEnabled" :disabled="!canDo.edit" />
                                        </UFormGroup>

                                        <UFormGroup name="discordSettings.roleName" :label="$t('common.role')">
                                            <UInput
                                                v-model="state.discordSettings.roleName"
                                                name="discordSettings.roleName"
                                                type="text"
                                                :placeholder="$t('common.role')"
                                                :disabled="!canDo.edit"
                                            />
                                        </UFormGroup>

                                        <UFormGroup
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
                                        </UFormGroup>
                                    </UContainer>
                                </template>

                                <template #label>
                                    <UContainer>
                                        <UFormGroup
                                            class="grid grid-cols-2 items-center gap-2"
                                            name="labelSyncEnabled"
                                            :label="$t('common.enabled')"
                                            :ui="{ container: '' }"
                                        >
                                            <UToggle v-model="state.labelSyncEnabled" :disabled="!canDo.edit" />
                                        </UFormGroup>

                                        <UFormGroup
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
                                        </UFormGroup>
                                    </UContainer>
                                </template>
                            </UAccordion>
                        </div>

                        <div>
                            <h2 class="text-gray-900 dark:text-white">
                                {{ $t('common.exam', 1) }}
                            </h2>

                            <UFormGroup
                                class="grid grid-cols-2 items-center gap-2"
                                name="examMode"
                                :label="$t('components.qualifications.exam_mode')"
                                :ui="{ container: '' }"
                            >
                                <ClientOnly>
                                    <USelectMenu v-model="state.examMode" :options="examModes" value-attribute="mode">
                                        <template #label>
                                            <span class="truncate">
                                                {{
                                                    $t(
                                                        `enums.qualifications.QualificationExamMode.${QualificationExamMode[state.examMode]}`,
                                                    )
                                                }}
                                            </span>
                                        </template>

                                        <template #option="{ option }">
                                            <span class="truncate">
                                                {{
                                                    $t(
                                                        `enums.qualifications.QualificationExamMode.${QualificationExamMode[option.mode]}`,
                                                    )
                                                }}
                                            </span>
                                        </template>

                                        <template #option-empty="{ query: search }">
                                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                        </template>

                                        <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormGroup>
                        </div>
                    </div>
                </template>

                <template #exam>
                    <div v-if="loading" class="flex flex-col gap-2">
                        <USkeleton v-for="idx in 6" :key="idx" class="size-24 w-full" />
                    </div>

                    <ExamEditor
                        v-else
                        v-model:settings="state.examSettings"
                        v-model:questions="state.exam"
                        class="overflow-y-auto"
                        :qualification-id="props.qualificationId"
                        @file-uploaded="(file) => state.files.push(file)"
                    />
                </template>
            </UTabs>
        </UDashboardPanelContent>
    </UForm>
</template>
