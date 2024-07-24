package vector

import (
    "math"
    "fmt"
)

type Vector struct {
    length int
    number []float64
}

func NewVector(num []float64) *Vector {
    length := len(num)
    if length <= 0 {
        panic("Invalid length")
    }

    v := WithSize(length)

    for i := 0; i < v.Length(); i++ {
        v.Set(i, num[i])
    }

    return v
}

func NewVectorByInt(num []int) *Vector {
    length := len(num)
    if length <= 0 {
        panic("Invalid length")
    }

    v := WithSize(length)

    for i := 0; i < v.Length(); i++ {
        v.Set(i, float64(num[i]))
    }

    return v
}

func WithSize(length int) *Vector {
    if length <= 0 {
        panic("Invalid length")
    }

    v := new(Vector)
    v.length = length
    v.number = make([]float64, length)

    return v
}

// The length of the Vector.
func (v *Vector) Length() int {
    return v.length
}

func (v *Vector) ToArray() (array []float64){
    array = make([]float64, 0)
    for i := 0; i < v.length; i++ {
        array = append(array, v.Get(i))
    }
    return 
}

func (v *Vector) ToIntArray() (array []int) {
    array = make([]int, 0)
    for i := 0; i < v.length; i++ {
        array = append(array, int(v.Get(i)))
    }
    return 
}

func (v *Vector) Copy() *Vector {
    result := WithSize(v.length)
    for index, value := range v.number {
        result.Set(index, value)
    }
    return result
}

// String 实现 Stringer 接口
func (v *Vector) String() string {
	return fmt.Sprintf("%v", v.number)
}

// Getter
func (v *Vector) Get(index int) float64 {
    if index < 0 || index >= v.Length() {
        panic("Invalid row size")
    }
    return v.number[index]
}

// Setter
func (v *Vector) Set(index int, data float64) {
    if index < 0 || index >= v.Length() {
        panic("Invalid row size")
    }
    v.number[index] = data
}

// Norm
func (v *Vector) Norm() float64 {
    sum := 0.0
    for i := 0; i < v.Length(); i++ {
        value := v.Get(i)
        sum += value * value
    }
    return math.Sqrt(sum)
}

// Add
func (v *Vector) Add(vector *Vector) {
    if v.length != vector.length {
        panic("Invalid vector size")
    }
    for i := 0; i < v.Length(); i++ {
        v.Set(i, v.Get(i) + vector.Get(i))
    }
}

func Add(v1 *Vector, v2 *Vector) *Vector {
    if v1.Length() != v2.Length() {
        panic("Invalid vector size")
    }
    result := WithSize(v1.Length())
    for i := 0; i < v1.Length(); i++ {
        result.Set(i, v1.Get(i) + v2.Get(i))
    }
    return result
}

// Dot
func (v1 *Vector) Dot(v2 *Vector) float64  {
    if v1.Length() != v2.Length() {
        panic("Invalid vector size")
    }
    sum := 0.0
    for i := 0; i < v1.Length(); i++ {
        sum += v1.Get(i) * v2.Get(i)
    }
    return sum
}

// Unit_vector
func (v *Vector) GetUnitVector() *Vector {
    if v.Norm() == 0 {
        return WithSize(v.length)
    }
    return Divide(v.Copy(), v.Norm())
}

// Subtract
func (v *Vector) Subtract(vector *Vector) {
    if v.Length() != vector.Length() {
        panic("Invalid vector size")
    }
    for i := 0; i < v.Length(); i++ {
        v.Set(i, v.Get(i) - vector.Get(i))
    }
}

func Subtract(v1 *Vector, v2 *Vector) *Vector {
    if v1.Length() != v2.Length() {
        panic("Invalid vector size")
    }
    result := WithSize(v1.Length())
    for i := 0; i < v1.Length(); i++ {
        result.Set(i, v1.Get(i) - v2.Get(i))
    }
    return result
}

// Devide
func Divide(v *Vector, value float64) *Vector {
    if value == 0 {
        panic("Cannot divide by zero")
    }
    for i := 0; i < v.Length(); i++ {
        v.Set(i, v.Get(i) / value)
    }
    return v
}

// Multiple
func (v *Vector) Multiple(value float64) *Vector {
    for i := 0; i < v.Length(); i++ {
        v.Set(i, v.Get(i) * value)
    }
    return v
}