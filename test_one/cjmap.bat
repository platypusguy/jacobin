cd ..
cd src

echo
echo ======================
echo TestJacobinHomeTempdir 
echo ======================
go test ./classloader -run TestJacobinHomeTempdir  -v

echo
echo ======================
echo TestJacobinHomeDefault
echo ======================
go test ./classloader -run TestJacobinHomeDefault -v

