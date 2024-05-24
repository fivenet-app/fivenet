<script lang="ts" setup>
defineProps<{
    qualificationId: string;
}>();

// TODO
</script>

<template>
    <UDashboardNavbar :title="$t('pages.qualifications.single.exam.edit.title')">
        <template #right>
            <UButton color="black" icon="i-mdi-arrow-back" :to="`/qualifications/${qualificationId}/exam`">
                {{ $t('common.back') }}
            </UButton>
        </template>
    </UDashboardNavbar>

    <!-- TODO -->
</template>
