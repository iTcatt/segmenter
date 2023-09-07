# avito-task

`Роутер - chi`

`СУБД - PostgresSQL`

# Запуск
make

./build/user-segments

# Формат запросов 

Запросы выполнял с помощью Postman
## Создание сегментов 

Пример:

`POST localhost:3000/api/segments`
```json
{
    "segments":[
        "AVITO_VOICE_MESSAGES",
        "AVITO_PERFORMANCE_VAS",
        "AVITO_DISCOUNT_30",
        "AVITO_DISCOUNT_50"
    ]
}

```

## Создание пользователей

Пример:

`POST localhost:3000/api/users`

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

## Добавление пользователя в сегменты 

Пример:

`POST localhost:3000/api/update`

```json
{
    "id": 1001,
    "add":[
        "AVITO_VOICE_MESSAGES",
        "AVITO_PERFORMANCE_VAS"
    ],
    "delete":[]
}
```

## Получение сегментов пользователя

user_id передавать как параметры запроса

Пример: 

`GET localhost:3000/api/segments?user_id=1001`