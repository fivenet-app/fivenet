<script lang="ts" setup>
import { z } from 'zod';
import { useCompletorStore } from '~/stores/completor';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import { QualificationExamMode, type QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { AccessLevelEnum, AccessType, MixedAccessEntry } from './helpers';

const props = withDefaults(
    defineProps<{
        disabled?: boolean;
        showRequired?: boolean;
        accessTypes: AccessType[];
        accessRoles?: AccessLevelEnum[];
        jobs?: Job[] | undefined;
    }>(),
    {
        disabled: false,
        showRequired: false,
        accessRoles: undefined,
        jobs: () => [],
    },
);

defineEmits<{
    (e: 'delete'): void;
}>();

const { $grpc } = useNuxtApp();

const entry = defineModel<MixedAccessEntry>({ required: true });

const completorStore = useCompletorStore();

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

const usersLoading = ref(false);
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
            const { response } = await $grpc.qualifications.qualifications.getQualification({
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
                entry.value.minimumGrade = grades[grades.length - 1]?.grade ?? 0;
            }
        }
    }
}

setFromProps();
watch(props, () => setFromProps());
</script>

<template>
    <UForm class="my-2 flex flex-row items-center gap-1" :schema="schema" :state="entry">
        <UCheckbox
            v-if="showRequired"
            v-model="entry.required"
            class="flex-initial"
            :disabled="disabled"
            :title="$t('common.require')"
            name="required"
        />

        <UFormGroup class="w-40 flex-initial">
            <UInput v-if="accessTypes.length === 1" type="text" disabled :model-value="accessTypes[0]?.name" />
            <ClientOnly v-else>
                <USelectMenu
                    v-model="entry.type"
                    :disabled="disabled"
                    :placeholder="$t('common.type')"
                    searchable
                    :search-attributes="['name']"
                    :searchable-placeholder="$t('common.search_field')"
                    value-attribute="type"
                    option-attribute="label"
                    :options="accessTypes"
                >
                    <template #label>
                        <span class="truncate">{{ accessTypes.find((t) => t.type === entry.type)?.name }}</span>
                    </template>

                    <template #option="{ option }">
                        <span class="truncate">{{ option.name }}</span>
                    </template>

                    <template #option-empty="{ query: search }">
                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                    </template>

                    <template #empty>
                        {{ $t('common.not_found', [$t('common.type')]) }}
                    </template>
                </USelectMenu>
            </ClientOnly>
        </UFormGroup>

        <template v-if="entry.type === 'user'">
            <UFormGroup class="flex-1" name="userId">
                <ClientOnly>
                    <USelectMenu
                        v-model="selectedUser"
                        class="flex-1"
                        :searchable="
                            async (q: string) => {
                                usersLoading = true;
                                const users = await completorStore.completeCitizens({
                                    search: q,
                                    userIds: entry.userId ? [entry.userId] : [],
                                });
                                usersLoading = false;
                                return users;
                            }
                        "
                        searchable-lazy
                        :search-attributes="['firstname', 'lastname']"
                        :searchable-placeholder="$t('common.search_field')"
                        :placeholder="$t('common.citizen', 1)"
                    >
                        <template #label>
                            <template v-if="selectedUser">
                                {{ usersToLabel([selectedUser]) }}
                            </template>
                        </template>

                        <template #option="{ option: user }">
                            {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                        </template>

                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormGroup>
        </template>

        <template v-else-if="entry.type === 'qualification'">
            <UFormGroup class="flex-1" name="qualificationId">
                <ClientOnly>
                    <USelectMenu
                        v-model="selectedQualification"
                        class="flex-1"
                        :searchable="
                            async (q: string) => {
                                const { response } = await $grpc.qualifications.qualifications.listQualifications({
                                    pagination: {
                                        offset: 0,
                                    },
                                    search: q,
                                });
                                return response?.qualifications ?? [];
                            }
                        "
                        searchable-lazy
                        :search-attributes="['abbreviation', 'title']"
                        :searchable-placeholder="$t('common.search_field')"
                        :placeholder="$t('common.qualification', 1)"
                    >
                        <template #label>
                            <template v-if="selectedQualification">
                                <span class="truncate">
                                    {{ selectedQualification.abbreviation }}: {{ selectedQualification.title }}
                                </span>
                            </template>
                        </template>

                        <template #option="{ option: qualification }">
                            {{ `${qualification?.abbreviation}: ${qualification?.title}` }}
                        </template>

                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.qualification', 2)]) }} </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormGroup>
        </template>

        <template v-else>
            <UFormGroup class="flex-1" name="job">
                <ClientOnly>
                    <USelectMenu
                        v-model="entry.job"
                        class="flex-1"
                        :disabled="disabled"
                        option-attribute="label"
                        searchable
                        :search-attributes="['name', 'label']"
                        value-attribute="name"
                        :options="jobs ?? []"
                        :placeholder="$t('common.job')"
                        :searchable-placeholder="$t('common.search_field')"
                    >
                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormGroup>

            <UFormGroup class="flex-1" name="minimumGrade">
                <ClientOnly>
                    <USelectMenu
                        class="flex-1"
                        :model-value="
                            jobs.find((j) => j.name === entry.job)?.grades.find((g) => g.grade === entry.minimumGrade)
                        "
                        :disabled="disabled || !entry.job"
                        option-attribute="label"
                        searchable
                        :search-attributes="['name', 'label']"
                        :options="jobs.find((j) => j.name === entry.job)?.grades ?? []"
                        :placeholder="$t('common.rank')"
                        :searchable-placeholder="$t('common.search_field')"
                        @update:model-value="entry.minimumGrade = $event?.grade ?? undefined"
                    >
                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormGroup>
        </template>

        <UFormGroup class="w-60 flex-initial" name="access">
            <ClientOnly>
                <USelectMenu
                    v-model="entry.access"
                    class="flex-1"
                    :disabled="disabled"
                    option-attribute="label"
                    value-attribute="value"
                    :options="accessRoles"
                    searchable
                    :search-attributes="['label']"
                    :placeholder="$t('common.na')"
                    :searchable-placeholder="$t('common.search_field')"
                >
                    <template #label>
                        {{ accessRoles.find((a) => a.value === entry.access)?.label ?? $t('common.na') }}
                    </template>

                    <template #option-empty="{ query: search }">
                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                    </template>

                    <template #empty> {{ $t('common.not_found', [$t('common.access', 2)]) }} </template>
                </USelectMenu>
            </ClientOnly>
        </UFormGroup>

        <UTooltip v-if="!disabled" :text="$t('components.access.remove_entry')">
            <UButton class="flex-initial" :ui="{ rounded: 'rounded-full' }" icon="i-mdi-close" @click="$emit('delete')" />
        </UTooltip>
    </UForm>
</template>
