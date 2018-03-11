package main

import (
	"bufio"
	//"dlog"
	"flag"
	"fmt"
	"genericsmrproto"
	"log"
	"masterproto"
	"math/rand"
	"net"
	"net/rpc"
	"runtime"
	"state"
	"time"
	"bytes"
)

var masterAddr *string = flag.String("maddr", "", "Master address. Defaults to localhost")
var masterPort *int = flag.Int("mport", 7087, "Master port.  Defaults to 7077.")
var reqsNb *int = flag.Int("q", 5000, "Total number of requests. Defaults to 5000.")
var writes *int = flag.Int("w", 100, "Percentage of updates (writes). Defaults to 100%.")
var noLeader *bool = flag.Bool("e", false, "Egalitarian (no leader). Defaults to false.")
var fast *bool = flag.Bool("f", false, "Fast Paxos: send message directly to all replicas. Defaults to false.")
var rounds *int = flag.Int("r", 1, "Split the total number of requests into this many rounds, and do rounds sequentially. Defaults to 1.")
var procs *int = flag.Int("p", 2, "GOMAXPROCS. Defaults to 2")
var check = flag.Bool("check", false, "Check that every expected reply was received exactly once.")
var eps *int = flag.Int("eps", 0, "Send eps more messages per round than the client will wait for (to discount stragglers). Defaults to 0.")
var conflicts *int = flag.Int("c", -1, "Percentage of conflicts. Defaults to 0%")
var s = flag.Float64("s", 2, "Zipfian s parameter")
var v = flag.Float64("v", 1, "Zipfian v parameter")

var N int

var successful []int

var rarray []int
var rsp []bool

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*procs)

	randObj := rand.New(rand.NewSource(42))
	zipf := rand.NewZipf(randObj, *s, *v, uint64(*reqsNb / *rounds + *eps))

	if *conflicts > 100 {
		log.Fatalf("Conflicts percentage must be between 0 and 100.\n")
	}

	master, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", *masterAddr, *masterPort))
	if err != nil {
		log.Fatalf("Error connecting to master\n")
	}

	rlReply := new(masterproto.GetReplicaListReply)

	err = master.Call("Master.GetReplicaList", new(masterproto.GetReplicaListArgs), rlReply)
	if err != nil {
		log.Fatalf("Error making the GetReplicaList RPC")
	}

	N = len(rlReply.ReplicaList) //number of servers or replicas
	servers := make([]net.Conn, N)
	var b bytes.Buffer
	//reader := bufio.NewReader(&b)
	//writer := make(bufio.Writer)
	writer := bufio.NewWriter(&b)

	rarray = make([]int, *reqsNb / *rounds + *eps)
	karray := make([]int64, *reqsNb / *rounds + *eps)
	put := make([]bool, *reqsNb / *rounds + *eps)
	//perReplicaCount := make([]int, N)
	test := make([]int, *reqsNb / *rounds + *eps) //a list of about a hundred elements. Value can be anything - a count of how many times it occurs
	for i := 0; i < len(rarray); i++ {
		r := rand.Intn(N) //random order of a particular slot
		rarray[i] = r
		/*
		if i < *reqsNb / *rounds {
			perReplicaCount[r]++ //per replica operation count
		}
		*/

		if *conflicts >= 0 {
			r = rand.Intn(100)
			if r < *conflicts {
				karray[i] = 42 //key is 42 if new random r is less than conflicts
			} else {
				karray[i] = int64(43 + i) //else 43+replica id
			}
			r = rand.Intn(100)
			if r < *writes {
				put[i] = true //it has been put/written
			} else {
				put[i] = false
			}
		} else {
			karray[i] = int64(zipf.Uint64())
			test[karray[i]]++ //test is assigned a value/count only when not put/written
		}
	}
	//log.Printf("\n Per replica count array : %v ", perReplicaCount)
	if *conflicts >= 0 {
		fmt.Println("Uniform distribution")
	} else {
		fmt.Println("Zipfian distribution:")
		fmt.Println("Test array 0-100 ",test[0:100]) //what was generated above as test[karray[i]]
	}
	//get input bidding

	//var reader = bufio.NewReader(os.Stdin)
	/*
	fmt.Println("Enter bidding value to be placed: ")
	var bid int
	var _,inperr=fmt.Scanf("%d",&bid)
	if inperr!=nil{
		fmt.Println("Bid value not gotten")
	}
	var server, err2 = rpc.DialHTTP("tcp",rlReply.ReplicaList[0])
	if err2!=nil{
		fmt.Println("Connection to server failed")
	}
	fmt.Println("Successfully connected to corresponding replica ")
	var reply = new(masterproto.GetServerReply)
	var bidargs = masterproto.GetBidPlacingArgs{BidValue:bid, BidReplica:0}
	server.Call("Server.placeBid",bidargs,reply)
	fmt.Println("Bid value sent")
	if reply.Status==true{
		fmt.Println("Bid value placed successfuly: ",bid)
	}
	*/
	
	//fmt.Println("rlReply is: ",rlReply)
	//fmt.Println("rlReply's replica list is ",rlReply.ReplicaList)

	fmt.Println("The value of N is: ",N)
	var bidarray [7][3] int
	fmt.Println("starting bid for 7 products")

	for i := 0;i<7;i++{

		for j := 0;j<N;j++{
			bidarray[i][j]=rand.Intn(N)
		}
	}
	fmt.Println(bidarray)
	var prodbid int

	for prod :=0; prod<7; prod++{
		var BidValue int
		fmt.Println("\n\n****************************************************************\n\n")
		fmt.Println("Bidding product ",prod+1)
		fmt.Println("\n\n****************************************************************\n\n")
	for rep := 0; rep < N; rep++ {
		var err error
		var i=bidarray[prod][rep]
		

		

		fmt.Println("Enter bid value to be place for product ",prod," by replica ",i+1)

		_,err=fmt.Scanf("%d",&BidValue)

		if err!=nil{
			fmt.Println("Error reading bid value")
		}

		if(BidValue<prodbid){
		fmt.Println("Bid value should be greater than previous bid. Place again: ")

		_,err=fmt.Scanf("%d",&BidValue)

		if err!=nil{
			fmt.Println("Error reading bid value")
		}
		}
		
		servers[i], err = net.Dial("tcp", rlReply.ReplicaList[i])
		if err != nil {
			log.Printf("Error connecting to replica %d\n", i)
		}
		//readers[i] = bufio.NewReader(servers[i])
		
		//reader=bufio.NewReader(servers[i])
		//writer=bufio.NewWriter(servers[i])
		

		//writers[i] = bufio.NewWriter(servers[i]) //servers are the ones that both read and write
	//}//change made here

	//successful = make([]int, N)
	leader := 0

	if *noLeader == false {
		reply := new(masterproto.GetLeaderReply)
		if err = master.Call("Master.GetLeader", new(masterproto.GetLeaderArgs), reply); err != nil {
			log.Fatalf("Error making the GetLeader RPC\n")
		}
		leader = reply.LeaderId
		log.Printf("The leader is replica %d\n", leader)
	}

	var id int32 = int32(i)
	//done := make(chan bool, N)
	args := genericsmrproto.Propose{id, state.Command{state.PUT, 0, 0}, 0} //This is something you might want to modify

	before_total := time.Now()

	//for j := 0; j < *rounds; j++ { //ignoring rounds

		//n := *reqsNb / *rounds
/*
		if *check {
			rsp = make([]bool, n)
			for j := 0; j < n; j++ {
				rsp[j] = false
			}
		}
*/
/* //Not working because of reader not being a list
		if *noLeader {
			//for i := 0; i < N; i++ {
				//go waitReplies(readers, i, perReplicaCount[i], done)
			go waitReplies(reader, i, 1, done)
			//}
		} else {
			go waitReplies(reader, leader, N, done)
		}
*/

		before := time.Now()

		//for i := 0; i < n+*eps; i++ {
			fmt.Printf("Sending proposal %d\n", id+1) //and make change here too
			args.CommandId = id
			if put[i] {
				args.Command.Op = state.PUT
			} else {
				args.Command.Op = state.GET
			}
			args.Command.K = state.Key(i)
			args.Command.V = state.Value(BidValue)

			//fmt.Println(args.Command.K,args.Command.V)
			//args.Timestamp = time.Now().UnixNano()
			if !*fast {
				if *noLeader {
					leader = rarray[i]
				}
				//writers[leader].WriteByte(genericsmrproto.PROPOSE) //some memory leak here
				//args.Marshal(writers[leader])
			} else {
				//send to everyone
				//for rep := 0; rep < N; rep++ {
					writer.WriteByte(genericsmrproto.PROPOSE)
					args.Marshal(writer)
					writer.Flush()
				}
			//}
			fmt.Println("Sent the bid placed by replica ", id+1)
			//id++
			writer.Flush()
			/*
			if i%100 == 0 {
				for i := 0; i < N; i++ {
					writers[i].Flush()
				}
			}
		//} //for loop i< n ** eps
		for i := 0; i < N; i++ {
			writers[i].Flush()
		}


		err := false
		if *noLeader {
			for i := 0; i < N; i++ {
				e := <-done
				err = e || err
			}
		} else {
			err = <-done
		}
		*/

		after := time.Now()

		fmt.Printf("Round took %v\n", after.Sub(before))
		/*
		if *check {
			for j := 0; j < n; j++ {
				if !rsp[j] {
					fmt.Println("Didn't receive", j)
				}
			}
		}
		*/

		if err != nil {
			if *noLeader {
				N = N - 1
			} else {
				reply := new(masterproto.GetLeaderReply)
				master.Call("Master.GetLeader", new(masterproto.GetLeaderArgs), reply)
				leader = reply.LeaderId
				log.Printf("New leader is replica %d\n", leader)
			}
		}
	//}

	after_total := time.Now()
	fmt.Printf("Test took %v\n", after_total.Sub(before_total))
	/*
	s := 0
	for _, succ := range successful {
		s += succ
	}

	fmt.Printf("Successful: %d\n", s)
	

	for _, client := range servers {
		if client != nil {
			//client.Close() //change made
		}
	}
	*/
	
	}
	prodbid=BidValue
	}//two for loops closed
	master.Close()
}

func waitReplies(readers []*bufio.Reader, leader int, n int, done chan bool) {
	e := false

	reply := new(genericsmrproto.ProposeReplyTS)
	for i := 0; i < n; i++ {
		if err := reply.Unmarshal(readers[leader]); err != nil {
			fmt.Println("Error when reading:", err)
			e = true
			continue
		}
		//fmt.Println("Reply value",reply.Value)
		if *check {
			if rsp[reply.CommandId] {
				fmt.Println("Duplicate reply", reply.CommandId)
			}
			rsp[reply.CommandId] = true
		}
		if reply.OK != 0 {
			successful[leader]++
		}
	}
	done <- e
}
