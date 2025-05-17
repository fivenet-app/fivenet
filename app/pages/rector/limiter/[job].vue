<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import AttrView from '~/components/rector/attrs/AttrView.vue';

definePageMeta({
    requiresAuth: true,
    permission: 'SuperUser',
    validate: async (route) => {
        route = route as TypedRouteFromName<'rector-limiter-job'>;
        // Check if the id is made up of digits
        if (typeof route.params.job !== 'string') {
            return false;
        }
        return route.params.job.length > 0 && route.params.job.length <= 20;
    },
});

const job = useRoute('rector-limiter-job').params.job;
</script>

<template>
    <div>
        <DataNoDataBlock v-if="!job" icon="i-mdi-select" :message="$t('common.none_selected', [$t('common.job', 2)])" />
        <AttrView v-else :job="job" @deleted="navigateTo({ name: 'rector-limiter' })" />
    </div>
</template>
