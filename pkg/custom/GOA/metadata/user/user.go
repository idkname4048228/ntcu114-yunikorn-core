package UserData

import (
	"sync"

	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
)

type UserData struct {
	UserCount      	int
	Requests        	[]string
	UserNames		[]string
	ResourceTypes  	[]string
	UserAsks	 	[][]float64

	sync.RWMutex
}

func NewUserData(ResourceTypes []string) *UserData {
	return &UserData{
		UserCount:  0,
		Requests:  make([]string, 0),
		UserNames: make([]string, 0),
		ResourceTypes:   ResourceTypes,
		UserAsks :make([][]float64, 0),
	}
}

// Parse the vcore and memory in node
func (userData *UserData) AddUser(ask *objects.AllocationAsk, app *objects.Application) {
	log.Log(log.Custom).Info("userdata add user")	

	userData.Lock()
	defer userData.Unlock()
	if ask == nil {
		log.Log(log.Custom).Info("request is nil")
		return 
	}

	userData.Requests = append(userData.Requests, ask.GetAllocationKey()) 
	userData.UserCount += 1;
	
	userAsk := make([]float64, len(userData.ResourceTypes))	

	curResource := ask.GetAllocatedResource().Resources
	for index, targetType := range userData.ResourceTypes {
		userAsk[index] += float64(curResource[targetType])	
	}

	userData.UserAsks = append(userData.UserAsks, userAsk)
}

func (userData *UserData) RemoveUser(index int) {
	userData.Requests = append(userData.Requests[:index], userData.Requests[index + 1:]...)
	userData.UserAsks = append(userData.UserAsks[:index], userData.UserAsks[index + 1:]...)
}
// make test easy 
func (userData *UserData) AddUserDirectly(userName string, userAsk []float64) {
	userData.UserCount += 1
	userData.UserNames = append(userData.UserNames, userName)
	userData.Requests = append(userData.Requests, userName) 

	userData.UserAsks = append(userData.UserAsks, userAsk)
}
