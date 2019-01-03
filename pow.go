package main

import (
	"fmt"
	"crypto/sha256"
	"encoding/hex"
	"./serialize"
	"log"
    "net"
    "strings"
    "strconv"
    "time"
)

const difficulty=5

//perform hashing
func encode(serialized string) string{
	// fmt.Println("String is"+serialized)
    h := sha256.New()
    h.Write([]byte(serialized))
    sha256_hash := hex.EncodeToString(h.Sum(nil))
    return sha256_hash
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}

func main() {
	
	addr:=GetOutboundIP().String()
	fmt.Println(addr)

	fmt.Println("Doing Proof of work")

	nonce,result:=Pow()
	fmt.Println()
	fmt.Println("***************")
	fmt.Println("Final nonce: "+strconv.Itoa(nonce))
	fmt.Println("Final hash: "+result)
}

func Pow()(nonce int, result string){

	nonce=0  //set this to previous nonce

	var startTime = time.Now()

	result=encode(serialize.Serialize(GetOutboundIP().String(),"32132132132132132113132123",0)+strconv.Itoa(nonce))
	
	for IsValidResult(result)==false{
		fmt.Println("Nonce was "+ strconv.Itoa(nonce))
		fmt.Println("Hash was "+ result)
		nonce++
		result=encode(serialize.Serialize(GetOutboundIP().String(),"32132132132132132113132123",0)+strconv.Itoa(nonce))
		fmt.Println()
	}
	fmt.Println()
	fmt.Println()
	fmt.Println()
	var duration = time.Since(startTime)
	fmt.Println(duration)
	return nonce, result

}

func IsValidResult(hash string) bool{
	prefix := strings.Repeat("0", difficulty)
return strings.HasPrefix(hash, prefix)
}