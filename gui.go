package main

import (
	"fmt"
	"github.com/pterm/pterm"
)

type ScreenLogType int32

const (
	Debug ScreenLogType = iota
	Info
	Success
	Warning
	Error
)

func ScreenLogo() {
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString("IDS / IPS Update")).Srender()
	pterm.DefaultCenter.Println(s) // Print BigLetters with the default CenterPrinter

	pterm.DefaultCenter.WithCenterEachLineSeparately().Println("Pfsense Update Server!\nCreated by Abolfazl Babamohammadi\n2023")
}

func ScreenDisplayHelp(ServerAddress string) {
	pterm.Info.Println("Go to Pfsense -> Services -> Suricata -> Global Setting\nUpdate input url box from below list\nThen use update button in Update tab for update")
	ETOpenRuleURL := fmt.Sprintf("http://%s/er/v6.0.4/emerging.rules.tar.gz", ServerAddress)
	CommunityRuleURL := fmt.Sprintf("http://%s/cr/v1.0/community-rules.tar.gz", ServerAddress)
	pterm.DefaultTable.WithHasHeader().WithData(pterm.TableData{
		{"Input Name", "URL"},
		{"ETOpen Custom Rule Download URL", ETOpenRuleURL},
		{"Snort GPLv2 Custom Rule Download URL", CommunityRuleURL},
	}).Render()

}

func ScreenNewLog(msg string, t ScreenLogType) {
	switch t {
	case Debug:
		pterm.Debug.Println(msg)
	case Info:
		pterm.Info.Println(msg)
	case Success:
		pterm.Success.Println(msg)
	case Warning:
		pterm.Warning.Println(msg)
	case Error:
		pterm.Error.Println(msg)
	}
}
