import Quill from 'quill';
import divider, { DEFAULT, Options } from './divider';

const Module = Quill.import('core/module');

class DividerToolbar extends Module {
    private quill: Quill;
    private readonly options: Options;

    constructor(quill: Quill, options: Options) {
        super(quill, options);
        this.options = Object.assign({}, DEFAULT, options);
        this.quill = quill;
        this.toolbar = quill.getModule('toolbar');
        this.toolbar.addHandler('divider', this.dividerHandler.bind(this));
        setTimeout(() => {
            const divider = document.querySelector('.ql-divider');
            if (divider) {
                divider.innerHTML = this.options.icon;
            }
        }, 75);
    }

    dividerHandler() {
        const getSelection = this.quill.getSelection() || { index: 0, length: 0 };
        const content = this.quill.getContents(getSelection.index, getSelection.length)?.ops || [];
        let selection = getSelection.index || this.quill.getLength();
        const [leaf] = this.quill.getLeaf(selection - 1);
        if (leaf instanceof divider) {
            this.quill.insertText(selection, '\n', Quill.sources.USER);
            selection++;
        }
        if (this.options.text) {
            this.quill.deleteText(getSelection.index, getSelection.length);
            this.options.text.children = typeof content[0]?.insert === 'string' ? content[0].insert : '';
        }
        this.quill.insertEmbed(selection, 'divider', this.options, Quill.sources.USER);
        if (getSelection.index === 0) {
            selection++;
            this.quill.insertText(selection, '\n', Quill.sources.USER);
        }
        this.quill.setSelection(selection + 2, 0);
    }
}

export default DividerToolbar;
