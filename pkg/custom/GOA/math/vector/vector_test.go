package vector

import (
	"testing"
	"fmt"

	"github.com/apache/yunikorn-core/pkg/log"
)

func TestAdd(t *testing.T) {
	v := NewVectorByInt([]int{4, 0})
	tmp := NewVectorByInt([]int{4, 0}) 
	gravity := NewVector([]float64{0, 0.010000000000000002})
	
	v.Add(gravity.Multiple(0.2))
	log.Log(log.Custom).Info(fmt.Sprintf("%v + %v is %v", tmp, gravity, v))
}