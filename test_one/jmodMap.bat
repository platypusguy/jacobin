@echo off
set batdir=%CD%
cd ..
cd src

echo
echo ======================
echo TestJacobinHomeTempdir 
echo ======================
go test  -v ./classloader -run TestJmodMapHomeTempdir

echo
echo ======================
echo TestJacobinHomeDefault
echo ======================
go test -v ./classloader -run TestJmodMapHomeDefault

REM Return to batch directory
cd %batdir%

