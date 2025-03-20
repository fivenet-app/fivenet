// @generated by protobuf-ts 2.9.6 with parameter optimize_speed,long_type_number,force_server_none
// @generated from protobuf file "resources/livemap/tracker.proto" (package "resources.livemap", syntax proto3)
// @ts-nocheck
import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import { WireType } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import { UnknownFieldHandler } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { reflectionMergePartial } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
import { UserMarker } from "./livemap";
/**
 * @generated from protobuf message resources.livemap.UsersUpdateEvent
 */
export interface UsersUpdateEvent {
    /**
     * @generated from protobuf field: repeated resources.livemap.UserMarker added = 1;
     */
    added: UserMarker[];
    /**
     * @generated from protobuf field: repeated resources.livemap.UserMarker removed = 2;
     */
    removed: UserMarker[];
}
// @generated message type with reflection information, may provide speed optimized methods
class UsersUpdateEvent$Type extends MessageType<UsersUpdateEvent> {
    constructor() {
        super("resources.livemap.UsersUpdateEvent", [
            { no: 1, name: "added", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => UserMarker },
            { no: 2, name: "removed", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => UserMarker }
        ]);
    }
    create(value?: PartialMessage<UsersUpdateEvent>): UsersUpdateEvent {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.added = [];
        message.removed = [];
        if (value !== undefined)
            reflectionMergePartial<UsersUpdateEvent>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UsersUpdateEvent): UsersUpdateEvent {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* repeated resources.livemap.UserMarker added */ 1:
                    message.added.push(UserMarker.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* repeated resources.livemap.UserMarker removed */ 2:
                    message.removed.push(UserMarker.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: UsersUpdateEvent, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* repeated resources.livemap.UserMarker added = 1; */
        for (let i = 0; i < message.added.length; i++)
            UserMarker.internalBinaryWrite(message.added[i], writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* repeated resources.livemap.UserMarker removed = 2; */
        for (let i = 0; i < message.removed.length; i++)
            UserMarker.internalBinaryWrite(message.removed[i], writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message resources.livemap.UsersUpdateEvent
 */
export const UsersUpdateEvent = new UsersUpdateEvent$Type();
