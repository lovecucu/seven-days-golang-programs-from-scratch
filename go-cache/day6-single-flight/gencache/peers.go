package gencache

// 根据key获取PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 根据group和key获取对应的缓存值
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
