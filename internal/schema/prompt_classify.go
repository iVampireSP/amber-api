package schema

type PromptClassify string
type QuestionLabel string

var (
	QuestionLabelSimpleTask       QuestionLabel = "Simple tasks"
	QuestionLabelComplexReasoning QuestionLabel = "complex reasoning"
	QuestionLabelCreativeWriting  QuestionLabel = "creative writing"
)

var QuestionLabels = []string{
	string(QuestionLabelCreativeWriting),
	string(QuestionLabelComplexReasoning),
	string(QuestionLabelSimpleTask),
}

var (
	PromptClassifyReAct PromptClassify = `
在回答问题时，使用以下输出
问题: 你必须要回答的问题
思考：你应该始终思考该做什么
操作：要采取的操作，你要是用什么工具，或者思考逻辑
动作输入：动作的输入
观察：行动的结果
思考：我现在知道最终答案了
最终答案：原始输入问题的最终答案
`
)

func (ql QuestionLabel) IsValid() bool {
	for _, label := range QuestionLabels {
		if label == string(ql) {
			return true
		}
	}

	return false
}

func (ql QuestionLabel) Prompt() string {
	if ql == QuestionLabelSimpleTask {
		return ""
	}

	if ql == QuestionLabelComplexReasoning {
		//return "请用更复杂的逻辑回答问题，不要直接输出答案，而是给出推理过程，并给出推理的步骤，最后给出答案。"
		return PromptClassifyReAct.String()
	}

	return ""
}

func (pc PromptClassify) String() string {
	return string(pc)
}
