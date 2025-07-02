<script lang="ts" setup>
import EmailCreateForm from '~/components/mailer/EmailCreateForm.vue';
import EmailList from '~/components/mailer/EmailList.vue';
import { canAccess } from '~/components/mailer/helpers';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useMailerStore } from '~/stores/mailer';
import type { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access';

useHead({
    title: 'pages.mailer.manage.title',
});

definePageMeta({
    title: 'pages.mailer.manage.title',
    requiresAuth: true,
    permission: 'mailer.MailerService/ListEmails',
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

const page = useRouteQuery('page', '1', { transform: Number });

const pagination = ref<PaginationResponse | undefined>();

async function listEmails(): Promise<void> {
    const response = await mailerStore.listEmails(isSuperuser.value, calculateOffset(page.value, pagination.value));
    pagination.value = response.pagination;
}

watch(page, async () => await listEmails());

const route = useRoute();

onBeforeMount(async () => {
    await listEmails();

    if (route.query?.tab === 'new') {
        if (emails.value.length === 0 || !hasPrivateEmail.value) {
            return;
        }
    }
});

// Disable create form when email is selected
watch(selectedEmail, async () => {
    if (!selectedEmail.value) {
        return;
    }

    loading.value = true;
    await mailerStore.getEmail(selectedEmail.value.id).finally(() => (loading.value = false));

    creating.value = false;
});

const canCreate = computed(
    () =>
        can('mailer.MailerService/CreateOrUpdateEmail').value &&
        attr('mailer.MailerService/CreateOrUpdateEmail', 'Fields', 'Job').value,
);

const loading = ref(false);
const creating = ref(false);
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel v-if="route.query?.tab === 'new' || getPrivateEmail?.deactivated === true" id="maileremaillist" grow>
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
                    v-else-if="getPrivateEmail?.deactivated === true"
                    :title="$t('errors.MailerService.ErrEmailDisabled.title')"
                    :message="$t('errors.MailerService.ErrEmailDisabled.content')"
                />

                <div v-else class="flex flex-1 flex-col items-center">
                    <div class="flex flex-1 flex-col items-center justify-center gap-2 text-gray-400 dark:text-gray-500">
                        <UIcon class="h-32 w-32" name="i-mdi-email-multiple" />

                        <div class="text-center text-gray-900 dark:text-white">
                            <h3 class="text-lg font-bold">{{ $t('components.mailer.manage.title') }}</h3>
                            <p class="text-bas">{{ $t('components.mailer.manage.subtitle') }}</p>
                        </div>

                        <EmailCreateForm
                            v-if="can('mailer.MailerService/CreateOrUpdateEmail').value"
                            personal-email
                            hide-label
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
                            @click="
                                creating = !creating;
                                selectedEmail = undefined;
                            "
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
                                <UIcon class="h-32 w-32" name="i-mdi-email-multiple" />
                                <EmailCreateForm v-if="canCreate" :personal-email="false" @refresh="creating = false" />
                            </div>
                        </div>
                    </UDashboardPanelContent>
                </template>

                <template v-else-if="selectedEmail">
                    <UDashboardNavbar>
                        <template #toggle>
                            <UDashboardNavbarToggle icon="i-mdi-close" />

                            <UDivider class="mx-1.5 lg:hidden" orientation="vertical" />
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
                                    selectedEmail.id !== 0 &&
                                    selectedEmail.job !== undefined &&
                                    canAccess(selectedEmail.access, selectedEmail.userId, AccessLevel.MANAGE)
                                "
                                color="error"
                                trailing-icon="i-mdi-delete"
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
                        <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.mail')])" />
                        <EmailCreateForm
                            v-else
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
                            <PartialsBackButton to="/mail" />
                        </template>
                    </UDashboardNavbar>

                    <UDashboardPanelContent>
                        <div
                            class="hidden flex-1 flex-col items-center justify-center gap-2 text-gray-400 lg:flex dark:text-gray-500"
                        >
                            <UIcon class="h-32 w-32" name="i-mdi-email-multiple" />
                            <p>{{ $t('common.none_selected', [$t('common.mail')]) }}</p>
                        </div>
                    </UDashboardPanelContent>
                </template>
            </UDashboardPanel>
        </template>
    </UDashboardPage>
</template>
