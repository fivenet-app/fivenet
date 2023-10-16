import Quill from 'quill';

const BlockEmbed = Quill.import('blots/block/embed');

export const DEFAULT = {
    cssText: 'border: none;border-bottom: 1px inset;',
    className: void 0,
    icon: '<svg class="icon" style="vertical-align: middle;fill: currentColor;overflow: hidden;" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg"><path class="ql-fill" d="M64 464h896v96H64v-96zM224 96v160h576V96h96v256H128V96h96z m576 832v-160H224v160H128v-256h768v256h-96z"></path></svg>',
};

export type Options = typeof DEFAULT & {
    text?: {
        orientation?: 'left' | 'right' | 'center';
        children: string;
        childrenStyle: string;
    };
};

class Divider extends BlockEmbed {
    static create(options: Options) {
        const parentNode = super.create();
        parentNode.setAttribute('contenteditable', false);
        if (options.text) {
            const width: (string | number)[] = ['100%', '100%'];
            if (options.text.orientation === 'left') {
                width[0] = 0;
            }
            if (options.text.orientation === 'right') {
                width[1] = 0;
            }
            const commonStyle = 'position: relative;top: 50%;transform: translateY(-50%);';
            const lineLeft = `<span contenteditable="false" style="width: ${width[0]};${commonStyle}${
                options.cssText || ''
            }" class="${options.className || ''}" ></span>`;
            const lineRight = `<span contenteditable="false" style="width: ${width[1]};${commonStyle}${
                options.cssText || ''
            }" class="${options.className || ''}" ></span>`;
            parentNode.style.display = 'flex';
            parentNode.style.textAlign = 'center';
            parentNode.innerHTML = `${lineLeft}<span style="display: inline-block;${
                options.text.childrenStyle
            }" contenteditable="true">${options.text.children || ''}</span>${lineRight}`;
        } else {
            parentNode.innerHTML = `<hr style="${options.cssText}" class="${options.className || ''}" >`;
        }
        return parentNode;
    }
}

Divider.blotName = 'divider';
Divider.tagName = 'div';
Divider.className = 'quill-hr';

export default Divider;
