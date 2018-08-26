/*
 * go-stable-toposort
 *
 * Copyright (C) 2018 SOFe
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package stableToposort

import (
	"testing"
	"log"
)

type testNode int

func (n testNode) Before(m Node) bool {
	type pair = [2]testNode
	switch (pair{n, m.(testNode)}) {
	case pair{5, 11}:
		return true
	case pair{7, 11}:
		return true
	case pair{7, 8}:
		return true
	case pair{3, 8}:
		return true
	case pair{3, 10}:
		return true
	case pair{11, 2}:
		return true
	case pair{11, 10}:
		return true
	case pair{11, 9}:
		return true
	case pair{8, 9}:
		return true
	default:
		return false
	}
}

func TestSort(t *testing.T) {
	for i := 0; i < 10000; i++ {
		sorted, bad := Sort([]Node{
			testNode(10),
			testNode(2),
			testNode(5),
			testNode(3),
			testNode(11),
			testNode(8),
			testNode(9),
			testNode(7),
		})

		if bad != nil {
			t.Errorf("returned bad: %v", bad)
			t.Fail()
		}

		expected := []testNode{5, 3, 7, 11, 8, 10, 2, 9}

		log.Printf("%v", sorted)
		for j := 0; j < len(sorted) || j < len(expected); j++ {
			if sorted[j] != expected[j] {
				t.Errorf("%v != %v", sorted[j], expected[j])
				t.Fail()
			}
		}
	}
}
