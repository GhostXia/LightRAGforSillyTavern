@echo off
REM ����Ƕ��� Python ������·��
set PYTHON_PATH=.\python\python.exe

set ENV_PATH=.\.env

REM ��������ǰ��
echo [��Ϣ] ������������ǰ����...
"%PYTHON_PATH%" ".\Gradio_web.py"
if errorlevel 1 (
    echo [����] ����ǰ��ʧ�ܣ���鿴�׳��Ĵ��󲢽����Ų顣
    pause
    exit /b
)

REM ��ʾ���
echo [�ɹ�] ���гɹ����ɿ�ʼʹ�ã�
pause