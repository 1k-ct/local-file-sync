package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/sys/windows"
)

var _ error = godotenv.Load(".env")
var FILEPATH = os.Getenv("FILE_PATH")

func main() {
	path, _ := windows.UTF16PtrFromString(FILEPATH + "*")
	win32finddata := &windows.Win32finddata{}
	h, err := windows.FindFirstFile(path, win32finddata)
	if err != nil {
		log.Fatal("FindFirstFile error:", err)
	}
	win32finddata2 := &windows.Win32finddata{
		FileAttributes:    0,
		CreationTime:      windows.Filetime{},
		LastAccessTime:    windows.Filetime{},
		LastWriteTime:     windows.Filetime{},
		FileSizeHigh:      0,
		FileSizeLow:       0,
		Reserved0:         0,
		Reserved1:         0,
		FileName:          win32finddata.FileName,
		AlternateFileName: [13]uint16{},
	}
	if err := windows.FindNextFile(h, win32finddata); err != nil {
		log.Fatal("FindNextFile", err)
	}
	log.Println(windows.UTF16ToString(win32finddata2.FileName[0:]))
	log.Println("this log is fileName: ", windows.UTF16ToString(win32finddata.AlternateFileName[:]))
	// log.Println(findFileHanlde)
	// log.Println(windows.UTF16ToString(win32finddata.LastAccessTime))
	// log.Println(win32finddata.LastWriteTime.Nanoseconds())
	// t := time.Unix(0, win32finddata.LastWriteTime.Nanoseconds())
	// log.Println(t.String())
	log.Println(windows.UTF16ToString(win32finddata.FileName[:]))

	// var windowsToken windows.Token = 0
	// dir, _ := windows.UTF16PtrFromString(FILEPATH)
	// // dirLen , _:= windows.
	// var u uint32 = 100
	// // token:= windows.GetCurrentProcessToken()
	// if err := windows.GetUserProfileDirectory(windowsToken, dir, &u); err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(windowsToken.GetUserProfileDirectory())
	// // log.Println(windows.UTF16ToString(win32finddata))
	// t := time.Unix(0, win32finddata.LastAccessTime.Nanoseconds())
	// log.Println(t.String())
}
func createFile() {
	name, _ := windows.UTF16PtrFromString(FILEPATH + "text2.txt")
	var access uint32 = windows.GENERIC_WRITE
	var mode uint32 = 0
	sa := &windows.SecurityAttributes{}
	var createmode uint32 = windows.CREATE_NEW
	var attrs uint32 = windows.FILE_ATTRIBUTE_NORMAL
	// var templatefile uintptr = &windows.Handle{}
	var templatefile windows.Handle = 0
	handle, err := windows.CreateFile(name, access, mode, sa, createmode, attrs, templatefile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(handle)
}
