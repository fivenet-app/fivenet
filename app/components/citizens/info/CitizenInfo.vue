<script lang="ts" setup>
import type { User } from '~~/gen/ts/resources/users/users';
import CitizenActions from './CitizenActions.vue';
import SetLabels from './props/SetLabels.vue';

const user = defineModel<User>({ required: true });

const { attr, can } = useAuth();

const isOpen = ref(false);
</script>

<template>
    <UDashboardPanel v-model:open="isOpen" class="max-w-90 flex-1" resizeable :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #body>
            <div class="flex flex-1 flex-col gap-4">
                <UPageCard :title="$t('common.action', 2)">
                    <!-- Register kbds for the citizens actions here as it will always be available not like the profile tab content -->
                    <CitizenActions
                        :user="user"
                        register-kbds
                        @update:wanted-status="user.props!.wanted = $event"
                        @update:job="
                            user.job = $event.job.name;
                            user.jobLabel = $event.job.label;
                            user.jobGrade = $event.grade.grade;
                            user.jobGradeLabel = $event.grade.label;
                        "
                        @update:traffic-infraction-points="user.props!.trafficInfractionPoints = $event"
                        @update:mug-shot="user.props!.mugshot = $event"
                    />
                </UPageCard>

                <UPageCard
                    v-if="
                        can('citizens.CitizensService/GetUser').value &&
                        attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.Labels').value
                    "
                    :title="$t('common.label', 2)"
                >
                    <SetLabels v-model="user.props!.labels" :user-id="user.userId" />
                </UPageCard>
            </div>
        </template>
    </UDashboardPanel>
</template>
