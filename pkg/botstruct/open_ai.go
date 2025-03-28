package botstruct

import "github.com/openimsdk/tools/errs"

type Request struct {
	Model              string      `json:"model"`
	Input              []InputItem `json:"input"`
	PreviousResponseID string      `json:"previous_response_id,omitempty"`
}

type InputItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	ID                 string        `json:"id"`
	Object             string        `json:"object"`
	CreatedAt          int64         `json:"created_at"`
	Status             string        `json:"status"`
	Error              interface{}   `json:"error"`
	IncompleteDetails  interface{}   `json:"incomplete_details"`
	Instructions       interface{}   `json:"instructions"`
	MaxOutputTokens    interface{}   `json:"max_output_tokens"`
	Model              string        `json:"model"`
	Output             []OutputItem  `json:"output"`
	ParallelToolCalls  bool          `json:"parallel_tool_calls"`
	PreviousResponseID interface{}   `json:"previous_response_id"`
	Reasoning          Reasoning     `json:"reasoning"`
	Store              bool          `json:"store"`
	Temperature        float64       `json:"temperature"`
	Text               TextFormat    `json:"text"`
	ToolChoice         string        `json:"tool_choice"`
	Tools              []interface{} `json:"tools"`
	TopP               float64       `json:"top_p"`
	Truncation         string        `json:"truncation"`
	Usage              Usage         `json:"usage"`
	User               interface{}   `json:"user"`
	Metadata           Metadata      `json:"metadata"`
}

func (r *Response) GetContentAndID() (string, string, error) {
	if len(r.Output) == 0 {
		return "", "", errs.New("no output").Wrap()
	}
	if len(r.Output[0].Content) == 0 {
		return "", "", errs.New("no content").Wrap()
	}
	return r.ID, r.Output[0].Content[0].Text, nil
}

type OutputItem struct {
	Type    string    `json:"type"`
	ID      string    `json:"id"`
	Status  string    `json:"status"`
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type        string        `json:"type"`
	Text        string        `json:"text"`
	Annotations []interface{} `json:"annotations"`
}

type Reasoning struct {
	Effort          interface{} `json:"effort"`
	GenerateSummary interface{} `json:"generate_summary"`
}

type TextFormat struct {
	Type string `json:"type"`
}

type Usage struct {
	InputTokens         int            `json:"input_tokens"`
	InputTokensDetails  map[string]int `json:"input_tokens_details"`
	OutputTokens        int            `json:"output_tokens"`
	OutputTokensDetails map[string]int `json:"output_tokens_details"`
	TotalTokens         int            `json:"total_tokens"`
}

type Metadata struct{}
