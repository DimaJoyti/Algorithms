package ml

import (
	"math"
	"testing"
)

func TestLinearRegression(t *testing.T) {
	lr := NewLinearRegression()

	// Test initial state
	if lr.Trained {
		t.Error("New model should not be trained")
	}

	// Test fit with simple linear data
	X := [][]float64{{1}, {2}, {3}, {4}, {5}}
	y := []float64{2, 4, 6, 8, 10}

	err := lr.Fit(X, y)
	if err != nil {
		t.Errorf("Fit should not return error: %v", err)
	}

	if !lr.Trained {
		t.Error("Model should be trained after fit")
	}

	// Check slope and intercept (y = 2x)
	if math.Abs(lr.Slope-2.0) > 1e-10 {
		t.Errorf("Slope should be 2.0, got %f", lr.Slope)
	}

	if math.Abs(lr.Intercept-0.0) > 1e-10 {
		t.Errorf("Intercept should be 0.0, got %f", lr.Intercept)
	}

	// Test predictions
	testX := [][]float64{{6}, {7}}
	predictions, err := lr.Predict(testX)
	if err != nil {
		t.Errorf("Predict should not return error: %v", err)
	}

	expected := []float64{12, 14}
	for i, pred := range predictions {
		if math.Abs(pred-expected[i]) > 1e-10 {
			t.Errorf("Prediction %d should be %f, got %f", i, expected[i], pred)
		}
	}

	// Test score
	score, err := lr.Score(X, y)
	if err != nil {
		t.Errorf("Score should not return error: %v", err)
	}

	if math.Abs(score-1.0) > 1e-10 {
		t.Errorf("Score should be 1.0 (perfect fit), got %f", score)
	}
}

func TestLinearRegressionErrors(t *testing.T) {
	lr := NewLinearRegression()

	// Test mismatched lengths
	X := [][]float64{{1}, {2}}
	y := []float64{1}
	err := lr.Fit(X, y)
	if err == nil {
		t.Error("Fit should return error for mismatched lengths")
	}

	// Test empty data
	err = lr.Fit([][]float64{}, []float64{})
	if err == nil {
		t.Error("Fit should return error for empty data")
	}

	// Test predict before training
	_, err = lr.Predict([][]float64{{1}})
	if err == nil {
		t.Error("Predict should return error before training")
	}
}

func TestKMeans(t *testing.T) {
	kmeans := NewKMeans(2, 100)

	// Test initial state
	if kmeans.Trained {
		t.Error("New model should not be trained")
	}

	// Test with simple 2D data (two clear clusters)
	X := [][]float64{
		{1, 1}, {1, 2}, {2, 1}, {2, 2}, // Cluster 1
		{8, 8}, {8, 9}, {9, 8}, {9, 9}, // Cluster 2
	}

	labels, err := kmeans.Fit(X)
	if err != nil {
		t.Errorf("Fit should not return error: %v", err)
	}

	if !kmeans.Trained {
		t.Error("Model should be trained after fit")
	}

	if len(labels) != len(X) {
		t.Errorf("Labels length should be %d, got %d", len(X), len(labels))
	}

	// Test that points in same cluster have same label
	cluster1Label := labels[0]
	for i := 0; i < 4; i++ {
		if labels[i] != cluster1Label {
			t.Errorf("Points 0-3 should have same cluster label")
		}
	}

	cluster2Label := labels[4]
	for i := 4; i < 8; i++ {
		if labels[i] != cluster2Label {
			t.Errorf("Points 4-7 should have same cluster label")
		}
	}

	// Test predictions
	testX := [][]float64{{1.5, 1.5}, {8.5, 8.5}}
	predictions, err := kmeans.Predict(testX)
	if err != nil {
		t.Errorf("Predict should not return error: %v", err)
	}

	if len(predictions) != 2 {
		t.Errorf("Predictions length should be 2, got %d", len(predictions))
	}

	// Test inertia
	inertia, err := kmeans.Inertia(X)
	if err != nil {
		t.Errorf("Inertia should not return error: %v", err)
	}

	if inertia < 0 {
		t.Errorf("Inertia should be non-negative, got %f", inertia)
	}
}

func TestKMeansErrors(t *testing.T) {
	kmeans := NewKMeans(5, 100)

	// Test with too few samples
	X := [][]float64{{1, 1}, {2, 2}} // Only 2 samples for k=5
	_, err := kmeans.Fit(X)
	if err == nil {
		t.Error("Fit should return error when k > number of samples")
	}

	// Test predict before training
	_, err = kmeans.Predict([][]float64{{1, 1}})
	if err == nil {
		t.Error("Predict should return error before training")
	}

	// Test inertia before training
	_, err = kmeans.Inertia([][]float64{{1, 1}})
	if err == nil {
		t.Error("Inertia should return error before training")
	}
}

func TestKNearestNeighbors(t *testing.T) {
	knn := NewKNearestNeighbors(3)

	// Test initial state
	if knn.Trained {
		t.Error("New model should not be trained")
	}

	// Test with simple classification data
	X := [][]float64{
		{1, 1}, {1, 2}, {2, 1}, // Class 0
		{5, 5}, {5, 6}, {6, 5}, // Class 1
	}
	y := []float64{0, 0, 0, 1, 1, 1}

	err := knn.Fit(X, y)
	if err != nil {
		t.Errorf("Fit should not return error: %v", err)
	}

	if !knn.Trained {
		t.Error("Model should be trained after fit")
	}

	// Test predictions
	testX := [][]float64{{1.5, 1.5}, {5.5, 5.5}}
	predictions, err := knn.Predict(testX)
	if err != nil {
		t.Errorf("Predict should not return error: %v", err)
	}

	if len(predictions) != 2 {
		t.Errorf("Predictions length should be 2, got %d", len(predictions))
	}

	// First point should be closer to class 0
	if math.Round(predictions[0]) != 0 {
		t.Errorf("First prediction should be close to 0, got %f", predictions[0])
	}

	// Second point should be closer to class 1
	if math.Round(predictions[1]) != 1 {
		t.Errorf("Second prediction should be close to 1, got %f", predictions[1])
	}

	// Test score
	score, err := knn.Score(X, y)
	if err != nil {
		t.Errorf("Score should not return error: %v", err)
	}

	if score != 1.0 {
		t.Errorf("Score should be 1.0 (perfect accuracy on training data), got %f", score)
	}

	// Test find neighbors
	neighbors := knn.FindNeighbors([]float64{1, 1})
	if len(neighbors) != 3 {
		t.Errorf("Should find 3 neighbors, got %d", len(neighbors))
	}

	// Neighbors should be sorted by distance
	for i := 1; i < len(neighbors); i++ {
		if neighbors[i].Distance < neighbors[i-1].Distance {
			t.Error("Neighbors should be sorted by distance")
		}
	}
}

func TestKNearestNeighborsErrors(t *testing.T) {
	knn := NewKNearestNeighbors(3)

	// Test mismatched lengths
	X := [][]float64{{1, 1}, {2, 2}}
	y := []float64{0}
	err := knn.Fit(X, y)
	if err == nil {
		t.Error("Fit should return error for mismatched lengths")
	}

	// Test predict before training
	_, err = knn.Predict([][]float64{{1, 1}})
	if err == nil {
		t.Error("Predict should return error before training")
	}
}

func TestDistance(t *testing.T) {
	knn := NewKNearestNeighbors(1)

	// Test Euclidean distance
	point1 := []float64{0, 0}
	point2 := []float64{3, 4}
	distance := knn.Distance(point1, point2)

	expected := 5.0 // 3-4-5 triangle
	if math.Abs(distance-expected) > 1e-10 {
		t.Errorf("Distance should be %f, got %f", expected, distance)
	}

	// Test distance with different dimensions
	point3 := []float64{0, 0, 0}
	distance = knn.Distance(point1, point3)
	if !math.IsInf(distance, 1) {
		t.Error("Distance between points of different dimensions should be infinity")
	}
}

func TestKMeansDistance(t *testing.T) {
	kmeans := NewKMeans(2, 100)

	// Test Euclidean distance
	point1 := []float64{0, 0}
	point2 := []float64{3, 4}
	distance := kmeans.Distance(point1, point2)

	expected := 5.0 // 3-4-5 triangle
	if math.Abs(distance-expected) > 1e-10 {
		t.Errorf("Distance should be %f, got %f", expected, distance)
	}
}

// Benchmark tests
func BenchmarkLinearRegressionFit(b *testing.B) {
	X := make([][]float64, 1000)
	y := make([]float64, 1000)
	for i := 0; i < 1000; i++ {
		X[i] = []float64{float64(i)}
		y[i] = float64(i) * 2
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lr := NewLinearRegression()
		lr.Fit(X, y)
	}
}

func BenchmarkKMeansFit(b *testing.B) {
	X := make([][]float64, 100)
	for i := 0; i < 100; i++ {
		X[i] = []float64{float64(i % 10), float64(i / 10)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kmeans := NewKMeans(3, 10)
		kmeans.Fit(X)
	}
}

func BenchmarkKNNPredict(b *testing.B) {
	X := make([][]float64, 100)
	y := make([]float64, 100)
	for i := 0; i < 100; i++ {
		X[i] = []float64{float64(i), float64(i)}
		y[i] = float64(i % 2)
	}

	knn := NewKNearestNeighbors(5)
	knn.Fit(X, y)

	testX := [][]float64{{50, 50}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		knn.Predict(testX)
	}
}
