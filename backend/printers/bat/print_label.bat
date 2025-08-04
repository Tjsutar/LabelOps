@echo off
copy /B "%1" "\\%COMPUTERNAME%\\Your_Printer_Share_Name"
echo Label sent to printer.
pause