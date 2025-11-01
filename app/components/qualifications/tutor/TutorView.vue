<script lang="ts" setup>
import { z } from 'zod';
import SelectMenu from '~/components/partials/SelectMenu.vue';
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
    users: z.coerce.number().array().max(5).default([]),
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
                    <UFormField class="flex-1" name="users" :label="$t('common.search')">
                        <SelectMenu
                            v-model="query.users"
                            multiple
                            :filter-fields="['firstname', 'lastname']"
                            :placeholder="$t('common.citizen', 2)"
                            :searchable="
                                async (q: string): Promise<UserShort[]> =>
                                    await completorStore.completeCitizens({
                                        search: q,
                                        userIds: query.users ? query.users : [],
                                    })
                            "
                            searchable-key="completor-citizens"
                            leading-icon="i-mdi-search"
                            value-key="userId"
                            class="w-full"
                        >
                            <template v-if="query.users.length > 0" #default>
                                {{ $t('common.selected', query.users.length) }}
                            </template>

                            <template #item-label="{ item: user }">
                                {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                            </template>

                            <template #empty> {{ $t('common.not_found', [$t('common.user', 2)]) }} </template>
                        </SelectMenu>
                    </UFormField>
                </UForm>
            </UDashboardToolbar>
        </template>

        <template #body>
            <UPageCard :title="$t('common.request', 2)" :ui="{ body: '' }">
                <template #default>
                    <RequestList
                        ref="requests"
                        class="-mx-4 -mb-4 sm:-mx-6 sm:-mb-6"
                        :qualification="qualification"
                        :exam-mode="qualification.examMode"
                        :search-query="query"
                        @refresh="async () => results?.refresh()"
                    />
                </template>
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

                <template #default>
                    <ResultList
                        ref="results"
                        class="-mx-4 -mb-4 sm:-mx-6 sm:-mb-6"
                        :qualification="qualification"
                        :exam-mode="qualification.examMode"
                        :search-query="query"
                        @refresh="async () => requests?.refresh()"
                    />
                </template>
            </UPageCard>
        </template>
    </UDashboardPanel>
</template>
