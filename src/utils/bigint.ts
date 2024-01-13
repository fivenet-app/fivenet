export function bigIntCeil(n: bigint, d: bigint): bigint {
    return n / d + (n % d ? 1n : 0n);
}
