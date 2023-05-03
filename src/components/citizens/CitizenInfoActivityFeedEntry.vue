<script setup lang="ts">
import { UserActivity } from '@fivenet/gen/resources/users/users_pb';
import { QuestionMarkCircleIcon, BellAlertIcon, BellSnoozeIcon } from '@heroicons/vue/24/outline'
import { FunctionalComponent } from 'vue';

const props = defineProps<{
    activity: UserActivity;
}>();

const { t } = useI18n();

const icon = ref<FunctionalComponent>(QuestionMarkCircleIcon);
const iconColor = ref<string>('text-neutral');
const actionText = ref<string>(props.activity.getKey());
const actionValue = ref<string>(`${props.activity.getOldvalue()} -> ${props.activity.getNewvalue()}`);

switch (props.activity.getKey()) {
    case 'UserProps.Wanted': {
        actionText.value = t('components.citizens.citizen_info_activity_feed_entry.set_citizen_as');

        if (props.activity.getNewvalue() === 'true') {
            icon.value = BellAlertIcon;
            iconColor.value = 'text-error-400'
            actionValue.value = props.activity.getNewvalue() === 'true' ? t('common.wanted') : `${t('common.not').toLowerCase()} ${t('common.wanted')}`;
        } else {
            icon.value = BellSnoozeIcon;
            iconColor.value = 'text-success-400'
            actionValue.value = props.activity.getNewvalue() === 'true' ? t('common.wanted') : `${t('common.not').toLowerCase()} ${t('common.wanted')}`;
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
                <h3 class="text-sm font-medium text-neutral">{{ activity.getSourceUser()?.getFirstname() }} {{
                    activity.getSourceUser()?.getLastname() }}</h3>
                <p class="text-sm text-gray-400">
                    {{ useLocaleTimeAgo(toDate(activity.getCreatedAt())!).value }}
                </p>
            </div>
            <p class="text-sm text-gray-300">{{ actionText }} <span class="font-bold">{{ actionValue }}</span></p>
        </div>
    </div>
</template>
