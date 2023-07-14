<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiClipboardPlus } from '@mdi/js';
import { useClipboard } from '@vueuse/core';
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';
import Time from '~/components/partials/elements/Time.vue';
import { useNotificationsStore } from '~/store/notifications';
import { AuditEntry, EVENT_TYPE } from '~~/gen/ts/resources/rector/audit';

const clipboard = useClipboard();

const { d } = useI18n();

const props = defineProps<{
    log: AuditEntry;
}>();

const notifications = useNotificationsStore();

async function addToClipboard(): Promise<void> {
    const user = props.log.user;
    let text = `**Audit Log Entry ${props.log.id} - ${d(toDate(props.log.createdAt)!, 'short')}**

`;
    if (user) {
        text += `User: ${user?.firstname}, ${user?.lastname} (${user?.userId}; ${user?.identifier})
`;
    }
    text += `Action: \`${props.log.service}/${props.log.method}\`
Event: \`${EVENT_TYPE[props.log.state]}\`
`;
    if (props.log.data) {
        text += `Data:
\`\`\`json
${JSON.stringify(JSON.parse(props.log.data!), null, 2)}
\`\`\`
`;
    } else {
        text += `Data: N/A
`;
    }

    notifications.dispatchNotification({
        title: { key: 'notifications.rector.audit_log.title', parameters: [] },
        content: { key: 'notifications.rector.audit_log.content', parameters: [] },
        type: 'info',
    });

    return clipboard.copy(text);
}
</script>

<template>
    <tr>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.id }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            <Time :value="log.createdAt" type="long" />
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            <NuxtLink :to="{ name: 'citizens-id', params: { id: log.userId.toString() } }">
                {{ log.user ? log.user?.firstname + ' ' + log.user?.lastname : 'N/A' }}
            </NuxtLink>
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.service }}/{{ log.method }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ EVENT_TYPE[log.state] }}
        </td>
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0 max-w-3xl">
            <span v-if="!log.data">N/A</span>
            <span v-else>
                <VueJsonPretty
                    :data="JSON.parse(props.log.data!) as any"
                    :showIcon="true"
                    :showLength="true"
                    :virtual="true"
                    :height="160"
                />
            </span>
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <button
                class="flex-initial text-primary-500 hover:text-primary-400"
                :title="$t('components.clipboard.clipboard_button.add')"
            >
                <SvgIcon class="w-6 h-auto ml-auto mr-2.5" @click="addToClipboard" type="mdi" :path="mdiClipboardPlus" />
            </button>
        </td>
    </tr>
</template>
