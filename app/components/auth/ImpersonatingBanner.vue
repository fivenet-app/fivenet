<script lang="ts" setup>
const props = defineProps<{
    job: string;
    jobGrade: number;
}>();

const completorStore = useCompletorStore();
const { listJobs } = completorStore;
const { jobs } = storeToRefs(completorStore);

const { isSuperuser } = useAuth();

const authStore = useAuthStore();
const { impersonateJob, setSuperuserMode } = authStore;

const foundJob = computed(() => jobs.value.find((j) => j.name === props.job));

onBeforeMount(async () => listJobs());
</script>

<template>
    <UBanner
        :title="
            $t('common.impersonation_active', {
                label: foundJob?.grades.find((g) => g.grade === jobGrade)?.label ?? jobGrade,
                number: jobGrade,
            })
        "
        icon="i-mdi-drama-masks"
        :ui="{ root: 'absolute top-0', container: 'h-5', title: 'text-xs', icon: 'size-4', close: 'text-xs' }"
        :close="{
            icon: undefined,
            trailingIcon: 'i-mdi-exit-run',
            label: $t('common.stop_impersonation'),
            ui: {
                trailingIcon: 'size-4',
            },
        }"
        @close="() => (isSuperuser ? setSuperuserMode(false) : impersonateJob(-1))"
    />
</template>
