<script setup lang="ts">
import type { ButtonProps } from '@nuxt/ui';
import { useSettingsStore } from '~/stores/settings';
import { isRoute } from '~/utils/route';

const { isHelpSlideoverOpen } = useDashboard();

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const { guides, startGuide } = useHelpGuides();
const { features } = useAppFeatures();
const route = useRoute();

const shortcuts = ref<boolean>(false);
const query = ref('');

const links = computed(() =>
    (
        [
            {
                label: t('common.shortcuts'),
                icon: 'i-mdi-key',
                trailingIcon: 'i-mdi-arrow-right',
                color: 'slate',
                onClick: () => {
                    shortcuts.value = true;
                },
            },
            !nuiEnabled.value
                ? {
                      label: t('common.help'),
                      icon: 'i-mdi-book-open-blank-variant-outline',
                      trailingIcon: 'i-mdi-external-link',
                      to: generateDerefURL('https://fivenet.app/getting-started'),
                      external: true,
                  }
                : undefined,
        ] as ButtonProps[]
    ).flatMap((item) => (item !== undefined ? [item] : [])),
);

const categories = computed(() => [
    {
        title: t('commandpalette.categories.general'),
        items: [
            { kbds: ['CTRL', 'K'], name: t('common.commandpalette') },
            { kbds: ['B'], name: t('common.notification', 2) },
            { kbds: ['?'], name: t('common.help') },
            { kbds: ['/'], name: t('common.search') },
        ],
    },
    {
        title: t('commandpalette.categories.navigation'),
        items: [
            { kbds: ['G', 'H'], name: t('common.goto_item', [t('common.overview')]) },
            { kbds: ['G', 'E'], name: t('common.goto_item', [t('common.mail')]) },
            { kbds: ['G', 'C'], name: t('common.goto_item', [t('common.citizen', 2)]) },
            { kbds: ['G', 'V'], name: t('common.goto_item', [t('common.vehicle', 2)]) },
            { kbds: ['G', 'D'], name: t('common.goto_item', [t('common.document', 2)]) },
            { kbds: ['G', 'J'], name: t('common.goto_item', [t('common.job')]) },
            { kbds: ['G', 'K'], name: t('common.goto_item', [t('common.calendar')]) },
            { kbds: ['G', 'Q'], name: t('common.goto_item', [t('common.qualification', 2)]) },
            { kbds: ['G', 'M'], name: t('common.goto_item', [t('common.livemap')]) },
            { kbds: ['G', 'W'], name: t('common.goto_item', [t('common.dispatch_center')]) },
            { kbds: ['G', 'L'], name: t('common.goto_item', [t('common.wiki')]) },
            { kbds: ['G', 'S'], name: t('common.goto_item', [t('common.control_panel')]) },
        ],
    },
    {
        title: t('pages.citizens.id.title'),
        items: [
            {
                kbds: ['C', 'W'],
                name: `${t('common.dialog')}: ${t('components.citizens.CitizenInfoProfile.revoke_wanted')}/ ${t('components.citizens.CitizenInfoProfile.set_wanted')}`,
            },
            { kbds: ['C', 'J'], name: `${t('common.dialog')}: ${t('components.citizens.CitizenInfoProfile.set_job')}` },
            {
                kbds: ['C', 'P'],
                name: `${t('common.dialog')}: ${t('components.citizens.CitizenInfoProfile.set_traffic_points')}`,
            },
            {
                kbds: ['C', 'M'],
                name: `${t('common.dialog')}: ${t('components.citizens.CitizenInfoProfile.set_mugshot')}`,
            },
            {
                kbds: ['C', 'D'],
                name: `${t('common.dialog')}: ${t('components.citizens.CitizenInfoProfile.create_new_document')}`,
            },
        ],
    },
    {
        title: t('common.document', 2),
        items: [
            { kbds: ['D', 'T'], name: `${t('common.open', 1)}/ ${t('common.close')}` },
            { kbds: ['D', 'E'], name: t('common.edit') },
            { kbds: ['D', 'R'], name: t('common.request', 2) },
            { kbds: ['D', 'A'], name: t('common.approvals', 2) },
        ],
    },
    {
        title: t('common.livemap'),
        items: [
            { kbds: ['M', 'D'], name: `${t('common.dialog')}: ${t('components.dispatch.take_dispatch.title')}` },
            { kbds: ['M', 'H'], name: `${t('common.mark')}: ${t('common.department_postal')}` },
            { kbds: ['C', 'U'], name: t('components.dispatch.update_unit_status.title') },
            { kbds: ['C', 'D'], name: t('components.dispatch.update_dispatch_status.title') },
        ],
    },
    {
        title: t('common.dispatch_center'),
        items: [{ kbds: ['C', 'Q'], name: `${t('common.dispatch_center')}: ${t('common.join')}/ ${t('common.leave')}` }],
    },
    {
        title: t('common.mail'),
        items: [
            { kbds: ['↑'], name: t('components.mailer.prev_thread') },
            { kbds: ['↓'], name: t('components.mailer.next_thread') },
        ],
    },
]);

const filteredCategories = computed(() => {
    return categories.value
        .map((category) => ({
            title: category.title,
            items: category.items.filter((item) => item.name.search(new RegExp(query.value, 'i')) !== -1),
        }))
        .filter((category) => !!category.items.length);
});

const guideFeatureGroups = computed(() =>
    features.value
        .map((feature) => ({
            id: feature.id,
            value: feature.id,
            label: feature.labelCount === undefined ? t(feature.label) : t(feature.label, feature.labelCount),
            icon: feature.icon,
            guides: guides.value.filter((guide) => guide.featureId === feature.id),
            ui: {
                trigger: (feature.activePaths ?? [feature.to]).some((path) => isRoute(route.path, path))
                    ? 'text-primary'
                    : undefined,
            },
        }))
        .filter((feature) => feature.guides.length > 0),
);
</script>

<template>
    <USlideover v-model:open="isHelpSlideoverOpen" :title="shortcuts ? $t('common.shortcuts') : $t('common.help')">
        <template #actions>
            <UTooltip v-if="shortcuts" :text="$t('common.back')">
                <UButton color="gray" icon="i-mdi-arrow-back" size="sm" variant="ghost" @click="shortcuts = false">
                    <span class="hidden truncate sm:block">
                        {{ $t('common.back') }}
                    </span>
                </UButton>
            </UTooltip>
        </template>

        <template #body>
            <div v-if="shortcuts" class="space-y-6">
                <UInput
                    v-model="query"
                    class="w-full"
                    icon="i-mdi-search"
                    :placeholder="$t('common.search_field')"
                    autofocus
                    color="neutral"
                />

                <USeparator />

                <div v-for="(category, index) in filteredCategories" :key="index">
                    <p class="mb-3 text-sm font-semibold text-highlighted">
                        {{ category.title }}
                    </p>

                    <div class="space-y-2">
                        <div v-for="(item, i) in category.items" :key="i" class="flex items-center justify-between">
                            <span class="text-sm text-muted">{{ item.name }}</span>

                            <div class="flex shrink-0 items-center justify-end gap-0.5">
                                <UKbd v-for="(kbd, j) in item.kbds" :key="j" :value="kbd" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div v-else class="flex flex-col gap-y-2">
                <template v-if="guides.length > 0">
                    <div class="space-y-1">
                        <p class="text-sm font-semibold text-highlighted">{{ $t('help_guides.title') }}</p>
                        <p class="text-sm text-muted">{{ $t('help_guides.description') }}</p>
                    </div>

                    <UAccordion
                        :items="guideFeatureGroups"
                        type="multiple"
                        :default-value="guideFeatureGroups.map((feature) => feature.id)"
                        :unmount-on-hide="true"
                        :ui="{ body: 'px-0 pt-0 pb-2' }"
                    >
                        <template #content="{ item: feature }">
                            <div class="grid grid-cols-1 gap-2 py-2 sm:grid-cols-2">
                                <UButton
                                    v-for="guide in feature.guides"
                                    :key="guide.id"
                                    color="neutral"
                                    :icon="guide.icon"
                                    :label="guide.title"
                                    :description="guide.description"
                                    variant="soft"
                                    :ui="{ base: 'flex-col', label: 'line-clamp-2 whitespace-normal overflow-visible' }"
                                    @click="startGuide(guide.id)"
                                />
                            </div>
                        </template>
                    </UAccordion>
                </template>
            </div>
        </template>

        <template #footer>
            <div class="flex w-full flex-col gap-y-2">
                <UButton v-for="(link, index) in links" :key="index" v-bind="link" />
            </div>
        </template>
    </USlideover>
</template>
