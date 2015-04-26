Orwell
======

P2P EC (eventually consistent) key-value storage with cryptographically verified ownership and update capabilities. Designed to run in an untrusted network, handling byzantine failures and malicious nodes.

Key features:
* value ownership - once you insert your value, only you can change it. Updates are signed with ECDSA.
* strong consistency for the keyspace, achieved by blockchain enforced consensus.
* eventual consistency for the values, with timestamps for value expiration dates.
* designed for mobile - value resolving mechanism requires a fixed ~8.55MiB/year, as compared to tens of GiB per year in BTC.  

Part of my MSc thesis about decentralized DNS-like systems. Work in progress.