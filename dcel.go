package main

type vertex struct {
	x, y int
}

type halfEdge struct {
	originVertex *vertex
	twinEdge     *halfEdge
	nextEdge     *halfEdge
}

type doublyConnectedEdgeList struct {
	vertices []vertex
	edges    []halfEdge
}

func (dcel *doublyConnectedEdgeList) addIsolatedVertex(x, y int) {
	dcel.vertices = append(dcel.vertices, vertex{x: x, y: y})
}

func (dcel *doublyConnectedEdgeList) addIsolatedEdge() {
	initialHalfEdge := halfEdge{originVertex: nil, twinEdge: nil, nextEdge: nil}
	twinHalfEdge := halfEdge{originVertex: nil, twinEdge: &initialHalfEdge, nextEdge: nil}
	initialHalfEdge.twinEdge = &twinHalfEdge
	dcel.edges = append(dcel.edges, initialHalfEdge)
	dcel.edges = append(dcel.edges, twinHalfEdge)
}
