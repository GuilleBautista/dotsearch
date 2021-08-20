package algorithms

import (
    "fmt"
    "container/list"
    s "github.com/guillebautista/dotsearch/structures"
)

type IDAstar struct {
    Graph s.Graph_t
    Solution s.Path_t
}

func (a IDAstar) generate_succesors(adj_matrix [][]int, path s.Path_t, last_v int, v []s.Node_t) *list.List {
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

//IDA star search. Set max depth to -1 for no limit.
func (a IDAstar) Solve(max_depth int)(s.Path_t, error){
    //Gather the data
    start := a.Graph.Start
    goal := a.Graph.Goal
    adj_matrix := a.Graph.Matrix
    v := a.Graph.V
    path_list := list.New()

    var max_cost int
    if max_depth < 0 {
        max_cost = 0
    }else {
        max_cost = max_depth
    }

    var path s.Path_t
    var last_v int
    
    //Main loop
    for ; true ; {
        min_newcost := 2147483647
        greater_cost_found := false

        path_list.PushBack( s.Path_t{
            Cost: 0,
            Length: 0,
            Node_names: []string{},
            Node_ids: []int{start},
        } )
    
        //Inside loop
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
                return path, nil
            }
            if path.Cost + v[last_v].H <= max_cost {
                //If it is in range, generate its succesors
                path_list.PushFrontList(a.generate_succesors(adj_matrix, path, last_v, v))
            }else if min_newcost > path.Cost + v[last_v].H{
                //Store the minimum cost encountered in the succesor expansion
                min_newcost = path.Cost + v[last_v].H
                greater_cost_found = true
            }
        }
        //This means the whole search space has been checked
        if !greater_cost_found {
            break
        }
        //After the main loop update the max_cost
        max_cost = min_newcost
    }
    fmt.Printf("Could not reach %s from %s.", v[goal].Name, v[start].Name)
    return s.Path_t{}, nil    
}
