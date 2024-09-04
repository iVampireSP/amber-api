package memory

import (
	"context"
	"rag-new/internal/schema"
)

func (s *Service) GenerateMemoryPrompt(ctx context.Context, userId schema.UserId) (string, error) {
	memories, err := s.GetMemories(ctx, userId)
	if err != nil {
		return "", err
	}

	var m = ""

	for _, memory := range memories {
		m += memory.Content
		m += "\n"
	}

	return m, nil
}

func (s *Service) updateMemoryPrompt(existingMemories string, memory string) string {
	return `You are an expert at merging, updating, and organizing memories. When provided with existing memories and new information, your task is to merge and update the memory list to reflect the most accurate and current information. You are also provided with the matching score for each existing memory to the new information. Make sure to leverage this information to make informed decisions about which memories to update or merge.

Guidelines:
- Eliminate duplicate memories and merge related memories to ensure a concise and updated list.
- If a memory is directly contradicted by new information, critically evaluate both pieces of information:
- If the new memory provides a more recent or accurate update, replace the old memory with new one.
- If the new memory seems inaccurate or less detailed, retain the old memory and discard the new one.
- Maintain a consistent and clear style throughout all memories, ensuring each entry is concise yet informative.
- If the new memory is a variation or extension of an existing memory, update the existing memory to reflect the new information.
- If there are duplicates in the existing memories, only the more specific one needs to be retained. For example: Like reading and Like reading book are duplicates.
- If there are more than 5 existing memories, some memories need to be merged.
- Save memories in Simplified Chinese

Here are the details of the task:
- Existing Memories:
` + existingMemories + `

- New Memory: ` + memory + `
`

}

func (s *Service) memoryDeductionPrompt(userInput string) string {
	return `Deduce the facts, preferences, and memories from the provided text.
Just return the facts, preferences, and memories in bullet points:
Natural language text: ` + userInput + `

Constraint for deducing facts, preferences, and memories:
- The facts, preferences, and memories should be concise and informative.
- Don't start by "The person likes Pizza". Instead, start with "Likes Pizza".
- Don't remember the user/agent details provided. Only remember the facts, preferences, and memories.
- Save memories in Simplified Chinese

Deduced facts, preferences, and memories:`
}

func (s *Service) memoryAnswerPrompt() string {
	return `You are an expert at answering questions based on the provided memories. Your task is to provide accurate and concise answers to the questions by leveraging the information given in the memories.

Guidelines:
- Extract relevant information from the memories based on the question.
- If no relevant information is found, make sure you don't say no information is found. Instead, accept the question and provide a general response.
- Ensure that the answers are clear, concise, and directly address the question.
- Save memories in Simplified Chinese

Here are the details of the task:`
}
