package main

import "fmt"
import "github.com/RafaelMarinheiro/paxos"

func main(){
	var a = paxos.Message{From: 0, To: 1, proposeNumber: 2, proposeValue: 3}
	fmt.Printf("Running paxos!\n")
}