MATCH (n) RETURN (n) # получить все вершины и ребра
MATCH (n: `go-remerge` {graph: 'file_dependency'}) RETURN (n) # получить граф file_dependency'