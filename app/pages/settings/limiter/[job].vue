<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import AttrView from '~/components/settings/attrs/AttrView.vue';

definePageMeta({
    requiresAuth: true,
    permission: 'internal.Superuser/JobAdmin',
    validate: async (route) => {
        route = route as TypedRouteFromName<'settings-limiter-job'>;
        // Check if the id is made up of digits
        if (typeof route.params.job !== 'string') {
            return false;
        }
        return route.params.job.length > 0 && route.params.job.length <= 20;
    },
});

defineEmits<{
    (e: 'deleted'): void;
}>();

const job = useRoute('settings-limiter-job').params.job;
</script>

<template>
    <div v-if="!job" class="p-4 sm:p-6">
        <DataNoDataBlock icon="i-mdi-select" :message="$t('common.none_selected', [$t('common.job', 2)])" :padded="false" />
    </div>
    <AttrView
        v-else
        :job="job"
        @deleted="
            async () => {
                navigateTo('/settings/limiter');
                $emit('deleted');
            }
        "
    />
</template>
