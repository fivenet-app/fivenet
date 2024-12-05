<script lang="ts" setup>
import EmailCreateForm from '~/components/mailer/EmailCreateForm.vue';
import EmailList from '~/components/mailer/EmailList.vue';
import { canAccess } from '~/components/mailer/helpers';
import BackButton from '~/components/partials/BackButton.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useMailerStore } from '~/store/mailer';
import type { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access';

useHead({
    title: 'pages.mailer.manage.title',
});
definePageMeta({
    title: 'pages.mailer.manage.title',
    requiresAuth: true,
    permission: 'MailerService.ListEmails',
});

const modal = useModal();

const mailerStore = useMailerStore();
const { emails, getPrivateEmail, hasPrivateEmail, loaded, error, selectedEmail } = storeToRefs(mailerStore);

const { attr, can, isSuperuser } = useAuth();

const isMailerPanelOpen = computed({
    get() {
        return !!selectedEmail.value || creating.value;
    },
    set(value: boolean) {
        if (!value) {
            selectedEmail.value = undefined;
            creating.value = false;
        }
    },
});

// Set email as query param for persistence between reloads
const route = useRoute();
const router = useRouter();

watch(selectedEmail, () => {
    if (!selectedEmail.value) {
        router.replace({ query: {} });
    } else {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { email: selectedEmail.value.id }, hash: '#' });
    }
});

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() => (pagination.value?.pageSize ?? 20) * (page.value - 1));
const pagination = ref<PaginationResponse | undefined>();

async function listEmails(): Promise<void> {
    const response = await mailerStore.listEmails(isSuperuser.value, offset.value);
    pagination.value = response.pagination;
}

watch(offset, async () => await listEmails());

onBeforeMount(async () => {
    await listEmails();

    if (!route.query.email) {
        return;
    }

    selectedEmail.value = await mailerStore.getEmail(route.query.email as string);
});

const canCreate = computed(
    () => can('MailerService.CreateOrUpdateEmail').value && attr('MailerService.CreateOrUpdateEmail', 'Fields', 'Job').value,
);

const creating = ref(false);
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel v-if="emails.length === 0 || !hasPrivateEmail" id="maileremaillist" grow>
            <UDashboardNavbar :title="$t('pages.mailer.manage.title')" />

            <UDashboardPanelContent>
                <DataErrorBlock
                    v-if="error"
                    :title="$t('common.unable_to_load', [$t('common.mail', 2)])"
                    :error="error"
                    :retry="async () => await listEmails()"
                />
                <DataPendingBlock v-else-if="!loaded" :message="$t('common.loading', [$t('common.mail', 2)])" />

                <DataErrorBlock
                    v-else-if="getPrivateEmail?.deactivated"
                    :title="$t('errors.MailerService.ErrEmailDisabled.title')"
                    :message="$t('errors.MailerService.ErrEmailDisabled.content')"
                />

                <div v-else class="flex flex-1 flex-col items-center">
                    <div class="flex flex-1 flex-col items-center justify-center gap-2 text-gray-400 dark:text-gray-500">
                        <UIcon name="i-mdi-email-multiple" class="h-32 w-32" />

                        <div class="text-center text-gray-900 dark:text-white">
                            <h3 class="text-lg font-bold">{{ $t('components.mailer.manage.title') }}</h3>
                            <p class="text-bas">{{ $t('components.mailer.manage.subtitle') }}</p>
                        </div>

                        <EmailCreateForm
                            v-if="can('MailerService.CreateOrUpdateEmail').value"
                            personal-email
                            hide-label
                            @refresh="async () => await navigateTo({ name: 'mail' })"
                        />
                    </div>
                </div>
            </UDashboardPanelContent>
        </UDashboardPanel>

        <template v-else>
            <UDashboardPanel id="maileremailslist" :width="450" :resizable="{ min: 325, max: 550 }">
                <UDashboardNavbar :title="$t('pages.mailer.manage.title')">
                    <template #right>
                        <UButton
                            v-if="canCreate"
                            :label="$t('common.create')"
                            trailing-icon="i-mdi-plus"
                            color="gray"
                            @click="creating = !creating"
                        />
                    </template>
                </UDashboardNavbar>

                <div class="relative flex-1 overflow-x-auto">
                    <EmailList v-model="selectedEmail" :emails="emails" :loaded="loaded">
                        <Pagination
                            v-if="emails.length > (pagination?.pageSize ?? 20)"
                            v-model="page"
                            :pagination="pagination"
                            :refresh="async () => await listEmails()"
                        />
                    </EmailList>
                </div>
            </UDashboardPanel>

            <UDashboardPanel id="maileremailsview" v-model="isMailerPanelOpen" collapsible grow side="right">
                <template v-if="creating">
                    <UDashboardNavbar :title="$t('pages.mailer.manage.title')">
                        <template #right>
                            <UButton
                                :label="$t('common.back')"
                                icon="i-mdi-arrow-back"
                                color="black"
                                @click="creating = false"
                            />
                        </template>
                    </UDashboardNavbar>

                    <UDashboardPanelContent>
                        <div class="flex flex-1 flex-col items-center">
                            <div
                                class="flex flex-1 flex-col items-center justify-center gap-2 text-gray-400 dark:text-gray-500"
                            >
                                <UIcon name="i-mdi-email-multiple" class="h-32 w-32" />
                                <EmailCreateForm v-if="canCreate" :personal-email="false" @refresh="creating = false" />
                            </div>
                        </div>
                    </UDashboardPanelContent>
                </template>

                <template v-else-if="selectedEmail">
                    <UDashboardNavbar>
                        <template #toggle>
                            <UDashboardNavbarToggle icon="i-mdi-close" />

                            <UDivider orientation="vertical" class="mx-1.5 lg:hidden" />
                        </template>

                        <template #right>
                            <UButton
                                class="hidden md:flex"
                                color="black"
                                icon="i-mdi-arrow-back"
                                @click="selectedEmail = undefined"
                            >
                                {{ $t('common.back') }}
                            </UButton>

                            <UButton
                                v-if="
                                    selectedEmail &&
                                    selectedEmail.id !== '0' &&
                                    selectedEmail.job !== undefined &&
                                    canAccess(selectedEmail.access, selectedEmail.userId, AccessLevel.MANAGE)
                                "
                                color="red"
                                trailing-icon="i-mdi-trash-can"
                                @click="
                                    modal.open(ConfirmModal, {
                                        confirm: async () =>
                                            selectedEmail?.id &&
                                            (await mailerStore.deleteEmail({
                                                id: selectedEmail.id,
                                            })),
                                    })
                                "
                            >
                                {{ $t('common.delete') }}
                            </UButton>
                        </template>
                    </UDashboardNavbar>

                    <UDashboardPanelContent>
                        <EmailCreateForm
                            v-model="selectedEmail"
                            :personal-email="selectedEmail.userId !== undefined"
                            :disabled="
                                !canAccess(selectedEmail.access, selectedEmail.userId, AccessLevel.MANAGE) ||
                                (!isSuperuser && selectedEmail.deactivated)
                            "
                        />
                    </UDashboardPanelContent>
                </template>

                <template v-else>
                    <UDashboardNavbar :title="$t('pages.mailer.manage.title')">
                        <template #right>
                            <BackButton to="/mail" />
                        </template>
                    </UDashboardNavbar>

                    <UDashboardPanelContent>
                        <div
                            class="hidden flex-1 flex-col items-center justify-center gap-2 text-gray-400 lg:flex dark:text-gray-500"
                        >
                            <UIcon name="i-mdi-email-multiple" class="h-32 w-32" />
                            <p>{{ $t('common.none_selected', [$t('common.mail')]) }}</p>
                        </div>
                    </UDashboardPanelContent>
                </template>
            </UDashboardPanel>
        </template>
    </UDashboardPage>
</template>
