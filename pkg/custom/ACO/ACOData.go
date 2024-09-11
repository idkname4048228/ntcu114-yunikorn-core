package ACO

import (
	"fmt"

	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
)

func (aco *ACO) AddNode(n *objects.Node) {
	aco.Lock()
	defer aco.Unlock()
	aco.metadata.AddNode(n)
	log.Log(log.Custom).Info(fmt.Sprintf("userCount is %v", aco.metadata.NodeData.NodeCount))
}

func (aco *ACO) AddUser(ask *objects.AllocationAsk) {
	aco.Lock()
	defer aco.Unlock()
	aco.metadata.AddUser(ask)
	log.Log(log.Custom).Info(fmt.Sprintf("userCount is %v", aco.metadata.UserData.UserCount))
}

func (aco *ACO) RemoveUser(index int) {
	log.Log(log.Custom).Info("removing user")
	aco.metadata.RemoveUser(index)
	log.Log(log.Custom).Info("removed user")

}