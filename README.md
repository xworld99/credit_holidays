# credit_holidays

## Запуск
Перед запуском необходимо установить docker, docker-compose, и убедиться, что порты 8080 и 5431 не заняты.
1. Клонировать репозиторий `git clone git@gitlab.com:ses011/credit_holidays.git`.
2. Перейти в репозиторий `cd credit_holidays`.
3. Для запуска ввести команду `make up`.
4. Для остановки ввести команду `make down`.


## TODO
1. local tests, tests -> gomock
2. squirrel instead of raw sql
3. API directory for swagger docs 
4. Postgres config
5. converting money to float64
