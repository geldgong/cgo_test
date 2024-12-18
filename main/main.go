package main

func main() {
	// 动态调用简单类型的函数，example.c 中定义的函数
	Handle1()

	// 动态调用结构体类型的函数，并使用反射进行释放动态链接库，person.c 中定义的函数
	Handle3()

	// 动态调用结构体类型的函数，并启用map进行管理和释放动态链接库, 复制person.中定义的函数
	Handle2()

}
