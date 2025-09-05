<script lang="ts" setup>
import { z } from 'zod';
import ColleagueName from '~/components/jobs/colleagues/ColleagueName.vue';
import InputMenu from '~/components/partials/InputMenu.vue';
import QualificationRequestList from '~/components/qualifications/tutor/QualificationRequestList.vue';
import QualificationResultList from '~/components/qualifications/tutor/QualificationResultList.vue';
import QualificationResultTutorModal from '~/components/qualifications/tutor/QualificationResultTutorModal.vue';
import { useCompletorStore } from '~/stores/completor';
import type { Qualification } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    qualification: Qualification;
}>();

const completorStore = useCompletorStore();

const overlay = useOverlay();
const qualificationResultTutorModal = overlay.create(QualificationResultTutorModal, {
    props: {
        qualificationId: props.qualification.id,
        onRefresh: undefined,
    },
});

const schema = z.object({
    user: z.number().optional(),
});

const query = useSearchForm('qualifications_tutor', schema);

const requests = ref<InstanceType<typeof QualificationRequestList> | null>(null);
const results = ref<InstanceType<typeof QualificationResultList> | null>(null);
</script>

<template>
    <div class="flex flex-1 flex-col gap-2">
        <UDashboardToolbar>
            <UForm :schema="schema" :state="query" class="mb-2 flex-1">
                <UFormField class="flex-1" name="user" :label="$t('common.search')">
                    <InputMenu
                        v-model="query.user"
                        :filter-fields="['firstname', 'lastname']"
                        :placeholder="$t('common.citizen', 1)"
                        trailing
                        :searchable="
                            async (q: string): Promise<UserShort[]> =>
                                await completorStore.completeCitizens({
                                    search: q,
                                    userIds: query.user ? [query.user] : [],
                                })
                        "
                        searchable-key="completor-citizens"
                        leading-icon="i-mdi-search"
                        value-key="userId"
                        class="w-full"
                    >
                        <template #item="{ item }">
                            <span class="truncate">
                                <ColleagueName :colleague="item" />
                            </span>
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.user', 2)]) }} </template>
                    </InputMenu>
                </UFormField>
            </UForm>
        </UDashboardToolbar>

        <div>
            <h2 class="text-sm text-highlighted">{{ $t('common.request', 2) }}</h2>

            <QualificationRequestList
                ref="requests"
                :qualification="qualification"
                :exam-mode="qualification.examMode"
                @refresh="async () => results?.refresh()"
            />
        </div>

        <div>
            <div class="flex flex-row justify-between gap-2">
                <h2 class="text-sm text-highlighted">{{ $t('common.result', 2) }}</h2>

                <UButton
                    icon="i-mdi-plus"
                    :label="$t('common.add')"
                    @click="
                        qualificationResultTutorModal.open({
                            qualificationId: qualification.id,
                            onRefresh: () => results?.refresh(),
                        })
                    "
                />
            </div>

            <QualificationResultList
                ref="results"
                :qualification="qualification"
                :exam-mode="qualification.examMode"
                @refresh="async () => requests?.refresh()"
            />
        </div>
    </div>
</template>
