<script lang="ts" setup>
import { z } from 'zod';
import ColleagueName from '~/components/jobs/colleagues/ColleagueName.vue';
import QualificationRequestList from '~/components/qualifications/tutor/QualificationRequestList.vue';
import QualificationResultList from '~/components/qualifications/tutor/QualificationResultList.vue';
import QualificationResultTutorModal from '~/components/qualifications/tutor/QualificationResultTutorModal.vue';
import { useCompletorStore } from '~/stores/completor';
import type { Qualification } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';

defineProps<{
    qualification: Qualification;
}>();

const modal = useOverlay();

const completorStore = useCompletorStore();

const schema = z.object({
    user: z.number().optional(),
});

const query = useSearchForm('qualifications_tutor', schema);

const usersLoading = ref(false);

const requests = ref<InstanceType<typeof QualificationRequestList> | null>(null);
const results = ref<InstanceType<typeof QualificationResultList> | null>(null);
</script>

<template>
    <div class="flex flex-1 flex-col gap-2">
        <UForm :schema="schema" :state="query">
            <UFormField class="flex-1" name="user" :label="$t('common.search')">
                <ClientOnly>
                    <UInputMenu
                        v-model="query.user"
                        nullable
                        :search-attributes="['firstname', 'lastname']"
                        :placeholder="$t('common.citizen', 1)"
                        block
                        trailing
                        :search="
                            async (q: string): Promise<UserShort[]> => {
                                usersLoading = true;
                                const users = await completorStore.completeCitizens({
                                    search: q,
                                    userIds: query.user ? [query.user] : [],
                                });
                                usersLoading = false;
                                return users;
                            }
                        "
                        search-lazy
                        :search-placeholder="$t('common.search_field')"
                        leading-icon="i-mdi-search"
                        value-key="userId"
                    >
                        <template #item-label="{ item }">
                            <span v-if="item" class="truncate">
                                {{ userToLabel(item) }}
                            </span>
                        </template>

                        <template #option="{ option: user }">
                            <span class="truncate">
                                <ColleagueName :colleague="user" />
                            </span>
                        </template>

                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.user', 2)]) }} </template>
                    </UInputMenu>
                </ClientOnly>
            </UFormField>
        </UForm>

        <div>
            <h2 class="text-highlighted text-sm">{{ $t('common.request', 2) }}</h2>

            <QualificationRequestList
                ref="requests"
                :qualification="qualification"
                :exam-mode="qualification.examMode"
                @refresh="async () => results?.refresh()"
            />
        </div>

        <div>
            <div class="flex flex-row justify-between gap-2">
                <h2 class="text-highlighted text-sm">{{ $t('common.result', 2) }}</h2>

                <UButton
                    icon="i-mdi-plus"
                    :label="$t('common.add')"
                    @click="
                        modal.open(QualificationResultTutorModal, {
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
