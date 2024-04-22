# Milestones

- [x] Basic P2P connection
- [ ] Edge cases for a P2P connection
- [ ] Basic Data synchronisation

# Basic Networking Protocol

## Initilizing a network pool

When two clients want to start a p2p connection, for initialization, either one of them would start the connection as a server. The client would then connect to this and send over a payload asking for the nodes client list and data. You would imagine it as two individuals exchanging each others contacts. Along with this the data would get shared over as well. This forms a p2p connection with 2 nodes.

![Formation of a network pool](https://i.imgur.com/X5qvPBa.png)

Suppose there exists a network pool and another peer is trying to enter this. As there would be no leader in this connection, they should be able to connect using either one of the nodes. The procedure would follow the same steps as before with an additional one

The entering node requests for the client list of its peers and the data. We assume that in a network, all the nodes have a synced up data. With this in mind the data is sent over to the node entering the network. They add each other to their client lists.


The node goes through the client list and attempts to connect to each one of its clients. It shares it client list, which would be the latest one, and informs them to add them to their client list, hence a connection is formed with each of the remaining clients in the pool. Again an assumption would be all the data would be in sync, hence there wouldn't be a need of sharing that part. In case of any inconsistency, we would rely on reconstructing the data during any other operation communication.

![Connecting with clients part of the client list](https://i.imgur.com/nDDMtx5.png)
