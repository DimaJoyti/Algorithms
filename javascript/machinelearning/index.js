// --- Directions
// Implement basic Machine Learning algorithms.
// This module contains implementations of fundamental ML algorithms
// including Linear Regression, K-Means Clustering, K-Nearest Neighbors,
// and Neural Network basics.
// --- Examples
//   const lr = new LinearRegression();
//   lr.fit([[1], [2], [3]], [2, 4, 6]);
//   lr.predict([[4]]); // ~8
// --- Complexity
//   Varies by algorithm

/**
 * Linear Regression implementation using least squares
 */
class LinearRegression {
  constructor() {
    this.slope = 0;
    this.intercept = 0;
    this.trained = false;
  }

  /**
   * Train the linear regression model
   * @param {number[][]} X - Training features (2D array)
   * @param {number[]} y - Training targets (1D array)
   */
  fit(X, y) {
    if (X.length !== y.length) {
      throw new Error('X and y must have the same length');
    }

    // Convert to 1D if needed (simple linear regression)
    const x = X.map(row => Array.isArray(row) ? row[0] : row);
    
    const n = x.length;
    const sumX = x.reduce((sum, val) => sum + val, 0);
    const sumY = y.reduce((sum, val) => sum + val, 0);
    const sumXY = x.reduce((sum, val, i) => sum + val * y[i], 0);
    const sumXX = x.reduce((sum, val) => sum + val * val, 0);

    // Calculate slope and intercept using least squares
    this.slope = (n * sumXY - sumX * sumY) / (n * sumXX - sumX * sumX);
    this.intercept = (sumY - this.slope * sumX) / n;
    this.trained = true;
  }

  /**
   * Make predictions
   * @param {number[][]} X - Features to predict
   * @returns {number[]} Predictions
   */
  predict(X) {
    if (!this.trained) {
      throw new Error('Model must be trained before making predictions');
    }

    return X.map(row => {
      const x = Array.isArray(row) ? row[0] : row;
      return this.slope * x + this.intercept;
    });
  }

  /**
   * Calculate R-squared score
   * @param {number[][]} X - Test features
   * @param {number[]} y - True values
   * @returns {number} R-squared score
   */
  score(X, y) {
    const predictions = this.predict(X);
    const yMean = y.reduce((sum, val) => sum + val, 0) / y.length;
    
    const totalSumSquares = y.reduce((sum, val) => sum + Math.pow(val - yMean, 2), 0);
    const residualSumSquares = y.reduce((sum, val, i) => 
      sum + Math.pow(val - predictions[i], 2), 0);
    
    return 1 - (residualSumSquares / totalSumSquares);
  }
}

/**
 * K-Means Clustering implementation
 */
class KMeans {
  constructor(k = 3, maxIterations = 100, tolerance = 1e-4) {
    this.k = k;
    this.maxIterations = maxIterations;
    this.tolerance = tolerance;
    this.centroids = [];
    this.labels = [];
    this.trained = false;
  }

  /**
   * Calculate Euclidean distance between two points
   * @param {number[]} point1 - First point
   * @param {number[]} point2 - Second point
   * @returns {number} Distance
   */
  distance(point1, point2) {
    return Math.sqrt(
      point1.reduce((sum, val, i) => sum + Math.pow(val - point2[i], 2), 0)
    );
  }

  /**
   * Initialize centroids randomly
   * @param {number[][]} X - Data points
   */
  initializeCentroids(X) {
    const features = X[0].length;
    this.centroids = [];

    for (let i = 0; i < this.k; i++) {
      const centroid = [];
      for (let j = 0; j < features; j++) {
        const min = Math.min(...X.map(point => point[j]));
        const max = Math.max(...X.map(point => point[j]));
        centroid.push(Math.random() * (max - min) + min);
      }
      this.centroids.push(centroid);
    }
  }

  /**
   * Assign points to nearest centroid
   * @param {number[][]} X - Data points
   * @returns {number[]} Cluster assignments
   */
  assignClusters(X) {
    return X.map(point => {
      let minDistance = Infinity;
      let cluster = 0;

      for (let i = 0; i < this.k; i++) {
        const dist = this.distance(point, this.centroids[i]);
        if (dist < minDistance) {
          minDistance = dist;
          cluster = i;
        }
      }

      return cluster;
    });
  }

  /**
   * Update centroids based on cluster assignments
   * @param {number[][]} X - Data points
   * @param {number[]} labels - Cluster assignments
   */
  updateCentroids(X, labels) {
    const newCentroids = [];

    for (let i = 0; i < this.k; i++) {
      const clusterPoints = X.filter((_, index) => labels[index] === i);
      
      if (clusterPoints.length === 0) {
        newCentroids.push([...this.centroids[i]]);
        continue;
      }

      const centroid = [];
      for (let j = 0; j < X[0].length; j++) {
        const mean = clusterPoints.reduce((sum, point) => sum + point[j], 0) / clusterPoints.length;
        centroid.push(mean);
      }
      newCentroids.push(centroid);
    }

    this.centroids = newCentroids;
  }

  /**
   * Check if centroids have converged
   * @param {number[][]} oldCentroids - Previous centroids
   * @returns {boolean} True if converged
   */
  hasConverged(oldCentroids) {
    for (let i = 0; i < this.k; i++) {
      const dist = this.distance(this.centroids[i], oldCentroids[i]);
      if (dist > this.tolerance) {
        return false;
      }
    }
    return true;
  }

  /**
   * Fit the K-means model
   * @param {number[][]} X - Training data
   */
  fit(X) {
    if (X.length < this.k) {
      throw new Error('Number of samples must be greater than k');
    }

    this.initializeCentroids(X);

    for (let iteration = 0; iteration < this.maxIterations; iteration++) {
      const oldCentroids = this.centroids.map(centroid => [...centroid]);
      
      this.labels = this.assignClusters(X);
      this.updateCentroids(X, this.labels);

      if (this.hasConverged(oldCentroids)) {
        break;
      }
    }

    this.trained = true;
  }

  /**
   * Predict cluster for new data points
   * @param {number[][]} X - Data points to predict
   * @returns {number[]} Cluster assignments
   */
  predict(X) {
    if (!this.trained) {
      throw new Error('Model must be trained before making predictions');
    }

    return this.assignClusters(X);
  }

  /**
   * Calculate within-cluster sum of squares (WCSS)
   * @param {number[][]} X - Data points
   * @returns {number} WCSS value
   */
  inertia(X) {
    if (!this.trained) {
      throw new Error('Model must be trained before calculating inertia');
    }

    let wcss = 0;
    for (let i = 0; i < X.length; i++) {
      const cluster = this.labels[i];
      wcss += Math.pow(this.distance(X[i], this.centroids[cluster]), 2);
    }

    return wcss;
  }
}

/**
 * K-Nearest Neighbors implementation
 */
class KNearestNeighbors {
  constructor(k = 3) {
    this.k = k;
    this.X_train = [];
    this.y_train = [];
    this.trained = false;
  }

  /**
   * Calculate Euclidean distance between two points
   * @param {number[]} point1 - First point
   * @param {number[]} point2 - Second point
   * @returns {number} Distance
   */
  distance(point1, point2) {
    return Math.sqrt(
      point1.reduce((sum, val, i) => sum + Math.pow(val - point2[i], 2), 0)
    );
  }

  /**
   * Train the KNN model (just store the data)
   * @param {number[][]} X - Training features
   * @param {number[]} y - Training labels
   */
  fit(X, y) {
    if (X.length !== y.length) {
      throw new Error('X and y must have the same length');
    }

    this.X_train = X.map(row => [...row]);
    this.y_train = [...y];
    this.trained = true;
  }

  /**
   * Find k nearest neighbors for a point
   * @param {number[]} point - Query point
   * @returns {Array} Array of {distance, label} objects
   */
  findNeighbors(point) {
    const distances = this.X_train.map((trainPoint, i) => ({
      distance: this.distance(point, trainPoint),
      label: this.y_train[i]
    }));

    return distances
      .sort((a, b) => a.distance - b.distance)
      .slice(0, this.k);
  }

  /**
   * Predict labels for new data points
   * @param {number[][]} X - Data points to predict
   * @returns {number[]} Predicted labels
   */
  predict(X) {
    if (!this.trained) {
      throw new Error('Model must be trained before making predictions');
    }

    return X.map(point => {
      const neighbors = this.findNeighbors(point);
      
      // For classification: majority vote
      const labelCounts = {};
      neighbors.forEach(neighbor => {
        labelCounts[neighbor.label] = (labelCounts[neighbor.label] || 0) + 1;
      });

      return parseInt(Object.keys(labelCounts).reduce((a, b) => 
        labelCounts[a] > labelCounts[b] ? a : b
      ));
    });
  }

  /**
   * Calculate accuracy score
   * @param {number[][]} X - Test features
   * @param {number[]} y - True labels
   * @returns {number} Accuracy score (0-1)
   */
  score(X, y) {
    const predictions = this.predict(X);
    const correct = predictions.reduce((count, pred, i) => 
      count + (pred === y[i] ? 1 : 0), 0);
    
    return correct / y.length;
  }
}

module.exports = {
  LinearRegression,
  KMeans,
  KNearestNeighbors
};
