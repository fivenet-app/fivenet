---
title: "Dispatch Center"
---

* Dispatch Center
    * Two views: One for disponents and the unit view on the livemap.
    * Disponents take "control" of the dispatch center and can "sign on"/ "sign off" from it at any time.
        * If no more disponents are signed on, the "fallback mode" is activated.
    * Job-based permissions on the main `Stream` permission (Attribute: job List)
    * Different Modes:
        * `MANUAL`: A human must do things.
        * `CENTRAL_COMMAND`: Only the human in the Leitstelle can assign dispatches.
        * `AUTO_ROUND_ROBIN`: Automatic assignment of dispatches to units that are either available, or busy (if none other available).
* Units
    * Created by the faction leaders via the control center.
    * Users can be assigned by disponents and self-assign into one unit of their own job only.
    * Users can set their own units status.
        * The "informal" status are ignored when getting the status (e.g., `USER_ADDED`, etc.)
    * Disponents
* Dispatches
    * Consist of a message, description, and a position (x and y coordinates).
    * Status of the dispatch is "shared" by multiple units, e.g., unit one sets status `EN_ROUTE` and unit two sets `AT_SCENE`, both are in the dispatch status log.
        * The "informal" status are ignored when getting the status (e.g., `USER_ADDED`, etc.)
    * Can be created manually via the dispatch center, livemap integration or "the phone" (for now the existing GKSPhone dispatch system is used).
    * Can have attributes which are a list of strings attached to them (e.g., `dangerous`, `gun shots`).
* Livemap
    * Quick action buttons for the different statuses.
        * E.g., `COMPLETED`, `NEED_ASSISTANCE`, dispatch require a reason to be given.
    * Show additional markers from the gameserver:
        * Restricted zones ("Sperrzonen")
* Panicbuttons
    * Are just dispatches that are sent to every unit (broadcast) of the job the panicbutton presser is part of.
