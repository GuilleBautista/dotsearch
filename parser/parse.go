package parser

import (
    "fmt"
    "regexp"
    "io/ioutil"
    "errors"
    "container/list"
    "strconv"
    s "github.com/guillebautista/dotsearch/structures"
)


func find(elem int, list []int) bool {
    for i := range list {
        if elem==list[i] { return true }
    }
    return false
}

type Parser interface {
    Parse() s.Graph_t
}

//Struct graph parser implements the parser interface
type GraphParser struct{
    File string
    //Allow transitions that start and end in the same node
    Allow_recursive bool
    Verbose bool
}

//This function runs when the parser detects a node
//it adds the correspondent node to the input graph
func add_node(new_node s.Node_t, graph map[s.Node_t]map[s.Node_t]int, from_newnode map[s.Node_t]int, directed bool, allow_recursive bool, verbose bool) (error){
    node_exists := false

    for n := range graph {
        if n.Name == new_node.Name {
            if new_node.H >= 0 && n.H != new_node.H {
                //Both the existing node and the new one have non negative heuristic value
                // And they differ
                return errors.New(fmt.Sprintf("Error: node %s heuristic specified multiple times.\n", new_node.Name))
            }
            //The new node has a negative heuristic value so assign the existing one to it
            new_node.H = n.H
            node_exists = true
            break
        }
    }

    //Check if recursive transitions are being added
    if !allow_recursive {
        for key, _ := range from_newnode {
            if key.Name==new_node.Name {
                from_newnode[key]=-1
                if verbose {
                    fmt.Sprintf("Warning: Removing recursive transition from %s...\n", key.Name)
                }
            }
        }
    }

    //Initialize node if it does not exist
    if !node_exists {
        graph[new_node]=from_newnode
    }

    //Add the destinations to the graph
    for key, value := range from_newnode{
        //If the node existed check that only new transitions are added
        if node_exists && graph[new_node][key]!=0 && graph[new_node][key] != value {
            return errors.New(fmt.Sprintf("duplicated edge %s->%s with different weights (%d, %d)", new_node, key, graph[new_node][key], value))
        }else{
            //If the destination node does not exist, add it
            if graph[key] == nil {
                err := add_node(key, graph, make(map[s.Node_t]int), true, allow_recursive, verbose)
                if err != nil{
                    panic(err)
                }
            }
            //If there is no error, add the transition
            if node_exists {
                //If the node did not exist the transition was previously added
                graph[new_node][key]=value
            }
        }
    }
    if !directed {
        for node, weight := range from_newnode {
            add_node(node, graph, map[s.Node_t]int{new_node: weight}, true, allow_recursive, verbose)
        }
    }
    return nil
}

//Function to get a map with the regexp groups assigned to their values (in str format)
func parse_re(re *regexp.Regexp, content string) ( map[string]string ) {
    match := re.FindStringSubmatch(string(content))
    paramsMap := make(map[string]string)
    for i, name := range re.SubexpNames() {
        if i > 0 && i <= len(match) {
            paramsMap[name] = match[i]
        }
    }
    return paramsMap
}

/*
Parse the origin node of a node declaration.
    Check if the node is the start or goal node.
        If it is, set the start node or goal node global variables to the node's name
    If it is and there was previously a goal node or start node and it was different from the current one, return error
Also parses multiple origin nodes
TODO: create proper errors
*/
func get_node_from_str(parse_node_re *regexp.Regexp, value string, parse_heuristic *regexp.Regexp, 
    parse_start *regexp.Regexp, parse_goal *regexp.Regexp, node_start *string, node_goal *string) (s.Node_t){
    node_data := parse_re(parse_node_re, value)

    //For each node in the node_origin class parse its name and attributes
    //Create an empty node
    node_origin := s.Node_t{
        Name: "",
        H: -1,
    }
    
    if len(node_data["node_name"])>0 {
        node_origin.Name=node_data["node_name"]
    } else {
        panic(errors.New(fmt.Sprintf("Node declaration error")))
        return node_origin
    }
    if len(node_data["node_attrs"])>0 {
        //With parse attribute get every possible attribute type
        h_data := parse_re(parse_heuristic, node_data["node_attrs"])
        start_data := parse_re(parse_start, node_data["node_attrs"])
        goal_data := parse_re(parse_goal, node_data["node_attrs"])
        
        //Not every node has to have all the types, and the types are not exclusive
        if len(start_data["start"])>0 {
            if len(*node_start)>0 && *node_start != node_origin.Name {
                panic(errors.New(fmt.Sprintf("Error: multiple start point specified: %s and %s.\n", *node_start, node_origin.Name)))
                return node_origin
            }
            *node_start = node_origin.Name
        }
        if len(goal_data["goal"])>0{
            if len(*node_goal)>0 && *node_goal != node_origin.Name {
                panic(errors.New(fmt.Sprintf("Error: multiple goals specified: %s and %s.\n", *node_goal, node_origin.Name)))
            }
            *node_goal = node_origin.Name
        }
        if len(h_data["heuristic"])>0{

            v, err := strconv.ParseInt(h_data["heuristic"], 10, 0)
            if err!=nil {
                panic(err)
            }
            node_origin.H = int(v)
        }
    }
    return node_origin
}

//This needs to exist to simplify changing graph's main data structure
func map_to_matrix(graph map[s.Node_t]map[s.Node_t]int, start string, goal string, directed bool) (s.Graph_t) {
    n := len(graph)

    v_list := make([]s.Node_t, n)
    v_dict := make(map[string]int)
    i:=0
    for k, _ := range graph{
        v_list[i]=k
        v_dict[k.Name]=i
        i++
    }

    adj_matrix := make([][]int, n)

    //The nodes must be added in order to avoid errors while translating
    //This is because how go iterates a dictionary
    for i:=0; i<len(v_dict); i++ {
        connection := graph[v_list[i]]
        adj_matrix[i] = make([]int, n)
        var exists_path = make([]bool, n)

        //Here order does not matter
        for destination, weight := range connection {
            var dest_n = v_dict[destination.Name]
            adj_matrix[i][dest_n] = weight
            exists_path[dest_n]=true
        }
        //Remove zeros
        for v := range adj_matrix[i]{
            if !exists_path[v] {
                adj_matrix[i][v]=-1
            }
        }
    }

    return s.Graph_t{
        V: v_list, 
        Matrix: adj_matrix,
        Start: v_dict[start],
        Goal: v_dict[goal],
        Directed: directed,
    } 
}

func (p GraphParser) Parse() (s.Graph_t) {
    node_start := ""
    node_goal := ""
    
    content, err := ioutil.ReadFile(p.File)

    if err != nil {
        panic(err)
    }

    graph := make(map[s.Node_t]map[s.Node_t]int)

    //Grammar:
    tab := regexp.MustCompile(`(\t| )+`)
    //\n \n\n \n
    enter := regexp.MustCompile(`(\n *)+`)
    //Hello
    str_re := `([A-Za-z0-9]+(_*\.*-*)*)+`
    
    parse_heuristic := fmt.Sprintf(`((H|h)(euristic)? ?= ?)?(\d+\.\d+|\.\d+|\d+)`)
    parse_heuristic_re := regexp.MustCompile(`((H|h)(euristic)? ?= ?)?(?P<heuristic>(\d+\.\d+|\.\d+|\d+))`)
    parse_start := fmt.Sprintf(`(S|s)tart`)
    parse_start_re := regexp.MustCompile(`(?P<start>(S|s)tart)`)
    parse_goal := fmt.Sprintf(`(G|g)oal`)
    parse_goal_re := regexp.MustCompile(`(?P<goal>(G|g)oal)`)
    //->[Weight=value]
    //TODO: -[Weight=value]
    parse_diedge_re := regexp.MustCompile(fmt.Sprintf(`-> ?(\[ ?((W|w)eight ?= ?)?(?P<weight>(\d+\.\d+|\.\d+|\d+)) ?\])?`))
    parse_edge_re := regexp.MustCompile(fmt.Sprintf(`-- ?(\[ ?((W|w)eight ?= ?)?(?P<weight>(\d+\.\d+|\.\d+|\d+)) ?\])?`))
    parse_diedge := fmt.Sprintf(`-> ?(\[ ?((W|w)eight ?= ?)?(\d+\.\d+|\.\d+|\d+) ?\])?`)
    parse_edge := fmt.Sprintf(`-- ?(\[ ?((W|w)eight ?= ?)?(\d+\.\d+|\.\d+|\d+) ?\])?`)
    //A[heuristic=value|Start|Goal]
    parse_attr := fmt.Sprintf(`(%s|%s|%s)`, parse_heuristic, parse_start, parse_goal)
    parse_node := fmt.Sprintf(`%s ?(\[ ?%s( ?, ?%s)* ?\])?`, str_re, parse_attr, parse_attr)
    parse_node_re := regexp.MustCompile(fmt.Sprintf(`(?P<node_name>%s) ?((?P<node_attrs>\[ ?%s( ?, ?%s)* ?\]))?`, str_re, parse_attr, parse_attr))
    parse_nodes := fmt.Sprintf(`{(%s ?(\[ ?%s( ?, ?%s)* ?\])?)( ?, ?%s ?(\[ ?%s ?(, ?%s)* ?\])?)*}`, str_re, parse_attr, parse_attr, str_re, parse_attr, parse_attr)
    //A[heuristic=value|Start|Goal] ->[weight=value] B[heuristic=value|Start|Goal]
    edge_declaration := regexp.MustCompile(fmt.Sprintf(`((?P<node_origin>%s)|(?P<nodes_origin>%s)) ?((?P<diedge>%s)|(?P<edge>%s)) ?((?P<node_dest>%s)|(?P<nodes_dest>%s)) ?;`, parse_node, parse_nodes, parse_diedge, parse_edge, parse_node, parse_nodes))
    //digraph G { A->B }
    
    content_notabs := tab.ReplaceAll(content, []byte(` `))
    
    //Used to be able to parse line by line?
    content_clean := enter.ReplaceAll(content_notabs, []byte("\n"))

    //End Grammar
    var directed bool = true

    for _, ed := range edge_declaration.FindAll(content_clean, -1){
        //Map containing all declaration attributes
        ed_attributes := parse_re(edge_declaration, string(ed))
        //Expecting node_origin class | nodes_origin class, edge class | diedge class, node_dest class
        nodes_origin := list.New()
        nodes_dest := list.New()

        var node_origin s.Node_t
        if len(ed_attributes["node_origin"]) > 0 {
            node_origin = get_node_from_str(parse_node_re, ed_attributes["node_origin"], parse_heuristic_re, parse_start_re, parse_goal_re, &node_start, &node_goal)
            nodes_origin.PushBack(node_origin)
        }else if len(ed_attributes["nodes_origin"]) > 0 {            
            for _, node := range parse_node_re.FindAll([]byte(ed_attributes["nodes_origin"]), -1){
                node_origin = get_node_from_str(parse_node_re, string(node), parse_heuristic_re, parse_start_re, parse_goal_re, &node_start, &node_goal)                
                nodes_origin.PushBack(node_origin)               
            }
        }

        var node_dest s.Node_t
        if len(ed_attributes["node_dest"]) > 0 {
            node_dest = get_node_from_str(parse_node_re, ed_attributes["node_dest"], parse_heuristic_re, parse_start_re, parse_goal_re, &node_start, &node_goal)
            nodes_dest.PushBack(node_dest)
        }else if len(ed_attributes["nodes_dest"]) > 0 {            
            for _, node := range parse_node_re.FindAll([]byte(ed_attributes["nodes_dest"]), -1){
                node_dest = get_node_from_str(parse_node_re, string(node), parse_heuristic_re, parse_start_re, parse_goal_re, &node_start, &node_goal)                
                nodes_dest.PushBack(node_dest)               
            }
        }
        //node_dest := get_node_from_str(parse_node_re, ed_attributes["node_dest"], parse_attribute, &node_start, &node_goal)
        //nodes_dest.PushBack(node_dest)

        var weight int64 = -1
        var err error = nil
        //Check if the edge is directed or simple
        if len(ed_attributes["edge"]) > 0 {
            weight, err = strconv.ParseInt(parse_re(parse_edge_re, ed_attributes["edge"])["weight"], 10, 0)
            directed = false
        }else if len(ed_attributes["diedge"]) > 0 {
            weight, err = strconv.ParseInt(parse_re(parse_diedge_re, ed_attributes["diedge"])["weight"], 10, 0)
        }
        if err != nil {
            if err.Error() == "strconv.ParseInt: parsing \"\": invalid syntax" {
                weight = 1
            }else{
                err_str := fmt.Sprintf("%s\nError parsing edge %s\n", err.Error(), ed)
                panic(errors.New(err_str))
            }
        }
        if weight < 0 {
            err_str := fmt.Sprintf("Error: negative weights are not allowed...\nError parsing edge %s\n", ed)
            panic(errors.New(err_str))
        
        }
        
        for o := nodes_origin.Front(); o != nil; o = o.Next() {
            node_origin := o.Value.(s.Node_t)

            from_newnode := make(map[s.Node_t]int)
            for d := nodes_dest.Front(); d != nil; d = d.Next() {
                node_dest := d.Value.(s.Node_t)
                from_newnode[node_dest] = int(weight)

            }
            //Add node origin and its destinations
            err = add_node(node_origin, graph, from_newnode, directed, p.Allow_recursive, p.Verbose)

            if err!=nil {
                panic(err)
            }
        }
    }
    return map_to_matrix(graph, node_start, node_goal, directed)
}