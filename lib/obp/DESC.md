# Orwell Binary Protocol

OBP is a general-purpose binary protocol for symmetric communication over a socket.
It assumes that the connection is octet-based and handles/corrects potential bit flip errors.
The design goal was to create **the simplest conceivable protocol** that can efficiently handle **two-way asynchronous request/response communication**.

## Messages

Peers communicate by exchanging *frames*. 
A frame consists of:

1. Context - the request/response pair identifier,
3. Payload - the frame contents.

### Context

Peers communicate by sending *request frames* and receiving *response frames*.
Each frame has an associated number, named the *frame context*.
A response to a given frame must have a context number matching the request context.

Matching request/response context pair is composed of numbers `{2k, 2k+1}`. 
In other words, request contexts are even and corresponding response contexts are odd, greater by exactly one.

Once the peer receives a response to a request, the context number can be reused.
Context numbers can be used in any order.

### Payload

Payload is a variable-length string of bytes.
The maximum payload length is protocol-dependant.

## Binary representation

All integers are represented as variable-length encoded unsigned 64-bit integers (varuints).
The following method for integer compression is used:

- If the number is less than 253, can be passed as-is, taking 1 byte.
- If the number is less than 2^16, it is prefixed with `0xFD` (253), taking 3 bytes.
- If the number is less than 2^32, it is prefixed with `0xFE` (254), taking 5 bytes.
- If the number is less than 2^64, it is prefixed with `0xFF` (255), taking 9 bytes.

It is not legal to pass smaller numbers with bigger prefixes. 
For instance, `0xFF` followed by eight `0x00` bytes is not a legal zero value.

The numbers are **obviously** passed as big-endian. If you like little-endian, visit your doctor immediately.

### Behaviour

OBP is simple, therefore errors are not tolerated. 
Any unexpected behaviour should result with immediate connection termination.

Examples include:

- sending requests with the same context numbers before the response to the first one is sent,
- sending responses to not sent requests,
- too long payload length.


### Frame structure

The frame is passed as a 3-tuple of the following structure:

1. Context, encoded as varuint,
3. Payload length, encoded as varuint,
4. Payload, a string of bytes of the length specified above.

And there you go. That's everything.