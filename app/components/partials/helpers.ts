export type imageSizes = 'sm' | 'md' | 'lg' | 'xl' | '6xl' | 'huge';

export function imageSize(size: imageSizes): string {
    switch (size) {
        case 'huge':
            return 'text-8xl h-44 w-44';
        case '6xl':
            return 'text-6xl h-32 w-32';
        case 'xl':
            return 'text-3xl h-20 w-20';
        case 'lg':
            return 'text-lg h-12 w-12';
        case 'sm':
            return 'text-sm h-8 w-8';
        case 'md':
        default:
            return 'text-base h-10 w-10';
    }
}
