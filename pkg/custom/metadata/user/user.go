package UserData

import (
	"math"
	// "fmt"
	"sync"

	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
)

type UserData struct {
	ResourceTypes  	[]string

	UserCount      	int
	UserNames		[]string
	UserAskMap 		map[string]*AskData

	sync.RWMutex
}

type AskData struct {
	AskCount		int
	UserAsks	 	[]float64
	Requests		[]*objects.AllocationAsk
}

func NewUserData(ResourceTypes []string) *UserData {
	return &UserData{
		UserCount:  0,
		UserNames: make([]string, 0),
		ResourceTypes:   ResourceTypes,
		UserAskMap : make(map[string]*AskData),
	}
}

// Parse the vcore and memory in node
func (userData *UserData) AddUser(ask *objects.AllocationAsk) {
	// log.Log(log.Custom).Info("userdata add user")	

	userData.Lock()
	defer userData.Unlock()
	if ask == nil {
		log.Log(log.Custom).Info("request is nil")
		return 
	}

	value, exist := userData.UserAskMap[ask.GetApplicationID()]

	if !exist {
		userData.UserCount += 1;

		askdata := &AskData{
			AskCount: 0, 
			UserAsks: make([]float64, len(userData.ResourceTypes)),
			Requests: make([]*objects.AllocationAsk, 0),
		}
		userData.UserAskMap[ask.GetApplicationID()] = askdata
		userData.UserNames = append(userData.UserNames, ask.GetApplicationID())

		value = userData.UserAskMap[ask.GetApplicationID()]
		value.UserAsks = userData.praseAskLimit(ask)
	}

	value.AskCount += 1
	value.Requests = append(value.Requests, ask)
}

func (userData *UserData) praseAskLimit(ask *objects.AllocationAsk) []float64{
	userAsk := make([]float64, len(userData.ResourceTypes))	

	curResource := ask.GetAllocatedResource().Resources
	for index, targetType := range userData.ResourceTypes {
		userAsk[index] += float64(curResource[targetType])	
	}
	return userAsk
}

func (userData *UserData) GetUserAsks() [][]float64{
	userData.Lock()
	defer userData.Unlock()
	names := userData.UserNames
	askMap := userData.UserAskMap
	userAsks := make([][]float64, 0)
	for _, name := range names {
		ask := askMap[name].UserAsks
		userAsks = append(userAsks, ask)	
	}
	return userAsks
}

func (userData *UserData) GetUserAskCount(name string) int{
	return userData.UserAskMap[name].AskCount
}

func (userData *UserData) GetName(index int) string {
	return userData.UserNames[index]
}

func (userData *UserData) Update() {
	names := userData.UserNames

	for i := len(names) - 1; i >= 0; i-- {
		name := names[i]
		if userData.UserAskMap[name].AskCount == 0 {
			delete(userData.UserAskMap, name)
			userData.RemoveUser(i)
		}
	}
}

func (userData *UserData) PopAsks(user string, amount int) []*objects.AllocationAsk{
	askData := userData.UserAskMap[user]

	asks := make([]*objects.AllocationAsk, 0)
	elements := int(math.Min(float64(amount), float64(askData.AskCount)))

	for len(asks) != elements {
		asks = append(asks, askData.Requests[0])
		askData.Requests = askData.Requests[1:]
		askData.AskCount -= 1;
	}

	// log.Log(log.Custom).Info(fmt.Sprintf("asks length is %v", len(asks)))

	return asks
}

func (userData *UserData) RemoveUser(index int) {
	userData.UserCount -= 1
	userData.UserNames = append(userData.UserNames[:index], userData.UserNames[index + 1:]...)
}

// make test easy 
func (userData *UserData) AddUserDirectly(userName string, userAsk []float64) {
	userData.UserCount += 1
	userData.UserNames = append(userData.UserNames, userName)
}
