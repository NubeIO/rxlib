package priority

import "fmt"

func NewPriority() *Priority {
	return &Priority{}
}

type Priority struct {
	Values [16]PriorityValue
}

func (p *Priority) SetValue(value PriorityValue, priorityNumber int) {
	if priorityNumber >= 1 && priorityNumber <= 16 {
		p.Values[priorityNumber-1] = value
	}
}

func (p *Priority) SetNull(priorityNumber int) {
	if priorityNumber >= 1 && priorityNumber <= 16 {
		p.Values[priorityNumber-1] = nil
	}
}

func (p *Priority) GetHighestPriorityValue() (PriorityValue, int) {
	for i, v := range p.Values {
		if v != nil {
			return v, i + 1
		}
	}
	return nil, 0
}

func (p *Priority) GetLowestPriorityValue() (PriorityValue, int) {
	for i := len(p.Values) - 1; i >= 0; i-- {
		if p.Values[i] != nil {
			return p.Values[i], i + 1
		}
	}
	return nil, 0
}

func (p *Priority) GetByPriorityNumber(priorityNumber int) PriorityValue {
	if priorityNumber >= 1 && priorityNumber <= 16 {
		return p.Values[priorityNumber-1]
	}
	return nil
}

func (p *Priority) ToMap() map[string]interface{} {
	jsonMap := make(map[string]interface{})
	for i, val := range p.Values {
		key := fmt.Sprintf("_%d", i+1) // Keys like _1, _2, ..., _16
		if val != nil {
			jsonMap[key] = val.GetValue()
		} else {
			jsonMap[key] = nil
		}
	}

	return jsonMap
}
