package epaxos

import (
	"epaxosproto"
	"fmt"
	"genericsmr"
	"state"
	"testing"
)

func initReplica() *Replica {
	peers := make([]string, 3)
	r := &Replica{genericsmr.NewReplica(0, peers, false, false, false),
		make(chan *Propose, CHAN_BUFFER_SIZE),
		make(chan *epaxosproto.Prepare, CHAN_BUFFER_SIZE),
		make(chan *epaxosproto.PreAccept, CHAN_BUFFER_SIZE),
		make(chan *epaxosproto.PreAccept, CHAN_BUFFER_SIZE),
		make(chan *epaxosproto.Accept, CHAN_BUFFER_SIZE),
		make(chan *epaxosproto.Commit, CHAN_BUFFER_SIZE),
		make(chan *epaxosproto.PrepareReply, CHAN_BUFFER_SIZE),
		make(chan *epaxosproto.PreAcceptReply, CHAN_BUFFER_SIZE),
		make(chan *epaxosproto.PreAcceptOK, CHAN_BUFFER_SIZE),
		make(chan *epaxosproto.AcceptReply, CHAN_BUFFER_SIZE),
		make([][]*Instance, 3),
		make([]int32, 3),
		false,
		0,
		make([]int32, 3),
		nil}

	for i := 0; i < r.N; i++ {
		r.InstanceSpace[i] = make([]*Instance, 1024*1024)
		r.crtInstance[i] = 0
		r.ExecedUpTo[i] = -1
	}

	r.exec = &Exec{r}

	return r
}

func (r *Replica) MakeInstance(q, i int, seq int32, deps [3]int32) {
	command := &state.Command{state.PUT, state.Key(q), state.Value(i)}
	r.InstanceSpace[q][i] = &Instance{command, 0, epaxosproto.COMMITTED, seq, deps, nil, 0, 0}
}

func TestExec(t *testing.T) {

	r := initReplica()

	r.MakeInstance(0, 0, 2, [3]int32{0, 0, 0})
	r.MakeInstance(1, 0, 1, [3]int32{0, 0, 0})
	r.MakeInstance(2, 0, 0, [3]int32{0, 0, 0})

	r.MakeInstance(0, 1, 0, [3]int32{0, 1, 0})
	r.MakeInstance(1, 1, 2, [3]int32{0, 0, 1})
	r.MakeInstance(2, 1, 0, [3]int32{1, 0, 0})

	r.MakeInstance(0, 2, 0, [3]int32{1, 1, 1})
	r.MakeInstance(1, 2, 0, [3]int32{1, 1, 2})
	r.MakeInstance(2, 2, 0, [3]int32{2, 1, 1})

	r.MakeInstance(0, 3, 0, [3]int32{2, 2, 2})
	r.MakeInstance(1, 3, 0, [3]int32{0, 0, 0})
	r.MakeInstance(2, 3, 0, [3]int32{0, 0, 0})

	r.MakeInstance(0, 4, 1, [3]int32{3, 5, 0})
	r.MakeInstance(1, 4, 2, [3]int32{0, 0, 0})
	r.MakeInstance(2, 4, 3, [3]int32{0, 0, 0})

	r.MakeInstance(0, 5, 4, [3]int32{4, 5, 5})
	r.MakeInstance(1, 5, 5, [3]int32{5, 5, 5})
	r.MakeInstance(2, 5, 6, [3]int32{5, 0, 5})

	r.exec.executeCommand(0, 5)
	r.exec.executeCommand(0, 5)

	'''
	executeCommand in epaxos-exec.go
	The structure Exec is defined in epaxos-exec.go and it is r *Replica (an arbitrary number of replicas)
	The strucure Replica is defined in genericsmr.go
	InstanceSpace of replica is created in initReplica() above
	The structure Instance is in epaxos.go
	Its main attribute seems to be Cmds[] which is state.Command which has operation, key and value
	The MakeInstance function is in this file and the operation is PUT, key is q and value is i
	Next, take a look at the function executeCommands() in epaxos.go
	Looked at executeCommands() - it marks off committed and executed commands with sleep in between and recovers if commit takes too long. It calls executeCommand() in epaxos-exec.go
	executeCommand() just returns true if the command is executed and false if it is NOT committed. 
	executeCommand() just looks like a sanity check for executed state. So I think executeCommands() in epaxos. go is where we would have to include the bidding system read and write logic after building a dataset
 	'''

	fmt.Println("Test ended\n")
}
