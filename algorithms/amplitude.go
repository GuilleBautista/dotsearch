package algorithms

import (
    "fmt"
    "container/list"
    s "github.com/guillebautista/dotsearch/structures"
    "errors"
)

type Amplitude struct {
    Graph s.Graph_t
    Solution s.Path_t
}

//Simple function to generate all succesors of a node without any further checks
func (a Amplitude) Generate_succesors(adj_matrix [][]int, path s.Path_t, last_v int, v []s.Node_t) *list.List {
    succesors:=list.New()
    var newpath s.Path_t
    for i := range adj_matrix {
        //If connected
        if adj_matrix[last_v][i] >= 0 {
            newpath.Cost = adj_matrix[last_v][i] + path.Cost
            newpath.Length=path.Length+1
            newpath.Node_names = make([]string, newpath.Length+1)
            newpath.Node_ids = make([]int, newpath.Length+1)

            copy(newpath.Node_ids, path.Node_ids)
            newpath.Node_ids[newpath.Length] = i
            succesors.PushBack( newpath )
        }
    }
    return succesors
}

//First in amplitude
func (a Amplitude) Solve()(s.Path_t, error){
    var expanded_nodes int
    //var time int
    //Gather the data
    start := a.Graph.Start
    goal := a.Graph.Goal
    adj_matrix := a.Graph.Matrix
    v := a.Graph.V
    path_list := list.New()
    
    path_list.PushBack( s.Path_t{
        Cost: 0,
        Length: 0,
        Node_names: []string{},
        Node_ids: []int{start},
    } )

    var path s.Path_t
    var last_v int
    
    for ; path_list.Len() > 0 ;  {
        //Take the first element of the queue
        path = path_list.Front().Value.(s.Path_t)
        path_list.Remove(path_list.Front())
        last_v = path.Node_ids[path.Length]

        //Check if it is the goal
        if last_v == goal {
            //Add the node names
            for i:=range(path.Node_ids) {
                path.Node_names[i] = v[path.Node_ids[i]].Name
            }
            a.Solution = path
            return path, nil
        }

        expanded_nodes++
        //If it is not, generate its succesors
        path_list.PushBackList(a.Generate_succesors(adj_matrix, path, last_v, v))

    }
    
    return s.Path_t{}, errors.New(fmt.Sprintf("Could not reach %s from %s.", v[goal].Name, v[start].Name))    
}
