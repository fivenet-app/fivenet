# tracker

**State**: Design needs to be finalized.

## Components

* Client
    * How to access users grouped by job?
    * Use Nats and/ or broker channels for accessing the event data.
    * Use key value store to access data by user ID. -> What about the user's job?
* Worker
    * Load locations from database and generate an event object (like with the broker channel)
    * Store state in key value store for accessing data by user IDs. -> What about the user's job?

## What is the state the tracker needs to keeps?

* User's job
* User's on-duty state
* User's unit id -> Centrum `UserUnitMapping` (`settings.proto`)
