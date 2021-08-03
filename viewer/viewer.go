package viewer

import (
    "fmt"
    "io/ioutil"
    "os/exec"
    "os"
    s "github.com/guillebautista/dotsearch/structures"
)

//There will be graph and path viewers
type Viewer interface {
    Print()
    View_png(dest_file string) error
}

type PathViewer struct {
    Path s.Path_t
}

type GraphViewer struct {
    Graph s.Graph_t
}


//Function to print a path into stdout
func (pv PathViewer) Print() {
    fmt.Printf("Cost: %f\nLength: %d\nNodes: %v\nIDs: %v \n", pv.Path.Cost, pv.Path.Length, pv.Path.Node_names, pv.Path.Node_ids)
}
/*
func Print_graph_map(graph map[s.Node_t]map[s.Node_t]int){
    for node, connections := range graph {
        fmt.Printf("%s %f: ", node.Name, node.H)
        for connection, weight := range connections{
            fmt.Printf("%s:%f ", connection.Name, weight)
        }

        fmt.Printf("\n")
    }
}
*/
func (gv GraphViewer) Print(){
    for v := range gv.Graph.V {
        if gv.Graph.V[v].H >= 0 {
            fmt.Printf("  %s h=%f", gv.Graph.V[v].Name, gv.Graph.V[v].H)
        }else{
            fmt.Printf("  %s", gv.Graph.V[v].Name)
        }
    }
    fmt.Printf("\n")
    for i := range gv.Graph.Matrix {
        fmt.Printf("%s: ", gv.Graph.V[i].Name)
        for j := range gv.Graph.Matrix[i]{
            if gv.Graph.Matrix[i][j] >= 0{
                fmt.Printf("%f ", gv.Graph.Matrix[i][j])
            }else{
                fmt.Printf("\t")
            }
        }
        fmt.Printf("\n")
    }
}

func (gv GraphViewer) View_png(dest_file string) error {
    if dest_file[4:]!=".png" {
        dest_file+=".png"
    }
    
    var result string
    if gv.Graph.Directed{
        result = "digraph G {\n"
    }else{
        result = "graph G{\n"
    }

    for v:=range gv.Graph.V {
        if gv.Graph.V[v].H >= 0 {
            result+=fmt.Sprintf("%s [label=\"%s, h=%f\"]\n", gv.Graph.V[v].Name, gv.Graph.V[v].Name, gv.Graph.V[v].H)
        }
    }

    for i := range gv.Graph.Matrix {
        if !gv.Graph.Directed{
            for j:=0 ; j<=i; j++ {
                if gv.Graph.Matrix[i][j] >= 0 {
                    result+=fmt.Sprintf("%s -- %s [label=\"%f\"]\n", gv.Graph.V[i].Name, gv.Graph.V[j].Name, gv.Graph.Matrix[i][j])
                }
            }
        }else{
            for j := range gv.Graph.Matrix[i] {
                if gv.Graph.Matrix[i][j] >= 0 {
                    result+=fmt.Sprintf("%s -> %s [label=\"%f\"]\n", gv.Graph.V[i].Name, gv.Graph.V[j].Name, gv.Graph.Matrix[i][j])
                }
            }
        }
    }

    result+="}\n"
    content := []byte(result)
	tmpfile, err := ioutil.TempFile(".", "tmp_dot_file")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name())

    err = os.WriteFile(tmpfile.Name(), content, 0666)
    if err != nil {
		panic(err)
	}

    //out, err := exec.Command(fmt.Sprintf("ls")).Output()

    out, err := exec.Command("dot",  "-Tpng", tmpfile.Name()).Output()

    err = os.WriteFile(dest_file, out, 0666)

    return err
}