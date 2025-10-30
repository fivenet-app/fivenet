/* eslint-disable @typescript-eslint/no-explicit-any */
// Extend Fabric.js to include a standalone class for creating curved text
import { classRegistry, FabricText, Group } from 'fabric';

export type FabricCurvedTextOptions = {
    fontSize?: number;
    fontFamily?: string;
    fill?: string;
    clockwise?: boolean;
    reverse?: boolean;
    arcAngleDeg?: number; // total arc span in degrees (e.g., 180)
};

export class FabricCurvedText extends Group {
    static override type = 'curved-text';

    text: string;
    radius: number;
    options: FabricCurvedTextOptions;

    constructor(text: string, radius: number = 100, opts: FabricCurvedTextOptions = {}) {
        const textObjects = FabricCurvedText.createTextObjects(text, radius, opts);

        super(textObjects, {
            originX: 'center',
            originY: 'center',
            selectable: true,
            evented: true,
            fill: opts.fill || '#c00',
        });

        this.text = text;
        this.radius = radius;
        this.options = opts;

        // Listen for scaling and rotating events
        this.on('scaling', this.handleScaling.bind(this));
        this.on('rotating', this.handleRotation.bind(this));
    }

    handleScaling(): void {
        const scaleX = this.scaleX || 1;
        const scaleY = this.scaleY || 1;

        // Reset scale to 1 and adjust radius based on scaling
        this.scaleX = 1;
        this.scaleY = 1;

        // Clamp the radius to prevent it from growing indefinitely
        const newRadius = Math.max(10, Math.min(this.radius * Math.max(scaleX, scaleY), 1000));

        // Update the layout with the new radius
        this.update(this.text, newRadius, this.options);
    }

    handleRotation(): void {
        // Ensure bounding box and coordinates are updated after rotation
        this.setCoords();

        // Normalize transformations to prevent cumulative growth
        this.scaleX = 1;
        this.scaleY = 1;

        // Update the bounding box dimensions
        this.updateBoundingBox();
    }

    updateBoundingBox(): void {
        // Dynamically recalculate the bounding box dimensions, excluding margins
        const boundingRect = this.getBoundingRect();

        // Dynamically clamp the bounding box dimensions based on the canvas size
        const maxWidth = 2000; // Example maximum width, adjust as needed
        const maxHeight = 2000; // Example maximum height, adjust as needed

        this.width = Math.max(10, Math.min(boundingRect.width, maxWidth));
        this.height = Math.max(10, Math.min(boundingRect.height, maxHeight));

        this.setCoords();
    }

    static createTextObjects(text: string, radius: number, opts: FabricCurvedTextOptions): FabricText[] {
        const {
            fontSize = 16,
            fontFamily = 'Arial',
            fill = '#c00',
            clockwise = true,
            reverse = false,
            arcAngleDeg = 180,
        } = opts;

        const chars = text.split('');
        const arcAngle = (arcAngleDeg * Math.PI) / 180;
        const angleStep = arcAngle / Math.max(chars.length - 1, 1);
        const startAngle = clockwise ? -arcAngle / 2 : Math.PI + arcAngle / 2;

        return chars.map((char, i) => {
            const angle = startAngle + i * angleStep * (clockwise ? 1 : -1);
            const x = radius * Math.cos(angle);
            const y = radius * Math.sin(angle);
            const rotation = (angle * 180) / Math.PI + (clockwise ? 90 : -90);

            return new FabricText(char, {
                left: x,
                top: y,
                fontSize,
                fontFamily,
                fill,
                angle: reverse ? rotation + 180 : rotation,
                originX: 'center',
                originY: 'center',
                selectable: false,
                evented: false,
                centeredRotation: true,
                centeredScaling: true,
            });
        });
    }

    update(text?: string, radius?: number, opts?: Partial<FabricCurvedTextOptions>): void {
        if (text !== undefined) this.text = text;
        if (radius !== undefined) this.radius = radius;
        if (opts !== undefined) this.options = { ...this.options, ...opts };

        // Store the current position and angle
        const currentLeft = this.left;
        const currentTop = this.top;
        const currentAngle = this.angle;

        const newTextObjects = FabricCurvedText.createTextObjects(this.text, this.radius, this.options);

        // Adjust positions of text objects relative to the group's center
        newTextObjects.forEach((obj) => {
            obj.left -= this.radius; // Center the text objects around the group's center
            obj.top -= this.radius;
            obj.originX = 'center';
            obj.originY = 'center';
        });

        this.removeAll();
        this.add(...newTextObjects);

        // Recalculate the group's dimensions
        this.updateBoundingBox();

        // Reapply the position and angle
        this.left = currentLeft;
        this.top = currentTop;
        this.angle = currentAngle;
        this.originX = 'center';
        this.originY = 'center';
    }

    override toObject(propertiesToInclude?: any[]): any {
        return {
            ...super.toObject(propertiesToInclude),
            text: this.text,
            radius: this.radius,
            options: this.options,
        };
    }
}

classRegistry.setClass(FabricCurvedText, 'curved-text');
