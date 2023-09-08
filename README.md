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

## Обновление сегментов пользователя 

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

Пример: 

`GET localhost:3000/api/segments?user_id=1001`

## Удаление сегментов 

`DELETE localhost:3000/api/segments?name=AVITO_VOICE_MESSAGES`

