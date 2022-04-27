package netFrame

type BNNet struct {
	Name string
	/*
			Values为各变量的取值范围，对于地震模型就是True False
			variable CardiacMixing {
		  		type discrete [ 4 ] { None, Mild, Complete, Transp. };
			}
			Values[CardiacMixing] = [None, Mild, Complete, Transp. ]
			Probability["Earthquake"] = 0.02 表示 Earthquake 为 True 的概率是 0.02 初始为 -1 表示他有祖先节点
			如果不在 | 左边进入 parse_probability_single 函数，内部设置对应值
	*/
	Variables []string
	Values    map[string][]string
	Parents   map[string][]string
	Children  map[string][]string
	// 受哪些条件变量取值的影响
	// 比如有 P(A | B, C), P(A | B, D)
	// 		InfluencedBy[A] = [
	// 			[B, C],
	// 			[B, D]
	// 		]
	// 		[B, C], [B, D] 为 sort 后的方便后面做 key
	InfluencedBy map[string][][]string
	// 先验概率
	PriorProb map[string][]float64
	// 条件概率
	// ConditionalProb[conditionKey][variable] 是 Values[variable] 各个取值对应的概率分布
	ConditionalProb map[string]map[string][]float64
	/*
		ConditionalProb[A]["B:True,C:True"] = P(A = True | B = True, C = True)
		因为还有 P(A | B, D), 那么我们在知道 BCD 取值的时候，去查 InfluencedBy
		根据 [B, C] 得到一个 key 比如 "B:False,C:True,A:True", 去查ConditionalProb[A][key] -> num1
		num1 为 A 为 True 的概率
		这里之所以 key 要有一个 A:True 是因为我们还要考虑 A 有三种取值的时候比如 [True, False, None]
		同样方法根据 [C, D] 得到一个 num2 -> slice := [num1, num2] / slice = append(slice, num2)
		此时 A 以 prob 的概率为 True
		其中 prob = 1 - (1 - num1) * (1 - num2)
	*/
}

var MainNet BNNet
