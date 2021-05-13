package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
)


var mp4ToTsGroup sync.WaitGroup
var concat sync.WaitGroup

var videoName="绽放吧百合"
var basePath =fmt.Sprintf("E:\\parseIqiyivVdeo\\%s",videoName)


func main() {
	files, _ := ioutil.ReadDir(basePath)
	episodeNum:=len(files)
	for i:=1;i<episodeNum+1;i++{
		dirname:=fmt.Sprintf(`%s\\%s第%d集`,basePath,videoName,i)
		files, err := ioutil.ReadDir(dirname)
		if err != nil {
			panic(err)
		}
		fmt.Printf("将 %s 文件夹下的MP4文件转成TS\n",dirname)
		concatString:="concat:"
		for j:=0;j<len(files);j++{
			mp4ToTsGroup.Add(1)
			go mp4ToTs(dirname, j)
			concatString+=fmt.Sprintf("%s\\%d.ts",dirname,j)
			if j!=len(files)-1{
				concatString+="|"
			}
		}
		mp4ToTsGroup.Wait()

		concat.Add(1)
		go concatVideo(i, concatString,dirname)
	}
	concat.Wait()


}

func concatVideo(episodeNum int, concatString string,dirname string) {
	defer concat.Done()
	concatName := fmt.Sprintf("%s\\%s第%d集.mp4", dirname,videoName,episodeNum)
	fmt.Printf("拼接生成 %s 文件\n",concatName)
	args:=[]string{ "-i", concatString, "-acodec", "copy", "-vcodec", "copy", "-absf", "aac_adtstoasc", concatName}
	cmd := exec.Command("ffmpeg",args...)
	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
	}
	oldpath:=concatName
	newpath:=fmt.Sprintf("%s\\%s第%d集.mp4",basePath,videoName,episodeNum)

	fmt.Printf("将 %s 文件 移动到 %s\n",oldpath,newpath)
	if err:=os.Rename(oldpath,newpath);err!=nil{
		fmt.Println(err.Error())
	}
	fmt.Printf("删除文件夹 %s\n",dirname)
	if err:=os.RemoveAll(dirname);err!=nil{
		fmt.Println(err.Error())
	}


}

func mp4ToTs(dirname string, j int) {
	defer mp4ToTsGroup.Done()
	originName := fmt.Sprintf("%s\\%d.mp4", dirname, j)
	toName := fmt.Sprintf("%s\\%d.ts", dirname, j)
	cmd := exec.Command("ffmpeg", "-i", originName, "-vcodec", "copy", "-acodec", "copy", "-vbsf", "h264_mp4toannexb", toName)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
	if err=os.Remove(originName);err!=nil{
		fmt.Println(err.Error())
	}
}