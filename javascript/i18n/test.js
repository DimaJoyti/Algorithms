const { I18n } = require('./index');

describe('Internationalization (i18n)', () => {
  let i18n;

  beforeEach(() => {
    i18n = new I18n('en');
  });

  describe('I18n class', () => {
    test('I18n class exists', () => {
      expect(typeof I18n).toEqual('function');
    });

    test('creates instance with default locale', () => {
      expect(i18n.getCurrentLocale()).toBe('en');
    });

    test('creates instance with custom locale', () => {
      const customI18n = new I18n('uk');
      expect(customI18n.getCurrentLocale()).toBe('uk');
    });
  });

  describe('Locale management', () => {
    test('getAvailableLocales returns array of locales', () => {
      const locales = i18n.getAvailableLocales();
      expect(Array.isArray(locales)).toBe(true);
      expect(locales.length).toBeGreaterThan(0);
      expect(locales).toContain('en');
    });

    test('setLocale changes current locale', () => {
      const result = i18n.setLocale('uk');
      expect(result).toBe(true);
      expect(i18n.getCurrentLocale()).toBe('uk');
    });

    test('setLocale returns false for invalid locale', () => {
      const result = i18n.setLocale('invalid');
      expect(result).toBe(false);
      expect(i18n.getCurrentLocale()).toBe('en'); // Should remain unchanged
    });

    test('getCurrentLocale returns current locale', () => {
      expect(i18n.getCurrentLocale()).toBe('en');
      i18n.setLocale('uk');
      expect(i18n.getCurrentLocale()).toBe('uk');
    });
  });

  describe('Translation functionality', () => {
    test('t function translates keys', () => {
      // Test with embedded translations
      const englishText = i18n.t('ui.start');
      expect(typeof englishText).toBe('string');
      expect(englishText.length).toBeGreaterThan(0);
    });

    test('t function falls back to key if translation not found', () => {
      const result = i18n.t('nonexistent.key');
      expect(result).toBe('nonexistent.key');
    });

    test('t function interpolates parameters', () => {
      // Add test translation with parameters
      i18n.addTranslations('en', {
        test: {
          greeting: 'Hello {{name}}!'
        }
      });

      const result = i18n.t('test.greeting', { name: 'World' });
      expect(result).toBe('Hello World!');
    });

    test('getTranslation returns translation for specific locale', () => {
      const translation = i18n.getTranslation('ui.start', 'en');
      expect(typeof translation).toBe('string');
    });

    test('getTranslation returns null for invalid locale', () => {
      const translation = i18n.getTranslation('ui.start', 'invalid');
      expect(translation).toBeNull();
    });
  });

  describe('Nested value access', () => {
    test('getNestedValue retrieves nested object values', () => {
      const obj = {
        level1: {
          level2: {
            value: 'test'
          }
        }
      };

      const result = i18n.getNestedValue(obj, 'level1.level2.value');
      expect(result).toBe('test');
    });

    test('getNestedValue returns null for invalid path', () => {
      const obj = { test: 'value' };
      const result = i18n.getNestedValue(obj, 'invalid.path');
      expect(result).toBeNull();
    });
  });

  describe('String interpolation', () => {
    test('interpolate replaces placeholders', () => {
      const text = 'Hello {{name}}, you have {{count}} messages';
      const params = { name: 'John', count: 5 };
      const result = i18n.interpolate(text, params);
      expect(result).toBe('Hello John, you have 5 messages');
    });

    test('interpolate handles missing parameters', () => {
      const text = 'Hello {{name}}, you have {{count}} messages';
      const params = { name: 'John' };
      const result = i18n.interpolate(text, params);
      expect(result).toBe('Hello John, you have {{count}} messages');
    });

    test('interpolate handles non-string input', () => {
      const result = i18n.interpolate(123, {});
      expect(result).toBe(123);
    });
  });

  describe('Algorithm information', () => {
    test('getAlgorithmInfo returns algorithm information', () => {
      const info = i18n.getAlgorithmInfo('sorting', 'bubbleSort');
      expect(typeof info).toBe('object');
      expect(info).toHaveProperty('category');
      expect(info).toHaveProperty('name');
      expect(info).toHaveProperty('description');
    });

    test('getLocalizedAlgorithms returns localized algorithms', () => {
      const algorithms = i18n.getLocalizedAlgorithms();
      expect(typeof algorithms).toBe('object');
    });
  });

  describe('Formatting functions', () => {
    test('formatNumber formats numbers', () => {
      const result = i18n.formatNumber(1234.56);
      expect(typeof result).toBe('string');
    });

    test('formatDate formats dates', () => {
      const date = new Date('2023-01-01');
      const result = i18n.formatDate(date);
      expect(typeof result).toBe('string');
    });

    test('formatNumber handles invalid locale gracefully', () => {
      i18n.setLocale('invalid');
      const result = i18n.formatNumber(1234);
      expect(typeof result).toBe('string');
    });
  });

  describe('Complexity information', () => {
    test('getComplexityInfo returns complexity information', () => {
      const info = i18n.getComplexityInfo('O(n)', 'O(1)');
      expect(info).toHaveProperty('time');
      expect(info).toHaveProperty('space');
      expect(info.time).toHaveProperty('label');
      expect(info.time).toHaveProperty('value');
      expect(info.time.value).toBe('O(n)');
      expect(info.space.value).toBe('O(1)');
    });
  });

  describe('UI and messages', () => {
    test('getUIText returns UI text', () => {
      const text = i18n.getUIText('start');
      expect(typeof text).toBe('string');
    });

    test('getMessage returns message', () => {
      // Add test message
      i18n.addTranslations('en', {
        messages: {
          test: 'Test message'
        }
      });

      const message = i18n.getMessage('test');
      expect(message).toBe('Test message');
    });

    test('getMessage interpolates parameters', () => {
      i18n.addTranslations('en', {
        messages: {
          greeting: 'Hello {{name}}'
        }
      });

      const message = i18n.getMessage('greeting', { name: 'World' });
      expect(message).toBe('Hello World');
    });
  });

  describe('Translation management', () => {
    test('addTranslations adds new translations', () => {
      i18n.addTranslations('en', {
        test: {
          newKey: 'New value'
        }
      });

      const result = i18n.t('test.newKey');
      expect(result).toBe('New value');
    });

    test('addTranslations merges with existing translations', () => {
      i18n.addTranslations('en', {
        ui: {
          newButton: 'New Button'
        }
      });

      // Should still have old translations
      expect(i18n.t('ui.start')).toBeTruthy();
      // Should have new translation
      expect(i18n.t('ui.newButton')).toBe('New Button');
    });

    test('removeLocale removes locale', () => {
      i18n.addTranslations('test', { key: 'value' });
      expect(i18n.getAvailableLocales()).toContain('test');
      
      const result = i18n.removeLocale('test');
      expect(result).toBe(true);
      expect(i18n.getAvailableLocales()).not.toContain('test');
    });

    test('removeLocale prevents removing fallback locale', () => {
      const result = i18n.removeLocale('en');
      expect(result).toBe(false);
      expect(i18n.getAvailableLocales()).toContain('en');
    });
  });

  describe('Locale switching', () => {
    test('translations change when locale is switched', () => {
      // Test with embedded translations
      i18n.setLocale('en');
      const englishText = i18n.t('ui.start');
      
      i18n.setLocale('uk');
      const ukrainianText = i18n.t('ui.start');
      
      expect(englishText).not.toBe(ukrainianText);
    });

    test('fallback works when translation missing in current locale', () => {
      // Add translation only in English
      i18n.addTranslations('en', {
        test: {
          onlyInEnglish: 'English only'
        }
      });

      i18n.setLocale('uk');
      const result = i18n.t('test.onlyInEnglish');
      expect(result).toBe('English only'); // Should fallback to English
    });
  });
});
