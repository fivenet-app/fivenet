<script lang="ts" setup>
import type { User } from '~~/gen/ts/resources/users/users';
import CitizenActions from './CitizenActions.vue';
import SetLabels from './props/SetLabels.vue';

const user = defineModel<User>({ required: true });

const { attr, can } = useAuth();

const isOpen = ref(false);
</script>

<template>
    <UDashboardPanel v-if="user" v-model:open="isOpen" class="max-w-72 flex-1" side="right">
        <template #body>
            <div class="flex flex-1 flex-col">
                <template v-if="user">
                    <UPageCard
                        :ui="{
                            wrapper: 'divide-y divide-transparent! space-y-0 *:pt-2 first:*:pt-2 first:*:pt-0 mb-6',
                        }"
                        :title="$t('common.action', 2)"
                    >
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
                        :ui="{
                            wrapper: 'divide-y divide-transparent! space-y-0 *:pt-2 first:*:pt-2 first:*:pt-0 mb-6',
                        }"
                        :title="$t('common.label', 2)"
                    >
                        <SetLabels v-model="user.props!.labels" :user-id="user.userId" />
                    </UPageCard>
                </template>
            </div>
        </template>
    </UDashboardPanel>
</template>
