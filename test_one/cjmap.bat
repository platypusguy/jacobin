@echo off
set batdir=%CD%
cd ..
cd src

echo
echo ======================
echo TestJacobinHomeTempdir 
echo ======================
go test  -v ./classloader -run TestJacobinHomeTempdir

echo
echo ======================
echo TestJacobinHomeDefault
echo ======================
go test -v ./classloader -run TestJacobinHomeDefault

REM Return to batch directory
cd %batdir%

