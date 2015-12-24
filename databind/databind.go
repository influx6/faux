package databind

const (
	lowerhex = "0123456789abcdef"

	packageDetails = `//Auto-generated from github.com/influx6/assets
// DO NOT CHANGE

// Package %s provides an auto-generated static embeding of data files within the specific directory %s
package %s

  `

	rootDir = `
// RootDirectory defines a directory root for these virtual files
var RootDirectory = NewDirCollector()

`

	rootInit = `
func init(){
%s
}

`

	subRegister = `
	dir.AddDirectory(%q,func() *VDir{
		return RootDirectory.Get(%q)
	})

`

	dirRegister = `
  RootDirectory.Set(%q,func() *VDir{
    var dir = NewVDir(%q,%q,%q,%t)

    // register the sub-directories
    {{ subs }}

    // register the files
    {{ files }}

    return dir
  }())
`

	debugFile = `
		dir.AddFile(NewVFile(%q,%q,%q,%d,%t,%t,%s))
	`

	prodRead = `func(v *VFile) ([]byte,error) {
	    return readData(v,[]byte(%s))
	  }`

	comfileRead = `func(v *VFile) ([]byte, error) {
			fo, err := os.Open(v.RealPath())
			if err != nil {
				return nil, fmt.Errorf("---> assets.readFile: Error reading file: %s at %s: %s\n", v.Name(), v.RealPath(), err)
			}

			defer fo.Close()

			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)

			_, err = io.Copy(gz, fo)
			gz.Close()

			if err != nil {
				return nil, fmt.Errorf("---> assets.readFile.gzip: Error gzipping file: %s at %s: %s\n", v.Name(), v.RealPath(), err)
			}

			return buf.Bytes(), nil
		}`

	fileRead = `func(v *VFile) ([]byte, error) {
			fo, err := os.Open(v.RealPath())
			if err != nil {
				return nil, fmt.Errorf("---> assets.readFile: Error reading file: %s at %s: %s\n", v.Name(), v.RealPath(), err)
			}

			defer fo.Close()

			var buf bytes.Buffer

			_, err = io.Copy(&buf,fo)
			if err != nil && err != io.EOF {
				return nil, err
			}

			return buf.Bytes(), nil
		}`

	comFunc = `
func readData(v *VFile,data []byte)([]byte,error){
	if %t {
		return readVData(v,data)
	}
	
	return readEData(v,data)
}
`
)
