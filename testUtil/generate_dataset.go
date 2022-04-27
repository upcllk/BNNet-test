package testUtil

import (
	"bufio"
	"fmt"
	"os"

	. "../netFrame"
)

func GenerateDataset(dataCount int) {
	fileName := "../data/dataset/" + MainNet.Name + ".data"
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("wrong while writing file,", err)
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	// 根据拓扑排序获得节点值的生成顺序
	order := get_generate_order()
	// fmt.Println(order)
	for i := 0; i < dataCount; i++ {
		data := generate_single_data(order)
		// fmt.Println(data)
		write.WriteString(data + "\n")
	}
	write.Flush()
}

// 返回顺序按照 MainNet.Variables 的变量顺序
func generate_single_data(order []string) string {
	// dic[variable] = index 表示当前 variable 取值 MainNet.Values[variables][index]为
	dic := make(map[string]int)
	for _, cur := range order {
		// 根节点 : priorProb，其他 : condProb
		// 类似前缀和一样的方法生成
		// 封装一个 getRandomValue(values[]string, prob[]float64)
		prob := make([][]float64, 0)
		if len(MainNet.Parents[cur]) == 0 {
			prob = append(prob, MainNet.PriorProb[cur])
		} else {
			/*
				如果 P(A = True | B = True, C = False) = 0.4,
					P(A = True | C = False, D = False) = 0.7,
				在已知 B = True, C = D = False 下，
				condKeys = [
					"A : True, C : False",
					"C : False, D : False"
				]
				对应 prob = [
					[0.4, 0.6],
					[0.7, 0.3]
				]
				当然 prob 顺序要跟 MainNet.Values 一一对应
				如果只有一个条件的话当然就是 [[0.4, 0.6]]
				然后扔给 get_random_index_from_status 产生符合条件的随机值
				产生 index 为下标，对应值为
			*/
			condKeys := get_keys_from_status(cur, dic)
			fmt.Println("***************", "condKeys : ", condKeys, "***************")
			for _, key := range condKeys {
				prob = append(prob, MainNet.ConditionalProb[cur][key])
			}
		}
		fmt.Println("CUR : ", cur, "\nPROB : ", prob)
		dic[cur] = get_random_index_from_status(cur, prob)
	}
	var result string
	return result
}

// 返回 slice 因为 InfluencedBy[variable] 可能不止一个，这时候就需要 1 - (1 - prob1) * (1 - prob2) * ...
func get_keys_from_status(cur string, dic map[string]int) []string {
	keys := make([]string, 0)
	for _, reasons := range MainNet.InfluencedBy[cur] {
		var temp string
		for _, single := range reasons {
			temp += (single + " : " + MainNet.Values[single][dic[single]] + " , ")
		}
		temp = temp[:len(temp)-3]
		keys = append(keys, temp)
	}
	return keys
}

func get_random_index_from_status(cur string, prob [][]float64) int {
	new_prob := merge_prob(prob)
	fmt.Println("merge prob : ", new_prob)
	return 0
}

func merge_prob(ori [][]float64) []float64 {
	if len(ori) == 0 {
		fmt.Println("error while merge prob")
		os.Exit(0)
	}
	result := make([]float64, len(ori[0]))
	for i := 0; i < len(result); i++ {
		var prob float64 = 1
		for j := 0; j < len(ori); j++ {
			prob *= (1 - ori[j][i])
		}
		prob = 1 - prob
		result[i] = prob
	}
	return result
}

func get_generate_order() []string {
	// 拓扑排序
	variables := make([]string, len(MainNet.Variables))
	// var children map[string][]string
	copy(variables, MainNet.Variables)
	visited := make([]bool, len(variables))
	inDgree := make(map[string]int)
	for _, child := range MainNet.Variables {
		inDgree[child] = len(MainNet.Parents[child])
	}
	var order []string
	for len(order) < len(variables) {
		for i := 0; i < len(variables); i++ {
			cur := variables[i]
			if visited[i] == true || inDgree[cur] != 0 {
				continue
			}
			order = append(order, cur)
			visited[i] = false
			for _, child := range MainNet.Children[cur] {
				inDgree[child]--
			}
		}
	}
	return order
}
