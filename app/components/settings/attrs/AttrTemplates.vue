<script lang="ts" setup>
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import {
    gangJobTemplate,
    neutralJobTemplate,
    policeJobTemplate,
    type AttributeTemplate,
    type PermissionTemplate,
} from './templates';

defineEmits<{
    (e: 'apply-template', permissions: PermissionTemplate[], attributes: AttributeTemplate[]): void;
}>();

const overlay = useOverlay();

const { t } = useI18n();

const templates = ref<
    {
        title: string;
        description?: string;
        permissions: PermissionTemplate[];
        attributes: AttributeTemplate[];
    }[]
>([
    {
        title: t('components.settings.attr_view.templates.no_perms.title'),
        description: t('components.settings.attr_view.templates.no_perms.description'),
        permissions: [],
        attributes: [],
    },
    {
        title: t('components.settings.attr_view.templates.job_police.title'),
        description: t('components.settings.attr_view.templates.job_police.description'),
        permissions: policeJobTemplate.permissions,
        attributes: policeJobTemplate.attributes,
    },
    {
        title: t('components.settings.attr_view.templates.job_neutral.title'),
        description: t('components.settings.attr_view.templates.job_neutral.description'),
        permissions: neutralJobTemplate.permissions,
        attributes: neutralJobTemplate.attributes,
    },
    {
        title: t('components.settings.attr_view.templates.job_gang.title'),
        description: t('components.settings.attr_view.templates.job_gang.description'),
        permissions: gangJobTemplate.permissions,
        attributes: gangJobTemplate.attributes,
    },
    {
        title: t('components.settings.attr_view.templates.full_perms.title'),
        description: t('components.settings.attr_view.templates.full_perms.description'),
        permissions: [{ category: 'Superuser', name: 'Superuser' }],
        attributes: [],
    },
]);

const confirmModal = overlay.create(ConfirmModal);
</script>

<template>
    <UCollapsible class="mb-4">
        <UButton :label="$t('common.template', 2)" block color="neutral" variant="subtle" trailing-icon="i-mdi-chevron-down" />

        <template #content>
            <div class="flex flex-col gap-2 p-2">
                <UAlert
                    icon="i-mdi-information-outline"
                    color="warning"
                    variant="subtle"
                    :title="$t('components.settings.attr_view.template_note.title')"
                    :description="$t('components.settings.attr_view.template_note.description')"
                />

                <UPageGrid class="grid-cols-1 sm:grid-cols-1 lg:grid-cols-2">
                    <UPageCard
                        v-for="(template, idx) in templates"
                        :key="idx"
                        :title="template.title"
                        :description="template.description"
                    >
                        <template #footer>
                            <UButton
                                :label="$t('common.apply')"
                                color="red"
                                variant="outline"
                                icon="i-mdi-plus"
                                @click="
                                    confirmModal.open({
                                        title: $t('components.settings.attr_view.template_apply.title', {
                                            name: template.title,
                                        }),
                                        description: $t('components.settings.attr_view.template_apply.description'),
                                        confirm: () => $emit('apply-template', template.permissions, template.attributes),
                                    })
                                "
                            />
                        </template>
                    </UPageCard>
                </UPageGrid>
            </div>
        </template>
    </UCollapsible>
</template>
