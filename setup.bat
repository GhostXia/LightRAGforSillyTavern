@echo off
REM ����Ƕ��� Python ������·��
set PYTHON_PATH=.\python\python.exe

REM ����Ƿ����Ƕ��� Python ������
if not exist "%PYTHON_PATH%" (
    echo [����] δ�ҵ�Ƕ��� Python ����������ȷ���ѽ� Python ��ȷ��ѹ�� LightRAG/python Ŀ¼��
    pause
    exit /b
)

REM ��� Python �汾
"%PYTHON_PATH%" --version >nul 2>&1
if errorlevel 1 (
    echo [����] δ��⵽��Ч�� Python ��װ������Ƕ��� Python �Ƿ�������
    pause
    exit /b
)

REM ��� requirements.txt �Ƿ����
if not exist "requirements.txt" (
    echo [����] δ�ҵ� requirements.txt �ļ�����ȷ���ļ����ڡ�
    pause
    exit /b
)

REM ʹ�� get-pip.py ��װ pip
    echo [��Ϣ] ����ʹ��Ƕ��� Python ��װ pip...
    "%PYTHON_PATH%" "python\get-pip.py"
    if errorlevel 1 (
        echo [����] pip ��װʧ�ܣ����� get-pip.py �ļ����������ӡ�
        pause
        exit /b
    )

REM ��װ��Ŀ����
echo [��Ϣ] ����ʹ��Ƕ��� Python ��װ��Ŀ���������Ժ�...
"%PYTHON_PATH%" -m pip install --no-deps -r requirements.txt
if errorlevel 1 (
    echo [����] ������װʧ�ܣ����� requirements.txt �ļ����������ӡ�
    pause
    exit /b
)

REM ��װ����ģʽ��
echo [��Ϣ] �����Կɱ༭ģʽ��װ��ǰ��Ŀ...
"%PYTHON_PATH%" -m pip install -e .
if errorlevel 1 (
    echo [����] �ɱ༭ģʽ��װʧ�ܣ����鵱ǰĿ¼�� Python ������
    pause
    exit /b
)

REM ��ʾ���
echo [�ɹ�] ��װ��ɣ���Ŀ��׼�������У�
pause
