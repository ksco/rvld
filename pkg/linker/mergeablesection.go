package linker

import "sort"

type MergeableSection struct {
	Parent      *MergedSection
	P2Align     uint8
	Strs        []string
	FragOffsets []uint32
	Fragments   []*SectionFragment
}

func (m *MergeableSection) GetFragment(offset uint32) (*SectionFragment, uint32) {
	pos := sort.Search(len(m.FragOffsets), func(i int) bool {
		return offset < m.FragOffsets[i]
	})

	if pos == 0 {
		return nil, 0
	}

	idx := pos - 1
	return m.Fragments[idx], offset - m.FragOffsets[idx]
}
