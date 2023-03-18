import slugify from 'slugify';

slugify.extend({ '.': '-' });

export default function slug(input: string): string {
    return slugify(input, {
        lower: true,
    });
}
