export type imageSizes = 'sm' | 'md' | 'lg' | 'xl';

export function imageSize(size: imageSizes): string {
    switch (size) {
        case 'xl':
            return 'text-lg h-20 w-20';
        case 'lg':
            return 'text-lg h-12 w-12';
        case 'sm':
            return 'text-sm h-8 w-8';
        case 'md':
        default:
            return 'text-base h-10 w-10';
    }
}
