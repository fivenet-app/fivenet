<script setup lang="ts">
import type { ButtonProps } from '@nuxt/ui';
import { useSettingsStore } from '~/stores/settings';

const { isHelpSlideoverOpen } = useDashboard();

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const shortcuts = ref(false);
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
            { kbds: ['G', 'P'], name: t('common.goto_item', [t('common.control_panel')]) },
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
        ],
    },
    {
        title: t('common.livemap'),
        items: [
            { kbds: ['M', 'D'], name: `${t('common.dialog')}: ${t('components.centrum.take_dispatch.title')}` },
            { kbds: ['M', 'H'], name: `${t('common.mark')}: ${t('common.department_postal')}` },
            { kbds: ['C', 'U'], name: t('components.centrum.update_unit_status.title') },
            { kbds: ['C', 'D'], name: t('components.centrum.update_dispatch_status.title') },
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
            items: category.items.filter((item) => {
                return item.name.search(new RegExp(query.value, 'i')) !== -1;
            }),
        }))
        .filter((category) => !!category.items.length);
});
</script>

<template>
    <USlideover v-model:open="isHelpSlideoverOpen" :title="shortcuts ? $t('common.shortcuts') : $t('common.help')">
        <template #actions>
            <UButton
                v-if="shortcuts"
                color="neutral"
                variant="ghost"
                size="sm"
                icon="i-mdi-arrow-left"
                @click="shortcuts = false"
            />
        </template>

        <template #body>
            <div v-if="shortcuts" class="space-y-6">
                <UInput
                    v-model="query"
                    icon="i-mdi-search"
                    :placeholder="$t('common.search_field')"
                    autofocus
                    color="neutral"
                    class="w-full"
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

            <div v-else class="flex flex-col gap-y-3">
                <UButton v-for="(link, index) in links" :key="index" v-bind="link" />
            </div>
        </template>
    </USlideover>
</template>
