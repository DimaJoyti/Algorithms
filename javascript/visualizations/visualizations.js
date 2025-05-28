/**
 * Algorithm Visualization System
 * Interactive demonstrations of various algorithms
 */

class AlgorithmVisualizer {
    constructor() {
        this.canvas = null;
        this.ctx = null;
        this.currentAlgorithm = null;
        this.isRunning = false;
        this.isPaused = false;
        this.animationSpeed = 5;
        this.stats = {
            comparisons: 0,
            swaps: 0,
            arrayAccess: 0,
            startTime: 0
        };
        
        // Algorithm data
        this.array = [];
        this.arraySize = 50;
        this.currentStep = 0;
        this.steps = [];
        
        // Colors
        this.colors = {
            default: '#3498db',
            comparing: '#e74c3c',
            swapping: '#f39c12',
            sorted: '#2ecc71',
            pivot: '#9b59b6',
            current: '#e67e22'
        };
        
        this.algorithms = {
            'bubble-sort': {
                name: 'Bubble Sort',
                description: 'Bubble Sort repeatedly steps through the list, compares adjacent elements and swaps them if they are in the wrong order.',
                timeComplexity: 'O(n²)',
                spaceComplexity: 'O(1)',
                visualize: this.visualizeBubbleSort.bind(this)
            },
            'quick-sort': {
                name: 'Quick Sort',
                description: 'Quick Sort picks a pivot element and partitions the array around it, then recursively sorts the sub-arrays.',
                timeComplexity: 'O(n log n) average, O(n²) worst',
                spaceComplexity: 'O(log n)',
                visualize: this.visualizeQuickSort.bind(this)
            },
            'binary-search': {
                name: 'Binary Search',
                description: 'Binary Search finds a target value in a sorted array by repeatedly dividing the search interval in half.',
                timeComplexity: 'O(log n)',
                spaceComplexity: 'O(1)',
                visualize: this.visualizeBinarySearch.bind(this)
            },
            'linear-regression': {
                name: 'Linear Regression',
                description: 'Linear Regression finds the best fit line through data points using least squares method.',
                timeComplexity: 'O(n)',
                spaceComplexity: 'O(1)',
                visualize: this.visualizeLinearRegression.bind(this)
            }
        };
    }
    
    init() {
        this.canvas = document.getElementById('visualizationCanvas');
        this.ctx = this.canvas.getContext('2d');
        
        this.setupEventListeners();
        this.generateRandomArray();
        this.draw();
    }
    
    setupEventListeners() {
        // Algorithm selection
        document.querySelectorAll('.algorithm-card').forEach(card => {
            card.addEventListener('click', (e) => {
                this.selectAlgorithm(card.dataset.algorithm);
            });
        });
        
        // Control buttons
        document.getElementById('startBtn').addEventListener('click', () => this.start());
        document.getElementById('pauseBtn').addEventListener('click', () => this.pause());
        document.getElementById('resetBtn').addEventListener('click', () => this.reset());
        document.getElementById('randomizeBtn').addEventListener('click', () => this.randomize());
        
        // Speed control
        document.getElementById('speedSlider').addEventListener('input', (e) => {
            this.animationSpeed = parseInt(e.target.value);
        });
    }
    
    selectAlgorithm(algorithmKey) {
        // Update UI
        document.querySelectorAll('.algorithm-card').forEach(card => {
            card.classList.remove('active');
        });
        document.querySelector(`[data-algorithm="${algorithmKey}"]`).classList.add('active');
        
        // Update algorithm info
        const algorithm = this.algorithms[algorithmKey];
        if (algorithm) {
            this.currentAlgorithm = algorithmKey;
            document.getElementById('algorithmTitle').textContent = algorithm.name;
            document.getElementById('algorithmDescription').textContent = algorithm.description;
            document.getElementById('timeComplexity').textContent = algorithm.timeComplexity;
            document.getElementById('spaceComplexity').textContent = algorithm.spaceComplexity;
        }
        
        this.reset();
    }
    
    generateRandomArray() {
        this.array = [];
        for (let i = 0; i < this.arraySize; i++) {
            this.array.push(Math.floor(Math.random() * 300) + 10);
        }
    }
    
    start() {
        if (!this.currentAlgorithm) {
            alert('Please select an algorithm first!');
            return;
        }
        
        if (this.isPaused) {
            this.isPaused = false;
            this.animate();
            return;
        }
        
        this.isRunning = true;
        this.stats.startTime = Date.now();
        this.resetStats();
        
        // Generate steps for the selected algorithm
        this.algorithms[this.currentAlgorithm].visualize();
        this.animate();
    }
    
    pause() {
        this.isPaused = !this.isPaused;
        document.getElementById('pauseBtn').textContent = this.isPaused ? '▶️ Resume' : '⏸️ Pause';
    }
    
    reset() {
        this.isRunning = false;
        this.isPaused = false;
        this.currentStep = 0;
        this.steps = [];
        this.resetStats();
        document.getElementById('pauseBtn').textContent = '⏸️ Pause';
        this.draw();
    }
    
    randomize() {
        this.generateRandomArray();
        this.reset();
    }
    
    resetStats() {
        this.stats.comparisons = 0;
        this.stats.swaps = 0;
        this.stats.arrayAccess = 0;
        this.updateStatsDisplay();
    }
    
    updateStatsDisplay() {
        document.getElementById('comparisons').textContent = this.stats.comparisons;
        document.getElementById('swaps').textContent = this.stats.swaps;
        document.getElementById('arrayAccess').textContent = this.stats.arrayAccess;
        
        if (this.stats.startTime) {
            const elapsed = Date.now() - this.stats.startTime;
            document.getElementById('timeElapsed').textContent = elapsed + 'ms';
        }
    }
    
    animate() {
        if (!this.isRunning || this.isPaused) return;
        
        if (this.currentStep < this.steps.length) {
            const step = this.steps[this.currentStep];
            this.executeStep(step);
            this.draw();
            this.updateStatsDisplay();
            this.currentStep++;
            
            setTimeout(() => this.animate(), 1000 / this.animationSpeed);
        } else {
            this.isRunning = false;
        }
    }
    
    executeStep(step) {
        switch (step.type) {
            case 'compare':
                this.stats.comparisons++;
                break;
            case 'swap':
                this.stats.swaps++;
                [this.array[step.i], this.array[step.j]] = [this.array[step.j], this.array[step.i]];
                break;
            case 'access':
                this.stats.arrayAccess++;
                break;
        }
    }
    
    draw() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        
        if (this.currentAlgorithm === 'linear-regression') {
            this.drawLinearRegression();
            return;
        }
        
        const barWidth = this.canvas.width / this.array.length;
        const maxHeight = this.canvas.height - 50;
        const maxValue = Math.max(...this.array);
        
        for (let i = 0; i < this.array.length; i++) {
            const height = (this.array[i] / maxValue) * maxHeight;
            const x = i * barWidth;
            const y = this.canvas.height - height;
            
            // Determine color based on current step
            let color = this.colors.default;
            if (this.currentStep < this.steps.length) {
                const step = this.steps[this.currentStep];
                if (step && (step.i === i || step.j === i)) {
                    color = step.type === 'compare' ? this.colors.comparing : this.colors.swapping;
                }
            }
            
            this.ctx.fillStyle = color;
            this.ctx.fillRect(x, y, barWidth - 1, height);
            
            // Draw value on top of bar
            this.ctx.fillStyle = '#2c3e50';
            this.ctx.font = '10px Arial';
            this.ctx.textAlign = 'center';
            this.ctx.fillText(this.array[i], x + barWidth/2, y - 5);
        }
    }
    
    drawLinearRegression() {
        // Generate sample data points
        const points = [];
        for (let i = 0; i < 20; i++) {
            const x = (i / 19) * this.canvas.width;
            const y = this.canvas.height/2 + (Math.random() - 0.5) * 200 + (i - 10) * 5;
            points.push({x, y});
        }
        
        // Draw points
        this.ctx.fillStyle = this.colors.comparing;
        points.forEach(point => {
            this.ctx.beginPath();
            this.ctx.arc(point.x, point.y, 4, 0, 2 * Math.PI);
            this.ctx.fill();
        });
        
        // Calculate and draw regression line
        const n = points.length;
        const sumX = points.reduce((sum, p) => sum + p.x, 0);
        const sumY = points.reduce((sum, p) => sum + p.y, 0);
        const sumXY = points.reduce((sum, p) => sum + p.x * p.y, 0);
        const sumXX = points.reduce((sum, p) => sum + p.x * p.x, 0);
        
        const slope = (n * sumXY - sumX * sumY) / (n * sumXX - sumX * sumX);
        const intercept = (sumY - slope * sumX) / n;
        
        // Draw regression line
        this.ctx.strokeStyle = this.colors.sorted;
        this.ctx.lineWidth = 3;
        this.ctx.beginPath();
        this.ctx.moveTo(0, intercept);
        this.ctx.lineTo(this.canvas.width, slope * this.canvas.width + intercept);
        this.ctx.stroke();
        
        // Draw equation
        this.ctx.fillStyle = '#2c3e50';
        this.ctx.font = '16px Arial';
        this.ctx.fillText(`y = ${slope.toFixed(2)}x + ${intercept.toFixed(2)}`, 20, 30);
    }
    
    // Algorithm implementations
    visualizeBubbleSort() {
        this.steps = [];
        const arr = [...this.array];
        
        for (let i = 0; i < arr.length - 1; i++) {
            for (let j = 0; j < arr.length - i - 1; j++) {
                this.steps.push({type: 'compare', i: j, j: j + 1});
                if (arr[j] > arr[j + 1]) {
                    this.steps.push({type: 'swap', i: j, j: j + 1});
                    [arr[j], arr[j + 1]] = [arr[j + 1], arr[j]];
                }
            }
        }
    }
    
    visualizeQuickSort() {
        this.steps = [];
        const arr = [...this.array];
        
        const quickSort = (low, high) => {
            if (low < high) {
                const pi = partition(low, high);
                quickSort(low, pi - 1);
                quickSort(pi + 1, high);
            }
        };
        
        const partition = (low, high) => {
            const pivot = arr[high];
            let i = low - 1;
            
            for (let j = low; j < high; j++) {
                this.steps.push({type: 'compare', i: j, j: high});
                if (arr[j] < pivot) {
                    i++;
                    this.steps.push({type: 'swap', i: i, j: j});
                    [arr[i], arr[j]] = [arr[j], arr[i]];
                }
            }
            
            this.steps.push({type: 'swap', i: i + 1, j: high});
            [arr[i + 1], arr[high]] = [arr[high], arr[i + 1]];
            return i + 1;
        };
        
        quickSort(0, arr.length - 1);
    }
    
    visualizeBinarySearch() {
        this.steps = [];
        // First sort the array for binary search
        this.array.sort((a, b) => a - b);
        const target = this.array[Math.floor(Math.random() * this.array.length)];
        
        let left = 0;
        let right = this.array.length - 1;
        
        while (left <= right) {
            const mid = Math.floor((left + right) / 2);
            this.steps.push({type: 'compare', i: mid, j: mid, target});
            
            if (this.array[mid] === target) {
                this.steps.push({type: 'found', i: mid});
                break;
            } else if (this.array[mid] < target) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
    }
    
    visualizeLinearRegression() {
        this.steps = [];
        // This will be handled in the draw function
    }
}
