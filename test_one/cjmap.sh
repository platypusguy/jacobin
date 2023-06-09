set -e
set batdir=`pwd`
cd ../src

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

