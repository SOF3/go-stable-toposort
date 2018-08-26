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
	"sort"
)

type Node interface {
	Before(other Node) bool
}

type nodeNumber int
type edge [2]nodeNumber
type edgeNumber int

type edgeIndex struct {
	slice []edge // probably
	index [2]map[nodeNumber]map[nodeNumber]edgeNumber
}

func newEdgeIndex() *edgeIndex {
	index := &edgeIndex{
		slice: make([]edge, 0),
	}
	for i := range index.index {
		index.index[i] = map[nodeNumber]map[nodeNumber]edgeNumber{}
	}
	return index
}
func (index *edgeIndex) add(edge edge) edgeNumber {
	number := edgeNumber(len(index.slice))
	index.slice = append(index.slice, edge)
	for pos := 0; pos < 2; pos++ {
		if _, exists := index.index[pos][edge[pos]]; !exists {
			index.index[pos][edge[pos]] = make(map[nodeNumber]edgeNumber)
		}
		index.index[pos][edge[pos]][edge[1-pos]] = number
	}
	return number
}
func (index edgeIndex) removeIndex(edge edge) {
	for pos := range [...]int{0, 1} {
		delete(index.index[pos][edge[pos]], edge[1-pos])
		if len(index.index[pos][edge[pos]]) == 0 {
			delete(index.index[pos], edge[pos])
		}
	}
}

// Sorts nodes by Kahn's algorithm
func Sort(nodes []Node) (output []Node, cycle []Node) {
	edges := newEdgeIndex()

	var i nodeNumber
	for i = 0; int(i) < len(nodes); i++ {
		var j nodeNumber
		for j = i + 1; int(j) < len(nodes); j++ {
			ij := nodes[i].Before(nodes[j])
			ji := nodes[j].Before(nodes[i])
			if ij && ji {
				return nil, []Node{nodes[i], nodes[j]}
			}
			if ij {
				edges.add(edge{i, j})
			} else if ji {
				edges.add(edge{j, i})
			}
		}
	}

	output = make([]Node, 0, len(nodes))

	roots := make([]nodeNumber, 0, len(nodes))
	{
		for mInt := range nodes {
			m := nodeNumber(mInt)
			if _, hasBefore := edges.index[1][m]; !hasBefore {
				roots = append(roots, m)
			}
		}
	}

	for len(roots) > 0 {
		n := roots[0]
		roots = roots[1:]
		output = append(output, nodes[n])

		var mSlice = make([]nodeNumber, 0, len(edges.index[0][n]))
		for m := range edges.index[0][n] {
			mSlice = append(mSlice, m)
		}
		sort.SliceStable(mSlice, func(i, j int) bool {
			return mSlice[i] < mSlice[j]
		}) // stabilize the output because we are using a map
		for _, m := range mSlice {
			e := edges.index[0][n][m]
			edges.removeIndex(edges.slice[e])
			if _, hasBefore := edges.index[1][m]; !hasBefore {
				roots = append(roots, m)
			}
		}
	}

	for pos := 0; pos < 2; pos++ {
		if len(edges.index[pos]) > 0 {
			cycle = make([]Node, 0, len(edges.index[0]))
			for n := range edges.index[pos] {
				cycle = append(cycle, nodes[n])
			}
			return nil, cycle
		}
	}

	return output, nil
}
