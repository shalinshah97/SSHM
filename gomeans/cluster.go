package gomeans

//Cluster struct composed of a center Point{X,Y float64} and of Points
type Cluster struct {
	Center Point
	Points []Point
}

func (cluster *Cluster) pointsXValues() (pointsXValues []float64) {
	for i := 0; i < len(cluster.Points); i++ {
		pointsXValues = append(pointsXValues, cluster.Points[i].Lat)
	}
	return
}

func (cluster *Cluster) pointsYValues() (pointsYValues []float64) {
	for i := 0; i < len(cluster.Points); i++ {
		pointsYValues = append(pointsYValues, cluster.Points[i].Lon)
	}
	return
}

func (cluster *Cluster) repositionCenter() {
	var x, y float64
	var clusterCount = len(cluster.Points)

	for i := 0; i < clusterCount; i++ {
		x = x + cluster.Points[i].Lat
		y = y + cluster.Points[i].Lon
	}
	cluster.Points = []Point{}
	cluster.Center = Point{x / float64(clusterCount), y / float64(clusterCount),"centroid"}
}
