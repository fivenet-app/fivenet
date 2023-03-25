<script setup lang="ts">
import { DocumentComment } from '@arpanet/gen/resources/documents/documents_pb';

const props = defineProps<{
    comments: DocumentComment[]
}>();

console.log(props.comments);

// TODO for adding/ editing a comment, use https://tailwindui.com/components/application-ui/forms/textareas#component-784309f82e9913989c2196a2d47eff4a
</script>

<template>
    <div class="flow-root px-4 rounded-lg bg-base-800 text-neutral">
        <ul role="list" class="divide-y divide-gray-200">
            <li v-for="comment in $props.comments" :key="comment.getId()" class="py-4">
                <div class="flex space-x-3">
                    <div class="flex-1 space-y-1">
                        <div class="flex items-center justify-between">
                            <router-link :to="{ name: 'Citizens: Info', params: { id: comment.getCreatorId() } }" class="text-sm font-medium text-primary-400 hover:text-primary-300">{{ comment.getCreator()?.getFirstname() }} {{ comment.getCreator()?.getLastname() }}</router-link>
                        </div>
                        <p class="text-sm">{{ comment.getComment() }}</p>
                    </div>
                </div>
            </li>
        </ul>
    </div>
</template>
