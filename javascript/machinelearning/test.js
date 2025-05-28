const {
  LinearRegression,
  KMeans,
  KNearestNeighbors
} = require('./index');

describe('Machine Learning Algorithms', () => {
  describe('LinearRegression', () => {
    test('LinearRegression class exists', () => {
      expect(typeof LinearRegression).toEqual('function');
    });

    test('creates instance with default values', () => {
      const lr = new LinearRegression();
      expect(lr.slope).toBe(0);
      expect(lr.intercept).toBe(0);
      expect(lr.trained).toBe(false);
    });

    test('fits simple linear data correctly', () => {
      const lr = new LinearRegression();
      const X = [[1], [2], [3], [4], [5]];
      const y = [2, 4, 6, 8, 10];

      lr.fit(X, y);

      expect(lr.trained).toBe(true);
      expect(lr.slope).toBeCloseTo(2, 5);
      expect(lr.intercept).toBeCloseTo(0, 5);
    });

    test('makes accurate predictions', () => {
      const lr = new LinearRegression();
      const X = [[1], [2], [3]];
      const y = [3, 5, 7];

      lr.fit(X, y);
      const predictions = lr.predict([[4], [5]]);

      expect(predictions[0]).toBeCloseTo(9, 1);
      expect(predictions[1]).toBeCloseTo(11, 1);
    });

    test('calculates R-squared score', () => {
      const lr = new LinearRegression();
      const X = [[1], [2], [3], [4]];
      const y = [2, 4, 6, 8];

      lr.fit(X, y);
      const score = lr.score(X, y);

      expect(score).toBeCloseTo(1, 5); // Perfect fit
    });

    test('throws error when X and y lengths differ', () => {
      const lr = new LinearRegression();
      expect(() => lr.fit([[1], [2]], [1])).toThrow();
    });

    test('throws error when predicting before training', () => {
      const lr = new LinearRegression();
      expect(() => lr.predict([[1]])).toThrow();
    });

    test('handles 1D input arrays', () => {
      const lr = new LinearRegression();
      const X = [1, 2, 3, 4];
      const y = [2, 4, 6, 8];

      lr.fit(X, y);
      expect(lr.trained).toBe(true);
    });
  });

  describe('KMeans', () => {
    test('KMeans class exists', () => {
      expect(typeof KMeans).toEqual('function');
    });

    test('creates instance with default parameters', () => {
      const kmeans = new KMeans();
      expect(kmeans.k).toBe(3);
      expect(kmeans.maxIterations).toBe(100);
      expect(kmeans.trained).toBe(false);
    });

    test('creates instance with custom parameters', () => {
      const kmeans = new KMeans(5, 50, 1e-3);
      expect(kmeans.k).toBe(5);
      expect(kmeans.maxIterations).toBe(50);
      expect(kmeans.tolerance).toBe(1e-3);
    });

    test('calculates distance correctly', () => {
      const kmeans = new KMeans();
      const dist = kmeans.distance([0, 0], [3, 4]);
      expect(dist).toBeCloseTo(5, 5);
    });

    test('fits and predicts on simple 2D data', () => {
      const kmeans = new KMeans(2);
      const X = [
        [1, 1], [1, 2], [2, 1], [2, 2], // Cluster 1
        [8, 8], [8, 9], [9, 8], [9, 9]  // Cluster 2
      ];

      kmeans.fit(X);

      expect(kmeans.trained).toBe(true);
      expect(kmeans.centroids.length).toBe(2);
      expect(kmeans.labels.length).toBe(8);

      const predictions = kmeans.predict([[1.5, 1.5], [8.5, 8.5]]);
      expect(predictions.length).toBe(2);
    });

    test('calculates inertia', () => {
      const kmeans = new KMeans(2);
      const X = [[1, 1], [2, 2], [8, 8], [9, 9]];

      kmeans.fit(X);
      const inertiaValue = kmeans.inertia(X);

      expect(typeof inertiaValue).toBe('number');
      expect(inertiaValue).toBeGreaterThanOrEqual(0);
    });

    test('throws error when k > number of samples', () => {
      const kmeans = new KMeans(5);
      const X = [[1, 1], [2, 2]]; // Only 2 samples

      expect(() => kmeans.fit(X)).toThrow();
    });

    test('throws error when predicting before training', () => {
      const kmeans = new KMeans();
      expect(() => kmeans.predict([[1, 1]])).toThrow();
    });

    test('handles edge case with empty clusters', () => {
      const kmeans = new KMeans(3);
      const X = [[1, 1], [1, 1], [1, 1]]; // All same points

      expect(() => kmeans.fit(X)).not.toThrow();
    });
  });

  describe('KNearestNeighbors', () => {
    test('KNearestNeighbors class exists', () => {
      expect(typeof KNearestNeighbors).toEqual('function');
    });

    test('creates instance with default k', () => {
      const knn = new KNearestNeighbors();
      expect(knn.k).toBe(3);
      expect(knn.trained).toBe(false);
    });

    test('creates instance with custom k', () => {
      const knn = new KNearestNeighbors(5);
      expect(knn.k).toBe(5);
    });

    test('calculates distance correctly', () => {
      const knn = new KNearestNeighbors();
      const dist = knn.distance([0, 0], [3, 4]);
      expect(dist).toBeCloseTo(5, 5);
    });

    test('fits and stores training data', () => {
      const knn = new KNearestNeighbors();
      const X = [[1, 1], [2, 2], [3, 3]];
      const y = [0, 0, 1];

      knn.fit(X, y);

      expect(knn.trained).toBe(true);
      expect(knn.X_train.length).toBe(3);
      expect(knn.y_train.length).toBe(3);
    });

    test('makes predictions correctly', () => {
      const knn = new KNearestNeighbors(3);
      const X = [
        [1, 1], [1, 2], [2, 1], // Class 0
        [5, 5], [5, 6], [6, 5]  // Class 1
      ];
      const y = [0, 0, 0, 1, 1, 1];

      knn.fit(X, y);
      const predictions = knn.predict([[1.5, 1.5], [5.5, 5.5]]);

      expect(predictions[0]).toBe(0);
      expect(predictions[1]).toBe(1);
    });

    test('finds neighbors correctly', () => {
      const knn = new KNearestNeighbors(2);
      const X = [[1, 1], [2, 2], [3, 3]];
      const y = [0, 1, 2];

      knn.fit(X, y);
      const neighbors = knn.findNeighbors([1.5, 1.5]);

      expect(neighbors.length).toBe(2);
      expect(neighbors[0].distance).toBeLessThanOrEqual(neighbors[1].distance);
    });

    test('calculates accuracy score', () => {
      const knn = new KNearestNeighbors(1);
      const X = [[1, 1], [2, 2], [3, 3], [4, 4]];
      const y = [0, 0, 1, 1];

      knn.fit(X, y);
      const score = knn.score(X, y);

      expect(score).toBe(1); // Perfect accuracy on training data
    });

    test('throws error when X and y lengths differ', () => {
      const knn = new KNearestNeighbors();
      expect(() => knn.fit([[1, 1], [2, 2]], [0])).toThrow();
    });

    test('throws error when predicting before training', () => {
      const knn = new KNearestNeighbors();
      expect(() => knn.predict([[1, 1]])).toThrow();
    });

    test('handles single neighbor case', () => {
      const knn = new KNearestNeighbors(1);
      const X = [[1, 1], [5, 5]];
      const y = [0, 1];

      knn.fit(X, y);
      const prediction = knn.predict([[2, 2]]);

      expect(prediction[0]).toBe(0); // Closer to [1,1]
    });
  });

  describe('Integration Tests', () => {
    test('all algorithms work together', () => {
      // Generate some test data
      const X = [];
      const y_regression = [];
      const y_classification = [];

      for (let i = 0; i < 20; i++) {
        const x = i / 2;
        X.push([x]);
        y_regression.push(2 * x + 1 + Math.random() * 0.1); // Linear with noise
        y_classification.push(x > 5 ? 1 : 0); // Binary classification
      }

      // Test Linear Regression
      const lr = new LinearRegression();
      lr.fit(X, y_regression);
      const lr_score = lr.score(X, y_regression);
      expect(lr_score).toBeGreaterThan(0.9);

      // Test KNN
      const knn = new KNearestNeighbors(3);
      knn.fit(X, y_classification);
      const knn_score = knn.score(X, y_classification);
      expect(knn_score).toBeGreaterThan(0.8);

      // Test K-Means
      const kmeans = new KMeans(2);
      const X_2d = X.map((x, i) => [x[0], y_regression[i]]);
      kmeans.fit(X_2d);
      expect(kmeans.trained).toBe(true);
    });

    test('handles edge cases gracefully', () => {
      // Test with minimal data
      const lr = new LinearRegression();
      lr.fit([[1], [2]], [1, 2]);
      expect(lr.trained).toBe(true);

      const knn = new KNearestNeighbors(1);
      knn.fit([[1, 1]], [0]);
      expect(knn.trained).toBe(true);

      const kmeans = new KMeans(1);
      kmeans.fit([[1, 1], [2, 2]]);
      expect(kmeans.trained).toBe(true);
    });
  });
});
