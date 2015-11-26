package main

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
)

//Circle struct
type Circle struct {
	Nodes Nodes
	sync.Mutex
}

//Node struct
type Node struct {
	ID     string
	HashID uint32
}

//Nodes array
type Nodes []*Node

//NewCircle create an empty circle object of type Nodes
func NewCircle() *Circle {
	return &Circle{Nodes: Nodes{}}
}

func hashID(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

//AddNode function
func (c *Circle) AddNode(id string) {
	c.Lock()
	defer c.Unlock()

	node := NewNode(id)
	c.Nodes = append(c.Nodes, node)

	sort.Sort(c.Nodes)
}

//search for the host name in the circle
func (c *Circle) searchhost(id string) int {
	searchHost := func(it int) bool {
		return c.Nodes[it].HashID >= hashID(id)
	}
	return sort.Search(c.Nodes.Len(), searchHost)
}

//Get server name function
func (c *Circle) Get(id string) string {
	i := c.searchhost(id)
	if i >= c.Nodes.Len() {
		i = 0
	}
	return c.Nodes[i].ID
}

//NewNode struct
func NewNode(id string) *Node {
	return &Node{
		ID:     id,
		HashID: hashID(id),
	}
}

func (n Nodes) Less(i, j int) bool {
	return n[i].HashID < n[j].HashID
}

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func main() {

	r := NewCircle()

	// array of servers
	hostname := []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:3002"}

	//array of input objects as key and value
	inputarr := []string{"1,a", "2,b", "3,c", "4,d", "5,e", "6,f", "7,g", "8,g", "9,h", "10,i"}

	// adding hostname in the ring
	for i := 0; i < len(hostname); i++ {
		r.AddNode(hostname[i])
	}

	// sending objects across servers.
	for i := 0; i < len(inputarr); i++ {
		splitkeyval := strings.Split(inputarr[i], ",")
		insertnode1 := r.Get(splitkeyval[0])
		putrestcall(insertnode1, splitkeyval[0], splitkeyval[1])
		fmt.Println()
		getidcall(insertnode1, splitkeyval[0])
	}

}

// Consuming server PUT call
func putrestcall(hostname string, ikey string, ivalue string) {
	urlreq := hostname + "/keys/" + ikey + "/" + ivalue
	fmt.Printf("\n PUT URL: %s", urlreq)
	fmt.Printf("\n Key: %s and value: %s is inserted in server %s", ikey, ivalue, hostname)
	req, _ := http.NewRequest("PUT", urlreq, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

//Consuming server GET by id call
func getidcall(hostname string, ikey string) {
	urlreq := hostname + "/keys/" + ikey
	fmt.Printf("\n GET URL: %s", urlreq)
	fmt.Printf("\n Fetching key: %s from server: %s", ikey, hostname)
	req, _ := http.NewRequest("GET", urlreq, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("\n Response:", string(body))
	defer resp.Body.Close()
}
