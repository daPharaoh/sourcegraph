package query

import "github.com/sourcegraph/sourcegraph/cmd/frontend/envvar"

func toRepoQuery(q Q) (*RepoQuery, bool) {
	return nil, false

}

func translate(q Q) InternalQuery {
	if _, ok := toRepoQuery(q); ok {
		if envvar.SourcegraphDotComMode() {
			return RepoQuery{DotComDefault{}}
		}
		return RepoQuery{Subset{}}
	}
	return GenericQuery{q}
}
