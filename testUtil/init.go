package testUtil

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	. "../netFrame"
)

func Init(fileName string) {
	init_net(&MainNet)
	text := readFile(fileName)
	parse_text(text)
	fmt.Println("Init .. OK")
}

func init_net(net *BNNet) {
	net.Variables = make([]string, 0)
	net.Values = make(map[string][]string)
	net.Children = make(map[string][]string)
	net.Parents = make(map[string][]string)
	net.InfluencedBy = make(map[string][][]string)
	net.PriorProb = make(map[string][]float64)
	net.ConditionalProb = make(map[string]map[string][]float64)
}

func readFile(fileName string) []string {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	// fmt.Println(string(bytes))
	// text := strings.Split(string(bytes), "\n")
	text := strings.Split(string(bytes), "\n")
	// fmt.Println(text)
	return text
}

func parse_text(text []string) {
	var left, right int = 0, 0
	for left < len(text) {
		for right < len(text) && !strings.HasSuffix(text[right], "}") {
			right++
		}
		if strings.HasPrefix(text[left], "network") {
			parse_netName(text[left : right+1])
		} else if strings.HasPrefix(text[left], "variable") {
			parse_variable(text[left : right+1])
		} else if strings.HasPrefix(text[left], "probability") {
			parse_probability(text[left : right+1])
		} else {
			// 一般是遇到了空行。。可以忽略？
			fmt.Println("error while read in line", left+1, " .. file format wrong, maybe empty line")
			// fmt.Println(text[left:right])
		}
		left = right + 1
		right = left
	}
}

// 应该是啥用没有
func parse_netName(text []string) {
	/*
		network unknown {
		}
	*/
	MainNet.Name = strings.Split(text[0], " ")[1]
}

// 记录各变量的取值范围
func parse_variable(text []string) {
	/*
		variable Earthquake {
			type discrete [ 2 ] { True, False };
		}
	*/
	variable := strings.Split(text[0], " ")[1]
	values_str := strings.Split(text[1], "{")[1]
	values_str = strings.Split(values_str, "}")[0]
	if !sliceContains(MainNet.Variables, variable) {
		MainNet.Variables = append(MainNet.Variables, variable)
	}
	// fmt.Println(possible_values)
	possible_values := strings.Split(values_str, ",")
	for index, value := range possible_values {
		possible_values[index] = strings.TrimSpace(value)
		// fmt.Printf("-%s-\n", possible_values[index])
	}
	// 注意这边 copy 也要留出足够空间
	// MainNet.Values[variable] = make([]string, 0)
	MainNet.Values[variable] = make([]string, len(possible_values))
	copy(MainNet.Values[variable], possible_values)
	make_if_not_maked(variable)
}

// 记录概率，包括先验概率和条件概率
func parse_probability(text []string) {
	/*
		probability ( Earthquake ) {
			table 0.02, 0.98;
		}
		---------------------------------------------------
		probability ( Alarm | Burglary, Earthquake ) {
			(True, True) 0.95, 0.05;
			(False, True) 0.29, 0.71;
			(True, False) 0.94, 0.06;
			(False, False) 0.001, 0.999;
		}
		---------------------------------------------------
		probability ( JohnCalls | Alarm ) {
			(True) 0.9, 0.1;
			(False) 0.05, 0.95;
		}
		... (True, False, True, True, ... )
	*/
	if strings.Contains(text[0], "|") {
		parse_probability_cond(text)
	} else {
		parse_probability_prior(text)
	}
}

// 处理先验概率
// probability ( Earthquake ) {
// 	table 0.02, 0.98;
// }
func parse_probability_prior(text []string) {
	// [Earthquake]
	variable := strings.TrimSpace(strings.Split(strings.Split(text[0], "(")[1], ")")[0])
	text[1] = strings.TrimSpace(text[1])
	//   table    0.01, 0.99;
	tempStr := strings.TrimSpace(strings.Split(strings.TrimSpace(strings.Split(text[1], ";")[0]), "table")[1])
	// "0.01, 0.99"
	slice := strings.Split(tempStr, ",")
	for index, value := range slice {
		prob, _ := strconv.ParseFloat(strings.TrimSpace(value), 64)
		MainNet.PriorProb[variable][index] = prob
	}
}

// 处理条件概率
// probability ( Alarm | Burglary, Earthquake ) {
// 	(True, True) 0.95, 0.05;
// 	(False, True) 0.29, 0.71;
// 	(True, False) 0.94, 0.06;
// 	(False, False) 0.001, 0.999;
// }
// ---------------------------------------------------
// probability ( JohnCalls | Alarm ) {
// 	(True) 0.9, 0.1;
// 	(False) 0.05, 0.95;
// }
func parse_probability_cond(text []string) {
	/*
		net.Probability[key] = probability
		关于概率的保存
		key使用get_key_from_map方法得到一个串，比如
		probability ( Alarm | Burglary, Earthquake ) {
			(True, True) 0.95, 0.05;
		key = "Burglary:True,Earthquake:True"
		顺序不用担心顺序问题比如Burglary:True,Earthquake:True
		[Bu, Ea]这个slice排个序就行
		对于多父节点独立影响的问题看稿纸先不写
	*/
	result, causes := parse_reason_cause(text[0])
	mainNet_add_edges(result, causes)
	copy_causes := make([]string, len(causes))
	copy(copy_causes, causes)
	sort.Slice(copy_causes, func(i, j int) bool {
		return copy_causes[i] < copy_causes[j]
	})
	MainNet.InfluencedBy[result] = append(MainNet.InfluencedBy[result], copy_causes)
	text = text[1 : len(text)-1]
	for _, single := range text {
		single = strings.TrimSpace(single)
		// "(True, True) 0.95, 0.05;"
		// "(True) 0.9, 0.1;"
		// P(left | right)
		// leftStr : 0.95, 0.05
		// rightStr : (True, True)
		rightSlice := strings.Split(strings.TrimSpace(strings.Split(strings.Split(single, ")")[0], "(")[1]), ",")
		leftSlice := strings.Split(strings.TrimSpace(strings.Split(strings.Split(single, ")")[1], ";")[0]), ",")
		// fmt.Println("left : \n", leftSlice)
		// fmt.Println("right : \n", rightSlice)
		// 为了方便观察下面两个 get_key_ 函数以及组成的条件概率 key 符号两边都有一个空格
		var_prob := get_probs_for_variable(result, leftSlice)
		cond_keys := get_keys_for_condition(causes, rightSlice)
		// fmt.Println("leftkeys :\n", left_keys)
		// fmt.Println("leftkeysProb :\n", left_keys_prob)
		// fmt.Println("rightkeys :\n", right_keys)
		MainNet.ConditionalProb[result][cond_keys] = make([]float64, len(var_prob))
		copy(MainNet.ConditionalProb[result][cond_keys], var_prob)
	}
}

func parse_reason_cause(str string) (string, []string) {
	/*
		probability ( Alarm | Burglary, Earthquake ) {
	*/
	slice := strings.Split(str, "|")
	var result = strings.TrimSpace(strings.Split(slice[0], "(")[1])
	var cause []string
	tempStr := strings.TrimSpace(strings.Split(slice[1], ")")[0])
	cause = strings.Split(tempStr, ",")
	for index, value := range cause {
		cause[index] = strings.TrimSpace(value)
	}
	// for _, single := range cause {
	// 	fmt.Printf("-%s-\n", single)
	// }
	// fmt.Printf("-%s-\n", result)
	return result, cause
}
