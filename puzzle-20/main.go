package main

// https://adventofcode.com/2023/day/20

import (
	"aoc/helper"
	"fmt"
	"regexp"
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
				module = &BroadcastModule{name: moduleName, outputs: outputs}
			} else if m[1] == "%" {
				module = &FlipFlopModule{name: moduleName, outputs: outputs}
			} else if m[1] == "&" {
				module = &ConjunctionModule{name: moduleName, outputs: outputs, inputs: make(map[string]bool)}
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
	Modules map[string]Module
}

type Module interface {
	Name() string
	ProcessPulse(from string, high bool) (bool, []string)
	Clone() Module
	EqualState(other Module) bool
	Outputs() []string
	Reset()
	StateStr() string
}

type BroadcastModule struct {
	name    string
	outputs []string
}

func (m *BroadcastModule) Name() string { return m.name }
func (m *BroadcastModule) ProcessPulse(from string, high bool) (bool, []string) {
	return high, m.outputs
}
func (m *BroadcastModule) Clone() Module {
	return &BroadcastModule{outputs: m.outputs}
}
func (m *BroadcastModule) EqualState(other Module) bool { return true }
func (m *BroadcastModule) Outputs() []string            { return m.outputs }
func (m *BroadcastModule) Reset()                       {}
func (m *BroadcastModule) StateStr() string             { return "" }

type FlipFlopModule struct {
	name    string
	outputs []string
	isOn    bool
}

func (m *FlipFlopModule) Name() string { return m.name }
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
func (m *FlipFlopModule) StateStr() string {
	if m.isOn {
		return "1"
	}
	return "0"
}

type ConjunctionModule struct {
	name    string
	outputs []string
	inputs  map[string]bool
}

func (m *ConjunctionModule) Name() string { return m.name }
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
	return &ConjunctionModule{outputs: m.outputs, inputs: helper.CloneMap(m.inputs)}
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
func (m *ConjunctionModule) StateStr() string {
	var str string
	helper.IterateMapInKeyOrder(m.inputs, func(k string, v bool) {
		if len(str) > 0 {
			str += ","
		}
		str += k + ":"
		if v {
			str += "1"
		} else {
			str += "0"
		}
	})
	return str
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

type Pulse struct {
	From, To string
	High     bool
}

func (s *System) SimulateSingleButtonPush() (int64, int64) {
	if _, ok := s.Modules["broadcaster"]; !ok {
		helper.ExitWithMessage("system has no broadcaster")
	}
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

func (s *System) StateStr() string {
	var str string
	helper.IterateMapInKeyOrder(s.Modules, func(k string, m Module) {
		if len(str) > 0 {
			str += "|"
		}
		str += m.StateStr()
	})
	return str
}

func (s *System) CountButtonPushesForRXLow() int64 {
	s.DetectLoopsForRX()
	return 0
}

func (s *System) DetectLoopsForRX() {
	mBroadcast := s.Modules["broadcaster"].(*BroadcastModule)
	moduleToRX := s.FindModuleToRX()
	mLoopEnd, ok := moduleToRX.(*ConjunctionModule)
	if !ok {
		helper.ExitWithMessage("module to rx is not ConjunctionModule")
	}
	loopLengths := make([]int64, 0)
	for _, m := range mBroadcast.Outputs() {
		mod := s.Modules[m]
		subSystem := s.BuildSubSystem(mod, mLoopEnd)
		loopLength := subSystem.FindLoopLength()
		loopLengths = append(loopLengths, int64(loopLength-1))
	}
	fmt.Println(helper.LeastCommonMultiple(loopLengths...) + 1)
}

func (s *System) FindModuleToRX() Module {
	var moduleToRX Module
	for _, m := range s.Modules {
		for _, r := range m.Outputs() {
			if r == "rx" {
				if moduleToRX != nil {
					helper.ExitWithMessage("multiple modules to rx found")
				}
				moduleToRX = m
			}
		}
	}
	if moduleToRX == nil {
		helper.ExitWithMessage("no module to rx found")
	}
	return moduleToRX
}

func (s *System) BuildSubSystem(start, end Module) *System {
	modules := make(map[string]Module, 0)
	s.collectModules(modules, start, end)
	modules["broadcaster"] = &BroadcastModule{
		name:    "broadcaster",
		outputs: []string{start.Name()},
	}
	return &System{Modules: modules}
}

func (s *System) collectModules(modules map[string]Module, start, end Module) {
	if _, ok := modules[start.Name()]; ok {
		return
	}
	if start == end {
		return
	}
	modules[start.Name()] = start
	for _, m := range start.Outputs() {
		mod := s.Modules[m]
		s.collectModules(modules, mod, end)
	}
}

func (s *System) FindLoopLength() int {
	stateIndices := make(map[string]int)
	for i := 0; ; i++ {
		stateStr := s.StateStr()
		if loopStartIndex, ok := stateIndices[stateStr]; ok {
			fmt.Println("known state at", loopStartIndex)
			return i
		}
		stateIndices[stateStr] = i
		s.SimulateSingleButtonPush()
	}
}
