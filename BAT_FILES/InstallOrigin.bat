@echo off
set chain=Amacoin
set multidir = "C:\Users\%USERNAME%\AppData\Roaming\Multichain"
echo Begining of the installation. The wizard is doing some old magic
PAUSE
if not exist "C:\Users\%USERNAME%\AppData\Roaming\Multichain" (GOTO :download) else (GOTO :init)

:init
echo Initializing the chain ...
PAUSE
chdir /D C:"\Users\%USERNAME%\AppData\Roaming\Multichain"
:: Partie temporaire, on va le remplacer par client.exe
if not exist "C:\Users\%USERNAME%\AppData\Roaming\Multichain\multichain-util.exe" (echo Error no exe found in the directory.PAUSE EXIT)
if not exist "C:\Users\%USERNAME%\AppData\Roaming\Multichain\%chain%" (start multichain-util.exe create %chain%) else (goto :starting)
echo The wizard installed successfully the files and created the chain.

:starting
echo starting...
PAUSE
if not exist "C:\Users\%USERNAME%\AppData\Roaming\Multichain\%chain%" (GOTO :init)
ECHO bonjour
PAUSE
if not exist "C:\Users\%USERNAME%\AppData\Roaming\Multichain\multichaind.exe" (ECHO Error no exe found)
ECHO bonsoir
PAUSE
start multichaind.exe %chain% -daemon
PAUSE
EXIT

:download
echo Downloading Multichain files ...
set curr_dir=%cd%
chdir /D "C:\Users\%USERNAME%\AppData\Roaming"
powershell.exe -Command "Invoke-WebRequest https://www.multichain.com/download/multichain-windows-1.0.4.zip -OutFile multichain.zip"
powershell.exe -nologo -noprofile -command "& { Add-Type -A 'System.IO.Compression.FileSystem'; [IO.Compression.ZipFile]::ExtractToDirectory('multichain.zip', 'Multichain'); }"
del "multichain.zip"
chdir /D %curr_dir%
goto init


::for /f "delims==  tokens=1,2" %%B in (%multi_dir%\params.dat) do set %%B=%%C
::for /f "delims=#  tokens=1,2" %%E in ("%default-rpc-port %") do set port=%%E
::echo MULTICHAIN_PORT=%port% >> %output%