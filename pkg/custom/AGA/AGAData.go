package AGA

import (
	"fmt"

	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
)

func (aga *AGA) AddNode(n *objects.Node) {
	aga.Lock()
	defer aga.Unlock()
	aga.metadata.AddNode(n)
	log.Log(log.Custom).Info(fmt.Sprintf("nodeCount is %v", aga.metadata.NodeData.NodeCount))
}

func (aga *AGA) AddUser(ask *objects.AllocationAsk, app *objects.Application) {
	aga.Lock()
	defer aga.Unlock()
	aga.metadata.AddUser(ask, app)
	log.Log(log.Custom).Info(fmt.Sprintf("userCount is %v", aga.metadata.UserData.UserCount))
}

func (aga *AGA) RemoveUser(index int) {
	log.Log(log.Custom).Info("removing user")
	aga.metadata.RemoveUser(index)
	log.Log(log.Custom).Info("removed user")

}