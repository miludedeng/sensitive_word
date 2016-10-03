package service

import (
	"testing"
)

func TestSensitiveFind(t *testing.T) {

	var tests = []struct {
		str  string
		want int
	}{
		{str: "精子分动物精子与植物精子。动物有性生殖过程中的雄性细胞，雄性动物的生殖细胞，异配生殖中的雄配子，由精子器产生的单倍体生殖细胞。而生活中精子更多是指男性成熟的生殖细胞，在精巢中形成。精液是一种有机物。精液中含有果糖和蛋白质、以及一些酶类物质、无机盐和有机盐。", want: 7},
		{str: "刘晓庆太平洋游泳美翻 长发暴乳好身材迷倒网友赞性感女神", want: 1},
	}
	for _, test := range tests {
		len := len(SensitiveFind([]rune(test.str)))
		if len != test.want {
			t.Errorf("SensitiveFind want %d,but result is %d", test.want, len)
		}
	}
}

func BenchmarkSensitiveFind(b *testing.B) {
	str := "精子分动物精子与植物精子。动物有性生殖过程中的雄性细胞，雄性动物的生殖细胞，异配生殖中的雄配子，由精子器产生的单倍体生殖细胞。而生活中精子更多是指男性成熟的生殖细胞，在精巢中形成。精液是一种有机物。精液中含有果糖和蛋白质、以及一些酶类物质、无机盐和有机盐。"
	for i := 0; i < b.N; i++ {
		SensitiveFind([]rune(str))
	}
}
