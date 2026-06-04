<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import LawBookEntry from '~/components/settings/laws/LawBookEntry.vue';
import { getSettingsLawsClient } from '~~/gen/ts/clients';
import type { Law, LawBook } from '~~/gen/ts/resources/laws/laws';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import { useDraggable } from 'vue-draggable-plus';

useHead({
    title: 'pages.settings.laws.title',
});

definePageMeta({
    title: 'pages.settings.laws.title',
    requiresAuth: true,
    permission: 'settings.LawsService/CreateOrUpdateLawBook',
});

const { isSuperuser } = useAuth();

const notifications = useNotificationsStore();

const settingsLawsClient = await getSettingsLawsClient();

const {
    data: lawBooks,
    status,
    refresh,
    error,
} = useLazyAsyncData(`lawbooks`, () => settingsLawsClient.listLawBooks({}).then((resp) => resp.response.books));

const bookOrderChanged = ref(false);
const lawBookListRef = useTemplateRef('lawBookListRef');

const { moveUp: moveLawBookUp, moveDown: moveLawBookDown } = useListReorder(lawBooks, {
    onMove: () => {
        bookOrderChanged.value = true;
    },
});

function deletedLawBook(id: number, deletedAt?: Timestamp): void {
    if (!lawBooks.value) return;

    const idx = lawBooks.value.findIndex((b) => b.id === id);
    if (idx === -1) return;

    if (!isSuperuser.value) {
        lawBooks.value.splice(idx, 1);
    } else {
        lawBooks.value[idx]!.deletedAt = deletedAt;
    }
}

const lastNewId = ref(-1);

function addLawBook(): void {
    lawBooks.value?.unshift({
        id: lastNewId.value,
        name: '',
        laws: [],
        sortOrder: 0,
    });
    lastNewId.value--;
}

const canSaveBookOrder = computed(
    () =>
        bookOrderChanged.value &&
        (lawBooks.value?.every((book) => book.deletedAt !== undefined || book.id > 0) ?? false) &&
        (lawBooks.value?.some((book) => book.deletedAt === undefined) ?? false),
);

async function reorderLawBooks(): Promise<void> {
    if (!lawBooks.value) return;

    const activeLawBooks = lawBooks.value.filter((book) => book.deletedAt === undefined);
    if (!activeLawBooks.every((book) => book.id > 0)) return;

    try {
        const call = settingsLawsClient.reorderLawBooks({
            lawBookIds: activeLawBooks.map((book) => book.id),
        });
        await call;

        bookOrderChanged.value = false;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function findLawBookByLawId(id: number): { book: LawBook; index: number } | undefined {
    if (!lawBooks.value) return undefined;

    for (const book of lawBooks.value) {
        const index = book.laws.findIndex((law) => law.id === id);
        if (index !== -1) {
            return { book, index };
        }
    }

    return undefined;
}

function updateLaw(event: { id: number; law: Law }): void {
    if (!lawBooks.value) return;

    const source = findLawBookByLawId(event.id) ?? findLawBookByLawId(event.law.id);
    const destinationBook = lawBooks.value.find((book) => book.id === event.law.lawbookId);
    if (destinationBook === undefined) return;

    if (source && source.book.id === destinationBook.id) {
        source.book.laws[source.index] = event.law;
        return;
    }

    if (source) {
        source.book.laws.splice(source.index, 1);
    }

    const destinationIndex = destinationBook.laws.findIndex((law) => law.id === event.id || law.id === event.law.id);
    if (destinationIndex === -1) {
        destinationBook.laws.push(event.law);
    } else {
        destinationBook.laws[destinationIndex] = event.law;
    }
}

onMounted(() => {
    useDraggable(lawBookListRef, lawBooks, {
        animation: 150,
        handle: '.law-book-handle',
        draggable: 'li',
        onUpdate: () => {
            bookOrderChanged.value = true;
        },
    });
});
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

                    <UTooltip v-if="bookOrderChanged" :text="$t('common.save', 1)">
                        <UButton
                            color="primary"
                            variant="outline"
                            icon="i-mdi-content-save"
                            :disabled="!canSaveBookOrder"
                            @click="() => reorderLawBooks()"
                        />
                    </UTooltip>

                    <RefreshButton :loading="isRequestPending(status)" @click="() => refresh()" />

                    <UButton
                        color="neutral"
                        variant="outline"
                        trailing-icon="i-mdi-plus"
                        :label="$t('pages.settings.laws.add_new_law_book')"
                        @click="addLawBook"
                    />
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

            <ul v-else ref="lawBookListRef" class="space-y-3" role="list">
                <li v-for="(book, idx) in lawBooks" :key="book.id">
                    <LawBookEntry
                        :id="`book-${lawBooks[idx]?.id}`"
                        v-model="lawBooks[idx]"
                        v-model:laws="book.laws"
                        :idx="idx"
                        :law-books="lawBooks"
                        :start-in-edit="book.id < 0"
                        :move-book-up="moveLawBookUp"
                        :move-book-down="moveLawBookDown"
                        @update:law="updateLaw($event)"
                        @deleted="deletedLawBook($event.id, $event.deletedAt)"
                    />
                </li>
            </ul>
        </template>
    </UDashboardPanel>
</template>
