# AvitoTechTask
Данный сервис решает проблему проведения пользовательских экспериментов внутри компании Avito: хранит данные обо всех проводимых и проведенных тестах, а также о пользвотелях, участвоваших в них.

## Содержание
- [Составные элементы](#составные-элементы)
- [Запуск проекта](#запуск-проекта)
- [Возможности сервиса](#возможности-сервиса)
- [Как устроена база данных](#как-устроена-база-данных)
- [Примеры запросов](#примеры-запросов)
- [Возникшие проблемы](#возникшие-проблемы)
- [Swagger](#swagger)

## Составные элементы
- Веб-сервис на языке Golang
- База данных PostgreSQL
- Два docker-контейнера
- Запросы: JSON-формат
- Конфигурация: ilyakaznacheev/cleanenv
- Логирование: sirupsen/logrus

## Запуск проекта
Все, что нужно для запуска приложения: скачать файл docker-compose и, находясь в директории с этим файлом, запустить в командной строке:

```sh
docker compose up
```
Есть нюанс в работе базы данных: при долгой работе с ней, то есть при отправке более чем 4-х запросов на сохранение данных, БД виснет и запускает checkpoit. 

Пока что не разобралась с причиной данного феномена, но если запрос не выполняется сразу же, то достаточно прервать работу запущенных конетйнеров и, без удаления данных об образах и контейнерах заново запустить приложение (docker compose up). 

Так естественно, добавленные и измененные в БД данные потеряны не будут. Пока что думаю над возможностью исправления данного бага.

При зависании HTTP-запросов:

```sh
Ctrl + C

> docker compose up
```

## Возможности сервиса

### Базовые возможности
- Создание сегмента на основе полученного имени сегмента (POST)
- Удаление сегмента на основе полученного имени сегмента (DELETE)
- Добавление пользователя в сегменты на основе списока имен сегментов и id пользователя (POST)
- Удаление пользователя из сегментов на основе списока имен сегментов и id пользователя (DELETE)
- Получение списка активных сегментов пользваотеля на основе id пользователя (GET)

### Дополнительные возможности

- Сохранение истории попадания/выбывания пользователя из сегмента с возможностью получения отчета по пользователю за определенный период в качестве ссылки на скачивание CSV-файла на основе id пользователя и периода, заданного годом и месяцем (GET)
- Возможность задавать TTL (время автоматического удаления пользователя из сегмента) на основе кол-ва дней его жизни в сегменте (POST)
- Возможность документирования через Swagger: не до конца реализована.

## Как устроена база данных

База данных состоит из двух сущностей: user_in_segment и segment. На рисунке ниже приведена диаграмма для сущностей. 
<p align="center">
<img width="600" alt="Снимок экрана 2023-08-31 в 23 45 20" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/acc90024-37b5-4068-a58a-d09dbcb2f9c0">
</p>

user_in_segment нужен, чтобы хранить пользователя в сегменте, а также дату его добавления в сегмент и удаления из сегмента. Это помогает решать сразу несколько проблем: 
- вытаскивать из БД для пользователя только активные сегменты: те, у которых дата удаления из  сегмента до текущей, либо у которых она вовсе отсутствует (по дефолту имеет значение null)
- задавать TTL: при добавления пользователя в сегменты с помощью POST-запроса, указывается значение поля "period" = кол-ву дней, которое пользователь будет существовать в сегменте, атким образом, когда время его "жизни" истечет, то есть текущая дата будет больше даты удаления, записанной на этапе вставки в БД, этот сегмент не будет выбран как активный
- при удалении пользователя из сегмента меняется лишь его дата удаления - если она null или ее значение больше текущего времени
- прежде чем добавлять пользователя в сегмент, нужно добавить сегмент с его иемнем в таблицу segment, иначе добавление пользователя не будет успешным и не имеет смысла

segment необходим для хранения имен сегментов, сопоставляемых id сегмента(который также используется в составном ключе для user_in_segment), а также для хранения информации о том, активен ли сегмент:
- при удалении сегмента из БД, он не удаляется физически, его поле "active" меняется с дефолтного "true" на "false", а также изменяется user_in_segment: все пользователи из удаленного сегмента должны получить новую дату удаления, если она текущая меньше сущетсвующей, либо если она отсутствует (равна null). Это нужно, чтобы не потерять данные об удалении пользователя из сегмента, которое произошло раньше удаления самого сегмента.

## Примеры запросов

В данном разделе разберем цепочку запорсосв, демонстрирующую логику работы сервиса.

```
Добавим сегмент в его таблицу 
```
```sh
POST
http://localhost:1234/segment/:name

input: 
    name (string in params)
        
output: 
    HTTP-code (int in header of responseWriter)
    answer (string) 
```

![image](https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/582f7cd9-daf2-4c3d-9266-18d20f1396a2)

```
Добавим пользователю 1000 этот сегмент
```
```sh
POST
http://localhost:1234/user/segments
input:
        {
            "user_id":1000,
            "segment_list":["AVITO_1","AVITO_3","AVITO_4"]
        }
output: 
            HTTP-code (int)
            answer (string)
```
![image](https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/68829186-3bb9-412c-8334-61bcf0b5583f)

```
Получим список активных сегментов пользователя 1000
```
```sh
GET
http://localhost:1234/user/:uid
http://localhost:1234/user/1000
input:
        uid (int in params)
output: 
       {
            "user_id":1000,
            "segment_list":["AVITO_1","AVITO_3","AVITO_4"]
        }
```
```
Получим список активных сегментов несущетсвующего в БД пользователя 10001
```
```sh
GET
http://localhost:1234/user/:uid
http://localhost:1234/user/10001
input:
        uid (int in params)
output: 
       {
            "user_id":10001,
            "segment_list":[]
        }
```

```
Удалим несуществующий сегмент с имененм AVITO_25
``` 
```sh
DELETE
http://localhost:1234/segment/:name
http://localhost:1234/segment/AVITO_25
input:
        name (string in params)
output: 
        HTTP-code (int)
        answer (string)
```
```
Получим историю пользователя 1000
```
ниже будет скрин CSV-файла, который автоматически загружается при запуске этого GET-запроса в браузере, POSTMAN не позволяет сразу это увидеть. 
```sh
GET
http://localhost:1234/history/:uid/:year/:month
http://localhost:1234/history/1000/2023/08
input:
        uid (int in params)
        year (int in params)
        month (int in params)
output:
       {
                UserId,SegmentName,Event,EventDate
                1000,AVITO_1,inserted,2023/08/31 19:03:33
                1000,AVITO_3,inserted,2023/08/31 19:03:33
                1000,AVITO_4,inserted,2023/08/31 19:03:33
        }
```

```
Удалим добавленный сегмент пользователя 1000 
``` 
(именно сегмент, а не пользователя из сегмента, то есть поменяется значение "active" для сегмента и значения out_date для пользователей в этом сегменте)
```sh
DELETE
http://localhost:1234/segment/:name
http://localhost:1234/segment/AVITO_1
input:
        name (string in params)
output: 
        HTTP-code (int)
        answer (string)
```

![image](https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/fd135e67-cbca-4976-8ca4-035874b8d0db)

```
Снова получим историю пользователя 1000
``` 
```sh
GET
http://localhost:1234/history/:uid/:year/:month
http://localhost:1234/history/1000/2023/08
input:
        uid (int in params)
        year (int in params)
        month (int in params)
output:
       {
                UserId,SegmentName,Event,EventDate
                1000,AVITO_1,inserted,2023/08/31 19:03:33
                1000,AVITO_3,inserted,2023/08/31 19:03:33
                1000,AVITO_4,inserted,2023/08/31 19:03:33
                1000,AVITO_3,deleted,2023/08/31 19:04:22
        }
```

```
Снова получим список активных сегментов пользователя 1000
```
```sh
GET
http://localhost:1234/user/:uid
http://localhost:1234/user/1000
input:
        uid (int in params)
output: 
       {
            "user_id":1000,
            "segment_list":["AVITO_1","AVITO_4"]
        }
```

```
Проделаем аналогичные операции, только с бОльшим кол-вом пользователей и сегментов
``` 
Кроме того, посмотрим на историю так, как она выглядит в скачанном CSV-файле, а не в JSON-ответе, полученном из ResponseWriter-а.
```
Данные о сегментах и пользователях уже добавлены в соответствующие таблиццы
```
приведены последовательные запросы:
![image](https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/c9062b56-49c7-4c60-946b-580460c2da7d)
![image](https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/79fd16c9-48ea-44ca-b20e-2cc70c17cd77)
![image](https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/c731fa82-aac5-4c5e-ae66-5368e99ac1f8)
![image](https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/db55cbb4-c1b4-47c1-bde4-02346d3b4576)

Как выглядят данные в скачанном СSV-файле::

<img width="348" alt="Снимок экрана 2023-08-31 в 23 31 06" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/bf4e5337-ff38-494d-a52a-402b4d50167f">


## Возникшие проблемы
1. Проблема со Swagger: НЕ СУЩЕСТВУЕТ обертки над swaggo/swager для разработки ручки на основе роутера от julienschmidt/httprouter. Однако все ручки задокументированы и сгенерированы yaml и json файлы для работы swagger-а, демонстрация полученных локально страниц в разделе "Swagger". Планирую исправить проблему подключения к Swagger в ближайшее время.
2. Есть нюанс в работе базы данных: при долгой работе с ней, то есть при отправке более чем 4-х запросов на сохранение данных, БД виснет и запускает checkpoint -- требуется просто остановить контейнеры и без удаления их и имеджей запустить заново docker compose.
3. Не успела качественно поработать с паролями для БД: сейчас они торчат наружу в конфигурации, необходимо аккуратно спрятать их в перемнных окружения, начала эту кропотливую работу, но еще не завершила.

## Swagger

<img width="1375" alt="Снимок экрана 2023-08-31 в 23 33 18" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/61105804-4b0c-4ed0-a327-82f16d59469f">
<img width="1376" alt="Снимок экрана 2023-08-31 в 23 34 01" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/6067f50f-3b9d-4678-9729-2e9d5a17af83">
подробнее запросы:
<img width="1374" alt="Снимок экрана 2023-08-31 в 23 34 28" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/8b8d6449-7976-476b-b8a7-5696895c1942">
<img width="1371" alt="Снимок экрана 2023-08-31 в 23 34 51" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/3c34beb5-281d-4c4f-89de-ec56f26ca953">
<img width="1371" alt="Снимок экрана 2023-08-31 в 23 35 39" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/d000bd6b-8a76-4a9f-b7e1-4874ed7ad589">
<img width="1366" alt="Снимок экрана 2023-08-31 в 23 36 02" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/f6e65dc0-60bc-485c-9710-9c9df6612397">
<img width="1371" alt="Снимок экрана 2023-08-31 в 23 36 18" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/8585e014-1fd7-469a-977f-a19a4a1baf44">
<img width="1376" alt="Снимок экрана 2023-08-31 в 23 36 33" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/dfec332b-47dc-426b-9263-1c28772bcc6d">
<img width="1373" alt="Снимок экрана 2023-08-31 в 23 36 50" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/2c658e7d-b631-44e7-b484-a01968922369">
<img width="1373" alt="Снимок экрана 2023-08-31 в 23 37 15" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/71bc8aad-de4d-4028-ab37-75313885dd72">

<img width="1368" alt="Снимок экрана 2023-08-31 в 23 37 33" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/5db10730-51b4-41d5-97da-ba4d3566e579">
<img width="1366" alt="Снимок экрана 2023-08-31 в 23 37 51" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/f4e319e7-6a33-402b-85c4-0df6810222e3">
<img width="1368" alt="Снимок экрана 2023-08-31 в 23 38 17" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/a4119614-3e51-4e90-9c71-77daa36df792">
<img width="1377" alt="Снимок экрана 2023-08-31 в 23 38 30" src="https://github.com/LinaAngelinaG/AvitoTechTask/assets/61655484/623a1f42-ead6-42ec-9fc0-bb5d63a157cd">




