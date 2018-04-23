@echo off
set chain=EcuChain-1
echo Begining of the installation. The wizard is doing some old magic
if not exist C:\Users\"%USERNAME%"\AppData\Roaming\Multichain (GOTO :download) else (GOTO :init)
PAUSE

:init
chdir /D "C:\Users\"%USERNAME%"\AppData\Roaming\Multichain"
start cmd.exe
EXIT


:download
echo Downloading Multichain files ...
set curr_dir=%cd%
chdir /D C:\Users\"%USERNAME%"\AppData\Roaming
powershell.exe -Command "Invoke-WebRequest https://www.multichain.com/download/multichain-windows-1.0.4.zip -OutFile multichain.zip"
powershell.exe -nologo -noprofile -command "& { Add-Type -A 'System.IO.Compression.FileSystem'; [IO.Compression.ZipFile]::ExtractToDirectory('multichain.zip', 'Multichain'); }"
del "multichain.zip"
chdir /D %curr_dir%
goto init




::for /f "delims==  tokens=1,2" %%B in (%multi_dir%\params.dat) do set %%B=%%C
::for /f "delims=#  tokens=1,2" %%E in ("%default-rpc-port %") do set port=%%E
::echo MULTICHAIN_PORT=%port% >> %output%