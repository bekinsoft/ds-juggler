package juggler

var (
	filters = []string{
		"fields",  // not implemented
		"include", // passed
		"limit",   // passed
		"order",   // passed
		"offset",  // or skip // passed
		"where",   // passed
	}
	operators = map[string]string{
		"=":   "=",   // passed
		"and": "and", // passed
		"or":  "or",  // passed
		"gt":  ">",   // passed
		"gte": ">=",  // passed
		"lt":  "<",   // passed
		"lte": "<=",  // passed
		// "between": "between",
		"inq":   "in",       // passed
		"nin":   "not in",   // passed
		"neq":   "<>",       // passed
		"like":  "like",     // passed
		"nlike": "not like", // passed
	}
)

// CRUD
const (
	Create     = "Create"
	Upsert     = "Upsert"
	Find       = "Find"
	FindByID   = "FindByID"
	FindOne    = "FindOne"
	Count      = "Count"
	Exists     = "Exists"
	Update     = "Update"
	DeleteByID = "DeleteByID"
)
