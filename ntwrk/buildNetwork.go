package ntwrk

import (
	"fmt"
    "flag"
    "strings"
    "../messages"
    "github.com/perlin-network/noise/crypto/ed25519"
    "github.com/perlin-network/noise/log"
    "github.com/perlin-network/noise/network"
    "github.com/perlin-network/noise/network/discovery"
    "github.com/perlin-network/noise/types/opcode"
)

type ClusterPlugin struct{ *network.Plugin }

func BuildNetwork(port int) (net *network.Network){

	// process other flags
    portFlag := flag.Int("port", port, "port to listen to")
    hostFlag := flag.String("host", "localhost", "host to listen to")
    protocolFlag := flag.String("protocol", "tcp", "protocol to use (kcp/tcp)")
    peersFlag := flag.String("peers", "", "peers to connect to")
    flag.Parse()
	
    port = int(*portFlag)
    host := *hostFlag
    protocol := *protocolFlag
    peers := strings.Split(*peersFlag, ",")

    keys := ed25519.RandomKeyPair()

    log.Info().Msgf("Private Key: %s", keys.PrivateKeyHex())
    log.Info().Msgf("Public Key: %s", keys.PublicKeyHex())

    opcode.RegisterMessageType(opcode.Opcode(1000), &messages.ClusterMessage{})
    opcode.RegisterMessageType(opcode.Opcode(1001), &messages.ClusterResponse{})

    builder := network.NewBuilder()
    builder.SetKeys(keys)
    builder.SetAddress(network.FormatAddress(protocol, host, uint16(port)))
    // Register peer discovery plugin.
    builder.AddPlugin(new(discovery.Plugin))

    fmt.Println("**********************************************************************")
    fmt.Println(peers)
    fmt.Println("**********************************************************************")

    // Add custom chat plugin.
    builder.AddPlugin(new(ClusterPlugin))

    net, err := builder.Build()
   
    if err != nil {
        log.Fatal().Err(err)
        return
    }

    go net.Listen()

    return net
}