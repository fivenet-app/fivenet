/* eslint-disable @typescript-eslint/no-explicit-any */
import { classRegistry, Rect, type ObjectEvents, type RectProps, type SerializedRectProps } from 'fabric';

type InputType = 'text' | 'number' | 'date' | 'datetime' | 'time' | 'checkbox' | 'select' | 'signature' | 'textarea';

type HtmlInputProps = {
    name?: string;
    inputType?: InputType;
    value?: string;
    placeholder?: string;
    fieldProps?: Record<string, any>;
} & Partial<RectProps>;

export class FabricHtmlInput<
    Props extends Partial<RectProps> = HtmlInputProps,
    SProps extends SerializedRectProps = SerializedRectProps,
    Events extends ObjectEvents = ObjectEvents,
> extends Rect<Props, SProps, Events> {
    static override type = 'html-input';

    name: string;
    inputType: string;
    value: string;
    placeholder: string;
    fieldProps: Record<string, any>;

    label: string;
    options: string[];
    variant: string;
    fontSize: number;
    fontFamily: string;
    textColor: string;

    constructor(opts: HtmlInputProps) {
        super({
            ...opts,
            width: opts.width ?? 150,
            height: opts.height ?? 30,
            fill: opts.fill ?? '#fff',
        } as any);

        this.name = opts.name || '';
        this.inputType = opts.inputType || 'text';
        this.value = opts.value || '';
        this.placeholder = opts.placeholder || '';
        this.fieldProps = opts.fieldProps || {};

        this.label = this.fieldProps.label || '';
        this.options = this.fieldProps.options || [];
        this.variant = this.fieldProps.variant || 'outline';
        this.fontSize = this.fieldProps.fontSize || 14;
        this.fontFamily = this.fieldProps.fontFamily || 'Arial';
        this.textColor = this.fieldProps.textColor || '#000';

        this.transparentCorners = true;
        this.padding = 4;
    }

    override render(ctx: CanvasRenderingContext2D) {
        super.render(ctx);

        if (this.label) {
            ctx.save();
            ctx.font = '12px sans-serif';
            ctx.fillStyle = '#333';
            ctx.textAlign = 'left';
            ctx.fillText(this.label, -this.width / 2 + 4, -this.height / 2 - 6);
            ctx.restore();
        }
    }

    override toObject(propertiesToInclude?: any): any {
        const base = super.toObject(propertiesToInclude);
        return {
            ...base,
            name: this.name,
            inputType: this.inputType || 'text',
            value: this.value,
            placeholder: this.placeholder,
            fieldProps: {
                label: this.label,
                options: this.options,
                variant: this.variant,
                fontSize: this.fontSize,
                fontFamily: this.fontFamily,
                textColor: this.textColor,
                ...this.fieldProps,
            },
        };
    }
}

classRegistry.setClass(FabricHtmlInput, 'html-input');
