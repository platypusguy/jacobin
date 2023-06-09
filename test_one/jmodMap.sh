set -e
set batdir=`pwd`
cd ../src

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

