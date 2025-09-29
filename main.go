package main

import (
    "fmt"
    "os/exec"

    "github.com/manifoldco/promptui"
)

var translations = map[string]map[string]string{
    "pt": {
        "select_refresh": "Selecione a taxa de atualização",
        "configured":     "Configurado para %s (min: %d, peak: %d)",
        "menu_60":        "60 Hz",
        "menu_90":        "90 Hz",
        "menu_120":       "120 Hz",
	"exit": "Sair",
        "lang_select":    "Selecione o idioma",
        "lang_pt":        "Português",
        "lang_en":        "Inglês",
    },
    "en": {
        "select_refresh": "Select refresh rate",
        "configured":     "Configured to %s (min: %d, peak: %d)",
        "menu_60":        "60 Hz",
        "menu_90":        "90 Hz",
        "menu_120":       "120 Hz",
	"exit": "Exit",
        "lang_select":    "Select language",
        "lang_pt":        "Portuguese",
        "lang_en":        "English",
    },
}

func main() {
    languages := []string{"Português", "English"}
    codeMap := map[string]string{"Português": "pt", "English": "en"}

    langPrompt := promptui.Select{
        Label: translations["en"]["lang_select"],
        Items: languages,
    }

    _, langChosen, err := langPrompt.Run()
    if err != nil {
        fmt.Printf("Prompt failed: %v\n", err)
        return
    }
    lang := codeMap[langChosen]
    t := translations[lang]

    options := []string{t["menu_60"], t["menu_90"], t["menu_120"], t["exit"]}
    hzPrompt := promptui.Select{
        Label: t["select_refresh"],
        Items: options,
    }

    idx, result, err := hzPrompt.Run()
    if err != nil {
        fmt.Printf("Prompt failed: %v\n", err)
        return
    }

    switch idx {
    case 0:
        exec.Command("settings", "put", "system", "min_refresh_rate", "60").Run()
        exec.Command("settings", "put", "system", "peak_refresh_rate", "90").Run()
        fmt.Printf(t["configured"]+"\n", result, 60, 90)
    case 1:
        exec.Command("settings", "put", "system", "min_refresh_rate", "90").Run()
        exec.Command("settings", "put", "system", "peak_refresh_rate", "90").Run()
        fmt.Printf(t["configured"]+"\n", result, 90, 90)
    case 2:
        exec.Command("settings", "put", "system", "min_refresh_rate", "90").Run()
        exec.Command("settings", "put", "system", "peak_refresh_rate", "120").Run()
        fmt.Printf(t["configured"]+"\n", result, 90, 120)
    case 3:
	fmt.Println(result)
    }
}
