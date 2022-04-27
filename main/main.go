package main

import (
	"../netFrame"
	"../testUtil"
)

func main() {
	// https://www.bnlearn.com/bnrepository/
	// 没有考虑多变量的情况即 P(A, B | C)
	var fileName string
	fileName = "../data/original/child.bif"
	fileName = "../data/original/earthquake.bif"
	fileName = "../data/original/knowledgeClip.bif"
	fileName = "../data/original/test.bif"
	// 增加了对空格处理的支持，理论上每一行可以随便造空格
	// 有个问题读入的时候如果是自己敲的比如test.bif要以 "\r\n" 来 split，直接下载的直接 "\n"，后续有时间再说
	testUtil.Init(fileName)
	testUtil.PrintNet(netFrame.MainNet)
	var datasetSize = 1
	testUtil.GenerateDataset(datasetSize)
}
