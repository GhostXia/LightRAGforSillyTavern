@echo off
REM ����Ƕ��� Python ������·��
set PYTHON_PATH=.\python\python.exe

REM �������ǰ��
echo [��Ϣ] ����������˷�����...
"%PYTHON_PATH%" ".\lightrag_api_openai_compatible.py"
if errorlevel 1 (
    echo [����] ���к��ʧ�ܣ���鿴�׳��Ĵ��󲢽����Ų顣
    pause
    exit /b
)




REM ��ʾ���
echo [�ɹ�] ���гɹ����ɿ�ʼʹ�ã�
pause