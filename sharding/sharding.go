package sharding

import (
    "context"
    "fmt"
    "time"
    "strconv"
    "../messages"
    "../clustering"
    "github.com/perlin-network/noise/network"
)

type ClusterPlugin struct{ *network.Plugin }

type node struct {
    port int
    shard int
    pubkey string
    privatekey string
}

var shardList,nodeList [4]node
var shardresp [10]int
var n int
var m int

func (state *ClusterPlugin) Receive(ctx *network.PluginContext) error {
    switch msg := ctx.Message().(type) {
   
   
    case *messages.ClusterMessage:
        index:=int(msg.Shard)
        shardresp[index-1]++
        fmt.Println("Cluster Message")
       
    case *messages.ClusterResponse:
        fmt.Println("Cluster Response")
        for i:=0 ; i<n ; i++ {
            if nodeList[i].port==int(msg.Port){
                nodeList[i].shard=int(msg.Shard)
                fmt.Println("Port", msg.Port)
                fmt.Println("Shard", msg.Shard)
                nodeList[i].pubkey=msg.Pubkey
                nodeList[i].privatekey=msg.Privatekey
                break
            }
        }
    }
    return nil
}

func Sharding(net *network.Network, port int) (shardList []node, shardNumber int){

    ctx := network.WithSignMessage(context.Background(), true)
    clusters := clustering.GetClusters()
    time.Sleep(1*time.Second)
    told:=false
    p:=time.Now()
    for {
       
        
        if n<4 && !told{
       
            fmt.Println(len(clusters))
            for i := 0; i < len(clusters); i++ {
                for j := 0; j < len(clusters[i]); j++ {
                    if port != clusters[i][j] {
                        //fmt.Println("Printing 2D array")
                        //fmt.Println(clusters[i][j])
                        nodeList[n]=node{port:clusters[i][j], shard:i+1}
                        n++
                        client, err := net.Client("tcp://localhost:"+strconv.Itoa(clusters[i][j]));
                        if err != nil {
                            fmt.Println(err)
                        }else{
                            //fmt.Println("Communicating Shard List")
                            client.Tell(ctx,&messages.ClusterMessage{Port: int32(clusters[i][j]),Shard: int32(i+1)})
                        }
                    } else{
                        nodeList[n]=node{port:port, shard:i+1, pubkey: net.GetKeys().PublicKeyHex(), privatekey: net.GetKeys().PrivateKeyHex()}
                        n++
                    }
                }
            }
            told=true
        }
       
        if time.Since(p)>30*time.Second{
           
            fmt.Println("Now calculating vote")
            max:=0
            for i := 1; i < 10; i++ {
                if shardresp[i]>shardresp[max]{
                    max=i
                }
            }
            for i:=0 ; i<n ; i++ {
                if nodeList[i].port==port{
                    nodeList[i].shard=max+1
                    break
                }
            }   
            shardNumber = max+1
            net.Broadcast(ctx,&messages.ClusterResponse{Port:int32(port),Shard:int32(shardNumber), Pubkey: net.GetKeys().PublicKeyHex(), Privatekey: net.GetKeys().PrivateKeyHex()})
            time.Sleep(15*time.Second)
            fmt.Println(nodeList)
           
            for i:=0 ; i<n ; i++ {
                if nodeList[i].shard==shardNumber{
                    shardList[m]=nodeList[i]
                    m++
                }
            }
            break
        }  
    }
    return shardList, shardNumber
}