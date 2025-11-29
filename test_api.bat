@echo off
chcp 65001 >nul

:: Создаем папку если ее нет
if not exist "storage\uploads" mkdir "storage\uploads"

:: Переходим в папку storage\uploads
cd storage\uploads

echo ====================================
echo File Translator API Test - FIXED
echo ====================================

echo.
echo 0. Creating user...
curl -X POST http://localhost:5500/users -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"password\":\"password\"}" > user_response.json
type user_response.json

echo.
echo 1. Creating test file...
echo Hello world. This is a test file for translation. > test.txt

echo.
echo 2. Logging in...
curl -X POST http://localhost:5500/sessions -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"password\":\"password\"}" -c cookies.txt > login_response.json
type login_response.json

echo.
echo 3. Uploading file...
curl -X POST http://localhost:5500/private/files -b cookies.txt -F "file=@test.txt" -o upload_response.json

echo.
echo 4. Getting file ID...
for /f "delims=" %%i in ('powershell -Command "(Get-Content upload_response.json | ConvertFrom-Json).id"') do set FILE_ID=%%i
echo File ID: %FILE_ID%

echo.
echo 5. Requesting translation...
curl -X POST http://localhost:5500/private/files/%FILE_ID%/translate -H "Content-Type: application/json" -b cookies.txt -d "{\"source_lang\":\"en\",\"target_lang\":\"ru\"}" -o translate_response.json

echo.
echo 6. Checking translations list...
curl -X GET http://localhost:5500/private/translations -b cookies.txt -o translations_list.json

echo.
echo 7. Downloading translated file...
curl -X GET http://localhost:5500/private/files/%FILE_ID%/download -b cookies.txt -o translated_file.txt

echo.
echo ====================================
echo Test completed!
echo ====================================

echo.
echo Checking results...
type translate_response.json
type translations_list.json

:: Возвращаемся в корневую папку
cd ..\..

echo.
pause