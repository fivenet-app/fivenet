import type { PaginationResponse } from '~~/gen/ts/resources/common/database/database';

export function calculateOffset(page: number, pagination?: PaginationResponse): number {
    return pagination?.pageSize ? pagination?.pageSize * (page - 1) : 0;
}
