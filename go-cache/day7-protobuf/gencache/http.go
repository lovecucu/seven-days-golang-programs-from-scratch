package gencache

import (
	"fmt"
	"gencache/consistenthash"
	pb "gencache/gencachepb"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self        string // 记录主机名/ip和端口
	basePath    string // 前缀地址，默认是/_geecache/
	mu          sync.Mutex
	peers       *consistenthash.Map    // 一致性哈希获取key的节点
	httpGetters map[string]*httpGetter // 远程节点和对应的httpGetter的映射
}

// New一个实例
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 判断请求路径是否包含basePath
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 { // 判断uri是否合法
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]
	group := GetGroup(groupName)
	if group == nil { // 判断group是否存在
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil { // 判断缓存获取是否失败
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 使用proto.Marshal()编码HTTP响应
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body) // 返回序列化后的proto格式数据
}

// 添加节点，同时初始化好httpGetters
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// 选取peer
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// 获取PeerGetter
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true // 返回peer对应的PeerGetter
	}
	return nil, false
}

type httpGetter struct {
	baseURL string
}

// 根据httpGetter获取缓存值
func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	// 拼接完整url
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)
	// 发送Get请求获取key的值
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 响应异常直接返回报错
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}

	// 读取返回值
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	// 使用proto.Unmarshal()解码HTTP响应
	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}

	return nil
}

// 判断httpGetter是否实现了PeerGetter接口
var _ PeerGetter = (*httpGetter)(nil)
