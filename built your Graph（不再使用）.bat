@echo off
REM ����Ƕ��� Python ������·��
set PYTHON_PATH=.\python\python.exe

REM ����֪ʶͼ��
echo [��Ϣ] ���ڹ���֪ʶͼ��...
"%PYTHON_PATH%" ".\built your Graph.py"
if errorlevel 1 (
    echo [����] ����֪ʶͼ��ʧ�ܣ���鿴�׳��Ĵ��󲢽����Ų飬���߼���������ӡ�
    pause
    exit /b
)

REM ��ʾ���
echo [�ɹ�] ������ɣ��ɿ�ʼ���з���
pause