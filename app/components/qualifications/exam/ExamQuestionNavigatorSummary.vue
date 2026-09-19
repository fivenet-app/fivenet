<script lang="ts" setup>
defineProps<{
    answered: number;
    unanswered: number;
    flagged: number;
}>();
</script>

<template>
    <div class="flex flex-wrap items-center gap-1">
        <UBadge
            color="success"
            variant="soft"
            icon="i-mdi-check-circle"
            :label="`${$t('components.qualifications.exam_view.answered')}: ${answered}`"
        />
        <UBadge
            color="neutral"
            variant="soft"
            icon="i-mdi-circle-outline"
            :label="`${$t('components.qualifications.exam_view.unanswered')}: ${unanswered}`"
        />
        <UBadge
            color="warning"
            variant="soft"
            icon="i-mdi-flag"
            :label="`${$t('components.qualifications.exam_view.flagged')}: ${flagged}`"
        />
    </div>
</template>
