package schema

//
//type PromptClassify string
//type QuestionLabel string
//
//var (
//	QuestionLabelSimpleTask       QuestionLabel = "Simple tasks"
//	QuestionLabelComplexReasoning QuestionLabel = "Complex reasoning"
//	QuestionLabelCreativeWriting  QuestionLabel = "Creative writing"
//	QuestionLabelMath             QuestionLabel = "Math"
//)
//
//var QuestionLabels = []string{
//	string(QuestionLabelCreativeWriting),
//	string(QuestionLabelComplexReasoning),
//	string(QuestionLabelSimpleTask),
//	string(QuestionLabelMath),
//}
//
//var (
//	PromptClassifyReAct PromptClassify = `
//在回答问题时，使用以下输出
//问题: 你必须要回答的问题
//思考：你应该始终思考该做什么
//操作：要采取的操作，你要是用什么工具，或者思考逻辑
//动作输入：动作的输入
//观察：行动的结果
//思考：我现在知道最终答案了
//最终答案：原始输入问题的最终答案
//
//如果你正在计算，你必须使用计算器工具，无论如何都不允许使用自己的知识或不计算进行输出，计算器永远比你正确的并且不会出错。
//`
//)
//
//func (ql QuestionLabel) IsValid() bool {
//	for _, label := range QuestionLabels {
//		if label == string(ql) {
//			return true
//		}
//	}
//
//	return false
//}
//
//func (ql QuestionLabel) Prompt() string {
//	var prompt = ""
//
//	switch ql {
//	case QuestionLabelComplexReasoning:
//		prompt = PromptClassifyReAct.String()
//	case QuestionLabelMath:
//		prompt = PromptClassifyReAct.String()
//	}
//
//	return prompt
//}
//
//func (pc PromptClassify) String() string {
//	return string(pc)
//}
