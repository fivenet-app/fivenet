<script lang="ts" setup>
import { z } from 'zod';
import ColleagueName from '~/components/jobs/colleagues/ColleagueName.vue';
import InputMenu from '~/components/partials/InputMenu.vue';
import RequestList from '~/components/qualifications/tutor/RequestList.vue';
import ResultList from '~/components/qualifications/tutor/ResultList.vue';
import ResultTutorModal from '~/components/qualifications/tutor/ResultTutorModal.vue';
import { useCompletorStore } from '~/stores/completor';
import type { Qualification } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    qualification: Qualification;
}>();

const completorStore = useCompletorStore();

const overlay = useOverlay();
const resultTutorModal = overlay.create(ResultTutorModal, {
    props: {
        qualificationId: props.qualification.id,
        onRefresh: undefined,
    },
});

const schema = z.object({
    user: z.coerce.number().optional(),
});

const query = useSearchForm('qualifications_tutor', schema);

const requests = ref<InstanceType<typeof RequestList> | null>(null);
const results = ref<InstanceType<typeof ResultList> | null>(null);
</script>

<template>
    <UDashboardPanel :ui="{ root: 'h-full min-h-0' }">
        <template #header>
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
                            <template #item-label="{ item }">
                                <ColleagueName :colleague="item" />
                            </template>

                            <template #empty> {{ $t('common.not_found', [$t('common.user', 2)]) }} </template>
                        </InputMenu>
                    </UFormField>
                </UForm>
            </UDashboardToolbar>
        </template>

        <template #body>
            <UPageCard :title="$t('common.request', 2)">
                <RequestList
                    ref="requests"
                    :qualification="qualification"
                    :exam-mode="qualification.examMode"
                    @refresh="async () => results?.refresh()"
                />
            </UPageCard>

            <UPageCard :ui="{ body: 'flex flex-col w-full', title: 'flex flex-row flex-1 items-center w-full' }">
                <template #title>
                    <div class="flex-1">{{ $t('common.result', 2) }}</div>

                    <UTooltip :text="$t('common.add')">
                        <UButton
                            icon="i-mdi-plus"
                            :label="$t('common.add')"
                            @click="
                                resultTutorModal.open({
                                    qualificationId: qualification.id,
                                    onRefresh: () => results?.refresh(),
                                })
                            "
                        />
                    </UTooltip>
                </template>

                <ResultList
                    ref="results"
                    :qualification="qualification"
                    :exam-mode="qualification.examMode"
                    @refresh="async () => requests?.refresh()"
                />
            </UPageCard>
        </template>
    </UDashboardPanel>
</template>
