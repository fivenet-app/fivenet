<script setup lang="ts">
const { isHelpSlideoverOpen } = useDashboard();
const { metaSymbol } = useShortcuts();

const { t } = useI18n();

const shortcuts = ref(false);
const query = ref('');

const links = [
    {
        label: t('common.shortcuts'),
        icon: 'i-heroicons-key',
        trailingIcon: 'i-heroicons-arrow-right-20-solid',
        color: 'gray',
        onClick: () => {
            shortcuts.value = true;
        },
    },
];

const categories = computed(() => [
    {
        title: t('command_palette.categories.general'),
        items: [
            { shortcuts: [metaSymbol.value, 'K'], name: t('common.command_palette') },
            { shortcuts: ['N'], name: t('common.notification', 2) },
            { shortcuts: ['?'], name: t('common.help') },
            { shortcuts: ['/'], name: t('common.search') },
        ],
    },
    {
        title: t('command_palette.categories.navigation'),
        items: [
            { shortcuts: ['G', 'H'], name: t('common.goto_item', [t('common.home')]) },
            { shortcuts: ['G', 'C'], name: t('common.goto_item', [t('common.citizen', 2)]) },
            { shortcuts: ['G', 'V'], name: t('common.goto_item', [t('common.vehicle', 2)]) },
            { shortcuts: ['G', 'D'], name: t('common.goto_item', [t('common.document', 2)]) },
            { shortcuts: ['G', 'J'], name: t('common.goto_item', [t('common.job')]) },
            { shortcuts: ['G', 'M'], name: t('common.goto_item', [t('common.livemap')]) },
            { shortcuts: ['G', 'W'], name: t('common.goto_item', [t('common.dispatch_center')]) },
            { shortcuts: ['G', 'P'], name: t('common.goto_item', [t('common.control_panel')]) },
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
                icon="i-heroicons-arrow-left-20-solid"
                @click="shortcuts = false"
            />

            {{ shortcuts ? 'Shortcuts' : 'Help & Support' }}
        </template>

        <div v-if="shortcuts" class="space-y-6">
            <UInput v-model="query" icon="i-heroicons-magnifying-glass" placeholder="Search..." autofocus color="gray" />

            <div v-for="(category, index) in filteredCategories" :key="index">
                <p class="mb-3 text-sm text-gray-900 dark:text-white font-semibold">
                    {{ category.title }}
                </p>

                <div class="space-y-2">
                    <div v-for="(item, i) in category.items" :key="i" class="flex items-center justify-between">
                        <span class="text-sm text-gray-500 dark:text-gray-400">{{ item.name }}</span>

                        <div class="flex items-center justify-end flex-shrink-0 gap-0.5">
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
