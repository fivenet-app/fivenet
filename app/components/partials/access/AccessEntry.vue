<script lang="ts" setup>
import { z } from 'zod';
import { useCompletorStore } from '~/stores/completor';
import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import { QualificationExamMode, type QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import SelectMenu from '../SelectMenu.vue';
import type { AccessLevelEnum, AccessType, MixedAccessEntry } from './helpers';

const props = withDefaults(
    defineProps<{
        disabled?: boolean;
        showRequired?: boolean;
        accessTypes: AccessType[];
        accessRoles?: AccessLevelEnum[];
        jobs?: Job[] | undefined;
        hideGrade?: boolean;
        hideJobs?: string[];
        name?: string;
    }>(),
    {
        disabled: false,
        showRequired: false,
        accessRoles: undefined,
        jobs: () => [],
        hideGrade: false,
        hideJobs: () => [],
        name: undefined,
    },
);

defineEmits<{
    (e: 'delete'): void;
}>();

const entry = defineModel<MixedAccessEntry>({ required: true });

const completorStore = useCompletorStore();

const { game } = useAppConfig();

const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

const schema = z.object({
    id: z.coerce.number(),
    type: z.coerce.number(),
    userId: z.coerce.number().optional(),
    user: z.custom<UserShort>().optional(),
    job: z.string().optional(),
    minimumGrade: z.coerce.number().optional(),
    qualificationId: z.coerce.number().optional(),
    access: z.coerce.number(),
    required: z.coerce.boolean().optional(),
});

const selectedUser = ref<UserShort | undefined>();
watch(selectedUser, () => {
    entry.value.user = selectedUser.value;
    entry.value.userId = selectedUser.value?.userId;
});

const selectedQualification = ref<QualificationShort | undefined>();
watch(selectedQualification, () => {
    entry.value.qualification = selectedQualification.value;
    entry.value.qualificationId = selectedQualification.value?.id;
});

async function findUser(userId?: number): Promise<UserShort[]> {
    if (userId === undefined) return [];

    return completorStore.completeCitizens({
        search: '',
        userIds: [userId],
    });
}

async function setFromProps(): Promise<void> {
    if (entry.value.type === 'user' && entry.value.userId !== undefined) {
        if (selectedUser.value?.userId === entry.value.userId) {
            return;
        }

        const users = await findUser(entry.value.userId);
        selectedUser.value = users.find((char) => char.userId === entry.value.userId);
    } else if (entry.value.type === 'qualification' && entry.value.qualificationId !== undefined) {
        if (selectedQualification.value?.id === entry.value.qualificationId || entry.value.qualificationId === undefined) {
            return;
        }

        try {
            const { response } = await qualificationsQualificationsClient.getQualification({
                qualificationId: entry.value.qualificationId,
            });
            selectedQualification.value = response.qualification;
        } catch (_) {
            // Fallback to show qualification id
            selectedQualification.value = {
                id: entry.value.qualificationId,
                job: '',
                weight: 0,
                abbreviation: 'N/A',
                title: 'N/A (ID: ' + entry.value.qualificationId + ')',
                closed: false,
                draft: false,
                public: false,
                creatorJob: '',
                examMode: QualificationExamMode.UNSPECIFIED,
                requirements: [],
            };
        }
    } else if (entry.value.type === 'job') {
        if (entry.value.minimumGrade === -1) {
            const grades = props.jobs.find((j) => j.name === entry.value.job)?.grades;
            if (grades) {
                entry.value.minimumGrade = grades[grades.length - 1]?.grade ?? game.startJobGrade;
            }
        }
    }
}

setFromProps();
watch(props, () => setFromProps());
if (props.hideGrade) {
    watch(
        () => entry.value.job,
        () => {
            if (!props.hideGrade) return;

            // If hide grade is true, we must set the minimumGrade to a sane default value
            if (entry.value.job && entry.value.minimumGrade === undefined) {
                const grades = props.jobs.find((j) => j.name === entry.value.job)?.grades;
                if (grades) {
                    entry.value.minimumGrade = grades[grades.length - 1]?.grade ?? game.startJobGrade;
                }
            }
        },
    );
}
</script>

<template>
    <div class="flex flex-1 flex-col gap-1 pb-2 md:flex-row md:pb-0">
        <div class="grid grid-cols-2 gap-2 md:flex md:flex-1">
            <div class="flex flex-initial flex-row items-center gap-2">
                <UFormField v-if="showRequired">
                    <UTooltip class="flex-initial" :text="$t('common.require')">
                        <UCheckbox v-model="entry.required" :disabled="disabled" name="required" />
                    </UTooltip>
                </UFormField>

                <UFormField
                    :name="`${name}.type`"
                    class="min-w-40 flex-initial"
                    :label="$t('common.type')"
                    :ui="{ label: 'md:hidden' }"
                >
                    <UInput v-if="accessTypes.length === 1" type="text" disabled :model-value="accessTypes[0]?.label" />
                    <ClientOnly v-else>
                        <USelectMenu
                            v-model="entry.type"
                            class="w-full"
                            :disabled="disabled"
                            :placeholder="$t('common.type')"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            value-key="value"
                            :items="accessTypes"
                        >
                            <template #default>
                                <span class="truncate">{{ accessTypes.find((t) => t.value === entry.type)?.label }}</span>
                            </template>

                            <template #item="{ item }">
                                <span class="truncate">{{ item.label }}</span>
                            </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('common.type')]) }}
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>
            </div>

            <UFormField
                v-if="entry.type === 'user'"
                class="flex-1"
                :name="`${name}.userId`"
                :label="$t('common.user')"
                :ui="{ label: 'md:hidden' }"
            >
                <SelectMenu
                    v-model="selectedUser"
                    class="w-full"
                    :searchable="
                        async (q: string) =>
                            await completorStore.completeCitizens({
                                search: q,
                                userIds: entry.userId ? [entry.userId] : [],
                            })
                    "
                    searchable-key="completor-citizens"
                    :filter-fields="['firstname', 'lastname']"
                    :search-input="{ placeholder: $t('common.search_field') }"
                    :placeholder="$t('common.citizen', 1)"
                >
                    <template v-if="selectedUser" #default>
                        {{ usersToLabel([selectedUser]) }}
                    </template>

                    <template #item="{ item }">
                        {{ `${item?.firstname} ${item?.lastname} (${item?.dateofbirth})` }}
                    </template>

                    <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                </SelectMenu>
            </UFormField>

            <UFormField
                v-else-if="entry.type === 'qualification'"
                class="flex-1"
                :name="`${name}.qualificationId`"
                :label="$t('common.qualification')"
                :ui="{ label: 'md:hidden' }"
            >
                <SelectMenu
                    v-model="selectedQualification"
                    class="w-full"
                    :searchable="
                        async (q: string) => {
                            const { response } = await qualificationsQualificationsClient.listQualifications({
                                pagination: {
                                    offset: 0,
                                },
                                search: q,
                            });
                            return (response?.qualifications ?? []) as QualificationShort[];
                        }
                    "
                    searchable-key="complete-qualifications"
                    :filter-fields="['abbreviation', 'title']"
                    :search-input="{ placeholder: $t('common.search_field') }"
                    :placeholder="$t('common.qualification', 1)"
                >
                    <template v-if="selectedQualification" #default>
                        <span class="truncate">
                            {{ selectedQualification.abbreviation }}: {{ selectedQualification.title }}
                        </span>
                    </template>

                    <template #item="{ item }">
                        {{ `${item?.abbreviation}: ${item?.title}` }}
                    </template>

                    <template #empty> {{ $t('common.not_found', [$t('common.qualification', 2)]) }} </template>
                </SelectMenu>
            </UFormField>

            <template v-else>
                <UFormField class="flex-1" :name="`${name}.job`" :label="$t('common.job')" :ui="{ label: 'md:hidden' }">
                    <ClientOnly>
                        <USelectMenu
                            v-model="entry.job"
                            class="w-full"
                            :disabled="disabled"
                            :filter-fields="['name', 'label']"
                            value-key="name"
                            :items="jobs?.filter((j) => hideJobs.length === 0 || !hideJobs.includes(j.name)) ?? []"
                            :placeholder="$t('common.job')"
                            :search-input="{ placeholder: $t('common.search_field') }"
                        >
                            <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>

                <UFormField
                    v-if="!hideGrade"
                    class="flex-1"
                    :name="`${name}.minimumGrade`"
                    :label="$t('common.rank')"
                    :ui="{ label: 'md:hidden' }"
                >
                    <ClientOnly>
                        <USelectMenu
                            class="w-full"
                            :model-value="
                                jobs.find((j) => j.name === entry.job)?.grades.find((g) => g.grade === entry.minimumGrade)
                            "
                            :disabled="disabled || !entry.job"
                            :filter-fields="['name', 'label']"
                            :items="jobs.find((j) => j.name === entry.job)?.grades ?? []"
                            :placeholder="$t('common.rank')"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            @update:model-value="entry.minimumGrade = $event?.grade ?? undefined"
                        >
                            <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>
            </template>

            <UFormField
                class="min-w-60 flex-initial"
                :name="`${name}.access`"
                :label="$t('common.access')"
                :ui="{ label: 'md:hidden' }"
            >
                <ClientOnly>
                    <USelectMenu
                        v-model="entry.access"
                        class="w-full"
                        :disabled="disabled"
                        value-key="value"
                        :items="accessRoles"
                        :filter-fields="['label']"
                        :placeholder="$t('common.na')"
                        :search-input="{ placeholder: $t('common.search_field') }"
                    >
                        <template #default>
                            {{ accessRoles.find((a) => a.value === entry.access)?.label ?? $t('common.na') }}
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.access', 2)]) }} </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormField>
        </div>

        <UFormField class="md:mt-1" :ui="{ container: 'flex justify-end-safe md:inline' }">
            <UTooltip v-if="!disabled" :text="$t('components.access.remove_entry')">
                <UButton
                    color="red"
                    class="flex-initial"
                    icon="i-mdi-close"
                    :label="$t('components.access.remove_entry')"
                    :ui="{ label: 'md:hidden' }"
                    @click="$emit('delete')"
                />
            </UTooltip>
        </UFormField>
    </div>
</template>
