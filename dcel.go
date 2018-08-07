package main

type vertex struct {
	x, y float64
}

type halfEdge struct {
	originVertex *vertex
	twinEdge     *halfEdge
	nextEdge     *halfEdge
}

type doublyConnectedEdgeList struct {
	vertices []*vertex
	edges    []*halfEdge
}

func (dcel *doublyConnectedEdgeList) addIsolatedVertex(x, y float64) *vertex {
	newVertex := vertex{x: x, y: y}
	dcel.vertices = append(dcel.vertices, &newVertex)
	return &newVertex
}

func (dcel *doublyConnectedEdgeList) addIsolatedEdge() *halfEdge {
	initialHalfEdge := halfEdge{originVertex: nil, twinEdge: nil, nextEdge: nil}
	twinHalfEdge := halfEdge{originVertex: nil, twinEdge: &initialHalfEdge, nextEdge: nil}
	initialHalfEdge.twinEdge = &twinHalfEdge
	dcel.edges = append(dcel.edges, &initialHalfEdge)
	dcel.edges = append(dcel.edges, &twinHalfEdge)
	return &initialHalfEdge
}

func getVertex() vertex {
	return vertex{}
}

func getVertexPointer() *vertex {
	return &vertex{}
}
