module.exports = {
  env: {
    browser: true,
    es2021: true,
    node: true,
    jest: true,
  },
  extends: [
    'eslint:recommended',
    'plugin:jest/recommended',
    'prettier'
  ],
  plugins: ['jest'],
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
  },
  rules: {
    // Code quality
    'no-console': 'warn',
    'no-debugger': 'error',
    'no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
    'no-var': 'error',
    'prefer-const': 'error',
    
    // Best practices
    'eqeqeq': ['error', 'always'],
    'curly': ['error', 'all'],
    'no-eval': 'error',
    'no-implied-eval': 'error',
    'no-new-func': 'error',
    
    // Style (handled by Prettier, but some logical rules)
    'max-len': ['warn', { code: 100, ignoreComments: true }],
    'max-lines': ['warn', { max: 300 }],
    'max-params': ['warn', 4],
    'complexity': ['warn', 10],
    
    // ES6+
    'arrow-spacing': 'error',
    'no-duplicate-imports': 'error',
    'prefer-arrow-callback': 'error',
    'prefer-template': 'error',
    
    // Jest specific
    'jest/expect-expect': 'error',
    'jest/no-disabled-tests': 'warn',
    'jest/no-focused-tests': 'error',
    'jest/prefer-to-have-length': 'warn',
    'jest/valid-expect': 'error',
  },
  overrides: [
    {
      files: ['**/*.test.js', '**/*.spec.js'],
      rules: {
        'no-console': 'off', // Allow console in tests
        'max-lines': 'off',  // Tests can be longer
      }
    }
  ]
};
