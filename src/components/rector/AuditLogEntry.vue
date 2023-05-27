<script lang="ts" setup>
import { AuditEntry } from '~~/gen/ts/resources/rector/audit';
import { ClipboardDocumentIcon } from '@heroicons/vue/24/solid';

const { d } = useI18n();

const props = defineProps<{
    log: AuditEntry,
}>();

async function addToClipboard(): Promise<void> {
    const user = props.log.user;
    const text = `**Audit Log Entry ${props.log.id} - ${d(toDate(props.log.createdAt)!, 'short')}**
User: ${user?.firstname}, ${user?.lastname} (${user?.userId}; ${user?.identifier})
Action: ${props.log.method}/${props.log.service}
Event: ${props.log.state}
Data:
\`\`\`
${props.log.data}
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
            {{ $d(toDate(log.createdAt)!, 'short') }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.user ? (log.user?.firstname + ' ' + log.user?.lastname) : 'N/A' }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.service }}: {{ log.method }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.state }}
        </td>
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ log.data ? log.data : 'N/A' }}
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <button class="flex-initial text-primary-500 hover:text-primary-400"
                :title="$t('components.clipboard.clipboard_button.add')">
                <ClipboardDocumentIcon class="w-6 h-auto ml-auto mr-2.5" @click="addToClipboard" />
            </button>
        </td>
    </tr>
</template>
