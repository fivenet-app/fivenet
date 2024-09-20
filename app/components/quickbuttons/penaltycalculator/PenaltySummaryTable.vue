<script lang="ts" setup>
import type { SelectedPenalty } from '~/components/quickbuttons/penaltycalculator/PenaltyCalculator.vue';
import type { LawBook } from '~~/gen/ts/resources/laws/laws';

const props = defineProps<{
    lawBooks: LawBook[];
    selectedLaws: SelectedPenalty[];
}>();

const { t } = useI18n();

function getNameForLawBookId(id: string): string | undefined {
    return props.lawBooks?.filter((b) => b.id === id)[0]?.name;
}

const columns = [
    {
        key: 'law',
        label: t('common.law'),
    },
    {
        key: 'fine',
        label: t('common.fine'),
    },
    {
        key: 'detentionTime',
        label: t('common.detention_time'),
    },
    {
        key: 'trafficInfractionPoints',
        label: t('common.traffic_infraction_points', 2),
    },
    {
        key: 'description',
        label: t('common.description'),
    },
    {
        key: 'count',
        label: t('common.count'),
    },
];
</script>

<template>
    <UButton
        v-if="selectedLaws.length === 0"
        disabled
        class="relative block w-full rounded-lg border border-dashed p-4 text-center"
    >
        <UIcon name="i-mdi-calculator" class="mx-auto size-12" />
        <span class="mt-2 block text-sm font-semibold">
            {{ $t('common.none_selected', [`${$t('common.crime')}`]) }}
        </span>
    </UButton>

    <UTable v-else :columns="columns" :rows="selectedLaws" class="max-w-full divide-y divide-base-600">
        <template #law-data="{ row: law }">
            <p class="whitespace-pre-line text-gray-900 dark:text-gray-300">
                {{ getNameForLawBookId(law.law.lawbookId) }} - {{ law.law.name }}
            </p>
        </template>
        <template #fine-data="{ row: law }"> ${{ law.law.fine ? law.law.fine * law.count : 0 }} </template>
        <template #detentionTime-data="{ row: law }">
            {{ law.law.detentionTime ? law.law.detentionTime * law.count : 0 }}
        </template>
        <template #trafficInfractionPoints-data="{ row: law }">
            {{ law.law.stvoPoints ? law.law.stvoPoints * law.count : 0 }}
        </template>
        <template #description-data="{ row: law }">
            <p class="line-clamp-2 w-full max-w-sm whitespace-normal break-all hover:line-clamp-none">
                {{ law.law.description }}
            </p>
        </template>
        <template #fine-count="{ row: law }">
            {{ law.count }}
        </template>
    </UTable>
</template>
