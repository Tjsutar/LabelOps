@echo off
setlocal

REM Check if ZPL file path is provided
if "%1"=="" (
    echo Error: No ZPL file path provided
    echo Usage: print_label.bat "path\to\file.zpl"
    pause
    exit /b 1
)

REM Check if file exists
if not exist "%1" (
    echo Error: ZPL file not found: %1
    pause
    exit /b 1
)

echo Printing ZPL file: %1
echo.

REM Option 1: Try to copy to LPT1 (parallel port)
echo Attempting to print via LPT1...
copy /B "%1" LPT1 >nul 2>&1
if %errorlevel%==0 (
    echo Successfully sent to LPT1
    goto :success
)

REM Option 2: Try to copy to a common USB printer port
echo LPT1 failed, attempting USB001...
copy /B "%1" USB001 >nul 2>&1
if %errorlevel%==0 (
    echo Successfully sent to USB001
    goto :success
)

REM Option 3: Try to copy to local printer share (if configured)
echo USB001 failed, attempting local printer share...
copy /B "%1" "\\localhost\ZebraPrinter" >nul 2>&1
if %errorlevel%==0 (
    echo Successfully sent to local printer share
    goto :success
)

REM If all methods fail, show the file content for debugging
echo.
echo All printing methods failed. ZPL file content:
echo ================================================
type "%1"
echo ================================================
echo.
echo Please configure your printer connection in this batch file.
echo Common options:
echo   - LPT1 (parallel port)
    - USB001 (USB printer)
echo   - \\\\computername\\printername (network share)
echo   - Configure Zebra Browser Print for web printing
goto :end

:success
echo Label sent to printer successfully!

:end
pause