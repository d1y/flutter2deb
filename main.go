package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/d1y/flutter2deb/pkg/gitconfig"
)

var _currentPath = ""

type templateData struct {
	Name        string
	Description string
	Version     string
	GitUsername string
	GitEmail    string
}

var tData *templateData

func genControlTemplate() string {
	var p = tData.Name

	// ** package name has characters that aren't lowercase alphanums or '-+.' **
	// the package support - | . | +
	// `pubspec.yaml` the package support `_`, replace
	p = strings.ReplaceAll(p, "_", "-")

	var b = fmt.Sprintf(`Source: %s
Section: unknown
Priority: optional
Maintainer: %s <%s>
Version: %s
Package: %s
Architecture: any
Description: %s

`, tData.Name, tData.GitUsername, tData.GitEmail, tData.Version, p, tData.Description)
	return b
}

func genDesktopTemplate(ext string) string {
	var b = fmt.Sprintf(`[Desktop Entry]
Name=%s
Comment=%s
Exec=/usr/bin/%s
Icon=/usr/share/icons/%s%s
Terminal=false
Type=Application
X-Ubuntu-Touch=true
Categories=Development`, tData.Name, tData.Description, tData.Name, tData.Name, ext)
	return b
}

// create template
func createTemplate() {
	// 1. 拿到需要操作的目录
	//   - 先从 os.Args, 若不存在就从 os.Getwd() 拿到
	//   - 判断模板是否已经存在(<debian> 目录是否已经存在)
	//   - 已经存在就跳过生成
	// 2. 生成目录
	//   - debian
	//   - ++DEBIAN
	//   - ++++control
	//   - ++usr
	//   - ++++bin
	//   - ++++share
	//   - ++++++applications
	//   - ++++++++name.desktop
	//   - ++++++icons
	//   - ++++++++name.png
	// 3. 从根目录中查找 `app.png` | `app.jpg` | `app.svg`
	// 4. **从`pubspec.yaml`中拿到模板字符串**
	//   - name
	//   - description
	//   - version
	// 5. 生成过程
	//
	// << `control` 文件 >>
	// Source: ${name}
	// Section: unknown
	// Priority: optional
	// Maintainer: ${username} <${email}>
	// Version: ${version}
	// Package: ${name}
	// Architecture: any
	// Description: ${description}
	//
	// << `*.desktop` 文件 >>
	// [Desktop Entry]
	// Name=${name}
	// Comment=${description}
	// Exec=/usr/bin/${name}
	// Icon=/usr/share/icons/${name}.png
	// Terminal=false
	// Type=Application
	// X-Ubuntu-Touch=true
	// Categories=Development
	//
	// 根据<3>将`app.png`文件复制到 share/icons/${name}

	icon, err := findAppIcon(_currentPath)
	checkPanic(err)

	var c = genControlTemplate()
	var d = genDesktopTemplate(path.Ext(icon.Name()))
	var __control = path.Join(_currentPath, "debian/DEBIAN/control")
	var __desktop = path.Join(_currentPath, "debian/usr/share/applications/"+tData.Name+".desktop")
	ensureDir(__control)
	ensureDir(__desktop)
	ioutil.WriteFile(__control, []byte(c), 0755)
	ioutil.WriteFile(__desktop, []byte(d), 0755)

	ensureDir(path.Join(_currentPath, "debian/usr/bin/.gitkeep"))

	var iconFilePath = path.Join(_currentPath, icon.Name())
	var iconDistPath = path.Join(_currentPath, "debian/usr/share/icons/"+tData.Name+path.Ext(icon.Name()))
	ensureDir(iconDistPath)
	CopyFile(iconFilePath, iconDistPath)

	var buildsh = path.Join(_currentPath, "build.sh")
	ioutil.WriteFile(buildsh, []byte(`flutter build linux

rm -rf debian/usr/bin/*

cp -rf build/linux/x64/release/bundle/* debian/usr/bin/

dpkg -b debian build/spark_store.deb`), 0755)
}

func main() {
	createTemplate()
}

func init() {
	username, err := gitconfig.Username()
	checkPanic(err)
	email, err := gitconfig.Email()
	checkPanic(err)
	_c, err := getCurrentPath()
	checkPanic(err)
	_currentPath = _c
	var yaml = path.Join(_currentPath, "pubspec.yaml")
	parse, err := NewParsePubFile(yaml)
	checkPanic(err)
	tData = &templateData{}
	tData.Name = parse.Pub.Name
	tData.Description = parse.Pub.Description
	tData.Version = parse.Pub.Version
	tData.GitUsername = username
	tData.GitEmail = email
}
