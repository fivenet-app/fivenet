<script lang="ts" setup>
import { z } from 'zod';
import QualificationResultTutorModal from '~/components/qualifications/tutor/QualificationResultTutorModal.vue';
import QualificationsRequestsList from '~/components/qualifications/tutor/QualificationsRequestsList.vue';
import QualificationsResultsList from '~/components/qualifications/tutor/QualificationsResultsList.vue';
import { useCompletorStore } from '~/store/completor';
import type { Qualification } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';

defineProps<{
    qualification: Qualification;
}>();

const modal = useModal();

const completorStore = useCompletorStore();

const schema = z.object({
    user: z.custom<UserShort>().optional(),
});

type Schema = z.output<typeof schema>;

const query = ref<Schema>({
    user: undefined,
});

const usersLoading = ref(false);

const requests = ref<InstanceType<typeof QualificationsRequestsList> | null>(null);
const results = ref<InstanceType<typeof QualificationsResultsList> | null>(null);
</script>

<template>
    <div class="flex flex-1 flex-col gap-2">
        <UForm :schema="schema" :state="query">
            <UFormGroup class="flex-1" name="user" :label="$t('common.search')">
                <ClientOnly>
                    <UInputMenu
                        v-model="query.user"
                        nullable
                        :search-attributes="['firstname', 'lastname']"
                        :placeholder="$t('common.citizen', 1)"
                        block
                        trailing
                        by="userId"
                        :search="
                            async (query: string): Promise<UserShort[]> => {
                                usersLoading = true;
                                const users = await completorStore.completeCitizens({
                                    search: query,
                                });
                                usersLoading = false;
                                return users;
                            }
                        "
                        search-lazy
                        :search-placeholder="$t('common.search_field')"
                    >
                        <template #label>
                            <span v-if="query.user" class="truncate">
                                {{ usersToLabel([query.user]) }}
                            </span>
                        </template>
                        <template #option="{ option: user }">
                            <span class="truncate">
                                {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                            </span>
                        </template>
                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>
                        <template #empty> {{ $t('common.not_found', [$t('common.creator', 2)]) }} </template>
                    </UInputMenu>
                </ClientOnly>
            </UFormGroup>
        </UForm>

        <div>
            <h2 class="text-sm text-gray-900 dark:text-white">{{ $t('common.request', 2) }}</h2>

            <QualificationsRequestsList
                ref="requests"
                :qualification-id="qualification.id"
                @refresh="async () => results?.refresh()"
            />
        </div>

        <div>
            <div class="flex flex-row justify-between gap-2">
                <h2 class="text-sm text-gray-900 dark:text-white">{{ $t('common.result', 2) }}</h2>

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

            <QualificationsResultsList
                ref="results"
                :qualification-id="qualification.id"
                :exam-mode="qualification.examMode"
                @refresh="async () => requests?.refresh()"
            />
        </div>
    </div>
</template>
