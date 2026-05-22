#!/usr/bin/env sh
cd "$(dirname "$0")" || exit 1

printf "===== EventProject: Остановка =====\n\n"

COMPOSE_CMD=""
if docker-compose version >/dev/null 2>&1; then
  COMPOSE_CMD="docker-compose"
elif docker compose version >/dev/null 2>&1; then
  COMPOSE_CMD="docker compose"
else
  printf "[ОШИБКА] Не найден docker-compose или docker compose.\n"
  exit 1
fi

if ! $COMPOSE_CMD stop; then
  printf "[ОШИБКА] Не удалось остановить контейнеры.\n"
  exit 1
fi

printf "\nКонтейнеры остановлены. Данные остаются в том же состоянии.\n"
printf "Для повторного запуска используйте ./start.sh или %s start\n" "$COMPOSE_CMD"
printf "Для полного удаления контейнеров и данных БД: %s down -v\n" "$COMPOSE_CMD"
