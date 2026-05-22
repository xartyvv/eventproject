#!/usr/bin/env sh
cd "$(dirname "$0")" || exit 1

printf "===== EventProject: Запуск Docker =====\n\n"

if ! command -v docker >/dev/null 2>&1; then
  printf "[ОШИБКА] Docker не найден. Установите Docker Desktop и запустите его.\n"
  exit 1
fi

docker info >/dev/null 2>&1
if [ $? -ne 0 ]; then
  printf "[ОШИБКА] Docker не запущен. Запустите Docker Desktop и попробуйте снова.\n"
  exit 1
fi

COMPOSE_CMD=""
if docker-compose version >/dev/null 2>&1; then
  COMPOSE_CMD="docker-compose"
elif docker compose version >/dev/null 2>&1; then
  COMPOSE_CMD="docker compose"
else
  printf "[ОШИБКА] Не найден docker-compose или docker compose.\n"
  exit 1
fi

printf "[1/3] Запуск PostgreSQL и бэкенда...\n"
if ! $COMPOSE_CMD up -d; then
  printf "[ОШИБКА] Не удалось запустить контейнеры.\n"
  exit 1
fi

printf "[2/3] Ожидание готовности PostgreSQL...\n"
if ! $COMPOSE_CMD ps postgres | grep -q "healthy" >/dev/null 2>&1; then
  printf "Подождите 10 секунд...\n"
  sleep 10
fi

printf "\n===== Готово! =====\n\n"
printf "Приложение доступно по адресу: http://localhost:8080\n"
printf "pgAdmin (опционально): http://localhost:5050 (admin@admin.com / admin)\n\n"
printf "Для просмотра логов: %s logs -f backend\n" "$COMPOSE_CMD"
printf "Для остановки: ./stop.sh или %s down\n" "$COMPOSE_CMD"
