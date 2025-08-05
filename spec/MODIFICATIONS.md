# gomatrixserverlib - Research Project Modifications

This repository is a modified fork of the original gomatrixserverlib that can be found [here](https://github.com/matrix-org/gomatrixeserverlib). The code here has been modified to serve as a component in the "Achieving User Pseudonymity in Federated End-to-End Encrypted Messaging Platforms" research project. For an overview of the project, see the [main project README](https://github.com/Ijnek117/anonymous-matrix). 

---

## 1. Introduction

gomatrixserverlib is a Go library of common functions needed by matrix servers that is used extensively throughout the dendrite homeserver. 

## 2. Summary of Changes

The following changes were made to this repository for the research project:

* Introduced the `EncryptedUserID` type and modifying/created structs and functions to support it.
* Added support to fetch a remote server's TLS certificate.

## 3. Rationale for Modifications

The new EncryptedUserID type is required as once encrypted the localpart of the User ID exceeds the length limit and doesn't follow the UserID's format.

## 4. Key Files Modified

This table provides a more detailed reference to the files and code that were changed or added.

| File Path | Change Description |
| :--- | :--- |
| `eventauth.go` | Modified event authorization logic to support the `org.matrix.msc4014` room version. Added the `membershipAllowedPseudo` function to handle membership events in these rooms. A check in `createEventAllowed` comparing the room ID domain to the sender's domain was commented out but should be added back as part of greater refactoring to link senderIDs to homeserves instead of User IDs |
| `fclient/federationclient.go` | Added a new `SendEncryptedInvite` method to the `FederationClient` interface and its implementation to send invites with an encrypted User ID. Also added a `GetServerTLSCertificate` method and implementation to fetch a remote server's TLS certificate by making a request to an existing endpoint to establish a TLS connection. |
| `handlejoin.go` | Added the `HandlePseudoSendJoin` function to process incoming join events for the `org.matrix.msc4014` room version. This function validates the `mxid_mapping` and its signature, and stores a mapping from the Sender ID to a fake User ID for the remote user. Should be modified to map to server domain instead. |
| `invite.go` | Updated the `FederatedInviteClient` interface with a new `SendEncryptedInvite`. |
| `performinvite.go` | Introduced `PerformEncryptedInvite` and its input struct `PerformEncryptedInviteInput` to handle the logic for sending an invite to an `EncryptedUserID`. |
| `performjoin.go` | Modified `PerformJoin` to use the user's `senderID` instead of their User ID when making a join request. For `org.matrix.msc4014` rooms, the `mxid_mapping` content was changed to send the server domain instead of the full User ID. The `storeMXIDMappings` helper was updated to handle these new mappings and to ignore invite events. |
| `spec/encrypteduserid.go` | **(New File)** Added the`EncryptedUserID` type to represent a User ID with an encrypted localpart. This includes parsing and validation logic, ensuring the localpart is a Base64-encoded string and the domain is a valid server name. |