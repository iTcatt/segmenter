# avito-task

# Запуск
make

./build/user-segments

# Формат запросов 

Запросы выполнял с помощью PostMan
## Создание сегментов 

Пример:

POST localhost:3000/api/segments
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

POST  localhost:3000/api/users

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

POST localhost:3000/api/update

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
