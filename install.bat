REM Copyright (c) 2014, Rob Thornton
REM All rights reserved.
REM This source code is governed by a Simplied BSD-License. Please see the
REM LICENSE included in this distribution for a copy of the full license
REM or, if one is not included, you may also find a copy at
REM http://opensource.org/licenses/BSD-2-Clause

go test ".\..."

IF %ERRORLEVEL% EQU 0 cd "calcc" & go install & cd ".."
copy calc.bat "%GOPATH%\bin"
