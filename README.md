# osiris
---

[![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](http://forthebadge.com)


# Микросервисы
## Auth Service 
+ OAuth2
+ Генерирует токен для доступа к микросервисам
### БД
1. id
2. first_name
3. last_name

## Alert Service 
+ создает тревоги 
+ получает тревоги
+ удаляет тревоги
Записывает в БД и передает по брокеру сообщений(alerts) следующему сервису
### БД
4. id
5. user_id
6. title
7. description
8. type
9. location
10. radius
11. status
12. created_at
13. updated_at

## Points Service 
+ CRUD точек
### БД
14. id
15. user_id
16. title
17. location
18. radius
19. created_at
20. updated_at

## Processing Service 
+ Обработка запросов
Получается нужные точки и тревогу из (alerts) 
Генерирует запросы и отсылает в брокер сообщений(notification) 

## Notification Service 
+ Отправка уведомлений
Получает запросы из брокера(notification) и отправляет сообщения