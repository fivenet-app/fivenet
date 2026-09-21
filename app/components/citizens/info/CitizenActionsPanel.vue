<script lang="ts" setup>
import Actions from '~/components/citizens/info/Actions.vue';
import SetLabels from '~/components/citizens/info/props/SetLabels.vue';
import type { File } from '~~/gen/ts/resources/file/file';
import type { Job, JobGrade } from '~~/gen/ts/resources/jobs/jobs';
import type { User } from '~~/gen/ts/resources/users/user';

defineProps<{
    user: User;
    registerKBDs?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:wantedStatus', value: boolean): void;
    (e: 'update:job', value: { job: Job; grade: JobGrade }): void;
    (e: 'update:trafficInfractionPoints', value: number): void;
    (e: 'update:mugshot', value?: File): void;
}>();

const { attr } = useAuth();
</script>

<template>
    <div class="flex flex-col gap-3">
        <Actions
            :user="user"
            :register-kbds="registerKBDs"
            @update:wanted-status="emit('update:wantedStatus', $event)"
            @update:job="emit('update:job', $event)"
            @update:traffic-infraction-points="emit('update:trafficInfractionPoints', $event)"
            @update:mug-shot="emit('update:mugshot', $event)"
        />

        <template v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.Labels').value">
            <USeparator />

            <div class="flex flex-col gap-2">
                <h2 class="flex min-w-0 items-center font-semibold text-highlighted">
                    {{ $t('common.label', 2) }}
                </h2>

                <SetLabels v-model="user.props!.labels" class="flex-1" :user-id="user.userId" />
            </div>
        </template>
    </div>
</template>
