package GOA

import (
	"fmt"

	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
	"github.com/apache/yunikorn-core/pkg/log"
)

var goa *GOA

func GetGOA() *GOA{
	return goa
}

func Init(){
	goa = NewGOA();
}

type GOA struct {
	nodes              *Nodes // use for round robin
}

func NewGOA() *GOA{
	return &GOA{
		nodes: NewNodes(),
	}
}

type Nodes []*objects.Node;

func NewNodes() *Nodes{
	tmp := make(Nodes, 0);
	return &tmp;
}

func (goa *GOA) AddNode(node *objects.Node){
	*goa.nodes = append(*goa.nodes, node);
	log.Log(log.Custom).Info(fmt.Sprintf("GOA got node: %v", node.NodeID));
}
