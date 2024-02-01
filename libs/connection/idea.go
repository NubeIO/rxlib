package connection

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
-----supervisor (supervisor will give APIs to each connection)
-------histories
---------connection-1 histories (edge-2)
---------connection-2 histories (edge-3)
---connection-1 (edge-2)
---connection-2 (edge-3) (will have a manager to manage its own histories)
----histManager
----hist
*/
