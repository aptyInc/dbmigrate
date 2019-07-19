package source

//FileReader utlilty helps you read file system
type FileReader interface {
	ReadDirs(root string) ([]string, error)
	ReadFilesWithExtension(root string, extension string) ([]string, error)
	ReadFileAsString(path string) (string, error)
}
