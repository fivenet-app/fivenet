export function calculatePage(totalCount?: number, pageSize?: number, currentPage?: number): number {
    if (totalCount === undefined || pageSize === undefined) {
        return 0;
    }

    const totalPages = Math.ceil(totalCount / pageSize);

    const pageC = currentPage ?? 1;
    if (pageC > totalPages) {
        return (totalPages - 1) * pageSize;
    } else if (pageC < 1) {
        return 0;
    }

    const o = pageSize * (pageC - 1);
    if (o < 0) {
        return 0;
    }

    return o;
}
