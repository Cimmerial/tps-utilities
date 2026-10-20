package img

// OverwriteDecider is called when dest already exists and is not the source path.
// Return true to replace dest, false to skip this file without error.
type OverwriteDecider func(dest, src string) (bool, error)

// ConvertOptions controls naming, overwrite, and delete-after-success.
type ConvertOptions struct {
	Basename          string // empty = derive from source
	DeleteOriginal    bool
	IfExists          OverwriteDecider // nil = overwrite without asking
	DuplicateBasename DuplicateBasenameMode
	BatchPlan         BatchPlan // filled by MassConvert / caller for multi-file
}

// ConvertStatus is the outcome of a single conversion.
type ConvertStatus int

const (
	ConvertOK       ConvertStatus = 0 // wrote or copied output
	ConvertNoChange ConvertStatus = 1 // same-format, same path, left untouched
	ConvertSkipped  ConvertStatus = 2 // user declined overwrite
)

// ConvertOutcome is returned from ConvertFile.
type ConvertOutcome struct {
	OutPath string
	Note    string
	Status  ConvertStatus
}
