# Биржа валют
Проект, который позволяет торговать валютами  

https://softree.group

# Менторы
Фронт - [Джахонгир Тулфоров](https://github.com/bin-umar)  
Бек - [Алексей Ларютин](https://github.com/ALarutin)

# Команда
- Фронт
  * [Быховец Марк](https://github.com/mark-by)
  * [Мурадян Артем](https://github.com/MuradyanArtem)
- Бек
  * [Трущелев Михаил](https://github.com/ThePsina)
  * [Макаров Тимофей](https://github.com/timofef)

# Как запустить сервер
1. Установить ПО
    * Docker
    * Docker-compose (можно через venv)

2. Из корня проекта выполнить команду:
    ```
    docker-compose up -d --build
    ```

3. Поднять базу из бекапа  
    * Получить id контейнера
    ```
    docker ps
    ```
    * Восстановить из бекапа
    ```
    docker exec <номер контейнера с бд> psql -U user -d moneycat -f /backup.sql
    ```

