package rubix

/*
Outgoing message
an object port will need to know if it has a remote rubix so it can publish its message twice on the eventbus,
once for any local connections and the 2nd for any remote connections, The reason is for the 2nd message is we could have 1000 object,
but we only do remote mapping for 1 object so we dont want to spam the local eventbus for nothing

Incoming message
we need a rest/mqtt API to send a message on the eventbus, there will be a service to handle an incoming message to publish,
the message on the eventbus that came via rest/mqtt


Network Manager
Will be a map[string]Networks

Network
- IP
- PORT
- Type MQTT/REST
- ConnectionList {uuid, connectionUUID}


*/

/*
- Lead Supervisor (in the cloud)
- Manager is an edge device that has ros-networks itself
- Network a network is a network to another ros instance

---Lead Supervisor (cloud)
-----GetAllMangerNetworks()
-----GetManager(uuid string)
----Manager Network (Edge-1)
------GetAllSch()
----Network-1 (Edge-2)
----Network-2 (Edge-3)
------hist
------sch


Edge (edge-1)
--------------
edge-1
--connectionManager
-----supervisor (supervisor will give APIs to each rubix)
-------histories
---------rubix-1 histories (edge-2)
---------rubix-2 histories (edge-3)
---rubix-1 (edge-2)
---rubix-2 (edge-3) (will have a manager to manage its own histories)
----histManager
----hist
*/
