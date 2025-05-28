// Package ml implements basic machine learning algorithms
package ml

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

// LinearRegression implements simple linear regression
type LinearRegression struct {
	Slope     float64
	Intercept float64
	Trained   bool
}

// NewLinearRegression creates a new linear regression model
func NewLinearRegression() *LinearRegression {
	return &LinearRegression{
		Slope:     0,
		Intercept: 0,
		Trained:   false,
	}
}

// Fit trains the linear regression model
// Time Complexity: O(n), Space Complexity: O(1)
func (lr *LinearRegression) Fit(X [][]float64, y []float64) error {
	if len(X) != len(y) {
		return errors.New("X and y must have the same length")
	}
	
	if len(X) == 0 {
		return errors.New("training data cannot be empty")
	}
	
	// Convert to 1D if needed (simple linear regression)
	x := make([]float64, len(X))
	for i, row := range X {
		if len(row) > 0 {
			x[i] = row[0]
		}
	}
	
	n := float64(len(x))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumXX := 0.0
	
	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumXX += x[i] * x[i]
	}
	
	// Calculate slope and intercept using least squares
	denominator := n*sumXX - sumX*sumX
	if math.Abs(denominator) < 1e-10 {
		return errors.New("cannot fit model: denominator too small")
	}
	
	lr.Slope = (n*sumXY - sumX*sumY) / denominator
	lr.Intercept = (sumY - lr.Slope*sumX) / n
	lr.Trained = true
	
	return nil
}

// Predict makes predictions on new data
func (lr *LinearRegression) Predict(X [][]float64) ([]float64, error) {
	if !lr.Trained {
		return nil, errors.New("model must be trained before making predictions")
	}
	
	predictions := make([]float64, len(X))
	for i, row := range X {
		x := 0.0
		if len(row) > 0 {
			x = row[0]
		}
		predictions[i] = lr.Slope*x + lr.Intercept
	}
	
	return predictions, nil
}

// Score calculates R-squared score
func (lr *LinearRegression) Score(X [][]float64, y []float64) (float64, error) {
	predictions, err := lr.Predict(X)
	if err != nil {
		return 0, err
	}
	
	// Calculate mean of y
	yMean := 0.0
	for _, val := range y {
		yMean += val
	}
	yMean /= float64(len(y))
	
	// Calculate total sum of squares and residual sum of squares
	totalSumSquares := 0.0
	residualSumSquares := 0.0
	
	for i := 0; i < len(y); i++ {
		totalSumSquares += math.Pow(y[i]-yMean, 2)
		residualSumSquares += math.Pow(y[i]-predictions[i], 2)
	}
	
	if totalSumSquares == 0 {
		return 1.0, nil // Perfect fit when all y values are the same
	}
	
	return 1 - (residualSumSquares / totalSumSquares), nil
}

// KMeans implements K-means clustering algorithm
type KMeans struct {
	K             int
	MaxIterations int
	Tolerance     float64
	Centroids     [][]float64
	Labels        []int
	Trained       bool
}

// NewKMeans creates a new K-means model
func NewKMeans(k, maxIterations int) *KMeans {
	return &KMeans{
		K:             k,
		MaxIterations: maxIterations,
		Tolerance:     1e-4,
		Centroids:     nil,
		Labels:        nil,
		Trained:       false,
	}
}

// Distance calculates Euclidean distance between two points
func (km *KMeans) Distance(point1, point2 []float64) float64 {
	if len(point1) != len(point2) {
		return math.Inf(1)
	}
	
	sum := 0.0
	for i := 0; i < len(point1); i++ {
		diff := point1[i] - point2[i]
		sum += diff * diff
	}
	
	return math.Sqrt(sum)
}

// InitializeCentroids randomly initializes centroids
func (km *KMeans) InitializeCentroids(X [][]float64) {
	if len(X) == 0 {
		return
	}
	
	features := len(X[0])
	km.Centroids = make([][]float64, km.K)
	
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < km.K; i++ {
		km.Centroids[i] = make([]float64, features)
		for j := 0; j < features; j++ {
			// Find min and max for this feature
			min := X[0][j]
			max := X[0][j]
			for _, point := range X {
				if point[j] < min {
					min = point[j]
				}
				if point[j] > max {
					max = point[j]
				}
			}
			km.Centroids[i][j] = min + rand.Float64()*(max-min)
		}
	}
}

// AssignClusters assigns points to nearest centroid
func (km *KMeans) AssignClusters(X [][]float64) []int {
	labels := make([]int, len(X))
	
	for i, point := range X {
		minDistance := math.Inf(1)
		cluster := 0
		
		for j, centroid := range km.Centroids {
			dist := km.Distance(point, centroid)
			if dist < minDistance {
				minDistance = dist
				cluster = j
			}
		}
		
		labels[i] = cluster
	}
	
	return labels
}

// UpdateCentroids updates centroids based on cluster assignments
func (km *KMeans) UpdateCentroids(X [][]float64, labels []int) {
	features := len(X[0])
	newCentroids := make([][]float64, km.K)
	
	for i := 0; i < km.K; i++ {
		newCentroids[i] = make([]float64, features)
		count := 0
		
		// Sum all points in this cluster
		for j, point := range X {
			if labels[j] == i {
				for k := 0; k < features; k++ {
					newCentroids[i][k] += point[k]
				}
				count++
			}
		}
		
		// Calculate average (new centroid)
		if count > 0 {
			for k := 0; k < features; k++ {
				newCentroids[i][k] /= float64(count)
			}
		} else {
			// Keep old centroid if no points assigned
			copy(newCentroids[i], km.Centroids[i])
		}
	}
	
	km.Centroids = newCentroids
}

// HasConverged checks if centroids have converged
func (km *KMeans) HasConverged(oldCentroids [][]float64) bool {
	for i := 0; i < km.K; i++ {
		dist := km.Distance(km.Centroids[i], oldCentroids[i])
		if dist > km.Tolerance {
			return false
		}
	}
	return true
}

// Fit trains the K-means model
func (km *KMeans) Fit(X [][]float64) ([]int, error) {
	if len(X) < km.K {
		return nil, errors.New("number of samples must be greater than k")
	}
	
	km.InitializeCentroids(X)
	
	for iteration := 0; iteration < km.MaxIterations; iteration++ {
		// Save old centroids for convergence check
		oldCentroids := make([][]float64, km.K)
		for i := range km.Centroids {
			oldCentroids[i] = make([]float64, len(km.Centroids[i]))
			copy(oldCentroids[i], km.Centroids[i])
		}
		
		km.Labels = km.AssignClusters(X)
		km.UpdateCentroids(X, km.Labels)
		
		if km.HasConverged(oldCentroids) {
			break
		}
	}
	
	km.Trained = true
	return km.Labels, nil
}

// Predict assigns new data points to clusters
func (km *KMeans) Predict(X [][]float64) ([]int, error) {
	if !km.Trained {
		return nil, errors.New("model must be trained before making predictions")
	}
	
	return km.AssignClusters(X), nil
}

// Inertia calculates within-cluster sum of squares (WCSS)
func (km *KMeans) Inertia(X [][]float64) (float64, error) {
	if !km.Trained {
		return 0, errors.New("model must be trained before calculating inertia")
	}
	
	wcss := 0.0
	for i, point := range X {
		cluster := km.Labels[i]
		dist := km.Distance(point, km.Centroids[cluster])
		wcss += dist * dist
	}
	
	return wcss, nil
}

// KNearestNeighbors implements K-Nearest Neighbors algorithm
type KNearestNeighbors struct {
	K       int
	XTrain  [][]float64
	YTrain  []float64
	Trained bool
}

// NewKNearestNeighbors creates a new KNN model
func NewKNearestNeighbors(k int) *KNearestNeighbors {
	return &KNearestNeighbors{
		K:       k,
		XTrain:  nil,
		YTrain:  nil,
		Trained: false,
	}
}

// Distance calculates Euclidean distance between two points
func (knn *KNearestNeighbors) Distance(point1, point2 []float64) float64 {
	if len(point1) != len(point2) {
		return math.Inf(1)
	}
	
	sum := 0.0
	for i := 0; i < len(point1); i++ {
		diff := point1[i] - point2[i]
		sum += diff * diff
	}
	
	return math.Sqrt(sum)
}

// Neighbor represents a neighbor with distance and label
type Neighbor struct {
	Distance float64
	Label    float64
}

// Fit trains the KNN model (just stores the data)
func (knn *KNearestNeighbors) Fit(X [][]float64, y []float64) error {
	if len(X) != len(y) {
		return errors.New("X and y must have the same length")
	}
	
	// Deep copy the training data
	knn.XTrain = make([][]float64, len(X))
	for i, row := range X {
		knn.XTrain[i] = make([]float64, len(row))
		copy(knn.XTrain[i], row)
	}
	
	knn.YTrain = make([]float64, len(y))
	copy(knn.YTrain, y)
	
	knn.Trained = true
	return nil
}

// FindNeighbors finds k nearest neighbors for a point
func (knn *KNearestNeighbors) FindNeighbors(point []float64) []Neighbor {
	neighbors := make([]Neighbor, len(knn.XTrain))
	
	for i, trainPoint := range knn.XTrain {
		neighbors[i] = Neighbor{
			Distance: knn.Distance(point, trainPoint),
			Label:    knn.YTrain[i],
		}
	}
	
	// Sort by distance
	for i := 0; i < len(neighbors)-1; i++ {
		for j := i + 1; j < len(neighbors); j++ {
			if neighbors[i].Distance > neighbors[j].Distance {
				neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
			}
		}
	}
	
	// Return k nearest neighbors
	k := knn.K
	if k > len(neighbors) {
		k = len(neighbors)
	}
	
	return neighbors[:k]
}

// Predict predicts labels for new data points
func (knn *KNearestNeighbors) Predict(X [][]float64) ([]float64, error) {
	if !knn.Trained {
		return nil, errors.New("model must be trained before making predictions")
	}
	
	predictions := make([]float64, len(X))
	
	for i, point := range X {
		neighbors := knn.FindNeighbors(point)
		
		// For regression: average of neighbor labels
		sum := 0.0
		for _, neighbor := range neighbors {
			sum += neighbor.Label
		}
		predictions[i] = sum / float64(len(neighbors))
	}
	
	return predictions, nil
}

// Score calculates accuracy score for classification
func (knn *KNearestNeighbors) Score(X [][]float64, y []float64) (float64, error) {
	predictions, err := knn.Predict(X)
	if err != nil {
		return 0, err
	}
	
	correct := 0
	for i := 0; i < len(y); i++ {
		// For classification, round predictions
		predicted := math.Round(predictions[i])
		actual := math.Round(y[i])
		if predicted == actual {
			correct++
		}
	}
	
	return float64(correct) / float64(len(y)), nil
}
