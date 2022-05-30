package dag

func CheckCircle(nodes map[string][]string) (bool, []string) {
	indegrees := make(map[string]int)
	for node := range nodes {
		indegrees[node] = 0
	}
	for _, dsts := range nodes {
		for _, dst := range dsts {
			indegrees[dst] += 1
		}
	}

	// topological sort
	for {
		found := false
		for node, degree := range indegrees {
			if degree == 0 {
				found = true
				for _, dst := range nodes[node] {
					indegrees[dst] -= 1
				}
				delete(indegrees, node)
			}
		}
		if !found {
			break
		}
	}

	if len(indegrees) > 0 {
		// has a circle
		var circle []string
		for node := range indegrees {
			circle = append(circle, node)
		}
		return true, circle
	}
	return false, nil
}
