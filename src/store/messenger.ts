import Dexie, { type Table } from 'dexie';
import type { Message } from '~~/gen/ts/resources/messenger/message';
import type { Thread } from '~~/gen/ts/resources/messenger/thread';

export class MessengerDexie extends Dexie {
    threads!: Table<Thread>;
    messages!: Table<Message>;

    constructor() {
        super('messenger');
        this.version(1).stores({
            threads: 'id',
            messages: 'id, threadId',
        });
    }
}

export const db = new MessengerDexie();
