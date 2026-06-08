<script lang="ts" setup>
import ListPage from '~/components/qualifications/ListPage.vue';
import RequestList from '~/components/qualifications/request/RequestList.vue';
import ResultList from '~/components/qualifications/result/ResultList.vue';

useHead({
    title: 'components.qualifications.user_qualifications',
});

definePageMeta({
    title: 'components.qualifications.user_qualifications',
    requiresAuth: true,
    permission: 'qualifications.QualificationsService/ListQualifications',
});
</script>

<template>
    <ListPage>
        <div class="flex flex-col gap-4">
            <ResultList />

            <RequestList />
        </div>
    </ListPage>
</template>
