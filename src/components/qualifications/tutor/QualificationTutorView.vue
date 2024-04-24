<script lang="ts" setup>
import { Qualification } from '~~/gen/ts/resources/qualifications/qualifications';
import QualificationsRequestsList from '~/components/qualifications/tutor/QualificationsRequestsList.vue';
import QualificationsResultsList from '~/components/qualifications/tutor/QualificationsResultsList.vue';

defineProps<{
    qualification: Qualification;
}>();

const requests = ref<InstanceType<typeof QualificationsRequestsList> | null>(null);
const results = ref<InstanceType<typeof QualificationsResultsList> | null>(null);
</script>

<template>
    <div>
        <div>
            <h2 class="text-lg text-gray-900 dark:text-white">{{ $t('common.request', 2) }}</h2>

            <QualificationsRequestsList
                ref="requests"
                :qualification-id="qualification.id"
                @refresh="async () => results?.refresh()"
            />
        </div>

        <div>
            <h2 class="text-lg text-gray-900 dark:text-white">{{ $t('common.result', 2) }}</h2>

            <QualificationsResultsList
                ref="results"
                :qualification-id="qualification.id"
                @refresh="async () => requests?.refresh()"
            />
        </div>
    </div>
</template>
