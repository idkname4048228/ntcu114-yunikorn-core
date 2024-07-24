package UserData

import (
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
	"github.com/apache/yunikorn-core/pkg/log"
)

type UserData struct {
	UserCount      	int
	AppIDs        	[]string
	UserNames		[]string
	ResourceTypes  	[]string
	UserAsks	 	[][]float64
}

func NewUserData(ResourceTypes []string) *UserData {
	return &UserData{
		UserCount:  0,
		AppIDs:  make([]string, 0),
		UserNames: make([]string, 0),
		ResourceTypes:   ResourceTypes,
		UserAsks :make([][]float64, 0),
	}
}

// Parse the vcore and memory in node
func (userData *UserData) AddUser(app *objects.Application) {
	userData.AppIDs = append(userData.AppIDs, app.ApplicationID) 
	userData.UserNames = append(userData.UserNames, app.GetUser().User)
	userData.UserCount += 1;
	
	if app.GetAllRequests() == nil {
		log.Log(log.Custom).Info("request is nil")
	}else{
		log.Log(log.Custom).Info("request exist")
	}

	userAsk := make([]float64, len(userData.ResourceTypes))	

	for _, request := range app.GetAllRequests(){
		curResource := request.GetAllocatedResource().Resources
		for index, targetType := range userData.ResourceTypes {
			userAsk[index] += float64(curResource[targetType])	
		}
	}

	userData.UserAsks = append(userData.UserAsks, userAsk)
}

// make test easy 
func (userData *UserData) AddUserDirectly(userName string, userAsk []float64) {
	userData.UserCount += 1
	userData.UserNames = append(userData.UserNames, userName)
	userData.AppIDs = append(userData.AppIDs, userName) 

	userData.UserAsks = append(userData.UserAsks, userAsk)
}
