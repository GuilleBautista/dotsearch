package algorithms

import (
	"fmt"
	"container/list"
    s "github.com/guillebautista/dotsearch/structures"
//    "github.com/guillebautista/dotsearch/viewer"

    "errors"
//	"time"
)

type Astar struct {
	Graph s.Graph_t
	Solution s.Path_t
}

//Simple function to generate all succesors of a node without any further checks
func (a Astar) generate_succesors(adj_matrix [][]int, path s.Path_t, last_v int, v []s.Node_t, closed []bool, open []int) []s.Path_t {
    var succesors []s.Path_t
    var newpath s.Path_t
	
	for i := range adj_matrix {
        //If connected and not already visited and the cost sum for the node is less than the previous one
        if !closed[i] && adj_matrix[last_v][i] >= 0 && (open[i] < 0 || adj_matrix[last_v][i]+path.Cost <= open[i]) {
			newpath.Cost = adj_matrix[last_v][i] + path.Cost
            newpath.Length = path.Length+1
			newpath.Node_names = make([]string, newpath.Length+1)
			newpath.Node_ids = make([]int, newpath.Length+1)

            copy(newpath.Node_ids, path.Node_ids)
            newpath.Node_ids[newpath.Length] = i
			
			//Insert new path into succesors
			//empty slice case
			if len(succesors)==0 {
				succesors = []s.Path_t{newpath}
				continue
			}
			low := 0
			high := len(succesors) - 1
			var tmp []s.Path_t

			//for loop to insert orderly
			for ;true; {
				if succesors[low].Cost + v[succesors[low].Node_ids[succesors[low].Length]].H >= newpath.Cost + v[i].H {
					succesors = append(append(append(tmp, succesors[:low]...), newpath), succesors[low:]...)
					break
				}else if succesors[high].Cost + v[succesors[high].Node_ids[succesors[high].Length]].H <= newpath.Cost + v[i].H {
					succesors = append(append(append(tmp, succesors[:high+1]...), newpath), succesors[high+1:]...)
					break
				}else if high-low  <= 1 {
					succesors = append(append(append(tmp, succesors[:low+1]...), newpath), succesors[low+1:]...)
					break
				}
				mid := int((high-low)/2)+low
				if succesors[mid].Cost + v[succesors[mid].Node_ids[succesors[mid].Length]].H > newpath.Cost + v[i].H {
					high = mid
				}else if succesors[mid].Cost + v[succesors[mid].Node_ids[succesors[mid].Length]].H < newpath.Cost + v[i].H {
					low = mid
				}else {
					succesors = append(append(append(tmp, succesors[:mid]...), newpath), succesors[mid:]...)
					break
				}
			}
        }
    }
    return succesors
}

//Dijkstra's algorithm
func (a Astar) Solve()(s.Path_t, error){
	//Gather the data
    start := a.Graph.Start
    goal := a.Graph.Goal
    adj_matrix := a.Graph.Matrix
    v := a.Graph.V
    path_list := list.New()
    open := make([]int, len(adj_matrix))
	for i := range open {
		open[i]=-1
	}
	closed := make([]bool, len(adj_matrix))

    path_list.PushBack( s.Path_t{
        Cost: 0,
        Length: 0,
        Node_names: []string{},
        Node_ids: []int{start},
    } )

    var path s.Path_t
    var last_v int
    
    for ; path_list.Len() > 0 ;  {
        //Take the first element of the open list
        path = path_list.Front().Value.(s.Path_t)

		//time.Sleep(1 * time.Second)
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

		closed[last_v]=true
        //If it is not, generate its succesors
        succesors := a.generate_succesors(adj_matrix, path, last_v, v, closed, open)

		i := 0
		marked := make([]bool, len(adj_matrix))
		for n := path_list.Front(); n!=nil; n=n.Next() {
			//Insert each succesor orderly
			for ;i < len(succesors) && n.Value.(s.Path_t).Cost + v[n.Value.(s.Path_t).Node_ids[n.Value.(s.Path_t).Length]].H > succesors[i].Cost + v[succesors[i].Node_ids[succesors[i].Length]].H; {
				succ_last_v := succesors[i].Node_ids[succesors[i].Length]
				
				//If it was not on the open list there is no node to delete 
				if open[succ_last_v] != -1 {
					//Mark for delete the previous path that ended in succ_last_v
					marked[succ_last_v]=true
				}
				//Update the open list
				open[succ_last_v] = succesors[i].Cost + v[succ_last_v].H
				path_list.InsertBefore(succesors[i], n)
				i++
			}
			//Delete the list element if it was marked
			n_last_v := n.Value.(s.Path_t).Node_ids[n.Value.(s.Path_t).Length]
			if marked[n_last_v] {
				path_list.Remove(n)
				marked[n_last_v] = false

			}
		}
		//If the end of the list has been reached there is no node to be marked
		// because the succesors must be all new
		for ;i < len(succesors); {
			path_list.PushBack(succesors[i])
			succ_last_v := succesors[i].Node_ids[succesors[i].Length]
			open[succ_last_v] = succesors[i].Cost
			i++
		}
	}
    
    return s.Path_t{}, errors.New(fmt.Sprintf("Could not reach %s from %s.", v[goal].Name, v[start].Name))
    
}
