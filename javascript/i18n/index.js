// --- Directions
// Implement internationalization (i18n) system for the algorithms collection.
// This module provides multi-language support for algorithm names, descriptions,
// and UI elements.
// --- Examples
//   const i18n = new I18n('en');
//   i18n.t('algorithms.sorting.bubbleSort.name'); // "Bubble Sort"
//   i18n.setLocale('uk');
//   i18n.t('algorithms.sorting.bubbleSort.name'); // "Бульбашкове сортування"

const fs = require('fs');
const path = require('path');

/**
 * Internationalization system for algorithms collection
 */
class I18n {
  constructor(defaultLocale = 'en') {
    this.currentLocale = defaultLocale;
    this.translations = new Map();
    this.fallbackLocale = 'en';
    
    // Load available locales
    this.loadLocales();
  }

  /**
   * Load all available locale files
   */
  loadLocales() {
    try {
      const localesDir = path.join(__dirname, '..', 'locales');
      const files = fs.readdirSync(localesDir);
      
      files.forEach(file => {
        if (file.endsWith('.json')) {
          const locale = path.basename(file, '.json');
          const filePath = path.join(localesDir, file);
          const content = fs.readFileSync(filePath, 'utf8');
          this.translations.set(locale, JSON.parse(content));
        }
      });
    } catch (error) {
      console.warn('Could not load locale files:', error.message);
      // Fallback to embedded translations
      this.loadEmbeddedTranslations();
    }
  }

  /**
   * Load embedded translations as fallback
   */
  loadEmbeddedTranslations() {
    // Basic English translations
    this.translations.set('en', {
      algorithms: {
        strings: { category: 'String Algorithms' },
        math: { category: 'Mathematical Algorithms' },
        sorting: { category: 'Sorting Algorithms' },
        search: { category: 'Search Algorithms' },
        graphs: { category: 'Graph Algorithms' },
        machineLearning: { category: 'Machine Learning' }
      },
      ui: {
        start: 'Start',
        pause: 'Pause',
        reset: 'Reset',
        speed: 'Speed'
      }
    });

    // Basic Ukrainian translations
    this.translations.set('uk', {
      algorithms: {
        strings: { category: 'Алгоритми для рядків' },
        math: { category: 'Математичні алгоритми' },
        sorting: { category: 'Алгоритми сортування' },
        search: { category: 'Алгоритми пошуку' },
        graphs: { category: 'Алгоритми графів' },
        machineLearning: { category: 'Машинне навчання' }
      },
      ui: {
        start: 'Почати',
        pause: 'Пауза',
        reset: 'Скинути',
        speed: 'Швидкість'
      }
    });
  }

  /**
   * Get available locales
   * @returns {string[]} Array of available locale codes
   */
  getAvailableLocales() {
    return Array.from(this.translations.keys());
  }

  /**
   * Set current locale
   * @param {string} locale - Locale code (e.g., 'en', 'uk')
   * @returns {boolean} True if locale was set successfully
   */
  setLocale(locale) {
    if (this.translations.has(locale)) {
      this.currentLocale = locale;
      return true;
    }
    console.warn(`Locale '${locale}' not found, keeping current locale '${this.currentLocale}'`);
    return false;
  }

  /**
   * Get current locale
   * @returns {string} Current locale code
   */
  getCurrentLocale() {
    return this.currentLocale;
  }

  /**
   * Translate a key to current locale
   * @param {string} key - Translation key (dot notation)
   * @param {Object} params - Parameters for interpolation
   * @returns {string} Translated text
   */
  t(key, params = {}) {
    const translation = this.getTranslation(key, this.currentLocale) || 
                       this.getTranslation(key, this.fallbackLocale) || 
                       key;
    
    return this.interpolate(translation, params);
  }

  /**
   * Get translation for specific locale
   * @param {string} key - Translation key
   * @param {string} locale - Locale code
   * @returns {string|null} Translation or null if not found
   */
  getTranslation(key, locale) {
    const translations = this.translations.get(locale);
    if (!translations) return null;

    return this.getNestedValue(translations, key);
  }

  /**
   * Get nested value from object using dot notation
   * @param {Object} obj - Object to search in
   * @param {string} path - Dot notation path
   * @returns {any} Value or null if not found
   */
  getNestedValue(obj, path) {
    return path.split('.').reduce((current, key) => {
      return current && current[key] !== undefined ? current[key] : null;
    }, obj);
  }

  /**
   * Interpolate parameters into translation string
   * @param {string} text - Text with placeholders
   * @param {Object} params - Parameters to interpolate
   * @returns {string} Interpolated text
   */
  interpolate(text, params) {
    if (typeof text !== 'string') return text;

    return text.replace(/\{\{(\w+)\}\}/g, (match, key) => {
      return params[key] !== undefined ? params[key] : match;
    });
  }

  /**
   * Get algorithm information in current locale
   * @param {string} category - Algorithm category
   * @param {string} algorithm - Algorithm name
   * @returns {Object} Algorithm information
   */
  getAlgorithmInfo(category, algorithm) {
    const categoryInfo = this.t(`algorithms.${category}`) || {};
    const algorithmInfo = this.t(`algorithms.${category}.${algorithm}`) || {};

    return {
      category: categoryInfo.category || category,
      name: algorithmInfo.name || algorithm,
      description: algorithmInfo.description || '',
      ...algorithmInfo
    };
  }

  /**
   * Get all algorithms with localized information
   * @returns {Object} Localized algorithms object
   */
  getLocalizedAlgorithms() {
    const result = {};
    const algorithmsKey = 'algorithms';
    const algorithms = this.getTranslation(algorithmsKey, this.currentLocale) || {};

    Object.keys(algorithms).forEach(category => {
      if (typeof algorithms[category] === 'object' && algorithms[category].category) {
        result[category] = {
          category: algorithms[category].category,
          description: algorithms[category].description || '',
          algorithms: {}
        };

        Object.keys(algorithms[category]).forEach(key => {
          if (key !== 'category' && key !== 'description' && 
              typeof algorithms[category][key] === 'object') {
            result[category].algorithms[key] = algorithms[category][key];
          }
        });
      }
    });

    return result;
  }

  /**
   * Format number according to current locale
   * @param {number} number - Number to format
   * @param {Object} options - Formatting options
   * @returns {string} Formatted number
   */
  formatNumber(number, options = {}) {
    try {
      return new Intl.NumberFormat(this.currentLocale, options).format(number);
    } catch (error) {
      return number.toString();
    }
  }

  /**
   * Format date according to current locale
   * @param {Date} date - Date to format
   * @param {Object} options - Formatting options
   * @returns {string} Formatted date
   */
  formatDate(date, options = {}) {
    try {
      return new Intl.DateTimeFormat(this.currentLocale, options).format(date);
    } catch (error) {
      return date.toString();
    }
  }

  /**
   * Get complexity information in current locale
   * @param {string} timeComplexity - Time complexity
   * @param {string} spaceComplexity - Space complexity
   * @returns {Object} Localized complexity information
   */
  getComplexityInfo(timeComplexity, spaceComplexity) {
    return {
      time: {
        label: this.t('complexity.time'),
        value: timeComplexity
      },
      space: {
        label: this.t('complexity.space'),
        value: spaceComplexity
      }
    };
  }

  /**
   * Get UI text in current locale
   * @param {string} key - UI element key
   * @returns {string} Localized UI text
   */
  getUIText(key) {
    return this.t(`ui.${key}`);
  }

  /**
   * Get message in current locale
   * @param {string} key - Message key
   * @param {Object} params - Parameters for interpolation
   * @returns {string} Localized message
   */
  getMessage(key, params = {}) {
    return this.t(`messages.${key}`, params);
  }

  /**
   * Add or update translations for a locale
   * @param {string} locale - Locale code
   * @param {Object} translations - Translation object
   */
  addTranslations(locale, translations) {
    const existing = this.translations.get(locale) || {};
    this.translations.set(locale, { ...existing, ...translations });
  }

  /**
   * Remove a locale
   * @param {string} locale - Locale code to remove
   * @returns {boolean} True if locale was removed
   */
  removeLocale(locale) {
    if (locale === this.fallbackLocale) {
      console.warn(`Cannot remove fallback locale '${locale}'`);
      return false;
    }
    
    return this.translations.delete(locale);
  }
}

// Create default instance
const defaultI18n = new I18n();

module.exports = {
  I18n,
  default: defaultI18n,
  t: (key, params) => defaultI18n.t(key, params),
  setLocale: (locale) => defaultI18n.setLocale(locale),
  getCurrentLocale: () => defaultI18n.getCurrentLocale(),
  getAvailableLocales: () => defaultI18n.getAvailableLocales(),
  getAlgorithmInfo: (category, algorithm) => defaultI18n.getAlgorithmInfo(category, algorithm),
  getLocalizedAlgorithms: () => defaultI18n.getLocalizedAlgorithms()
};
