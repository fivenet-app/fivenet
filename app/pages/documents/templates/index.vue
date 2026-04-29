<script lang="ts" setup>
import { z } from 'zod';
import List from '~/components/documents/templates/List.vue';
import Pagination from '~/components/partials/Pagination.vue';
import type { TemplateShort } from '~~/gen/ts/resources/documents/templates/templates';

useHead({
    title: 'pages.documents.templates.title',
});

definePageMeta({
    title: 'pages.documents.templates.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService/ListTemplates',
});

const { can } = useAuth();

const schema = z.object({
    title: z.string().optional(),
});

const query = useSearchForm('documents-templates', schema);

async function selected(t: TemplateShort | undefined): Promise<void> {
    if (!t) return;

    await navigateTo({ name: 'documents-templates-id', params: { id: t.id } });
}

const templatesListRef = useTemplateRef('templatesListRef');

const inputRef = useTemplateRef('inputRef');

defineShortcuts({
    '/': () => inputRef.value?.inputRef?.focus(),
});
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('pages.documents.templates.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton to="/documents" />

                    <UButton
                        v-if="can('TODOService/TODOMethod').value"
                        to="/documents/templates/forms"
                        icon="i-mdi-form"
                        truncate
                    >
                        <span class="hidden truncate sm:block">
                            {{ $t('common.form', 2) }}
                        </span>
                    </UButton>

                    <UTooltip v-if="can('documents.DocumentsService/CreateTemplate').value" :text="$t('common.create')">
                        <UButton to="/documents/templates/create" color="neutral" variant="outline" trailing-icon="i-mdi-plus">
                            <span class="hidden truncate sm:block">
                                {{ $t('common.template') }}
                            </span>
                        </UButton>
                    </UTooltip>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <UForm ref="formRef" class="my-2 flex w-full flex-1 flex-col gap-2" :schema="schema" :state="query">
                    <UFormField class="flex-1" name="title" :label="$t('common.search')">
                        <UInput
                            ref="inputRef"
                            v-model="query.title"
                            class="w-full"
                            type="text"
                            name="title"
                            :placeholder="$t('common.title')"
                            leading-icon="i-mdi-search"
                        >
                            <template #trailing>
                                <UKbd value="/" />
                            </template>
                        </UInput>
                    </UFormField>
                </UForm>
            </UDashboardToolbar>
        </template>

        <template #body>
            <List ref="templatesListRef" link :search-title="query.title" @selected="selected($event)" />
        </template>

        <template #footer>
            <Pagination
                :status="templatesListRef?.status ?? 'pending'"
                :refresh="templatesListRef?.refresh"
                hide-buttons
                hide-text
            />
        </template>
    </UDashboardPanel>
</template>
