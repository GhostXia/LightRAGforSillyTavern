@echo off
REM ����Ƕ��� Python ������·��
set PYTHON_PATH=.\python\python.exe

REM ����֪ʶͼ��
echo [��Ϣ] ���ڹ���֪ʶͼ��...
"%PYTHON_PATH%" lightrag_api_openai_compatible.py
if errorlevel 1 (
    echo [����] ���з���ʧ�ܣ���鿴�׳��Ĵ��󲢽����Ų飬���߼���������ӡ�
    pause
    exit /b
)

REM ��ʾ���
echo [�ɹ�] ���гɹ����ɿ�ʼ���Ӿƹݣ�
echo [��ʾ] ��ķ���˿�Ϊ���鿴README���ھƹ��������URL
pause