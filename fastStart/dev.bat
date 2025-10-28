@echo off
chcp 65001
title Go Server Auto-Reload

echo 🔥 Запуск Go сервера с авто-перезагрузкой...
echo.

:start
echo 📅 [%time%] Запуск сервера...
go run main.go
echo.
echo ⚠️  Сервер остановлен. Перезапуск через 3 секунды...
timeout /t 3 >nul
goto start