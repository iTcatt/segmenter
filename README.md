# segmenter

Требуется реализовать сервис, хранящий пользователя и сегменты, в которых он состоит
(создание, изменение, удаление сегментов, а также добавление и удаление
пользователей в сегменты)

# Запуск

В файле ./configs/config.yaml установить нужные значения или оставить их по умолчанию.

```bash
make doker-build
make docker-run
```

# Формат запросов 

Запросы выполнял с помощью Postman

## Создание сегментов 

Условие задания не запрещало запросы, в которых добавляемый сегмент повторяется несколько раз.
Я решил, что на этот сегмент будет выводится только один ответ, который был первым.

* если сегмент был создан, то вернется сообщение `created`.

* если сегмент уже был добавлен ранее, то вернется сообщение `already exist`.

* если при создании сегмента возникнет ошибка, то вернется сообщение `not created`.

### Пример запроса:

`POST localhost:3000/api/segment`

```json
{
  "segments":[
      "AVITO_VOICE_MESSAGES",
      "AVITO_PERFORMANCE_VAS",
      "AVITO_DISCOUNT_30",
      "AVITO_VOICE_MESSAGES"
    ]
}
```

### Пример ответа от сервера:

```json
{
    "AVITO_DISCOUNT_30": "created",
    "AVITO_PERFORMANCE_VAS": "created",
    "AVITO_VOICE_MESSAGES": "created"
}
```

### Пример запроса:

`POST localhost:3000/api/segment`

```json
{ 
    "segments":[
        "AVITO_VOICE_MESSAGES",
        "AVITO_DISCOUNT_30",
        "AVITO_DISCOUNT_30",
        "AVITO_DISCOUNT_50"
    ]
}
```

### Ответ от сервера:
```json
{ 
    "AVITO_DISCOUNT_30": "already exist",
    "AVITO_DISCOUNT_50": "created",
    "AVITO_VOICE_MESSAGES": "already exist"
}
```

## Создание пользователей

Схема ответа такая же, как при создании сегмента. 

* если пользователь будет создан, то вернется `created`.
* если пользователь уже был создан ранее, то вернется `already exist`
* если при создании пользователя возникнет ошибка, то вернется `not created`

### Пример запроса:

`POST localhost:3000/api/user`

```json
{
    "users":[
        32,
        64, 
        128, 
        128
    ]
}
```

### Ответ от сервера:

```json
{
    "128": "created",
    "32": "created",
    "64": "created"
}
```

### Пример запроса:

`POST localhost:3000/api/user`

```json
{
    "users":[
      32,
      128,
      256
    ]
}
```

### Ответ от сервера:

```json
{
    "128": "already exist",
    "256": "created",
    "32": "already exist"
}
```
    
## Обновление сегментов пользователя 
ID пользователя передается как URL параметр.
Eсли такого пользователя несуществует, вернется код ответа `404`

В body запроса передаются два поля: 
1) add_segments - список названий сегментов, в которых нужно добавить пользователя
2) delete_segments - список названий сегментов, из которых нужно удалить пользователя.

В ответе возвращается структура пользователя: его ID и сегменты, в которых он состоит после проведенных изменений.

### Пример запроса:

`PATCH localhost:3000/api/user/32`

```json
{
    "add_segments":[
        "AVITO_VOICE_MESSAGES",
        "AVITO_PERFORMANCE_VAS"
    ],
    "delete_segments":[
        "NO_SEGMENT"
    ]
}
```

### Ответ от сервера:

```json
{
    "id": 32,
    "segments": [
        "AVITO_PERFORMANCE_VAS",
        "AVITO_VOICE_MESSAGES"
    ]
}
```

### Пример запроса:

`PATCH localhost:3000/api/user/10000`

```json
{
    "add":[
        "AVITO_VOICE_MESSAGES",
        "AVITO_PERFORMANCE_VAS"
    ],
    "delete":[
        "NO_SEGMENT"
    ]
}
```

### Ответ от сервера: 

Вернется код ответа `404 NotFound`

### Пример запроса:

`PATCH localhost:3000/api/user/32`

```json
{
    "add":[
        "AVITO_DISCOUNT_30"
    ],
    "delete":[
        "AVITO_PERFORMANCE_VAS"
    ]
}
```

### Ответ от сервера:

```json
{
    "id": 32,
    "segments": [
        "AVITO_VOICE_MESSAGES",
        "AVITO_DISCOUNT_30"
    ]
}
```

## Получение сегментов пользователя

ID пользователя передается через URL параметры. 

* если id не число, то вернется код ответа `400` и сообщение об ошибке
* если такой пользователь существует, то вернется список его сегментов
* если такого пользователя не существует, то вернется код ответа `404`

### Пример запроса: 

`GET localhost:3000/api/user/32`

### Ответ от сервера:
```json
{
    "id": 32,
    "segments": [
        "AVITO_VOICE_MESSAGES",
        "AVITO_DISCOUNT_30"
    ]
}
```

### Пример некорректного запроса: 

`GET localhost:3000/api/user/number`

Такой запрос вернет код ответа `400 - Bad Request`

### Пример для пользователя не входящего ни в один сегмент:

`GET localhost:3000/api/segments/64`

### Ответ от сервера:

```json
{
    "id": 64,
    "segments": []
}
```

### Пример для несозданного пользователя: 

`GET localhost:3000/api/user/100000`

### Ответ от сервера:

Вернется код ответа `404`

## Удаление сегмента

* если сегмент существует, то после удаления будет возвращен код ответа `204`

* если сегмент отсутствует, то вернется код ответа `404`

### Пример удаления существующего сегмента:

`DELETE localhost:3000/api/segment/AVITO_VOICE_MESSAGES`

Такой запрос вернет код ответа `204`. 

### Пример удаления несозданного сегмента:

`DELETE localhost:3000/api/segment/AVITO_VOICE_MESSAGES`

Т.к. на предыдущем шаге сегмент был удален, то вернется код ответа `404`.

## Удаление пользователя 

* если пользователь существует, то после удаления будет возвращен код ответа `204`.

* если пользователь отстутствует, то вернется код ответа `404`.

* если в параметрах запроса передали не число, то вернется код ответа `400`, а также сообщение об ошибке.

### Пример успешного удаления:

`DELETE localhost:3000/api/user/128`

Такой запрос вернет код ответа `204`.

### Пример удаления несозданного пользователя:

`DELETE localhost:3000/api/user/128`

Т.к. на предыдущем шаге пользователь был удален, то вернется код ответа `404`.

### Пример неправильного запроса:

`DELETE localhost:3000/api/user/A`
```json
{
    "error": "validation error"
}
```
Такой запрос вернет код ответа `400`.
