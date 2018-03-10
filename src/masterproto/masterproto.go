package masterproto

type RegisterArgs struct {
	Addr string
	Port int
}

type RegisterReply struct {
	ReplicaId int
	NodeList  []string
	Ready     bool
}

type GetLeaderArgs struct {
}

type GetLeaderReply struct {
	LeaderId int
}

type GetServerReply struct{
	Status bool
}

type GetReplicaListArgs struct {
}

type GetReplicaListReply struct {
	ReplicaList []string
	Ready       bool
}

type GetBidPlacingArgs struct{
	BidValue int
	BidReplica int
}
