<script setup lang="ts">
import { UserActivity } from '@fivenet/gen/resources/users/users_pb';
import { QuestionMarkCircleIcon, BellAlertIcon, BellSnoozeIcon, AtSymbolIcon } from '@heroicons/vue/24/outline'
import { FunctionalComponent } from 'vue';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';

const props = defineProps<{
    activity: UserActivity;
}>();

const { t } = useI18n();

const icon = ref<FunctionalComponent>(QuestionMarkCircleIcon);
const iconColor = ref<string>('text-neutral');
const actionText = ref<string>(props.activity.getKey());
const actionValue = ref<string>(`${props.activity.getOldValue()} -> ${props.activity.getNewValue()}`);
const actionLink = ref<undefined | RoutesNamedLocations>();
const actionLinkText = ref<string>();

switch (props.activity.getKey()) {
    case 'UserProps.Wanted': {
        actionText.value = t('components.citizens.citizen_info_activity_feed_entry.set_citizen_as');

        if (props.activity.getNewValue() === 'true') {
            icon.value = BellAlertIcon;
            iconColor.value = 'text-error-400';
            actionValue.value = props.activity.getNewValue() === 'true' ? t('common.wanted') : `${t('common.not').toLowerCase()} ${t('common.wanted')}`;
        } else {
            icon.value = BellSnoozeIcon;
            iconColor.value = 'text-success-400';
            actionValue.value = props.activity.getNewValue() === 'true' ? t('common.wanted') : `${t('common.not').toLowerCase()} ${t('common.wanted')}`;
        }

        break;
    };
    case 'DocStore.Relation': {
        actionText.value = t('components.citizens.citizen_info_activity_feed_entry.document_relation');
        icon.value = AtSymbolIcon;
        actionLink.value = { name: 'documents-id', params: { id: 0 } };
        actionLinkText.value = t('common.document', 1);

        if (props.activity.getNewValue() !== '') {
            iconColor.value = 'text-info-600';
            actionLink.value.params.id = props.activity.getNewValue();
        } else if (props.activity.getOldValue() !== '') {
            iconColor.value = 'text-base-600';
            actionLink.value.params.id = props.activity.getOldValue();
        }

        break;
    };
}
</script>

<template>
    <div class="flex space-x-3">
        <div class="h-10 w-10 rounded-full flex items-center justify-center my-auto">
            <component :class="[iconColor, 'w-full h-full']" :is="icon" />
        </div>
        <div class="flex-1 space-y-1">
            <div class="flex items-center justify-between">
                <h3 class="text-sm font-medium text-neutral">
                    {{ activity.getSourceUser()?.getFirstname() }} {{ activity.getSourceUser()?.getLastname() }}
                </h3>
                <p class="text-sm text-gray-400">
                    {{ useLocaleTimeAgo(toDate(activity.getCreatedAt())!).value }}
                </p>
            </div>
            <p class="text-sm text-gray-300">{{ actionText }} <span class="font-bold">
                    <NuxtLink v-if="actionLink" :to="actionLink">{{ actionLinkText }}</NuxtLink>
                    <span v-else v-html="actionValue"></span>
                </span>
            </p>
        </div>
    </div>
</template>
