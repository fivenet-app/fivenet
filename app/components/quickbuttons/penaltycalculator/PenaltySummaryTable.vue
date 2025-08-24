<script lang="ts" setup>
import type { SelectedPenalty } from '~/components/quickbuttons/penaltycalculator/PenaltyCalculator.vue';
import type { LawBook } from '~~/gen/ts/resources/laws/laws';

const props = defineProps<{
    lawBooks: LawBook[];
    selectedLaws: SelectedPenalty[];
    reduction: number;
}>();

const { t } = useI18n();

function getNameForLawBookId(id: number): string | undefined {
    return props.lawBooks?.filter((b) => b.id === id)[0]?.name;
}

const columns = [
    {
        accessorKey: 'law',
        label: t('common.law'),
    },
    {
        accessorKey: 'fine',
        label: t('common.fine'),
    },
    {
        accessorKey: 'detentionTime',
        label: t('common.detention_time'),
    },
    {
        accessorKey: 'trafficInfractionPoints',
        label: t('common.traffic_infraction_points', 2),
    },
    {
        accessorKey: 'description',
        label: t('common.description'),
    },
    {
        accessorKey: 'count',
        label: t('common.count'),
    },
];

const leeway = computed(() => props.reduction / 100);
</script>

<template>
    <UButton v-if="selectedLaws.length === 0" class="relative block w-full p-4 text-center" disabled variant="outline">
        <UIcon class="mx-auto size-12" name="i-mdi-calculator" />
        <span class="mt-2 block text-sm font-semibold">
            {{ $t('common.none_selected', [`${$t('common.crime')}`]) }}
        </span>
    </UButton>

    <UTable v-else class="divide-base-600 max-w-full divide-y" :columns="columns" :data="selectedLaws">
        <template #law-cell="{ row: law }">
            <div class="inline-flex items-center gap-2">
                <p class="whitespace-pre-line text-highlighted">
                    {{ getNameForLawBookId(law.law.lawbookId) }} - {{ law.law.name }}
                </p>

                <UTooltip v-if="law.law.hint" :text="law.law.hint">
                    <UIcon class="size-5" name="i-mdi-information-outline" />
                </UTooltip>
            </div>
        </template>

        <template #fine-cell="{ row: law }">
            ${{ law.law.fine ? law.law.fine * law.count : 0 }}
            <span v-if="leeway > 0 && law.law.fine * law.count > 0">
                ($-{{ (law.law.fine * law.count * leeway).toFixed(0) }})
            </span>
        </template>

        <template #detentionTime-cell="{ row: law }">
            {{ law.law.detentionTime ? law.law.detentionTime * law.count : 0 }}
            <span v-if="leeway > 0 && law.law.detentionTime * law.count > 0">
                (-{{ (law.law.detentionTime * law.count * leeway).toFixed(0) }})
            </span>
        </template>

        <template #trafficInfractionPoints-cell="{ row: law }">
            {{ law.law.stvoPoints ? law.law.stvoPoints * law.count : 0 }}
            <span v-if="leeway > 0 && law.law.stvoPoints * law.count > 0">
                (-{{ (law.law.stvoPoints * law.count * leeway).toFixed(0) }})
            </span>
        </template>

        <template #description-cell="{ row: law }">
            <p class="line-clamp-2 w-full max-w-sm break-all whitespace-normal hover:line-clamp-none">
                {{ law.law.description }}
            </p>
        </template>

        <template #fine-count="{ row: law }">
            {{ law.count }}
        </template>
    </UTable>
</template>
