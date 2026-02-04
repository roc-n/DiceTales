package tmp

import "fmt"

func SliceTrouble() {
	// a := []int{1, 2, 3, 4, 5}
	// b := a[1:2]

	// b[0] = 10
	// fmt.Println(a) // 输出: [1 10 3 4 5]

	// a := []int{1, 2, 3}
	// b := a[1:]
	// a = append(a, 4, 5, 6) // 扩容，a 指向了新的底层数组
	// b[0] = 99
	// fmt.Println(a) // [1 2 3 4 5 6]，b 修改没影响 a

	s := []int{0, 2, 3} // 值传递
	f(s)
	fmt.Println(s) // [99 2 3]
}

func f(s []int) {
	s[0] = 99    // 会修改原 slice
	s = []int{1} // 不会影响调用者
	fmt.Println(s)
}
