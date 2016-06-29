package databind

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//AssetMap provides a map of paths that contain assets of the specific filepaths
type AssetMap map[string]string

// ReloadAssetMap reloads the files into the map skipping the already found ones
func ReloadAssetMap(tree AssetMap, dir string, ext []string, skip []string) error {
	var stat os.FileInfo
	var err error

	//do the path exists
	if stat, err = os.Stat(dir); err != nil {
		return err
	}

	//do we have a directory?
	if !stat.IsDir() {

		if tree.Has(filepath.ToSlash(dir)) {
			return nil
		}

		var fext string
		var rel = filepath.Base(dir)

		if strings.Index(rel, ".") != -1 {
			fext = filepath.Ext(rel)
		}

		if len(ext) > 0 {
			if hasExt(ext, fext) {
				tree[rel] = filepath.ToSlash(dir)
			}
		} else {
			tree[rel] = filepath.ToSlash(dir)
		}

	} else {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

			//if info is nil or is a directory when we skip
			if info == nil || info.IsDir() {
				return nil
			}

			repath := filepath.ToSlash(path)

			if tree.Has(repath) {
				return nil
			}

			if strings.Contains(repath, ".git") {
				return nil
			}

			if strings.Index(repath, ".git") != -1 {
				return nil
			}

			if hasIn(skip, repath) {
				return nil
			}

			var rel string
			var rerr error

			//is this path relative to the current one, if not,return err
			if rel, rerr = filepath.Rel(dir, path); rerr != nil {
				return rerr
			}

			var fext string

			if strings.Index(rel, ".") != -1 {
				fext = filepath.Ext(rel)
			}

			if len(ext) > 0 {
				if hasExt(ext, fext) {
					tree[rel] = filepath.ToSlash(path)
				}
			} else {
				tree[rel] = filepath.ToSlash(path)
			}

			return nil
		})
	}
	return nil
}

// AssetTree provides a map tree of files across the given directory that match the filenames being used
func AssetTree(dir string, ext, skip []string) (AssetMap, error) {
	//do the path exists
	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}

	var tree = make(AssetMap)

	err := ReloadAssetMap(tree, dir, ext, skip)

	return tree, err
}

// AssetLoader is a function type which returns the content in []byte of a specific asset
type AssetLoader func(string) ([]byte, error)

// Has returns true/false if the filename without its extension exists
func (am AssetMap) Has(name string) bool {
	_, ok := am[name]
	return ok
}

// Load returns the data of the specific file with the name
func (am AssetMap) Load(name string) ([]byte, error) {
	if !am.Has(name) {
		return nil, fmt.Errorf("AssetMap: %s", fmt.Sprintf("%s unknown", name))
	}
	return ioutil.ReadFile(am[name])
}

// PathValidator provides a type for validating a path and info set
type PathValidator func(string, os.FileInfo) bool

//PathMux provides a function type for mixing a new file path
type PathMux func(string, os.FileInfo) string

func defaultValidator(_ string, _ os.FileInfo) bool {
	return true
}

func defaultMux(n string, _ os.FileInfo) string {
	return n
}

// BasicAssetTree represent a directory structure and its corresponding assets
type BasicAssetTree struct {
	Dir      string
	ModDir   string
	AbsDir   string
	Info     os.FileInfo
	Tree     *MapWriter
	Ml       sync.RWMutex
	Children []*BasicAssetTree
	root     bool
}

// String returns a formatted string response of how a BasicAssetTree should look when printed
func (b *BasicAssetTree) String() string {
	return fmt.Sprintf("AssetTree for -> %s at %s with Total Files of %d \n", b.Dir, b.AbsDir, b.Tree.Size())
}

// Add adds a BasicAssetTree into the children lists
func (b *BasicAssetTree) Add(bs *BasicAssetTree) {
	b.Ml.Lock()
	defer b.Ml.Unlock()
	b.Children = append(b.Children, bs)
}

// Flush empties out the tree
func (b *BasicAssetTree) Flush() {
	b.Tree.Flush()
	b.Children = nil
}

// Delete removes this item from the lists
func (b *BasicAssetTree) Delete(bs *BasicAssetTree) {
	var index = -1
	var size int

	b.Ml.RLock()
	size = len(b.Children)
	for ind, item := range b.Children {
		if item == bs {
			index = ind
			break
		}
	}
	b.Ml.RUnlock()

	// not found,so we skip
	if index < 0 {
		return
	}

	//we found one,so we cut and reset children list
	b.Ml.Lock()
	copy(b.Children[:index], b.Children[index+1:])
	b.Children = b.Children[:(size - 1)]
	b.Ml.Unlock()
}

// EachChild iterates through the children dir in this tree
func (b *BasicAssetTree) EachChild(fx func(bs *BasicAssetTree)) {
	if fx == nil {
		return
	}
	b.Ml.RLock()
	defer b.Ml.RUnlock()
	for _, child := range b.Children {
		fx(child)
	}
}

// EmptyAssetTree returns a new AssetTree based of the given path
func EmptyAssetTree(path, mod string, info os.FileInfo, abs string) *BasicAssetTree {
	as := BasicAssetTree{
		Dir:      path,
		ModDir:   mod,
		AbsDir:   abs,
		Info:     info,
		Tree:     NewMapWriter(make(AssetMap)),
		Children: make([]*BasicAssetTree, 0),
	}

	return &as
}

// BuildAssetPath reloads the files into the map skipping the already found ones
func BuildAssetPath(base string, files []os.FileInfo, dirs *TreeMapWriter, pathtree *BasicAssetTree, validator PathValidator, mux PathMux) error {
	// ws.Add(1)
	// defer ws.Done()

	pathtree.Tree.Flush()

	for _, pitem := range files {

		//get the file path using the provided base
		dir := filepath.Join(base, pitem.Name())

		if !validator(dir, pitem) {
			continue
		}

		// create a BasicAssetTree for this path

		if !pitem.IsDir() {
			// log.Printf("will mux path: %+s ", mux)
			pmx := mux(dir, pitem)
			if !pathtree.Tree.Has(pmx) {
				pathtree.Tree.Add(pmx, dir)
			}
			// pathtree.Tree[dir] = mux(dir,pitem)
			continue
		}

		var target *BasicAssetTree
		var muxdir = mux(dir, pitem)

		//open up the filepath since its a directory, read and sort
		dfiles, err := getDirListings(dir)

		// did an erro occur while trying to get the directories,if so remove if exists and continue
		if err != nil {
			target = dirs.Get(muxdir)
			dirs.Delete(muxdir)

			if target != nil {
				target.Flush()
				target = nil
			}

			continue
		}

		if dirs.Has(muxdir) {
			target = dirs.Get(muxdir)
		} else {
			tabsDir, _ := filepath.Abs(dir)
			target = EmptyAssetTree(dir, muxdir, pitem, tabsDir)

			//add into the global dir listings
			dirs.Add(muxdir, target)

			//add to parenttree as a root dir
			pathtree.Add(target)
		}

		// var directories []os.FileInfo

		//do another build but send into go-routine
		BuildAssetPath(dir, dfiles, dirs, target, validator, mux)
	}

	return nil
}

// LoadTree loads the path into the assettree updating any found trees along
func LoadTree(dir string, tree *TreeMapWriter, fx PathValidator, fxm PathMux) error {
	if fx == nil {
		fx = defaultValidator
	}

	// use defaultMux if non is set
	if fxm == nil {
		fxm = defaultMux
	}

	var st os.FileInfo
	var err error

	if st, err = os.Stat(dir); err != nil {
		return err
	}

	if !st.IsDir() {
		return os.ErrInvalid
	}

	//grab the absolute path of the dir
	// abs, _ := filepath.Abs(dir)

	files, err := getDirListings(dir)

	//unable to retrieve directory list, return tree and error
	if err != nil {
		return err
	}

	//get the absolute path for the path
	absdir, _ := filepath.Abs(dir)

	var cur *BasicAssetTree
	var muxcur = fxm(dir, st)

	if tree.Has(muxcur) {
		cur = tree.Get(muxcur)
	} else {
		//create the assettree for this path
		cur = EmptyAssetTree(dir, muxcur, st, absdir)
		cur.root = true

		//register and mux the path as requried to the super tree as a directory tree
		tree.Add(muxcur, cur)
		// tree.Add("/", cur)

		// tree[fxm(dir, st)] = cur
		//register into the parent tree as  child tree
		// parentTree.Add(cur)
	}

	// log.Printf("cur at %s -> %s", cur, muxcur)

	//ok lets build the path children
	if err := BuildAssetPath(dir, files, tree, cur, fx, fxm); err != nil {
		return err
	}

	return nil
}

// Assets returns a tree map listing of a specific path and if an error occured before the map was created it will return a nil and an error but if the path was valid but its children were faced with an invalid state then it returns the map and an error
func Assets(dir string, fx PathValidator, fxm PathMux) (*TreeMapWriter, error) {
	// use defaultValidator if non is set
	if fx == nil {
		fx = defaultValidator
	}

	// use defaultMux if non is set
	if fxm == nil {
		fxm = defaultMux
	}

	//since we will be using some go-routines to improve performance,lets create a waitgroup
	// var wsg = new(sync.WaitGroup)
	var tmap = make(AssetTreeMap)
	var tree = NewTreeMapWriter(tmap)

	err := LoadTree(dir, tree, fx, fxm)

	// wsg.Wait()

	if err != nil {
		return nil, err
	}

	return tree, nil
}

// DirListing provides a struct that can generate a map of listings of a path
type DirListing struct {
	dir       string
	Listings  *TreeMapWriter
	validator PathValidator
	mux       PathMux
}

// DirListings returns a new DirListings for the set path or returns an error if that path does not exists or is not a directory path
func DirListings(path string, valid PathValidator, mux PathMux) (*DirListing, error) {
	tree, err := Assets(path, valid, mux)

	if err != nil {
		return nil, err
	}

	dls := DirListing{
		dir:       path,
		Listings:  tree,
		validator: valid,
		mux:       mux,
	}

	return &dls, nil
}

// EachDir calls the internal listings EachDir function
func (dls *DirListing) EachDir(fx func(b *BasicAssetTree, realPath string)) {
	dls.Listings.Each(fx)
}

// EachDirFiles calls the internal listings EachDirFiles function
func (dls *DirListing) EachDirFiles(fx func(b *BasicAssetTree, moddedPath, realPath string)) {
	dls.Listings.EachDirFiles(fx)
}

// Dir returns the path of this listings
func (dls *DirListing) Dir() string {
	return dls.dir
}

// Reload reloads this directory and its listing structures with updates
func (dls *DirListing) Reload() error {
	return LoadTree(dls.dir, dls.Listings, dls.validator, dls.mux)
}
