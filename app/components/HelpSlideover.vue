<script setup lang="ts">
import type { ButtonColor } from '#ui/types';
import { useSettingsStore } from '~/store/settings';

const { isHelpSlideoverOpen } = useDashboard();
const { metaSymbol } = useShortcuts();

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const shortcuts = ref(false);
const query = ref('');

const links = computed(() =>
    [
        {
            label: t('common.shortcuts'),
            icon: 'i-mdi-key',
            trailingIcon: 'i-mdi-arrow-right',
            color: 'gray' as ButtonColor,
            onClick: () => {
                shortcuts.value = true;
            },
        },
        !nuiEnabled.value
            ? {
                  label: t('common.help'),
                  icon: 'i-mdi-book-open-blank-variant-outline',
                  trailingIcon: 'i-mdi-external-link',
                  to: generateDerefedURL('https://fivenet.app/getting-started'),
                  external: true,
              }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const categories = computed(() => [
    {
        title: t('commandpalette.categories.general'),
        items: [
            { shortcuts: [metaSymbol.value, 'K'], name: t('common.commandpalette') },
            { shortcuts: ['B'], name: t('common.notification', 2) },
            { shortcuts: ['?'], name: t('common.help') },
            { shortcuts: ['/'], name: t('common.search') },
        ],
    },
    {
        title: t('commandpalette.categories.navigation'),
        items: [
            { shortcuts: ['G', 'H'], name: t('common.goto_item', [t('common.overview')]) },
            { shortcuts: ['G', 'E'], name: t('common.goto_item', [t('common.mail')]) },
            { shortcuts: ['G', 'C'], name: t('common.goto_item', [t('common.citizen', 2)]) },
            { shortcuts: ['G', 'V'], name: t('common.goto_item', [t('common.vehicle', 2)]) },
            { shortcuts: ['G', 'D'], name: t('common.goto_item', [t('common.document', 2)]) },
            { shortcuts: ['G', 'J'], name: t('common.goto_item', [t('common.job')]) },
            { shortcuts: ['G', 'K'], name: t('common.goto_item', [t('common.calendar')]) },
            { shortcuts: ['G', 'Q'], name: t('common.goto_item', [t('common.qualification', 2)]) },
            { shortcuts: ['G', 'M'], name: t('common.goto_item', [t('common.livemap')]) },
            { shortcuts: ['G', 'W'], name: t('common.goto_item', [t('common.dispatch_center')]) },
            { shortcuts: ['G', 'L'], name: t('common.goto_item', [t('common.wiki')]) },
            { shortcuts: ['G', 'I'], name: t('common.goto_item', [t('common.internet')]) },
            { shortcuts: ['G', 'P'], name: t('common.goto_item', [t('common.control_panel')]) },
        ],
    },
    {
        title: t('pages.citizens.id.title'),
        items: [
            {
                shortcuts: ['C', 'W'],
                name: `${t('common.dialog')}: ${t('components.citizens.CitizenInfoProfile.revoke_wanted')}/ ${t('components.citizens.CitizenInfoProfile.set_wanted')}`,
            },
            { shortcuts: ['C', 'J'], name: `${t('common.dialog')}: ${t('components.citizens.CitizenInfoProfile.set_job')}` },
            {
                shortcuts: ['C', 'P'],
                name: `${t('common.dialog')}: ${t('components.citizens.CitizenInfoProfile.set_traffic_points')}`,
            },
            {
                shortcuts: ['C', 'M'],
                name: `${t('common.dialog')}: ${t('components.citizens.CitizenInfoProfile.set_mug_shot')}`,
            },
            {
                shortcuts: ['C', 'D'],
                name: `${t('common.dialog')}: ${t('components.citizens.CitizenInfoProfile.create_new_document')}`,
            },
        ],
    },
    {
        title: t('common.document', 2),
        items: [
            { shortcuts: ['D', 'T'], name: `${t('common.open', 1)}/ ${t('common.close')}` },
            { shortcuts: ['D', 'E'], name: t('common.edit') },
            { shortcuts: ['D', 'R'], name: t('common.request', 2) },
        ],
    },
    {
        title: t('common.livemap'),
        items: [
            { shortcuts: ['M', 'D'], name: `${t('common.dialog')}: ${t('components.centrum.take_dispatch.title')}` },
            { shortcuts: ['M', 'H'], name: `${t('common.mark')}: ${t('common.department_postal')}` },
            { shortcuts: ['C', 'U'], name: t('components.centrum.update_unit_status.title') },
            { shortcuts: ['C', 'D'], name: t('components.centrum.update_dispatch_status.title') },
        ],
    },
    {
        title: t('common.dispatch_center'),
        items: [{ shortcuts: ['C', 'Q'], name: `${t('common.dispatch_center')}: ${t('common.join')}/ ${t('common.leave')}` }],
    },
    {
        title: t('common.mail'),
        items: [
            { shortcuts: ['↑'], name: t('components.mailer.prev_thread') },
            { shortcuts: ['↓'], name: t('components.mailer.next_thread') },
        ],
    },
]);

const filteredCategories = computed(() => {
    return categories.value
        .map((category) => ({
            title: category.title,
            items: category.items.filter((item) => {
                return item.name.search(new RegExp(query.value, 'i')) !== -1;
            }),
        }))
        .filter((category) => !!category.items.length);
});
</script>

<template>
    <UDashboardSlideover v-model="isHelpSlideoverOpen">
        <template #title>
            <UButton
                v-if="shortcuts"
                color="gray"
                variant="ghost"
                size="sm"
                icon="i-mdi-arrow-left"
                @click="shortcuts = false"
            />

            {{ shortcuts ? $t('common.shortcuts') : $t('common.help') }}
        </template>

        <div v-if="shortcuts" class="space-y-6">
            <UInput v-model="query" icon="i-mdi-search" :placeholder="$t('common.search_field')" autofocus color="gray" />

            <div v-for="(category, index) in filteredCategories" :key="index">
                <p class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">
                    {{ category.title }}
                </p>

                <div class="space-y-2">
                    <div v-for="(item, i) in category.items" :key="i" class="flex items-center justify-between">
                        <span class="text-sm text-gray-500 dark:text-gray-400">{{ item.name }}</span>

                        <div class="flex shrink-0 items-center justify-end gap-0.5">
                            <UKbd v-for="(shortcut, j) in item.shortcuts" :key="j">
                                {{ shortcut }}
                            </UKbd>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div v-else class="flex flex-col gap-y-3">
            <UButton v-for="(link, index) in links" :key="index" v-bind="link" />
        </div>
    </UDashboardSlideover>
</template>
