package main

import(
	"fmt"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
	"math"
)
var shardName string ="Shard "
//var shard2 string ="Shard 2"

var currentBlock int=0
//shard 1
var nodes [8]string
// var node1 string="127.0.0.1 8080"
// var node2 string="127.0.0.1 8081"
// var node3 string="127.0.0.1 8082"
// var node4 string="127.0.0.1 8083"

// //shard 2
// var node5 string="127.0.0.1 8084"
// var node6 string="127.0.0.1 8085"
// var node7 string="127.0.0.1 8086"
// var node8 string="127.0.0.1 8087"
var nodesHash [8]string
var age map[string]float64
//var age [8]float64
var hashScore map[string]int



func main(){
// func LeaderElection(nodes []node, shardNumber int) {

	nodes[0]="127.0.0.1 8080"
	nodes[1]="127.0.0.1 8081"
	nodes[2]="127.0.0.1 8082"
	nodes[3]="127.0.0.1 8083"
	nodes[4]="127.0.0.1 8084"
	nodes[5]="127.0.0.1 8085"
	nodes[6]="127.0.0.1 8086"
	nodes[7]="127.0.0.1 8087"
	hashScore = make(map[string]int)
	age = make(map[string]float64)

	// hashScore = make(map[node]int)
	// age = make(map[node]float64)	

	h := sha256.New()
	for i:=0;i<8;i++{
		age[nodes[i]]=1.0
		h.Write([]byte(nodes[i]))
		nodesHash[i]=hex.EncodeToString(h.Sum(nil))
	}

	startTime:=time.Now()
	
	//t:=int(time.Now().Unix())
	
	
	
	for{

		if time.Since(startTime)>10*time.Second{
			fmt.Println("New Leader Epoch")
			fmt.Println("Electing Leader")
			h.Write([]byte(shardName + shardNumber + strconv.Itoa(currentBlock)))
			identity1:=hex.EncodeToString(h.Sum(nil))
			
			for i:=0;i<8;i++{
				hashScore[nodes[i]]=64-levenshtein(nodesHash[i],identity1)
			}
			// h.Write([]byte(shard2+strconv.Itoa(currentBlock)))
			// identity2:=hex.EncodeToString(h.Sum(nil))
			// var m1 map[string]int
			// m1 = make(map[string]int)
			// m1[node1]=levenshtein(node1Hash,identity1)
			// m1[node2]=levenshtein(node2Hash,identity1)
			// m1[node3]=levenshtein(node3Hash,identity1)
			// m1[node4]=levenshtein(node4Hash,identity1)
			// m1[node5]=levenshtein(node5Hash,identity1)
			// m1[node6]=levenshtein(node6Hash,identity1)
			// m1[node7]=levenshtein(node7Hash,identity1)
			// m1[node8]=levenshtein(node8Hash,identity1)

			fmt.Println("Similarity score is:",hashScore)
			leader:=decideLeader()
			// /fmt.Println()
			//communicateLeader(leader)
			fmt.Println("Leader is ",leader)
			currentBlock++
			startTime=time.Now()
			fmt.Println()
		}
	}
}



func levenshtein(str1, str2 string) int {

    s1len := len(str1)
	s2len := len(str2)
	
    column := make([]int, len(str1)+1)
 
    for y := 1; y <= s1len; y++ {
        column[y] = y
    }
    for x := 1; x <= s2len; x++ {
        column[0] = x
        lastkey := x - 1
        for y := 1; y <= s1len; y++ {
            oldkey := column[y]
            var incr int
            if str1[y-1] != str2[x-1] {
                incr = 1
            }
 
            column[y] = minimum(column[y]+1, column[y-1]+1, lastkey+incr)
            lastkey = oldkey
        }
    }
    return column[s1len]
}
 
func minimum(a, b, c int) int {
    if a < b {
        if a < c {
            return a
        }
    } else {
        if b < c {
            return b
        }
    }
    return c
}

func decideLeader() string{

	//var cummulativeScore [8]float64
	var cummulativeScore map[string]float64
	
	cummulativeScore = make(map[string]float64)
	for i:=0;i<8;i++{
		cummulativeScore[nodes[i]]=toFixed(0.4*float64(hashScore[nodes[i]])+0.6*float64(age[nodes[i]]),4)
	}
	fmt.Println("Age:",age)
	fmt.Println("Leader Competence Score:",cummulativeScore)
	var max float64=0
	var index int=-1
	for i:=0;i<8;i++{
		if cummulativeScore[nodes[i]]>max{
			max=cummulativeScore[nodes[i]]
			index=i
		}	else if cummulativeScore[nodes[i]]==max{
			if age[nodes[i]]>age[nodes[index]]{
				index=i
			}	
		}
	}
	for i:=0;i<8;i++{
		if i!=index{
			age[nodes[i]]=toFixed(age[nodes[i]]+float64(0.1),2)
		} else{
			age[nodes[i]]=0.5
		}
	}
	//fmt.Println()
	//fmt.Println("Age:",age)

	return nodes[index]
}

//rounding up values
func round(num float64) int {
    return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
    output := math.Pow(10, float64(precision))
    return float64(round(num * output)) / output
}

func communicateLeader(leader node) {
	
}