package main_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"

	"github.com/joho/godotenv"
	"golang.org/x/sys/windows"
)

var _ error = godotenv.Load(".env")
var FILEPATH = os.Getenv("FILE_PATH")

func TestMain(t *testing.T) {
	f, err := os.Open(FILEPATH)
	if err != nil {
		t.Fatal("this is OpenFile error: ", err)
	}
	list, err := f.ReadDir(-1)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range list {
		t.Log(l.Name())
	}
	f.Close()
}

type _FileInfo struct {
	windows.Win32finddata
	handle windows.Handle
}

func (fi *_FileInfo) Name() string {
	return windows.UTF16ToString(fi.FileName[:])
}

func findFirst(pattern string) (*_FileInfo, error) {
	pattern16, err := windows.UTF16PtrFromString(pattern)
	if err != nil {
		return nil, err
	}
	this := new(_FileInfo)
	this.handle, err = windows.FindFirstFile(pattern16, &this.Win32finddata)
	if err != nil {
		return nil, err
	}
	return this, nil
}

type windows32finddata struct{ syscall.Win32finddata }

func (fi *_FileInfo) close() {
	windows.FindClose(fi.handle)
}
func (wd *windows32finddata) Attribute() uint32 {
	return wd.FileAttributes
}

func (wd *windows32finddata) IsDir() bool {
	return (wd.Attribute() & windows.FILE_ATTRIBUTE_DIRECTORY) != 0
}
func (wd *windows32finddata) IsFile() bool {
	return (wd.Attribute() & windows.FILE_ATTRIBUTE_NORMAL) != 0
}

type FileInfo = _FileInfo

func (fi *_FileInfo) clone() *FileInfo {
	return &_FileInfo{fi.Win32finddata, fi.handle}
}
func (fi *_FileInfo) findNext() error {
	return windows.FindNextFile(fi.handle, &fi.Win32finddata)
}
func Walk(pattern string) error {
	return walk(pattern)
}
func walk(pattern string) error {
	this, err := findFirst(pattern)
	if err != nil {
		return err
	}
	// _pattern := strings.ToUpper(filepath.Base(pattern))
	defer this.close()
	for {
		// if ctx != nil {
		// 	select {
		// 	case <-ctx.Done():
		// 		return ctx.Err()
		// 	default:
		// 	}
		// }
		// _name := strings.ToUpper(this.Name())
		// matched, err := filepath.Match(_pattern, _name)
		// if err == nil && matched {
		// 	if !callback(this.clone()) {
		// 		return nil
		// 	}
		// }
		log.Println(this.Name())
		if err := this.findNext(); err != nil {
			return nil
		}
	}
}
func TestWalk(t *testing.T) {
	// path := FILEPATH + "*"
	// Walk(context.TODO(),)
	Walk("*")
}
func TestCall(t *testing.T) {
	path := windows.StringToUTF16Ptr(FILEPATH)
	data := &syscall.Win32finddata{}
	handle, err := syscall.FindFirstFile(path, data)
	if err != nil {
		t.Fatal("FindFirstFile error:", err)
	}
	win32FinddataList := []syscall.Win32finddata{}
	win32FinddataList = append(win32FinddataList, *data)
	for {
		win32FinddataList = append(win32FinddataList, *data)
		if err := syscall.FindNextFile(handle, data); err != nil {
			break
			// t.Fatal("FindNextFile error:", err)
		}
	}
	for _, w := range win32FinddataList {
		wd := windows32finddata{w}
		log.Println(windows.UTF16ToString(w.FileName[:]), wd.IsDir())
	}
}
func TestDigDirTree(t *testing.T) {
	path := FILEPATH
	// w32fdl, err := digDirTree(path)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// b, err := json.MarshalIndent(w32fdl, "", "  ")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(string(b))
	// for _, v := range w32fdl.Win32finddata {
	// 	fmt.Println(v.FileNameStr)
	// }
	folder, err := diggingDirectory(path)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(folder, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}
func diggingDirectory(path string) (*Folder, error) {
	resultFolder := &Folder{}
	currentFolder := &Folder{}
	// path := ""
	for {
		win32FindDataList, err := digDirTree(path)
		path := win32FindDataList.Path
		if err != nil {
			err := fmt.Errorf("digDirTree error: %v", err)
			return nil, err
		}
		for _, data := range win32FindDataList.Win32finddata {
			if data.IsDir() {
				fileName := windows.UTF16ToString(data.FileName[:])
				fileName = path + "/" + fileName

				newFolder := NewFolder(windows.UTF16ToString(data.FileName[:]) /*, data*/)
				resultFolder.AppendFolder(newFolder)
				currentFolder = newFolder
				path = fileName
			}
			if !data.IsDir() {
				file := NewFile(windows.UTF16ToString(data.FileName[:]) /*, data*/)
				resultFolder.AppendFile(file)
				// currentFolder =
			}
			// currentFolder.AppendFolder()
		}
		// if win32FindDataList.Win32finddata
	}
	fmt.Println(currentFolder)
	return nil, nil
}

func (wd *win32FindData) Attribute() uint32 {
	return wd.FileAttributes
}

func (wd *win32FindData) IsDir() bool {
	return (wd.Attribute() & windows.FILE_ATTRIBUTE_DIRECTORY) != 0
}
func (wd *win32FindData) IsFile() bool {
	return (wd.Attribute() & windows.FILE_ATTRIBUTE_NORMAL) != 0
}

type win32FindData struct {
	syscall.Win32finddata
	FileNameStr string
}
type win32FindDataList struct {
	Path          string
	Win32finddata []win32FindData
}

func digDirTree(path string) (*win32FindDataList, error) {
	pathUint := windows.StringToUTF16Ptr(path + "/*")
	data := syscall.Win32finddata{}
	// win32fd := win32FindData{}
	win32fdl := &win32FindDataList{}
	handle, err := syscall.FindFirstFile(pathUint, &data)
	if err != nil {
		err = fmt.Errorf("FindFirstFile error: %v", err)
		return nil, err
	}
	// win32FinddataList := []syscall.Win32finddata{}

	for {
		if err := syscall.FindNextFile(handle, &data); err != nil {
			break
			// t.Fatal("FindNextFile error:", err)
		}

		name := windows.UTF16ToString(data.FileName[0:])
		if name == "." || name == ".." {
			continue
		}

		w32f := &win32FindData{
			Win32finddata: data,
			FileNameStr:   name,
		}
		win32fdl.Win32finddata = append(win32fdl.Win32finddata, *w32f)
		// fileName := windows.UTF16ToString(data.FileName[:])
		// fileName = path + "/" + fileName
		win32fdl.Path = path

		// _, err = syscall.FindFirstFile(windows.StringToUTF16Ptr(fileName), &data)
		// if err != nil {
		// 	err = fmt.Errorf("FindFirstFile error: %v", err)
		// 	return nil, nil, err
		// }
		// _, win32finddataListInner, err := ff.dirTreeWalk(fileName)
		// if err != nil {
		// 	return nil, nil, err
		// }
		// // win32FinddataList = append(win32FinddataList, *data)
		// // log.Println("this is file name: ", fileName)

		// for _, w := range win32finddataListInner {
		// 	// wd := windows32finddata{w}
		// 	// log.Println(windows.UTF16ToString(w.FileName[:]), wd.IsDir())
		// 	win32FinddataList = append(win32FinddataList, w)

		// }
		// newFolder := NewFolder(windows.UTF16ToString(data.FileName[:]) /*, data*/)
		// currentFolder.AppendFolder(newFolder)
		// currentFolder = newFolder
		// continue

	}
	// for _, w := range win32FinddataList {
	// 	wd := windows32finddata{w}
	// 	log.Println(windows.UTF16ToString(w.FileName[:]), wd.IsDir())
	// }
	// return ff, win32FinddataList, nil
	return win32fdl, nil
}

// var ff *Folder

var currentFolder *Folder = &Folder{}

func (ff *Folder) dirTreeWalk(path string) (*Folder, []syscall.Win32finddata, error) {
	// currentFolder := &Folder{}
	// currentFolder := ff
	pathUint := windows.StringToUTF16Ptr(path + "/*")
	data := syscall.Win32finddata{}
	handle, err := syscall.FindFirstFile(pathUint, &data)
	if err != nil {
		err = fmt.Errorf("FindFirstFile error: %v", err)
		return nil, nil, err
	}
	win32FinddataList := []syscall.Win32finddata{}
	// win32FinddataList = append(win32FinddataList, *data)
	// pathMask := ""
	for {
		if err := syscall.FindNextFile(handle, &data); err != nil {
			break
			// t.Fatal("FindNextFile error:", err)
		}

		name := windows.UTF16ToString(data.FileName[0:])
		if name == "." || name == ".." {
			continue
		}

		w32f := &windows32finddata{data}
		// if w32f.IsFile() {
		if !w32f.IsDir() {
			file := NewFile(windows.UTF16ToString(data.FileName[:]) /*, data*/)
			currentFolder.AppendFile(file)
			currentFolder = ff
			win32FinddataList = append(win32FinddataList, data)
			continue
		}
		if w32f.IsDir() {
			fileName := windows.UTF16ToString(data.FileName[:])
			fileName = path + "/" + fileName
			_, err = syscall.FindFirstFile(windows.StringToUTF16Ptr(fileName), &data)
			if err != nil {
				err = fmt.Errorf("FindFirstFile error: %v", err)
				return nil, nil, err
			}
			_, win32finddataListInner, err := ff.dirTreeWalk(fileName)
			if err != nil {
				return nil, nil, err
			}
			// win32FinddataList = append(win32FinddataList, *data)
			// log.Println("this is file name: ", fileName)

			for _, w := range win32finddataListInner {
				// wd := windows32finddata{w}
				// log.Println(windows.UTF16ToString(w.FileName[:]), wd.IsDir())
				win32FinddataList = append(win32FinddataList, w)

			}
			newFolder := NewFolder(windows.UTF16ToString(data.FileName[:]) /*, data*/)
			currentFolder.AppendFolder(newFolder)
			currentFolder = newFolder
			continue
		}
	}
	// for _, w := range win32FinddataList {
	// 	wd := windows32finddata{w}
	// 	log.Println(windows.UTF16ToString(w.FileName[:]), wd.IsDir())
	// }
	// return ff, win32FinddataList, nil
	return currentFolder, win32FinddataList, nil
}
func TestTreeWalk(t *testing.T) {
	path := FILEPATH
	// win32FinddataList, err := dirWalk(path)
	// data := syscall.Win32finddata{}
	f := NewFolder("root" /*, data*/)
	// currentFolder,
	_, win32FinddataList, err := f.dirTreeWalk(path)
	if err != nil {
		t.Fatal(err)
	}
	// log.Println(currentFolder)
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
	// log.Println(win32Finddata)
	for _, w := range win32FinddataList {
		wd := windows32finddata{w}
		log.Println(windows.UTF16ToString(w.FileName[:]), wd.IsDir())
		w32f := &windows32finddata{w}
		if w32f.IsDir() {

		}
	}

}
func dirWalk(path string) ([]syscall.Win32finddata, error) {
	// ourrentFolder := &Folder{}
	pathUint := windows.StringToUTF16Ptr(path + "/*")
	data := &syscall.Win32finddata{}
	handle, err := syscall.FindFirstFile(pathUint, data)
	if err != nil {
		err = fmt.Errorf("FindFirstFile error: %v", err)
		return nil, err
	}
	win32FinddataList := []syscall.Win32finddata{}
	// win32FinddataList = append(win32FinddataList, *data)
	for {
		if err := syscall.FindNextFile(handle, data); err != nil {
			break
			// t.Fatal("FindNextFile error:", err)
		}

		name := windows.UTF16ToString(data.FileName[0:])
		if name == "." || name == ".." {
			continue
		}

		win32FinddataList = append(win32FinddataList, *data)
	}
	return win32FinddataList, nil
}

type Folder struct {
	Name string
	// Win32filedata syscall.Win32finddata
	Folders []*Folder
	Files   []*File
}

// type File struct {
// 	Name string
// 	// Win32filedata syscall.Win32finddata
// }

func NewFolder(name string /*, win32filedata syscall.Win32finddata*/) *Folder {
	return &Folder{
		Name: name,
		// Win32filedata: win32filedata,
		Folders: []*Folder{},
		Files:   []*File{},
	}
}

func (f *Folder) FindFolder(ff *Folder) *Folder {
	for _, folder := range f.Folders {
		if folder.Name == ff.Name {
			return folder
		}
	}
	return nil
}

func (f *Folder) AppendFolder(ff *Folder) {
	f.Folders = append(f.Folders, ff)
}

func (f *Folder) AppendFile(fl *File) {
	f.Files = append(f.Files, fl)
}

func NewFile(name string /*, win32filedata syscall.Win32finddata*/) *File {
	return &File{
		Name: name,
		// Win32filedata: win32filedata,
	}
}
func DeepEqualJSON(j1, j2 string) (error, bool) {
	var err error

	var d1 interface{}
	err = json.Unmarshal([]byte(j1), &d1)

	if err != nil {
		return err, false
	}

	var d2 interface{}
	err = json.Unmarshal([]byte(j2), &d2)

	if err != nil {
		return err, false
	}

	if reflect.DeepEqual(d1, d2) {
		return nil, true
	} else {
		return nil, false
	}
}

type Directory struct {
	Name    string
	Path    string
	Files   []*File
	Folders []*Directory
}

type File struct {
	Name string
}

func buildTree(root string) (*Directory, error) {
	rootUtf16 := syscall.StringToUTF16Ptr(root + "\\*")
	dirHandle, err := syscall.FindFirstFile(rootUtf16, &syscall.Win32finddata{})
	if dirHandle == syscall.InvalidHandle {
		return nil, err
	}
	defer syscall.FindClose(dirHandle)

	d := &Directory{Name: root, Path: root}
	for {
		findData := &syscall.Win32finddata{}
		// findDataPointer := uintptr(unsafe.Pointer(findData))
		if err := syscall.FindNextFile(dirHandle, findData); err != nil {
			return nil, err
		}
		if findData.FileName[0] == '.' {
			continue
		}
		if findData.FileAttributes&syscall.FILE_ATTRIBUTE_DIRECTORY != 0 {
			path := root + "\\" + syscall.UTF16ToString(findData.FileName[:])
			subDir, err := buildTree(path)
			if err != nil {
				return nil, err
			}
			d.Folders = append(d.Folders, subDir)
		} else {
			d.Files = append(d.Files, &File{Name: syscall.UTF16ToString(findData.FileName[:])})
		}
	}
	return d, nil
}
func TestGetFiles(t *testing.T) {
	// d, err := buildTree(FILEPATH)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(d)
	// files, err := getFiles(FILEPATH + "*")
	root, err := getFiles("./test")
	if err != nil {
		fmt.Println(err)
	}
	printTree(root, "")
}
func getFiles(path string) (*Node, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	path = filepath.Clean(path)
	var handle syscall.Handle
	var findData syscall.Win32finddata
	handle, err = syscall.FindFirstFile(syscall.StringToUTF16Ptr(path), &findData)
	if err != nil {
		return nil, err
	}
	defer syscall.FindClose(handle)

	node := &Node{Name: syscall.UTF16ToString(findData.FileName[:]), IsDir: findData.FileAttributes&syscall.FILE_ATTRIBUTE_DIRECTORY != 0}
	for {
		err := syscall.FindNextFile(handle, &findData)
		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}
			return nil, err
		}
		if findData.FileName[0] == '.' {
			continue
		}
		child := &Node{Name: syscall.UTF16ToString(findData.FileName[:]), IsDir: findData.FileAttributes&syscall.FILE_ATTRIBUTE_DIRECTORY != 0}
		if child.IsDir {
			childPath := path + "\\" + child.Name
			grandchild, err := getFiles(childPath)
			if err != nil {
				return nil, err
			}
			child.Children = grandchild.Children
		}
		node.Children = append(node.Children, child)
	}
	return node, nil
}

type Node struct {
	Name     string
	Children []*Node
	IsDir    bool
}

func printTree(node *Node, prefix string) {
	if node.IsDir {
		fmt.Println(prefix + "|-- " + node.Name)
		for i, child := range node.Children {
			if i == len(node.Children)-1 {
				printTree(child, prefix+" ")
			} else {
				printTree(child, prefix+"| ")
			}
		}
	} else {
		fmt.Println(prefix + "|-- " + node.Name)
	}
}
