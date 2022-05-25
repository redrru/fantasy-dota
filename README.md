#### Разработка

#### Makefile

Основные команды:

1. `make test` - запустить тесты
2. `make style` - запустить gofmt, gomod, golangci-lint
3. `make up` - собрать и запустить проект в докере
4. `make down` - остановить проект
5. `make codegen` - сгенерировать сервер из openapi
6. `make run-compile-daemon` - запустить проект под CompileDaemon

#### Запуск

Поднять проект в докере `make up`.

В проекте используется live-reload на основе CompileDaemon. После локальных изменений файлов проекта, в докере запустится новая сборка проекта.

##### Fetcher

Для переодического опроса url есть либа [Fetcher](https://github.com/redrru/fantasy-dota/blob/master/pkg/fetcher/fetcher.go).

Чтобы добавить новый фетчер необходимо:
1. Создать новый класс в `/internal/fantasy-dota/fetchers/`, который будет удовлетворять интерфейсу [Handler](https://github.com/redrru/fantasy-dota/blob/f7467a4bdd7d8168e7399108bc7220a0a81b58ff/pkg/fetcher/fetcher.go#L15), [Пример](https://github.com/redrru/fantasy-dota/blob/f7467a4bdd7d8168e7399108bc7220a0a81b58ff/internal/fantasy-dota/fetchers/example.go#L18).
2. Зарегестрировать обработчик в приложении, [пример](https://github.com/redrru/fantasy-dota/blob/master/cmd/fantasy-dota/main.go):
    ```go
    app.RegisterFetchers(fetchers.NewExample())
    ```

#### Http

Для обработки http запросов используется роутер [echo](https://github.com/labstack/echo), дефолтный порт 8080ю

Чтобы добавить новый обработчик необходимо:

1. Добавить описание в [openapi.yaml](https://github.com/redrru/fantasy-dota/blob/master/api/http/openapi.yaml)
2. Сгенерировать сервер `make codegen`
3. Добавить метод в класс [Server](https://github.com/redrru/fantasy-dota/blob/f7467a4bdd7d8168e7399108bc7220a0a81b58ff/internal/gateways/http/server.go#L9), [пример](https://github.com/redrru/fantasy-dota/blob/master/internal/gateways/http/example.go)

#### Трассировка запросов

Для трассировки используется [OpenTelemetry](https://opentelemetry.io) + [Jaeger](https://www.jaegertracing.io).

При запуске проекта поднимется контейнер `jaeger`, gui доступен по `http://localhost:16686`

В ответ на любой http запрос будет добавлен хедер `trace_id`, по нему можно найти весь трейс в gui jaeger.

```bash
# curl -v localhost:8080/example
*   Trying 127.0.0.1:8080...
...
< HTTP/1.1 500 Internal Server Error
< Content-Type: application/json; charset=UTF-8
< Trace_id: f39e2c20f7176eda04ca9ca3cd3d05af <------ Trace_id
< Date: Mon, 23 May 2022 10:11:33 GMT
< Content-Length: 36
```
