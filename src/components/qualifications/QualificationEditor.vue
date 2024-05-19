<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import { useNotificatorStore } from '~/store/notificator';
import {
    AccessLevel,
    QualificationAccess,
    type Qualification,
    QualificationRequirement,
    QualificationShort,
    QualificationDiscordSettings,
} from '~~/gen/ts/resources/qualifications/qualifications';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import type {
    CreateQualificationResponse,
    UpdateQualificationResponse,
} from '~~/gen/ts/services/qualifications/qualifications';
import QualificationAccessEntry from '~/components/qualifications/QualificationAccessEntry.vue';
import QualificationRequirementEntry from '~/components/qualifications/QualificationRequirementEntry.vue';
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import DocEditor from '~/components/partials/DocEditor.vue';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    qualificationId?: string;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const completorStore = useCompletorStore();

const maxAccessEntries = 10;

const canEdit = ref(true);

const canDo = computed(() => ({
    edit: true,
    access: true,
}));

const schema = z.object({
    weight: z.number(),
    abbreviation: z.string().min(3).max(20),
    title: z.string().min(3).max(255),
    description: z.union([z.string().min(3).max(512), z.string().length(0).optional()]),
    content: z.string().min(20).max(750000),
    closed: z.boolean(),
    discordSettings: z.object({
        syncEnabled: z.boolean(),
        roleName: z.string().max(64).optional(),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    weight: 0,
    abbreviation: '',
    title: '',
    description: '',
    content: '',
    closed: false,
    discordSettings: {
        syncEnabled: false,
        roleName: '',
    },
});

const access = ref<
    Map<
        string,
        {
            id: string;
            type: number;
            values: {
                job?: string;
                accessRole?: AccessLevel;
                minimumGrade?: number;
            };
        }
    >
>(new Map());
const qualiAccess = ref<QualificationAccess>();
const qualiRequirements = ref<QualificationRequirement[]>([]);

onMounted(async () => {
    if (props.qualificationId) {
        try {
            const call = $grpc.getQualificationsClient().getQualification({
                qualificationId: props.qualificationId,
            });
            const { response } = await call;

            const qualification = response.qualification;
            qualiAccess.value = response.qualification?.access;

            if (qualification) {
                state.abbreviation = qualification.abbreviation;
                state.title = qualification.title;
                state.description = qualification.description;
                state.content = qualification.content;
                state.closed = qualification.closed;
                state.abbreviation = qualification.abbreviation;
                state.discordSettings = qualification.discordSettings ?? {
                    syncEnabled: false,
                    roleName: '',
                };

                qualiRequirements.value = qualification.requirements;
            }

            if (response.qualification?.access) {
                let accessId = 0;

                response.qualification?.access.jobs.forEach((job) => {
                    const id = accessId.toString();
                    access.value.set(id, {
                        id,
                        type: 0,
                        values: {
                            job: job.job,
                            accessRole: job.access,
                            minimumGrade: job.minimumGrade,
                        },
                    });
                    accessId++;
                });
            }
        } catch (e) {
            $grpc.handleError(e as RpcError);

            await navigateTo({ name: 'qualifications' });

            return;
        }
    } else {
        const accessId = 0;
        access.value.set(accessId.toString(), {
            id: accessId.toString(),
            type: 0,
            values: {
                job: activeChar.value?.job,
                minimumGrade: 1,
                accessRole: AccessLevel.EDIT,
            },
        });
    }

    canEdit.value = true;
});

async function createQualification(values: Schema): Promise<CreateQualificationResponse> {
    const req = {
        qualification: {
            id: '0',
            job: '',
            weight: 0,
            closed: values.closed,
            abbreviation: values.abbreviation,
            title: values.title,
            description: values.description,
            content: values.content,
            creatorId: 0,
            creatorJob: '',
            requirements: qualiRequirements.value,
            access: {
                jobs: [],
            } as QualificationAccess,
            discordSettings: values.discordSettings,
        },
    };
    access.value.forEach((entry) => {
        if (entry.values.accessRole === undefined) {
            return;
        }

        if (entry.type === 0) {
            if (!entry.values.job) {
                return;
            }

            req.qualification.access.jobs.push({
                id: '0',
                qualificationId: '0',
                job: entry.values.job,
                minimumGrade: entry.values.minimumGrade ? entry.values.minimumGrade : 0,
                access: entry.values.accessRole,
            });
        }
    });

    try {
        const call = $grpc.getQualificationsClient().createQualification(req);
        const { response } = await call;

        await navigateTo({
            name: 'qualifications-id',
            params: { id: response.qualificationId },
        });

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function updateQualification(values: Schema): Promise<UpdateQualificationResponse> {
    const req = {
        qualification: {
            id: props.qualificationId!,
            job: '',
            weight: 0,
            closed: values.closed,
            abbreviation: values.abbreviation,
            title: values.title,
            description: values.description,
            content: values.content,
            creatorId: 0,
            creatorJob: '',
            requirements: qualiRequirements.value,
            access: {
                jobs: [],
            } as QualificationAccess,
            discordSettings: values.discordSettings,
        },
    };
    access.value.forEach((entry) => {
        if (entry.values.accessRole === undefined) {
            return;
        }

        if (entry.type === 0) {
            if (!entry.values.job) {
                return;
            }

            req.qualification.access.jobs.push({
                id: '0',
                qualificationId: props.qualificationId!,
                job: entry.values.job,
                minimumGrade: entry.values.minimumGrade ? entry.values.minimumGrade : 0,
                access: entry.values.accessRole,
            });
        }
    });

    try {
        const call = $grpc.getQualificationsClient().updateQualification(req);

        const { response } = await call;

        await navigateTo({
            name: 'qualifications-id',
            params: { id: response.qualificationId },
        });

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    if (props.qualificationId === undefined) {
        await createQualification(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
    } else {
        await updateQualification(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
    }
}, 1000);

function addAccessEntry(): void {
    if (access.value.size > maxAccessEntries - 1) {
        notifications.add({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            description: {
                key: 'notifications.max_access_entry.content',
                parameters: { max: maxAccessEntries.toString() },
            },
            type: NotificationType.ERROR,
        });
        return;
    }

    const id = access.value.size > 0 ? parseInt([...access.value.keys()]?.pop() ?? '1') + 1 : 0;
    access.value.set(id.toString(), {
        id: id.toString(),
        type: 0,
        values: {},
    });
}

function removeAccessEntry(event: { id: string }): void {
    access.value.delete(event.id);
}

function updateAccessEntryType(event: { id: string; type: number }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.type = event.type;
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryName(event: { id: string; job?: Job; req?: Qualification }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    if (event.job) {
        accessEntry.values.job = event.job.name;
    }

    access.value.set(event.id, accessEntry);
}

function updateAccessEntryRank(event: { id: string; rank: JobGrade }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.minimumGrade = event.rank.grade;
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryAccess(event: { id: string; access: AccessLevel }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.accessRole = event.access;
    access.value.set(event.id, accessEntry);
}

function updateQualificationRequirement(idx: number, qualification?: QualificationShort): void {
    if (!qualification) {
        return;
    }

    qualiRequirements.value[idx].qualificationId = props.qualificationId ?? '0';
    qualiRequirements.value[idx].targetQualificationId = qualification.id;
}

const { data: jobs } = useAsyncData('completor-jobs', () => completorStore.listJobs());
</script>

<template>
    <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
        <UDashboardNavbar :title="$t('pages.qualifications.edit.title')">
            <template #right>
                <UButtonGroup class="inline-flex">
                    <UButton
                        color="black"
                        icon="i-mdi-arrow-back"
                        :to="
                            qualificationId ? { name: 'qualifications-id', params: { id: qualificationId } } : '/qualifications'
                        "
                    >
                        {{ $t('common.back') }}
                    </UButton>

                    <UButton
                        type="submit"
                        trailing-icon="i-mdi-content-save"
                        :disabled="!canEdit || !canSubmit"
                        :loading="!canSubmit"
                    >
                        <template v-if="!qualificationId">
                            {{ $t('common.create') }}
                        </template>
                        <template v-else>
                            {{ $t('common.save') }}
                        </template>
                    </UButton>
                </UButtonGroup>
            </template>
        </UDashboardNavbar>

        <UDashboardToolbar>
            <template #default>
                <div class="flex w-full flex-col gap-2">
                    <div class="flex w-full flex-row gap-2">
                        <UFormGroup name="abbreviation" :label="$t('common.abbreviation')" class="max-w-48 shrink" required>
                            <UInput
                                v-model="state.abbreviation"
                                name="abbreviation"
                                type="text"
                                size="xl"
                                :placeholder="$t('common.abbreviation')"
                                :disabled="!canEdit || !canDo.edit"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup name="title" :label="$t('common.title')" class="flex-1" required>
                            <UInput
                                v-model="state.title"
                                name="title"
                                type="text"
                                size="xl"
                                :placeholder="$t('common.title')"
                                :disabled="!canEdit || !canDo.edit"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>
                    </div>

                    <div class="flex w-full flex-row gap-2">
                        <UFormGroup name="description" :label="$t('common.description')" class="flex-1">
                            <UTextarea
                                v-model="state.description"
                                name="description"
                                block
                                :placeholder="$t('common.description')"
                                :disabled="!canEdit || !canDo.edit"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup name="closed" :label="`${$t('common.close', 2)}?`" class="flex-initial">
                            <USelectMenu
                                v-model="state.closed"
                                :disabled="!canEdit || !canDo.edit"
                                :options="[
                                    { label: $t('common.open', 2), closed: false },
                                    { label: $t('common.close', 2), closed: true },
                                ]"
                                value-attribute="closed"
                                :searchable-placeholder="$t('common.search_field')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            >
                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>
                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.close', 1)]) }}
                                </template>
                            </USelectMenu>
                        </UFormGroup>
                    </div>
                </div>
            </template>
        </UDashboardToolbar>

        <template v-if="canDo.edit">
            <UFormGroup name="content">
                <ClientOnly>
                    <DocEditor v-model="state.content" :disabled="!canEdit || !canDo.edit" />
                </ClientOnly>
            </UFormGroup>
        </template>

        <div class="mt-2 flex flex-col gap-2 px-2">
            <div>
                <h2 class="text- text-gray-900 dark:text-white">
                    {{ $t('common.access') }}
                </h2>

                <QualificationAccessEntry
                    v-for="entry in access.values()"
                    :key="entry.id"
                    :init="entry"
                    :read-only="!canDo.access"
                    :jobs="jobs"
                    @type-change="updateAccessEntryType($event)"
                    @name-change="updateAccessEntryName($event)"
                    @rank-change="updateAccessEntryRank($event)"
                    @access-change="updateAccessEntryAccess($event)"
                    @delete-request="removeAccessEntry($event)"
                />

                <UButton
                    :ui="{ rounded: 'rounded-full' }"
                    :title="$t('components.documents.document_editor.add_permission')"
                    :disabled="!canEdit || !canDo.access"
                    icon="i-mdi-plus"
                    @click="addAccessEntry()"
                />
            </div>

            <div>
                <h2 class="text- text-gray-900 dark:text-white">
                    {{ $t('common.requirements') }}
                </h2>

                <QualificationRequirementEntry
                    v-for="(requirement, idx) in qualiRequirements"
                    :key="requirement.id"
                    :requirement="requirement"
                    @update-qualification="updateQualificationRequirement(idx, $event)"
                    @remove="qualiRequirements.splice(idx, 1)"
                />

                <UButton
                    :ui="{ rounded: 'rounded-full' }"
                    :disabled="!canSubmit"
                    icon="i-mdi-plus"
                    @click="qualiRequirements.push({ id: '0', qualificationId: '0', targetQualificationId: '0' })"
                />
            </div>

            <div>
                <h2 class="text- text-gray-900 dark:text-white">
                    {{ $t('common.discord') }}
                </h2>

                <UAccordion :items="[{ slot: 'discord', label: $t('common.discord'), icon: 'i-simple-icons-discord' }]">
                    <template #discord>
                        <UContainer>
                            <UFormGroup name="discordSettings.enabled" :label="$t('common.enabled')">
                                <UToggle v-model="state.discordSettings.syncEnabled" :disabled="!canEdit || !canDo.edit">
                                    <span class="sr-only">
                                        {{ $t('common.enabled') }}
                                    </span>
                                </UToggle>
                                <span class="ml-3 text-sm font-medium">{{ $t('common.enabled') }}</span>
                            </UFormGroup>

                            <UFormGroup name="discordSettings.roleName" :label="$t('common.role')">
                                <UInput
                                    v-model="state.discordSettings.roleName"
                                    name="discordSettings.roleName"
                                    type="text"
                                    :placeholder="$t('common.role')"
                                    :disabled="!canEdit || !canDo.edit"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </UFormGroup>
                        </UContainer>
                    </template>
                </UAccordion>
            </div>
        </div>
    </UForm>
</template>
