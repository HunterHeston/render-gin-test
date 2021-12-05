# Value Store

The top level directory contains an interface for a value store. All sub-directories implement that interface. 

## Current implementations
* In memory
  * Simple implementation using a map as the back end data store. 
* Firebase
  * Uses Firebase Firestore as the back end storage system.