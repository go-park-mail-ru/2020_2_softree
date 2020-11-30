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

2. Поднять базу:
    ```
    docker-compose up -d --build postgres
    ```

3. Из корня проекта запустить выполнить команду:
    ```
    make run
    ```
