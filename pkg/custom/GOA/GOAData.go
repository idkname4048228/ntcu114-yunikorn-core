package GOA

import (
	"fmt"

	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
)

func (goa *GOA) AddNode(n *objects.Node) {
	goa.Lock()
	defer goa.Unlock()
	goa.metadata.NodeData.NodeCount += 1
	goa.metadata.AddNode(n)
	log.Log(log.Custom).Info(fmt.Sprintf("userCount is %v", goa.metadata.NodeData.NodeCount))
}

func (goa *GOA) AddUser(ask *objects.AllocationAsk, app *objects.Application) {
	goa.Lock()
	defer goa.Unlock()
	goa.metadata.UserData.UserCount += 1
	goa.metadata.AddUser(ask)
	log.Log(log.Custom).Info(fmt.Sprintf("userCount is %v", goa.metadata.UserData.UserCount))
}

func (goa *GOA) RemoveUser(index int) {
	log.Log(log.Custom).Info("removing user")
	goa.metadata.UserData.UserCount -= 1
	goa.metadata.RemoveUser(index)
	log.Log(log.Custom).Info("removed user")

}