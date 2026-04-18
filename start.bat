@echo off
setlocal enabledelayedexpansion
cd /d "%~dp0"

echo ===== EventProject: Запуск Docker =====
echo.

REM Проверяем, что Docker доступен
docker info >nul 2>&1
if errorlevel 1 (
    echo [ОШИБКА] Docker не запущен. Запустите Docker Desktop и попробуйте снова.
    pause
    exit /b 1
)

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

echo [1/3] Запуск PostgreSQL и бэкенда...
%COMPOSE_CMD% up -d
if errorlevel 1 (
    echo [ОШИБКА] Не удалось запустить контейнеры.
    pause
    exit /b 1
)

echo [2/3] Ожидание готовности PostgreSQL...
%COMPOSE_CMD% ps postgres | findstr /R /C:"healthy" >nul 2>&1
if errorlevel 1 (
    timeout /t 10 /nobreak >nul
)

echo.
echo ===== Готово! =====
echo.
echo Приложение доступно по адресу: http://localhost:8080
echo pgAdmin (опционально): http://localhost:5050 (admin@admin.com / admin)
echo.
echo Для просмотра логов: %COMPOSE_CMD% logs -f backend
echo Для остановки: stop.bat или %COMPOSE_CMD% down
echo.
pause
