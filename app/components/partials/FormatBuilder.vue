<script setup lang="ts">
import { VueDraggable } from 'vue-draggable-plus';
import ReorderButtons from './ReorderButtons.vue';
import DraggableHandle from './DraggableHandle.vue';

const props = defineProps<{
    modelValue: string;
    extensions: Extension[];
    disabled?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

// Type for a placeholder/token extension
interface Extension {
    label: string;
    value: string;
}

// Types for our block structure
type Block = { id: string; type: 'text'; value: string } | { id: string; type: 'token'; value: string };

// Parse the format string into an array of blocks
function parseBlocks(str: string): Block[] {
    if (!str) return [{ id: randomUUID(), type: 'text', value: '' }];
    const regex = /%([a-zA-Z0-9_]+)%/g;
    let result: RegExpExecArray | null,
        lastIndex = 0;
    const blocks: Block[] = [];
    while ((result = regex.exec(str)) !== null) {
        // Add any text before the matched token
        if (result.index > lastIndex) {
            blocks.push({
                id: randomUUID(),
                type: 'text',
                value: str.slice(lastIndex, result.index),
            });
        }
        // Add the token block
        blocks.push({
            id: randomUUID(),
            type: 'token',
            value: result[1] ?? '',
        });
        lastIndex = regex.lastIndex;
    }
    // Add trailing text, if any
    if (lastIndex < str.length) {
        blocks.push({
            id: randomUUID(),
            type: 'text',
            value: str.slice(lastIndex),
        });
    }
    // Always at least one block
    if (blocks.length === 0) {
        blocks.push({ id: randomUUID(), type: 'text', value: '' });
    }
    return blocks;
}

// Convert blocks back into a format string
function blocksToString(blocks: Block[]): string {
    return blocks
        .map(
            (b) =>
                b.type === 'text'
                    ? b.value // leave as-is
                    : `%${b.value}%`, // tokens get percent-wrapped
        )
        .join('');
}

const blocks = ref<Block[]>(parseBlocks(props.modelValue));

watch(
    () => props.modelValue,
    (val) => {
        if (blocksToString(blocks.value) !== val) {
            blocks.value = parseBlocks(val);
        }
    },
);

function emitUpdate() {
    emit('update:modelValue', blocksToString(blocks.value));
}

watch(blocks, emitUpdate, { deep: true });

function handleTextInput(index: number, value: string) {
    const regex = /%([a-zA-Z0-9_]+)%/g;
    let result: RegExpExecArray | null,
        lastIndex = 0;
    const originalId = blocks.value[index]?.id;
    const newBlocks: Block[] = [];

    let firstText = true;

    while ((result = regex.exec(value)) !== null) {
        if (result.index > lastIndex) {
            newBlocks.push({
                id: firstText ? (originalId ?? randomUUID()) : randomUUID(),
                type: 'text',
                value: value.slice(lastIndex, result.index),
            });
            firstText = false;
        }
        newBlocks.push({
            id: randomUUID(),
            type: 'token',
            value: result[1] ?? '',
        });
        lastIndex = regex.lastIndex;
    }
    if (lastIndex < value.length) {
        newBlocks.push({
            id: firstText ? (originalId ?? randomUUID()) : randomUUID(),
            type: 'text',
            value: value.slice(lastIndex),
        });
    }
    // Replace the single text block at this index with all the new blocks
    blocks.value.splice(index, 1, ...newBlocks);

    ensureTrailingTextBlock();
}

function insertToken(tokenValue: string) {
    blocks.value.push({
        id: randomUUID(),
        type: 'token',
        value: tokenValue,
    });
}

function insertTextBlock() {
    blocks.value.push({
        id: randomUUID(),
        type: 'text',
        value: '',
    });
}

function removeBlock(idx: number) {
    blocks.value.splice(idx, 1);
    if (blocks.value.length === 0) {
        insertTextBlock();
    }
}

function findLabel(val: string): string | undefined {
    const ext = props.extensions.find((e) => e.value === val);
    return ext?.label ?? `%${val}%`;
}

function ensureTrailingTextBlock() {
    const last = blocks.value[blocks.value.length - 1];
    if (!last || last.type !== 'text') {
        blocks.value.push({
            id: randomUUID(),
            type: 'text',
            value: '',
        });
    }
}

const { moveUp, moveDown } = useListReorder(blocks);
</script>

<template>
    <div class="flex flex-col gap-2">
        <VueDraggable
            v-model="blocks"
            class="flex min-h-[48px] flex-wrap items-center gap-2 rounded-sm bg-neutral-100 p-2 dark:bg-neutral-800"
            :item-key="'id'"
            handle=".handle"
            :ghost-class="'opacity-50'"
            :disabled="disabled"
        >
            <template v-for="(element, index) in blocks" :key="element.id">
                <div class="flex items-center gap-1">
                    <UBadge class="flex items-center gap-1 font-mono" color="primary" variant="soft" size="md">
                        <UInput
                            v-if="element.type === 'text'"
                            class="w-28"
                            :model-value="element.value"
                            size="xs"
                            type="text"
                            :placeholder="$t('common.text')"
                            :disabled="disabled"
                            square
                            @update:model-value="handleTextInput(index, $event)"
                        />

                        <span v-else>{{ findLabel(element.value) || element.value }}</span>

                        <div class="inline-flex items-center gap-1">
                            <UButton
                                type="button"
                                icon="i-mdi-close-circle"
                                size="xs"
                                color="error"
                                variant="link"
                                tabindex="-1"
                                :disabled="disabled"
                                @click="removeBlock(index)"
                            />

                            <DraggableHandle size="xs" orientation="horizontal" :disabled="disabled" />

                            <ReorderButtons
                                size="xs"
                                :idx="index"
                                :move-up="moveUp"
                                :move-down="moveDown"
                                orientation="horizontal"
                            />
                        </div>
                    </UBadge>
                </div>
            </template>
        </VueDraggable>

        <div class="flex flex-wrap items-center gap-2">
            <UButton
                v-for="ext in extensions"
                :key="ext.value"
                class="font-mono"
                type="button"
                icon="i-mdi-plus-circle"
                size="xs"
                color="primary"
                variant="soft"
                :label="ext.label"
                :disabled="disabled"
                @click="insertToken(ext.value)"
            />

            <UButton
                type="button"
                size="xs"
                :disabled="disabled"
                color="primary"
                variant="soft"
                icon="i-mdi-plus"
                :label="$t('common.text')"
                @click="insertTextBlock"
            />
        </div>
    </div>
</template>
