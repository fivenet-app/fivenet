<script lang="ts" setup>
const props = defineProps<{
    job?: string;
    grade?: number;
}>();

const completorStore = useCompletorStore();

const { data: jobsList } = useAsyncData('completor-jobs', () => completorStore.listJobs());

const job = computed(() => jobsList.value?.find((j) => j.name === props.job));

const grade = computed(() => job.value?.grades.find((g) => g.grade === props.grade));
</script>

<template>
    <span>
        {{ job?.label ?? props.job }}

        <span
            v-if="props.grade !== undefined && props.grade >= 0"
            :title="`${job?.label ?? props.job} - ${$t('common.rank')} ${grade?.grade ?? props.grade}`"
        >
            ({{ grade?.label ?? props.grade }})</span
        >
    </span>
</template>
