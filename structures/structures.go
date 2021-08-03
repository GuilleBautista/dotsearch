package structures

type Node_t struct {  
    Name string
    H int
}

type Graph_t struct {
    V []Node_t
    Matrix [][]int
    Start int
    Goal int
    Directed bool
}

type Path_t struct {
    Cost int
    Length int
    Node_names []string
    Node_ids []int
}
