package testUtil

/*
testUtil中用到的功能性子函数
*/

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	. "../netFrame"
)

// 切片slice中是否含有str串，不是长串匹配
func sliceContains(slice []string, str string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	return false
}

// 记录节点图的边
func mainNet_add_edges(result string, causes []string) {
	for _, single := range causes {
		if !sliceContains(MainNet.Parents[result], single) {
			MainNet.Parents[result] = append(MainNet.Parents[result], single)
		}
		if !sliceContains(MainNet.Children[single], result) {
			MainNet.Children[single] = append(MainNet.Children[single], result)
		}
	}
}

// go 语法相关
func make_if_not_maked(variable string) {
	// 两者长度保持一致
	MainNet.PriorProb[variable] = make([]float64, len(MainNet.Values[variable]))
	if MainNet.Children[variable] == nil {
		MainNet.Children[variable] = make([]string, 0)
	}
	if MainNet.Parents[variable] == nil {
		MainNet.Parents[variable] = make([]string, 0)
	}
	if MainNet.InfluencedBy[variable] == nil {
		MainNet.InfluencedBy[variable] = make([][]string, 0)
	}
	if MainNet.ConditionalProb[variable] == nil {
		MainNet.ConditionalProb[variable] = make(map[string][]float64)
	}
}

func get_probs_for_variable(variable string, probSlice []string) []float64 {
	var prob []float64
	for _, value := range probSlice {
		temp, _ := strconv.ParseFloat(strings.TrimSpace(value), 64)
		prob = append(prob, temp)
	}
	return prob
}

func get_keys_for_condition(causes []string, rightSlice []string) string {
	// 记得排序避免顺序不同被认作不同条件
	var keys string
	mp := make(map[string]string)
	cause_vec := make([]string, 0)
	for i := 0; i < len(causes); i++ {
		mp[strings.TrimSpace(causes[i])] = strings.TrimSpace(rightSlice[i])
		cause_vec = append(cause_vec, strings.TrimSpace(causes[i]))
	}
	sort.Slice(cause_vec, func(i, j int) bool {
		return cause_vec[i] < cause_vec[j]
	})
	for _, single_cause := range cause_vec {
		keys += (single_cause + " : " + mp[single_cause] + " , ")
	}
	return keys[:len(keys)-3]
}

func PrintNet(net BNNet) {
	fmt.Println("--values--\n", net.Values)
	fmt.Println("--parents--\n", net.Parents)
	fmt.Println("--childred--\n", net.Children)
	fmt.Println("--influenced--\n", net.InfluencedBy)
	fmt.Println("--priorprob--\n", net.PriorProb)
	fmt.Println("--condprob--")
	for index, value := range net.ConditionalProb {
		fmt.Println(index)
		for key, prob := range value {
			fmt.Println(key, prob)
		}
	}
}
