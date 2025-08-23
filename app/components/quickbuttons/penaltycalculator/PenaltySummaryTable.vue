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

const leeway = computed(() => props.reduction / 100);
</script>

<template>
    <UButton v-if="selectedLaws.length === 0" class="relative block w-full p-4 text-center" disabled variant="outline">
        <UIcon class="mx-auto size-12" name="i-mdi-calculator" />
        <span class="mt-2 block text-sm font-semibold">
            {{ $t('common.none_selected', [`${$t('common.crime')}`]) }}
        </span>
    </UButton>

    <UTable v-else class="divide-base-600 max-w-full divide-y" :columns="columns" :rows="selectedLaws">
        <template #law-data="{ row: law }">
            <div class="inline-flex items-center gap-2">
                <p class="text-highlighted whitespace-pre-line">
                    {{ getNameForLawBookId(law.law.lawbookId) }} - {{ law.law.name }}
                </p>

                <UTooltip v-if="law.law.hint" :text="law.law.hint">
                    <UIcon class="size-5" name="i-mdi-information-outline" />
                </UTooltip>
            </div>
        </template>

        <template #fine-data="{ row: law }">
            ${{ law.law.fine ? law.law.fine * law.count : 0 }}
            <span v-if="leeway > 0 && law.law.fine * law.count > 0">
                ($-{{ (law.law.fine * law.count * leeway).toFixed(0) }})
            </span>
        </template>

        <template #detentionTime-data="{ row: law }">
            {{ law.law.detentionTime ? law.law.detentionTime * law.count : 0 }}
            <span v-if="leeway > 0 && law.law.detentionTime * law.count > 0">
                (-{{ (law.law.detentionTime * law.count * leeway).toFixed(0) }})
            </span>
        </template>

        <template #trafficInfractionPoints-data="{ row: law }">
            {{ law.law.stvoPoints ? law.law.stvoPoints * law.count : 0 }}
            <span v-if="leeway > 0 && law.law.stvoPoints * law.count > 0">
                (-{{ (law.law.stvoPoints * law.count * leeway).toFixed(0) }})
            </span>
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
