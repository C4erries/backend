#!/bin/bash

# Логирование
LOG_FILE="./startup.log"
exec > >(tee -a "$LOG_FILE") 2>&1

echo "=== Запуск скрипта: $(date) ==="

# Устанавливаем необходимые пакеты
echo "Устанавливаю необходимые пакеты..."
sudo apt-get update
sudo apt-get install -y docker.io docker-compose git

# Путь к проекту (относительно расположения скрипта)
PROJECT_DIR="./backend3"

# Проверяем, существует ли директория проекта
if [ -d "$PROJECT_DIR" ]; then
  echo "Директория проекта $PROJECT_DIR существует. Проверяю обновления..."
  cd "$PROJECT_DIR"

  # Сохраняем текущий хэш коммита
  OLD_COMMIT=$(git rev-parse HEAD)

  # Обновляем репозиторий
  echo "Выполняю git pull..."
  git pull origin main

  # Сравниваем хэши коммитов
  NEW_COMMIT=$(git rev-parse HEAD)
  if [ "$OLD_COMMIT" == "$NEW_COMMIT" ]; then
    echo "Изменений нет. Запускаю существующие контейнеры..."
    sudo docker-compose up -d
  else
    echo "Обнаружены изменения. Пересобираю и запускаю контейнеры..."
    sudo docker-compose up -d --build
  fi
else
  echo "Директория проекта $PROJECT_DIR не существует. Клонирую репозиторий..."
  git clone https://github.com/c4erries/backend.git "$PROJECT_DIR"
  cd "$PROJECT_DIR"
  echo "Запускаю сборку, а затем контейнеры..."
  sudo docker-compose up -d --build
fi

echo "=== Завершение скрипта: $(date) ==="
