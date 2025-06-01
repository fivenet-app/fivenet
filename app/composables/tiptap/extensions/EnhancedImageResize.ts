// EnhancedImageResize.ts
//
// Combines:
//   • Tiptap's <img> node behaviour
//   • Resize handles + alignment toolbar (left / center / right)
//   • Extra `fileId` attribute for your upload / meta.files workflow
//
// Drop this file into `extensions/` and register it in the editor
// instead of the default Image extension.
//
// ---------------------------------------------------------------------

// Modified version of [GitHub bae-sh/tiptap-extension-resize-image](https://github.com/bae-sh/tiptap-extension-resize-image/blob/main/lib/imageResize.ts),
// which is licensed under the [MIT License](https://github.com/bae-sh/tiptap-extension-resize-image/blob/main/LICENSE).
import { EnhancedImage } from './EnhancedImage';

// From https://icon-sets.iconify.design/mdi/?icon-filter=format+align
const alignLeftIcon =
    '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M3 3h18v2H3zm0 4h12v2H3zm0 4h18v2H3zm0 4h12v2H3zm0 4h18v2H3z"/></svg>';
const alignCenterIcon =
    '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M3 3h18v2H3zm4 4h10v2H7zm-4 4h18v2H3zm4 4h10v2H7zm-4 4h18v2H3z"/></svg>';
const alignRightIcon =
    '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M3 3h18v2H3zm6 4h12v2H9zm-6 4h18v2H3zm6 4h12v2H9zm-6 4h18v2H3z"/></svg>';

export const EnhancedImageResize = EnhancedImage.extend({
    addAttributes() {
        return {
            ...this.parent?.(),
            style: {
                default: 'width: 100%; height: auto; cursor: pointer;',
                parseHTML: (element) => {
                    const width = element.getAttribute('width');
                    return width ? `width: ${width}px; height: auto; cursor: pointer;` : `${element.style.cssText}`;
                },
            },
        };
    },
    addNodeView() {
        return ({ node, editor, getPos }) => {
            const {
                view,
                options: { editable },
            } = editor;
            const { style } = node.attrs;
            const $wrapper = document.createElement('div');
            const $container = document.createElement('div');
            const $img = document.createElement('img');
            const iconStyle = 'width: 24px; height: 24px; cursor: pointer;';

            const dispatchNodeView = () => {
                if (typeof getPos === 'function') {
                    const newAttrs = {
                        ...node.attrs,
                        style: `${$img.style.cssText}`,
                    };
                    view.dispatch(view.state.tr.setNodeMarkup(getPos(), null, newAttrs));
                }
            };
            const paintPositionContoller = () => {
                const $postionController = document.createElement('div');

                const $leftController = document.createElement('div');
                const $centerController = document.createElement('div');
                const $rightController = document.createElement('div');

                const controllerMouseOver = (e: MouseEvent) => {
                    // @ts-expect-error target should always be valid for the mouse event
                    e.target.style.opacity = 0.3;
                };

                const controllerMouseOut = (e: MouseEvent) => {
                    // @ts-expect-error target should always be valid for the mouse event
                    e.target.style.opacity = 1;
                };

                $postionController.setAttribute(
                    'class',
                    'absolute top-0 left-1/2 h-7 z-[999] flex justify-between items-center px-1 py-0.5 gap-2 cursor-pointer bg-gray-200 dark:bg-gray-600',
                );
                $postionController.setAttribute(
                    'style',
                    'border-radius: 4px; border: 2px solid #6C6C6C; transform: translate(-50%, -50%);',
                );

                $leftController.innerHTML = alignLeftIcon;
                $leftController.setAttribute('style', iconStyle);
                $leftController.addEventListener('mouseover', controllerMouseOver);
                $leftController.addEventListener('mouseout', controllerMouseOut);

                $centerController.innerHTML = alignCenterIcon;
                $centerController.setAttribute('style', iconStyle);
                $centerController.addEventListener('mouseover', controllerMouseOver);
                $centerController.addEventListener('mouseout', controllerMouseOut);

                $rightController.innerHTML = alignRightIcon;
                $rightController.setAttribute('style', iconStyle);
                $rightController.addEventListener('mouseover', controllerMouseOver);
                $rightController.addEventListener('mouseout', controllerMouseOut);

                $leftController.addEventListener('click', () => {
                    $img.setAttribute('style', `${$img.style.cssText} margin: 0 auto 0 0;`);
                    dispatchNodeView();
                });
                $centerController.addEventListener('click', () => {
                    $img.setAttribute('style', `${$img.style.cssText} margin: 0 auto;`);
                    dispatchNodeView();
                });
                $rightController.addEventListener('click', () => {
                    $img.setAttribute('style', `${$img.style.cssText} margin: 0 0 0 auto;`);
                    dispatchNodeView();
                });

                $postionController.appendChild($leftController);
                $postionController.appendChild($centerController);
                $postionController.appendChild($rightController);

                $container.appendChild($postionController);
            };

            $wrapper.setAttribute('style', `display: flex;`);
            $wrapper.appendChild($container);

            $container.setAttribute('style', `${style}`);
            $container.appendChild($img);

            Object.entries(node.attrs).forEach(([key, value]) => {
                if (value === undefined || value === null) return;
                $img.setAttribute(key, value);
            });

            if (!editable) return { dom: $container };
            const isMobile = document.documentElement.clientWidth < 768;
            const dotPosition = isMobile ? '-8px' : '-4px';
            const dotsPosition = [
                `top: ${dotPosition}; left: ${dotPosition}; cursor: nwse-resize;`,
                `top: ${dotPosition}; right: ${dotPosition}; cursor: nesw-resize;`,
                `bottom: ${dotPosition}; left: ${dotPosition}; cursor: nesw-resize;`,
                `bottom: ${dotPosition}; right: ${dotPosition}; cursor: nwse-resize;`,
            ];

            let isResizing = false;
            let startX: number, startWidth: number;

            $container.addEventListener('click', (_) => {
                //remove remaining dots and position controller
                const isMobile = document.documentElement.clientWidth < 768;
                isMobile && (document.querySelector('.ProseMirror-focused') as HTMLElement)?.blur();

                if ($container.childElementCount > 3) {
                    for (let i = 0; i < 5; i++) {
                        $container.removeChild($container.lastChild as Node);
                    }
                }

                paintPositionContoller();

                $container.setAttribute('style', `position: relative; border: 1px dashed #6C6C6C; ${style} cursor: pointer;`);

                Array.from({ length: 4 }, (_, index) => {
                    const $dot = document.createElement('div');
                    $dot.setAttribute(
                        'style',
                        `position: absolute; width: ${isMobile ? 16 : 9}px; height: ${isMobile ? 16 : 9}px; border: 1.5px solid #6C6C6C; border-radius: 50%; ${dotsPosition[index]}`,
                    );

                    $dot.addEventListener('mousedown', (e) => {
                        e.preventDefault();
                        isResizing = true;
                        startX = e.clientX;
                        startWidth = $container.offsetWidth;

                        const onMouseMove = (e: MouseEvent) => {
                            if (!isResizing) return;
                            const deltaX = index % 2 === 0 ? -(e.clientX - startX) : e.clientX - startX;

                            const newWidth = startWidth + deltaX;

                            $container.style.width = newWidth + 'px';

                            $img.style.width = newWidth + 'px';
                        };

                        const onMouseUp = () => {
                            if (isResizing) {
                                isResizing = false;
                            }
                            dispatchNodeView();

                            document.removeEventListener('mousemove', onMouseMove);
                            document.removeEventListener('mouseup', onMouseUp);
                        };

                        document.addEventListener('mousemove', onMouseMove);
                        document.addEventListener('mouseup', onMouseUp);
                    });

                    $dot.addEventListener(
                        'touchstart',
                        (e) => {
                            e.cancelable && e.preventDefault();
                            isResizing = true;
                            startX = e.touches[0]!.clientX;
                            startWidth = $container.offsetWidth;

                            const onTouchMove = (e: TouchEvent) => {
                                if (!isResizing) return;
                                const deltaX =
                                    index % 2 === 0 ? -(e.touches[0]!.clientX - startX) : e.touches[0]!.clientX - startX;

                                const newWidth = startWidth + deltaX;

                                $container.style.width = newWidth + 'px';

                                $img.style.width = newWidth + 'px';
                            };

                            const onTouchEnd = () => {
                                if (isResizing) {
                                    isResizing = false;
                                }
                                dispatchNodeView();

                                document.removeEventListener('touchmove', onTouchMove);
                                document.removeEventListener('touchend', onTouchEnd);
                            };

                            document.addEventListener('touchmove', onTouchMove);
                            document.addEventListener('touchend', onTouchEnd);
                        },
                        { passive: false },
                    );
                    $container.appendChild($dot);
                });
            });

            document.addEventListener('click', (e: MouseEvent) => {
                const $target = e.target as HTMLElement;
                const isClickInside = $container.contains($target) || $target.style.cssText === iconStyle;

                if (!isClickInside) {
                    const containerStyle = $container.getAttribute('style');
                    const newStyle = containerStyle?.replace('border: 1px dashed #6C6C6C;', '');
                    $container.setAttribute('style', newStyle as string);

                    if ($container.childElementCount > 3) {
                        for (let i = 0; i < 5; i++) {
                            $container.removeChild($container.lastChild as Node);
                        }
                    }
                }
            });

            return {
                dom: $wrapper,
            };
        };
    },
});
