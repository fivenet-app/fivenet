<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useCompletorStore } from '~/store/completor';
import LawBookEntry from '~/components/rector/laws/LawBookEntry.vue';
import type { Law } from '~~/gen/ts/resources/laws/laws';
import GenericContainer from '~/components/partials/elements/GenericContainer.vue';

const completorStore = useCompletorStore();

const { data: lawBooks, pending, refresh, error } = useLazyAsyncData(`lawbooks`, () => completorStore.listLawBooks());

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
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="w-full sm:flex-auto">
                    <UButton
                        class="w-full rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                        @click="addLawBook"
                    >
                        {{ $t('pages.rector.laws.add_new_law_book') }}
                    </UButton>
                </div>
            </div>
            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 py-2 align-middle">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.law', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.law', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="lawBooks === null || lawBooks.length === 0"
                            icon="i-mdi-gavel"
                            :type="$t('common.law', 2)"
                        />
                        <template v-else>
                            <ul role="list" class="space-y-3">
                                <li v-for="(book, idx) in lawBooks" :key="book.id">
                                    <GenericContainer>
                                        <LawBookEntry
                                            v-model="lawBooks[idx]"
                                            v-model:laws="lawBooks[idx].laws"
                                            :start-in-edit="parseInt(book.id) < 0"
                                            @update:law="updateLaw($event)"
                                            @deleted="deletedLawBook($event)"
                                        />
                                    </GenericContainer>
                                </li>
                            </ul>
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
