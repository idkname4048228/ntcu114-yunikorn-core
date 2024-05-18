package custom

var roundRobin RoundRobin;

func Init(){
	roundRobin = *NewRoundRobin();
}

func GetRoundRobin() *RoundRobin{
	return &roundRobin;
}