<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums, type AccessType } from '~/components/partials/access/helpers';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import QualificationRequirementEntry from '~/components/qualifications/QualificationRequirementEntry.vue';
import { useNotificatorStore } from '~/stores/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { QualificationJobAccess } from '~~/gen/ts/resources/qualifications/access';
import { AccessLevel } from '~~/gen/ts/resources/qualifications/access';
import type { ExamQuestions } from '~~/gen/ts/resources/qualifications/exam';
import type {
    QualificationExamSettings,
    QualificationRequirement,
    QualificationShort,
} from '~~/gen/ts/resources/qualifications/qualifications';
import { QualificationExamMode } from '~~/gen/ts/resources/qualifications/qualifications';
import type {
    CreateQualificationRequest,
    CreateQualificationResponse,
    UpdateQualificationRequest,
    UpdateQualificationResponse,
} from '~~/gen/ts/services/qualifications/qualifications';
import ExamEditor from './exam/ExamEditor.vue';

const props = defineProps<{
    qualificationId?: number;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { attr, can, activeChar } = useAuth();

const notifications = useNotificatorStore();

const { maxAccessEntries } = useAppConfig();

const canDo = computed(() => ({
    edit: !props.qualificationId
        ? can('qualifications.QualificationsService.CreateQualification').value
        : can('qualifications.QualificationsService.UpdateQualification').value,
    access: true,
    public: attr('qualifications.QualificationsService.CreateQualification', 'Fields', 'Public').value,
}));

const loading = ref(props.qualificationId !== undefined);

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
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    weight: 0,
    abbreviation: '',
    title: '',
    description: '',
    content: '',
    closed: false,
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
    },
    exam: {
        questions: [],
    },
    access: {
        jobs: [],
    },
    labelSyncEnabled: false,
    labelSyncFormat: '%abbr%: %name%',
});
const qualiRequirements = ref<QualificationRequirement[]>([]);

async function getQualification(qualificationId: number): Promise<void> {
    try {
        const call = $grpc.qualifications.qualifications.getQualification({
            qualificationId: qualificationId,
            withExam: true,
        });
        const { response } = await call;

        const qualification = response.qualification;
        if (qualification) {
            state.abbreviation = qualification.abbreviation;
            state.title = qualification.title;
            state.description = qualification.description;
            state.content = qualification.content?.rawContent ?? '';
            state.closed = qualification.closed;
            state.public = qualification.public;
            state.abbreviation = qualification.abbreviation;
            state.discordSyncEnabled = qualification.discordSyncEnabled;
            state.discordSettings = qualification.discordSettings ?? {
                roleName: '',
                roleFormat: '',
            };
            state.examMode = qualification.examMode;
            if (qualification.examSettings) {
                state.examSettings = qualification.examSettings;
            }
            if (qualification.exam) {
                qualification.exam.questions.forEach((q) => {
                    if (q.answer === undefined) {
                        q.answer = {
                            answerKey: '',
                        };
                    }
                });
                state.exam = qualification.exam;
            }
            if (qualification.access) {
                state.access = qualification.access;
            }

            qualiRequirements.value = qualification.requirements;
        }

        loading.value = false;
    } catch (e) {
        handleGRPCError(e as RpcError);

        await navigateTo({ name: 'qualifications' });
        throw e;
    }
}

onMounted(async () => {
    if (props.qualificationId) {
        await getQualification(props.qualificationId);
    } else {
        state.access.jobs.push({
            id: 0,
            targetId: 0,
            job: activeChar.value!.job,
            minimumGrade: -1,
            access: AccessLevel.EDIT,
        });
    }
});

async function createQualification(values: Schema): Promise<CreateQualificationResponse> {
    const req: CreateQualificationRequest = {
        qualification: {
            id: 0,
            job: '',
            weight: 0,
            closed: values.closed,
            public: values.public,
            abbreviation: values.abbreviation,
            title: values.title,
            description: values.description,
            content: {
                rawContent: values.content,
            },
            creatorId: activeChar.value!.userId,
            creatorJob: activeChar.value!.job,
            requirements: qualiRequirements.value,
            discordSyncEnabled: values.discordSyncEnabled,
            discordSettings: values.discordSettings,
            examMode: values.examMode,
            examSettings: values.examSettings,
            exam: values.exam,
            access: values.access,
            labelSyncEnabled: values.labelSyncEnabled,
            labelSyncFormat: values.labelSyncFormat,
        },
    };

    try {
        const call = $grpc.qualifications.qualifications.createQualification(req);
        const { response } = await call;

        await navigateTo({
            name: 'qualifications-id',
            params: { id: response.qualificationId },
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function updateQualification(values: Schema): Promise<UpdateQualificationResponse> {
    const req: UpdateQualificationRequest = {
        qualification: {
            id: props.qualificationId!,
            job: '',
            weight: 0,
            closed: values.closed,
            public: values.public,
            abbreviation: values.abbreviation,
            title: values.title,
            description: values.description,
            content: {
                rawContent: values.content,
            },
            creatorId: activeChar.value!.userId,
            creatorJob: activeChar.value!.job,
            requirements: qualiRequirements.value,
            discordSyncEnabled: values.discordSyncEnabled,
            discordSettings: values.discordSettings,
            examMode: values.examMode,
            examSettings: values.examSettings,
            exam: values.exam,
            access: values.access,
            labelSyncEnabled: values.labelSyncEnabled,
            labelSyncFormat: values.labelSyncFormat,
        },
    };

    try {
        const call = $grpc.qualifications.qualifications.updateQualification(req);

        const { response } = await call;

        await navigateTo({
            name: 'qualifications-id',
            params: { id: response.qualificationId },
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
    if (props.qualificationId === undefined) {
        await createQualification(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
    } else {
        await updateQualification(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
    }

    notifications.add({
        title: { key: 'notifications.action_successfull.title', parameters: {} },
        description: { key: 'notifications.action_successfull.content', parameters: {} },
        type: NotificationType.SUCCESS,
    });
}, 1000);

const accessTypes: AccessType[] = [{ type: 'job', name: t('common.job', 2) }];

function updateQualificationRequirement(idx: number, qualification?: QualificationShort): void {
    if (!qualification || !qualiRequirements.value[idx]) {
        return;
    }

    qualiRequirements.value[idx]!.qualificationId = props.qualificationId ?? 0;
    qualiRequirements.value[idx]!.targetQualificationId = qualification.id;
}

const items = [
    {
        slot: 'edit',
        label: t('common.edit'),
        icon: 'i-mdi-pencil',
    },
    {
        slot: 'details',
        label: t('common.detail', 2),
        icon: 'i-mdi-details',
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
</script>

<template>
    <UForm
        class="min-h-dscreen flex w-full max-w-full flex-1 flex-col overflow-y-auto"
        :schema="schema"
        :state="state"
        @submit="onSubmitThrottle"
    >
        <UDashboardNavbar :title="qualificationId ? $t('pages.qualifications.edit.title') : $t('pages.qualifications.create')">
            <template #right>
                <UButton
                    color="black"
                    icon="i-mdi-arrow-back"
                    :to="qualificationId ? { name: 'qualifications-id', params: { id: qualificationId } } : '/qualifications'"
                >
                    {{ $t('common.back') }}
                </UButton>

                <UButton
                    type="submit"
                    trailing-icon="i-mdi-content-save"
                    :disabled="!canDo.edit || !canSubmit"
                    :loading="!canSubmit"
                >
                    <span class="hidden truncate sm:block">
                        <template v-if="!qualificationId">
                            {{ $t('common.create') }}
                        </template>
                        <template v-else>
                            {{ $t('common.save') }}
                        </template>
                    </span>
                </UButton>
            </template>
        </UDashboardNavbar>

        <UDashboardPanelContent class="p-0 sm:pb-0">
            <UTabs
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
                <template #edit>
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
                                    class="mx-auto w-full max-w-screen-xl flex-1 overflow-y-hidden"
                                    :disabled="!canDo.edit"
                                    rounded="rounded-none"
                                />
                            </ClientOnly>
                        </UFormGroup>
                    </template>
                </template>

                <template #details>
                    <div class="flex flex-col gap-2 overflow-y-auto px-2">
                        <div>
                            <h2 class="text- text-gray-900 dark:text-white">
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

                        <div>
                            <h2 class="text- text-gray-900 dark:text-white">
                                {{ $t('common.requirements', 2) }}
                            </h2>

                            <QualificationRequirementEntry
                                v-for="(requirement, idx) in qualiRequirements"
                                :key="requirement.id"
                                :requirement="requirement"
                                @update-qualification="updateQualificationRequirement(idx, $event)"
                                @remove="qualiRequirements.splice(idx, 1)"
                            />

                            <UTooltip :text="$t('components.qualifications.add_requirement')">
                                <UButton
                                    :ui="{ rounded: 'rounded-full' }"
                                    :disabled="!canSubmit"
                                    icon="i-mdi-plus"
                                    @click="qualiRequirements.push({ id: 0, qualificationId: 0, targetQualificationId: 0 })"
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
                                            name="discordSettings.enabled"
                                            :label="$t('common.enabled')"
                                            :ui="{ container: 'inline-flex gap-2' }"
                                        >
                                            <UToggle v-model="state.discordSyncEnabled" :disabled="!canDo.edit">
                                                <span class="sr-only">
                                                    {{ $t('common.enabled') }}
                                                </span>
                                            </UToggle>
                                            <span class="text-sm font-medium">{{ $t('common.enabled') }}</span>
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
                                            <UInput
                                                v-model="state.discordSettings.roleFormat"
                                                name="discordSettings.roleFormat"
                                                type="text"
                                                :placeholder="
                                                    $t(
                                                        'components.settings.job_props.discord_sync_settings.qualifications_role_format.title',
                                                    )
                                                "
                                                :disabled="!canDo.edit"
                                            />
                                        </UFormGroup>
                                    </UContainer>
                                </template>

                                <template #label>
                                    <UContainer>
                                        <UFormGroup
                                            name="labelSyncEnabled"
                                            :label="$t('common.enabled')"
                                            :ui="{ container: 'inline-flex gap-2' }"
                                        >
                                            <UToggle v-model="state.labelSyncEnabled" :disabled="!canDo.edit">
                                                <span class="sr-only">
                                                    {{ $t('common.enabled') }}
                                                </span>
                                            </UToggle>
                                            <span class="text-sm font-medium">{{ $t('common.enabled') }}</span>
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
                                            <UInput
                                                v-model="state.labelSyncFormat"
                                                name="labelSyncFormat"
                                                type="text"
                                                :placeholder="
                                                    $t('components.qualifications.qualification_editor.label_sync_format.label')
                                                "
                                                :disabled="!canDo.edit"
                                            />
                                        </UFormGroup>
                                    </UContainer>
                                </template>
                            </UAccordion>
                        </div>

                        <div>
                            <h2 class="text- text-gray-900 dark:text-white">
                                {{ $t('common.exam', 1) }}
                            </h2>

                            <UFormGroup name="examMode">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.examMode"
                                        class="w-40 max-w-40"
                                        :options="examModes"
                                        value-attribute="mode"
                                    >
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
                    />
                </template>
            </UTabs>
        </UDashboardPanelContent>
    </UForm>
</template>
