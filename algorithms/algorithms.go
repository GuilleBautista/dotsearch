package algorithms

import (
    "container/list"
    s "github.com/guillebautista/dotsearch/structures"
)

//Each algorithm will implement the solver interface
type Solver interface {
    Solve() s.Path_t
    Generate_succesors() *list.List
}

//Common find function for all the algorithms
func find(elem int, list []int) bool {
    for i := range list {
        if elem==list[i] { return true }
    }
    return false
}
