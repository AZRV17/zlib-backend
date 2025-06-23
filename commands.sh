# Создадим каталог для скриптов, если его еще нет
mkdir -p /home/az/GolandProjects/ZLib/ZLib-backend/scripts

# Скопируем скрипт из предыдущего местоположения
cp /home/az/GolandProjects/ZLib/ZLib-backend/scripts/seed_additional_data.sql /home/az/Документы/seed_additional_data.sql

# Выполнение SQL скрипта с использованием аутентификации по паролю
PGPASSWORD=ваш_пароль psql -U postgres -h localhost -d zlib -f /home/az/GolandProjects/ZLib/ZLib-backend/scripts/seed_additional_data.sql

# Если выше не сработало, попробуйте один из этих вариантов:

# Вариант 1: Запустить от имени пользователя postgres
sudo -u postgres psql -d zlib -f /home/az/GolandProjects/ZLib/ZLib-backend/scripts/seed_additional_data.sql

# Вариант 2: Подключиться как текущий пользователь системы (возможно, у вас есть доступ)
psql -d zlib -f /home/az/GolandProjects/ZLib/ZLib-backend/scripts/seed_additional_data.sql

# Вариант 3: Использовать сервис psql
PGPASSWORD=ваш_пароль psql "postgresql://postgres:ваш_пароль@localhost:5432/zlib" -f /home/az/GolandProjects/ZLib/ZLib-backend/scripts/seed_additional_data.sql
