import type { Transaction } from '@tiptap/pm/state';
import type { Editor } from '@tiptap/vue-3';

type TextAlign = 'left' | 'center' | 'right' | 'justify';

export function useTiptapToolbar(editor: () => Editor | null | undefined) {
    const ui = reactive({
        // Active states
        bold: false,
        italic: false,
        underline: false,
        strike: false,
        code: false,
        superscript: false,
        subscript: false,
        highlight: false,
        link: false,
        paragraph: false,
        bulletList: false,
        orderedList: false,
        taskList: false,
        checkboxStandalone: false,
        codeBlock: false,
        blockquote: false,

        // Attributes / derived
        textAlign: 'left' as TextAlign,
        headingLevel: 0 as 0 | 1 | 2 | 3 | 4 | 5 | 6,
        fontColor: null as string | null,
        highlightColor: null as string | null,

        // Abilities (for :disabled)
        canBold: true,
        canItalic: true,
        canUnderline: true,
        canStrike: true,
        canCode: true,
        canUndo: false,
        canRedo: false,
    });

    // runs at most every animation frame + debounced for safety
    const refreshNow = () => {
        const ed = editor();
        if (!ed) return;
        // ACTIVE STATES (cheap single reads)
        ui.bold = ed.isActive('bold');
        ui.italic = ed.isActive('italic');
        ui.underline = ed.isActive('underline');
        ui.strike = ed.isActive('strike');
        ui.code = ed.isActive('code');
        ui.superscript = ed.isActive('superscript');
        ui.subscript = ed.isActive('subscript');
        ui.highlight = ed.isActive('highlight');
        ui.link = ed.isActive('link');
        ui.paragraph = ed.isActive('paragraph');
        ui.bulletList = ed.isActive('bulletList');
        ui.orderedList = ed.isActive('orderedList');
        ui.taskList = ed.isActive('taskList');
        ui.checkboxStandalone = ed.isActive('checkboxStandalone');
        ui.codeBlock = ed.isActive('codeBlock');
        ui.blockquote = ed.isActive('blockquote');

        // textAlign: check attrs once, don’t pass object literals in template
        if (ed.isActive({ textAlign: 'center' })) ui.textAlign = 'center';
        else if (ed.isActive({ textAlign: 'right' })) ui.textAlign = 'right';
        else if (ed.isActive({ textAlign: 'justify' })) ui.textAlign = 'justify';
        else ui.textAlign = 'left';

        // Heading level (0 means "no heading")
        let level: 0 | 1 | 2 | 3 | 4 | 5 | 6 = 0;
        for (let l = 1 as 1 | 2 | 3 | 4 | 5 | 6; l <= 6; l++) {
            if (ed.isActive('heading', { level: l })) {
                level = l;
                break;
            }
        }
        ui.headingLevel = level;

        // === Current colors ===
        // Color extension sets a textStyle mark
        const textStyleAttrs = ed.getAttributes('textStyle');
        ui.fontColor = (textStyleAttrs?.color as string | undefined) ?? null;

        // Highlight color lives on the highlight mark as `color` (when configured)
        const highlightAttrs = ed.getAttributes('highlight');
        ui.highlightColor = (highlightAttrs?.color as string | undefined) ?? null;

        // ABILITIES (compute once; avoid chain() in render)
        ui.canBold = !!ed.can().toggleBold?.();
        ui.canItalic = !!ed.can().toggleItalic?.();
        ui.canUnderline = !!ed.can().toggleUnderline?.();
        ui.canStrike = !!ed.can().toggleStrike?.();
        ui.canCode = !!ed.can().toggleCode?.();
        ui.canUndo = !!ed.can().undo?.();
        ui.canRedo = !!ed.can().redo?.();
    };

    // Debounce to avoid storms; 0ms still coalesces microtasks, or use ~16ms.
    const refresh = useDebounceFn(refreshNow, 0);

    const onSelectionUpdate = () => refresh();
    const onTransaction = ({ transaction }: { transaction: Transaction }) => {
        if (transaction?.selectionSet) refresh();
    };
    const onFocusBlur = () => refresh(); // keep UI in sync on focus changes

    // Wire editor events (call once when editor is ready)
    function attach() {
        const ed = editor();
        if (!ed) return;
        ed.on('selectionUpdate', onSelectionUpdate);
        ed.on('transaction', onTransaction);
        ed.on('focus', onFocusBlur);
        ed.on('blur', onFocusBlur);
        // Initial paint
        refresh();
    }

    function detach() {
        const ed = editor();
        if (!ed) return;
        ed.off('selectionUpdate', onSelectionUpdate);
        ed.off('transaction', onTransaction);
        ed.off('focus', onFocusBlur);
        ed.off('blur', onFocusBlur);
    }

    return { ui, attach, detach, refresh };
}
