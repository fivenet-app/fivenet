<script lang="ts" setup>
import { AuditEntry } from '~~/gen/ts/resources/rector/audit';
import { EVENT_TYPE_Util } from '~~/gen/ts/resources/rector/audit.pb_enums';
import { ClipboardDocumentIcon } from '@heroicons/vue/24/solid';

const { d } = useI18n();

const props = defineProps({
    log: {
        type: AuditEntry,
        required: true,
    }
});

async function addToClipboard(): Promise<void> {
    const user = props.log.getUser();
    const text = `**Audit Log Entry ${props.log.id} - ${d(props.log.createdAt?.timestamp?.toDate()!, 'short')}**
User: ${user?.firstname}, ${user?.lastname} (${user?.userId}; ${user?.getIdentifier()})
Action: ${props.log.getMethod()}/${props.log.getService()}
Event: ${EVENT_TYPE_Util.toEnumKey(props.log.getState())}
Data:
\`\`\`
${props.log.getData()}
\`\`\`
`;

    return navigator.clipboard.writeText(text);
}
</script>

<template>
    <tr>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.id }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ $d(log.createdAt?.timestamp?.toDate()!, 'short') }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.hasUser() ? (log.getUser()?.firstname + ' ' + log.getUser()?.lastname) : 'N/A' }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.getService() }}: {{ log.getMethod() }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ EVENT_TYPE_Util.toEnumKey(log.getState()) }}
        </td>
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.getData() ? log.getData() : 'N/A' }}
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <button class="flex-initial text-primary-500 hover:text-primary-400"
                :title="$t('components.clipboard.clipboard_button.add')">
                <ClipboardDocumentIcon class="w-6 h-auto ml-auto mr-2.5" @click="addToClipboard" />
            </button>
        </td>
    </tr>
</template>
