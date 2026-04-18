@echo off
setlocal enabledelayedexpansion
cd /d "%~dp0"

echo ===== EventProject: Остановка =====
echo.

REM Выбираем доступную команду Compose
docker-compose version >nul 2>&1
if errorlevel 0 (
    set "COMPOSE_CMD=docker-compose"
) else (
    docker compose version >nul 2>&1
    if errorlevel 0 (
        set "COMPOSE_CMD=docker compose"
    ) else (
        echo [ОШИБКА] Не найден docker-compose или docker compose.
        pause
        exit /b 1
    )
)

%COMPOSE_CMD% down
if errorlevel 1 (
    echo [ОШИБКА] Не удалось остановить контейнеры.
    pause
    exit /b 1
)

echo.
echo Контейнеры остановлены.
echo Для удаления данных БД: %COMPOSE_CMD% down -v
echo.
pause
