<script lang="ts" setup>
import MailerList from '~/components/mailer/MailerList.vue';
import MailerThread from '~/components/mailer/MailerThread.vue';
import ThreadCreateOrUpdateModal from '~/components/mailer/ThreadCreateOrUpdateModal.vue';
import { canAccessThread } from '~/components/mailer/helpers';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { mailerDB, useMailerStore } from '~/store/mailer';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access';

useHead({
    title: 'common.mailer',
});
definePageMeta({
    title: 'common.mailer',
    requiresAuth: true,
    permission: 'MailerService.ListThreads',
});

const { t } = useI18n();

const { can } = useAuth();

const modal = useModal();

const mailerStore = useMailerStore();
const { selectedThread } = storeToRefs(mailerStore);

const tabItems = [
    {
        label: t('common.all'),
    },
    {
        label: t('common.unread'),
    },
];
const selectedTab = ref(0);

const dropdownItems = computed(() =>
    [
        [
            can('MailerService.CreateOrUpdateThread').value &&
            canAccessThread(selectedThread.value?.access, selectedThread.value?.creator, AccessLevel.MANAGE)
                ? {
                      label: t('common.edit'),
                      icon: 'i-mdi-pencil-outline',
                      click: () => {
                          if (!selectedThread.value) {
                              return;
                          }

                          modal.open(ThreadCreateOrUpdateModal, {
                              thread: selectedThread.value,
                          });
                      },
                  }
                : {
                      label: t('common.leave'),
                      icon: 'i-mdi-pencil-outline',
                      click: () => {
                          if (!selectedThread.value) {
                              return;
                          }

                          modal.open(ConfirmModal, {
                              confirm: async () => mailerStore.leaveThread(selectedThread.value!.id),
                          });
                      },
                  },
        ].flatMap((item) => (item !== undefined ? [item] : [])),
        [
            can('MailerService.DeleteThread').value &&
            canAccessThread(selectedThread.value?.access, selectedThread.value?.creator, AccessLevel.ADMIN)
                ? {
                      label: t('common.delete'),
                      icon: 'i-mdi-trash-can-outline',
                      click: async () => {
                          if (!selectedThread.value) {
                              return;
                          }

                          modal.open(ConfirmModal, {
                              confirm: async () => mailerStore.deleteThread({ threadId: selectedThread.value!.id }),
                          });
                      },
                  }
                : undefined,
        ].flatMap((item) => (item !== undefined ? [item] : [])),
    ].flatMap((items) => (items.length > 0 ? [items] : [])),
);

onBeforeMount(async () => {
    const count = await mailerDB.threads.count();
    const call = getGRPCMailerClient().listThreads({
        pagination: {
            offset: 0,
        },
        after: count > 0 ? undefined : toTimestamp(),
    });
    const { response } = await call;

    mailerDB.threads.bulkPut(response.threads);
});

const threads = useDexieLiveQuery(() => mailerDB.threads.toArray().then((threads) => ({ threads, loaded: true })), {
    initialValue: { threads: [], loaded: false },
});

// Filter mails based on the selected tab
const filteredThreads = computed(() => {
    if (selectedTab.value === 1) {
        return threads.value.threads.filter((thread) => !!thread.userState?.lastRead);
    }

    return threads.value.threads;
});

const threadUserState = computed(() => selectedThread.value?.userState);

const isMailerPanelOpen = computed({
    get() {
        return !!selectedThread.value;
    },
    set(value: boolean) {
        if (!value) {
            selectedThread.value = undefined;
        }
    },
});

// Reset selected mail if it's not in the filtered mails
watch(filteredThreads, () => {
    if (!filteredThreads.value.find((thread) => thread.id === selectedThread.value?.id)) {
        selectedThread.value = undefined;
    }
});

// Set thread as query param for persistence between reloads
const route = useRoute();
const router = useRouter();

watch(selectedThread, () => {
    if (!selectedThread.value) {
        router.replace({ query: {} });
    } else {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { thread: selectedThread.value.id }, hash: '#' });
    }
});
onMounted(async () => {
    if (!route.query.thread) {
        return;
    }

    selectedThread.value = await mailerStore.getThread(route.query.thread as string);
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel id="mailerthreadlist" :width="425" :resizable="{ min: 300, max: 500 }">
            <UDashboardNavbar :title="$t('common.mailer')" :badge="filteredThreads.length">
                <template #right>
                    <UTabs
                        v-model="selectedTab"
                        :items="tabItems"
                        :ui="{ wrapper: '', list: { height: 'h-9', tab: { height: 'h-7', size: 'text-[13px]' } } }"
                    />
                </template>
            </UDashboardNavbar>

            <div class="relative flex-1 overflow-x-auto">
                <MailerList v-model="selectedThread" :threads="filteredThreads" :loaded="threads.loaded" />
            </div>

            <UDashboardToolbar class="flex justify-between border-t border-gray-200 px-3 py-3.5 dark:border-gray-700">
                <template #right>
                    <UButtonGroup class="inline-flex">
                        <UButton
                            v-if="can('MailerService.CreateOrUpdateThread').value"
                            color="gray"
                            trailing-icon="i-mdi-plus"
                            @click="() => modal.open(ThreadCreateOrUpdateModal, {})"
                        >
                            {{ $t('components.mailer.create_thread') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UDashboardToolbar>
        </UDashboardPanel>

        <UDashboardPanel id="mailerthreadview" v-model="isMailerPanelOpen" collapsible grow side="right">
            <template v-if="selectedThread">
                <UDashboardNavbar>
                    <template #toggle>
                        <UDashboardNavbarToggle icon="i-mdi-close" />

                        <UDivider orientation="vertical" class="mx-1.5 lg:hidden" />
                    </template>

                    <template #left>
                        <UTooltip :text="$t('components.mailer.mark_unread')">
                            <UButton
                                :icon="!threadUserState?.unread ? 'i-mdi-check-circle-outline' : 'i-mdi-check-circle'"
                                color="gray"
                                variant="ghost"
                                @click="
                                    async () =>
                                        (selectedThread!.userState = await mailerStore.setThreadUserState({
                                            threadId: selectedThread!.id,
                                            unread: !threadUserState?.unread,
                                        }))
                                "
                            />
                        </UTooltip>

                        <UTooltip :text="$t('components.mailer.mark_important')">
                            <UButton
                                :icon="!threadUserState?.important ? 'i-mdi-alert-circle-outline' : 'i-mdi-alert-circle'"
                                color="gray"
                                variant="ghost"
                                @click="
                                    async () =>
                                        (selectedThread!.userState = await mailerStore.setThreadUserState({
                                            threadId: selectedThread!.id,
                                            important: !threadUserState?.important,
                                        }))
                                "
                            />
                        </UTooltip>
                    </template>

                    <template #right>
                        <UTooltip :text="$t('components.mailer.star_thread')">
                            <UButton
                                :icon="!threadUserState?.favorite ? 'i-mdi-star-circle-outline' : 'i-mdi-star-circle'"
                                color="gray"
                                variant="ghost"
                                @click="
                                    async () =>
                                        (selectedThread!.userState = await mailerStore.setThreadUserState({
                                            threadId: selectedThread!.id,
                                            favorite: !threadUserState?.favorite,
                                        }))
                                "
                            />
                        </UTooltip>

                        <UTooltip :text="$t('components.mailer.mute_thread')">
                            <UButton
                                :icon="!threadUserState?.muted ? 'i-mdi-pause-circle-outline' : 'i-mdi-pause-circle'"
                                color="gray"
                                variant="ghost"
                                @click="
                                    async () =>
                                        (selectedThread!.userState = await mailerStore.setThreadUserState({
                                            threadId: selectedThread!.id,
                                            muted: !threadUserState?.muted,
                                        }))
                                "
                            />
                        </UTooltip>

                        <UDivider orientation="vertical" class="mx-1.5" />

                        <UDropdown :items="dropdownItems">
                            <UButton icon="i-mdi-ellipsis-vertical" color="gray" variant="ghost" />
                        </UDropdown>
                    </template>
                </UDashboardNavbar>

                <MailerThread :thread-id="selectedThread.id" />
            </template>
            <div v-else class="hidden flex-1 items-center justify-center lg:flex">
                <UIcon name="i-mdi-conversation" class="h-32 w-32 text-gray-400 dark:text-gray-500" />
            </div>
        </UDashboardPanel>
    </UDashboardPage>
</template>
