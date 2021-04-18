package sCache

import pb "sCache/sCache/sCachepb"

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter,ok bool)
}

type PeerGetter interface {
	Get(in *pb.Request,out *pb.Response) error
}
