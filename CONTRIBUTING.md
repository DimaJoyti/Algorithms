# 🤝 Contributing Guide

Thank you for your interest in contributing to the Algorithms project! This document contains instructions and guidelines for contributors.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Types of Contributions](#types-of-contributions)
- [Development Process](#development-process)
- [Code Standards](#code-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)

## 🤝 Code of Conduct

This project adheres to the [Contributor Covenant](https://www.contributor-covenant.org/). By participating, you agree to uphold this code.

## 🚀 Getting Started

### 1. Environment Setup

```bash
# Fork the repository on GitHub
# Clone your fork
git clone https://github.com/your-username/algorithms.git
cd algorithms

# Add upstream remote
git remote add upstream https://github.com/original-owner/algorithms.git

# Install dependencies
cd javascript && npm install
cd ../python && pip install -r requirements.txt
```

### 2. Creating a Branch

```bash
# Update main branch
git checkout main
git pull upstream main

# Create new branch for your feature
git checkout -b feature/algorithm-name
# or
git checkout -b fix/issue-description
```

## 🎯 Types of Contributions

### ✅ What we welcome

- **New algorithms** with complete implementation and tests
- **Bug fixes** in existing code
- **Performance improvements** of algorithms
- **Adding tests** for existing algorithms
- **Documentation improvements** and comments
- **Documentation translations**
- **Algorithm visualizations**

### ❌ What we don't accept

- Incomplete implementations without tests
- Code without comments and documentation
- Duplicates of existing algorithms without improvements
- Code that doesn't follow project standards

## 🔄 Процес розробки

### Додавання нового алгоритму

#### JavaScript
```bash
# Створіть папку для алгоритму
mkdir javascript/algorithm-name
cd javascript/algorithm-name

# Створіть файли
touch index.js test.js
```

**Структура index.js:**
```javascript
// --- Directions
// Опис алгоритму та його призначення
// --- Examples
//   algorithmName(input) --> expected output

function algorithmName(input) {
  // Ваша реалізація тут
}

module.exports = algorithmName;
```

**Структура test.js:**
```javascript
const algorithmName = require('./index');

test('algorithm function exists', () => {
  expect(typeof algorithmName).toEqual('function');
});

test('handles basic case', () => {
  expect(algorithmName(input)).toEqual(expectedOutput);
});

// Додайте більше тестів для edge cases
```

#### Python
```python
"""
Назва алгоритму

Опис алгоритму та його призначення.

Time Complexity: O(?)
Space Complexity: O(?)

Examples:
    >>> algorithm_name(input)
    expected_output
"""

def algorithm_name(input):
    """
    Детальний опис функції.

    Args:
        input: Опис параметра

    Returns:
        Опис повертаємого значення

    Raises:
        ValueError: Коли виникає помилка
    """
    # Ваша реалізація тут
    pass

if __name__ == "__main__":
    # Приклади використання
    print(algorithm_name(test_input))
```

## 📏 Стандарти коду

### JavaScript
- Використовуйте ES6+ синтаксис
- Дотримуйтесь ESLint правил
- Використовуйте camelCase для змінних та функцій
- Додавайте JSDoc коментарі для функцій

### Python
- Дотримуйтесь PEP 8
- Використовуйте type hints
- Додавайте docstrings для всіх функцій
- Використовуйте snake_case для змінних та функцій

### Загальні правила
- Код має бути читабельним та зрозумілим
- Додавайте коментарі для складних частин
- Уникайте магічних чисел
- Використовуйте описові назви змінних

## 🧪 Тестування

### JavaScript тести
```bash
# Запуск всіх тестів
npm test

# Запуск конкретного тесту
npm test algorithm-name

# Запуск з покриттям
npm test -- --coverage
```

### Python тести
```bash
# Запуск всіх тестів
pytest

# Запуск конкретного файлу
pytest test_algorithm.py

# Запуск з покриттям
pytest --cov=.
```

### Вимоги до тестів
- Мінімум 90% покриття коду
- Тести для edge cases
- Тести для невалідних входів
- Performance тести для великих даних

## 📖 Документація

### Обов'язкові елементи:
- Опис алгоритму та його призначення
- Часова та просторова складність (Big O)
- Приклади використання
- Пояснення підходу
- Посилання на джерела (якщо є)

### Формат коментарів:
```javascript
// --- Directions
// Детальний опис того, що робить алгоритм
// --- Examples
//   input -> output
//   input2 -> output2
// --- Complexity
//   Time: O(n)
//   Space: O(1)
```

## 🔀 Pull Request процес

### 1. Перед створенням PR
- [ ] Код пройшов всі тести
- [ ] Додано нові тести для нової функціональності
- [ ] Оновлено документацію
- [ ] Код відформатовано згідно стандартів
- [ ] Немає конфліктів з main гілкою

### 2. Створення PR
- Використовуйте описову назву
- Заповніть шаблон PR
- Додайте скріншоти (якщо потрібно)
- Позначте відповідні labels

### 3. Шаблон PR
```markdown
## Опис
Короткий опис змін

## Тип змін
- [ ] Новий алгоритм
- [ ] Виправлення бага
- [ ] Покращення продуктивності
- [ ] Оновлення документації

## Тестування
- [ ] Всі тести проходять
- [ ] Додано нові тести
- [ ] Перевірено на edge cases

## Checklist
- [ ] Код відповідає стандартам проекту
- [ ] Додано/оновлено документацію
- [ ] Немає breaking changes
```

## 🏷️ Система міток

- `good first issue` - Підходить для новачків
- `help wanted` - Потрібна допомога спільноти
- `bug` - Виправлення помилки
- `enhancement` - Нова функціональність
- `documentation` - Покращення документації
- `performance` - Оптимізація продуктивності

## 🆘 Отримання допомоги

- Створіть issue для обговорення
- Приєднайтесь до Discussions
- Зверніться до мейнтейнерів
- Перегляньте існуючі PR для прикладів

## 🎉 Визнання

Всі контрибутори будуть додані до списку в README.md. Дякуємо за ваш внесок!

---

**Пам'ятайте: Немає занадто маленького внеску. Кожна допомога цінується!** 🙏
