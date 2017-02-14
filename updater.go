package main

import "log"
import "github.com/probandula/figlet4go"
import "runtime"
import "fmt"
import "io/ioutil"
import "strings"
import "os"
import "os/exec"
import "github.com/kardianos/osext"

import (
	"io"
	"net/http"
)

func downloadFile(filepath string, url string) (err error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// Create the file
		out, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer out.Close()

		// Get the data
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Writer the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func doCommand(cmd string, args []string) {
	fmt.Println("C>", cmd, args)
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "IO> %v", string(out))
		fmt.Fprintf(os.Stderr, "E> %v", err)
		//os.Exit(1)
	}
	if string(out) != "" {
		fmt.Fprintf(os.Stderr, "O> %v\n\n", string(out))
	}
}

func buildGithub(repo string) {
	cmd := "go"
	args := []string{"build", repo}
	fmt.Printf("I> Building %v\n", repo)
	doCommand(cmd, args)
}

func installGithub(repo string) {
	cmd := "go"
	args := []string{"get", "-u", repo}
	fmt.Printf("I> Installing %v\n", repo)
	doCommand(cmd, args)
}

func loadRepos(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		//Do something
	}
	lines := strings.Split(string(content), "\n")
	return lines
}

func unPackGoMacOSX(folderPath string) {
	if runtime.GOOS == "darwin" {
		doCommand("xar", []string{"-xf", "go1.7.5.darwin-amd64.pkg"})
		doCommand("sh", []string{"-c", "cat com.googlecode.go.pkg/Payload | gunzip -dc | cpio -i"})
		os.Setenv("GOROOT", fmt.Sprintf("%v/usr/local/go/", folderPath))
		os.Setenv("PATH", fmt.Sprintf("%v/usr/local/go/bin/:%v", folderPath, os.Getenv("PATH")))
		doCommand("go", []string{"version"})
	}
}

func buildGo() {
	figlet("COMPILING GO")
	cwd, _ := os.Getwd()
	fmt.Println("I> Deleting directory golangCompiler")
	//doCommand("rm", []string{"-r", "golangCompiler"})
	doCommand("git", []string{"clone", "https://go.googlesource.com/go", "golangCompiler"})
	os.Chdir("golangCompiler/src")

	doCommand("git", []string{"checkout", "go1.7.5"})

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		doCommand("bash", []string{"all.bash"})
		os.Chdir(cwd)
		os.Setenv("GOROOT", fmt.Sprintf("%v/golangCompiler/", cwd))
		os.Setenv("PATH", fmt.Sprintf("%v/golangCompiler/bin/:%v", cwd, os.Getenv("PATH")))
	} else {
		doCommand("all.bat", []string{})
		os.Chdir(cwd)
		os.Setenv("GOROOT", fmt.Sprintf("%v\\golangCompiler\\", cwd))
		os.Setenv("PATH", fmt.Sprintf("%v\\golangCompiler\\bin\\:%v", cwd, os.Getenv("PATH")))
	}

}

func printEnv() {
	fmt.Println("I> GOROOT_BOOTSTRAP", os.Getenv("GOROOT_BOOTSTRAP"))
	fmt.Println("I> GOPATH", os.Getenv("GOPATH"))
	fmt.Println("I> PATH", os.Getenv("PATH"))
	fmt.Println("I> GOROOT", runtime.GOROOT())
}

func figlet(s string) string {
	ascii := figlet4go.NewAsciiRender()

	// Adding the colors to RenderOptions
	options := figlet4go.NewRenderOptions()
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		options.FontColor = []figlet4go.Color{
			// Colors can be given by default ansi color codes...
			figlet4go.ColorGreen,
			figlet4go.ColorYellow,
			figlet4go.ColorCyan,
			// ...or by an hex string...
			//figlet4go.NewTrueColorFromHexString("885DBA"),
			// ...or by an TrueColor object with rgb values
			//figlet4go.TrueColor{136, 93, 186},
		}
	}

	renderStr, _ := ascii.RenderOpts(s, options)
	return renderStr
}

func makeWith(optName, srcDir, libName string) string {
	return fmt.Sprintf("--with-%v=%v/%v", optName, srcDir, libName)
}

func makeOpt(optName, optVal string) string {
	return fmt.Sprintf("--%v=%v", optName, optVal)
}

func unTgzLib(lib string) {
	doCommand("tar", []string{"-xzvf", fmt.Sprintf("zips/%v.tar.gz", lib)})
}

func unBzLib(lib string) {
	doCommand("tar", []string{"-xjvf", fmt.Sprintf("zips/%v.tar.bz2", lib)})
}

func standardConfigureBuild(name, buildDir string, args []string) {
	fmt.Println(figlet(name))
	cwd, _ := os.Getwd()
    configurePath := fmt.Sprintf("%v/%v/%v", cwd, name, "configure")
	unBzLib(name)
	unTgzLib(name)
	os.Chdir(name)
	os.Chdir(buildDir)
	doCommand("chmod", []string{"a+rwx", "configure"})
	doCommand(configurePath, args)
	doCommand("make", []string{})
	doCommand("make", []string{"install"})
	os.Chdir(cwd)
	doCommand("rm", []string{"-r", name})
}

func buildGcc(path string) {
	arch := "x86_64"
	targetDir := fmt.Sprintf("%v/fakeRoot", path)
	//srcDir := fmt.Sprintf("%v/src", path)
	os.Chdir(path)
	fmt.Println(figlet("GMP"))
	//doCommand("git", []string{"clone", "https://github.com/bw-oss/gmp"})
	gmpName := "gmp-6.1.2"
	mpfrName := "mpfr-3.1.5"
	mpcName := "mpc-1.0.3"
	gccName := "gcc-6.3.0"
    islName := "isl-0.15"

	standardConfigureBuild(gmpName, ".", []string{"--disable-shared", "--enable-static", makeOpt("prefix", targetDir), makeOpt("build", arch)})
	standardConfigureBuild(mpfrName, ".", []string{"--disable-shared", "--enable-static", makeWith("gmp", targetDir, ""), makeOpt("prefix", targetDir)})
	standardConfigureBuild(mpcName, ".", []string{"--disable-shared", "--enable-static", makeOpt("with-gmp", targetDir), makeOpt("with-mpfr", targetDir), makeOpt("prefix", targetDir)})
	standardConfigureBuild(islName, ".", []string{"--disable-shared", "--enable-static", makeOpt("with-gmp-prefix", targetDir), makeOpt("prefix", targetDir)})
	standardConfigureBuild(gccName, "gcc/objdir", []string{"--enable-languages=c,c++,go", "--disable-shared", "--enable-static", "--disable-multilib", "--disable-shared", "--enable-static", makeWith("gmp", targetDir, ""), makeWith("mpfr", targetDir, ""), makeWith("mpc", targetDir, ""), makeWith("isl", targetDir, ""), makeOpt("prefix", targetDir)})
}

func unSevenZ(SzPath, file string) {
	fmt.Println(SzPath, file)
	doCommand(SzPath, []string{"x", file})
}


/*
func unzipWithPathMake(zipName) {
	fmt.Println(figlet(zipName))
	cwd, _ := os.Getwd()
	os.Chdir(srcDir)
	unSevenZ(SzPath, "../zips/zeromq-4.2.1.zip")
	os.Chdir("zeromq-4.2.1/zeromq-4.2.1/builds/mingw32")
	doCommand("make", []string{})
	os.Chdir(zipName)
}
*/

func fetchBuild(targetDir, name, zip, url, plan string) {
    downloadFile(fmt.Sprintf("zips/%v", zip), url)
    standardConfigureBuild(name, ".", []string{makeOpt("prefix", targetDir), makeOpt("with-sysroot", targetDir) })
}


func main() {
	printEnv()
	fmt.Println(figlet(runtime.GOOS))
    os.Setenv("CFLAGS", "-D_XOPEN_SOURCE=1")
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		os.Exit(1)
	}
	myDir := fmt.Sprintf("%v/goFiles", folderPath)
	zipsDir := fmt.Sprintf("%v/zips", folderPath)
	rootDir := fmt.Sprintf("%v/fakeRoot", folderPath)
	srcDir := fmt.Sprintf("%v/src", folderPath)
	SzDir := fmt.Sprintf("%v/7zip", folderPath)
	SzPath := fmt.Sprintf("%v/7zip/7z.exe", folderPath)
	goDir := fmt.Sprintf("%v/golangCompiler", folderPath)
	fmt.Println("I> Creating", myDir)
	os.Mkdir(myDir, os.ModeDir|0777)
	os.Mkdir(zipsDir, os.ModeDir|0777)
	os.Mkdir(rootDir, os.ModeDir|0777)
	os.Mkdir(SzDir, os.ModeDir|0777)
	os.Mkdir(srcDir, os.ModeDir|0777)
	fmt.Println("Creating ", goDir)
	os.Mkdir(goDir, os.ModeDir|0777)

	downloadFile("zips/zeromq-4.2.1.zip", "https://github.com/zeromq/libzmq/releases/download/v4.2.1/zeromq-4.2.1.zip")
    fetchBuild(rootDir, "libelf-0.8.13", "libelf-0.8.13.tar.gz", "http://www.mr511.de/software/libelf-0.8.13.tar.gz", "standardConfigure")
    fetchBuild(rootDir, "libunwind-1.2", "libunwind-1.2.tar.gz", "http://download.savannah.gnu.org/releases/libunwind/libunwind-1.2.tar.gz", "standardConfigure")
    fetchBuild(rootDir, "zeromq-4.2.1", "zeromq-4.2.1.tar.gz", "https://github.com/zeromq/libzmq/releases/download/v4.2.1/zeromq-4.2.1.tar.gz", "standardConfigure")

	downloadFile("zips/go1.7.5.windows-amd64.zip", "https://storage.googleapis.com/golang/go1.7.5.windows-amd64.zip")
	fmt.Println(figlet("GCC COMPILER"))
	//os.Exit(0)
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		downloadFile("zips/gmp-6.1.2.tar.bz2", "https://gmplib.org/download/gmp/gmp-6.1.2.tar.bz2")
		buildGcc(folderPath)
	} else {
		fmt.Println(figlet("DOWNLOADING"))
		downloadFile("zips/nuwen-14.1.7z", "https://nuwen.net/files/mingw/components-14.1.7z")
		downloadFile("zips/gcc-5.1.0-tdm64-1-core.zip", "https://kent.dl.sourceforge.net/project/tdm-gcc/TDM-GCC%205%20series/5.1.0-tdm64-1/gcc-5.1.0-tdm64-1-core.zip")
		downloadFile("zips/7z1604.exe", "http://www.7-zip.org/a/7z1604.exe")
		doCommand("zips/7z1604.exe", []string{"/S", fmt.Sprintf("/D=%v", SzDir)})
		doCommand("7zip/7z.exe", []string{"x", "zips/nuwen-14.1.7z"})
		os.Chdir("components-14.1")
		files, err := ioutil.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), "7z") {
				unSevenZ(SzPath, file.Name())
			}
		}
		os.Chdir(folderPath)
		os.Setenv("PATH", fmt.Sprintf("%v/components-14.1/bin/;%v", folderPath, os.Getenv("PATH")))
		printEnv()

	}

    os.Exit(0)

	fmt.Println(figlet("GO COMPILER"))
	os.Mkdir(goDir, os.ModeDir|0777)

	os.Setenv("GOPATH", myDir)
	os.Setenv("GOROOT_BOOTSTRAP", runtime.GOROOT())
	printEnv()
	if runtime.GOOS == "darwin" {
		unPackGoMacOSX(folderPath)
	} else if runtime.GOOS == "windows" {
		os.Chdir(goDir)
		unSevenZ(SzPath, "../zips/go1.7.5.windows-amd64.zip")
		os.Chdir(folderPath)
	} else {
		os.Setenv("GOROOT", goDir)
		buildGo()
	}
	printEnv()

	fmt.Println(figlet("LIBRARIES"))
	repos := loadRepos("libs")
	for _, v := range repos {
		installGithub(v)
	}

	fmt.Println(figlet("APPLICATIONS"))
	repos = loadRepos("apps")
	for _, v := range repos {
		installGithub(v)
	}

	fmt.Println(figlet("DO THIS"))
	fmt.Printf("\nNow set your path with one of the following commands\n\n")

	newPath := fmt.Sprintf("%v/usr/local/go/bin/", folderPath)
	fmt.Printf(setCommand(newPath))
	newPath = fmt.Sprintf("%v/bin/", myDir)
	fmt.Printf(setCommand(newPath))
	fmt.Printf("Job's a good'un, boss\n")
}

func setCommand(p string) string {
	return fmt.Sprintf("set -x PATH %v $PATH\nexport PATH=%v/:$PATH\n\n\n", p, p)
}
