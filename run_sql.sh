#!/bin/bash
# Создаём временный скрипт для выполнения SQL файла

# Вариант с sudo
sudo -u postgres psql -d zlib -f /home/az/GolandProjects/ZLib/ZLib-backend/scripts/seed_additional_data.sql

# Делаем скрипт исполняемым
chmod +x /home/az/GolandProjects/ZLib/ZLib-backend/run_sql.sh

# Запуск скрипта
./run_sql.sh
