<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import LawBookEntry from '~/components/settings/laws/LawBookEntry.vue';
import { useCompletorStore } from '~/stores/completor';
import type { Law } from '~~/gen/ts/resources/laws/laws';

useHead({
    title: 'pages.settings.laws.title',
});

definePageMeta({
    title: 'pages.settings.laws.title',
    requiresAuth: true,
    permission: 'settings.LawsService/CreateOrUpdateLawBook',
});

const completorStore = useCompletorStore();

const { data: lawBooks, status, refresh, error } = useLazyAsyncData(`lawbooks`, () => completorStore.listLawBooks());

function deletedLawBook(id: number): void {
    if (!lawBooks.value) {
        return;
    }

    const idx = lawBooks.value.findIndex((b) => b.id === id);
    if (idx > -1) {
        lawBooks.value.splice(idx, 1);
    }
}

const lastNewId = ref(-1);

function addLawBook(): void {
    lawBooks.value?.unshift({
        id: lastNewId.value,
        name: '',
        laws: [],
    });
    lastNewId.value--;
}

function updateLaw(event: { id: number; law: Law }): void {
    const book = lawBooks.value?.find((b) => b.id === event.law.lawbookId);
    if (book === undefined) {
        return;
    }

    const idx = book?.laws.findIndex((l) => l.id === event.law.id || l.id === event.id);
    if (idx === -1) {
        return;
    }

    book.laws[idx] = event.law;
}
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('pages.settings.laws.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/settings" />

                    <UButton color="neutral" trailing-icon="i-mdi-plus" @click="addLawBook">
                        {{ $t('pages.settings.laws.add_new_law_book') }}
                    </UButton>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.law', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.law', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!lawBooks || lawBooks.length === 0" icon="i-mdi-gavel" :type="$t('common.law', 2)" />

            <template v-else>
                <ul class="space-y-3" role="list">
                    <li v-for="(book, idx) in lawBooks" :key="book.id">
                        <LawBookEntry
                            :id="`book-${lawBooks[idx]?.id}`"
                            v-model="lawBooks[idx]"
                            v-model:laws="book.laws"
                            :start-in-edit="book.id < 0"
                            @update:law="updateLaw($event)"
                            @deleted="deletedLawBook($event)"
                        />
                    </li>
                </ul>
            </template>
        </template>
    </UDashboardPanel>
</template>
