<script lang="ts" setup>
import { GavelIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useCompletorStore } from '~/store/completor';
import LawBookEntry from '~/components/rector/laws/LawBookEntry.vue';
import type { Law } from '~~/gen/ts/resources/laws/laws';

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

function updateLaw(law: Law): void {
    const book = lawBooks.value?.find((b) => b.id === law.lawbookId);
    if (book === undefined) {
        return;
    }

    const idx = book?.laws.findIndex((l) => l.id === law.id);
    if (idx === -1) {
        return;
    }

    book.laws[idx] = law;
}
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto w-full">
                    <button
                        type="button"
                        class="w-full px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                        @click="addLawBook"
                    >
                        {{ $t('pages.rector.laws.add_new_law_book') }}
                    </button>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="mx-0 -my-2 overflow-x-auto">
                    <div class="inline-block min-w-full py-2 align-middle px-1">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.vehicle', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.vehicle', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="lawBooks === null || lawBooks.length === 0"
                            :icon="GavelIcon"
                            :type="$t('common.law', 2)"
                        />
                        <div v-else>
                            <ul role="list" class="space-y-3 divide-base-600 divide-y">
                                <li v-for="(book, idx) in lawBooks" :key="book.id">
                                    <LawBookEntry
                                        v-model="lawBooks[idx]"
                                        v-model:laws="lawBooks[idx].laws"
                                        :start-in-edit="parseInt(book.id, 10) < 0"
                                        @update:law="updateLaw($event)"
                                        @deleted="deletedLawBook($event)"
                                    />
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
