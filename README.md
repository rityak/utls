# utls

TLS-клиент с отпечатками браузеров, на нём работает и reality.

Форк [metacubex/utls](https://github.com/metacubex/utls) для [miray](https://github.com/rityak/miray). Наш код — в ветке `miray`.

## Что изменено

- Все TLS-записи одного `Write` уходят одним системным вызовом, а не по одному на каждые 16 КБ.
