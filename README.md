# avito-task

# Запуск
make

./user-segments

# Формат запросов 

Запросы выполнял с помощью PostMan
## Создание сегментов 

Пример:

POST localhost:3000/api/create

```json:
{
    "segments":[
        "AVITO_VOICE_MESSAGES",
        "AVITO_PERFORMANCE_VAS",
        "AVITO_DISCOUNT_30",
        "AVITO_DISCOUNT_50"
    ]
}

```
