package query

import "github.com/sourcegraph/sourcegraph/cmd/frontend/envvar"

func toRepoQuery(q Q) (*RepoQuery, bool) {
	includeRepos, excludeRepos := FindFields(q, "repo")            // enforce isNotNegated
	includeRepoGroups := FindField(q, "repogroup")                 // enforce singular, isNotNegated
	fork := ParseYesNoOnly(FindPositiveField(q, "fork"))           // enforce singular
	archived := ParseYesNoOnly(FindPositiveField(q, "archived"))   // enforce singular
	privacy := ParseVisibility(FindPositiveField(q, "visibility")) // ... etc..
	// TODO: deal with repoHasCommitAfter.
	contextID := FindPositiveField(q, "context")
	// TODO: version context from URL parameter.

	return nil, false
}

func translate(q Q) InternalQuery {
	if repoQuery, ok := toRepoQuery(q); ok {
		if envvar.SourcegraphDotComMode() {
			return RepoQuery{DotComDefault{}}
		}
		return repoQuery
	}
	return GenericQuery{q}
}
