<script lang="ts" setup>
import EmailCreateForm from '~/components/mailer/EmailCreateForm.vue';
import { useMailerStore } from '~/store/mailer';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';

const { isOpen } = useModal();

const { attr, can } = useAuth();

const mailerStore = useMailerStore();
const { selectedEmail } = storeToRefs(mailerStore);

const {
    data: emails,
    pending: loading,
    error,
    refresh,
} = useLazyAsyncData(`mailer-emails-jobonly`, () => mailerStore.listEmails());

const canCreate = computed(
    () => can('MailerService.CreateOrUpdateEmail').value && attr('MailerService.CreateOrUpdateEmail', 'Fields', 'Job').value,
);

// TODO
const creating = ref(false);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard
            :ui="{
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                base: 'flex flex-1 flex-col',
                body: { base: 'flex flex-1 flex-col' },
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.mail', 2) }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div class="flex max-w-screen-xl flex-col gap-2">
                <template v-if="!creating">
                    <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.mail', 2)])" />
                    <DataErrorBlock
                        v-else-if="error"
                        :title="$t('common.unable_to_load', [$t('common.mail', 2)])"
                        :retry="refresh"
                    />
                    <DataNoDataBlock
                        v-if="!emails || emails?.length === 1"
                        :type="$t('common.mail')"
                        icon="i-mdi-file-outline"
                        :actions="
                            canCreate
                                ? [{ label: $t('common.create'), icon: 'i-mdi-plus', click: () => (creating = true) }]
                                : []
                        "
                    />

                    <template v-else>
                        <UButton v-if="canCreate" :label="$t('common.create')" @click="creating = true" />

                        <UPageGrid :ui="{ wrapper: 'grid-cols-1 md:grid-cols-3 xl:grid-cols-4' }">
                            <UCard v-for="email in emails" :key="email.id">
                                <template #header>
                                    {{ email.email }}
                                </template>
                                {{
                                    (email?.label && email?.label !== '' ? email?.label : undefined) ??
                                    (email?.userId ? $t('common.personal_email') : undefined) ??
                                    $t('common.none', [$t('common.label')])
                                }}
                            </UCard>
                        </UPageGrid>
                    </template>
                </template>

                <div v-else>
                    <UButton :label="$t('common.cancel')" color="red" @click="creating = false" />

                    <EmailCreateForm />
                </div>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton block class="flex-1" color="black" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
