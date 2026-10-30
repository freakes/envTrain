# Outbox consumer

Самостоятельный шаблон Kafka-консьюмера: получил событие → выполнил бизнес-логику
в PostgreSQL → подтвердил Kafka offset. Читает envelope, опубликованный relay,
но не зависит от его Go-модуля.

```sh
go mod download
make check
make run
```

Перед запуском: `make -C .. up migrate`.
[Общий сценарий и конфигурация](../README.md).

```text
cmd/outbox-consumer       конфигурация, сигналы и exit code
internal/app/start.go    сборка зависимостей и Kafka reader
internal/app/tasks/
  domain/                event, Task, контракты и валидация
  delivery/kafka/        Fetch → Handle → CommitMessages
  usecase/               проверка события и запуск бизнес-сценария
  repository/postgres/   inbox + изменение projection в одной транзакции
internal/config/         переменные окружения
internal/storage/        подключение PostgreSQL
migrations/              inbox и демонстрационная projection
```

Бизнес-пример — обновление `task_projection`. Уникальное событие увеличивает
`updates`; повтор с тем же ID ничего не меняет. Inbox INSERT и projection
фиксируются одной транзакцией. Конкурирующие дубли сериализуются уникальным
ключом inbox. При DB-ошибке inbox также откатывается и событие можно повторить.

Offset подтверждается синхронно только после успешной DB-транзакции. При ошибке
обработки/commit сервис завершается до чтения следующего сообщения. Это важно:
commit более позднего offset мог бы пропустить неуспешное сообщение той же
partition. При перезапуске возможен дубль, который отсекает inbox.

Пример понимает `task.created.v1`; неизвестные типы, повреждённый JSON и ошибки
валидации приводят к остановке без намеренного пропуска события. DLQ и retry
policy добавляйте явно под свой контракт. Для запуска задайте `DATABASE_URL`,
`KAFKA_BROKERS`, `KAFKA_TOPIC`, `KAFKA_GROUP`, `OPERATION_TIMEOUT` либо используйте
локальные значения по умолчанию.

Inbox не делает внешние HTTP-вызовы идемпотентными. Гарантии, ограничения,
конфигурация и сквозной тест описаны в [общем README](../README.md).
