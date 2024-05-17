## BTC Handshake

This program establishes a handshake with any other BTC peer that handles the [Version Handshake](https://en.bitcoin.it/wiki/Version_Handshake) described in the Bitcoin wiki:

> On connect, version and verack messages are exchanged, in order to ensure compatibility between peers.

Usage:

- Starts a connection and initiates the handshake with a peer:

To initiate the handshake, you need to have the peer addr and the port of the peer you want to connect with, you can have a list of peers address from here [bitnodes.io](https://bitnodes.io/nodes/?q=United%20States), or you can install [btcd](https://github.com/btcsuite/btcd) and quickly bootstrap a local btcd node executin `btcd` in your shell, which will listen on `0.0.0.0:8333`

```sh
go run ./cmd/... --peer-addr=143.110.175.248 --peer-port=8333
```

After that, we will send a `Version` message and the peer will respond us with its message as well as a `VerAck`, the handshake process is specified here: https://en.bitcoin.it/wiki/Version_Handshake

The output look like this:

```sh
sent 127 bytes (total 127)...

received from remote (150 bytes): 0xf9beb4d976657273696f6e000000000066000000726a8eb57f1101000d040000000000002ea8466600000000010000000000000000000000000000000000ffff8f6eaff8208d0d04000000000000000000000000000000000000000000000000c82f5ba45c0b663c102f5361746f7368693a302e32302e312f04e00c0001f9beb4d976657261636b000000000000000000005df6e0e2
remote's message:
[magic=Main] [command=version] < [number=70015] [services=1037] [ts=1715906606] [recv=< [services=00000001] [ip=143.110.175.248] [port=8333] >] [from=< [services=10000001101] [ip=0.0.0.0] [port=0] >] [nonce=4352178582422499272] [user-agent=/Satoshi:0.20.1/] [start-height=843780] [relay=true] >

remote's message
[magic=Main] [command=verack] <  >

sent 24 bytes (total 24)...
```

- Waits for a handshake and respond it:

Running the project without any infor about a external peer it will start listen for active connections and for handshakes

```sh
go run ./cmd/...
```

The project will start listening on TCP port 8080, so you can bootstrap a [btcd](https://github.com/btcsuite/btcd) node locally with the command `btcd -a 0.0.0.0:8080` (the flag `-a` add a peer to connect with at startup) then it will, at startup, start a version handshake process with our node, the [btcd](https://github.com/btcsuite/btcd) output logs will appear a line like this:

```sh
2024-05-16 21:09:34.014 [INF] SYNC: New valid peer 0.0.0.0:8080 (outbound) (btc/eclesios-node)
```

and ours output logs will look like this:

```sh
listening on 0.0.0.0:8080

received from remote (137 bytes): 0xf9beb4d976657273696f6e00000000007100000039d26161801101004d040000000000004eae466600000000000000000000000000000000000000000000ffff000000001f904d04000000000000000000000000000000000000000000000000b1107ac9995652291b2f627463776972653a302e352e302f627463643a302e32342e322f3e80030001
remote's message:
[magic=Main] [command=version] < [number=70016] [services=1101] [ts=1715908174] [recv=< [services=00000000] [ip=0.0.0.0] [port=8080] >] [from=< [services=10001001101] [ip=0.0.0.0] [port=0] >] [nonce=2977537522155524273] [user-agent=/btcwire:0.5.0/btcd:0.24.2/] [start-height=229438] [relay=true] >

sent 127 bytes (total 127)...
sent 24 bytes (total 24)...
```
