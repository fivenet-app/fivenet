export interface ProtobufMessage<T> {
    create(): T;
    fromBinary(bytes: Uint8Array): T;
    toBinary(m: T): Uint8Array;
}
