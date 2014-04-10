go test ".\..."

IF %ERRORLEVEL% EQU 0 cd "calcc" & go install & cd ".."
copy calc.bat "%GOPATH%\bin"
