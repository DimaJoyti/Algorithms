module.exports = {
  // Basic formatting
  semi: true,
  trailingComma: 'es5',
  singleQuote: true,
  printWidth: 80,
  tabWidth: 2,
  useTabs: false,
  
  // JavaScript specific
  quoteProps: 'as-needed',
  jsxSingleQuote: true,
  bracketSpacing: true,
  bracketSameLine: false,
  arrowParens: 'avoid',
  
  // Line endings
  endOfLine: 'lf',
  
  // Override for specific files
  overrides: [
    {
      files: '*.json',
      options: {
        printWidth: 120,
      },
    },
    {
      files: '*.md',
      options: {
        printWidth: 100,
        proseWrap: 'always',
      },
    },
  ],
};
