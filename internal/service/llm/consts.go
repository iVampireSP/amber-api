package llm

// 强制停止（如果连续函数调用超过 4 次，则强制停止输出）
const forceStopCount = 6

// 警告次数（如果 LLM 连续调用超过 3 次，则警告输出）
const warningCount = 4

// 警告 LLM 调用太多次工具， 要求停止
const warningMessage = "[Warning]You are attempting to call the tool/function repeatedly, please use the tool/function properly and stop response. If you continue to call repeatedly, the chat will be forcibly terminated."
const forceStopSystemMessage = "[Force Stop]You have still repeatedly called the tool/function many times, and the chat has been forcibly terminated."
