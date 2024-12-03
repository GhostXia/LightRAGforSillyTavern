@echo off
echo ��ӭʹ�û����������ù��ߣ�
echo Ŀǰ��ֱ��֧��OpenAIģ�ͣ�����ģ����ʹ����תϵͳ��֧��OpenAI��׼��ʽ����Ӧ
echo ��Ҫ����ĳ�������Ļ���������skip��ֱ�ӻس���������ַ�������
echo.

:: Step 1: ȷ��û���ظ���� file_DIR
setlocal enabledelayedexpansion
set /p file_DIR="): "�������ʹ�õ��ı��ļ���·�� (Ĭ�� ./text/book.txt������ skip ����:)"
if "%file_DIR%"=="skip" (
    echo �������� file_DIR
) else (
    if "%file_DIR%"=="" set file_DIR=./text/book.txt
    echo ���ڸ��� file_DIR ��������...
    findstr /v "file_DIR=" .env > .env.tmp
    echo file_DIR=%file_DIR% >> .env.tmp
    move /Y .env.tmp .env
    echo �������� file_DIR ������Ϊ��%file_DIR%
)

:: Step 2: ȷ��û���ظ���� RAG_DIR
set /p rag_DIR="������֪ʶͼ���ļ���·�� (Ĭ�� ./file/test������ skip ����:) "
if "%rag_DIR%"=="skip" (
    echo �������� RAG_DIR
) else (
    if "%rag_DIR%"=="" set rag_DIR=./file/test
    echo ���ڸ��� RAG_DIR ��������...
    findstr /v "RAG_DIR=" .env > .env.tmp
    echo RAG_DIR=%rag_DIR% >> .env.tmp
    move /Y .env.tmp .env
    echo �������� RAG_DIR ������Ϊ��%rag_DIR%
)

:: Step 3: �����û����� OPENAI_BASE_URL
set /p OPENAI_BASE_URL="������ OpenAI API ���� URL (���� skip ����): "
if "%OPENAI_BASE_URL%"=="skip" (
    echo �������� OPENAI_BASE_URL
) else (
    echo ���ڸ��� OPENAI_BASE_URL ��������...
    findstr /v "OPENAI_BASE_URL=" .env > .env.tmp
    echo OPENAI_BASE_URL=%OPENAI_BASE_URL% >> .env.tmp
    move /Y .env.tmp .env
    echo �������� OPENAI_BASE_URL ������Ϊ��%OPENAI_BASE_URL%
)

:: Step 4: �����û����� OPENAI_API_KEY
set /p OPENAI_API_KEY="������ OpenAI API ��Կ (���� skip ����): "
if "%OPENAI_API_KEY%"=="skip" (
    echo �������� OPENAI_API_KEY
) else (
    echo ���ڸ��� OPENAI_API_KEY ��������...
    findstr /v "OPENAI_API_KEY=" .env > .env.tmp
    echo OPENAI_API_KEY=%OPENAI_API_KEY% >> .env.tmp
    move /Y .env.tmp .env
    echo �������� OPENAI_API_KEY ������Ϊ��%OPENAI_API_KEY%
)

:: Step 5: �����û����� LLM_MODEL
set /p LLM_MODEL="������ LLM ģ�� (���� skip ����): "
if "%LLM_MODEL%"=="skip" (
    echo �������� LLM_MODEL
) else (
    echo ���ڸ��� LLM_MODEL ��������...
    findstr /v "LLM_MODEL=" .env > .env.tmp
    echo LLM_MODEL=%LLM_MODEL% >> .env.tmp
    move /Y .env.tmp .env
    echo �������� LLM_MODEL ������Ϊ��%LLM_MODEL%
)

:: Step 6: �����û����� EMBEDDING_MODEL
set /p EMBEDDING_MODEL="������Ƕ��ģ�� (���� skip ����): "
if "%EMBEDDING_MODEL%"=="skip" (
    echo �������� EMBEDDING_MODEL
) else (
    echo ���ڸ��� EMBEDDING_MODEL ��������...
    findstr /v "EMBEDDING_MODEL=" .env > .env.tmp
    echo EMBEDDING_MODEL=%EMBEDDING_MODEL% >> .env.tmp
    move /Y .env.tmp .env
    echo �������� EMBEDDING_MODEL ������Ϊ��%EMBEDDING_MODEL%
)

:: Step 7: �����û����� EMBEDDING_MAX_TOKEN_SIZE
set /p EMBEDDING_MAX_TOKEN_SIZE="������Ƕ��ģ�����tokens���� (Ĭ�� 2046������ skip ����): "
if "%EMBEDDING_MAX_TOKEN_SIZE%"=="skip" (
    echo �������� EMBEDDING_MAX_TOKEN_SIZE
) else (
    if "%EMBEDDING_MAX_TOKEN_SIZE%"=="" set EMBEDDING_MAX_TOKEN_SIZE=2046
    echo ���ڸ��� EMBEDDING_MAX_TOKEN_SIZE ��������...
    findstr /v "EMBEDDING_MAX_TOKEN_SIZE=" .env > .env.tmp
    echo EMBEDDING_MAX_TOKEN_SIZE=%EMBEDDING_MAX_TOKEN_SIZE% >> .env.tmp
    move /Y .env.tmp .env
    echo �������� EMBEDDING_MAX_TOKEN_SIZE ������Ϊ��%EMBEDDING_MAX_TOKEN_SIZE%
)

:: Step 8: �����û����� API_port
set /p API_port="���������˿� (Ĭ�� 8020������ skip ����): "
if "%API_port%"=="skip" (
    echo �������� API_port
) else (
    if "%API_port%"=="" set API_port=8020
    echo ���ڸ��� API_port ��������...
    findstr /v "API_port=" .env > .env.tmp
    echo API_port=%API_port% >> .env.tmp
    move /Y .env.tmp .env
    echo �������� API_port ������Ϊ��%API_port%
)

echo.
echo ���л���������������ϣ�
echo ������Setup Your Graph.bat�Խ���֪ʶͼ�׵Ĺ�����������ϻ����Ѿ���֪ʶͼ�����������Run API.bat
