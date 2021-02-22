package toolkit

// 隐藏11位手机号
//var phoneHidden = regexp.MustCompile("(\\d{3})\\d{4}(\\d{4})")

func HiddenPhone(phone string) string {
	var result string
	for i := 0; i < len(phone); i++ {
		temp := phone[i]
		if i >= 3 && i <= 6 {
			temp = byte('*')
		}
		result += string(temp)
	}
	return result
}
