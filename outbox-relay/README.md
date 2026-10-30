# Outbox relay

Самостоятельный шаблон сервиса, который читает PostgreSQL outbox и публикует
события в Kafka. Никаких зависимостей от API или второго Go-модуля.

```sh
go mod download
make check
make run
```

Перед запуском поднимите инфраструктуру и примените миграции:
`make -C .. up migrate`. [Общий сценарий и конфигурация](../README.md).

```text
cmd/outbox-relay          конфигурация, сигналы и exit code
internal/app/start.go    сборка зависимостей и polling loop
internal/app/tasks/
  domain/                envelope и порты Repository / Batch / Publisher
  usecase/               publish ACK → mark published
  repository/postgres/   транзакция и FOR UPDATE SKIP LOCKED
internal/services/kafka/ синхронный publisher, RequireAll, key = aggregate ID
internal/config/         переменные окружения
internal/storage/        подключение PostgreSQL
migrations/              outbox и пример исходной бизнес-таблицы
```

Relay проверяет только envelope, не интерпретирует бизнес-payload. `tasks` —
название демонстрационной фичи; при адаптации его можно заменить на `outbox`.
Несколько реплик пропускают заблокированные строки. Отметки сохраняются только
после успешного batch; при ошибке весь batch откатывается и отправленные события
могут повториться. Повторяется исходный event ID, новый не генерируется.
Ошибки повторяются после poll interval; SIGTERM отменяет текущую операцию и
закрывает подключения. Идемпотентный consumer обязателен.

Изменение исходных данных и INSERT outbox должны быть одной транзакцией —
см. `../scripts/seed.sql`. Из корня этого модуля `make run` использует локальные
настройки стенда; для своего окружения задайте `DATABASE_URL`, `KAFKA_BROKERS`,
`KAFKA_TOPIC`, `BATCH_SIZE`, `POLL_INTERVAL`, `OPERATION_TIMEOUT`.

Подробности гарантий, ограничений порядка, retry и сквозного теста — в
[общем README](../README.md).
