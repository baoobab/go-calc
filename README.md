# Калькулятор примитивных арифметических выражений

## Инфо

Работает по `http`, принимает входящие запросы на эндпоинт `/api/v1/calculate`.
Стандартный хост и порт: `localhost:8080`

## Структура входящих запросов

**POST** запрос. Пример тела запроса (структура `type CalcRequest`):
<br>`{ expression: "2+2/4" }`

## Структура ответов

- При успешной обработке
<br>Статус 200, объект с полем `result` (структура `type CalcResponse`).
<br>Пример успешного ответа: `{ "result": 2.5 } `

- Иначе 
<br>Статус 422 или 500, объект с полем `error` (структура `type CalcResponse`).
<br> Пример неудачного ответа: `{ "error": "Expression is not valid" } `

## Разрешённые символы

- Знаки операций (только бинарные): `+`, `-`, `*`, `/`
- Знаки приоритизации: `(`, `)`
- Рациональные числа (либо целые, либо через `.`)

## Подготовка

1. Проверьте наличие языка Go на устройстве, и всех необходимых пакетов
2. ```git clone https://github.com/baoobab/go-calc```

## Билд

```
go build
```

## Запуск

```
go run main.go
```

## Примеры взаимодействия через curl

1. **Корректный** запрос:
    ```
    curl --location 'localhost:8080/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{
      "expression": "2+2*2"
    }'
    ```
   Ответ: `{ "result": 6 }` Статус: `200`

2. **Некорректный** запрос:
    ```
    curl --location 'localhost:8080/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data '{
      "expression": "bibi bebe"
    }'
    ```
   Ответ: `{ "error": "Expression is not valid" }` Статус `422`

3. Внутренняя ошибка сервера (запрос роли не играет):
   ```
    curl --location 'localhost:8080/api/v1/calculate' \
    --header 'Content-Type: application/json' \
    --data ''
    ```
   Ответ: `{ "error": "Internal server error" }` Статус `500`
