<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import LawBookEntry from '~/components/rector/laws/LawBookEntry.vue';
import { useCompletorStore } from '~/store/completor';
import type { Law } from '~~/gen/ts/resources/laws/laws';

const completorStore = useCompletorStore();

const { data: lawBooks, pending: loading, refresh, error } = useLazyAsyncData(`lawbooks`, () => completorStore.listLawBooks());

function deletedLawBook(id: string): void {
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
        id: lastNewId.value.toString(),
        name: '',
        laws: [],
    });
    lastNewId.value--;
}

function updateLaw(event: { id: string; law: Law }): void {
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
    <UDashboardNavbar :title="$t('pages.rector.laws.title')">
        <template #right>
            <UButton color="black" icon="i-mdi-arrow-back" to="/rector">
                {{ $t('common.back') }}
            </UButton>

            <UButton color="gray" trailing-icon="i-mdi-plus" @click="addLawBook">
                {{ $t('pages.rector.laws.add_new_law_book') }}
            </UButton>
        </template>
    </UDashboardNavbar>

    <UDashboardPanelContent>
        <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.law', 2)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.law', 2)])" :retry="refresh" />
        <DataNoDataBlock v-else-if="!lawBooks || lawBooks.length === 0" icon="i-mdi-gavel" :type="$t('common.law', 2)" />

        <template v-else>
            <ul role="list" class="space-y-3">
                <li v-for="(book, idx) in lawBooks" :key="book.id">
                    <LawBookEntry
                        v-model="lawBooks[idx]"
                        v-model:laws="book.laws"
                        :start-in-edit="parseInt(book.id) < 0"
                        @update:law="updateLaw($event)"
                        @deleted="deletedLawBook($event)"
                    />
                </li>
            </ul>
        </template>
    </UDashboardPanelContent>
</template>
