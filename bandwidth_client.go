// Code references : https://github.com/netsec-ethz/scion-homeworks/blob/master/bottleneck_bw_est/v1_bw_est_client.go 

package main

import (
	"flag"
	"encoding/binary"
	"fmt" //importing fmt package for printing
	"log"
	"math/rand" //importing for mathematical operations

	"sort"  //importing package for sorting
	"time"

	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"  //importing snet packages for the scion connections

	"github.com/scionproto/scion/go/lib/spath" //importing packages for finding path
	"github.com/scionproto/scion/go/lib/spath/spathmeta"
)


type Checkpoint struct {
	sent, receive int64 // sent and received
}

var (

	recvHash map[uint64]*Checkpoint // intialising datatype map to make hash in golang
	udpConnection *snet.Conn
)

func check(e1 error) {
	if e1 != nil {
		log.Fatal(e1)
	}
}

func averageBW() (float64, float64) { // finding average bottleneck bandwidth for sent and received packets

sorted := make([]*Checkpoint, ) // sorting checkpoints and 3 is number of packets sent
	i := 0
	for _, c := range recvHash {
		if c.receive != 0 {
			sorted[i] = c
			i += 1
		}
	}
  sort.Slice(sorted, func(i, j int) bool { return sorted[i].sent < sorted[j].sent })//sorting in an order

  var Sent_val, Receive_val int64 //intialising variables

  for i := 1; i < 3; i+=1 {
		Sent_val += (sorted[i].sent - sorted[i-1].sent)// finding difference intervals between sent packets
		Receive_val += (sorted[i].receive - sorted[i-1].receive)// finding difference intervals between received packets
	}

  sentbw := float64(1000*8*1e3) / (float64(Sent_val) / float64(2)) // finding bandwidth using equation size/time and conversion to MBps, 1000 is the packet size and 2 is iterations-1
  receivedbw := float64(1000*8*1e3) / (float64(Receive_val) / float64(2))

  return sentbw, receivedbw
  }

  func sendPackets() {
  var e1 error

  sendPacketBuffer := make([]byte, 1000)

  seed := rand.NewSource(time.Now().UnixNano()) //capturing time when packet is sent
	iters := 0
	for iters < 3 { //keeping 3 as number of iterations for accuracy
		iters += 1

		id := rand.New(seed).Uint64()
		_ = binary.PutUvarint(sendPacketBuffer, id)

		recvMap[id] = &Checkpoint{time.Now().UnixNano(), 0}
		_, err = udpConnection.Write(sendPacketBuffer)  // writing packets to server
		check(e1)
	}
}
func receivePackets() int {

	var e1 error
	receivePacketBuffer := make([]byte, 1000) //packet size =1000
	count := 0
  	for count < 3 {

 	 _, _, e = udpConnection.ReadFrom(receivePacketBuffer) // receiving message from server
  	check(e1)
  	time_received, _ := binary.Varint(receivePacketBuffer[n:]) // taking the time recived from received packet
 	count += 1
 }
 	return count
 }

 func main() {
	var (
		addressclient string
		addressserver string

    	e2 error
		source  *snet.Addr
		destination *snet.Addr
  )

 	flag.StringVar(&addressclient, "s", "", "Source SCION Address")
	flag.StringVar(&addressserver, "d", "", "Destination SCION Address")
	flag.Parse()

	if  len(addressclient) > 0 {
		source, err = snet.AddrFromString(addressclient)// creating udp connection
		check(e2)
	}
	if  len(addressserver) > 0 {
    	destination, err = snet.AddrFromString(addressserver)
		check(e2)
	}

    dAddr := "/run/shm/dispatcher/default.sock"
	snet.Init(source.IA, sciond.GetDefaultSCIONDPath(nil), dAddr)

  	var pathEntry *sciond.PathReplyEntry //path chosing
	var options spathmeta.AppPathSet
	options = snet.DefNetwork.PathResolver().Query(source.IA, destination.IA)

	var biggest string
	for k, entry := range options {
 	if k.String() > biggest {
    pathEntry = entry.Entry /* Choose the first random one. */
  	}
}

	fmt.Println("\nPath:", pathEntry.Path.String())
	remote.Path = spath.New(pathEntry.Path.FwdPath)
	remote.Path.InitOffsets()
	remote.NextHopHost = pathEntry.HostInfo.Host()
	remote.NextHopPort = pathEntry.HostInfo.Port 

	udpConnection, err = snet.DialSCION("udp4", source, destination)//dialling scion
	check(e2)

  	recvHash = make(map[uint64]*Checkpoint) // creating hash table using checkpoints

    sendPackets()
	count := receivePackets()

  	sentbw, receivedbw := averageBW()

  	fmt.Printf("\nSource: %s\nDestination: %s\n", addressclient, addressserver);
	fmt.Println("Bandwidth:")
	fmt.Printf("\tBW - %.3fMbps\n", sentbw)
	fmt.Println("Bottleneck Bandwidth estimate:")
	fmt.Printf("\tBW - %.3fMbps\n", receivedbw)
}