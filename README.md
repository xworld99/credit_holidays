# credit_holidays

## Запуск
Перед запуском необходимо установить docker, docker-compose, и убедиться, что порты 8080 и 5431 не заняты.
1. Клонировать репозиторий `git clone git@gitlab.com:ses011/credit_holidays.git`.
2. Перейти в репозиторий `cd credit_holidays`.
3. Для запуска ввести команду `make up`.
4. Для остановки ввести команду `make down`.
5. Тесты запускаются командой `make test`.

## Описание

Документация, при запущенном приложении, доступна [тут](http://0.0.0.0:8080/swagger/index.html)
Пользователь (`user`) создается при первом изменении его баланса, имеет два типа баланса: обычный и зарезервированный.

Под услугой (`service`) понимается любое действие, изменяющее баланс пользователя. Услуга имеет две специфичные характеристики:
1. service_type - тип услуги (`accrual` - зачисление средств (увеличение баланса пользователя), `withdraw` - списание средств (уменьшение баланса)).
2. confirmation_needed - нужно ли услугу дополнительно подтверждать (например, услуга "зачисление средств на баланс пользователя" не требует подтверждение, средства начисляются сразу на активный баланс пользователя, услуга "вывод средств на внешний источник" требует подтверждения, поэтому изначально взаимодействует с зарезервированным балансом).

Любой запрос пользователя на услугу называется заказом (`order`).

Отчет (`report`) для бухгалтерии является csv файлом, сохраняемом на сервере. В ответ на запрос на создание отчета возвращается его название, отчет по названию можно получить по URL'у статической директории.

Используемая для рассчетов валюта - условная единица, равная 1/100 от принятой в регионе валюты (конкретно в моей реализации для рассчетов используется сотая часть рубля, то есть копейка). Это необходимо учитывать при работе с балансом пользователя и создаваемыми заказами.