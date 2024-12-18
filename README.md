# cgo_test
cgo test，加载静态，和动态链接库 （linux环境）

# 编译动态库
gcc -shared -fPIC -o libexample.so example.c
gcc -shared -fPIC -o libperson.so person.c

# 动态库放置路径
放在和go文件同一目录下，路径读取：soPath := "./libexample.so"
如果没有放在同一个目录下，需要指定路径，如：soPath := "/usr/lib/libexample.so"