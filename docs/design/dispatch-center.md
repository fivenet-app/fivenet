---
title: "Dispatch Center"
---

## Dispatch Center

* Two views (only for their own job; no job list attribute):
    * One for disponents.
    * Unit view on the livemap as a sidebar with quick buttons (see below "Livemap" point).
* Disponents take "control" of the dispatch center and can "sign on"/ "sign off" from it at any time.
    * If no more disponents are signed on, the "fallback mode" is activated.
* Different Modes:
    * `MANUAL`: Disponents **and** units can assign dispatches themselves.
    * `CENTRAL_COMMAND`: Only the disponents can assign dispatches.
    * `AUTO_ROUND_ROBIN`: Automatic assignment of dispatches to units that are either available, or busy (if none other available).
* Settings of dispatch center:
    * `Enabled`: If dispatch center is enabled or not.
    * `Mode`: Which mode it is run in by default.
    * `Fallback Mode`: In case the dispatch center is empty, which mode it will fallback into.
* Status:
    * `INACTIVE` - No one in the dispatch center (and not run by a bot).
    * `ACTIVE` - Dispatch center has "humans" signed on to it.
    * `AUTO` - An automatic system is running the dispatch center, requires the "mode" to be set to an "automatic" mode.
* Actions:
    * `GetSettings` - Retrieve the dispatch center settings.
    * `UpdateSettings` - Faction leaders set the dispatch center settings via this method.
    * `TakeControl` - Disponents "sign on"/"sign off" from their dispatch center duties using this.
    * `Stream` - Data stream for the changes happening.

## Units

* Created by the faction leaders via the control center.
* Users can be assigned by disponents and self-assign into one unit of their own job only.
* Users can set their own units status.
    * The "informal" status are ignored in the "current status" when getting the status (e.g., `USER_ADDED`, etc.)
* Disponents can update unit status as they please.
* Unit Status:
    * `UNKNOWN` - Dummy status should something go (very) wrong.
    * `USER_ADDED` - User added to unit.
    * `USER_REMOVED` - User removed from unit.
    * `UNAVAILABLE` - Unit is unavailable for dispatches.
    * `AVAILABLE` - Unit is available for dispatches.
    * `ON_BREAK` - Unit is on break (will still get pings about broadcast dispatches).
    * `BUSY` - Unit is busy (normally means en route to/on scene dispatch, etc.)
* Actions:
    * `ListUnits` - List all units that the person has access to.
    * `CreateOrUpdateUnit` - Create unit or update unit info (e.g., name, initials, description).
    * `DeleteUnit` - Delete unit.
    * `UpdateUnitStatus` - Update Unit status. Own unit can be updated by users, disponents can update any unit of their job.
    * `AssignUnit` - Disponent function to assign an user to an unit.
    * `JoinUnit` - For "normal" users to join and leave an unit via the Centrum Sidebar.
    * `ListUnitActivity` - List own or others units activity.

## Dispatches

* Consist of a message, description, and a position (x and y coordinates).
    * If specified or added/updated later on, attributes (list of strings) can be added.
* Status of the dispatch is "shared" by multiple units, e.g., unit one sets status `EN_ROUTE` and unit two sets `ON_SCENE`, both are in the dispatch status log.
    * The "informal" status are ignored in the "current status" when getting the status (e.g., `UNIT_ASSIGNED`, etc.)
* Can be created manually via the dispatch center, livemap integration or "the phone" (for now the existing GKSPhone dispatch system is used).
* Can have attributes which are a list of strings attached to them (e.g., `dangerous`, `gun shots`).
* Disponents can update dispatches, including their status, as they please.
* Dispatches assigned to an unit can either be "forced" or "requested".
    * "Requested" means that the assigned units get a popup with a timeout of 20 seconds that requires each unit to give their "approval or denial" to taking the dispatch.
        * If a dispatch is not accepted, it will be marked as unassigned if no other unit was assigend to it. "Last one" to not accept, will cause the dispatch status to go to `UNASSIGNED`
* If an user leaves an unit and it will be empty, or an unit marks itself as `UNAVAILABLE`, the dispatches that are assigned to the unit that aren't `COMPLETED` or `CANCELLED`, will be set as `UNASSIGNED` (when no other unit is assigned then).
* Dispatch Status:
    * `NEW` - Freshly created dispatch.
    * `UNASSIGNED` - Waiting for "triage" by the disponents or units (depending on the mode).
    * `UPDATED` - Dispatch has been updated.
    * `UNIT_ASSIGNED` - Unit assigned to dispatch.
    * `UNIT_UNASSIGNED` - Unit unassigned from dispatch.
    * `EN_ROUTE` - Unit en route to dispatch.
    * `ON_SCENE` - Unit arrived at dispatch location/on scene.
    * `NEED_ASSISTANCE` - Assistance needed at dispatch.
    * `COMPLETED` - Dispatch completed.
    * `CANCELLED` - Dispatch cancelled.
* Actions:
    * `ListDispatches` - List all dispatches that the person has access to.
    * `CreateDispatch` - Create a dispatch (mainly used via the livemap page or dispatch center livemap).
    * `UpdateDispatch` - Update dispatch details.
    * `UpdateDispatchStatus` - Update dispatch status.
    * `ListDispatchActivity` - List own or other dispatches activity.
    * `AssignDispatch` - Assign dispatch to units.
    * `TakeDispatch` - Units take/ "accept" dispatches either when they have been assigned one and accept it, or when they self assign a dispatch.
        * Response status:
            * `TIMEOUT`
            * `ACCEPTED`
            * `DECLINED`

## Livemap

* Quick action buttons for the different statuses.
    * E.g., `COMPLETED`, `NEED_ASSISTANCE`, dispatch require a reason to be given.
* Show additional markers from the gameserver:
    * Restricted zones ("Sperrzonen") as a [Leaflet Circle Markers (in the example the red circle)](https://leafletjs.com/examples/quick-start/#markers-circles-and-polygons).

## Panicbuttons

* Are "converted" to dispatches that are sent to every unit (broadcast) of the job the panicbutton presser is part of.

## Limitations

* All dispatches that are not completed or archived are loaded/ sent to the users.
