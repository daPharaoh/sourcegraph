package query

type repoPrivacy string

const (
	AnyPrivacy repoPrivacy = "anyPrivacy"
	Private    repoPrivacy = "private"
	Public     repoPrivacy = "public"
)

type repoState string

const (
	Cloned    repoState = "cloned"
	NotCloned repoState = "notCloned"
)

type repoLabelName string

const (
	AnyLabel repoLabelName = "anyLabel"
	Fork     repoLabelName = "fork"
	Archive  repoLabelName = "archive"
)

type repoLabel struct {
	name    repoLabelName
	negated bool
}

type scope struct {
	privacy         []repoPrivacy
	state           []repoState
	label           []repoLabel
	includePatterns []string
	excludePattern  string
}

type limit_offset struct {
	limit  int
	offset int
}

type options struct {
	scope        scope
	limit_offset limit_offset
}

/* repoSetID variants */
type repoSetID interface {
	repoSetIDValue()
}

func (VersionContext) repoSetIDValue() {}
func (Context) repoSetIDValue()        {}

type VersionContext struct {
	ID string
}

type Context struct {
	ID string
}

/* repoSet variants */
type repoSet interface {
	repoSetValue()
}

func (DotComDefault) repoSetValue() {}
func (GlobalSet) repoSetValue()     {}
func (LabeledSubset) repoSetValue() {}
func (Subset) repoSetValue()        {}

type DotComDefault struct {
	Options options
}

type GlobalSet struct {
	Options options
}

type LabeledSubset struct {
	ID      repoSetID
	Options options
}

type Subset struct {
	Options options
}

/* Internal Query */
type InternalQuery interface {
	internalQueryValue()
}

func (RepoQuery) internalQueryValue()    {}
func (GenericQuery) internalQueryValue() {}

type RepoQuery struct {
	repoSet
}

type GenericQuery struct {
	Q
}
