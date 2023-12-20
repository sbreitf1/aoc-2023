package main

// https://adventofcode.com/2023/day/20

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"time"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	system := ParseSystem(lines)
	highCount, lowCount := system.CountPulsesForButtonPushes(1000)
	solution1 := highCount * lowCount
	fmt.Println("-> part 1:", solution1)

	system.Reset()
	solution2 := system.CountButtonPushesForRXLow()
	fmt.Println("-> part 2:", solution2)
}

func ParseSystem(lines []string) *System {
	pattern := regexp.MustCompile(`^([%&]?)([a-z]+)\s*->\s*(.*)$`)

	modules := make(map[string]Module)
	for _, line := range lines {
		m := pattern.FindStringSubmatch(line)
		if len(m) == 4 {
			moduleName := m[2]
			outputs := helper.SplitAndTrim(m[3], ",")
			var module Module
			if m[1] == "" {
				module = &BroadcastModule{outputs: outputs}
			} else if m[1] == "%" {
				module = &FlipFlopModule{outputs: outputs}
			} else if m[1] == "&" {
				module = &ConjunctionModule{outputs: outputs, inputs: make(map[string]bool)}
			} else {
				helper.ExitWithMessage("unsupported module %q", m[1])
			}
			modules[moduleName] = module
		}
	}
	for from, m := range modules {
		for _, o := range m.Outputs() {
			mo := modules[o]
			if mo, ok := mo.(*ConjunctionModule); ok {
				mo.inputs[from] = false
			}
		}
	}
	return &System{Modules: modules}
}

type System struct {
	Modules    map[string]Module
	RXLowCount int64
}

type Module interface {
	ProcessPulse(from string, high bool) (bool, []string)
	Clone() Module
	EqualState(other Module) bool
	Outputs() []string
	Reset()
}

type BroadcastModule struct {
	outputs []string
}

func (m *BroadcastModule) ProcessPulse(from string, high bool) (bool, []string) {
	return high, m.outputs
}
func (m *BroadcastModule) Clone() Module {
	return &BroadcastModule{outputs: m.outputs}
}
func (m *BroadcastModule) EqualState(other Module) bool { return true }
func (m *BroadcastModule) Outputs() []string            { return m.outputs }
func (m *BroadcastModule) Reset()                       {}

type FlipFlopModule struct {
	outputs []string
	isOn    bool
}

func (m *FlipFlopModule) ProcessPulse(from string, high bool) (bool, []string) {
	if high {
		return false, nil
	}
	m.isOn = !m.isOn
	return m.isOn, m.outputs
}
func (m *FlipFlopModule) Clone() Module {
	return &FlipFlopModule{outputs: m.outputs, isOn: m.isOn}
}
func (m *FlipFlopModule) EqualState(other Module) bool {
	return m.isOn == other.(*FlipFlopModule).isOn
}
func (m *FlipFlopModule) Outputs() []string { return m.outputs }
func (m *FlipFlopModule) Reset() {
	m.isOn = false
}

type ConjunctionModule struct {
	outputs []string
	inputs  map[string]bool
}

func (m *ConjunctionModule) ProcessPulse(from string, high bool) (bool, []string) {
	m.inputs[from] = high
	for _, v := range m.inputs {
		if !v {
			return true, m.outputs
		}
	}
	return false, m.outputs
}
func (m *ConjunctionModule) Clone() Module {
	clonedInputs := make(map[string]bool, len(m.inputs))
	for k, v := range m.inputs {
		clonedInputs[k] = v
	}
	return &ConjunctionModule{outputs: m.outputs, inputs: clonedInputs}
}
func (m *ConjunctionModule) EqualState(other Module) bool {
	if len(m.inputs) != len(other.(*ConjunctionModule).inputs) {
		return false
	}
	for k, v := range m.inputs {
		ov := other.(*ConjunctionModule).inputs[k]
		if ov != v {
			return false
		}
	}
	return true
}
func (m *ConjunctionModule) Outputs() []string { return m.outputs }
func (m *ConjunctionModule) Reset() {
	for k := range m.inputs {
		m.inputs[k] = false
	}
}

func (s *System) CountPulsesForButtonPushes(pushCount int64) (int64, int64) {
	var lowCount, highCount int64
	for i := int64(0); i < pushCount; i++ {
		l, h := s.SimulateSingleButtonPush()
		lowCount += l
		highCount += h
	}
	return lowCount, highCount
}

func (s *System) CountButtonPushesForRXLow() int64 {
	lastPrint := time.Now()
	for i := int64(1); ; i++ {
		s.SimulateSingleButtonPush()
		if s.RXLowCount > 0 {
			return i
		}

		if time.Since(lastPrint) > time.Second {
			fmt.Println(i)
			lastPrint = time.Now()
		}
	}
}

type Pulse struct {
	From, To string
	High     bool
}

func (s *System) SimulateSingleButtonPush() (int64, int64) {
	var lowCount, highCount int64
	pulses := []Pulse{
		{From: "button", To: "broadcaster", High: false},
	}
	for len(pulses) > 0 {
		p := pulses[0]
		pulses = pulses[1:]

		if p.High {
			//fmt.Println(p.From, "[high]", "->", p.To)
			highCount++
		} else {
			//fmt.Println(p.From, "[low]", "->", p.To)
			lowCount++
		}

		if p.To == "rx" && !p.High {
			s.RXLowCount++
		}

		if m, ok := s.Modules[p.To]; ok {
			outHigh, receivers := m.ProcessPulse(p.From, p.High)
			for _, r := range receivers {
				pulses = append(pulses, Pulse{From: p.To, To: r, High: outHigh})
			}
		}
	}
	return lowCount, highCount
}

func (s *System) Reset() {
	for _, m := range s.Modules {
		m.Reset()
	}
}
