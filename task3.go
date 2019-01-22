// reference link is https://github.com/netsec-ethz/scion-homeworks/blob/master/latency/controlplane_client.go

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/hpkt"
	"github.com/scionproto/scion/go/lib/overlay"
	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/scmp"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sock/reliable"
	"github.com/scionproto/scion/go/lib/spath"
	//"github.com/scionproto/scion/go/lib/spath/spathmeta"
	"github.com/scionproto/scion/go/lib/spkt"
)

var Seed rand.Source // seeding random nos from source

// to display incase of any issues during runtime
func exceptioncheck(e1 error) {
	if e1 != nil {
		log.Fatal(e1)
	}
}

//scmp packet created here has a header indicating if the packet is to be processed by every router on the path and an error flag
// error flag shows up if packet contains any error message

func packetcheck(pkt *spkt.ScnPkt, id uint64) (*scmp.Hdr, *scmp.InfoEcho, error) {

	//scmp mandatory header carrying length of the passing scmp message along with its type and time stamp of message creation
	scmpHdr := pkt.L4.(*scmp.Hdr)

	//scmp payload carrying actual content of the message

	scmpPld := pkt.Pld.(*scmp.Payload)
	
	info := scmpPld.Info.(*scmp.InfoEcho)
	
	return scmpHdr, info, nil
}

func main() {
	var (
		addresssource string
		addressdestination string

		e2    error
		local  *snet.Addr
		remote *snet.Addr

		scmpConnection *reliable.Conn
	)

	// Fetch arguments from command line
	flag.StringVar(&addresssource, "c", "", "Source SCION Address")
	flag.StringVar(&addressdestination, "s", "", "Destination SCION Address")
	flag.Parse()

	// Create the SCION UDP socket
	if len(addresssource) > 0 {
		local, e2 = snet.AddrFromString(addresssource)
		exceptioncheck(e2)
	} 
	
	if len(addressdestination) > 0 {
		remote, e2 = snet.AddrFromString(addressdestination)
		exceptioncheck(e2)
	} 
//border router on router path 
	dispatcherAddr := "/run/shm/dispatcher/default.sock"
	snet.Init(local.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)
//source side  AS
	localAppAddr := &reliable.AppAddr{Addr: local.Host, Port: local.L4Port}
	//creating source scmp connection

	scmpConnection, _, e2 = reliable.Register(dispatcherAddr, local.IA, localAppAddr, nil, addr.SvcNone)
	exceptioncheck(e2)

	//  Path declaration to Remote
	var entryroute *sciond.PathReplyEntry
	//var choice spathmeta.AppPathSet
	//selecting path to destination from available paths b/w src and dest
	//choice := snet.DefNetwork.PathResolver().Query(local.IA, remote.IA)

//printing travel path
	fmt.Println("Path:", entryroute.Path.String())

	remote.Path = spath.New(entryroute.Path.FwdPath)
	remote.Path.InitOffsets()
	remote.NextHopHost = entryroute.HostInfo.Host()
	remote.NextHopPort = entryroute.HostInfo.Port

//dest or remote side AS creation with address and port
	remoteAppAddr := &reliable.AppAddr{Addr: remote.NextHopHost, Port: remote.NextHopPort}
	if remote.NextHopHost == nil {
		remoteAppAddr = &reliable.AppAddr{Addr: remote.Host, Port: overlay.EndhostPort}
	}

	Seed = rand.NewSource(time.Now().UnixNano())//creating seed with random ids from the  new source
	var ans int64 = 0
	
	buff := make(common.RawBytes, entryroute.Path.Mtu)
	//Now path established between src and remote
	// now we need to send and receive echo packets and note down the time

		//  created packet called inside main
		id, pkt := createrequest(local, remote)
		pktLen, e2 := hpkt.WriteScnPkt(pkt, buff)
		exceptioncheck(e2)


		t1 := time.Now()//fetching the time at which packet is sent
		_, e2 = scmpConnection.WriteTo(buff[:pktLen], remoteAppAddr)
		exceptioncheck(e2)

		n, e2 := scmpConnection.Read(buff)
		t2 := time.Now()//reading receive time after packet reaches back

		recvpkt := &spkt.ScnPkt{}
		e2 = hpkt.ParseScnPkt(recvpkt, buff[:n])
		exceptioncheck(e2)
		_, info, e2 := packetcheck(recvpkt, id)
		exceptioncheck(e2)

//calculation is done when the sent id and received matches 
		if info.Id == id { 
			 difference := (t2.UnixNano() - t1.UnixNano())
			ans = difference
		}
	var result float64 = float64(ans)

	fmt.Printf("\nSource:%s\n", addresssource);
	fmt.Printf("\nDestination:%s\n",addressdestination);
	fmt.Println("Results Obtained as follows:")
	// Print in ms, so divide by 1e6 from nano
	fmt.Printf("\tRTT - %.3fms\n", result/1e6)
	fmt.Printf("\tLatency - %.3fms\n", result/2e6)//since we take RTT as 2x Latency
}

// a function to make scmp echo request packet
func createrequest(local *snet.Addr, remote *snet.Addr) (uint64, *spkt.ScnPkt) {
	id := rand.New(Seed).Uint64() // id generation using random number func with validation

	// details of the scmp packets like id, sequence are loaded to info variable with sequence number = 0
	info := &scmp.InfoEcho{Id: id, Seq: 0}
 	//data inside scmp
	scmpMeta := scmp.Meta{InfoLen: uint8(info.Len() / common.LineLen)}

	pld := make(common.RawBytes, scmp.MetaLen+info.Len())//payload having actual content
	scmpMeta.Write(pld)//writing the message

	//mandatory header creation  with  its class, type and length of payload
	scmpHdr := scmp.NewHdr(scmp.ClassType{Class: scmp.C_General, Type: scmp.T_G_EchoRequest}, len(pld))

	pkt := &spkt.ScnPkt{ // packet contents as below
		
		//the variables seen on right side of the declarations/initializations will be used in main func 
		DstIA:   remote.IA,
		SrcIA:   local.IA,
		DstHost: remote.Host,
		SrcHost: local.Host,
		Path:    remote.Path,
		L4:      scmpHdr,
		Pld:     pld,
	}

	return id, pkt
}
