package pathfinder

type NodeList struct {
	Source      string
	Destination string
	Cost        int
}

func GetPath(nodes []NodeList, pathSource string, pathDestination string) []string {

	graph := newGraph()

	//Considerare l'inserimento parallelo dei nodi nel caso di una lista con molti nodi
	for _, node := range nodes {
		graph.addEdge(node.Source, node.Destination, node.Cost)
	}

	_, path := graph.getPath(pathSource, pathDestination)

	//path puo essere nil
	return path

}
